---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Phase 1 Implementation Guide: Foundation & Core Infrastructure

**Duration**: 1.5 weeks
**Goal**: Establish intent boundary contract types, IntentRouter, and refactor root model
**Status**: Ready for implementation

---

## Overview

Phase 1 consists of three main tasks:

1. **Task 1.1**: Implement intent boundary contract types (`internal/cli/intents/contract.go`)
2. **Task 1.2**: Implement IntentRouter (`internal/cli/intents/router.go`)
3. **Task 1.3**: Refactor root model to use IntentRouter (`internal/cli/app/app.go`)

Each task builds on the previous one and includes concrete Go code scaffolds ready for implementation.

---

## Task 1.1: Intent Boundary Contract Types

### File Location
`internal/cli/intents/contract.go`

### Overview
This file defines the core types that all intents must implement and use for communication. It establishes the type-safe boundary contract.

### Key Types to Implement

```go
// IntentStatus - enumeration with four states
type IntentStatus string
const (
    StatusCompleted IntentStatus = "completed"
    StatusCancelled IntentStatus = "cancelled"
    StatusFailed    IntentStatus = "failed"
    StatusPartial   IntentStatus = "partial"
)

// IntentError - error with code, message, and cause
type IntentError struct {
    Code    string
    Message string
    Cause   error
}

// IntentResult[T] - generic result with status, data, error, metadata
type IntentResult[T any] struct {
    Status   IntentStatus
    Data     T
    Error    *IntentError
    Metadata map[string]interface{}
}

// Intent - interface all intents must implement
type Intent interface {
    Init(ctx context.Context) tea.Cmd
    Update(msg tea.Msg) (Intent, tea.Cmd)
    View() string
    Result() *IntentResult[interface{}]
}

// ModalEditResult[T] - result from modal sub-flows
type ModalEditResult[T any] struct {
    Original T
    Modified T
    Accepted bool
    Changes  map[string]interface{}
}

// GlobalContext - non-mutable shared metadata
type GlobalContext struct {
    UserPreferences   map[string]interface{}
    LastSelectedItems map[string]interface{}
    AppConfig         map[string]interface{}
}
```

### Helper Methods Required

For `IntentResult[T]`:
- `WithMetadata(key string, value interface{})` - add metadata
- `GetMetadata(key string)` - retrieve metadata
- `IsSuccess()` - check if completed
- `IsPartialSuccess()` - check if partial
- `IsFailed()` - check if failed
- `IsCancelled()` - check if cancelled

For `ModalEditResult[T]`:
- `WithChange(field string, value interface{})` - track change
- `Accept()` - mark as accepted
- `Reject()` - mark as rejected

For `GlobalContext`:
- `GetPreference(key string)` - get user preference
- `GetLastSelectedItem(key string)` - get transient state
- `SetLastSelectedItem(key string, value interface{})` - set transient state
- `GetConfig(key string)` - get app config

### Implementation Checklist

- [ ] Create file `internal/cli/intents/contract.go`
- [ ] Implement `IntentStatus` constants
- [ ] Implement `IntentError` struct with Error() method
- [ ] Implement `IntentResult[T]` struct with all helper methods
- [ ] Implement `ModalEditResult[T]` struct with all helper methods
- [ ] Implement `Intent` interface
- [ ] Implement `GlobalContext` struct with helper methods
- [ ] Write comprehensive unit tests (target >95% coverage)
- [ ] Verify with `go fmt`, `go vet`, `golangci-lint`

---

## Task 1.2: IntentRouter Implementation

### File Location
`internal/cli/intents/router.go`

### Overview
This file implements the IntentRouter, which orchestrates intent activation, switching, and result handling. It maintains the navigation history and manages back navigation.

### Key Methods to Implement

```go
type IntentRouter struct {
    intents       map[string]Intent
    current       Intent
    history       []Intent
    ctx           context.Context
    resultHandler func(*IntentResult[interface{}])
    mu            sync.RWMutex
}

// Constructor
func NewIntentRouter(ctx context.Context) *IntentRouter

// Core operations
func (ir *IntentRouter) Register(name string, intent Intent) error
func (ir *IntentRouter) Activate(intentName string) tea.Cmd
func (ir *IntentRouter) Update(msg tea.Msg) tea.Cmd
func (ir *IntentRouter) View() string
func (ir *IntentRouter) Back() tea.Cmd

// Accessors (for testing)
func (ir *IntentRouter) GetCurrentIntent() Intent
func (ir *IntentRouter) GetHistory() []Intent
func (ir *IntentRouter) SetResultHandler(handler func(*IntentResult[interface{}]))
```

### Behavior Specifications

**Register**:
- Validate intent doesn't already exist
- Store in intents map
- Return error if duplicate

**Activate**:
- Find intent by name
- Push current intent to history
- Set new intent as current
- Call Init() on new intent
- Return init command

**Update**:
- Delegate message to current intent
- Update current if intent changed
- Check for result and call handler
- Return command from intent

**View**:
- Return current intent's view
- Handle nil current gracefully

**Back**:
- Pop previous intent from history
- Set as current
- Call Init() on restored intent
- Return init command
- Return nil if no history

### Thread Safety Requirements

- Use `sync.RWMutex` to protect concurrent access
- Lock when reading/writing intents map
- Lock when modifying history
- Lock when changing current intent

### Implementation Checklist

- [ ] Create file `internal/cli/intents/router.go`
- [ ] Implement `IntentRouter` struct with mutex
- [ ] Implement `NewIntentRouter()` constructor
- [ ] Implement `Register()` with validation
- [ ] Implement `Activate()` with history management
- [ ] Implement `Update()` with result handling
- [ ] Implement `View()` with nil check
- [ ] Implement `Back()` with history restoration
- [ ] Implement `GetCurrentIntent()` and `GetHistory()`
- [ ] Implement `SetResultHandler()`
- [ ] Write unit tests for all methods
- [ ] Write integration tests for workflows
- [ ] Test thread safety with `go test -race`
- [ ] Verify with `go fmt`, `go vet`, `golangci-lint`

---

## Task 1.3: Root Model Refactoring

### File Location
`internal/cli/app/app.go`

### Overview
Refactor the root Bubble Tea model to use IntentRouter instead of inline state management. The root model becomes a thin wrapper around the router that handles global shortcuts and delegates all intent-specific logic.

### Key Changes

**Before (Current)**:
- Inline state for each intent
- Mixed message handling
- Direct model-to-model navigation
- Type-unsafe state management

**After (Phase 1)**:
- Single `router *IntentRouter` field
- Router delegates all intent logic
- Global shortcuts at root level
- Type-safe via IntentResult

### Root Model Structure

```go
type Model struct {
    router        *intents.IntentRouter
    globalContext *intents.GlobalContext
    lastError     error
    theme         *Theme
    isQuitting    bool
}

func NewModel(ctx context.Context) *Model
func (m *Model) Init() tea.Cmd
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m *Model) View() string
func (m *Model) handleIntentResult(result *intents.IntentResult[interface{}])
```

### Global Shortcuts to Handle

- `Ctrl+C` - Quit application
- `Ctrl+H` - Show help
- `Ctrl+Home` - Return to main menu
- `Esc` - Back navigation (or quit if no history)

### Implementation Checklist

**Analysis**:
- [ ] Review current `internal/cli/app/app.go`
- [ ] Document all inline state fields
- [ ] Document all message types
- [ ] Document current Update() flow
- [ ] Document current View() flow
- [ ] Identify all global shortcuts
- [ ] Identify all global UI chrome

**Refactoring**:
- [ ] Create new Model struct with router and globalContext
- [ ] Keep only global UI state (theme, isQuitting, lastError)
- [ ] Implement NewModel() constructor
  - [ ] Create IntentRouter
  - [ ] Create GlobalContext
  - [ ] Register placeholder intents
  - [ ] Set result handler
- [ ] Implement Init() to activate initial intent
- [ ] Implement Update() to handle global shortcuts
  - [ ] Ctrl+C → quit
  - [ ] Ctrl+H → show help
  - [ ] Ctrl+Home → main menu
  - [ ] Esc → back or quit
  - [ ] Delegate other messages to router
- [ ] Implement View() to combine intent view with global chrome
- [ ] Implement handleIntentResult() to process intent results
- [ ] Implement helper methods for UI chrome rendering

**Testing**:
- [ ] Write unit tests for Model
- [ ] Write integration tests for navigation
- [ ] Verify all existing tests still pass
- [ ] Test global shortcuts
- [ ] Test error handling
- [ ] Verify no regressions

---

## Code Quality Standards

### For All Code
- [ ] Follows Go conventions and idioms
- [ ] Properly formatted with `gofmt`
- [ ] Passes `golangci-lint`
- [ ] Comprehensive comments on exported functions
- [ ] >90% test coverage
- [ ] No unused code or imports
- [ ] Clear, descriptive variable names

### For Tests
- [ ] Use Ginkgo for organization
- [ ] Use Gomega for assertions
- [ ] Test both happy path and error cases
- [ ] Test edge cases
- [ ] Mock external dependencies
- [ ] >90% code coverage

---

## Integration Testing Scenarios

### Scenario 1: Intent Activation
```
1. Create router
2. Register intent
3. Activate intent
4. Verify intent is current
5. Verify Init was called
```

### Scenario 2: Back Navigation
```
1. Create router
2. Register two intents
3. Activate intent 1
4. Activate intent 2 (should push intent 1 to history)
5. Call Back()
6. Verify intent 1 is current
7. Verify intent 2 is in history
```

### Scenario 3: Result Handling
```
1. Create router with result handler
2. Register intent that returns result
3. Activate intent
4. Update with message that triggers result
5. Verify result handler was called
6. Verify result status is correct
```

### Scenario 4: Global Shortcuts
```
1. Create root model
2. Send Ctrl+C message
3. Verify isQuitting is true
4. Send Esc message
5. Verify Back() was called on router
```

---

## Validation Checklist

Before Phase 1 is complete:

- [ ] All three files compile without errors
- [ ] All unit tests pass
- [ ] All integration tests pass
- [ ] Code coverage >90%
- [ ] `go fmt` passes (no formatting issues)
- [ ] `go vet` passes (no type issues)
- [ ] `golangci-lint` passes (no linting issues)
- [ ] `go test -race` passes (no race conditions)
- [ ] Code review approved
- [ ] Documentation updated in AGENTS.md
- [ ] No breaking changes to existing CLI
- [ ] All existing tests still pass

---

## Success Criteria

Phase 1 is complete when:

✅ Intent boundary contract types are fully implemented and tested
✅ IntentRouter is fully functional with comprehensive tests
✅ Root model successfully delegates to IntentRouter
✅ All tests pass with >90% coverage
✅ No breaking changes to existing CLI
✅ Code review approved
✅ Team is ready to begin Phase 2 (CaptureEventIntent)

---

## Questions During Implementation

### Q1: How should existing components integrate with intents?
**A**: Intents can reuse existing components. Refactor components to fit intent pattern as needed.

### Q2: Should all intents be registered immediately?
**A**: Register placeholder intents for now. Implement actual intents in Phase 2+.

### Q3: How to handle persistent state across intent sessions?
**A**: Use GlobalContext for session state. Persist to database for permanent storage.

### Q4: What about error recovery in intents?
**A**: Intents return IntentError in result. Root model displays and handles appropriately.

### Q5: Can intents communicate with each other?
**A**: No. All communication is through IntentResult and router callbacks.

---

## Next Phase

Once Phase 1 is complete:

1. Create MainMenuIntent stub
2. Create CaptureEventIntent template
3. Implement all state transitions and views
4. Use as template for other intents
5. Begin Phase 2 implementation

---

*Follow this guide step-by-step. Use the scaffold code as starting points. Refer to the PRD for acceptance criteria. Good luck!*

