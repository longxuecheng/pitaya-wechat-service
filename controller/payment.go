package controller

import (
	"gotrue/facility/utils"
	"gotrue/middle_ware"
	"gotrue/service"
	"gotrue/wechat"

	"github.com/gin-gonic/gin"
)

func Pay(c *gin.Context) {
	userID := middle_ware.MustGetCurrentUser(c)
	user, err := service.UserServiceInstance().GetUserByID(userID)
	utils.CheckAndPanic(err)
	result, err := wechat.WechatService().Pay(user.OpenID)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, result)
}
