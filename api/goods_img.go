package api

import "gotrue/dto"

type IGoodsImgService interface {
	GetByGoodsID(goodsID int64) ([]*dto.GoodsGalleryDTO, error)
}
