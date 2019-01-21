package dto

import (
	"github.com/shopspring/decimal"
)

type GoodsItemDTO struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	PicURL      string          `json:"picUrl"`
	RetailPrice decimal.Decimal `json:"retailPrice"`
}
