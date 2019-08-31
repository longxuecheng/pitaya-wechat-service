package response

import (
	"github.com/shopspring/decimal"
)

// GoodsStockDTO 商品库存数据DTO模型
type GoodsStockDTO struct {
	ID                    int64           `json:"id"`
	GoodsID               int64           `json:"goodsId"`
	SaleUnitPrice         decimal.Decimal `json:"saleUnitPrice"`
	AvailableQuantity     decimal.Decimal `json:"goods_number"`
	SpecificationIDs      []int64         `json:"specificationIds"`
	GoodsSpecificationIDs string          `json:"goods_specification_ids"`
}
