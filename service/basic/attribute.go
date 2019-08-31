package basic

import (
	"gotrue/dao"
	"gotrue/dto/response"
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

func (cs *Attribute) GetByIDs(IDs []int64) ([]*response.AttributeDTO, error) {
	attributes, err := cs.dao.SelectByIDs(IDs)
	if err != nil {
		return nil, err
	}
	return buildAttributeDTOs(attributes), nil
}

func installAttributeDTO(model *model.Attribute) *response.AttributeDTO {
	dto := new(response.AttributeDTO)
	dto.ID = model.ID
	dto.Name = model.Name
	dto.Value = model.Value
	return dto
}

func buildAttributeDTOs(models []*model.Attribute) []*response.AttributeDTO {
	if models == nil || len(models) == 0 {
		return nil
	}
	dtos := make([]*response.AttributeDTO, len(models))
	for i, model := range models {
		dtos[i] = installAttributeDTO(model)
	}
	return dtos
}
