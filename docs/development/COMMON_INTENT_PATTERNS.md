# Common Intent Patterns & Refactoring Guide

**Last Updated**: 2026-01-14
**Status**: Living Document - Updated as patterns emerge
**Purpose**: Identify and extract reusable patterns across all intents

---

## Table of Contents

1. [Overview](#overview)
2. [Extracted Patterns](#extracted-patterns)
3. [Candidate Patterns for Extraction](#candidate-patterns-for-extraction)
4. [Pattern Usage Matrix](#pattern-usage-matrix)
5. [Refactoring Guidelines](#refactoring-guidelines)
6. [Implementation Examples](#implementation-examples)

---

## Overview

As the KaRiya TUI architecture matures, we've identified common patterns that appear across multiple intents. Extracting these into reusable interfaces and utilities:

- ✅ **Reduces code duplication** (DRY principle)
- ✅ **Ensures consistency** across intents
- ✅ **Simplifies testing** (test once, use everywhere)
- ✅ **Eases maintenance** (fix once, benefit everywhere)
- ✅ **Accelerates development** (less code to write)

---

## Extracted Patterns

### 1. FilterBehavior Interface ✅

**Status**: Implemented
**File**: `internal/cli/intents/filter_behavior.go`
**Used By**: BrowseTimeline, ManageSkills

#### Interface Definition

```go
type FilterBehavior interface {
    // HasActiveFilters returns true if any non-default filters are active
    HasActiveFilters() bool
    
    // ClearFilters resets filters in FIFO order (most recent first)
    ClearFilters()
    
    // ApplyFilters applies current filter state to the data
    ApplyFilters()
    
    // RefreshData reloads/refreshes the filtered data
    RefreshData() tea.Cmd
}
```

#### When to Use

Implement FilterBehavior when your intent has:
- Search functionality (text-based filtering)
- Filter options (dropdown/checkbox filtering)
- Sort options (ordering data)
- Clear filters capability ('x' key to reset)

#### Benefits

- ✅ Consistent FIFO clearing order across all intents
- ✅ Standardized 'x' key behavior
- ✅ Clear contract for filter operations
- ✅ Easy to test (mock interface)

#### Helper Functions

```go
// SearchableText checks if text contains query (case-insensitive)
func SearchableText(text, query string) bool

// SearchableFields checks if ANY field matches query
func SearchableFields(query string, fields ...string) bool
```

#### Implementation Example

See: `internal/cli/intents/manage_skills_intent.go` lines 1336-1390

---

### 2. MessageInterceptor ✅

**Status**: Implemented
**File**: `internal/cli/intents/view_helpers.go`
**Used By**: ManageSkills (7 handlers), other intents

#### Pattern Definition

```go
type MessageInterceptor struct {
    backHandler GlobalKeyHandler
    quitHandler GlobalKeyHandler
    helpHandler GlobalKeyHandler
}

// Usage
interceptor := NewMessageInterceptor().
    OnQuit(StandardQuitHandler()).
    OnHelp(StandardHelpHandler(i.BaseIntent)).
    OnBack(func() tea.Cmd {
        // Custom back logic
        return nil
    }).
    InterceptOr(msg, fallbackHandler)
```

#### When to Use

Use MessageInterceptor in **every key handler** to ensure:
- Global keys (q, ?, Esc) work consistently
- Priority ordering: global keys → custom keys
- Standard behaviors for quit/help

#### Benefits

- ✅ Consistent global key handling
- ✅ Prevents key shadowing
- ✅ Reduces boilerplate (4 lines → 1 chainable call)
- ✅ Standard handlers available (StandardQuitHandler, StandardHelpHandler)

#### Standard Handlers

```go
// StandardQuitHandler - Returns quit message
func StandardQuitHandler() GlobalKeyHandler

// StandardHelpHandler - Shows help modal
func StandardHelpHandler(baseIntent *BaseIntent) GlobalKeyHandler
```

#### Implementation Example

```go
func (i *YourIntent) handleListKeys(msg tea.KeyMsg) tea.Cmd {
    return NewMessageInterceptor().
        OnQuit(StandardQuitHandler()).
        OnHelp(StandardHelpHandler(i.BaseIntent)).
        OnBack(func() tea.Cmd {
            return i.cancel()
        }).
        InterceptOr(msg, func() tea.Cmd {
            // Custom key handling
            switch msg.String() {
            case "enter":
                return i.selectItem()
            }
            return nil
        })
}
```

---

## Candidate Patterns for Extraction

### 3. Screen Result Handling 🔄

**Status**: Pattern Identified, Not Yet Extracted
**Used By**: GenerateCV, ManageSkills, BrowseTimeline

#### Current Pattern

Every screen-based intent has similar code:

```go
func (i *Intent) handleScreenResult(result screens.ScreenResult) tea.Cmd {
    switch r := result.(type) {
    case *screens.NavigateResult:
        return i.handleNavigateResult(r)
    case *screens.CancelResult:
        return i.handleCancelResult(r)
    case *screens.SubmitResult:
        return i.handleSubmitResult(r)
    case *screens.ErrorResult:
        return i.handleErrorResult(r)
    default:
        return nil
    }
}
```

#### Proposed Interface

```go
// ScreenResultHandler defines how an intent handles screen results
type ScreenResultHandler interface {
    HandleNavigate(result *screens.NavigateResult) tea.Cmd
    HandleCancel(result *screens.CancelResult) tea.Cmd
    HandleSubmit(result *screens.SubmitResult) tea.Cmd
    HandleError(result *screens.ErrorResult) tea.Cmd
}

// DefaultScreenResultDispatcher provides default routing
type DefaultScreenResultDispatcher struct {
    handler ScreenResultHandler
}

func (d *DefaultScreenResultDispatcher) Dispatch(result screens.ScreenResult) tea.Cmd {
    // Type switch with automatic routing
}
```

#### Benefits

- ✅ Eliminates repetitive type switching
- ✅ Forces implementation of all result handlers
- ✅ Compile-time safety
- ✅ Easy to test (mock interface)

#### Estimated Impact

- **Lines saved per intent**: ~30 lines
- **Intents affected**: 3 (GenerateCV, ManageSkills, BrowseTimeline)
- **Total savings**: ~90 lines

---

### 4. Modal Lifecycle Management 🔄

**Status**: Pattern Identified, Not Yet Extracted
**Used By**: ManageSkills (3 modals), BrowseTimeline (5 modals)

#### Current Pattern

Every modal has similar update logic:

```go
func (i *Intent) handleXModalUpdate(msg tea.Msg) tea.Cmd {
    var cmd tea.Cmd
    i.xModal, cmd = i.xModal.Update(msg)
    
    if i.xModal.WasAccepted() {
        data := i.xModal.GetData()
        i.xModal = nil
        return i.applyModalData(data)
    }
    
    if i.xModal.WasCancelled() {
        i.xModal = nil
        return nil
    }
    
    return cmd
}
```

#### Proposed Interface

```go
// ModalHandler defines the lifecycle of a modal
type ModalHandler interface {
    Update(msg tea.Msg) tea.Cmd
    WasAccepted() bool
    WasCancelled() bool
    GetData() interface{}
}

// ModalLifecycleManager handles common modal update logic
type ModalLifecycleManager struct {
    onAccept func(data interface{}) tea.Cmd
    onCancel func() tea.Cmd
}

func (m *ModalLifecycleManager) HandleUpdate(
    modal ModalHandler, 
    msg tea.Msg,
) (ModalHandler, tea.Cmd)
```

#### Benefits

- ✅ Reduces modal boilerplate (15 lines → 3 lines)
- ✅ Consistent modal lifecycle
- ✅ Handles cleanup automatically
- ✅ Easy to test

#### Estimated Impact

- **Lines saved per modal**: ~12 lines
- **Modals affected**: 8 (3 ManageSkills + 5 BrowseTimeline)
- **Total savings**: ~96 lines

---

### 5. State Transition Helpers 🔄

**Status**: Pattern Identified, Not Yet Extracted
**Used By**: ManageSkills, potentially others

#### Current Pattern

```go
func (i *ManageSkillsIntent) transitionToListScreen() tea.Cmd {
    i.currentState = SkillsStateList
    i.selectedSkill = nil
    i.eventDetail = nil
    // ... more cleanup
    return i.reloadSkills()
}

func (i *ManageSkillsIntent) transitionToDetailScreen() tea.Cmd {
    i.currentState = SkillsStateDetail
    return i.loadSkillEvents(i.selectedSkill.ID)
}
```

#### Proposed Pattern

```go
// StateTransition encapsulates a state change with cleanup
type StateTransition struct {
    targetState interface{}
    cleanup     []func()
    initCmd     tea.Cmd
}

func (i *Intent) Transition(transition StateTransition) tea.Cmd {
    // Run cleanup functions
    for _, fn := range transition.cleanup {
        fn()
    }
    
    // Set new state
    i.currentState = transition.targetState
    
    // Return initialization command
    return transition.initCmd
}
```

#### Benefits

- ✅ Declarative state transitions
- ✅ Consistent cleanup patterns
- ✅ Easier to track state changes
- ✅ Testable state machine

#### Estimated Impact

- **Lines saved per intent**: ~20-40 lines
- **Intents affected**: All intents with multiple states
- **Total savings**: ~200+ lines across project

---

### 6. Error Handling Pattern 🔄

**Status**: Pattern Identified, Not Yet Extracted
**Used By**: All intents

#### Current Pattern

```go
func (i *Intent) handleErrorInternal(err error) tea.Cmd {
    i.state = StateError
    i.error = err
    return nil
}

func (i *Intent) handleErrorResult(result *screens.ErrorResult) tea.Cmd {
    i.state = StateError
    i.error = result.Err
    return nil
}
```

#### Proposed Interface

```go
// ErrorHandler defines how an intent handles errors
type ErrorHandler interface {
    SetError(err error)
    GetError() error
    ClearError()
    IsErrorState() bool
}

// StandardErrorBehavior provides default error handling
type StandardErrorBehavior struct {
    currentError error
    setState     func(state interface{})
}

func (s *StandardErrorBehavior) HandleError(err error) tea.Cmd {
    s.currentError = err
    s.setState(StateError)
    return nil
}
```

#### Benefits

- ✅ Consistent error handling
- ✅ Error state management
- ✅ Recovery patterns
- ✅ Error logging/tracking

#### Estimated Impact

- **Lines saved per intent**: ~10-15 lines
- **Intents affected**: All (11 intents)
- **Total savings**: ~130 lines

---

## Pattern Usage Matrix

| Pattern | Browse Timeline | Manage Skills | Generate CV | Capture Event | Others |
|---------|-----------------|---------------|-------------|---------------|--------|
| **FilterBehavior** | ✅ Implemented | ✅ Implemented | ❌ N/A | ❌ N/A | ❌ N/A |
| **MessageInterceptor** | ⚠️ Partial | ✅ Full (7 handlers) | ⚠️ Partial | ⚠️ Partial | ⚠️ Partial |
| **Screen Result** | ✅ Has pattern | ✅ Has pattern | ✅ Has pattern | ❌ N/A | ⚠️ Some |
| **Modal Lifecycle** | ✅ 5 modals | ✅ 3 modals | ❌ N/A | ⚠️ 3 modals | ⚠️ Some |
| **State Transition** | ⚠️ Simple | ✅ Complex (4 helpers) | ⚠️ Simple | ⚠️ Medium | ⚠️ Varies |
| **Error Handling** | ✅ Has pattern | ✅ Has pattern | ✅ Has pattern | ✅ Has pattern | ✅ All |

**Legend**:
- ✅ Pattern implemented/applicable
- ⚠️ Pattern partially implemented or applicable
- ❌ Pattern not applicable or not implemented

---

## Refactoring Guidelines

### When to Extract a Pattern

Extract a pattern when:
1. ✅ Used in **3+ intents** (Rule of Three)
2. ✅ Logic is **identical or nearly identical** (>80% similar)
3. ✅ Pattern has **clear boundaries** (well-defined inputs/outputs)
4. ✅ **Tests exist** for the pattern (extract tested code)
5. ✅ Extraction **reduces complexity** (doesn't add abstraction overhead)

### When NOT to Extract

Avoid extraction when:
- ❌ Used in only 1-2 places (premature abstraction)
- ❌ Logic differs significantly between uses (false similarity)
- ❌ Extraction adds more complexity than it removes
- ❌ Pattern is still evolving (wait for stability)

### Extraction Process

1. **Identify Pattern**: Document usage across intents
2. **Design Interface**: Define clear contract
3. **Write Tests**: Test interface in isolation
4. **Implement**: Create interface and default implementation
5. **Migrate One**: Update one intent to use new pattern
6. **Verify**: Run all tests, ensure no regressions
7. **Migrate All**: Update remaining intents
8. **Clean Up**: Remove old code
9. **Document**: Update this guide and workflow docs

### Refactoring Priority

**High Priority** (Extract Next):
1. Screen Result Handling (3 intents, clear pattern)
2. Modal Lifecycle Management (8 modals, high duplication)

**Medium Priority**:
3. State Transition Helpers (complex, but varies by intent)
4. Error Handling Pattern (all intents, but simple)

**Low Priority**:
5. Future patterns as they emerge

---

## Implementation Examples

### Example 1: Implementing FilterBehavior

```go
// 1. Add interface compliance check
var _ FilterBehavior = (*YourIntent)(nil)

// 2. Implement HasActiveFilters
func (i *YourIntent) HasActiveFilters() bool {
    return i.filters != nil && (
        i.filters.SearchText != "" ||
        i.filters.Category != "" ||
        i.filters.SortBy != ""
    )
}

// 3. Implement ClearFilters (FIFO order)
func (i *YourIntent) ClearFilters() {
    if i.filters == nil {
        return
    }
    
    // Clear in FIFO order: search → filter → sort
    if i.filters.SearchText != "" {
        i.filters.SearchText = ""
        return
    }
    
    if i.filters.Category != "" {
        i.filters.Category = ""
        return
    }
    
    i.filters.SortBy = ""
}

// 4. Implement ApplyFilters
func (i *YourIntent) ApplyFilters() {
    if i.filters != nil && i.filters.SearchText != "" {
        i.data = i.applySearchFilter(i.data, i.filters.SearchText)
    }
}

// 5. Implement RefreshData
func (i *YourIntent) RefreshData() tea.Cmd {
    return i.reloadData()
}

// 6. Add 'x' key handler
func (i *YourIntent) handleListKeys(msg tea.KeyMsg) tea.Cmd {
    switch msg.String() {
    case "x":
        if i.HasActiveFilters() {
            i.ClearFilters()
            return i.RefreshData()
        }
    }
    return nil
}
```

### Example 2: Using MessageInterceptor

```go
func (i *YourIntent) handleDetailKeys(msg tea.KeyMsg) tea.Cmd {
    // BEFORE: Manual global key handling (verbose)
    switch msg.String() {
    case "q", "ctrl+c":
        return i.quit()
    case "?", "h":
        return i.showHelp()
    case "esc":
        return i.goBack()
    case "e":
        return i.edit()
    case "d":
        return i.delete()
    }
    return nil
    
    // AFTER: MessageInterceptor (concise, consistent)
    return NewMessageInterceptor().
        OnQuit(StandardQuitHandler()).
        OnHelp(StandardHelpHandler(i.BaseIntent)).
        OnBack(func() tea.Cmd { return i.goBack() }).
        InterceptOr(msg, func() tea.Cmd {
            switch msg.String() {
            case "e":
                return i.edit()
            case "d":
                return i.delete()
            }
            return nil
        })
}
```

---

## Future Patterns

As new patterns emerge, add them here:

### Breadcrumb Generation Pattern (Potential)

All intents generate breadcrumbs from state. Could extract:

```go
type BreadcrumbProvider interface {
    GetBreadcrumbs() []string
}
```

### Help Footer Generation Pattern (Potential)

Most intents use `getContextHelp()`. Could extract:

```go
type ContextualHelp interface {
    GetHelpForState(state interface{}) string
}
```

---

## Metrics & Impact

### Current Extractions

| Pattern | Intents Using | Lines Saved | Consistency Gain |
|---------|---------------|-------------|------------------|
| FilterBehavior | 2 | ~50 per intent | High |
| MessageInterceptor | ~5 (partial) | ~10 per handler | High |
| **Total** | - | **~150 lines** | - |

### Projected Savings (After All Extractions)

| Pattern | Intents Affected | Lines per Intent | Total Savings |
|---------|------------------|------------------|---------------|
| Screen Result Handling | 3 | 30 | 90 |
| Modal Lifecycle | 2 | 96 | 192 |
| State Transitions | 8 | 25 | 200 |
| Error Handling | 11 | 12 | 132 |
| **Total Projected** | - | - | **~614 lines** |

---

## References

- **FilterBehavior Interface**: `internal/cli/intents/filter_behavior.go`
- **MessageInterceptor**: `internal/cli/intents/view_helpers.go` lines 446-632
- **Intent Patterns Library**: `docs/development/INTENT_PATTERNS_LIBRARY.md`
- **BrowseTimeline (Reference)**: `internal/cli/intents/browse_timeline_intent.go`
- **ManageSkills (Example)**: `internal/cli/intents/manage_skills_intent.go`

---

**Last Updated**: 2026-01-14
**Next Review**: After next intent refactoring
**Owner**: KaRiya TUI Architecture Team
