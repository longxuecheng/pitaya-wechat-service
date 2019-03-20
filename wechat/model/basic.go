package wechat_model

type BasicResponse struct {
	ErrorCode int32  `json:"errcode"`
	ErrorMsg  string `json:"errmsg"`
}

type Code2SessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	BasicResponse
}
