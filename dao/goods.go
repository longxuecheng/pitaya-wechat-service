package dao

import (
	"gotrue/model"

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

func (dao *Goods) QueyAll() (model.GoodsSet, error) {
	goods := model.GoodsSet{}
	pred := sq.Eq{"is_delete": false}
	err := dao.db.SelectDSL(&goods, dao.columns, dao.table, pred, "display_weight DESC", "click_count DESC", "id DESC")
	if err != nil {
		return nil, err
	}
	return goods, nil
}

func (dao *Goods) QueryAllByStatus(status string) (model.GoodsList, error) {
	goods := model.GoodsList{}
	pred := sq.Eq{"status": status, "is_delete": false}
	err := dao.db.SelectDSL(&goods, dao.columns, dao.table, pred, "display_weight DESC", "click_count DESC", "id DESC")
	if err != nil {
		return nil, err
	}
	return goods, nil
}

func (dao *Goods) QueryAllByCategory(categoryID int64) (model.GoodsSet, error) {
	goods := model.GoodsSet{}
	pred := sq.Eq{"category_id": categoryID, "is_delete": false}
	err := dao.db.SelectDSL(&goods, dao.columns, dao.table, pred, "id DESC")
	if err != nil {
		return nil, err
	}
	return goods, nil
}

func (dao *Goods) QueryByID(id int64) (*model.Goods, error) {
	goods := &model.Goods{}
	err := dao.db.SelectOneDSL(goods, dao.columns, dao.table, sq.Eq{"id": id, "is_delete": false})
	if err != nil {
		return nil, err
	}
	return goods, nil
}

func (dao *Goods) SelectOnSaleByID(ID int64) (*model.Goods, error) {
	goods := &model.Goods{}
	err := dao.db.SelectOneDSL(goods, dao.columns, dao.table, sq.Eq{"id": ID, "is_delete": false, "status": model.GoodsStatusOnSale})
	if err != nil {
		return nil, err
	}
	return goods, nil
}

func (dao *Goods) SelectByIDs(IDs []int64) (model.GoodsSet, error) {
	goods := model.GoodsSet{}
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

func (dao *Goods) IncreaseClickCount(goodsID int64) error {
	return dao.db.ExecuteSQLTx(nil, "UPDATE goods SET click_count = click_count + 1 WHERE id = ?", goodsID)
}
