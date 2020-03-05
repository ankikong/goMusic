package bilibili

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/ankikong/goMusic/tool"
)

const (
	url       = "https://interface.bilibili.com/v2/playurl?"
	appkey    = "iVGUTjsxvpLeuDCf"
	secretKey = "aHRmhWMLkdeMuILqORnYZocwMBpMEOdt"
	params    = "otype=json&qn=80&quality=80&type=&platform=flash&cid=%s&appkey=%s"
)

// Pages 记录一个av下有分p的信息
type Pages struct {
	CID  int64  `json:"cid"`
	Page int32  `json:"page"`
	Name string `json:"part"`
}
type cidRes struct {
	Code int `json:"code"`
	Data struct {
		Pgs []Pages `json:"pages"`
	} `json:"data"`
}

func getSign(param string) string {
	ps := strings.Split(param, "&")
	sort.Slice(ps, func(i, j int) bool {
		return ps[i] < ps[j]
	})
	param = strings.Join(ps, "&")
	return param + "&sign=" + tool.MD5(param+secretKey)
}

type durl struct {
	URL  string  `json:"url"`
	Size float64 `json:"size"`
}
type bilibiliRet struct {
	Format string `json:"format"`
	URLs   []durl `json:"durl"`
}

// GetVideoURL 根据所给CID获取视频链接
func GetVideoURL(CID string) string {
	ps := fmt.Sprintf(params, CID, appkey)
	ps = getSign(ps)
	rs, _ := tool.DoHTTP("GET", url+ps, "", "", "", "")
	var data bilibiliRet
	json.Unmarshal([]byte(rs), &data)

	return data.URLs[0].URL
}

// GetCID 根据aid获取视频所有分P的cid
func GetCID(aid string) []Pages {
	url := "https://api.bilibili.com/x/web-interface/view?aid=" + aid
	rs, err := tool.DoHTTP("GET", url, "", "", "", "")
	if err != nil {
		panic("contact develop to fix(or retry):" + err.Error())
	}
	var data cidRes
	err = json.Unmarshal([]byte(rs), &data)
	if err != nil {
		panic("contact develop to fix(or retry):" + err.Error())
	}
	if data.Code != 0 {
		panic("contact develop to fix(or retry):" + err.Error())
	}
	return data.Data.Pgs
}

// Deal 处理链接并询问下载
func Deal(url string) {
	reg := regexp.MustCompile(`av(\d+)`)
	aid := string(reg.Find([]byte(url)))
	cid := GetCID(aid[2:])
	fmt.Printf("Get %d P video:\n", len(cid))
	for j, i := range cid {
		fmt.Printf("\t%3d: %s\n", j, i.Name)
	}
	fmt.Println("which page to download?(-1 for all, use space to split cids, none to quit):")
	reader := bufio.NewReader(os.Stdin)
	inbuf, _ := reader.ReadBytes([]byte("\n")[0])
	in := string(inbuf)
	in = strings.TrimSpace(in)
	cids := strings.Split(in, " ")
	if len(cids) == 0 {
		return
	}
	fmt.Println(cids)
	for _, i := range cids {
		if ind, err := strconv.Atoi(i); err == nil && ind < len(cid) {
			url := GetVideoURL(fmt.Sprint(cid[ind].CID))
			fmt.Println(url)
			tool.Download(url, fmt.Sprintf("%s-%d.%s", aid, ind, "flv"), "")
		} else {
			fmt.Printf("error input, no a digit: \"%s\", skip", i)
		}
	}
}
