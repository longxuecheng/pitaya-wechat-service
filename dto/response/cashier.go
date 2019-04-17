package response

import "github.com/shopspring/decimal"

// Cashier is a statistic information of goods amount and price
type Cashier struct {
	ExpressFee      string        `json:"expressFee"`
	GoodsTotalPrice string        `json:"goodsTotalPrice"`
	OrderTotalPrice string        `json:"orderTotalPrice"`
	Items           []CashierItem `json:"items"`
}

// CashierItem is the checkout unit of cashier
type CashierItem struct {
	// GoodsSN              string          `json:"goods_sn"`
	StockID              int64           `json:"product_id"`
	GoodsID              int64           `json:"goods_id"`
	GoodsName            string          `json:"goods_name"`
	RetailPrice          decimal.Decimal `json:"retail_price"`
	Quantity             decimal.Decimal `json:"number"`
	GoodsSpecDescription string          `json:"goods_specifition_name_value"`
	GoodsSpecIDs         string          `json:"goods_specification_ids"`
	ListPicURL           string          `json:"list_pic_url"`
}
