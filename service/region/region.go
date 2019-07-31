package region

import (
	"gotrue/dao"
	"gotrue/dto/response"
	"gotrue/model"
	"strings"
)

var RegionService *regionService

func Init() {
	RegionService = &regionService{
		dao.RegionDaoInstance(),
	}
}

type regionService struct {
	dao *dao.RegionDao
}

func (rs *regionService) FullName(regionIDs []int) (string, error) {
	regions, err := rs.dao.SelectByIDs(regionIDs)
	if err != nil {
		return "", err
	}
	regionNames := []string{}
	for _, r := range regions {
		regionNames = append(regionNames, r.Name)
	}
	return strings.Join(regionNames, "-"), nil
}

func (s *regionService) GetRegionsByParentID(parentID int) ([]response.Region, error) {
	regions, err := s.dao.SelectByParentID(parentID)
	if err != nil {
		return nil, err
	}
	return regionResponses(regions), err
}

func regionResponse(model model.Region) response.Region {
	dto := response.Region{}
	dto.ID = model.ID
	dto.ParentID = model.ParentID
	dto.Name = model.Name
	dto.Type = int(model.Type)
	dto.AgencyID = model.AgencyID
	return dto
}

func regionResponses(models []model.Region) []response.Region {
	if models == nil || len(models) == 0 {
		return nil
	}
	dtos := make([]response.Region, len(models))
	for i, model := range models {
		dtos[i] = regionResponse(model)
	}
	return dtos
}
