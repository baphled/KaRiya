package models

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burst_fact"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmBurstMsg is sent when user confirms a burst suggestion
type ConfirmBurstMsg struct {
	Burst *career.Burst
}

// BurstProcessingCompleteMsg is sent when burst suggestion workflow is done
type BurstProcessingCompleteMsg struct {
	ConfirmedCount int
	RejectedCount  int
}

// BurstSuggestionEditField indicates which field is being edited
type BurstSuggestionEditField int

const (
	BurstSuggestionNameField BurstSuggestionEditField = iota
	BurstSuggestionDescField
)

// BurstSuggestionModel represents the burst suggestion review and confirmation screen
type BurstSuggestionModel struct {
	*BaseStandardModel
	service       *careerservice.Service
	ctx           context.Context
	suggestions   []burstfact.BurstSuggestion
	currentIdx    int
	confirmed     []burstfact.BurstSuggestion
	rejected      []burstfact.BurstSuggestion
	editing       bool
	editField     BurstSuggestionEditField
	editedNames   map[int]string    // Map of suggestion index to edited name
	editedDescs   map[int]string    // Map of suggestion index to edited description
	inputs        []textinput.Model // For editing name and description
	focusIndex    int
	relatedEvents map[int][]*career.CareerEvent // Cache of related events
	width         int
	height        int
	helpFooter    components.HelpFooterModel // Help footer
}

// NewBurstSuggestionModel creates a new burst suggestion model
func NewBurstSuggestionModel(svc *careerservice.Service, suggestions []burstfact.BurstSuggestion, ctx context.Context) *BurstSuggestionModel {
	// Create input fields for editing name and description
	nameInput := textinput.New()
	nameInput.Placeholder = "Burst name (optional)"
	nameInput.Width = 60
	nameInput.Focus()

	descInput := textinput.New()
	descInput.Placeholder = "Burst description (optional)"
	descInput.Width = 60

	return &BurstSuggestionModel{
		BaseStandardModel: NewBaseStandardModel(),
		service:           svc,
		ctx:               ctx,
		suggestions:       suggestions,
		currentIdx:        0,
		confirmed:         []burstfact.BurstSuggestion{},
		rejected:          []burstfact.BurstSuggestion{},
		editing:           false,
		editField:         BurstSuggestionNameField,
		editedNames:       make(map[int]string),
		editedDescs:       make(map[int]string),
		inputs:            []textinput.Model{nameInput, descInput},
		focusIndex:        0,
		relatedEvents:     make(map[int][]*career.CareerEvent),
		width:             80,
		height:            24,
	}
}

// Init initializes the model
func (m *BurstSuggestionModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *BurstSuggestionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
func (m *BurstSuggestionModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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

// handleEditKeyMsg handles keyboard input while editing name/description
func (m *BurstSuggestionModel) handleEditKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyTab:
		m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
		m.updateEditInputFocus()
		return m, nil

	case tea.KeyShiftTab:
		m.focusIndex = (m.focusIndex - 1 + len(m.inputs)) % len(m.inputs)
		m.updateEditInputFocus()
		return m, nil

	case tea.KeyEnter:
		// Save edits and return to review
		return m.saveEdits()

	case tea.KeyEsc:
		// Cancel editing
		m.editing = false
		m.focusIndex = 0
		return m, nil
	}

	// Update the focused input
	var cmd tea.Cmd
	m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
	return m, cmd
}

// startEdit begins editing the current suggestion's name/description
func (m *BurstSuggestionModel) startEdit() (tea.Model, tea.Cmd) {
	if len(m.suggestions) == 0 {
		return m, nil
	}

	m.editing = true
	m.focusIndex = 0
	current := m.suggestions[m.currentIdx]

	// Set input values
	if name, exists := m.editedNames[m.currentIdx]; exists {
		m.inputs[0].SetValue(name)
	} else {
		m.inputs[0].SetValue(current.Name)
	}

	if desc, exists := m.editedDescs[m.currentIdx]; exists {
		m.inputs[1].SetValue(desc)
	} else {
		m.inputs[1].SetValue(current.Description)
	}

	m.updateEditInputFocus()
	return m, nil
}

// saveEdits saves the edited name/description
func (m *BurstSuggestionModel) saveEdits() (tea.Model, tea.Cmd) {
	m.editedNames[m.currentIdx] = m.inputs[0].Value()
	m.editedDescs[m.currentIdx] = m.inputs[1].Value()

	m.editing = false
	m.focusIndex = 0
	return m, nil
}

// confirmCurrent confirms the current suggestion and moves to next
func (m *BurstSuggestionModel) confirmCurrent() (tea.Model, tea.Cmd) {
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
func (m *BurstSuggestionModel) rejectCurrent() (tea.Model, tea.Cmd) {
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

// updateEditInputFocus updates focus state for input fields
func (m *BurstSuggestionModel) updateEditInputFocus() {
	for i := 0; i < len(m.inputs); i++ {
		if i == m.focusIndex {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
}

// View renders the burst suggestion screen
func (m *BurstSuggestionModel) View() string {
	if len(m.suggestions) == 0 {
		return styles.ErrorBox.Render("No burst suggestions available")
	}

	if m.editing {
		return m.renderEditView()
	}

	return m.renderReviewView()
}

// renderReviewView renders the suggestion review view
func (m *BurstSuggestionModel) renderReviewView() string {
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

// renderEditView renders the editing view
func (m *BurstSuggestionModel) renderEditView() string {
	var parts []string

	header := styles.HeaderMain.Render("Edit Burst Name & Description")
	parts = append(parts, header)

	// Name input
	nameLbl := styles.InputLabel.Render("Name:")
	nameInput := m.inputs[0].View()
	nameField := lipgloss.JoinHorizontal(lipgloss.Top, nameLbl, "  ", nameInput)
	parts = append(parts, nameField)

	// Description input
	descLbl := styles.InputLabel.Render("Description:")
	descInput := m.inputs[1].View()
	descField := lipgloss.JoinHorizontal(lipgloss.Top, descLbl, "  ", descInput)
	parts = append(parts, descField)

	// Help text
	helpText := lipgloss.NewStyle().Render("Tab navigate • Enter save • Esc cancel")
	parts = append(parts, helpText)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// renderProgressBar renders the progress through suggestions
func (m *BurstSuggestionModel) renderProgressBar() string {
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
func (m *BurstSuggestionModel) renderConfidenceScore(score float64) string {
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
func (m *BurstSuggestionModel) renderBurstDetails(suggestion burstfact.BurstSuggestion) string {
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
func (m *BurstSuggestionModel) renderRelatedEvents(suggestion burstfact.BurstSuggestion) string {
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
func (m *BurstSuggestionModel) GetConfirmed() []burstfact.BurstSuggestion {
	return m.confirmed
}

// GetRejected returns the list of rejected suggestions
func (m *BurstSuggestionModel) GetRejected() []burstfact.BurstSuggestion {
	return m.rejected
}

// IsDone returns true if all suggestions have been processed
func (m *BurstSuggestionModel) IsDone() bool {
	return len(m.confirmed)+len(m.rejected) == len(m.suggestions)
}

// createBurstFromSuggestion converts a BurstSuggestion into a Burst domain object
func (m *BurstSuggestionModel) createBurstFromSuggestion(suggestion burstfact.BurstSuggestion) *career.Burst {
	// Generate a name if none provided
	name := suggestion.Name
	if name == "" {
		name = fmt.Sprintf("Burst of %d events", len(suggestion.EventIDs))
	}

	// Infer competency focus from events
	competencyFocus := m.service.InferCompetencyForBurst(m.ctx, suggestion.EventIDs)

	return &career.Burst{
		Name:            name,
		Description:     suggestion.Description,
		EventIDs:        suggestion.EventIDs,
		CompetencyFocus: competencyFocus,
	}
}
