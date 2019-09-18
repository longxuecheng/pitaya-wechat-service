package model

import (
	"database/sql"
	"gotrue/facility/utils"

	"github.com/dgrijalva/jwt-go"
)

type UserType int

const (
	UserTypeCustomer UserType = iota
	UserTypeAdmin
	UserTypePartener
)

func (ut UserType) Int() int {
	return int(ut)
}

type User struct {
	ID        int64          `db:"id"`
	UserType  UserType       `db:"user_type"`
	Name      sql.NullString `db:"name"`
	PhoneNo   sql.NullString `db:"phone_no"`
	Email     sql.NullString `db:"email"`
	NickName  string         `db:"nick_name"`
	AvatarURL string         `db:"avatar_url"`
	WechatID  string         `db:"wechat_id"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) Columns() []string {
	return utils.TagValues(u, "db")
}

type UserClaims struct {
	UserID int64 `json:"userId"`
	jwt.StandardClaims
}
