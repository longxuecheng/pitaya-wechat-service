package test

import (
	"fmt"
	"gotrue/dao"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	dao.Init()
	os.Exit(m.Run())
}

func TestCreateUser(t *testing.T) {
	setMap := map[string]interface{}{
		"wechat_id": "wechat id",
		"nick_name": "我爱国家\xF0\x9F\x90\x91",
	}
	_, err := dao.UserDaoSingleton.CreateUser(setMap)
	if err != nil {
		t.Error(err)
	}
}

func TestNickName(t *testing.T) {
	fmt.Println("\xF0\x9F\x90\x91")
}

func TestDBWriteLock(t *testing.T) {
	user, err := dao.UserDaoSingleton.SelectByID(0)
	if err != nil {
		t.Error(err)
		return
	}
	if user != nil {
		fmt.Printf("%+v\n", user)
	}
}
