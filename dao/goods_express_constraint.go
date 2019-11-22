package dao

import (
	"database/sql"
	"gotrue/facility/utils"
	"gotrue/model"
	

	"github.com/Masterminds/squirrel"
)

// GoodsExpressConstraintDao is a singleton of goods dao
var GoodsExpressConstraintDao *GoodsExpressConstraint

func InitGoodsExpressConstraintDao() {
	m := &model.GoodsExpressConstraint{}
	GoodsExpressConstraintDao = &GoodsExpressConstraint{
		db:      GetEasyDB(),
		table:   m.TableName(),
		columns: m.Columns(),
	}
}

// GoodsExpressConstraint is dao
type GoodsExpressConstraint struct {
	db      *EasyDB
	table   string
	columns []string
}

func (d *GoodsExpressConstraint) QueryByGoodsID(goodsID int64) ([]*model.GoodsExpressConstraint, error) {
	data := []*model.GoodsExpressConstraint{}
	err := d.db.SelectDSL(&data, d.columns, d.table, squirrel.Eq{"goods_id": goodsID})
	return data, err
}

func (d *GoodsExpressConstraint) QueryByStockAndProvince(stockID int64, provinceID int) (*model.GoodsExpressConstraint, error) {
	data := &model.GoodsExpressConstraint{}
	err := d.db.SelectOneDSL(&data, d.columns, d.table, squirrel.Eq{"stock_id": stockID, "province_id": provinceID})
	return data, err
}

func (d *GoodsExpressConstraint) CreateConstraint(c *model.GoodsExpressConstraint, tx *sql.Tx) (int64, error) {
	setMap := utils.InsertMap(c, "db")
	_, id, err := d.db.Insert(d.table, setMap, nil)
	return id, err
}
