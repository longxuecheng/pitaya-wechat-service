package goods

import (
	"database/sql"
	"gotrue/dao"
	"gotrue/dto/response"
	"gotrue/model"
	"gotrue/service/api"
	"gotrue/service/basic"
	"strings"
)

var GoodsService *Goods

func initGoodsService() {
	if GoodsService != nil {
		return
	}

	GoodsService = &Goods{
		goodsDao:             dao.GoodsDao,
		userAddressDao:       dao.UserAddressDao,
		expressConstraintDao: dao.GoodsExpressConstraintDao,
		goodsAttributeDao:    dao.GoodsAttributeDao,
		goodsSpecDao:         dao.GoodsSpecDao,
		attributeService:     basic.AttributeService,
	}
}

type Goods struct {
	goodsDao             *dao.Goods
	userAddressDao       *dao.UserAddress
	expressConstraintDao *dao.GoodsExpressConstraint
	goodsAttributeDao    *dao.GoodsAttribute
	goodsSpecDao         *dao.GoodsSpec
	attributeService     api.IAttributeService
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
	return &response.GoodsExpressConstraint{
		IsFree:         constraint.IsFree,
		UnitExpressFee: constraint.ExpressFee,
		BusinessResult: result,
	}, nil
}

func (s *Goods) GetGoodsByCategory(categoryID int64) ([]*response.GoodsItemDTO, error) {
	goods, err := s.goodsDao.SelectByCategory(categoryID)
	if err != nil {
		return nil, err
	}
	return buildGoodsDTOs(goods), nil
}

func (s *Goods) Gallery(goodsID int64) ([]response.GoodsGalleryDTO, error) {
	return nil, nil
}

func (s *Goods) Info(goodsID int64) (*response.GoodsInfoDTO, error) {
	goods, err := s.goodsDao.SelectByID(goodsID)
	if err != nil {
		return nil, err
	}
	dto := installGoodsInfoDTO(goods)
	return dto, nil
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

func (s *Goods) HotGoods() ([]*response.GoodsItemDTO, error) {
	goodsList, err := s.goodsDao.SelectAllByStatus(model.GoodsStatusOnSale)
	if err != nil {
		return nil, err
	}
	return buildGoodsDTOs(goodsList), nil
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

func installGoodsInfoDTO(model *model.Goods) *response.GoodsInfoDTO {
	data := new(response.GoodsInfoDTO)
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
	return data
}

func installGoodsDTO(model *model.Goods) *response.GoodsItemDTO {
	data := new(response.GoodsItemDTO)
	data.ID = model.ID
	data.Name = model.Name
	data.PicURL = model.ListPicURL.String
	data.RetailPrice = model.RetailPrice
	return data
}

func buildGoodsDTOs(models []*model.Goods) []*response.GoodsItemDTO {
	if models == nil || len(models) == 0 {
		return nil
	}
	dtos := make([]*response.GoodsItemDTO, len(models))
	for i, model := range models {
		dtos[i] = installGoodsDTO(model)
	}
	return dtos
}
