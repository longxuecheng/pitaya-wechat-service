package dao

import (
	"database/sql"
	"gotrue/facility/utils"
	"gotrue/model"
	"gotrue/sys"

	sq "github.com/Masterminds/squirrel"
)

// SaleDetailDaoSingleton is a singleton of goods dao
var SaleDetailDaoSingleton *SaleDetailDao

func SaleDetailDaoInstance() *SaleDetailDao {
	if SaleDetailDaoSingleton == nil {
		SaleDetailDaoSingleton = new(SaleDetailDao)
		SaleDetailDaoSingleton.db = sys.GetEasyDB()
	}
	return SaleDetailDaoSingleton
}

var columns_sale_detail_all = []string{"id", "order_id", "stock_id", "goods_id", "goods_name", "quantity", "sale_unit_price", "cost_unit_price", "goods_spec_ids", "goods_spec_description", "list_pic_url"}

// SaleDetailDao is dao
type SaleDetailDao struct {
	db *sys.EasyDB
}

func (dao *SaleDetailDao) ExecTx(execFunc func(tx *sql.Tx) error) {
	dao.db.ExecTx(execFunc)
}

func (dao *SaleDetailDao) Create(detail model.SaleDetail, tx ...*sql.Tx) (id int64, err error) {
	setMap := utils.StructToMap(detail, "db", "exclude")
	_, id, err = dao.db.Insert(model.Table_Sale_Detail, setMap, tx...)
	return
}

func (dao *SaleDetailDao) SelectByOrderIDs(orderIDs ...int64) ([]*model.SaleDetail, error) {
	details := []*model.SaleDetail{}
	err := dao.db.SelectDSL(&details, columns_sale_detail_all, model.Table_Sale_Detail, sq.Eq{"order_id": orderIDs})
	return details, err
}

func (dao *SaleDetailDao) SelectByOrderID(orderID int64) ([]*model.SaleDetail, error) {
	details := []*model.SaleDetail{}
	err := dao.db.SelectDSL(&details, columns_sale_detail_all, model.Table_Sale_Detail, sq.Eq{"order_id": orderID})
	return details, err
}
