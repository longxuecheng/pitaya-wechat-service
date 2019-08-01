package controller

import (
	"gotrue/dto/pagination"
	"gotrue/dto/request"
	"gotrue/facility/utils"
	"gotrue/middle_ware"
	"gotrue/service/express"
	"gotrue/service/order"
	"gotrue/service/supplier"
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
	orderID, err := order.SaleOrderService.Create(userID, req)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, orderID)
}

func QuickSubmitOrder(c *gin.Context) {
	req := request.SaleOrderQuickAddRequest{}
	err := c.BindJSON(&req)
	utils.CheckAndPanic(err)
	userID := middle_ware.MustGetCurrentUser(c)
	orderID, err := order.SaleOrderService.QuickCreate(userID, req)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, orderID)
}

// ListSaleOrders 获取用户的所有订单
func ListSaleOrders(c *gin.Context) {
	req := pagination.PaginationRequest{}
	err := c.BindJSON(&req)
	utils.CheckAndPanic(err)
	userID := middle_ware.MustGetCurrentUser(c)
	result, err := order.SaleOrderService.List(userID, req)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, result)
}

// ListSupplierOrders 给发货商列出订单
func ListSupplierOrders(c *gin.Context) {
	userID := middle_ware.MustGetCurrentUser(c)
	supplier, err := supplier.SupplierService.QueryByAdmin(userID)
	utils.CheckAndPanic(err)
	orderList, err := order.SaleOrderService.ListSupplierOrders(supplier.ID)
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

func WechatPayResult(c *gin.Context) {
	orderID := bindSaleOrderIDFromQuery(c)
	orderInfo, err := order.SaleOrderService.Info(orderID)
	utils.CheckAndPanic(err)
	payResult, err := wechat.WechatService().QueryPayResult(orderInfo.OrderNo)
	utils.CheckAndPanic(err)
	err = order.SaleOrderService.UpdateByWechatPayResult(orderID, payResult)
	utils.CheckAndPanic(err)
}

func UpdateExpressInfo(c *gin.Context) {
	userID := middle_ware.MustGetCurrentUser(c)
	req := &request.OrderExpressUpdate{}
	err := c.BindJSON(req)
	if err != nil {
		middle_ware.BadRequet(c, "请求参数不合法")
		return
	}
	supplier, err := supplier.SupplierService.QueryByAdmin(userID)
	utils.CheckAndPanic(err)
	req.SupplierID = supplier.ID
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
