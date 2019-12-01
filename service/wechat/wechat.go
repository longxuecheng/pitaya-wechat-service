package wechat

import (
	"encoding/json"
	"fmt"
	"gotrue/facility/http_util"
	"gotrue/service/wechat/payment"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	defaultWechatService *wechatService
)

func InitWechatService() {
	defaultWechatService = &wechatService{
		NewTokenManager(true),
	}
}

// WechatService returns default wechat service instance
func WechatService() *wechatService {
	return defaultWechatService
}

type wechatService struct {
	*TokenManager
}

func (s *wechatService) UserInfo(code string) (Code2SessionResponse, error) {
	url := fmt.Sprintf(authorize_url, appID, secret, code)
	response := Code2SessionResponse{}
	err := http_util.Get(url, &response)
	return response, err
}

func (s *wechatService) AccessToken() string {
	return s.TokenManager.AccessToken()
}

func (s *wechatService) SendTemplateMessage(req *TemplateMsgRequest) error {
	bytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	response, err := http_util.Send(http.MethodPost, sendmessage_url, strings.NewReader(string(bytes)), func(r *http.Request) {
		r.Header.Set("Content-Type", "application/json")
		r.Form.Add("access_token", s.AccessToken())
		r.URL.RawQuery = r.Form.Encode()
	})
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(response.Body)
	result := &NotifyResponse{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}
	if ok, err := result.IsOK(); !ok {
		return err
	}
	return nil
}

func (s *wechatService) SendUniformMessage(req *UniformMsgRequest) error {
	bytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	response, err := http_util.Send(http.MethodPost, uniformmessage_url, strings.NewReader(string(bytes)), func(r *http.Request) {
		r.Header.Set("Content-Type", "application/json")
		r.Form.Add("access_token", s.AccessToken())
		r.URL.RawQuery = r.Form.Encode()
	})
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(response.Body)
	result := &NotifyResponse{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}
	if ok, err := result.IsOK(); !ok {
		return err
	}
	return nil
}

// Pay 调用微信接口进行预支付
func (s *wechatService) PrePay(r *PrepayRequest) (*payment.PrepayReponse, error) {

	app := payment.AppBasic{
		AppID:      appID,
		MerchantID: merchantID,
	}
	sign := payment.SignBasic{
		SignType: payment.SignTypeMD5,
		NonceStr: strconv.FormatInt(time.Now().Unix(), 10),
	}
	req := payment.UnifiedOrderRequest{
		AppBasic:   app,
		SignBasic:  sign,
		OutTradeNo: r.OrderNo,
		NotifyURL:  "https://www.geluxiya.com",
		ServerIP:   "192.168.1.123",
		TradeType:  payment.TradeTypeJSAPI,
		TotalFee:   r.TotalFee,
		Desc:       r.Desc,
		OpenID:     r.OpenID,
	}
	fmt.Println(req.OutTradeNo)
	// req.OpenID = "ovxEC5YTWQk6Vv5FJdN_30gkBr-g"
	signedStr := req.SignParams()
	req.Sign = signedStr
	xmlText, err := req.ToXml()
	if err != nil {
		return nil, err
	}
	result := &payment.UnifiedOrderResponse{}
	err = http_util.PostXml(result, unifiedorder_url, strings.NewReader(xmlText))
	if err != nil {
		return nil, err
	}
	ok, err := result.IsOK()
	if !ok {
		return nil, err
	}
	preypayResponse := &payment.PrepayReponse{
		AppID:     appID,
		NonceStr:  strconv.FormatInt(time.Now().Unix(), 10),
		SignType:  payment.SignTypeMD5,
		TimeStamp: strconv.FormatInt(time.Now().Unix(), 10),
		Package:   "prepay_id=" + result.PrepayID,
	}
	preypayResponse.SignParams()
	preypayResponse.PrepayID = result.PrepayID
	return preypayResponse, nil
}

func (s *wechatService) QueryPayResult(orderNo string) (*payment.QueryOrderResponse, error) {
	app := payment.AppBasic{
		AppID:      appID,
		MerchantID: merchantID,
	}
	sign := payment.SignBasic{
		SignType: payment.SignTypeMD5,
		NonceStr: strconv.FormatInt(time.Now().Unix(), 10),
	}
	req := payment.OrderQueryRequest{
		AppBasic:   app,
		SignBasic:  sign,
		OutTradeNo: orderNo,
	}
	xml, err := req.ToXml()
	if err != nil {
		return nil, err
	}
	result := &payment.QueryOrderResponse{}
	err = http_util.PostXml(result, payment.OrderQueryURL, strings.NewReader(xml))
	if err != nil {
		return nil, err
	}
	ok, err := result.IsOK()
	if !ok {
		return nil, err
	}
	return result, nil
}
