package controller

import (
	"gotrue/facility/utils"
	"gotrue/middle_ware"
	"gotrue/service"

	"github.com/gin-gonic/gin"
)

func BannerList(c *gin.Context) {
	banners, err := service.BannerService.OnlineBanners()
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, gin.H{
		"banners": banners,
	})
}
