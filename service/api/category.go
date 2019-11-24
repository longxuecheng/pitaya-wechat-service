package api

import "gotrue/dto"

// ICategoryService is a category service's interface
type ICategoryService interface {
	GetTopList() ([]*dto.CategoryDTO, error)

	GetInternalTopList() ([]*InternalCategory, error)

	GetCategoryTree() ([]*dto.CategoryDTO, error)

	GetCategoryByID(ID int) (*dto.CategoryDTO, error)

	// GetCategoriesByParentID 通过一个类别的ID获取与该类别同级的类别列表
	GetCategoriesByParentID(parentID int) ([]*dto.CategoryDTO, error)

	Channels() ([]dto.CategoryChannel, error)
}

type InternalCategory struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
