package tool

import (
	"bytes"
	"strings"
	"testing"
)

func TestCBC(t *testing.T) {
	data := []byte("raw")
	encrypt := AesEncryptCBC(data, []byte("0102030405060708"), []byte("0CoJUm6Qyw8W8jud"))
	decrypt := AesDecryptCBC(encrypt, []byte("0102030405060708"), []byte("0CoJUm6Qyw8W8jud"))
	if bytes.Compare(data, decrypt) != 0 {
		t.Error("CBC error", data, decrypt)
	}
}

func TestMD5(t *testing.T) {
	rs := MD5("aaa")
	rs = strings.ToUpper(rs)
	if strings.Compare(rs, "47BCE5C74F589F4867DBD57E9CA9F808") != 0 {
		t.Error("aaa", rs, "47BCE5C74F589F4867DBD57E9CA9F808")
	}
}
