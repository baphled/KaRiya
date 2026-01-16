package models

import (
	"time"

	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// CaptureForm wraps a huh form for capturing career events.
type CaptureForm struct {
	*BaseStandardModel
	cliService *service.CLIEventService
	formData   *forms.CaptureEventFormData
	form       *huh.Form
	strategy   string
	width      int
	height     int
}

// NewCaptureForm creates a new capture form.
func NewCaptureForm(cliService *service.CLIEventService) *CaptureForm {
	formData := forms.NewCaptureEventFormData()

	m := &CaptureForm{
		BaseStandardModel: NewBaseStandardModel(),
		cliService:        cliService,
		formData:          formData,
		strategy:          "manual",
		width:             80,
		height:            24,
	}

	m.rebuildForm()
	return m
}

// rebuildForm creates a new form with current settings.
func (m *CaptureForm) rebuildForm() {
	m.form = forms.NewCaptureEventForm(
		m.formData,
		m.strategy,
		m.width-4,
		forms.DefaultFormHeight(m.height),
	)
}

// Init initializes the form.
func (m *CaptureForm) Init() tea.Cmd {
	return m.form.Init()
}

// Update handles messages.
func (m *CaptureForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.form = m.form.WithHeight(forms.DefaultFormHeight(m.height)).WithWidth(m.width - 4)
		return m, nil

	case tea.KeyMsg:
		// Handle escape BEFORE delegating to form
		// This allows the parent intent to handle back navigation
		if msg.String() == "esc" {
			// Signal back navigation to parent intent
			return m, nil // Parent intent will check for escape via HandleGlobalKeys
		}

		if msg.String() == "ctrl+s" {
			m.formData.SubmitConfirmed = true
			return m, m.submitForm()
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	// When form is completed (user pressed Submit button), trigger submission
	if m.form.State == huh.StateCompleted {
		return m, m.submitForm()
	}

	return m, cmd
}

// View renders the form.
func (m *CaptureForm) View() string {
	return m.form.View()
}

// submitForm creates a SubmitMsg.
func (m *CaptureForm) submitForm() tea.Cmd {
	return func() tea.Msg {
		var eventDate time.Time
		var err error

		if m.formData.Date == "" {
			eventDate = time.Now()
		} else {
			eventDate, err = forms.ParseDateString(m.formData.Date)
			if err != nil {
				return SubmitMsg{Event: nil, Err: err}
			}
		}

		event := &career.CareerEvent{
			Text:       m.formData.Text,
			Date:       eventDate,
			Company:    m.formData.Company,
			Project:    m.formData.Project,
			Tags:       m.formData.Tags,
			Categories: m.formData.Categories,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		return SubmitMsg{Event: event, Err: nil}
	}
}

// SetStrategy updates the strategy.
func (m *CaptureForm) SetStrategy(strategy string) {
	m.strategy = strategy
	m.rebuildForm()
}

// LoadEventForEditing populates the form with existing event data.
func (m *CaptureForm) LoadEventForEditing(event *career.CareerEvent) {
	if event == nil {
		return
	}

	m.formData.Text = event.Text
	m.formData.Date = event.Date.Format("2006-01-02")
	m.formData.Company = event.Company
	m.formData.Project = event.Project

	if event.Tags != nil {
		m.formData.Tags = make([]string, len(event.Tags))
		copy(m.formData.Tags, event.Tags)
	}
	if event.Categories != nil {
		m.formData.Categories = make([]string, len(event.Categories))
		copy(m.formData.Categories, event.Categories)
	}

	m.rebuildForm()
}

// GetStrategy returns the current strategy.
func (m *CaptureForm) GetStrategy() string {
	return m.strategy
}

// SubmitForm triggers form submission (for compatibility).
func (m *CaptureForm) SubmitForm() tea.Cmd {
	m.formData.SubmitConfirmed = true
	return m.submitForm()
}
