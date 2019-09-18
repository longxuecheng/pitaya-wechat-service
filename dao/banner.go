package dao

import (
	"gotrue/model"
	"gotrue/sys"

	"github.com/Masterminds/squirrel"
)

var BannerDao *Banner

func initBannerDao() {
	banner := &model.Banner{}
	BannerDao = &Banner{
		table:   banner.TableName(),
		columns: banner.Columns(),
		db:      sys.GetEasyDB(),
	}
}

type Banner struct {
	table   string
	columns []string
	db      *sys.EasyDB
}

func (d *Banner) QueryOnlineBanners() ([]*model.Banner, error) {
	banners := []*model.Banner{}
	return banners, d.db.SelectDSL(&banners, d.columns, d.table, squirrel.Eq{"is_online": true}, "priority ASC")
}
