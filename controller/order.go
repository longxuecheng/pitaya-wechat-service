package controller

import (
	"gotrue/dto/request"
	"gotrue/facility/utils"
	"gotrue/middle_ware"
	"gotrue/service/api"
	"gotrue/service/express"
	"gotrue/service/order"
	"gotrue/service/wechat"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ReceiveOrder(c *gin.Context) {
	str := c.Query("orderId")
	if str == "" {
		middle_ware.BadRequest(c, "")
		return
	}
	orderID, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		middle_ware.BadRequest(c, "")
		return
	}
	data, err := order.SaleOrderService.ConfirmReceive(c.Request.Context(), orderID)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, gin.H{
		"order": data,
	})
}

func CancelOrder(c *gin.Context) {
	str := c.Query("orderId")
	if str == "" {
		middle_ware.BadRequest(c, "")
		return
	}
	orderID, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		middle_ware.BadRequest(c, "")
		return
	}
	data, err := order.SaleOrderService.Cancel(c.Request.Context(), orderID)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, gin.H{
		"order": data,
	})
}

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
	orderID, err := order.SaleOrderService.CreateFromStock(c.Request.Context(), req)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, orderID)
}

func SaleOrderList4User(c *gin.Context) {
	req := api.OrderListRequest{}
	err := c.BindJSON(&req)
	utils.CheckAndPanic(err)
	result, err := order.SaleOrderService.List(c.Request.Context(), req)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, result)
}

func SaleOrderList4SupplierAdmin(c *gin.Context) {
	req := api.OrderListRequest{}
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
	middle_ware.SetResponseData(c, gin.H{
		"orderInfo": info,
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

func Prepay(c *gin.Context) {
	orderID := bindSaleOrderIDFromQuery(c)
	userID := middle_ware.MustGetCurrentUser(c)
	result, err := order.SaleOrderService.Prepay(userID, orderID)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, result)
}

// GetPayResult is used for query wechat pay result from mini-program
func GetPayResult(c *gin.Context) {
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
	err = order.SaleOrderService.SendExpress(req)
	utils.CheckAndPanic(err)
}

func bindSaleOrderIDFromQuery(c *gin.Context) int64 {
	orderIDStr := c.Query("orderId")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	utils.CheckAndPanic(err)
	return orderID
}
