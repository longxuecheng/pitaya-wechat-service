package wechat

const (
	appID            string = "wx91eb35616f475761"
	secret           string = "41f8cb23aab151c6aa8ca64a0cfb596c"
	merchantID       string = "1540271371"
	authorize_url    string = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	accessToken_url  string = "https://api.weixin.qq.com/cgi-bin/token?grant_type=%s&appid=%s&secret=%s"
	unifiedorder_url string = "https://api.mch.weixin.qq.com/pay/unifiedorder"
)

type IWechat interface {
	AccessToken() error
}
