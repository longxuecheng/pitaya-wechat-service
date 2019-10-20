package api

import (
	"context"
	"gotrue/dto/request"
	"gotrue/dto/response"

	"github.com/shopspring/decimal"
)

type ICutService interface {
	MyActivatedCutOrder(req *request.CutOrder) (*response.CutOrder, error)

	GetCutOrderByCutNo(cutNo string) (*response.CutOrder, error)

	CutDetails(cutOrderID int64) (decimal.Decimal, []*response.CutDetail, error)

	CreateCutOrder(ctx context.Context, req *request.CutOrder) (*response.CutOrder, error)

	CreateCutDetail(req *request.CreateCutDetail) (decimal.Decimal, error)

	ConsumeCutOrder(ctx context.Context, req *request.ConsumeCutOrder) error
}
