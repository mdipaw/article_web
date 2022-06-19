package rest

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaginatorV2[T any, ToMap any] struct {
	QueryPage int
	PerPage   int
}
type Response[T any, ToMap any] struct {
	Error      error `json:"-"`
	Page       int   `json:"page"`
	TotalPage  int   `json:"totalPage"`
	TotalCount int   `json:"totalCount"`
	Data       T     `json:"data"`
}
type MappedResponse[T any] struct {
	Error      error `json:"-"`
	Page       int   `json:"page"`
	TotalPage  int   `json:"totalPage"`
	TotalCount int   `json:"totalCount"`
	Data       T     `json:"data"`
}

func PaginatorFromContext[T any, ToMap any](c *gin.Context) PaginatorV2[T, ToMap] {
	var page, perPage int
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	if page < 1 {
		page = 1
	}

	perPage, err = strconv.Atoi(c.Query("perPage"))
	if err != nil {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	} else if perPage < 1 {
		perPage = 1
	}

	return PaginatorV2[T, ToMap]{
		QueryPage: page,
		PerPage:   perPage,
	}
}

func (p *PaginatorV2[T, X]) PaginateQuery(q *gorm.DB) *Response[T, X] {
	var offset int
	if p.QueryPage == 1 {
		offset = 0
	} else {
		offset = (p.QueryPage - 1) * p.PerPage
	}

	var results T
	if err := q.Session(&gorm.Session{}).Limit(p.PerPage).Offset(offset).Find(&results).Error; err != nil {
		return &Response[T, X]{Error: err}
	}

	var totalCount int64
	if countQuery := q.Count(&totalCount); countQuery.Error != nil {
		return &Response[T, X]{Error: countQuery.Error}
	}
	pageCount := int(math.Ceil(float64(totalCount) / float64(p.PerPage)))

	return &Response[T, X]{
		Page:       p.QueryPage,
		TotalPage:  pageCount,
		TotalCount: int(totalCount),
		Data:       results,
	}
}

func (r *Response[T, X]) Map(mapFunc func(x T) X) *MappedResponse[X] {
	if r.Error != nil {
		return &MappedResponse[X]{Error: r.Error}
	}
	mappedData := mapFunc(r.Data)
	return &MappedResponse[X]{
		Page:       r.Page,
		TotalPage:  r.TotalPage,
		TotalCount: r.TotalCount,
		Data:       mappedData,
	}
}
