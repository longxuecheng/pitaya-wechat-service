package api

import "gotrue/dto/response"

// ISettlementService is settlement service
type ISettlementService interface {

	// ShowSettlementForSupplier show settle money for supplier
	ShowSettlementForSupplier(supplierID int64) (*response.SupplierSettlement, error)
}
