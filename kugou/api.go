package kugou

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ankikong/goMusic/songBean"
)

type kugouSongUrl struct {
	SongBr   int      `json:"bitRate"`
	SongSize int      `json:"fileSize"`
	Urls     []string `json:"url"`
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

func (kg kugouSearchPerResult) GetUrl(br int) songBean.SongInfo {
	var rs songBean.SongInfo
	if br == 990 && kg.SQHash != "" {
		rs = GetSongUrl([]string{kg.SQHash})[0]
	}

	if br == 320 || (br > 320 && len(rs.SongUrl) < 5) {
		rs = GetSongUrl([]string{kg.HQHash})[0]
	}

	if br == 192 || len(rs.SongUrl) < 5 {
		rs = GetSongUrl([]string{kg.LQHash})[0]
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

func doGet(url string) []byte {
	rs, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rs.Body.Close()
	tmpBuf := make([]byte, 65536)
	len, _ := rs.Body.Read(tmpBuf)
	return tmpBuf[:len]
}

func GetSongUrl(ids []string) []songBean.SongInfo {
	ansRet := make([]songBean.SongInfo, len(ids))
	index := 0
	for _, id := range ids {
		tmpHash := MD5(id + "kgcloudv2")
		api := `http://trackercdn.kugou.com/i/v2/?key=` + tmpHash + `&hash=` + id + `&br=hq&appid=1005&pid=2&cmd=25&behavior=play`
		tmpBuf := doGet(api)
		var song kugouSongUrl
		json.Unmarshal(tmpBuf, &song)
		ansRet[index].SongBr = song.SongBr
		ansRet[index].SongName = song.SongName
		ansRet[index].SongSize = song.SongSize

		if len(song.Urls) > 0 {
			ansRet[index].SongUrl = song.Urls[0]
		}
		index++
	}
	return ansRet
}

func Search(word string) []kugouSearchPerResult {
	Url := `http://songsearch.kugou.com/song_search_v2?pagesize=5&keyword=` + word
	rs := doGet(Url)
	// fmt.Println(string(rs))
	var ans kugouSearchResult
	json.Unmarshal(rs, &ans)
	for i := 0; i < len(ans.Data.List); i++ {
		ans.Data.List[i].Source = "KuGou"
	}
	return ans.Data.List
}
