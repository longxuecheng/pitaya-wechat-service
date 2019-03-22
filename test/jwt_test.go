package test

import (
	"pitaya-wechat-service/model"
	"pitaya-wechat-service/service"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var hmacSampleSecret = []byte("my_secret_key")

func TestBuildToken(t *testing.T) {
	token, err := service.BuildToken(100, time.Now().Add(time.Second).Unix())
	t.Logf("token is %s err is %v", token, err)
	parseCustom(token, t)
}

// sample token is expired.  override time so it parses as valid
func parseCustom(ss string, t *testing.T) {
	token, err := jwt.ParseWithClaims(ss, &model.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return hmacSampleSecret, nil
	})
	time.Sleep(4 * time.Second)
	if claims, ok := token.Claims.(*model.UserClaims); ok && token.Valid {
		expiresAt := claims.StandardClaims.ExpiresAt
		t.Logf("user id is %v expires at %v", claims.UserID, expiresAt)
		if time.Now().Unix() >= expiresAt {
			t.Logf("now is %d and token is now expired already %d !", time.Now().Unix(), expiresAt)
		}
	} else {
		t.Fatal(err)
	}
}
