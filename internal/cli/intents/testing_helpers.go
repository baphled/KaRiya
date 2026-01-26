package intents

import (
	"github.com/baphled/kariya/internal/cli/themes"
	tea "github.com/charmbracelet/bubbletea"
)

// MockIntent is a mock implementation of Intent for testing.
type MockIntent struct {
	InitCalled   bool
	UpdateCalled int
	ViewCalled   bool
	resultValue  *IntentResult[interface{}]
	Messages     []tea.Msg
}

// NewMockIntent creates a new MockIntent for testing.
func NewMockIntent() *MockIntent {
	return &MockIntent{
		Messages: make([]tea.Msg, 0),
	}
}

// Init implements Intent.Init.
func (m *MockIntent) Init() tea.Cmd {
	m.InitCalled = true
	return nil
}

// Update implements Intent.Update.
func (m *MockIntent) Update(msg tea.Msg) tea.Cmd {
	m.UpdateCalled++
	m.Messages = append(m.Messages, msg)
	return nil
}

// View implements Intent.View.
func (m *MockIntent) View() string {
	m.ViewCalled = true
	return "Mock Intent View"
}

// Result implements Intent.Result.
func (m *MockIntent) Result() *IntentResult[interface{}] {
	return m.resultValue
}

// SetResult sets the result for this mock intent.
func (m *MockIntent) SetResult(result *IntentResult[interface{}]) {
	m.resultValue = result
}

// ThemeAwareMockIntent is a mock intent that supports theme management.
type ThemeAwareMockIntent struct {
	*MockIntent
	themeManager *themes.ThemeManager
}

// NewThemeAwareMockIntent creates a new ThemeAwareMockIntent for testing.
func NewThemeAwareMockIntent() *ThemeAwareMockIntent {
	return &ThemeAwareMockIntent{
		MockIntent: NewMockIntent(),
	}
}

// SetThemeManager implements ThemeAware.SetThemeManager.
func (m *ThemeAwareMockIntent) SetThemeManager(tm *themes.ThemeManager) {
	m.themeManager = tm
}

// GetThemeManager implements ThemeAware.GetThemeManager.
func (m *ThemeAwareMockIntent) GetThemeManager() *themes.ThemeManager {
	return m.themeManager
}

// MockIntentWithSelection extends MockIntent with selection state for testing.
type MockIntentWithSelection struct {
	*MockIntent
	SelectedIndex int
}

// NewMockIntentWithSelection creates a new MockIntentWithSelection for testing.
func NewMockIntentWithSelection() *MockIntentWithSelection {
	return &MockIntentWithSelection{
		MockIntent:    NewMockIntent(),
		SelectedIndex: 0,
	}
}

// GetSelectedIndex returns the selected index.
func (m *MockIntentWithSelection) GetSelectedIndex() int {
	return m.SelectedIndex
}

// SetSelectedIndex sets the selected index.
func (m *MockIntentWithSelection) SetSelectedIndex(idx int) {
	m.SelectedIndex = idx
}
