package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Focus indices for burst editor navigation
const (
	BurstNameFieldIdx = iota
	BurstDescriptionFieldIdx
	BurstCompetencyFocusFieldIdx
	BurstSaveButtonIdx
	BurstCancelButtonIdx
	BurstEditorFieldCount // Total number of focus positions
)

// BurstEditorModel represents the burst editor form state
type BurstEditorModel struct {
	*BaseStandardModel
	burst                *career.Burst
	originalBurst        *career.Burst // For reverting changes
	service              *careerservice.Service
	ctx                  context.Context
	nameInput            textinput.Model
	descriptionInput     textinput.Model
	competencyFocusInput textinput.Model
	focusIndex           int
	err                  error
	submitted            bool
	cancelled            bool
	fieldErrors          map[int]string
	width                int
	height               int
	helpFooter           components.HelpFooterModel // Help footer for keyboard shortcuts
}

// NewBurstEditorModel creates a new burst editor model
func NewBurstEditorModel(burst *career.Burst, service *careerservice.Service, ctx context.Context) *BurstEditorModel {
	// Create a copy of the burst for reverting
	burstCopy := *burst

	// Create name input
	nameInput := textinput.New()
	nameInput.Placeholder = "Burst name"
	nameInput.SetValue(burst.Name)
	nameInput.Width = 70
	nameInput.Focus()

	// Create description input
	descriptionInput := textinput.New()
	descriptionInput.Placeholder = "Burst description (optional)"
	descriptionInput.SetValue(burst.Description)
	descriptionInput.Width = 70

	// Create competency focus input
	competencyFocusInput := textinput.New()
	competencyFocusInput.Placeholder = "Primary competency focus"
	competencyFocusInput.SetValue(burst.CompetencyFocus)
	competencyFocusInput.Width = 70

	return &BurstEditorModel{
		BaseStandardModel:    NewBaseStandardModel(),
		burst:                burst,
		originalBurst:        &burstCopy,
		service:              service,
		ctx:                  ctx,
		nameInput:            nameInput,
		descriptionInput:     descriptionInput,
		competencyFocusInput: competencyFocusInput,
		focusIndex:           0,
		err:                  nil,
		submitted:            false,
		cancelled:            false,
		fieldErrors:          make(map[int]string),
		helpFooter:           components.NewHelpFooter("burst_editor", 80),
	}
}

// Init initializes the model
func (m *BurstEditorModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *BurstEditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	// Update active text input
	switch m.focusIndex {
	case BurstNameFieldIdx:
		var cmd tea.Cmd
		m.nameInput, cmd = m.nameInput.Update(msg)
		m.burst.Name = m.nameInput.Value()
		return m, cmd
	case BurstDescriptionFieldIdx:
		var cmd tea.Cmd
		m.descriptionInput, cmd = m.descriptionInput.Update(msg)
		m.burst.Description = m.descriptionInput.Value()
		return m, cmd
	case BurstCompetencyFocusFieldIdx:
		var cmd tea.Cmd
		m.competencyFocusInput, cmd = m.competencyFocusInput.Update(msg)
		m.burst.CompetencyFocus = m.competencyFocusInput.Value()
		return m, cmd
	}

	return m, nil
}

func (m *BurstEditorModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		// Allow quit from burst editor
		return m, func() tea.Msg { return QuitMsg{} }

	case "esc":
		// Cancel editor and return to parent
		m.cancelled = true
		return m, nil

	case "tab":
		m.clearFieldErrors()
		m.focusIndex = (m.focusIndex + 1) % BurstEditorFieldCount
		m.updateInputFocus()
		return m, nil

	case "shift+tab":
		m.clearFieldErrors()
		m.focusIndex = (m.focusIndex - 1 + BurstEditorFieldCount) % BurstEditorFieldCount
		m.updateInputFocus()
		return m, nil

	case "enter":
		if m.focusIndex == BurstSaveButtonIdx {
			return m.submitBurst()
		} else if m.focusIndex == BurstCancelButtonIdx {
			m.cancelled = true
			return m, nil
		}
		return m, nil
	}

	return m, nil
}

func (m *BurstEditorModel) updateInputFocus() {
	m.nameInput.Blur()
	m.descriptionInput.Blur()
	m.competencyFocusInput.Blur()

	switch m.focusIndex {
	case BurstNameFieldIdx:
		m.nameInput.Focus()
	case BurstDescriptionFieldIdx:
		m.descriptionInput.Focus()
	case BurstCompetencyFocusFieldIdx:
		m.competencyFocusInput.Focus()
	}
}

func (m *BurstEditorModel) clearFieldErrors() {
	m.fieldErrors = make(map[int]string)
}

func (m *BurstEditorModel) submitBurst() (tea.Model, tea.Cmd) {
	m.clearFieldErrors()

	// Validate name
	if strings.TrimSpace(m.burst.Name) == "" {
		m.fieldErrors[BurstNameFieldIdx] = "Burst name cannot be empty"
		m.err = fmt.Errorf("validation error: burst name is required")
		return m, nil
	}

	if len(m.burst.Name) > 200 {
		m.fieldErrors[BurstNameFieldIdx] = "Burst name cannot exceed 200 characters"
		m.err = fmt.Errorf("validation error: burst name exceeds 200 characters")
		return m, nil
	}

	// Validate description (optional, but if provided check length)
	if len(m.burst.Description) > 1000 {
		m.fieldErrors[BurstDescriptionFieldIdx] = "Description cannot exceed 1000 characters"
		m.err = fmt.Errorf("validation error: description exceeds 1000 characters")
		return m, nil
	}

	// Validate competency focus (optional)
	if len(m.burst.CompetencyFocus) > 100 {
		m.fieldErrors[BurstCompetencyFocusFieldIdx] = "Competency focus cannot exceed 100 characters"
		m.err = fmt.Errorf("validation error: competency focus exceeds 100 characters")
		return m, nil
	}

	// Update timestamp
	m.burst.UpdatedAt = time.Now()

	// Save the burst using the repository
	burstRepo := m.service.GetBurstRepository()
	if burstRepo == nil {
		m.err = fmt.Errorf("burst repository not available")
		return m, nil
	}
	
	err := burstRepo.Update(m.ctx, m.burst)
	if err != nil {
		m.err = fmt.Errorf("failed to save burst: %w", err)
		return m, nil
	}

	m.submitted = true
	return m, nil
}

// GetBurst returns the edited burst
func (m *BurstEditorModel) GetBurst() *career.Burst {
	return m.burst
}

// IsSubmitted returns true if changes were saved
func (m *BurstEditorModel) IsSubmitted() bool {
	return m.submitted
}

// IsCancelled returns true if operation was cancelled
func (m *BurstEditorModel) IsCancelled() bool {
	return m.cancelled
}

// Revert reverts changes to the original burst
func (m *BurstEditorModel) Revert() {
	m.burst.Name = m.originalBurst.Name
	m.burst.Description = m.originalBurst.Description
	m.burst.CompetencyFocus = m.originalBurst.CompetencyFocus
	m.nameInput.SetValue(m.originalBurst.Name)
	m.descriptionInput.SetValue(m.originalBurst.Description)
	m.competencyFocusInput.SetValue(m.originalBurst.CompetencyFocus)
}

// GetError returns the current error
func (m *BurstEditorModel) GetError() error {
	return m.err
}

// View renders the editor UI
func (m *BurstEditorModel) View() string {
	// Render form content
	formContent := m.renderFormContent()

	// Use header and footer components
	headerView := components.NewHeader("Burst Editor", m.width).View()
	footerView := components.NewFooter(m.width).View()

	// Render help footer
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

// renderFormContent renders all form fields
func (m *BurstEditorModel) renderFormContent() string {
	// Create form container with responsive layout
	formContainer := components.NewFormContainer().
		SetWidth(m.width).
		SetHeight(m.height).
		SetLayout(components.Responsive).
		SetPadding(2).
		SetVerticalSpacing(2)

	// Add name field
	nameFieldErr := ""
	if err, ok := m.fieldErrors[BurstNameFieldIdx]; ok {
		nameFieldErr = err
	}

	formContainer.AddField(components.FormField{
		Label:      "Burst Name (required):",
		Input:      m.nameInput.View(),
		Error:      nameFieldErr,
		IsFocused:  m.focusIndex == BurstNameFieldIdx,
		IsRequired: true,
		FullWidth:  true,
	})

	// Add description field
	descriptionFieldErr := ""
	if err, ok := m.fieldErrors[BurstDescriptionFieldIdx]; ok {
		descriptionFieldErr = err
	}

	formContainer.AddField(components.FormField{
		Label:     "Description (optional):",
		Input:     m.descriptionInput.View(),
		Error:     descriptionFieldErr,
		IsFocused: m.focusIndex == BurstDescriptionFieldIdx,
		FullWidth: true,
	})

	// Add competency focus field
	competencyFieldErr := ""
	if err, ok := m.fieldErrors[BurstCompetencyFocusFieldIdx]; ok {
		competencyFieldErr = err
	}

	formContainer.AddField(components.FormField{
		Label:     "Competency Focus (optional):",
		Input:     m.competencyFocusInput.View(),
		Error:     competencyFieldErr,
		IsFocused: m.focusIndex == BurstCompetencyFocusFieldIdx,
		FullWidth: true,
	})

	// Add buttons field
	saveBtn := "[ Save ]"
	cancelBtn := "[ Cancel ]"

	if m.focusIndex == BurstSaveButtonIdx {
		saveBtn = styles.ButtonPrimaryFocused.Render(saveBtn)
		cancelBtn = styles.ButtonSecondary.Render(cancelBtn)
	} else if m.focusIndex == BurstCancelButtonIdx {
		saveBtn = styles.ButtonPrimary.Render(saveBtn)
		cancelBtn = styles.ButtonSecondaryFocused.Render(cancelBtn)
	} else {
		saveBtn = styles.ButtonPrimary.Render(saveBtn)
		cancelBtn = styles.ButtonSecondary.Render(cancelBtn)
	}

	buttonsStr := strings.Join([]string{saveBtn, cancelBtn}, "  ")

	formContainer.AddField(components.FormField{
		Input:     buttonsStr,
		IsFocused: m.focusIndex >= BurstSaveButtonIdx,
		FullWidth: true,
	})

	// Get rendered form
	formContent := formContainer.Render()

	// Add model-level error if present
	if m.err != nil {
		formContent += "\n\n" + styles.ErrorBox.Render(m.err.Error())
	}

	return formContent
}

