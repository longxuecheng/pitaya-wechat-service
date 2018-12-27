package controller

import (
	"net/http"
	"pitaya-wechat-service/api"
	"pitaya-wechat-service/service"

	"github.com/gin-gonic/gin"
)

var userServiceReference api.UserService = service.UserServiceSingleton

// GetUserListByConditions 获取用户列表
func GetUserListByConditions(c *gin.Context) {
	users, err := userServiceReference.GetList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
	}
	c.JSON(http.StatusOK, users)
}
