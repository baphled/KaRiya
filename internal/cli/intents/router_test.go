package intents

import (
	"testing"

	"github.com/baphled/kariya/internal/cli/themes"
	tea "github.com/charmbracelet/bubbletea"
)

// MockIntent is a mock implementation of Intent for testing.
type MockIntent struct {
	initCalled   bool
	updateCalled int
	viewCalled   bool
	result       *IntentResult[interface{}]
	messages     []tea.Msg
}

func NewMockIntent() *MockIntent {
	return &MockIntent{
		messages: make([]tea.Msg, 0),
	}
}

func (m *MockIntent) Init() tea.Cmd {
	m.initCalled = true
	return nil
}

func (m *MockIntent) Update(msg tea.Msg) tea.Cmd {
	m.updateCalled++
	m.messages = append(m.messages, msg)
	return nil
}

func (m *MockIntent) View() string {
	m.viewCalled = true
	return "Mock Intent View"
}

func (m *MockIntent) Result() *IntentResult[interface{}] {
	return m.result
}

// SetResult sets the result for this mock intent.
func (m *MockIntent) SetResult(result *IntentResult[interface{}]) {
	m.result = result
}

func TestDefaultIntentRouter_RegisterIntent(t *testing.T) {
	router := NewDefaultIntentRouter()
	factory := func() Intent { return NewMockIntent() }

	err := router.RegisterIntent("test_intent", factory)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Try to register the same intent again
	err = router.RegisterIntent("test_intent", factory)
	if err == nil {
		t.Errorf("expected error when registering duplicate intent")
	}
}

func TestDefaultIntentRouter_ActivateIntent(t *testing.T) {
	router := NewDefaultIntentRouter()
	factory := func() Intent { return NewMockIntent() }

	_ = router.RegisterIntent("test_intent", factory) // nolint: errcheck

	_, err := router.ActivateIntent("test_intent", nil)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	active := router.GetActiveIntent()
	if active == nil {
		t.Errorf("expected active intent to be set")
	}

	mockIntent := active.(*MockIntent)
	if !mockIntent.initCalled {
		t.Errorf("expected intent.Init() to be called")
	}

}

func TestDefaultIntentRouter_ActivateIntent_NotFound(t *testing.T) {
	router := NewDefaultIntentRouter()

	_, err := router.ActivateIntent("nonexistent", nil)
	if err == nil {
		t.Errorf("expected error when activating nonexistent intent")
	}
}

func TestDefaultIntentRouter_GetActiveIntent(t *testing.T) {
	router := NewDefaultIntentRouter()
	mockIntent := NewMockIntent()
	factory := func() Intent { return mockIntent }

	_ = router.RegisterIntent("test_intent", factory) // nolint: errcheck
	_, _ = router.ActivateIntent("test_intent", nil)  // nolint: errcheck

	active := router.GetActiveIntent()
	if active == nil {
		t.Errorf("expected GetActiveIntent() to return an intent")
	}
}

func TestDefaultIntentRouter_HandleMessage(t *testing.T) {
	router := NewDefaultIntentRouter()
	factory := func() Intent { return NewMockIntent() }

	_ = router.RegisterIntent("test_intent", factory) // nolint: errcheck
	_, _ = router.ActivateIntent("test_intent", nil)  // nolint: errcheck

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	_, result := router.HandleMessage(msg)

	active := router.GetActiveIntent().(*MockIntent)
	if active.updateCalled != 1 {
		t.Errorf("expected intent.Update() to be called once, got %d", active.updateCalled)
	}

	if result != nil {
		t.Errorf("expected no result when intent hasn't completed")
	}
}

func TestDefaultIntentRouter_View(t *testing.T) {
	router := NewDefaultIntentRouter()
	factory := func() Intent { return NewMockIntent() }

	_ = router.RegisterIntent("test_intent", factory) // nolint: errcheck
	_, _ = router.ActivateIntent("test_intent", nil)  // nolint: errcheck

	view := router.View()
	if view != "Mock Intent View" {
		t.Errorf("expected view from active intent, got %q", view)
	}

	active := router.GetActiveIntent().(*MockIntent)
	if !active.viewCalled {
		t.Errorf("expected intent.View() to be called")
	}
}

func TestDefaultIntentRouter_Back(t *testing.T) {
	router := NewDefaultIntentRouter()
	factory1 := func() Intent { return NewMockIntent() }
	factory2 := func() Intent { return NewMockIntent() }

	_ = router.RegisterIntent("intent1", factory1) // nolint: errcheck
	_ = router.RegisterIntent("intent2", factory2) // nolint: errcheck

	_, _ = router.ActivateIntent("intent1", nil) // nolint: errcheck
	_, _ = router.ActivateIntent("intent2", nil) // nolint: errcheck

	if router.GetActiveIntent() == nil {
		t.Errorf("expected intent2 to be active")
	}

	_, err := router.Back()
	if err != nil {
		t.Errorf("expected no error when going back, got %v", err)
	}

	if router.GetHistoryDepth() != 1 {
		t.Errorf("expected history depth 1 after going back")
	}
}

func TestDefaultIntentRouter_Back_NoHistory(t *testing.T) {
	router := NewDefaultIntentRouter()

	_, err := router.Back()
	if err == nil {
		t.Errorf("expected error when going back with no history")
	}
}

func TestDefaultIntentRouter_GetHistory(t *testing.T) {
	router := NewDefaultIntentRouter()
	factory1 := func() Intent { return NewMockIntent() }
	factory2 := func() Intent { return NewMockIntent() }

	_ = router.RegisterIntent("intent1", factory1) // nolint: errcheck
	_ = router.RegisterIntent("intent2", factory2) // nolint: errcheck

	_, _ = router.ActivateIntent("intent1", nil) // nolint: errcheck
	_, _ = router.ActivateIntent("intent2", nil) // nolint: errcheck

	history := router.GetHistory()
	if len(history) != 1 {
		t.Errorf("expected history length 1, got %d", len(history))
	}
}

func TestDefaultIntentRouter_HandleMessage_NoActiveIntent(t *testing.T) {
	router := NewDefaultIntentRouter()

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	_, result := router.HandleMessage(msg)

	if result != nil {
		t.Errorf("expected nil when no active intent")
	}
}

func TestDefaultIntentRouter_View_NoActiveIntent(t *testing.T) {
	router := NewDefaultIntentRouter()

	view := router.View()
	if view != "No active intent" {
		t.Errorf("expected default message when no active intent, got %q", view)
	}
}

func TestDefaultIntentRouter_GetHistoryDepth(t *testing.T) {
	router := NewDefaultIntentRouter()
	factory1 := func() Intent { return NewMockIntent() }
	factory2 := func() Intent { return NewMockIntent() }

	_ = router.RegisterIntent("intent1", factory1) // nolint: errcheck
	_ = router.RegisterIntent("intent2", factory2) // nolint: errcheck

	if router.GetHistoryDepth() != 0 {
		t.Errorf("expected depth 0 when no intent active")
	}

	_, _ = router.ActivateIntent("intent1", nil) // nolint: errcheck
	if router.GetHistoryDepth() != 1 {
		t.Errorf("expected depth 1 after first activation")
	}

	_, _ = router.ActivateIntent("intent2", nil) // nolint: errcheck
	if router.GetHistoryDepth() != 2 {
		t.Errorf("expected depth 2 after second activation")
	}
}

func TestDefaultIntentRouter_HandleMessage_WithResult(t *testing.T) {
	router := NewDefaultIntentRouter()
	mockIntent := NewMockIntent()
	factory := func() Intent { return mockIntent }

	_ = router.RegisterIntent("test_intent", factory) // nolint: errcheck
	_, _ = router.ActivateIntent("test_intent", nil)  // nolint: errcheck

	// Set a result on the intent
	expectedResult := NewCompletedResult[interface{}]("test data")
	mockIntent.SetResult(expectedResult)

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	_, result := router.HandleMessage(msg)

	if result == nil {
		t.Errorf("expected result when intent has completed")
	}

	intentResult, ok := result.(*IntentResult[interface{}])
	if !ok {
		t.Errorf("expected IntentResult type")
	}

	if intentResult.Status != Completed {
		t.Errorf("expected Completed status, got %s", intentResult.Status)
	}
}

// =============================================================================
// Theme Management Tests
// =============================================================================

// ThemeAwareMockIntent is a mock intent that supports theme management.
type ThemeAwareMockIntent struct {
	*MockIntent
	themeManager *themes.ThemeManager
}

func NewThemeAwareMockIntent() *ThemeAwareMockIntent {
	return &ThemeAwareMockIntent{
		MockIntent: NewMockIntent(),
	}
}

func (m *ThemeAwareMockIntent) SetThemeManager(tm *themes.ThemeManager) {
	m.themeManager = tm
}

func (m *ThemeAwareMockIntent) GetThemeManager() *themes.ThemeManager {
	return m.themeManager
}

func TestDefaultIntentRouter_ThemeManager_InitializedByDefault(t *testing.T) {
	router := NewDefaultIntentRouter()

	tm := router.GetThemeManager()
	if tm == nil {
		t.Error("expected theme manager to be initialized by default")
	}
}

func TestDefaultIntentRouter_SetThemeManager(t *testing.T) {
	router := NewDefaultIntentRouter()
	customTM := themes.NewThemeManager()

	router.SetThemeManager(customTM)

	if router.GetThemeManager() != customTM {
		t.Error("expected custom theme manager to be set")
	}
}

func TestDefaultIntentRouter_Theme_ReturnsActiveTheme(t *testing.T) {
	router := NewDefaultIntentRouter()

	theme := router.Theme()
	if theme == nil {
		t.Error("expected Theme() to return the active theme")
	}

	// Default theme should be "default"
	if theme.Name() != "default" {
		t.Errorf("expected default theme name 'default', got '%s'", theme.Name())
	}
}

func TestDefaultIntentRouter_Theme_NilWhenNoThemeManager(t *testing.T) {
	router := NewDefaultIntentRouter()
	router.SetThemeManager(nil)

	theme := router.Theme()
	if theme != nil {
		t.Error("expected Theme() to return nil when no theme manager")
	}
}

func TestDefaultIntentRouter_PropagatesThemeToIntent(t *testing.T) {
	router := NewDefaultIntentRouter()
	mockIntent := NewThemeAwareMockIntent()
	factory := func() Intent { return mockIntent }

	_ = router.RegisterIntent("test_intent", factory)
	_, _ = router.ActivateIntent("test_intent", nil)

	// Get the actual intent that was activated (factory creates a new instance)
	active := router.GetActiveIntent()
	themeAware, ok := active.(*ThemeAwareMockIntent)
	if !ok {
		t.Fatal("expected ThemeAwareMockIntent")
	}

	if themeAware.GetThemeManager() == nil {
		t.Error("expected theme manager to be propagated to intent")
	}
}

func TestDefaultIntentRouter_SetThemeManager_PropagatestoActiveIntent(t *testing.T) {
	router := NewDefaultIntentRouter()
	mockIntent := NewThemeAwareMockIntent()
	factory := func() Intent { return mockIntent }

	_ = router.RegisterIntent("test_intent", factory)
	_, _ = router.ActivateIntent("test_intent", nil)

	// Set a new theme manager
	newTM := themes.NewThemeManager()
	router.SetThemeManager(newTM)

	// The active intent should have received the new theme manager
	active := router.GetActiveIntent()
	themeAware, ok := active.(*ThemeAwareMockIntent)
	if !ok {
		t.Fatal("expected ThemeAwareMockIntent")
	}

	if themeAware.GetThemeManager() != newTM {
		t.Error("expected new theme manager to be propagated to active intent")
	}
}
