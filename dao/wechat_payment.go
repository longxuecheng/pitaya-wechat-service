package dao

import (
	"database/sql"
	"gotrue/facility/utils"
	"gotrue/model"
	"gotrue/sys"

	sq "github.com/Masterminds/squirrel"
)

// WechatPaymentDao is a singleton of goods dao
var WechatPaymentDao *WechatPayment

func initWechatPaymentDao() *WechatPayment {
	wp := &model.WechatPayment{}
	WechatPaymentDao = &WechatPayment{
		db:        sys.GetEasyDB(),
		columns:   utils.TagValues(wp, "db"),
		tableName: wp.TableName(),
	}
	return WechatPaymentDao
}

// WechatPayment is dao
type WechatPayment struct {
	db        *sys.EasyDB
	columns   []string
	tableName string
}

func (dao *WechatPayment) Create(order *model.WechatPayment, tx *sql.Tx) (id int64, err error) {
	setMap := utils.InsertMap(order, "db")
	_, id, err = dao.db.Insert(dao.tableName, setMap, tx)
	return
}

func (dao *WechatPayment) UpdateByID(id int64, updateMap map[string]interface{}, tx *sql.Tx) error {
	_, err := dao.db.UpdateTx(tx, dao.tableName, updateMap, sq.Eq{"id": id})
	return err
}

func (dao *WechatPayment) SelectByOrderNo(orderNo string, txnType model.TransactionType) ([]*model.WechatPayment, error) {
	txns := []*model.WechatPayment{}
	err := dao.db.SelectDSL(&txns, dao.columns, dao.tableName,
		sq.Eq{"sale_order_no": orderNo, "transaction_type": txnType})
	return txns, err
}
