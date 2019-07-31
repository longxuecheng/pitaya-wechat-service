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
	ID                   int64           `json:"id"`
	GoodsName            string          `json:"goods_name"`
	Quantity             decimal.Decimal `json:"number"`
	ListPicURL           string          `json:"list_pic_url"`
	RetailPrice          decimal.Decimal `json:"retail_price"`
	GoodsSpecDescription string          `json:"goods_specifition_name_value"`
}

type SaleOrderInfoDTO struct {
	ID            int64           `json:"id"`
	CreatedAt     string          `json:"add_time"`
	OrderNo       string          `json:"order_sn"`
	OrderAmt      decimal.Decimal `json:"actual_price"`
	Status        string          `json:"order_status_text"`
	FullRegion    string          `json:"full_region"`
	Address       string          `json:"address"`
	Consignee     string          `json:"consignee"`
	Mobile        string          `json:"mobile"`
	GoodsAmt      decimal.Decimal `json:"goods_amt"`
	ExpressMethod string          `json:"express_method"`
	ExpressNo     string          `json:"express_no"`
	ExpressFee    decimal.Decimal `json:"express_fee"`
}
