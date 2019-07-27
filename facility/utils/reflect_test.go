package utils

import (
	"fmt"
	"testing"
)

type tag struct {
	A string `db:"a"`
	B string `db:"b"`
}

func TestTagValues(t *testing.T) {
	fmt.Println(TagValues(&tag{}, "db"))
}
