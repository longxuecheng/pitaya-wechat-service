package order

import (
	"gotrue/facility/errors"
	"gotrue/model"

	"github.com/looplab/fsm"
)

var (
	ErrorOrderInvalid = errors.NewWithCodef("OrderInvalid", "当前订单不支持该操作")
)

type saleOrderFSM struct {
	fsm *fsm.FSM
}

func newOrderFSM(order *model.SaleOrder) *saleOrderFSM {
	return &saleOrderFSM{
		fsm: fsm.NewFSM(
			order.Status.String(),
			fsm.Events{
				{Name: actionPay.String(), Src: []string{model.Created.String(), model.PayFailed.String(), model.Paying.String()}, Dst: model.Paid.String()},
				{Name: actionCancel.String(), Src: []string{model.Created.String()}, Dst: model.Cancel.String()},
				{Name: actionSend.String(), Src: []string{model.Paid.String(), model.Sent.String()}, Dst: model.Sent.String()},
				{Name: actionReceive.String(), Src: []string{model.Sent.String()}, Dst: model.Finish.String()},
				{Name: actionApplySupport.String(), Src: []string{model.Finish.String()}, Dst: model.PostSaleFinished.String()},
				{Name: actionApplyRefund.String(), Src: []string{model.Sent.String(), model.Paid.String()}, Dst: model.Refund.String()},
			},
			fsm.Callbacks{},
		),
	}
}

func (sof *saleOrderFSM) can(action action) error {
	if !sof.fsm.Can(action.String()) {
		return ErrorOrderInvalid
	}
	return nil
}

func (sof *saleOrderFSM) current() model.OrderStatus {
	return model.OrderStatus(sof.fsm.Current())
}
