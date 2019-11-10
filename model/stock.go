package model

import (
	"database/sql"
	"gotrue/facility/utils"
	"strings"

	"github.com/shopspring/decimal"
)

const (
	StockStatusOnSale  = "ON_SALE"
	StockStatusOffSale = "OFF_SALE"
)

type StockMinMax struct {
	MinSalePrice decimal.Decimal `db:"min_sale_price"`
	MaxSalePrice decimal.Decimal `db:"max_sale_price"`
}

// Stock 商品库存数据模型
type Stock struct {
	ID                int64           `db:"id"`
	SupplierID        int64           `db:"supplier_id"`
	Name              sql.NullString  `db:"name"`
	ProfitPrice       decimal.Decimal `db:"profit_price" virtual:"true"`
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
	return utils.TagValues(s, "db", "virtual")
}

func (s *Stock) SpecIDs() []int64 {
	specStrings := strings.Split(s.Specification.String, "_")
	specIn64IDs := []int64{}
	for _, specString := range specStrings {
		specInt64ID, _ := utils.ParseInt64(specString)
		specIn64IDs = append(specIn64IDs, specInt64ID)
	}
	return specIn64IDs
}

type StockSet struct {
	stocks []*Stock
}

func NewStockSet(stocks []*Stock) *StockSet {
	return &StockSet{
		stocks: stocks,
	}
}

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

func (s *StockSet) SpecMap() StockSpecMap {
	m := StockSpecMap{}
	for _, item := range s.stocks {
		m[item.ID] = item.SpecIDs()
	}
	return m
}

func (s *StockSet) GetByGoods(goodsID int64) []*Stock {
	gStocks := []*Stock{}
	for _, item := range s.stocks {
		if goodsID == item.GoodsID {
			gStocks = append(gStocks, item)
		}
	}
	return gStocks
}

type StockSpecMap map[int64][]int64

func (m StockSpecMap) GetSpecs(stockID int64) []int64 {
	return m[stockID]
}
