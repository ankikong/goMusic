package zzzfun

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/ankikong/goMusic/tool"
)

const (
	epsURL    = "http://service-agbhuggw-1259251677.gz.apigw.tencentcs.com/android/video/list_ios?videoId="
	playURL   = "http://service-agbhuggw-1259251677.gz.apigw.tencentcs.com/android/video/newplay"
	secretKey = "534697"
)

func getSign(playID string) string {
	return tool.MD5(playID + secretKey)
}

// Episode 分集信息
type Episode struct {
	PlayID string `json:"playid"`
	Index  string `json:"ji"`
	Title  string
}

type epsLocalStruct struct {
	Eps []Episode `json:"list"`
}

type epsStruct struct {
	Code int    `json:"errorCode"`
	Msg  string `json:"errorMsg"`
	Data struct {
		VideoName string           `json:"videoName"`
		Videos    []epsLocalStruct `json:"videoSets"`
	} `json:"data"`
}

// GetEps 获取所有分集
func GetEps(videoID string) (ret []Episode) {
	url := epsURL + videoID
	raw, err := tool.DoHTTP("GET", url, "", "", "", "")
	if err != nil {
		panic(err)
	}
	var data epsStruct
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		panic(err)
	}
	if data.Code != 0 {
		panic(data.Msg)
	}
	videoName := data.Data.VideoName
	ret = data.Data.Videos[0].Eps
	for i := range ret {
		ret[i].Title = fmt.Sprintf("[%s][%d]", videoName, i)
	}
	return
}

type urlStruct struct {
	Code int    `json:"errorCode"`
	Msg  string `json:"errorMsg"`
	Data struct {
		URL string `json:"videoplayurl"`
	}
}

// GetVideoURL 通过playID获取视频链接
func GetVideoURL(playID string) string {
	// 新的sing算法在 org.daimhim.zzzfun.data.remote.HttpRequestManager的invokeSuspend方法里
	// map=int(time.time()*1000)
	// sing=md5(map的值+zandroidzz)
	// 例如：map=1486876988464，sing=md5("zandroidzz1486876988464")
	// 但是他们没有验证时间戳的时间和playid，所以直接定死就完事了
	data := "playid=" + playID + "&sing=a47bdac30dd237e18f187cee332b3d2a&map=1486876988464"
	// fmt.Println(data)
	//data = data + "&" + getSign(data)
	raw, err := tool.DoHTTP("POST", playURL, data, "application/x-www-form-urlencoded", "", "")
	if err != nil {
		panic(err)
	}
	var rs urlStruct
	if err := json.Unmarshal([]byte(raw), &rs); err != nil {
		panic(err)
	}
	if rs.Code != 0 {
		panic(rs.Msg)
	}
	return rs.Data.URL
}

// GetByURL 通过url分析数据
func GetByURL(url string) {
	reg := regexp.MustCompile(`\d+`)
	videoID := string(reg.Find([]byte(url)))
	eps := GetEps(videoID)
	for i, j := range eps {
		fmt.Printf("%3d-%s\n", i, j.Title)
	}
	fmt.Print("input index to download(-1 for all):")
	var ind int
	_, err := fmt.Scan(&ind)
	if err != nil || ind < -1 || ind >= len(eps) {
		fmt.Println("error input")
		return
	}
	if ind == -1 {
		for _, j := range eps {
			url := GetVideoURL(j.PlayID)
			tool.Download(url, j.Title, "", true, true)
		}
	} else {
		url := GetVideoURL(eps[ind].PlayID)
		tool.Download(url, eps[ind].Title, "", true, true)
	}

}
