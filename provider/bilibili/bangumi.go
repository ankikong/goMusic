package bilibili

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ankikong/goMusic/tool"
)

const (
	epsURL        = "https://api.bilibili.com/pgc/web/season/section?season_id="
	originURL     = "https://bangumi.bilibili.com/player/web_api/playurl/?"
	bangumiParams = "appkey=%s&cid=%d&module=bangumi&otype=json&platform=flash&player=1&qn=%s&type="
)

// Episode 番剧每一集的信息
type Episode struct {
	Aid   int    `json:"aid"`
	Cid   int    `json:"cid"`
	Title string `json:"long_title"`
	Index string `json:"title"`
}

// URL 保存某个番剧的一集的某一个块
// 因为发起的是flash请求，所以应该只有一个完整的块
// Length 视频的时间长度
// Order 第几个分块
// Size 此分块的物理大小
// URL 视频的链接
type URL struct {
	Length int    `json:"length"`
	Order  int    `json:"order"`
	Size   int    `json:"size"`
	URL    string `json:"url"`
	Format string
}

type rawURL struct {
	Code       int    `json:"code"`
	Format     string `json:"format"`
	TimeLength int    `json:"timelength"`
	DURL       []URL  `json:"durl"`
	Message    string `json:"message"`
}
type epsStruct struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  struct {
		Data struct {
			Episodes []Episode `json:"episodes"`
		} `json:"main_section"`
	} `json:"result"`
}

// SeasonInfo 一季番剧的信息
type SeasonInfo struct {
	SeasonID    int    `json:"season_id"`
	SeasonTitle string `json:"season_title"`
	Title       string `json:"title"`
}

type midStrut struct {
	MediaInfo struct {
		Seasons []SeasonInfo `json:"seasons"`
	} `json:"mediaInfo"`
}

// GetEpisodesBySeasonID 通过season ID获取该番剧的所有分集
func GetEpisodesBySeasonID(seasonID string) []Episode {
	rs, err := tool.DoHTTP("GET", epsURL+seasonID, "", "", "", "")
	if err != nil {
		panic(err)
	}
	var data epsStruct
	if err := json.Unmarshal([]byte(rs), &data); err != nil {
		panic(err)
	}
	if data.Code != 0.0 {
		panic(data.Message)
	}
	return data.Result.Data.Episodes
}

// GetBangumiURL 通过cid获取番剧的实际地址
func GetBangumiURL(cid int) []URL {
	param := fmt.Sprintf(bangumiParams, appkey, cid, "112")
	// param = getSign(param)
	// fmt.Println(originURL + param)
	rs, err := tool.DoHTTP("GET", originURL+param, "", "", "", "")
	if err != nil {
		panic(err)
	}
	var data rawURL
	if err := json.Unmarshal([]byte(rs), &data); err != nil {
		panic(err)
	}
	if data.Code != 0 {
		panic(data.Message)
	}
	format := data.Format
	for i := range data.DURL {
		data.DURL[i].Format = format
	}
	return data.DURL
}

// GetAllSeasonByMIDURL 通过总集页获取番剧所有季的信息
func GetAllSeasonByMIDURL(url string) []SeasonInfo {
	rs, err := tool.DoHTTP("GET", url, "", "", "", "")
	if err != nil {
		panic(err)
	}
	reg := regexp.MustCompile(`__INITIAL_STATE__=(.*?);\(`)
	data := reg.Find([]byte(rs))
	data = data[18 : len(data)-2]
	var rawData midStrut
	if err := json.Unmarshal(data, &rawData); err != nil {
		panic(err)
	}
	return rawData.MediaInfo.Seasons
}

// GetSeasonIDByPlayPage 通过播放页获取此番剧的season ID
func GetSeasonIDByPlayPage(url string) string {
	rs, err := tool.DoHTTP("GET", url, "", "", "", "")
	if err != nil {
		panic(err)
	}
	reg := regexp.MustCompile(`"ssId":\d+`)
	seasonID := string(reg.Find([]byte(rs)))[7:]
	return seasonID
}

// BangumiDeal 处理番剧链接
func BangumiDeal(url string) {
	var seasonID string
	if strings.Contains(url, "media") {
		data := GetAllSeasonByMIDURL(url)
		for j, i := range data {
			fmt.Printf("%d:%s-%s\n", j, i.SeasonTitle, i.Title)
		}
		var ind int
		fmt.Printf("select which season:")
		fmt.Scan(&ind)
		if ind > 0 && ind < len(data) {
			seasonID = fmt.Sprint(data[ind].SeasonID)
		} else {
			panic("error input")
		}
	} else if strings.Contains(url, "play") {
		seasonID = GetSeasonIDByPlayPage(url)
	} else {
		panic("unsupport url")
	}
	rs := GetEpisodesBySeasonID(seasonID)
	for i, j := range rs {
		fmt.Printf("%d:%s\n", i, j.Title)
	}
	fmt.Printf("select which episode to download:")
	var ind int
	fmt.Scan(&ind)
	// ind := 0
	if ind >= 0 && ind < len(rs) {
		tmp := GetBangumiURL(rs[ind].Cid)
		if len(tmp) == 0 {
			fmt.Println("get url error with length 0")
			return
		}
		fmt.Println("has", len(tmp), "part(s)")
		for _, j := range tmp {
			tool.Download(j.URL, fmt.Sprintf("%s-%d", seasonID, j.Order), "", false, false)
		}
		if len(tmp) == 1 {
			os.Rename(fmt.Sprintf("%s-%d", seasonID, 1),
				fmt.Sprintf("%s-%d.flv", seasonID, ind))
		} else {
			var input []string
			for i := range tmp {
				input = append(input, fmt.Sprintf("%s-%d", seasonID, i+1))
			}
			tool.MergeFLV(fmt.Sprintf("%s-%d.flv", seasonID, ind), input)
		}
	}
}
