package songBean

// 这个是用来放回歌曲文件的解析度，文件名，
// 文件大小，和文件的链接
type SongInfo struct {
	SongName string
	SongUrl  string
	SongBr   int
	SongSize int
}

// 这个是统一的搜索返回结果都要实现的接口，方便统一管理
type SongData interface {
	// 获取歌曲链接
	GetUrl(br int) SongInfo
	// 获取歌曲名称
	GetFileName() string
	// 获取歌手名字
	GetArtistName() string
	// 获取歌单名字
	GetAlbumName() string
	// 获取歌曲来源
	GetSource() string
}
