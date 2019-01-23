package api

import "pitaya-wechat-service/dto"

// ISpecificationService 规格属性服务
type ISpecificationService interface {
	GetByIDs(IDs []int64) ([]*dto.SpecificationDTO, error)
}
