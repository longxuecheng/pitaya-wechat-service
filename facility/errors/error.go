package errors

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

type Error struct {
	cause   error
	code    string
	message string
	pcs     []uintptr
}

func (e *Error) Error() string {
	msg := ""
	if e.cause == nil {
		msg = fmt.Sprintf("Code is [ %s ] Message is [ %s ]", e.code, e.message)
	} else {
		msg = e.cause.Error()
	}
	return e.stackTrace() + msg
}

func (e *Error) Code() string {
	return e.code
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) Cause() error {
	return e.cause
}

func (e *Error) stackTrace() string {
	sb := strings.Builder{}
	defer sb.Reset()
	frames := runtime.CallersFrames(e.pcs)
	for {
		frame, more := frames.Next()
		sb.WriteString(fmt.Sprintf("| %s file:%s line:%d entry address:%d\n", frame.Function, frame.File, frame.Line, frame.Entry))
		if !more {
			break
		}
	}
	return sb.String()
}

func (e *Error) PrintStackTrace() {
	frames := runtime.CallersFrames(e.pcs)
	for {
		frame, more := frames.Next()
		log.Printf("| %s file:%s line:%d entry address:%d\n", frame.Function, frame.File, frame.Line, frame.Entry)
		if !more {
			break
		}
	}
}

type readable interface {
	Code() string
	Cause() error
	Message() string
	Error() string
	PrintStackTrace()
}

func Readable(err error) (readable, bool) {
	r, ok := err.(readable)
	return r, ok
}

// NewWithCodef use error.Error as result, will contains error code
// to distinguish differences between errors
func NewWithCodef(code, format string, args ...interface{}) error {

	return &Error{
		code:    code,
		message: fmt.Sprintf(format, args...),
		pcs:     stack(),
	}
}

func CauseWithCodef(cause error, code, format string, args ...interface{}) error {
	return &Error{
		cause:   cause,
		code:    code,
		message: fmt.Sprintf(format, args...),
		pcs:     stack(),
	}
}

func stack() []uintptr {
	pc := make([]uintptr, 10)
	n := runtime.Callers(2, pc)
	pc = pc[:n]
	return pc
}
