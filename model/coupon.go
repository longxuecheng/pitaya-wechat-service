package model

import (
	"errors"
	"gotrue/facility/slices"
	"gotrue/facility/utils"
	"time"

	"github.com/shopspring/decimal"
)

// ValidCouponTypes is valid coupon type list
var ValidCouponTypes = []CouponType{CouponTypeAllCategory, CouponTypeSpecialCategory, CouponTypeSpecialGoods}

const (
	CouponTypeAllCategory     CouponType = "All_Category"
	CouponTypeSpecialCategory CouponType = "Special_Category"
	CouponTypeSpecialGoods    CouponType = "Special_Goods"
)

const (
	ComposableTypeNone   ComposableType = "None"
	ComposableTypeCutoff ComposableType = "Cutoff"
)

type CouponType string

func NewCouponTypeFromString(str string) (CouponType, error) {
	has := false
	couponType := CouponType(str)
	for _, ct := range ValidCouponTypes {
		if couponType == ct {
			has = true
			break
		}
	}
	if !has {
		return couponType, errors.New("Invalid coupon type")
	}
	return couponType, nil
}

func (c CouponType) String() string {
	return string(c)
}

func (c CouponType) Title() string {
	switch c {
	case CouponTypeAllCategory:
		return "全品类"
	case CouponTypeSpecialCategory:
		return "指定品类"
	case CouponTypeSpecialGoods:
		return "指定商品"
	default:
		return "未知"
	}
}

// ComposableType is type of discount that can use together
type ComposableType string

type Coupon struct {
	ID             int64           `db:"id" insert:"false"`
	ActivityID     int64           `db:"activity_id"`
	CouponNo       string          `db:"coupon_no"`
	Price          decimal.Decimal `db:"price"`
	CreateTime     time.Time       `db:"create_time"`
	Consumed       bool            `db:"consumed"`
	Received       bool            `db:"received"`
	ConsumeTime    NullUTC8Time    `db:"consumed_time"`
	SaleOrderID    int64           `db:"sale_order_id"`
	ExpireTime     time.Time       `db:"expire_time"`
	Type           CouponType      `db:"type"`
	CategoryID     int64           `db:"category_id"`
	GoodsID        int64           `db:"goods_id"`
	ComposableType ComposableType  `db:"composable_type"`
	UserID         int64           `db:"user_id"`
}

func (c Coupon) TableName() string {
	return "coupon"
}

func (c *Coupon) Columns() []string {
	return utils.TagValues(c, "db")
}

type CouponList []*Coupon

func (l CouponList) CategoryIDList() []int64 {
	idList := make([]int64, len(l))
	for i, c := range l {
		idList[i] = c.ID
	}
	return slices.DistinctInt64(idList)
}

func (l CouponList) GoodsIDList() []int64 {
	idList := make([]int64, len(l))
	for i, c := range l {
		idList[i] = c.GoodsID
	}
	return slices.DistinctInt64(idList)
}
