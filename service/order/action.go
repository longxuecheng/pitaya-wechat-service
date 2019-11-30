package order

import (
	"gotrue/model"
	"gotrue/service/api"
)

const (
	actionPay          action = "pay"
	actionSend         action = "send"
	actionCancel       action = "cancel"
	actionApplyRefund  action = "apply_refund"
	actionRemindSent   action = "remind_sent"
	actionApplySupport action = "apply_support"
	actionRemove       action = "remove"
	actionReceive      action = "receive"
	actionExpress      action = "express"
	actionComment      action = "comment"
)

var (
	defaultActionMapper = actionMapper{
		model.Created:   actions{actionPay, actionCancel},
		model.Paid:      actions{actionRemindSent},
		model.PayFailed: actions{actionPay, actionCancel},
		model.Sent:      actions{actionReceive, actionExpress}, // 暂时不让申请退款
		// model.Sent:      actions{actionReceive, actionApplyRefund},
		model.Finish: actions{actionComment, actionApplySupport, actionExpress},
	}
)

type action string

type actions []action

func (acts actions) ActionMap() api.SaleActions {
	saleActions := api.SaleActions{}
	for _, act := range acts {
		saleActions[act.String()] = true
	}
	return saleActions
}

func (a action) String() string {
	return string(a)
}

type actionMapper map[model.OrderStatus]actions

func (m actionMapper) push(status model.OrderStatus, actions ...action) {
	m[status] = append(m[status], actions...)
}

func (m actionMapper) getAPIActions(status model.OrderStatus) api.SaleActions {
	return m[status].ActionMap()
}
