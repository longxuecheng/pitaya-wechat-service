package dto

import (
	"github.com/shopspring/decimal"
)

type GoodsItemDTO struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	PicURL      string          `json:"picUrl"`
	RetailPrice decimal.Decimal `json:"retailPrice"`
}

type GoodsInfoDTO struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	BriefDesc   string          `json:"briefDesc"`
	RetailPrice decimal.Decimal `json:"retailPrice"`
	Unit        string          `json:"unit"`
	Description string          `json:"description"`
}

type GoodsGalleryDTO struct {
	ID        int64  `json:"id"`
	GoodsID   int64  `json:"goodsId"`
	ImgURL    string `json:"imgUrl"`
	ImgDesc   string `json:"imgDesc"`
	SortOrder int    `json:"sortOrder"`
}
