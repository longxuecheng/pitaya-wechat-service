package model

import (
	"database/sql"
	"gotrue/facility/utils"
	"time"

	"github.com/shopspring/decimal"
)

type GoodsStatus string

const (
	GoodsStatusOnSale  GoodsStatus = "ON_SALE"
	GoodsStatusOffSale GoodsStatus = "OFF_SALE"
)

type Goods struct {
	ID               int64           `db:"id"`
	Name             string          `db:"name"`
	ProducingArea    string          `db:"producing_area"`
	CategoryID       int             `db:"category_id"`
	BriefDescription sql.NullString  `db:"brief_description"`
	Description      sql.NullString  `db:"description"`
	Status           string          `db:"status"`
	IsDelete         int8            `db:"is_delete"`
	SupplierID       int64           `db:"supplier_id"`
	CreateTime       *time.Time      `db:"create_time"`
	ListPicURL       sql.NullString  `db:"list_pic_url"`
	RetailPrice      decimal.Decimal `db:"retail_price"`
}

func (g *Goods) TableName() string {
	return "goods"
}

func (g *Goods) Columns() []string {
	return utils.TagValues(g, "db")
}

type GoodsSet struct {
	items []*Goods
}

func NewGoodsSet(items []*Goods) *GoodsSet {
	return &GoodsSet{
		items: items,
	}
}

func (s *GoodsSet) Map() map[int64]*Goods {
	goodsMap := map[int64]*Goods{}
	for _, item := range s.items {
		goodsMap[item.ID] = item
	}
	return goodsMap
}
