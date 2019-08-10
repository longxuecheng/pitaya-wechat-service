package payment

import (
	"database/sql/driver"
	"fmt"

	"go.planetmeican.com/meican.x/manage/facility/errors"
)

const success = "SUCCESS"

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
type tradeState string

func (ts tradeState) String() string {
	return string(ts)
}

func (ts tradeState) Value() (driver.Value, error) {
	return ts.String(), nil
}

const (
	Success  tradeState = "SUCCESS"
	Refund   tradeState = "REFUND"
	NotPay   tradeState = "NOTPAY"
	CLOSED   tradeState = "CLOSED"
	Revoked  tradeState = "REVOKED"
	Paying   tradeState = "USERPAYING"
	PayError tradeState = "PAYERROR"
)
