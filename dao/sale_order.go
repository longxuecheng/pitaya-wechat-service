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
	SaleOrderDao = &SaleOrder{
		db: sys.GetEasyDB(),
	}
}

var columns_sale_order = []string{
	"id",
	"order_no",
	"create_time",
	"user_id",
	"status",
	"receiver",
	"province_id",
	"city_id",
	"district_id",
	"address",
	"phone_no",
	"supplier_id",
	"order_amt",
	"goods_amt",
	"express_method",
	"express_order_no",
	"express_fee",
}

// SaleOrder is dao
type SaleOrder struct {
	db *sys.EasyDB
}

func (dao *SaleOrder) Create(order model.SaleOrder, tx ...*sql.Tx) (id int64, err error) {
	setMap := utils.StructToMap(order, "db", "exclude")
	_, id, err = dao.db.Insert(model.Table_Sale_Order, setMap, tx...)
	return
}

func (dao *SaleOrder) SelectByUserID(userID int64) ([]model.SaleOrder, error) {
	orders := []model.SaleOrder{}
	err := dao.db.SelectDSL(&orders, columns_sale_order, model.Table_Sale_Order, sq.Eq{"user_id": userID})
	return orders, err
}

func (dao *SaleOrder) SelectByID(ID int64) (*model.SaleOrder, error) {
	order := &model.SaleOrder{}
	err := dao.db.SelectOneDSL(order, columns_sale_order, model.Table_Sale_Order, sq.Eq{"id": ID})
	return order, err
}

func (dao *SaleOrder) SelectByParentID(parentID int64) ([]model.SaleOrder, error) {
	orders := []model.SaleOrder{}
	err := dao.db.SelectDSL(&orders, columns_sale_order, model.Table_Sale_Order, sq.Eq{"parent_id": parentID})
	return orders, err
}

func (dao *SaleOrder) SelectByOrderNo(orderNo string) (*model.SaleOrder, error) {
	order := &model.SaleOrder{}
	err := dao.db.SelectOneDSL(order, columns_sale_order, model.Table_Sale_Order, sq.Eq{"order_no": orderNo})
	return order, err
}

func (dao *SaleOrder) SelectByUserIDWitPagination(userID int64, offset uint64, limit uint64) ([]model.SaleOrder, int64, error) {
	orderList := []model.SaleOrder{}
	totalRecords, err := dao.db.SelectPagination(&orderList, columns_sale_order, model.Table_Sale_Order, offset, limit, sq.Eq{"user_id": userID})
	return orderList, totalRecords, err
}

func (dao *SaleOrder) UpdateByID(orderID int64, updateMap map[string]interface{}, tx *sql.Tx) error {
	_, err := dao.db.UpdateTx(tx, model.Table_Sale_Order, updateMap, sq.Eq{"id": orderID})
	return err
}
