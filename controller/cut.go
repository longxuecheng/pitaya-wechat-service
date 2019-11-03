package controller

import (
	"gotrue/dto/request"
	"gotrue/facility/utils"
	"gotrue/middle_ware"
	"gotrue/service/cut"

	"github.com/gin-gonic/gin"
)

func CutOrderInfo(c *gin.Context) {
	cutNo := c.Query("cutNo")
	cutOrder, err := cut.GetCutService().GetCutOrderByCutNo(cutNo)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, gin.H{
		"cutOrder": cutOrder,
	})
}

func CreateCutOrder(c *gin.Context) {
	req := &request.CutOrder{}
	err := c.BindJSON(req)
	utils.CheckAndPanic(err)
	userID := middle_ware.MustGetCurrentUser(c)
	req.UserID = userID
	cutOrder, err := cut.GetCutService().CreateCutOrder(c.Request.Context(), req)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, gin.H{
		"cutOrder": cutOrder,
	})
}

func CreateCutDetail(c *gin.Context) {
	userID := middle_ware.MustGetCurrentUser(c)
	req := &request.CreateCutDetail{}
	err := c.BindJSON(req)
	utils.CheckAndPanic(err)
	req.UserID = userID
	cutoff, err := cut.GetCutService().CreateCutDetail(req)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, gin.H{
		"cutoffPrice": cutoff,
	})
}
