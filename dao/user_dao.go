package dao

import (
	"pitaya-wechat-service/model"
	"pitaya-wechat-service/sys"

	"github.com/jmoiron/sqlx"
)

var UserDaoSingleton *UserDao

func init() {
	UserDaoSingleton = new(UserDao)
	UserDaoSingleton.dbDriver = sys.DBConnection()
}

type UserDao struct {
	dbDriver *sqlx.DB
}

func (dao *UserDao) SelectAll() ([]*model.User, error) {
	users := []*model.User{}
	err := dao.dbDriver.Select(&users, "SELECT id,name,phone_no,email FROM user ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	return users, err
}

