package dao

import (
	"database/sql"
	"gotrue/facility/utils"
	"gotrue/model"
	"gotrue/sys"

	"go.planetmeican.com/manage/paperwork-facility/reflect_util"

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
	setMap := utils.StructToMap(order, "db", "insert", "count")
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

func (dao *SaleOrder) SelectAllByUserWithPagination(userID int64, offset uint64, limit uint64) ([]model.SaleOrder, int64, error) {
	orderList := []model.SaleOrder{}
	c := sys.PaginationCondition{
		Columns:   dao.columns,
		TableName: dao.table,
		Offset:    offset,
		Limit:     limit,
		WherePred: sq.Eq{"user_id": userID},
	}
	totalRecords, err := dao.db.SelectPagination(&orderList, c)
	return orderList, totalRecords, err
}

func (dao *SaleOrder) SelectByUserAndStatus(userID int64, statusList []model.OrderStatus, offset uint64, limit uint64) ([]model.SaleOrder, int64, error) {
	orderList := []model.SaleOrder{}
	c := sys.PaginationCondition{
		Columns:   dao.columns,
		TableName: dao.table,
		Offset:    offset,
		Limit:     limit,
		WherePred: sq.Eq{"user_id": userID, "status": statusList},
	}
	totalRecords, err := dao.db.SelectPagination(&orderList, c)
	return orderList, totalRecords, err
}

func (dao *SaleOrder) SelectAllBySupplierWithPagination(supplierID int64, offset uint64, limit uint64) ([]model.SaleOrder, int64, error) {
	orderList := []model.SaleOrder{}
	c := sys.PaginationCondition{
		Columns:   dao.columns,
		TableName: dao.table,
		Offset:    offset,
		Limit:     limit,
		WherePred: sq.Eq{"supplier_id": supplierID},
	}
	totalRecords, err := dao.db.SelectPagination(&orderList, c)
	return orderList, totalRecords, err
}

func (dao *SaleOrder) SelectBySupplierAndStatus(supplierID int64, stats []model.OrderStatus, offset uint64, limit uint64) ([]model.SaleOrder, int64, error) {
	orderList := []model.SaleOrder{}
	c := sys.PaginationCondition{
		Columns:   dao.columns,
		TableName: dao.table,
		Offset:    offset,
		Limit:     limit,
		WherePred: sq.Eq{"supplier_id": supplierID, "status": stats},
	}
	totalRecords, err := dao.db.SelectPagination(&orderList, c)
	return orderList, totalRecords, err
}

func (dao *SaleOrder) UpdateByID(order *model.SaleOrder, tx *sql.Tx) error {
	updateMap := reflect_util.StructToMap(order, "db", "pk", "count")
	_, err := dao.db.UpdateTx(tx, dao.table, updateMap, sq.Eq{"id": order.ID})
	return err
}

func (dao *SaleOrder) SelectBySupplier(supplierID int64) ([]model.SaleOrder, error) {
	orders := []model.SaleOrder{}
	err := dao.db.SelectDSL(&orders, dao.columns, dao.table, sq.Eq{"supplier_id": supplierID})
	return orders, err
}
