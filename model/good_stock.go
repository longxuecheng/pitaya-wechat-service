package model

import (
	"database/sql"

	"github.com/shopspring/decimal"
)

var Table_Stock = "stock"

// GoodsStock 商品库存数据模型
type GoodsStock struct {
	ID                int64           `db:"id"`
	Name              sql.NullString  `db:"name"`
	SaleUnitPrice     decimal.Decimal `db:"sale_unit_price"`
	CostUnitPrice     decimal.Decimal `db:"cost_unit_price"`
	AvailableQuantity decimal.Decimal `db:"available_quantity"`
	Specification     sql.NullString  `db:"specification"`
	GoodsID           int64           `db:"goods_id"`
}

// NewStockMap transform an array to map key is id value is a pointer to item
// 用在这里是因为这个函数于业务不相关，指示作为一个基础的model转换
func NewStockMap(items []*GoodsStock) map[int64]*GoodsStock {
	stockMap := map[int64]*GoodsStock{}
	for _, item := range items {
		stockMap[item.ID] = item
	}
	return stockMap
}
