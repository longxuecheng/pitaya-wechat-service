package wechat

type BasicResponse struct {
	ErrorCode int32  `json:"errcode"`
	ErrorMsg  string `json:"errmsg"`
}

func (b *BasicResponse) isOk() bool {
	return b.ErrorCode == 0
}

type Code2SessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	BasicResponse
}

type AccessTokenResonse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	BasicResponse
}
