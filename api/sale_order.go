package api

import (
	"pitaya-wechat-service/dto/pagination"
	"pitaya-wechat-service/dto/request"
	"pitaya-wechat-service/dto/response"
)

// ISaleOrderService is a order service interface
type ISaleOrderService interface {
	Create(req request.SaleOrderAddRequest) (id int64, err error)
	List(userID int64, req pagination.PaginationRequest) (page pagination.PaginationResonse, err error)
	Info(orderID int64) (response.SaleOrderInfoDTO, error)
	ListGoods(orderID int64) ([]response.SaleOrderGoodsDTO, error)
}
