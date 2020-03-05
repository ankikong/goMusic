package kugou

import (
	"encoding/json"

	"github.com/ankikong/goMusic/tool"

	"github.com/ankikong/goMusic/provider/songbean"
)

type kugouSongURL struct {
	SongBr   int      `json:"bitRate"`
	SongSize int      `json:"fileSize"`
	URLs     []string `json:"URL"`
	SongName string   `json:"fileName"`
}

// SearchResult 酷狗搜索的结果
type SearchResult struct {
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

// GetFileName 生成歌曲文件名
func (kg SearchResult) GetFileName() string {
	return kg.FileName
}

// GetArtistName 获取歌手名字
func (kg SearchResult) GetArtistName() string {
	return kg.ArtistName
}

// GetAlbumName 获取专辑名
func (kg SearchResult) GetAlbumName() string {
	return kg.AlbumName
}

// GetSource 获取来源
func (kg SearchResult) GetSource() string {
	return kg.Source
}

// GetURL 搜索结果中是不包含歌曲链接的，所以要获取链接就必须调用此方法
func (kg SearchResult) GetURL(br int) songbean.SongInfo {
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
	List []SearchResult `json:"lists"`
}
type kugouSearchResult struct {
	Data kugouData
}

// GetSongURL 根据id来获取歌曲链接
func GetSongURL(ids []string) []songbean.SongInfo {
	ansRet := make([]songbean.SongInfo, len(ids))
	index := 0
	for _, id := range ids {
		tmpHash := tool.MD5(id + "kgcloudv2")
		api := `http://trackercdn.kugou.com/i/v2/?key=` + tmpHash + `&hash=` + id + `&br=hq&appid=1005&pid=2&cmd=25&behavior=play`
		tmpBuf, err := tool.DoHTTP("GET", api, "", "", "", "")
		if err != nil {
			return nil
		}
		var song kugouSongURL
		err = json.Unmarshal([]byte(tmpBuf), &song)
		if err != nil {
			return nil
		}
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

// Search 根据word搜索
func Search(word string) (ret []SearchResult) {
	ret = *new([]SearchResult)
	URL := `http://songsearch.kugou.com/song_search_v2?pagesize=5&keyword=` + word
	rs, err := tool.DoHTTP("GET", URL, "", "", "", "")
	if err != nil {
		return
	}
	var ans kugouSearchResult
	err = json.Unmarshal([]byte(rs), &ans)
	for i := 0; i < len(ans.Data.List); i++ {
		ans.Data.List[i].Source = "KuGou"
	}
	return ans.Data.List
}
