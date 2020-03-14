package tool

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

var (
	currentTimestamp uint32 = 0
)

type flv struct {
	path     string
	Fb       *os.File
	Header   []byte
	Amf2     map[string]interface{}
	restData []byte
}

func byte2int(da []byte) uint32 {
	var offset uint32 = 0
	for _, j := range da {
		offset = (offset << 8) | uint32(j)
	}
	return offset
}

func (f *flv) Init(path string) {
	var err error
	f.Fb, err = os.Open(path)
	if err != nil {
		panic(err)
	}
	f.Header = make([]byte, 13) // header + preTagSize
	leng, err := f.Fb.Read(f.Header)
	if err != nil || leng != 13 {
		panic("read file failed")
	}
	if string(f.Header[0:3]) != "FLV" {
		panic("No a valid flv")
	}
	offset := byte2int(f.Header[5:9])
	//fmt.Println(offset)
	if offset > 9 {
		buff := make([]byte, offset-9)
		_, _ = f.Fb.Read(buff)
	}
}

func (f *flv) GetFirstTag() {
	f.Amf2 = make(map[string]interface{})
	//dis := make([]byte, 4)
	//_, _ = f.Fb.Read(dis)
	tHeader := make([]byte, 11)
	ln, err := f.Fb.Read(tHeader)
	if err != nil || ln != 11 {
		panic("read file tagHeader error")
	}
	if tHeader[0] != 0x12 {
		panic("file format error")
	}
	dataSize := byte2int(tHeader[1:4])
	fTag := make([]byte, dataSize+4) // body + preTagSize
	fLen, err := f.Fb.Read(fTag)
	if uint32(fLen) != dataSize+4 || err != nil {
		panic("error")
	}
	cnt := byte2int(fTag[13:14])
	if cnt != 0x08 {
		panic("contact develop")
	}
	cnt = byte2int(fTag[14:18])

	var pos uint32 = 18
	for i := 0; uint32(i) < cnt; i++ {
		kLen := byte2int(fTag[pos : pos+2])
		pos += 2
		key := string(fTag[pos : pos+kLen])
		pos += kLen
		valType := byte2int(fTag[pos : pos+1])
		pos += 1
		var value interface{}
		if valType == 0 {
			tp := binary.BigEndian.Uint64(fTag[pos : pos+8])
			value = math.Float64frombits(tp)
			pos += 8
		} else if valType == 1 {
			value = bool(fTag[pos] == 0x01)
			pos += 1
		} else if valType == 2 {
			sLen := uint32(binary.BigEndian.Uint16(fTag[pos : pos+2]))
			pos += 2
			value = string(fTag[pos : pos+sLen])
			pos += sLen
		} else {
			//fmt.Println(fTag[fLen-7:fLen])
			f.restData = fTag[pos : fLen-4]
			return
			//panic("unknown type")
		}
		f.Amf2[key] = value
	}
}

func (f *flv) Print() {
	fmt.Printf("header:%v\n", f.Header)
	fmt.Printf("amf2:%v\n", f.Amf2)
}

func (f *flv) Close() {
	f.Fb.Close()
}

func writeHeaderAndTag(fb *os.File, f flv) {
	//fmt.Println(f.Amf2)
	_, err := fb.Write(f.Header) // 正常的话，一个视频的分块，header应该是一致的，随便找一个写进去进去就行了
	if err != nil {
		panic(err)
	}

	// 生成scriptTag 的 data
	buf := bytes.Buffer{}
	tmp := make([]byte, 8)
	// 写入onMetaData
	buf.WriteByte(2)
	binary.BigEndian.PutUint16(tmp[0:2], 10)
	buf.Write(tmp[0:2])
	buf.WriteString("onMetaData")

	// 写入ECMA
	buf.WriteByte(8)
	binary.BigEndian.PutUint32(tmp[0:4], uint32(len(f.Amf2)))
	buf.Write(tmp[0:4])
	for k, v := range f.Amf2 {
		kLen := uint32(len(k))
		binary.BigEndian.PutUint32(tmp[0:4], kLen)
		buf.Write(tmp[0:4])
		buf.WriteString(k)
		switch v.(type) {
		case float64:
			buf.WriteByte(0)
			binary.BigEndian.PutUint64(tmp, math.Float64bits(v.(float64)))
			buf.Write(tmp[:])
		case bool:
			buf.WriteByte(1)
			if v.(bool) {
				buf.WriteByte(1)
			} else {
				buf.WriteByte(0)
			}
			break
		case string:
			buf.WriteByte(2)
			tp := v.(string)
			vLen := uint16(len(tp))
			binary.BigEndian.PutUint16(tmp[0:2], vLen)
			buf.Write(tmp[0:2])
			buf.WriteString(tp)
		}
	}
	// 写入剩余数据
	//buf.Write(f.restData)
	buf.Write([]byte{0, 0, 9})
	data := buf.Bytes()
	buf.Reset()
	// script tag:12
	buf.WriteByte(12)
	// data length
	binary.BigEndian.PutUint32(tmp[0:4], uint32(len(data)))
	buf.Write(tmp[1:4])
	// timestamp
	buf.Write([]byte{0, 0, 0, 0})
	// streamID
	buf.Write([]byte{0, 0, 0})
	_, err = fb.Write(buf.Bytes())
	if err != nil {
		panic(err)
	}
	_, err = fb.Write(data)
	if err != nil {
		panic(err)
	}
	tagLen := len(data) + buf.Len()
	binary.BigEndian.PutUint32(tmp[0:4], uint32(tagLen))
	_, _ = fb.Write(tmp[0:4])
}

func writeBody(output, input *os.File) {
	tHeader := make([]byte, 11)
	tmp := make([]byte, 8)
	var nowTs uint32 = 0
	for true {
		hLen, err := input.Read(tHeader)
		if err != nil || hLen != 11 {
			break
		}
		dataSize := byte2int(tHeader[1:4])
		body := make([]byte, dataSize+4) // body + preTagSize
		bLen, err := input.Read(body)
		if err != nil || uint32(bLen) != dataSize+4 {
			panic("error")
		}

		for i, j := 1, 4; i < 5; i++ {
			tmp[i%4] = tHeader[j]
			j++
		}
		nowTs = binary.BigEndian.Uint32(tmp[0:4])
		newTs := nowTs + currentTimestamp
		binary.BigEndian.PutUint32(tmp[0:4], newTs)
		for i, j := 1, 4; i < 5; i++ {
			tHeader[j] = tmp[i%4]
			j++
		}
		output.Write(tHeader)
		output.Write(body)
	}
	currentTimestamp += nowTs
}

// MergeFLV 只支持合并bilibili的flv，其他flv未知
// 仍有bug: 无法快进
func MergeFLV(output string, input []string) {
	if len(input) < 2 {
		return
	}
	var fls []flv
	for i, j := range input {
		nw := flv{
			path:   j,
			Fb:     nil,
			Header: nil,
			Amf2:   nil,
		}
		nw.Init(j)
		nw.GetFirstTag()
		//nw.Print()
		fls = append(fls, nw)
		if i != 0 {
			edit := []string{"audiosize", "datasize", "filesize", "videosize", "lasttimestamp", "duration"}
			for _, j := range edit {
				add := nw.Amf2[j].(float64)
				ori := fls[0].Amf2[j].(float64)
				ori += add
				fls[0].Amf2[j] = ori

			}
		}
		if i == len(input)-1 {
			fileSize := fls[0].Amf2["filesize"].(float64)
			fileSize -= nw.Amf2["filesize"].(float64)
			fls[0].Amf2["lastkeyframelocation"] = fileSize + nw.Amf2["lastkeyframelocation"].(float64)
			lkfts := fls[0].Amf2["duration"].(float64) - nw.Amf2["duration"].(float64)
			lkfts += nw.Amf2["lastkeyframetimestamp"].(float64)
			fls[0].Amf2["lastkeyframetimestamp"] = lkfts
		}
	}
	fb, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0755)
	//os.Open("E:/tmp/nw.flv")
	if err != nil {
		panic(err)
	}

	writeHeaderAndTag(fb, fls[0])
	for i, _ := range fls {
		//_, err := io.Copy(fb, fls[i].Fb)
		//if err != nil {
		//	panic(err)
		//}
		writeBody(fb, fls[i].Fb)
	}
	for i, _ := range fls {
		fls[i].Fb.Close()
	}
	fb.Close()
}
