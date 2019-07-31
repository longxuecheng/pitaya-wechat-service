package dao

import (
	"gotrue/model"
	"gotrue/sys"

	sq "github.com/Masterminds/squirrel"
)

// SpecificationDao is a singleton of goods dao
var SpecificationDao *Specification

func initSpecificationDao() {
	SpecificationDao = &Specification{
		db: sys.GetEasyDB(),
	}
}

var columns_specification = []string{"id", "name"}

// Specification is dao
type Specification struct {
	db *sys.EasyDB
}

func (dao *Specification) SelectByIDs(IDs []int64) ([]*model.Specification, error) {
	specs := []*model.Specification{}
	err := dao.db.SelectDSL(&specs, columns_specification, model.Table_Specification, sq.Eq{"id": IDs})
	if err != nil {
		return nil, err
	}
	return specs, nil
}
