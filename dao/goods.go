package dao

import (
	"fmt"
	"gotrue/model"
	"gotrue/sys"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

// GoodsDao is a singleton of goods dao
var GoodsDao *Goods

func initGoodsDao() {
	goods := &model.Goods{}
	GoodsDao = &Goods{
		table:   goods.TableName(),
		columns: goods.Columns(),
		db:      sys.GetEasyDB(),
	}
}

// Goods is dao
type Goods struct {
	table   string
	columns []string
	db      *sys.EasyDB
}

func (dao *Goods) SelectByCategory(categoryID int64) ([]*model.Goods, error) {
	goods := []*model.Goods{}
	err := dao.db.Select(&goods, fmt.Sprintf("SELECT %s FROM goods WHERE category_id = ? ORDER BY id ASC", strings.Join(dao.columns, ",")), categoryID)
	if err != nil {
		return nil, err
	}
	return goods, nil
}

func (dao *Goods) SelectByID(ID int64) (*model.Goods, error) {
	goods := new(model.Goods)
	err := dao.db.SelectOne(goods, fmt.Sprintf("SELECT %s FROM goods WHERE id = ?", strings.Join(dao.columns, ",")), ID)
	if err != nil {
		return nil, err
	}
	return goods, nil
}

func (dao *Goods) SelectByIDs(IDs []int64) ([]*model.Goods, error) {
	goods := []*model.Goods{}
	err := dao.db.SelectDSL(&goods, dao.columns, dao.table, sq.Eq{"id": IDs})
	if err != nil {
		return nil, err
	}
	return goods, nil
}

func (dao *Goods) SelectAllByStatus(status model.GoodsStatus) ([]*model.Goods, error) {
	goods := []*model.Goods{}
	err := dao.db.SelectDSL(&goods, dao.columns, dao.table, sq.Eq{"status": string(status)})
	if err != nil {
		return nil, err
	}
	return goods, nil
}
