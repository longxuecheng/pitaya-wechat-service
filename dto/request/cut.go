package request

type CutoffInfo struct {
	GoodsID int64  `json:"goodsId" form:"goodsId"`
	UserID  string `json:"userId" form:"userId"`
}

type CutOrder struct {
	GoodsID int64
	UserID  int64
}
type ConsumeCutOrder struct {
	SaleOrderID int64
	CutOrderID  int64
}

type AssistCutoff struct {
	UserID   string `json:"userId"`
	HelperID int64  `json:"-"`
	GoodsID  int64  `json:"goodsId"`
}
