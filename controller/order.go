package controller

import (
	"gotrue/dto/pagination"
	"gotrue/dto/request"
	"gotrue/facility/utils"
	"gotrue/middle_ware"
	"gotrue/service"
	"gotrue/service/express"
	"gotrue/service/wechat"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SubmitSaleOrder 提交销售订单
func SubmitSaleOrder(c *gin.Context) {
	req := request.SaleOrderAddRequest{}
	err := c.BindJSON(&req)
	utils.CheckAndPanic(err)
	userID := middle_ware.MustGetCurrentUser(c)
	orderID, err := service.SaleOrderServiceInstance().Create(userID, req)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, orderID)
}

func QuickSubmitOrder(c *gin.Context) {
	req := request.SaleOrderQuickAddRequest{}
	err := c.BindJSON(&req)
	utils.CheckAndPanic(err)
	userID := middle_ware.MustGetCurrentUser(c)
	orderID, err := service.SaleOrderServiceInstance().QuickCreate(userID, req)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, orderID)
}

// ListSaleOrders 获取用户的所有订单
func ListSaleOrders(c *gin.Context) {
	req := pagination.PaginationRequest{}
	err := c.BindJSON(&req)
	utils.CheckAndPanic(err)
	userID := middle_ware.MustGetCurrentUser(c)
	result, err := service.SaleOrderServiceInstance().List(userID, req)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, result)
}

// SaleOrderInfo 获取订单详情
func SaleOrderInfo(c *gin.Context) {
	orderID := bindSaleOrderIDFromQuery(c)
	info, err := service.SaleOrderServiceInstance().Info(orderID)
	utils.CheckAndPanic(err)
	goodsList, err := service.SaleOrderServiceInstance().ListGoods(orderID)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, map[string]interface{}{
		"orderInfo":  info,
		"orderGoods": goodsList,
	})
}

func SaleOrderExpressInfo(c *gin.Context) {
	orderID := bindSaleOrderIDFromQuery(c)
	info, err := service.SaleOrderServiceInstance().Info(orderID)
	utils.CheckAndPanic(err)
	expressInfo, err := express.ExpressService.ExpressInfo(express.ExpressType(info.ExpressMethod), info.ExpressNo)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, expressInfo)
}

func WechatPrePay(c *gin.Context) {
	orderID := bindSaleOrderIDFromQuery(c)
	userID := middle_ware.MustGetCurrentUser(c)
	result, err := service.SaleOrderServiceInstance().WechatPrepay(userID, orderID)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, result)
}

func WechatPayResult(c *gin.Context) {
	orderID := bindSaleOrderIDFromQuery(c)
	order, err := service.SaleOrderServiceInstance().Info(orderID)
	utils.CheckAndPanic(err)
	payResult, err := wechat.WechatService().QueryPayResult(order.OrderNo)
	utils.CheckAndPanic(err)
	err = service.SaleOrderServiceInstance().UpdateByWechatPayResult(orderID, payResult)
	utils.CheckAndPanic(err)
}

func bindSaleOrderIDFromQuery(c *gin.Context) int64 {
	orderIDStr := c.Query("orderId")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	utils.CheckAndPanic(err)
	return orderID
}
