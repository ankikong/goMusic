package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ankikong/goMusic/tool"

	"github.com/ankikong/goMusic/provider/bilibili"

	"github.com/ankikong/goMusic/provider/qq"

	"github.com/ankikong/goMusic/provider/kugou"
	"github.com/ankikong/goMusic/provider/songbean"

	"github.com/ankikong/goMusic/provider/netease"
	"github.com/jedib0t/go-pretty/table"
)

func search(text string) {
	var result []songbean.SongData
	for _, rs := range netease.Search(text) {
		result = append(result, rs)
	}
	for _, rs := range kugou.Search(text) {
		result = append(result, rs)
	}
	for _, rs := range qq.Search(text) {
		result = append(result, rs)
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"编号", "歌名", "歌手", "专辑", "来源"})
	// fmt.Println("歌名\t歌手\t专辑")
	// fmt.Printf("%-26s%-26s%-26s\n", "歌名", "歌手", "专辑")
	index := 0
	for _, rs := range result {
		t.AppendRow(table.Row{index, rs.GetFileName(), rs.GetArtistName(), rs.GetAlbumName(), rs.GetSource()})
		index++
		// fmt.Printf("%-26s%-26s%-26s\n", rs.GetFileName(), rs.GetArtistName(), rs.GetAlbumName())
	}
	var input string
	var num uint64
	// stdin := bufio.NewReader(os.Stdin)
	t.Render()
	for {
		fmt.Print("please input integer(input q to quit):")
		// fmt.Println(input)
		fmt.Scan(&input)
		// fmt.Fscan(stdin, input)
		input = strings.TrimSpace(input)
		input = strings.ToLower(input)
		if input == "q" {
			return
		}
		nums, err := strconv.ParseUint(input, 10, 20)
		if err != nil {
			continue
		}
		num = nums
		break
	}
	if num >= 0 && num < uint64(len(result)) {
		rss := result[num].GetURL(320)
		fmt.Println(rss.SongURL)
		tool.Download(rss.SongURL, result[num].GetFileName(), "", true)
	} else {
		panic("index out of range")
	}
}

func main() {
	var (
		URL     string
		keyword string
	)
	flag.StringVar(&URL, "url", "", "url of song")
	flag.StringVar(&keyword, "kw", "", "search keyword")
	flag.Parse()
	if len(keyword) > 0 {
		search(keyword)
	} else {
		if strings.Contains(URL, "music.163.com") {
			netease.GetByURL(URL)
		} else if strings.Contains(URL, "qq.com") {
			qq.GetByURL(URL)
		} else if strings.Contains(URL, "bilibili") {
			if strings.Contains(URL, "video") {
				bilibili.Deal(URL)
			} else if strings.Contains(URL, "bangumi") {
				bilibili.BangumiDeal(URL)
			} else {
				fmt.Println("unsupport url")
			}
		}
	}
}
