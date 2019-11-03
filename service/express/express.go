package express

import (
	"encoding/json"
	"fmt"
	"gotrue/dto/response"
	"gotrue/facility/errors"
	"gotrue/facility/http_util"
	"gotrue/facility/strings"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	baiduPartnerID   string      = "4001"
	baiduExpressURL  string      = "https://sp0.baidu.com/9_Q4sjW91Qh3otqbppnN2DJv/pae/channel/data/asyncqury"
	ems              expressType = "ems"
	youzheng         expressType = "youzhengguonei"
	bsht             expressType = "huitongkuaidi"
	sto              expressType = "shentong"
	yunda            expressType = "yunda"
	zto              expressType = "zhongtong"
	yto              expressType = "yuantong"
	tt               expressType = "tiantian"
	expressErrorTemp string      = "ExressError%s"
)

var baiduExpressMap = map[ExpressMethod]expressType{
	ExpressMethodZTO:  zto,
	ExpressMethodSTO:  sto,
	ExpressMethodYTO:  yto,
	ExpressMethodEMS:  ems,
	ExpressMethodYZ:   youzheng,
	ExpressMethodYDA:  yunda,
	ExpressMethodBSHT: bsht,
	ExpressMethodTT:   tt,
}

type expressType string

func (et expressType) String() string {
	return string(et)
}

type ExpressBaseResponse struct {
	Msg       string          `json:"msg"`
	Status    string          `json:"status"`
	ErrorCode string          `json:"error_code"`
	Data      json.RawMessage `json:"data"`
}

func (r *ExpressBaseResponse) IsOK() bool {
	return r.Status == "0"
}

func (r *ExpressBaseResponse) Error() error {
	return errors.NewWithCodef(fmt.Sprintf(expressErrorTemp, r.ErrorCode), r.Msg)
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

func (s *expressService) ExpressList() []*response.Express {
	expressList := []*response.Express{}
	for k := range expressMethodMap {
		exp := &response.Express{
			Method: k.String(),
			Name:   k.Name(),
		}
		expressList = append(expressList, exp)
	}
	return expressList
}

func (s *expressService) GetExpressFromChinaPost(expressNo string) (*ExpressSummary, error) {
	slideDecoder := NewSlideDecoder()
	summary := &ExpressSummary{
		ExpressNo: expressNo,
		Company:   "中国邮政",
	}
	err := slideDecoder.LoadVerifyCode()
	if err != nil {
		summary.Traces = []*ExpressTrace{
			&ExpressTrace{
				Time: "",
				Desc: "从邮政获取信息异常，请不用担心",
			},
		}
		return summary, nil
	}
	err = slideDecoder.CheckStartPosition()
	if err != nil {
		summary.Traces = []*ExpressTrace{
			&ExpressTrace{
				Time: "",
				Desc: "解析邮政验证码错误，30秒后重新尝试查询，请不用担心",
			},
		}
		return summary, nil
	}
	chinaPostTraces, err := slideDecoder.QueryExpress(expressNo)
	commonTraces := make([]*ExpressTrace, len(chinaPostTraces))
	for i, chinaPostTrace := range chinaPostTraces {
		commonTraces[i] = chinaPostTrace.ExpressTrace()
	}
	summary.Traces = commonTraces
	return summary, nil
}

func (s *expressService) ExpressInfo(expressCom ExpressMethod, expressNo string) (*ExpressSummary, error) {
	expressType, ok := baiduExpressMap[expressCom]
	if !ok {
		return nil, errors.NewWithCodef("ExpressNotSupport", "当前不支持%s", expressCom.Name())
	}
	if strings.IsEmpty(expressNo) {
		return nil, nil
	}
	if expressType == youzheng {
	}
	resp, err := http_util.Send(http.MethodGet, baiduExpressURL, nil, func(r *http.Request) {
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
		r.Form.Set("appid", baiduPartnerID)
		r.Form.Set("com", expressType.String())
		r.Form.Set("nu", expressNo)
		r.URL.RawQuery = r.Form.Encode()
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
	if baseResult.IsOK() {
		expressInfo := &ExpressBody{}
		err = json.Unmarshal(baseResult.Data, expressInfo)
		if err != nil {
			return nil, err
		}
		expressInfo.Info.ExpressNo = expressNo
		return expressInfo.Info, nil
	}
	return nil, baseResult.Error()
}
