package model

type Category struct {
	ID           int    `db:"id"`
	Name         string `db:"name"`
	ParentID     int    `db:"parent_id"`
	WapBannerURL string `db:"wap_banner_url"`
	BannerURL    string `db:"banner_url"`
	FrontName    string `db:"front_name"`
}
