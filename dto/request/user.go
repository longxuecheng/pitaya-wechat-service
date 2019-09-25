package request

import "gotrue/facility/errors"

type UserAddressAddRequest struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	ProvinceID int    `json:"province_id"`
	CityID     int    `json:"city_id"`
	DistrictID int    `json:"district_id"`
	Address    string `json:"address"`
	Mobile     string `json:"mobile"`
	IsDefault  bool   `json:"is_default"`
	UserID     int64  `json:"-"`
}

type WechatLogin struct {
	Code       string `json:"code"`
	WechatUser `json:"userInfo"`
}

func (u *WechatLogin) Validate() error {
	if u.Code == "" {
		return errors.NewWithCodef("InvalidCode", "invalid source")
	}
	if u.AvatarURL == "" && u.NickName == "" {
		return errors.NewWithCodef("InvalidUser", "user invalid")
	}
	return nil
}

type WechatUser struct {
	AvatarURL string `json:"avatarUrl"`
	NickName  string `json:"nickName"`
}
