package bilibili

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ankikong/goMusic/tool"
)

const (
	url       = "https://interface.bilibili.com/v2/playurl?"
	appkey    = "iVGUTjsxvpLeuDCf"
	secretKey = "aHRmhWMLkdeMuILqORnYZocwMBpMEOdt"
	params    = "otype=json&qn=80&quality=80&type=&platform=flash&cid=%s&appkey=%s"
)

type pages struct {
}

func getSign(param string) string {
	ps := strings.Split(param, "&")
	sort.Slice(ps, func(i, j int) bool {
		return ps[i] < ps[j]
	})
	param = strings.Join(ps, "&")
	return param + "&sign=" + tool.MD5(param+secretKey)
}

// GetVideoURL 根据所给CID获取视频链接
func GetVideoURL(CID string) string {
	ps := fmt.Sprintf(params, CID, appkey)
	ps = getSign(ps)
	return url + ps
}

// GetCID 根据aid获取视频所有分P的cid
func GetCID(aid string) string {
	url := "https://api.bilibili.com/x/web-interface/view?aid=" + aid
	rs := tool.DoHTTP("GET", )
}
