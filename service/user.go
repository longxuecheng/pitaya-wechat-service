package service

import (
	"pitaya-wechat-service/dao"
	"pitaya-wechat-service/dto"
	"pitaya-wechat-service/model"
)

var UserServiceSingleton *UserServiceImplement

// init 在此实现spring中类似注入的功能
func init() {
	UserServiceSingleton = new(UserServiceImplement)
	UserServiceSingleton.userDao = dao.UserDaoSingleton
}

type UserServiceImplement struct {
	userDao *dao.UserDao
}

func (userService *UserServiceImplement) GetList() ([]*dto.UserDTO, error) {
	users, err := userService.userDao.SelectAll()
	if err != nil {
		return nil, err
	}
	return buildUserDTOs(users), nil
}

func installUserDTO(model *model.User) *dto.UserDTO {
	userDto := new(dto.UserDTO)
	userDto.Name = model.Name
	userDto.PhoneNo = model.PhoneNo
	userDto.Email = model.Email
	return userDto
}

func buildUserDTOs(models []*model.User) []*dto.UserDTO {
	if models == nil || len(models) == 0 {
		return nil
	}
	dtos := make([]*dto.UserDTO, len(models))
	for i, model := range models {
		dtos[i] = installUserDTO(model)
	}
	return dtos
}
