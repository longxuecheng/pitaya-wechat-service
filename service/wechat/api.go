package wechat

import (
	"fmt"
	"gotrue/facility/errors"
	"gotrue/service/wechat/payment"
)

const (
	appID            string = "wx91eb35616f475761"
	secret           string = "41f8cb23aab151c6aa8ca64a0cfb596c"
	merchantID       string = "1540271371"
	authorize_url    string = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	accessToken_url  string = "https://api.weixin.qq.com/cgi-bin/token?grant_type=%s&appid=%s&secret=%s"
	unifiedorder_url string = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	sendmessage_url  string = "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send"
)

type IWechat interface {
	AccessToken() error
	Pay(*PrepayRequest) (*payment.PrepayReponse, error)
	QueryPayResult(orderNo string) (*payment.QueryOrderResponse, error)
	SendNotification(*NotifyRequest) error
}

type NotifyRequest struct {
	// AccessToken string                 `access_token`
	ToUser     string                 `json:"touser"`
	TemplateID string                 `json:"template_id"`
	Page       string                 `json:"page"`
	FormID     string                 `json:"form_id"`
	Data       map[string]interface{} `json:"data"`
}

type NotifyResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (nr *NotifyResponse) IsOK() (bool, error) {
	if nr.ErrCode == 0 {
		return true, nil
	}
	return false, errors.NewWithCodef(fmt.Sprintf("WechatNotifyError_%d", nr.ErrCode), nr.ErrMsg)
}

type PrepayRequest struct {
	OpenID   string
	OrderNo  string
	TotalFee int64
	Desc     string
}

type BasicResponse struct {
	ErrorCode int32  `json:"errcode"`
	ErrorMsg  string `json:"errmsg"`
}

func (b *BasicResponse) isOk() bool {
	return b.ErrorCode == 0
}

type Code2SessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	BasicResponse
}

type AccessTokenResonse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	BasicResponse
}
