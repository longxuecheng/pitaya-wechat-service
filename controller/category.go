package controller

import (
	"pitaya-wechat-service/api"
	"pitaya-wechat-service/facility/utils"
	"pitaya-wechat-service/middle_ware"
	"pitaya-wechat-service/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

var categoryServiceReference api.ICategoryService = service.CategoryServiceSingleton

// GetCatogoriesTree 获取顶级类目树列表
func GetCatogoriesTree(c *gin.Context) {
	categories, err := categoryServiceReference.GetCategoryTree()
	if err != nil {
		panic(err)
	}
	c.Set("data", map[string]interface{}{"categoryList": categories})
}

func GetCategoryInfo(c *gin.Context) {
	ID := c.Query("id")
	IDInt64, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		panic(err)
	}
	category, err := categoryServiceReference.GetCategoryByID(int(IDInt64))
	if err != nil {
		panic(err)
	}
	brothers, err := categoryServiceReference.GetCategoriesByParentID(category.ParentID)
	if err != nil {
		panic(err)
	}
	defer c.Set("data", map[string]interface{}{"currentCategory": category, "brotherCategory": brothers})
}

func GetCategoryChannels(c *gin.Context) {
	channels, err := categoryServiceReference.Channels()
	defer middle_ware.SetResponseData(c, channels)
	utils.CheckAndPanic(err)
}

func GetTopCategories(c *gin.Context) {
	categories, err := categoryServiceReference.GetTopList()
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, map[string]interface{}{
		"categories": categories,
	})
}
