package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	cliservice "github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/cli/validation"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Focus indices for navigation
const (
	MetadataDateFieldIdx = iota
	MetadataCompanyFieldIdx
	MetadataProjectFieldIdx
	MetadataTagsFieldIdx
	MetadataCategoriesFieldIdx
	MetadataSaveButtonIdx
	MetadataCancelButtonIdx
	MetadataFieldCount // Total number of focus positions
)

// MetadataEditorModel represents the metadata editor form state
type MetadataEditorModel struct {
	*BaseStandardModel
	event            *career.CareerEvent
	originalEvent    *career.CareerEvent // For reverting changes
	service          *careerservice.Service
	cliService       *cliservice.CLIEventService // For persisting metadata changes
	ctx              context.Context
	inputs           []textinput.Model
	focusIndex       int
	tagIndex         int // Index for tag navigation within tags field
	categoryIndex    int // Index for category navigation within categories field
	err              error
	submitted        bool
	cancelled        bool
	tagSelector      *components.TagSelector
	categorySelector *components.CategorySelector
	fieldErrors      map[int]string // Map of field index to error message
	width            int
	height           int
	validator        *validation.MetadataValidator
	calculator       *careerservice.DataQualityCalculator
	helpFooter       components.HelpFooterModel // Help footer
}

// NewMetadataEditorModel creates a new metadata editor model
func NewMetadataEditorModel(event *career.CareerEvent, service *careerservice.Service, cliSvc *cliservice.CLIEventService, ctx context.Context) *MetadataEditorModel {
	// Create input fields (3 fields: date, company, project)
	inputs := make([]textinput.Model, 3)

	// Date input
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "YYYY-MM-DD or 'today', '1 week ago'"
	inputs[0].SetValue(event.Date.Format("2006-01-02"))
	inputs[0].Width = 60
	inputs[0].Focus()

	// Company input
	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Company name (optional)"
	inputs[1].SetValue(event.Company)
	inputs[1].Width = 60

	// Project input
	inputs[2] = textinput.New()
	inputs[2].Placeholder = "Project name (optional)"
	inputs[2].SetValue(event.Project)
	inputs[2].Width = 60

	// Create a copy of the event for reverting
	eventCopy := *event

	// Create tag and category selectors
	tagSelector := components.NewTagSelector()
	tagSelector.SetSelectedTags(event.Tags)

	categorySelector := components.NewCategorySelector()
	categorySelector.SetSelected(event.Categories)

	return &MetadataEditorModel{
		BaseStandardModel: NewBaseStandardModel(),
		event:             event,
		originalEvent:     &eventCopy,
		service:           service,
		cliService:        cliSvc,
		ctx:               ctx,
		inputs:            inputs,
		focusIndex:        0,
		tagIndex:          0,
		categoryIndex:     0,
		err:               nil,
		submitted:         false,
		cancelled:         false,
		tagSelector:       tagSelector,
		categorySelector:  categorySelector,
		fieldErrors:       make(map[int]string),
		validator:         validation.NewMetadataValidator(),
		calculator:        careerservice.NewDataQualityCalculator(),
	}
}

// Init initializes the model
func (m *MetadataEditorModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *MetadataEditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m *MetadataEditorModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab":
		m.focusIndex = (m.focusIndex + 1) % MetadataFieldCount
		m.updateInputFocus()
		return m, nil

	case "shift+tab":
		m.focusIndex = (m.focusIndex - 1 + MetadataFieldCount) % MetadataFieldCount
		m.updateInputFocus()
		return m, nil

	case "up":
		if m.focusIndex == MetadataTagsFieldIdx {
			if m.tagIndex > 0 {
				m.tagIndex--
			}
		} else if m.focusIndex == MetadataCategoriesFieldIdx {
			if m.categoryIndex > 0 {
				m.categoryIndex--
			}
		} else {
			m.focusIndex = (m.focusIndex - 1 + MetadataFieldCount) % MetadataFieldCount
			m.updateInputFocus()
		}
		return m, nil

	case "down":
		if m.focusIndex == MetadataTagsFieldIdx {
			availableTags := m.tagSelector.AvailableTags()
			if m.tagIndex < len(availableTags)-1 {
				m.tagIndex++
			}
		} else if m.focusIndex == MetadataCategoriesFieldIdx {
			availableCategories := m.categorySelector.AvailableCategories()
			if m.categoryIndex < len(availableCategories)-1 {
				m.categoryIndex++
			}
		} else {
			m.focusIndex = (m.focusIndex + 1) % MetadataFieldCount
			m.updateInputFocus()
		}
		return m, nil

	case "space":
		if m.focusIndex == MetadataTagsFieldIdx {
			availableTags := m.tagSelector.AvailableTags()
			if m.tagIndex < len(availableTags) {
				tag := availableTags[m.tagIndex]
				_ = m.tagSelector.ToggleTag(tag)
			}
		} else if m.focusIndex == MetadataCategoriesFieldIdx {
			availableCategories := m.categorySelector.AvailableCategories()
			if m.categoryIndex < len(availableCategories) {
				category := availableCategories[m.categoryIndex]
				_ = m.categorySelector.ToggleCategory(category)
			}
		}
		return m, nil

	case "enter":
		if m.focusIndex == MetadataSaveButtonIdx {
			return m.saveChanges()
		} else if m.focusIndex == MetadataCancelButtonIdx {
			m.cancelled = true
			return m, tea.Quit
		}
		return m, nil

	case "esc":
		m.cancelled = true
		return m, tea.Quit

	default:
		// Handle text input for text fields
		if m.focusIndex < MetadataTagsFieldIdx {
			m.inputs[m.focusIndex].SetValue(m.inputs[m.focusIndex].Value() + msg.String())
		}
	}

	return m, nil
}

func (m *MetadataEditorModel) updateInputFocus() {
	for i := range m.inputs {
		if i == m.focusIndex && m.focusIndex < MetadataTagsFieldIdx {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
}

func (m *MetadataEditorModel) saveChanges() (tea.Model, tea.Cmd) {
	// Update event with new values
	m.event.Company = strings.TrimSpace(m.inputs[MetadataCompanyFieldIdx].Value())
	m.event.Project = strings.TrimSpace(m.inputs[MetadataProjectFieldIdx].Value())
	m.event.Tags = m.tagSelector.SelectedTags()
	m.event.Categories = m.categorySelector.SelectedCategories()

	// Parse and validate date
	dateStr := strings.TrimSpace(m.inputs[MetadataDateFieldIdx].Value())
	if dateStr != "" {
		parsedDate, err := m.parseDate(dateStr)
		if err != nil {
			m.fieldErrors[MetadataDateFieldIdx] = err.Error()
			m.focusIndex = MetadataDateFieldIdx
			return m, nil
		}
		m.event.Date = parsedDate
	}

	// Validate all fields
	if err := m.event.Validate(); err != nil {
		m.err = err
		return m, nil
	}

	// Update timestamps
	m.event.UpdatedAt = time.Now()

	// Persist metadata changes to service
	if m.cliService != nil {
		if err := m.cliService.UpdateEventMetadata(m.ctx, m.event); err != nil {
			m.err = fmt.Errorf("failed to save metadata: %w", err)
			return m, nil
		}
	}

	m.submitted = true
	return m, tea.Quit
}

func (m *MetadataEditorModel) parseDate(dateStr string) (time.Time, error) {
	// Try ISO format first
	t, err := time.Parse("2006-01-02", dateStr)
	if err == nil {
		return t, nil
	}

	// Try "today"
	if strings.ToLower(dateStr) == "today" {
		return time.Now(), nil
	}

	// For now, return the original date if parsing fails
	return m.event.Date, nil
}

// View renders the metadata editor

// View renders the metadata editor using FormFieldContainers
func (m *MetadataEditorModel) View() string {
	if m.width == 0 {
		m.width = 80
	}
	if m.height == 0 {
		m.height = 24
	}

	// Render form content using FormFieldContainers
	formContent := m.renderFormContentWithContainers()

	// Use header and footer components
	headerView := components.NewHeader("Edit Event Metadata", m.width).View()
	footerView := components.NewFooter(m.width).View()

	// Help footer with keyboard shortcuts
	m.helpFooter.SetWidth(m.width)
	helpFooterContent := m.helpFooter.View()

	// Combine all sections
	fullContent := strings.Join([]string{
		headerView,
		"",
		formContent,
		"",
		footerView,
		"",
		helpFooterContent,
	}, "\n")

	return fullContent
}

// renderFormContentWithContainers renders all form fields using FormFieldContainers
func (m *MetadataEditorModel) renderFormContentWithContainers() string {
	var content []string

	// Add all fields
	content = append(content,
		m.renderDateFieldWithContainer(),
		m.renderCompanyFieldWithContainer(),
		m.renderProjectFieldWithContainer(),
		m.renderTagsFieldWithContainer(),
		m.renderCategoriesFieldWithContainer(),
		m.renderButtonsWithContainer(),
	)

	// Add model-level error if present
	if m.err != nil {
		content = append(content, styles.ErrorBox.Render("Error: "+m.err.Error()))
	}

	// Add field-specific errors
	for fieldIdx, errMsg := range m.fieldErrors {
		content = append(content, styles.ErrorBox.Render(fmt.Sprintf("Field %d: %s", fieldIdx, errMsg)))
	}

	return strings.Join(content, "\n\n")
}

// renderDateFieldWithContainer renders the date field using FormFieldContainer
func (m *MetadataEditorModel) renderDateFieldWithContainer() string {
	focused := m.focusIndex == MetadataDateFieldIdx
	fieldErr := ""
	if err, ok := m.fieldErrors[MetadataDateFieldIdx]; ok {
		fieldErr = err
	}

	fieldContent := components.NewFormFieldContainer().
		SetLabel("Date (required):").
		SetInput(m.inputs[MetadataDateFieldIdx].View()).
		SetError(fieldErr).
		SetFocused(focused).
		Render()

	return m.addFocusIndicatorToField(fieldContent, focused)
}

// renderCompanyFieldWithContainer renders the company field using FormFieldContainer
func (m *MetadataEditorModel) renderCompanyFieldWithContainer() string {
	focused := m.focusIndex == MetadataCompanyFieldIdx
	fieldErr := ""
	if err, ok := m.fieldErrors[MetadataCompanyFieldIdx]; ok {
		fieldErr = err
	}

	fieldContent := components.NewFormFieldContainer().
		SetLabel("Company:").
		SetInput(m.inputs[MetadataCompanyFieldIdx].View()).
		SetError(fieldErr).
		SetFocused(focused).
		Render()

	return m.addFocusIndicatorToField(fieldContent, focused)
}

// renderProjectFieldWithContainer renders the project field using FormFieldContainer
func (m *MetadataEditorModel) renderProjectFieldWithContainer() string {
	focused := m.focusIndex == MetadataProjectFieldIdx
	fieldErr := ""
	if err, ok := m.fieldErrors[MetadataProjectFieldIdx]; ok {
		fieldErr = err
	}

	fieldContent := components.NewFormFieldContainer().
		SetLabel("Project:").
		SetInput(m.inputs[MetadataProjectFieldIdx].View()).
		SetError(fieldErr).
		SetFocused(focused).
		Render()

	return m.addFocusIndicatorToField(fieldContent, focused)
}

// renderTagsFieldWithContainer renders the tags field using FormFieldContainer
// renderTagsFieldWithContainer renders the tags field using FormFieldContainer
func (m *MetadataEditorModel) renderTagsFieldWithContainer() string {
	focused := m.focusIndex == MetadataTagsFieldIdx
	fieldErr := ""
	if err, ok := m.fieldErrors[MetadataTagsFieldIdx]; ok {
		fieldErr = err
	}

	// Render tags as badges
	availableTags := m.tagSelector.AvailableTags()
	selectedTags := m.tagSelector.SelectedTags()
	selectedMap := make(map[string]bool)
	for _, tag := range selectedTags {
		selectedMap[tag] = true
	}

	var badges []string
	for i, tag := range availableTags {
		isFocused := m.focusIndex == MetadataTagsFieldIdx && i == m.tagIndex
		indicator := "[ ]"
		if selectedMap[tag] {
			indicator = "[✓]"
		}

		badge := fmt.Sprintf("%s %s", indicator, tag)
		if isFocused {
			badge = styles.ButtonPrimaryFocused.Render(badge)
		} else if selectedMap[tag] {
			badge = styles.ButtonPrimary.Render(badge)
		}
		badges = append(badges, badge)
	}
	badgesStr := strings.Join(badges, "  ")

	hint := fmt.Sprintf("Navigate with ↑↓ | Select with Space | %d selected", len(selectedTags))

	fieldContent := components.NewFormFieldContainer().
		SetLabel("Tags:").
		SetInput(badgesStr).
		SetHint(hint).
		SetError(fieldErr).
		SetFocused(focused).
		Render()

	return m.addFocusIndicatorToField(fieldContent, focused)
}

// renderCategoriesFieldWithContainer renders the categories field using FormFieldContainer
func (m *MetadataEditorModel) renderCategoriesFieldWithContainer() string {
	focused := m.focusIndex == MetadataCategoriesFieldIdx
	fieldErr := ""
	if err, ok := m.fieldErrors[MetadataCategoriesFieldIdx]; ok {
		fieldErr = err
	}

	// Render categories as badges
	availableCategories := m.categorySelector.AvailableCategories()
	selectedCategories := m.categorySelector.SelectedCategories()
	selectedMap := make(map[string]bool)
	for _, cat := range selectedCategories {
		selectedMap[cat] = true
	}

	var badges []string
	for i, cat := range availableCategories {
		isFocused := m.focusIndex == MetadataCategoriesFieldIdx && i == m.categoryIndex
		indicator := "[ ]"
		if selectedMap[cat] {
			indicator = "[✓]"
		}

		badge := fmt.Sprintf("%s %s", indicator, cat)
		if isFocused {
			badge = styles.ButtonPrimaryFocused.Render(badge)
		} else if selectedMap[cat] {
			badge = styles.ButtonPrimary.Render(badge)
		}
		badges = append(badges, badge)
	}
	badgesStr := strings.Join(badges, "  ")

	hint := fmt.Sprintf("Navigate with ↑↓ | Select with Space | %d selected", len(selectedCategories))

	fieldContent := components.NewFormFieldContainer().
		SetLabel("Categories:").
		SetInput(badgesStr).
		SetHint(hint).
		SetError(fieldErr).
		SetFocused(focused).
		Render()

	return m.addFocusIndicatorToField(fieldContent, focused)
}
func (m *MetadataEditorModel) renderButtonsWithContainer() string {
	focused := m.focusIndex >= MetadataSaveButtonIdx

	saveBtn := "[ Save ]"
	cancelBtn := "[ Cancel ]"

	if m.focusIndex == MetadataSaveButtonIdx {
		saveBtn = styles.ButtonPrimaryFocused.Render(saveBtn)
		cancelBtn = styles.ButtonSecondary.Render(cancelBtn)
	} else if m.focusIndex == MetadataCancelButtonIdx {
		saveBtn = styles.ButtonPrimary.Render(saveBtn)
		cancelBtn = styles.ButtonSecondaryFocused.Render(cancelBtn)
	} else {
		saveBtn = styles.ButtonPrimary.Render(saveBtn)
		cancelBtn = styles.ButtonSecondary.Render(cancelBtn)
	}

	buttonsStr := strings.Join([]string{saveBtn, cancelBtn}, "  ")

	fieldContent := components.NewFormFieldContainer().
		SetInput(buttonsStr).
		SetFocused(focused).
		Render()

	return m.addFocusIndicatorToField(fieldContent, focused)
}

// addFocusIndicatorToField adds a focus indicator to the rendered field
func (m *MetadataEditorModel) addFocusIndicatorToField(fieldContent string, focused bool) string {
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

// Revert reverts changes to original event
func (m *MetadataEditorModel) Revert() {
	*m.event = *m.originalEvent
}

// IsSubmitted returns whether the form was submitted
func (m *MetadataEditorModel) IsSubmitted() bool {
	return m.submitted
}

// IsCancelled returns whether the form was cancelled
func (m *MetadataEditorModel) IsCancelled() bool {
	return m.cancelled
}

// GetEvent returns the edited event
func (m *MetadataEditorModel) GetEvent() *career.CareerEvent {
	return m.event
}
