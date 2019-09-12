package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	err := NewWithCodef("ErrorCodeTest", "Testing error")
	if err, ok := err.(error); ok {
		if readable, ok := Readable(err); ok {
			readable.PrintStackTrace()
		}
	} else {
		fmt.Println("Error is not readable")
	}
	fmt.Printf("error %+v\n", err)

	err = CauseWithCodef(errors.New("Raw error"), "Code1", "Test cause")

	fmt.Printf("error %+v\n", err)
}
