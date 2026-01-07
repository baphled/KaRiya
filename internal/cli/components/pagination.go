package components

import "fmt"

// PaginationHelper manages pagination state and calculations for list models.
// It provides methods for calculating pages, navigating between pages,
// and getting page bounds for data slicing.
type PaginationHelper struct {
	currentPage int // Current page number (1-based)
	pageSize    int // Number of items per page
	totalCount  int // Total number of items
}

// NewPaginationHelper creates a new pagination helper
func NewPaginationHelper(pageSize int) *PaginationHelper {
	return &PaginationHelper{
		currentPage: 1,
		pageSize:    pageSize,
		totalCount:  0,
	}
}

// SetTotalCount sets the total number of items
func (ph *PaginationHelper) SetTotalCount(count int) *PaginationHelper {
	ph.totalCount = count
	// Ensure current page is still valid
	if ph.currentPage > ph.GetTotalPages() && ph.totalCount > 0 {
		ph.currentPage = ph.GetTotalPages()
	}
	if ph.totalCount == 0 {
		ph.currentPage = 1
	}
	return ph
}

// GetCurrentPage returns the current page number
func (ph *PaginationHelper) GetCurrentPage() int {
	return ph.currentPage
}

// GetPageSize returns the page size
func (ph *PaginationHelper) GetPageSize() int {
	return ph.pageSize
}

// GetTotalCount returns the total number of items
func (ph *PaginationHelper) GetTotalCount() int {
	return ph.totalCount
}

// GetTotalPages calculates the total number of pages
func (ph *PaginationHelper) GetTotalPages() int {
	if ph.totalCount == 0 {
		return 1
	}
	return (ph.totalCount + ph.pageSize - 1) / ph.pageSize
}

// GetPageStartIndex returns the starting index for the current page (0-based)
func (ph *PaginationHelper) GetPageStartIndex() int {
	return (ph.currentPage - 1) * ph.pageSize
}

// GetPageEndIndex returns the ending index for the current page (exclusive, 0-based)
func (ph *PaginationHelper) GetPageEndIndex() int {
	return ph.GetPageStartIndex() + ph.pageSize
}

// GetPaginationInfo returns formatted pagination info for display
// Example: "Showing 1-10 of 25 items"
func (ph *PaginationHelper) GetPaginationInfo(itemName string) string {
	if ph.totalCount == 0 {
		return fmt.Sprintf("Showing 0 of 0 %s", itemName)
	}
	startIdx := ph.GetPageStartIndex() + 1
	endIdx := ph.GetPageEndIndex()
	if endIdx > ph.totalCount {
		endIdx = ph.totalCount
	}
	return fmt.Sprintf("Showing %d-%d of %d %s", startIdx, endIdx, ph.totalCount, itemName)
}

// NextPage moves to the next page if available
func (ph *PaginationHelper) NextPage() *PaginationHelper {
	if ph.currentPage < ph.GetTotalPages() {
		ph.currentPage++
	}
	return ph
}

// PrevPage moves to the previous page if available
func (ph *PaginationHelper) PrevPage() *PaginationHelper {
	if ph.currentPage > 1 {
		ph.currentPage--
	}
	return ph
}

// GoToFirstPage moves to the first page
func (ph *PaginationHelper) GoToFirstPage() *PaginationHelper {
	ph.currentPage = 1
	return ph
}

// GoToLastPage moves to the last page
func (ph *PaginationHelper) GoToLastPage() *PaginationHelper {
	ph.currentPage = ph.GetTotalPages()
	return ph
}

// IsFirstPage returns true if on the first page
func (ph *PaginationHelper) IsFirstPage() bool {
	return ph.currentPage == 1
}

// IsLastPage returns true if on the last page
func (ph *PaginationHelper) IsLastPage() bool {
	return ph.currentPage >= ph.GetTotalPages()
}

// CanGoNext returns true if there's a next page
func (ph *PaginationHelper) CanGoNext() bool {
	return ph.currentPage < ph.GetTotalPages()
}

// CanGoPrev returns true if there's a previous page
func (ph *PaginationHelper) CanGoPrev() bool {
	return ph.currentPage > 1
}
