package pagination

// PaginationResonse 用于分页的数据响应传输对象
type PaginationResonse struct {
	PaginationRequest
	Count      int64       `json:"count"`
	TotalPages int64       `json:"totalPages"`
	Data       interface{} `json:"data"`
}

func (p PaginationResonse) SetCount(count int64) {
	p.Count = count
	pageNumber := count % int64(p.PageSize)
	if pageNumber == 0 {
		p.TotalPages = pageNumber
	} else {
		p.TotalPages = pageNumber + 1
	}
}

// PaginationRequest 用于分页的数据请求传输对象
type PaginationRequest struct {
	PageSize    uint64 `json:"pageSize"`
	CurrentPage uint64 `json:"currentPage"`
}

func (pr PaginationRequest) Offet() uint64 {
	if pr.PageSize == 0 {
		pr.PageSize = 10
	}
	offset := (pr.CurrentPage - 1) * pr.PageSize
	return offset
}

func (pr PaginationRequest) Limit() uint64 {
	return pr.PageSize
}
