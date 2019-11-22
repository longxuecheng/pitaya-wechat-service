package model

const Table_Category = "goods_category"

type Category struct {
	ID           int64  `db:"id"`
	Name         string `db:"name"`
	ParentID     int64  `db:"parent_id"`
	WapBannerURL string `db:"wap_banner_url"`
	BannerURL    string `db:"banner_url"`
	IconUrl      string `db:"icon_url"`
	FrontName    string `db:"front_name"`
}

type CategoryMap map[int64]*Category

func (m CategoryMap) Get(id int64) *Category {
	return m[id]
}

type CategoryList []*Category

func (c CategoryList) CategoryMap() CategoryMap {
	categoryMap := CategoryMap{}
	for _, category := range c {
		categoryMap[category.ID] = category
	}
	return categoryMap
}
