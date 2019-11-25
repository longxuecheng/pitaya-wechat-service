package coupon

import (
	"context"
	"gotrue/dao"
	"gotrue/service/api"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"
)

func TestSendCoupon(t *testing.T) {
	dao.Init()
	price, _ := decimal.NewFromString("5.00")
	couponService := GetCouponService()
	req := &api.SendCouponRequest{
		TargetUserID:         12,
		CouponType:           "All_Category",
		CouponPrice:          price,
		CouponQuantity:       2,
		CategoryID:           0,
		GoodsID:              0,
		ExpireUnixTime:       time.Now().AddDate(0, 0, 1).Unix(),
		ComposableWithCutoff: false,
	}
	err := couponService.SendCouponToUser(context.Background(), req)
	assert.Equal(t, nil, err, "断言赠送优惠券成功")
}
