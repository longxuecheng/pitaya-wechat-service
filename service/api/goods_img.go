package api

import "gotrue/dto/response"

type IGoodsImgService interface {
	GetByGoodsID(goodsID int64) ([]*response.GoodsGalleryDTO, error)
}
