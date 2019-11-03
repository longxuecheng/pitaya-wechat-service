package router

import (
	"gotrue/controller"
	"gotrue/middle_ware"

	"github.com/gin-gonic/gin"
)

func apiRouter(r *gin.Engine) {
	root := r.Group("/api")
	root.GET("/banner/list", controller.BannerList)
	categoryGroup := root.Group("/category")
	categoryGroup.GET("/tree", controller.GetCatogoriesTree)
	categoryGroup.GET("/info", controller.GetCategoryInfo)
	categoryGroup.GET("/channel", controller.GetCategoryChannels)
	categoryGroup.GET("/top/list", controller.GetTopCategories)

	authGroup := root.Group("/auth")
	authGroup.POST("/login", controller.LoginByWechat)

	goodsGroup := root.Group("/goods")
	goodsGroup.GET("/list", controller.GetGoodsListByCategory)
	goodsGroup.GET("/constraint/express", controller.GoodsExpressConstraint)
	goodsGroup.POST("/constraint/import", controller.ImportExpressConstraints)
	goodsGroup.GET("/detail", controller.GetGoodsInfo)
	goodsGroup.GET("/hot", controller.GetHotGoods)
	goodsGroup.GET("/card/list", controller.GetOnlineGoodsCards)

	cartGroup := root.Group("/cart", middle_ware.ValidateAuthorization)
	cartGroup.POST("add", controller.AddCart)
	cartGroup.GET("index", controller.CartIndex)
	cartGroup.POST("checked", controller.CartItemCheck)
	cartGroup.GET("goodscount", controller.CartGoodsCount)

	cashierGroup := root.Group("/cashier", middle_ware.ValidateAuthorization)
	cashierGroup.GET("preview/cart", controller.PreviewCashierFromCart)
	cashierGroup.POST("preview/stock", controller.PreviewCashierFromStock)

	userGroup := root.Group("/user", middle_ware.ValidateAuthorization)
	userGroup.POST("/channel/bind", controller.BindChannelUser)

	regionGroup := root.Group("/region")
	regionGroup.GET("list", controller.RegionList)
	regionGroup.GET("tips", controller.AddressTips)

	addressGroup := root.Group("/address", middle_ware.ValidateAuthorization)
	addressGroup.GET("list", controller.UserAddressList)
	addressGroup.POST("save", controller.AddNewAddress)
	addressGroup.GET("detail", controller.UserAddressDetail)

	{
		orderGroup := root.Group("/order", middle_ware.ValidateAuthorization)
		orderGroup.POST("submit", controller.SubmitSaleOrder)
		orderGroup.POST("submit/quick", controller.QuickSubmitOrder)
		orderGroup.POST("list", controller.ListSaleOrders)
		orderGroup.POST("list/supplier", controller.ListSupplierOrders)
		orderGroup.GET("detail", controller.SaleOrderInfo)
		orderGroup.GET("express", controller.SaleOrderExpressInfo)
		orderGroup.POST("express/update", controller.UpdateExpressInfo)
		orderGroup.GET("express/list", controller.ExpressList)
		orderGroup.GET("prepay", controller.WechatPrePay)
		orderGroup.GET("pay/result", controller.WechatPayResult)
	}
	{
		settleGroup := root.Group("/settlement", middle_ware.ValidateAuthorization)
		settleGroup.GET("supplier/show", controller.SettlementForSupplier)
	}
	{
		cutoffGroup := root.Group("/cutoff", middle_ware.ValidateAuthorization)
		cutoffGroup.POST("/start", controller.CreateCutOrder)
		cutoffGroup.GET("/info", controller.CutOrderInfo)
		cutoffGroup.POST("/assist", controller.CreateCutDetail)
	}
}
