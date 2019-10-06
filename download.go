package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func Download(url, name, path string) error {
	path = strings.TrimSpace(path)
	if len(path) == 0 {
		path = "./"
	} else if path[len(path)-1] != '/' {
		path = path + "/"
	}
	var ext string
	if strings.Contains(url, "mp3") {
		ext = "mp3"
	} else if strings.Contains(url, "flac") {
		ext = "flac"
	} else if strings.Contains(url, "aac") {
		ext = "aac"
	}
	name = name + "." + ext
	// reg, _ := regexp.Compile(`(/\|<>:*?")`)
	for _, i := range `/\|<>:*?"` {
		name = strings.ReplaceAll(name, string(i), " ")
	}

	rs, err := http.Get(url)
	if err != nil {
		log.Println("download fail:", err.Error())
		return err
	}
	file, err := os.Create(path + name)
	if err != nil {
		log.Println("create file fail:", err.Error())
		return err
	}
	fmt.Println("start download:" + name)
	buf := make([]byte, 262144)
	for {
		len, err := rs.Body.Read(buf)
		file.Write(buf[:len])
		if err != nil {
			break
		}
		fmt.Print(".")
	}
	file.Close()
	rs.Body.Close()
	return nil
}
