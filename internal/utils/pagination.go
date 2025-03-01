package utils

import (
	"math"
	"strconv"
)

// Pagination represents pagination metadata
type Pagination struct {
	CurrentPage int
	PerPage     int
	TotalItems  int
	TotalPages  int
	HasPrev     bool
	HasNext     bool
	PrevPage    int
	NextPage    int
	Pages       []int
}

// NewPagination creates a new pagination instance
func NewPagination(page, perPage, totalItems int) *Pagination {
	// Ensure we have valid values
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))
	if totalPages < 1 {
		totalPages = 1
	}

	// Ensure page is within bounds
	if page > totalPages {
		page = totalPages
	}

	// Calculate prev/next page
	hasPrev := page > 1
	hasNext := page < totalPages
	prevPage := page - 1
	nextPage := page + 1

	// Generate page numbers to display
	var pages []int
	if totalPages <= 7 {
		// If we have 7 or fewer pages, show all
		for i := 1; i <= totalPages; i++ {
			pages = append(pages, i)
		}
	} else {
		// Show pages with ellipsis
		if page <= 3 {
			// Near the beginning
			for i := 1; i <= 5; i++ {
				pages = append(pages, i)
			}
			pages = append(pages, totalPages)
		} else if page >= totalPages-2 {
			// Near the end
			pages = append(pages, 1)
			for i := totalPages - 4; i <= totalPages; i++ {
				pages = append(pages, i)
			}
		} else {
			// Middle - show current page with 2 pages on each side
			pages = append(pages, 1)
			for i := page - 2; i <= page+2; i++ {
				pages = append(pages, i)
			}
			pages = append(pages, totalPages)
		}
	}

	return &Pagination{
		CurrentPage: page,
		PerPage:     perPage,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		HasPrev:     hasPrev,
		HasNext:     hasNext,
		PrevPage:    prevPage,
		NextPage:    nextPage,
		Pages:       pages,
	}
}

// GetPageFromQuery extracts page number from query parameters
func GetPageFromQuery(query map[string][]string, defaultPage int) int {
	if pageStr, ok := query["page"]; ok && len(pageStr) > 0 {
		if page, err := strconv.Atoi(pageStr[0]); err == nil && page > 0 {
			return page
		}
	}
	return defaultPage
}

// GetPerPageFromQuery extracts per page from query parameters
func GetPerPageFromQuery(query map[string][]string, defaultPerPage int) int {
	if perPageStr, ok := query["per_page"]; ok && len(perPageStr) > 0 {
		if perPage, err := strconv.Atoi(perPageStr[0]); err == nil && perPage > 0 {
			if perPage > 100 {
				return 100 // Max limit
			}
			return perPage
		}
	}
	return defaultPerPage
}