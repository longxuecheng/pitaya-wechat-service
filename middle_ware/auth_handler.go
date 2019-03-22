package middle_ware

import (
	"net/http"
	"pitaya-wechat-service/service"

	"github.com/gin-gonic/gin"
)

type CurrentUser struct {
	WechatID string
	UserID   int64
}

const (
	accessTokenKey = "X-Pitaya-Token"
	currentUserKey = "current-user"
)

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. get token from header
		accessToken := c.GetHeader(accessTokenKey)
		if accessToken == "" {
			c.AbortWithStatusJSON(http.StatusOK, newResponseDTO(nil, "user unauthorized", STATUS_UNAUTHORIZED))
			return
		}
		err := service.ValidateToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, newResponseDTO(nil, "invalid token or expired", STATUS_TOKEN_EXP))
			return
		}
		c.Next()
	}
}

func MustGetCurrentUser(c *gin.Context) CurrentUser {
	return c.MustGet(currentUserKey).(CurrentUser)
}
