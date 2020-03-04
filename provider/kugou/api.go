package kugou

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ankikong/goMusic/songbean"
)

type kugouSongURL struct {
	SongBr   int      `json:"bitRate"`
	SongSize int      `json:"fileSize"`
	URLs     []string `json:"URL"`
	SongName string   `json:"fileName"`
}

type kugouSearchPerResult struct {
	LQHash     string `json:"FileHash"`
	LQSize     int    `json:"FileSize"`
	HQHash     string `json:"HQFileHash"`
	HQSize     int    `json:"HQFileSize"`
	SQHash     string `json:"SQFileHash"`
	SQSize     int    `json:"SQFileSize"`
	FileName   string `json:"FileName"`
	ArtistName string `json:"SingerName"`
	AlbumName  string `json:"AlbumName"`
	Source     string
}

func (kg kugouSearchPerResult) GetFileName() string {
	return kg.FileName
}

func (kg kugouSearchPerResult) GetArtistName() string {
	return kg.ArtistName
}

func (kg kugouSearchPerResult) GetAlbumName() string {
	return kg.AlbumName
}

func (kg kugouSearchPerResult) GetSource() string {
	return kg.Source
}

func (kg kugouSearchPerResult) GetURL(br int) songbean.SongInfo {
	var rs songbean.SongInfo
	if br == 990 && kg.SQHash != "" {
		rs = GetSongURL([]string{kg.SQHash})[0]
	}

	if br == 320 || (br > 320 && len(rs.SongURL) < 5) {
		rs = GetSongURL([]string{kg.HQHash})[0]
	}

	if br == 192 || len(rs.SongURL) < 5 {
		rs = GetSongURL([]string{kg.LQHash})[0]
	}
	return rs
}

type kugouData struct {
	List []kugouSearchPerResult `json:"lists"`
}
type kugouSearchResult struct {
	Data kugouData
}

func MD5(text string) string {
	data := []byte(text)
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

func doGet(URL string) []byte {
	rs, err := http.Get(URL)
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

func GetSongURL(ids []string) []songbean.SongInfo {
	ansRet := make([]songbean.SongInfo, len(ids))
	index := 0
	for _, id := range ids {
		tmpHash := MD5(id + "kgcloudv2")
		api := `http://trackercdn.kugou.com/i/v2/?key=` + tmpHash + `&hash=` + id + `&br=hq&appid=1005&pid=2&cmd=25&behavior=play`
		tmpBuf := doGet(api)
		var song kugouSongURL
		json.Unmarshal(tmpBuf, &song)
		ansRet[index].SongBr = song.SongBr
		ansRet[index].SongName = song.SongName
		ansRet[index].SongSize = song.SongSize

		if len(song.URLs) > 0 {
			ansRet[index].SongURL = song.URLs[0]
		}
		index++
	}
	return ansRet
}

func Search(word string) []kugouSearchPerResult {
	URL := `http://songsearch.kugou.com/song_search_v2?pagesize=5&keyword=` + word
	rs := doGet(URL)
	// fmt.Println(string(rs))
	var ans kugouSearchResult
	json.Unmarshal(rs, &ans)
	for i := 0; i < len(ans.Data.List); i++ {
		ans.Data.List[i].Source = "KuGou"
	}
	return ans.Data.List
}
