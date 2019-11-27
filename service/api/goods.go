package api

import (
	"gotrue/dto/response"
	"io"
)

type IGoodsService interface {
	GetGoodsByCategory(categoryID int64) ([]*response.GoodsItem, error)

	SpecificationDesc(goodsID int64, specIDs []int64, sep string) (string, error)

	GetInternalGoodsByCategory(categoryID int64) ([]*InternalGoods, error)

	ExpressConstraint(stockID int64, provinceID int64) (*response.GoodsExpressConstraint, error)

	GoodsInfo(goodsID int64) (*GoodsInfoRespone, error)

	Attributes(goodsID int64) ([]*response.AttributeDTO, error)

	Specifications(goodsID int64) ([]*response.GoodsSpecificationDTO, error)

	GoodsList(categoryID int64) ([]IndexGoodsResponse, error)

	OneSaleGoodsCards() ([]*response.GoodsCard, error)

	ImportExpressConstraints(stockID int64, r io.Reader) error
}

type InternalGoods struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type StockPrice struct {
	Spec      string `json:"spec"`
	SalePrice string `json:"salePrice"`
}

type IndexGoodsResponse struct {
	ID            int64        `json:"id"`
	Status        string       `json:"status"`
	Available     bool         `json:"available"`
	Name          string       `json:"name"`
	PicURL        string       `json:"picUrl"`
	ProducingArea string       `json:"producingArea"`
	StockPrices   []StockPrice `json:"stockPrices"`
}

type GoodsInfoRespone struct {
	ID            int64  `json:"id"`
	SupplierID    int64  `json:"supplierId"`
	Name          string `json:"name"`
	BriefDesc     string `json:"briefDesc"`
	Status        string `json:"status"`
	MinPrice      string `json:"minPrice"`
	MaxPrice      string `json:"maxPrice"`
	Available     bool   `json:"available"`
	Description   string `json:"description"`
	ListPicURL    string `json:"listPicUrl"`
	ProducingArea string `json:"producingArea"`
}
