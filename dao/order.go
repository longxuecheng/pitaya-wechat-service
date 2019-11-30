package dao

import (
	"database/sql"
	"gotrue/facility/utils"
	"gotrue/model"

	"go.planetmeican.com/manage/paperwork-facility/reflect_util"

	sq "github.com/Masterminds/squirrel"
)

// SaleOrderDao is a singleton of goods dao
var SaleOrderDao *SaleOrder

func initSaleOrderDao() {
	m := &model.SaleOrder{}
	SaleOrderDao = &SaleOrder{
		db:      GetEasyDB(),
		table:   m.TableName(),
		columns: m.Columns(),
	}
}

// SaleOrder is dao
type SaleOrder struct {
	db      *EasyDB
	columns []string
	table   string
}

func (dao *SaleOrder) Create(order *model.SaleOrder, tx *sql.Tx) (id int64, err error) {
	setMap := utils.StructToMap(order, "db", "insert", "count")
	_, id, err = dao.db.Insert(dao.table, setMap, tx)
	return
}

func (dao *SaleOrder) SelectByUserID(userID int64) (model.SaleOrderList, error) {
	orders := model.SaleOrderList{}
	err := dao.db.SelectDSL(&orders, dao.columns, dao.table, sq.Eq{"user_id": userID})
	return orders, err
}

func (dao *SaleOrder) QueryByUserAndID(userID int64, orderID int64) (*model.SaleOrder, error) {
	order := &model.SaleOrder{}
	err := dao.db.SelectOneDSL(order, dao.columns, dao.table, sq.Eq{"id": orderID, "user_id": userID})
	return order, err
}

func (dao *SaleOrder) SelectByID(ID int64) (*model.SaleOrder, error) {
	order := &model.SaleOrder{}
	err := dao.db.SelectOneDSL(order, dao.columns, dao.table, sq.Eq{"id": ID})
	return order, err
}

func (dao *SaleOrder) SelectByParentID(parentID int64) (model.SaleOrderList, error) {
	orders := model.SaleOrderList{}
	err := dao.db.SelectDSL(&orders, dao.columns, dao.table, sq.Eq{"parent_id": parentID})
	return orders, err
}

func (dao *SaleOrder) SelectByOrderNo(orderNo string) (*model.SaleOrder, error) {
	order := &model.SaleOrder{}
	err := dao.db.SelectOneDSL(order, dao.columns, dao.table, sq.Eq{"order_no": orderNo})
	return order, err
}

func (dao *SaleOrder) SelectAllByUserWithPagination(userID int64, offset uint64, limit uint64) (model.SaleOrderList, int64, error) {
	orderList := model.SaleOrderList{}
	c := PaginationCondition{
		Columns:     dao.columns,
		TableName:   dao.table,
		Offset:      offset,
		Limit:       limit,
		WherePred:   sq.Eq{"user_id": userID},
		OrderbyPred: "id DESC",
	}
	totalRecords, err := dao.db.SelectPagination(&orderList, c)
	return orderList, totalRecords, err
}

func (dao *SaleOrder) SelectByUserAndStatus(userID int64, statusList []model.OrderStatus, offset uint64, limit uint64) (model.SaleOrderList, int64, error) {
	orderList := model.SaleOrderList{}
	c := PaginationCondition{
		Columns:     dao.columns,
		TableName:   dao.table,
		Offset:      offset,
		Limit:       limit,
		WherePred:   sq.Eq{"user_id": userID, "status": statusList},
		OrderbyPred: "id DESC",
	}
	totalRecords, err := dao.db.SelectPagination(&orderList, c)
	return orderList, totalRecords, err
}

// SelectBySupplierWithPagination query orders for single supplier
func (dao *SaleOrder) SelectBySupplierWithPagination(supplierID int64, offset uint64, limit uint64) (model.SaleOrderList, int64, error) {
	orderList := model.SaleOrderList{}
	c := PaginationCondition{
		Columns:     dao.columns,
		TableName:   dao.table,
		Offset:      offset,
		Limit:       limit,
		WherePred:   sq.Eq{"supplier_id": supplierID},
		OrderbyPred: "id DESC",
	}
	totalRecords, err := dao.db.SelectPagination(&orderList, c)
	return orderList, totalRecords, err
}

// SelectAllBySuppliersWithPagination query orders for multiple suppliers
func (dao *SaleOrder) SelectAllBySuppliersWithPagination(supplierIDs []int64, offset uint64, limit uint64) (model.SaleOrderList, int64, error) {
	orderList := model.SaleOrderList{}
	c := PaginationCondition{
		Columns:     dao.columns,
		TableName:   dao.table,
		Offset:      offset,
		Limit:       limit,
		WherePred:   sq.Eq{"supplier_id": supplierIDs},
		OrderbyPred: "id DESC",
	}
	totalRecords, err := dao.db.SelectPagination(&orderList, c)
	return orderList, totalRecords, err
}

// SelectBySupplierAndStatus query orders by status
func (dao *SaleOrder) SelectBySupplierAndStatus(supplierIDs []int64, stats []model.OrderStatus, offset uint64, limit uint64) (model.SaleOrderList, int64, error) {
	orderList := model.SaleOrderList{}
	c := PaginationCondition{
		Columns:     dao.columns,
		TableName:   dao.table,
		Offset:      offset,
		Limit:       limit,
		WherePred:   sq.Eq{"supplier_id": supplierIDs, "status": stats},
		OrderbyPred: "id DESC",
	}
	totalRecords, err := dao.db.SelectPagination(&orderList, c)
	return orderList, totalRecords, err
}

func (dao *SaleOrder) UpdateByID(order *model.SaleOrder, tx *sql.Tx) error {
	updateMap := reflect_util.StructToMap(order, "db", "pk", "count")
	_, err := dao.db.UpdateTx(tx, dao.table, updateMap, sq.Eq{"id": order.ID})
	return err
}

func (dao *SaleOrder) SelectBySupplier(supplierID int64) (model.SaleOrderList, error) {
	orders := model.SaleOrderList{}
	err := dao.db.SelectDSL(&orders, dao.columns, dao.table, sq.Eq{"supplier_id": supplierID})
	return orders, err
}

func (dao *SaleOrder) QueryUnSettledOrdersBySupplier(supplierID int64) (*model.SaleOrderSet, error) {
	orders := model.SaleOrderList{}
	err := dao.db.SelectDSL(&orders, dao.columns, dao.table, sq.Eq{"supplier_id": supplierID, "settlement_id": 0, "status": model.Paid})
	return &model.SaleOrderSet{
		Items: orders,
	}, err
}
