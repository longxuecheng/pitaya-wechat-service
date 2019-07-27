package model

import (
	"database/sql/driver"
	"time"

	"github.com/shopspring/decimal"
)

type TransactionType string

func (tt TransactionType) Value() (driver.Value, error) {
	return string(tt), nil
}

const (
	TransactionTypePay    TransactionType = "D"
	TransactionTypeRefund TransactionType = "C"
)

type WechatPayment struct {
	ID             int64           `db:"id"`
	SaleOrderID    int64           `db:"sale_order_id"`
	SaleOrderNo    string          `db:"sale_order_no"`
	TransactionID  *string         `db:"transaction_id"`
	TransationType TransactionType `db:"transaction_type"`
	Amount         decimal.Decimal `db:"amount"`
	Status         string          `db:"status"`
	CreateTime     time.Time       `db:"create_time"`
	Description    *string         `db:"description"`
	UpdateTime     time.Time       `db:"update_time" omitinsert:"true"`
}

func (payment *WechatPayment) TableName() string {
	return "wechat_payment"
}
