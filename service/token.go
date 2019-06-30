package service

import (
	"gotrue/model"
	"gotrue/sys"
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
	// cacheCurrentUser(ss, userID)
	return ss, nil
}

func cacheCurrentUser(token string, userID int64) {
	sys.UserCache().Add(token, 0, userID)
}

// ParseToken validate signed token and return claims
func ParseToken(ss string) (*model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(ss, &model.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signSecrete, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*model.UserClaims); ok {
		if token.Valid {
			expiresAt := claims.StandardClaims.ExpiresAt
			if time.Now().Unix() >= expiresAt {
				return nil, jwt.NewValidationError("token has already expired", jwt.ValidationErrorExpired)
			}
			return claims, nil
		}
		return nil, jwt.NewValidationError("token invalid", jwt.ValidationErrorNotValidYet)

	}
	return nil, jwt.NewValidationError("token claims invalid", jwt.ValidationErrorClaimsInvalid)
}
