package goods

import (
	"gotrue/dao"
	"gotrue/dto"
	"gotrue/model"
)

var GoodsImgService *GoodsImg

func initGoodsImgService() {
	if GoodsImgService != nil {
		return
	}
	GoodsImgService = &GoodsImg{
		dao: dao.GoodsImgDao,
	}
}

// GoodsImg 作为类目服务，实现了api.GoodsImg接口
type GoodsImg struct {
	dao *dao.GoodsImg
}

func (s *GoodsImg) GetByGoodsID(goodsID int64) ([]*dto.GoodsGalleryDTO, error) {
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
