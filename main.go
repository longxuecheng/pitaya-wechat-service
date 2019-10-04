package main

import (
	"gotrue/dao"
	"gotrue/router"
	"gotrue/service"
	_ "net/http/pprof"
	"time"
)

func init() {
	localloc := time.FixedZone("Asia/Beijing", 3600*8)
	time.Local = localloc
	dao.Init()
	service.Init()
}

func main() {
	router.ListenAndServe("8082")
}
