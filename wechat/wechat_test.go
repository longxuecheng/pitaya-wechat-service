package wechat

import (
	"fmt"
	"gotrue/dto/request"
	"testing"
)

func TestPay(t *testing.T) {
	s, err := WechatService().Pay()
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("Result xml text is %s \n", s)
}

func TestOpenID(t *testing.T) {
	req := request.WechatLogin{
		Code: "071YwLU32wEmVP0l4lW32wdHU32YwLUh",
	}
	resp, err := WechatService().UserInfo(req.Code)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("open id is  %s \n", resp.OpenID)
}
