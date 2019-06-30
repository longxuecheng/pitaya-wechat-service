package test

import (
	"fmt"
	"runtime/pprof"
	"testing"
	"time"
)

func TestProfiles(t *testing.T) {
	pfs := pprof.Profiles()
	for _, pf := range pfs {
		fmt.Printf("- pprof profile: %s \n", pf.Name())
	}
}

func TestLookUp(t *testing.T) {
	a := func() {
		pf := pprof.Lookup("goroutine")
		fmt.Printf("- pporf goroutine profile is %v \n", pf.Count())
	}
	for i := 0; i < 1000; i++ {
		go a()
	}
	time.Sleep(time.Second)
}
