package dao

import (
	"fmt"
	"pitaya-wechat-service/model"
	"pitaya-wechat-service/sys"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

// GoodsDaoSingleton is a singleton of goods dao
var GoodsDaoSingleton *GoodsDao

func init() {
	GoodsDaoSingleton = new(GoodsDao)
	GoodsDaoSingleton.db = sys.GetEasyDB()
}

var columns_goods = []string{"id", "name", "category_id", "description", "brief_description", "status", "is_delete", "supplier_id", "create_time", "list_pic_url", "retail_price"}

// GoodsDao is dao
type GoodsDao struct {
	db *sys.EasyDB
}

func (dao *GoodsDao) SelectByCategory(categoryID int64) ([]*model.Goods, error) {
	goods := []*model.Goods{}
	err := dao.db.Select(&goods, fmt.Sprintf("SELECT %s FROM goods WHERE category_id = ? ORDER BY id ASC", strings.Join(columns_goods, ",")), categoryID)
	if err != nil {
		return nil, err
	}
	return goods, nil
}

func (dao *GoodsDao) SelectByID(ID int64) (*model.Goods, error) {
	goods := new(model.Goods)
	err := dao.db.SelectOne(goods, fmt.Sprintf("SELECT %s FROM goods WHERE id = ?", strings.Join(columns_goods, ",")), ID)
	if err != nil {
		return nil, err
	}
	return goods, nil
}

func (dao *GoodsDao) SelectByIDs(IDs []int64) ([]*model.Goods, error) {
	goods := []*model.Goods{}
	err := dao.db.SelectDSL(&goods, columns_goods, model.Table_Goods, sq.Eq{"id": IDs})
	if err != nil {
		return nil, err
	}
	return goods, nil
}

func (dao *GoodsDao) SelectAllByStatus(status model.GoodsStatus) ([]*model.Goods, error) {
	goods := []*model.Goods{}
	err := dao.db.SelectDSL(&goods, columns_goods, model.Table_Goods, sq.Eq{"status": string(status)})
	if err != nil {
		return nil, err
	}
	return goods, nil
}
