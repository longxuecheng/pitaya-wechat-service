package controller

import (
	"gotrue/dto/request"
	"gotrue/facility/utils"
	"gotrue/middle_ware"
	"gotrue/service"

	"github.com/gin-gonic/gin"
)

var (
	cashierServiceRef     = service.CashierServiceServiceInstance()
	userServiceRf4Cashier = service.UserServiceInstance()
)

func PreviewCashierFromCart(c *gin.Context) {
	userID := middle_ware.MustGetCurrentUser(c)
	cashier, err := cashierServiceRef.CartCheckout(userID)
	utils.CheckAndPanic(err)
	address, err := userServiceRf4Cashier.DefaultAddress(userID)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, gin.H{
		"cashier": cashier,
		"address": address,
	})
}

func PreviewCashierFromStock(c *gin.Context) {
	userID := middle_ware.MustGetCurrentUser(c)
	req := &request.CashierPreview{}
	err := c.BindJSON(req)
	utils.CheckAndPanic(err)
	cashier, err := cashierServiceRef.QuickCheckout(*req)
	utils.CheckAndPanic(err)
	address, err := userServiceRf4Cashier.DefaultAddress(userID)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, gin.H{
		"cashier": cashier,
		"address": address,
	})
}
