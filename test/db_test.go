package test

import (
	"database/sql"
	"gotrue/dao"
	"gotrue/sys"
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

func TestDBTx(t *testing.T) {
	edb := sys.GetEasyDB()
	setMap1 := map[string]interface{}{
		"user_id":  1,
		"goods_id": 1,
		"goods_sn": "SN-01005",
		"stock_id": 100,
	}
	setMap2 := map[string]interface{}{
		"user_id":  1,
		"goods_id": 1,
		"goods_sn": "SN-01006",
		"stock_id": 100,
	}

	edb.ExecTx(func(tx *sql.Tx) error {
		ra, id, err := edb.Insert("cart", setMap1, tx)
		if err != nil {
			return err
		}
		// err = errors.New("raise by me")
		// if err != nil {
		// 	return err
		// }
		t.Logf("affected rows %d , last inserted id = %d ", ra, id)
		ra, id, err = edb.Insert("cart", setMap2, tx)
		if err != nil {
			return err
		}
		t.Logf("affected rows %d , last inserted id = %d ", ra, id)
		return nil
	})

}
