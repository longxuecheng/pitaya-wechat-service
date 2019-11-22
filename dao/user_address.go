package dao

import (
	"database/sql"
	"gotrue/model"
	

	sq "github.com/Masterminds/squirrel"
)

var UserAddressDao *UserAddress

func initUserAddressDao() {
	UserAddressDao = &UserAddress{
		db: GetEasyDB(),
	}
}

var columns_user_address_all = []string{"id", "name", "user_id", "country_id", "province_id", "city_id", "district_id", "address", "mobile", "is_default"}

type UserAddress struct {
	db *EasyDB
}

func (dao *UserAddress) SelectByUserID(userID int64) ([]*model.UserAddress, error) {
	usrAds := []*model.UserAddress{}
	err := dao.db.SelectDSL(&usrAds, columns_user_address_all, model.Table_User_Address, sq.Eq{"user_id": userID})
	return usrAds, err
}

func (dao *UserAddress) SelectByID(ID int64) (*model.UserAddress, error) {
	a := &model.UserAddress{}
	err := dao.db.SelectOneDSL(a, columns_user_address_all, model.Table_User_Address, sq.Eq{"id": ID})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return a, err
}

func (dao *UserAddress) Create(tx *sql.Tx, setMap map[string]interface{}) (int64, error) {
	_, id, err := dao.db.Insert(model.Table_User_Address, setMap, tx)
	return id, err
}

func (dao *UserAddress) UpdateByUserID(tx *sql.Tx, userID int64, setMap map[string]interface{}) error {
	_, err := dao.db.UpdateTx(tx, model.Table_User_Address, setMap, sq.Eq{"user_id": userID})
	return err
}

func (dao *UserAddress) UpdateByID(tx *sql.Tx, id int64, setMap map[string]interface{}) error {
	_, err := dao.db.UpdateTx(tx, model.Table_User_Address, setMap, sq.Eq{"id": id})
	return err
}
