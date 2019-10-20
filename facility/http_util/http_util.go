package http_util

import (
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var client = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
	Timeout: 30 * time.Second,
}

func JsonHeader(r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
}

func UnmarshalBody(r *http.Response, dst interface{}) error {
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	// fmt.Println(string(bs))
	return json.Unmarshal(bs, dst)
}

type Handler func(*http.Request)

func Get(url string, dst interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Http get %s with status code %d", url, response.StatusCode)
	}
	bb, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bb, dst)
	if err != nil {
		return err
	}
	return nil
}

func DoGet(dst interface{}, url string, interceptors ...Handler) error {
	rq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if interceptors != nil {
		for _, interceptor := range interceptors {
			interceptor(rq)
		}
	}
	rp, err := client.Do(rq)
	if err != nil {
		return err
	}
	if rp.StatusCode != http.StatusOK {
		return fmt.Errorf("Http get %s with status code %d", url, rp.StatusCode)
	}
	bb, err := ioutil.ReadAll(rp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bb, dst)
	if err != nil {
		return err
	}
	return nil
}

func PostXml(dst interface{}, url string, data io.Reader) error {
	resp, err := http.Post(url, "application/xml", data)
	if err != nil {
		return err
	}
	fmt.Printf("HTTP post xml status code is %s \n", resp.Status)
	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("Response body is %s \n", string(bb))
	if err := xml.Unmarshal([]byte(bb), dst); err != nil {
		return err
	}
	return nil
}

func Send(method, target string, data io.Reader, handlers ...Handler) (*http.Response, error) {
	req, err := http.NewRequest(method, target, data)
	if err != nil {
		return nil, err
	}
	if handlers != nil {
		if len(req.Form) == 0 {
			req.Form = url.Values{}
		}
		for _, handler := range handlers {
			handler(req)
		}
		if http.MethodGet == method {
			req.Form = nil
		}

	}
	return client.Do(req)
}
