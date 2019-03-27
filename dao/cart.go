package dao

import (
	"pitaya-wechat-service/model"
	"pitaya-wechat-service/sys"

	sq "github.com/Masterminds/squirrel"
)

// CartDaoSingleton is a singleton of goods dao
var CartDaoSingleton *CartDao

func init() {
	CartDaoSingleton = new(CartDao)
	CartDaoSingleton.db = sys.GetEasyDB()
}

var columns_cart_all = []string{"id", "user_id", "supplier_id", "session_id", "goods_id", "goods_sn", "stock_id", "goods_name", "quantity", "goods_spec_description", "goods_spec_ids", "checked", "list_pic_url"}

// CartDao is dao
type CartDao struct {
	db *sys.EasyDB
}

func (dao *CartDao) AddCart(setMap map[string]interface{}) (id int64, err error) {
	_, id, err = dao.db.Insert(model.Table_Cart, setMap)
	return
}

func (dao *CartDao) SelectByUserID(userID int64) ([]model.Cart, error) {
	carts := []model.Cart{}
	err := dao.db.SelectDSL(&carts, columns_cart_all, model.Table_Cart, sq.Eq{"user_id": userID})
	if err != nil {
		return nil, err
	}
	return carts, nil
}

func (dao *CartDao) SelectChecked(userID int64) ([]model.Cart, error) {
	carts := []model.Cart{}
	err := dao.db.SelectDSL(&carts, columns_cart_all, model.Table_Cart, sq.Eq{"user_id": userID, "checked": true})
	if err != nil {
		return nil, err
	}
	return carts, nil
}

func (dao *CartDao) SelectCountByUserID(userID int64) (total int64, err error) {
	count := new(model.Count)
	err = dao.db.SelectOneDSL(count, []string{"count(1) as count"}, model.Table_Cart, sq.Eq{"user_id": userID})
	total = count.Count
	return
}

func (dao *CartDao) UpdateByID(ID int64, updateMap map[string]interface{}) error {

	_, err := dao.db.Update(model.Table_Cart, updateMap, sq.Eq{"id": ID})
	if err != nil {
		return err
	}
	return nil
}

func (dao *CartDao) UpdateByStockID(stockID int64, updateMap map[string]interface{}) error {

	_, err := dao.db.Update(model.Table_Cart, updateMap, sq.Eq{"stock_id": stockID})
	if err != nil {
		return err
	}
	return nil
}
