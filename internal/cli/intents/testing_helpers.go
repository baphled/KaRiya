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

// NewMockIntent creates a new MockIntent for testing with all tracking
//
// Returns:
//   - A fully initialized MockIntent ready for use.
//
// Side effects:
//   - None.
func NewMockIntent() *MockIntent {
	return &MockIntent{
		Messages: make([]tea.Msg, 0),
	}
}

// Init implements Intent.Init by recording that initialization was invoked.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *MockIntent) Init() tea.Cmd {
	m.InitCalled = true
	return nil
}

// Update implements Intent.Update by recording the received message for
// later assertion in tests.
//
// Expected:
//   - msg is the Bubble Tea message to record.
//
// Returns:
//   - Always nil; no commands are dispatched.
//
// Side effects:
//   - Increments UpdateCalled counter.
//   - Appends msg to the Messages slice.
func (m *MockIntent) Update(msg tea.Msg) tea.Cmd {
	m.UpdateCalled++
	m.Messages = append(m.Messages, msg)
	return nil
}

// View implements Intent.View by producing a fixed placeholder string
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *MockIntent) View() string {
	m.ViewCalled = true
	return "Mock Intent View"
}

// Result implements Intent.Result by exposing the pre-configured result
//
// Returns:
//   - A fully initialized IntentResult[interface{}] ready for use.
//
// Side effects:
//   - None.
func (m *MockIntent) Result() *IntentResult[interface{}] {
	return m.resultValue
}

// SetResult configures the result that this mock intent will return
//
// Expected:
//   - intentresult[interface{}] must be valid.
//
// Side effects:
//   - None.
func (m *MockIntent) SetResult(result *IntentResult[interface{}]) {
	m.resultValue = result
}

// ThemeAwareMockIntent is a mock intent that supports theme management.
type ThemeAwareMockIntent struct {
	*MockIntent
	themeManager *themes.ThemeManager
}

// NewThemeAwareMockIntent creates a new ThemeAwareMockIntent that embeds a
//
// Returns:
//   - A fully initialized ThemeAwareMockIntent ready for use.
//
// Side effects:
//   - None.
func NewThemeAwareMockIntent() *ThemeAwareMockIntent {
	return &ThemeAwareMockIntent{
		MockIntent: NewMockIntent(),
	}
}

// SetThemeManager implements ThemeAware.SetThemeManager by storing the
//
// Expected:
//   - th must be a valid theme instance (can be nil).
//
// Side effects:
//   - None.
func (m *ThemeAwareMockIntent) SetThemeManager(tm *themes.ThemeManager) {
	m.themeManager = tm
}

// GetThemeManager implements ThemeAware.GetThemeManager by exposing the
//
// Returns:
//   - A fully initialized themes.ThemeManager ready for use.
//
// Side effects:
//   - None.
func (m *ThemeAwareMockIntent) GetThemeManager() *themes.ThemeManager {
	return m.themeManager
}

// MockIntentWithSelection extends MockIntent with selection state for testing.
type MockIntentWithSelection struct {
	*MockIntent
	SelectedIndex int
}

// NewMockIntentWithSelection creates a new MockIntentWithSelection that
//
// Returns:
//   - A fully initialized MockIntentWithSelection ready for use.
//
// Side effects:
//   - None.
func NewMockIntentWithSelection() *MockIntentWithSelection {
	return &MockIntentWithSelection{
		MockIntent:    NewMockIntent(),
		SelectedIndex: 0,
	}
}

// GetSelectedIndex provides the current selection position for test
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (m *MockIntentWithSelection) GetSelectedIndex() int {
	return m.SelectedIndex
}

// SetSelectedIndex configures the selection position to simulate
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
func (m *MockIntentWithSelection) SetSelectedIndex(idx int) {
	m.SelectedIndex = idx
}
