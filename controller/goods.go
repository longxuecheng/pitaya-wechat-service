package controller

import (
	"gotrue/dto/response"
	"gotrue/facility/utils"
	"gotrue/middle_ware"
	"gotrue/service/basic"
	"gotrue/service/goods"
	"gotrue/service/stock"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func GoodsExpressConstraint(c *gin.Context) {
	stockID := c.Query("stockId")
	addressId := c.Query("addressId")
	if stockID == "" || addressId == "" {
		middle_ware.BadRequest(c, "请求参数有误")
		return
	}

	stockIDInt64, err := utils.ParseInt64(stockID)
	utils.CheckAndPanic(err)
	addressIDInt64, err := utils.ParseInt64(addressId)
	utils.CheckAndPanic(err)
	constraint, err := goods.GoodsService.ExpressConstraint(stockIDInt64, addressIDInt64)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, constraint)
}

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
	goods, err := goods.GoodsService.GetGoodsByCategory(categoryID)

	if err != nil {
		panic(err)
	}
	c.Set("data", map[string]interface{}{"goodsList": goods})
}

// GetGoodsInfo 获取指定商品详情
// 包括（1.商品信息 2.图片列表 3.商品属性 4. 商品库存 5. 商品规格）
func GetGoodsInfo(c *gin.Context) {
	ID := c.Query("id")
	goodsID, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		panic(err)
	}
	goodsInfo, err := goods.GoodsService.Info(goodsID)
	if err != nil {
		panic(err)
	}
	gallery, err := goods.GoodsImgService.GetByGoodsID(goodsID)
	if err != nil {
		panic(err)
	}
	attributes, err := goods.GoodsService.Attributes(goodsID)
	if err != nil {
		panic(err)
	}
	goodsSpecDTOs, err := goods.GoodsService.Specifications(goodsID)
	if err != nil {
		panic(err)
	}
	goodsSpecSet := response.NewGoodsSpecSet(goodsSpecDTOs)
	stockDTOs, err := stock.StockService.GetStocksByGoodsID(goodsID)
	if err != nil {
		panic(err)
	}
	specDTOs, err := basic.SpecificationService.GetByIDs(goodsSpecSet.DistinctSpecIDs())
	if err != nil {
		panic(err)
	}
	goodsSpecSet.SetSpecs(specDTOs)
	specTree := goodsSpecSet.SpecTree()
	c.Set("data", map[string]interface{}{
		"info":              goodsInfo,
		"gallery":           gallery,
		"attributes":        attributes,
		"specificationList": specTree,
		"productList":       stockDTOs,
	})
}

// GetHotGoods 获取热门商品
func GetHotGoods(c *gin.Context) {
	hotGoods, err := goods.GoodsService.HotGoods()
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, map[string]interface{}{
		"hotGoods": hotGoods,
	})
}
