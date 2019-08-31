package order

import (
	"database/sql"
	"gotrue/dao"
	"gotrue/dto"
	"gotrue/dto/pagination"
	"gotrue/dto/request"
	"gotrue/dto/response"
	"gotrue/facility/errors"
	"gotrue/model"
	"gotrue/service/api"
	"gotrue/service/cart"
	"gotrue/service/express"
	"gotrue/service/goods"
	"gotrue/service/region"
	"gotrue/service/user"
	"gotrue/service/wechat"
	"gotrue/service/wechat/payment"
	"gotrue/sys"
	"log"
	"strconv"
	"time"

	"github.com/looplab/fsm"

	"github.com/bwmarrin/snowflake"
	"github.com/shopspring/decimal"
)

var SaleOrderService *SaleOrder

var (
	OrderInvalidError = errors.NewWithCodef("OrderInvalid", "当前订单不支持该操作")
)

func beforeInit() {
	cart.Init()
	goods.Init()
}

func initSaleOrderService() {
	if SaleOrderService != nil {
		return
	}
	beforeInit()
	SaleOrderService = &SaleOrder{
		dao:              dao.SaleOrderDao,
		stockDao:         dao.StockDao,
		goodsDao:         dao.GoodsDao,
		saleDetailDao:    dao.SaleDetailDao,
		wechatPaymentDao: dao.WechatPaymentDao,
		regionService:    region.RegionService,
		cartService:      cart.CartService,
		userService:      user.UserService,
		goodsService:     goods.GoodsService,
	}

}

type saleOrderFSM struct {
	fsm *fsm.FSM
}

func newOrderFSM(order *model.SaleOrder) *saleOrderFSM {
	return &saleOrderFSM{
		fsm: fsm.NewFSM(
			order.Status.String(),
			fsm.Events{
				{Name: "pay", Src: []string{model.Created.String(), model.PayFailed.String(), model.Paying.String()}, Dst: model.Paid.String()},
				{Name: "cancel", Src: []string{model.Created.String()}, Dst: model.Canceled.String()},
				{Name: "send", Src: []string{model.Paid.String(), model.Sent.String()}, Dst: model.Sent.String()},
				{Name: "post_sale", Src: []string{model.Received.String()}, Dst: model.PostSaleFinished.String()},
				{Name: "refund", Src: []string{model.Paid.String()}, Dst: model.Refound.String()},
			},
			fsm.Callbacks{},
		),
	}
}

func (sof *saleOrderFSM) can(event string) error {
	if !sof.fsm.Can(event) {
		return OrderInvalidError
	}
	return nil
}

func (sof *saleOrderFSM) orderStatus() model.OrderStatus {
	return model.OrderStatus(sof.fsm.Current())
}

// SaleOrder 作为销售订单服务，实现了api.IOrderService
type SaleOrder struct {
	dao              *dao.SaleOrder
	saleDetailDao    *dao.SaleDetail
	stockDao         *dao.Stock
	goodsDao         *dao.Goods
	wechatPaymentDao *dao.WechatPayment
	goodsService     *goods.Goods
	cartService      *cart.Cart
	userService      *user.User
	regionService    api.IRegionService
}

// UpdateExpressInfo will send express
func (s *SaleOrder) UpdateExpressInfo(req *request.OrderExpressUpdate) error {
	if err := express.IsSupport(req.ExpressMethod); err != nil {
		return err
	}
	order, err := s.dao.SelectByID(req.OrderID)
	if err != nil {
		return err
	}
	// TODO 检查订单状态
	sof := newOrderFSM(order)
	err = sof.can("send")
	if err != nil {
		return err
	}
	sof.fsm.Event("send")
	// 修改订单状态为已经发货
	order.ExpressMethod = &req.ExpressMethod
	order.ExpressNo = &req.ExpressNo
	order.Status = model.OrderStatus(sof.fsm.Current())

	if err := s.dao.UpdateByID(order, nil); err != nil {
		return err
	}
	return nil
}

func (s *SaleOrder) payStatus(req *payment.QueryOrderResponse) model.OrderStatus {
	var orderStatus model.OrderStatus
	if req.TradeState == payment.Success {
		orderStatus = model.Paid
	}
	if req.TradeState == payment.Paying || req.TradeState == payment.NotPay {
		orderStatus = model.Paying
	}
	if req.TradeState == payment.CLOSED {
		orderStatus = model.Closed
	}
	if req.TradeState == payment.PayError {
		orderStatus = model.PayFailed
	}
	return orderStatus
}

func (s *SaleOrder) notifyFarmer(prepayID string, order *model.SaleOrder) {
	// if paid success then notify tenant to send goods
	if model.Paid != order.Status {
		return
	}
	openIDLXC := "ovxEC5YTWQk6Vv5FJdN_30gkBr-g"
	notificationReq := &wechat.NotifyRequest{
		ToUser:     openIDLXC,
		TemplateID: "F56_89H1A2SiyEmnwUSGNw_kyTIcdFLBELFaU2sFUhU",
		FormID:     prepayID,
		Data: map[string]interface{}{
			"keyword1": map[string]string{
				"value": order.OrderNo,
			},
			"keyword2": map[string]string{
				"value": "",
			},
			"keyword3": map[string]interface{}{
				"value": order.OrderAmt,
			},
			"keyword4": map[string]interface{}{
				"value": "商品名称",
			},
			"keyword5": map[string]interface{}{
				"value": order.Status.Name(),
			},
			"keyword6": map[string]interface{}{
				"value": "2019-08-20",
			},
		},
	}
	wechat.WechatService().SendNotification(notificationReq)
}

// UpdateByWechatPayResult 通过微信支付查询结果更新订单状态和交易状态
func (s *SaleOrder) UpdateByWechatPayResult(r *request.QueryWechatPayResult, req *payment.QueryOrderResponse) error {
	order, err := s.dao.SelectByID(r.OrderID)
	if err != nil {
		return err
	}
	// 查找支付交易
	txns, err := s.wechatPaymentDao.SelectByOrderNo(order.OrderNo, model.TransactionTypePay)
	if err != nil {
		return err
	}
	orderStatus := s.payStatus(req)
	if orderStatus == "" {
		return nil
	}
	order.Status = orderStatus
	// 如果当前订单不是主单
	if !order.IsMaster() {
		sys.GetEasyDB().ExecTx(func(tx *sql.Tx) error {
			err := s.dao.UpdateByID(order, nil)
			if err != nil {
				return err
			}
			// 一个或者多个订单对应一笔支付
			if len(txns) == 1 {
				updateMap := map[string]interface{}{
					"status": req.TradeState,
				}
				err = s.wechatPaymentDao.UpdateByID(txns[0].ID, updateMap, nil)
				if err != nil {
					return err
				}
			}
			return nil
		})
		return nil
	}
	subOrders, err := s.dao.SelectByParentID(order.ID)
	if err != nil {
		return err
	}

	sys.GetEasyDB().ExecTx(func(tx *sql.Tx) error {
		err := s.dao.UpdateByID(order, tx)
		if err != nil {
			return err
		}
		for _, subOrder := range subOrders {
			subOrder.Status = orderStatus
			err = s.dao.UpdateByID(&subOrder, tx)
			if err != nil {
				return err
			}
		}
		// 一个或者多个订单对应一笔支付
		if len(txns) == 1 {
			updateMap := map[string]interface{}{
				"status": req.TradeState,
			}
			return s.wechatPaymentDao.UpdateByID(txns[0].ID, updateMap, nil)
		}

		return nil
	})
	go s.notifyFarmer(r.PrepayID, order)
	return nil
}

// Create 从购物车创建订单
// 1. 创建订单
// 2. 创建订单明细
// 3. 删除购物车所选中的项目x
// 4. 获取邮费及配送范围
func (s *SaleOrder) Create(userID int64, req request.SaleOrderAddRequest) (id int64, err error) {
	checkedItems, err := s.cartService.CheckedItems(userID)
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
				id, err = s.saveOrders(splittedSupplierOrders, tx)
				if err != nil {
					return err
				}
			} else {
				id, err = s.save(so, tx)
				if err != nil {
					return err
				}
			}

		}
		return nil
	})
	return
}

// Cancel will cancel order
func (s *SaleOrder) Cancel(orderID int64) (*response.SaleOrderInfo, error) {
	saleOrder, err := s.dao.SelectByID(orderID)
	if err != nil {
		return nil, err
	}
	sof := newOrderFSM(saleOrder)
	err = sof.can("cancel")
	if err != nil {
		return nil, err
	}
	err = sof.fsm.Event("cancel")
	if err != nil {
		return nil, err
	}
	saleOrder.Status = sof.orderStatus()
	err = s.dao.UpdateByID(saleOrder, nil)
	if err != nil {
		return nil, err
	}
	return s.Info(orderID)
}

// QuickCreate 快速下单
func (s *SaleOrder) QuickCreate(userID int64, req request.SaleOrderQuickAddRequest) (id int64, err error) {
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
	expressConstraint, err := s.goodsService.ExpressConstraint(stock.ID, int64(address.ID))
	if !expressConstraint.IsOK() {
		return 0, expressConstraint.Error()
	}
	expressFee := expressConstraint.UnitExpressFee
	// 由于目前水果的包装都是按照一个快递单来进行的，所以目前对每一个项目创建一个订单
	sys.GetEasyDB().ExecTx(func(tx *sql.Tx) error {
		ss := []supplierStock{
			supplierStock{
				SupplierID:     goods.SupplierID,
				Quantity:       req.Quantity,
				Stock:          stock,
				Goods:          goods,
				UnitExpressFee: expressFee,
			},
		}
		so := &supplierOrder{
			supplierID:     goods.SupplierID,
			supplierStocks: ss,
		}
		so.bindBasically(userID, address)
		if so.splitable() {
			splittedSupplierOrders := so.split()
			id, err = s.saveOrders(splittedSupplierOrders, tx)
			if err != nil {
				return err
			}
		} else {
			id, err = s.save(so, tx)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return
}

func (s *SaleOrder) saveOrders(orderList []*supplierOrder, tx *sql.Tx) (int64, error) {
	var masterID int64
	if len(orderList) > 1 {
		// save master order
		masterOrder := orderList[0]
		masterID, err := s.save(masterOrder, tx)
		if err != nil {
			return 0, err
		}
		// children order
		orderList = orderList[1:]
		for _, order := range orderList {
			order.parentID = masterID
			_, err := s.save(order, tx)
			if err != nil {
				return 0, err
			}
		}
		return masterID, nil
	}
	masterID, err := s.save(orderList[0], tx)
	if err != nil {
		return 0, err
	}
	return masterID, nil
}

// save 将业务模型转换成数据库模型并持久化
func (s *SaleOrder) save(so *supplierOrder, tx *sql.Tx) (int64, error) {
	saleOrder, saleDetails, err := so.transfer()
	if err != nil {
		return 0, err
	}
	orderID, err := s.dao.Create(saleOrder, tx)
	if err != nil {
		return 0, err
	}
	for _, detail := range saleDetails {
		detail.OrderID = orderID
		_, err := s.saleDetailDao.Create(detail)
		if err != nil {
			return 0, err
		}
	}
	return orderID, nil
}

// ListSupplierOrders list orders for a supplier's admin
func (s *SaleOrder) ListSupplierOrders(supplierID int64, req request.OrderListRequest) (page *pagination.Page, err error) {
	var orders []model.SaleOrder
	var total int64
	if req.Type == request.All {
		orders, total, err = s.dao.SelectAllBySupplierWithPagination(supplierID, req.Offet(), req.Limit())
	} else {
		stats := s.mappingRequestStatus(req)
		orders, total, err = s.dao.SelectBySupplierAndStatus(supplierID, stats, req.Offet(), req.Limit())
	}
	if err != nil {
		return nil, err
	}
	orderSet := newSaleOrderSet(orders)
	details, err := s.saleDetailDao.SelectByOrderIDs(orderSet.orderIDList()...)
	if err != nil {
		return page, err
	}
	orderSet.setSaleDetails(details)
	page = req.Page
	page.SetCount(total)
	page.Data = orderSet.orderDTOs()
	return
}

// mappingRequestStatus map status from request to db
func (s *SaleOrder) mappingRequestStatus(req request.OrderListRequest) []model.OrderStatus {
	stats := []model.OrderStatus{}
	if req.Type == request.Created {
		stats = append(stats, model.Created)
	}
	if req.Type == request.Finished {
		stats = append(stats, model.Finish)
	}
	if req.Type == request.Sent {
		stats = append(stats, model.Sent)
	}
	if req.Type == request.Paid {
		stats = append(stats, model.Paid)
	}
	return stats
}

// List will list orders for a user
func (s *SaleOrder) List(userID int64, req request.OrderListRequest) (*pagination.Page, error) {
	var orders []model.SaleOrder
	var total int64
	var err error
	if req.Type == request.All {
		orders, total, err = s.dao.SelectAllByUserWithPagination(userID, req.Offet(), req.Limit())
	} else {
		orderStatusList := s.mappingRequestStatus(req)
		orders, total, err = s.dao.SelectByUserAndStatus(userID, orderStatusList, req.Offet(), req.Limit())
		if err != nil {
			return nil, err
		}
	}
	orderSet := newSaleOrderSet(orders)
	details, err := s.saleDetailDao.SelectByOrderIDs(orderSet.orderIDList()...)
	if err != nil {
		return nil, err
	}
	page := req.Page
	orderSet.setSaleDetails(details)
	page.SetCount(total)
	page.Data = orderSet.orderDTOs()
	return page, nil
}

// WechatPrepay will prepay an order
func (s *SaleOrder) WechatPrepay(userID, orderID int64) (*payment.PrepayReponse, error) {
	order, err := s.dao.SelectByID(orderID)
	if err != nil {
		return nil, err
	}
	sof := newOrderFSM(order)
	err = sof.can("pay")
	if err != nil {
		return nil, err
	}
	totalPrice := decimal.Zero
	if order.IsMaster() {
		subOrders, err := s.dao.SelectByParentID(order.ID)
		if err != nil {
			return nil, err
		}
		orderSet := newSaleOrderSet(subOrders)
		totalPrice = orderSet.sumOrderPrice()
		totalPrice = totalPrice.Add(order.OrderAmt)
	} else {
		totalPrice = order.OrderAmt
	}
	log.Printf("total price is %s\n", totalPrice.String())
	user, err := s.userService.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	prepayReq := &wechat.PrepayRequest{
		OpenID:   user.OpenID,
		OrderNo:  order.OrderNo,
		TotalFee: totalPrice.Mul(decimal.New(1, 2)).IntPart(),
		Desc:     order.OrderNo,
	}
	result, err := wechat.WechatService().PrePay(prepayReq)

	if err != nil {
		return nil, err
	}
	wp := &model.WechatPayment{
		SaleOrderID:    order.ID,
		SaleOrderNo:    order.OrderNo,
		Amount:         totalPrice,
		Status:         model.Paying.String(),
		CreateTime:     time.Now(),
		TransationType: model.TransactionTypePay,
	}
	_, err = s.wechatPaymentDao.Create(wp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Info shows sale order info but sale details are not included
func (s *SaleOrder) Info(orderID int64) (*response.SaleOrderInfo, error) {
	saleOrder, err := s.dao.SelectByID(orderID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if saleOrder == nil {
		return nil, nil
	}
	fullName, err := s.regionService.FullName(saleOrder.RegionIDs())
	if err != nil {
		return nil, err
	}
	dto := s.installSaleInfoDTO(saleOrder)
	dto.FullRegion = fullName
	return dto, nil
}

// ListGoods list sale details for a sale order
func (s *SaleOrder) ListGoods(orderID int64) ([]response.SaleOrderGoodsDTO, error) {
	goodsList, err := s.saleDetailDao.SelectByOrderID(orderID)
	if err != nil {
		return nil, err
	}
	for _, v := range goodsList {
		specDesc, err := s.goodsService.SpecificationDesc(v.GoodsID, v.SpecIDs(), "/")
		if err != nil {
			return nil, err
		}
		v.GoodsSpecDescription = specDesc
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

func (s *SaleOrder) installSaleInfoDTO(order *model.SaleOrder) *response.SaleOrderInfo {
	data := &response.SaleOrderInfo{}
	data.ID = order.ID
	data.OrderNo = order.OrderNo
	data.Status = order.Status.Name()
	data.CreatedAt = order.CreateTime.Format("2006-01-02 15:04:05")
	data.Consignee = order.Receiver
	data.Mobile = order.PhoneNo
	data.FullRegion = "TODO"
	data.Address = order.Address
	data.GoodsAmt = order.GoodsAmt

	if order.ExpressMethod != nil {
		data.ExpressMethod = *order.ExpressMethod
	}
	if order.ExpressNo != nil {
		data.ExpressNo = *order.ExpressNo
	}
	data.ExpressFee = order.ExpressFee
	data.OrderAmt = order.OrderAmt
	return data
}

func installSaleOrderItemDTO(model model.SaleOrder) response.SaleOrderItemDTO {
	dto := response.SaleOrderItemDTO{}
	dto.ID = model.ID
	dto.OrderNo = model.OrderNo
	dto.OrderAmt = model.OrderAmt
	dto.Status = model.Status.Name()
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

func installSaleDetailDTO(model *model.SaleDetail) response.SaleOrderGoodsDTO {
	dto := response.SaleOrderGoodsDTO{}
	dto.ID = model.ID
	dto.GoodsName = model.GoodsName
	dto.Quantity = model.Quantity
	dto.RetailPrice = model.SaleUnitPrice
	dto.ListPicURL = model.ListPicURL.String
	dto.GoodsSpecDescription = model.GoodsSpecDescription
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
	// 单件商品库存的运费
	UnitExpressFee decimal.Decimal
	Stock          *model.GoodsStock
	Goods          *model.Goods
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
// 一个供应商有多个库存，每个库存可能数量大于1
// 数量大于1的库存可以拆分为1的单位库存
type supplierOrder struct {
	parentID       int64
	supplierID     int64
	userID         int64
	address        *dto.UserAddress
	supplierStocks []supplierStock
}

func (so *supplierOrder) bindBasically(userID int64, address *dto.UserAddress) {
	so.userID = userID
	so.address = address
}

func (so *supplierOrder) validate() error {
	if so.userID == 0 || so.address.ID == 0 {
		return errors.NewWithCodef("AddressBoundError", "user id and address must be bound first")
	}
	return nil
}

// transfer a supplier order to a sale order db model
func (so *supplierOrder) transfer() (model.SaleOrder, []model.SaleDetail, error) {
	err := so.validate()
	if err != nil {
		return model.SaleOrder{}, nil, err
	}
	orderNo, err := generateOrderNumber(1)
	if err != nil {
		return model.SaleOrder{}, nil, err
	}
	// assemble data
	saleOrder := model.SaleOrder{
		ParentID:   so.parentID,
		OrderNo:    orderNo,
		UserID:     so.userID,
		ProvinceID: so.address.ProvinceID,
		CityID:     so.address.CityID,
		DistrictID: so.address.DistrictID,
		Address:    so.address.Address,
		Receiver:   so.address.Name,
		PhoneNo:    so.address.Mobile,
		SupplierID: so.supplierID,
	}
	saleOrder.OrderAmt = decimal.Zero
	saleOrder.GoodsAmt = decimal.Zero
	// 总邮费
	orderExpressFee := decimal.Zero
	// 每一个库存都会有运费，整个订单可能有多个库存，需要累加所有库存运费
	// 目前每个订单只有一个库存
	saleDetails := []model.SaleDetail{}
	for _, ss := range so.supplierStocks {
		sum := ss.Stock.SaleUnitPrice.Mul(ss.Quantity)
		orderExpressFee = orderExpressFee.Add(ss.UnitExpressFee)
		saleOrder.GoodsAmt = saleOrder.GoodsAmt.Add(sum)
		saleDetail := ss.saleDetail()
		saleDetails = append(saleDetails, saleDetail)
	}
	saleOrder.OrderAmt = saleOrder.GoodsAmt.Add(orderExpressFee)
	saleOrder.ExpressFee = orderExpressFee
	return saleOrder, saleDetails, nil
}

// split 按照库存单位1对一个供应商的原始订单进行拆分
func (so *supplierOrder) split() []*supplierOrder {
	sos := []*supplierOrder{}
	for _, ss := range so.supplierStocks {
		for i := 0; i < int(ss.Quantity.IntPart()); i++ {
			ssList := []supplierStock{
				supplierStock{
					SupplierID:     so.supplierID,
					Quantity:       decimal.NewFromFloat32(1.0),
					Stock:          ss.Stock,
					Goods:          ss.Goods,
					UnitExpressFee: ss.UnitExpressFee,
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
	goodsList []*model.SaleDetail
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

func (set *SaleOrderSet) sumOrderPrice() decimal.Decimal {
	sum := decimal.Zero
	if len(set.orders) == 0 {
		return sum
	}
	for _, order := range set.orders {
		sum = sum.Add(order.OrderAmt)
	}
	return sum
}

func (set *SaleOrderSet) setSaleDetails(details []*model.SaleDetail) {
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
