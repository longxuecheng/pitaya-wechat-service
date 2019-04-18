package router

import (
	"pitaya-wechat-service/controller"
	"pitaya-wechat-service/middle_ware"

	"github.com/gin-gonic/gin"
)

func apiRouter(r *gin.Engine) {

	root := r.Group("/api")
	categoryGroup := root.Group("/category")
	categoryGroup.GET("/tree", controller.GetCatogoriesTree)
	categoryGroup.GET("/info", controller.GetCategoryInfo)
	categoryGroup.GET("/channel", controller.GetCategoryChannels)
	categoryGroup.GET("/top/list", controller.GetTopCategories)
	authGroup := root.Group("/auth")
	authGroup.POST("/login", controller.LoginByWechat)
	goodsGroup := root.Group("/goods")
	goodsGroup.GET("/list", controller.GetGoodsListByCategory)
	goodsGroup.GET("/detail", controller.GetGoodsInfo)
	goodsGroup.GET("/hot", controller.GetHotGoods)
	cartGroup := root.Group("/cart", middle_ware.AuthCheck())
	cartGroup.POST("add", controller.AddCart)
	cartGroup.GET("index", controller.CartIndex)
	cartGroup.POST("checked", controller.CartItemCheck)
	cashierGroup := root.Group("/cashier")
	cashierGroup.GET("preview/cart", controller.PreviewCashierFromCart)
	cashierGroup.POST("preview/stock", controller.PreviewCashierFromStock)
	regionGroup := root.Group("/region")
	regionGroup.GET("list", controller.RegionList)
	addressGroup := root.Group("/address", middle_ware.AuthCheck())
	addressGroup.GET("list", controller.UserAddressList)
	addressGroup.POST("save", controller.AddNewAddress)
	orderGroup := root.Group("/order", middle_ware.AuthCheck())
	orderGroup.POST("submit", controller.SubmitSaleOrder)
	orderGroup.POST("list", controller.ListSaleOrders)
	orderGroup.GET("detail", controller.SaleOrderInfo)
}