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
	StockID  int64           `json:"stockId"`
	Quantity decimal.Decimal `json:"quantity"`
}

type OrderExpressUpdate struct {
	SupplierID    int64  `json:"-"`
	OrderID       int64  `json:"orderId"`
	ExpressMethod string `json:"expressMethod"`
	ExpressNo     string `json:"expressNo"`
}
