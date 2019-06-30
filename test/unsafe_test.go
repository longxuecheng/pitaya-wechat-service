package test

import (
	"math"
	"testing"
	"unsafe"
)

type s struct {
	A string
	B string
}

type k struct {
	A string
	B string
}

func TestUnsafe(t *testing.T) {
	var f float64 = 1000000 // +1.000000e+006
	println(f)
	i := math.Float64bits(f)
	println(i) // 4696837146684686336
	s1 := s{
		"A",
		"B",
	}
	p := unsafe.Pointer(&s1.A)
	println("up of s1.A is ", p)
	println("uintptr of s1 is ", uintptr(p))
	p1 := unsafe.Pointer(uintptr(unsafe.Pointer(&s1)) + unsafe.Offsetof(s1.A))
	println("up of s1 + s1.A is ", p1)
	println("off set of s1.A is  ", unsafe.Offsetof(s1.A))
	println("off set of s1.B is  ", unsafe.Offsetof(s1.B))

}

func TestUnsafeConverting(t *testing.T) {
	// This test can convert type s to k 绕过类型的校验
	s1 := s{
		"A",
		"B",
	}
	t.Log("s1 is ", s1)
	p := unsafe.Pointer(&s1.A)
	k1 := *(*k)(p)
	t.Log("k1 is ", k1)
	// if k1, ok := s1.(k); ok {
	// 	t.Log("assert s to k is ", k1)
	// } else {
	// 	t.Log("can not convert type s to k")
	// }

}
