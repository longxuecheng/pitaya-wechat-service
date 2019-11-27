package service

import (
	"gotrue/service/basic"
	"gotrue/service/cart"
	"gotrue/service/cashier"
	"gotrue/service/goods"
	"gotrue/service/order"
	"gotrue/service/region"
	"gotrue/service/stock"
	"gotrue/service/supplier"
	"gotrue/service/tencloud"
	"gotrue/service/user"
	"gotrue/service/wechat"
)

func Init() {
	wechat.InitWechatService()
	user.Init()
	basic.Init()
	region.Init()
	stock.Init()
	goods.Init()
	cart.Init()
	cashier.Init()
	order.Init()
	supplier.Init()
	InitSettlementService()
	InitBannerService()
	InitGaodeMapService()
	tencloud.InitCosService()
}
