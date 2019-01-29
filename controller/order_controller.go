package controller

import (
	"pitaya-wechat-service/api"
	"pitaya-wechat-service/dto/pagination"
	"pitaya-wechat-service/dto/request"
	"pitaya-wechat-service/facility/utils"
	"pitaya-wechat-service/service"

	"github.com/gin-gonic/gin"
)

var (
	saleOrderServiceRf api.ISaleOrderService = service.SaleOrderServiceInstance()
)

// SubmitSaleOrder 提交销售订单
func SubmitSaleOrder(c *gin.Context) {
	req := request.SaleOrderAddRequest{}
	req.UserID = 0
	err := c.BindJSON(&req)
	utils.CheckAndPanic(err)
	orderID, err := saleOrderServiceRf.Create(req)
	utils.CheckAndPanic(err)
	c.Set("data", orderID)
}

func ListSaleOrders(c *gin.Context) {
	req := pagination.PaginationRequest{}
	err := c.BindJSON(&req)
	utils.CheckAndPanic(err)
	result, err := saleOrderServiceRf.List(0, req)
	utils.CheckAndPanic(err)
	c.Set("data", result)
}
