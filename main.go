package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
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
	rss := result[num].GetURL(320)
	fmt.Println(rss.SongURL)
	tool.Download(rss.SongURL, result[num].GetFileName(), "")
}

func GetByNeteaseId(URL string) {
	reg, _ := regexp.Compile(`\Wid=\d+`)
	ids := reg.FindAllString(URL, -1)
	var id string
	if len(ids) == 0 {
		tmp := strings.Split(URL, "/")
		for _, i := range tmp {
			if _, err := strconv.ParseInt(i, 10, 32); err == nil {
				id = i
				break
			}
		}
	} else {
		id = ids[0][4:]
	}
	rs := netease.GetSongURL([]string{fmt.Sprint(id)}, 320)[0]
	fmt.Println("开始下载:", rs.SongName)
	tool.Download(rs.SongURL, rs.SongName, "")
}

func GetByQQId(URL string) {
	reg, _ := regexp.Compile(`songid=\d+`)
	ids := reg.FindAllString(URL, -1)
	var id string
	if len(ids) != 0 {
		id = ids[0][7:]
	} else {
		fmt.Println("error: ", URL)
		return
	}
	rs := qq.GetSongURL(id, "320")
	fmt.Println("开始下载", rs.SongName)
	tool.Download(rs.SongURL, rs.SongName, "")
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
			GetByNeteaseId(URL)
		} else if strings.Contains(URL, "qq.com") {
			GetByQQId(URL)
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
