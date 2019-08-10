package test

import (
	"errors"
	"fmt"
	"runtime"
	"testing"
)

func TestFrames(t *testing.T) {
	// FuncA()
	fmt.Printf("error is %v\n", FuncError())
}

func FuncError() error {
	return errors.New("calling function error ")
}

func FuncB() {
	fmt.Println("function B")
	FuncC()
}

func FuncA() {
	fmt.Println("function A")
	FuncB()
}

func FuncC() {
	// 向runtime.Callers请求最多10个pcs，包括runtime.Callers调用本身
	pc := make([]uintptr, 10)
	// 填充调用goroutine栈的函数调用程序计数器。skip参数用来指定在pc中存储栈帧之前跳过多少个栈帧
	// 如果是0代表调用函数本身的栈帧，如果是1代表当前执行函数的调用函数
	n := runtime.Callers(1, pc)
	if n == 0 {
		// No pcs available. Stop now.
		// This can happen if the first argument to runtime.Callers is large.
		return
	}

	pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
	// 使用Callers返回的程序计数器列表并且准备返回形如function/file/line的信息
	// 在处理完栈帧之前不要修改pc
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

/*
function A
function B
- more:true | runtime.Callers file:/usr/local/go/src/runtime/extern.go line:208 entry address:16820464
- more:true | gotrue/test.FuncC file:/Users/lxc/GoDevelop/src/gotrue/test/runtime_test.go line:28 entry address:22169264
- more:true | gotrue/test.FuncB file:/Users/lxc/GoDevelop/src/gotrue/test/runtime_test.go line:15 entry address:22168944
- more:true | gotrue/test.FuncA file:/Users/lxc/GoDevelop/src/gotrue/test/runtime_test.go line:20 entry address:22169104
- more:true | gotrue/test.TestFrames file:/Users/lxc/GoDevelop/src/gotrue/test/runtime_test.go line:10 entry address:22168896
- more:true | testing.tRunner file:/usr/local/go/src/testing/testing.go line:865 entry address:18046384
- more:false | runtime.goexit file:/usr/local/go/src/runtime/asm_amd64.s line:1337 entry address:17193808
*/
