package netease

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/ankikong/goMusic/songBean"
)

type NeteaseSong struct {
	Id   int    `json:"id"`
	Url  string `json:"url"`
	Br   int    `json:"br"`
	Size int    `json:"size"`
	// SongType string `json"type"`
}
type NeteaseSongRes struct {
	Data []NeteaseSong `json:"data"`
	Code int           `json:"code"`
}

func doHttp(method, Url, data, encryptoMethod string) string {
	rs := make(map[string]string)
	if encryptoMethod == "web" {
		rs, _ = weapi(data)
	} else if encryptoMethod == "linux" {
		rs = linuxApi(data)
	} else {
		rs = nil
	}
	str := new(strings.Builder)
	if rs != nil {
		for k, v := range rs {
			str.WriteString(k + "=" + v + "&")
		}
		data = str.String()
		data = data[:len(data)-1]
	}
	req, _ := http.NewRequest(method, Url, strings.NewReader(data))
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36")
	req.Header.Add("origin", "https://music.163.com")
	if rs != nil {
		req.Header.Add("content-type", "application/x-www-form-urlencoded")
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	defer res.Body.Close()
	resString := new(strings.Builder)
	tmpBuf := make([]byte, 4096)
	for {
		leng, err := res.Body.Read(tmpBuf)
		resString.Write(tmpBuf[:leng])
		if err != nil {
			break
		}
	}
	return resString.String()
}

// 获取网易的歌曲播放地址
// ids是歌曲的id，br是码率，取990,320,192,128
func GetSongUrl(ids []string, br int) []songBean.SongInfo {
	input := `{"method":"POST","url":"https://music.163.com/api/song/enhance/player/url","params":{"ids":"[` +
		strings.Join(ids, ",") + `]","br":` + fmt.Sprintf("%d", br*1000) + `}}`
	rs := doHttp("POST", "https://music.163.com/api/linux/forward", input, "linux")
	var ans NeteaseSongRes
	json.Unmarshal([]byte(rs), &ans)
	ansRet := make([]songBean.SongInfo, len(ids))
	reg, _ := regexp.Compile(`"title": "[^"]+`)
	index := 0
	for _, id := range ids {
		tmpRs := doHttp("GET", "https://music.163.com/song?id="+id, "", "null")
		songname := reg.FindAllString(tmpRs, 1)[0][10:]
		ansRet[index].SongBr = ans.Data[index].Br
		ansRet[index].SongUrl = ans.Data[index].Url
		ansRet[index].SongName = string(songname)
		ansRet[index].SongSize = ans.Data[index].Size
		index++
	}
	return ansRet
}

type neteaseSearchForm struct {
	S      string `json:"s"`
	Type   int    `json:"type"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Csrf   string `json:"csrf_token"`
}

type neteaseSearchArist struct {
	Name string `json:"name"`
}
type neteaseSearchPerResult struct {
	Album struct {
		AlbumName string `json:"name"`
	} `json:"album"`
	Artists  []neteaseSearchArist `json:"artists"`
	SongName string               `json:"name"`
	SongId   int                  `json:"id"`
}

type neteaseSearch struct {
	Result struct {
		Results []neteaseSearchPerResult `json:"songs"`
	} `json:"result"`
}

type neteaseSearchResult struct {
	FileName   string
	ArtistName string
	AlbumName  string
	SongId     int
}

func (kg neteaseSearchResult) GetFileName() string {
	return kg.FileName
}

func (kg neteaseSearchResult) GetArtistName() string {
	return kg.ArtistName
}

func (kg neteaseSearchResult) GetAlbumName() string {
	return kg.AlbumName
}

func (nt neteaseSearchResult) GetUrl(br int) songBean.SongInfo {
	songId := fmt.Sprint(nt.SongId)
	if br != 990 && br != 320 && br != 192 && br != 128 {
		br = 320
	}
	return GetSongUrl([]string{songId}, br)[0]
}

func Search(text string) []neteaseSearchResult {
	query := new(neteaseSearchForm)
	query.S = text
	query.Type = 1
	query.Csrf = ""
	query.Limit = 10
	query.Offset = 0
	rs, _ := json.Marshal(query)
	ans := doHttp("POST", "https://music.163.com/weapi/search/get", string(rs), "web")
	var tmpAns neteaseSearch
	err := json.Unmarshal([]byte(ans), &tmpAns)
	if err != nil {
		log.Println(err)
	}
	ansRet := make([]neteaseSearchResult, len(tmpAns.Result.Results))
	index := 0
	for _, result := range tmpAns.Result.Results {
		ansRet[index].AlbumName = result.Album.AlbumName
		arName := new(strings.Builder)
		for _, name := range result.Artists {
			arName.WriteString(name.Name)
			arName.WriteString("/")
		}
		tmpName := arName.String()
		if len(tmpName) > 1 {
			tmpName = tmpName[:len(tmpName)-1]
		}
		ansRet[index].ArtistName = tmpName
		ansRet[index].FileName = result.SongName
		ansRet[index].SongId = result.SongId
		index++
	}
	return ansRet
}
