package model

import (
	"database/sql/driver"
	"gotrue/facility/utils"
	"time"

	"github.com/shopspring/decimal"
)

const (
	Created          OrderStatus = "CREATED"
	Cancel           OrderStatus = "CANCEL"
	Paid             OrderStatus = "PAID"
	Paying           OrderStatus = "PAYING"
	PayFailed        OrderStatus = "PAY_FAILED"
	Sent             OrderStatus = "SENT"
	Refunding        OrderStatus = "REFUNDING"
	RefundRefused    OrderStatus = "REFUND_REFUSED"
	Refund           OrderStatus = "REFUND"
	PostSale         OrderStatus = "PST_SALE"
	PostSaleFinished OrderStatus = "PST_SALE_FNS"
	Finish           OrderStatus = "FINISHED"
	Closed           OrderStatus = "CLOSED"
)

const (
	DiscountTypeNone   DiscountType = "None"
	DiscountTypeCutoff DiscountType = "Cutoff"
	DiscountTypeScore  DiscountType = "Score"
)

var orderStatusMap = map[OrderStatus]string{
	Created:          "待付款",
	Paying:           "支付中",
	Cancel:           "已取消",
	Paid:             "已付款", // 待收货
	PayFailed:        "付款失败",
	Sent:             "待收货",
	Refunding:        "退款中",
	RefundRefused:    "拒绝退款",
	Refund:           "退款成功",
	PostSale:         "售后处理中",
	PostSaleFinished: "售后完成",
	Finish:           "交易完成",
}

type OrderStatus string

func (os OrderStatus) Value() (driver.Value, error) {
	return string(os), nil
}

func (os OrderStatus) String() string {
	return string(os)
}

func (so OrderStatus) Name() string {
	if statusName, ok := orderStatusMap[so]; ok {
		return statusName
	}
	return "未知"
}

type DiscountType string

func (d DiscountType) Value() (driver.Value, error) {
	return string(d), nil
}

type SaleOrderList []*SaleOrder

func (l SaleOrderList) IDList() []int64 {
	ids := make([]int64, len(l))
	for i, o := range l {
		ids[i] = o.ID
	}
	return ids
}

func (l SaleOrderList) TotalPrice() decimal.Decimal {
	sum := decimal.Zero
	for _, order := range l {
		sum = sum.Add(order.OrderAmt)
	}
	return sum
}

type SaleOrder struct {
	ID            int64           `db:"id" insert:"false" pk:"true"`
	ParentID      int64           `db:"parent_id"`
	OrderNo       string          `db:"order_no"`
	CreateTime    time.Time       `db:"create_time" insert:"false"`
	UserID        int64           `db:"user_id"`
	Status        OrderStatus     `db:"status" insert:"false"`
	Receiver      string          `db:"receiver"`
	ProvinceID    int             `db:"province_id"`
	CityID        int             `db:"city_id"`
	DistrictID    int             `db:"district_id"`
	Address       string          `db:"address"`
	PhoneNo       string          `db:"phone_no"`
	SupplierID    int64           `db:"supplier_id"`
	DiscountAmt   decimal.Decimal `db:"discount_amt"`
	DiscountType  DiscountType    `db:"discount_type"`
	OrderAmt      decimal.Decimal `db:"order_amt"`
	GoodsAmt      decimal.Decimal `db:"goods_amt"`
	CostAmt       decimal.Decimal `db:"cost_amt"`
	ExpressMethod *string         `db:"express_method"`
	ExpressNo     *string         `db:"express_order_no"`
	ExpressFee    decimal.Decimal `db:"express_fee"`
	SettlementID  int64           `db:"settlement_id"`
	Count         int64           `db:"count" count:"true"`
}

func (so *SaleOrder) TableName() string {
	return "sale_order"
}

func (so *SaleOrder) Columns() []string {
	return utils.TagValues(so, "db", "count")
}

func (so *SaleOrder) OrderNo12() string {
	return so.OrderNo[0:10] + ".."
}

func (so *SaleOrder) RegionIDs() []int {
	return []int{so.ProvinceID, so.CityID, so.DistrictID}
}

// IsMaster tells wether a sale order is master
func (so *SaleOrder) IsMaster() bool {
	return so.ParentID == 0
}

type SaleOrderSet struct {
	Items          []*SaleOrder
	costPrice      decimal.Decimal
	profitPrice    decimal.Decimal
	totalSalePrice decimal.Decimal
	settlePrice    decimal.Decimal
}

func (s *SaleOrderSet) Sum() {
	totalPrice := decimal.Zero
	goodsPrice := decimal.Zero
	costPrice := decimal.Zero
	settlePrice := decimal.Zero
	expressPrice := decimal.Zero
	for _, item := range s.Items {
		totalPrice = totalPrice.Add(item.OrderAmt)
		goodsPrice = goodsPrice.Add(item.GoodsAmt)
		costPrice = costPrice.Add(item.CostAmt)
		expressPrice = expressPrice.Add(item.ExpressFee)
	}
	// express fee + cost price is the money we should pay for supplier
	settlePrice = costPrice.Add(expressPrice)
	profitPrice := totalPrice.Sub(settlePrice)
	s.totalSalePrice = totalPrice
	s.costPrice = costPrice
	s.settlePrice = settlePrice
	s.profitPrice = profitPrice
}

func (s *SaleOrderSet) CostPrice() decimal.Decimal {
	return s.costPrice
}

func (s *SaleOrderSet) SalePrice() decimal.Decimal {
	return s.totalSalePrice
}

func (s *SaleOrderSet) ProfitPrice() decimal.Decimal {
	return s.profitPrice
}

func (s *SaleOrderSet) SettlePrice() decimal.Decimal {
	return s.settlePrice
}

func (s *SaleOrderSet) Size() int {
	return len(s.Items)
}
