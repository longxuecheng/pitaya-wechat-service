package response

import "github.com/shopspring/decimal"

// SupplierSettlement is settlement for supplier
type SupplierSettlement struct {
	TotalOrder       int             `json:"totalOrder"`
	TotalSalePrice   decimal.Decimal `json:"totalSalePrice"`
	TotalCostPrice   decimal.Decimal `json:"totalCostPrice"`
	TotalProfitPrice decimal.Decimal `json:"totalProfitPrice"`
	TotalSettlePrice decimal.Decimal `json:"totalSettlePrice"`
}
