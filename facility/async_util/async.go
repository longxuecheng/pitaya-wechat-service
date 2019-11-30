package async_util

import (
	"log"
)

// RunAsyncWithRecovery run a function with rececovery protection
func RunAsyncWithRecovery(f func() error) {
	go func() {
		defer catchPanic()
		err := f()
		if err != nil {
			log.Printf("RunAsync function err %+v\n", err)
		}
	}()
}

func catchPanic() {
	if e := recover(); e != nil {
		log.Printf("RunAsync catch panic %+v", e)
	}
}
