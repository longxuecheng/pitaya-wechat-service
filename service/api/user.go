package api

import (
	"gotrue/dto/request"
	"gotrue/dto/response"
	"gotrue/model"
)

// IUserService is a user service interface
type IUserService interface {
	GetList() ([]*response.UserDTO, error)
	AddressList(userID int64) ([]*response.UserAddress, error)
	DefaultAddress(userID int64) (response.UserAddress, error)
	CreateAddress(userID int64, req request.UserAddressAddRequest) (id int64, err error)
	Login(openID string, nickName string, avatarURL string) (user *model.User, err error)
}
