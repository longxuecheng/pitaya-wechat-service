package cut

import (
	"context"
	"gotrue/dao"
	"gotrue/dto/request"
	"gotrue/dto/response"
	"gotrue/facility/errors"
	"gotrue/facility/utils"
	"gotrue/model"
	"gotrue/service/api"
	"time"

	"github.com/google/uuid"

	"github.com/shopspring/decimal"
)

var cutServiceIns api.ICutService

type Cut struct {
	cutOrderDao  *dao.CutOrder
	cutDetailDao *dao.CutDetail
	userDao      *dao.UserDao
	stockDao     *dao.Stock
	goodsSpecDao *dao.GoodsSpec
}

func GetCutService() api.ICutService {
	if cutServiceIns == nil {
		cutServiceIns = NewCutService()
	}
	return cutServiceIns
}

func NewCutService() *Cut {
	return &Cut{
		cutOrderDao:  dao.CutOrderDao,
		cutDetailDao: dao.CutDetailDao,
		userDao:      dao.UserDaoSingleton,
		stockDao:     dao.StockDao,
		goodsSpecDao: dao.GoodsSpecDao,
	}
}

// MyActivatedCutOrder will get an activated cut order for a given stock
func (s *Cut) MyActivatedCutOrder(req *request.CutOrder) (*response.CutOrder, error) {
	cutOrder, err := s.cutOrderDao.QueryByUserAndStock(req.UserID, req.StockID)
	if err != nil {
		return nil, err
	}
	if cutOrder == nil {
		return nil, nil
	}
	stock, err := s.stockDao.SelectByID(req.StockID)
	if err != nil {
		return nil, err
	}
	totalCutoff, apiItems, err := s.CutDetails(cutOrder.ID)
	if err != nil {
		return nil, err
	}
	apiCutOrder := cutOrder.ResponseCutOrder()
	apiCutOrder.Items = apiItems
	apiCutOrder.OriginPriceString = stock.SaleUnitPrice.StringFixed(2)
	apiCutOrder.SetCutoffPrice(totalCutoff)
	return apiCutOrder, nil
}

func (s *Cut) ConsumeCutOrder(ctx context.Context, req *request.ConsumeCutOrder) error {
	cutOrder, err := s.cutOrderDao.QueryByID(req.CutOrderID)
	if err != nil {
		return err
	}
	if cutOrder == nil {
		return nil
	}
	cutOrder.Consumed = true
	cutOrder.SaleOrderID = req.SaleOrderID
	cutOrder.ConsumeTime = model.NullUTC8Time{
		Time:  time.Now(),
		Valid: true,
	}
	tx := utils.GetTx(ctx)
	err = s.cutOrderDao.UpdateCutOrder(cutOrder, tx)
	if err != nil {
		return err
	}
	return nil
}

func (s *Cut) GetCutOrderByCutNo(cutNo string) (*response.CutOrder, error) {
	cutOrder, err := s.cutOrderDao.QueryByCutNo(cutNo)
	if err != nil {
		return nil, err
	}
	if cutOrder == nil {
		return nil, errors.NewWithCodef("InvalidCutOrder", "当前砍价已经结束啦")
	}
	stock, err := s.stockDao.SelectByID(cutOrder.StockID)
	if err != nil {
		return nil, err
	}
	total, apiDetails, err := s.CutDetails(cutOrder.ID)
	if err != nil {
		return nil, err
	}
	goodsSpecs, err := s.goodsSpecDao.SelectByGoodsIDs([]int64{stock.GoodsID})
	if err != nil {
		return nil, err
	}
	goodsSpecMap := goodsSpecs.Map()
	apiCutOrder := cutOrder.ResponseCutOrder()
	apiCutOrder.Items = apiDetails
	apiCutOrder.SpecName = goodsSpecMap.SpecName(stock.SpecIDs())
	apiCutOrder.OriginPriceString = stock.SaleUnitPrice.StringFixed(2)
	apiCutOrder.CurrentPriceString = stock.SaleUnitPrice.Sub(total).StringFixed(2)
	apiCutOrder.SetCutoffPrice(total)
	return apiCutOrder, nil
}

// CutDetails will get detail list for a given cut order
func (s *Cut) CutDetails(cutOrderID int64) (decimal.Decimal, []*response.CutDetail, error) {
	total := decimal.Zero
	detailSet, err := s.cutDetailDao.QueryByCutOrder(cutOrderID)
	if err != nil {
		return total, nil, err
	}
	userCollection, err := s.userDao.SelectByIDs(detailSet.UserIDs())
	userMap := userCollection.ToMap()
	apiDetails := make([]*response.CutDetail, detailSet.Size())
	for i, detail := range detailSet.Items {
		apiDetails[i] = detail.ResponseCutDetail(userMap[detail.UserID])
	}
	total = detailSet.TotalCutoffPrice()
	return total, apiDetails, nil
}

func (s *Cut) CreateCutOrder(ctx context.Context, req *request.CutOrder) (*response.CutOrder, error) {
	cutOrder, err := s.cutOrderDao.QueryByUserAndStock(req.UserID, req.StockID)
	if err != nil {
		return nil, err
	}
	if cutOrder != nil {
		return nil, errors.NewWithCodef("DuplicateCutOrder", "当前商品正在砍价")
	}
	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	cutOrderReq := &model.CutOrder{
		UserID:     req.UserID,
		CutNo:      uuid.String(),
		GoodsID:    req.GoodsID,
		StockID:    req.StockID,
		ExpireTime: time.Now().AddDate(0, 0, 7),
	}
	id, err := s.cutOrderDao.CreateCutOrder(cutOrderReq, nil)
	if err != nil {
		return nil, err
	}
	cutOrder, err = s.cutOrderDao.QueryByID(id)
	if err != nil {
		return nil, err
	}
	return cutOrder.ResponseCutOrder(), nil
}

// CreateCutDetail creates a detailfor a cut order
func (s *Cut) CreateCutDetail(req *request.CreateCutDetail) (decimal.Decimal, error) {
	cutOrder, err := s.cutOrderDao.QueryByCutNo(req.CutNo)
	randomCutoff := decimal.Zero
	if err != nil {
		return randomCutoff, err
	}
	if cutOrder == nil {
		return randomCutoff, errors.NewWithCodef("ConsumedCutOrder", "您的朋友已经购买了商品")
	}
	// check duplicate cut
	historyDetail, err := s.cutDetailDao.QueryByCutOrderAndUser(cutOrder.ID, req.UserID)
	if err != nil {
		return randomCutoff, err
	}
	if historyDetail != nil {
		return randomCutoff, errors.NewWithCodef("DuplicateCut", "您已经帮您的朋友砍过价了")
	}
	detailSet, err := s.cutDetailDao.QueryByCutOrder(cutOrder.ID)
	if err != nil {
		return randomCutoff, err
	}
	stock, err := s.stockDao.SelectByID(cutOrder.StockID)
	if err != nil {
		return randomCutoff, err
	}
	discount := detailSet.TotalCutoffPrice()
	cutCalc := model.NewCutCalculator(stock.CostUnitPrice, stock.SaleUnitPrice, discount)
	randomCutoff = cutCalc.RandomCut()
	cutDetail := &model.CutDetail{
		CutOrderID: cutOrder.ID,
		UserID:     req.UserID,
		CutPrice:   randomCutoff,
	}
	_, err = s.cutDetailDao.CreateCutDetail(cutDetail, nil)
	if err != nil {
		return randomCutoff, err
	}
	return randomCutoff, nil
}
