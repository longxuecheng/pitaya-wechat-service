package dao

import (
	"database/sql"
	"gotrue/model"

	"github.com/Masterminds/squirrel"
)

var ActivityCouponDao *ActivityCoupon

func initActivityCouponDao() {
	m := &model.ActivityCoupon{}
	ActivityCouponDao = &ActivityCoupon{
		table:   m.TableName(),
		columns: m.Columns(),
		db:      GetEasyDB(),
	}
}

type ActivityCoupon struct {
	table   string
	columns []string
	db      *EasyDB
}

func (d *ActivityCoupon) QueryByActivityID(activityID int64) (*model.ActivityCoupon, error) {
	ac := &model.ActivityCoupon{}
	return ac, d.db.SelectOneDSL(ac, d.columns, d.table, squirrel.Eq{"activity_id": activityID})
}

func (d *ActivityCoupon) QueryByID(id int64) (*model.ActivityCoupon, error) {
	ac := &model.ActivityCoupon{}
	return ac, d.db.SelectOneDSL(ac, d.columns, d.table, squirrel.Eq{"id": id})
}

func (d *ActivityCoupon) DecreaseAvailQuantityByID(id int64, tx *sql.Tx) error {
	return d.db.ExecuteSQLTx(tx, "Update activity_coupon SET available_quantity = available_quantity - 1 WHERE id = ? AND available_quantity > 0", id)
}
