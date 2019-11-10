package controller

import (
	"gotrue/dto/request"
	"gotrue/facility/utils"
	"gotrue/middle_ware"
	"gotrue/service/cut"

	"github.com/gin-gonic/gin"
)

func CutoffInfo(c *gin.Context) {
	req := &request.CutoffInfo{}
	err := c.BindQuery(req)
	if err != nil {
		middle_ware.BadRequest(c, "")
		return
	}
	cutoff, err := cut.GetCutService().GetCutoffInfo(req.UserID, req.GoodsID)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, gin.H{
		"cutoffInfo": cutoff,
	})
}

func AssistCutoff(c *gin.Context) {
	userID := middle_ware.MustGetCurrentUser(c)
	req := &request.AssistCutoff{}
	err := c.BindJSON(req)
	utils.CheckAndPanic(err)
	req.HelperID = userID
	cutoff, err := cut.GetCutService().AssistCutoff(c.Request.Context(), req)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, gin.H{
		"cutoffInfo": cutoff,
	})
}
