package api

import "pitaya-wechat-service/dto"

type IGoodsImgService interface {
	GetByGoodsID(goodsID int64) ([]*dto.GoodsGalleryDTO, error)
}
