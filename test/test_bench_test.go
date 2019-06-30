package test

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strconv"
	"testing"
	"text/template"
	"time"
)

func TestFoo(t *testing.T) {
	// <setup code>
	t.Run("A=1", func(t *testing.T) {
		t.Log("TestFoo sub test A=1")
	})
	t.Run("A=2", func(t *testing.T) {
		t.Log("TestFoo sub test A=2")
	})
	t.Run("B=1", func(t *testing.T) {
		t.Log("TestFoo sub test B=1")
	})
	// <tear-down code>
}

func TestFlag(t *testing.T) {
	println("Flag name from command-line is ", *name)
}

// Test parallism
func TestGroupedParallel(t *testing.T) {
	name := flag.String("name", "default name ", "[usage] -name to pass the params")
	println(name)
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()
	type x struct {
		Name string
	}
	tests := make([]x, 10)
	for i := 0; i < 10; i++ {
		tests[i] = x{
			"TEST" + strconv.FormatInt(int64(i), 10),
		}
	}
	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			t.Log(tc.Name)
		})
	}
}

func TestTeardownParallel(t *testing.T) {
	parallelTest1 := func(t *testing.T) {
		t.Log("parallelTest1")
	}
	parallelTest2 := func(t *testing.T) {
		t.Log("parallelTest2")
	}
	parallelTest3 := func(t *testing.T) {
		t.Log("parallelTest3")
	}
	// This Run will not return until the parallel tests finish.
	t.Run("group", func(t *testing.T) {
		t.Run("Test1", parallelTest1)
		t.Run("Test2", parallelTest2)
		t.Run("Test3", parallelTest3)
	})
	t.Log("Now runing tear down functions...")
}

func BenchmarkHello(b *testing.B) {
	mockLongTimeOp()
	b.ResetTimer() // 忽略长时间的初始化操作重新开始计时
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("hello")
	}
}

func mockLongTimeOp() {
	time.Sleep(1000 * time.Millisecond)
}

func BenchmarkTemplateParallel(b *testing.B) {
	templ := template.Must(template.New("test").Parse("Hello, {{.}}!"))
	b.RunParallel(func(pb *testing.PB) {
		var buf bytes.Buffer
		for pb.Next() {
			buf.Reset()
			templ.Execute(&buf, "World")
		}
	})
}

// Examples

func ExampleHello() {
	fmt.Println("hello")
	// Output: hello1
}

func ExampleSalutations() {
	fmt.Println("hello, and")
	fmt.Println("goodbye")
	// Output:
	// hello, and
	// goodbye
}

func ExamplePerm() {
	for _, value := range []int{0, 1, 2, 3, 4} {
		fmt.Println(value)
	}
	// Unordered output: 4
	// 2
	// 1
	// 3
	// 0
}

var name *string

// TestMain是主测试方法会在所有的测试方法之前运行
func TestMain(m *testing.M) {
	name = flag.String("name", "default name ", "[usage] -name to pass the params")
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()
	println("Flag got by TestMain is ", *name)
	os.Exit(m.Run())
}
