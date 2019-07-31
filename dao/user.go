package dao

import (
	"database/sql"
	"gotrue/model"
	"gotrue/sys"

	"github.com/Masterminds/squirrel"
)

var UserDaoSingleton *UserDao

func initUserDao() {
	UserDaoSingleton = &UserDao{
		db: sys.GetEasyDB(),
	}
}

type UserDao struct {
	db *sys.EasyDB
}

var columns_user = []string{"id", "name", "phone_no", "email", "wechat_id", "avatar_url", "nick_name", "user_type"}

func (dao *UserDao) SelectAll() ([]*model.User, error) {
	users := []*model.User{}
	err := dao.db.SelectDSL(&users, columns_user, model.Table_User, nil)
	if err != nil {
		return nil, err
	}
	return users, err
}

func (dao *UserDao) SelectByWechatID(wechatID string) (*model.User, error) {
	users := new(model.User)
	err := dao.db.SelectOneDSL(users, columns_user, model.Table_User, squirrel.Eq{"wechat_id": wechatID})
	if err != nil {
		if sql.ErrNoRows == err {
			return nil, nil
		}
		return nil, err
	}
	return users, err
}

func (dao *UserDao) SelectByID(userID int64) (*model.User, error) {
	users := new(model.User)
	err := dao.db.SelectOneDSL(users, columns_user, model.Table_User, squirrel.Eq{"id": userID})
	if err != nil {
		if sql.ErrNoRows == err {
			return nil, nil
		}
		return nil, err
	}
	return users, err
}

func (dao *UserDao) CreateUser(setMap map[string]interface{}) (int64, error) {
	_, id, err := dao.db.Insert(model.Table_User, setMap)
	return id, err
}
