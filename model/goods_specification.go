package model

var Table_Goods_Specification = "goods_specification"

// GoodsSpecification 定义了某商品的规格
type GoodsSpecification struct {
	ID              int64  `db:"id"`
	GoodsID         int64  `db:"goods_id"`
	SpecificationID int64  `db:"specification_id"`
	Value           string `db:"value"`
	PicURL          string `db:"pic_url"`
}
