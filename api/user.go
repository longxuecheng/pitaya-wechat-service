package api

import (
	"pitaya-wechat-service/dto"
	"pitaya-wechat-service/dto/request"
)

// IUserService is a user service interface
type IUserService interface {
	GetList() ([]*dto.UserDTO, error)
	AddressList(userID int64) ([]dto.UserAddressDTO, error)
	CreateAddress(req request.UserAddressAddRequest) (id int64, err error)
}
