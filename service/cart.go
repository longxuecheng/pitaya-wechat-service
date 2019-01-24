package service

import (
	"pitaya-wechat-service/dao"
	"pitaya-wechat-service/dto/request"
	"pitaya-wechat-service/dto/response"
	"pitaya-wechat-service/facility/utils"
	"pitaya-wechat-service/model"
	"strings"
)

// CartServiceSingleton 是CartService的一个单例
var CartServiceSingleton *CartService

// init 在此实现spring中类似注入的功能
func init() {
	CartServiceSingleton = new(CartService)
	CartServiceSingleton.dao = dao.CartDaoSingleton
	CartServiceSingleton.goodsService = GoodsServiceInstance()
	CartServiceSingleton.stockService = StockServiceInstance()
	CartServiceSingleton.goodsImgService = GoodsImgServiceInstance()
}

// CartService 作为规格服务，实现了api.ICartService
type CartService struct {
	dao             *dao.CartDao
	goodsService    *GoodsService
	goodsImgService *GoodsImgService
	stockService    *GoodsStockService
}

func (s *CartService) AddGoods(request request.CartAddRequest) (id int64, err error) {
	goods, err := s.goodsService.Info(request.GoodsID) // 商品信息
	if err != nil {
		return
	}
	specs, err := s.goodsService.Specifications(request.GoodsID) // 商品规格
	if err != nil {
		return
	}
	gallery, err := s.goodsImgService.GetByGoodsID(request.GoodsID)
	if err != nil {
		return
	}
	var listPic = "" // 列表图片
	if len(gallery) > 0 {
		listPic = gallery[0].ImgURL
	}
	stock, err := s.stockService.GetByID(request.StockID)
	if err != nil {
		return
	}
	specNames := []string{}
	stockSpecIDs, err := utils.ParseIntArray(stock.GoodsSpecificationIDs, "_", 10, 64)
	if err != nil {
		return
	}
	for _, stockSpecID := range stockSpecIDs {
		for _, spec := range specs {
			if stockSpecID == spec.ID {
				specNames = append(specNames, spec.Value)
				break
			}
		}
	}

	specDesc := strings.Join(specNames, ";") // 商品规格组合描述

	setMap := map[string]interface{}{
		"goods_id":               goods.ID,
		"stock_id":               request.StockID,
		"quantity":               request.Quantity,
		"goods_name":             goods.Name,
		"market_price":           goods.RetailPrice,
		"sale_unit_price":        goods.RetailPrice,
		"goods_spec_description": specDesc,
		"goods_spec_ids":         stock.GoodsSpecificationIDs,
		"list_pic_url":           listPic,
	}
	id, err = s.dao.AddCart(setMap)
	if err != nil {
		return
	}
	return
}

func (s *CartService) ListCart4User(userID int64) ([]response.CartItemDTO, error) {
	carts, err := s.dao.SelectByUserID(userID)
	if err != nil {
		return nil, err
	}
	return buildCartItemDTOs(carts), nil
}

func (s *CartService) GoodsCount(userID int64) (count int64, err error) {
	return s.dao.SelectCountByUserID(userID)
}

func (s *CartService) CheckItem(req request.CartCheckRequest) error {
	setMap := map[string]interface{}{
		"checked": req.IsChecked,
	}
	return s.dao.UpdateByID(req.ID, setMap)
}

func installCartItemDTO(model model.Cart) response.CartItemDTO {
	dto := response.CartItemDTO{}
	dto.ID = model.ID
	dto.GoodsName = model.GoodsName
	dto.GoodsSN = model.GoodsSN
	dto.GoodsSpecDescription = model.GoodsSpecDescription
	dto.GoodsSpecIDs = model.GoodsSpecIDs
	dto.MarketPrice = model.MarketPrice
	dto.RetailPrice = model.RetailPrice
	dto.Quantity = model.Quantity
	dto.StockID = model.StockID
	dto.GoodsID = model.GoodsID
	dto.GoodsSpecIDs = model.GoodsSpecIDs
	dto.ListPicURL = model.ListPicURL
	dto.SessionID = model.SessionID
	dto.Checked = model.Checked
	return dto
}

func buildCartItemDTOs(models []model.Cart) []response.CartItemDTO {
	if models == nil || len(models) == 0 {
		return nil
	}
	dtos := make([]response.CartItemDTO, len(models))
	for i, model := range models {
		dtos[i] = installCartItemDTO(model)
	}
	return dtos
}
