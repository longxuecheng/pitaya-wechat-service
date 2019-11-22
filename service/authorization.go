package service

import (
	"gotrue/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var signSecrete = []byte("geluxiya")
var issuer = "geluxiya-access-token"

var IgnoreEXP = func(token *jwt.Token, ve *jwt.ValidationError) {
	if token.Valid {
		return
	}
	if ve.Errors&(jwt.ValidationErrorExpired) != 0 {
		token.Valid = true
	}
}

var IgnoreNBF = func(token *jwt.Token, ve *jwt.ValidationError) {
	if token.Valid {
		return
	}
	if ve.Errors&(jwt.ValidationErrorNotValidYet) != 0 {
		token.Valid = true
	}
}

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
	return ss, nil
}

// ParseToken validate signed token and return claims
func ParseToken(ss string, skipExp bool) (*model.UserClaims, error) {
	claims := &model.UserClaims{}
	token, err := jwt.ParseWithClaims(ss, claims, func(token *jwt.Token) (interface{}, error) {
		return signSecrete, nil
	})
	if err == nil {
		return claims, nil
	}
	ve, ok := err.(*jwt.ValidationError)
	if !ok {
		return nil, err
	}
	if ve.Errors == jwt.ValidationErrorExpired {
		if skipExp {
			token.Valid = true
		}
	}
	if token.Valid {
		return claims, nil
	}
	return nil, err
}
