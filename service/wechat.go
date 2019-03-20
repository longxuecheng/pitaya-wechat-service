package service

import (
	"fmt"
	"pitaya-wechat-service/facility/http_util"
	wechat_model "pitaya-wechat-service/wechat/model"
)

const (
	appID         string = "wx48b125b125deff5e"
	secret        string = "49e671ecef555f8fc09062ada5a6b829"
	authorize_url string = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

func GetWechatUserInfo(code string) (wechat_model.Code2SessionResponse, error) {
	url := fmt.Sprintf(authorize_url, appID, secret, code)
	response := wechat_model.Code2SessionResponse{}
	err := http_util.Get(url, &response)
	return response, err
}
