package service

import (
	"pitaya-wechat-service/dao"
	"pitaya-wechat-service/dto"
	"pitaya-wechat-service/dto/request"
	"pitaya-wechat-service/model"
)

var UserServiceSingleton *UserService

func UserServiceInstance() *UserService {
	if UserServiceSingleton == nil {
		UserServiceSingleton = new(UserService)
		UserServiceSingleton.userDao = dao.UserDaoSingleton
		UserServiceSingleton.addressDao = dao.UserAddressDaoInstance()
	}
	return UserServiceSingleton
}

type UserService struct {
	userDao    *dao.UserDao
	addressDao *dao.UserAddressDao
}

func (s *UserService) GetList() ([]*dto.UserDTO, error) {
	users, err := s.userDao.SelectAll()
	if err != nil {
		return nil, err
	}
	return buildUserDTOs(users), nil
}

func (s *UserService) DefaultAddress(userID int64) (dto.UserAddressDTO, error) {
	var uad = dto.UserAddressDTO{}
	ads, err := s.addressDao.SelectByUserID(userID)
	if err != nil {
		return uad, err
	}
	for _, ad := range ads {
		if ad.IsDefault {
			uad = installUserAddressDTO(ad)
			break
		}
	}
	return uad, nil
}

func (s *UserService) AddressList(userID int64) ([]dto.UserAddressDTO, error) {
	ads, err := s.addressDao.SelectByUserID(userID)
	if err != nil {
		return nil, err
	}
	dtos := make([]dto.UserAddressDTO, len(ads))
	for i, ad := range ads {
		dtos[i] = installUserAddressDTO(ad)
	}
	return dtos, nil
}

func (s *UserService) GetAddressByID(ID int64) (dto dto.UserAddressDTO, err error) {
	uad, err := s.addressDao.SelectByID(ID)
	if err != nil {
		return
	}
	return installUserAddressDTO(uad), nil
}

func (s *UserService) CreateAddress(req request.UserAddressAddRequest) (id int64, err error) {
	setMap := map[string]interface{}{
		"name":        req.Name,
		"mobile":      req.Mobile,
		"province_id": req.ProvinceID,
		"city_id":     req.CityID,
		"district_id": req.DistrictID,
		"address":     req.Address,
		"is_default":  req.IsDefault,
		"user_id":     req.UserID,
	}
	return s.addressDao.Create(setMap)
}

func (s *UserService) Login(openID string, nickName string, avatarURL string) (*model.User, error) {
	user, err := s.userDao.SelectByWechatID(openID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}
	setmap := map[string]interface{}{
		"wechat_id":  openID,
		"nick_name":  nickName,
		"avatar_url": avatarURL,
	}
	id, err := s.userDao.CreateUser(setmap)
	if err != nil {
		return nil, err
	}
	user = &model.User{
		ID:        id,
		NickName:  nickName,
		AvatarURL: avatarURL,
	}
	return user, nil
}

func installUserAddressDTO(ad model.UserAddress) dto.UserAddressDTO {
	dto := dto.UserAddressDTO{}
	dto.ID = ad.ID
	dto.Name = ad.Name
	dto.IsDefault = ad.IsDefault
	dto.Mobile = ad.Mobile
	dto.Address = ad.Address
	dto.ProvinceID = ad.ProvinceID
	dto.CityID = ad.CityID
	dto.DistrictID = ad.DistricID
	return dto
}

func installUserDTO(model *model.User) *dto.UserDTO {
	userDto := new(dto.UserDTO)
	userDto.Name = model.Name.String
	userDto.PhoneNo = model.PhoneNo.String
	userDto.Email = model.Email.String
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
