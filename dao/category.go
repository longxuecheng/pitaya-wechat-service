package dao

import (
	"fmt"
	"gotrue/model"
	"strings"

	"github.com/Masterminds/squirrel"
)

// CategoryDao is a singleton of category dao
var CategoryDao *Category

func initCategoryDao() {
	CategoryDao = &Category{db: GetEasyDB()}
}

var columns_category = []string{"id", "name", "parent_id", "wap_banner_url", "banner_url", "front_name", "icon_url"}

// Category is dao
type Category struct {
	db *EasyDB
}

func (dao *Category) SelectAllTopCategories() ([]*model.Category, error) {
	categories := []*model.Category{}
	err := dao.db.Select(&categories, fmt.Sprintf("SELECT %s FROM goods_category WHERE parent_id = 0 ORDER BY id ASC", strings.Join(columns_category, ",")))
	if err != nil {
		return nil, err
	}
	return categories, err
}

func (dao *Category) SelectAll() ([]*model.Category, error) {
	categories := []*model.Category{}
	err := dao.db.Select(&categories, fmt.Sprintf("SELECT %s FROM goods_category ORDER BY id ASC", strings.Join(columns_category, ",")))
	if err != nil {
		return nil, err
	}
	return categories, err
}

func (dao *Category) SelectByID(ID int) (*model.Category, error) {
	category := new(model.Category)
	err := dao.db.SelectOne(category, fmt.Sprintf("SELECT %s FROM goods_category WHERE id = ?", strings.Join(columns_category, ",")), ID)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (dao *Category) SelectByParentID(parentID int64) ([]*model.Category, error) {
	categories := []*model.Category{}
	err := dao.db.Select(&categories, fmt.Sprintf("SELECT %s FROM goods_category WHERE parent_id = ?", strings.Join(columns_category, ",")), parentID)
	if err != nil {
		return nil, err
	}
	return categories, err
}

func (dao *Category) SelectChannels() ([]*model.Category, error) {
	categories := []*model.Category{}
	err := dao.db.SelectDSL(&categories, columns_category, model.Table_Category, squirrel.Eq{"is_show": true})
	if err != nil {
		return nil, err
	}
	return categories, err
}

func (dao *Category) QueyMapByCategoryIDList(idList []int64) (model.CategoryMap, error) {
	categoryList := model.CategoryList{}
	err := dao.db.SelectDSL(&categoryList, columns_category, model.Table_Category, squirrel.Eq{"id": idList})
	if err != nil {
		return nil, err
	}
	return categoryList.CategoryMap(), err
}
