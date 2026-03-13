package event

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
)

// ReviewEnrichmentConfig holds configuration for the review enrichment view.
type ReviewEnrichmentConfig struct {
	AvailableTags       []string
	AvailableCategories []string
	AvailableSkills     []string
	TerminalWidth       int
	TerminalHeight      int
}

// ReviewEnrichmentResult holds the data returned when enrichment is submitted.
type ReviewEnrichmentResult struct {
	Event    display.Event
	FormData *forms.MetadataFormData
}

// ReviewEnrichment is a view for editing event metadata during review.
type ReviewEnrichment struct {
	widgets.BaseView
	event         display.Event
	originalEvent display.Event
	form          forms.Form
	formData      *forms.MetadataFormData
	submitted     bool
	cancelled     bool
	err           error
}

// NewReviewEnrichment creates a ReviewEnrichment view for the given event.
//
// Expected:
//   - event must be valid.
//   - config must be a valid configuration object.
//
// Returns:
//   - A fully initialized ReviewEnrichment ready for use.
//
// Side effects:
//   - None.
func NewReviewEnrichment(evt display.Event, config ReviewEnrichmentConfig) *ReviewEnrichment {
	width := config.TerminalWidth
	height := config.TerminalHeight
	if width == 0 {
		width = 80
	}
	if height == 0 {
		height = 40
	}
	eventCopy := evt
	formData := metadataFormDataFromDisplayEvent(evt)
	form := forms.NewMetadataForm(formData, forms.MetadataFormConfig{
		AvailableTags:       config.AvailableTags,
		AvailableCategories: config.AvailableCategories,
		AvailableSkills:     nil,
		Width:               width,
		Height:              height,
	})
	re := &ReviewEnrichment{
		event:         evt,
		originalEvent: eventCopy,
		form:          form,
		formData:      formData,
	}
	re.SetTerminalInfo(width, height)
	return re
}

// Init returns the initial command for the form.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *ReviewEnrichment) Init() tea.Cmd {
	return m.form.Init()
}

// Update handles messages and returns a ViewResult on user action.
//
// Expected:
//   - msg is a valid tea.Msg (key press, window resize, or custom message).
//
// Returns:
//   - A tea.Cmd and widgets.ViewResult pair.
//
// Side effects:
//   - Updates internal view state based on message type.
func (m *ReviewEnrichment) Update(msg tea.Msg) (tea.Cmd, widgets.ViewResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetTerminalInfo(msg.Width, msg.Height)
		return nil, nil
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			m.cancelled = true
			return nil, &widgets.CancelViewResult{}
		}
		if msg.Type == tea.KeyCtrlS {
			if err := applyMetadataFormDataToDisplayEvent(&m.event, m.formData); err != nil {
				m.err = err
				return nil, nil
			}
			m.event.UpdatedAt = time.Now()
			m.submitted = true
			return nil, &widgets.SubmitViewResult{FormData: ReviewEnrichmentResult{Event: m.event, FormData: m.formData}}
		}
	}
	var cmd tea.Cmd
	m.form, cmd = forms.Update(m.form, msg)
	if forms.IsCompleted(m.form) {
		if err := applyMetadataFormDataToDisplayEvent(&m.event, m.formData); err != nil {
			m.err = err
			return nil, nil
		}
		m.event.UpdatedAt = time.Now()
		m.submitted = true
		return nil, &widgets.SubmitViewResult{FormData: ReviewEnrichmentResult{Event: m.event, FormData: m.formData}}
	}
	if forms.IsAborted(m.form) {
		m.cancelled = true
		return nil, &widgets.CancelViewResult{}
	}
	return cmd, nil
}

// RenderContent returns the rendered form content.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *ReviewEnrichment) RenderContent() string {
	return m.form.View()
}

// HelpText returns key binding instructions.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *ReviewEnrichment) HelpText() string {
	return "Enter: Confirm | Esc: Cancel | Tab: Next Field | Shift+Tab: Previous"
}

// GetEvent returns the current event being edited.
//
// Returns:
//   - A display.Event value.
//
// Side effects:
//   - None.
func (m *ReviewEnrichment) GetEvent() display.Event {
	return m.event
}

// IsSubmitted returns whether the form was submitted.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *ReviewEnrichment) IsSubmitted() bool {
	return m.submitted
}

// IsCancelled returns whether the form was cancelled.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *ReviewEnrichment) IsCancelled() bool {
	return m.cancelled
}

// Revert restores the event to its original state before editing.
//
// Side effects:
//   - None.
func (m *ReviewEnrichment) Revert() {
	m.event = m.originalEvent
	m.formData = metadataFormDataFromDisplayEvent(m.originalEvent)
	m.form = forms.NewMetadataForm(m.formData, forms.MetadataFormConfig{
		AvailableTags:       nil,
		AvailableCategories: nil,
		AvailableSkills:     nil,
		Width:               m.GetTerminalWidth(),
		Height:              m.GetTerminalHeight(),
	})
}

// GetError returns any validation error from the last submit attempt.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (m *ReviewEnrichment) GetError() error {
	return m.err
}

// GetTitle returns the view title.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *ReviewEnrichment) GetTitle() string {
	return "Edit Event Metadata"
}

// GetContent returns the rendered form content.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *ReviewEnrichment) GetContent() string {
	return m.form.View()
}

// GetFooter returns the help text footer.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *ReviewEnrichment) GetFooter() string {
	return m.HelpText()
}

func metadataFormDataFromDisplayEvent(evt display.Event) *forms.MetadataFormData {
	date := ""
	if !evt.Date.IsZero() {
		date = evt.Date.Format("2006-01-02")
	}

	return &forms.MetadataFormData{
		Date:       date,
		Company:    evt.Company,
		Project:    evt.Project,
		Tags:       append([]string(nil), evt.Tags...),
		Categories: append([]string(nil), evt.Categories...),
		Skills:     append([]string(nil), evt.Skills...),
	}
}

func applyMetadataFormDataToDisplayEvent(evt *display.Event, data *forms.MetadataFormData) error {
	if evt == nil || data == nil {
		return nil
	}

	if data.Date != "" {
		parsedDate, err := forms.ParseDateString(data.Date)
		if err != nil {
			return err
		}
		evt.Date = parsedDate
	}

	evt.Company = data.Company
	evt.Project = data.Project
	evt.Tags = append([]string(nil), data.Tags...)
	evt.Categories = append([]string(nil), data.Categories...)
	evt.Skills = append([]string(nil), data.Skills...)

	return nil
}
