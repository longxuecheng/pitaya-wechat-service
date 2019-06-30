package dao

import (
	"fmt"
	"gotrue/model"
	"gotrue/sys"
	"strings"

	"github.com/Masterminds/squirrel"
)

// CategoryDaoSingleton is a singleton of category dao
var CategoryDaoSingleton *CategoryDao

func init() {
	CategoryDaoSingleton = new(CategoryDao)
	CategoryDaoSingleton.db = sys.GetEasyDB()
}

var columns_category = []string{"id", "name", "parent_id", "wap_banner_url", "banner_url", "front_name", "icon_url"}

// CategoryDao is dao
type CategoryDao struct {
	db *sys.EasyDB
}

func (dao *CategoryDao) SelectAllTopCategories() ([]*model.Category, error) {
	categories := []*model.Category{}
	err := dao.db.Select(&categories, fmt.Sprintf("SELECT %s FROM goods_category WHERE parent_id = 0 ORDER BY id ASC", strings.Join(columns_category, ",")))
	if err != nil {
		return nil, err
	}
	return categories, err
}

func (dao *CategoryDao) SelectAll() ([]*model.Category, error) {
	categories := []*model.Category{}
	err := dao.db.Select(&categories, fmt.Sprintf("SELECT %s FROM goods_category ORDER BY id ASC", strings.Join(columns_category, ",")))
	if err != nil {
		return nil, err
	}
	return categories, err
}

func (dao *CategoryDao) SelectByID(ID int) (*model.Category, error) {
	category := new(model.Category)
	err := dao.db.SelectOne(category, fmt.Sprintf("SELECT %s FROM goods_category WHERE id = ?", strings.Join(columns_category, ",")), ID)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (dao *CategoryDao) SelectByParentID(parentID int) ([]*model.Category, error) {
	categories := []*model.Category{}
	err := dao.db.Select(&categories, fmt.Sprintf("SELECT %s FROM goods_category WHERE parent_id = ?", strings.Join(columns_category, ",")), parentID)
	if err != nil {
		return nil, err
	}
	return categories, err
}

func (dao *CategoryDao) SelectChannels() ([]*model.Category, error) {
	categories := []*model.Category{}
	err := dao.db.SelectDSL(&categories, columns_category, model.Table_Category, squirrel.Eq{"is_show": true})
	if err != nil {
		return nil, err
	}
	return categories, err
}
