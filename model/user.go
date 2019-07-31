package model

import (
	"database/sql"

	"github.com/dgrijalva/jwt-go"
)

const Table_User = "user"

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

type UserClaims struct {
	UserID int64 `json:"userId"`
	jwt.StandardClaims
}
