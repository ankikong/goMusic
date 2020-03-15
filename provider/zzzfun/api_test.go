package zzzfun

import (
	"fmt"
	"testing"
)

func TestGetEps(t *testing.T) {
	rs := GetEps("193")
	fmt.Println(rs)
}
func TestGetVideoURL(t *testing.T) {
	rs := GetVideoURL("513-1-a")
	fmt.Println(rs)
}
