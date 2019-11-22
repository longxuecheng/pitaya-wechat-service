package dao

import (
	"fmt"
	"gotrue/model"

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
		db:      GetEasyDB(),
	}
}

// Goods is dao
type Goods struct {
	table   string
	columns []string
	db      *EasyDB
}

func (dao *Goods) SelectAllByStatus(status model.GoodsStatus) ([]*model.Goods, error) {
	goods := []*model.Goods{}
	err := dao.db.SelectDSL(&goods, dao.columns, dao.table, sq.Eq{"status": string(status)}, "id DESC")
	if err != nil {
		return nil, err
	}
	return goods, nil
}

func (dao *Goods) SelectByCategory(categoryID int64) ([]*model.Goods, error) {
	goods := []*model.Goods{}
	pred := sq.Eq{"category_id": categoryID, "status": model.GoodsStatusOnSale}
	err := dao.db.SelectDSL(&goods, dao.columns, dao.table, pred, "id DESC")
	if err != nil {
		return nil, err
	}
	return goods, nil
}

func (dao *Goods) SelectByID(ID int64) (*model.Goods, error) {
	goods := new(model.Goods)
	err := dao.db.SelectOne(goods, fmt.Sprintf("SELECT %s FROM goods WHERE id = ? AND status = ?", strings.Join(dao.columns, ",")), ID, model.GoodsStatusOnSale)
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

func (dao *Goods) QueryMapByIDs(idList []int64) (model.GoodsMap, error) {
	goods := model.GoodsList{}
	err := dao.db.SelectDSL(&goods, dao.columns, dao.table, sq.Eq{"id": idList})
	if err != nil {
		return nil, err
	}
	return goods.GoodsMap(), nil
}
