package api

import "pitaya-wechat-service/dto"

// UserService is a user service interface
type UserService interface {
	GetList() ([]*dto.UserDTO, error)
}
