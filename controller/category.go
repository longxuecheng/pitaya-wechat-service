package controller

import (
	"gotrue/facility/utils"
	"gotrue/middle_ware"
	"gotrue/service/basic"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCatogoriesTree 获取顶级类目树列表
func GetCatogoriesTree(c *gin.Context) {
	categories, err := basic.CategoryService.GetCategoryTree()
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
	category, err := basic.CategoryService.GetCategoryByID(int(IDInt64))
	if err != nil {
		panic(err)
	}
	brothers, err := basic.CategoryService.GetCategoriesByParentID(category.ParentID)
	if err != nil {
		panic(err)
	}
	defer c.Set("data", map[string]interface{}{"currentCategory": category, "brotherCategory": brothers})
}

func GetCategoryChannels(c *gin.Context) {
	channels, err := basic.CategoryService.Channels()
	defer middle_ware.SetResponseData(c, channels)
	utils.CheckAndPanic(err)
}

func GetTopCategories(c *gin.Context) {
	categories, err := basic.CategoryService.GetTopList()
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, map[string]interface{}{
		"categories": categories,
	})
}
