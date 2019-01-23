package model

type GoodsImg struct {
	ID           int64  `db:"id"`
	GoodsID      int64  `db:"goods_id"`
	Name         string `db:"name"`
	URL          string `db:"path"`
	DisplayOrder int    `db:"display_order"`
}
