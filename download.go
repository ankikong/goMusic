package main

import (
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
	tmp := strings.Split(url, ".")
	ext := tmp[len(tmp)-1]
	name = name + "." + ext
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
	buf := make([]byte, 262144)
	for {
		len, err := rs.Body.Read(buf)
		file.Write(buf[:len])
		if err != nil {
			break
		}
	}
	file.Close()
	rs.Body.Close()
	return nil
}
