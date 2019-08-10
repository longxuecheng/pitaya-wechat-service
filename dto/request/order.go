package request

import (
	"gotrue/dto/pagination"

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

type orderType string

const (
	All      orderType = "all"
	Finished orderType = "finished"
	Created  orderType = "created"
	Sent     orderType = "sent"
)

type OrderListRequest struct {
	pagination.PaginationRequest
	Type orderType `json:"type"`
}

type OrderExpressUpdate struct {
	SupplierID    int64  `json:"-"`
	OrderID       int64  `json:"orderId"`
	ExpressMethod string `json:"expressMethod"`
	ExpressNo     string `json:"expressNo"`
}

type QueryWechatPayResult struct {
	OrderID  int64  `form:"orderId"`
	PrepayID string `form:"prepayId"`
	UserID   int64  `form:"-"`
}
