package api

import "gotrue/dto/response"

type IAttributeService interface {
	GetByIDs(IDs []int64) ([]*response.AttributeDTO, error)
}
