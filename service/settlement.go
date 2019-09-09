package service

import (
	"gotrue/dao"
	"gotrue/dto/response"
	"gotrue/service/api"
)

var SettlementService api.ISettlementService

func InitSettlementService() {
	SettlementService = &Settlement{
		orderDao: dao.SaleOrderDao,
	}
}

type Settlement struct {
	orderDao *dao.SaleOrder
}

func (s *Settlement) ShowSettlementForSupplier(supplierID int64) (*response.SupplierSettlement, error) {
	orderSet, err := s.orderDao.QueryUnSettledOrdersBySupplier(supplierID)
	if err != nil {
		return nil, err
	}
	orderSet.Sum()
	return &response.SupplierSettlement{
		TotalOrder:       orderSet.Size(),
		TotalSalePrice:   orderSet.SalePrice(),
		TotalCostPrice:   orderSet.CostPrice(),
		TotalProfitPrice: orderSet.ProfitPrice(),
		TotalSettlePrice: orderSet.SettlePrice(),
	}, nil
}
