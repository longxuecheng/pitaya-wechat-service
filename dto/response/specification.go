package response

type SpecificationDTO struct {
	ID        int64                    `json:"specification_id"`
	Name      string                   `json:"name"`
	ValueList []*GoodsSpecificationDTO `json:"valueList"`
}

type GoodsSpecificationDTO struct {
	ID      int64  `json:"id"`
	GoodsID int64  `json:"goodsId"`
	SpecID  int64  `json:"specification_id"`
	Value   string `json:"value"`
	PicURL  string `json:"picUrl"`
}

// GoodsSpecificationSet 是商品定义的集合，用于进行简单的数据操作
type GoodsSpecificationSet struct {
	itemMap          map[int64]*GoodsSpecificationDTO
	specificationIDs map[int64]bool
	specNameMap      map[int64]string
	specs            []*SpecificationDTO
}

func NewGoodsSpecSet(items []*GoodsSpecificationDTO) *GoodsSpecificationSet {
	set := new(GoodsSpecificationSet)
	itemMap := map[int64]*GoodsSpecificationDTO{}
	distinctSpecIDMap := map[int64]bool{}
	for _, item := range items {
		itemMap[item.ID] = item
		distinctSpecIDMap[item.SpecID] = true
	}
	set.itemMap = itemMap
	set.specificationIDs = distinctSpecIDMap
	return set
}

func (s *GoodsSpecificationSet) SetSpecs(specs []*SpecificationDTO) {
	s.specs = specs
}

func (s *GoodsSpecificationSet) DistinctSpecIDs() []int64 {
	specIDs := make([]int64, 0, len(s.specificationIDs))
	for key := range s.specificationIDs {
		specIDs = append(specIDs, key)
	}
	return specIDs
}

func (s *GoodsSpecificationSet) SpecTree() []*SpecificationDTO {
	for _, spec := range s.specs {
		values := []*GoodsSpecificationDTO{}
		for _, item := range s.itemMap {
			if item.SpecID == spec.ID {
				values = append(values, item)
			}
		}

		spec.ValueList = values
	}
	return s.specs
}
