package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// Pagination parameters.
type Pagination struct {
	Page    int
	PerPage int
	Limit   int
	Offset  int
}

// GetPagination parses pagination parameters from query string,
// validates them, and calculates Limit and Offset.
// This centralizes pagination logic to enforce SRP in handlers.
func GetPagination(c *gin.Context, defaultPage, defaultPerPage int) Pagination {
	p := Pagination{
		Page:    defaultPage,
		PerPage: defaultPerPage,
	}

	if pageStr := c.Query("page"); pageStr != "" {
		fmt.Sscanf(pageStr, "%d", &p.Page)
	}
	if perPageStr := c.Query("per_page"); perPageStr != "" {
		fmt.Sscanf(perPageStr, "%d", &p.PerPage)
	}

	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage < 1 || p.PerPage > 100 {
		p.PerPage = defaultPerPage
	}

	p.Limit = p.PerPage
	p.Offset = (p.Page - 1) * p.PerPage

	return p
}

// NewMeta calculates pagination metadata.
func NewMeta(page, perPage int, total int64) *Meta {
	totalPages := total / int64(perPage)
	if total%int64(perPage) != 0 {
		totalPages++
	}

	return &Meta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}
}
