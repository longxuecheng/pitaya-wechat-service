package test

import (
	"fmt"
	"runtime"
	"testing"
)

func TestFrames(t *testing.T) {
	c := func() {
		// 向runtime.Callers请求最多10个pcs，包括runtime.Callers调用本身
		pc := make([]uintptr, 10)
		n := runtime.Callers(0, pc)
		if n == 0 {
			// No pcs available. Stop now.
			// This can happen if the first argument to runtime.Callers is large.
			return
		}

		pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
		frames := runtime.CallersFrames(pc)

		// Loop to get frames.
		// A fixed number of pcs can expand to an indefinite number of Frames.
		for {
			frame, more := frames.Next()
			// To keep this example's output stable
			// even if there are changes in the testing package,
			// stop unwinding when we leave package runtime.
			// if !strings.Contains(frame.File, "runtime/") {
			// 	break
			// }
			fmt.Printf("- more:%v | %s file:%s line:%d entry address:%d\n", more, frame.Function, frame.File, frame.Line, frame.Entry)
			if !more {
				break
			}
		}
	}

	b := func() { c() }
	a := func() { b() }

	a()
}
