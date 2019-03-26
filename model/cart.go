package model

import (
	"github.com/shopspring/decimal"
)

var Table_Cart string = "cart"

// Cart 购物车数据库模型
type Cart struct {
	ID                   int64           `db:"id"`
	UserID               int64           `db:"user_id"`
	SessionID            string          `db:"session_id"`
	GoodsID              int64           `db:"goods_id"`
	GoodsSN              string          `db:"goods_sn"`
	StockID              int64           `db:"stock_id"`
	GoodsName            string          `db:"goods_name"`
	Quantity             decimal.Decimal `db:"quantity"`
	GoodsSpecDescription string          `db:"goods_spec_description"`
	GoodsSpecIDs         string          `db:"goods_spec_ids"`
	Checked              int8            `db:"checked"`
	ListPicURL           string          `db:"list_pic_url"`
}
