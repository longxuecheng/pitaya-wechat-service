package order

import (
	"gotrue/dto/response"
	"gotrue/facility/errors"
	"gotrue/model"

	"github.com/shopspring/decimal"
)

var one = decimal.NewFromFloat32(1.0)

type StockOrder struct {
	SaleOrder   *model.SaleOrder
	SaleDetails []*model.SaleDetail
}

func (s *StockOrder) AddDetail(detail *model.SaleDetail) {
	if s.SaleDetails == nil {
		s.SaleDetails = []*model.SaleDetail{}
		s.SaleDetails = append(s.SaleDetails, detail)
	} else {
		s.SaleDetails = append(s.SaleDetails, detail)
	}
}

type StockOrderBuilder struct {
	Quantity       decimal.Decimal
	Goods          *model.Goods
	Stock          *model.Stock
	UserID         int64
	UnitExpressFee decimal.Decimal
	Address        *response.UserAddress
}

func (s *StockOrderBuilder) validate() error {
	if s.Address == nil || s.Address.ID == 0 {
		return errors.NewWithCodef("AddressNotSelected", "请选择下单的地址")
	}
	if s.UserID == 0 {
		return errors.NewWithCodef("UserNotExist", "无归属的订单")
	}
	return nil
}

func (s *StockOrderBuilder) BuildSingleOrder() (*model.SaleOrder, error) {
	err := s.validate()
	if err != nil {
		return nil, err
	}
	sale_order := &model.SaleOrder{}
	sale_order.ParentID = 0
	sale_order.UserID = s.UserID
	sale_order.SupplierID = s.Goods.SupplierID
	orderNo, err := generateOrderNumber(1)
	if err != nil {
		return nil, err
	}
	goods_price := s.Stock.SaleUnitPrice.Mul(s.Quantity)
	express_fee := decimal.Zero
	if !s.splitable() {
		express_fee = s.UnitExpressFee
	} else {
		express_fee = s.UnitExpressFee.Mul(s.Quantity)
	}
	order_price := goods_price.Add(express_fee)
	cost_price := s.Stock.CostUnitPrice.Mul(s.Quantity)
	sale_order.OrderNo = orderNo
	sale_order.PhoneNo = s.Address.Mobile
	sale_order.DistrictID = s.Address.DistrictID
	sale_order.ProvinceID = s.Address.ProvinceID
	sale_order.CityID = s.Address.CityID
	sale_order.Address = s.Address.Address
	sale_order.Receiver = s.Address.Name
	sale_order.GoodsAmt = goods_price
	sale_order.ExpressFee = express_fee
	sale_order.OrderAmt = order_price
	sale_order.CostAmt = cost_price
	return sale_order, nil
}

func (s *StockOrderBuilder) splitable() bool {
	return s.Stock.Splitable
}

func (s *StockOrderBuilder) Build() ([]*StockOrder, error) {
	stock_orders := []*StockOrder{}
	if !s.splitable() || s.Quantity.Equals(one) {
		order, err := s.BuildSingleOrder()
		if err != nil {
			return nil, err
		}
		stock_order := &StockOrder{
			SaleOrder: order,
		}
		stock_order.AddDetail(s.BuildSaleDetail())
		stock_orders = append(stock_orders, stock_order)
	} else {
		for i := 0; i < int(s.Quantity.IntPart()); i++ {
			sb := &StockOrderBuilder{
				Goods:          s.Goods,
				Quantity:       one,
				Stock:          s.Stock,
				UserID:         s.UserID,
				UnitExpressFee: s.UnitExpressFee,
				Address:        s.Address,
			}
			order, err := sb.BuildSingleOrder()
			if err != nil {
				return nil, err
			}
			stock_order := &StockOrder{
				SaleOrder: order,
			}
			stock_order.AddDetail(sb.BuildSaleDetail())
			stock_orders = append(stock_orders, stock_order)
		}
	}
	return stock_orders, nil
}

func (s *StockOrderBuilder) BuildSaleDetail() *model.SaleDetail {
	return &model.SaleDetail{
		OrderID:       0,
		GoodsID:       s.Stock.GoodsID,
		GoodsName:     s.Goods.Name,
		Quantity:      s.Quantity,
		StockID:       s.Stock.ID,
		SaleUnitPrice: s.Stock.SaleUnitPrice,
		CostUnitPrice: s.Stock.CostUnitPrice,
		GoodsSpecIDs:  s.Stock.Specification.String,
		ListPicURL:    s.Goods.ListPicURL,
	}

}
