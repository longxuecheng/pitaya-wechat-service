package express

import (
	"encoding/base64"
	"errors"
	"fmt"
	"gotrue/facility/http_util"
	"image/png"
	"log"
	"net/http"
	"net/url"

	"strconv"
	"strings"
	"time"
)

const (
	loadVerifyCodeURL  = "http://yjcx.chinapost.com.cn/qps/showPicture/verify/slideVerifyLoad?t=%d"
	checkVeriryCodeURL = "http://yjcx.chinapost.com.cn/qps/showPicture/verify/slideVerifyCheck"
)

const whiteColor uint32 = 65535

type SlideImage struct {
	MainBase64PNGImage string `json:"YYPng_base64"`
	CutBase64PNGImage  string `json:"CutPng_base64"`
	UUID               string `json:"uuid"`
}

type SlideDecoder struct {
	startX     int
	startY     int
	slideImage *SlideImage
}

func NewSlideDecoder() *SlideDecoder {
	return &SlideDecoder{}
}

func NewSlideDecoderWitStart(start int, uuid string) *SlideDecoder {
	return &SlideDecoder{
		startX: start,
		slideImage: &SlideImage{
			UUID: uuid,
		},
	}
}

func (s *SlideDecoder) LoadVerifyCode() error {
	slideImg := &SlideImage{}
	url := fmt.Sprintf(loadVerifyCodeURL, time.Now().Unix())
	resp, err := http_util.Send(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	err = http_util.UnmarshalBody(resp, slideImg)
	if err != nil {
		return err
	}
	s.slideImage = slideImg
	return nil
}

func (s *SlideDecoder) CheckStartPosition() error {
	pngImage := strings.NewReader(s.slideImage.MainBase64PNGImage)
	pngReader := base64.NewDecoder(base64.StdEncoding, pngImage)
	img, err := png.Decode(pngReader)
	if err != nil {
		return err
	}
	leftTop := img.Bounds().Min
	rightBotom := img.Bounds().Max
	found := false

	for x := leftTop.X; x < rightBotom.X; x++ {
		for y := leftTop.Y; y < rightBotom.Y; y++ {
			f1 := false
			f2 := false

			r, g, b, a := img.At(x, y).RGBA()

			if r == whiteColor && g == whiteColor && b == whiteColor && a == whiteColor {
				f1 = true
			}
			r1, g1, b1, a1 := img.At(x, y+10).RGBA()
			if r1 == whiteColor && g1 == whiteColor && b1 == whiteColor && a1 == whiteColor {
				f2 = true
			}
			if f1 && f2 {
				log.Printf("Found start x position %d y position %d\n", x, y)
				s.startX = x
				s.startY = y
				found = true
				break
			}

		}
		if found {
			break
		}
	}
	return nil
}

func (s *SlideDecoder) StartX() int64 {
	return int64(s.startX)
}

func (s *SlideDecoder) UUID() string {
	return s.slideImage.UUID
}

func (s *SlideDecoder) QueryExpress(number string) ([]*ChinaPostExpressTrace, error) {
	movedEndX := strconv.FormatInt(s.StartX(), 10)
	postData := url.Values{}
	postData.Add("uuid", s.UUID())
	postData.Add("moveEnd_X", movedEndX)
	postData.Add("text[]", number)
	postData.Add("selectType", "1")
	resp, err := http_util.Send(http.MethodPost, checkVeriryCodeURL, strings.NewReader(postData.Encode()), func(r *http.Request) {
		r.Header.Set("Cookie", "JSESSIONID=1DC89A777760DD5AD3320F3FAE865A82; hibext_instdsigdipv2=1")
		r.Header.Set("Referer", "http://yjcx.chinapost.com.cn/qps/yjcx")
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	})
	if err != nil {
		return nil, err
	}
	traceList := []*ChinaPostExpressTrace{}
	err = http_util.UnmarshalBody(resp, &traceList)
	if err != nil {
		return nil, ErrorNeedRetry
	}
	return traceList, nil
}

var ErrorNeedRetry = errors.New("Analyze error,please retry")

func (s *SlideDecoder) NewJSESSIONCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "JSESSIONID",
		Value:    "1DC89A777760DD5AD3320F3FAE865A82",
		Domain:   "yjcx.chinapost.com.cn",
		Path:     "/",
		Expires:  time.Now().AddDate(1, 0, 0),
		HttpOnly: false,
		Secure:   false,
	}
}
func (s *SlideDecoder) NewAnotherCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "hibext_instdsigdipv2",
		Value:    "1",
		Domain:   "yjcx.chinapost.com.cn",
		Path:     "/",
		Expires:  time.Now().AddDate(1, 0, 0),
		HttpOnly: false,
		Secure:   false,
	}
}

type ChinaPostExpressTrace struct {
	TraceNo         string `json:"traceNo"`
	OpTime          string `json:"opTime"`
	OpOrgCode       string `json:"opOrgCode"`
	OpOrgName       string `json:"opOrgName"`
	OpName          string `json:"opName"`
	OpCode          string `json:"opCode"`
	OpOrgCity       string `json:"opOrgCity"`
	BizProductName  string `json:"bizProductName"`
	BaseProductName string `json:"baseProductName"`
	PostDate        string `json:"postDate"`
	StatusDesc      string `json:"statusDesc"`
	Destination     string `json:"destination"`
	Date            string `json:"date"`
	Time            string `json:"time"`
	OpCodeStatus    string `json:"opCodeStatus"`
}

func (c *ChinaPostExpressTrace) ExpressTrace() *ExpressTrace {
	return &ExpressTrace{
		Time: c.OpTime,
		Desc: c.StatusDesc,
	}
}
