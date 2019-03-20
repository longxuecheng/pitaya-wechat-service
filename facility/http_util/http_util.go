package http_util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Get(url string, resonseptr interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Http get %s with status code %s", url, response.StatusCode))
	}
	body_bytes,err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body_bytes, resonseptr)
	if err != nil {
		return err
	}
	return nil
}
