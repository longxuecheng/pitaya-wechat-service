package errors

import (
	"fmt"
	"testing"
)

func TestWrafError(t *testing.T) {
	err := NewWithCodef("ErrorCode1", "test error message %s", "ok")
	if e, ok := err.(*Error); ok {
		fmt.Printf("[ ERROR-Manage ] detail is : %+v\n", e.err)
	}

	err1 := Newf(nil, "test error message %s", "ok")
	fmt.Printf("[ ERROR-Manage ] detail is : %+v\n", err1)

}
