package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FormField represents the different fields in the form
type FormField int

const (
	TextField FormField = iota
	DateField
	CompanyField
	ProjectField
	TagsField
	CategoriesField
	ModeField
	SubmitButton
)

// FormModel represents the event capture form state
type FormModel struct {
	*BaseStandardModel
	cliService       *service.CLIEventService
	inputs           []textinput.Model
	focusIndex       int
	modeIndex        int // Index for capture mode selection
	tagIndex         int // Index for tag navigation when TagsField is focused
	categoryIndex    int // Index for category navigation when CategoriesField is focused
	modes            []careerservice.EventCaptureMode
	err              error
	submitted        bool
	event            *career.CareerEvent
	charCount        int
	maxChars         int
	tagSelector      *components.TagSelector
	categorySelector *components.CategorySelector
	fieldErrors      map[FormField]string       // Track field-level validation errors
	editMode         bool                       // True if editing an existing event
	editEventID      string                     // ID of event being edited
	helpFooter       components.HelpFooterModel // Help footer for keyboard shortcuts
	header           components.HeaderModel     // Header component
	footer           components.FooterModel     // Footer component
	breadcrumbs      []string                   // Navigation breadcrumb trail
}

// NewFormModel creates a new form model with the required fields
func NewFormModel(cliService *service.CLIEventService) *FormModel {
	// Create input fields
	inputs := make([]textinput.Model, 4)

	// Text input (required)
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Enter event description (required, max 2000 chars)"
	inputs[0].Focus()
	// inputs[0].CharLimit = 2000 // Allow detection of exceeding limit
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
		BaseStandardModel: NewBaseStandardModel(),
		cliService:        cliService,
		inputs:            inputs,
		focusIndex:        0,
		modeIndex:         0,
		tagIndex:          0,
		categoryIndex:     0,
		modes:             modes,
		err:               nil,
		submitted:         false,
		maxChars:          2000,
		tagSelector:       components.NewTagSelector(),
		categorySelector:  components.NewCategorySelector(),
		fieldErrors:       make(map[FormField]string),
		helpFooter:        components.NewHelpFooter("form", 80),
		header:            components.NewHeader("Capture Career Event", 80),
		footer:            components.NewFooter(80),
	}
}

// Init initializes the form model
func (m *FormModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates the form state
func (m *FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.header.SetWidth(msg.Width)
		m.footer.SetWidth(msg.Width)
		m.helpFooter.SetWidth(msg.Width)

		return m, nil
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
		case "esc":
			// Signal back navigation to parent
			return m, func() tea.Msg { return BackMsg{} }
		case "ctrl+c", "q":
			// Signal quit to parent
			return m, func() tea.Msg { return QuitMsg{} }

		case " ":
			// Handle space key for tag/category selection
			if m.focusIndex == int(TagsField) {
				availableTags := m.tagSelector.AvailableTags()
				if m.tagIndex < len(availableTags) {
					tag := availableTags[m.tagIndex]
					if m.tagSelector.IsSelected(tag) {
						m.tagSelector.DeselectTag(tag)
					} else {
						m.tagSelector.SelectTag(tag)
					}
				}
				return m, nil
			}
			if m.focusIndex == int(CategoriesField) {
				availableCategories := m.categorySelector.AvailableCategories()
				if m.categoryIndex < len(availableCategories) {
					category := availableCategories[m.categoryIndex]
					if m.categorySelector.IsSelected(category) {
						m.categorySelector.DeselectCategory(category)
					} else {
						m.categorySelector.SelectCategory(category)
					}
				}
				return m, nil
			}

		case "j", "k":
			// Handle j/k navigation for tags, categories, and modes
			// Only when focused on these fields, not in text input
			if m.focusIndex == int(TagsField) {
				availableTags := m.tagSelector.AvailableTags()
				if msg.String() == "j" {
					m.tagIndex++
					if m.tagIndex >= len(availableTags) {
						m.tagIndex = 0
					}
				} else if msg.String() == "k" {
					m.tagIndex--
					if m.tagIndex < 0 {
						m.tagIndex = len(availableTags) - 1
					}
				}
				return m, nil
			}

			if m.focusIndex == int(CategoriesField) {
				availableCategories := m.categorySelector.AvailableCategories()
				if msg.String() == "j" {
					m.categoryIndex++
					if m.categoryIndex >= len(availableCategories) {
						m.categoryIndex = 0
					}
				} else if msg.String() == "k" {
					m.categoryIndex--
					if m.categoryIndex < 0 {
						m.categoryIndex = len(availableCategories) - 1
					}
				}
				return m, nil
			}

			if m.focusIndex == int(ModeField) {
				if msg.String() == "j" {
					m.modeIndex++
					if m.modeIndex >= len(m.modes) {
						m.modeIndex = 0
					}
				} else if msg.String() == "k" {
					m.modeIndex--
					if m.modeIndex < 0 {
						m.modeIndex = len(m.modes) - 1
					}
				}
				return m, nil
			}
			// If not in a navigation field, fall through to text input

		case "T", "t":
			// Jump to tags field for selection (only if KeyType is not KeyRunes)
			if msg.Type != tea.KeyRunes {
				// Only if not typing in text field
				if m.focusIndex != int(TextField) && m.focusIndex != int(CompanyField) && m.focusIndex != int(ProjectField) {
					m.focusIndex = int(TagsField)
					m.tagIndex = 0
					return m, nil
				}
			}

		case "G", "g":
			// Jump to categories field for selection (only if KeyType is not KeyRunes)
			if msg.Type != tea.KeyRunes {
				// Only if not typing in text field
				if m.focusIndex != int(TextField) && m.focusIndex != int(CompanyField) && m.focusIndex != int(ProjectField) {
					m.focusIndex = int(CategoriesField)
					m.categoryIndex = 0
					return m, nil
				}
			}

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Validate current field before moving away
			if s == "tab" || s == "shift+tab" || s == "enter" {
				m.validateCurrentField()
			}

			// Handle tag selection when on tags field
			if m.focusIndex == int(TagsField) {
				availableTags := m.tagSelector.AvailableTags()
				if s == "up" {
					m.tagIndex--
					if m.tagIndex < 0 {
						m.tagIndex = len(availableTags) - 1
					}
					return m, nil
				}
				if s == "down" {
					m.tagIndex++
					if m.tagIndex >= len(availableTags) {
						m.tagIndex = 0
					}
					return m, nil
				}
			}

			// Handle category selection when on categories field
			if m.focusIndex == int(CategoriesField) {
				availableCategories := m.categorySelector.AvailableCategories()
				if s == "up" {
					m.categoryIndex--
					if m.categoryIndex < 0 {
						m.categoryIndex = len(availableCategories) - 1
					}
					return m, nil
				}
				if s == "down" {
					m.categoryIndex++
					if m.categoryIndex >= len(availableCategories) {
						m.categoryIndex = 0
					}
					return m, nil
				}
			}

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

			// Wrap around
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

	// Validate text field as user types
	m.validateTextField()

	return m, cmd
}

// validateCurrentField validates the currently focused field
func (m *FormModel) validateCurrentField() {
	switch FormField(m.focusIndex) {
	case TextField:
		m.validateTextField()
	case DateField:
		m.validateDateField()
	}
}

// validateTextField validates the text field
func (m *FormModel) validateTextField() {
	text := m.inputs[0].Value()
	if len(text) > m.maxChars {
		m.fieldErrors[TextField] = fmt.Sprintf("Text exceeds 2000 characters limit")
	} else {
		delete(m.fieldErrors, TextField)
	}
}

// validateDateField validates the date field
func (m *FormModel) validateDateField() {
	dateStr := strings.TrimSpace(m.inputs[1].Value())
	if dateStr == "" {
		delete(m.fieldErrors, DateField)
		return
	}

	eventDate, err := m.parseDate(dateStr)
	if err != nil {
		m.fieldErrors[DateField] = fmt.Sprintf("Invalid date format")
		return
	}

	if eventDate.After(time.Now()) {
		m.fieldErrors[DateField] = fmt.Sprintf("Date cannot be in the future")
		return
	}

	delete(m.fieldErrors, DateField)
}

// View renders the form
func (m *FormModel) View() string {
	// Render form content using FormFieldContainers
	formContent := m.renderFormContentWithContainers()

	// Wrap in a card
	formCard := styles.CardBase.
		Width(styles.MaxWidth(80) - 4).
		Render(formContent)

	// Use header and footer components
	headerView := m.header.View()
	footerView := m.footer.View()

	// Help footer with keyboard shortcuts
	m.helpFooter.SetWidth(styles.MaxWidth(80))
	helpFooterContent := m.helpFooter.View()

	// Combine all sections
	fullContent := lipgloss.JoinVertical(
		lipgloss.Left,
		headerView,
		"",
		formCard,
		"",
		footerView,
		"",
		helpFooterContent,
	)

	return fullContent
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

		// Validate date based on capture mode (only for new events)
		if !m.editMode {
			mode := m.modes[m.modeIndex]
			if mode == careerservice.TimelineJournaling {
				thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
				if eventDate.Before(thirtyDaysAgo) {
					return SubmitMsg{Err: fmt.Errorf("timeline journaling events must be within the last 30 days")}
				}
			}
		}

		// Get optional fields
		company := strings.TrimSpace(m.inputs[2].Value())
		project := strings.TrimSpace(m.inputs[3].Value())

		fmt.Printf("Submitting form: text=%q, date=%v, company=%q, project=%q, mode=%v\n", text, eventDate, company, project, m.modes[m.modeIndex])

		// Get selected tags and categories
		tags := m.tagSelector.SelectedTags()
		categories := m.categorySelector.SelectedCategories()

		// Prepare options
		ctx := context.Background()
		opts := []service.Option{}
		if company != "" {
			opts = append(opts, service.WithCompany(company))
		}
		if project != "" {
			opts = append(opts, service.WithProject(project))
		}
		if len(tags) > 0 {
			opts = append(opts, service.WithTags(tags))
		}
		if len(categories) > 0 {
			opts = append(opts, service.WithCategories(categories))
		}

		// Handle edit mode vs create mode
		if m.editMode {
			// Update existing event
			if err := m.cliService.UpdateEvent(ctx, m.editEventID, text, eventDate, opts...); err != nil {
				return SubmitMsg{Err: fmt.Errorf("failed to update event: %w", err)}
			}
		} else {
			// Create new event
			mode := m.modes[m.modeIndex]
			if err := m.cliService.CaptureEvent(ctx, text, eventDate, mode, opts...); err != nil {
				return SubmitMsg{Err: fmt.Errorf("failed to capture event: %w", err)}
			}
		}

		// Build event for display (this is just for UI purposes)
		event := &career.CareerEvent{
			ID:         m.editEventID, // Will be empty for new events
			Text:       text,
			Date:       eventDate,
			Company:    company,
			Project:    project,
			Tags:       tags,
			Categories: categories,
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

// TagSelector returns the tag selector instance
func (m *FormModel) TagSelector() *components.TagSelector {
	return m.tagSelector
}

// Reset resets the form to its initial state
func (m *FormModel) Reset() {
	m.focusIndex = 0
	m.modeIndex = 0
	m.submitted = false
	m.event = nil
	m.err = nil
	m.charCount = 0
	m.fieldErrors = make(map[FormField]string)

	for i := range m.inputs {
		m.inputs[i].SetValue("")
	}

	// Reset tag and category selectors
	m.tagSelector.Reset()
	m.categorySelector.Clear()

	m.inputs[0].Focus()
}

// SetInitialMode sets the initial capture mode for the form
func (m *FormModel) SetInitialMode(mode string) {
	for i, m_mode := range m.modes {
		if string(m_mode) == mode {
			m.modeIndex = i
			break
		}
	}
}

// LoadEventForEditing populates the form with an existing event's data
func (m *FormModel) LoadEventForEditing(event *career.CareerEvent) {
	m.editMode = true
	m.editEventID = event.ID
	m.event = event

	// Populate text field
	m.inputs[0].SetValue(event.Text)
	m.charCount = len(event.Text)

	// Populate date field
	m.inputs[1].SetValue(event.Date.Format("2006-01-02"))

	// Populate company field
	if event.Company != "" {
		m.inputs[2].SetValue(event.Company)
	}

	// Populate project field
	if event.Project != "" {
		m.inputs[3].SetValue(event.Project)
	}

	// Set tags in tag selector
	if len(event.Tags) > 0 {
		m.tagSelector.SetSelectedTags(event.Tags)
	}

	// Set categories in category selector
	if len(event.Categories) > 0 {
		m.categorySelector.SetSelected(event.Categories)
	}

	// Focus first input
	m.inputs[0].Focus()
}

// IsEditMode returns whether the form is in edit mode
func (m *FormModel) IsEditMode() bool {
	return m.editMode
}

// GetEditEventID returns the ID of the event being edited
func (m *FormModel) GetEditEventID() string {
	return m.editEventID
}

// CategorySelector returns the category selector instance
func (m *FormModel) CategorySelector() *components.CategorySelector {
	return m.categorySelector
}

// renderTagSelector renders the tag selector with available tags
func (m *FormModel) renderTagSelector() string {
	var b strings.Builder
	availableTags := m.tagSelector.AvailableTags()

	for i, tag := range availableTags {
		selected := ""
		if i == m.tagIndex {
			selected = "►"
		}

		isSelected := m.tagSelector.IsSelected(tag)
		checkbox := "☐"
		if isSelected {
			checkbox = "☑"
		}

		b.WriteString(fmt.Sprintf("║   %s %s %-20s║\n", selected, checkbox, tag))
	}

	return b.String()
}

// renderCategorySelector renders the category selector with available categories
func (m *FormModel) renderCategorySelector() string {
	var b strings.Builder
	availableCategories := m.categorySelector.AvailableCategories()

	for i, category := range availableCategories {
		selected := ""
		if i == m.categoryIndex {
			selected = "►"
		}

		isSelected := m.categorySelector.IsSelected(category)
		checkbox := "☐"
		if isSelected {
			checkbox = "☑"
		}

		b.WriteString(fmt.Sprintf("║   %s %s %-20s║\n", selected, checkbox, category))
	}

	return b.String()
}

// SetBreadcrumbs sets breadcrumb trail for display in header

// renderTextFieldWithContainer renders the text field using FormFieldContainer
// Helper to add focus indicator to rendered field
func (m *FormModel) addFocusIndicatorToField(fieldContent string, focused bool) string {
	if focused {
		// Add focus indicator before the first line
		lines := strings.Split(fieldContent, "\n")
		if len(lines) > 0 {
			lines[0] = "► " + lines[0]
			return strings.Join(lines, "\n")
		}
		return "► " + fieldContent
	}
	// Add space to align with focused fields
	lines := strings.Split(fieldContent, "\n")
	if len(lines) > 0 {
		lines[0] = "  " + lines[0]
		return strings.Join(lines, "\n")
	}
	return "  " + fieldContent
}

func (m *FormModel) renderTextFieldWithContainer() string {
	focused := m.focusIndex == int(TextField)
	fieldErr := ""
	if err, ok := m.fieldErrors[TextField]; ok {
		fieldErr = err
	}

	charInfo := fmt.Sprintf("Characters: %d/%d %s",
		m.charCount, m.maxChars, m.getCharCountIndicator())

	fieldContent := components.NewFormFieldContainer().
		SetLabel("Event Text (required):").
		SetInput(m.inputs[0].View()).
		SetHint(charInfo).
		SetError(fieldErr).
		SetFocused(focused).
		Render()

	return m.addFocusIndicatorToField(fieldContent, focused)
}

// renderDateFieldWithContainer renders the date field using FormFieldContainer
func (m *FormModel) renderDateFieldWithContainer() string {
	focused := m.focusIndex == int(DateField)
	fieldErr := ""
	if err, ok := m.fieldErrors[DateField]; ok {
		fieldErr = err
	}

	fieldContent := components.NewFormFieldContainer().
		SetLabel("Date (optional):").
		SetInput(m.inputs[1].View()).
		SetError(fieldErr).
		SetFocused(focused).
		Render()

	return m.addFocusIndicatorToField(fieldContent, focused)
}

// renderCompanyFieldWithContainer renders the company field using FormFieldContainer
func (m *FormModel) renderCompanyFieldWithContainer() string {
	focused := m.focusIndex == int(CompanyField)
	fieldErr := ""
	if err, ok := m.fieldErrors[CompanyField]; ok {
		fieldErr = err
	}

	fieldContent := components.NewFormFieldContainer().
		SetLabel("Company (optional):").
		SetInput(m.inputs[2].View()).
		SetError(fieldErr).
		SetFocused(focused).
		Render()

	return m.addFocusIndicatorToField(fieldContent, focused)
}

// renderProjectFieldWithContainer renders the project field using FormFieldContainer
func (m *FormModel) renderProjectFieldWithContainer() string {
	focused := m.focusIndex == int(ProjectField)
	fieldErr := ""
	if err, ok := m.fieldErrors[ProjectField]; ok {
		fieldErr = err
	}

	fieldContent := components.NewFormFieldContainer().
		SetLabel("Project (optional):").
		SetInput(m.inputs[3].View()).
		SetError(fieldErr).
		SetFocused(focused).
		Render()

	return m.addFocusIndicatorToField(fieldContent, focused)
}

// renderTagsFieldWithContainer renders the tags field using FormFieldContainer
func (m *FormModel) renderTagsFieldWithContainer() string {
	focused := m.focusIndex == int(TagsField)
	var tagsDisplay string
	if focused {
		tagsDisplay = m.renderTagSelector()
	} else {
		selectedTags := m.tagSelector.SelectedTags()
		if len(selectedTags) > 0 {
			for _, tag := range selectedTags {
				tagsDisplay += styles.TagBase.Render(tag) + " "
			}
		} else {
			tagsDisplay = styles.InfoText.Render("(none selected)")
		}
	}

	hint := "Use Up/Down to navigate, Space to toggle"
	if !focused {
		hint = ""
	}

	fieldContent := components.NewFormFieldContainer().
		SetLabel("Tags:").
		SetInput(tagsDisplay).
		SetHint(hint).
		SetFocused(focused).
		Render()

	return m.addFocusIndicatorToField(fieldContent, focused)
}

// renderCategoriesFieldWithContainer renders the categories field using FormFieldContainer
func (m *FormModel) renderCategoriesFieldWithContainer() string {
	focused := m.focusIndex == int(CategoriesField)
	var categoriesDisplay string
	if focused {
		categoriesDisplay = m.renderCategorySelector()
	} else {
		selectedCategories := m.categorySelector.SelectedCategories()
		if len(selectedCategories) > 0 {
			for _, category := range selectedCategories {
				categoriesDisplay += styles.TagBase.Render(category) + " "
			}
		} else {
			categoriesDisplay = styles.InfoText.Render("(none selected)")
		}
	}

	hint := "Use Up/Down to navigate, Space to toggle"
	if !focused {
		hint = ""
	}

	fieldContent := components.NewFormFieldContainer().
		SetLabel("Categories:").
		SetInput(categoriesDisplay).
		SetHint(hint).
		SetFocused(focused).
		Render()

	return m.addFocusIndicatorToField(fieldContent, focused)
}

// renderModeFieldWithContainer renders the mode field using FormFieldContainer
func (m *FormModel) renderModeFieldWithContainer() string {
	focused := m.focusIndex == int(ModeField)
	modeContent := m.renderModeSelector()

	hint := "Use Up/Down to navigate"
	if !focused {
		hint = ""
	}

	fieldContent := components.NewFormFieldContainer().
		SetLabel("Capture Mode:").
		SetInput(modeContent).
		SetHint(hint).
		SetFocused(focused).
		Render()

	return m.addFocusIndicatorToField(fieldContent, focused)
}

// renderSubmitButtonWithContainer renders the submit button
func (m *FormModel) renderSubmitButtonWithContainer() string {
	focused := m.focusIndex == int(SubmitButton)
	var submitBtn string
	if focused {
		submitBtn = styles.ButtonPrimary.Render("[ > Submit < ]")
	} else {
		submitBtn = styles.ButtonPrimary.Render("[ Submit ]")
	}
	return submitBtn
}

// renderFormContentWithContainers renders all form fields using FormFieldContainers
func (m *FormModel) renderFormContentWithContainers() string {
	var content []string

	// Add all fields
	content = append(content,
		m.renderTextFieldWithContainer(),
		m.renderDateFieldWithContainer(),
		m.renderCompanyFieldWithContainer(),
		m.renderProjectFieldWithContainer(),
		m.renderTagsFieldWithContainer(),
		m.renderCategoriesFieldWithContainer(),
		m.renderModeFieldWithContainer(),
		m.renderSubmitButtonWithContainer(),
	)

	// Add model-level error if present
	if m.err != nil {
		content = append(content, styles.ErrorText.Render(m.err.Error()))
	}

	return strings.Join(content, "\n\n")
}

func (m *FormModel) SetBreadcrumbs(crumbs []string) {
	m.breadcrumbs = crumbs
	m.header.SetBreadcrumbs(crumbs)
}
