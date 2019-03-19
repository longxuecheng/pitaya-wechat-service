package api

import "pitaya-wechat-service/dto"

// IGoodsService 商品服务
type IGoodsService interface {
	GetGoodsByCategory(categoryID int64) ([]*dto.GoodsItemDTO, error)
	Gallery(goodsID int64) ([]dto.GoodsGalleryDTO, error)
	Info(goodsID int64) (*dto.GoodsInfoDTO, error)
	Attributes(goodsID int64) ([]*dto.AttributeDTO, error)
	Specifications(goodsID int64) ([]*dto.GoodsSpecificationDTO, error)
	// HotGoods get hot saled goods
	HotGoods() ([]*dto.GoodsItemDTO, error)
}
