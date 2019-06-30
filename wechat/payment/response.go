package payment

import "fmt"

// StatusResponse 代表接口状态结果，属于接口层结果
type StatusResponse struct {
	Code string `xml:"return_code"`
	Msg  string `xml:"return_msg"`
}

// BusinessResponse 在StatusResponse.Code为Success时返回，属于业务层结果
type BusinessResponse struct {
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrMsg     string `xml:"err_code_des"`
}

// UnifiedOrderResponse 在StatusResponse.Code为成功且
// BusinessResponse.ResultCode为成功时才有返回
type UnifiedOrderResponse struct {
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
}

func (pr *PrepayReponse) SignParams() string {
	pr.PaySign = signStruct(pr, "json")
	fmt.Printf("Prepay response sign string is %s", pr.PaySign)
	return pr.PaySign
}
