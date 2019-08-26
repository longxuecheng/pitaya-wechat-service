package pagination

// Page 用于分页的数据响应传输对象
type Page struct {
	PageSize    uint64      `json:"pageSize"`
	CurrentPage uint64      `json:"currentPage"`
	Count       int64       `json:"count"`
	TotalPages  int64       `json:"totalPages"`
	Data        interface{} `json:"data"`
}

func (p *Page) SetCount(count int64) {
	p.Count = count
	if count%int64(p.PageSize) == 0 {
		p.TotalPages = count / int64(p.PageSize)
	} else {
		p.TotalPages = count/int64(p.PageSize) + 1
	}
}

func (p Page) Offet() uint64 {
	if p.PageSize == 0 {
		p.PageSize = 10
	}
	offset := (p.CurrentPage - 1) * p.PageSize
	return offset
}

func (p Page) Limit() uint64 {
	return p.PageSize
}
