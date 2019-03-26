package controller

import (
	"fmt"
	"log"
	"net/http"
	"pitaya-wechat-service/api"
	"pitaya-wechat-service/dto/request"
	"pitaya-wechat-service/facility/utils"
	"pitaya-wechat-service/middle_ware"
	"pitaya-wechat-service/service"

	"github.com/gin-gonic/gin"
)

var userServiceRf api.IUserService = service.UserServiceInstance()

// GetUserListByConditions 获取用户列表
func GetUserListByConditions(c *gin.Context) {
	users, err := userServiceRf.GetList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
	}
	c.JSON(http.StatusOK, users)
}

func UserAddressList(c *gin.Context) {
	userID := middle_ware.MustGetCurrentUser(c)
	addresses, err := userServiceRf.AddressList(userID)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, addresses)
}

func AddNewAddress(c *gin.Context) {
	req := request.UserAddressAddRequest{}
	err := c.BindJSON(&req)
	utils.CheckAndPanic(err)
	userID := middle_ware.MustGetCurrentUser(c)
	id, err := userServiceRf.CreateAddress(userID, req)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, id)
}

// LoginByWechat 微信登录
func LoginByWechat(c *gin.Context) {
	req := request.WechatLogin{}
	utils.CheckAndPanic(c.BindJSON(&req))
	wechatResp, err := service.GetWechatUserInfo(req.Code)
	utils.CheckAndPanic(err)
	user, err := userServiceRf.Login(wechatResp.OpenID, req.NickName, req.AvatarURL)
	utils.CheckAndPanic(err)
	log.Println(fmt.Sprintf("LoginByWechat response code is %d jscode is %s nickname %s avatar url %s", wechatResp.ErrorCode, req.Code, req.NickName, req.AvatarURL))
	utils.CheckAndPanic(err)
	accessToken, err := service.BuildToken(user.ID, 3600)
	utils.CheckAndPanic(err)
	wechatUser := request.WechatUser{
		NickName:  user.NickName,
		AvatarURL: user.AvatarURL,
	}
	middle_ware.SetResponseData(c, map[string]interface{}{
		"token":    accessToken,
		"userInfo": wechatUser,
	})
}
