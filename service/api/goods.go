package api

import (
	"gotrue/dto/response"
)

// IGoodsService 商品服务
type IGoodsService interface {
	GetGoodsByCategory(categoryID int64) ([]*response.GoodsItemDTO, error)

	Gallery(goodsID int64) ([]response.GoodsGalleryDTO, error)

	ExpressConstraint(stockID int64, provinceID int) (*response.GoodsExpressConstraint, error)

	Info(goodsID int64) (*response.GoodsInfoDTO, error)

	Attributes(goodsID int64) ([]*response.AttributeDTO, error)

	Specifications(goodsID int64) ([]*response.GoodsSpecificationDTO, error)

	// HotGoods get hot saled goods
	HotGoods() ([]*response.GoodsItemDTO, error)
}
