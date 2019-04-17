package request

import (
	"github.com/shopspring/decimal"
)

// CashierPreview 结算台结算预览请求
type CashierPreview struct {
	StockID   int64           `json:"stockId"`
	Quantity  decimal.Decimal `json:"quantity"`
	AddressID int64           `json:"addressId"`
}
