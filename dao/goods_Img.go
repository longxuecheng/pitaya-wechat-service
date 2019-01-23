package dao

import (
	"fmt"
	"pitaya-wechat-service/model"
	"pitaya-wechat-service/sys"
	"strings"
)

// GoodsImgDaoSingleton is a singleton of goods dao
var GoodsImgDaoSingleton *GoodsImgDao

func init() {
	GoodsImgDaoSingleton = new(GoodsImgDao)
	GoodsImgDaoSingleton.db = sys.GetEasyDB()
}

var columns_goods_img = []string{"id", "goods_id", "name", "path", "display_order"}

// GoodsImgDao is dao
type GoodsImgDao struct {
	db *sys.EasyDB
}

func (dao *GoodsImgDao) SelectByGoodsID(goodsID int64) ([]*model.GoodsImg, error) {
	imgs := []*model.GoodsImg{}
	err := dao.db.Select(&imgs, fmt.Sprintf("SELECT %s FROM goods_photo WHERE goods_id = ? ORDER BY display_order ASC", strings.Join(columns_goods_img, ",")), goodsID)
	if err != nil {
		return nil, err
	}
	return imgs, nil
}
