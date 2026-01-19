package models

import (
	"context"
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/navigation"
	cliservice "github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/layout"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// MetadataEditorModelNew represents the metadata editor form state using huh.
// This model uses the huh library for form handling, providing:
// - Automatic focus management (no more manual tab handling)
// - Built-in validation with custom validators
// - Catppuccin theming
// - Consistent keyboard navigation
// - MultiSelect for tags and categories
type MetadataEditorModelNew struct {
	*BaseStandardModel
	event            *career.CareerEvent
	originalEvent    *career.CareerEvent // For reverting changes
	service          *careerservice.Service
	cliService       *cliservice.CLIEventService // For persisting metadata changes
	ctx              context.Context
	form             *huh.Form
	formData         *forms.MetadataFormData
	tagSelector      *components.TagSelector
	categorySelector *components.CategorySelector
	skillSelector    *components.SkillSelector
	err              error
	submitted        bool
	cancelled        bool
	width            int
	height           int
	theme            themes.Theme
}

// NewMetadataEditorModelNew creates a new metadata editor model using huh forms.
func NewMetadataEditorModelNew(event *career.CareerEvent, service *careerservice.Service, cliSvc *cliservice.CLIEventService, ctx context.Context) *MetadataEditorModelNew {
	// Create a copy of the event for reverting
	eventCopy := *event

	// Create tag and category selectors to get available options
	tagSelector := components.NewTagSelector()
	tagSelector.SetSelectedTags(event.Tags)
	availableTags := tagSelector.AvailableTags()

	categorySelector := components.NewCategorySelector()
	_ = categorySelector.SetSelected(event.Categories) // Error ignored: existing event categories should be valid
	availableCategories := categorySelector.AvailableCategories()

	// Load all available skills from repository
	allSkills, err := service.GetSkillRepository().List(ctx, nil)
	if err != nil {
		allSkills = []*career.Skill{} // If error, use empty list
	}

	skillSelector := components.NewSkillSelector(allSkills)
	// Pre-select skills from event
	for _, skillID := range event.Skills {
		_ = skillSelector.SelectSkill(skillID) // Error ignored: event skills should be valid
	}
	availableSkills := skillSelector.AvailableSkills()

	// Extract form data from event
	formData := forms.GetMetadataFormData(event)

	// Create huh form with available tags, categories, and skills
	form := forms.NewMetadataEditorFormWithData(formData, availableTags, availableCategories, availableSkills)

	return &MetadataEditorModelNew{
		BaseStandardModel: NewBaseStandardModel(),
		event:             event,
		originalEvent:     &eventCopy,
		service:           service,
		cliService:        cliSvc,
		ctx:               ctx,
		form:              form,
		formData:          formData,
		tagSelector:       tagSelector,
		categorySelector:  categorySelector,
		skillSelector:     skillSelector,
		err:               nil,
		submitted:         false,
		cancelled:         false,
		width:             80,
		height:            24,
		theme:             themes.NewDefaultTheme(),
	}
}

// Init initializes the model
func (m *MetadataEditorModelNew) Init() tea.Cmd {
	return m.form.Init()
}

// Update handles messages
func (m *MetadataEditorModelNew) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Handle escape BEFORE delegating to form
		// This ensures the parent intent can navigate back
		if msg.String() == "esc" {
			m.cancelled = true
			return m, nil
		}

		// Handle quit
		if msg.String() == "q" || msg.String() == "ctrl+c" {
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
		return m.handleFormCompletion()
	}

	if forms.IsAborted(m.form) {
		m.cancelled = true
		return m, nil
	}

	return m, cmd
}

// handleFormCompletion processes the completed form and saves the metadata.
func (m *MetadataEditorModelNew) handleFormCompletion() (tea.Model, tea.Cmd) {
	// Check if user confirmed via the submit button
	// If they selected "Cancel" on the confirm, treat as cancelled
	if !m.formData.SubmitConfirmed {
		m.cancelled = true
		return m, nil
	}

	// Apply form data to event
	err := forms.ApplyMetadataFormData(m.event, m.formData)
	if err != nil {
		m.err = fmt.Errorf("failed to apply form data: %w", err)
		return m, nil
	}

	// Validate event
	if err := m.event.Validate(); err != nil {
		m.err = fmt.Errorf("validation failed: %w", err)
		return m, nil
	}

	// Update timestamp
	m.event.UpdatedAt = time.Now()

	// Persist metadata changes to service
	if m.cliService != nil {
		if err := m.cliService.UpdateEventMetadata(m.ctx, m.event); err != nil {
			m.err = fmt.Errorf("failed to save metadata: %w", err)
			return m, nil
		}
	}

	m.submitted = true
	return m, nil
}

// GetEvent returns the edited event
func (m *MetadataEditorModelNew) GetEvent() *career.CareerEvent {
	return m.event
}

// IsSubmitted returns true if changes were saved
func (m *MetadataEditorModelNew) IsSubmitted() bool {
	return m.submitted
}

// IsCancelled returns true if operation was cancelled
func (m *MetadataEditorModelNew) IsCancelled() bool {
	return m.cancelled
}

// Revert reverts changes to the original event
func (m *MetadataEditorModelNew) Revert() {
	*m.event = *m.originalEvent
	// Update form data
	m.formData = forms.GetMetadataFormData(m.originalEvent)
}

// GetError returns the current error
func (m *MetadataEditorModelNew) GetError() error {
	return m.err
}

// GetTitle returns the modal title for overlay rendering.
func (m *MetadataEditorModelNew) GetTitle() string {
	return "Edit Event Metadata"
}

// GetContent returns just the form content without header/footer.
// This allows parent intents to compose the modal as an overlay.
func (m *MetadataEditorModelNew) GetContent() string {
	formView := m.form.View()

	// Add error if present
	if m.err != nil {
		errorColor := m.theme.ErrorColor()
		errorStyle := lipgloss.NewStyle().
			Foreground(errorColor).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(errorColor).
			Padding(1, 2).
			MarginTop(1)

		formView += "\n\n" + errorStyle.Render(m.err.Error())
	}

	return formView
}

// GetFooter returns the footer instructions for the modal.
func (m *MetadataEditorModelNew) GetFooter() string {
	return "Enter: Confirm | Esc: Cancel | Tab: Next Field | Shift+Tab: Previous"
}

// View renders the editor UI
func (m *MetadataEditorModelNew) View() string {
	// Render form using huh
	formView := m.form.View()

	// Add error if present
	if m.err != nil {
		errorColor := m.theme.ErrorColor()
		errorStyle := lipgloss.NewStyle().
			Foreground(errorColor).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(errorColor).
			Padding(1, 2).
			MarginTop(1)

		formView += "\n\n" + errorStyle.Render(m.err.Error())
	}

	// Use UIKit layout components
	headerView := layout.NewHeader("Edit Event Metadata", m.width).
		WithTheme(m.theme).
		View()
	footerView := layout.NewFooter(m.width).
		WithTheme(m.theme).
		WithHelp(navigation.GetContextualHelp("metadata_editor")).
		View()

	// Combine all sections
	contentStyle := lipgloss.NewStyle().
		Width(m.width).
		Padding(1, 2)

	fullContent := lipgloss.JoinVertical(
		lipgloss.Left,
		headerView,
		"",
		contentStyle.Render(formView),
		"",
		footerView,
	)

	return fullContent
}
