package dao

import (
	"fmt"
	"pitaya-wechat-service/model"
	"pitaya-wechat-service/sys"
	"strings"
)

// GoodsAttributeDaoSingleton is a singleton of goods dao
var GoodsAttributeDaoSingleton *GoodsAttributeDao

func init() {
	GoodsAttributeDaoSingleton = new(GoodsAttributeDao)
	GoodsAttributeDaoSingleton.db = sys.GetEasyDB()
}

var columns_goods_attribute = []string{"id", "goods_id", "attribute_id", "value"}

// GoodsAttributeDao is dao
type GoodsAttributeDao struct {
	db *sys.EasyDB
}

func (dao *GoodsAttributeDao) SelectByGoodsID(goodsID int64) ([]*model.GoodsAttribute, error) {
	attributes := []*model.GoodsAttribute{}
	err := dao.db.Select(&attributes, fmt.Sprintf("SELECT %s FROM goods_attribute WHERE goods_id = ?", strings.Join(columns_goods_attribute, ",")), goodsID)
	if err != nil {
		return nil, err
	}
	return attributes, nil
}
