package dao

import (
	"gotrue/model"
	"gotrue/sys"

	sq "github.com/Masterminds/squirrel"
)

// AttributeDaoSingleton is a singleton of goods dao
var AttributeDaoSingleton *AttributeDao

func init() {
	AttributeDaoSingleton = new(AttributeDao)
	AttributeDaoSingleton.db = sys.GetEasyDB()
}

var columns_attribute = []string{"id", "name", "value", "sort_order"}

// AttributeDao is dao
type AttributeDao struct {
	db *sys.EasyDB
}

func (dao *AttributeDao) SelectByIDs(IDs []int64) ([]*model.Attribute, error) {
	abs := []*model.Attribute{}
	sql, args, err := sq.Select(columns_attribute...).From("attribute").Where(sq.Eq{"id": IDs}).ToSql()
	if err != nil {
		return nil, err
	}
	err = dao.db.Select(&abs, sql, args...)
	if err != nil {
		return nil, err
	}
	return abs, nil
}
