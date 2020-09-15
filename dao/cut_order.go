package dao

import (
	"database/sql"
	"gotrue/facility/utils"
	"gotrue/model"
	

	sq "github.com/Masterminds/squirrel"
	"/manage/paperwork-facility/reflect_util"
)

var CutOrderDao *CutOrder

func initCutOrderDao() {
	m := &model.CutOrder{}
	CutOrderDao = &CutOrder{
		db:      GetEasyDB(),
		table:   m.TableName(),
		columns: m.Columns(),
	}
}

type CutOrder struct {
	db      *EasyDB
	table   string
	columns []string
}

func (d *CutOrder) CreateCutOrder(order *model.CutOrder, tx *sql.Tx) (int64, error) {
	setMap := utils.StructToMap(order, "db", "insert", "count")
	_, id, err := d.db.Insert(d.table, setMap, tx)
	return id, err
}

func (d *CutOrder) UpdateCutOrder(order *model.CutOrder, tx *sql.Tx) error {
	updateMap := reflect_util.StructToMap(order, "db", "pk", "count")
	_, err := d.db.UpdateTx(tx, d.table, updateMap, sq.Eq{"id": order.ID, "consumed": false})
	return err
}

func (d *CutOrder) QueryByCutNo(cutNo string) (*model.CutOrder, error) {
	order := &model.CutOrder{}
	err := d.db.SelectOneDSL(order, d.columns, d.table, sq.Eq{"cut_no": cutNo, "consumed": false})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return order, err
}

func (d *CutOrder) QueryByID(id int64) (*model.CutOrder, error) {
	order := &model.CutOrder{}
	err := d.db.SelectOneDSL(order, d.columns, d.table, sq.Eq{"id": id, "consumed": false})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return order, err
}

func (d *CutOrder) QueryByUserAndStock(userID int64, stockID int64) (*model.CutOrder, error) {
	order := &model.CutOrder{}
	err := d.db.SelectOneDSL(order, d.columns, d.table, sq.Eq{"user_id": userID, "stock_id": stockID, "consumed": false})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return order, err
}

func (d *CutOrder) QueryByUserAndGoods(userID int64, goodsID int64) (*model.CutOrder, error) {
	order := &model.CutOrder{}
	err := d.db.SelectOneDSL(order, d.columns, d.table, sq.Eq{"user_id": userID, "goods_id": goodsID, "consumed": false})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return order, err
}
