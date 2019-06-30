package api

import (
	"gotrue/dto/pagination"
	"gotrue/dto/request"
	"gotrue/dto/response"
)

// ISaleOrderService is a order service interface
type ISaleOrderService interface {
	// Create order from selected cart items
	Create(userID int64, req request.SaleOrderAddRequest) (id int64, err error)
	// QuickCreate create order from selected goods
	QuickCreate(req request.SaleOrderQuickAddRequest) (int64, error)
	List(userID int64, req pagination.PaginationRequest) (page pagination.PaginationResonse, err error)
	Info(orderID int64) (response.SaleOrderInfoDTO, error)
	ListGoods(orderID int64) ([]response.SaleOrderGoodsDTO, error)
}
