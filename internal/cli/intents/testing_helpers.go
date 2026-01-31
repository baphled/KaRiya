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
// fields initialized to their zero values and an empty message slice.
//
// Returns:
//   - A ready-to-use MockIntent instance.
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
//   - Always nil; no commands are dispatched.
//
// Side effects:
//   - Sets InitCalled to true.
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
// and recording that the view was rendered.
//
// Returns:
//   - The static string "Mock Intent View".
//
// Side effects:
//   - Sets ViewCalled to true.
func (m *MockIntent) View() string {
	m.ViewCalled = true
	return "Mock Intent View"
}

// Result implements Intent.Result by exposing the pre-configured result
// value that was set via SetResult.
//
// Returns:
//   - The stored IntentResult, or nil if none was configured.
//
// Side effects:
//   - None.
func (m *MockIntent) Result() *IntentResult[interface{}] {
	return m.resultValue
}

// SetResult configures the result that this mock intent will return
// from subsequent calls to Result.
//
// Expected:
//   - result is the IntentResult to store; may be nil to clear.
//
// Side effects:
//   - Replaces the previously stored result value.
func (m *MockIntent) SetResult(result *IntentResult[interface{}]) {
	m.resultValue = result
}

// ThemeAwareMockIntent is a mock intent that supports theme management.
type ThemeAwareMockIntent struct {
	*MockIntent
	themeManager *themes.ThemeManager
}

// NewThemeAwareMockIntent creates a new ThemeAwareMockIntent that embeds a
// fresh MockIntent and supports theme manager injection for testing.
//
// Returns:
//   - A ready-to-use ThemeAwareMockIntent instance.
//
// Side effects:
//   - None.
func NewThemeAwareMockIntent() *ThemeAwareMockIntent {
	return &ThemeAwareMockIntent{
		MockIntent: NewMockIntent(),
	}
}

// SetThemeManager implements ThemeAware.SetThemeManager by storing the
// provided theme manager for later retrieval.
//
// Expected:
//   - tm is the ThemeManager to associate with this intent; may be nil.
//
// Side effects:
//   - Replaces the previously stored theme manager reference.
func (m *ThemeAwareMockIntent) SetThemeManager(tm *themes.ThemeManager) {
	m.themeManager = tm
}

// GetThemeManager implements ThemeAware.GetThemeManager by exposing the
// theme manager that was previously injected via SetThemeManager.
//
// Returns:
//   - The stored ThemeManager, or nil if none was set.
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
// embeds a fresh MockIntent and initializes the selection index to zero.
//
// Returns:
//   - A ready-to-use MockIntentWithSelection instance.
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
// assertions against expected navigation state.
//
// Returns:
//   - The zero-based index of the currently selected item.
//
// Side effects:
//   - None.
func (m *MockIntentWithSelection) GetSelectedIndex() int {
	return m.SelectedIndex
}

// SetSelectedIndex configures the selection position to simulate
// navigation state changes in tests.
//
// Expected:
//   - idx is the zero-based index to select.
//
// Side effects:
//   - Replaces the previously stored selection index.
func (m *MockIntentWithSelection) SetSelectedIndex(idx int) {
	m.SelectedIndex = idx
}
