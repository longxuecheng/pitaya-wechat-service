package middle_ware

import (
	"gotrue/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	accessTokenKey = "Authorization"
	currentUserKey = "current-user"
)

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader(accessTokenKey)
		if accessToken == "" {
			c.AbortWithStatusJSON(http.StatusOK, newResponse(nil, "user unauthorized", STATUS_UNAUTHORIZED))
			return
		}
		userClaims, err := service.ParseToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, newResponse(nil, "invalid token or expired", STATUS_TOKEN_EXP))
			return
		}
		c.Set(currentUserKey, userClaims.UserID)
		c.Next()
	}
}

func MustGetCurrentUser(c *gin.Context) int64 {
	return c.MustGet(currentUserKey).(int64)
}
