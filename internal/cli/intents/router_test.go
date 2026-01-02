package intents

import (
	"testing"

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

	router.RegisterIntent("test_intent", factory)

	cmd, err := router.ActivateIntent("test_intent", nil)
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

	if cmd == nil {
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

	router.RegisterIntent("test_intent", factory)
	router.ActivateIntent("test_intent", nil)

	active := router.GetActiveIntent()
	if active == nil {
		t.Errorf("expected GetActiveIntent() to return an intent")
	}
}

func TestDefaultIntentRouter_HandleMessage(t *testing.T) {
	router := NewDefaultIntentRouter()
	factory := func() Intent { return NewMockIntent() }

	router.RegisterIntent("test_intent", factory)
	router.ActivateIntent("test_intent", nil)

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

	router.RegisterIntent("test_intent", factory)
	router.ActivateIntent("test_intent", nil)

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

	router.RegisterIntent("intent1", factory1)
	router.RegisterIntent("intent2", factory2)

	router.ActivateIntent("intent1", nil)
	router.ActivateIntent("intent2", nil)

	if router.GetActiveIntent() == nil {
		t.Errorf("expected intent2 to be active")
	}

	cmd, err := router.Back()
	if err != nil {
		t.Errorf("expected no error when going back, got %v", err)
	}

	if cmd == nil {
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

	router.RegisterIntent("intent1", factory1)
	router.RegisterIntent("intent2", factory2)

	router.ActivateIntent("intent1", nil)
	router.ActivateIntent("intent2", nil)

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

	router.RegisterIntent("intent1", factory1)
	router.RegisterIntent("intent2", factory2)

	if router.GetHistoryDepth() != 0 {
		t.Errorf("expected depth 0 when no intent active")
	}

	router.ActivateIntent("intent1", nil)
	if router.GetHistoryDepth() != 1 {
		t.Errorf("expected depth 1 after first activation")
	}

	router.ActivateIntent("intent2", nil)
	if router.GetHistoryDepth() != 2 {
		t.Errorf("expected depth 2 after second activation")
	}
}

func TestDefaultIntentRouter_HandleMessage_WithResult(t *testing.T) {
	router := NewDefaultIntentRouter()
	mockIntent := NewMockIntent()
	factory := func() Intent { return mockIntent }

	router.RegisterIntent("test_intent", factory)
	router.ActivateIntent("test_intent", nil)

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
