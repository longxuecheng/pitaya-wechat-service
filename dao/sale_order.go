package dao

import (
	"database/sql"
	"pitaya-wechat-service/model"
	"pitaya-wechat-service/sys"

	sq "github.com/Masterminds/squirrel"
)

// saleOrderDaoSingleton is a singleton of goods dao
var saleOrderDaoSingleton *SaleOrderDao

func SaleOrderDaoInstance() *SaleOrderDao {
	if saleOrderDaoSingleton == nil {
		saleOrderDaoSingleton = new(SaleOrderDao)
		saleOrderDaoSingleton.db = sys.GetEasyDB()
	}
	return saleOrderDaoSingleton
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
	"express_fee",
}

// SaleOrderDao is dao
type SaleOrderDao struct {
	db *sys.EasyDB
}

func (dao *SaleOrderDao) ExecTx(execFunc func(tx *sql.Tx) error) {
	dao.db.ExecTx(execFunc)
}

func (dao *SaleOrderDao) Create(setMap map[string]interface{}, tx ...*sql.Tx) (id int64, err error) {
	_, id, err = dao.db.Insert(model.Table_Sale_Order, setMap, tx...)
	return
}

func (dao *SaleOrderDao) SelectByUserID(userID int64) ([]model.SaleOrder, error) {
	orders := []model.SaleOrder{}
	err := dao.db.SelectDSL(&orders, columns_sale_order, model.Table_Sale_Order, sq.Eq{"user_id": userID})
	return orders, err
}

func (dao *SaleOrderDao) SelectByID(ID int64) (model.SaleOrder, error) {
	order := model.SaleOrder{}
	err := dao.db.SelectDSL(&order, columns_sale_order, model.Table_Sale_Order, sq.Eq{"id": ID})
	return order, err
}

func (dao *SaleOrderDao) SelectByUserIDWitPagination(userID int64, offset uint64, limit uint64) ([]model.SaleOrder, error) {
	orderList := []model.SaleOrder{}
	err := dao.db.SelectPagination(&orderList, columns_sale_order, model.Table_Sale_Order, offset, limit, sq.Eq{"user_id": userID})
	return orderList, err
}
