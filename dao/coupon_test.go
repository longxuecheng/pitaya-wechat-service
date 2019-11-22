package dao

import (
	"gotrue/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/shopspring/decimal"
)

func TestCreateCoupon(t *testing.T) {
	initCouponDao()
	price, _ := decimal.NewFromString("22.00")
	coupon := &model.Coupon{
		Price:          price,
		ActivityID:     1,
		CouponNo:       "coupon-0001",
		CreateTime:     time.Now(),
		ExpireTime:     time.Now().AddDate(0, 0, 7),
		Type:           model.CounponTypeSpecialCategory,
		CategoryID:     1,
		GoodsID:        0,
		ComposableType: model.ComposableTypeNone,
		UserID:         1,
	}
	id, err := CouponDao.CreateCoupon(coupon, nil)
	assert.Equal(t, nil, err)
	assert.True(t, id > 0)
}
