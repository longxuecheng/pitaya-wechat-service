package main

import (
	"gotrue/dao"
	"gotrue/router"
	_ "net/http/pprof"
)

func init() {
	dao.Init()
}

func main() {
	router.ListenAndServe("8081")
}
