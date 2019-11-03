package request

type CutOrder struct {
	GoodsID int64 `json:"goodsId"`
	StockID int64 `json:"stockId"`
	UserID  int64 `json:"-"`
}

type ConsumeCutOrder struct {
	SaleOrderID int64
	CutOrderID  int64
}

type CreateCutDetail struct {
	CutNo  string `json:"cutNo"`
	UserID int64  `json:"-"`
}
