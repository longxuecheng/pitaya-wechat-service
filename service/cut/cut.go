package cut

import (
	"context"
	"database/sql"
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
	goodsDao     *dao.Goods
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
		goodsDao:     dao.GoodsDao,
	}
}

// MyActivatedCutOrder will get an activated cut order for a given stock
func (s *Cut) MyActivatedCutOrder(req *request.CutOrder) (*response.CutOrder, error) {
	cutOrder, err := s.cutOrderDao.QueryByUserAndGoods(req.UserID, req.GoodsID)
	if err != nil {
		return nil, err
	}
	if cutOrder == nil {
		return nil, nil
	}
	detailSet, err := s.cutDetailDao.QueryByCutOrder(cutOrder.ID)
	if err != nil {
		return nil, err
	}
	cutoffPrice := detailSet.TotalCutoffPrice()
	apiCutOrder := cutOrder.ResponseCutOrder()
	apiCutOrder.SetCutoffPrice(cutoffPrice)
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

func (s *Cut) GetCutoffInfo(userID string, goodsID int64) (*response.CutOrder, error) {
	userInt64ID := utils.DecodeIn64(userID)
	cutOrder, err := s.cutOrderDao.QueryByUserAndGoods(userInt64ID, goodsID)
	if err != nil {
		return nil, err
	}
	user, err := s.userDao.SelectByID(userInt64ID)
	if err != nil {
		return nil, err
	}
	if cutOrder == nil {
		return &response.CutOrder{
			Avatar:            user.AvatarURL,
			CutoffPriceString: decimal.Zero.StringFixed(2),
		}, nil
	}
	total, apiDetails, err := s.CutDetails(cutOrder.ID)
	if err != nil {
		return nil, err
	}

	apiCutOrder := cutOrder.ResponseCutOrder()
	apiCutOrder.SetCutoffPrice(total)
	apiCutOrder.Items = apiDetails
	apiCutOrder.Avatar = user.AvatarURL
	apiCutOrder.NickName = user.NickName
	return apiCutOrder, nil
}

func (s *Cut) AssistCutoff(ctx context.Context, req *request.AssistCutoff) (*response.CutOrder, error) {
	userID := utils.DecodeIn64(req.UserID)
	if userID == req.HelperID {
		return nil, errors.NewWithCodef("InvalidCutoff", "自己不可以帮自己砍价哦")
	}
	stocks, err := s.stockDao.SelectByGoodsIDWithPriceASC(req.GoodsID)
	if err != nil {
		return nil, err
	}
	if len(stocks) == 0 {
		return nil, errors.NewWithCodef("NonCutoffGoods", "该商品不支持砍价")
	}
	cutOrder, err := s.cutOrderDao.QueryByUserAndGoods(userID, req.GoodsID)
	if err != nil {
		return nil, err
	}
	randomCutoff := decimal.Zero
	var minPriceStock = stocks[0]
	var discount = decimal.Zero

	if cutOrder == nil {
		dao.GetEasyDB().ExecTx(func(tx *sql.Tx) error {
			orderID, err := s.createCutoffOrder(userID, req.GoodsID, tx)
			if err != nil {
				return err
			}
			calc := model.NewCutCalculator(minPriceStock.CostUnitPrice, minPriceStock.SaleUnitPrice, discount)
			randomCutoff = calc.RandomCut()
			return s.createCutoffDetail(req.HelperID, orderID, randomCutoff, tx)
		})
		return s.GetCutoffInfo(req.UserID, req.GoodsID)
	}
	// check duplicate cut
	historyDetail, err := s.cutDetailDao.QueryByCutOrderAndUser(cutOrder.ID, req.HelperID)
	if err != nil {
		return nil, err
	}
	if historyDetail != nil {
		return nil, errors.NewWithCodef("DuplicateCutoff", "您已经帮您的朋友砍过价了")
	}
	detailSet, err := s.cutDetailDao.QueryByCutOrder(cutOrder.ID)
	if err != nil {
		return nil, err
	}
	discount = detailSet.TotalCutoffPrice()
	calc := model.NewCutCalculator(minPriceStock.CostUnitPrice, minPriceStock.SaleUnitPrice, discount)
	randomCutoff = calc.RandomCut()
	err = s.createCutoffDetail(req.HelperID, cutOrder.ID, randomCutoff, nil)
	if err != nil {
		return nil, err
	}
	return s.GetCutoffInfo(req.UserID, req.GoodsID)
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

func (s *Cut) createCutoffOrder(userID, goodsID int64, tx *sql.Tx) (int64, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return 0, err
	}
	req := &model.CutOrder{
		UserID:     userID,
		CutNo:      uuid.String(),
		GoodsID:    goodsID,
		StockID:    0,
		ExpireTime: time.Now().AddDate(0, 0, 7),
	}
	return s.cutOrderDao.CreateCutOrder(req, tx)
}

func (s *Cut) createCutoffDetail(userID, cutOrderID int64, cutoff decimal.Decimal, tx *sql.Tx) error {
	cutDetail := &model.CutDetail{
		CutOrderID: cutOrderID,
		UserID:     userID,
		CutPrice:   cutoff,
	}
	_, err := s.cutDetailDao.CreateCutDetail(cutDetail, nil)
	return err
}
