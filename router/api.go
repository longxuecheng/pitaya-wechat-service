package router

import (
	"pitaya-wechat-service/controller"

	"github.com/gin-gonic/gin"
)

func apiRouter(r *gin.Engine) {
	root := r.Group("/api")
	root.GET("/user/all", controller.GetUserListByConditions)
	root.GET("/category/tree", controller.GetCatogoriesTree)
	root.GET("/category/info", controller.GetCategoryInfo)
	root.GET("category/channel", controller.GetCategoryChannels)
	authGroup := root.Group("/auth")
	authGroup.POST("/login", controller.LoginByWechat)
	goodsGroup := root.Group("/goods")
	goodsGroup.GET("/list", controller.GetGoodsListByCategory)
	goodsGroup.GET("/detail", controller.GetGoodsInfo)
	goodsGroup.GET("/hot", controller.GetHotGoods)
	cartGroup := root.Group("/cart")
	cartGroup.POST("add", controller.AddCart)
	cartGroup.GET("index", controller.CartIndex)
	cartGroup.POST("checked", controller.CartItemCheck)
	cartGroup.GET("checkout", controller.CartCheckout)
	regionGroup := root.Group("/region")
	regionGroup.GET("list", controller.RegionList)
	addressGroup := root.Group("/address")
	addressGroup.GET("list", controller.UserAddressList)
	addressGroup.POST("save", controller.AddNewAddress)
	orderGroup := root.Group("/order")
	orderGroup.POST("submit", controller.SubmitSaleOrder)
	orderGroup.POST("list", controller.ListSaleOrders)
	orderGroup.GET("detail", controller.SaleOrderInfo)
}
