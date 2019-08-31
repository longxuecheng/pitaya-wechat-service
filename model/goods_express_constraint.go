package model

import (
	"github.com/shopspring/decimal"
	"go.planetmeican.com/manage/paperwork-facility/reflect_util"
)

// GoodsExpressConstraint is constraint of goods express
type GoodsExpressConstraint struct {
	ID         int64           `db:"id"`
	GoodsID    int64           `db:"goods_id"`
	StockID    int64           `db:"stock_id"`
	IsFree     bool            `db:"is_free"`
	ProvinceID int             `db:"province_id"`
	ExpressFee decimal.Decimal `db:"express_fee"`
}

func (g *GoodsExpressConstraint) TableName() string {
	return "goods_express_constraint"
}

func (g *GoodsExpressConstraint) Columns() []string {
	return reflect_util.TagValues(g, "db")
}
