package cut

import (
	"fmt"
	"gotrue/model"
	"testing"

	"github.com/shopspring/decimal"
)

func TestCutCalculator(t *testing.T) {
	low := decimal.New(50, 0)
	origin := decimal.New(100, 0)
	discount := decimal.New(40, 0)
	calc := model.NewCutCalculator(low, origin, discount)
	for i := 1; i <= 50; i++ {
		cutoff := calc.RandomCut()
		fmt.Printf("%d time cut off %s total cut off %s \n", i, cutoff, calc.TotalCutoff())
	}
	fmt.Printf("final price is %s\n", calc.CurrentPrice())
}
