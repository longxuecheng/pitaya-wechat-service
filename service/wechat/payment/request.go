package payment

import (
	"bytes"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"gotrue/facility/utils"
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
