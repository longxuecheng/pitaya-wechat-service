package service

import (
	"pitaya-wechat-service/dao"
	"pitaya-wechat-service/dto"
	"pitaya-wechat-service/model"
)

var GoodsServiceSingleton *GoodsService

// init 在此实现spring中类似注入的功能
func init() {
	GoodsServiceSingleton = new(GoodsService)
	GoodsServiceSingleton.dao = dao.GoodsDaoSingleton
}

// GoodsService 作为类目服务，实现了api.GoodsService接口
type GoodsService struct {
	dao *dao.GoodsDao
}

func (s *GoodsService) GetGoodsByCategory(categoryID int64) ([]*dto.GoodsItemDTO, error) {
	goods, err := s.dao.SelectByCategory(categoryID)
	if err != nil {
		return nil, err
	}
	return buildGoodsDTOs(goods), nil
}

func installGoodsDTO(model *model.Goods) *dto.GoodsItemDTO {
	dto := new(dto.GoodsItemDTO)
	dto.ID = model.ID
	dto.Name = model.Name
	dto.PicURL = model.ListPicURL.String
	dto.RetailPrice = model.RetailPrice
	return dto
}

func buildGoodsDTOs(models []*model.Goods) []*dto.GoodsItemDTO {
	if models == nil || len(models) == 0 {
		return nil
	}
	dtos := make([]*dto.GoodsItemDTO, len(models))
	for i, model := range models {
		dtos[i] = installGoodsDTO(model)
	}
	return dtos
}
