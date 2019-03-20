package service

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/muesli/cache2go"
)

var hmacSampleSecret = []byte("my_secret_key")
var TokenCache *cache2go.CacheTable

func init() {
	TokenCache = cache2go.Cache("token_cache")
}

func Authorize() (string, error) {
	// following token
	jwtMap := jwt.MapClaims{
		"name": "lxc",
		"id":   1,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtMap)
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		return "", err
	}
	TokenCache.Add(tokenString, 5*time.Second, jwtMap)
	return tokenString, nil
}
