package model

import "gotrue/facility/utils"

type SupplierAdmin struct {
	ID         int64 `db:"id" omitinsert:"true"`
	UserID     int64 `db:"user_id"`
	Priority   int   `db:"priority"`
	SupplierID int64 `db:"supplier_id"`
}

func (s *SupplierAdmin) TableName() string {
	return "supplier_admin"
}

func (s *SupplierAdmin) Columns() []string {
	return utils.TagValues(s, "db")
}

type SupplierAdminCollection struct {
	Items []*SupplierAdmin
}

func (c *SupplierAdminCollection) SupplierIDs() []int64 {
	distinctMap := map[int64]bool{}
	for _, item := range c.Items {
		distinctMap[item.SupplierID] = true
	}
	ids := []int64{}
	for v := range distinctMap {
		ids = append(ids, v)
	}
	return ids
}
