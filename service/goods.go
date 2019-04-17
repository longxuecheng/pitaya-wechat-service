package service

import (
	"pitaya-wechat-service/api"
	"pitaya-wechat-service/dao"
	"pitaya-wechat-service/dto"
	"pitaya-wechat-service/model"
	"strings"
)

var GoodsServiceSingleton *GoodsService

func GoodsServiceInstance() *GoodsService {
	if GoodsServiceSingleton != nil {
		return GoodsServiceSingleton
	}
	GoodsServiceSingleton = new(GoodsService)
	GoodsServiceSingleton.goodsDao = dao.GoodsDaoSingleton
	GoodsServiceSingleton.goodsAttributeDao = dao.GoodsAttributeDaoSingleton
	GoodsServiceSingleton.goodsSpecDao = dao.GoodsSpecificationDaoSingleton
	GoodsServiceSingleton.attributeService = AttributeServiceSingleton
	return GoodsServiceSingleton
}

// GoodsService 作为类目服务，实现了api.GoodsService接口
// 服务依赖 (1. attributeService)
type GoodsService struct {
	goodsDao          *dao.GoodsDao
	goodsAttributeDao *dao.GoodsAttributeDao
	goodsSpecDao      *dao.GoodsSpecificationDao
	attributeService  api.IAttributeService
}

func (s *GoodsService) GetGoodsByCategory(categoryID int64) ([]*dto.GoodsItemDTO, error) {
	goods, err := s.goodsDao.SelectByCategory(categoryID)
	if err != nil {
		return nil, err
	}
	return buildGoodsDTOs(goods), nil
}

func (s *GoodsService) Gallery(goodsID int64) ([]dto.GoodsGalleryDTO, error) {
	return nil, nil
}

func (s *GoodsService) Info(goodsID int64) (*dto.GoodsInfoDTO, error) {
	goods, err := s.goodsDao.SelectByID(goodsID)
	if err != nil {
		return nil, err
	}
	dto := installGoodsInfoDTO(goods)
	return dto, nil
}

func (s *GoodsService) Attributes(goodsID int64) ([]*dto.AttributeDTO, error) {
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

func (s *GoodsService) Specifications(goodsID int64) ([]*dto.GoodsSpecificationDTO, error) {
	goodsSpecs, err := s.goodsSpecDao.SelectByGoodsID(goodsID)
	if err != nil {
		return nil, err
	}
	return buildGoodsSpecificationDTOs(goodsSpecs), nil
}

func (s *GoodsService) SpecificationDesc(goodsID int64, specIDs []int64, sep string) (string, error) {
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

func (s *GoodsService) HotGoods() ([]*dto.GoodsItemDTO, error) {
	goodsList, err := s.goodsDao.SelectAllByStatus(model.GoodsStatusOnSale)
	if err != nil {
		return nil, err
	}
	return buildGoodsDTOs(goodsList), nil
}

func buildGoodsSpecificationDTOs(models []*model.GoodsSpecification) []*dto.GoodsSpecificationDTO {
	dtos := make([]*dto.GoodsSpecificationDTO, len(models))
	for i, spec := range models {
		dto := new(dto.GoodsSpecificationDTO)
		dto.ID = spec.ID
		dto.SpecID = spec.SpecificationID
		dto.Value = spec.Value
		dto.PicURL = spec.PicURL
		dto.GoodsID = spec.GoodsID
		dtos[i] = dto
	}
	return dtos

}

func installGoodsInfoDTO(model *model.Goods) *dto.GoodsInfoDTO {
	dto := new(dto.GoodsInfoDTO)
	dto.ID = model.ID
	dto.SupplierID = model.SupplierID
	dto.Name = model.Name
	dto.ListPicURL = model.ListPicURL.String
	if model.Description.Valid {
		dto.Description = model.Description.String
	}
	if model.BriefDescription.Valid {
		dto.BriefDesc = model.BriefDescription.String
	}
	dto.RetailPrice = model.RetailPrice
	return dto
}

func installGoodsDTO(model *model.Goods) *dto.GoodsItemDTO {
	dto := new(dto.GoodsItemDTO)
	dto.ID = model.ID
	dto.Name = model.Name
	dto.PicURL = model.ListPicURL.String
	dto.RetailPrice = model.RetailPrice
	return dto
}

func buildGoodsDTOs(models []*model.Goods) []*dto.GoodsItemDTO {
	if models == nil || len(models) == 0 {
		return nil
	}
	dtos := make([]*dto.GoodsItemDTO, len(models))
	for i, model := range models {
		dtos[i] = installGoodsDTO(model)
	}
	return dtos
}
