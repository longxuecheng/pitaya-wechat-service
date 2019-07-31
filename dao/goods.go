package dao

import (
	"fmt"
	"gotrue/model"
	"gotrue/sys"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

// GoodsDao is a singleton of goods dao
var GoodsDao *Goods

func initGoodsDao() {
	GoodsDao = &Goods{
		db: sys.GetEasyDB(),
	}
}

var columns_goods = []string{"id", "name", "category_id", "description", "brief_description", "status", "is_delete", "supplier_id", "create_time", "list_pic_url", "retail_price"}

// Goods is dao
type Goods struct {
	db *sys.EasyDB
}

func (dao *Goods) SelectByCategory(categoryID int64) ([]*model.Goods, error) {
	goods := []*model.Goods{}
	err := dao.db.Select(&goods, fmt.Sprintf("SELECT %s FROM goods WHERE category_id = ? ORDER BY id ASC", strings.Join(columns_goods, ",")), categoryID)
	if err != nil {
		return nil, err
	}
	return goods, nil
}

func (dao *Goods) SelectByID(ID int64) (*model.Goods, error) {
	goods := new(model.Goods)
	err := dao.db.SelectOne(goods, fmt.Sprintf("SELECT %s FROM goods WHERE id = ?", strings.Join(columns_goods, ",")), ID)
	if err != nil {
		return nil, err
	}
	return goods, nil
}

func (dao *Goods) SelectByIDs(IDs []int64) ([]*model.Goods, error) {
	goods := []*model.Goods{}
	err := dao.db.SelectDSL(&goods, columns_goods, model.Table_Goods, sq.Eq{"id": IDs})
	if err != nil {
		return nil, err
	}
	return goods, nil
}

func (dao *Goods) SelectAllByStatus(status model.GoodsStatus) ([]*model.Goods, error) {
	goods := []*model.Goods{}
	err := dao.db.SelectDSL(&goods, columns_goods, model.Table_Goods, sq.Eq{"status": string(status)})
	if err != nil {
		return nil, err
	}
	return goods, nil
}
