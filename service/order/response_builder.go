package order

import (
	"gotrue/facility/utils"
	"gotrue/model"
	"gotrue/service/api"
)

type responseBuilder struct {
	orders    model.SaleOrderList
	orderIDs  []int64
	goodsList []*model.SaleDetail
}

func newResponseBuilder() *responseBuilder {
	return &responseBuilder{}
}

func (b *responseBuilder) setOrders(l model.SaleOrderList) *responseBuilder {
	b.orders = l
	return b
}

func (b *responseBuilder) setDetails(details []*model.SaleDetail) *responseBuilder {
	b.goodsList = details
	return b
}

func (b *responseBuilder) buildList() []*api.SaleOrderResponse {
	dtos := make([]*api.SaleOrderResponse, len(b.orders))
	for i, model := range b.orders {
		dto := installInfo(model)
		goodsList := []api.SaleDetailResponse{}
		for _, goods := range b.goodsList {
			if model.ID == goods.OrderID {
				goodsList = append(goodsList, installDetail(goods))
				break
			}
		}
		dto.Details = goodsList
		dtos[i] = dto
	}
	return dtos
}

func installInfo(order *model.SaleOrder) *api.SaleOrderResponse {
	data := &api.SaleOrderResponse{}
	data.ID = order.ID
	data.OrderNo = order.OrderNo
	data.Status = order.Status.Name()
	data.CreatedAt = utils.FormatTime(order.CreateTime, utils.TimePrecision_Seconds)
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
	data.Actions = defaultActionMapper.getAPIActions(order.Status)
	return data
}

func installDetail(model *model.SaleDetail) api.SaleDetailResponse {
	dto := api.SaleDetailResponse{}
	dto.ID = model.ID
	dto.GoodsName = model.GoodsName
	dto.Quantity = model.Quantity
	dto.RetailPrice = model.SaleUnitPrice
	dto.ListPicURL = model.ListPicURL.String
	dto.GoodsSpecDescription = model.GoodsSpecDescription
	return dto
}
