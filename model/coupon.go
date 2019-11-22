package model

import (
	"gotrue/facility/slices"
	"gotrue/facility/utils"
	"time"

	"github.com/shopspring/decimal"
)

const (
	CounponTypeAllCategory     CounponType = "All_Category"
	CounponTypeSpecialCategory CounponType = "Special_Category"
	CounponTypeSpecialGoods    CounponType = "Special_Goods"
)

const (
	ComposableTypeNone   ComposableType = "None"
	ComposableTypeCutoff ComposableType = "Cutoff"
)

type CounponType string

func (c CounponType) String() string {
	return string(c)
}

func (c CounponType) Title() string {
	switch c {
	case CounponTypeAllCategory:
		return "全品类"
	case CounponTypeSpecialCategory:
		return "指定品类"
	case CounponTypeSpecialGoods:
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
	ConsumeTime    NullUTC8Time    `db:"consumed_time"`
	SaleOrderID    int64           `db:"sale_order_id"`
	ExpireTime     time.Time       `db:"expire_time"`
	Type           CounponType     `db:"type"`
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
