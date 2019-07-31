package cart

import (
	"gotrue/dao"
	"gotrue/dto"
	"gotrue/dto/request"
	"gotrue/dto/response"
	"gotrue/facility/utils"
	"gotrue/model"
	"gotrue/service/goods"
	"gotrue/service/stock"

	"github.com/shopspring/decimal"
)

// CartService 是Cart的一个单例
var CartService *Cart

func beforeInit() {
	goods.Init()
	stock.Init()
}

func initCartService() {
	if CartService != nil {
		return
	}
	beforeInit()
	CartService = &Cart{
		dao:             dao.CartDao,
		stockDao:        dao.StockDao,
		goodsService:    goods.GoodsService,
		goodsImgService: goods.GoodsImgService,
		stockService:    stock.StockService,
	}
}

// Cart 作为规格服务，实现了api.ICart
type Cart struct {
	dao             *dao.Cart
	stockDao        *dao.Stock
	goodsService    *goods.Goods
	goodsImgService *goods.GoodsImg
	stockService    *stock.Stock
}

func (s *Cart) AddGoods(request request.CartAddRequest) (id int64, err error) {
	goods, err := s.goodsService.Info(request.GoodsID) // 商品信息
	if err != nil {
		return
	}
	stock, err := s.stockService.GetByID(request.StockID)
	if err != nil {
		return
	}

	stockSpecIDs, err := utils.ParseIntArray(stock.GoodsSpecificationIDs, "_", 10, 64)
	if err != nil {
		return
	}
	specDesc, err := s.goodsService.SpecificationDesc(request.GoodsID, stockSpecIDs, "/")
	if err != nil {
		return
	}
	cartCreator := newCartCreator(goods, stock, request.UserID, request.Quantity)
	id, err = s.dao.AddCart(cartCreator.create(specDesc))
	if err != nil {
		return
	}
	return
}

func (s *Cart) List(userID int64) ([]response.CartItem, error) {
	cartItems, err := s.dao.SelectByUserID(userID)
	if err != nil {
		return nil, err
	}
	cartSet := model.NewCartSet(cartItems)
	stockIDs := cartSet.StockIDs()
	stocks, err := s.stockDao.SelectByIDs(stockIDs)
	if err != nil {
		return nil, err
	}
	stockMap := model.NewStockSet(stocks).Map()
	wrapper := newCartResponseWrapper(cartItems, stockMap)
	return wrapper.DTOItems(), nil
}

func (s *Cart) GoodsCount(userID int64) (count int64, err error) {
	return s.dao.SelectCountByUserID(userID)
}

func (s *Cart) CheckItem(req request.CartCheckRequest) error {
	setMap := map[string]interface{}{
		"checked": req.IsChecked,
	}
	return s.dao.UpdateByID(req.ID, setMap)
}

func (s *Cart) CheckedItems(userID int64) ([]model.Cart, error) {
	checkedItems, err := s.dao.SelectChecked(userID)
	if err != nil {
		return nil, err
	}
	return checkedItems, nil
}

type cartResposneWrapper struct {
	items    []model.Cart
	stockMap map[int64]*model.GoodsStock
}

func newCartResponseWrapper(items []model.Cart, stockMap map[int64]*model.GoodsStock) *cartResposneWrapper {
	return &cartResposneWrapper{
		items:    items,
		stockMap: stockMap,
	}
}

func (set *cartResposneWrapper) DTOItems() []response.CartItem {
	dtos := make([]response.CartItem, len(set.items))
	for i, model := range set.items {
		dto := response.CartItem{}
		dto.ID = model.ID
		dto.GoodsName = model.GoodsName
		dto.GoodsSN = model.GoodsSN
		dto.GoodsSpecDescription = model.GoodsSpecDescription
		dto.GoodsSpecIDs = model.GoodsSpecIDs
		if stock, ok := set.stockMap[model.StockID]; ok {
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

type cartCreator struct {
	goods    *dto.GoodsInfoDTO
	stock    *dto.GoodsStockDTO
	userID   int64
	quantity decimal.Decimal
}

func newCartCreator(goods *dto.GoodsInfoDTO,
	stock *dto.GoodsStockDTO, userID int64, quantity decimal.Decimal) *cartCreator {
	cc := new(cartCreator)
	cc.goods = goods
	cc.stock = stock
	cc.userID = userID
	cc.quantity = quantity
	return cc
}

func (cc *cartCreator) create(specDesc string) model.Cart {
	return model.Cart{
		UserID:               cc.userID,
		GoodsID:              cc.goods.ID,
		SupplierID:           cc.goods.SupplierID,
		StockID:              cc.stock.ID,
		Quantity:             cc.quantity,
		GoodsName:            cc.goods.Name,
		GoodsSpecDescription: specDesc,
		GoodsSpecIDs:         cc.stock.GoodsSpecificationIDs,
		ListPicURL:           cc.goods.ListPicURL,
	}
}