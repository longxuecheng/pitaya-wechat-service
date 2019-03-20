package controller

import (
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
	addresses, err := userServiceRf.AddressList(0)
	utils.CheckAndPanic(err)
	c.Set("data", addresses)
}

func AddNewAddress(c *gin.Context) {
	req := request.UserAddressAddRequest{}
	err := c.BindJSON(&req)
	utils.CheckAndPanic(err)
	req.UserID = 0
	id, err := userServiceRf.CreateAddress(req)
	utils.CheckAndPanic(err)
	c.Set("data", id)
}

func LoginByWechat(c *gin.Context) {
	code := request.WechatLogin{}
	utils.CheckAndPanic(c.BindJSON(&code))
	resp, err := service.GetWechatUserInfo(code.Code)
	log.Println(resp)
	utils.CheckAndPanic(err)
	accessToken, err := service.Authorize()
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, map[string]interface{}{
		"token":    accessToken,
		"userInfo": code.WechatUser,
	})
}
