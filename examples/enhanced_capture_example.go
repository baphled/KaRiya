package intents

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
)

// EnhancedCaptureEventModel demonstrates best practices for using lipgloss and bubbles
// in the KaRiya TUI. This model shows:
// - Centralized style usage
// - Bubbles components for interactivity
// - Responsive layout
// - Proper focus management
// - Component composition
type EnhancedCaptureEventModel struct {
	// State management
	state FormState

	// Bubbles interactive components
	titleInput       textinput.Model
	descriptionInput textinput.Model
	companyInput     textinput.Model
	projectInput     textinput.Model

	// UI state
	focusedField int
	width        int
	height       int
	errorMsg     string

	// Form data
	formData *FormData
}

// FormState represents the current form state
type FormState string

const (
	FormStateInput  FormState = "input"
	FormStateReview FormState = "review"
	FormStateSubmit FormState = "submit"
)

// FormData holds the captured form data
type FormData struct {
	Title       string
	Description string
	Company     string
	Project     string
}

// NewEnhancedCaptureEventModel creates and initializes the model
func NewEnhancedCaptureEventModel() *EnhancedCaptureEventModel {
	m := &EnhancedCaptureEventModel{
		state:        FormStateInput,
		focusedField: 0,
		formData:     &FormData{},
	}

	// Initialize all text input fields with bubbles
	m.initializeInputFields()

	return m
}

// initializeInputFields sets up all text input bubbles components
func (m *EnhancedCaptureEventModel) initializeInputFields() {
	// Title input
	m.titleInput = textinput.New()
	m.titleInput.Placeholder = "e.g., Led API redesign project"
	m.titleInput.Focus()
	m.titleInput.PromptStyle = lipgloss.NewStyle().Foreground(styles.ColorAccentTeal)
	m.titleInput.TextStyle = lipgloss.NewStyle().Foreground(styles.ColorTextPrimary)

	// Description input
	m.descriptionInput = textinput.New()
	m.descriptionInput.Placeholder = "e.g., Redesigned REST API for 3x performance improvement"
	m.descriptionInput.PromptStyle = lipgloss.NewStyle().Foreground(styles.ColorAccentTeal)
	m.descriptionInput.TextStyle = lipgloss.NewStyle().Foreground(styles.ColorTextPrimary)

	// Company input
	m.companyInput = textinput.New()
	m.companyInput.Placeholder = "e.g., Acme Corp"
	m.companyInput.PromptStyle = lipgloss.NewStyle().Foreground(styles.ColorAccentTeal)
	m.companyInput.TextStyle = lipgloss.NewStyle().Foreground(styles.ColorTextPrimary)

	// Project input
	m.projectInput = textinput.New()
	m.projectInput.Placeholder = "e.g., Backend Infrastructure"
	m.projectInput.PromptStyle = lipgloss.NewStyle().Foreground(styles.ColorAccentTeal)
	m.projectInput.TextStyle = lipgloss.NewStyle().Foreground(styles.ColorTextPrimary)
}

// Init initializes the model and returns an initial command
func (m *EnhancedCaptureEventModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles incoming messages
func (m *EnhancedCaptureEventModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		return m.handleKeyInput(msg)
	}

	return nil
}

// handleKeyInput processes keyboard input
func (m *EnhancedCaptureEventModel) handleKeyInput(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c":
		return tea.Quit

	case "tab":
		m.focusedField = (m.focusedField + 1) % 4
		m.updateFieldFocus()

	case "shift+tab":
		m.focusedField = (m.focusedField - 1 + 4) % 4
		m.updateFieldFocus()

	case "enter":
		if m.state == FormStateInput {
			if m.validateForm() {
				m.syncFormData()
				m.state = FormStateReview
			}
		} else if m.state == FormStateReview {
			m.state = FormStateSubmit
		}

	case "esc":
		if m.state == FormStateReview {
			m.state = FormStateInput
		}
	}

	// Delegate input to the focused field
	return m.updateFocusedField(msg)
}

// updateFieldFocus manages which field is currently focused
func (m *EnhancedCaptureEventModel) updateFieldFocus() {
	// Blur all fields
	m.titleInput.Blur()
	m.descriptionInput.Blur()
	m.companyInput.Blur()
	m.projectInput.Blur()

	// Focus the selected field
	switch m.focusedField {
	case 0:
		m.titleInput.Focus()
	case 1:
		m.descriptionInput.Focus()
	case 2:
		m.companyInput.Focus()
	case 3:
		m.projectInput.Focus()
	}
}

// updateFocusedField delegates input to the currently focused field
func (m *EnhancedCaptureEventModel) updateFocusedField(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd

	switch m.focusedField {
	case 0:
		m.titleInput, cmd = m.titleInput.Update(msg)
	case 1:
		m.descriptionInput, cmd = m.descriptionInput.Update(msg)
	case 2:
		m.companyInput, cmd = m.companyInput.Update(msg)
	case 3:
		m.projectInput, cmd = m.projectInput.Update(msg)
	}

	return cmd
}

// validateForm checks if the form has required fields filled
func (m *EnhancedCaptureEventModel) validateForm() bool {
	if m.titleInput.Value() == "" {
		m.errorMsg = "Title is required"
		return false
	}
	if m.companyInput.Value() == "" {
		m.errorMsg = "Company is required"
		return false
	}
	m.errorMsg = ""
	return true
}

// syncFormData copies input values to the form data struct
func (m *EnhancedCaptureEventModel) syncFormData() {
	m.formData.Title = m.titleInput.Value()
	m.formData.Description = m.descriptionInput.Value()
	m.formData.Company = m.companyInput.Value()
	m.formData.Project = m.projectInput.Value()
}

// View renders the current state
func (m *EnhancedCaptureEventModel) View() string {
	switch m.state {
	case FormStateInput:
		return m.viewInputForm()
	case FormStateReview:
		return m.viewReview()
	case FormStateSubmit:
		return m.viewSuccess()
	default:
		return ""
	}
}

// viewInputForm renders the input form with styled fields
func (m *EnhancedCaptureEventModel) viewInputForm() string {
	// Header with emoji and styling
	headerText := "📝 Capture Career Event"
	header := styles.CardHeader.Copy().
		Foreground(styles.ColorTextPrimary).
		MarginBottom(1).
		Render(headerText)

	// Form fields with visual feedback
	titleField := m.renderFormField("Title", m.titleInput.View(), m.focusedField == 0)
	descField := m.renderFormField("Description", m.descriptionInput.View(), m.focusedField == 1)
	companyField := m.renderFormField("Company", m.companyInput.View(), m.focusedField == 2)
	projectField := m.renderFormField("Project", m.projectInput.View(), m.focusedField == 3)

	// Combine form fields vertically
	formContent := lipgloss.JoinVertical(
		lipgloss.Left,
		titleField,
		descField,
		companyField,
		projectField,
	)

	// Error message if validation failed
	var errorDisplay string
	if m.errorMsg != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(styles.ColorError).
			MarginTop(1)
		errorDisplay = errorStyle.Render("✗ " + m.errorMsg)
		formContent = lipgloss.JoinVertical(lipgloss.Left, formContent, errorDisplay)
	}

	// Help footer with keyboard shortcuts
	footerText := "Tab: Next • Shift+Tab: Prev • Enter: Review • Ctrl+C: Exit"
	footer := styles.CardFooter.Copy().
		Foreground(styles.ColorTextSecondary).
		MarginTop(2).
		Render(footerText)

	// Use CardContainer for consistent styling
	card := components.NewCardContainer().
		SetHeader(header).
		SetBody(formContent).
		SetFooter(footer).
		Render()

	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(card)
}

// renderFormField creates a styled form field with label and input
func (m *EnhancedCaptureEventModel) renderFormField(label, input string, focused bool) string {
	// Label styling - consistent width for alignment
	labelStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextPrimary).
		Width(13).
		Bold(true)

	// Input wrapper - highlight border when focused
	borderColor := styles.ColorBorder
	if focused {
		borderColor = styles.ColorBorderActive
	}

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1).
		MarginBottom(1)

	// Combine label and input horizontally
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		labelStyle.Render(label),
		inputStyle.Render(input),
	)
}

// viewReview renders the review screen showing captured data
func (m *EnhancedCaptureEventModel) viewReview() string {
	// Header
	header := styles.CardHeader.Copy().
		Foreground(styles.ColorTextPrimary).
		MarginBottom(1).
		Render("✓ Review Event")

	// Display captured data with consistent styling
	dataStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextPrimary).
		MarginBottom(1)

	dataLines := []string{
		fmt.Sprintf("Title:       %s", m.formData.Title),
		fmt.Sprintf("Description: %s", m.formData.Description),
		fmt.Sprintf("Company:     %s", m.formData.Company),
		fmt.Sprintf("Project:     %s", m.formData.Project),
	}

	data := dataStyle.Render(lipgloss.JoinVertical(lipgloss.Left, dataLines...))

	// Action buttons with consistent styling
	confirmBtn := styles.ButtonPrimary.Render("✓ Confirm")
	editBtn := styles.ButtonSecondary.Render("✎ Edit")

	actions := lipgloss.JoinHorizontal(
		lipgloss.Top,
		confirmBtn,
		editBtn,
	)

	footer := styles.CardFooter.Copy().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1).
		Render("Enter: Confirm • Esc: Edit")

	// Use CardContainer for consistent styling
	card := components.NewCardContainer().
		SetHeader(header).
		SetBody(lipgloss.JoinVertical(lipgloss.Left, data, "", actions)).
		SetFooter(footer).
		Render()

	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(card)
}

// viewSuccess renders the success confirmation screen
func (m *EnhancedCaptureEventModel) viewSuccess() string {
	successMsg := "✓ Event captured successfully!"
	successStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSuccess).
		Bold(true)

	message := successStyle.Render(successMsg)

	footer := styles.CardFooter.Copy().
		Foreground(styles.ColorTextSecondary).
		Render("Press any key to continue...")

	card := components.NewCardContainer().
		SetHeader("Success").
		SetBody(message).
		SetFooter(footer).
		Render()

	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(card)
}

// ============================================================================
// Key Patterns Demonstrated
// ============================================================================

// Pattern 1: Centralized Style Usage
// All styles come from internal/cli/styles package, not defined inline
// This makes it easy to maintain consistent theming

// Pattern 2: Bubbles Component Integration
// Uses textinput.Model for each form field, delegating keyboard handling
// Bubbles components handle cursor, selection, and text manipulation

// Pattern 3: Focus Management
// Clear updateFieldFocus() method manages which component is active
// Each field is properly blurred/focused based on user navigation

// Pattern 4: Component Composition
// Uses CardContainer component for consistent card styling
// Avoids repeating border, padding, and layout logic

// Pattern 5: Responsive Layout
// Uses lipgloss.JoinVertical and JoinHorizontal for flexible layouts
// Can easily adapt to different terminal sizes

// Pattern 6: Error Handling
// Validation errors displayed inline with error color
// User gets immediate feedback on form issues

// Pattern 7: State-Based Rendering
// Different views for input, review, and success states
// Clear state transitions with keyboard commands

// ============================================================================
// Usage Example
// ============================================================================

// To use this model in your Bubble Tea application:
//
// model := NewEnhancedCaptureEventModel()
// p := tea.NewProgram(model)
// if _, err := p.Run(); err != nil {
//     log.Fatal(err)
// }

