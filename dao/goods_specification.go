package dao

import (
	"pitaya-wechat-service/model"
	"pitaya-wechat-service/sys"

	sq "github.com/Masterminds/squirrel"
)

// GoodsSpecificationDaoSingleton is a singleton of goods dao
var GoodsSpecificationDaoSingleton *GoodsSpecificationDao

func init() {
	GoodsSpecificationDaoSingleton = new(GoodsSpecificationDao)
	GoodsSpecificationDaoSingleton.db = sys.GetEasyDB()
}

var columns_goods_specs = []string{"id", "goods_id", "specification_id", "value", "pic_url"}

// GoodsSpecificationDao is dao
type GoodsSpecificationDao struct {
	db *sys.EasyDB
}

func (dao *GoodsSpecificationDao) SelectByGoodsID(goodsID int64) ([]*model.GoodsSpecification, error) {
	goodsSpecs := []*model.GoodsSpecification{}
	err := dao.db.SelectDSL(&goodsSpecs, columns_goods_specs, model.Table_Goods_Specification, sq.Eq{"goods_id": goodsID})
	if err != nil {
		return nil, err
	}
	return goodsSpecs, nil
}
