# KaRiya Testing Patterns Guide

**A comprehensive guide for understanding and writing tests in the KaRiya project**

---

## Table of Contents

1. [Overview](#overview)
2. [Test Types](#test-types)
3. [Test Infrastructure](#test-infrastructure)
4. [Mock Usage Guide](#mock-usage-guide)
5. [Writing Tests](#writing-tests)
6. [Suite Organization](#suite-organization)
7. [Common Patterns](#common-patterns)
8. [Troubleshooting](#troubleshooting)

---

## Overview

KaRiya uses **Ginkgo v2** with **Gomega** for all testing. The testing strategy follows a clear hierarchy:

```
Unit Tests → Navigation Tests → E2E Tests
   ↓              ↓                ↓
 Fast         Medium           Slower
 Isolated     UI-focused       Persistence-focused
```

### Quick Reference

| Test Type | Location | Setup Function | Purpose |
|-----------|----------|----------------|---------|
| Unit | Same package as code | N/A | Test isolated functions |
| Navigation | `internal/cli/intents/` | `e2e.SetupWithMemory()` | Test keyboard/UI interactions |
| E2E | `internal/testutil/e2e/` | `e2e.Setup()` | Test data persistence |

---

## Test Types

### Unit Tests

**Purpose**: Test individual functions/methods in isolation

**Location**: Same package as the code being tested (e.g., `internal/cli/intents/contract_test.go`)

**When to Use**:
- Testing pure functions
- Testing struct methods in isolation
- Testing type conversions and utilities
- Testing error conditions

**Example Files**:
- `internal/cli/intents/contract_test.go` - BaseIntent tests
- `internal/cli/intents/result_test.go` - IntentResult tests
- `internal/cli/intents/router_test.go` - Router logic tests

**Example**:
```go
var _ = Describe("IntentResult", func() {
    Describe("NewSuccessResult", func() {
        It("should create a result with StatusCompleted", func() {
            result := NewSuccessResult("test data")
            Expect(result.Status).To(Equal(StatusCompleted))
            Expect(result.Data).To(Equal("test data"))
        })

        It("should return nil error", func() {
            result := NewSuccessResult("data")
            Expect(result.Error).To(BeNil())
        })
    })
})
```

### Navigation Tests

**Purpose**: Test keyboard navigation and UI state transitions

**Location**: `internal/cli/intents/` with `*_navigation_test.go` suffix

**When to Use**:
- Testing keyboard input handling (arrow keys, vim keys, Enter, Esc)
- Testing screen transitions
- Testing view rendering
- Testing UI component interactions

**Setup**: Uses `e2e.SetupWithMemory(GinkgoT())` for fast in-memory testing

**Example Files**:
- `internal/cli/intents/browse_navigation_test.go` - Timeline navigation
- `internal/cli/intents/burst_management_navigation_test.go` - Burst list navigation

**Example**:
```go
var _ = Describe("Browse Navigation", func() {
    var env *e2e.TestEnv

    BeforeEach(func() {
        env = e2e.SetupWithMemory(GinkgoT()) // Fast in-memory setup
    })

    AfterEach(func() {
        env.Cleanup()
    })

    It("should navigate down with 'j' key", func() {
        env.SelectIntentByName("browse_timeline")
        env.PressKeyRune('j')
        view := env.GetView()
        Expect(view).NotTo(BeEmpty())
    })

    It("should return to main menu on Escape", func() {
        env.SelectIntentByName("browse_timeline")
        env.Cancel() // Sends Esc
        env.AssertViewContains("Capture Event")
    })
})
```

### E2E Tests

**Purpose**: Test data persistence across intent boundaries and app restarts

**Location**: `internal/testutil/e2e/` with `*_e2e_test.go` suffix

**When to Use**:
- Testing database read/write operations
- Testing data integrity after workflows
- Testing persistence across app restarts
- Testing repository integration

**Setup**: Uses `e2e.Setup(GinkgoT())` for SQLite database setup

**Out of Scope** (use Navigation tests instead):
- Keyboard navigation
- View rendering
- State transitions without persistence

**Example Files**:
- `internal/testutil/e2e/browse_e2e_test.go` - Event persistence
- `internal/testutil/e2e/burst_management_e2e_test.go` - Burst persistence

**Example**:
```go
var _ = Describe("E2E Browse Workflow", func() {
    var env *e2e.TestEnv

    BeforeEach(func() {
        env = e2e.Setup(GinkgoT()) // SQLite setup
    })

    AfterEach(func() {
        env.Cleanup()
    })

    It("should persist events after restart", func() {
        env.PopulateTestData(3, 0, 0)
        env.AssertEventCount(3)

        env.SimulateRestart() // Key: Tests persistence

        env.AssertEventCount(3)
    })
})
```

---

## Test Infrastructure

### TestEnv

The `e2e.TestEnv` provides a complete test environment:

```go
type TestEnv struct {
    app          *app.App
    repo         career.Repository
    burstRepo    career.BurstRepository
    factRepo     career.FactRepository
    service      *career.Service
    globalCtx    *context.GlobalContext
}
```

**Key Methods**:

| Method | Description |
|--------|-------------|
| `Setup(GinkgoT())` | Create environment with SQLite persistence |
| `SetupWithMemory(GinkgoT())` | Create environment with in-memory storage |
| `Cleanup()` | Clean up resources (call in AfterEach) |
| `PopulateTestData(events, bursts, facts int)` | Add test data |
| `SimulateRestart()` | Simulate app restart (tests persistence) |

**Navigation Methods**:

| Method | Description |
|--------|-------------|
| `SelectIntentByName(name)` | Navigate to specific intent |
| `PressKey(tea.Key)` | Send keyboard key |
| `PressKeyRune(rune)` | Send character key |
| `Confirm()` | Press Enter |
| `Cancel()` | Press Escape |
| `GoBack()` | Press 'b' for back |
| `Quit()` | Press 'q' to quit |

**Assertion Methods**:

| Method | Description |
|--------|-------------|
| `AssertViewContains(text)` | View contains exact text |
| `AssertViewContainsAny(texts...)` | View contains any of texts |
| `AssertEventCount(n)` | Verify event count |
| `GetView()` | Get current view string |
| `GetEvents()` | Get all events |

### Ginkgo Suite Files

Each test package needs a suite file:

```go
// internal/cli/intents/intents_suite_test.go
package intents_test

import (
    "testing"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

func TestIntentsSuite(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Intents Suite")
}
```

---

## Mock Usage Guide

Mocks are defined in `internal/cli/intents/testing_helpers.go`.

### MockIntent

**Purpose**: Basic intent mock for testing router and intent lifecycle

**Use When**:
- Testing intent registration
- Testing router dispatch
- Testing basic intent interface compliance

```go
mock := NewMockIntent()
mock.Init()
Expect(mock.InitCalled).To(BeTrue())

mock.Update(tea.KeyMsg{Type: tea.KeyEnter})
Expect(mock.UpdateCalled).To(Equal(1))
Expect(mock.Messages).To(HaveLen(1))
```

**Fields**:
- `InitCalled bool` - Init was invoked
- `UpdateCalled int` - Number of Update calls
- `ViewCalled bool` - View was invoked
- `Messages []tea.Msg` - All received messages
- `Result_ *IntentResult[interface{}]` - Set result for testing

### ThemeAwareMockIntent

**Purpose**: Test theme management integration

**Use When**:
- Testing theme injection
- Testing theme-aware rendering
- Testing theme manager lifecycle

```go
mock := NewThemeAwareMockIntent()
tm := themes.NewThemeManager()
mock.SetThemeManager(tm)
Expect(mock.GetThemeManager()).To(Equal(tm))
```

### MockIntentWithSelection

**Purpose**: Test selection state in list-based intents

**Use When**:
- Testing list navigation state
- Testing selection preservation
- Testing cursor position

```go
mock := NewMockIntentWithSelection()
mock.SetSelectedIndex(5)
Expect(mock.GetSelectedIndex()).To(Equal(5))
```

### MockIntentWithStateMachine

**Purpose**: Test state machine patterns and track state transitions

**Use When**:
- Testing intent state transitions
- Verifying state history
- Testing complex state flows
- Debugging state machine behavior

```go
mock := NewMockIntentWithStateMachine("initial")
Expect(mock.CurrentState).To(Equal("initial"))

// Simple transition
mock.TransitionTo("working")
Expect(mock.CurrentState).To(Equal("working"))

// Transition with trigger tracking
mock.TransitionToWithTrigger("complete", "enter_key")
transitions := mock.GetTransitions()
Expect(transitions[0].From).To(Equal("working"))
Expect(transitions[0].To).To(Equal("complete"))
Expect(transitions[0].Trigger).To(Equal("enter_key"))

// Get full state history
history := mock.GetStateHistory()
Expect(history).To(Equal([]string{"initial", "working", "complete"}))

// Reset for new test scenario
mock.Reset("initial")
```

**Fields**:
- `CurrentState string` - Current state
- `StateHistory []string` - Complete state history
- `Transitions []StateTransition` - Detailed transition records

**StateTransition Structure**:
```go
type StateTransition struct {
    From      string    // Previous state
    To        string    // New state
    Trigger   string    // What caused the transition
    Timestamp time.Time // When it happened
}
```

### MockIntentWithErrors

**Purpose**: Test error handling by injecting errors into intent lifecycle

**Use When**:
- Testing error recovery flows
- Testing error display in UI
- Testing error boundary behavior
- Verifying graceful degradation

```go
mock := NewMockIntentWithErrors()

// Inject error on Init
mock.SetInitError(errors.New("init failed"))
cmd := mock.Init()
msg := cmd().(ErrorMsg)
Expect(msg.Err.Error()).To(Equal("init failed"))

// Inject error on specific Update call
mock.SetUpdateError(errors.New("update failed"), 2) // Fail on 3rd update
mock.Update(tea.KeyMsg{}) // OK
mock.Update(tea.KeyMsg{}) // OK
cmd = mock.Update(tea.KeyMsg{}) // Fails
Expect(cmd).NotTo(BeNil())

// Track error count
Expect(mock.GetErrorCount()).To(Equal(1))

// Clear all errors
mock.ClearErrors()
```

**Fields**:
- `ErrorOnInit error` - Error to return on Init()
- `ErrorOnUpdate error` - Error to return on Update()
- `ErrorOnNthUpdate int` - Which update call should fail (0-indexed, -1 for immediate)
- `ErrorCount int` - Number of errors triggered

**ErrorMsg Type**:
```go
type ErrorMsg struct {
    Err error
}
```

---

## Writing Tests

### TDD Workflow

Always follow Red-Green-Refactor:

1. **RED**: Write failing test first
2. **GREEN**: Write minimal implementation to pass
3. **REFACTOR**: Clean up without changing behavior

### Test Structure

Use nested `Describe` and `Context` blocks for organization:

```go
var _ = Describe("BrowseTimelineIntent", func() {
    // Common setup
    var intent *BrowseTimelineIntent

    BeforeEach(func() {
        intent = NewBrowseTimelineIntent(/* deps */)
    })

    Describe("Init", func() {
        It("should set initial state", func() {
            cmd := intent.Init()
            Expect(intent.state).To(Equal(StateTimeline))
            Expect(cmd).NotTo(BeNil())
        })
    })

    Describe("Update", func() {
        Context("when in timeline state", func() {
            BeforeEach(func() {
                intent.state = StateTimeline
            })

            It("should handle Enter key", func() {
                cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
                Expect(intent.state).To(Equal(StateEventDetail))
            })

            It("should handle Escape key", func() {
                cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
                result := intent.Result()
                Expect(result.Status).To(Equal(StatusCancelled))
            })
        })

        Context("when in detail state", func() {
            // Different state tests
        })
    })

    Describe("View", func() {
        It("should render without panics", func() {
            view := intent.View()
            Expect(view).NotTo(BeEmpty())
            Expect(view).NotTo(ContainSubstring("panic"))
        })
    })
})
```

### Testing Navigation

For keyboard navigation tests:

```go
Describe("Vim-style Navigation", func() {
    BeforeEach(func() {
        env = e2e.SetupWithMemory(GinkgoT())
        env.PopulateTestData(5, 0, 0)
        env.SelectIntentByName("browse_timeline")
    })

    AfterEach(func() {
        env.Cleanup()
    })

    It("should navigate down with 'j' key", func() {
        env.PressKeyRune('j')
        view := env.GetView()
        Expect(view).NotTo(BeEmpty())
    })

    It("should not crash at boundaries", func() {
        // Try to go past beginning
        env.PressKeyRune('k')
        env.PressKeyRune('k')
        view := env.GetView()
        Expect(view).NotTo(BeEmpty())

        // Try to go past end
        for i := 0; i < 10; i++ {
            env.PressKeyRune('j')
        }
        view = env.GetView()
        Expect(view).NotTo(BeEmpty())
    })
})
```

### Testing Persistence (E2E)

For data persistence tests:

```go
Describe("Session Persistence", func() {
    BeforeEach(func() {
        env = e2e.Setup(GinkgoT()) // Use SQLite
    })

    AfterEach(func() {
        env.Cleanup()
    })

    It("should persist data after restart", func() {
        env.PopulateTestData(3, 0, 0)
        eventsBefore := env.GetEvents()
        Expect(len(eventsBefore)).To(Equal(3))

        env.SimulateRestart()

        eventsAfter := env.GetEvents()
        Expect(len(eventsAfter)).To(Equal(3))
        Expect(eventsAfter[0].ID).To(Equal(eventsBefore[0].ID))
    })
})
```

---

## Suite Organization

### Current Suite Structure

KaRiya uses the following test suite organization:

```
internal/
├── cli/
│   └── intents/
│       └── intents_suite_test.go      # Main suite (all intent tests)
├── testutil/
│   └── e2e/
│       └── e2e_suite_test.go          # E2E suite (persistence tests)
└── [other packages]/
    └── *_suite_test.go                # Per-package suites
```

### When to Create New Suite

Create a new suite file when:

1. **New subpackage**: Every package with tests needs its own suite
2. **Large spec count**: Consider splitting at > 500 specs if they can run independently
3. **Logical separation**: Different concerns warrant separate test runs
4. **Parallel execution**: Tests can benefit from running in separate processes

**Do NOT create a new suite when:**
- Tests are in the same package (use single suite)
- Tests share common setup (use shared `BeforeEach`)
- Tests have interdependencies

### File Naming Convention

| Pattern | Purpose | Location |
|---------|---------|----------|
| `{package}_suite_test.go` | Main suite file (one per package) | Same directory as code |
| `{feature}_test.go` | Feature-specific unit tests | Same directory as code |
| `{feature}_navigation_test.go` | Navigation/UI tests | `internal/cli/intents/` |
| `{feature}_e2e_test.go` | End-to-end persistence tests | `internal/testutil/e2e/` |
| `{feature}_integration_test.go` | Integration tests | Same directory as code |

### Suite File Template

```go
package yourpackage_test

import (
    "testing"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

func TestYourPackageSuite(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "YourPackage Suite")
}
```

### Avoiding Suite Sprawl

**Guidelines to prevent test organization degradation:**

1. **One suite per package**: Resist creating multiple suites in same package
2. **Use Describe blocks**: Group related tests with nested `Describe`/`Context`
3. **Use test files**: Separate features into different `*_test.go` files
4. **Shared setup**: Use `BeforeSuite`/`AfterSuite` for expensive setup
5. **Focus on behavior**: Name tests by behavior, not implementation

### Running Specific Suites

```bash
# Run all tests in a package
go test ./internal/cli/intents/...

# Run specific suite
go test ./internal/testutil/e2e/...

# Run with Ginkgo CLI for better output
ginkgo -r ./internal/cli/intents/

# Run focused tests
ginkgo -r --focus="Browse" ./internal/cli/intents/

# Skip slow tests
ginkgo -r --skip="E2E" ./...
```

---

## Common Patterns

### Testing State Machines

```go
Describe("State Transitions", func() {
    var intent *YourIntent

    BeforeEach(func() {
        intent = NewYourIntent(/* deps */)
        intent.Init()
    })

    DescribeTable("valid transitions",
        func(fromState, trigger, expectedState string) {
            intent.state = YourState(fromState)
            intent.Update(triggerToMsg(trigger))
            Expect(string(intent.state)).To(Equal(expectedState))
        },
        Entry("initial -> working on Enter", "initial", "enter", "working"),
        Entry("working -> final on confirm", "working", "confirm", "final"),
        Entry("working -> initial on cancel", "working", "cancel", "initial"),
    )
})
```

### Testing View Rendering

```go
Describe("View", func() {
    It("should render without panics", func() {
        view := intent.View()
        Expect(view).NotTo(BeEmpty())
        Expect(view).NotTo(ContainSubstring("panic"))
    })

    It("should show expected content", func() {
        view := intent.View()
        Expect(view).To(ContainSubstring("Expected Title"))
    })

    It("should not have duplicate elements", func() {
        view := intent.View()
        enterCount := strings.Count(view, "Enter")
        Expect(enterCount).To(BeNumerically("<=", 2))
    })
})
```

### Testing Error Handling

```go
Describe("Error Handling", func() {
    Context("when service returns error", func() {
        BeforeEach(func() {
            mockService.SetError(errors.New("service error"))
        })

        It("should show error message", func() {
            intent.loadData()
            view := intent.View()
            Expect(view).To(ContainSubstring("Error"))
        })

        It("should allow retry", func() {
            intent.loadData()
            intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Retry
            Expect(mockService.LoadCalled).To(Equal(2))
        })
    })
})
```

---

## Troubleshooting

### Common Issues

**Multiple Ginkgo Entry Points**
```
Error: multiple Ginkgo entry points detected
```
Solution: Ensure only one `*_suite_test.go` file per package.

**Race Conditions**
```
Error: DATA RACE detected
```
Solution: Run `go test -race ./...` and add proper synchronization.

**Test Pollution**
```
Error: unexpected state in test
```
Solution: Always use `BeforeEach`/`AfterEach` for setup/cleanup.

### Running Tests

```bash
# All tests
go test -v ./...

# With race detector
go test -race ./...

# Specific package
go test -v ./internal/cli/intents/...

# Run Ginkgo with focus
ginkgo -r --focus="Browse" ./internal/cli/intents/

# Skip slow tests
ginkgo -r --skip="E2E" ./...

# Verbose output
ginkgo -r -v ./internal/cli/intents/
```

---

## Related Documentation

- [Go Guidelines](rules/go-guidelines.md) - Go testing standards
- [TUI Developer Guide](TUI_DEVELOPER_GUIDE.md) - TUI component testing
- [Senior Engineer Guidelines](rules/senior-engineer-guidelines.md) - TDD protocol
- [Master Task Prompt](rules/master-task-prompt.md) - Development workflow

---

**Document Version**: 1.0
**Created**: 2026-01-11
**Last Updated**: 2026-01-11
