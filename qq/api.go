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

	"github.com/ankikong/goMusic/songBean"
)

type qqVkeyStruct struct {
	Data struct {
		Items []struct {
			Vkey string `json:"vkey"`
		} `json:"items"`
	} `json:"data"`
}

type QQMusicSearchResult struct {
	Albummid  string `json:"albummid"`
	AlbumName string `json:"albumname"`
	Songmid   string `json:"songmid"`
	Songname  string `json:"songname"`
	Singers   []struct {
		Name string `json:"name"`
	} `json:"singer"`
}

func (q QQMusicSearchResult) GetUrl(_ int) songBean.SongInfo {
	// res := new(songBean.SongInfo)
	// res.SongBr = 128
	// res.SongName = q.Songname
	// res.SongSize = 0
	return GetSongUrl(q.Songmid)
	// return *res
}

func (q QQMusicSearchResult) GetFileName() string {
	return q.GetArtistName() + q.Songname
}

func (q QQMusicSearchResult) GetArtistName() string {
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

func (q QQMusicSearchResult) GetAlbumName() string {
	return q.AlbumName
}

func (q QQMusicSearchResult) GetSource() string {
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
	_api := "https://c.y.qq.com/base/fcgi-bin/fcg_music_express_mobile3.fcg?format=json" +
		"&cid=205361747&uin=0&songmid=%s&filename=M800%s.mp4&guid=3655047200&platform=yqq"
	url := fmt.Sprintf(_api, mid, mid)
	res := doGet(url)
	var ans qqVkeyStruct
	json.Unmarshal(res, &ans)
	return ans.Data.Items[0].Vkey
}

type qqMusicIdTmp struct {
	SongName  string `json:"songname"`
	Albumname string `json:"albumname"`
	SongMid   string `json:"songmid"`
	Singer    []struct {
		Name string `json:"name"`
	}
}

func GetSongUrl(songMid string) songBean.SongInfo {
	// return fmt.Sprintf("http://aqqmusic.tc.qq.com/amobile.music.tc.qq.com/M500%s.mp3?vkey=%s&guid=3655047200&uin=0&fromtag=8", songMid, vkey)
	var tmpUrl string
	if _, err := strconv.ParseInt(songMid, 10, 30); err != nil {
		tmpUrl = "https://y.qq.com/n/yqq/song/" + songMid + ".html"
	} else {
		tmpUrl = "https://y.qq.com/n/yqq/song/" + songMid + "_num.html?ADTAG=h5_playsong&no_redirect=1"
	}
	rs := doGet(tmpUrl)
	reg := regexp.MustCompile(`g_SongData = ([^;]+)`)
	rss := reg.FindAllString(string(rs), -1)[0][13:]
	var ans qqMusicIdTmp
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
	ret := new(songBean.SongInfo)
	ret.SongBr = 128
	ret.SongName = tmpName + " - " + ans.SongName
	ret.SongSize = 0
	vkey := getVkey(ans.SongMid)
	ret.SongUrl = fmt.Sprintf("http://aqqmusic.tc.qq.com/amobile.music.tc.qq.com/M500%s.mp3?vkey=%s&guid=3655047200&uin=0&fromtag=8", ans.SongMid, vkey)
	return *ret
}

type qqMusicSearchResults struct {
	Data struct {
		Song struct {
			List []QQMusicSearchResult `json:"list"`
		} `json:"Song"`
	} `json:"data"`
}

func Search(text string) []QQMusicSearchResult {
	url := "https://c.y.qq.com/soso/fcgi-bin/client_search_cp?n=5&w=" + text + "&p=1"
	res := doGet(url)
	tmp := string(res)
	tmp = tmp[9 : len(tmp)-1]
	res = []byte(tmp)
	var ans qqMusicSearchResults
	json.Unmarshal(res, &ans)
	return ans.Data.Song.List
}
