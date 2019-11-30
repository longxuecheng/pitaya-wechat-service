package request

import (
	"gotrue/facility/errors"

	"github.com/shopspring/decimal"
)

type SaleOrderAddRequest struct {
	AddressID int64 `json:"addressId"`
	CouponID  int64 `json:"couponId"`
}

func (s *SaleOrderAddRequest) Validate() error {
	if s.AddressID == 0 {
		return errors.NewWithCodef("InvalidAddress", "地址无效")
	}
	return nil
}

type SaleOrderQuickAddRequest struct {
	SaleOrderAddRequest
	StockID  int64           `json:"stockId"`
	Quantity decimal.Decimal `json:"quantity"`
}

type OrderExpressUpdate struct {
	OrderID       int64  `json:"orderId"`
	ExpressMethod string `json:"expressMethod"`
	ExpressNo     string `json:"expressNo"`
}

type QueryWechatPayResult struct {
	OrderID  int64  `form:"orderId"`
	PrepayID string `form:"prepayId"`
	UserID   int64  `form:"-"`
}
