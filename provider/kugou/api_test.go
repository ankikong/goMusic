package kugou

import (
	"fmt"
	"testing"

	"github.com/ankikong/goMusic/provider/songbean"
)

func TestGetSongURL(t *testing.T) {
	rs := GetSongURL([]string{"970570A464E93978D3CDB13B4CB5EA2A", "40AD169093CDE5523A13DA8E7A09066B"})
	fmt.Println(rs)
}

func TestSearch(t *testing.T) {
	rs := Search("claris")
	fmt.Println(rs)
	rs[0].GetURL(320)
	var val songbean.SongData = rs[0]
	fmt.Println(val.GetURL(320))
}
