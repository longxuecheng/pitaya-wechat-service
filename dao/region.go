package dao

import (
	"pitaya-wechat-service/model"
	"pitaya-wechat-service/sys"

	sq "github.com/Masterminds/squirrel"
)

// RegionDaoSingleton is a singleton of goods dao
var RegionDaoSingleton *RegionDao

func RegionDaoInstance() *RegionDao {
	if RegionDaoSingleton != nil {
		return RegionDaoSingleton
	}
	RegionDaoSingleton = new(RegionDao)
	RegionDaoSingleton.db = sys.GetEasyDB()
	return RegionDaoSingleton
}

var columns_region_all = []string{"id", "parent_id", "name", "type", "agency_id"}

// RegionDao is dao
type RegionDao struct {
	db *sys.EasyDB
}

func (dao *RegionDao) SelectByParentID(parentID int) ([]model.Region, error) {
	regions := []model.Region{}
	err := dao.db.SelectDSL(&regions, columns_region_all, model.Table_Region, sq.Eq{"parent_id": parentID})
	return regions, err
}
