package behaviors

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// FilterBehavior defines the contract for intents with search/filter/sort functionality.
// This interface establishes ManageSkills as the source of truth for filter patterns.
//
// All intents with data tables (lists) should implement this interface to ensure
// consistent user experience across search, filter, and sort operations.
type FilterBehavior interface {
	// HasActiveFilters returns true if any non-default filters are active.
	// This determines whether the "x: Clear filters" badge should be shown.
	HasActiveFilters() bool

	// ClearFilters resets filters in FIFO order (most recent filter first).
	// Implementation should follow this order:
	//   1. Clear search text if present
	//   2. Clear active filters (companies, categories, etc.) if present
	//   3. Reset sort to default if modified
	ClearFilters()

	// ApplyFilters applies current filter state to the data.
	// This should be called after any filter/sort/search change.
	ApplyFilters()

	// RefreshData reloads/refreshes the filtered data.
	// Returns a tea.Cmd that triggers data reload.
	RefreshData() tea.Cmd
}

// SearchableText returns true if text contains the search query (case-insensitive).
// Returns true if query is empty (no filter applied).
//
// Example:
//
//	SearchableText("Backend Developer", "backend") // true
//	SearchableText("Backend Developer", "frontend") // false
//	SearchableText("Backend Developer", "") // true (no filter)
func SearchableText(text, query string) bool {
	if query == "" {
		return true
	}
	return strings.Contains(strings.ToLower(text), strings.ToLower(query))
}

// SearchableFields returns true if ANY field matches the search query (case-insensitive).
// Returns true if query is empty (no filter applied).
//
// This is the recommended pattern for search implementation across all intents.
//
// Example:
//
//	SearchableFields("backend", "Backend Developer", "TechCorp", "Engineering")
//	// Returns true because "Backend Developer" contains "backend"
//
//	SearchableFields("frontend", "Backend Developer", "TechCorp", "Engineering")
//	// Returns false because none of the fields contain "frontend"
func SearchableFields(query string, fields ...string) bool {
	if query == "" {
		return true
	}
	searchLower := strings.ToLower(query)
	for _, field := range fields {
		if strings.Contains(strings.ToLower(field), searchLower) {
			return true
		}
	}
	return false
}

// FilterStack represents a stack of active filters for FIFO clearing.
// This allows users to progressively remove filters by pressing 'x' multiple times.
type FilterStack struct {
	// Layers represents the filter layers in FIFO order.
	// Index 0 is the most recently applied filter (cleared first).
	Layers []FilterLayer
}

// FilterLayer represents a single filter type that can be cleared.
type FilterLayer string

const (
	// FilterLayerSearch represents search text filtering
	FilterLayerSearch FilterLayer = "search"

	// FilterLayerCategory represents category-based filtering
	FilterLayerCategory FilterLayer = "category"

	// FilterLayerCompany represents company-based filtering
	FilterLayerCompany FilterLayer = "company"

	// FilterLayerProject represents project-based filtering
	FilterLayerProject FilterLayer = "project"

	// FilterLayerTags represents tag-based filtering
	FilterLayerTags FilterLayer = "tags"

	// FilterLayerSort represents custom sort (non-default)
	FilterLayerSort FilterLayer = "sort"
)

// NewFilterStack creates an empty filter stack.
func NewFilterStack() *FilterStack {
	return &FilterStack{
		Layers: []FilterLayer{},
	}
}

// Push adds a filter layer to the stack (most recent).
func (s *FilterStack) Push(layer FilterLayer) {
	// Add to front (FIFO - first in, first out when clearing)
	s.Layers = append([]FilterLayer{layer}, s.Layers...)
}

// Pop removes and returns the most recent filter layer.
// Returns empty string if stack is empty.
func (s *FilterStack) Pop() FilterLayer {
	if len(s.Layers) == 0 {
		return ""
	}
	layer := s.Layers[0]
	s.Layers = s.Layers[1:]
	return layer
}

// IsEmpty returns true if no filters are active.
func (s *FilterStack) IsEmpty() bool {
	return len(s.Layers) == 0
}

// Clear removes all filter layers.
func (s *FilterStack) Clear() {
	s.Layers = []FilterLayer{}
}
