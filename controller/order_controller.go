package controller

import (
	"pitaya-wechat-service/api"
	"pitaya-wechat-service/dto/pagination"
	"pitaya-wechat-service/dto/request"
	"pitaya-wechat-service/facility/utils"
	"pitaya-wechat-service/middle_ware"
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
	err := c.BindJSON(&req)
	utils.CheckAndPanic(err)
	userID := middle_ware.MustGetCurrentUser(c)
	orderID, err := saleOrderServiceRf.Create(userID, req)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, orderID)
}

// ListSaleOrders 获取用户的所有订单
func ListSaleOrders(c *gin.Context) {
	req := pagination.PaginationRequest{}
	err := c.BindJSON(&req)
	utils.CheckAndPanic(err)
	userID := middle_ware.MustGetCurrentUser(c)
	result, err := saleOrderServiceRf.List(userID, req)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, result)
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
	middle_ware.SetResponseData(c, map[string]interface{}{
		"orderInfo":  info,
		"orderGoods": goodsList,
	})
}
