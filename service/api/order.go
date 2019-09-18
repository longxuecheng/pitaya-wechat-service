package api

import (
	"gotrue/dto/pagination"
	"gotrue/dto/request"
	"gotrue/dto/response"
	"gotrue/service/wechat/payment"
)

// ISaleOrderService is a order service interface
type ISaleOrderService interface {
	// CreateFromCart order from selected cart items
	CreateFromCart(userID int64, req request.SaleOrderAddRequest) (id int64, err error)

	// CreateFromStock create order from selected goods
	CreateFromStock(req request.SaleOrderQuickAddRequest) (int64, error)

	List(userID int64, req pagination.Page) (page pagination.Page, err error)

	// ListManagedOrders list orders for a user that manages serveral suppliers
	ListManagedOrders(userID int64, req request.OrderListRequest) (page *pagination.Page, err error)

	Info(orderID int64) (*response.SaleOrderInfo, error)

	ListGoods(orderID int64) ([]response.SaleOrderGoodsDTO, error)

	Cancel(orderID int64) (*response.SaleOrderInfo, error)

	PayResult(r *request.QueryWechatPayResult, req *payment.QueryOrderResponse) error
}
