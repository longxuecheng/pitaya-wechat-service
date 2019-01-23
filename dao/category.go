package dao

import (
	"fmt"
	"pitaya-wechat-service/model"
	"pitaya-wechat-service/sys"
	"strings"
)

// CategoryDaoSingleton is a singleton of category dao
var CategoryDaoSingleton *CategoryDao

func init() {
	CategoryDaoSingleton = new(CategoryDao)
	CategoryDaoSingleton.db = sys.GetEasyDB()
}

var columns = []string{"id", "name", "parent_id", "wap_banner_url", "banner_url", "front_name"}

// CategoryDao is dao
type CategoryDao struct {
	db *sys.EasyDB
}

func (dao *CategoryDao) SelectAllTopCategories() ([]*model.Category, error) {
	categories := []*model.Category{}
	err := dao.db.Select(&categories, fmt.Sprintf("SELECT %s FROM goods_category WHERE parent_id = 0 ORDER BY id ASC", strings.Join(columns, ",")))
	if err != nil {
		return nil, err
	}
	return categories, err
}

func (dao *CategoryDao) SelectAll() ([]*model.Category, error) {
	categories := []*model.Category{}
	err := dao.db.Select(&categories, fmt.Sprintf("SELECT %s FROM goods_category ORDER BY id ASC", strings.Join(columns, ",")))
	if err != nil {
		return nil, err
	}
	return categories, err
}

func (dao *CategoryDao) SelectByID(ID int) (*model.Category, error) {
	category := new(model.Category)
	err := dao.db.SelectOne(category, fmt.Sprintf("SELECT %s FROM goods_category WHERE id = ?", strings.Join(columns, ",")), ID)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (dao *CategoryDao) SelectByParentID(parentID int) ([]*model.Category, error) {
	categories := []*model.Category{}
	err := dao.db.Select(&categories, fmt.Sprintf("SELECT %s FROM goods_category WHERE parent_id = ?", strings.Join(columns, ",")), parentID)
	if err != nil {
		return nil, err
	}
	return categories, err
}
