package payment

import (
	"bytes"
	"crypto/md5"
	"database/sql/driver"
	"encoding/xml"
	"fmt"
	"gotrue/facility/utils"
	"gotrue/model"
	"sort"
	"strconv"
	"strings"

	"go.planetmeican.com/meican.x/manage/facility/errors"
)

const (
	SignTypeMD5         string = "MD5"
	SignTypeHMAC_SHA256 string = "HMAC_SHA256"
)

type TradeType string

const (
	TradeTypeJSAPI TradeType = "JSAPI"
)

type AppBasic struct {
	AppID      string `xml:"appid"`
	MerchantID string `xml:"mch_id"`
}

type SignBasic struct {
	NonceStr string `xml:"nonce_str"`
	Sign     string `xml:"sign"`
	SignType string `xml:"sign_type"`
}

// UnifiedOrderRequest 统一下单请求
type UnifiedOrderRequest struct {
	AppBasic
	SignBasic
	OpenID     string    `xml:"openid"`
	Desc       string    `xml:"body"`
	OutTradeNo string    `xml:"out_trade_no"`
	TotalFee   int64     `xml:"total_fee"`
	ServerIP   string    `xml:"spbill_create_ip"`
	NotifyURL  string    `xml:"notify_url"`
	TradeType  TradeType `xml:"trade_type"`
}

func (u *UnifiedOrderRequest) SignParams() string {
	return signStruct(u, "xml")
}

func (u *UnifiedOrderRequest) ToXml() (string, error) {
	if u.Sign == "" {
		return "", errors.NewWithCodef("PaySignError", "request param must be signed first")
	}
	var b bytes.Buffer
	if err := xml.NewEncoder(&b).EncodeElement(u, xml.StartElement{Name: xml.Name{Local: "xml"}}); err != nil {
		return "", err
	}
	return string(b.Bytes()), nil
}

// OrderQueryRequest 订单查询请求
type OrderQueryRequest struct {
	AppBasic
	SignBasic
	OutTradeNo string `xml:"out_trade_no"`
}

func (oqr *OrderQueryRequest) ToXml() (string, error) {
	oqr.Sign = signStruct(oqr, "xml")
	if oqr.Sign == "" {
		return "", errors.NewWithCodef("PaySignError", "request param must be signed first")
	}
	var b bytes.Buffer
	if err := xml.NewEncoder(&b).EncodeElement(oqr, xml.StartElement{Name: xml.Name{Local: "xml"}}); err != nil {
		return "", err
	}
	return string(b.Bytes()), nil
}

// signStruct 微信支付的签名算法
func signStruct(input interface{}, tag string) string {
	fieldMap := utils.StructToMap(input, tag, "exlude")
	fs := []string{}
	for k := range fieldMap {
		fs = append(fs, k)
	}
	sort.Slice(fs, func(i, j int) bool {
		return fs[i] < fs[j]
	})
	fmt.Println("ordered string key list is ", fs)
	sb := strings.Builder{}
	for _, item := range fs {
		if v, ok := fieldMap[item]; ok {
			// if value is not an empty string
			if s, ok := v.(string); ok {
				if s != "" {
					sb.WriteString(item + "=")
					sb.WriteString(s)
					sb.WriteString("&")
				}
			}
			if t, ok := v.(TradeType); ok {
				if t != "" {
					sb.WriteString(item + "=")
					sb.WriteString(string(t))
					sb.WriteString("&")
				}
			}
			if n, ok := v.(int64); ok {
				if n != 0 {
					sb.WriteString(item + "=")
					sb.WriteString(strconv.FormatInt(n, 10))
					sb.WriteString("&")
				}
			}
		}
	}
	// 将key写入前面字符串最后
	sb.WriteString("key=")
	sb.WriteString("geluxiyaguozhen18911792314789012")
	fmt.Printf("pre string for sign is %s \n", sb.String())
	prestring := sb.String()
	// MD5加密
	sum := md5.Sum([]byte(prestring))
	signedStr := strings.ToUpper(fmt.Sprintf("%x", sum))
	sb.Reset()
	fmt.Printf("Sign string is %s \n", signedStr)
	return signedStr
}

const success = "SUCCESS"

const (
	Success  tradeState = "SUCCESS"
	Refund   tradeState = "REFUND"
	NotPay   tradeState = "NOTPAY"
	CLOSED   tradeState = "CLOSED"
	Revoked  tradeState = "REVOKED"
	Paying   tradeState = "USERPAYING"
	PayError tradeState = "PAYERROR"
)

type tradeState string

func (ts tradeState) String() string {
	return string(ts)
}

func (ts tradeState) Value() (driver.Value, error) {
	return ts.String(), nil
}

// baseResponse 代表接口状态结果，属于接口层结果
type baseResponse struct {
	Code       string `xml:"return_code"`
	Msg        string `xml:"return_msg"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrMsg     string `xml:"err_code_des"`
}

// IsOK check whether response is succeed
func (sr *baseResponse) IsOK() (bool, error) {
	if sr.Code != success {
		return false, errors.NewWithCodef(sr.Code, sr.Msg)
	}
	if sr.ResultCode != success {
		return false, errors.NewWithCodef(sr.ErrCode, sr.ErrMsg)
	}
	return true, nil
}

// UnifiedOrderResponse 在StatusResponse.Code为成功且
// BusinessResponse.ResultCode为成功时才有返回
type UnifiedOrderResponse struct {
	baseResponse
	AppInfo   AppBasic
	SignInfo  SignBasic
	TradeType TradeType `xml:"trade_type"`
	PrepayID  string    `xml:"prepay_id"`
	CodeURL   string    `xml:"code_url"`
}

type PrepayReponse struct {
	AppID     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign" exclude:"true"`
	PrepayID  string `json:"prepayId" exclude:"true"`
}

func (pr *PrepayReponse) SignParams() string {
	pr.PaySign = signStruct(pr, "json")
	fmt.Printf("Prepay response sign string is %s", pr.PaySign)
	return pr.PaySign
}

type trade struct {
	TradeState     tradeState `xml:"trade_state"`
	TradeType      string     `xml:"trade_type"`
	OutTradeNo     string     `xml:"out_trade_no"`
	Attach         string     `xml:"attach"`
	TimeEnd        string     `xml:"time_end"` // yyyyMMddHHmmss
	TradeStateDesc string     `xml:"trade_state_desc"`
}

type price struct {
	BankeType          string `xml:"bank_type"`
	TotalFee           int64  `xml:"total_fee"`
	SettlementTotalFee int64  `xml:"settlement_total_fee"`
	FeeType            string `xml:"fee_type"`
	CashFee            int64  `xml:"cash_fee"`
	CashFeeType        string `xml:"cash_fee_type"`
}

type couponSummary struct {
	CouponFee   int64 `xml:"coupon_fee"`
	CouponCount int64 `xml:"coupon_count"`
}

type QueryOrderResponse struct {
	baseResponse
	AppBasic
	SignBasic
	trade
	price
	couponSummary
}

func (q *QueryOrderResponse) OrderStatus() model.OrderStatus {
	var orderStatus model.OrderStatus
	if q.TradeState == Success {
		orderStatus = model.Paid
	}
	if q.TradeState == Paying || q.TradeState == NotPay {
		orderStatus = model.Paying
	}
	if q.TradeState == CLOSED {
		orderStatus = model.Closed
	}
	if q.TradeState == PayError {
		orderStatus = model.PayFailed
	}

	return orderStatus
}
