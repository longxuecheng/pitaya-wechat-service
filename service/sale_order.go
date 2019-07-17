package service

import (
	"database/sql"
	"errors"
	"gotrue/dao"
	"gotrue/dto"
	"gotrue/dto/pagination"
	"gotrue/dto/request"
	"gotrue/dto/response"
	"gotrue/model"
	"gotrue/sys"
	"strconv"

	"github.com/bwmarrin/snowflake"
	"github.com/shopspring/decimal"
)

var saleOrderServiceSingleton *SaleOrderService

// SaleOrderServiceInstance get a service instance of singleton
func SaleOrderServiceInstance() *SaleOrderService {
	if saleOrderServiceSingleton == nil {
		saleOrderServiceSingleton = new(SaleOrderService)
		saleOrderServiceSingleton.dao = dao.SaleOrderDaoInstance()
		saleOrderServiceSingleton.stockDao = dao.GoodsStockDaoSingleton
		saleOrderServiceSingleton.goodsDao = dao.GoodsDaoSingleton
		saleOrderServiceSingleton.saleDetailDao = dao.SaleDetailDaoInstance()
		saleOrderServiceSingleton.cartService = CartServiceInstance()
		saleOrderServiceSingleton.userService = UserServiceInstance()
	}
	return saleOrderServiceSingleton
}

// SaleOrderService 作为销售订单服务，实现了api.IOrderService
type SaleOrderService struct {
	dao           *dao.SaleOrderDao
	saleDetailDao *dao.SaleDetailDao
	stockDao      *dao.GoodsStockDao
	goodsDao      *dao.GoodsDao
	cartService   *CartService
	userService   *UserService
}

// QuickCreate 快速下单
func (s *SaleOrderService) QuickCreate(req request.SaleOrderQuickAddRequest) (int64, error) {
	stock, err := s.stockDao.SelectByID(req.StockID)
	if err != nil {
		return 0, err
	}
	address, err := s.userService.GetAddressByID(req.AddressID)
	if err != nil {
		return 0, err
	}
	goods, err := s.goodsDao.SelectByID(stock.GoodsID)
	if err != nil {
		return 0, err
	}
	var orderID int64
	// 由于目前水果的包装都是按照一个快递单来进行的，所以目前对每一个项目创建一个订单
	sys.GetEasyDB().ExecTx(func(tx *sql.Tx) error {
		ss := []supplierStock{
			supplierStock{
				SupplierID: goods.SupplierID,
				Quantity:   req.Quantity,
				Stock:      stock,
				Goods:      goods,
			},
		}
		so := &supplierOrder{
			supplierID:     goods.SupplierID,
			supplierStocks: ss,
		}
		so.bindBasically(req.UserID, address)

		if so.splitable() {
			splittedSupplierOrders := so.split()
			for _, splittedSupplierOrder := range splittedSupplierOrders {
				err := s.save(splittedSupplierOrder, tx)
				if err != nil {
					return err
				}
			}
		} else {
			err := s.save(so, tx)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return orderID, nil
}

// save 将业务模型转换成数据库模型并持久化
func (s *SaleOrderService) save(so *supplierOrder, tx *sql.Tx) error {
	saleOrder, saleDetails, err := so.transfer()
	if err != nil {
		return err
	}
	orderID, err := s.dao.Create(saleOrder, tx)
	if err != nil {
		return err
	}
	for _, detail := range saleDetails {
		detail.OrderID = orderID
		_, err := s.saleDetailDao.Create(detail)
		if err != nil {
			return err
		}
	}
	return nil
}

// Create 从购物车创建订单
// 1. 创建订单
// 2. 创建订单明细
// 3. 删除购物车所选中的项目
func (s *SaleOrderService) Create(userID int64, req request.SaleOrderAddRequest) (id int64, err error) {
	checkedItems, err := s.cartService.checkedItems(userID)
	if err != nil {
		return
	}
	address, err := s.userService.GetAddressByID(req.AddressID)
	if err != nil {
		return
	}
	cartSet := model.NewCartSet(checkedItems)
	goodsIDs := cartSet.GoodsIDs()
	goodsList, err := s.goodsDao.SelectByIDs(goodsIDs)
	if err != nil {
		return
	}
	stocks, err := s.stockDao.SelectByIDs(cartSet.StockIDs())
	if err != nil {
		return
	}
	orderCreator := newSaleOrderCreator(cartSet)
	orderCreator.setGoods(goodsList)
	orderCreator.bindNecessary(stocks, userID, address)
	supplierOrders := orderCreator.rawSupplierOrders()
	sys.GetEasyDB().ExecTx(func(tx *sql.Tx) error {
		for _, so := range supplierOrders {
			if so.splitable() {
				splittedSupplierOrders := so.split()
				for _, splittedSupplierOrder := range splittedSupplierOrders {
					err := s.save(splittedSupplierOrder, tx)
					if err != nil {
						return err
					}
				}
			} else {
				err := s.save(so, tx)
				if err != nil {
					return err
				}
			}

		}
		return nil
	})
	return
}

func (s *SaleOrderService) List(userID int64, req pagination.PaginationRequest) (page pagination.PaginationResonse, err error) {
	page = pagination.PaginationResonse{
		PaginationRequest: req,
	}
	orders, total, err := s.dao.SelectByUserIDWitPagination(userID, req.Offet(), req.Limit())
	if err != nil {
		return
	}
	orderSet := newSaleOrderSet(orders)
	details, err := s.saleDetailDao.SelectByOrderIDs(orderSet.orderIDList()...)
	if err != nil {
		return
	}
	orderSet.setSaleDetails(details)
	page.SetCount(total)
	page.Data = orderSet.orderDTOs()
	return
}

func (s *SaleOrderService) Info(orderID int64) (response.SaleOrderInfoDTO, error) {
	orderInfo := response.SaleOrderInfoDTO{}
	saleOrder, err := s.dao.SelectByID(orderID)
	if err != nil {
		return orderInfo, err
	}
	if err != nil {
		return orderInfo, nil
	}
	return s.installSaleInfoDTO(saleOrder), nil
}

func (s *SaleOrderService) ListGoods(orderID int64) ([]response.SaleOrderGoodsDTO, error) {
	goodsList, err := s.saleDetailDao.SelectByOrderID(orderID)
	if err != nil {
		return nil, err
	}
	dtos := make([]response.SaleOrderGoodsDTO, len(goodsList))
	for i, goods := range goodsList {
		dtos[i] = installSaleDetailDTO(goods)
	}
	return dtos, nil
}

func generateOrderNumber(nodeNo int64) (string, error) {
	// Create a new Node with a Node number of 1
	node, err := snowflake.NewNode(nodeNo)
	if err != nil {
		return "", err
	}
	// Generate a snowflake ID.
	id := node.Generate()
	return strconv.FormatInt(id.Int64(), 10), nil
}

func (s *SaleOrderService) installSaleInfoDTO(order model.SaleOrder) response.SaleOrderInfoDTO {
	dto := response.SaleOrderInfoDTO{}
	dto.ID = order.ID
	dto.OrderNo = order.OrderNo.String
	dto.Status = order.Status
	dto.CreatedAt = order.CreateTime.Format("2006-01-02 15:04:05")
	dto.Consignee = order.Receiver
	dto.Mobile = order.PhoneNo
	dto.FullRegion = "TODO"
	dto.Address = order.Address
	dto.GoodsAmt = order.GoodsAmt
	dto.ExpressFee = order.ExpressFee
	dto.OrderAmt = order.OrderAmt
	return dto
}

func installSaleOrderItemDTO(model model.SaleOrder) response.SaleOrderItemDTO {
	dto := response.SaleOrderItemDTO{}
	dto.ID = model.ID
	dto.OrderNo = model.OrderNo.String
	dto.OrderAmt = model.OrderAmt
	dto.Status = model.Status
	return dto
}

func buildSaleOrderItemDTOs(models []model.SaleOrder) []response.SaleOrderItemDTO {
	if models == nil || len(models) == 0 {
		return nil
	}
	dtos := make([]response.SaleOrderItemDTO, len(models))
	for i, model := range models {
		dtos[i] = installSaleOrderItemDTO(model)
	}
	return dtos
}

func installSaleDetailDTO(model model.SaleDetail) response.SaleOrderGoodsDTO {
	dto := response.SaleOrderGoodsDTO{}
	dto.ID = model.ID
	dto.GoodsName = model.GoodsName
	dto.Quantity = model.Quantity
	dto.ListPicURL = model.ListPicURL.String
	return dto
}

// installSaleDetailFromStock 从库存创建一个订单明细,且数量为1个用于目前的供应商
func installSaleDetailFromStock(orderID int64, stock *model.GoodsStock, goods *model.Goods) model.SaleDetail {
	sd := model.SaleDetail{
		OrderID:       orderID,
		GoodsID:       stock.GoodsID,
		GoodsName:     goods.Name,
		ListPicURL:    goods.ListPicURL,
		Quantity:      decimal.NewFromFloat32(1.0),
		StockID:       stock.ID,
		SaleUnitPrice: stock.SaleUnitPrice,
		CostUnitPrice: stock.CostUnitPrice,
		GoodsSpecIDs:  stock.Specification.String,
	}
	return sd
}

// supplierStock 是供应商的购买库存信息
type supplierStock struct {
	SupplierID int64
	Quantity   decimal.Decimal
	Stock      *model.GoodsStock
	Goods      *model.Goods
}

func (ss supplierStock) saleDetail() model.SaleDetail {
	return model.SaleDetail{
		OrderID:       0,
		GoodsID:       ss.Stock.GoodsID,
		GoodsName:     ss.Goods.Name,
		Quantity:      ss.Quantity,
		StockID:       ss.Stock.ID,
		SaleUnitPrice: ss.Stock.SaleUnitPrice,
		GoodsSpecIDs:  ss.Stock.Specification.String,
		ListPicURL:    ss.Goods.ListPicURL,
	}

}

// supplierOrder 是供应商的订单信息
type supplierOrder struct {
	supplierID     int64
	userID         int64
	address        *dto.UserAddress
	supplierStocks []supplierStock
}

func (so *supplierOrder) bindBasically(userID int64, address *dto.UserAddress) {
	so.userID = userID
	so.address = address
}

func (so *supplierOrder) transfer() (model.SaleOrder, []model.SaleDetail, error) {
	if so.userID == 0 || so.address.ID == 0 {
		return model.SaleOrder{}, nil, errors.New("user id and address must be bound first")
	}
	orderNo, err := generateOrderNumber(1)
	if err != nil {
		return model.SaleOrder{}, nil, err
	}
	saleOrder := model.SaleOrder{
		OrderNo:    sql.NullString{Valid: true, String: orderNo},
		UserID:     so.userID,
		ProvinceID: so.address.ProvinceID,
		CityID:     so.address.CityID,
		DistrictID: so.address.DistrictID,
		Address:    so.address.Address,
		Receiver:   so.address.Name,
		PhoneNo:    so.address.Mobile,
	}
	saleOrder.OrderAmt = decimal.Zero
	saleOrder.GoodsAmt = decimal.Zero
	saleDetails := []model.SaleDetail{}
	for _, ss := range so.supplierStocks {
		sum := ss.Stock.SaleUnitPrice.Mul(ss.Quantity)
		saleOrder.GoodsAmt = saleOrder.GoodsAmt.Add(sum)
		saleDetail := ss.saleDetail()
		saleDetails = append(saleDetails, saleDetail)
	}
	saleOrder.OrderAmt = saleOrder.GoodsAmt
	return saleOrder, saleDetails, nil
}

// split 按照库存单位1对一个供应商的原始订单进行拆分
func (so *supplierOrder) split() []*supplierOrder {
	sos := []*supplierOrder{}
	for _, ss := range so.supplierStocks {
		for i := 0; i < int(ss.Quantity.IntPart()); i++ {
			ssList := []supplierStock{
				supplierStock{
					SupplierID: so.supplierID,
					Quantity:   decimal.NewFromFloat32(1.0),
					Stock:      ss.Stock,
					Goods:      ss.Goods,
				},
			}
			sorder := &supplierOrder{
				supplierID:     so.supplierID,
				supplierStocks: ssList,
			}
			sorder.bindBasically(so.userID, so.address)
			sos = append(sos, sorder)
		}
	}
	return sos
}

// splitable 判断订单是否可以拆分，目前只有供应商ID为1时可以拆
func (so *supplierOrder) splitable() bool {
	return so.supplierID == 1
}

// SaleOrderCreator 是订单服务中特有的购物车管理
// 若进行为服务拆分，那么购物车可能会作为一个单独的服务提供数据，
// 该对象对返回的购物车数据进行业务处理
type saleOrderCreator struct {
	goodsMap map[int64]*model.Goods
	// 供应商购物车项目
	supplierCart map[int64][]model.Cart
	// 供应商订单
	supplierOrders []*supplierOrder
	stocks         map[int64]*model.GoodsStock
}

func newSaleOrderCreator(cartSet *model.CartSet) *saleOrderCreator {
	creator := new(saleOrderCreator)
	supplierCart := map[int64][]model.Cart{}
	supplierIDs := cartSet.SupplierIDs()
	// Classify cart items of different suppliers
	for _, supplierID := range supplierIDs {
		carts := []model.Cart{}
		for _, cart := range cartSet.List() {
			if cart.SupplierID == supplierID {
				carts = append(carts, cart)
			}
		}
		supplierCart[supplierID] = carts
	}
	creator.supplierCart = supplierCart
	return creator
}

func (c *saleOrderCreator) setGoods(goods []*model.Goods) {
	goodsSet := model.NewGoodsSet(goods)
	c.goodsMap = goodsSet.Map()
}

func (c *saleOrderCreator) bindNecessary(stocks []*model.GoodsStock, userID int64, address *dto.UserAddress) {
	stockSet := model.NewStockSet(stocks)
	stockMap := stockSet.Map()
	supplierOrders := []*supplierOrder{}
	for supplierID, carts := range c.supplierCart {
		supplierOrder := &supplierOrder{
			supplierID: supplierID,
		}
		supplierOrder.bindBasically(userID, address)
		supplierStocks := make([]supplierStock, len(carts))
		for idx, cart := range carts {
			supplierStock := supplierStock{
				SupplierID: supplierID,
				Quantity:   cart.Quantity,
			}
			if stock, ok := stockMap[cart.StockID]; ok {
				supplierStock.Stock = stock
			}
			if goods, ok := c.goodsMap[cart.GoodsID]; ok {
				supplierStock.Goods = goods
			}
			supplierStocks[idx] = supplierStock
		}
		supplierOrder.supplierStocks = supplierStocks
		supplierOrders = append(supplierOrders, supplierOrder)
	}
	c.supplierOrders = supplierOrders
}

func (c *saleOrderCreator) rawSupplierOrders() []*supplierOrder {
	return c.supplierOrders
}

type SaleOrderSet struct {
	orders    []model.SaleOrder
	orderIDs  []int64
	goodsList []model.SaleDetail
}

func newSaleOrderSet(orders []model.SaleOrder) *SaleOrderSet {
	set := &SaleOrderSet{}
	set.orders = orders
	orderIds := []int64{}
	for _, order := range orders {
		orderIds = append(orderIds, order.ID)
	}
	set.orderIDs = orderIds
	return set
}

func (set *SaleOrderSet) setSaleDetails(details []model.SaleDetail) {
	set.goodsList = details
}

func (set *SaleOrderSet) orderIDList() []int64 {
	return set.orderIDs
}

func (set *SaleOrderSet) orderDTOs() []response.SaleOrderItemDTO {
	if len(set.orders) == 0 {
		return nil
	}
	dtos := make([]response.SaleOrderItemDTO, len(set.orders))
	for i, model := range set.orders {
		dto := installSaleOrderItemDTO(model)
		goodsList := []response.SaleOrderGoodsDTO{}
		for _, goods := range set.goodsList {
			if model.ID == goods.OrderID {
				goodsList = append(goodsList, installSaleDetailDTO(goods))
			}
		}
		dto.GoodsList = goodsList
		dtos[i] = dto
	}
	return dtos
}
