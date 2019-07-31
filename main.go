package main

import (
	"gotrue/dao"
	"gotrue/router"
	"gotrue/service"
	_ "net/http/pprof"
)

func init() {
	dao.Init()
	service.Init()
}

func main() {
	router.ListenAndServe("8082")
}
