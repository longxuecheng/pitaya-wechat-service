package controller

import (
	"pitaya-wechat-service/api"
	"pitaya-wechat-service/dto/pagination"
	"pitaya-wechat-service/dto/request"
	"pitaya-wechat-service/facility/utils"
	"pitaya-wechat-service/service"
	"strconv"

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

// ListSaleOrders 获取用户的所有订单
func ListSaleOrders(c *gin.Context) {
	req := pagination.PaginationRequest{}
	err := c.BindJSON(&req)
	utils.CheckAndPanic(err)
	result, err := saleOrderServiceRf.List(0, req)
	utils.CheckAndPanic(err)
	c.Set("data", result)
}

// SaleOrderInfo 获取订单详情
func SaleOrderInfo(c *gin.Context) {
	orderIDStr := c.Query("orderId")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	utils.CheckAndPanic(err)
	info, err := saleOrderServiceRf.Info(orderID)
	utils.CheckAndPanic(err)
	goodsList, err := saleOrderServiceRf.ListGoods(orderID)
	utils.CheckAndPanic(err)
	c.Set("data", map[string]interface{}{
		"orderInfo":  info,
		"orderGoods": goodsList,
	})
}
