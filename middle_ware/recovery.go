package middle_ware

import (
	"fmt"
	"gotrue/facility/errors"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gin-gonic/gin"
)

// Recovery returns a middleware that recovers from any panics and writes a error response
func Recovery(c *gin.Context) {
	defer func() {
		var response Response
		if err := recover(); err != nil {
			logRequest(c)
			log.Printf("[Error] %+v", err)
			if err, ok := err.(error); ok {
				if readable, ok := errors.Readable(err); ok {
					response = newResponse(nil, readable.Code(), readable.Message())
				} else {
					response = newResponse(nil, UNKNOWN_ERROR, err.Error())
				}
				c.AbortWithStatusJSON(http.StatusOK, response)
			} else {
				response := newResponse(nil, UNKNOWN_ERROR, fmt.Sprintf("%s", err))
				c.AbortWithStatusJSON(http.StatusOK, response)
			}
		}
	}()
	c.Next()
}

func logRequest(c *gin.Context) {
	httprequest, _ := httputil.DumpRequest(c.Request, false)
	log.Printf("[Recovery] %s request:\n%s", timeFormat(time.Now()), string(httprequest))
}

func timeFormat(t time.Time) string {
	var timeString = t.Format("2006/01/02 - 15:04:05")
	return timeString
}
