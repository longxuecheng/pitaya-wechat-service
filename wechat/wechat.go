package wechat

import (
	"fmt"
	"gotrue/facility/http_util"
	"gotrue/wechat/payment"
	"log"
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

func (s *wechatService) Pay(openID string) (*payment.PrepayReponse, error) {
	req := payment.UnifiedOrderRequest{}
	app := payment.AppBasic{
		AppID:      appID,
		MerchantID: merchantID,
	}
	sign := payment.SignBasic{
		SignType: payment.SignTypeMD5,
		NonceStr: strconv.FormatInt(time.Now().Unix(), 10),
	}
	req.AppBasic = app
	req.SignBasic = sign
	req.OutTradeNo = "orderNo" + req.NonceStr
	req.NotifyURL = "https://www.geluxiya.com"
	req.ServerIP = "192.168.1.123"
	req.TradeType = payment.TradeTypeJSAPI
	req.TotalFee = 10
	req.Desc = "商品描述" + req.NonceStr
	// req.OpenID = "ovxEC5YTWQk6Vv5FJdN_30gkBr-g"
	req.OpenID = openID
	signedStr := req.SignParams()
	req.Sign = signedStr
	xmlText, err := req.ToXml()
	if err != nil {
		return nil, err
	}
	ur := &payment.UnifiedOrderResponse{}
	err = http_util.PostXml(ur, unifiedorder_url, strings.NewReader(xmlText))
	if err != nil {
		return nil, err
	}
	pr := &payment.PrepayReponse{}
	pr.AppID = appID
	pr.NonceStr = strconv.FormatInt(time.Now().Unix(), 10)
	pr.SignType = payment.SignTypeMD5
	pr.TimeStamp = strconv.FormatInt(time.Now().Unix(), 10)
	pr.Package = "prepay_id=" + ur.PrepayID
	pr.SignParams()
	return pr, nil
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
	m.crontab.MustAddJob("* * * * *", m.refreshAccessToken)
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
	logger.Printf("access token refresh result is %+v\n", act)
}
