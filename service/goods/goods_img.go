package goods

import (
	"gotrue/dao"
	"gotrue/dto/response"
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

func (s *GoodsImg) GetByGoodsID(goodsID int64) ([]*response.GoodsGalleryDTO, error) {
	imgs, err := s.dao.SelectByGoodsID(goodsID)
	if err != nil {
		return nil, err
	}
	return buildGalleryDTOs(imgs), nil
}

func installGalleryDTO(model *model.GoodsImg) *response.GoodsGalleryDTO {
	data := new(response.GoodsGalleryDTO)
	data.ID = model.ID
	data.ImgDesc = model.Name
	data.ImgURL = model.URL
	data.SortOrder = model.DisplayOrder
	return data
}

func buildGalleryDTOs(models []*model.GoodsImg) []*response.GoodsGalleryDTO {
	if models == nil || len(models) == 0 {
		return nil
	}
	dtos := make([]*response.GoodsGalleryDTO, len(models))
	for i, model := range models {
		dtos[i] = installGalleryDTO(model)
	}
	return dtos
}
