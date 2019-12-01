package dao

import (
	"database/sql"
	"gotrue/facility/utils"
	"gotrue/model"
	"time"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	"go.planetmeican.com/manage/paperwork-facility/reflect_util"
)

var CouponDao *Coupon

func initCouponDao() {
	m := &model.Coupon{}
	CouponDao = &Coupon{
		table:   m.TableName(),
		columns: m.Columns(),
		db:      GetEasyDB(),
	}
}

type Coupon struct {
	table   string
	columns []string
	db      *EasyDB
}

func (d *Coupon) QueryByID(id int64) (*model.Coupon, error) {
	coupon := &model.Coupon{}
	return coupon, d.db.SelectOneDSL(coupon, d.columns, d.table, squirrel.Eq{"id": id})
}

func (d *Coupon) QueryByCouponNo(couponNo string) (*model.Coupon, error) {
	coupon := &model.Coupon{}
	return coupon, d.db.SelectOneDSL(coupon, d.columns, d.table, squirrel.Eq{"coupon_no": couponNo, "consumed": false})
}

func (d *Coupon) QueryByUserID(userID int64) (model.CouponList, error) {
	coupons := model.CouponList{}
	return coupons, d.db.SelectDSL(&coupons, d.columns, d.table, squirrel.Eq{"user_id": userID, "consumed": false})
}

func (d *Coupon) QueryByUserAndCouponID(userID, couponID int64) (*model.Coupon, error) {
	coupon := &model.Coupon{}
	return coupon, d.db.SelectOneDSL(coupon, d.columns, d.table, squirrel.Eq{"id": couponID, "user_id": userID})
}

func (d *Coupon) QueryAvailableCouponList(userID int64) (model.CouponList, error) {
	couponList := model.CouponList{}
	err := d.db.Select(&couponList, "SELECT * FROM coupon WHERE user_id = ? AND received IS true AND expire_time > ?", userID, time.Now())
	if err != nil {
		return nil, err
	}
	return couponList, nil
}

func (d *Coupon) QueryReceivableCouponList(userID int64) (model.CouponList, error) {
	couponList := model.CouponList{}
	err := d.db.Select(&couponList, "SELECT * FROM coupon WHERE user_id = ? AND received IS false AND expire_time > ?", userID, time.Now())
	if err != nil {
		return nil, err
	}
	return couponList, nil
}

func (d *Coupon) QueryExpiredCouponList(userID int64) (model.CouponList, error) {
	couponList := model.CouponList{}
	err := d.db.Select(&couponList, "SELECT * FROM coupon WHERE user_id = ? AND expire_time <= ?", userID, time.Now())
	if err != nil {
		return nil, err
	}
	return couponList, nil
}

func (d *Coupon) QueryByUserAndActivity(userID, activityID int64) (*model.Coupon, error) {
	coupon := &model.Coupon{}
	return coupon, d.db.SelectOneDSL(coupon, d.columns, d.table, squirrel.Eq{"user_id": userID, "activity_id": activityID})
}

func (d *Coupon) CreateCoupon(coupon *model.Coupon, tx *sql.Tx) (int64, error) {
	setMap := utils.StructToMap(coupon, "db", "insert", "count")
	_, id, err := d.db.Insert(d.table, setMap, tx)
	return id, err
}

func (d *Coupon) Update(coupon *model.Coupon, tx *sql.Tx) error {
	updateMap := reflect_util.StructToMap(coupon, "db", "pk", "count")
	_, err := d.db.UpdateTx(tx, d.table, updateMap, sq.Eq{"id": coupon.ID})
	return err
}
