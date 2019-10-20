package service

import (
	"bytes"
	"context"
	"gotrue/service/api"
	"net/http"
	"net/url"

	"github.com/tencentyun/cos-go-sdk-v5"
)

const (
	cosSecretID  = "AKIDve9ONsSD4f0BNMsapiL8i26uToQn7rWv"
	cosSecretKey = "JNCm7F8mAIAoOimcD9pGDvNZWoqeCUYp"
	bucketURL    = "https://wxacode-1258625730.cos.ap-chengdu.myqcloud.com"
	stage        = "prod/"
)

var CosService api.ITencloudCos

type Cos struct{}

func NewCosService() *Cos {
	return &Cos{}
}

func InitCosService() {
	CosService = NewCosService()
}

func (s *Cos) PushImageObject(name string, data []byte) error {
	u, _ := url.Parse(bucketURL)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cosSecretID,
			SecretKey: cosSecretKey,
		},
	})
	name = stage + name
	_, err := c.Object.Put(context.Background(), name, bytes.NewReader(data), nil)
	if err != nil {
		return err
	}
	return nil
}
