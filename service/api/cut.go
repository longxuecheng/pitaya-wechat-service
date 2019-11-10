package api

import (
	"context"
	"gotrue/dto/request"
	"gotrue/dto/response"

	"github.com/shopspring/decimal"
)

type ICutService interface {
	MyActivatedCutOrder(req *request.CutOrder) (*response.CutOrder, error)

	GetCutoffInfo(userID string, goodsID int64) (*response.CutOrder, error)

	CutDetails(cutOrderID int64) (decimal.Decimal, []*response.CutDetail, error)

	AssistCutoff(ctx context.Context, req *request.AssistCutoff) (*response.CutOrder, error)

	ConsumeCutOrder(ctx context.Context, req *request.ConsumeCutOrder) error
}
