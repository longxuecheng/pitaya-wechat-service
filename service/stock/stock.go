package stock

import (
	"gotrue/dao"
	"gotrue/dto/response"
	"gotrue/model"
)

var StockService *Stock

func initStockService() {
	if StockService != nil {
		return
	}
	StockService = &Stock{
		dao: dao.StockDao,
	}
}

// Stock 作为类目服务，实现了api.Stock接口
// 服务依赖 (1. attributeService)
type Stock struct {
	dao *dao.Stock
}

func (s *Stock) GetByID(ID int64) (*response.GoodsStockDTO, error) {
	stock, err := s.dao.SelectByID(ID)
	if err != nil {
		return nil, err
	}
	return installGoodsStockDTO(stock), nil
}

func (s *Stock) GetStocksByGoodsID(goodsID int64) ([]*response.GoodsStockDTO, error) {
	stocks, err := s.dao.SelectByGoodsID(goodsID)
	if err != nil {
		return nil, err
	}
	return buildGoodsStockDTOs(stocks), nil
}

func installGoodsStockDTO(model *model.Stock) *response.GoodsStockDTO {
	dto := new(response.GoodsStockDTO)
	dto.ID = model.ID
	dto.GoodsID = model.GoodsID
	dto.GoodsSpecificationIDs = model.Specification.String
	dto.SaleUnitPrice = model.SaleUnitPrice
	dto.AvailableQuantity = model.AvailableQuantity
	return dto
}

func buildGoodsStockDTOs(models []*model.Stock) []*response.GoodsStockDTO {
	if models == nil || len(models) == 0 {
		return nil
	}
	dtos := make([]*response.GoodsStockDTO, len(models))
	for i, model := range models {
		dtos[i] = installGoodsStockDTO(model)
	}
	return dtos
}
