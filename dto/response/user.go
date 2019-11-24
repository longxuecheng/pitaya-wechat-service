package response

type User struct {
	AvatarURL string `json:"avatarUrl"`
	NickName  string `json:"nickName"`
	UserType  int    `json:"userType"`
	UserID    string `json:"userId"`
}

type InternalUser struct {
	ID        int64  `json:"id"`
	AvatarURL string `json:"avatarUrl"`
	NickName  string `json:"nickname"`
}

type UserDTO struct {
	Name      string `json:"name"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarUrl"`
	PhoneNo   string `json:"phoneNo"`
	Email     string `json:"email"`
	OpenID    string `json:"-"`
}

type UserAddress struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	FullRegion string `json:"full_region"`
	Address    string `json:"address"`
	Mobile     string `json:"mobile"`
	IsDefault  bool   `json:"is_default"`
	ProvinceID int
	CityID     int
	DistrictID int
}
