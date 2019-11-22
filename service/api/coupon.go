package api

import (
	"context"
	"gotrue/model"

	"github.com/shopspring/decimal"
)

type ICouponService interface {
	CouponActivityInfo(ctx context.Context, activityID int64) (*CouponActivityResponse, error)

	GetCouponByCouponNo(ctx context.Context, couponNo string) (*model.Coupon, error)

	GetCouponListByUser(ctx context.Context) ([]*CouponResponse, error)

	GrabCoupon(ctx context.Context, activityID int64) error
}

type CouponResponse struct {
	ID           int64           `json:"-"`
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
