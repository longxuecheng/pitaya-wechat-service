package model

import (
	"github.com/shopspring/decimal"
)

var Table_Cart string = "cart"

// Cart 购物车数据库模型
type Cart struct {
	ID                   int64           `db:"id" exclude:"true"`
	SupplierID           int64           `db:"supplier_id"`
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

type CartSet struct {
	items       []Cart
	supplierIDs []int64
	goodsIDs    []int64
	stockIDs    []int64
}

func NewCartSet(items []Cart) *CartSet {
	set := new(CartSet)
	set.items = items
	supplierIDMap := map[int64]bool{}
	goodsIDMap := map[int64]bool{}
	stockIDMap := map[int64]bool{}
	for _, item := range items {
		supplierIDMap[item.SupplierID] = true
		goodsIDMap[item.GoodsID] = true
		stockIDMap[item.StockID] = true
	}
	supplierIDs := []int64{}
	goodsIDs := []int64{}
	stockIDs := []int64{}
	for v := range supplierIDMap {
		supplierIDs = append(supplierIDs, v)
	}
	for v := range goodsIDMap {
		goodsIDs = append(goodsIDs, v)
	}
	for v := range stockIDMap {
		stockIDs = append(stockIDs, v)
	}
	set.goodsIDs = goodsIDs
	set.supplierIDs = supplierIDs
	set.stockIDs = stockIDs
	return set
}

func (s *CartSet) SupplierIDs() []int64 {
	return s.supplierIDs
}

func (s *CartSet) GoodsIDs() []int64 {
	return s.goodsIDs
}

func (s *CartSet) StockIDs() []int64 {
	return s.stockIDs
}

func (s *CartSet) List() []Cart {
	return s.items
}
