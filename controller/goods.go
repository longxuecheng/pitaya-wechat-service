package controller

import (
	"fmt"
	"gotrue/dto/response"
	"gotrue/facility/errors"
	"gotrue/facility/utils"
	"gotrue/middle_ware"
	"gotrue/service/basic"
	"gotrue/service/goods"
	"gotrue/service/stock"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GoodsExpressConstraint(c *gin.Context) {
	stockID := c.Query("stockId")
	addressID := c.Query("addressId")
	if stockID == "" || addressID == "" {
		middle_ware.BadRequest(c, "请求参数有误")
		return
	}
	stockIDInt64, err := utils.ParseInt64(stockID)
	utils.CheckAndPanic(err)
	addressIDInt64, err := utils.ParseInt64(addressID)
	utils.CheckAndPanic(err)
	constraint, err := goods.GoodsService.ExpressConstraint(stockIDInt64, addressIDInt64)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, constraint)
}

func ImportExpressConstraints(c *gin.Context) {
	stockIDString := c.Query("stockId")
	stockID, err := utils.ParseInt64(stockIDString)
	utils.CheckAndPanic(err)
	mf, err := c.MultipartForm()
	if err == http.ErrNotMultipart {
		middle_ware.BadRequest(c, "文件上传方式有误")
		return
	}
	if len(mf.File) == 0 {
		panic(errors.NewWithCodef("FileNotFound", "文件不能为空"))
	}
	if mf.File["constraint"] == nil {
		middle_ware.BadRequest(c, "file must be put to key:(constraint)")
		return
	}
	fh := mf.File["constraint"][0]
	f, err := fh.Open()
	defer f.Close()
	utils.CheckAndPanic(err)
	err = goods.GoodsService.ImportExpressConstraints(stockID, f)
	utils.CheckAndPanic(err)
}

// GetGoodsListByCategory 按照商品分类获取商品列表
func GetGoodsListByCategory(c *gin.Context) {
	categoryIDString := c.Query("categoryId")
	if categoryIDString == "" {
		middle_ware.BadRequest(c, fmt.Sprintf("param categoryId %s is invalid,need a not null number", categoryIDString))
		return
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
	categoryID := c.Query("categoryId")
	categoryInt64ID, err := strconv.ParseInt(categoryID, 10, 64)
	if err != nil {
		categoryInt64ID = 0
	}
	hotGoods, err := goods.GoodsService.HotGoods(categoryInt64ID)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, map[string]interface{}{
		"hotGoods": hotGoods,
	})
}

// GetOnlineGoodsCards ...
func GetOnlineGoodsCards(c *gin.Context) {
	cards, err := goods.GoodsService.OneSaleGoodsCards()
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, gin.H{
		"cards": cards,
	})
}
