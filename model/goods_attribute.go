package model

type GoodsAttribute struct {
	ID          int64  `db:"id"`
	GoodsID     int64  `db:"goods_id"`
	AttributeID int64  `db:"attribute_id"`
	Value       string `db:"value"`
}
