package api

import "gotrue/dto"

// IGoodsStockService 是商品库存的内部接口
// 方便以后以RPC或者其他方式进行服务拆分和依赖规划
type IGoodsStockService interface {
	GetStocksByGoodsID(goodsID int64) ([]*dto.GoodsStockDTO, error)
}
