package domain

import (
	"math"
)

type PaginationRequest struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

func NewPaginationRequest(page, pageSize int) PaginationRequest {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return PaginationRequest{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func (p PaginationRequest) GetLimit() int {
	return p.PageSize
}

type PaginationResponse struct {
	CurrentPage  int  `json:"current_page"`
	PageSize     int  `json:"page_size"`
	TotalRecords int  `json:"total_records"`
	TotalPages   int  `json:"total_pages"`
	HasNext      bool `json:"has_next"`
	HasPrev      bool `json:"has_prev"`
}

func NewPaginationResponse(req PaginationRequest, totalRecords int) PaginationResponse {
	totalPages := int(math.Ceil(float64(totalRecords) / float64(req.PageSize)))

	return PaginationResponse{
		CurrentPage:  req.Page,
		PageSize:     req.PageSize,
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		HasNext:      req.Page < totalPages,
		HasPrev:      req.Page > 1,
	}
}

type PageResult struct {
	Data       interface{}        `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

func NewPageResult(data interface{}, req PaginationRequest, totalRecords int) PageResult {
	return PageResult{
		Data:       data,
		Pagination: NewPaginationResponse(req, totalRecords),
	}
}
