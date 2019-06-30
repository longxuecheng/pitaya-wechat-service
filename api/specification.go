package api

import "gotrue/dto"

// ISpecificationService 规格属性服务
type ISpecificationService interface {
	GetByIDs(IDs []int64) ([]*dto.SpecificationDTO, error)
}
