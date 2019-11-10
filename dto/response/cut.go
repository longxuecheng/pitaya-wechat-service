package response

import "github.com/shopspring/decimal"

type CutOrder struct {
	ID                int64           `json:"id"`
	Avatar            string          `json:"avatar"`
	NickName          string          `json:"nickName"`
	CutNo             string          `json:"cutNo"`
	GoodsID           int64           `json:"goodsId"`
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
	Avatar   string          `json:"avatar"`
	NickName string          `json:"nickName"`
	CutTime  string          `json:"cutTime"`
}
