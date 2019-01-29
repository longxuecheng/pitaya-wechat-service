package response

import (
	"github.com/shopspring/decimal"
)

type SaleOrderItemDTO struct {
	ID        int64               `json:"id"`
	OrderNo   string              `json:"order_sn"`
	Status    string              `json:"order_status_text"`
	OrderAmt  decimal.Decimal     `json:"actual_price"`
	GoodsList []SaleOrderGoodsDTO `json:"goodsList"`
}

type SaleOrderGoodsDTO struct {
	ID         int64           `json:"id"`
	GoodsName  string          `json:"goods_name"`
	Quantity   decimal.Decimal `json:"number"`
	ListPicURL string          `json:"list_pic_url"`
}

type SaleOrderInfoDTO struct {
}
