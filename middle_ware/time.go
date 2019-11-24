package middle_ware

import (
	"gotrue/facility/log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// before request
		c.Next()
		// after request
		elapsed := time.Since(t)
		log.Log.Debug("Time used is %0.2f seconds", elapsed.Seconds())
	}
}
