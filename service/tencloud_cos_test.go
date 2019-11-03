package service

import (
	"bufio"
	"os"
	"testing"
)

func TestPushImageObject(t *testing.T) {
	cos := NewCosService()
	f, err := os.Open("./wxacode.jpeg")
	defer f.Close()
	if err != nil {
		t.Error(err)
	}
	info, _ := f.Stat()
	buff := make([]byte, info.Size())
	reader := bufio.NewReader(f)
	reader.Read(buff)
	err = cos.PushImageObject("wxacode.jpeg", buff)
	if err != nil {
		t.Error(err)
	}
}
