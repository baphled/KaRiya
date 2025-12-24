package models

import (
	"time"

	careerrepo "github.com/baphled/kariya/internal/repository/career"
)

// FilterModel manages event filtering with date range and tags
type FilterModel struct {
	tags       []string
	startDate  *time.Time
	endDate    *time.Time
	focusIndex int
	err        error
}

// NewFilterModel creates a new filter model with empty filters
func NewFilterModel() *FilterModel {
	return &FilterModel{
		tags:       []string{},
		startDate:  nil,
		endDate:    nil,
		focusIndex: 0,
	}
}

// ToggleTag adds or removes a tag from the filter
func (m *FilterModel) ToggleTag(tag string) {
	for i, t := range m.tags {
		if t == tag {
			// Remove tag
			m.tags = append(m.tags[:i], m.tags[i+1:]...)
			return
		}
	}
	// Add tag
	m.tags = append(m.tags, tag)
}

// GetTags returns the currently selected tags
func (m *FilterModel) GetTags() []string {
	return m.tags
}

// SetStartDate sets the start date for the date range filter
func (m *FilterModel) SetStartDate(date *time.Time) {
	m.startDate = date
}

// GetStartDate returns the start date filter
func (m *FilterModel) GetStartDate() *time.Time {
	return m.startDate
}

// SetEndDate sets the end date for the date range filter
func (m *FilterModel) SetEndDate(date *time.Time) {
	m.endDate = date
}

// GetEndDate returns the end date filter
func (m *FilterModel) GetEndDate() *time.Time {
	return m.endDate
}

// ToListFilters converts the filter model to a repository ListFilters struct
func (m *FilterModel) ToListFilters() careerrepo.ListFilters {
	filters := careerrepo.ListFilters{
		Tags:      m.tags,
		StartDate: m.startDate,
		EndDate:   m.endDate,
	}
	return filters
}

// Reset clears all filters
func (m *FilterModel) Reset() {
	m.tags = []string{}
	m.startDate = nil
	m.endDate = nil
	m.focusIndex = 0
	m.err = nil
}

// IsActive returns true if any filters are set
func (m *FilterModel) IsActive() bool {
	return len(m.tags) > 0 || m.startDate != nil || m.endDate != nil
}

// Next moves focus to the next field
func (m *FilterModel) Next() {
	m.focusIndex++
}

// Previous moves focus to the previous field
func (m *FilterModel) Previous() {
	if m.focusIndex > 0 {
		m.focusIndex--
	}
}

// GetError returns any error from filter operations
func (m *FilterModel) GetError() error {
	return m.err
}

// SetError sets an error state for the filter model
func (m *FilterModel) SetError(err error) {
	m.err = err
}

// IsTagSelected returns true if the specified tag is selected
func (m *FilterModel) IsTagSelected(tag string) bool {
	for _, t := range m.tags {
		if t == tag {
			return true
		}
	}
	return false
}
