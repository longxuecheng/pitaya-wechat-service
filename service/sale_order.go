package service

import (
	"database/sql"
	"pitaya-wechat-service/dao"
	"pitaya-wechat-service/dto/pagination"
	"pitaya-wechat-service/dto/request"
	"pitaya-wechat-service/dto/response"
	"pitaya-wechat-service/model"

	"github.com/shopspring/decimal"
)

var saleOrderServiceSingleton *SaleOrderService

func SaleOrderServiceInstance() *SaleOrderService {
	if saleOrderServiceSingleton == nil {
		saleOrderServiceSingleton = new(SaleOrderService)
		saleOrderServiceSingleton.dao = dao.SaleOrderDaoInstance()
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
	cartService   *CartService
	userService   *UserService
}

// Create 创建订单有 1. 创建订单 2. 创建订单明细 3. 删除购物车所选中的项目
func (s *SaleOrderService) Create(req request.SaleOrderAddRequest) (id int64, err error) {
	allCarts, err := s.cartService.ListCart4User(req.UserID)
	if err != nil {
		return
	}
	address, err := s.userService.GetAddressByID(req.AddressID)
	if err != nil {
		return
	}
	cartManager := newCartManger(allCarts)
	setMap := map[string]interface{}{
		"order_no":    nil,
		"user_id":     req.UserID,
		"receiver":    address.Name,
		"phone_no":    address.Mobile,
		"goods_amt":   cartManager.totalGoodsPrice(),
		"express_fee": decimal.Zero,
		"order_amt":   cartManager.totalGoodsPrice().Add(decimal.Zero),
		"province_id": address.ProvinceID,
		"city_id":     address.CityID,
		"district_id": address.DistrictID,
		"address":     address.Address,
	}
	s.dao.ExecTx(func(tx *sql.Tx) error {
		id, err = s.dao.Create(setMap)
		if err != nil {
			return err
		}
		for _, detail := range cartManager.installSaleDetails(id) {
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

func (s *SaleOrderService) installSaleInfoDTO(model model.SaleOrder) response.SaleOrderInfoDTO {
	dto := response.SaleOrderInfoDTO{}
	dto.ID = model.ID
	dto.OrderNo = model.OrderNo.String
	dto.Status = model.Status
	dto.CreatedAt = model.CreateTime.Format("2006-01-02 15:04:05")
	dto.Consignee = model.Receiver
	dto.Mobile = model.PhoneNo
	dto.FullRegion = "TODO"
	dto.Address = model.Address
	dto.GoodsAmt = model.GoodsAmt
	dto.ExpressFee = model.ExpressFee
	dto.OrderAmt = model.OrderAmt
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

// cartManger 是订单服务中特有的购物车管理
// 若进行为服务拆分，那么购物车可能会作为一个单独的服务提供数据，
// 该对象对返回的购物车数据进行业务处理
type cartManger struct {
	checkedItems         []response.CartItemDTO
	checkedGoodsAmt      decimal.Decimal
	checkedGoodsQuantity decimal.Decimal
}

func newCartManger(totalItems []response.CartItemDTO) *cartManger {
	manager := new(cartManger)
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

func (cm *cartManger) totalGoodsPrice() decimal.Decimal {
	return cm.checkedGoodsAmt
}

func (cm *cartManger) totalGoodsQuantity() decimal.Decimal {
	return cm.checkedGoodsQuantity
}

func (cm *cartManger) installSaleDetails(orderID int64) []model.SaleDetail {
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
