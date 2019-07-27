package wechat

import (
	"fmt"
	"gotrue/dto/request"
	"testing"
)

func TestPay(t *testing.T) {
	openIDLXC := "ovxEC5YTWQk6Vv5FJdN_30gkBr-g"
	s, err := WechatService().Pay(openIDLXC)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("Result xml text is %s \n", s)
}

func TestQueryPayResult(t *testing.T) {
	result, err := WechatService().QueryPayResult("orderNo1563775302")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf(" result is %+v\n", result)
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
