package api

import (
	"gotrue/dto/response"
)

// IRegionService 是区域相关服务
type IRegionService interface {
	GetRegionsByParentID(parentID int) ([]response.RegionDTO, error)
}
