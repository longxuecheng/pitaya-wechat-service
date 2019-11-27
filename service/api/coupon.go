package api

import (
	"context"
	"gotrue/facility/errors"
	"gotrue/model"

	"github.com/shopspring/decimal"
)

type ICouponService interface {
	CouponActivityInfo(ctx context.Context, activityID int64) (*CouponActivityResponse, error)

	GetCouponByCouponNo(ctx context.Context, couponNo string) (*model.Coupon, error)

	GetAvailableCouponList(ctx context.Context) ([]*CouponResponse, error)

	GetReceivableCouponList(ctx context.Context) ([]*CouponResponse, error)

	GetExpiredCouponList(ctx context.Context) ([]*CouponResponse, error)

	GrabCoupon(ctx context.Context, activityID int64) error

	SendCouponToUser(ctx context.Context, req *SendCouponRequest) error
}

type CouponResponse struct {
	ID           int64           `json:"-"`
	Status       string          `json:"status"`
	CouponNo     string          `json:"couponNo"`
	ExpireTime   string          `json:"expireTime"`
	PriceString  string          `json:"price"`
	CouponType   string          `json:"couponType"`
	CategoryName string          `json:"categoryName"`
	CategoryID   int64           `json:"categoryId"`
	GoodsName    string          `json:"goodsName"`
	GoodsID      int64           `json:"goodsId"`
	Price        decimal.Decimal `json:"-"`
}

type CouponActivityResponse struct {
	Title           string `json:"title"`
	BannerURL       string `json:"bannerUrl"`
	BgURL           string `json:"bgUrl"`
	StartTime       string `json:"startTime"`
	ExpireTime      string `json:"expireTime"`
	CouponType      string `json:"couponType"`
	CouponPrice     string `json:"couponPrice"`
	TotalCoupon     int64  `json:"totalCoupon"`
	AvailableCoupon int64  `json:"availableCoupon"`
}

type SendCouponRequest struct {
	CategoryID           int64           `json:"categoryId"`
	GoodsID              int64           `json:"goodsId"`
	CouponType           string          `json:"couponType"`
	CouponPrice          decimal.Decimal `json:"couponPrice"`
	CouponQuantity       int             `json:"couponQuantity"`
	TargetUserID         int64           `json:"targetUserID"`
	ExpireTime           string          `json:"expireTime"`
	ComposableWithCutoff bool            `json:"composableWithCuffoff"`
}

func (r *SendCouponRequest) Validate() error {
	if r.TargetUserID == 0 {
		return errors.NewWithCodef("MissingTargetUser", "还没有指定用户")
	}
	if r.CouponPrice.Equals(decimal.Zero) {
		return errors.NewWithCodef("InvalidCouponPrice", "优惠券金额不能为0")
	}
	if r.CouponQuantity == 0 {
		return errors.NewWithCodef("InvalidCouponQuantity", "优惠券数量不能为0")
	}
	if r.ExpireTime == "" {
		return errors.NewWithCodef("MissingCouponExpireTime", "优惠券过期时间错误")
	}
	return nil
}
