package model

import (
	"gotrue/dto/response"
	"gotrue/facility/utils"
	"math/rand"
	"sync"
	"time"

	"/manage/paperwork-facility/reflect_util"

	"github.com/shopspring/decimal"
)

type CutOrder struct {
	ID          int64        `db:"id" insert:"true" pk:"true"`
	UserID      int64        `db:"user_id"`
	CutNo       string       `db:"cut_no"`
	GoodsID     int64        `db:"goods_id"`
	StockID     int64        `db:"stock_id"`
	CreateTime  time.Time    `db:"create_time" insert:"true"`
	ExpireTime  time.Time    `db:"expire_time"`
	Consumed    bool         `db:"consumed" insert:"true"`
	SaleOrderID int64        `db:"sale_order_id"`
	ConsumeTime NullUTC8Time `db:"consume_time"`
}

func (c *CutOrder) TableName() string {
	return "cut_order"
}

func (c *CutOrder) Columns() []string {
	return reflect_util.TagValues(c, "db")
}

func (c *CutOrder) ResponseCutOrder() *response.CutOrder {
	return &response.CutOrder{
		ID:         c.ID,
		CutNo:      c.CutNo,
		GoodsID:    c.GoodsID,
		CreateTime: utils.FormatTime(c.CreateTime, utils.TimePrecision_Seconds),
		ExpireTime: utils.FormatTime(c.ExpireTime, utils.TimePrecision_Seconds),
	}
}

type CutDetail struct {
	ID         int64           `db:"id" insert:"true"`
	CutOrderID int64           `db:"cut_order_id"`
	UserID     int64           `db:"user_id"`
	CutPrice   decimal.Decimal `db:"cut_price"`
	CreateTime time.Time       `db:"create_time" insert:"true"`
}

func (c *CutDetail) TableName() string {
	return "cut_detail"
}

func (c *CutDetail) Columns() []string {
	return reflect_util.TagValues(c, "db")
}

func (c *CutDetail) ResponseCutDetail(u *User) *response.CutDetail {
	detail := &response.CutDetail{
		CutPrice: c.CutPrice,
		CutTime:  utils.FormatTime(c.CreateTime, utils.TimePrecision_Seconds),
	}
	if u != nil {
		detail.Avatar = u.AvatarURL
		detail.NickName = u.NickName
	}
	return detail
}

type CutDetailSet struct {
	Items []*CutDetail
	size  int
	total decimal.Decimal
}

func NewCutDetailSet(items []*CutDetail) *CutDetailSet {
	total := decimal.Zero
	for _, item := range items {
		total = total.Add(item.CutPrice)
	}
	return &CutDetailSet{
		Items: items,
		total: total,
		size:  len(items),
	}
}

func (c *CutDetailSet) TotalCutoffPrice() decimal.Decimal {
	return c.total
}

func (c *CutDetailSet) UserIDs() []int64 {
	userIDs := []int64{}
	for _, detail := range c.Items {
		userIDs = append(userIDs, detail.UserID)
	}
	return userIDs
}

func (c *CutDetailSet) Size() int {
	return c.size
}

type CutCalculator struct {
	remain   decimal.Decimal
	cutoff   decimal.Decimal
	current  decimal.Decimal
	discount decimal.Decimal
	mutex    sync.Mutex
	cutCount int
}

func NewCutCalculator(low, origin, discount decimal.Decimal) *CutCalculator {
	return &CutCalculator{
		remain:   origin.Sub(discount).Sub(low),
		discount: discount,
		current:  origin.Sub(discount),
		cutoff:   decimal.Zero,
	}
}

func (c *CutCalculator) RemainPrice() decimal.Decimal {
	return c.remain
}

func (c *CutCalculator) CurrentPrice() decimal.Decimal {
	return c.current.Sub(c.cutoff)
}

func (c *CutCalculator) RandomCut() decimal.Decimal {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.remain.Equal(c.cutoff) {
		c.cutoff = c.remain
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var randomInt64 int64
	for {
		randomInt64 = r.Int63n(100)
		if randomInt64 > 0 {
			break
		}
	}
	mulFactor := decimal.New(randomInt64, 0)
	cutoffPrice := decimal.NewFromFloat(0.01).Mul(mulFactor)
	c.cutoff = c.cutoff.Add(cutoffPrice)
	c.cutCount++
	if c.remain.LessThan(c.cutoff) {
		c.cutoff = c.remain
		return decimal.Zero
	}
	return cutoffPrice
}

func (c *CutCalculator) TotalCutoff() decimal.Decimal {
	return c.cutoff
}
