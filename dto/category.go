package dto

// CategoryDTO 是商品分类的DTO模型
type CategoryDTO struct {
	ID              int            `json:"id"`
	Name            string         `json:"name"`
	BannerURL       string         `json:"bannerUrl"`
	WapBannerURL    string         `json:"wapBannerUrl"`
	FrontName       string         `json:"frontName"`
	SubCategoryList []*CategoryDTO `json:"subCategoryList"`
	ParentID        int            `json:"parentId"`
}

// CategoryChannel 是商城首页展示分类频道的信息
type CategoryChannel struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	IconURL string `json:"iconUrl"`
}
