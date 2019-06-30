package dao

import (
	"gotrue/model"
	"gotrue/sys"

	sq "github.com/Masterminds/squirrel"
)

var UserAddressDaoSingleton *UserAddressDao

func initUserAddressDao() {
	if UserAddressDaoSingleton == nil {
		UserAddressDaoSingleton = new(UserAddressDao)
		UserAddressDaoSingleton.db = sys.GetEasyDB()
	}
}

func UserAddressDaoInstance() *UserAddressDao {
	initUserAddressDao()
	return UserAddressDaoSingleton
}

var columns_user_address_all = []string{"id", "name", "user_id", "country_id", "province_id", "city_id", "district_id", "address", "mobile", "is_default"}

type UserAddressDao struct {
	db *sys.EasyDB
}

func (dao *UserAddressDao) SelectByUserID(userID int64) ([]model.UserAddress, error) {
	usrAds := []model.UserAddress{}
	err := dao.db.SelectDSL(&usrAds, columns_user_address_all, model.Table_User_Address, sq.Eq{"user_id": userID})
	return usrAds, err
}

func (dao *UserAddressDao) SelectByID(ID int64) (model.UserAddress, error) {
	uad := model.UserAddress{}
	err := dao.db.SelectOneDSL(&uad, columns_user_address_all, model.Table_User_Address, sq.Eq{"id": ID})
	return uad, err
}

func (dao *UserAddressDao) Create(setMap map[string]interface{}) (int64, error) {
	_, id, err := dao.db.Insert(model.Table_User_Address, setMap)
	return id, err
}
