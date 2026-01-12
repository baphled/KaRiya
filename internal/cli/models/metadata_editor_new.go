package models

import (
	"context"
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/forms"
	cliservice "github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/styles"
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
	err              error
	submitted        bool
	cancelled        bool
	width            int
	height           int
	helpFooter       components.HelpFooterModel
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

	// Extract form data from event
	formData := forms.GetMetadataFormData(event)

	// Create huh form with available tags and categories
	form := forms.NewMetadataEditorFormWithData(formData, availableTags, availableCategories)

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
		err:               nil,
		submitted:         false,
		cancelled:         false,
		width:             80,
		height:            24,
		helpFooter:        components.NewHelpFooter("metadata_editor", 80),
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
		errorStyle := lipgloss.NewStyle().
			Foreground(styles.ColorError).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorError).
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
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F38BA8")). // Catppuccin Red
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#F38BA8")).
			Padding(1, 2).
			MarginTop(1)

		formView += "\n\n" + errorStyle.Render(m.err.Error())
	}

	// Use header and footer components
	headerView := components.NewHeader("Edit Event Metadata", m.width).View()
	footerView := components.NewFooter(m.width).View()

	// Render help footer
	m.helpFooter.SetWidth(m.width)
	helpFooterContent := m.helpFooter.View()

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
		"",
		helpFooterContent,
	)

	return fullContent
}
