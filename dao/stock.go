package dao

import (
	"gotrue/model"
	"gotrue/sys"

	sq "github.com/Masterminds/squirrel"
)

// GoodsStockDaoSingleton is a singleton of goods dao
var GoodsStockDaoSingleton *GoodsStockDao

func initStockDao() {
	GoodsStockDaoSingleton = &GoodsStockDao{
		db: sys.GetEasyDB(),
	}
}

var columns_goods_stock = []string{"id", "supplier_id", "goods_id", "sale_unit_price", "cost_unit_price", "available_quantity", "specification"}

// GoodsStockDao is dao
type GoodsStockDao struct {
	db *sys.EasyDB
}

func (dao *GoodsStockDao) SelectByID(ID int64) (*model.GoodsStock, error) {
	stock := new(model.GoodsStock)
	err := dao.db.SelectOneDSL(stock, columns_goods_stock, model.Table_Stock, sq.Eq{"id": ID})
	if err != nil {
		return nil, err
	}
	return stock, nil
}

func (dao *GoodsStockDao) SelectByGoodsID(goodsID int64) ([]*model.GoodsStock, error) {
	stocks := []*model.GoodsStock{}
	err := dao.db.SelectDSL(&stocks, columns_goods_stock, model.Table_Stock, sq.Eq{"goods_id": goodsID})
	if err != nil {
		return nil, err
	}
	return stocks, nil
}

func (dao *GoodsStockDao) SelectByIDs(ids []int64) ([]*model.GoodsStock, error) {
	stocks := []*model.GoodsStock{}
	err := dao.db.SelectDSL(&stocks, columns_goods_stock, model.Table_Stock, sq.Eq{"id": ids})
	if err != nil {
		return nil, err
	}
	return stocks, nil
}
