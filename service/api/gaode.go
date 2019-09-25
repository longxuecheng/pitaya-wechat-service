package api

import "gotrue/dto/response"

type IGaodeMapService interface {
	AddressLocation(city string, address string) (*response.Location, error)
	AddressTips(keywords, lng, lat string) ([]string, error)
	ReverseGeo() error
}
