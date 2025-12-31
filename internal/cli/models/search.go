package models

import (
	"strings"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

// HighlightPosition represents a start and end position for highlighting
type HighlightPosition struct {
	Start int
	End   int
}

// SearchModel manages event search with debouncing and highlighting
type SearchModel struct {
	*BaseStandardModel
	query          string
	lastUpdate     time.Time
	debounceDelay  time.Duration
	err            error
	lastSearchTime time.Time
}

// NewSearchModel creates a new search model with default debounce delay
func NewSearchModel() *SearchModel {
	return &SearchModel{
		BaseStandardModel: NewBaseStandardModel(),
		query:             "",
		lastUpdate:        time.Now(),
		debounceDelay:     300 * time.Millisecond,
		err:               nil,
		lastSearchTime:    time.Now(),
	}
}

// SetQuery sets the search query and clears any previous errors
func (m *SearchModel) SetQuery(query string) {
	m.query = query
	m.lastUpdate = time.Now()
	m.err = nil
	m.BaseStandardModel.ClearError()
}

// GetQuery returns the current search query
func (m *SearchModel) GetQuery() string {
	return m.query
}

// Matches checks if an event matches the current search query
func (m *SearchModel) Matches(event *career.CareerEvent) bool {
	if m.query == "" {
		return true
	}

	searchLower := strings.ToLower(m.query)

	// Check event text
	if strings.Contains(strings.ToLower(event.Text), searchLower) {
		return true
	}

	// Check company
	if strings.Contains(strings.ToLower(event.Company), searchLower) {
		return true
	}

	// Check project
	if strings.Contains(strings.ToLower(event.Project), searchLower) {
		return true
	}

	return false
}

// IsActive returns true if a search query is active (non-empty)
func (m *SearchModel) IsActive() bool {
	return m.query != ""
}

// GetHighlightPositions returns the start and end positions of all matches in text
func (m *SearchModel) GetHighlightPositions(text string) []HighlightPosition {
	if m.query == "" {
		return []HighlightPosition{}
	}

	var positions []HighlightPosition
	searchLower := strings.ToLower(m.query)
	textLower := strings.ToLower(text)

	start := 0
	for {
		idx := strings.Index(textLower[start:], searchLower)
		if idx == -1 {
			break
		}

		actualStart := start + idx
		actualEnd := actualStart + len(m.query)

		positions = append(positions, HighlightPosition{
			Start: actualStart,
			End:   actualEnd,
		})

		start = actualEnd
	}

	return positions
}

// SetError sets an error state for the search
func (m *SearchModel) SetError(err error) {
	m.err = err
	// Also set in the base standard model for unified error tracking
	m.BaseStandardModel.SetError(err)
}

// GetError returns any error from search operations
func (m *SearchModel) GetError() error {
	return m.err
}

// ClearError clears any error state
func (m *SearchModel) ClearError() {
	m.err = nil
	m.BaseStandardModel.ClearError()
}

// Reset clears all search state
func (m *SearchModel) Reset() {
	m.query = ""
	m.err = nil
	m.lastUpdate = time.Now()
	m.lastSearchTime = time.Now()
	m.BaseStandardModel.Reset()
}

// GetLastUpdate returns the time of the last query update
func (m *SearchModel) GetLastUpdate() time.Time {
	return m.lastUpdate
}

// GetDebounceDelay returns the debounce delay duration
func (m *SearchModel) GetDebounceDelay() time.Duration {
	return m.debounceDelay
}

// ShouldSearch returns true if enough time has passed since the last update
// to perform a new search (respecting debounce delay)
func (m *SearchModel) ShouldSearch() bool {
	elapsed := time.Since(m.lastUpdate)
	return elapsed >= m.debounceDelay
}

// MarkSearched marks that a search has been performed
func (m *SearchModel) MarkSearched() {
	m.lastSearchTime = time.Now()
}

// GetLastSearchTime returns the time of the last performed search
func (m *SearchModel) GetLastSearchTime() time.Time {
	return m.lastSearchTime
}
