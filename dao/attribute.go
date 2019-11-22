package dao

import (
	"gotrue/model"

	sq "github.com/Masterminds/squirrel"
)

// AttributeDao is a singleton of goods dao
var AttributeDao *Attribute

func initAttributeDao() {
	AttributeDao = &Attribute{
		db: GetEasyDB(),
	}
}

var columns_attribute = []string{"id", "name", "value", "sort_order"}

// AttributeDao is dao
type Attribute struct {
	db *EasyDB
}

func (dao *Attribute) SelectByIDs(IDs []int64) ([]*model.Attribute, error) {
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
