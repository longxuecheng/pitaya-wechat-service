package sms

import (
	"fmt"
	"testing"
)

func TestSendNotificationMsg(t *testing.T) {
	err := SendPayNotificationMsg(&MultiSendRequest{
		Mobiles: []string{"18911792314", "15201389366"},
		Params:  []string{"果真管理员", "test3940401."},
	})
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}
