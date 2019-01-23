package model

import (
	"database/sql"

	"github.com/shopspring/decimal"
)

var Table_Stock = "stock"

// GoodsStock 商品库存数据模型
type GoodsStock struct {
	ID                int64           `db:"id"`
	SaleUnitPrice     decimal.Decimal `db:"sale_unit_price"`
	AvailableQuantity decimal.Decimal `db:"available_quantity"`
	Specification     sql.NullString  `db:"specification"`
	GoodsID           int64           `db:"goods_id"`
}
