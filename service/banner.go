package service

import (
	"gotrue/dao"
	"gotrue/dto/response"
)

var BannerService *Banner

type Banner struct {
	dao *dao.Banner
}

func InitBannerService() {
	BannerService = &Banner{
		dao: dao.BannerDao,
	}
}

func (s *Banner) OnlineBanners() ([]*response.Banner, error) {
	banners, err := s.dao.QueryOnlineBanners()
	if err != nil {
		return nil, err
	}
	apiBanners := make([]*response.Banner, len(banners))
	for i, banner := range banners {
		apiBanners[i] = &response.Banner{
			Name: banner.Name,
			Src:  banner.Src,
			Link: banner.Link,
		}
	}
	return apiBanners, nil
}
