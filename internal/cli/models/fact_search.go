package models

import (
	"strings"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
)

// FactSearchModel manages fact search with advanced filtering capabilities
type FactSearchModel struct {
	*BaseStandardModel
	query              string
	competencyFilter   string
	roleFitFilter      career.RoleFit
	audienceFilter     string
	lastUpdate         time.Time
	debounceDelay      time.Duration
	err                error
	lastSearchTime     time.Time
}

// NewFactSearchModel creates a new fact search model with default debounce delay
func NewFactSearchModel() *FactSearchModel {
	return &FactSearchModel{
		BaseStandardModel: NewBaseStandardModel(),
		query:             "",
		competencyFilter:  "",
		roleFitFilter:     "",
		audienceFilter:    "",
		lastUpdate:        time.Now(),
		debounceDelay:     300 * time.Millisecond,
		err:               nil,
		lastSearchTime:    time.Now(),
	}
}

// SetQuery sets the search query and clears any previous errors
func (m *FactSearchModel) SetQuery(query string) {
	m.query = query
	m.lastUpdate = time.Now()
	m.err = nil
	m.BaseStandardModel.ClearError()
}

// GetQuery returns the current search query
func (m *FactSearchModel) GetQuery() string {
	return m.query
}

// SetCompetencyFilter sets the competency filter
func (m *FactSearchModel) SetCompetencyFilter(competency string) {
	m.competencyFilter = competency
	m.lastUpdate = time.Now()
}

// GetCompetencyFilter returns the current competency filter
func (m *FactSearchModel) GetCompetencyFilter() string {
	return m.competencyFilter
}

// SetRoleFitFilter sets the role fit filter
func (m *FactSearchModel) SetRoleFitFilter(roleFit career.RoleFit) {
	m.roleFitFilter = roleFit
	m.lastUpdate = time.Now()
}

// GetRoleFitFilter returns the current role fit filter
func (m *FactSearchModel) GetRoleFitFilter() career.RoleFit {
	return m.roleFitFilter
}

// SetAudienceFilter sets the audience filter
func (m *FactSearchModel) SetAudienceFilter(audience string) {
	m.audienceFilter = audience
	m.lastUpdate = time.Now()
}

// GetAudienceFilter returns the current audience filter
func (m *FactSearchModel) GetAudienceFilter() string {
	return m.audienceFilter
}

// Matches checks if a fact matches the current search criteria
func (m *FactSearchModel) Matches(fact *career.Fact) bool {
	// Check text query match (case-insensitive)
	if m.query != "" {
		searchLower := strings.ToLower(m.query)
		textLower := strings.ToLower(fact.Text)

		if !strings.Contains(textLower, searchLower) {
			return false
		}
	}

	// Check competency filter
	if m.competencyFilter != "" {
		hasCompetency := false
		for _, comp := range fact.CompetencyCategories {
			if strings.EqualFold(comp, m.competencyFilter) {
				hasCompetency = true
				break
			}
		}
		if !hasCompetency {
			return false
		}
	}

	// Check role fit filter
	if m.roleFitFilter != "" && fact.RoleFit != m.roleFitFilter {
		return false
	}

	// Check audience filter
	if m.audienceFilter != "" {
		hasAudience := false
		for _, aud := range fact.AudienceRelevance {
			if strings.EqualFold(aud, m.audienceFilter) {
				hasAudience = true
				break
			}
		}
		if !hasAudience {
			return false
		}
	}

	return true
}

// IsActive returns true if any search filter is active
func (m *FactSearchModel) IsActive() bool {
	return m.query != "" ||
		m.competencyFilter != "" ||
		m.roleFitFilter != "" ||
		m.audienceFilter != ""
}

// ClearFilters clears all search filters
func (m *FactSearchModel) ClearFilters() {
	m.query = ""
	m.competencyFilter = ""
	m.roleFitFilter = ""
	m.audienceFilter = ""
	m.err = nil
	m.lastUpdate = time.Now()
	m.BaseStandardModel.ClearError()
}

// SetError sets an error state for the search
func (m *FactSearchModel) SetError(err error) {
	m.err = err
	m.BaseStandardModel.SetError(err)
}

// GetError returns any error from search operations
func (m *FactSearchModel) GetError() error {
	return m.err
}

// ClearError clears any error state
func (m *FactSearchModel) ClearError() {
	m.err = nil
	m.BaseStandardModel.ClearError()
}

// Reset clears all search state
func (m *FactSearchModel) Reset() {
	m.query = ""
	m.competencyFilter = ""
	m.roleFitFilter = ""
	m.audienceFilter = ""
	m.err = nil
	m.lastUpdate = time.Now()
	m.lastSearchTime = time.Now()
	m.BaseStandardModel.Reset()
}

// GetLastUpdate returns the time of the last filter update
func (m *FactSearchModel) GetLastUpdate() time.Time {
	return m.lastUpdate
}

// GetDebounceDelay returns the debounce delay duration
func (m *FactSearchModel) GetDebounceDelay() time.Duration {
	return m.debounceDelay
}

// ShouldSearch returns true if enough time has passed since the last update
// to perform a new search (respecting debounce delay)
func (m *FactSearchModel) ShouldSearch() bool {
	elapsed := time.Since(m.lastUpdate)
	return elapsed >= m.debounceDelay
}

// MarkSearched marks that a search has been performed
func (m *FactSearchModel) MarkSearched() {
	m.lastSearchTime = time.Now()
}

// GetLastSearchTime returns the time of the last performed search
func (m *FactSearchModel) GetLastSearchTime() time.Time {
	return m.lastSearchTime
}

