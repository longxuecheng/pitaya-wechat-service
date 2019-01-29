package api

import (
	"pitaya-wechat-service/dto/pagination"
	"pitaya-wechat-service/dto/request"
)

// ISaleOrderService is a order service interface
type ISaleOrderService interface {
	Create(req request.SaleOrderAddRequest) (id int64, err error)
	List(userID int64, req pagination.PaginationRequest) (page pagination.PaginationResonse, err error)
}
