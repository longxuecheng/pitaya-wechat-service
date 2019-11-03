package user

import (
	"database/sql"
	"gotrue/dao"
	"gotrue/dto/request"
	"gotrue/dto/response"
	"gotrue/facility/errors"
	"gotrue/model"
	"gotrue/service/api"
	"gotrue/service/region"
	"gotrue/sys"
	"log"
	"time"

	"github.com/google/uuid"
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

func (a *address) userAddressDTO(fullRegion string) *response.UserAddress {
	dto := installUserAddress(a.data)
	dto.FullRegion = fullRegion
	return dto
}

func (s *User) GetList() ([]*response.UserDTO, error) {
	users, err := s.userDao.SelectAll()
	if err != nil {
		return nil, err
	}
	return buildUserDTOs(users), nil
}

func (s *User) DefaultAddress(userID int64) (*response.UserAddress, error) {
	var address = &response.UserAddress{}
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

func (s *User) AddressList(userID int64) ([]*response.UserAddress, error) {
	ads, err := s.addressDao.SelectByUserID(userID)
	if err != nil {
		return nil, err
	}
	dtos := make([]*response.UserAddress, len(ads))
	for i, ad := range ads {
		userAddress := installUserAddress(ad)
		fullRegion, err := s.regionService.FullName(ad.RegionIDs())
		if err != nil {
			return nil, err
		}
		userAddress.FullRegion = fullRegion
		dtos[i] = userAddress
	}
	return dtos, nil
}

func (s *User) BindChannelPerson(userID int64, channelCode string) error {
	channelPerson, err := s.userDao.SelectByChannelCode(channelCode)
	if err != nil {
		return errors.NewWithCodef("ChannelCodeInvalid", "渠道码不合法")
	}
	if userID == channelPerson.ID {
		return errors.NewWithCodef("ChannelPersonInvalid", "渠道人不合法")
	}
	user, err := s.userDao.SelectByID(userID)
	if err != nil {
		return err
	}
	if user.ChannelUserID > 0 {
		return errors.NewWithCodef("ChannelPersonUnchangeable", "渠道人不可变更")
	}
	user.ChannelUserID = channelPerson.ID
	user.BindChannelTime = model.NullUTC8Time{
		Time:  time.Now(),
		Valid: true,
	}
	return s.userDao.UpdateByID(user)
}

func (s *User) GetAddressByID(ID int64) (dto *response.UserAddress, err error) {
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

func (s *User) GetUserByID(userID int64) (dto *response.UserDTO, err error) {
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
	log.Printf("openID %s nickName %s", openID, nickName)
	user, err := s.userDao.SelectByWechatID(openID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		if user.NickName == "" {
			user.NickName = nickName
			user.AvatarURL = avatarURL
			return user, s.userDao.UpdateByID(user)
		}
		return user, nil
	}
	uuid := uuid.New()
	setmap := map[string]interface{}{
		"wechat_id":    openID,
		"nick_name":    nickName,
		"avatar_url":   avatarURL,
		"channel_code": uuid.String(),
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

func installUserAddress(ad *model.UserAddress) *response.UserAddress {
	dto := &response.UserAddress{}
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

func installUserDTO(model *model.User) *response.UserDTO {
	userDto := new(response.UserDTO)
	userDto.Name = model.Name.String
	userDto.PhoneNo = model.PhoneNo.String
	userDto.Email = model.Email.String
	userDto.OpenID = model.WechatID
	return userDto
}

func buildUserDTOs(models []*model.User) []*response.UserDTO {
	if models == nil || len(models) == 0 {
		return nil
	}
	dtos := make([]*response.UserDTO, len(models))
	for i, model := range models {
		dtos[i] = installUserDTO(model)
	}
	return dtos
}
