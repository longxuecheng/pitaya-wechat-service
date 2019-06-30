package main

import (
	"gotrue/router"
	_ "net/http/pprof"
)

func main() {
	router.ListenAndServe("8081")
}
