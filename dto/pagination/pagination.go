package pagination

// PaginationResonse 用于分页的数据响应传输对象
type PaginationResonse struct {
	PaginationRequest
	Count      int64       `json:"count"`
	TotalPages int64       `json:"totalPages"`
	Data       interface{} `json:"data"`
}

func (p *PaginationResonse) SetCount(count int64) {
	p.Count = count
	if count%int64(p.PageSize) == 0 {
		p.TotalPages = count / int64(p.PageSize)
	} else {
		p.TotalPages = count/int64(p.PageSize) + 1
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
