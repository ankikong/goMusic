package tool

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

// Download 下载url的文件到路径path
func Download(url, name, path string, inferExt bool) error {
	path = strings.TrimSpace(path)
	if len(path) == 0 {
		path = "./"
	} else if path[len(path)-1] != '/' {
		path = path + "/"
	}
	if inferExt {
		exts := []string{"mp3", "flac", "acc", "flv", "mp4"}
		var ext string
		for _, j := range exts {
			if strings.Contains(url, j) {
				ext = j
				break
			}
		}
		name = name + "." + ext
	}
	// reg, _ := regexp.Compile(`(/\|<>:*?")`)
	for _, i := range `/\|<>:*?"` {
		name = strings.ReplaceAll(name, string(i), " ")
	}

	rs, err := http.Get(url)

	if err != nil {
		log.Println("download fail:", err.Error())
		return err
	}

	file, err := os.Create(path + name)
	if err != nil {
		log.Println("create file fail:", err.Error())
		return err
	}

	lng, err := strconv.ParseInt(rs.Header.Get("content-length"), 0, 60)
	if err != nil {
		fmt.Println("start download(unknown size):" + name)
		buf := make([]byte, 262144)
		for {
			len, err := rs.Body.Read(buf)
			file.Write(buf[:len])
			if err != nil {
				break
			}
			fmt.Print(".")
		}
	} else {
		bar := pb.Full.Start64(lng)
		barReader := bar.NewProxyReader(rs.Body)
		io.Copy(file, barReader)
		bar.Finish()
	}
	rs.Body.Close()
	file.Close()
	return nil
}

// DoHTTP 以method发起对URL的请求
func DoHTTP(method, URL, data, dataFormat, origin, referer string) (ret string, err error) {

	req, _ := http.NewRequest(method, URL, strings.NewReader(data))
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36")
	if len(origin) > 0 {
		req.Header.Add("origin", origin)
	}
	if len(referer) > 0 {
		req.Header.Add("referer", referer)
	}
	if len(data) > 0 && len(dataFormat) > 0 {
		req.Header.Add("content-type", dataFormat)
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer res.Body.Close()
	resString := new(strings.Builder)
	tmpBuf := make([]byte, 4096)
	for {
		leng, ierr := res.Body.Read(tmpBuf)
		resString.Write(tmpBuf[:leng])
		if ierr == io.EOF {
			break
		} else if ierr != nil {
			err = ierr
			return
		}
	}
	return resString.String(), nil

}

// MapToURLParams 将map转换为url参数格式
func MapToURLParams(in map[string]string) (ret string) {
	param := url.Values{}
	for k, v := range in {
		param.Add(k, v)
	}
	return param.Encode()
}
