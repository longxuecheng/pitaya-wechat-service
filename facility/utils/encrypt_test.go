package utils

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	InitEncryptor()
	fmt.Println(EncodeIn64(3))
}
