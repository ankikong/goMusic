package qq

import (
	"fmt"
	"testing"

	"github.com/ankikong/goMusic/songBean"
)

func TestGetVkey(t *testing.T) {
	fmt.Println(GetSongUrl("105603683"))
	fmt.Println(GetSongUrl("0031Jhwu0ryf6Q"))
}

func TestSearch(t *testing.T) {
	rs := Search("claris")
	fmt.Println(rs)
	var tmp songBean.SongData = rs[0]
	fmt.Println(rs[0].GetUrl(0))
	fmt.Println(tmp)
}
