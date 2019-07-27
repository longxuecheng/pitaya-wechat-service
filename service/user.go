package service

import (
	"database/sql"
	"gotrue/dao"
	"gotrue/dto"
	"gotrue/dto/request"
	"gotrue/model"
	"gotrue/sys"
	"strings"
)

var UserServiceSingleton *UserService

func UserServiceInstance() *UserService {
	if UserServiceSingleton == nil {
		UserServiceSingleton = &UserService{
			userDao:    dao.UserDaoSingleton,
			addressDao: dao.UserAddressDaoSingleton,
			regionDao:  dao.RegionDaoInstance(),
		}
	}
	return UserServiceSingleton
}

type UserService struct {
	userDao    *dao.UserDao
	addressDao *dao.UserAddressDao
	regionDao  *dao.RegionDao
}

type address struct {
	data *model.UserAddress
}

func newAddress(data *model.UserAddress) *address {
	return &address{
		data,
	}
}

func (a *address) regionIDs() []int {
	ids := []int{}
	ids = append(ids, a.data.ProvinceID)
	ids = append(ids, a.data.CityID)
	ids = append(ids, a.data.DistricID)
	return ids
}

func (a *address) userAddressDTO(regions []*model.Region) *dto.UserAddress {
	regionNames := []string{}
	for _, r := range regions {
		regionNames = append(regionNames, r.Name)
	}
	dto := installUserAddress(a.data)
	dto.FullRegion = strings.Join(regionNames, "-")
	return dto
}

func (s *UserService) GetList() ([]*dto.UserDTO, error) {
	users, err := s.userDao.SelectAll()
	if err != nil {
		return nil, err
	}
	return buildUserDTOs(users), nil
}

func (s *UserService) DefaultAddress(userID int64) (*dto.UserAddress, error) {
	var address = &dto.UserAddress{}
	ads, err := s.addressDao.SelectByUserID(userID)
	if err != nil {
		return address, err
	}
	for _, ad := range ads {
		if ad.IsDefault {
			address = installUserAddress(ad)
			break
		}
	}
	return address, nil
}

func (s *UserService) DeleteAddressByID(id int64) error {
	return nil
}

func (s *UserService) AddressList(userID int64) ([]*dto.UserAddress, error) {
	ads, err := s.addressDao.SelectByUserID(userID)
	if err != nil {
		return nil, err
	}
	dtos := make([]*dto.UserAddress, len(ads))
	for i, ad := range ads {
		dtos[i] = installUserAddress(ad)
	}
	return dtos, nil
}

func (s *UserService) GetAddressByID(ID int64) (dto *dto.UserAddress, err error) {
	a, err := s.addressDao.SelectByID(ID)
	if err != nil {
		return
	}
	address := newAddress(a)
	regions, err := s.regionDao.SelectByIDs(address.regionIDs())
	if err != nil {
		return
	}
	return address.userAddressDTO(regions), nil
}

func (s *UserService) GetUserByID(userID int64) (dto *dto.UserDTO, err error) {
	user, err := s.userDao.SelectByID(userID)
	if err != nil {
		return
	}
	return installUserDTO(user), nil
}

// CreateAddress create or update an user address
func (s *UserService) CreateAddress(userID int64, req request.UserAddressAddRequest) (id int64, err error) {
	a, err := s.addressDao.SelectByID(req.ID)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	setMap := map[string]interface{}{
		"name":        req.Name,
		"mobile":      req.Mobile,
		"province_id": req.ProvinceID,
		"city_id":     req.CityID,
		"district_id": req.DistrictID,
		"address":     req.Address,
		"is_default":  req.IsDefault,
		"user_id":     userID,
	}
	sys.GetEasyDB().ExecTx(func(tx *sql.Tx) error {
		if req.IsDefault {
			updateMap := map[string]interface{}{
				"is_default": false,
			}
			err = s.addressDao.UpdateByUserID(tx, userID, updateMap)
			if err != nil {
				return err
			}
		}
		if a != nil {
			// update address
			err = s.addressDao.UpdateByID(tx, req.ID, setMap)
			if err != nil {
				return err
			}
		} else {
			id, err = s.addressDao.Create(tx, setMap)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return
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

func installUserAddress(ad *model.UserAddress) *dto.UserAddress {
	dto := &dto.UserAddress{}
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
	userDto.OpenID = model.WechatID
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
