package api

import (
	"gotrue/dto/response"
	"io"
)

// IGoodsService 商品服务
type IGoodsService interface {
	GetGoodsByCategory(categoryID int64) ([]*response.GoodsItem, error)

	GetInternalGoodsByCategory(categoryID int64) ([]*InternalGoods, error)

	Gallery(goodsID int64) ([]response.GoodsGalleryDTO, error)

	ExpressConstraint(stockID int64, provinceID int) (*response.GoodsExpressConstraint, error)

	Info(goodsID int64) (*response.GoodsInfo, error)

	Attributes(goodsID int64) ([]*response.AttributeDTO, error)

	Specifications(goodsID int64) ([]*response.GoodsSpecificationDTO, error)

	// HotGoods get hot saled goods
	HotGoods() ([]response.HotGoods, error)

	OneSaleGoodsCards() ([]*response.GoodsCard, error)

	ImportExpressConstraints(stockID int64, r io.Reader) error
}

type InternalGoods struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
