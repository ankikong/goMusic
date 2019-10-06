package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ankikong/goMusic/kugou"
)

func TestKuGou(t *testing.T) {
	rs := kugou.Search("Hello Alone")
	for _, per := range rs {
		fmt.Println(per.GetUrl(320))
	}
	fmt.Println(rs)
}

type tt struct {
	Key string `json:"key"`
}

func TestSearch(t *testing.T) {
	a := new(tt)
	a.Key = "claris"
	rs, _ := json.Marshal(a)
	fmt.Println(string(rs), rs)
	// search("トゥルーエンド プレイヤー")
}

func TestQQId(t *testing.T) {
	GetByQQId("https://y.qq.com/n/yqq/song/0031Jhwu0ryf6Q.html")
	// GetByQQId("https://i.y.qq.com/v8/playsong.html?songid=105603683&source=yqq#wechat_redirect")
}
