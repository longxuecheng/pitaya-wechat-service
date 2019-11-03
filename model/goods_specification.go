package model

var Table_Goods_Specification = "goods_specification"

// GoodsSpecification 定义了某商品的规格
type GoodsSpecification struct {
	ID              int64  `db:"id"`
	GoodsID         int64  `db:"goods_id"`
	SpecificationID int64  `db:"specification_id"`
	Value           string `db:"value"`
	Is_Delete       bool   `db:"is_delete"`
	PicURL          string `db:"pic_url"`
}

type GoodsSpecMap map[int64]*GoodsSpecification

func (m GoodsSpecMap) Get(specID int64) *GoodsSpecification {
	return m[specID]
}

type GoodsSpecSet struct {
	items []*GoodsSpecification
}

func NewGoodsSpecSet(items []*GoodsSpecification) *GoodsSpecSet {
	return &GoodsSpecSet{
		items: items,
	}
}

func (s *GoodsSpecSet) Map() GoodsSpecMap {
	result := GoodsSpecMap{}
	for _, item := range s.items {
		result[item.ID] = item
	}
	return result
}
