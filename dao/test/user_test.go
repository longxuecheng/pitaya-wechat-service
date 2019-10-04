package test

import (
	"fmt"
	"gotrue/dao"
	"testing"
)

func TestCreateUser(t *testing.T) {
	dao.Init()
	setMap := map[string]interface{}{
		"wechat_id": "wechat id",
		"nick_name": "\xF0\x9F\x90\x91",
	}
	_, err := dao.UserDaoSingleton.CreateUser(setMap)
	if err != nil {
		t.Error(err)
	}
}

func TestNickName(t *testing.T) {
	fmt.Println("\xF0\x9F\x90\x91")
}
