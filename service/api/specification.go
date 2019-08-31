package api

import "gotrue/dto/response"

// ISpecificationService 规格属性服务
type ISpecificationService interface {
	GetByIDs(IDs []int64) ([]*response.SpecificationDTO, error)
}
