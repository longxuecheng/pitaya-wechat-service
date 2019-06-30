package http_util

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

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

func DoGet(dst interface{}, url string, modifyFn func(r *http.Request)) error {
	rq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if modifyFn != nil {
		modifyFn(rq)
	}
	rp, err := http.DefaultClient.Do(rq)
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
