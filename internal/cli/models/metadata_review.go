package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/charmbracelet/bubbletea"
)

// BulkOperationsMsg is sent to open bulk operations for selected events
type BulkOperationsMsg struct {
	Events []*career.CareerEvent
}

// MetadataReviewModel represents the metadata review screen
type MetadataReviewModel struct {
	service       *careerservice.Service
	calculator    *careerservice.DataQualityCalculator
	ctx           context.Context
	events        []*career.CareerEvent
	qualityScores map[string]*careerservice.QualityScore
	selectedIdx   int
	width         int
	height        int
	err           error
	expandedIdx   int    // Index of expanded event (-1 if none)
	filterMode    string // "all", "incomplete"
	sortBy        string // "date", "company", "quality"
	importedEventIDs map[string]bool // IDs of recently imported events
	isImportReview   bool            // True if reviewing only imported events
	fieldOrigins map[string]map[string]bool // eventID -> field -> isFromCSV
	parsingWarnings map[string][]string // eventID -> warnings
	duplicateStatus map[string]string // eventID -> original event ID (empty if not duplicate)
}

// NewMetadataReviewModel creates a new metadata review model
func NewMetadataReviewModel(svc *careerservice.Service, ctx context.Context) *MetadataReviewModel {
	calculator := careerservice.NewDataQualityCalculator()
	model := &MetadataReviewModel{
		service:       svc,
		calculator:    calculator,
		ctx:           ctx,
		selectedIdx:   0,
		expandedIdx:   -1,
		filterMode:    "all",
		sortBy:        "quality",
		qualityScores: make(map[string]*careerservice.QualityScore),
			fieldOrigins: make(map[string]map[string]bool),
		parsingWarnings: make(map[string][]string),
		duplicateStatus: make(map[string]string),
}

	// Load events
	model.loadEvents()

	return model
}

// loadEvents loads events and calculates quality scores
// NewMetadataReviewModelForImport creates a metadata review model for imported events
func NewMetadataReviewModelForImport(svc *careerservice.Service, ctx context.Context, importedEventIDs []string) *MetadataReviewModel {
	calculator := careerservice.NewDataQualityCalculator()
	
	// Convert slice to map for O(1) lookup
	importedMap := make(map[string]bool)
	for _, id := range importedEventIDs {
		importedMap[id] = true
	}
	
	model := &MetadataReviewModel{
		service:          svc,
		calculator:       calculator,
		ctx:              ctx,
		selectedIdx:      0,
		expandedIdx:      -1,
		filterMode:       "all",
		sortBy:           "quality",
		qualityScores:    make(map[string]*careerservice.QualityScore),
		importedEventIDs: importedMap,
		isImportReview:   true,
			fieldOrigins: make(map[string]map[string]bool),
		parsingWarnings: make(map[string][]string),
		duplicateStatus: make(map[string]string),
}

	// Load events (will be filtered to only imported)
	model.loadEvents()

	return model
}


func (m *MetadataReviewModel) loadEvents() {
	filters := careerrepo.ListFilters{
		SortBy:    "date",
		SortOrder: "desc",
		Limit:     100,
	}

	events, err := m.service.ListEvents(m.ctx, filters)
	if err != nil {
		m.err = err
		m.events = []*career.CareerEvent{}
		return
	}

	// Calculate quality scores for all events
	m.qualityScores = make(map[string]*careerservice.QualityScore)
	for _, event := range events {
		score := m.calculator.CalculateQuality(event)
		m.qualityScores[event.ID] = &score
	}

	// Apply filtering
	m.events = m.filterEvents(events)

	// Apply sorting
	m.sortEvents()
}

// filterEvents filters events based on current filter mode
func (m *MetadataReviewModel) filterEvents(events []*career.CareerEvent) []*career.CareerEvent {
	// Filter by imported events if in import review mode
	if m.isImportReview && len(m.importedEventIDs) > 0 {
		filtered := make([]*career.CareerEvent, 0)
		for _, event := range events {
			if m.importedEventIDs[event.ID] {
				filtered = append(filtered, event)
			}
		}
		events = filtered
	}

	if m.filterMode == "incomplete" {
		filtered := make([]*career.CareerEvent, 0)
		for _, event := range events {
			score := m.qualityScores[event.ID]
			if score != nil && (score.Level == careerservice.QualityIncomplete || score.Level == careerservice.QualityBasic) {
				filtered = append(filtered, event)
			}
		}
		return filtered
	}
	return events
}

// sortEvents sorts events based on current sort mode
func (m *MetadataReviewModel) sortEvents() {
	switch m.sortBy {
	case "quality":
		// Sort by quality score (ascending - worst first)
		for i := 0; i < len(m.events)-1; i++ {
			for j := i + 1; j < len(m.events); j++ {
				scoreI := m.qualityScores[m.events[i].ID]
				scoreJ := m.qualityScores[m.events[j].ID]
				if scoreI != nil && scoreJ != nil && scoreI.Score > scoreJ.Score {
					m.events[i], m.events[j] = m.events[j], m.events[i]
				}
			}
		}
	case "company":
		// Sort by company name
		for i := 0; i < len(m.events)-1; i++ {
			for j := i + 1; j < len(m.events); j++ {
				if m.events[i].Company > m.events[j].Company {
					m.events[i], m.events[j] = m.events[j], m.events[i]
				}
			}
		}
	}
}

// Init initializes the model
func (m *MetadataReviewModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *MetadataReviewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "backspace", "q", "esc":
			return m, nil
		case "ctrl+c":
			return m, nil
		case "up", "k":
			m.prevItem()
		case "down", "j":
			m.nextItem()
		case "home", "g":
			m.selectedIdx = 0
		case "end", "G":
			m.selectedIdx = len(m.events) - 1
		case "space":
			// Toggle expand/collapse
			if m.expandedIdx == m.selectedIdx {
				m.expandedIdx = -1
			} else {
				m.expandedIdx = m.selectedIdx
			}
		case "f":
			// Toggle filter mode
			if m.filterMode == "all" {
				m.filterMode = "incomplete"
			} else {
				m.filterMode = "all"
			}
			m.loadEvents()
			m.selectedIdx = 0
		case "s":
			// Cycle through sort modes
			switch m.sortBy {
			case "quality":
				m.sortBy = "date"
			case "date":
				m.sortBy = "company"
			case "company":
				m.sortBy = "quality"
			}
			m.sortEvents()
		case "enter":
			// Edit the selected event
			if m.selectedIdx < len(m.events) {
				return m, func() tea.Msg {
					return EditEventMsg{Event: m.events[m.selectedIdx]}
				}
			}
		case "b":
			// Trigger bulk operations on all events
			return m, func() tea.Msg {
				return BulkOperationsMsg{Events: m.events}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the model
func (m *MetadataReviewModel) View() string {
	var content []string

	// Header
	header := styles.HeaderMain.Render("Metadata Review")
	content = append(content, header)

	// Status bar
	statusText := fmt.Sprintf("Showing %d events | Filter: %s | Sort: %s | Press 'f' to filter, 's' to sort",
		len(m.events), m.filterMode, m.sortBy)
	content = append(content, styles.InputHint.Render(statusText))
	content = append(content, "")

	// Error handling
	if m.err != nil {
		content = append(content, styles.ErrorText.Render(fmt.Sprintf("Error loading events: %v", m.err)))
		return strings.Join(content, "\n")
	}

	// Empty state
	if len(m.events) == 0 {
		content = append(content, styles.InfoText.Render("No events to review. Start capturing events to improve their metadata."))
		return strings.Join(content, "\n")
	}

	// Events list
	for i, event := range m.events {
		var eventContent string

		if i == m.selectedIdx {
			eventContent = m.renderSelectedEvent(event, i)
		} else {
			eventContent = m.renderEventItem(event, i)
		}

		content = append(content, eventContent)

		// Show expanded view if selected
		if i == m.expandedIdx {
			content = append(content, m.renderExpandedEvent(event))
		}
	}

	// Footer
	footer := styles.InputHint.Render("↑/↓ navigate | Space expand | f filter | s sort | b bulk | Backspace back")
	content = append(content, "")
	content = append(content, footer)

	return strings.Join(content, "\n")
}

// renderEventItem renders a compact event item
func (m *MetadataReviewModel) renderEventItem(event *career.CareerEvent, idx int) string {
	score := m.qualityScores[event.ID]
	if score == nil {
		return ""
	}

	// Truncate text to 60 chars
	text := event.Text
	if len(text) > 60 {
		text = text[:57] + "..."
	}

	// Format date
	dateStr := event.Date.Format("2006-01-02")

	// Quality level
	quality := string(score.Level)

	// Build item
	item := fmt.Sprintf("  %s | %s | %s | %s (%d%%)",
		text,
		dateStr,
		event.Company,
		quality,
		score.Score,
	)

	return styles.ListItem.Render(item)
}

// renderSelectedEvent renders a selected event with highlight
func (m *MetadataReviewModel) renderSelectedEvent(event *career.CareerEvent, idx int) string {
	score := m.qualityScores[event.ID]
	if score == nil {
		return ""
	}

	// Truncate text to 60 chars
	text := event.Text
	if len(text) > 60 {
		text = text[:57] + "..."
	}

	// Format date
	dateStr := event.Date.Format("2006-01-02")

	// Quality level
	quality := string(score.Level)

	// Build item with selection marker
	item := fmt.Sprintf("▶ %s | %s | %s | %s (%d%%)",
		text,
		dateStr,
		event.Company,
		quality,
		score.Score,
	)

	return styles.ListItemSelected.Render(item)
}

// renderExpandedEvent renders the full event details
func (m *MetadataReviewModel) renderExpandedEvent(event *career.CareerEvent) string {
	score := m.qualityScores[event.ID]
	if score == nil {
		return ""
	}

	var details []string
	details = append(details, "")
	details = append(details, "  Full Details:")
	details = append(details, fmt.Sprintf("    Text: %s", event.Text))
	details = append(details, fmt.Sprintf("    Date: %s", event.Date.Format("2006-01-02")))
	details = append(details, fmt.Sprintf("    Company: %s", event.Company))
	details = append(details, fmt.Sprintf("    Project: %s", event.Project))

	if len(event.Tags) > 0 {
		details = append(details, fmt.Sprintf("    Tags: %s", strings.Join(event.Tags, ", ")))
	}

	if len(event.Categories) > 0 {
		details = append(details, fmt.Sprintf("    Categories: %s", strings.Join(event.Categories, ", ")))
	}

	// Quality score details
	details = append(details, "")
	details = append(details, "  Quality Score:")
	details = append(details, fmt.Sprintf("    Text: %d/20", score.TextScore))
	details = append(details, fmt.Sprintf("    Date: %d/20", score.DateScore))
	details = append(details, fmt.Sprintf("    Company: %d/10", score.CompanyScore))
	details = append(details, fmt.Sprintf("    Project: %d/10", score.ProjectScore))
	details = append(details, fmt.Sprintf("    Tags: %d/15", score.TagsScore))
	details = append(details, fmt.Sprintf("    Categories: %d/15", score.CategoriesScore))
	details = append(details, fmt.Sprintf("    Match: %d/10", score.MatchScore))

	return styles.CardContent.Render(strings.Join(details, "\n"))
}

// prevItem moves selection to previous item
func (m *MetadataReviewModel) prevItem() {
	if m.selectedIdx > 0 {
		m.selectedIdx--
		m.expandedIdx = -1 // Collapse on navigation
	}
}

// nextItem moves selection to next item
func (m *MetadataReviewModel) nextItem() {
	if m.selectedIdx < len(m.events)-1 {
		m.selectedIdx++
		m.expandedIdx = -1 // Collapse on navigation
	}
}

// GetSelectedEvent returns the currently selected event
func (m *MetadataReviewModel) GetSelectedEvent() *career.CareerEvent {
	if m.selectedIdx >= 0 && m.selectedIdx < len(m.events) {
		return m.events[m.selectedIdx]
	}
	return nil
}

// GetEvents returns the current list of events
func (m *MetadataReviewModel) GetEvents() []*career.CareerEvent {
	return m.events
}

// Refresh reloads events from service
func (m *MetadataReviewModel) Refresh() {
	m.loadEvents()
	if m.selectedIdx >= len(m.events) {
		m.selectedIdx = len(m.events) - 1
	}
	if m.selectedIdx < 0 {
		m.selectedIdx = 0
	}
}

// SetFieldOrigins sets the field origins for an imported event
func (m *MetadataReviewModel) SetFieldOrigins(eventID string, origins map[string]bool) {
	if m.fieldOrigins == nil {
		m.fieldOrigins = make(map[string]map[string]bool)
	}
	m.fieldOrigins[eventID] = origins
}

// GetFieldOrigins returns the field origins for an event
func (m *MetadataReviewModel) GetFieldOrigins(eventID string) map[string]bool {
	if m.fieldOrigins == nil {
		return nil
	}
	return m.fieldOrigins[eventID]
}

// IsFieldFromCSV returns whether a field came from CSV (true) or is a default (false)
func (m *MetadataReviewModel) IsFieldFromCSV(eventID, field string) bool {
	origins := m.GetFieldOrigins(eventID)
	if origins == nil {
		return false // Default assumption if no origin info
	}
	fromCSV, exists := origins[field]
	return exists && fromCSV
}

// GetDefaultFields returns the list of fields that are defaults for an event
func (m *MetadataReviewModel) GetDefaultFields(eventID string) []string {
	origins := m.GetFieldOrigins(eventID)
	if origins == nil {
		return []string{}
	}
	var defaults []string
	for field, fromCSV := range origins {
		if !fromCSV {
			defaults = append(defaults, field)
		}
	}
	return defaults
}

// HasDefaultFields checks if an event has any default fields
func (m *MetadataReviewModel) HasDefaultFields(eventID string) bool {
	return len(m.GetDefaultFields(eventID)) > 0
}

// SetParsingWarnings sets the parsing warnings for an imported event
func (m *MetadataReviewModel) SetParsingWarnings(eventID string, warnings []string) {
	if m.parsingWarnings == nil {
		m.parsingWarnings = make(map[string][]string)
	}
	m.parsingWarnings[eventID] = warnings
}

// GetParsingWarnings returns the parsing warnings for an event
func (m *MetadataReviewModel) GetParsingWarnings(eventID string) []string {
	if m.parsingWarnings == nil {
		return []string{}
	}
	warnings, exists := m.parsingWarnings[eventID]
	if !exists {
		return []string{}
	}
	return warnings
}

// HasParsingWarnings checks if an event has parsing warnings
func (m *MetadataReviewModel) HasParsingWarnings(eventID string) bool {
	return len(m.GetParsingWarnings(eventID)) > 0
}

// SetDuplicateStatus sets the duplicate status for an event
func (m *MetadataReviewModel) SetDuplicateStatus(eventID string, isDuplicate bool, originalEventID string) {
	if m.duplicateStatus == nil {
		m.duplicateStatus = make(map[string]string)
	}
	if isDuplicate {
		m.duplicateStatus[eventID] = originalEventID
	} else {
		m.duplicateStatus[eventID] = ""
	}
}

// GetDuplicateStatus returns whether an event is a duplicate and its original event ID
func (m *MetadataReviewModel) GetDuplicateStatus(eventID string) (bool, string) {
	if m.duplicateStatus == nil {
		return false, ""
	}
	originalID, exists := m.duplicateStatus[eventID]
	if !exists || originalID == "" {
		return false, ""
	}
	return true, originalID
}

// IsDuplicate checks if an event is a duplicate
func (m *MetadataReviewModel) IsDuplicate(eventID string) bool {
	isDuplicate, _ := m.GetDuplicateStatus(eventID)
	return isDuplicate
}

// GetOriginalEventID returns the original event ID if this is a duplicate
func (m *MetadataReviewModel) GetOriginalEventID(eventID string) string {
	_, originalID := m.GetDuplicateStatus(eventID)
	return originalID
}
