package controller

import (
	"gotrue/facility/utils"
	"gotrue/service/region"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegionList(c *gin.Context) {
	parentIDStr := c.Query("parentId")
	parentID, err := strconv.ParseInt(parentIDStr, 10, 64)
	utils.CheckAndPanic(err)
	regions, err := region.RegionService.GetRegionsByParentID(int(parentID))
	utils.CheckAndPanic(err)
	c.Set("data", regions)
}
