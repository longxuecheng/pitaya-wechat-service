package sms

import (
	"bytes"
	"encoding/json"
	"gotrue/facility/errors"
	"gotrue/facility/http_util"
	"net/http"
)

const (
	smsHost      = "http://localhost:7788"
	payNotifyURI = "/sendsms"
)

type MultiSendRequest struct {
	Mobiles []string `json:"mobiles"`
	Params  []string `json:"params"`
}

type response struct {
	ResultCode string `json:"resultCode"`
	resultDesc string `json:"resultDesc"`
}

func (r *response) IsOK() bool {
	return r.ResultCode == "OK"
}

func SendPayNotificationMsg(req *MultiSendRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	resp, err := http_util.Send(http.MethodPost, smsHost+payNotifyURI, bytes.NewReader(data), http_util.JsonHeader)
	if err != nil {
		return err
	}
	result := &response{}
	err = http_util.UnmarshalBody(resp, result)
	if err != nil {
		return err
	}
	if !result.IsOK() {
		return errors.NewWithCodef(result.ResultCode, result.resultDesc)
	}
	return nil
}
