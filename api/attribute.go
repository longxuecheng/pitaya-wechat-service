package api

import "pitaya-wechat-service/dto"

type IAttributeService interface {
	GetByIDs(IDs []int64) ([]*dto.AttributeDTO, error)
}
