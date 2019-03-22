package service

import (
	"fmt"
	"pitaya-wechat-service/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var signSecrete = []byte("geluxiya")
var issuer = "geluxiya-access-token"

// BuildToken build a jwt token with a total life in seconds
func BuildToken(userID int64, ttl int64) (string, error) {
	claim := model.UserClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + ttl,
			Issuer:    issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	ss, err := token.SignedString(signSecrete)
	if err != nil {
		return "", err
	}
	fmt.Printf("%v %v", ss, err)
	return ss, nil
}

// ValidateToken validate signed string
func ValidateToken(ss string) error {
	token, err := jwt.ParseWithClaims(ss, &model.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signSecrete, nil
	})
	if err != nil {
		return err
	}
	if claims, ok := token.Claims.(*model.UserClaims); ok && token.Valid {
		expiresAt := claims.StandardClaims.ExpiresAt
		if time.Now().Unix() >= expiresAt {
			return jwt.NewValidationError("token has already expired", jwt.ValidationErrorExpired)
		}
	}
	return nil
}
