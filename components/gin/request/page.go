package request

import (
	"gorm.io/gorm"
	"math"
)

type Pager struct {
	Page     int `json:"page"     query:"page" binding:"required" `
	PageSize int `json:"pageSize" query:"pageSize" binding:"required" `

	TotalCount int64
}

// PagerSummary is summary of the query.
type PagerSummary struct {
	Page       int   `json:"page"        query:"page"`
	PerPage    int   `json:"pageSize"    query:"pageSize"`
	TotalCount int64 `json:"totalCount" query:"totalCount"`
	TotalPage  int   `json:"totalPage"  query:"totalPage"`
}

// Summary returns a PagerSummary.
func (pager *Pager) Summary() *PagerSummary {
	return &PagerSummary{
		Page:       pager.Page,
		PerPage:    pager.PageSize,
		TotalCount: pager.TotalCount,
		TotalPage:  int(math.Ceil(float64(pager.TotalCount) / float64(pager.PageSize))),
	}
}

// Scope returns a GORM scope.
func (pager *Pager) Scope() func(db *gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {
		// 纠正异常参数
		if pager.Page <= 0 {
			pager.Page = 1
		}

		switch {
		case pager.PageSize > 100:
			pager.PageSize = 100
		case pager.PageSize <= 0:
			pager.PageSize = 10
		}

		offset := (pager.Page - 1) * pager.PageSize
		return db.Offset(offset).Limit(pager.PageSize)
	}
}
