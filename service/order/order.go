package order

import (
	"context"
	"database/sql"
	"gotrue/dao"
	"gotrue/dto/pagination"
	"gotrue/dto/request"
	"gotrue/facility/context_util"
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
	"time"

	"github.com/rs/zerolog/log"

	"github.com/shopspring/decimal"
)

var SaleOrderService api.ISaleOrderService

var (
	ErrInvalidOption = errors.NewWithCodef("InvalidOption", "非法操作")
)

func initSaleOrderService() {
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
	goodsService     api.IGoodsService
	cartService      *cart.Cart
	userService      *user.User
	regionService    api.IRegionService
	cutService       api.ICutService
}

// SendExpress will send express
func (s *SaleOrder) SendExpress(req *request.OrderExpressUpdate) error {
	if err := express.IsSupport(req.ExpressMethod); err != nil {
		return err
	}
	order, err := s.dao.SelectByID(req.OrderID)
	if err != nil {
		return err
	}
	sof := newOrderFSM(order)
	err = sof.can(actionSend)
	if err != nil {
		return err
	}
	sof.fsm.Event(actionSend.String())
	order.ExpressMethod = &req.ExpressMethod
	order.ExpressNo = &req.ExpressNo
	order.Status = sof.current()
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
		dao.GetEasyDB().ExecTx(func(tx *sql.Tx) error {
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
	dao.GetEasyDB().ExecTx(func(tx *sql.Tx) error {
		err := s.dao.UpdateByID(order, tx)
		if err != nil {
			return err
		}
		for _, subOrder := range subOrders {
			subOrder.Status = orderStatus
			err = s.dao.UpdateByID(subOrder, tx)
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
			if me.ChannelUserID == 0 {
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

func (s *SaleOrder) ConfirmReceive(ctx context.Context, orderID int64) (*api.SaleOrderResponse, error) {
	userID, _ := context_util.GetUserID(ctx)
	order, err := s.dao.QueryByUserAndID(userID, orderID)
	if err == sql.ErrNoRows {
		return nil, ErrInvalidOption
	}
	if err != nil {
		return nil, err
	}
	sof := newOrderFSM(order)
	err = sof.can(actionReceive)
	if err != nil {
		return nil, err
	}
	err = sof.fsm.Event(actionReceive.String())
	if err != nil {
		return nil, err
	}
	order.Status = sof.current()
	apiDetails, err := s.installDetails(orderID)
	if err != nil {
		return nil, err
	}
	err = s.dao.UpdateByID(order, nil)
	if err != nil {
		return nil, err
	}
	apiData := installInfo(order)
	apiData.Details = apiDetails
	return apiData, nil
}

// Cancel will cancel order
func (s *SaleOrder) Cancel(ctx context.Context, orderID int64) (*api.SaleOrderResponse, error) {
	userID, _ := context_util.GetUserID(ctx)
	order, err := s.dao.QueryByUserAndID(userID, orderID)
	if err == sql.ErrNoRows {
		return nil, ErrInvalidOption
	}
	if err != nil {
		return nil, err
	}
	sof := newOrderFSM(order)
	err = sof.can(actionCancel)
	if err != nil {
		return nil, err
	}
	err = sof.fsm.Event(actionCancel.String())
	if err != nil {
		return nil, err
	}
	order.Status = sof.current()
	// TODO 取消的后续操作
	apiDetails, err := s.installDetails(orderID)
	if err != nil {
		return nil, err
	}
	err = s.dao.UpdateByID(order, nil)
	if err != nil {
		return nil, err
	}
	apiData := installInfo(order)
	apiData.Details = apiDetails
	return apiData, nil
}

// CreateFromStock create order from stock
func (s *SaleOrder) CreateFromStock(ctx context.Context, req request.SaleOrderQuickAddRequest) (id int64, err error) {
	userID, _ := context_util.GetUserID(ctx)
	err = req.Validate()
	if err != nil {
		return
	}
	stock, err := s.stockDao.SelectByID(req.StockID)
	if err != nil {
		return 0, err
	}
	goods, err := s.goodsDao.SelectOnSaleByID(stock.GoodsID)
	if err == sql.ErrNoRows {
		return 0, errors.NewWithCodef("GoodsNotAvailable", "该商品当前不可售")
	}
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
	dao.GetEasyDB().ExecTx(func(tx *sql.Tx) error {
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
func (s *SaleOrder) ListManagedOrders(userID int64, req api.OrderListRequest) (page *pagination.Page, err error) {
	admins, err := s.supplierAdminDao.QueryByUserID(userID)
	if err != nil {
		return nil, err
	}
	var orders model.SaleOrderList
	var total int64
	if req.IsAllStatus() {
		orders, total, err = s.dao.SelectAllBySuppliersWithPagination(admins.SupplierIDs(), req.Offet(), req.Limit())
	} else {
		stats := []model.OrderStatus{apiStatus2ModelStatusMap[req.Type]}
		orders, total, err = s.dao.SelectBySupplierAndStatus(admins.SupplierIDs(), stats, req.Offet(), req.Limit())
	}
	if err != nil {
		return nil, err
	}
	details, err := s.saleDetailDao.QueryByIDs(orders.IDList())
	if err != nil {
		return page, err
	}
	page = req.Page
	page.SetCount(total)
	page.Data = newResponseBuilder().setOrders(orders).setDetails(details).buildList()
	return
}

var apiStatus2ModelStatusMap = map[api.OrderStatus]model.OrderStatus{
	api.Created:  model.Created,
	api.Finished: model.Finish,
	api.Sent:     model.Sent,
	api.Paid:     model.Paid,
}

// List will list orders for a user
func (s *SaleOrder) List(ctx context.Context, req api.OrderListRequest) (*pagination.Page, error) {
	userID, _ := context_util.GetUserID(ctx)
	var orders model.SaleOrderList
	var total int64
	var err error
	if req.IsAllStatus() {
		orders, total, err = s.dao.SelectAllByUserWithPagination(userID, req.Offet(), req.Limit())
	} else {
		stats := []model.OrderStatus{apiStatus2ModelStatusMap[req.Type]}
		orders, total, err = s.dao.SelectByUserAndStatus(userID, stats, req.Offet(), req.Limit())
		if err != nil {
			return nil, err
		}
	}
	details, err := s.saleDetailDao.QueryByIDs(orders.IDList())
	if err != nil {
		return nil, err
	}
	page := req.Page
	page.SetCount(total)
	page.Data = newResponseBuilder().setOrders(orders).setDetails(details).buildList()
	return page, nil
}

// Prepay will prepay an order
func (s *SaleOrder) Prepay(userID, orderID int64) (*payment.PrepayReponse, error) {
	order, err := s.dao.SelectByID(orderID)
	if err != nil {
		return nil, err
	}
	sof := newOrderFSM(order)
	err = sof.can(actionPay)
	if err != nil {
		return nil, err
	}
	totalPrice := decimal.Zero
	if order.IsMaster() {
		subOrders, err := s.dao.SelectByParentID(order.ID)
		if err != nil {
			return nil, err
		}
		totalPrice = subOrders.TotalPrice()
		totalPrice = totalPrice.Add(order.OrderAmt)
	} else {
		totalPrice = order.OrderAmt
	}
	log.Printf("total price is %s\n", totalPrice.String())
	user, err := s.userService.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	result, err := wechat.WechatService().PrePay(&wechat.PrepayRequest{
		OpenID:   user.OpenID,
		OrderNo:  order.OrderNo,
		TotalFee: totalPrice.Mul(decimal.New(1, 2)).IntPart(),
		Desc:     order.OrderNo,
	})

	if err != nil {
		return nil, err
	}
	_, err = s.wechatPaymentDao.Create(&model.WechatPayment{
		SaleOrderID:    order.ID,
		SaleOrderNo:    order.OrderNo,
		Amount:         totalPrice,
		Status:         model.Paying.String(),
		CreateTime:     time.Now(),
		TransationType: model.TransactionTypePay,
	}, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Info shows sale order info but sale details are not included
func (s *SaleOrder) Info(orderID int64) (*api.SaleOrderResponse, error) {
	order, err := s.dao.SelectByID(orderID)
	if err == sql.ErrNoRows {
		return nil, errors.NewWithCodef("OrderNotExists", "订单不存在")
	}
	if err != nil {
		return nil, err
	}
	fullName, err := s.regionService.FullName(order.RegionIDs())
	if err != nil {
		return nil, err
	}
	apiDetails, err := s.installDetails(orderID)
	if err != nil {
		return nil, err
	}
	apiOrder := installInfo(order)
	apiOrder.FullRegion = fullName
	apiOrder.Details = apiDetails
	return apiOrder, nil
}

func (s *SaleOrder) installDetails(orderID int64) ([]api.SaleDetailResponse, error) {
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
	apiDataList := make([]api.SaleDetailResponse, len(goodsList))
	for i, goods := range goodsList {
		apiDataList[i] = installDetail(goods)
	}
	return apiDataList, nil
}
