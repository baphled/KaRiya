package models

import (
	"strings"

	"github.com/baphled/kariya/internal/domain/career"
)

// SortModel manages event sorting options
type SortModel struct {
	sortBy    string
	sortOrder string
}

// NewSortModel creates a new sort model with default sort (date, descending)
func NewSortModel() *SortModel {
	return &SortModel{
		sortBy:    "date",
		sortOrder: "desc",
	}
}

// GetSortBy returns the current sort field
func (m *SortModel) GetSortBy() string {
	return m.sortBy
}

// SetSortBy sets the sort field (date, created_at, or text)
func (m *SortModel) SetSortBy(field string) {
	field = strings.ToLower(field)
	validFields := []string{"date", "created_at", "text"}

	for _, valid := range validFields {
		if field == valid {
			m.sortBy = field
			return
		}
	}
	// Keep current if invalid
}

// GetSortOrder returns the current sort order (asc or desc)
func (m *SortModel) GetSortOrder() string {
	return m.sortOrder
}

// SetSortOrder sets the sort order (asc or desc)
func (m *SortModel) SetSortOrder(order string) {
	order = strings.ToLower(order)
	if order == "asc" || order == "desc" {
		m.sortOrder = order
	}
	// Keep current if invalid
}

// ToggleSortOrder toggles between ascending and descending
func (m *SortModel) ToggleSortOrder() {
	if m.sortOrder == "asc" {
		m.sortOrder = "desc"
	} else {
		m.sortOrder = "asc"
	}
}

// GetSortFields returns available sort fields
func (m *SortModel) GetSortFields() []string {
	return []string{"date", "created_at", "text"}
}

// GetSortOrders returns available sort orders
func (m *SortModel) GetSortOrders() []string {
	return []string{"asc", "desc"}
}

// Reset resets to default sort (date, descending)
func (m *SortModel) Reset() {
	m.sortBy = "date"
	m.sortOrder = "desc"
}

// ShouldReorder determines if two events need reordering based on current sort
func (m *SortModel) ShouldReorder(event1, event2 *career.CareerEvent) bool {
	if event1 == nil || event2 == nil {
		return false
	}

	var compare int

	switch m.sortBy {
	case "date":
		if event1.Date.Before(event2.Date) {
			compare = -1
		} else if event1.Date.After(event2.Date) {
			compare = 1
		} else {
			compare = 0
		}

	case "created_at":
		if event1.CreatedAt.Before(event2.CreatedAt) {
			compare = -1
		} else if event1.CreatedAt.After(event2.CreatedAt) {
			compare = 1
		} else {
			compare = 0
		}

	case "text":
		compare = strings.Compare(event1.Text, event2.Text)

	default:
		// Default to date descending
		if event1.Date.Before(event2.Date) {
			compare = -1
		} else if event1.Date.After(event2.Date) {
			compare = 1
		} else {
			compare = 0
		}
	}

	// For descending order, return true if comparison indicates we need to swap
	if m.sortOrder == "desc" {
		return compare < 0
	}
	// For ascending order, return true if comparison indicates we need to swap
	return compare > 0
}
