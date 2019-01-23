package dao

import (
	"pitaya-wechat-service/model"
	"pitaya-wechat-service/sys"

	sq "github.com/Masterminds/squirrel"
)

// SpecificationDaoSingleton is a singleton of goods dao
var SpecificationDaoSingleton *SpecificationDao

func init() {
	SpecificationDaoSingleton = new(SpecificationDao)
	SpecificationDaoSingleton.db = sys.GetEasyDB()
}

var columns_specification = []string{"id", "name"}

// SpecificationDao is dao
type SpecificationDao struct {
	db *sys.EasyDB
}

func (dao *SpecificationDao) SelectByIDs(IDs []int64) ([]*model.Specification, error) {
	specs := []*model.Specification{}
	err := dao.db.SelectDSL(&specs, columns_specification, model.Table_Specification, sq.Eq{"id": IDs})
	if err != nil {
		return nil, err
	}
	return specs, nil
}
