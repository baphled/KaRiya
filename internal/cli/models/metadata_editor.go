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
		event:            event,
		originalEvent:    &eventCopy,
		service:          service,
		cliService:       cliSvc,
		ctx:              ctx,
		inputs:           inputs,
		focusIndex:       0,
		tagIndex:         0,
		categoryIndex:    0,
		err:              nil,
		submitted:        false,
		cancelled:        false,
		tagSelector:      tagSelector,
		categorySelector: categorySelector,
		fieldErrors:      make(map[int]string),
		validator:        validation.NewMetadataValidator(),
		calculator:       careerservice.NewDataQualityCalculator(),
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
func (m *MetadataEditorModel) View() string {
	if m.width == 0 {
		m.width = 80
	}
	if m.height == 0 {
		m.height = 24
	}

	var sb strings.Builder

	// Header
	sb.WriteString(styles.HeaderMain.Render("Edit Event Metadata") + "\n\n")

	// Date field
	sb.WriteString(m.renderDateField())
	sb.WriteString("\n")

	// Company field
	sb.WriteString(m.renderCompanyField())
	sb.WriteString("\n")

	// Project field
	sb.WriteString(m.renderProjectField())
	sb.WriteString("\n")

	// Tags field
	sb.WriteString(m.renderTagsField())
	sb.WriteString("\n")

	// Categories field
	sb.WriteString(m.renderCategoriesField())
	sb.WriteString("\n\n")

	// Buttons
	sb.WriteString(m.renderButtons())
	sb.WriteString("\n")

	// Error messages
	if m.err != nil {
		sb.WriteString(styles.ErrorBox.Render("Error: " + m.err.Error()))
		sb.WriteString("\n")
	}

	// Field errors
	for fieldIdx, errMsg := range m.fieldErrors {
		sb.WriteString(styles.ErrorBox.Render(fmt.Sprintf("Field %d: %s", fieldIdx, errMsg)))
		sb.WriteString("\n")
	}

	return sb.String()
}

func (m *MetadataEditorModel) renderDateField() string {
	label := "Date"
	if m.focusIndex == MetadataDateFieldIdx {
		label = styles.InputLabel.Render("► " + label)
	} else {
		label = styles.InputLabel.Render(label)
	}
	return fmt.Sprintf("%s\n%s", label, m.inputs[MetadataDateFieldIdx].View())
}

func (m *MetadataEditorModel) renderCompanyField() string {
	label := "Company"
	if m.focusIndex == MetadataCompanyFieldIdx {
		label = styles.InputLabel.Render("► " + label)
	} else {
		label = styles.InputLabel.Render(label)
	}
	return fmt.Sprintf("%s\n%s", label, m.inputs[MetadataCompanyFieldIdx].View())
}

func (m *MetadataEditorModel) renderProjectField() string {
	label := "Project"
	if m.focusIndex == MetadataProjectFieldIdx {
		label = styles.InputLabel.Render("► " + label)
	} else {
		label = styles.InputLabel.Render(label)
	}
	return fmt.Sprintf("%s\n%s", label, m.inputs[MetadataProjectFieldIdx].View())
}

func (m *MetadataEditorModel) renderTagsField() string {
	label := "Tags"
	if m.focusIndex == MetadataTagsFieldIdx {
		label = styles.InputLabel.Render("► " + label)
	} else {
		label = styles.InputLabel.Render(label)
	}

	var sb strings.Builder
	sb.WriteString(label + "\n")

	availableTags := m.tagSelector.AvailableTags()
	selectedTags := m.tagSelector.SelectedTags()
	selectedMap := make(map[string]bool)
	for _, tag := range selectedTags {
		selectedMap[tag] = true
	}

	for i, tag := range availableTags {
		var tagStr string
		if selectedMap[tag] {
			tagStr = "[✓] " + tag
		} else {
			tagStr = "[ ] " + tag
		}

		if m.focusIndex == MetadataTagsFieldIdx && i == m.tagIndex {
			tagStr = styles.InputLabel.Render("► " + tagStr)
		} else if m.focusIndex == MetadataTagsFieldIdx {
			tagStr = "  " + tagStr
		}

		sb.WriteString(tagStr + "  ")
		if (i+1)%3 == 0 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (m *MetadataEditorModel) renderCategoriesField() string {
	label := "Categories"
	if m.focusIndex == MetadataCategoriesFieldIdx {
		label = styles.InputLabel.Render("► " + label)
	} else {
		label = styles.InputLabel.Render(label)
	}

	var sb strings.Builder
	sb.WriteString(label + "\n")

	availableCategories := m.categorySelector.AvailableCategories()
	selectedCategories := m.categorySelector.SelectedCategories()
	selectedMap := make(map[string]bool)
	for _, cat := range selectedCategories {
		selectedMap[cat] = true
	}

	for i, category := range availableCategories {
		var catStr string
		if selectedMap[category] {
			catStr = "[✓] " + category
		} else {
			catStr = "[ ] " + category
		}

		if m.focusIndex == MetadataCategoriesFieldIdx && i == m.categoryIndex {
			catStr = styles.InputLabel.Render("► " + catStr)
		} else if m.focusIndex == MetadataCategoriesFieldIdx {
			catStr = "  " + catStr
		}

		sb.WriteString(catStr + "  ")
		if (i+1)%2 == 0 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (m *MetadataEditorModel) renderButtons() string {
	saveBtn := "Save"
	if m.focusIndex == MetadataSaveButtonIdx {
		saveBtn = styles.ButtonFocused.Render(saveBtn)
	} else {
		saveBtn = styles.ButtonPrimary.Render(saveBtn)
	}

	cancelBtn := "Cancel"
	if m.focusIndex == MetadataCancelButtonIdx {
		cancelBtn = styles.ButtonFocused.Render(cancelBtn)
	} else {
		cancelBtn = styles.ButtonSecondary.Render(cancelBtn)
	}

	return fmt.Sprintf("%s  %s", saveBtn, cancelBtn)
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
