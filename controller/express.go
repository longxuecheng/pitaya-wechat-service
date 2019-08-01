package controller

import (
	"gotrue/middle_ware"
	"gotrue/service/express"

	"github.com/gin-gonic/gin"
)

func ExpressList(c *gin.Context) {
	expressList := express.ExpressService.ExpressList()
	middle_ware.SetResponseData(c, expressList)
}
