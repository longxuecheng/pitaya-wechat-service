package service

import (
	"pitaya-wechat-service/dao"
	"pitaya-wechat-service/dto/request"
	"pitaya-wechat-service/dto/response"
	"pitaya-wechat-service/facility/utils"
	"pitaya-wechat-service/model"
	"strings"
)

// cartServiceSingleton 是CartService的一个单例
var cartServiceSingleton *CartService

func CartServiceInstance() *CartService {
	if cartServiceSingleton == nil {
		cartServiceSingleton = new(CartService)
		cartServiceSingleton.dao = dao.CartDaoSingleton
		cartServiceSingleton.stockDao = dao.GoodsStockDaoSingleton
		cartServiceSingleton.goodsService = GoodsServiceInstance()
		cartServiceSingleton.stockService = StockServiceInstance()
		cartServiceSingleton.goodsImgService = GoodsImgServiceInstance()
	}
	return cartServiceSingleton
}

// CartService 作为规格服务，实现了api.ICartService
type CartService struct {
	dao             *dao.CartDao
	stockDao        *dao.GoodsStockDao
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
		"user_id":                request.UserID,
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
	cartItems, err := s.dao.SelectByUserID(userID)
	if err != nil {
		return nil, err
	}
	cartSet := newCartSet(cartItems)
	stockIDs := cartSet.stockIDs()
	stocks, err := s.stockDao.SelectByIDs(stockIDs)
	cartSet.bindStocks(stocks)
	if err != nil {
		return nil, err
	}
	return cartSet.DTOItems(), nil
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

type cartSet struct {
	items    []model.Cart
	stockMap map[int64]*model.GoodsStock
}

func newCartSet(items []model.Cart) *cartSet {
	return &cartSet{
		items: items,
	}
}

func (set *cartSet) bindStocks(stocks []*model.GoodsStock) {
	set.stockMap = model.NewStockMap(stocks)
}

func (set *cartSet) stockIDs() []int64 {
	stockIDs := []int64{}
	stockIDMap := map[int64]bool{}
	for _, item := range set.items {
		stockIDMap[item.StockID] = true
	}
	for k := range stockIDMap {
		stockIDs = append(stockIDs, k)
	}
	return stockIDs
}

func (set *cartSet) DTOItems() []response.CartItemDTO {
	dtos := make([]response.CartItemDTO, len(set.items))
	for i, model := range set.items {
		dto := response.CartItemDTO{}
		dto.ID = model.ID
		dto.GoodsName = model.GoodsName
		dto.GoodsSN = model.GoodsSN
		dto.GoodsSpecDescription = model.GoodsSpecDescription
		dto.GoodsSpecIDs = model.GoodsSpecIDs
		if stock, ok := set.stockMap[model.ID]; ok {
			dto.MarketPrice = stock.SaleUnitPrice
			dto.RetailPrice = stock.SaleUnitPrice
		}
		dto.Quantity = model.Quantity
		dto.StockID = model.StockID
		dto.GoodsID = model.GoodsID
		dto.GoodsSpecIDs = model.GoodsSpecIDs
		dto.ListPicURL = model.ListPicURL
		dto.SessionID = model.SessionID
		dto.Checked = model.Checked
		dtos[i] = dto
	}
	return dtos
}
