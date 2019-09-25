package api

import (
	"gotrue/dto/response"
	"io"
)

// IGoodsService 商品服务
type IGoodsService interface {
	GetGoodsByCategory(categoryID int64) ([]*response.GoodsItem, error)

	Gallery(goodsID int64) ([]response.GoodsGalleryDTO, error)

	ExpressConstraint(stockID int64, provinceID int) (*response.GoodsExpressConstraint, error)

	Info(goodsID int64) (*response.GoodsInfo, error)

	Attributes(goodsID int64) ([]*response.AttributeDTO, error)

	Specifications(goodsID int64) ([]*response.GoodsSpecificationDTO, error)

	// HotGoods get hot saled goods
	HotGoods() ([]*response.GoodsItem, error)

	OneSaleGoodsCards() ([]*response.GoodsCard, error)

	ImportExpressConstraints(stockID int64, r io.Reader) error
}
