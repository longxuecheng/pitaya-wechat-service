package payment

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"gotrue/model"
	"testing"
)

func TestXMLMarshal(t *testing.T) {
	app := AppBasic{
		AppID:      "xxxxxx",
		MerchantID: "m-dfjdfkasjdkf",
	}
	ur := UnifiedOrderRequest{
		AppBasic: app,
	}
	bytes, err := xml.Marshal(ur)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(bytes))
}

func TestXMLUnmarshal(t *testing.T) {
	blob := `<xml>
	<appid>xxxxxx</appid>
	<mch_id>m-dfjdfkasjdkf</mch_id>
	<nonce_str></nonce_str>
	<sign></sign>
	<sign_type></sign_type>
	<body></body>
	<out_trade_no>orderNo-132123</out_trade_no>
	<total_fee>100000</total_fee>
	<spbill_create_ip></spbill_create_ip>
	<notify_url></notify_url>
	<trade_type>JSAPI</trade_type>
	<err>dfjkd</err>
</xml>
`
	v1 := UnifiedOrderRequest{}
	if err := xml.Unmarshal([]byte(blob), &v1); err != nil {
		t.Error(err.Error())
	}
	fmt.Println(v1)
}

var app = AppBasic{
	AppID:      "xxxxxx",
	MerchantID: "m-dfjdfkasjdkf",
}
var sign = SignBasic{
	NonceStr: "随机数",
	SignType: SignTypeMD5,
}
var v = UnifiedOrderRequest{
	AppBasic:   app,
	OutTradeNo: "orderNo-132123",
	TotalFee:   100000,
	Desc:       "商品描述",
	TradeType:  TradeTypeJSAPI,
	SignBasic:  sign,
	NotifyURL:  "https://www.geluxiya.com",
	ServerIP:   "192.193.123.231",
}

func TestCustomerMarshal(t *testing.T) {
	t.Run("sign", TestSign)
	var b bytes.Buffer
	if err := xml.NewEncoder(&b).EncodeElement(v, xml.StartElement{Name: xml.Name{Local: "xml"}}); err != nil {
		t.Error(err.Error())
	}
	xmlText := string(b.Bytes())
	fmt.Println(xmlText)

}

func TestSign(t *testing.T) {
	v.Sign = v.SignParams()
}

func TestQueryOrderStatus(t *testing.T) {
	r := &QueryOrderResponse{
		trade: trade{
			TradeState: Success,
		},
	}
	orderstatus := r.OrderStatus()
	if orderstatus == model.Paid {
		fmt.Println(orderstatus)
	}
}
