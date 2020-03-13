package bilibili

import (
	"fmt"
	"testing"
)

func TestGetEpisodesBySeasonID(t *testing.T) {
	rs := GetEpisodesBySeasonID("28937")
	fmt.Println(rs)
}

func TestGetBangumiURL(t *testing.T) {
	rs := GetBangumiURL(139676062)
	fmt.Println(rs)
}

func TestGetAllByPlayPage(t *testing.T) {
	GetSeasonIDByPlayPage("https://www.bilibili.com/bangumi/play/ep307065")
}

func TestGetAllSeasonByMID(t *testing.T) {
	rs := GetAllSeasonByMIDURL("https://www.bilibili.com/bangumi/media/md28224394")
	fmt.Println(rs)
}

func TestBangumiDeal(t *testing.T) {
	BangumiDeal("https://www.bilibili.com/bangumi/play/ep307065")
}