package service

import (
	"pitaya-wechat-service/dao"
	"pitaya-wechat-service/dto/response"
	"pitaya-wechat-service/model"
)

// RegionServiceSingleton 是RegionService的一个单例
var RegionServiceSingleton *RegionService

// init 在此实现spring中类似注入的功能
func init() {
	instantiateRegionService()
}

func instantiateRegionService() {
	RegionServiceSingleton = new(RegionService)
	RegionServiceSingleton.dao = dao.RegionDaoInstance()
}

func RegionServiceInstance() *RegionService {
	if RegionServiceSingleton == nil {
		instantiateRegionService()
	}
	return RegionServiceSingleton
}

// RegionService 作为规格服务，实现了api.IRegionService
type RegionService struct {
	dao *dao.RegionDao
}

func (s *RegionService) GetRegionsByParentID(parentID int) ([]response.RegionDTO, error) {
	regions, err := s.dao.SelectByParentID(parentID)
	if err != nil {
		return nil, err
	}
	return buildRegionDTOs(regions), err
}

func installRegionDTO(model model.Region) response.RegionDTO {
	dto := response.RegionDTO{}
	dto.ID = model.ID
	dto.ParentID = model.ParentID
	dto.Name = model.Name
	dto.Type = int(model.Type)
	dto.AgencyID = model.AgencyID
	return dto
}

func buildRegionDTOs(models []model.Region) []response.RegionDTO {
	if models == nil || len(models) == 0 {
		return nil
	}
	dtos := make([]response.RegionDTO, len(models))
	for i, model := range models {
		dtos[i] = installRegionDTO(model)
	}
	return dtos
}
