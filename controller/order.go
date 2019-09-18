package controller

import (
	"gotrue/dto/request"
	"gotrue/facility/utils"
	"gotrue/middle_ware"
	"gotrue/service/express"
	"gotrue/service/order"
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
	orderID, err := order.SaleOrderService.CreateFromCart(userID, req)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, orderID)
}

func QuickSubmitOrder(c *gin.Context) {
	req := request.SaleOrderQuickAddRequest{}
	err := c.BindJSON(&req)
	utils.CheckAndPanic(err)
	userID := middle_ware.MustGetCurrentUser(c)
	orderID, err := order.SaleOrderService.CreateFromStock(userID, req)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, orderID)
}

// ListSaleOrders 获取用户的所有订单
func ListSaleOrders(c *gin.Context) {
	req := request.OrderListRequest{}
	err := c.BindJSON(&req)
	utils.CheckAndPanic(err)
	userID := middle_ware.MustGetCurrentUser(c)
	result, err := order.SaleOrderService.List(userID, req)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, result)
}

// ListSupplierOrders 给供应商管理员用户列出所有的订单
func ListSupplierOrders(c *gin.Context) {
	req := request.OrderListRequest{}
	err := c.BindJSON(&req)
	utils.CheckAndPanic(err)
	userID := middle_ware.MustGetCurrentUser(c)
	orderList, err := order.SaleOrderService.ListManagedOrders(userID, req)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, orderList)
}

// SaleOrderInfo 获取订单详情
func SaleOrderInfo(c *gin.Context) {
	orderID := bindSaleOrderIDFromQuery(c)
	info, err := order.SaleOrderService.Info(orderID)
	utils.CheckAndPanic(err)
	goodsList, err := order.SaleOrderService.ListGoods(orderID)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, map[string]interface{}{
		"orderInfo":  info,
		"orderGoods": goodsList,
	})
}

func SaleOrderExpressInfo(c *gin.Context) {
	orderID := bindSaleOrderIDFromQuery(c)
	info, err := order.SaleOrderService.Info(orderID)
	utils.CheckAndPanic(err)
	expressInfo, err := express.ExpressService.ExpressInfo(express.ExpressMethod(info.ExpressMethod), info.ExpressNo)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, expressInfo)
}

func WechatPrePay(c *gin.Context) {
	orderID := bindSaleOrderIDFromQuery(c)
	userID := middle_ware.MustGetCurrentUser(c)
	result, err := order.SaleOrderService.WechatPrepay(userID, orderID)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, result)
}

// WechatPayResult is used for query wechat pay result from mini-program
func WechatPayResult(c *gin.Context) {
	// orderID := bindSaleOrderIDFromQuery(c)
	req := &request.QueryWechatPayResult{}
	err := c.BindQuery(req)
	utils.CheckAndPanic(err)
	orderInfo, err := order.SaleOrderService.Info(req.OrderID)
	utils.CheckAndPanic(err)
	payResult, err := wechat.WechatService().QueryPayResult(orderInfo.OrderNo)
	utils.CheckAndPanic(err)
	err = order.SaleOrderService.PayResult(req, payResult)
	utils.CheckAndPanic(err)
}

// UpdateExpressInfo is used to update express information by tenant
func UpdateExpressInfo(c *gin.Context) {
	// userID := middle_ware.MustGetCurrentUser(c)
	req := &request.OrderExpressUpdate{}
	err := c.BindJSON(req)
	if err != nil {
		middle_ware.BadRequest(c, "请求参数不合法")
		return
	}
	utils.CheckAndPanic(err)
	err = order.SaleOrderService.UpdateExpressInfo(req)
	utils.CheckAndPanic(err)
}

func bindSaleOrderIDFromQuery(c *gin.Context) int64 {
	orderIDStr := c.Query("orderId")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	utils.CheckAndPanic(err)
	return orderID
}
