package api

import (
	"pitaya-wechat-service/dto"
	"pitaya-wechat-service/dto/request"
	"pitaya-wechat-service/model"
)

// IUserService is a user service interface
type IUserService interface {
	GetList() ([]*dto.UserDTO, error)
	AddressList(userID int64) ([]dto.UserAddress, error)
	DefaultAddress(userID int64) (dto.UserAddress, error)
	CreateAddress(userID int64, req request.UserAddressAddRequest) (id int64, err error)
	Login(openID string, nickName string, avatarURL string) (user *model.User, err error)
}
