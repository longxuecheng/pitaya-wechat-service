package dao

import (
	"database/sql"
	"gotrue/facility/utils"
	"gotrue/model"
	"gotrue/sys"

	sq "github.com/Masterminds/squirrel"
)

// SaleDetailDao is a singleton of goods dao
var SaleDetailDao *SaleDetail

func initSaleDetailDao() {
	SaleDetailDao = &SaleDetail{
		db: sys.GetEasyDB(),
	}
}

var columns_sale_detail_all = []string{"id", "order_id", "stock_id", "goods_id", "goods_name", "quantity", "sale_unit_price", "cost_unit_price", "goods_spec_ids", "goods_spec_description", "list_pic_url"}

// SaleDetail is dao
type SaleDetail struct {
	db *sys.EasyDB
}

func (dao *SaleDetail) ExecTx(execFunc func(tx *sql.Tx) error) {
	dao.db.ExecTx(execFunc)
}

func (dao *SaleDetail) Create(detail *model.SaleDetail, tx ...*sql.Tx) (id int64, err error) {
	setMap := utils.StructToMap(detail, "db", "exclude")
	_, id, err = dao.db.Insert(model.Table_Sale_Detail, setMap, tx...)
	return
}

func (dao *SaleDetail) SelectByOrderIDs(orderIDs ...int64) ([]*model.SaleDetail, error) {
	details := []*model.SaleDetail{}
	err := dao.db.SelectDSL(&details, columns_sale_detail_all, model.Table_Sale_Detail, sq.Eq{"order_id": orderIDs})
	return details, err
}

func (dao *SaleDetail) SelectByOrderID(orderID int64) ([]*model.SaleDetail, error) {
	details := []*model.SaleDetail{}
	err := dao.db.SelectDSL(&details, columns_sale_detail_all, model.Table_Sale_Detail, sq.Eq{"order_id": orderID})
	return details, err
}
