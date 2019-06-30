package api

import "gotrue/dto"

type IAttributeService interface {
	GetByIDs(IDs []int64) ([]*dto.AttributeDTO, error)
}
