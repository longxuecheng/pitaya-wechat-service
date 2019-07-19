package controller

import (
	"gotrue/dto"
	"gotrue/facility/utils"
	"gotrue/middle_ware"
	"gotrue/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

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
	goods, err := service.GoodsServiceInstance().GetGoodsByCategory(categoryID)

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
	goodsInfo, err := service.GoodsServiceInstance().Info(goodsID)
	if err != nil {
		panic(err)
	}
	gallery, err := service.GoodsImgServiceSingleton.GetByGoodsID(goodsID)
	if err != nil {
		panic(err)
	}
	attributes, err := service.GoodsServiceInstance().Attributes(goodsID)
	if err != nil {
		panic(err)
	}
	goodsSpecDTOs, err := service.GoodsServiceInstance().Specifications(goodsID)
	if err != nil {
		panic(err)
	}
	goodsSpecSet := dto.NewGoodsSpecSet(goodsSpecDTOs)
	stockDTOs, err := service.StockServiceInstance().GetStocksByGoodsID(goodsID)
	if err != nil {
		panic(err)
	}
	specDTOs, err := service.SpecificationServiceSingleton.GetByIDs(goodsSpecSet.DistinctSpecIDs())
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
	hotGoods, err := service.GoodsServiceInstance().HotGoods()
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, map[string]interface{}{
		"hotGoods": hotGoods,
	})
}
