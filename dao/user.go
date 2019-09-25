package dao

import (
	"database/sql"
	"gotrue/model"
	"gotrue/sys"

	"github.com/Masterminds/squirrel"
)

var UserDaoSingleton *UserDao

func initUserDao() {
	u := &model.User{}
	UserDaoSingleton = &UserDao{
		table:   u.TableName(),
		columns: u.Columns(),
		db:      sys.GetEasyDB(),
	}
}

type UserDao struct {
	table   string
	columns []string
	db      *sys.EasyDB
}

func (dao *UserDao) SelectAll() ([]*model.User, error) {
	users := []*model.User{}
	err := dao.db.SelectDSL(&users, dao.columns, dao.table, nil)
	if err != nil {
		return nil, err
	}
	return users, err
}

func (dao *UserDao) SelectByWechatID(wechatID string) (*model.User, error) {
	users := new(model.User)
	err := dao.db.SelectOneDSL(users, dao.columns, dao.table, squirrel.Eq{"wechat_id": wechatID})
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
	err := dao.db.SelectOneDSL(users, dao.columns, dao.table, squirrel.Eq{"id": userID})
	if err != nil {
		if sql.ErrNoRows == err {
			return nil, nil
		}
		return nil, err
	}
	return users, err
}

func (dao *UserDao) SelectByIDs(userIDs []int64) (*model.UserCollection, error) {
	users := []*model.User{}
	err := dao.db.SelectDSL(&users, dao.columns, dao.table, squirrel.Eq{"id": userIDs})
	if err != nil {
		return nil, err
	}
	return &model.UserCollection{
		Items: users,
	}, err
}

func (dao *UserDao) CreateUser(setMap map[string]interface{}) (int64, error) {
	_, id, err := dao.db.Insert(dao.table, setMap)
	return id, err
}
