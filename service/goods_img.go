package service

import (
	"pitaya-wechat-service/dao"
	"pitaya-wechat-service/dto"
	"pitaya-wechat-service/model"
)

var GoodsImgServiceSingleton *GoodsImgService

// init 在此实现spring中类似注入的功能
func init() {
	GoodsImgServiceSingleton = new(GoodsImgService)
	GoodsImgServiceSingleton.dao = dao.GoodsImgDaoSingleton
}

// GoodsImgService 作为类目服务，实现了api.GoodsImgService接口
type GoodsImgService struct {
	dao *dao.GoodsImgDao
}

func (s *GoodsImgService) GetByGoodsID(goodsID int64) ([]*dto.GoodsGalleryDTO, error) {
	imgs, err := s.dao.SelectByGoodsID(goodsID)
	if err != nil {
		return nil, err
	}
	return buildGalleryDTOs(imgs), nil
}

func installGalleryDTO(model *model.GoodsImg) *dto.GoodsGalleryDTO {
	dto := new(dto.GoodsGalleryDTO)
	dto.ID = model.ID
	dto.ImgDesc = model.Name
	dto.ImgURL = model.URL
	dto.SortOrder = model.DisplayOrder
	return dto
}

func buildGalleryDTOs(models []*model.GoodsImg) []*dto.GoodsGalleryDTO {
	if models == nil || len(models) == 0 {
		return nil
	}
	dtos := make([]*dto.GoodsGalleryDTO, len(models))
	for i, model := range models {
		dtos[i] = installGalleryDTO(model)
	}
	return dtos
}
