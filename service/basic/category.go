package basic

import (
	"gotrue/dao"
	"gotrue/dto"
	"gotrue/model"
)

var CategoryService *Category

// init 在此实现spring中类似注入的功能
func initCategoryService() {
	CategoryService = new(Category)
	CategoryService.categoryDao = dao.CategoryDao
}

// Category 作为类目服务，实现了api.Category接口
type Category struct {
	categoryDao *dao.Category
}

func (cs *Category) GetCategoryTree() ([]*dto.CategoryDTO, error) {
	categories, err := cs.categoryDao.SelectAll()
	if err != nil {
		return nil, err
	}
	return buildCategoryTreeDTOs(categories), nil
}

func (cs *Category) GetTopList() ([]*dto.CategoryDTO, error) {
	list, err := cs.categoryDao.SelectAllTopCategories()
	if err != nil {
		return nil, err
	}
	return buildCategoryDTOs(list), nil
}

func (cs *Category) GetCategoryByID(ID int) (*dto.CategoryDTO, error) {
	parent, err := cs.categoryDao.SelectByID(ID)
	if err != nil {
		return nil, err
	}
	if parent.ParentID > 0 {
		return installCategoryDTO(parent), nil
	}
	children, err := cs.categoryDao.SelectByParentID(parent.ID)
	if err != nil {
		return nil, err
	}
	parentDTO := installCategoryDTO(parent)
	childrenDTOs := buildCategoryDTOs(children)
	parentDTO.SubCategoryList = childrenDTOs
	parent = nil
	children = nil
	return parentDTO, nil
}

func (cs *Category) GetCategoriesByParentID(parentID int) ([]*dto.CategoryDTO, error) {
	children, err := cs.categoryDao.SelectByParentID(parentID)
	if err != nil {
		return nil, err
	}
	childrenDTOs := buildCategoryDTOs(children)
	return childrenDTOs, nil
}

func (cs *Category) Channels() ([]dto.CategoryChannel, error) {
	list, err := cs.categoryDao.SelectChannels()
	if err != nil {
		return nil, err
	}
	channels := buildChannels(list)
	return channels, nil
}

func installCategoryDTO(model *model.Category) *dto.CategoryDTO {
	dto := new(dto.CategoryDTO)
	dto.ID = model.ID
	dto.Name = model.Name
	dto.ParentID = model.ParentID
	dto.BannerURL = model.BannerURL
	dto.WapBannerURL = model.WapBannerURL
	dto.FrontName = model.FrontName
	return dto
}

func buildChannels(models []*model.Category) []dto.CategoryChannel {
	if models == nil || len(models) == 0 {
		return nil
	}
	dtos := make([]dto.CategoryChannel, len(models))
	for i, model := range models {
		dtos[i] = dto.CategoryChannel{
			ID:      model.ID,
			Name:    model.Name,
			IconURL: model.IconUrl,
		}
	}
	return dtos
}

func buildCategoryDTOs(models []*model.Category) []*dto.CategoryDTO {
	if models == nil || len(models) == 0 {
		return nil
	}
	dtos := make([]*dto.CategoryDTO, len(models))
	for i, model := range models {
		dtos[i] = installCategoryDTO(model)
	}
	return dtos
}

func buildCategoryTreeDTOs(models []*model.Category) []*dto.CategoryDTO {
	dtos := buildCategoryDTOs(models)
	if dtos == nil || len(dtos) == 0 {
		return nil
	}
	roots := []*dto.CategoryDTO{}
	for _, item := range dtos {
		if item.ParentID == 0 {
			root := item
			children := []*dto.CategoryDTO{}
			for _, item1 := range dtos {
				if item1.ParentID == item.ID {
					children = append(children, item1)
				}
			}
			root.SubCategoryList = children
			roots = append(roots, root)
		}
	}
	return roots
}
