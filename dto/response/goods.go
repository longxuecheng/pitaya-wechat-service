package response

import (
	"gotrue/facility/errors"

	"github.com/shopspring/decimal"
)

type AttributeDTO struct {
	ID    int64  `json:"-"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type BusinessResult struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (r *BusinessResult) IsOK() bool {
	return r.Code == "OK"
}

func (r *BusinessResult) Error() error {
	if !r.IsOK() {
		return errors.NewWithCodef(r.Code, r.Message)
	}
	return nil
}

type GoodsItem struct {
	ID            int64           `json:"id"`
	Name          string          `json:"name"`
	PicURL        string          `json:"picUrl"`
	RetailPrice   decimal.Decimal `json:"retailPrice"`
	ProducingArea string          `json:"producingArea"`
}

type GoodsCard struct {
	ID            int64           `json:"id"`
	Name          string          `json:"name"`
	CardPicURL    string          `json:"cardPicUrl"`
	Description   string          `json:"description"`
	RetailPrice   decimal.Decimal `json:"retailPrice"`
	Status        string          `json:"status"`
	ProducingArea string          `json:"producingArea"`
	Location      Location        `json:"location"`
}

type Location struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

type GoodsExpressConstraint struct {
	IsFree          bool            `json:"isFree"`
	UnitExpressFee  decimal.Decimal `json:"unitExpressFee"`
	TotalExpressFee decimal.Decimal `json:"expressFee"`
	BusinessResult
}

func (g *GoodsExpressConstraint) CalculateTotalExpressFee(quantity decimal.Decimal) {
	g.TotalExpressFee = g.UnitExpressFee.Mul(quantity)
}

type GoodsGalleryDTO struct {
	ID        int64  `json:"id"`
	GoodsID   int64  `json:"goodsId"`
	ImgURL    string `json:"imgUrl"`
	ImgDesc   string `json:"imgDesc"`
	SortOrder int    `json:"sortOrder"`
}
