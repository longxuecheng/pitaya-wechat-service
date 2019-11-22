package dao

import (
	"fmt"
	"gotrue/model"
	
	"strings"
)

// GoodsAttributeDao is a singleton of goods dao
var GoodsAttributeDao *GoodsAttribute

func initGoodsAttributeDao() {
	GoodsAttributeDao = &GoodsAttribute{
		db: GetEasyDB(),
	}
}

var columns_goods_attribute = []string{"id", "goods_id", "attribute_id", "value"}

// GoodsAttribute is dao
type GoodsAttribute struct {
	db *EasyDB
}

func (dao *GoodsAttribute) SelectByGoodsID(goodsID int64) ([]*model.GoodsAttribute, error) {
	attributes := []*model.GoodsAttribute{}
	err := dao.db.Select(&attributes, fmt.Sprintf("SELECT %s FROM goods_attribute WHERE goods_id = ?", strings.Join(columns_goods_attribute, ",")), goodsID)
	if err != nil {
		return nil, err
	}
	return attributes, nil
}
