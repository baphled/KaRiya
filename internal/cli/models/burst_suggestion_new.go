package models

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burst_fact"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// BurstSuggestionModelNew represents the burst suggestion review and confirmation screen using huh forms.
// This model uses the huh library for the edit mode, providing:
// - Automatic focus management during editing
// - Built-in validation
// - Catppuccin theming
// - Consistent keyboard navigation
type BurstSuggestionModelNew struct {
	*BaseStandardModel
	service       *careerservice.Service
	ctx           context.Context
	suggestions   []burstfact.BurstSuggestion
	currentIdx    int
	confirmed     []burstfact.BurstSuggestion
	rejected      []burstfact.BurstSuggestion
	editing       bool
	editForm      *huh.Form
	editFormData  *forms.BurstSuggestionFormData
	editedNames   map[int]string                // Map of suggestion index to edited name
	editedDescs   map[int]string                // Map of suggestion index to edited description
	relatedEvents map[int][]*career.CareerEvent // Cache of related events
	width         int
	height        int
}

// NewBurstSuggestionModelNew creates a new burst suggestion model using huh forms for editing.
func NewBurstSuggestionModelNew(svc *careerservice.Service, suggestions []burstfact.BurstSuggestion, ctx context.Context) *BurstSuggestionModelNew {
	return &BurstSuggestionModelNew{
		BaseStandardModel: NewBaseStandardModel(),
		service:           svc,
		ctx:               ctx,
		suggestions:       suggestions,
		currentIdx:        0,
		confirmed:         []burstfact.BurstSuggestion{},
		rejected:          []burstfact.BurstSuggestion{},
		editing:           false,
		editForm:          nil,
		editFormData:      nil,
		editedNames:       make(map[int]string),
		editedDescs:       make(map[int]string),
		relatedEvents:     make(map[int][]*career.CareerEvent),
		width:             80,
		height:            24,
	}
}

// Init initializes the model
func (m *BurstSuggestionModelNew) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *BurstSuggestionModelNew) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	}

	return m, nil
}

// handleKeyMsg processes keyboard input
func (m *BurstSuggestionModelNew) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.editing {
		return m.handleEditKeyMsg(msg)
	}

	switch msg.Type {
	case tea.KeyUp:
		if len(m.suggestions) > 0 {
			m.currentIdx = (m.currentIdx - 1 + len(m.suggestions)) % len(m.suggestions)
			// Clear the related events cache for the new suggestion to load fresh
			delete(m.relatedEvents, m.currentIdx)
		}
	case tea.KeyDown:
		if len(m.suggestions) > 0 {
			m.currentIdx = (m.currentIdx + 1) % len(m.suggestions)
			delete(m.relatedEvents, m.currentIdx)
		}
	case tea.KeyEsc:
		// Exit suggestion review (parent will handle navigation)
		return m, func() tea.Msg { return BackMsg{} }
	}

	// Handle character input for confirm/reject/edit
	if msg.Type == tea.KeyRunes {
		for _, r := range msg.Runes {
			switch r {
			case 'y', 'Y':
				return m.confirmCurrent()
			case 'n', 'N':
				return m.rejectCurrent()
			case 'e', 'E':
				return m.startEdit()
			}
		}
	}

	return m, nil
}

// handleEditKeyMsg handles keyboard input while editing name/description using huh form
func (m *BurstSuggestionModelNew) handleEditKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.editForm == nil {
		// Shouldn't happen, but safety check
		m.editing = false
		return m, nil
	}

	// Update the form
	form, cmd := m.editForm.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.editForm = f
	}

	// Check if form was completed
	if forms.IsCompleted(m.editForm) {
		return m.saveEdits()
	}

	// Check if form was aborted (Esc)
	if forms.IsAborted(m.editForm) {
		m.editing = false
		m.editForm = nil
		m.editFormData = nil
		return m, nil
	}

	return m, cmd
}

// startEdit begins editing the current suggestion's name/description using huh form
func (m *BurstSuggestionModelNew) startEdit() (tea.Model, tea.Cmd) {
	if len(m.suggestions) == 0 {
		return m, nil
	}

	m.editing = true
	current := m.suggestions[m.currentIdx]

	// Create form data with current or edited values
	m.editFormData = &forms.BurstSuggestionFormData{
		Name:        current.Name,
		Description: current.Description,
	}

	// Use edited values if they exist
	if name, exists := m.editedNames[m.currentIdx]; exists {
		m.editFormData.Name = name
	}
	if desc, exists := m.editedDescs[m.currentIdx]; exists {
		m.editFormData.Description = desc
	}

	// Create huh form
	m.editForm = forms.NewBurstSuggestionEditFormWithData(m.editFormData)

	return m, m.editForm.Init()
}

// saveEdits saves the edited name/description from the huh form
func (m *BurstSuggestionModelNew) saveEdits() (tea.Model, tea.Cmd) {
	if m.editFormData != nil {
		m.editedNames[m.currentIdx] = m.editFormData.Name
		m.editedDescs[m.currentIdx] = m.editFormData.Description
	}

	m.editing = false
	m.editForm = nil
	m.editFormData = nil
	return m, nil
}

// confirmCurrent confirms the current suggestion and moves to next
func (m *BurstSuggestionModelNew) confirmCurrent() (tea.Model, tea.Cmd) {
	if len(m.suggestions) == 0 {
		return m, nil
	}

	current := m.suggestions[m.currentIdx]

	// Apply edits if any
	if name, exists := m.editedNames[m.currentIdx]; exists {
		current.Name = name
	}
	if desc, exists := m.editedDescs[m.currentIdx]; exists {
		current.Description = desc
	}

	m.confirmed = append(m.confirmed, current)

	// Create a burst from the suggestion
	burst := m.createBurstFromSuggestion(current)

	// Send confirm message for this burst
	confirmCmd := func() tea.Msg {
		return ConfirmBurstMsg{Burst: burst}
	}

	// Move to next suggestion or complete
	if m.currentIdx < len(m.suggestions)-1 {
		m.currentIdx++
		return m, confirmCmd
	} else {
		// All suggestions processed, send completion message
		completeCmd := tea.Batch(
			confirmCmd,
			func() tea.Msg {
				return BurstProcessingCompleteMsg{
					ConfirmedCount: len(m.confirmed),
					RejectedCount:  len(m.rejected),
				}
			},
		)
		return m, completeCmd
	}
}

// rejectCurrent rejects the current suggestion and moves to next
func (m *BurstSuggestionModelNew) rejectCurrent() (tea.Model, tea.Cmd) {
	if len(m.suggestions) == 0 {
		return m, nil
	}

	current := m.suggestions[m.currentIdx]
	m.rejected = append(m.rejected, current)

	// Send reject message for this suggestion
	rejectCmd := func() tea.Msg {
		return RejectBurstSuggestionMsg{Suggestion: current}
	}

	// Move to next suggestion or complete
	if m.currentIdx < len(m.suggestions)-1 {
		m.currentIdx++
		return m, rejectCmd
	} else {
		// All suggestions processed, send completion message
		completeCmd := tea.Batch(
			rejectCmd,
			func() tea.Msg {
				return BurstProcessingCompleteMsg{
					ConfirmedCount: len(m.confirmed),
					RejectedCount:  len(m.rejected),
				}
			},
		)
		return m, completeCmd
	}
}

// View renders the burst suggestion screen
func (m *BurstSuggestionModelNew) View() string {
	if len(m.suggestions) == 0 {
		return styles.ErrorBox.Render("No burst suggestions available")
	}

	if m.editing {
		return m.renderEditView()
	}

	return m.renderReviewView()
}

// renderReviewView renders the suggestion review view
func (m *BurstSuggestionModelNew) renderReviewView() string {
	current := m.suggestions[m.currentIdx]

	var parts []string

	// Header
	header := styles.HeaderMain.Render(fmt.Sprintf("Burst Suggestion %d of %d", m.currentIdx+1, len(m.suggestions)))
	parts = append(parts, header)

	// Progress bar
	progressBar := m.renderProgressBar()
	parts = append(parts, progressBar)

	// Confidence score with visualization
	confidenceVis := m.renderConfidenceScore(current.ConfidenceScore)
	parts = append(parts, confidenceVis)

	// Burst details
	burstDetails := m.renderBurstDetails(current)
	parts = append(parts, burstDetails)

	// Related events preview
	relatedEventsView := m.renderRelatedEvents(current)
	parts = append(parts, relatedEventsView)

	// Navigation help
	helpText := lipgloss.NewStyle().Render("↑/↓ navigate • y confirm • n reject • e edit name/desc • Esc back")
	parts = append(parts, helpText)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// renderEditView renders the editing view using huh form
func (m *BurstSuggestionModelNew) renderEditView() string {
	if m.editForm == nil {
		return styles.ErrorBox.Render("Edit form not initialized")
	}

	var parts []string

	header := styles.HeaderMain.Render("Edit Burst Name & Description")
	parts = append(parts, header)

	// Render huh form
	formView := m.editForm.View()
	parts = append(parts, formView)

	// Help text
	helpText := lipgloss.NewStyle().Render("Tab navigate • Enter save • Esc cancel")
	parts = append(parts, helpText)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// renderProgressBar renders the progress through suggestions
func (m *BurstSuggestionModelNew) renderProgressBar() string {
	total := len(m.suggestions)
	current := m.currentIdx + 1
	progress := float64(current) / float64(total)
	barWidth := 30

	filledWidth := int(math.Round(progress * float64(barWidth)))
	bar := strings.Repeat("█", filledWidth) + strings.Repeat("░", barWidth-filledWidth)

	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		fmt.Sprintf("[%s] %d/%d", bar, current, total),
	)
}

// renderConfidenceScore renders the confidence score with color coding
func (m *BurstSuggestionModelNew) renderConfidenceScore(score float64) string {
	percentage := int(score * 100)
	var scoreStyle lipgloss.Style

	switch {
	case score >= 0.8:
		scoreStyle = styles.SuccessBox
	case score >= 0.6:
		scoreStyle = styles.InfoBox
	case score >= 0.4:
		scoreStyle = styles.WarningBox
	default:
		scoreStyle = styles.ErrorBox
	}

	scoreText := fmt.Sprintf("Confidence: %d%% ", percentage)
	scoreBar := strings.Repeat("█", percentage/5) + strings.Repeat("░", 20-percentage/5)

	content := fmt.Sprintf("%s[%s]", scoreText, scoreBar)
	return scoreStyle.Render(content)
}

// renderBurstDetails renders the burst name and description
func (m *BurstSuggestionModelNew) renderBurstDetails(suggestion burstfact.BurstSuggestion) string {
	var parts []string

	// Name or suggested name
	name := suggestion.Name
	if editedName, exists := m.editedNames[m.currentIdx]; exists {
		name = editedName
	}

	if name != "" {
		nameBox := styles.CardBase.Render(fmt.Sprintf("Name: %s", name))
		parts = append(parts, nameBox)
	} else {
		nameBox := styles.WarningBox.Render("Name: (not set - will be auto-generated)")
		parts = append(parts, nameBox)
	}

	// Description
	desc := suggestion.Description
	if editedDesc, exists := m.editedDescs[m.currentIdx]; exists {
		desc = editedDesc
	}

	if desc != "" {
		descBox := styles.CardBase.Render(fmt.Sprintf("Description: %s", desc))
		parts = append(parts, descBox)
	}

	// Event count
	eventCountBox := styles.InfoBox.Render(fmt.Sprintf("Related Events: %d", len(suggestion.EventIDs)))
	parts = append(parts, eventCountBox)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// renderRelatedEvents renders preview of related events
func (m *BurstSuggestionModelNew) renderRelatedEvents(suggestion burstfact.BurstSuggestion) string {
	if len(suggestion.EventIDs) == 0 {
		return styles.WarningBox.Render("No related events")
	}

	var events []*career.CareerEvent

	// Try to load from cache first
	if cachedEvents, exists := m.relatedEvents[m.currentIdx]; exists {
		events = cachedEvents
	} else {
		// Load events from service
		for _, eventID := range suggestion.EventIDs {
			event, err := m.service.GetEventByID(m.ctx, eventID)
			if err == nil && event != nil {
				events = append(events, event)
			}
		}
		// Cache the loaded events
		m.relatedEvents[m.currentIdx] = events
	}

	var eventLines []string
	for i, event := range events {
		if i >= 3 { // Show max 3 events
			eventLines = append(eventLines, fmt.Sprintf("  ... and %d more", len(events)-3))
			break
		}

		// Truncate text to 60 chars
		text := event.Text
		if len(text) > 60 {
			text = text[:57] + "..."
		}

		eventLine := fmt.Sprintf("  • %s (%s)", text, event.Date.Format("2006-01-02"))
		eventLines = append(eventLines, eventLine)
	}

	header := "Related Events:"
	content := header + "\n" + strings.Join(eventLines, "\n")
	return styles.CardBase.Render(content)
}

// GetConfirmed returns the list of confirmed suggestions
func (m *BurstSuggestionModelNew) GetConfirmed() []burstfact.BurstSuggestion {
	return m.confirmed
}

// GetRejected returns the list of rejected suggestions
func (m *BurstSuggestionModelNew) GetRejected() []burstfact.BurstSuggestion {
	return m.rejected
}

// IsDone returns true if all suggestions have been processed
func (m *BurstSuggestionModelNew) IsDone() bool {
	return len(m.confirmed)+len(m.rejected) == len(m.suggestions)
}

// createBurstFromSuggestion converts a BurstSuggestion into a Burst domain object
func (m *BurstSuggestionModelNew) createBurstFromSuggestion(suggestion burstfact.BurstSuggestion) *career.Burst {
	// Generate a name if none provided
	name := suggestion.Name
	if name == "" {
		name = fmt.Sprintf("Burst of %d events", len(suggestion.EventIDs))
	}

	return &career.Burst{
		Name:        name,
		Description: suggestion.Description,
		EventIDs:    suggestion.EventIDs,
	}
}
