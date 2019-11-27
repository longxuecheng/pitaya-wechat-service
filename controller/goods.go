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

func GetInternalGoodsListByCategory(c *gin.Context) {
	categoryIDString := c.Query("categoryId")
	if categoryIDString == "" {
		middle_ware.BadRequest(c, "")
		return
	}
	categoryID, err := strconv.ParseInt(categoryIDString, 10, 64)
	if err != nil {
		middle_ware.BadRequest(c, "")
		return
	}
	goods, err := goods.GoodsService.GetGoodsByCategory(categoryID)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, gin.H{
		"goodsList": goods,
	})
}

// GetGoodsInfo 获取指定商品详情
// 包括（1.商品信息 2.图片列表 3.商品属性 4. 商品库存 5. 商品规格）
func GetGoodsInfo(c *gin.Context) {
	ID := c.Query("id")
	goodsID, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		middle_ware.BadRequest(c, "")
		return
	}
	goodsInfo, err := goods.GoodsService.GoodsInfo(goodsID)
	utils.CheckAndPanic(err)
	gallery, err := goods.GoodsImgService.GetByGoodsID(goodsID)
	utils.CheckAndPanic(err)
	attributes, err := goods.GoodsService.Attributes(goodsID)
	utils.CheckAndPanic(err)
	goodsSpecDTOs, err := goods.GoodsService.Specifications(goodsID)
	utils.CheckAndPanic(err)
	goodsSpecSet := response.NewGoodsSpecSet(goodsSpecDTOs)
	stockDTOs, err := stock.StockService.GetStocksByGoodsID(goodsID)
	utils.CheckAndPanic(err)
	specDTOs, err := basic.SpecificationService.GetByIDs(goodsSpecSet.DistinctSpecIDs())
	utils.CheckAndPanic(err)
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

func GetGoodsList(c *gin.Context) {
	categoryID := c.Query("categoryId")
	categoryInt64ID, err := strconv.ParseInt(categoryID, 10, 64)
	if err != nil {
		categoryInt64ID = 0
	}
	goodsList, err := goods.GoodsService.GoodsList(categoryInt64ID)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, gin.H{
		"goodsList": goodsList,
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
