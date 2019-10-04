package express

import (
	"fmt"
	"testing"
)

func TestExpressInfo(t *testing.T) {
	expressInfo, err := ExpressService.ExpressInfo("EMS", "9896437530078")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("express detail is %+v", expressInfo)
	for i, e := range expressInfo.Traces {
		fmt.Printf("express trace %d is % v", i, e)
	}
}
