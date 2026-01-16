package primitives

import (
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Input is a theme-aware text input component wrapping bubbles textinput.
// It provides a fluent API for configuration and themed rendering.
//
// Example:
//
//	input := primitives.NewInput(theme).
//	    Label("Email").
//	    Placeholder("user@example.com").
//	    Width(50)
//	input.Focus()
//	value := input.GetValue()
type Input struct {
	theme.Aware
	textInput   textinput.Model
	label       string
	errorMsg    string
	width       int
	labelStyle  lipgloss.Style
	errorStyle  lipgloss.Style
	borderStyle lipgloss.Style
}

// NewInput creates a new input component with the given theme.
// If theme is nil, the default theme is used.
func NewInput(th theme.Theme) *Input {
	ti := textinput.New()
	ti.CharLimit = 256

	i := &Input{
		textInput: ti,
		label:     "",
		errorMsg:  "",
		width:     40, // Default width
	}

	if th != nil {
		i.SetTheme(th)
	}

	i.applyTheming()
	return i
}

// Label sets the label displayed above the input.
// Returns the input for method chaining.
func (i *Input) Label(label string) *Input {
	i.label = label
	return i
}

// Placeholder sets the placeholder text shown when input is empty.
// Returns the input for method chaining.
func (i *Input) Placeholder(placeholder string) *Input {
	i.textInput.Placeholder = placeholder
	return i
}

// Value sets the initial value of the input.
// Returns the input for method chaining.
func (i *Input) Value(value string) *Input {
	i.textInput.SetValue(value)
	return i
}

// Error sets an error message to display below the input.
// Pass empty string to clear the error.
// Returns the input for method chaining.
func (i *Input) Error(msg string) *Input {
	i.errorMsg = msg
	return i
}

// Width sets the width of the input field.
// Returns the input for method chaining.
func (i *Input) Width(w int) *Input {
	i.width = w
	i.textInput.Width = w
	return i
}

// Focus returns a command to focus the input.
func (i *Input) Focus() tea.Cmd {
	return i.textInput.Focus()
}

// Blur removes focus from the input.
func (i *Input) Blur() {
	i.textInput.Blur()
}

// GetValue returns the current value of the input.
func (i *Input) GetValue() string {
	return i.textInput.Value()
}

// Update handles input events and updates the internal state.
// Implements the tea.Model interface for Bubble Tea integration.
func (i *Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	i.textInput, cmd = i.textInput.Update(msg)
	return i, cmd
}

// View returns the rendered input as a string.
// This is an alias for Render() to support Bubble Tea's tea.Model interface.
func (i *Input) View() string {
	return i.Render()
}

// Init implements tea.Model interface. Returns nil as no initialization is needed.
func (i *Input) Init() tea.Cmd {
	return nil
}

// Render returns the styled input as a string.
// The input is rendered with:
// - Label above (if set)
// - The input field with themed border
// - Error message below (if set)
func (i *Input) Render() string {
	var parts []string

	// Render label if set
	if i.label != "" {
		labelText := i.labelStyle.Render(i.label)
		parts = append(parts, labelText)
	}

	// Render input field with border
	inputView := i.textInput.View()

	// Apply border based on focus state
	borderColor := i.BorderColor()
	if i.textInput.Focused() {
		borderColor = i.AccentColor()
	}
	if i.errorMsg != "" {
		borderColor = i.ErrorColor()
	}

	borderedInput := i.borderStyle.
		BorderForeground(borderColor).
		Render(inputView)
	parts = append(parts, borderedInput)

	// Render error message if set
	if i.errorMsg != "" {
		errorText := i.errorStyle.Render(i.errorMsg)
		parts = append(parts, errorText)
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// applyTheming applies theme colors to the input styles.
func (i *Input) applyTheming() {
	// Label styling
	i.labelStyle = lipgloss.NewStyle().
		Foreground(i.SecondaryColor()).
		Bold(true).
		MarginBottom(0)

	// Error styling
	i.errorStyle = lipgloss.NewStyle().
		Foreground(i.ErrorColor()).
		MarginTop(0)

	// Border styling
	i.borderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1)

	// Apply theme to textinput
	i.textInput.TextStyle = lipgloss.NewStyle().Foreground(i.Theme().ForegroundColor())
	i.textInput.PlaceholderStyle = lipgloss.NewStyle().Foreground(i.MutedColor())
	i.textInput.PromptStyle = lipgloss.NewStyle().Foreground(i.PrimaryColor())
	i.textInput.Cursor.Style = lipgloss.NewStyle().Foreground(i.AccentColor())
}
