package main

import (
	"gotrue/dao"
	"gotrue/facility/utils"
	"gotrue/router"
	"gotrue/service"
	_ "net/http/pprof"
	"time"
)

func init() {
	localloc := time.FixedZone("Asia/Beijing", 3600*8)
	time.Local = localloc
	utils.InitEncryptor()
	dao.Init()
	service.Init()
}

func main() {
	router.ListenAndServe("8082")
}
