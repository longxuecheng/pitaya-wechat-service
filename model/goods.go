package model

import (
	"database/sql"
	"database/sql/driver"
	"gotrue/facility/utils"
	"time"

	"github.com/shopspring/decimal"
)

type GoodsStatus string

func (s GoodsStatus) String() string {
	return string(s)
}

func (s GoodsStatus) Value() (driver.Value, error) {
	return s.String(), nil
}

const (
	GoodsStatusOnSale  GoodsStatus = "ON_SALE"
	GoodsStatusOffSale GoodsStatus = "OFF_SALE"
	GoodsStatusPreSale GoodsStatus = "PRE_SALE"
)

var GoodsStatusMap = map[GoodsStatus]string{
	GoodsStatusOnSale:  "限时热销",
	GoodsStatusOffSale: "下架",
	GoodsStatusPreSale: "敬请期待"}

type Goods struct {
	ID               int64           `db:"id"`
	Name             string          `db:"name"`
	ProducingArea    string          `db:"producing_area"`
	ProdAreaLng      float32         `db:"prod_area_lng"`
	ProdAreaLat      float32         `db:"prod_area_lat"`
	CategoryID       int             `db:"category_id"`
	BriefDescription sql.NullString  `db:"brief_description"`
	Description      sql.NullString  `db:"description"`
	Status           GoodsStatus     `db:"status"`
	IsDelete         int8            `db:"is_delete"`
	SupplierID       int64           `db:"supplier_id"`
	CreateTime       *time.Time      `db:"create_time"`
	ListPicURL       sql.NullString  `db:"list_pic_url"`
	CardPicURL       string          `db:"card_pic_url"`
	RetailPrice      decimal.Decimal `db:"retail_price"`
}

func (g *Goods) TableName() string {
	return "goods"
}

func (g *Goods) Columns() []string {
	return utils.TagValues(g, "db")
}

func (g *Goods) StatusName() string {
	return GoodsStatusMap[g.Status]
}

type GoodsMap map[int64]*Goods

func (m GoodsMap) Get(id int64) *Goods {
	return m[id]
}

type GoodsList []*Goods

func (l GoodsList) GoodsMap() GoodsMap {
	goodsMap := GoodsMap{}
	for _, goods := range l {
		goodsMap[goods.ID] = goods
	}
	return goodsMap
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

func (s *GoodsSet) GoodsIDs() []int64 {
	distinctMap := map[int64]bool{}
	for _, item := range s.items {
		distinctMap[item.ID] = true
	}
	goodsIDs := []int64{}
	for v := range distinctMap {
		goodsIDs = append(goodsIDs, v)
	}
	return goodsIDs
}
