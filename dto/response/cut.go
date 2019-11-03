package response

import "github.com/shopspring/decimal"

type CutOrder struct {
	ID                int64           `json:"id"`
	UserID            int64           `json:"userId"`
	CutNo             string          `json:"cutNo"`
	GoodsID           int64           `json:"goodsId"`
	StockID           int64           `json:"stockId"`
	OriginPriceString string          `json:"originPrice"`
	CutoffPriceString string          `json:"cutoffPrice"`
	CutoffPrice       decimal.Decimal `json:"-"`
	CreateTime        string          `json:"createTime"`
	ExpireTime        string          `json:"expireTime"`
	Items             []*CutDetail    `json:"items"`
}

func (c *CutOrder) SetCutoffPrice(cutoff decimal.Decimal) {
	c.CutoffPriceString = cutoff.StringFixed(2)
	c.CutoffPrice = cutoff
}

type CutDetail struct {
	CutPrice decimal.Decimal `json:"cutPrice"`
	Avartar  string          `json:"avartar"`
	NickName string          `json:"nickName"`
	CutTime  string          `json:"cutTime"`
}
