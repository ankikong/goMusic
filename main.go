package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/ankikong/goMusic/kugou"
	"github.com/ankikong/goMusic/songBean"

	"github.com/ankikong/goMusic/netease"
	"github.com/jedib0t/go-pretty/table"
)

func search(text string) {
	var result []songBean.SongData
	for _, rs := range netease.Search(text) {
		result = append(result, rs)
	}
	for _, rs := range kugou.Search(text) {
		result = append(result, rs)
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"编号", "歌名", "歌手", "专辑"})
	// fmt.Println("歌名\t歌手\t专辑")
	// fmt.Printf("%-26s%-26s%-26s\n", "歌名", "歌手", "专辑")
	index := 0
	for _, rs := range result {
		t.AppendRow(table.Row{index, rs.GetFileName(), rs.GetArtistName(), rs.GetAlbumName()})
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
	rss := result[num].GetUrl(320)
	fmt.Println(rss.SongUrl)
	Download(rss.SongUrl, rss.SongName, "")
}

func GetByNeteaseId(url string) {
	reg, _ := regexp.Compile(`\Wid=\d+`)
	ids := reg.FindAllString(url, -1)
	var id string
	if len(ids) == 0 {
		tmp := strings.Split(url, "/")
		for _, i := range tmp {
			if _, err := strconv.ParseInt(i, 10, 32); err == nil {
				id = i
				break
			}
		}
	} else {
		id = ids[0][4:]
	}
	rs := netease.GetSongUrl([]string{fmt.Sprint(id)}, 320)[0]
	fmt.Println("开始下载:", rs.SongName)
	Download(rs.SongUrl, rs.SongName, "")
}

func main() {
	var (
		url     string
		keyword string
	)
	flag.StringVar(&url, "url", "", "url of song")
	flag.StringVar(&keyword, "kw", "", "search keyword")
	flag.Parse()
	if len(keyword) > 0 {
		search(keyword)
	} else {
		if strings.Contains(url, "music.163.com") {
			GetByNeteaseId(url)
		}
	}
}
