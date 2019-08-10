package middle_ware

import (
	"log"
	"net/http"
	"net/http/httputil"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				err := err.(error)
				logRequest(c)
				printStackTrace(err)
				response := newResponse(nil, err.Error(), STATUS_ERROR)
				c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			}
		}()
		c.Next()
	}
}

func logRequest(c *gin.Context) {
	httprequest, _ := httputil.DumpRequest(c.Request, false)
	log.Printf("[Recovery] %s request:\n%s", timeFormat(time.Now()), string(httprequest))
}

func printStackTrace(err error) {
	log.Printf("error is %+v\n", err) // verb is +v so stack trace can be printed
	pc := make([]uintptr, 10)
	n := runtime.Callers(4, pc)
	if n == 0 {
		return
	}
	pc = pc[:n]
	frames := runtime.CallersFrames(pc)
	for {
		frame, more := frames.Next()
		log.Printf("| %s file:%s line:%d entry address:%d\n", frame.Function, frame.File, frame.Line, frame.Entry)
		if !more {
			break
		}
	}
}

func timeFormat(t time.Time) string {
	var timeString = t.Format("2006/01/02 - 15:04:05")
	return timeString
}
