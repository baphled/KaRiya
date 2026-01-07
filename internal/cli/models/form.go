package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
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
	SubmitButton
)

// FormModel represents the event capture form state
type FormModel struct {
	*BaseStandardModel
	cliService         *service.CLIEventService
	inputs             []textinput.Model
	focusIndex         int
	tagIndex           int // Index for tag navigation when TagsField is focused
	categoryIndex      int // Index for category navigation when CategoriesField is focused
	err                error
	submitted          bool
	event              *career.CareerEvent
	charCount          int
	maxChars           int
	tagSelector        *components.TagSelector
	categorySelector   *components.CategorySelector
	fieldErrors        map[FormField]string       // Track field-level validation errors
	editMode           bool                       // True if editing an existing event
	editEventID        string                     // ID of event being edited
	helpFooter         components.HelpFooterModel // Help footer for keyboard shortcuts
	header             components.HeaderModel     // Header component
	footer             components.FooterModel     // Footer component
	breadcrumbs        []string                   // Navigation breadcrumb trail
	width              int                        // Available terminal width
	height             int                        // Available terminal height
	strategy           string                     // Capture strategy: "quick" or "manual"
	showOptionalFields bool                       // Toggle for optional field visibility (manual mode only)
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
	inputs[0].Width = 60 // Will be updated dynamically based on terminal size

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

	return &FormModel{
		BaseStandardModel:  NewBaseStandardModel(),
		cliService:         cliService,
		inputs:             inputs,
		focusIndex:         0,
		tagIndex:           0,
		categoryIndex:      0,
		err:                nil,
		submitted:          false,
		maxChars:           2000,
		tagSelector:        components.NewTagSelector(),
		categorySelector:   components.NewCategorySelector(),
		fieldErrors:        make(map[FormField]string),
		helpFooter:         components.NewHelpFooter("form", 80),
		header:             components.NewHeader("Capture Career Event", 80),
		footer:             components.NewFooter(80),
		strategy:           "manual", // Default to manual mode (show all fields)
		showOptionalFields: true,     // Show optional fields by default
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
		m.width = msg.Width
		m.height = msg.Height
		m.helpFooter.SetWidth(msg.Width)

		// Update input field widths adaptively
		adaptiveWidth := m.getAdaptiveFieldWidth()
		for i := range m.inputs {
			m.inputs[i].Width = adaptiveWidth
		}

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

			// If not in a navigation field, fall through to text input

		case "ctrl+o":
			// Toggle optional fields visibility (Ctrl+O works in any mode, any field)
			// Note: This replaces the old 't' key toggle to avoid conflicts with text input
			if m.strategy == "manual" {
				m.ToggleOptionalFields()
				// If hiding fields, ensure focus is on a visible field
				if !m.showOptionalFields && !m.isFieldVisible(FormField(m.focusIndex)) {
					m.focusIndex = int(TextField)
					cmds := m.updateFocus()
					return m, tea.Batch(cmds...)
				}
			}
			return m, nil

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

			// Handle navigation
			if s == "enter" && m.focusIndex == int(SubmitButton) {
				return m, m.submitForm()
			}

			// Navigate to next/previous field, skipping hidden fields
			direction := 1
			if s == "up" || s == "shift+tab" {
				direction = -1
			}

			// Move focus and skip hidden fields
			startIndex := m.focusIndex
			for {
				m.focusIndex += direction

				// Wrap around
				if m.focusIndex > int(SubmitButton) {
					m.focusIndex = 0
				} else if m.focusIndex < 0 {
					m.focusIndex = int(SubmitButton)
				}

				// If this field is visible, or we've looped back to start, stop
				if m.isFieldVisible(FormField(m.focusIndex)) || m.focusIndex == startIndex {
					break
				}
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

// isFieldVisible returns whether a field should be displayed based on current strategy
func (m *FormModel) isFieldVisible(field FormField) bool {
	// TextField and SubmitButton are always visible
	if field == TextField || field == SubmitButton {
		return true
	}

	// In quick mode, optional fields are hidden
	// In manual mode, respect the showOptionalFields toggle
	if m.strategy == "quick" {
		return false
	}

	// In manual mode or edit mode, check the toggle
	if field == DateField || field == CompanyField || field == ProjectField || field == TagsField || field == CategoriesField {
		return m.showOptionalFields
	}

	return true
}

// getAdaptiveFieldWidth returns an appropriate field width based on terminal size
func (m *FormModel) getAdaptiveFieldWidth() int {
	if m.width == 0 {
		return 60 // Default width
	}

	// Use 80% of available width, with min/max bounds
	width := int(float64(m.width) * 0.8)

	// Minimum width of 40 characters
	if width < 40 {
		return 40
	}

	// Maximum width of 80 characters for readability
	if width > 80 {
		return 80
	}

	return width
}

// isCompactMode returns true if the terminal is too small for full rendering
func (m *FormModel) isCompactMode() bool {
	return m.height > 0 && m.height < 25
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
		m.fieldErrors[TextField] = "Text exceeds 2000 characters limit"
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
		m.fieldErrors[DateField] = "Invalid date format"
		return
	}

	if eventDate.After(time.Now()) {
		m.fieldErrors[DateField] = "Date cannot be in the future"
		return
	}

	delete(m.fieldErrors, DateField)
}

// View renders the form
func (m *FormModel) View() string {
	// Render form content using FormFieldContainers
	formContent := m.renderFormContentWithContainers()

	// Wrap in a card (no title, no navigation helper)
	formCard := styles.CardBase.
		Width(styles.MaxWidth(80) - 4).
		Render(formContent)

	return formCard
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

		// Get optional fields
		company := strings.TrimSpace(m.inputs[2].Value())
		project := strings.TrimSpace(m.inputs[3].Value())

		// Get selected tags and categories
		tags := m.tagSelector.SelectedTags()
		categories := m.categorySelector.SelectedCategories()

		// NOTE: Form does NOT save the event. The intent is responsible for saving
		// the event after the review phase is complete.
		// The form only validates and collects the data.

		// Build event object with collected data
		event := &career.CareerEvent{
			ID:         m.editEventID, // Empty for new events, set for edits
			Text:       text,
			Date:       eventDate,
			Company:    company,
			Project:    project,
			Tags:       tags,
			Categories: categories,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		// Return the collected data without saving to database
		return SubmitMsg{Event: event, Err: nil}
	}
}

// SubmitForm triggers form submission
// This is called by the intent when user presses Ctrl+S or clicks the Submit button
// It validates all form fields and returns a SubmitMsg with the collected data
func (m *FormModel) SubmitForm() tea.Cmd {
	return m.submitForm()
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

// SetStrategy sets the capture strategy and initializes field visibility accordingly
func (m *FormModel) SetStrategy(strategy string) {
	m.strategy = strategy
	// Quick mode: hide optional fields (show only required event text)
	// Manual mode: show all fields by default
	m.showOptionalFields = (strategy == "manual")
}

// ToggleOptionalFields toggles the visibility of optional fields in manual mode
func (m *FormModel) ToggleOptionalFields() {
	if m.strategy == "manual" {
		m.showOptionalFields = !m.showOptionalFields
	}
}

// ShowOptionalFields returns whether optional fields should be displayed
func (m *FormModel) ShowOptionalFields() bool {
	return m.showOptionalFields
}

// GetStrategy returns the current capture strategy
func (m *FormModel) GetStrategy() string {
	return m.strategy
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

// renderFormContentWithContainers renders all form fields using the smart FormContainer
func (m *FormModel) renderFormContentWithContainers() string {
	var formParts []string
	spacing := "\n\n"
	if m.isCompactMode() {
		spacing = "\n"
	}

	// Render each field/section manually for full control over layout

	// 1. Event Text (required, full width)
	textFieldErr := ""
	if err, ok := m.fieldErrors[TextField]; ok {
		textFieldErr = err
	}
	charInfo := fmt.Sprintf("Characters: %d/%d %s",
		m.charCount, m.maxChars, m.getCharCountIndicator())

	textLabel := styles.Label.Render("Event Text (required):")
	textInput := m.inputs[0].View()
	textField := textLabel + "\n" + textInput
	if textFieldErr != "" {
		textField += "\n" + styles.ErrorText.Render(textFieldErr)
	}
	if charInfo != "" {
		textField += "\n" + styles.InputHint.Render(charInfo)
	}
	formParts = append(formParts, textField)

	// 2. Date (optional, full width)
	if m.isFieldVisible(DateField) {
		dateFieldErr := ""
		if err, ok := m.fieldErrors[DateField]; ok {
			dateFieldErr = err
		}
		dateLabel := styles.Label.Render("Date (optional):")
		dateInput := m.inputs[1].View()
		dateField := dateLabel + "\n" + dateInput
		if dateFieldErr != "" {
			dateField += "\n" + styles.ErrorText.Render(dateFieldErr)
		}
		formParts = append(formParts, dateField)
	}

	// 3. Company and Project side-by-side
	if m.isFieldVisible(CompanyField) || m.isFieldVisible(ProjectField) {
		var companyPart, projectPart string

		// Calculate half-width for side-by-side layout
		// Each column gets approximately half the available width minus spacing
		columnWidth := (m.getAdaptiveFieldWidth() / 2) - 2
		if columnWidth < 25 {
			columnWidth = 25 // Minimum column width
		}

		if m.isFieldVisible(CompanyField) {
			companyFieldErr := ""
			if err, ok := m.fieldErrors[CompanyField]; ok {
				companyFieldErr = err
			}
			companyLabel := styles.Label.Render("Company (optional):")
			companyInput := m.inputs[2].View()
			companyPart = companyLabel + "\n" + companyInput
			if companyFieldErr != "" {
				companyPart += "\n" + styles.ErrorText.Render(companyFieldErr)
			}
			// Apply fixed width to the entire column
			companyPart = lipgloss.NewStyle().Width(columnWidth).Render(companyPart)
		}

		if m.isFieldVisible(ProjectField) {
			projectFieldErr := ""
			if err, ok := m.fieldErrors[ProjectField]; ok {
				projectFieldErr = err
			}
			projectLabel := styles.Label.Render("Project (optional):")
			projectInput := m.inputs[3].View()
			projectPart = projectLabel + "\n" + projectInput
			if projectFieldErr != "" {
				projectPart += "\n" + styles.ErrorText.Render(projectFieldErr)
			}
			// Apply fixed width to the entire column
			projectPart = lipgloss.NewStyle().Width(columnWidth).Render(projectPart)
		}

		// Combine side-by-side if both visible
		if companyPart != "" && projectPart != "" {
			combined := lipgloss.JoinHorizontal(lipgloss.Top, companyPart, "    ", projectPart)
			formParts = append(formParts, combined)
		} else if companyPart != "" {
			formParts = append(formParts, companyPart)
		} else if projectPart != "" {
			formParts = append(formParts, projectPart)
		}
	}

	// 4. Tags and Categories side-by-side
	if m.isFieldVisible(TagsField) || m.isFieldVisible(CategoriesField) {
		var tagsPart, categoriesPart string

		// Use the same column width as Company/Project for consistency
		columnWidth := (m.getAdaptiveFieldWidth() / 2) - 2
		if columnWidth < 25 {
			columnWidth = 25 // Minimum column width
		}

		if m.isFieldVisible(TagsField) {
			var tagsDisplay string
			if m.focusIndex == int(TagsField) {
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
			tagsLabel := styles.Label.Render("Tags:")
			tagsPart = tagsLabel + "\n" + tagsDisplay
			// Apply fixed width to the entire column
			tagsPart = lipgloss.NewStyle().Width(columnWidth).Render(tagsPart)
		}

		if m.isFieldVisible(CategoriesField) {
			var categoriesDisplay string
			if m.focusIndex == int(CategoriesField) {
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
			categoriesLabel := styles.Label.Render("Categories:")
			categoriesPart = categoriesLabel + "\n" + categoriesDisplay
			// Apply fixed width to the entire column
			categoriesPart = lipgloss.NewStyle().Width(columnWidth).Render(categoriesPart)
		}

		// Combine side-by-side if both visible
		if tagsPart != "" && categoriesPart != "" {
			combined := lipgloss.JoinHorizontal(lipgloss.Top, tagsPart, "    ", categoriesPart)
			formParts = append(formParts, combined)
		} else if tagsPart != "" {
			formParts = append(formParts, tagsPart)
		} else if categoriesPart != "" {
			formParts = append(formParts, categoriesPart)
		}
	}

	// 5. Submit button
	var submitBtn string
	if m.focusIndex == int(SubmitButton) {
		submitBtn = styles.ButtonPrimary.Render("[ > Submit < ]")
	} else {
		submitBtn = styles.ButtonPrimary.Render("[ Submit ]")
	}
	formParts = append(formParts, submitBtn)

	// Join all parts with spacing
	formContent := strings.Join(formParts, spacing)

	// Add toggle hint in manual mode
	if m.strategy == "manual" {
		var toggleHint string
		if m.showOptionalFields {
			toggleHint = styles.InfoHint.Render("Ctrl+O to hide optional fields")
		} else {
			toggleHint = styles.InfoHint.Render("Ctrl+O to show optional fields")
		}
		formContent += "\n\n" + toggleHint
	}

	// Add model-level error if present
	if m.err != nil {
		formContent += "\n\n" + styles.ErrorBox.Render(m.err.Error())
	}

	return formContent
}

func (m *FormModel) SetBreadcrumbs(crumbs []string) {
	m.breadcrumbs = crumbs
	// Note: header.SetBreadcrumbs removed - breadcrumbs now handled by StandardView
}
