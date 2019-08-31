package wechat

import (
	"encoding/json"
	"fmt"
	"gotrue/facility/http_util"
	"gotrue/service/wechat/payment"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mileusna/crontab"
)

var (
	logger               = log.New(os.Stderr, "[Wechat-service] ", log.LstdFlags)
	defaultWechatService = &wechatService{
		newTokenManager(true),
	}
)

// WechatService returns default wechat service instance
func WechatService() *wechatService {
	return defaultWechatService
}

type wechatService struct {
	*wechatTokenManager
}

func (s *wechatService) UserInfo(code string) (Code2SessionResponse, error) {
	url := fmt.Sprintf(authorize_url, appID, secret, code)
	response := Code2SessionResponse{}
	err := http_util.Get(url, &response)
	return response, err
}

func (s *wechatService) AccessToken() string {
	return s.accessToken()
}

func (s *wechatService) SendNotification(req *NotifyRequest) error {
	bytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	response, err := http_util.Send(http.MethodPost, sendmessage_url, strings.NewReader(string(bytes)), func(r *http.Request) error {
		r.Header.Set("Content-Type", "application/json")
		r.Form.Add("access_token", s.accessToken())
		r.URL.RawQuery = r.Form.Encode()
		return nil
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

type wechatTokenManager struct {
	at      string
	atExpIn int64
	crontab *crontab.Crontab
}

func newTokenManager(startSchedule bool) *wechatTokenManager {
	m := &wechatTokenManager{}
	m.crontab = crontab.New()
	if startSchedule {
		m.scheduleTasks()
	}
	return m
}

// crontab syntax https://github.com/mileusna/crontab
func (m *wechatTokenManager) scheduleTasks() {
	m.crontab.MustAddJob("*/10 * * * *", m.refreshAccessToken)
	// run imediately when start
	m.crontab.RunAll()
}

func (m *wechatTokenManager) accessToken() string {
	return m.at
}

func (m *wechatTokenManager) refreshAccessToken() {
	act := AccessTokenResonse{}
	url := fmt.Sprintf(accessToken_url, "client_credential", appID, secret)
	err := http_util.DoGet(&act, url, nil)
	if err != nil {
		logger.Printf("access token refresh error %+v\n", err)
	}
	if act.isOk() {
		m.at = act.AccessToken
		m.atExpIn = act.ExpiresIn
	}
	logger.Printf("Access token : %s\n", act.AccessToken)
}
