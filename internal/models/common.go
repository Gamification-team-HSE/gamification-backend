package models

type Pagination struct {
	Page int
	Size int
}

func (p *Pagination) ToRepo() *RepoPagination {
	return &RepoPagination{
		Offset: (p.Page - 1) * p.Size,
		Limit:  p.Size,
	}
}

type RepoPagination struct {
	Offset int
	Limit  int
}
