package middle_ware

import (
	"face2face/backend/facility/errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gin-gonic/gin"
)

// Recovery returns a middleware that recovers from any panics and writes a error response
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			var response Response
			if err := recover(); err != nil {
				logRequest(c)
				if err, ok := err.(error); ok {
					if readable, ok := errors.Readable(err); ok {
						readable.StackTrace()
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
}

func logRequest(c *gin.Context) {
	httprequest, _ := httputil.DumpRequest(c.Request, false)
	log.Printf("[Recovery] %s request:\n%s", timeFormat(time.Now()), string(httprequest))
}

// func printStackTrace(err error) {
// pc := make([]uintptr, 10)
// n := runtime.Callers(4, pc)
// if n == 0 {
// 	return
// }
// pc = pc[:n]
// frames := runtime.CallersFrames(pc)
// for {
// 	frame, more := frames.Next()
// 	log.Printf("| %s file:%s line:%d entry address:%d\n", frame.Function, frame.File, frame.Line, frame.Entry)
// 	if !more {
// 		break
// 	}
// }
// }

func timeFormat(t time.Time) string {
	var timeString = t.Format("2006/01/02 - 15:04:05")
	return timeString
}
