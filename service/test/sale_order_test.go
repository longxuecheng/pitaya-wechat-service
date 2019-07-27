package test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/bwmarrin/snowflake"
)

func TestGenerateOrderNo(t *testing.T) {
	// Create a new Node with a Node number of 1
	node, err := snowflake.NewNode(1)
	if err != nil {
		t.Error(err)
	}
	// Generate a snowflake ID.
	id := node.Generate()
	s := strconv.FormatInt(id.Int64(), 10)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(s)
}
