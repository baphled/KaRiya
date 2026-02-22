package captureevent

import (
	"context"
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/cli/forms"
	cliservice "github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/selectors"
	"github.com/baphled/kariya/internal/domain/capture"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ReviewEnrichmentModel represents the review enrichment form state using huh.
// It allows users to edit event metadata during the event-review flow,
// providing:
// - Automatic focus management (no more manual tab handling)
// - Built-in validation with custom validators
// - Catppuccin theming
// - Consistent keyboard navigation
// - MultiSelect for tags and categories.
type ReviewEnrichmentModel struct {
	forms.EditorFields
	event            *career.Event
	originalEvent    *career.Event
	service          *careerservice.Service
	cliService       *cliservice.CLIEventService
	ctx              context.Context
	formData         *forms.MetadataFormData
	tagSelector      *selectors.TagSelector
	categorySelector *selectors.CategorySelector
	skillSelector    *selectors.SkillSelector
	submitted        bool
}

// ReviewEnrichmentDimensions holds terminal dimensions for the review enrichment modal.
// Pass these from the intent so the form sizes correctly inside the overlay.
type ReviewEnrichmentDimensions struct {
	TerminalWidth  int
	TerminalHeight int
}

// NewReviewEnrichmentModel creates a new review enrichment model using huh forms.
//
// dims may be nil, in which case defaults (80x40) are used.
//
// Expected:
//   - ctx must be a valid context.
//   - event must be a valid career.Event pointer.
//   - service must be a valid careerservice.Service pointer.
//   - cliSvc may be nil, but if provided must be a valid CLIEventService pointer.
//   - dims may be nil, in which case defaults are used.
//
// Returns:
//   - A fully initialized ReviewEnrichmentModel ready for use.
//
// Side effects:
//   - May query skill repository to load available skills.
func NewReviewEnrichmentModel(
	ctx context.Context, event *career.Event, service *careerservice.Service,
	cliSvc *cliservice.CLIEventService, dims *ReviewEnrichmentDimensions,
) *ReviewEnrichmentModel {
	// Apply defaults for nil dimensions.
	termWidth := 80
	termHeight := 40
	if dims != nil {
		if dims.TerminalWidth > 0 {
			termWidth = dims.TerminalWidth
		}
		if dims.TerminalHeight > 0 {
			termHeight = dims.TerminalHeight
		}
	}

	// Create a copy of the event for reverting.
	eventCopy := *event

	// Create tag and category selectors to get available options.
	tagSelector := selectors.NewTagSelector()
	tagSelector.SetSelectedTags(event.Tags)
	availableTags := tagSelector.AvailableTags()

	categorySelector := selectors.NewCategorySelector()
	if err := categorySelector.SetSelected(event.Categories); err != nil {
		// Existing event categories should be valid, ignore.
		_ = err
	}
	availableCategories := categorySelector.AvailableCategories()

	// Load all available skills from repository.
	allSkills, err := service.GetSkillRepository().List(ctx, nil)
	if err != nil {
		allSkills = []*career.Skill{}
	}

	skillSelector := selectors.NewSkillSelector(allSkills)
	// Pre-select skills from event.
	for _, skillID := range event.Skills {
		if err := skillSelector.SelectSkill(skillID); err != nil {
			// Event skills should be valid, ignore.
			_ = err
		}
	}
	availableSkills := skillSelector.AvailableSkills()

	// Extract form data from event.
	formData := forms.GetMetadataFormData(event)

	// Calculate form dimensions for the modal overlay context.
	// The modal has a fixed width of 80 with border+padding chrome.
	const modalWidth = 80
	formWidth := forms.ModalFormWidth(modalWidth)
	formHeight := forms.ModalFormHeight(termHeight)

	form := forms.NewMetadataForm(
		formData, forms.MetadataFormConfig{
			AvailableTags:       availableTags,
			AvailableCategories: availableCategories,
			AvailableSkills:     availableSkills,
			Width:               formWidth,
			Height:              formHeight,
		},
	)

	return &ReviewEnrichmentModel{
		EditorFields: forms.EditorFields{
			Form:   form,
			Width:  termWidth,
			Height: termHeight,
			Theme:  themes.NewDefaultTheme(),
		},
		event:            event,
		originalEvent:    &eventCopy,
		service:          service,
		cliService:       cliSvc,
		ctx:              ctx,
		formData:         formData,
		tagSelector:      tagSelector,
		categorySelector: categorySelector,
		skillSelector:    skillSelector,
	}
}

// Init initializes the model
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *ReviewEnrichmentModel) Init() tea.Cmd {
	return m.Form.Init()
}

// Update handles messages.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - tea.Model: the updated model.
//   - tea.Cmd: command to execute.
//
// Side effects:
//   - May update internal state based on message type.
func (m *ReviewEnrichmentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return forms.EditorUpdate(&m.EditorFields, m, msg, m.handleFormCompletion, func() tea.Msg { return QuitMsg{} })
}

// handleFormCompletion processes the completed form and saves the metadata.
func (m *ReviewEnrichmentModel) handleFormCompletion() (tea.Model, tea.Cmd) {
	// Check if user confirmed via the submit button
	// If they selected "Cancel" on the confirm, treat as cancelled
	if !m.formData.SubmitConfirmed {
		m.Cancelled = true
		return m, nil
	}

	// Apply form data to event
	err := forms.ApplyMetadataFormData(m.event, m.formData)
	if err != nil {
		m.Err = fmt.Errorf("failed to apply form data: %w", err)
		return m, nil
	}

	// Validate event
	if err := m.event.Validate(); err != nil {
		m.Err = fmt.Errorf("validation failed: %w", err)
		return m, nil
	}

	// Update timestamp
	m.event.UpdatedAt = time.Now()

	// Persist metadata changes to service
	if m.cliService != nil {
		if err := m.cliService.UpdateEventMetadata(m.ctx, m.event); err != nil {
			m.Err = fmt.Errorf("failed to save metadata: %w", err)
			return m, nil
		}
	}

	m.submitted = true
	return m, nil
}

// ExtractInput returns the current form data as a pure domain MetadataInput.
// This bridges the TUI form data to the domain layer without any Huh dependencies.
//
// Returns:
//   - A capture.MetadataInput struct.
//
// Side effects:
//   - None.
func (m *ReviewEnrichmentModel) ExtractInput() capture.MetadataInput {
	return capture.MetadataInput{
		Date:       m.formData.Date,
		Company:    m.formData.Company,
		Project:    m.formData.Project,
		Tags:       m.formData.Tags,
		Categories: m.formData.Categories,
		Skills:     m.formData.Skills,
	}
}

// GetEvent returns the edited event
//
// Returns:
//   - A fully initialized career.Event ready for use.
//
// Side effects:
//   - None.
func (m *ReviewEnrichmentModel) GetEvent() *career.Event {
	return m.event
}

// IsSubmitted returns true if changes were saved
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *ReviewEnrichmentModel) IsSubmitted() bool {
	return m.submitted
}

// IsCancelled returns true if operation was cancelled
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *ReviewEnrichmentModel) IsCancelled() bool {
	return m.Cancelled
}

// Revert reverts changes to the original event
//
// Side effects:
//   - None.
func (m *ReviewEnrichmentModel) Revert() {
	*m.event = *m.originalEvent
	// Update form data
	m.formData = forms.GetMetadataFormData(m.originalEvent)
}

// GetError returns the current error
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (m *ReviewEnrichmentModel) GetError() error {
	return m.Err
}

// GetTitle returns the modal title for overlay rendering.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *ReviewEnrichmentModel) GetTitle() string {
	return "Edit Event Metadata"
}

// GetContent returns just the form content without header/footer.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *ReviewEnrichmentModel) GetContent() string {
	formView := m.Form.View()

	// Add error if present
	if m.Err != nil {
		errorColor := m.Theme.ErrorColor()
		errorStyle := lipgloss.NewStyle().
			Foreground(errorColor).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(errorColor).
			Padding(1, 2).
			MarginTop(1)

		formView += "\n\n" + errorStyle.Render(m.Err.Error())
	}

	return formView
}

// GetFooter returns the footer instructions for the modal.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *ReviewEnrichmentModel) GetFooter() string {
	return "Enter: Confirm | Esc: Cancel | Tab: Next Field | Shift+Tab: Previous"
}

// View renders the editor UI
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *ReviewEnrichmentModel) View() string {
	return forms.RenderEditorView(&m.EditorFields, "Edit Event Metadata", "review_enrichment")
}
