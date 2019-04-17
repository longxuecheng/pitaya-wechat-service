package api

import (
	"pitaya-wechat-service/dto/request"
	"pitaya-wechat-service/dto/response"
)

// ICartService 是购物车相关服务
type ICartService interface {
	// AddGoods just add your goods into cart
	AddGoods(request request.CartAddRequest) (cartID int64, err error)
	List(userID int64) ([]response.CartItemDTO, error)
	GoodsCount(userID int64) (count int64, err error)
	CheckItem(req request.CartCheckRequest) error
}
