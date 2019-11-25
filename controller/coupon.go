package controller

import (
	"gotrue/facility/utils"
	"gotrue/middle_ware"
	"gotrue/service/api"
	"gotrue/service/coupon"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GrabCoupon(c *gin.Context) {
	couponService := coupon.GetCouponService()
	activityID := c.Query("activityId")
	if activityID == "" {
		middle_ware.BadRequest(c, "")
		return
	}
	activityIDInt64, err := strconv.ParseInt(activityID, 10, 64)
	if err != nil {
		middle_ware.BadRequest(c, "")
		return
	}
	err = couponService.GrabCoupon(c.Request.Context(), activityIDInt64)
	utils.CheckAndPanic(err)
}

func SendCoupon(c *gin.Context) {
	req := &api.SendCouponRequest{}
	err := c.BindJSON(req)
	if err != nil {
		middle_ware.BadRequest(c, "")
		return
	}
	err = coupon.GetCouponService().SendCouponToUser(c.Request.Context(), req)
	utils.CheckAndPanic(err)
}

func GetCouponListForUser(c *gin.Context) {
	couponService := coupon.GetCouponService()
	couponList, err := couponService.GetCouponListByUser(c.Request.Context())
	utils.CheckAndPanic(err)
	middle_ware.SetResponseData(c, gin.H{
		"coupons": couponList,
	})
}
