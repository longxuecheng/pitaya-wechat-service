package basic

import (
	"gotrue/dao"
	"gotrue/dto"
	"gotrue/model"
)

var AttributeService *Attribute

func initAttributeService() {
	AttributeService = new(Attribute)
	AttributeService.dao = dao.AttributeDao
}

// Attribute 作为类目服务，实现了api.IAttribute
type Attribute struct {
	dao *dao.Attribute
}

func (cs *Attribute) GetByIDs(IDs []int64) ([]*dto.AttributeDTO, error) {
	attributes, err := cs.dao.SelectByIDs(IDs)
	if err != nil {
		return nil, err
	}
	return buildAttributeDTOs(attributes), nil
}

func installAttributeDTO(model *model.Attribute) *dto.AttributeDTO {
	dto := new(dto.AttributeDTO)
	dto.ID = model.ID
	dto.Name = model.Name
	dto.Value = model.Value
	return dto
}

func buildAttributeDTOs(models []*model.Attribute) []*dto.AttributeDTO {
	if models == nil || len(models) == 0 {
		return nil
	}
	dtos := make([]*dto.AttributeDTO, len(models))
	for i, model := range models {
		dtos[i] = installAttributeDTO(model)
	}
	return dtos
}
