package dao

import (
	"pitaya-wechat-service/model"
	"pitaya-wechat-service/sys"

	sq "github.com/Masterminds/squirrel"
)

// GoodsStockDaoSingleton is a singleton of goods dao
var GoodsStockDaoSingleton *GoodsStockDao

func init() {
	GoodsStockDaoSingleton = new(GoodsStockDao)
	GoodsStockDaoSingleton.db = sys.GetEasyDB()
}

var columns_goods_stock = []string{"id", "goods_id", "sale_unit_price", "available_quantity", "specification"}

// GoodsStockDao is dao
type GoodsStockDao struct {
	db *sys.EasyDB
}

func (dao *GoodsStockDao) SelectByGoodsID(goodsID int64) ([]*model.GoodsStock, error) {
	stocks := []*model.GoodsStock{}
	err := dao.db.SelectDSL(&stocks, columns_goods_stock, model.Table_Stock, sq.Eq{"goods_id": goodsID})
	if err != nil {
		return nil, err
	}
	return stocks, nil
}
