package dao

import (
	"database/sql"
	"gotrue/facility/utils"
	"gotrue/model"
	"gotrue/sys"

	sq "github.com/Masterminds/squirrel"
)

var CutDetailDao *CutDetail

func initCutDetailDao() {
	m := &model.CutDetail{}
	CutDetailDao = &CutDetail{
		db:      sys.GetEasyDB(),
		table:   m.TableName(),
		columns: m.Columns(),
	}
}

type CutDetail struct {
	db      *sys.EasyDB
	table   string
	columns []string
}

func (d *CutDetail) CreateCutDetail(detail *model.CutDetail, tx *sql.Tx) (int64, error) {
	setMap := utils.StructToMap(detail, "db", "insert", "count")
	_, id, err := d.db.Insert(d.table, setMap, tx)
	return id, err
}

func (d *CutDetail) QueryByCutOrder(orderID int64) (*model.CutDetailSet, error) {
	details := []*model.CutDetail{}
	err := d.db.SelectDSL(&details, d.columns, d.table, sq.Eq{"cut_order_id": orderID, "is_delete": false})
	return model.NewCutDetailSet(details), err
}

func (d *CutDetail) QueryByCutOrderAndUser(orderID int64, userID int64) (*model.CutDetail, error) {
	detail := &model.CutDetail{}
	err := d.db.SelectOneDSL(detail, d.columns, d.table, sq.Eq{"cut_order_id": orderID, "user_id": userID, "is_delete": false})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return detail, err
}
