package model

import (
	"gotrue/facility/utils"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ActivityCoupon struct {
	ID                int64           `db:"id"`
	ActivityID        int64           `db:"activity_id"`
	CouponType        CouponType      `db:"coupon_type"`
	ComposableType    ComposableType  `db:"composable_type"`
	CouponPrice       decimal.Decimal `db:"coupon_price"`
	TotalQuantity     int64           `db:"total_quantity"`
	AvailableQuantity int64           `db:"available_quantity"`
	CreateTime        time.Time       `db:"create_time"`
	CategoryID        int64           `db:"category_id"`
	GoodsID           int64           `db:"goods_id"`
}

func (c *ActivityCoupon) TableName() string {
	return "activity_coupon"
}

func (c *ActivityCoupon) Columns() []string {
	return utils.TagValues(c, "db")
}

func (c *ActivityCoupon) NewCoupon(userID int64) *Coupon {
	return &Coupon{
		UserID:         userID,
		CouponNo:       uuid.New().String(),
		Price:          c.CouponPrice,
		ActivityID:     c.ActivityID,
		Type:           c.CouponType,
		CategoryID:     c.CategoryID,
		GoodsID:        c.GoodsID,
		ComposableType: c.ComposableType,
		CreateTime:     time.Now(),
		ExpireTime:     time.Now().AddDate(0, 0, 7),
	}
}

func (c *ActivityCoupon) IsDrained() bool {
	return c.AvailableQuantity == 0
}
