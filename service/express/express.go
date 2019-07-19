package express

import (
	"encoding/json"
	"gotrue/facility/http_util"
	"gotrue/facility/strings"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	url  string      = "https://sp0.baidu.com/9_Q4sjW91Qh3otqbppnN2DJv/pae/channel/data/asyncqury"
	EMS  ExpressType = "ems"
	BSHT ExpressType = "huitongkuaidi"
	ST   ExpressType = "shentong"
	YD   ExpressType = "yunda"
	ZT   ExpressType = "zhongtong"
	YT   ExpressType = "yuantong"
)

type ExpressType string

func (et ExpressType) String() string {
	return string(et)
}

type ExpressBaseResponse struct {
	Msg       string          `json:"msg"`
	Status    string          `json:"status"`
	ErrorCode string          `json:"error_code"`
	Data      json.RawMessage `json:"data"`
}

type ExpressBody struct {
	Info *ExpressSummary `json:"info"`
}

type ExpressSummary struct {
	ExpressNo      string          `json:"express_no"`
	Status         string          `json:"status"`
	Company        string          `json:"com"`
	State          string          `json:"state"`
	SendTime       string          `json:"send_time"`
	DepartureCity  string          `json:"departure_city"`
	ArrivalCity    string          `json:"arrival_city"`
	LatestProgress string          `json:"latest_progress"`
	Traces         []*ExpressTrace `json:"context"`
}

type ExpressTrace struct {
	Time string `json:"time"`
	Desc string `json:"desc"`
}

var ExpressService *expressService

type expressService struct{}

func init() {
	ExpressService = &expressService{}
}
func (s *expressService) ExpressInfo(expressCom ExpressType, expressNo string) (*ExpressSummary, error) {
	if strings.IsEmpty(expressNo) {
		return nil, nil
	}
	resp, err := http_util.Send(http.MethodGet, url, nil, func(r *http.Request) error {
		c := &http.Cookie{
			Name:     "BAIDUID",
			Value:    "363248B1E700BC951CA9F586683F104D:FG=1",
			Domain:   "baidu.com",
			Path:     "/",
			Expires:  time.Now().AddDate(1, 0, 0),
			HttpOnly: false,
			Secure:   false,
		}
		r.AddCookie(c)
		r.Form.Set("appid", "4001")
		r.Form.Set("com", expressCom.String())
		r.Form.Set("nu", expressNo)
		r.URL.RawQuery = r.Form.Encode()
		return nil
	})
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	baseResult := &ExpressBaseResponse{}
	err = json.Unmarshal(bytes, baseResult)
	if err != nil {
		return nil, err
	}
	if baseResult.Status == "0" {
		expressInfo := &ExpressBody{}
		err = json.Unmarshal(baseResult.Data, expressInfo)
		if err != nil {
			return nil, err
		}
		expressInfo.Info.ExpressNo = expressNo
		return expressInfo.Info, nil
	}
	return nil, nil
}
