package test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/google/uuid"
)

func TestEscape(t *testing.T) {
	fmt.Println(url.QueryEscape("Asia/Shanghai"))
}

func TestCreateCode(t *testing.T) {
	uuid := uuid.New()
	fmt.Println(uuid.String())
}
