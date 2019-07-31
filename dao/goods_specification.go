package dao

import (
	"gotrue/model"
	"gotrue/sys"

	sq "github.com/Masterminds/squirrel"
)

// GoodsSpecDao is a singleton of goods dao
var GoodsSpecDao *GoodsSpec

func initGoodsSpecDao() {
	GoodsSpecDao = &GoodsSpec{
		db: sys.GetEasyDB(),
	}
}

var columns_goods_specs = []string{"id", "goods_id", "specification_id", "value", "pic_url"}

// GoodsSpec is dao
type GoodsSpec struct {
	db *sys.EasyDB
}

func (dao *GoodsSpec) SelectByGoodsID(goodsID int64) ([]*model.GoodsSpecification, error) {
	goodsSpecs := []*model.GoodsSpecification{}
	err := dao.db.SelectDSL(&goodsSpecs, columns_goods_specs, model.Table_Goods_Specification, sq.Eq{"goods_id": goodsID})
	if err != nil {
		return nil, err
	}
	return goodsSpecs, nil
}
