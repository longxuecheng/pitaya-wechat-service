package model

import (
	"database/sql"
	"gotrue/facility/utils"
	"time"

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
	ID              int64          `db:"id" pk:"true"`
	UserType        UserType       `db:"user_type"`
	Name            sql.NullString `db:"name"`
	PhoneNo         sql.NullString `db:"phone_no"`
	Email           sql.NullString `db:"email"`
	NickName        string         `db:"nick_name"`
	AvatarURL       string         `db:"avatar_url"`
	WechatID        string         `db:"wechat_id"`
	ChannelUserID   int64          `db:"channel_user_id"`
	BindChannelTime NullUTC8Time   `db:"bind_channel_time"`
	ChannelCode     string         `db:"channel_code"`
	CreateTime      time.Time      `db:"create_time"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) Columns() []string {
	return utils.TagValues(u, "db")
}

func (u *User) HasMobile() bool {
	return u.PhoneNo.Valid
}

type UserCollection struct {
	Items []*User
}

func (c *UserCollection) ToMap() map[int64]*User {
	userMap := map[int64]*User{}
	for _, item := range c.Items {
		userMap[item.ID] = item
	}
	return userMap
}

func (c *UserCollection) UserMobiles() []string {
	mobiles := []string{}
	for _, item := range c.Items {
		if item.PhoneNo.String != "" {
			mobiles = append(mobiles, item.PhoneNo.String)
		}
	}
	return mobiles
}

type UserClaims struct {
	UserID int64 `json:"userId"`
	jwt.StandardClaims
}
