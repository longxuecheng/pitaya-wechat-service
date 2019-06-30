package service

import (
	"gotrue/dao"
	"gotrue/dto"
	"gotrue/model"
)

var SpecificationServiceSingleton *SpecificationService

// init 在此实现spring中类似注入的功能
func init() {
	SpecificationServiceSingleton = new(SpecificationService)
	SpecificationServiceSingleton.dao = dao.SpecificationDaoSingleton
}

// SpecificationService 作为规格服务，实现了api.ISpecificationService
type SpecificationService struct {
	dao *dao.SpecificationDao
}

func (cs *SpecificationService) GetByIDs(IDs []int64) ([]*dto.SpecificationDTO, error) {
	Specifications, err := cs.dao.SelectByIDs(IDs)
	if err != nil {
		return nil, err
	}
	return buildSpecificationDTOs(Specifications), nil
}

func installSpecificationDTO(model *model.Specification) *dto.SpecificationDTO {
	dto := new(dto.SpecificationDTO)
	dto.ID = model.ID
	dto.Name = model.Name
	return dto
}

func buildSpecificationDTOs(models []*model.Specification) []*dto.SpecificationDTO {
	if models == nil || len(models) == 0 {
		return nil
	}
	dtos := make([]*dto.SpecificationDTO, len(models))
	for i, model := range models {
		dtos[i] = installSpecificationDTO(model)
	}
	return dtos
}
