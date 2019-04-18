package controller

import (
	"pitaya-wechat-service/api"
	"pitaya-wechat-service/facility/utils"
	"pitaya-wechat-service/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	regionServiceRf api.IRegionService = service.RegionServiceInstance()
)

func RegionList(c *gin.Context) {
	parentIDStr := c.Query("parentId")
	parentID, err := strconv.ParseInt(parentIDStr, 10, 64)
	utils.CheckAndPanic(err)
	regions, err := regionServiceRf.GetRegionsByParentID(int(parentID))
	utils.CheckAndPanic(err)
	c.Set("data", regions)
}
