package tool

import "testing"

func TestMergeFLV(t *testing.T) {
	MergeFLV("E:/tmp/out.flv", []string{"E:/tmp/1.flv", "E:/tmp/2.flv"})
}
