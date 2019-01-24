package service

import (
	"pitaya-wechat-service/dao"
	"pitaya-wechat-service/dto"
	"pitaya-wechat-service/model"
)

var GoodsStockServiceSingleton *GoodsStockService

// init 在此实现spring中类似注入的功能
func init() {
	GoodsStockServiceSingleton = new(GoodsStockService)
	GoodsStockServiceSingleton.stockDao = dao.GoodsStockDaoSingleton
}

func StockServiceInstance() *GoodsStockService {
	if GoodsStockServiceSingleton != nil {
		return GoodsStockServiceSingleton
	}
	GoodsStockServiceSingleton = new(GoodsStockService)
	GoodsStockServiceSingleton.stockDao = dao.GoodsStockDaoSingleton
	return GoodsStockServiceSingleton
}

// GoodsStockService 作为类目服务，实现了api.GoodsStockService接口
// 服务依赖 (1. attributeService)
type GoodsStockService struct {
	stockDao *dao.GoodsStockDao
}

func (s *GoodsStockService) GetByID(ID int64) (*dto.GoodsStockDTO, error) {
	stock, err := s.stockDao.SelectByID(ID)
	if err != nil {
		return nil, err
	}
	return installGoodsStockDTO(stock), nil
}

func (s *GoodsStockService) GetStocksByGoodsID(goodsID int64) ([]*dto.GoodsStockDTO, error) {
	stocks, err := s.stockDao.SelectByGoodsID(goodsID)
	if err != nil {
		return nil, err
	}
	return buildGoodsStockDTOs(stocks), nil
}

func installGoodsStockDTO(model *model.GoodsStock) *dto.GoodsStockDTO {
	dto := new(dto.GoodsStockDTO)
	dto.ID = model.ID
	dto.GoodsID = model.GoodsID
	dto.GoodsSpecificationIDs = model.Specification.String
	dto.SaleUnitPrice = model.SaleUnitPrice
	dto.AvailableQuantity = model.AvailableQuantity
	return dto
}

func buildGoodsStockDTOs(models []*model.GoodsStock) []*dto.GoodsStockDTO {
	if models == nil || len(models) == 0 {
		return nil
	}
	dtos := make([]*dto.GoodsStockDTO, len(models))
	for i, model := range models {
		dtos[i] = installGoodsStockDTO(model)
	}
	return dtos
}
