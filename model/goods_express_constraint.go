package model

import (
	"gotrue/facility/utils"

	"github.com/shopspring/decimal"
	"go.planetmeican.com/manage/paperwork-facility/reflect_util"
)

// GoodsExpressConstraint is constraint of goods express
type GoodsExpressConstraint struct {
	ID         int64           `db:"id" omitinsert:"true"`
	GoodsID    int64           `db:"goods_id"`
	StockID    int64           `db:"stock_id"`
	IsFree     bool            `db:"is_free"`
	ProvinceID int             `db:"province_id"`
	Reachable  bool            `db:"reachable"`
	ExpressFee decimal.Decimal `db:"express_fee"`
}

func (g *GoodsExpressConstraint) TableName() string {
	return "goods_express_constraint"
}

func (g *GoodsExpressConstraint) Columns() []string {
	return reflect_util.TagValues(g, "db")
}

func (g *GoodsExpressConstraint) SetFeeFromCsv(record []string) error {
	isFree := false
	reachable := false
	provinceID, err := utils.ParseInt64(record[0])
	if err != nil {
		return err
	}
	g.ProvinceID = int(provinceID)
	if record[4] == "1" || record[4] == "是" {
		isFree = true
	}
	g.IsFree = isFree
	expressFee, err := decimal.NewFromString(record[2])
	if err != nil {
		return err
	}
	g.ExpressFee = expressFee
	if record[3] == "1" || record[3] == "是" {
		reachable = true
	}
	g.Reachable = reachable
	return nil
}
