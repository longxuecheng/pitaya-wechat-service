package service

import (
	"database/sql"
	"pitaya-wechat-service/dao"
	"pitaya-wechat-service/dto/pagination"
	"pitaya-wechat-service/dto/request"
	"pitaya-wechat-service/dto/response"
	"pitaya-wechat-service/model"
	"pitaya-wechat-service/sys"
	"strconv"

	"github.com/bwmarrin/snowflake"
	"github.com/shopspring/decimal"
)

var saleOrderServiceSingleton *SaleOrderService

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

func (s *SaleOrderService) QuickCreate(req request.SaleOrderQuickAddRequest) error {
	stock, err := s.stockDao.SelectByID(req.StockID)
	if err != nil {
		return err
	}
	address, err := s.userService.GetAddressByID(req.AddressID)
	if err != nil {
		return err
	}
	goods, err := s.goodsDao.SelectByID(stock.GoodsID)
	if err != nil {
		return err
	}
	sys.GetEasyDB().ExecTx(func(tx *sql.Tx) error {
		for i := 0; i < req.Amount; i++ {
			orderNo, err := generateOrderNumber(1)
			if err != nil {
				return err
			}
			totalGoodsPrice := decimal.Zero
			setMap := map[string]interface{}{
				"order_no":    orderNo,
				"user_id":     req.UserID,
				"receiver":    address.Name,
				"phone_no":    address.Mobile,
				"goods_amt":   totalGoodsPrice,
				"express_fee": decimal.Zero,
				"order_amt":   totalGoodsPrice.Add(decimal.Zero),
				"province_id": address.ProvinceID,
				"city_id":     address.CityID,
				"district_id": address.DistrictID,
				"address":     address.Address,
			}
			orderID, err := s.dao.Create(setMap)
			if err != nil {
				return err
			}
			saleDetail := installSaleDetailFromStock(orderID, stock, goods)
			_, err = s.saleDetailDao.Create(saleDetail, tx)
			if err != nil {
				return err
			}

		}

		return nil
	})

	return nil
}

// Create 从购物车创建订单
// 1. 创建订单
// 2. 创建订单明细
// 3. 删除购物车所选中的项目
func (s *SaleOrderService) Create(userID int64, req request.SaleOrderAddRequest) (id int64, err error) {
	allCarts, err := s.cartService.ListCart4User(userID)
	if err != nil {
		return
	}
	address, err := s.userService.GetAddressByID(req.AddressID)
	if err != nil {
		return
	}
	orderCreator := newSaleOrderCreator(allCarts)
	orderNo, err := generateOrderNumber(1)
	if err != nil {
		return
	}
	setMap := map[string]interface{}{
		"order_no":    orderNo,
		"user_id":     userID,
		"receiver":    address.Name,
		"phone_no":    address.Mobile,
		"goods_amt":   orderCreator.totalGoodsPrice(),
		"express_fee": decimal.Zero,
		"order_amt":   orderCreator.totalGoodsPrice().Add(decimal.Zero),
		"province_id": address.ProvinceID,
		"city_id":     address.CityID,
		"district_id": address.DistrictID,
		"address":     address.Address,
	}
	sys.GetEasyDB().ExecTx(func(tx *sql.Tx) error {
		id, err = s.dao.Create(setMap)
		if err != nil {
			return err
		}
		for _, detail := range orderCreator.installSaleItemsFromCart(id) {
			_, err := s.saleDetailDao.Create(detail)
			if err != nil {
				return err
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

// SaleOrderCreator 是订单服务中特有的购物车管理
// 若进行为服务拆分，那么购物车可能会作为一个单独的服务提供数据，
// 该对象对返回的购物车数据进行业务处理
type SaleOrderCreator struct {
	checkedItems         []response.CartItemDTO
	checkedGoodsAmt      decimal.Decimal
	checkedGoodsQuantity decimal.Decimal
}

func newSaleOrderCreator(totalItems []response.CartItemDTO) *SaleOrderCreator {
	manager := new(SaleOrderCreator)
	checkedCarts := []response.CartItemDTO{}
	checkedGoodsQuantity := decimal.Zero
	checkedGoodsAmt := decimal.Zero
	for _, cart := range totalItems {
		if cart.Checked == 1 {
			checkedCarts = append(checkedCarts, cart)
			checkedGoodsQuantity = checkedGoodsQuantity.Add(cart.Quantity)
			checkedGoodsAmt = checkedGoodsAmt.Add(cart.Quantity.Mul(cart.RetailPrice))
		}
	}
	manager.checkedItems = checkedCarts
	manager.checkedGoodsAmt = checkedGoodsAmt
	manager.checkedGoodsQuantity = checkedGoodsQuantity
	return manager
}

func (cm *SaleOrderCreator) totalGoodsPrice() decimal.Decimal {
	return cm.checkedGoodsAmt
}

func (cm *SaleOrderCreator) totalGoodsQuantity() decimal.Decimal {
	return cm.checkedGoodsQuantity
}

func (cm *SaleOrderCreator) installSaleItemsFromCart(orderID int64) []model.SaleDetail {
	saleDetails := make([]model.SaleDetail, len(cm.checkedItems))
	for i, item := range cm.checkedItems {
		saleDetail := model.SaleDetail{
			OrderID:              orderID,
			GoodsID:              item.GoodsID,
			GoodsName:            item.GoodsName,
			Quantity:             item.Quantity,
			StockID:              item.StockID,
			SaleUnitPrice:        item.RetailPrice,
			GoodsSpecIDs:         item.GoodsSpecIDs,
			GoodsSpecDescription: item.GoodsSpecDescription,
			ListPicURL:           sql.NullString{String: item.ListPicURL},
		}
		saleDetails[i] = saleDetail
	}
	return saleDetails
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
