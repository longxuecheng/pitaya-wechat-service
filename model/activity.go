package model

import (
	"database/sql"
	"gotrue/facility/utils"
	"time"
)

const (
	ActivityTypeCoupon ActivityType = "GrabCoupon"
)

type ActivityType string

type Activity struct {
	ID         int64          `db:"id"`
	Title      string         `db:"title"`
	BannerURL  string         `db:"banner_url"`
	BgURL      sql.NullString `db:"bg_url"`
	PageRoute  string         `db:"page_route"`
	Type       ActivityType   `db:"type"`
	StartTime  time.Time      `db:"start_time"`
	ExpireTime time.Time      `db:"expire_time"`
	IsDelete   bool           `db:"is_delete"`
	IsOnline   bool           `db:"is_online"`
}

func (c *Activity) TableName() string {
	return "activity"
}

func (c *Activity) Columns() []string {
	return utils.TagValues(c, "db")
}

func (c *Activity) IsAvailable() bool {
	if c.StartTime.After(time.Now()) {
		return false
	}
	if c.ExpireTime.Before(time.Now()) {
		return false
	}
	return true
}
