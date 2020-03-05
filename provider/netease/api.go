package netease

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/ankikong/goMusic/provider/songbean"
	"github.com/ankikong/goMusic/tool"
)

const (
	iv          = "0102030405060708"
	presetKey   = "0CoJUm6Qyw8W8jud"
	linuxapiKey = "rFgB&h#%2?^eDg:Q"
	base62      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	publicKey   = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDgtQn2JZ34ZC28NWYpAUd98iZ37BUrX/aKzmFbt7clFSs6sXqHauqKWqdtLkF2KexO40H1YTX8z2lSgBBOAxLsvaklV8k4cBFK9snQXE9/DDaFt6Rr7iVZMldczhC0JNgTz+SHXT6CBHuX3e9SdB1Ua44oncaTWz7OBGLbCiK45wIDAQAB\n-----END PUBLIC KEY-----"
	eapiKey     = "e82ckenh8dichen8"
	modulus     = "00e0b509f6259df8642dbc35662901477df22677ec152b5ff68ace615bb7" +
		"b725152b3ab17a876aea8a5aa76d2e417629ec4ee341f56135fccf695280" +
		"104e0312ecbda92557c93870114af6c9d05c4f7f0c3685b7a46bee255932" +
		"575cce10b424d813cfe4875d3e82047b97ddef52741d546b8e289dc6935b" +
		"3ece0462db0a22b8e7"
)

// module

func weapi(text string) (rs map[string]string) {
	secretKey := make([]byte, 16)
	for i := 0; i < 16; i++ {
		secretKey[i] = byte(base62[rand.Int31n(62)])
	}
	param := base64.StdEncoding.EncodeToString(tool.AesEncryptCBC([]byte(text), []byte(iv), []byte(presetKey)))
	param = base64.StdEncoding.EncodeToString(tool.AesEncryptCBC([]byte(param), []byte(iv), secretKey))
	for i, j := 0, 15; i < j; i++ {
		secretKey[i], secretKey[j] = secretKey[j], secretKey[i]
		j--
	}
	data := tool.RsaEncrypt(secretKey, modulus)
	rs = make(map[string]string)
	rs["params"], rs["encSecKey"] = param, data
	return rs
}

func linuxAPI(text string) map[string]string {
	rs := tool.AesEncryptECB([]byte(text), []byte(linuxapiKey))
	ret := make(map[string]string)
	ret["eparams"] = strings.ToUpper(hex.EncodeToString(rs))
	return ret
}

type neteaseSong struct {
	ID   int    `json:"id"`
	URL  string `json:"url"`
	Br   int    `json:"br"`
	Size int    `json:"size"`
	// SongType string `json"type"`
}
type neteaseSongRes struct {
	Data []neteaseSong `json:"data"`
	Code int           `json:"code"`
}

func doHTTP(method, URL, data, encryptoMethod string) string {
	rs := make(map[string]string)
	if encryptoMethod == "web" {
		rs = weapi(data)
	} else if encryptoMethod == "linux" {
		rs = linuxAPI(data)
	} else {
		rs = nil
	}
	if rs != nil {
		param := url.Values{}
		for k, v := range rs {
			param.Add(k, v)
		}
		data = param.Encode()
	}
	req, _ := http.NewRequest(method, URL, strings.NewReader(data))
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

// GetSongURL 的歌曲播放地址
// ids是歌曲的id，br是码率，取990,320,192,128
func GetSongURL(ids []string, br int) []songbean.SongInfo {
	input := `{"method":"POST","url":"https://music.163.com/api/song/enhance/player/url","params":{"ids":"[` +
		strings.Join(ids, ",") + `]","br":` + fmt.Sprintf("%d", br*1000) + `}}`
	// rs := doHTTP("POST", "https://music.163.com/api/linux/forward", input, "linux")

	postData := tool.MapToURLParams(linuxAPI(input))

	rs, err := tool.DoHTTP("POST", "https://music.163.com/api/linux/forward", postData, "application/x-www-form-urlencoded", "https://music.163.com", "")
	if err != nil {
		return nil
	}
	var ans neteaseSongRes
	json.Unmarshal([]byte(rs), &ans)
	ansRet := make([]songbean.SongInfo, len(ids))
	reg, _ := regexp.Compile(`"title": "[^"]+`)
	index := 0
	for _, id := range ids {
		// tmpRs := doHTTP("GET", "https://music.163.com/song?id="+id, "", "null")
		tmpRs, err := tool.DoHTTP("GET", "https://music.163.com/song?id="+id, "", "", "", "https://music.163.com")
		if err != nil {
			return nil
		}
		songname := reg.FindAllString(tmpRs, 1)[0][10:]
		ansRet[index].SongBr = ans.Data[index].Br
		ansRet[index].SongURL = ans.Data[index].URL
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
	SongID   int                  `json:"id"`
}

type neteaseSearch struct {
	Result struct {
		Results []neteaseSearchPerResult `json:"songs"`
	} `json:"result"`
}

// SearchResult 网易云歌曲详情
type SearchResult struct {
	FileName   string
	ArtistName string
	AlbumName  string
	Source     string
	SongID     int
}

// GetFileName 返回文件名
func (nt SearchResult) GetFileName() string {
	return nt.FileName
}

// GetArtistName 放回歌手名
func (nt SearchResult) GetArtistName() string {
	return nt.ArtistName
}

// GetAlbumName 返回专辑名
func (nt SearchResult) GetAlbumName() string {
	return nt.AlbumName
}

// GetSource 返回歌曲来源
func (nt SearchResult) GetSource() string {
	return nt.Source
}

// GetURL 搜索结果中不含歌曲名称，所以要根据ID获取歌曲实际地址
func (nt SearchResult) GetURL(br int) songbean.SongInfo {
	SongID := fmt.Sprint(nt.SongID)
	if br != 990 && br != 320 && br != 192 && br != 128 {
		br = 320
	}
	var rs songbean.SongInfo
	rs = GetSongURL([]string{SongID}, br)[0]
	if len(rs.SongURL) < 5 {
		rs = GetSongURL([]string{SongID}, 128)[0]
	}
	rs.SongName = nt.FileName
	return rs
}

// Search 网易云搜索
func Search(text string) []SearchResult {
	query := new(neteaseSearchForm)
	query.S = text
	query.Type = 1
	query.Csrf = ""
	query.Limit = 10
	query.Offset = 0
	rs, _ := json.Marshal(query)
	// ans := doHTTP("POST", "https://music.163.com/weapi/search/get", string(rs), "web")
	postData := tool.MapToURLParams(weapi(string(rs)))
	ans, err := tool.DoHTTP("POST", "https://music.163.com/weapi/search/get", postData, "application/x-www-form-urlencoded", "https://music.163.com", "")
	var tmpAns neteaseSearch
	err = json.Unmarshal([]byte(ans), &tmpAns)
	if err != nil {
		log.Println(ans)
		log.Println(err)
	}
	ansRet := make([]SearchResult, len(tmpAns.Result.Results))
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
		ansRet[index].Source = "netease"
		ansRet[index].ArtistName = tmpName
		ansRet[index].FileName = tmpName + " - " + result.SongName
		ansRet[index].SongID = result.SongID
		index++
	}
	return ansRet
}
