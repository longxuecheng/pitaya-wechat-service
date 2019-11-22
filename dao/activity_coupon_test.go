package dao

import (
	"fmt"
	"testing"
	"time"
)

func TestDecreaseAvailQuantityByID(t *testing.T) {
	initActivityCouponDao()
	for i := 0; i < 100; i++ {
		go func() {
			err := ActivityCouponDao.DecreaseAvailQuantityByID(1, nil)
			if err != nil {
				fmt.Printf("%+v\n", err)
			}
		}()
	}
	time.Sleep(5 * time.Second)
}
