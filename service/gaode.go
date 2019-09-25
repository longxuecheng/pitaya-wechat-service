package service

import (
	"encoding/json"
	"gotrue/dto/response"
	"gotrue/facility/errors"
	"gotrue/facility/http_util"
	"gotrue/service/api"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	key          = "dd528118a358c80c97b9507c7683af86"
	geoURL       = "https://restapi.amap.com/v3/geocode/geo"
	inputTipsURL = "https://restapi.amap.com/v3/assistant/inputtips"
)

var GaodeMapService api.IGaodeMapService

type GeoRequest struct {
	Key      string
	Address  string
	City     string
	Batch    bool
	Sign     string
	Output   string
	Callback string
}

func (g *GeoRequest) EncodeURL() string {
	vals := url.Values{}
	vals.Set("key", g.Key)
	vals.Set("address", g.Address)
	vals.Set("city", g.City)
	vals.Set("batch", strconv.FormatBool(g.Batch))
	vals.Set("sig", g.Sign)
	vals.Set("output", "JSON")
	vals.Set("callback", "")
	return vals.Encode()
}

type BaseResonse struct {
	Status   string `json:"status"`
	Count    string `json:"count"`
	Info     string `json:"info"`
	InfoCode string `json:"infocode"`
}

func (r *BaseResonse) IsOK() bool {
	return r.Status == "1"
}

func (r *BaseResonse) Error() error {
	return errors.NewWithCodef(r.InfoCode, r.Info)
}

func (r *BaseResonse) HasResult() bool {
	return r.Count != "0"
}

type Location string

func (l Location) Coordinate() (lng float32, lat float32) {
	coordinates := strings.Split(string(l), ",")
	lng64, _ := strconv.ParseFloat(coordinates[0], 32)
	lat64, _ := strconv.ParseFloat(coordinates[1], 32)
	lng = float32(lng64)
	lat = float32(lat64)
	return
}

type GeoResponse struct {
	BaseResonse
	GeoCodes []GeoCodeResponse `json:"geocodes"`
}

type GeoCodeResponse struct {
	FormattedAddress string   `json:"formatted_address"`
	Country          string   `json:"country"`
	Province         string   `json:"province"`
	City             string   `json:"city"`
	CityCode         string   `json:"citycode"`
	District         []string `json:"dictrict"`
	Street           []string `json:"street"`
	Number           []string `json:"number"`
	AddressCode      string   `json:"adcode"` // 区域编码感觉像邮编
	Location         Location `json:"location"`
	Level            string   `json:"level"`
}

type dataType string

const (
	all     dataType = "all"
	poi     dataType = "poi"
	bus     dataType = "bus"
	busline dataType = "busline"
)

type InputTipRequest struct {
	Key       string
	Keywords  string
	PoiType   string
	Location  string
	City      string
	CityLimit bool
	DataType  string
	Sign      string
	Output    string
	Callback  string
}

func (r *InputTipRequest) EncodeURL() string {
	vals := url.Values{}
	vals.Set("key", r.Key)
	vals.Set("keywords", r.Keywords)
	vals.Set("type", r.PoiType)
	vals.Set("location", r.Location)
	vals.Set("city", r.City)
	vals.Set("citylimit", strconv.FormatBool(r.CityLimit))
	vals.Set("datatype", r.DataType)
	vals.Set("sig", r.Sign)
	vals.Set("output", "JSON")
	vals.Set("callback", "")
	return vals.Encode()
}

type TipResponse struct {
	BaseResonse
	Tips []TipUnit
}

func (r *TipResponse) AddressTips() ([]string, error) {
	tips := []string{}
	for _, tip := range r.Tips {
		ads, err := tip.UnmarshalAddress()
		if err != nil {
			return nil, err
		}
		tips = append(tips, ads...)
	}
	return tips, nil
}

type TipUnit struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	District    string          `json:"ditrict"`
	AddressCode string          `json:"adcode"`
	Location    Location        `json:"location"`
	Address     json.RawMessage `json:"address"`
}

func (r *TipUnit) UnmarshalAddress() ([]string, error) {
	var ad string
	if err := json.Unmarshal(r.Address, &ad); err == nil {
		return []string{ad}, nil
	}
	adlist := []string{}
	if err := json.Unmarshal(r.Address, &adlist); err != nil {
		return nil, err
	}
	return adlist, nil
}

func InitGaodeMapService() {
	GaodeMapService = &GaodeMap{}
}

type GaodeMap struct{}

func (s *GaodeMap) AddressLocation(city string, address string) (*response.Location, error) {
	req := &GeoRequest{
		Key:     key,
		City:    city,
		Address: address,
	}
	data := &GeoResponse{}
	err := http_util.DoGet(data, geoURL, func(r *http.Request) {
		r.URL.RawQuery = req.EncodeURL()
	})
	if err != nil {
		return nil, err
	}
	if !data.IsOK() {
		return nil, data.Error()
	}
	if !data.HasResult() {
		return nil, nil
	}
	lng, lat := data.GeoCodes[0].Location.Coordinate()
	return &response.Location{
		Longitude: lng,
		Latitude:  lat,
	}, nil
}

func (s *GaodeMap) AddressTips(keywords, lng, lat string) ([]string, error) {
	req := &InputTipRequest{
		Key:      key,
		Keywords: keywords,
	}
	if lng != "" && lat != "" {
		req.Location = lng + "," + lat
	}
	data := &TipResponse{}
	err := http_util.DoGet(data, inputTipsURL, func(r *http.Request) {
		r.URL.RawQuery = req.EncodeURL()
	})
	if err != nil {
		return nil, err
	}
	if !data.IsOK() {
		return nil, data.Error()
	}
	return data.AddressTips()
}

func (s *GaodeMap) ReverseGeo() error {
	return nil
}
