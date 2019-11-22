package dao

import (
	"database/sql"
	"gotrue/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/shopspring/decimal"
)

func newCoupon() *model.Coupon {
	price, _ := decimal.NewFromString("22.00")
	return &model.Coupon{
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
}

func TestTxExecutor(t *testing.T) {
	initCouponDao()
	initActivityCouponDao()
	txExecutor, err := GetEasyDB().NewTxExecutor()
	assert.Equal(t, nil, err)
	grabCouponFunc := func(input interface{}, tx *sql.Tx) (interface{}, error) {
		return nil, ActivityCouponDao.DecreaseAvailQuantityByID(1, tx)
	}
	createCouponFunc := func(input interface{}, tx *sql.Tx) (interface{}, error) {
		coupon := newCoupon()
		return CouponDao.CreateCoupon(coupon, tx)
	}
	txExecutor.AppendFunc(grabCouponFunc)
	txExecutor.AppendFunc(createCouponFunc)
	err = txExecutor.Execute()
	assert.Equal(t, ErrorNowRowAffected, err, "断言没有抢到优惠券")
}
