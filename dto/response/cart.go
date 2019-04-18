package response

import (
	"github.com/shopspring/decimal"
)

// CartItem 购物车列表条目
type CartItem struct {
	ID                   int64           `json:"id"`
	UserID               int64           `json:"user_id"`
	SessionID            string          `json:"session_id"`
	GoodsSN              string          `json:"goods_sn"`
	StockID              int64           `json:"product_id"`
	GoodsID              int64           `json:"goods_id"`
	GoodsName            string          `json:"goods_name"`
	MarketPrice          decimal.Decimal `json:"market_price"`
	RetailPrice          decimal.Decimal `json:"retail_price"`
	Quantity             decimal.Decimal `json:"number"`
	GoodsSpecDescription string          `json:"goods_specifition_name_value"`
	GoodsSpecIDs         string          `json:"goods_specification_ids"`
	Checked              int8            `json:"checked"`
	ListPicURL           string          `json:"list_pic_url"`
}

type CartSummary struct {
	CartList  []CartItem `json:"cartList"`
	CartTotal CartTotal  `json:"cartTotal"`
}

type CartTotal struct {
	GoodsCount         int64  `json:"goodsCount"`
	CheckedGoodsCount  int    `json:"checkedGoodsCount"`
	CheckedGoodsAmount string `json:"checkedGoodsAmount"`
}
