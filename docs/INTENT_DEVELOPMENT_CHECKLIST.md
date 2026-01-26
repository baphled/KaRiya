# Intent Development Checklist

**Purpose**: Ensure new intents follow KaRiya TUI standards  
**Last Updated**: 2026-01-06

---

## Pre-Development

### 1. Design Phase
- [ ] Define all states the intent will have
- [ ] Identify root state (first state user sees)
- [ ] Identify async operation states (if any)
- [ ] Map state transitions (which state leads to which)
- [ ] Design data structures (`Context`, `Result`, `Model`)

### 2. Review Standards
- [ ] Read `docs/TUI_STANDARDS.md`
- [ ] Review escape key patterns
- [ ] Check existing intent implementations for reference
- [ ] Understand the 5 state patterns (root, intermediate, async, error, completion)

---

## Implementation Checklist

### File Structure (Subdirectory - REQUIRED)

Create the subdirectory structure:
```bash
mkdir -p internal/cli/intents/{feature}/
mkdir -p internal/cli/screens/{feature}/
mkdir -p internal/cli/screens/{feature}/modals/  # If >2 modals
```

**Core Files (5 Required)**:
- [ ] Create `internal/cli/intents/{feature}/context.go` (IntentContext + Validate())
- [ ] Create `internal/cli/intents/{feature}/result.go` (Result struct)
- [ ] Create `internal/cli/intents/{feature}/constants.go` (State enum ONLY)
- [ ] Create `internal/cli/intents/{feature}/messages.go` (ALL *Msg types)
- [ ] Create `internal/cli/intents/{feature}/intent.go` (Implementation)

**Optional Files (for larger intents)**:
- [ ] Create `types.go` (Intent struct if intent.go > 300 lines)
- [ ] Create `handlers.go` (ScreenResultHandler methods)
- [ ] Create `helpers.go` (Helper methods)
- [ ] Create `filters.go` (Domain-specific filter logic)
- [ ] Create `interfaces.go` (Service interfaces)

**Test Files**:
- [ ] Create `{feature}_suite_test.go` (Ginkgo suite)
- [ ] Create `intent_test.go` (Unit tests)
- [ ] Create `{feature}_escape_test.go` (Escape behavior tests)

### Data Structures

#### Context (in context.go)
```go
// File: intents/{feature}/context.go
package myfeature

type IntentContext struct {
    Items   []*domain.Item
    Service ItemService  // Use interface for dependency injection
}

func (c *IntentContext) Validate() error {
    if c.Items == nil {
        c.Items = make([]*domain.Item, 0)
    }
    return nil
}

// Domain types can also be here
type Filters struct {
    SearchText string
    Categories []string
}
```

- [ ] Define `IntentContext` struct in `context.go`
- [ ] Add `Validate()` method
- [ ] Document all fields with comments
- [ ] Include domain types (Filters, etc.) if needed

#### Result (in result.go)
```go
// File: intents/{feature}/result.go
package myfeature

type Result struct {
    SelectedItem *domain.Item
    FinalFilters *Filters
}
```

- [ ] Define `Result` struct in `result.go`
- [ ] Include all data needed by parent/caller
- [ ] Document what each field represents

#### Intent Struct (in intent.go or types.go)
```go
// File: intents/{feature}/types.go (or intent.go if small)
package myfeature

type Intent struct {
    *intents.BaseIntent
    context      *IntentContext
    state        State
    active       bool
    result       *intents.IntentResult[*Result]
    activeScreen screens.Screen
    // Modal fields...
    modalRegistry *intents.ModalRegistry
}
```

- [ ] Define `Intent` struct (in `intent.go` or `types.go`)
- [ ] Embed `*intents.BaseIntent`
- [ ] Include `context`, `state`, `active`, `result` fields
- [ ] Add `modalRegistry` for unified modal handling
- [ ] Add state-specific fields as needed

### State Machine

#### State Enum (in constants.go)
```go
// File: intents/{feature}/constants.go
package myfeature

type State string

const (
    StateList   State = "list"
    StateDetail State = "detail"
)
```

- [ ] Define State type in `constants.go` (NOT in intent.go)
- [ ] Define all states as constants
- [ ] Use descriptive state names
- [ ] Document state transitions in comments
- [ ] Do NOT put error constants here (use Go errors package)

### Intent Interface Implementation

#### Init Method
```go
func (i *YourIntentModel) Init() tea.Cmd {
    i.active = true
    i.state = StateInitial
    // Initialize components
    return nil
}
```

- [ ] Set `active = true`
- [ ] Set initial state
- [ ] Initialize any UI components
- [ ] Return appropriate command (or nil)

#### Update Method
```go
func (i *YourIntentModel) Update(msg tea.Msg) tea.Cmd {
    if !i.active {
        return nil
    }
    
    switch i.state {
    case StateInitial:
        return i.updateInitial(msg)
    case StateIntermediate:
        return i.updateIntermediate(msg)
    // ... other states ...
    }
    return nil
}
```

- [ ] Check `active` flag first
- [ ] Use switch on `i.state`
- [ ] Delegate to state-specific update methods
- [ ] Return appropriate command

#### View Method
```go
func (i *YourIntentModel) View() string {
    switch i.state {
    case StateInitial:
        return i.viewInitial()
    case StateIntermediate:
        return i.viewIntermediate()
    // ... other states ...
    }
    return ""
}
```

- [ ] Use switch on `i.state`
- [ ] Delegate to state-specific view methods
- [ ] Return empty string for unknown states

#### Result Method
```go
func (i *YourIntentModel) Result() *IntentResult[interface{}] {
    if i.result == nil {
        return nil
    }
    
    return &IntentResult[interface{}]{
        Status:   i.result.Status,
        Data:     i.result.Data,
        Error:    i.result.Error,
        Metadata: i.result.Metadata,
    }
}
```

- [ ] Return nil if no result set
- [ ] Convert typed result to `interface{}`
- [ ] Preserve status, data, error, metadata

### State Update Methods

For EACH state, implement:

#### Root State Pattern
```go
func (i *YourIntentModel) updateRootState(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "esc":
            i.setCancelled()  // Cancel and return to menu
            return nil
        case "m":
            i.setCancelled()  // Same as esc
            return nil
        case "q", "ctrl+c":
            i.setCancelled()  // Quit
            return nil
        // ... state-specific keys ...
        }
    }
    return nil
}
```

- [ ] Add `esc` handler (cancel intent)
- [ ] Add `m` handler (same as esc)
- [ ] Add `q`/`ctrl+c` handler (quit)
- [ ] Add state-specific navigation keys
- [ ] Transition to next state on action

#### Intermediate State Pattern
```go
func (i *YourIntentModel) updateIntermediateState(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "esc":
            i.state = PreviousState  // Go back
            return nil
        case "m":
            i.setCancelled()  // Main menu
            return nil
        case "q", "ctrl+c":
            i.setCancelled()  // Quit
            return nil
        // ... state-specific keys ...
        }
    }
    return nil
}
```

- [ ] Add `esc` handler (go to previous state)
- [ ] Add `m` handler (cancel intent)
- [ ] Add `q`/`ctrl+c` handler (quit)
- [ ] Add state-specific navigation keys
- [ ] Transition to next/previous state appropriately

#### Async Operation State Pattern
```go
func (i *YourIntentModel) updateAsyncState(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "esc":
            // Let operation complete in background
            i.state = PreviousState
            return nil
        case "m":
            i.setCancelled()  // Cancel immediately
            return nil
        case "q", "ctrl+c":
            i.setCancelled()  // Quit immediately
            return nil
        }
    case YourProgressMsg:
        // Handle progress updates
        return nil
    case YourCompleteMsg:
        i.state = CompleteState
        return nil
    }
    return nil
}
```

- [ ] Add `esc` handler (background completion)
- [ ] Add `m` handler (immediate cancel)
- [ ] Add `q`/`ctrl+c` handler (immediate quit)
- [ ] Handle progress messages
- [ ] Handle completion messages
- [ ] Transition to completion state on success

#### Error State Pattern
```go
func (i *YourIntentModel) updateErrorState(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "esc":
            i.state = PreviousState
            // Keep error visible!
            return nil
        case "m":
            i.setCancelled()
            return nil
        case "r":
            // Retry logic
            return i.retryOperation()
        }
    }
    return nil
}
```

- [ ] Add `esc` handler (go back, keep error visible)
- [ ] Add `m` handler (cancel intent)
- [ ] Add `r` handler for retry (if applicable)
- [ ] DO NOT clear error on escape

### State View Methods

For EACH state, implement:

```go
func (i *YourIntentModel) viewStateName() string {
    var content strings.Builder
    
    // Build view content
    content.WriteString("State Title\n\n")
    content.WriteString("Content here...\n")
    
    // Footer with key options
    footer := footerStyle.Render(
        "Enter: Confirm | Esc: Back | m: Main menu | q: Quit"
    )
    
    return lipgloss.JoinVertical(lipgloss.Left, content.String(), footer)
}
```

- [ ] Build content using `strings.Builder`
- [ ] Apply consistent styling with lipgloss
- [ ] Always include footer with available keys
- [ ] Show appropriate keys for state type:
  - Root: `"Esc: Cancel | m: Main menu | q: Quit"`
  - Intermediate: `"Esc: Back | m: Main menu | q: Quit"`
  - Async: `"Esc: Background | m: Cancel | q: Quit"`
  - Error: `"Esc: Back (error visible) | m: Main menu | r: Retry"`

### Helper Methods

```go
func (i *YourIntentModel) setCompleted() {
    i.result = &IntentResult[*YourIntentResult]{
        Status: Completed,
        Data:   &YourIntentResult{/* ... */},
        Metadata: map[string]interface{}{
            "timestamp": time.Now(),
        },
    }
    i.active = false
}

func (i *YourIntentModel) setCancelled() {
    i.result = &IntentResult[*YourIntentResult]{
        Status: Cancelled,
    }
    i.active = false
}

func (i *YourIntentModel) setFailed(code, message string, cause error) {
    i.result = &IntentResult[*YourIntentResult]{
        Status: Failed,
        Error: &IntentError{
            Code:    code,
            Message: message,
            Cause:   cause,
        },
    }
    i.active = false
}
```

- [ ] Implement `setCompleted()` method
- [ ] Implement `setCancelled()` method
- [ ] Implement `setFailed()` method
- [ ] Set `active = false` in all helpers

---

## Testing Checklist

### Unit Tests (`{intent_name}_test.go`)

For EACH state:
- [ ] Test state initialization
- [ ] Test valid state transitions
- [ ] Test invalid state transitions
- [ ] Test view rendering
- [ ] Test data validation
- [ ] Test error handling

### Escape Behavior Tests (`{intent_name}_escape_test.go`)

For EACH state:
- [ ] Test `esc` key behavior
  - Root: should cancel intent
  - Intermediate: should go back to previous state
  - Async: should allow background completion
- [ ] Test `m` key behavior (should always cancel intent)
- [ ] Test `q` key behavior (should always cancel intent)
- [ ] Test that view shows correct footer with key options
- [ ] Test that errors stay visible on escape (if error state)

Example test structure:
```go
var _ = Describe("YourIntent - Escape Key Behavior", func() {
    Describe("StateInitial (Root)", func() {
        It("should cancel intent when escape is pressed", func() {
            // Test implementation
        })
        
        It("should cancel intent when 'm' is pressed", func() {
            // Test implementation
        })
        
        It("should show 'm' key in footer", func() {
            // Test implementation
        })
    })
    
    // ... repeat for all states ...
})
```

### Integration Tests
- [ ] Test complete happy path (all states)
- [ ] Test cancellation from each state
- [ ] Test error recovery paths
- [ ] Test async operation completion

---

## Documentation Checklist

### Code Documentation
- [ ] Add godoc comments to all exported types
- [ ] Document state machine transitions
- [ ] Document context fields
- [ ] Document result fields
- [ ] Add examples in comments

### Intent-Specific Docs
- [ ] Create user guide section for intent
- [ ] Document all keyboard shortcuts
- [ ] Document error messages and recovery
- [ ] Add workflow examples

### Update Standards
- [ ] Add intent to `docs/TUI_STANDARDS.md` implementation table
- [ ] Document any new patterns introduced
- [ ] Update `docs/USER_GUIDE_NAVIGATION.md` with new workflows

---

## Pre-Commit Checklist

### Code Quality
- [ ] Run `go fmt ./...`
- [ ] Run `go vet ./...`
- [ ] Run `golangci-lint run ./...` (if available)
- [ ] No compiler warnings

### Tests
- [ ] All unit tests pass: `go test ./internal/cli/intents/...`
- [ ] All escape tests pass: `ginkgo --focus="YourIntent.*Escape"`
- [ ] Run with race detector: `go test -race ./internal/cli/intents/...`
- [ ] No test failures introduced
- [ ] Test coverage ≥ 80%

### Build
- [ ] Application builds: `go build ./cmd/cli`
- [ ] No build errors or warnings
- [ ] Binary runs without crashes

### Manual Testing
- [ ] Test all state transitions manually
- [ ] Test escape key from every state
- [ ] Test 'm' key from every state
- [ ] Test error recovery paths
- [ ] Test async operations (if any)
- [ ] Verify footer shows correct keys

---

## Review Checklist

### Code Review
- [ ] Follows established patterns (see existing intents)
- [ ] No code duplication
- [ ] Proper error handling
- [ ] Clear variable names
- [ ] Appropriate comments

### UX Review
- [ ] Consistent with other intents
- [ ] Clear navigation flow
- [ ] Helpful error messages
- [ ] Keyboard shortcuts work as expected
- [ ] Footer always shows available keys

### Testing Review
- [ ] Tests cover all states
- [ ] Tests cover all error paths
- [ ] Tests are clear and maintainable
- [ ] No flaky tests
- [ ] Good test descriptions

---

## Common Pitfalls to Avoid

### ❌ Don't Do This
1. **Forgetting `active` check** in `Update()` method
2. **Missing escape handlers** in any state
3. **Clearing errors** when navigating back
4. **No async escape handler** in long-running operations
5. **Inconsistent footer text** across states
6. **Not setting `active = false`** in result helpers
7. **Type asserting `Result()` output** (should return `interface{}`)

### ✅ Do This Instead
1. Always check `if !i.active` first in `Update()`
2. Every state has `esc`, `m`, and `q` handlers
3. Keep errors visible when pressing escape
4. Async states allow background completion on escape
5. Use standard footer format from examples
6. Set `active = false` in all result helpers
7. Convert typed result to `IntentResult[interface{}]`

---

## Reference Examples

### New Subdirectory Structure (RECOMMENDED)
See: `internal/cli/intents/browse_timeline/` - Complete subdirectory implementation with:
- 5 core files (context.go, result.go, constants.go, messages.go, intent.go)
- Optional files (types.go, handlers.go, helpers.go, filters.go, interfaces.go)
- Feature-specific screens (`screens/timeline/`)
- Feature-specific modals (`screens/timeline/modals/`)

### Legacy Flat Structure (To Be Migrated)
See: `internal/cli/intents/capture_event_intent.go`
See: `internal/cli/intents/generate_cv_intent.go`
See: `internal/cli/intents/burst_management_intent.go`

**Note**: Legacy intents should be migrated to subdirectory structure.
See: `docs/guides/INTENT_MIGRATION_TO_SUBDIRECTORY.md`

---

## Quick Start Template

Copy this template to start a new intent:

```bash
# Create subdirectory structure
mkdir -p internal/cli/intents/my_feature
mkdir -p internal/cli/screens/my_feature
mkdir -p internal/cli/screens/my_feature/modals  # If >2 modals

# Create core files (5 required)
touch internal/cli/intents/my_feature/context.go
touch internal/cli/intents/my_feature/result.go
touch internal/cli/intents/my_feature/constants.go
touch internal/cli/intents/my_feature/messages.go
touch internal/cli/intents/my_feature/intent.go

# Create optional files (for larger intents)
touch internal/cli/intents/my_feature/types.go      # Intent struct
touch internal/cli/intents/my_feature/handlers.go   # ScreenResultHandler
touch internal/cli/intents/my_feature/helpers.go    # Helper methods

# Create test files
touch internal/cli/intents/my_feature/my_feature_suite_test.go
touch internal/cli/intents/my_feature/intent_test.go

# Copy from reference implementation
# RECOMMENDED: intents/browse_timeline/ - Complete subdirectory example
#              screens/timeline/ - Screen and modal patterns
```

**Reference Implementation**: `intents/browse_timeline/` demonstrates the new subdirectory pattern.

---

## Completion Criteria

An intent is considered complete when:
- ✅ All states have escape/m/q handlers
- ✅ All view methods show correct footers
- ✅ All unit tests pass (≥80% coverage)
- ✅ All escape behavior tests pass
- ✅ Application builds without errors
- ✅ Manual testing confirms all workflows work
- ✅ Documentation is updated
- ✅ Code review approved

---

**Last Updated**: 2026-01-24  
**Based on**: Subdirectory Structure Migration + Escape Key Standardization  
**Reference Implementation**: `intents/browse_timeline/` and `screens/timeline/`  
**Migration Guide**: `docs/guides/INTENT_MIGRATION_TO_SUBDIRECTORY.md`
