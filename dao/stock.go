package dao

import (
	"gotrue/model"
	"gotrue/sys"

	sq "github.com/Masterminds/squirrel"
)

// StockDao is a singleton of goods dao
var StockDao *Stock

func initStockDao() {
	stock := &model.Stock{}
	StockDao = &Stock{
		Table:   stock.TableName(),
		Columns: stock.Columns(),
		db:      sys.GetEasyDB(),
	}
}

// Stock is dao
type Stock struct {
	Table   string
	Columns []string
	db      *sys.EasyDB
}

func (dao *Stock) SelectByID(ID int64) (*model.Stock, error) {
	stock := new(model.Stock)
	err := dao.db.SelectOneDSL(stock, dao.Columns, dao.Table, sq.Eq{"id": ID})
	if err != nil {
		return nil, err
	}
	return stock, nil
}

func (dao *Stock) SelectByGoodsID(goodsID int64) ([]*model.Stock, error) {
	stocks := []*model.Stock{}
	err := dao.db.SelectDSL(&stocks, dao.Columns, dao.Table, sq.Eq{"goods_id": goodsID})
	if err != nil {
		return nil, err
	}
	return stocks, nil
}

func (dao *Stock) SelectByGoodsIDs(goodsIDs []int64) (*model.StockSet, error) {
	stocks := []*model.Stock{}
	err := dao.db.SelectDSL(&stocks, dao.Columns, dao.Table, sq.Eq{"goods_id": goodsIDs})
	if err != nil {
		return nil, err
	}
	return model.NewStockSet(stocks), nil
}

func (dao *Stock) SelectByIDs(ids []int64) ([]*model.Stock, error) {
	stocks := []*model.Stock{}
	err := dao.db.SelectDSL(&stocks, dao.Columns, dao.Table, sq.Eq{"id": ids})
	if err != nil {
		return nil, err
	}
	return stocks, nil
}
