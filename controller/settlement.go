package controller

import (
	"gotrue/facility/utils"
	"gotrue/middle_ware"
	"gotrue/service"
	"gotrue/service/supplier"

	"github.com/gin-gonic/gin"
)

func SettlementForSupplier(c *gin.Context) {
	userID := middle_ware.MustGetCurrentUser(c)
	sup, err := supplier.SupplierService.QueryByAdmin(userID)
	utils.CheckAndPanic(err)
	settlement, err := service.SettlementService.ShowSettlementForSupplier(sup.ID)
	middle_ware.SetResponseData(c, gin.H{
		"settlement": settlement,
	})
}
