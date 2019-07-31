package dao

import (
	"fmt"
	"gotrue/model"
	"gotrue/sys"
	"strings"
)

// GoodsImgDao is a singleton of goods dao
var GoodsImgDao *GoodsImg

func initGoodsImageDao() {
	GoodsImgDao = &GoodsImg{
		db: sys.GetEasyDB(),
	}
}

var columns_goods_img = []string{"id", "goods_id", "name", "path", "display_order"}

// GoodsImg is dao
type GoodsImg struct {
	db *sys.EasyDB
}

func (dao *GoodsImg) SelectByGoodsID(goodsID int64) ([]*model.GoodsImg, error) {
	imgs := []*model.GoodsImg{}
	err := dao.db.Select(&imgs, fmt.Sprintf("SELECT %s FROM goods_photo WHERE goods_id = ? ORDER BY display_order ASC", strings.Join(columns_goods_img, ",")), goodsID)
	if err != nil {
		return nil, err
	}
	return imgs, nil
}
