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
	a.Key = "你好"
	rs, _ := json.Marshal(a)
	fmt.Println(string(rs), rs)
	// search("トゥルーエンド プレイヤー")
}
