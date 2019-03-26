package controller

import (
	"pitaya-wechat-service/api"
	"pitaya-wechat-service/dto"
	"pitaya-wechat-service/dto/request"
	"pitaya-wechat-service/dto/response"
	"pitaya-wechat-service/facility/utils"
	"pitaya-wechat-service/middle_ware"
	"pitaya-wechat-service/service"

	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

var (
	cartServiceRf      api.ICartService = service.CartServiceInstance()
	cartRf_UserService api.IUserService = service.UserServiceInstance()
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
	_, err = cartServiceRf.AddGoods(req)
	if err != nil {
		panic(err)
	}
	total, err := cartServiceRf.GoodsCount(req.UserID)
	if err != nil {
		panic(err)
	}
	cartSummary := response.CartSummaryDTO{}
	cartTotal := response.CartTotalDTO{
		GoodsCount: total,
	}
	cartSummary.CartTotal = cartTotal
	c.Set("data", cartSummary)
}

// CartIndex 获取某个用户下的购物车列表
func CartIndex(c *gin.Context) {
	userID := middle_ware.MustGetCurrentUser(c)
	carts, err := cartServiceRf.ListCart4User(userID)
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
	err = cartServiceRf.CheckItem(req)
	if err != nil {
		panic(err)
	}
	userID := middle_ware.MustGetCurrentUser(c)
	carts, err := cartServiceRf.ListCart4User(userID)
	middle_ware.SetResponseData(c, summaryCart(carts))
}

// CartCheckout 结算台信息
func CartCheckout(c *gin.Context) {
	userID := middle_ware.MustGetCurrentUser(c)
	carts, err := cartServiceRf.ListCart4User(userID)
	utils.CheckAndPanic(err)
	cartsFiltered := filterCartItem(carts, func(cart response.CartItemDTO) bool {
		return cart.Checked == 1
	})
	cartSum := summaryCart(cartsFiltered)
	expressFee := decimal.NewFromFloat32(20.35)
	goodsTotalPrice, _ := decimal.NewFromString(cartSum.CartTotal.CheckedGoodsAmount)
	orderTotalPrice := goodsTotalPrice.Add(expressFee)
	addressList, err := cartRf_UserService.AddressList(userID)
	var checkedAddress = dto.UserAddressDTO{}
	for _, address := range addressList {
		if address.IsDefault {
			checkedAddress = address
			break
		}
	}
	utils.CheckAndPanic(err)

	resultmap := map[string]interface{}{
		"checkedGoodsList": cartsFiltered,
		"expressFee":       expressFee.StringFixed(2),
		"goodsTotalPrice":  goodsTotalPrice.StringFixed(2),
		"orderTotalPrice":  orderTotalPrice.StringFixed(2),
		"actualPrice":      orderTotalPrice.StringFixed(2),
		"checkedAddress":   checkedAddress,
	}
	c.Set("data", resultmap)
}

func filterCartItem(carts []response.CartItemDTO, filterFunc func(input response.CartItemDTO) bool) []response.CartItemDTO {
	results := []response.CartItemDTO{}
	for _, cart := range carts {
		if filterFunc(cart) {
			results = append(results, cart)
		}
	}
	return results
}

func summaryCart(carts []response.CartItemDTO) response.CartSummaryDTO {
	checkedGoodsCount := 0
	checkedGoodsAmount := decimal.Zero
	for _, cart := range carts {
		if cart.Checked == 1 {
			checkedGoodsCount++
			checkedGoodsAmount = checkedGoodsAmount.Add(cart.RetailPrice.Mul(cart.Quantity))
		}
	}
	cartSummary := response.CartSummaryDTO{}
	cartTotal := response.CartTotalDTO{
		GoodsCount:         int64(len(carts)),
		CheckedGoodsCount:  checkedGoodsCount,
		CheckedGoodsAmount: checkedGoodsAmount.StringFixed(2),
	}
	cartSummary.CartList = carts
	cartSummary.CartTotal = cartTotal
	return cartSummary
}
