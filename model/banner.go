package model

import "gotrue/facility/utils"

type Banner struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	Src      string `db:"src"`
	Type     string `db:"type"`
	Link     string `db:"link"`
	IsOnline bool   `db:"is_online"`
}

func (b *Banner) TableName() string {
	return "banner"
}

func (b *Banner) Columns() []string {
	return utils.TagValues(b, "db")
}
