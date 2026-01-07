package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/navigation"
	"github.com/baphled/kariya/internal/cli/styles"
)

// HelpFooterModel renders a footer with contextual help text showing available keyboard shortcuts
type HelpFooterModel struct {
	width   int
	height  int
	context string // Context for contextual help (form, list, metadata_review, etc.)
	keys    []navigation.NavigationKey
}

// NewHelpFooter creates a new help footer with the given context
// Context can be "form", "list", "metadata_review", "bulk_operations", "home", or "" for default
func NewHelpFooter(context string, width int) HelpFooterModel {
	return HelpFooterModel{
		width:   width,
		height:  1,
		context: context,
		keys:    []navigation.NavigationKey{},
	}
}

// NewHelpFooterWithKeys creates a help footer with specific keys
func NewHelpFooterWithKeys(keys []navigation.NavigationKey, width int) HelpFooterModel {
	return HelpFooterModel{
		width:   width,
		height:  1,
		context: "",
		keys:    keys,
	}
}

// Init implements the BubbleTea Model interface
func (m HelpFooterModel) Init() tea.Cmd {
	return nil
}

// Update implements the BubbleTea Model interface
// Help footer doesn't handle any messages - it's display-only
func (m HelpFooterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}
	return m, nil
}

// View implements the BubbleTea Model interface
func (m HelpFooterModel) View() string {
	if m.width <= 0 {
		return ""
	}

	// Determine which keys to display
	var helpText string
	if len(m.keys) > 0 {
		helpText = navigation.GetHelpTextCompact(m.keys, true)
	} else {
		helpText = navigation.GetContextualHelp(m.context)
	}

	// Apply styling
	styledText := styles.InfoHint.Render(helpText)

	// Truncate if necessary for narrow terminals
	if len(styledText) > m.width {
		// Try to create a shortened version
		if len(m.keys) > 0 {
			helpText = m.getTruncatedHelp(m.keys, m.width-4) // -4 for styling overhead
		} else {
			helpText = m.getTruncatedContextualHelp(m.context, m.width-4)
		}
		styledText = styles.InfoHint.Render(helpText)
	}

	return styledText
}

// SetWidth updates the footer width
func (m *HelpFooterModel) SetWidth(width int) {
	m.width = width
}

// SetContext changes the help context
func (m *HelpFooterModel) SetContext(context string) {
	m.context = context
}

// SetKeys sets specific keys to display
func (m *HelpFooterModel) SetKeys(keys []navigation.NavigationKey) {
	m.keys = keys
}

// GetHeight returns the height of the footer (always 1)
func (m HelpFooterModel) GetHeight() int {
	return m.height
}

// getTruncatedHelp returns a truncated version of help text for narrow terminals
func (m HelpFooterModel) getTruncatedHelp(keys []navigation.NavigationKey, maxWidth int) string {
	if len(keys) == 0 {
		return ""
	}

	// Try with fewer keys
	var truncated []navigation.NavigationKey
	helpText := ""

	// Start with the first 3 keys and add more as space allows
	for i := 0; i < len(keys) && i < 5; i++ {
		test := append(truncated, keys[i])
		testText := navigation.GetHelpTextCompact(test, true)
		if len(testText) <= maxWidth {
			truncated = test
			helpText = testText
		} else {
			break
		}
	}

	// If we still can't fit, just show minimal help
	if helpText == "" || len(helpText) > maxWidth {
		return "? Help"
	}

	return helpText
}

// getTruncatedContextualHelp returns truncated contextual help for narrow terminals
func (m HelpFooterModel) getTruncatedContextualHelp(context string, maxWidth int) string {
	contextKeys := map[string][]navigation.NavigationKey{
		"form": {
			navigation.KeyUp,
			navigation.KeyDown,
			navigation.KeySelect,
			navigation.KeyBack,
		},
		"list": {
			navigation.KeyUp,
			navigation.KeyDown,
			navigation.KeySelect,
			navigation.KeyEdit,
		},
		"metadata_review": {
			navigation.KeyUp,
			navigation.KeyDown,
			navigation.KeySelect,
			navigation.KeyEdit,
			navigation.KeyBulk,
		},
		"bulk_operations": {
			navigation.KeyUp,
			navigation.KeyDown,
			navigation.KeyToggle,
			navigation.KeySelect,
		},
		"home": {
			navigation.KeyCapture,
			navigation.KeyList,
			navigation.KeyHelp,
		},
	}

	keys, exists := contextKeys[context]
	if !exists {
		keys = contextKeys["list"]
	}

	return m.getTruncatedHelp(keys, maxWidth)
}

// Style returns the styled footer text
func (m HelpFooterModel) Style() lipgloss.Style {
	return styles.InfoHint
}

// RenderForWidth renders the footer constrained to a specific width
func (m HelpFooterModel) RenderForWidth(width int) string {
	m.width = width
	return m.View()
}

// WithBorder renders the footer with a border
func (m HelpFooterModel) WithBorder() string {
	view := m.View()
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(m.width).
		Render(view)
	return style
}

// WithPadding renders the footer with padding
func (m HelpFooterModel) WithPadding(vertical, horizontal int) string {
	view := m.View()
	style := lipgloss.NewStyle().
		Padding(vertical, horizontal).
		Render(view)
	return style
}
