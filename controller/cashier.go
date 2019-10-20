package controller

import (
	"gotrue/dto/request"
	"gotrue/dto/response"
	"gotrue/facility/utils"
	"gotrue/middle_ware"
	"gotrue/service/cashier"
	"gotrue/service/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PreviewCashierFromCart 从购物车进入结算台预览
func PreviewCashierFromCart(c *gin.Context) {
	userID := middle_ware.MustGetCurrentUser(c)
	addressIDString := c.Query("addressId")
	if addressIDString == "" {
		middle_ware.SetResponseDataWithStatus(c, nil, http.StatusBadRequest)
		return
	}
	addressID, err := utils.ParseInt64(addressIDString)
	utils.CheckAndPanic(err)
	cashier, err := cashier.CashierService.CartCheckout(userID)
	utils.CheckAndPanic(err)
	address, err := getUserAddress(addressID, userID)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, gin.H{
		"cashier": cashier,
		"address": address,
	})
}

// PreviewCashierFromStock 从商品直接进入结算台
func PreviewCashierFromStock(c *gin.Context) {
	userID := middle_ware.MustGetCurrentUser(c)
	req := &request.CashierPreview{}
	err := c.BindJSON(req)
	utils.CheckAndPanic(err)
	address, err := getUserAddress(req.AddressID, userID)
	utils.CheckAndPanic(err)
	req.AddressID = address.ID
	req.UserID = userID
	cashier, err := cashier.CashierService.StockCheckout(req)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, gin.H{
		"cashier": cashier,
		"address": address,
	})
}

func getUserAddress(addressID, userID int64) (*response.UserAddress, error) {
	if addressID == 0 {
		return user.UserService.DefaultAddress(userID)
	}
	return user.UserService.GetAddressByID(addressID)
}
