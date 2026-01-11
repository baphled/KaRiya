package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// HuhFormModel is a huh-based implementation of the capture event form.
// It provides the same external interface as FormModel but uses huh for form handling.
type HuhFormModel struct {
	*BaseStandardModel
	cliService  *service.CLIEventService
	form        *huh.Form
	formData    *forms.CaptureEventFormData
	strategy    string
	width       int
	height      int
	submitted   bool
	event       *career.CareerEvent
	err         error
	editMode    bool
	editEventID string
}

// NewHuhFormModel creates a new huh-based form model.
func NewHuhFormModel(cliService *service.CLIEventService) *HuhFormModel {
	formData := forms.NewCaptureEventFormData()

	m := &HuhFormModel{
		BaseStandardModel: NewBaseStandardModel(),
		cliService:        cliService,
		formData:          formData,
		strategy:          "manual",
		width:             80,
		height:            24,
	}

	// Create initial form
	m.rebuildForm()

	return m
}

// rebuildForm creates a new form with current settings.
func (m *HuhFormModel) rebuildForm() {
	m.form = forms.NewCaptureEventForm(
		m.formData,
		m.strategy,
		m.width-4, // Leave margin
		forms.DefaultFormHeight(m.height),
	)
}

// Init initializes the form model.
func (m *HuhFormModel) Init() tea.Cmd {
	return m.form.Init()
}

// Update handles messages and updates the form state.
func (m *HuhFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Update form dimensions without losing state
		m.form = m.form.
			WithHeight(forms.DefaultFormHeight(m.height)).
			WithWidth(m.width - 4)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return BackMsg{} }
		case "ctrl+c", "q":
			return m, func() tea.Msg { return QuitMsg{} }
		}
	}

	// Update the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	// Check form state
	if forms.IsCompleted(m.form) {
		return m, m.submitForm()
	}

	if forms.IsAborted(m.form) {
		return m, func() tea.Msg { return BackMsg{} }
	}

	return m, cmd
}

// View renders the form.
func (m *HuhFormModel) View() string {
	return m.form.View()
}

// submitForm handles form submission.
func (m *HuhFormModel) submitForm() tea.Cmd {
	return func() tea.Msg {
		// Check if user cancelled via confirm button
		if !m.formData.SubmitConfirmed {
			return BackMsg{}
		}

		// Validate text
		text := strings.TrimSpace(m.formData.Text)
		if text == "" {
			return SubmitMsg{Err: fmt.Errorf("event text is required")}
		}
		if len(text) < 10 {
			return SubmitMsg{Err: fmt.Errorf("event text must be at least 10 characters")}
		}
		if len(text) > 2000 {
			return SubmitMsg{Err: fmt.Errorf("event text cannot exceed 2000 characters")}
		}

		// Parse date
		eventDate, err := parseDate(m.formData.Date)
		if err != nil {
			return SubmitMsg{Err: fmt.Errorf("invalid date: %w", err)}
		}

		// Validate date is not in future
		if eventDate.After(time.Now()) {
			return SubmitMsg{Err: fmt.Errorf("date cannot be in the future")}
		}

		// Build event object
		event := &career.CareerEvent{
			ID:         m.editEventID,
			Text:       text,
			Date:       eventDate,
			Company:    strings.TrimSpace(m.formData.Company),
			Project:    strings.TrimSpace(m.formData.Project),
			Tags:       m.formData.Tags,
			Categories: m.formData.Categories,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		m.submitted = true
		m.event = event

		return SubmitMsg{Event: event, Err: nil}
	}
}

// SubmitForm triggers form submission (for Ctrl+S shortcut).
func (m *HuhFormModel) SubmitForm() tea.Cmd {
	// Mark submit as confirmed and trigger submission
	m.formData.SubmitConfirmed = true
	return m.submitForm()
}

// SetStrategy sets the capture strategy and rebuilds the form.
func (m *HuhFormModel) SetStrategy(strategy string) {
	m.strategy = strategy
	m.rebuildForm()
}

// GetStrategy returns the current capture strategy.
func (m *HuhFormModel) GetStrategy() string {
	return m.strategy
}

// Reset resets the form to its initial state.
func (m *HuhFormModel) Reset() {
	m.formData = forms.NewCaptureEventFormData()
	m.submitted = false
	m.event = nil
	m.err = nil
	m.editMode = false
	m.editEventID = ""
	m.rebuildForm()
}

// LoadEventForEditing populates the form with an existing event's data.
func (m *HuhFormModel) LoadEventForEditing(event *career.CareerEvent) {
	m.editMode = true
	m.editEventID = event.ID
	m.event = event

	m.formData.Text = event.Text
	m.formData.Date = event.Date.Format("2006-01-02")
	m.formData.Company = event.Company
	m.formData.Project = event.Project
	m.formData.Tags = event.Tags
	m.formData.Categories = event.Categories

	m.rebuildForm()
}

// IsEditMode returns whether the form is in edit mode.
func (m *HuhFormModel) IsEditMode() bool {
	return m.editMode
}

// GetEditEventID returns the ID of the event being edited.
func (m *HuhFormModel) GetEditEventID() string {
	return m.editEventID
}

// IsSubmitted returns whether the form has been submitted.
func (m *HuhFormModel) IsSubmitted() bool {
	return m.submitted
}

// GetEvent returns the submitted event.
func (m *HuhFormModel) GetEvent() *career.CareerEvent {
	return m.event
}

// GetError returns any error that occurred during submission.
func (m *HuhFormModel) GetError() error {
	return m.err
}

// ShowOptionalFields returns whether optional fields are shown (always true for manual, false for quick).
func (m *HuhFormModel) ShowOptionalFields() bool {
	return m.strategy == "manual"
}

// ToggleOptionalFields is a no-op for huh forms (visibility is controlled by strategy).
func (m *HuhFormModel) ToggleOptionalFields() {
	// No-op: huh forms handle visibility via WithHide
}

// Helper function to parse date strings.
func parseDate(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" || strings.ToLower(dateStr) == "today" {
		return time.Now(), nil
	}

	// Try parsing as YYYY-MM-DD
	t, err := time.Parse("2006-01-02", dateStr)
	if err == nil {
		return t, nil
	}

	// Try parsing relative dates like "1 week ago", "2 days ago"
	if strings.HasSuffix(dateStr, " ago") {
		return parseRelativeDate(dateStr)
	}

	return time.Time{}, fmt.Errorf("invalid date format: %s", dateStr)
}

// parseRelativeDate parses dates like "1 week ago", "2 days ago".
func parseRelativeDate(dateStr string) (time.Time, error) {
	parts := strings.Fields(dateStr)
	if len(parts) < 3 {
		return time.Time{}, fmt.Errorf("invalid relative date: %s", dateStr)
	}

	var num int
	_, err := fmt.Sscanf(parts[0], "%d", &num)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid number in relative date: %s", dateStr)
	}

	unit := strings.ToLower(parts[1])
	now := time.Now()

	switch {
	case strings.HasPrefix(unit, "day"):
		return now.AddDate(0, 0, -num), nil
	case strings.HasPrefix(unit, "week"):
		return now.AddDate(0, 0, -num*7), nil
	case strings.HasPrefix(unit, "month"):
		return now.AddDate(0, -num, 0), nil
	case strings.HasPrefix(unit, "year"):
		return now.AddDate(-num, 0, 0), nil
	default:
		return time.Time{}, fmt.Errorf("unknown time unit: %s", unit)
	}
}
