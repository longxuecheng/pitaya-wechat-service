package dto

// PaginationDTO 用于分页的数据传输对象
type PaginationDTO struct {
	PageSize    int64       `json:"pageSize"`
	CurrentPage int64       `json:"currentPage"`
	TotalCount  int64       `json:"totalCount"`
	TotalPages  int64       `json:"totalPages"`
	Data        interface{} `json:"data"`
}
