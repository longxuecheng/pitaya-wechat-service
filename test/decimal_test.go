package test

import (
	"testing"

	"github.com/shopspring/decimal"
)

type Struct struct {
	Number decimal.Decimal
}

func TestDecimalAdd(t *testing.T) {
	n, err := decimal.NewFromString("-123.343")
	if err != nil {
		t.Error(err)
	}
	t.Log(n.String())
	d2, _ := decimal.NewFromString("200.955")
	sum := n.Add(d2)
	t.Logf("d1 add d2 = %s", sum.String())
}
