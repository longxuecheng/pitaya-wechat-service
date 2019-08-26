package model

import (
	"database/sql/driver"
	"time"

	"github.com/shopspring/decimal"
)

const (
	Created   OrderStatus = "CREATED"
	Canceled  OrderStatus = "CANCELED"
	Paid      OrderStatus = "PAID"
	Paying    OrderStatus = "PAYING"
	PayFailed OrderStatus = "PAY_FAILED"
	//Sent 商家已发货
	Sent             OrderStatus = "SENT"
	Received         OrderStatus = "RECEIVED"
	Refunding        OrderStatus = "REFUNDING"
	RefundRefused    OrderStatus = "REFUND_REFUSED"
	Refound          OrderStatus = "REFOUND"
	PostSale         OrderStatus = "PST_SALE"
	PostSaleFinished OrderStatus = "PST_SALE_FNS"
	Finish           OrderStatus = "FINISHED"
	Closed           OrderStatus = "CLOSED"
)

var orderStatuMap = map[OrderStatus]string{
	Created:          "已创建",
	Canceled:         "已取消",
	Paying:           "待付款",
	Paid:             "已付款", // 待收货
	PayFailed:        "付款失败",
	Sent:             "已发货",
	Received:         "已收货",
	Refunding:        "退款中",
	RefundRefused:    "拒绝退款",
	Refound:          "退款成功",
	PostSale:         "售后处理中",
	PostSaleFinished: "售后完成",
	Finish:           "完成",
}

type OrderStatus string

func (os OrderStatus) Value() (driver.Value, error) {
	return string(os), nil
}

func (os OrderStatus) String() string {
	return string(os)
}

func (so OrderStatus) Name() string {
	if statusName, ok := orderStatuMap[so]; ok {
		return statusName
	}
	return "未知"
}

type SaleOrder struct {
	ID            int64           `db:"id" exclude:"true"`
	ParentID      int64           `db:"parent_id"`
	OrderNo       string          `db:"order_no"`
	CreateTime    time.Time       `db:"create_time" exclude:"true"`
	UserID        int64           `db:"user_id"`
	Status        OrderStatus     `db:"status" exclude:"true"`
	Receiver      string          `db:"receiver"`
	ProvinceID    int             `db:"province_id"`
	CityID        int             `db:"city_id"`
	DistrictID    int             `db:"district_id"`
	Address       string          `db:"address"`
	PhoneNo       string          `db:"phone_no"`
	SupplierID    int64           `db:"supplier_id"`
	OrderAmt      decimal.Decimal `db:"order_amt"`
	GoodsAmt      decimal.Decimal `db:"goods_amt"`
	ExpressMethod *string         `db:"express_method"`
	ExpressNo     *string         `db:"express_order_no"`
	ExpressFee    decimal.Decimal `db:"express_fee"`
	Count         int64           `db:"count" exclude:"true"`
}

func (so *SaleOrder) TableName() string {
	return "sale_order"
}

func (so *SaleOrder) Columns() []string {
	return []string{
		"id",
		"order_no",
		"create_time",
		"user_id",
		"status",
		"receiver",
		"province_id",
		"city_id",
		"district_id",
		"address",
		"phone_no",
		"supplier_id",
		"order_amt",
		"goods_amt",
		"express_method",
		"express_order_no",
		"express_fee",
	}
}

func (so *SaleOrder) RegionIDs() []int {
	return []int{so.ProvinceID, so.CityID, so.DistrictID}
}

// IsMaster tells wether a sale order is master
func (so *SaleOrder) IsMaster() bool {
	return so.ParentID == 0
}
