package api

import (
	"pitaya-wechat-service/dto/pagination"
	"pitaya-wechat-service/dto/request"
	"pitaya-wechat-service/dto/response"
)

// ISaleOrderService is a order service interface
type ISaleOrderService interface {
	Create(userID int64, req request.SaleOrderAddRequest) (id int64, err error)
	QuickCreate(req request.SaleOrderQuickAddRequest) (int64, error)
	// PreCreate show detail about order to be created
	PreCreate(req request.SaleOrderQuickAddRequest)
	List(userID int64, req pagination.PaginationRequest) (page pagination.PaginationResonse, err error)
	Info(orderID int64) (response.SaleOrderInfoDTO, error)
	ListGoods(orderID int64) ([]response.SaleOrderGoodsDTO, error)
}
