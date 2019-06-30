package dto

type UserDTO struct {
	Name    string `json:"name"`
	PhoneNo string `json:"phoneNo"`
	Email   string `json:"email"`
	OpenID  string `json:"-"`
}

type UserAddress struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	FullRegion string `json:"full_region"`
	Address    string `json:"address"`
	Mobile     string `json:"mobile"`
	IsDefault  bool   `json:"is_default"`
	ProvinceID int
	CityID     int
	DistrictID int
}
