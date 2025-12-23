package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// FormField represents the different fields in the form
type FormField int

const (
	TextField FormField = iota
	DateField
	CompanyField
	ProjectField
	ModeField
	SubmitButton
)

// FormModel represents the event capture form state
type FormModel struct {
	cliService *service.CLIEventService
	inputs     []textinput.Model
	focusIndex int
	modeIndex  int // Index for capture mode selection
	modes      []careerservice.EventCaptureMode
	err        error
	submitted  bool
	event      *career.CareerEvent
	charCount  int
	maxChars   int
}

// NewFormModel creates a new form model with the required fields
func NewFormModel(cliService *service.CLIEventService) *FormModel {
	// Create input fields
	inputs := make([]textinput.Model, 4)

	// Text input (required)
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Enter event description (required, max 2000 chars)"
	inputs[0].Focus()
	inputs[0].CharLimit = 2000
	inputs[0].Width = 60

	// Date input (optional)
	inputs[1] = textinput.New()
	inputs[1].Placeholder = "YYYY-MM-DD or 'today', '1 week ago' (defaults to today)"
	inputs[1].Width = 60

	// Company input (optional)
	inputs[2] = textinput.New()
	inputs[2].Placeholder = "Company name (optional)"
	inputs[2].Width = 60

	// Project input (optional)
	inputs[3] = textinput.New()
	inputs[3].Placeholder = "Project name (optional)"
	inputs[3].Width = 60

	// Capture modes
	modes := []careerservice.EventCaptureMode{
		careerservice.TimelineJournaling,
		careerservice.CVBackfill,
		careerservice.ManualEntry,
	}

	return &FormModel{
		cliService: cliService,
		inputs:     inputs,
		focusIndex: 0,
		modeIndex:  0,
		modes:      modes,
		err:        nil,
		submitted:  false,
		maxChars:   2000,
	}
}

// Init initializes the form model
func (m *FormModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates the form state
func (m *FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SubmitMsg:
		// Handle form submission result
		if msg.Err != nil {
			m.err = msg.Err
			m.submitted = false
			m.event = nil
		} else {
			m.err = nil
			m.submitted = true
			m.event = msg.Event
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Handle mode selection when on mode field
			if m.focusIndex == int(ModeField) {
				if s == "up" {
					m.modeIndex--
					if m.modeIndex < 0 {
						m.modeIndex = len(m.modes) - 1
					}
					return m, nil
				}
				if s == "down" {
					m.modeIndex++
					if m.modeIndex >= len(m.modes) {
						m.modeIndex = 0
					}
					return m, nil
				}
			}

			// Handle navigation
			if s == "enter" && m.focusIndex == int(SubmitButton) {
				return m, m.submitForm()
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > int(SubmitButton) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = int(SubmitButton)
			}

			// Update focus state
			cmds := m.updateFocus()
			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and update the focused input
	cmd := m.updateInputs(msg)
	m.charCount = len(m.inputs[0].Value())

	return m, cmd
}

// View renders the form
func (m *FormModel) View() string {
	var b strings.Builder

	b.WriteString("\n╔═══════════════════════════════════════════════════════════════════════════════╗\n")
	b.WriteString("║                         Capture Career Event                                  ║\n")
	b.WriteString("╠═══════════════════════════════════════════════════════════════════════════════╣\n")
	b.WriteString("║                                                                               ║\n")

	// Text field
	focused := m.focusIndex == int(TextField)
	b.WriteString(fmt.Sprintf("║ Event Text (required): %s%-41s║\n",
		m.getFocusIndicator(focused), ""))
	b.WriteString(fmt.Sprintf("║ %s%-74s║\n", m.inputs[0].View(), ""))
	b.WriteString(fmt.Sprintf("║ Characters: %d/%d %s%-54s║\n",
		m.charCount, m.maxChars, m.getCharCountIndicator(), ""))
	b.WriteString("║                                                                               ║\n")

	// Date field
	focused = m.focusIndex == int(DateField)
	b.WriteString(fmt.Sprintf("║ Date (optional): %s%-55s║\n",
		m.getFocusIndicator(focused), ""))
	b.WriteString(fmt.Sprintf("║ %s%-74s║\n", m.inputs[1].View(), ""))
	b.WriteString("║                                                                               ║\n")

	// Company field
	focused = m.focusIndex == int(CompanyField)
	b.WriteString(fmt.Sprintf("║ Company (optional): %s%-52s║\n",
		m.getFocusIndicator(focused), ""))
	b.WriteString(fmt.Sprintf("║ %s%-74s║\n", m.inputs[2].View(), ""))
	b.WriteString("║                                                                               ║\n")

	// Project field
	focused = m.focusIndex == int(ProjectField)
	b.WriteString(fmt.Sprintf("║ Project (optional): %s%-52s║\n",
		m.getFocusIndicator(focused), ""))
	b.WriteString(fmt.Sprintf("║ %s%-74s║\n", m.inputs[3].View(), ""))
	b.WriteString("║                                                                               ║\n")

	// Mode selector
	focused = m.focusIndex == int(ModeField)
	b.WriteString(fmt.Sprintf("║ Capture Mode: %s%-58s║\n",
		m.getFocusIndicator(focused), ""))
	b.WriteString(m.renderModeSelector())
	b.WriteString("║                                                                               ║\n")

	// Submit button
	focused = m.focusIndex == int(SubmitButton)
	buttonText := "[ Submit ]"
	if focused {
		buttonText = "[ > Submit < ]"
	}
	b.WriteString(fmt.Sprintf("║ %s%-72s║\n", buttonText, ""))
	b.WriteString("║                                                                               ║\n")

	// Error message
	if m.err != nil {
		b.WriteString(fmt.Sprintf("║ Error: %-70s║\n", m.err.Error()))
		b.WriteString("║                                                                               ║\n")
	}

	// Help text
	b.WriteString("║ Navigation: Tab/Shift+Tab to move, Up/Down for mode, Enter to submit, Esc to cancel ║\n")
	b.WriteString("╚═══════════════════════════════════════════════════════════════════════════════╝\n")

	return b.String()
}

// renderModeSelector renders the capture mode selector
func (m *FormModel) renderModeSelector() string {
	var b strings.Builder

	for i, mode := range m.modes {
		selected := ""
		if i == m.modeIndex {
			selected = "►"
		}
		modeName := m.getModeName(mode)
		modeDesc := m.getModeDescription(mode)
		b.WriteString(fmt.Sprintf("║   %s %-20s - %-49s║\n", selected, modeName, modeDesc))
	}

	return b.String()
}

// getModeName returns a human-readable name for a capture mode
func (m *FormModel) getModeName(mode careerservice.EventCaptureMode) string {
	switch mode {
	case careerservice.TimelineJournaling:
		return "Timeline Journaling"
	case careerservice.CVBackfill:
		return "CV Backfill"
	case careerservice.ManualEntry:
		return "Manual Entry"
	default:
		return string(mode)
	}
}

// getModeDescription returns a description for a capture mode
func (m *FormModel) getModeDescription(mode careerservice.EventCaptureMode) string {
	switch mode {
	case careerservice.TimelineJournaling:
		return "Real-time logging (last 30 days)"
	case careerservice.CVBackfill:
		return "Import from existing CV (any date)"
	case careerservice.ManualEntry:
		return "Manual entry (any date)"
	default:
		return ""
	}
}

// getFocusIndicator returns a visual indicator for focused fields
func (m *FormModel) getFocusIndicator(focused bool) string {
	if focused {
		return "►"
	}
	return " "
}

// getCharCountIndicator returns a visual indicator for character count
func (m *FormModel) getCharCountIndicator() string {
	if m.charCount > m.maxChars-100 {
		return "⚠"
	}
	return ""
}

// updateFocus updates the focus state of all inputs
func (m *FormModel) updateFocus() []tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := 0; i < len(m.inputs); i++ {
		if i == m.focusIndex {
			cmds[i] = m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}

	return cmds
}

// updateInputs handles input for the currently focused field
func (m *FormModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only update the focused input
	if m.focusIndex < len(m.inputs) {
		m.inputs[m.focusIndex], cmds[m.focusIndex] = m.inputs[m.focusIndex].Update(msg)
	}

	return tea.Batch(cmds...)
}

// SubmitMsg is sent when form submission completes
type SubmitMsg struct {
	Event *career.CareerEvent
	Err   error
}

// submitForm handles form submission
func (m *FormModel) submitForm() tea.Cmd {
	return func() tea.Msg {
		// Validate and parse inputs
		text := strings.TrimSpace(m.inputs[0].Value())
		if text == "" {
			return SubmitMsg{Err: fmt.Errorf("event text is required")}
		}

		// Validate text length
		if len(text) > 2000 {
			return SubmitMsg{Err: fmt.Errorf("event text cannot exceed 2000 characters")}
		}

		// Parse date
		dateStr := strings.TrimSpace(m.inputs[1].Value())
		eventDate, err := m.parseDate(dateStr)
		if err != nil {
			return SubmitMsg{Err: fmt.Errorf("invalid date: %w", err)}
		}

		// Validate date is not in future
		if eventDate.After(time.Now()) {
			return SubmitMsg{Err: fmt.Errorf("date cannot be in the future")}
		}

		// Validate date based on capture mode
		mode := m.modes[m.modeIndex]
		if mode == careerservice.TimelineJournaling {
			thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
			if eventDate.Before(thirtyDaysAgo) {
				return SubmitMsg{Err: fmt.Errorf("timeline journaling events must be within the last 30 days")}
			}
		}

		// Get optional fields
		company := strings.TrimSpace(m.inputs[2].Value())
		project := strings.TrimSpace(m.inputs[3].Value())

		// Capture event through service
		ctx := context.Background()
		opts := []service.Option{}
		if company != "" {
			opts = append(opts, service.WithCompany(company))
		}
		if project != "" {
			opts = append(opts, service.WithProject(project))
		}

		if err := m.cliService.CaptureEvent(ctx, text, eventDate, mode, opts...); err != nil {
			return SubmitMsg{Err: fmt.Errorf("failed to capture event: %w", err)}
		}

		// Build event for display (this is just for UI purposes)
		event := &career.CareerEvent{
			Text:    text,
			Date:    eventDate,
			Company: company,
			Project: project,
			Tags:    []string{}, // Will be handled in tag selector
		}

		return SubmitMsg{Event: event, Err: nil}
	}
}

// parseDate parses a date string into a time.Time
func (m *FormModel) parseDate(dateStr string) (time.Time, error) {
	if dateStr == "" || strings.ToLower(dateStr) == "today" {
		return time.Now(), nil
	}

	// Try parsing as YYYY-MM-DD
	t, err := time.Parse("2006-01-02", dateStr)
	if err == nil {
		return t, nil
	}

	// Try parsing relative dates like "1 week ago", "2 days ago"
	// This is a simple implementation - could be expanded
	if strings.HasSuffix(dateStr, " ago") {
		parts := strings.Fields(dateStr)
		if len(parts) >= 3 {
			// Simple parsing for "N unit ago"
			amount := 0
			fmt.Sscanf(parts[0], "%d", &amount)
			unit := parts[1]

			switch unit {
			case "day", "days":
				return time.Now().AddDate(0, 0, -amount), nil
			case "week", "weeks":
				return time.Now().AddDate(0, 0, -amount*7), nil
			case "month", "months":
				return time.Now().AddDate(0, -amount, 0), nil
			case "year", "years":
				return time.Now().AddDate(-amount, 0, 0), nil
			}
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

// Submitted returns whether the form has been successfully submitted
func (m *FormModel) Submitted() bool {
	return m.submitted
}

// Event returns the captured event
func (m *FormModel) Event() *career.CareerEvent {
	return m.event
}

// Error returns the last error
func (m *FormModel) Error() error {
	return m.err
}

// GetInputValue returns the value of an input field (for testing)
func (m *FormModel) GetInputValue(index int) string {
	if index >= 0 && index < len(m.inputs) {
		return m.inputs[index].Value()
	}
	return ""
}

// Reset resets the form to its initial state
func (m *FormModel) Reset() {
	m.focusIndex = 0
	m.modeIndex = 0
	m.submitted = false
	m.event = nil
	m.err = nil
	m.charCount = 0

	for i := range m.inputs {
		m.inputs[i].SetValue("")
	}

	m.inputs[0].Focus()
}
