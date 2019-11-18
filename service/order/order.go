package order

import (
	"database/sql"
	"gotrue/dao"
	"gotrue/dto/pagination"
	"gotrue/dto/request"
	"gotrue/dto/response"
	"gotrue/facility/errors"
	"gotrue/facility/utils"
	"gotrue/model"
	"gotrue/service/api"
	"gotrue/service/cart"
	"gotrue/service/cut"
	"gotrue/service/express"
	"gotrue/service/goods"
	"gotrue/service/region"
	"gotrue/service/sms"
	"gotrue/service/user"
	"gotrue/service/wechat"
	"gotrue/service/wechat/payment"
	"gotrue/sys"
	"time"

	"github.com/looplab/fsm"
	"github.com/rs/zerolog/log"

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
		supplierAdminDao: dao.SupplierAdminDao,
		goodsDao:         dao.GoodsDao,
		userDao:          dao.UserDaoSingleton,
		saleDetailDao:    dao.SaleDetailDao,
		wechatPaymentDao: dao.WechatPaymentDao,
		regionService:    region.RegionService,
		cartService:      cart.CartService,
		userService:      user.UserService,
		goodsService:     goods.GoodsService,
		cutService:       cut.GetCutService(),
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
	supplierDao      *dao.Supplier
	supplierAdminDao *dao.SupplierAdmin
	stockDao         *dao.Stock
	goodsDao         *dao.Goods
	userDao          *dao.UserDao
	wechatPaymentDao *dao.WechatPayment
	goodsService     *goods.Goods
	cartService      *cart.Cart
	userService      *user.User
	regionService    api.IRegionService
	cutService       api.ICutService
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

// notifyFarmer deprecated!!!
// eg: PersonA paid a bill and wechat mp do not support mp send template message to
// PersonB
func (s *SaleOrder) notifyFarmer(prepayID string, order *model.SaleOrder) {
	// if paid success then notify tenant to send goods
	if model.Paid != order.Status {
		return
	}
	supplierName := ""
	openID := "ovxEC5YTWQk6Vv5FJdN_30gkBr-g"
	supplier, err := s.supplierDao.SelectByID(order.SupplierID)
	if err != nil {
		log.Printf("query supplier %d error %+v\n", order.SupplierID, err)
	}
	supplierName = supplier.Name
	supplierAdmin, err := s.userDao.SelectByID(0)
	if err != nil {
		log.Printf("query supplier admin %d error %+v\n", supplier.ID, err)
	}
	if supplierAdmin != nil {
		openID = supplierAdmin.WechatID
	}
	log.Printf("notify to open id %s\n", openID)
	data := map[string]interface{}{
		"keyword1": map[string]string{
			"value": order.OrderNo,
		},
		"keyword2": map[string]string{ // 订单详情
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
		"keyword6": map[string]interface{}{ // 下单时间
			"value": utils.FormatTime(order.CreateTime, utils.TimePrecision_Seconds),
		},
		"keyword7": map[string]interface{}{ // 地址
			"value": "",
		},
		"keyword8": map[string]interface{}{
			"value": "请及时发货，有问题随时沟通客户",
		},
		"keyword9": map[string]interface{}{
			"value": supplierName,
		},
		"keyword10": map[string]interface{}{
			"value": order.PhoneNo,
		},
	}
	notificationReq := &wechat.TemplateMsgRequest{
		ToUser:     openID,
		TemplateID: templateID,
		FormID:     prepayID,
		Data:       data,
	}
	err = wechat.WechatService().SendTemplateMessage(notificationReq)
	if err != nil {
		log.Printf("[Error] Send template message to %s error %+v\n", openID, err)
	}
	uniformMsgReq := &wechat.UniformMsgRequest{
		ToUser: openID,
		WebAppTemplateMsg: wechat.WebAppTemplateMsg{
			TemplateID:      templateID,
			FormID:          prepayID,
			Data:            data,
			EmphasisKeyword: "keyword1.DATA",
		},
	}
	err = wechat.WechatService().SendUniformMessage(uniformMsgReq)
	if err != nil {
		log.Printf("[Error] Send uniform message to %s error %+v\n", openID, err)
	}
}

// PayResult 通过微信支付查询结果更新订单状态和交易状态
func (s *SaleOrder) PayResult(r *request.QueryWechatPayResult, req *payment.QueryOrderResponse) error {
	order, err := s.dao.SelectByID(r.OrderID)
	if err != nil {
		return err
	}
	// 查找支付交易
	txns, err := s.wechatPaymentDao.SelectByOrderNo(order.OrderNo, model.TransactionTypePay)
	if err != nil {
		return err
	}
	orderStatus := req.OrderStatus()
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
				err = s.wechatPaymentDao.UpdateByID(txns[0].ID,
					map[string]interface{}{
						"status": req.TradeState,
					}, nil)
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
	if order.Status == model.Paid {
		go func() {
			admins, err := s.supplierAdminDao.QueryBySupplierID(order.SupplierID)
			if err != nil {
				return
			}
			users, err := s.userDao.SelectByIDs(admins.UserIDs())
			if err != nil {
				return
			}
			err = sms.SendPayNotificationMsg(&sms.MultiSendRequest{
				Mobiles: users.UserMobiles(),
				Params:  []string{"果真管理员", order.OrderNo12()},
			})
			if err != nil {
				log.Printf("Send notification to admins err %+v\n", err)
			}
			me, err := s.userDao.SelectByID(order.UserID)
			if err != nil {
				return
			}
			agency, err := s.userDao.SelectByChannelUserID(me.ChannelUserID)
			if err != nil {
				return
			}
			if !agency.HasMobile() {
				return
			}
			err = sms.SendPayNotificationMsg(&sms.MultiSendRequest{
				Mobiles: []string{agency.PhoneNo.String},
				Params:  []string{"果真Gotrue代理", order.OrderNo12()},
			})
			if err != nil {
				log.Printf("Send notification to agency err %+v\n", err)
			}
		}()
	}
	return nil
}

func (s *SaleOrder) CreateFromCart(userID int64, req request.SaleOrderAddRequest) (id int64, err error) {
	// checkedItems, err := s.cartService.CheckedItems(userID)
	// if err != nil {
	// 	return
	// }
	// address, err := s.userService.GetAddressByID(req.AddressID)
	// if err != nil {
	// 	return
	// }
	// cartSet := model.NewCartSet(checkedItems)
	// goodsIDs := cartSet.GoodsIDs()
	// goodsList, err := s.goodsDao.SelectByIDs(goodsIDs)
	// if err != nil {
	// 	return
	// }
	// stocks, err := s.stockDao.SelectByIDs(cartSet.StockIDs())
	// if err != nil {
	// 	return
	// }
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

// CreateFromStock create order from stock
func (s *SaleOrder) CreateFromStock(userID int64, req request.SaleOrderQuickAddRequest) (id int64, err error) {
	err = req.Validate()
	if err != nil {
		return
	}
	stock, err := s.stockDao.SelectByID(req.StockID)
	if err != nil {
		return 0, err
	}
	cutReq := &request.CutOrder{
		GoodsID: stock.GoodsID,
		UserID:  userID,
	}
	cutorder, err := s.cutService.MyActivatedCutOrder(cutReq)
	if err != nil {
		return 0, err
	}
	address, err := s.userService.GetAddressByID(req.AddressID)
	if err != nil {
		return 0, err
	}
	goods, err := s.goodsDao.SelectByID(stock.GoodsID)
	if err == sql.ErrNoRows {
		return 0, errors.NewWithCodef("GoodsInvalidStatus", "商品不可售")
	}
	if err != nil {
		return 0, err
	}
	expressConstraint, err := s.goodsService.ExpressConstraint(stock.ID, int64(address.ID))
	if !expressConstraint.IsOK() {
		return 0, expressConstraint.Error()
	}
	expressFee := expressConstraint.UnitExpressFee
	sb := &StockOrderBuilder{
		Goods:          goods,
		Quantity:       req.Quantity,
		Stock:          stock,
		UserID:         userID,
		UnitExpressFee: expressFee,
		Address:        address,
		CutOrder:       cutorder,
	}
	err = sb.Build()
	if err != nil {
		return 0, err
	}
	sb.CuttoffFirst()
	sys.GetEasyDB().ExecTx(func(tx *sql.Tx) error {
		id, err = s.createStockOrders(sb.StockOrders(), tx)
		if err != nil {
			return err
		}
		if cutorder != nil {
			consumeReq := &request.ConsumeCutOrder{
				SaleOrderID: id,
				CutOrderID:  cutorder.ID,
			}
			return s.cutService.ConsumeCutOrder(utils.ContextWithTx(tx), consumeReq)
		}
		return nil
	})
	return
}

func (s *SaleOrder) createOrder(order *model.SaleOrder, details []*model.SaleDetail, tx *sql.Tx) (int64, error) {
	orderID, err := s.dao.Create(order, tx)
	if err != nil {
		return 0, err
	}
	for _, detail := range details {
		detail.OrderID = orderID
		_, err := s.saleDetailDao.Create(detail, nil)
		if err != nil {
			return 0, err
		}
	}
	return orderID, nil
}

// createStockOrders will create orders and return mater order id
func (s *SaleOrder) createStockOrders(orders []*StockOrder, tx *sql.Tx) (masterID int64, err error) {
	masterOrder := orders[0]
	if len(orders) > 1 {
		// save master order
		masterID, err = s.createOrder(masterOrder.SaleOrder, masterOrder.SaleDetails, tx)
		if err != nil {
			return 0, err
		}
		// children order
		children := orders[1:]
		for _, stockOrder := range children {
			order := stockOrder.SaleOrder
			order.ParentID = masterID
			_, err := s.createOrder(order, stockOrder.SaleDetails, tx)
			if err != nil {
				return 0, err
			}
		}
	} else {
		masterID, err = s.createOrder(masterOrder.SaleOrder, masterOrder.SaleDetails, tx)
		if err != nil {
			return 0, err
		}
	}
	return masterID, nil
}

// ListManagedOrders list orders for a supplier's admin
func (s *SaleOrder) ListManagedOrders(userID int64, req request.OrderListRequest) (page *pagination.Page, err error) {
	admins, err := s.supplierAdminDao.QueryByUserID(userID)
	if err != nil {
		return nil, err
	}
	var orders []model.SaleOrder
	var total int64
	if req.Type == request.All {
		orders, total, err = s.dao.SelectAllBySuppliersWithPagination(admins.SupplierIDs(), req.Offet(), req.Limit())
	} else {
		stats := s.mappingRequestStatus(req)
		orders, total, err = s.dao.SelectBySupplierAndStatus(admins.SupplierIDs(), stats, req.Offet(), req.Limit())
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
	page.Data = orderSet.apiOrders()
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
	page.Data = orderSet.apiOrders()
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
	_, err = s.wechatPaymentDao.Create(wp, nil)
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
	data.OrderAmt = order.OrderPrice()
	return data
}

func installSaleOrderItemDTO(model model.SaleOrder) response.SaleOrderItemDTO {
	dto := response.SaleOrderItemDTO{}
	dto.ID = model.ID
	dto.OrderNo = model.OrderNo
	dto.OrderAmt = model.OrderPrice()
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
func installSaleDetailFromStock(orderID int64, stock *model.Stock, goods *model.Goods) model.SaleDetail {
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
	Stock          *model.Stock
	Goods          *model.Goods
}

func (ss supplierStock) saleDetail() *model.SaleDetail {
	return &model.SaleDetail{
		OrderID:       0,
		GoodsID:       ss.Stock.GoodsID,
		GoodsName:     ss.Goods.Name,
		Quantity:      ss.Quantity,
		StockID:       ss.Stock.ID,
		SaleUnitPrice: ss.Stock.SaleUnitPrice,
		CostUnitPrice: ss.Stock.CostUnitPrice,
		GoodsSpecIDs:  ss.Stock.Specification.String,
		ListPicURL:    ss.Goods.ListPicURL,
	}

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

func (set *SaleOrderSet) apiOrders() []response.SaleOrderItemDTO {
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
