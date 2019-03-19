package main

import (
	_ "net/http/pprof"
	"pitaya-wechat-service/router"
)

func main() {
	router.ListenAndServe("8081")
}
