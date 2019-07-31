package dao

import (
	"gotrue/model"
	"gotrue/sys"

	sq "github.com/Masterminds/squirrel"
)

// RegionDao is a singleton of goods dao
var RegionDao *Region

func initRegionDao() {
	RegionDao = &Region{
		db: sys.GetEasyDB(),
	}
}

var columns_region_all = []string{"id", "parent_id", "name", "type", "agency_id"}

// Region is dao
type Region struct {
	db *sys.EasyDB
}

func (dao *Region) SelectByParentID(parentID int) ([]model.Region, error) {
	regions := []model.Region{}
	err := dao.db.SelectDSL(&regions, columns_region_all, model.Table_Region, sq.Eq{"parent_id": parentID})
	return regions, err
}

func (dao *Region) SelectByIDs(ids []int) ([]*model.Region, error) {
	regions := []*model.Region{}
	return regions, dao.db.SelectDSL(&regions, columns_region_all, model.Table_Region, sq.Eq{"id": ids}, "type asc")
}
