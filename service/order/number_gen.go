package order

import (
	"strconv"

	"github.com/bwmarrin/snowflake"
)

var DefaultGenerator = NewOrderNumberGenerator()

type OrderNumberGenerator struct {
	Node *snowflake.Node
}

func NewOrderNumberGenerator() *OrderNumberGenerator {
	node, _ := snowflake.NewNode(1)
	return &OrderNumberGenerator{
		Node: node,
	}
}

func (g *OrderNumberGenerator) GenerateNumber() string {
	return strconv.FormatInt(g.Node.Generate().Int64(), 10)
}
