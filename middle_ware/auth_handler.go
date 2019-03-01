package middle_ware

import (
	"github.com/gin-gonic/gin"
)

var authTokenKey = "Auth-Token"

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. get token from header
		// 2. parse token and get cache key
		// 3. get cache data by cache key parsed in 2nd step
		// token := c.Request.Header.Get(authTokenKey)
		c.Next()

	}
}
