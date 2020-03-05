package qq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/ankikong/goMusic/tool"

	"github.com/ankikong/goMusic/provider/songbean"
)

type qqVkeyStruct struct {
	Data struct {
		Items []struct {
			Vkey string `json:"vkey"`
		} `json:"items"`
	} `json:"data"`
}

//MusicSearchResult 返回搜索结果
type MusicSearchResult struct {
	Albummid  string `json:"albummid"`
	AlbumName string `json:"albumname"`
	Songmid   string `json:"songmid"`
	Songname  string `json:"songname"`
	Singers   []struct {
		Name string `json:"name"`
	} `json:"singer"`
	Size128  float64 `json:"size128"`
	Size320  float64 `json:"size320"`
	SizeFlac float64 `json:"sizeflac"`
}

//GetURL 默认搜索结果中不会包含URL，所以想要获取链接就一定要调用这个方法
func (q MusicSearchResult) GetURL(_ int) (res songbean.SongInfo) {
	res.SongName = q.Songname
	var br string
	if int(q.Size320) != 0 {
		br = "320"
		res.SongBr = 320
		res.SongSize = int(q.Size320)
	} else if int(q.SizeFlac) != 0 {
		br = "flac"
		res.SongBr = 990
		res.SongSize = int(q.SizeFlac)
	} else if int(q.Size128) != 0 {
		br = "128"
		res.SongBr = 128
		res.SongSize = int(q.Size128)
	} else {
		return
	}
	res.SongURL = GetSongURL(q.Songmid, br).SongURL
	return
}

//GetFileName 生成文件名，可自定义
func (q MusicSearchResult) GetFileName() string {
	return q.GetArtistName() + q.Songname
}

//GetArtistName 生成歌手名字，其实就是把Singer的数据进行合并
func (q MusicSearchResult) GetArtistName() string {
	sb := new(strings.Builder)
	for _, a := range q.Singers {
		sb.WriteString(a.Name)
		sb.WriteString("/")
	}
	res := sb.String()
	if len(res) > 0 {
		res = res[:len(res)-1]
	}
	return res
}

//GetAlbumName 返回AlbumName
func (q MusicSearchResult) GetAlbumName() string {
	return q.AlbumName
}

// GetSource 显示来源
func (q MusicSearchResult) GetSource() string {
	return "QQ"
}

func doGet(url string) []byte {
	rs, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rs.Body.Close()
	buf := new(bytes.Buffer)
	tmpBuf := make([]byte, 4096)
	for {
		len, err := rs.Body.Read(tmpBuf)
		buf.Write(tmpBuf[:len])
		if err != nil {
			break
		}
	}
	return buf.Bytes()
}

func getVkey(mid string) string {
	_api := "http://c.y.qq.com/base/fcgi-bin/fcg_music_express_mobile3.fcg?" +
		"g_tk=1722049047&loginUin=956581739&needNewCode=0&cid=205361747" +
		"&uin=323&songmid=%s&filename=M500%s.mp3&guid=11451466"
	url := fmt.Sprintf(_api, mid, mid)
	res := doGet(url)
	var ans qqVkeyStruct
	json.Unmarshal(res, &ans)
	return ans.Data.Items[0].Vkey
}

type qqMusicIDTmp struct {
	SongName  string `json:"songname"`
	Albumname string `json:"albumname"`
	SongMid   string `json:"songmid"`
	Singer    []struct {
		Name string `json:"name"`
	}
}

type qqMusicURLData struct {
	Code float64 `json:"result"`
	URL  string  `json:"data"`
	Msg  string  `json:"errMsg"`
}

//GetSongURL 根据mid获取歌曲URL
func GetSongURL(songMid, br string) songbean.SongInfo {
	// return fmt.Sprintf("http://aqqmusic.tc.qq.com/amobile.music.tc.qq.com/M500%s.mp3?vkey=%s&guid=3655047200&uin=0&fromtag=8", songMid, vkey)
	var tmpURL string
	if _, err := strconv.ParseInt(songMid, 10, 30); err != nil {
		tmpURL = "https://y.qq.com/n/yqq/song/" + songMid + ".html"
	} else {
		tmpURL = "https://y.qq.com/n/yqq/song/" + songMid + "_num.html?ADTAG=h5_playsong&no_redirect=1"
	}
	rs := doGet(tmpURL)
	reg := regexp.MustCompile(`g_SongData = ([^;]+)`)
	rss := reg.FindAllString(string(rs), -1)[0][13:]
	var ans qqMusicIDTmp
	err := json.Unmarshal([]byte(rss), &ans)
	if err != nil {
		log.Panic(err)
	}
	sb := new(strings.Builder)
	for _, per := range ans.Singer {
		sb.WriteString(per.Name)
		sb.WriteString("/")
	}
	tmpName := sb.String()
	if len(tmpName) > 0 {
		tmpName = tmpName[:len(tmpName)-1]
	}
	ret := new(songbean.SongInfo)
	ret.SongBr = 128
	ret.SongName = tmpName + " - " + ans.SongName
	ret.SongSize = 0
	// vkey := getVkey(ans.SongMid)
	// ret.SongURL = fmt.Sprintf("http://122.226.161.16/amobile.music.tc.qq.com/M500%s.mp3?vkey=%s&guid=11451466&uin=323&fromtag=66", ans.SongMid, vkey)
	// return *ret
	url := fmt.Sprintf("https://api.qq.jsososo.com/song/url?type=%s&id=%s", br, songMid)
	rss, err = tool.DoHTTP("GET", url, "", "", "", "")
	if err != nil {

	} else {
		var data qqMusicURLData
		err = json.Unmarshal([]byte(rss), &data)
		if err != nil || data.Code != 100 {

		} else {
			ret.SongURL = data.URL
		}
	}
	return *ret
}

type qqMusicSearchResults struct {
	Data struct {
		Song struct {
			List []MusicSearchResult `json:"list"`
		} `json:"Song"`
	} `json:"data"`
}

//Search 根据keyword搜索qq上的歌曲
func Search(text string) []MusicSearchResult {
	url := "https://c.y.qq.com/soso/fcgi-bin/client_search_cp?n=5&w=" + text + "&p=1"
	res := doGet(url)
	tmp := string(res)
	tmp = tmp[9 : len(tmp)-1]
	res = []byte(tmp)
	var ans qqMusicSearchResults
	json.Unmarshal(res, &ans)
	return ans.Data.Song.List
}
