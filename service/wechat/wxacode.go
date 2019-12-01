package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gotrue/facility/http_util"
	"io/ioutil"
	"net/http"
)

// GetParameterizedMPCode params is like a=1&b=2
// page is page address like pages/goods/goods
// width is image size
func (s *wechatService) GetParameterizedMPCode(params, page string, width int) ([]byte, error) {
	url := fmt.Sprintf(wxacode_url, s.AccessToken())
	req := WxAcodeUnlimitedRequest{
		Scene: params,
		Page:  page,
		Width: width,
	}
	data, _ := json.Marshal(req)
	resp, err := http_util.Send(http.MethodPost, url, bytes.NewReader(data), http_util.JsonHeader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
