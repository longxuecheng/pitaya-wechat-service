package test

import (
	"pitaya-wechat-service/dao"
	"pitaya-wechat-service/sys"
	"testing"
)

func TestDBInsert(t *testing.T) {
	edb := sys.GetEasyDB()
	setMap := map[string]interface{}{
		"user_id":  1,
		"goods_id": 1,
		"goods_sn": "SN-01001",
		"stock_id": 100,
	}
	ra, id, err := edb.Insert("cart", setMap)
	if err != nil {
		t.Error(err)
	}
	t.Logf("affected rows %d , last inserted id = %d ", ra, id)
}

func TestDbCount(t *testing.T) {
	total, err := dao.CartDaoSingleton.SelectCountByUserID(0)
	if err != nil {
		t.Error(err)
	}
	t.Logf("total cart items was %d", total)
}

func TestDbUpdate(t *testing.T) {
	setMap := map[string]interface{}{
		"checked": 0,
	}
	err := dao.CartDaoSingleton.UpdateByID(1, setMap)
	if err != nil {
		t.Error(err)
	}

}
