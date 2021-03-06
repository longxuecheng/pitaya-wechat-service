package cashier

import (
	"database/sql"
	"gotrue/dao"
	"gotrue/dto/request"
	"gotrue/dto/response"
	"gotrue/facility/errors"
	"gotrue/facility/utils"
	"gotrue/model"
	"gotrue/service/api"
	"gotrue/service/cart"
	"gotrue/service/cut"
	"gotrue/service/goods"
	"gotrue/service/user"

	"github.com/shopspring/decimal"
)

const (
	specSep     = "_"
	specDescSep = "/"
)

var one = decimal.NewFromFloat32(1.0)

var CashierService *Cashier

func beforeInit() {
	goods.Init()
	cart.Init()
}

func initCashierService() {
	if CashierService != nil {
		return
	}
	beforeInit()
	CashierService = &Cashier{
		stockDao:     dao.StockDao,
		goodsDao:     dao.GoodsDao,
		goodsService: goods.GoodsService,
		cartService:  cart.CartService,
		userService:  user.UserService,
		cutService:   cut.GetCutService(),
	}
}

// Cashier is checkout service
type Cashier struct {
	stockDao     *dao.Stock
	goodsDao     *dao.Goods
	goodsService api.IGoodsService
	cartService  *cart.Cart
	userService  *user.User
	cutService   api.ICutService
}

// CartCheckout is 从购物车结算
func (s *Cashier) CartCheckout(userID int64) (*response.Cashier, error) {
	items, err := s.cartService.List(userID)
	if err != nil {
		return nil, err
	}
	cc := newCartCashier(items)
	return cc.summary(), nil
}

// StockCheckout is 从库存进行快速结算
func (s *Cashier) StockCheckout(req *request.CashierPreview) (*response.Cashier, error) {
	stock, err := s.stockDao.SelectByID(req.StockID)
	if err != nil {
		return nil, err
	}
	goods, err := s.goodsDao.SelectOnSaleByID(stock.GoodsID)
	if err == sql.ErrNoRows {
		return nil, errors.NewWithCodef("GoodsNotAvailable", "该商品当前不可售")
	}
	if err != nil {
		return nil, err
	}
	stockSpecIDs, err := utils.ParseIntArray(stock.Specification.String, specSep, 10, 64)
	if err != nil {
		return nil, err
	}
	specDesc, err := s.goodsService.SpecificationDesc(stock.GoodsID, stockSpecIDs, specDescSep)
	if err != nil {
		return nil, err
	}
	express, err := s.goodsService.ExpressConstraint(req.StockID, req.AddressID)
	if err != nil {
		return nil, err
	}
	cutReq := &request.CutOrder{
		GoodsID: goods.ID,
		UserID:  req.UserID,
	}
	cutOrder, err := s.cutService.MyActivatedCutOrder(cutReq)
	if err != nil {
		return nil, err
	}
	gc := newGoodsCashier(stock, goods, specDesc, req.Quantity)
	gc.express = express
	gc.setCutoff(cutOrder)
	cs := gc.summary()
	return cs, nil
}

type cashier interface {
	summary() *response.Cashier
}

// goodsCashier 是商品结算台
type goodsCashier struct {
	quantity decimal.Decimal
	stock    *model.Stock
	goods    *model.Goods
	express  *response.GoodsExpressConstraint
	cutoff   *response.CutOrder
	specText string
}

func newGoodsCashier(stock *model.Stock, goods *model.Goods, specDesc string, quantity decimal.Decimal) *goodsCashier {
	return &goodsCashier{
		stock:    stock,
		goods:    goods,
		quantity: quantity,
		specText: specDesc,
	}
}

func (gc *goodsCashier) setCutoff(cutoff *response.CutOrder) {
	gc.cutoff = cutoff
}

func (gc *goodsCashier) summary() *response.Cashier {
	goodsTotalPrice := gc.stock.SaleUnitPrice.Mul(gc.quantity)
	orderTotalPrice := goodsTotalPrice
	if gc.cutoff != nil {
		orderTotalPrice = goodsTotalPrice.Sub(gc.cutoff.CutoffPrice)
	}
	// total express fee = (unit expresss fee) * quantity
	if !gc.stock.Splitable {
		gc.express.CalculateTotalExpressFee(one)
	} else {
		if gc.express != nil {
			gc.express.CalculateTotalExpressFee(gc.quantity)
		}
	}
	cc := &response.Cashier{
		GoodsTotalPrice:        goodsTotalPrice.StringFixed(2),
		OrderTotalPrice:        orderTotalPrice.Add(gc.express.TotalExpressFee).StringFixed(2),
		GoodsExpressConstraint: gc.express,
		CutOrder:               gc.cutoff,
	}

	ci := response.CashierItem{
		GoodsID:              gc.goods.ID,
		GoodsName:            gc.goods.Name,
		GoodsSpecDescription: gc.specText,
		GoodsSpecIDs:         gc.stock.Specification.String,
		ListPicURL:           gc.goods.ListPicURL.String,
		StockID:              gc.stock.ID,
		RetailPrice:          gc.stock.SaleUnitPrice,
		Quantity:             gc.quantity,
	}
	cc.Items = []response.CashierItem{
		ci,
	}

	return cc
}

// cartCashier购物车结算台
type cartCashier struct {
	items []response.CartItem
}

func newCartCashier(items []response.CartItem) *cartCashier {
	filter := func(carts []response.CartItem, filterFunc func(input response.CartItem) bool) []response.CartItem {
		results := []response.CartItem{}
		for _, cart := range carts {
			if filterFunc(cart) {
				results = append(results, cart)
			}
		}
		return results
	}
	cc := &cartCashier{
		items: filter(items, func(input response.CartItem) bool {
			return input.Checked == 1
		}),
	}
	return cc
}

func (cc *cartCashier) summary() *response.Cashier {
	checkedGoodsCount := decimal.Zero
	checkedGoodsPrice := decimal.Zero
	cashierItems := make([]response.CashierItem, len(cc.items))
	for i, cart := range cc.items {
		checkedGoodsCount = checkedGoodsCount.Add(cart.Quantity)
		checkedGoodsPrice = checkedGoodsPrice.Add(cart.RetailPrice.Mul(cart.Quantity))
		ci := response.CashierItem{}
		ci.GoodsID = cart.GoodsID
		ci.GoodsName = cart.GoodsName
		ci.GoodsSpecDescription = cart.GoodsSpecDescription
		ci.GoodsSpecIDs = cart.GoodsSpecIDs
		ci.ListPicURL = cart.ListPicURL
		ci.Quantity = cart.Quantity
		ci.RetailPrice = cart.RetailPrice
		ci.StockID = cart.StockID
		cashierItems[i] = ci
	}
	cashier := &response.Cashier{
		GoodsTotalPrice: checkedGoodsPrice.StringFixed(2),
		// OrderTotalPrice: checkedGoodsPrice.Add(defaultExpressFee).StringFixed(2),
		Items: cashierItems,
	}
	return cashier
}
