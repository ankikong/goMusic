package qq

import (
	"fmt"
	"testing"

	"github.com/ankikong/goMusic/provider/songbean"
)

func TestGetVkey(t *testing.T) {
	fmt.Println(GetSongURL("105603683", "320"))
	fmt.Println(GetSongURL("0031Jhwu0ryf6Q", "320"))
}

func TestSearch(t *testing.T) {
	rs := Search("claris")
	fmt.Println(rs)
	var tmp songbean.SongData = rs[0]
	fmt.Println(rs[0].GetURL(0))
	fmt.Println(tmp)
}
