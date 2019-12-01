package wechat

import (
	"fmt"
	"gotrue/dto/request"
	"testing"
	"time"
)

// func TestPay(t *testing.T) {
// 	openIDLXC := "ovxEC5YTWQk6Vv5FJdN_30gkBr-g"
// 	req := &PrepayRequest{
// 		OpenID:   openIDLXC,
// 		OrderNo:  "order-xckasdklfdlsa1",
// 		TotalFee: 10,
// 		Desc:     "测试支付通知",
// 	}
// 	s, err := WechatService().PrePay(req)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	notificationReq := &NotifyRequest{
// 		ToUser:     openIDLXC,
// 		TemplateID: "F56_89H1A2SiyEmnwUSGNw_kyTIcdFLBELFaU2sFUhU",
// 		FormID:     s.PrepayID,
// 		Data: map[string]interface{}{
// 			"keyword1": map[string]string{
// 				"value": "北京平谷桃源",
// 			},
// 			"keyword2": map[string]string{
// 				"value": "大桃",
// 			},
// 		},
// 	}

// 	err = WechatService().SendNotification(notificationReq)

// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

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

func TestTokenManager(t *testing.T) {
	tm := NewTokenManager(true)
	tm.ScheduleTasks()
	time.Sleep(time.Hour)
}
