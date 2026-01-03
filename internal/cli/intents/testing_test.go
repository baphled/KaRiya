package intents

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTestIntentFactory(t *testing.T) {
	t.Run("NewTestIntentFactory creates valid factory", func(t *testing.T) {
		factory := NewTestIntentFactory()
		if factory == nil {
			t.Fatal("expected factory, got nil")
		}
	})

	t.Run("Register and Create work together", func(t *testing.T) {
		factory := NewTestIntentFactory()
		mock := NewMockIntent()
		factory.Register("test", func() Intent {
			return mock
		})

		intent := factory.Create("test")
		if intent == nil {
			t.Fatal("expected intent, got nil")
		}
		if intent != mock {
			t.Fatalf("expected mock intent, got %v", intent)
		}
	})

	t.Run("Create returns nil for unregistered intent", func(t *testing.T) {
		factory := NewTestIntentFactory()
		intent := factory.Create("nonexistent")
		if intent != nil {
			t.Fatalf("expected nil, got %v", intent)
		}
	})

	t.Run("supports multiple registered intents", func(t *testing.T) {
		factory := NewTestIntentFactory()
		mock1 := NewMockIntent()
		mock2 := NewMockIntent()

		factory.Register("intent1", func() Intent { return mock1 })
		factory.Register("intent2", func() Intent { return mock2 })

		intent1 := factory.Create("intent1")
		intent2 := factory.Create("intent2")

		if intent1 != mock1 {
			t.Fatalf("expected mock1, got %v", intent1)
		}
		if intent2 != mock2 {
			t.Fatalf("expected mock2, got %v", intent2)
		}
	})
}

func TestIntentWithState(t *testing.T) {
	t.Run("NewIntentWithState creates valid wrapper", func(t *testing.T) {
		mock := NewMockIntent()
		state := "test state"
		wrapper := NewIntentWithState(mock, state)

		if wrapper == nil {
			t.Fatal("expected wrapper, got nil")
		}
		if wrapper.GetState() != state {
			t.Fatalf("expected state %q, got %q", state, wrapper.GetState())
		}
	})

	t.Run("Delegation to wrapped intent works", func(t *testing.T) {
		mock := NewMockIntent()
		wrapper := NewIntentWithState(mock, "state")

		wrapper.Init()
		if !mock.initCalled {
			t.Fatal("expected initCalled to be true")
		}

		wrapper.Update(tea.KeyMsg{})
		if mock.updateCalled <= 0 {
			t.Fatalf("expected updateCalled > 0, got %d", mock.updateCalled)
		}

		wrapper.View()
		if !mock.viewCalled {
			t.Fatal("expected viewCalled to be true")
		}

		// Result returns whatever the mock returns (which is nil by default)
		result := wrapper.Result()
		// Just verify it doesn't panic - result can be nil
		_ = result
	})

	t.Run("SetState changes the state", func(t *testing.T) {
		mock := NewMockIntent()
		wrapper := NewIntentWithState(mock, "initial")

		wrapper.SetState("updated")
		if wrapper.GetState() != "updated" {
			t.Fatalf("expected state 'updated', got %q", wrapper.GetState())
		}
	})
}

func TestIntentTestHarness(t *testing.T) {
	t.Run("NewIntentTestHarness creates valid harness", func(t *testing.T) {
		mock := NewMockIntent()
		harness := NewIntentTestHarness(t, mock)

		if harness == nil {
			t.Fatal("expected harness, got nil")
		}
	})

	t.Run("Init calls the intent's Init method", func(t *testing.T) {
		mock := NewMockIntent()
		harness := NewIntentTestHarness(t, mock)

		harness.Init()
		if !mock.initCalled {
			t.Fatal("expected initCalled to be true")
		}
	})

	t.Run("SendMessage calls the intent's Update method", func(t *testing.T) {
		mock := NewMockIntent()
		harness := NewIntentTestHarness(t, mock)

		harness.SendMessage(tea.KeyMsg{})
		if mock.updateCalled <= 0 {
			t.Fatalf("expected updateCalled > 0, got %d", mock.updateCalled)
		}
	})

	t.Run("GetView returns the intent's view", func(t *testing.T) {
		mock := NewMockIntent()
		harness := NewIntentTestHarness(t, mock)

		view := harness.GetView()
		if view != "Mock Intent View" {
			t.Fatalf("expected 'Mock Intent View', got %q", view)
		}
	})

	t.Run("GetResult returns the intent's result", func(t *testing.T) {
		mock := NewMockIntent()
		expectedResult := &IntentResult[interface{}]{Status: Completed}
		mock.SetResult(expectedResult)

		harness := NewIntentTestHarness(t, mock)
		result := harness.GetResult()

		if result != expectedResult {
			t.Fatalf("expected %v, got %v", expectedResult, result)
		}
	})

	t.Run("AssertResultCompleted passes for completed result", func(t *testing.T) {
		mock := NewMockIntent()
		mock.SetResult(&IntentResult[interface{}]{Status: Completed})

		harness := NewIntentTestHarness(t, mock)
		// Should not panic
		harness.AssertResultCompleted()
	})

	t.Run("AssertResultCancelled passes for cancelled result", func(t *testing.T) {
		mock := NewMockIntent()
		mock.SetResult(&IntentResult[interface{}]{Status: Cancelled})

		harness := NewIntentTestHarness(t, mock)
		// Should not panic
		harness.AssertResultCancelled()
	})

	t.Run("AssertViewContains checks for content", func(t *testing.T) {
		mock := NewMockIntent()
		harness := NewIntentTestHarness(t, mock)

		// Should not panic
		harness.AssertViewContains("Mock")
	})
}

func TestIntentRouterTestHelper(t *testing.T) {
	t.Run("NewIntentRouterTestHelper creates valid helper", func(t *testing.T) {
		router := NewDefaultIntentRouter()
		helper := NewIntentRouterTestHelper(t, router)

		if helper == nil {
			t.Fatal("expected helper, got nil")
		}
	})

	t.Run("ActivateIntent activates intent", func(t *testing.T) {
		router := NewDefaultIntentRouter()
		helper := NewIntentRouterTestHelper(t, router)

		mock := NewMockIntent()
		_ = router.RegisterIntent("test", func() Intent { // nolint: errcheck
			return mock
		})

		contextMap := map[string]interface{}{"test": "value"}
		helper.ActivateIntent("test", contextMap)

		active := helper.GetActiveIntent()
		if active != mock {
			t.Fatalf("expected mock intent, got %v", active)
		}
	})

	t.Run("GetHistory returns history", func(t *testing.T) {
		router := NewDefaultIntentRouter()
		helper := NewIntentRouterTestHelper(t, router)

		mock1 := NewMockIntent()
		mock2 := NewMockIntent()

		_ = router.RegisterIntent("test1", func() Intent { return mock1 }) // nolint: errcheck
		_ = router.RegisterIntent("test2", func() Intent { return mock2 }) // nolint: errcheck

		contextMap := map[string]interface{}{}
		helper.ActivateIntent("test1", contextMap)
		helper.ActivateIntent("test2", contextMap)

		history := helper.GetHistory()
		// History contains previous intents, not the current one
		// After activating test1 and test2, history should have test1
		if len(history) != 1 {
			t.Fatalf("expected history length 1, got %d", len(history))
		}
	})

	t.Run("GoBack navigates back", func(t *testing.T) {
		router := NewDefaultIntentRouter()
		helper := NewIntentRouterTestHelper(t, router)

		mock1 := NewMockIntent()
		mock2 := NewMockIntent()

		_ = router.RegisterIntent("test1", func() Intent { return mock1 }) // nolint: errcheck
		_ = router.RegisterIntent("test2", func() Intent { return mock2 }) // nolint: errcheck

		contextMap := map[string]interface{}{}
		helper.ActivateIntent("test1", contextMap)
		helper.ActivateIntent("test2", contextMap)

		helper.GoBack()

		active := helper.GetActiveIntent()
		if active != mock1 {
			t.Fatalf("expected mock1 after back, got %v", active)
		}
	})

	t.Run("AssertIntentActive verifies active intent", func(t *testing.T) {
		router := NewDefaultIntentRouter()
		helper := NewIntentRouterTestHelper(t, router)

		mock := NewMockIntent()
		_ = router.RegisterIntent("test", func() Intent { return mock }) // nolint: errcheck

		contextMap := map[string]interface{}{}
		helper.ActivateIntent("test", contextMap)
		// Should not panic
		helper.AssertIntentActive()
	})

	t.Run("AssertHistoryLength verifies history length", func(t *testing.T) {
		router := NewDefaultIntentRouter()
		helper := NewIntentRouterTestHelper(t, router)

		mock1 := NewMockIntent()
		mock2 := NewMockIntent()

		_ = router.RegisterIntent("test1", func() Intent { return mock1 }) // nolint: errcheck
		_ = router.RegisterIntent("test2", func() Intent { return mock2 }) // nolint: errcheck

		contextMap := map[string]interface{}{}
		helper.ActivateIntent("test1", contextMap)
		helper.ActivateIntent("test2", contextMap)

		// After activating test1 and test2, history should have 1 (test1)
		helper.AssertHistoryLength(1)
	})
}

func TestAssertionHelpers(t *testing.T) {
	t.Run("contains finds substring", func(t *testing.T) {
		text := "Hello World"
		result := contains(text, "World")
		if !result {
			t.Fatal("expected 'World' to be in 'Hello World'")
		}
	})

	t.Run("contains returns false for missing substring", func(t *testing.T) {
		text := "Hello World"
		result := contains(text, "Goodbye")
		if result {
			t.Fatal("expected 'Goodbye' not to be in 'Hello World'")
		}
	})
}
