package request

import (
	"github.com/shopspring/decimal"
)

type SaleOrderAddRequest struct {
	AddressID int64 `json:"addressId"`
	CouponID  int64 `json:"couponId"`
}

type SaleOrderQuickAddRequest struct {
	SaleOrderAddRequest
	UserID   int64           `json:"-"`
	StockID  int64           `json:"stockId"`
	Quantity decimal.Decimal `json:"quantity"`
}