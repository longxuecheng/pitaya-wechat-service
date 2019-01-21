package controller

import (
	"pitaya-wechat-service/api"
	"pitaya-wechat-service/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var goodsServiceReference api.IGoodsService = service.GoodsServiceSingleton

// GetGoodsListByCategory 按照商品分类获取商品列表
func GetGoodsListByCategory(c *gin.Context) {
	categoryIDString := c.Query("categoryId")
	if categoryIDString == "" {
		panic(errors.Errorf("param categoryId %s is invalid,need a not null number", categoryIDString))
	}
	categoryID, err := strconv.ParseInt(categoryIDString, 10, 64)
	if err != nil {
		panic(err)
	}
	goods, err := goodsServiceReference.GetGoodsByCategory(categoryID)
	if err != nil {
		panic(err)
	}
	c.Set("data", map[string]interface{}{"goodsList": goods})
}

func GetGoodsInfo(c *gin.Context) {
	ID := c.Query("id")
	IDInt64, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		panic(err)
	}
	category, err := categoryServiceReference.GetCategoryByID(int(IDInt64))
	if err != nil {
		panic(err)
	}
	defer c.Set("data", map[string]interface{}{"currentCategory": category})
}
