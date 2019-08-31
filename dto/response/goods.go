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
	Message string `json:"messsage"`
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

type GoodsItemDTO struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	PicURL      string          `json:"picUrl"`
	RetailPrice decimal.Decimal `json:"retailPrice"`
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

type GoodsInfoDTO struct {
	ID          int64           `json:"id"`
	SupplierID  int64           `json:"supplierId"`
	Name        string          `json:"name"`
	BriefDesc   string          `json:"briefDesc"`
	RetailPrice decimal.Decimal `json:"retailPrice"`
	Unit        string          `json:"unit"`
	Description string          `json:"description"`
	ListPicURL  string          `json:"listPicUrl"`
}

type GoodsGalleryDTO struct {
	ID        int64  `json:"id"`
	GoodsID   int64  `json:"goodsId"`
	ImgURL    string `json:"imgUrl"`
	ImgDesc   string `json:"imgDesc"`
	SortOrder int    `json:"sortOrder"`
}
