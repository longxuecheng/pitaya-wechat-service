package test

import (
	"fmt"
	"gotrue/dao"
	"testing"
)

func TestGoodsExpressConstraint(t *testing.T) {
	dao.InitGoodsExpressConstraintDao()
	list, err := dao.GoodsExpressConstraintDao.QueryByGoodsID(1)
	if err != nil {
		t.Error(err)
	}
	for _, item := range list {
		fmt.Printf("province %d is free %v express fee %s\n", item.ProvinceID, item.IsFree, item.ExpressFee.String())
	}
}
