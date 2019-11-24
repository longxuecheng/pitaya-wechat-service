package goods

import (
	"database/sql"
	"encoding/csv"
	"gotrue/dao"
	"gotrue/dto/response"
	"gotrue/model"
	"gotrue/service/api"
	"gotrue/service/basic"
	"io"
	"log"
	"strings"
)

var GoodsService *Goods

func initGoodsService() {
	if GoodsService != nil {
		return
	}

	GoodsService = &Goods{
		goodsDao:             dao.GoodsDao,
		stockDao:             dao.StockDao,
		userAddressDao:       dao.UserAddressDao,
		expressConstraintDao: dao.GoodsExpressConstraintDao,
		goodsAttributeDao:    dao.GoodsAttributeDao,
		goodsSpecDao:         dao.GoodsSpecDao,
		attributeService:     basic.AttributeService,
	}
}

type Goods struct {
	goodsDao             *dao.Goods
	stockDao             *dao.Stock
	userAddressDao       *dao.UserAddress
	expressConstraintDao *dao.GoodsExpressConstraint
	goodsAttributeDao    *dao.GoodsAttribute
	goodsSpecDao         *dao.GoodsSpec
	attributeService     api.IAttributeService
}

func (s *Goods) ImportExpressConstraints(stockID int64, r io.Reader) error {
	stock, err := s.stockDao.SelectByID(stockID)
	if err != nil {
		return err
	}
	// Parse the file
	csvReader := csv.NewReader(r)
	// Iterate through the records
	constrains := []*model.GoodsExpressConstraint{}
	for {
		// Read each record from csv
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		constraint := &model.GoodsExpressConstraint{
			StockID: stockID,
			GoodsID: stock.GoodsID,
		}
		err = constraint.SetFeeFromCsv(record)
		if err != nil {
			return err
		}
		constrains = append(constrains, constraint)
	}
	dao.GetEasyDB().ExecTx(func(tx *sql.Tx) error {
		for _, c := range constrains {
			_, err = s.expressConstraintDao.CreateConstraint(c, tx)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// ExpressConstraint gets express constraint for a given goods' stock
// eg: express fee definition and destination reachable
func (s *Goods) ExpressConstraint(stockID int64, addressID int64) (*response.GoodsExpressConstraint, error) {
	address, err := s.userAddressDao.SelectByID(addressID)
	if err != nil {
		return nil, err
	}
	result := response.BusinessResult{
		Code: "OK",
	}
	if address == nil {
		result.Message = "配送地址无效"
		result.Code = "InvalidAddress"
		return &response.GoodsExpressConstraint{
			BusinessResult: result,
		}, nil
	}
	constraint, err := s.expressConstraintDao.QueryByStockAndProvince(stockID, address.ProvinceID)
	if err != nil {
		if err == sql.ErrNoRows {
			result.Message = "不支持的配送范围"
			result.Code = "NoneSupportedArea"
		} else {
			return nil, err
		}
	}
	if !constraint.Reachable {
		result.Message = "不支持的配送范围"
		result.Code = "NoneSupportedArea"
	}
	return &response.GoodsExpressConstraint{
		IsFree:         constraint.IsFree,
		UnitExpressFee: constraint.ExpressFee,
		BusinessResult: result,
	}, nil
}

func (s *Goods) GetGoodsByCategory(categoryID int64) ([]*response.GoodsItem, error) {
	goods, err := s.goodsDao.SelectByCategory(categoryID)
	if err != nil {
		return nil, err
	}
	return buildApiGoods(goods), nil
}

func (s *Goods) GetInternalGoodsByCategory(categoryID int64) ([]*api.InternalGoods, error) {
	goods, err := s.goodsDao.SelectByCategory(categoryID)
	if err != nil {
		return nil, err
	}
	apiGoods := make([]*api.InternalGoods, len(goods))
	for i, g := range goods {
		apiGoods[i] = &api.InternalGoods{
			ID:   g.ID,
			Name: g.Name,
		}
	}
	return apiGoods, nil
}

func (s *Goods) Gallery(goodsID int64) ([]response.GoodsGalleryDTO, error) {
	return nil, nil
}

func (s *Goods) Info(goodsID int64) (*response.GoodsInfo, error) {
	goods, err := s.goodsDao.SelectByID(goodsID)
	if err != nil {
		return nil, err
	}
	return installApiGoodsInfo(goods), nil
}

func (s *Goods) Attributes(goodsID int64) ([]*response.AttributeDTO, error) {
	// 获取商品属性
	goodsAttributes, err := s.goodsAttributeDao.SelectByGoodsID(goodsID)
	if err != nil {
		return nil, err
	}
	attributeIDs := make([]int64, len(goodsAttributes))
	for i, goodsAttribute := range goodsAttributes {
		attributeIDs[i] = goodsAttribute.AttributeID
	}
	// 获取属性定义单元
	attributes, err := s.attributeService.GetByIDs(attributeIDs)
	if err != nil {
		return nil, err
	}
	// 设置商品属性对应
	for _, attribute := range attributes {
		for _, goodsAttribute := range goodsAttributes {
			if attribute.ID == goodsAttribute.AttributeID {
				attribute.Value = goodsAttribute.Value
			}
		}
	}
	return attributes, nil
}

func (s *Goods) Specifications(goodsID int64) ([]*response.GoodsSpecificationDTO, error) {
	goodsSpecs, err := s.goodsSpecDao.SelectByGoodsID(goodsID)
	if err != nil {
		return nil, err
	}
	return buildGoodsSpecificationDTOs(goodsSpecs), nil
}

func (s *Goods) SpecificationDesc(goodsID int64, specIDs []int64, sep string) (string, error) {
	if len(specIDs) == 0 {
		return "", nil
	}
	specs, err := s.Specifications(goodsID)
	if err != nil {
		return "", err
	}
	specNames := []string{}
	for _, specID := range specIDs {
		for _, spec := range specs {
			if specID == spec.ID {
				specNames = append(specNames, spec.Value)
				break
			}
		}
	}
	return strings.Join(specNames, sep), nil // 商品规格组合描述
}

func (s *Goods) HotGoods(categoryID int64) ([]response.HotGoods, error) {
	var goodsList []*model.Goods
	var err error
	if categoryID == 0 {
		goodsList, err = s.goodsDao.SelectAllByStatus(model.GoodsStatusOnSale)
		if err != nil {
			return nil, err
		}
	} else {
		goodsList, err = s.goodsDao.SelectByCategory(categoryID)
		if err != nil {
			return nil, err
		}
	}
	goodsSet := model.NewGoodsSet(goodsList)
	stockSet, err := s.stockDao.SelectByGoodsIDs(goodsSet.GoodsIDs())
	if err != nil {
		return nil, err
	}
	goodsSpecSet, err := s.goodsSpecDao.SelectByGoodsIDs(goodsSet.GoodsIDs())
	if err != nil {
		return nil, err
	}
	return s.buildApiHotGoods(goodsList, stockSet, goodsSpecSet), nil
}

func (s *Goods) OneSaleGoodsCards() ([]*response.GoodsCard, error) {
	goodsList, err := s.goodsDao.SelectAllByStatus(model.GoodsStatusOnSale)
	if err != nil {
		return nil, err
	}
	cards := []*response.GoodsCard{}
	for _, goods := range goodsList {
		if goods.ProdAreaLng == 0 && goods.ProdAreaLat == 0 {
			continue
		}
		cards = append(cards, &response.GoodsCard{
			ID:            goods.ID,
			Name:          goods.Name,
			CardPicURL:    goods.CardPicURL,
			RetailPrice:   goods.RetailPrice,
			Description:   goods.BriefDescription.String,
			Status:        goods.StatusName(),
			ProducingArea: goods.ProducingArea,
			Location: response.Location{
				Longitude: goods.ProdAreaLng,
				Latitude:  goods.ProdAreaLat,
			},
		})
	}
	return cards, nil
}

func buildGoodsSpecificationDTOs(models []*model.GoodsSpecification) []*response.GoodsSpecificationDTO {
	dtos := make([]*response.GoodsSpecificationDTO, len(models))
	for i, spec := range models {
		data := new(response.GoodsSpecificationDTO)
		data.ID = spec.ID
		data.SpecID = spec.SpecificationID
		data.Value = spec.Value
		data.PicURL = spec.PicURL
		data.GoodsID = spec.GoodsID
		dtos[i] = data
	}
	return dtos

}

func installApiGoodsInfo(model *model.Goods) *response.GoodsInfo {
	data := new(response.GoodsInfo)
	data.ID = model.ID
	data.SupplierID = model.SupplierID
	data.Name = model.Name
	data.ListPicURL = model.ListPicURL.String
	if model.Description.Valid {
		data.Description = model.Description.String
	}
	if model.BriefDescription.Valid {
		data.BriefDesc = model.BriefDescription.String
	}
	data.RetailPrice = model.RetailPrice
	data.ProducingArea = model.ProducingArea
	return data
}

func installApiGoods(model *model.Goods) *response.GoodsItem {
	data := new(response.GoodsItem)
	data.ID = model.ID
	data.Name = model.Name
	data.PicURL = model.ListPicURL.String
	data.RetailPrice = model.RetailPrice
	data.ProducingArea = model.ProducingArea
	return data
}

func (s *Goods) buildApiHotGoods(models []*model.Goods, stockSet *model.StockSet, goodsSpecSet *model.GoodsSpecSet) []response.HotGoods {
	if models == nil || len(models) == 0 {
		return nil
	}
	stockSpecMap := stockSet.SpecMap()
	goodsSpecMap := goodsSpecSet.Map()
	dtos := make([]response.HotGoods, len(models))
	// goods => stocks => specifications
	for i, model := range models {
		data := response.HotGoods{
			ID:            model.ID,
			Name:          model.Name,
			PicURL:        model.ListPicURL.String,
			RetailPrice:   model.RetailPrice,
			ProducingArea: model.ProducingArea,
		}
		stocks := stockSet.GetByGoods(model.ID)
		stockPriceList := make([]response.StockPrice, len(stocks))
		for i, stock := range stocks {
			specIDList := stockSpecMap.GetSpecs(stock.ID)
			specNames := []string{}
			for _, specID := range specIDList {
				spec := goodsSpecMap.Get(specID)
				if spec != nil {
					specNames = append(specNames, spec.Value)
				}
			}
			stockPriceList[i] = response.StockPrice{
				Spec:      strings.Join(specNames, "|"),
				SalePrice: stock.SaleUnitPrice.StringFixed(2),
			}
		}
		data.StockPrices = stockPriceList
		dtos[i] = data
	}
	return dtos
}

func buildApiGoods(models []*model.Goods) []*response.GoodsItem {
	if models == nil || len(models) == 0 {
		return nil
	}
	dtos := make([]*response.GoodsItem, len(models))
	for i, model := range models {
		dtos[i] = installApiGoods(model)
	}
	return dtos
}
