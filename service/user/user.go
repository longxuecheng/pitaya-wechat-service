package user

import (
	"database/sql"
	"gotrue/service/api"
	"gotrue/dao"
	"gotrue/dto"
	"gotrue/dto/request"
	"gotrue/model"
	"gotrue/service/region"
	"gotrue/sys"
)

var UserService *User

func beforeInit() {
	region.Init()
}

func initUserService() {
	if UserService != nil {
		return
	}
	beforeInit()
	UserService = &User{
		userDao:       dao.UserDaoSingleton,
		addressDao:    dao.UserAddressDao,
		regionService: region.RegionService,
	}
}

type User struct {
	userDao       *dao.UserDao
	addressDao    *dao.UserAddress
	regionService api.IRegionService
}

type address struct {
	data *model.UserAddress
}

func newAddress(data *model.UserAddress) *address {
	return &address{
		data,
	}
}

func (a *address) userAddressDTO(fullRegion string) *dto.UserAddress {
	dto := installUserAddress(a.data)
	dto.FullRegion = fullRegion
	return dto
}

func (s *User) GetList() ([]*dto.UserDTO, error) {
	users, err := s.userDao.SelectAll()
	if err != nil {
		return nil, err
	}
	return buildUserDTOs(users), nil
}

func (s *User) DefaultAddress(userID int64) (*dto.UserAddress, error) {
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

func (s *User) DeleteAddressByID(id int64) error {
	return nil
}

func (s *User) AddressList(userID int64) ([]*dto.UserAddress, error) {
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

func (s *User) GetAddressByID(ID int64) (dto *dto.UserAddress, err error) {
	a, err := s.addressDao.SelectByID(ID)
	if err != nil {
		return
	}
	address := newAddress(a)
	fullRegion, err := s.regionService.FullName(a.RegionIDs())
	if err != nil {
		return nil, err
	}
	return address.userAddressDTO(fullRegion), nil
}

func (s *User) GetUserByID(userID int64) (dto *dto.UserDTO, err error) {
	user, err := s.userDao.SelectByID(userID)
	if err != nil {
		return
	}
	return installUserDTO(user), nil
}

// CreateAddress create or update an user address
func (s *User) CreateAddress(userID int64, req request.UserAddressAddRequest) (id int64, err error) {
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

func (s *User) Login(openID string, nickName string, avatarURL string) (*model.User, error) {
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
