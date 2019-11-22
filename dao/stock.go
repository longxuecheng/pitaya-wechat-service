package dao

import (
	"gotrue/model"
	

	sq "github.com/Masterminds/squirrel"
)

// StockDao is a singleton of goods dao
var StockDao *Stock

func initStockDao() {
	stock := &model.Stock{}
	StockDao = &Stock{
		Table:   stock.TableName(),
		Columns: stock.Columns(),
		db:      GetEasyDB(),
	}
}

// Stock is dao
type Stock struct {
	Table   string
	Columns []string
	db      *EasyDB
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

func (dao *Stock) SelectMinMaxSalePriceByGoodsID(goodsID int64) (*model.StockMinMax, error) {
	stocks := &model.StockMinMax{}
	err := dao.db.SelectDSL(&stocks, []string{"MIN(sale_unit_price) AS min_sale_price", "MAX(sale_unit_price) AS max_sale_price"}, dao.Table, sq.Eq{"goods_id": goodsID})
	if err != nil {
		return nil, err
	}
	return stocks, nil
}

func (dao *Stock) SelectByGoodsIDWithPriceASC(goodsID int64) ([]*model.Stock, error) {
	stocks := []*model.Stock{}
	columns := dao.Columns
	columns = append(columns, "(sale_unit_price - cost_unit_price) AS profit_price")
	err := dao.db.SelectDSL(&stocks, columns, dao.Table, sq.Eq{"goods_id": goodsID}, "profit_price ASC")
	if err != nil {
		return nil, err
	}
	return stocks, nil
}

func (dao *Stock) SelectByGoodsIDs(goodsIDs []int64) (*model.StockSet, error) {
	stocks := []*model.Stock{}
	err := dao.db.SelectDSL(&stocks, dao.Columns, dao.Table, sq.Eq{"goods_id": goodsIDs, "status": model.StockStatusOnSale})
	if err != nil {
		return nil, err
	}
	return model.NewStockSet(stocks), nil
}

func (dao *Stock) SelectByIDs(ids []int64) ([]*model.Stock, error) {
	stocks := []*model.Stock{}
	err := dao.db.SelectDSL(&stocks, dao.Columns, dao.Table, sq.Eq{"id": ids, "status": model.StockStatusOnSale})
	if err != nil {
		return nil, err
	}
	return stocks, nil
}
