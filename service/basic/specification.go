package basic

import (
	"gotrue/dao"
	"gotrue/dto/response"
	"gotrue/model"
)

var SpecificationService *Specification

func initSpecificationService() {
	SpecificationService = &Specification{
		dao: dao.SpecificationDao,
	}
}

// Specification 作为规格服务，实现了api.ISpecification
type Specification struct {
	dao *dao.Specification
}

func (cs *Specification) GetByIDs(IDs []int64) ([]*response.SpecificationDTO, error) {
	Specifications, err := cs.dao.SelectByIDs(IDs)
	if err != nil {
		return nil, err
	}
	return buildSpecificationDTOs(Specifications), nil
}

func installSpecificationDTO(model *model.Specification) *response.SpecificationDTO {
	dto := new(response.SpecificationDTO)
	dto.ID = model.ID
	dto.Name = model.Name
	return dto
}

func buildSpecificationDTOs(models []*model.Specification) []*response.SpecificationDTO {
	if models == nil || len(models) == 0 {
		return nil
	}
	dtos := make([]*response.SpecificationDTO, len(models))
	for i, model := range models {
		dtos[i] = installSpecificationDTO(model)
	}
	return dtos
}
