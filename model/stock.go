package model

import (
	"database/sql"
	"gotrue/facility/utils"

	"github.com/shopspring/decimal"
)

// Stock 商品库存数据模型
type Stock struct {
	ID                int64           `db:"id"`
	SupplierID        int64           `db:"supplier_id"`
	Name              sql.NullString  `db:"name"`
	SaleUnitPrice     decimal.Decimal `db:"sale_unit_price"`
	CostUnitPrice     decimal.Decimal `db:"cost_unit_price"`
	AvailableQuantity decimal.Decimal `db:"available_quantity"`
	Specification     sql.NullString  `db:"specification"`
	GoodsID           int64           `db:"goods_id"`
	Splitable         bool            `db:"splitable"`
}

func (s *Stock) TableName() string {
	return "stock"
}

func (s *Stock) Columns() []string {
	return utils.TagValues(s, "db")
}

type StockSet struct {
	stocks []*Stock
}

func NewStockSet(stocks []*Stock) *StockSet {
	return &StockSet{
		stocks: stocks,
	}
}

// Map 是库存的集合
// 用在这里是因为这个函数于业务不相关，指示作为一个基础的model转换
func (s *StockSet) Map() map[int64]*Stock {
	stockMap := map[int64]*Stock{}
	for _, item := range s.stocks {
		stockMap[item.ID] = item
	}
	return stockMap
}

// GoodsIDs get distinct goods id list
func (s *StockSet) GoodsIDs() []int64 {
	distinctMap := map[int64]bool{}
	for _, item := range s.stocks {
		distinctMap[item.GoodsID] = true
	}
	goodsIDs := []int64{}
	for v := range distinctMap {
		goodsIDs = append(goodsIDs, v)
	}
	return goodsIDs
}
