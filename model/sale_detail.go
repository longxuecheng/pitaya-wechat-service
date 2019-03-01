package model

import (
	"database/sql"

	"github.com/shopspring/decimal"
)

var Table_Sale_Detail string = "sale_detail"

// SaleDetail 是model.SaleOrder的订单明细
type SaleDetail struct {
	ID                   int64           `db:"id" exclude:"true"`
	OrderID              int64           `db:"order_id"`
	StockID              int64           `db:"stock_id"`
	GoodsID              int64           `db:"goods_id"`
	GoodsName            string          `db:"goods_name"`
	Quantity             decimal.Decimal `db:"quantity"`
	CostUnitPrice        decimal.Decimal `db:"cost_unit_price"`
	SaleUnitPrice        decimal.Decimal `db:"sale_unit_price"`
	GoodsSpecIDs         string          `db:"goods_spec_ids"`
	GoodsSpecDescription string          `db:"goods_spec_description"`
	ListPicURL           sql.NullString  `db:"list_pic_url"`
}
