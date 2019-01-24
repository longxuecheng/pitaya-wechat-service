package api

import (
	"pitaya-wechat-service/dto/response"
)

// IRegionService 是区域相关服务
type IRegionService interface {
	GetRegionsByParentID(parentID int) ([]response.RegionDTO, error)
}
