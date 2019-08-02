package dao

import (
	"database/sql"
	"gotrue/facility/utils"
	"gotrue/model"
	"gotrue/sys"

	sq "github.com/Masterminds/squirrel"
)

// SaleOrderDao is a singleton of goods dao
var SaleOrderDao *SaleOrder

func initSaleOrderDao() {
	m := &model.SaleOrder{}
	SaleOrderDao = &SaleOrder{
		db:      sys.GetEasyDB(),
		table:   m.TableName(),
		columns: m.Columns(),
	}
}

// SaleOrder is dao
type SaleOrder struct {
	db      *sys.EasyDB
	columns []string
	table   string
}

func (dao *SaleOrder) Create(order model.SaleOrder, tx ...*sql.Tx) (id int64, err error) {
	setMap := utils.StructToMap(order, "db", "exclude")
	_, id, err = dao.db.Insert(dao.table, setMap, tx...)
	return
}

func (dao *SaleOrder) SelectByUserID(userID int64) ([]model.SaleOrder, error) {
	orders := []model.SaleOrder{}
	err := dao.db.SelectDSL(&orders, dao.columns, dao.table, sq.Eq{"user_id": userID})
	return orders, err
}

func (dao *SaleOrder) SelectByID(ID int64) (*model.SaleOrder, error) {
	order := &model.SaleOrder{}
	err := dao.db.SelectOneDSL(order, dao.columns, dao.table, sq.Eq{"id": ID})
	return order, err
}

func (dao *SaleOrder) SelectByParentID(parentID int64) ([]model.SaleOrder, error) {
	orders := []model.SaleOrder{}
	err := dao.db.SelectDSL(&orders, dao.columns, dao.table, sq.Eq{"parent_id": parentID})
	return orders, err
}

func (dao *SaleOrder) SelectByOrderNo(orderNo string) (*model.SaleOrder, error) {
	order := &model.SaleOrder{}
	err := dao.db.SelectOneDSL(order, dao.columns, dao.table, sq.Eq{"order_no": orderNo})
	return order, err
}

func (dao *SaleOrder) SelectByUserIDWitPagination(userID int64, offset uint64, limit uint64) ([]model.SaleOrder, int64, error) {
	orderList := []model.SaleOrder{}
	totalRecords, err := dao.db.SelectPagination(&orderList, dao.columns, dao.table, offset, limit, sq.Eq{"user_id": userID})
	return orderList, totalRecords, err
}

func (dao *SaleOrder) SelectBySupplierWitPagination(supplierID int64, offset uint64, limit uint64) ([]model.SaleOrder, int64, error) {
	orderList := []model.SaleOrder{}
	totalRecords, err := dao.db.SelectPagination(&orderList, dao.columns, dao.table, offset, limit, sq.Eq{"supplier_id": supplierID})
	return orderList, totalRecords, err
}

func (dao *SaleOrder) UpdateByID(orderID int64, updateMap map[string]interface{}, tx *sql.Tx) error {
	_, err := dao.db.UpdateTx(tx, dao.table, updateMap, sq.Eq{"id": orderID})
	return err
}

func (dao *SaleOrder) SelectBySupplier(supplierID int64) ([]model.SaleOrder, error) {
	orders := []model.SaleOrder{}
	err := dao.db.SelectDSL(&orders, dao.columns, dao.table, sq.Eq{"supplier_id": supplierID})
	return orders, err
}
