package api

import (
	"context"
	"gotrue/dto/pagination"
	"gotrue/dto/request"
	"gotrue/service/wechat/payment"

	"github.com/shopspring/decimal"
)

// ISaleOrderService is a order service interface
type ISaleOrderService interface {
	CreateFromCart(userID int64, req request.SaleOrderAddRequest) (id int64, err error)

	CreateFromStock(ctx context.Context, req request.SaleOrderQuickAddRequest) (id int64, err error)

	List(ctx context.Context, req OrderListRequest) (*pagination.Page, error)

	ListManagedOrders(userID int64, req OrderListRequest) (page *pagination.Page, err error)

	Info(orderID int64) (*SaleOrderResponse, error)

	Cancel(ctx context.Context, orderID int64) (*SaleOrderResponse, error)

	ConfirmReceive(ctx context.Context, orderID int64) (*SaleOrderResponse, error)

	Prepay(userID, orderID int64) (*payment.PrepayReponse, error)

	PayResult(r *request.QueryWechatPayResult, req *payment.QueryOrderResponse) error

	SendExpress(req *request.OrderExpressUpdate) error
}

type OrderStatus string

const (
	All      OrderStatus = "all"
	Finished OrderStatus = "finished"
	Created  OrderStatus = "created"
	Sent     OrderStatus = "sent"
	Paid     OrderStatus = "paid"
)

type OrderListRequest struct {
	*pagination.Page
	Type OrderStatus `json:"type"`
}

func (r OrderListRequest) IsAllStatus() bool {
	return r.Type == All
}

type SaleDetailResponse struct {
	ID                   int64           `json:"id"`
	GoodsName            string          `json:"goods_name"`
	Quantity             decimal.Decimal `json:"number"`
	ListPicURL           string          `json:"list_pic_url"`
	RetailPrice          decimal.Decimal `json:"retail_price"`
	GoodsSpecDescription string          `json:"spec_desc"`
}

type SaleActions map[string]bool

type SaleOrderResponse struct {
	ID            int64                `json:"id"`
	CreatedAt     string               `json:"created_at"`
	OrderNo       string               `json:"order_sn"`
	OrderAmt      decimal.Decimal      `json:"actual_price"`
	GoodsAmt      decimal.Decimal      `json:"goods_amt"`
	Status        string               `json:"order_status_text"`
	ExpressMethod string               `json:"express_method"`
	ExpressNo     string               `json:"express_no"`
	ExpressFee    decimal.Decimal      `json:"express_fee"`
	FullRegion    string               `json:"full_region"`
	Address       string               `json:"address"`
	Consignee     string               `json:"consignee"`
	Mobile        string               `json:"mobile"`
	Actions       SaleActions          `json:"actions"`
	Details       []SaleDetailResponse `json:"details"`
}
