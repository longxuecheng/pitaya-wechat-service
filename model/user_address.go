package model

var Table_User_Address string = "user_address"

type UserAddress struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	UserID     int64  `db:"user_id"`
	CountryID  int    `db:"country_id"`
	ProvinceID int    `db:"province_id"`
	CityID     int    `db:"city_id"`
	DistricID  int    `db:"district_id"`
	Address    string `db:"address"`
	Mobile     string `db:"mobile"`
	IsDefault  bool   `db:"is_default"`
}

func (ua *UserAddress) RegionIDs() []int {
	return []int{ua.ProvinceID, ua.CityID, ua.DistricID}
}
