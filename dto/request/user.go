package request

type UserAddressAddRequest struct {
	ID         int    `json:"id"`
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

type WechatUser struct {
	AvatarURL string `json:"avatarUrl"`
	NickName  string `json:"nickName"`
}
