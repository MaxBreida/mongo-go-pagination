package mongopagination

import (
	"context"
	"math"
)

// Paginator struct for holding pagination info
type Paginator struct {
	TotalRecord int `json:"total_record"`
	TotalPage   int `json:"total_page"`
	Offset      int `json:"offset"`
	Limit       int `json:"limit"`
	Page        int `json:"page"`
	PrevPage    int `json:"prev_page"`
	NextPage    int `json:"next_page"`
}

// PaginationData struct for returning pagination stat
type PaginationData struct {
	Total         int `json:"total"`
	Page          int `json:"page"`
	PerPage       int `json:"perPage"`
	Prev          int `json:"prev"`
	Next          int `json:"next"`
	TotalPages    int `json:"totalPages"`
	RecordsOnPage int `json:"recordsOnPage"`
}

// PaginationData returns PaginationData struct which
// holds information of all stats needed for pagination
func (p *Paginator) PaginationData() *PaginationData {
	data := PaginationData{
		Total:      p.TotalRecord,
		Page:       p.Page,
		PerPage:    p.Limit,
		Prev:       0,
		Next:       0,
		TotalPages: p.TotalPage,
	}
	if p.Page != p.PrevPage && p.TotalRecord > 0 {
		data.Prev = p.PrevPage
	}
	if p.Page != p.NextPage && p.TotalRecord > 0 && p.Page <= p.TotalPage {
		data.Next = p.NextPage
	}

	return &data
}

// Paging returns Paginator struct which hold pagination
// stats
func Paging(p *PagingQuery) *Paginator {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 10
	}
	var paginator Paginator
	var offset int
	total, _ := p.Collection.CountDocuments(context.Background(), p.Filter)

	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.Limit
	}
	paginator.TotalRecord = int(total)
	paginator.Page = p.Page
	paginator.Offset = offset
	paginator.Limit = p.Limit
	paginator.TotalPage = int(math.Ceil(float64(total) / float64(p.Limit)))
	if p.Page > 1 {
		paginator.PrevPage = p.Page - 1
	} else {
		paginator.PrevPage = p.Page
	}
	if p.Page == paginator.TotalPage {
		paginator.NextPage = p.Page
	} else {
		paginator.NextPage = p.Page + 1
	}
	return &paginator
}
