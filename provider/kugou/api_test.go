package kugou

import (
	"fmt"
	"testing"

	"github.com/ankikong/goMusic/songBean"
)

func TestGetSongUrl(t *testing.T) {
	rs := GetSongUrl([]string{"4D870D0DEB13AA6700BEECA513C6B03C", "40AD169093CDE5523A13DA8E7A09066B"})
	fmt.Println(rs)
}

func TestSearch(t *testing.T) {
	rs := Search("claris")
	fmt.Println(rs)
	rs[0].GetUrl(320)
	var val songBean.SongData = rs[0]
	fmt.Println(val.GetUrl(320))
}
