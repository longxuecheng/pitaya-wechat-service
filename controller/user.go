package controller

import (
	"fmt"
	"gotrue/dto/request"
	"gotrue/dto/response"
	"gotrue/facility/utils"
	"gotrue/middle_ware"
	"gotrue/service"
	"gotrue/service/wechat"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserListByConditions 获取用户列表
func GetUserListByConditions(c *gin.Context) {
	users, err := service.UserServiceInstance().GetList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
	}
	c.JSON(http.StatusOK, users)
}

func UserAddressDelete(c *gin.Context) {
}

func UserAddressList(c *gin.Context) {
	userID := middle_ware.MustGetCurrentUser(c)
	addresses, err := service.UserServiceInstance().AddressList(userID)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, addresses)
}

func UserAddressDetail(c *gin.Context) {
	id := c.Query("id")
	addressID, err := utils.ParseInt64(id)
	utils.CheckAndPanic(err)
	address, err := service.UserServiceInstance().GetAddressByID(addressID)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, address)
}

func AddNewAddress(c *gin.Context) {
	req := request.UserAddressAddRequest{}
	err := c.BindJSON(&req)
	utils.CheckAndPanic(err)
	userID := middle_ware.MustGetCurrentUser(c)
	id, err := service.UserServiceInstance().CreateAddress(userID, req)
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, id)
}

// LoginByWechat 微信登录
func LoginByWechat(c *gin.Context) {
	req := request.WechatLogin{}
	utils.CheckAndPanic(c.BindJSON(&req))
	wechatResp, err := wechat.WechatService().UserInfo(req.Code)
	utils.CheckAndPanic(err)
	user, err := service.UserServiceInstance().Login(wechatResp.OpenID, req.NickName, req.AvatarURL)
	utils.CheckAndPanic(err)
	log.Println(fmt.Sprintf("LoginByWechat response code is %d jscode is %s nickname %s avatar url %s", wechatResp.ErrorCode, req.Code, req.NickName, req.AvatarURL))
	utils.CheckAndPanic(err)
	accessToken, err := service.BuildToken(user.ID, 3600)
	utils.CheckAndPanic(err)
	wechatUser := response.User{
		NickName:  user.NickName,
		AvatarURL: user.AvatarURL,
		UserType:  user.UserType.Int(),
	}

	middle_ware.SetResponseData(c, gin.H{
		"token":    accessToken,
		"userInfo": wechatUser,
	})
}
