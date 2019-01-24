package model

var Table_Region = "region"

// RegionType 区域类型
type RegionType int

var (
	RegionType_Country  RegionType = 0 // country
	RegionType_City     RegionType = 1 // city
	RegionType_District RegionType = 2 // ditrict
)

type Region struct {
	ID       int        `db:"id"`
	ParentID int        `db:"parent_id"`
	Name     string     `db:"name"`
	Type     RegionType `db:"type"`
	AgencyID int        `db:"agency_id"`
}
