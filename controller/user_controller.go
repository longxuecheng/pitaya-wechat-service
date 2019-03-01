package controller

import (
	"context"
	"net/http"
	"pitaya-wechat-service/api"
	"pitaya-wechat-service/dto/request"
	"pitaya-wechat-service/facility/utils"
	"pitaya-wechat-service/rpc"
	"pitaya-wechat-service/rpc/client"
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

func Login(c *gin.Context) {
	cc := client.AuthRPCConnection()
	authClient := rpc.NewAuthorizationServiceClient(cc)
	authRequest := &rpc.AuthRequest{
		Name:    "lxc",
		Code:    "code-dfdf",
		PhoneNo: "18911792314",
	}
	token, err := authClient.Authorize(context.Background(), authRequest)
	utils.CheckAndPanic(err)
	c.Set("data", token)
}
