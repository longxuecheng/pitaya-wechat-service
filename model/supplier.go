package model

import "gotrue/facility/utils"

type Supplier struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func (s *Supplier) TableName() string {
	return "supplier"
}

func (s *Supplier) Columns() []string {
	return utils.TagValues(s, "db")
}
