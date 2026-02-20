package captureevent

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	domcapture "github.com/baphled/kariya/internal/domain/capture"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burstfact"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// BurstSuggestionModelNew represents the burst suggestion review and confirmation screen using huh forms.
//
// This model uses the huh library for the edit mode, providing:
//   - Automatic focus management during editing
//   - Built-in validation
//   - Catppuccin theming
//   - Consistent keyboard navigation.
type BurstSuggestionModelNew struct {
	service       *careerservice.Service
	ctx           context.Context
	suggestions   []burstfact.BurstSuggestion
	currentIdx    int
	confirmed     []burstfact.BurstSuggestion
	rejected      []burstfact.BurstSuggestion
	editing       bool
	editForm      forms.Form
	editFormData  *forms.BurstSuggestionFormData
	editedNames   map[int]string
	editedDescs   map[int]string
	relatedEvents map[int][]*career.Event
	width         int
	height        int
}

// NewBurstSuggestionModelNew creates a new burst suggestion model using huh forms for editing.
//
// Expected:
//   - ctx must be a valid context.
//   - svc must be a valid careerservice.Service pointer.
//   - suggestions must be a valid slice of BurstSuggestion.
//
// Returns:
//   - A fully initialized BurstSuggestionModelNew ready for use.
//
// Side effects:
//   - None.
func NewBurstSuggestionModelNew(
	ctx context.Context, svc *careerservice.Service, suggestions []burstfact.BurstSuggestion,
) *BurstSuggestionModelNew {
	return &BurstSuggestionModelNew{
		service:       svc,
		ctx:           ctx,
		suggestions:   suggestions,
		currentIdx:    0,
		confirmed:     []burstfact.BurstSuggestion{},
		rejected:      []burstfact.BurstSuggestion{},
		editing:       false,
		editForm:      nil,
		editFormData:  nil,
		editedNames:   make(map[int]string),
		editedDescs:   make(map[int]string),
		relatedEvents: make(map[int][]*career.Event),
		width:         80,
		height:        24,
	}
}

// Init initializes the model.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *BurstSuggestionModelNew) Init() tea.Cmd {
	return nil
}

// Update handles messages.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - tea.Model: the updated model.
//   - tea.Cmd: command to execute.
//
// Side effects:
//   - May update internal state based on message type.
func (m *BurstSuggestionModelNew) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleBurstKeyMsg(msg)
	}

	return m, nil
}

// handleBurstKeyMsg processes keyboard input.
func (m *BurstSuggestionModelNew) handleBurstKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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

	// Handle character input for confirm/reject/edit and vim navigation
	if msg.Type == tea.KeyRunes {
		for _, r := range msg.Runes {
			switch r {
			case 'j':
				if len(m.suggestions) > 0 {
					m.currentIdx = (m.currentIdx + 1) % len(m.suggestions)
					delete(m.relatedEvents, m.currentIdx)
				}
			case 'k':
				if len(m.suggestions) > 0 {
					m.currentIdx = (m.currentIdx - 1 + len(m.suggestions)) % len(m.suggestions)
					delete(m.relatedEvents, m.currentIdx)
				}
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

// handleEditKeyMsg handles keyboard input while editing name/description using huh form.
func (m *BurstSuggestionModelNew) handleEditKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.editForm == nil {
		// Shouldn't happen, but safety check
		m.editing = false
		return m, nil
	}

	// Update the form
	form, cmd := m.editForm.Update(msg)
	if f, ok := form.(forms.Form); ok {
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

// startEdit begins editing the current suggestion's name/description using huh form.
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
	m.editForm = forms.NewBurstSuggestionForm(m.editFormData)

	return m, m.editForm.Init()
}

// saveEdits saves the edited name/description from the huh form.
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

// confirmCurrent confirms the current suggestion and moves to next.
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
	}

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

// rejectCurrent rejects the current suggestion and moves to next.
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
	}

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

// View renders the burst suggestion screen
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *BurstSuggestionModelNew) View() string {
	if len(m.suggestions) == 0 {
		return containers.NewBox(m.getTheme()).
			Variant(containers.BoxDestructive).
			Content("No burst suggestions available").
			Render()
	}

	if m.editing {
		return m.renderEditView()
	}

	return m.renderReviewView()
}

// renderReviewView renders the suggestion review view.
func (m *BurstSuggestionModelNew) renderReviewView() string {
	current := m.suggestions[m.currentIdx]
	theme := m.getTheme()

	var parts []string

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.AccentColor())
	header := headerStyle.Render(fmt.Sprintf("Burst Suggestion %d of %d", m.currentIdx+1, len(m.suggestions)))
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

// renderEditView renders the editing view using huh form.
func (m *BurstSuggestionModelNew) renderEditView() string {
	theme := m.getTheme()

	if m.editForm == nil {
		return containers.NewBox(theme).
			Variant(containers.BoxDestructive).
			Content("Edit form not initialized").
			Render()
	}

	var parts []string

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.AccentColor())
	header := headerStyle.Render("Edit Burst Name & Description")
	parts = append(parts, header)

	// Render huh form
	formView := m.editForm.View()
	parts = append(parts, formView)

	// Help text
	helpText := lipgloss.NewStyle().Render("Tab navigate • Enter save • Esc cancel")
	parts = append(parts, helpText)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// renderProgressBar renders the progress through suggestions.
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

// renderConfidenceScore renders the confidence score with color coding.
func (m *BurstSuggestionModelNew) renderConfidenceScore(score float64) string {
	theme := m.getTheme()
	percentage := int(score * 100)

	var variant containers.BoxVariant
	switch {
	case score >= 0.8:
		variant = containers.BoxSuccess
	case score >= 0.6:
		variant = containers.BoxInfo
	case score >= 0.4:
		variant = containers.BoxWarning
	default:
		variant = containers.BoxDestructive
	}

	scoreText := fmt.Sprintf("Confidence: %d%% ", percentage)
	scoreBar := strings.Repeat("█", percentage/5) + strings.Repeat("░", 20-percentage/5)

	content := fmt.Sprintf("%s[%s]", scoreText, scoreBar)
	return containers.NewBox(theme).Variant(variant).Content(content).Render()
}

// renderBurstDetails renders the burst name and description.
func (m *BurstSuggestionModelNew) renderBurstDetails(suggestion burstfact.BurstSuggestion) string {
	theme := m.getTheme()
	var parts []string

	// Name or suggested name
	name := suggestion.Name
	if editedName, exists := m.editedNames[m.currentIdx]; exists {
		name = editedName
	}

	if name != "" {
		nameBox := containers.NewBox(theme).Content("Name: " + name).Render()
		parts = append(parts, nameBox)
	} else {
		nameBox := containers.NewBox(theme).
			Variant(containers.BoxWarning).
			Content("Name: (not set - will be auto-generated)").
			Render()
		parts = append(parts, nameBox)
	}

	// Description
	desc := suggestion.Description
	if editedDesc, exists := m.editedDescs[m.currentIdx]; exists {
		desc = editedDesc
	}

	if desc != "" {
		descBox := containers.NewBox(theme).Content("Description: " + desc).Render()
		parts = append(parts, descBox)
	}

	// Event count
	eventCountBox := containers.NewBox(theme).
		Variant(containers.BoxInfo).
		Content(fmt.Sprintf("Related Events: %d", len(suggestion.EventIDs))).
		Render()
	parts = append(parts, eventCountBox)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// renderRelatedEvents renders preview of related events.
func (m *BurstSuggestionModelNew) renderRelatedEvents(suggestion burstfact.BurstSuggestion) string {
	theme := m.getTheme()

	if len(suggestion.EventIDs) == 0 {
		return containers.NewBox(theme).
			Variant(containers.BoxWarning).
			Content("No related events").
			Render()
	}

	var events []*career.Event

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
		if i >= 3 {
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
	return containers.NewBox(theme).Content(content).Render()
}

// GetConfirmed returns the list of confirmed suggestions
//
// Returns:
//   - A []burstfact.BurstSuggestion value.
//
// Side effects:
//   - None.
func (m *BurstSuggestionModelNew) GetConfirmed() []burstfact.BurstSuggestion {
	return m.confirmed
}

// GetRejected returns the list of rejected suggestions
//
// Returns:
//   - A []burstfact.BurstSuggestion value.
//
// Side effects:
//   - None.
func (m *BurstSuggestionModelNew) GetRejected() []burstfact.BurstSuggestion {
	return m.rejected
}

// IsDone returns true if all suggestions have been processed
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *BurstSuggestionModelNew) IsDone() bool {
	return len(m.confirmed)+len(m.rejected) == len(m.suggestions)
}

// IsEditing returns true if the model is currently in edit mode.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *BurstSuggestionModelNew) IsEditing() bool {
	return m.editing
}

// SetEditedName sets an edited name for a suggestion at the given index.
// This is primarily for testing purposes.
//
// Expected:
//   - idx must be a valid index in the suggestions slice.
//   - name is the new name to set.
//
// Returns:
//   - None.
//
// Side effects:
//   - Modifies the internal editedNames map.
func (m *BurstSuggestionModelNew) SetEditedName(idx int, name string) {
	m.editedNames[idx] = name
}

// SetEditedDescription sets an edited description for a suggestion at the given index.
// This is primarily for testing purposes.
//
// Expected:
//   - idx must be a valid index in the suggestions slice.
//   - desc is the new description to set.
//
// Returns:
//   - None.
//
// Side effects:
//   - Modifies the internal editedDescs map.
func (m *BurstSuggestionModelNew) SetEditedDescription(idx int, desc string) {
	m.editedDescs[idx] = desc
}

// ExitEditMode forces the model out of edit mode without saving.
// This is primarily for testing purposes.
//
// Returns:
//   - None.
//
// Side effects:
//   - Sets editing to false and clears edit form state.
func (m *BurstSuggestionModelNew) ExitEditMode() {
	m.editing = false
	m.editForm = nil
	m.editFormData = nil
}

// GetTitle returns the modal title for overlay rendering.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *BurstSuggestionModelNew) GetTitle() string {
	if m.editing {
		return "Edit Burst Name & Description"
	}
	return fmt.Sprintf("Burst Suggestion %d of %d", m.currentIdx+1, len(m.suggestions))
}

// GetContent returns the view content without wrapper for overlay rendering.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *BurstSuggestionModelNew) GetContent() string {
	if len(m.suggestions) == 0 {
		return "No burst suggestions available"
	}

	if m.editing {
		if m.editForm == nil {
			return "Edit form not initialized"
		}
		return m.editForm.View()
	}

	// Return review view content
	current := m.suggestions[m.currentIdx]
	var parts []string

	// Progress bar
	progressBar := m.renderProgressBar()
	parts = append(parts, progressBar)

	// Confidence score
	confidenceVis := m.renderConfidenceScore(current.ConfidenceScore)
	parts = append(parts, confidenceVis)

	// Burst details
	burstDetails := m.renderBurstDetails(current)
	parts = append(parts, burstDetails)

	// Related events preview
	relatedEventsView := m.renderRelatedEvents(current)
	parts = append(parts, relatedEventsView)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// GetFooter returns the footer instructions for the modal.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *BurstSuggestionModelNew) GetFooter() string {
	if m.editing {
		return "Tab: Navigate | Enter: Save | Esc: Cancel"
	}
	return "Up/Down: Navigate | y: Confirm | n: Reject | e: Edit | Esc: Back"
}

// getTheme returns the theme or a default theme if none is set.
func (m *BurstSuggestionModelNew) getTheme() themes.Theme {
	return themes.NewDefaultTheme()
}

// createBurstFromSuggestion converts a BurstSuggestion into a Burst domain object.
// Delegates to the pure domain function capture.CreateBurstFromSuggestion.
func (m *BurstSuggestionModelNew) createBurstFromSuggestion(suggestion burstfact.BurstSuggestion) *career.Burst {
	return domcapture.CreateBurstFromSuggestion(domcapture.BurstSuggestionInput{
		Name:        suggestion.Name,
		Description: suggestion.Description,
		EventIDs:    suggestion.EventIDs,
	})
}
