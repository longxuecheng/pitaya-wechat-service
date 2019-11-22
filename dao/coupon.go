package dao

import (
	"database/sql"
	"gotrue/facility/utils"
	"gotrue/model"

	"github.com/Masterminds/squirrel"
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
	return coupons, d.db.SelectDSL(coupons, d.columns, d.table, squirrel.Eq{"user_id": userID, "consumed": false})
}

func (d *Coupon) CreateCoupon(coupon *model.Coupon, tx *sql.Tx) (int64, error) {
	setMap := utils.StructToMap(coupon, "db", "insert", "count")
	_, id, err := d.db.Insert(d.table, setMap, tx)
	return id, err
}
