package coupon

import (
	"context"
	"database/sql"
	"gotrue/dao"
	"gotrue/facility/context_util"
	"gotrue/facility/errors"
	"gotrue/facility/utils"
	"gotrue/model"
	"gotrue/service/api"
)

var (
	ErrorActivityOffline      = errors.NewWithCodef("ActivityOfflin", "活动下线啦")
	ErrorActivityNotAvailable = errors.NewWithCodef("ActivityNotAvailable", "不在活动期间")
	ErrorCouponDrained        = errors.NewWithCodef("CouponDrained", "优惠券抢完啦")
	ErrorDuplicateGrabCoupon  = errors.NewWithCodef("DuplicateGrabCoupon", "不可以重复抢哦")
)

var couponServiceIns api.ICouponService

type Cut struct {
	cutOrderDao  *dao.CutOrder
	cutDetailDao *dao.CutDetail
	userDao      *dao.UserDao
	stockDao     *dao.Stock
	goodsSpecDao *dao.GoodsSpec
	goodsDao     *dao.Goods
}

func GetCouponService() api.ICouponService {
	if couponServiceIns == nil {
		couponServiceIns = NewCouponService()
	}
	return couponServiceIns
}

func NewCouponService() *CouponService {
	return &CouponService{
		couponDao:         dao.CouponDao,
		activityCouponDao: dao.ActivityCouponDao,
		activityDao:       dao.ActivityDao,
		categoryDao:       dao.CategoryDao,
		goodsDao:          dao.GoodsDao,
	}
}

type CouponService struct {
	couponDao         *dao.Coupon
	activityCouponDao *dao.ActivityCoupon
	activityDao       *dao.Activity
	categoryDao       *dao.Category
	goodsDao          *dao.Goods
}

func (s *CouponService) CouponActivityInfo(ctx context.Context, activityID int64) (*api.CouponActivityResponse, error) {
	activity, err := s.activityDao.QueryByID(activityID)
	if err == sql.ErrNoRows {
		return nil, ErrorActivityOffline
	}
	if err != nil {
		return nil, err
	}
	couponActivity, err := s.activityCouponDao.QueryByActivityID(activityID)
	if err != nil {
		return nil, err
	}
	apiCouponActivity := &api.CouponActivityResponse{
		Title:           activity.Title,
		BannerURL:       activity.BannerURL,
		BgURL:           activity.BgURL.String,
		StartTime:       utils.FormatTime(activity.StartTime, utils.TimePrecision_Seconds),
		ExpireTime:      utils.FormatTime(activity.ExpireTime, utils.TimePrecision_Seconds),
		CouponType:      couponActivity.CouponType.Title(),
		CouponPrice:     couponActivity.CouponPrice.StringFixed(2),
		TotalCoupon:     couponActivity.TotalQuantity,
		AvailableCoupon: couponActivity.AvailableQuantity,
	}
	return apiCouponActivity, nil
}

func (s *CouponService) GetCouponByCouponNo(ctx context.Context, couponNo string) (*model.Coupon, error) {
	coupon, err := s.couponDao.QueryByCouponNo(couponNo)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return coupon, nil
}

func (s *CouponService) GetCouponListByUser(ctx context.Context) ([]*api.CouponResponse, error) {
	userID, err := context_util.GetUserID(ctx)
	if err != nil {
		return nil, err
	}
	coupons, err := s.couponDao.QueryByUserID(userID)
	if err != nil {
		return nil, err
	}
	categoryMap, err := s.categoryDao.QueyMapByCategoryIDList(coupons.CategoryIDList())
	if err != nil {
		return nil, err
	}
	goodsMap, err := s.goodsDao.QueryMapByIDs(coupons.GoodsIDList())
	if err != nil {
		return nil, err
	}
	apiCoupons := make([]*api.CouponResponse, len(coupons))
	for i, coupon := range coupons {
		apiCoupons[i] = s.newAPICouponResponse(categoryMap, goodsMap, coupon)
	}
	return apiCoupons, nil
}

func (s *CouponService) newAPICouponResponse(categoryMap model.CategoryMap, goodsMap model.GoodsMap, coupon *model.Coupon) *api.CouponResponse {
	var categoryName string
	var goodsName string
	if cat := categoryMap.Get(coupon.CategoryID); cat != nil {
		categoryName = cat.Name
	}
	if goods := goodsMap.Get(coupon.GoodsID); goods != nil {
		goodsName = goods.Name
	}
	return &api.CouponResponse{
		ID:           coupon.ID,
		CategoryName: categoryName,
		CategoryID:   coupon.CategoryID,
		CouponNo:     coupon.CouponNo,
		CouponType:   coupon.Type.Title(),
		GoodsID:      coupon.GoodsID,
		GoodsName:    goodsName,
		Price:        coupon.Price,
		PriceString:  coupon.Price.StringFixed(2),
		ExpireTime:   utils.FormatTime(coupon.ExpireTime, utils.TimePrecision_Date),
	}
}

func (s *CouponService) GrabCoupon(ctx context.Context, activityID int64) error {
	// TODO 检查一个用户是否已经在该活动下抢到了优惠券，避免一个人抢多张
	userID, err := context_util.GetUserID(ctx)
	if err != nil {
		return err
	}
	_, err = s.couponDao.QueryByUserAndActivity(userID, activityID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == nil {
		return ErrorDuplicateGrabCoupon
	}
	activity, err := s.activityDao.QueryByID(activityID)
	if err == sql.ErrNoRows {
		return ErrorActivityOffline
	}
	if err != nil {
		return err
	}
	if !activity.IsAvailable() {
		return ErrorActivityNotAvailable
	}
	couponActivity, err := s.activityCouponDao.QueryByActivityID(activityID)
	if err != nil {
		return err
	}
	// double check coupon over status
	if couponActivity.IsDrained() {
		return ErrorCouponDrained
	}
	txExecutor, err := dao.GetEasyDB().NewTxExecutor()
	if err != nil {
		return err
	}
	grabCouponFunc := func(input interface{}, tx *sql.Tx) (interface{}, error) {
		return nil, s.activityCouponDao.DecreaseAvailQuantityByID(couponActivity.ID, tx)
	}
	createCouponFunc := func(input interface{}, tx *sql.Tx) (interface{}, error) {
		coupon := couponActivity.NewCoupon(userID)
		return s.couponDao.CreateCoupon(coupon, tx)
	}
	txExecutor.AppendFunc(grabCouponFunc)
	txExecutor.AppendFunc(createCouponFunc)
	err = txExecutor.Execute()
	if err == dao.ErrorNowRowAffected {
		return ErrorCouponDrained
	}
	// TODO return coupon info
	return err
}
