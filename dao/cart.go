package dao

import (
	"gotrue/facility/utils"
	"gotrue/model"
	"gotrue/sys"

	sq "github.com/Masterminds/squirrel"
)

// CartDao is a singleton of goods dao
var CartDao *Cart

func initCartDao() {
	CartDao = &Cart{
		db: sys.GetEasyDB(),
	}
}

var columns_cart_all = []string{"id", "user_id", "supplier_id", "session_id", "goods_id", "goods_sn", "stock_id", "goods_name", "quantity", "goods_spec_description", "goods_spec_ids", "checked", "list_pic_url"}

// Cart is dao
type Cart struct {
	db *sys.EasyDB
}

func (dao *Cart) AddCart(cart model.Cart) (id int64, err error) {
	setMap := utils.StructToMap(cart, "db", "exclude")
	_, id, err = dao.db.Insert(model.Table_Cart, setMap)
	return
}

func (dao *Cart) SelectByUserID(userID int64) ([]model.Cart, error) {
	carts := []model.Cart{}
	err := dao.db.SelectDSL(&carts, columns_cart_all, model.Table_Cart, sq.Eq{"user_id": userID})
	if err != nil {
		return nil, err
	}
	return carts, nil
}

func (dao *Cart) SelectChecked(userID int64) ([]model.Cart, error) {
	carts := []model.Cart{}
	err := dao.db.SelectDSL(&carts, columns_cart_all, model.Table_Cart, sq.Eq{"user_id": userID, "checked": true})
	if err != nil {
		return nil, err
	}
	return carts, nil
}

func (dao *Cart) SelectCountByUserID(userID int64) (total int64, err error) {
	count := new(model.Count)
	err = dao.db.SelectOneDSL(count, []string{"count(1) as count"}, model.Table_Cart, sq.Eq{"user_id": userID})
	total = count.Count
	return
}

func (dao *Cart) UpdateByID(ID int64, updateMap map[string]interface{}) error {

	_, err := dao.db.Update(model.Table_Cart, updateMap, sq.Eq{"id": ID})
	if err != nil {
		return err
	}
	return nil
}

func (dao *Cart) UpdateByStockID(stockID int64, updateMap map[string]interface{}) error {

	_, err := dao.db.Update(model.Table_Cart, updateMap, sq.Eq{"stock_id": stockID})
	if err != nil {
		return err
	}
	return nil
}
