package api

import "pitaya-wechat-service/dto"

type IGoodsService interface {
	GetGoodsByCategory(categoryID int64) ([]*dto.GoodsItemDTO, error)
}
