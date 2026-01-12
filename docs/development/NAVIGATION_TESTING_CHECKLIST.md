# Navigation Testing Checklist

**Last Updated**: 2026-01-12  
**Purpose**: Quick reference for implementing navigation tests  
**Related**: [Navigation Testing Guide](NAVIGATION_TESTING_GUIDE.md) (complete guide)

---

## Quick Start

**Before implementing a new intent or modifying navigation**:
1. ✅ Read [TUI Standards](../TUI_STANDARDS.md) for navigation requirements
2. ✅ Review [Navigation Testing Guide](NAVIGATION_TESTING_GUIDE.md) for patterns
3. ✅ Use this checklist to ensure complete test coverage

**Test file naming**:
- Navigation tests: `internal/cli/intents/<intent_name>_navigation_test.go`
- Escape key tests: `internal/cli/intents/<intent_name>_escape_test.go`
- E2E tests: `internal/testutil/e2e/<intent_name>_e2e_test.go`

---

## Pre-Implementation Checklist

Use this checklist **before writing any code**:

### 1. Define State Machine

- [ ] Listed all states (root, intermediate, async, final)
- [ ] Documented all valid state transitions
- [ ] Identified which state is the root state
- [ ] Identified which states are intermediate states
- [ ] Identified which states are async operations
- [ ] Identified which states are final states (success/error)

**Example**:
```
CaptureEvent States:
- ChooseStrategy (root)
- Form (intermediate)
- Review (intermediate)
- Submit (final)
```

### 2. Define Escape Key Behavior

- [ ] Documented Escape behavior for root state (cancel → main menu)
- [ ] Documented Escape behavior for each intermediate state (go back → previous)
- [ ] Documented Escape behavior for async states (let complete → go back)
- [ ] Documented Escape behavior for final states (retry/go back)

**Example**:
```
ChooseStrategy (root) + Esc → Cancel intent, return to menu
Form (intermediate) + Esc → Go back to ChooseStrategy
Review (intermediate) + Esc → Go back to Form
Submit (final) + Esc → Go back to Review (preserve error if exists)
```

### 3. Identify Context to Preserve

- [ ] Listed what context must be preserved on back navigation
  - [ ] User selections (profile, audience, format, etc.)
  - [ ] Form data (text, date, company, etc.)
  - [ ] Error messages
  - [ ] Scroll positions
  - [ ] Filter settings

---

## Test Coverage Requirements

Use this checklist to verify **100% navigation test coverage**:

### A. Escape Key Tests (CRITICAL)

Create file: `<intent_name>_escape_test.go`

- [ ] **Root state escape test** (cancel intent)
  ```go
  It("should cancel intent when esc is pressed from root state", func() {
      intent.state.currentState = StateRoot
      intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
      Expect(intent.active).To(BeFalse())
      Expect(intent.result.Status).To(Equal(Cancelled))
  })
  ```

- [ ] **Intermediate state escape tests** (one per state)
  ```go
  It("should go back to StateX when esc is pressed from StateY", func() {
      intent.state.currentState = StateY
      intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
      Expect(intent.state.currentState).To(Equal(StateX))
      Expect(intent.active).To(BeTrue())
  })
  ```

- [ ] **Async state escape tests**
  ```go
  It("should go back from async state when esc is pressed", func() {
      intent.state.currentState = StateGenerating
      intent.state.isGenerating = true
      intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
      Expect(intent.state.currentState).To(Equal(StatePrevious))
  })
  ```

- [ ] **Final state escape tests** (error recovery)
  ```go
  It("should go back from error state when esc is pressed", func() {
      intent.state.currentState = StateFinal
      intent.state.error = &IntentError{Code: "TEST"}
      intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
      Expect(intent.state.currentState).To(Equal(StatePrevious))
  })
  
  It("should preserve error when going back", func() {
      originalError := intent.state.error
      intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
      Expect(intent.state.error).To(Equal(originalError))
  })
  ```

**Target**: 1 escape test per state = N tests (where N = number of states)

### B. Forward Navigation Tests

Create file: `<intent_name>_navigation_test.go`

- [ ] **Menu selection test**
  ```go
  It("should show IntentName as menu item", func() {
      env.AssertViewContains("Intent Name")
  })
  
  It("should navigate to IntentName when selected", func() {
      env.SelectIntentByName("intent_name")
      env.AssertViewContainsAny("Root State", "Initial View")
  })
  ```

- [ ] **State transition tests** (one per transition)
  ```go
  It("should transition from StateA to StateB", func() {
      env.SelectIntentByName("intent_name")
      env.Confirm()  // Or other action
      env.AssertViewContainsAny("StateB", "indicators")
  })
  ```

- [ ] **Complete workflow test**
  ```go
  It("should navigate through complete workflow", func() {
      env.SelectIntentByName("intent_name")
      env.AssertViewContains("StateRoot")
      
      env.Confirm()
      env.AssertViewContains("StateIntermediate")
      
      // ... continue through all states ...
  })
  ```

### C. Back Navigation Tests

- [ ] **Back from each state test**
  ```go
  It("should navigate back from StateY to StateX", func() {
      env.SelectIntentByName("intent_name")
      env.Confirm()  // Go to StateY
      env.Cancel()   // Go back to StateX
      env.AssertViewContains("StateX")
  })
  ```

- [ ] **Back to main menu test**
  ```go
  It("should return to main menu when cancelling from root state", func() {
      env.SelectIntentByName("intent_name")
      env.Cancel()
      env.AssertViewContains("Capture Event")
  })
  ```

- [ ] **Context preservation test**
  ```go
  It("should preserve selections when navigating back", func() {
      env.SelectIntentByName("intent_name")
      env.Confirm()  // Select something
      env.Cancel()   // Go back
      env.Confirm()  // Re-select
      // Verify selection is preserved
  })
  ```

### D. Universal Shortcuts Tests

- [ ] **'m' key test** (return to main menu)
  ```go
  It("should return to main menu when 'm' is pressed", func() {
      env.SelectIntentByName("intent_name")
      env.PressKeyRune('m')
      env.AssertViewContains("Capture Event")
  })
  ```

- [ ] **'q' key test** (quit application)
  ```go
  It("should quit application when 'q' is pressed", func() {
      env.SelectIntentByName("intent_name")
      env.Quit()
      // App terminates - no panic
  })
  ```

### E. List Navigation Tests (if intent has lists)

- [ ] **Arrow key navigation**
  ```go
  It("should navigate down with down arrow", func() {
      env.SelectIntentByName("intent_name")
      env.PressKey(tea.KeyDown)
      // Verify navigation
  })
  
  It("should navigate up with up arrow", func() {
      env.PressKey(tea.KeyUp)
      // Verify navigation
  })
  ```

- [ ] **Vim-style navigation**
  ```go
  It("should navigate down with 'j' key", func() {
      env.PressKeyRune('j')
      // Verify navigation
  })
  
  It("should navigate up with 'k' key", func() {
      env.PressKeyRune('k')
      // Verify navigation
  })
  ```

- [ ] **Boundary tests**
  ```go
  It("should not go below first item", func() {
      for i := 0; i < 5; i++ {
          env.PressKey(tea.KeyUp)
      }
      view := env.GetView()
      Expect(view).NotTo(ContainSubstring("panic"))
  })
  
  It("should not go above last item", func() {
      for i := 0; i < 10; i++ {
          env.PressKey(tea.KeyDown)
      }
      view := env.GetView()
      Expect(view).NotTo(ContainSubstring("panic"))
  })
  ```

### F. Error Recovery Tests

- [ ] **Empty state test**
  ```go
  It("should handle empty state gracefully", func() {
      env.SelectIntentByName("intent_name")
      view := env.GetView()
      Expect(view).NotTo(ContainSubstring("panic"))
  })
  ```

- [ ] **Error state navigation test**
  ```go
  It("should allow navigation from error state", func() {
      // Trigger error
      env.Cancel()
      env.AssertViewNotContains("panic")
  })
  ```

- [ ] **Rapid key press test**
  ```go
  It("should handle rapid navigation", func() {
      for i := 0; i < 20; i++ {
          env.PressKeyRune('j')
          env.PressKeyRune('k')
      }
      view := env.GetView()
      Expect(view).NotTo(ContainSubstring("panic"))
  })
  ```

### G. E2E Workflow Tests

Create file: `internal/testutil/e2e/<intent_name>_e2e_test.go`

- [ ] **Database persistence test** (if applicable)
  ```go
  It("should persist data after workflow completion", func() {
      env := e2e.Setup(GinkgoT())
      defer env.Cleanup()
      
      // Complete workflow
      env.SelectIntentByName("intent_name")
      // ... complete workflow ...
      
      env.AssertEventCount(1)  // Or other data
  })
  ```

- [ ] **Cancel before persistence test**
  ```go
  It("should not persist data when cancelled", func() {
      env.SelectIntentByName("intent_name")
      // ... partial workflow ...
      env.Cancel()
      env.AssertEventCount(0)
  })
  ```

- [ ] **Session restart test** (if applicable)
  ```go
  It("should maintain data after restart", func() {
      env.PopulateTestData(3, 0, 0)
      env.SimulateRestart()
      env.AssertEventCount(3)
  })
  ```

---

## Test Template

### Escape Key Test Template

**File**: `internal/cli/intents/<intent_name>_escape_test.go`

```go
package intents

import (
    tea "github.com/charmbracelet/bubbletea"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

var _ = Describe("IntentName - Escape Key Behavior", func() {
    var intent *YourIntent

    BeforeEach(func() {
        // Setup intent with test context
        ctx := &YourIntentContext{
            // ... context setup ...
        }
        var err error
        intent, err = NewYourIntent(ctx)
        Expect(err).NotTo(HaveOccurred())
        intent.Init()
    })

    Describe("RootState (Cancel to Menu)", func() {
        BeforeEach(func() {
            intent.state.currentState = StateRoot
        })

        It("should cancel intent when esc is pressed", func() {
            intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

            Expect(intent.active).To(BeFalse())
            Expect(intent.result).NotTo(BeNil())
            Expect(intent.result.Status).To(Equal(Cancelled))
        })
    })

    Describe("IntermediateState (Go Back)", func() {
        BeforeEach(func() {
            intent.state.currentState = StateIntermediate
            // ... setup required state data ...
        })

        It("should go back to RootState when esc is pressed", func() {
            intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

            Expect(intent.state.currentState).To(Equal(StateRoot))
            Expect(intent.active).To(BeTrue())
        })
    })

    // Add describe block for EACH state in your intent
    // ...
})
```

### Navigation Test Template

**File**: `internal/cli/intents/<intent_name>_navigation_test.go`

```go
package intents_test

import (
    "github.com/baphled/kariya/internal/testutil/e2e"
    tea "github.com/charmbracelet/bubbletea"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

var _ = Describe("IntentName Navigation", func() {
    var env *e2e.TestEnv

    Describe("Navigation to IntentName", func() {
        BeforeEach(func() {
            env = e2e.SetupWithMemory(GinkgoT())
        })

        AfterEach(func() {
            env.Cleanup()
        })

        It("should show IntentName as menu item", func() {
            env.AssertViewContains("Intent Name")
        })

        It("should navigate to IntentName when selected", func() {
            env.SelectIntentByName("intent_name")
            env.AssertViewContainsAny("Root State", "indicators")
        })
    })

    Describe("Forward Navigation", func() {
        BeforeEach(func() {
            env = e2e.SetupWithMemory(GinkgoT())
            env.SelectIntentByName("intent_name")
        })

        AfterEach(func() {
            env.Cleanup()
        })

        It("should navigate to StateB when confirmed from StateA", func() {
            env.Confirm()
            env.AssertViewContainsAny("StateB", "indicators")
        })
    })

    Describe("Back Navigation", func() {
        BeforeEach(func() {
            env = e2e.SetupWithMemory(GinkgoT())
        })

        AfterEach(func() {
            env.Cleanup()
        })

        It("should navigate back through all states to main menu", func() {
            env.SelectIntentByName("intent_name")
            env.Confirm()  // Go to StateB
            
            env.Cancel()   // Back to StateA
            env.AssertViewContains("StateA")
            
            env.Cancel()   // Back to menu
            env.AssertViewContains("Capture Event")
        })
    })
})
```

### E2E Test Template

**File**: `internal/testutil/e2e/<intent_name>_e2e_test.go`

```go
package e2e_test

import (
    "github.com/baphled/kariya/internal/testutil/e2e"
    . "github.com/onsi/ginkgo/v2"
)

var _ = Describe("E2E IntentName Workflow", func() {
    var env *e2e.TestEnv

    Describe("Complete Workflow", func() {
        BeforeEach(func() {
            env = e2e.Setup(GinkgoT())
        })

        AfterEach(func() {
            env.Cleanup()
        })

        It("should complete workflow and persist data", func() {
            env.SelectIntentByName("intent_name")
            // ... complete workflow ...
            
            env.AssertEventCount(1)  // Verify persistence
        })
    })

    Describe("Cancellation", func() {
        BeforeEach(func() {
            env = e2e.Setup(GinkgoT())
        })

        AfterEach(func() {
            env.Cleanup()
        })

        It("should not persist data when cancelled", func() {
            env.SelectIntentByName("intent_name")
            // ... partial workflow ...
            env.Cancel()
            
            env.AssertEventCount(0)  // No data persisted
        })
    })
})
```

---

## Pre-Commit Checklist

Before committing navigation changes:

### Code Changes

- [ ] Added/modified escape key handlers in `Update()` method
- [ ] Verified escape behavior for each state type
- [ ] Added 'm' key handler (return to main menu)
- [ ] Added 'q' key handler (quit application)
- [ ] Added boundary checks for list navigation
- [ ] Preserved context on back navigation

### Test Changes

- [ ] Created/updated `<intent_name>_escape_test.go`
- [ ] Created/updated `<intent_name>_navigation_test.go`
- [ ] Created/updated E2E test file (if applicable)
- [ ] All escape key tests pass
- [ ] All navigation tests pass
- [ ] All E2E tests pass

### Documentation

- [ ] Updated workflow diagram (if state machine changed)
- [ ] Updated AGENTS.md (if new intent added)
- [ ] Added navigation section to intent documentation

### Verification

- [ ] Run: `go test ./internal/cli/intents/...`
- [ ] Run: `go test ./internal/testutil/e2e/...`
- [ ] Run: `go test -race ./...` (check for race conditions)
- [ ] Run: `make check-compliance`
- [ ] Manual test in actual TUI application

---

## Test Coverage Goals

### Minimum Coverage (Required)

- ✅ **Escape key**: 1 test per state (N tests)
- ✅ **Forward navigation**: 1 test per state transition (N-1 tests)
- ✅ **Back navigation**: 1 test for complete backward flow
- ✅ **Universal shortcuts**: 'm' and 'q' tests

**Total Minimum**: ~(2N + 2) tests per intent

### Recommended Coverage

- ✅ Minimum coverage (above)
- ✅ List navigation (if applicable): 6 tests (up, down, j, k, boundaries)
- ✅ Error recovery: 3 tests (empty state, error state, rapid keys)
- ✅ E2E workflow: 2 tests (complete + cancel)

**Total Recommended**: ~(2N + 13) tests per intent

### Excellent Coverage (Aspirational)

- ✅ Recommended coverage (above)
- ✅ Context preservation: 1 test per context item
- ✅ Multi-intent navigation: 2 tests
- ✅ Session restart: 1 test (if applicable)

**Total Excellent**: ~(2N + 20) tests per intent

### Coverage Examples

**CaptureEvent** (4 states):
- Minimum: ~10 tests
- Recommended: ~21 tests
- Actual: **304 tests** (well above excellent)

**GenerateCV** (10 states):
- Minimum: ~22 tests
- Recommended: ~33 tests
- Actual: **19 integration tests + numerous unit tests**

---

## Common Mistakes to Avoid

### ❌ Missing Escape Handler

```go
// WRONG - no escape handler
func (i *YourIntent) Update(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch i.state {
        case StateIntermediate:
            switch msg.Type {
            case tea.KeyEnter:
                // ... handle enter ...
            // ❌ Missing case tea.KeyEsc
            }
        }
    }
}
```

**Fix**: Add escape handler for every state

### ❌ Cancelling Instead of Going Back

```go
// WRONG - cancels from intermediate state
case StateIntermediate:
    switch msg.Type {
    case tea.KeyEsc:
        i.active = false  // ❌ Should go back, not cancel
        i.result = &IntentResult{Status: Cancelled}
    }
```

**Fix**: Go back to previous state, only cancel from root state

### ❌ Losing Context on Back Navigation

```go
// WRONG - loses error context
case StateError:
    switch msg.Type {
    case tea.KeyEsc:
        i.state.error = nil  // ❌ Error context lost
        i.state.currentState = StatePrevious
    }
```

**Fix**: Preserve error, selections, and other context

### ❌ No Boundary Checks

```go
// WRONG - can go out of bounds
case tea.KeyDown:
    i.selectedIndex++  // ❌ Can exceed list length
```

**Fix**: Add boundary checks

### ❌ Missing 'm' Key Handler

```go
// WRONG - 'm' only works in some states
case StateIntermediate:
    switch msg.Type {
    case tea.KeyRunes:
        if msg.Runes[0] == 'm' {
            // Handle 'm'
        }
    }
// ❌ Other states don't handle 'm'
```

**Fix**: Handle 'm' globally before state-specific logic

### ❌ Not Testing E2E

```go
// WRONG - only unit tests
var _ = Describe("IntentName", func() {
    // Only tests individual methods
    // ❌ Never tests complete workflow
})
```

**Fix**: Add E2E tests for complete workflows

---

## Quick Reference Commands

### Run Tests

```bash
# All tests
make test

# Intent tests only
go test ./internal/cli/intents/...

# E2E tests only
go test ./internal/testutil/e2e/...

# Specific intent
go test ./internal/cli/intents/ -run CaptureNavigation

# With race detector
go test -race ./...

# With coverage
go test -cover ./...

# Focus specific test
# Use FIt("test name") in test file, then:
ginkgo -r --focus="test name" ./internal/cli/intents/
```

### Compliance Checks

```bash
# Full compliance check
make check-compliance

# Just tests
make test

# Just linting
make lint

# Just formatting
make fmt
```

---

## Related Documentation

- **[Navigation Testing Guide](NAVIGATION_TESTING_GUIDE.md)** - Complete guide (this checklist's companion)
- **[TUI Standards](../TUI_STANDARDS.md)** - UI/UX navigation requirements
- **[TUI Developer Guide](../TUI_DEVELOPER_GUIDE.md)** - Creating TUI components
- **[Keyboard Shortcuts Guide](../KEYBOARD_SHORTCUTS_GUIDE.md)** - User-facing shortcuts
- **[Keyboard System Guide](KEYBOARD_SYSTEM_GUIDE.md)** - Developer keyboard implementation
- **[Workflow Guides](../workflows/README.md)** - User workflow documentation

---

## Next Steps

1. ✅ Use this checklist for every intent you create or modify
2. ✅ Read [Navigation Testing Guide](NAVIGATION_TESTING_GUIDE.md) for detailed examples
3. ✅ Review existing tests as examples:
   - `internal/cli/intents/capture_navigation_test.go`
   - `internal/cli/intents/capture_event_escape_test.go`
   - `internal/testutil/e2e/capture_e2e_test.go`
4. ✅ Run `make test` to verify all navigation tests pass

---

*Last Updated: 2026-01-12*  
*Maintainer: KaRiya Development Team*
