package middle_ware

import (
	"gotrue/facility/context_util"
	"gotrue/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	accessTokenKey = "Authorization"
	currentUserKey = "current-user"
)

// ValidateAuthorization check whether a request is authorized
func ValidateAuthorization(c *gin.Context) {
	accessToken := c.GetHeader(accessTokenKey)
	if accessToken == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, newResponse(nil, "user unauthorized", UNAUTHORIZED))
		return
	}
	userClaims, err := service.ParseToken(accessToken, true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, newResponse(nil, "invalid token or expired", TOKEN_EXP))
		return
	}
	c.Set(currentUserKey, userClaims.UserID)
	InjectUserID(c, userClaims.UserID)
	c.Next()
}

func InjectUserID(c *gin.Context, userID int64) {
	ctx := context_util.WithUserID(c.Request.Context(), userID)
	copy := c.Request.WithContext(ctx)
	c.Request = copy
}

func MustGetCurrentUser(c *gin.Context) int64 {
	return c.MustGet(currentUserKey).(int64)
}
