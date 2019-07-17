package controller

import (
	"gotrue/dto/request"
	"gotrue/dto/response"
	"gotrue/facility/utils"
	"gotrue/middle_ware"
	"gotrue/service"

	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

// AddCart 向购物车添加商品
func AddCart(c *gin.Context) {
	req := request.CartAddRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		panic(err)
	}
	userID := middle_ware.MustGetCurrentUser(c)
	req.UserID = userID
	_, err = service.CartServiceInstance().AddGoods(req)
	if err != nil {
		panic(err)
	}
	total, err := service.CartServiceInstance().GoodsCount(req.UserID)
	if err != nil {
		panic(err)
	}
	cartSummary := response.CartSummary{}
	cartTotal := response.CartTotal{
		GoodsCount: total,
	}
	cartSummary.CartTotal = cartTotal
	middle_ware.SetResponseData(c, cartSummary)
}

// CartIndex 获取某个用户下的购物车列表
func CartIndex(c *gin.Context) {
	userID := middle_ware.MustGetCurrentUser(c)
	carts, err := service.CartServiceInstance().List(userID)
	if err != nil {
		panic(err)
	}
	middle_ware.SetResponseData(c, summaryCart(carts))
}

// CartItemCheck 操作购物车条目的选择
func CartItemCheck(c *gin.Context) {
	req := request.CartCheckRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		panic(err)
	}
	err = service.CartServiceInstance().CheckItem(req)
	if err != nil {
		panic(err)
	}
	userID := middle_ware.MustGetCurrentUser(c)
	carts, err := service.CartServiceInstance().List(userID)
	middle_ware.SetResponseData(c, summaryCart(carts))
}

func CartGoodsCount(c *gin.Context) {
	userID := middle_ware.MustGetCurrentUser(c)
	count, err := service.CartServiceInstance().GoodsCount(userID)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, count)
}

func summaryCart(carts []response.CartItem) response.CartSummary {
	checkedGoodsCount := 0
	checkedGoodsAmount := decimal.Zero
	for _, cart := range carts {
		if cart.Checked == 1 {
			checkedGoodsCount++
			checkedGoodsAmount = checkedGoodsAmount.Add(cart.RetailPrice.Mul(cart.Quantity))
		}
	}
	cartSummary := response.CartSummary{}
	cartTotal := response.CartTotal{
		GoodsCount:         int64(len(carts)),
		CheckedGoodsCount:  checkedGoodsCount,
		CheckedGoodsAmount: checkedGoodsAmount.StringFixed(2),
	}
	cartSummary.CartList = carts
	cartSummary.CartTotal = cartTotal
	return cartSummary
}
