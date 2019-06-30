package middle_ware

import (
	"net/http"
	"gotrue/service"

	"github.com/gin-gonic/gin"
)

const (
	accessTokenKey = "X-Pitaya-Token"
	currentUserKey = "current-user"
)

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader(accessTokenKey)
		if accessToken == "" {
			c.AbortWithStatusJSON(http.StatusOK, newResponseDTO(nil, "user unauthorized", STATUS_UNAUTHORIZED))
			return
		}
		userClaims, err := service.ParseToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, newResponseDTO(nil, "invalid token or expired", STATUS_TOKEN_EXP))
			return
		}
		c.Set(currentUserKey, userClaims.UserID)
		c.Next()
	}
}

func MustGetCurrentUser(c *gin.Context) int64 {
	return c.MustGet(currentUserKey).(int64)
}
