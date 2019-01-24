package request

import (
	"github.com/shopspring/decimal"
)

// CartAddRequest is a request DTO of adding goods to user's cart
type CartAddRequest struct {
	UserID   int64           `json:"-"`
	GoodsID  int64           `json:"goodsId"`
	Quantity decimal.Decimal `json:"number"`
	StockID  int64           `json:"productId"`
}

type CartCheckRequest struct {
	ID        int64 `json:"id"`
	IsChecked int   `json:"isChecked"`
}
