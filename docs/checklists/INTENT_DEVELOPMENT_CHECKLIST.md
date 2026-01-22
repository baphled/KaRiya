# Intent Development Checklist

**Purpose**: Mandatory checklist for all intent development  
**Enforcement**: `make check-intent-architecture` (automated)  
**Status**: REQUIRED - Must pass before commit

---

## Pre-Commit Validation

Run this command **before every commit**:

```bash
make check-intent-architecture
```

This checklist is **automatically enforced** by the pre-commit hook. Violations will **block your commit**.

---

## Architecture Requirements (BLOCKING)

### 1. BaseIntent Embedding ✅

**Rule**: All intents MUST embed `*BaseIntent`

```go
type MyIntent struct {
    *BaseIntent  // REQUIRED
    
    // ... other fields
}
```

**Why**: Provides terminal info, theme, logo, and help modal infrastructure.

**Checked by**: `check-intent-architecture.sh` (Check #6)

---

### 2. Typed State Enum ✅

**Rule**: State fields MUST use typed enums, not raw strings

```go
// REQUIRED: Define typed state
type MyIntentState string

const (
    StateList   MyIntentState = "list"
    StateDetail MyIntentState = "detail"
)

type MyIntent struct {
    *BaseIntent
    state MyIntentState  // Typed, not string
}
```

**Why**: Type safety prevents typos and enables compile-time checks.

**Checked by**: `check-intent-architecture.sh` (Check #1)

---

### 3. Flattened State Model ✅

**Rule**: Intent state MUST be flattened directly into the intent struct, not wrapped in a separate model

```go
// ❌ BAD: Wrapped state
type MyIntent struct {
    *BaseIntent
    state *MyIntentModel  // WRONG
}

// ✅ GOOD: Flattened state
type MyIntent struct {
    *BaseIntent
    
    // State machine
    state      MyIntentState
    active     bool
    result     *IntentResult[*MyResult]
    
    // Screens
    listScreen   *myfeature.ListScreen
    detailScreen *myfeature.DetailScreen
    
    // Modals
    deleteModal *feedback.ConfirmModal
    
    // Data (flattened)
    items         []*Item
    selectedItem  *Item
    filters       *FilterConfig
}
```

**Why**: Reduces indirection, matches architecture guide patterns.

**Checked by**: `check-intent-architecture.sh` (Check #2)

---

### 4. Context Usage ✅

**Rule**: Use `i.getContext()` instead of `context.Background()`

```go
// ❌ BAD
ctx := context.Background()

// ✅ GOOD
ctx := i.getContext()
```

**Why**: Enables proper cancellation and timeout support.

**Checked by**: `check-intent-architecture.sh` (Check #3)

---

### 5. No Dead Code ✅

**Rule**: Remove code marked "should not be reached" or create cleanup task

```go
// ❌ BAD
case *timeline.TimelineEventDetailScreen:
    // LEGACY: This case should not be reached
    return screen.View()
```

**Action**: Either remove the code or create a task: `make new-bug BUG="Remove legacy code"`

**Checked by**: `check-intent-architecture.sh` (Check #4)

---

### 6. ScreenResultHandler Implementation ✅

**Rule**: Intents using screens MUST implement `ScreenResultHandler`

```go
// REQUIRED: Interface compliance check
var _ ScreenResultHandler = (*MyIntent)(nil)

// Implement all methods:
func (i *MyIntent) HandleCancel(*screens.CancelResult) tea.Cmd
func (i *MyIntent) HandleNavigate(*screens.NavigateResult) tea.Cmd
func (i *MyIntent) HandleSubmit(*screens.SubmitResult) tea.Cmd
func (i *MyIntent) HandleError(*screens.ErrorResult) tea.Cmd
```

**Checked by**: `check-intent-architecture.sh` (Check #9)

---

### 7. Modal Overlay Pattern ✅

**Rule**: Use `behaviors.RenderModalOverlay()` for modal rendering

```go
// ✅ GOOD
func (i *MyIntent) View() string {
    baseView := i.listScreen.View()
    
    if i.deleteModal != nil && i.deleteModal.IsVisible() {
        return behaviors.RenderModalOverlay(i.deleteModal, baseView)
    }
    
    return baseView
}
```

**Checked by**: `check-intent-architecture.sh` (Check #8)

---

## Code Quality (WARNINGS)

### 8. Godoc Completeness ⚠️

**Rule**: All exported functions should have godoc comments

```go
// HandleCancel handles screen cancellation (back/escape).
//
// Implements ScreenResultHandler interface.
func (i *MyIntent) HandleCancel(result *screens.CancelResult) tea.Cmd {
    // ...
}
```

**Why**: Documentation is code. Complete sentences required.

**Checked by**: `check-intent-architecture.sh` (Check #7), `golangci-lint`

---

### 9. Screen Management Pattern ⚠️

**Recommendation**: Declare typed screen fields for clarity

```go
// ⚠️ ACCEPTABLE but not ideal
type MyIntent struct {
    *BaseIntent
    activeScreen screens.Screen  // Generic
}

// ✅ BETTER
type MyIntent struct {
    *BaseIntent
    listScreen   *myfeature.ListScreen
    detailScreen *myfeature.DetailScreen
    activeScreen screens.Screen  // Points to one of the above
}
```

**Checked by**: `check-intent-architecture.sh` (Check #5)

---

## Testing Requirements

### E2E Tests

**Rule**: All intents MUST have E2E tests with Happy and Sad paths

```bash
# Check file exists
ls internal/cli/intents/my_intent_e2e_test.go

# Must contain:
grep "Happy Paths" internal/cli/intents/my_intent_e2e_test.go
grep "Sad Paths" internal/cli/intents/my_intent_e2e_test.go
```

**Checked by**: `scripts/check-patterns-strict.sh` (Check #7)

---

### Coverage Requirement

**Rule**: Package coverage ≥ 95%

```bash
make coverage
```

**Checked by**: `scripts/check-patterns-strict.sh` (Check #9)

---

## Architectural Constraints

### Layer Dependencies

**NEVER** violate these rules (enforced by `.golangci.yml`):

| Layer | Can Import | NEVER Import |
|-------|------------|--------------|
| `intents/` | screens, uikit, behaviors, components | - |
| `screens/` | uikit, behaviors | **intents** (FORBIDDEN) |
| `uikit/` | themes only | **screens, intents** (FORBIDDEN) |
| `behaviors/` | uikit, themes | **screens, intents** (FORBIDDEN) |

**Checked by**: `golangci-lint` (depguard rules)

---

### Forms Architecture

**Rule**: `huh` can ONLY be imported in `forms/` package

```go
// ❌ NEVER in intents
import "github.com/charmbracelet/huh"

// ✅ ALWAYS use wrappers
import "github.com/baphled/kariya/internal/cli/models"
import "github.com/baphled/kariya/internal/cli/forms"
```

**Checked by**: `golangci-lint` (depguard rules)

---

## Documentation Requirements

### New Intents

When creating a new intent, update:

1. **Workflow documentation** (`docs/workflows/`)
2. **Intent patterns** (`docs/development/INTENT_PATTERNS_LIBRARY.md`)
3. **Architecture guide** (`docs/INTENT_ARCHITECTURE_GUIDE.md` - add to reference list)

**Checked by**: `scripts/check-patterns-strict.sh` (Check #8)

---

## Quick Reference Commands

| Command | Purpose |
|---------|---------|
| `make check-intent-architecture` | Run all architecture checks |
| `make check-compliance` | Full compliance (includes architecture) |
| `make golangci-lint` | Comprehensive linting |
| `make check-patterns-strict` | Strict pattern validation |
| `make pre-commit` | Quick pre-commit checks |

---

## Example: Compliant Intent Structure

```go
package intents

import (
    "github.com/baphled/kariya/internal/cli/behaviors"
    "github.com/baphled/kariya/internal/cli/screens"
    "github.com/baphled/kariya/internal/cli/screens/myfeature"
    "github.com/baphled/kariya/internal/cli/uikit/feedback"
    tea "github.com/charmbracelet/bubbletea"
)

// MyIntentState defines the state machine states.
type MyIntentState string

const (
    StateList   MyIntentState = "list"
    StateDetail MyIntentState = "detail"
)

// Ensure interface compliance.
var _ ScreenResultHandler = (*MyIntent)(nil)

// MyIntent implements the Intent interface for [description].
type MyIntent struct {
    *BaseIntent
    
    // State machine
    state  MyIntentState
    active bool
    result *IntentResult[*MyResult]
    
    // Screens
    listScreen   *myfeature.ListScreen
    detailScreen *myfeature.DetailScreen
    
    // Modals
    deleteModal *feedback.ConfirmModal
    
    // Data
    items         []*Item
    selectedItem  *Item
}

// NewMyIntent creates a new MyIntent.
func NewMyIntent(ctx context.Context, items []*Item) (*MyIntent, error) {
    return &MyIntent{
        BaseIntent: NewBaseIntent(),
        state:      StateList,
        active:     true,
        items:      items,
    }, nil
}

// Init initializes the intent.
func (i *MyIntent) Init() tea.Cmd {
    i.listScreen = myfeature.NewListScreen(i.items)
    i.listScreen.SetTerminalInfo(i.GetTerminalInfo())
    i.listScreen.SetTheme(i.Theme())
    return nil
}

// Update processes messages.
func (i *MyIntent) Update(msg tea.Msg) tea.Cmd {
    if !i.active {
        return nil
    }
    
    // Modal handling
    if i.deleteModal != nil && i.deleteModal.IsVisible() {
        model, cmd := i.deleteModal.Update(msg)
        i.deleteModal = model.(*feedback.ConfirmModal)
        return cmd
    }
    
    // Screen delegation
    switch i.state {
    case StateList:
        cmd, result := i.listScreen.Update(msg)
        if result != nil {
            return i.handleScreenResult(result)
        }
        return cmd
    }
    
    return nil
}

// View renders the intent.
func (i *MyIntent) View() string {
    var baseView string
    
    switch i.state {
    case StateList:
        baseView = i.listScreen.View()
    case StateDetail:
        baseView = i.detailScreen.View()
    }
    
    // Overlay modals
    if i.deleteModal != nil && i.deleteModal.IsVisible() {
        return behaviors.RenderModalOverlay(i.deleteModal, baseView)
    }
    
    return baseView
}

// Result returns the intent result.
func (i *MyIntent) Result() *IntentResult[interface{}] {
    // Type erasure for interface compliance
    if i.result == nil {
        return nil
    }
    return &IntentResult[interface{}]{
        Status: i.result.Status,
        Data:   i.result.Data,
    }
}

// HandleCancel handles screen cancellation.
func (i *MyIntent) HandleCancel(result *screens.CancelResult) tea.Cmd {
    i.setCancelled()
    return nil
}

// HandleNavigate handles screen navigation.
func (i *MyIntent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
    // State transitions here
    return nil
}

// HandleSubmit handles form submission.
func (i *MyIntent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
    return nil
}

// HandleError handles errors from screens.
func (i *MyIntent) HandleError(result *screens.ErrorResult) tea.Cmd {
    return nil
}

// getContext returns context for service calls.
func (i *MyIntent) getContext() context.Context {
    // Use proper context with cancellation
    return i.BaseIntent.GetContext()
}

// setCancelled marks the intent as cancelled.
func (i *MyIntent) setCancelled() {
    i.result = &IntentResult[*MyResult]{
        Status: Cancelled,
    }
    i.active = false
}
```

---

## Enforcement Summary

| Check | Type | Severity | Script |
|-------|------|----------|--------|
| BaseIntent embedding | Automated | 🔴 BLOCKING | check-intent-architecture.sh |
| Typed state enum | Automated | 🔴 BLOCKING | check-intent-architecture.sh |
| Flattened state | Automated | 🔴 BLOCKING | check-intent-architecture.sh |
| Context usage | Automated | 🔴 BLOCKING | check-intent-architecture.sh |
| Dead code | Automated | 🔴 BLOCKING | check-intent-architecture.sh |
| ScreenResultHandler | Automated | 🔴 BLOCKING | check-intent-architecture.sh |
| Modal overlay | Automated | 🔴 BLOCKING | check-intent-architecture.sh |
| Layer dependencies | Automated | 🔴 BLOCKING | golangci-lint (depguard) |
| Forms architecture | Automated | 🔴 BLOCKING | golangci-lint (depguard) |
| Godoc completeness | Automated | 🟡 WARNING | check-intent-architecture.sh |
| Screen management | Automated | 🟡 WARNING | check-intent-architecture.sh |
| E2E tests | Automated | 🔴 BLOCKING | check-patterns-strict.sh |
| Coverage ≥95% | Automated | 🔴 BLOCKING | check-patterns-strict.sh |

---

**See Also**:
- [Intent Architecture Guide](../INTENT_ARCHITECTURE_GUIDE.md) - Full architecture documentation
- [Intent Patterns Library](../development/INTENT_PATTERNS_LIBRARY.md) - Pattern examples
- [AGENTS.md](../../AGENTS.md) - AI agent refusal rules

---

**Last Updated**: 2026-01-22  
**Status**: Production - Enforced by pre-commit hook
