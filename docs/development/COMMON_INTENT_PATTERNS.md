# Common Intent Patterns & Refactoring Guide

**Last Updated**: 2026-01-20
**Status**: Living Document - Updated as patterns emerge
**Purpose**: Identify and extract reusable patterns across all intents

> **UIKit Migration Note**: New implementations should use UIKit components.
> See [UIKIT_GUIDE.md](../UIKIT_GUIDE.md) for the standardized component library.

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

### 3. Screen Result Handling ✅

**Status**: Implemented
**File**: `internal/cli/intents/screen_result_behavior.go`
**Used By**: GenerateCV, ManageSkills, BrowseTimeline

#### Interface Definition

```go
// ScreenResultHandler defines how an intent handles screen results
type ScreenResultHandler interface {
	HandleNavigate(result *screens.NavigateResult) tea.Cmd
	HandleCancel(result *screens.CancelResult) tea.Cmd
	HandleSubmit(result *screens.SubmitResult) tea.Cmd
	HandleError(result *screens.ErrorResult) tea.Cmd
}
```

#### When to Use

Implement ScreenResultHandler when your intent uses screen-based architecture and needs to:
- Handle navigation results (user selected something)
- Handle cancel results (user pressed Escape)
- Handle submit results (user submitted a form)
- Handle error results (screen encountered an error)

#### Benefits

- ✅ Eliminates repetitive type switching (20+ lines → 1 line)
- ✅ Forces implementation of all result handlers (compile-time safety)
- ✅ Clear contract for screen result handling
- ✅ Easy to test (mock interface)

#### Dispatcher

```go
// ScreenResultDispatcher routes screen results to appropriate handler methods
type ScreenResultDispatcher struct {
	handler ScreenResultHandler
}

func (d *ScreenResultDispatcher) Dispatch(result screens.ScreenResult) tea.Cmd {
	// Automatically routes to correct handler based on result type
}
```

#### Implementation Example

```go
// 1. Add interface compliance check
var _ ScreenResultHandler = (*YourIntent)(nil)

// 2. Implement all handler methods
func (i *YourIntent) HandleNavigate(result *screens.NavigateResult) tea.Cmd { ... }
func (i *YourIntent) HandleCancel(result *screens.CancelResult) tea.Cmd { ... }
func (i *YourIntent) HandleSubmit(result *screens.SubmitResult) tea.Cmd { ... }
func (i *YourIntent) HandleError(result *screens.ErrorResult) tea.Cmd { ... }

// 3. Replace handleScreenResult with dispatcher
func (i *YourIntent) handleScreenResult(result screens.ScreenResult) tea.Cmd {
	return NewScreenResultDispatcher(i).Dispatch(result)
}
```

#### Actual Impact

- **Lines saved per intent**: ~30 lines
- **Intents migrated**: 3 (GenerateCV, ManageSkills, BrowseTimeline)
- **Total savings**: ~90 lines
- **Tests added**: 17 comprehensive tests

---

---

## Documented Patterns (Convention-Based)

These patterns have been analyzed and determined to be better suited for documentation and convention rather than code extraction. The patterns are already consistent across the codebase and extraction would add more complexity than value.

### 4. State Transition Helpers 📖

**Status**: Documented as Best Practice
**Documentation**: `docs/development/STATE_TRANSITION_PATTERNS.md`
**Used By**: ManageSkills (4 helpers), others use inline transitions

#### Pattern Description

State transition helpers encapsulate the logic for moving between intent states, including:
- Setting the new state
- Initializing resources for the new state
- Cleaning up old state resources
- Returning initialization commands

#### When to Use

**Extract helpers when**:
- Intent has **4+ states**
- Transitions require initialization or cleanup
- Same transition is used from multiple places
- Transition logic is **>15 lines**

**Use inline when**:
- Intent has **≤3 states**
- Transitions are simple (**<10 lines**)
- Each transition is unique

#### Reference Implementation

See `ManageSkillsIntent` (`manage_skills_intent.go` lines 2189-2264):
```go
func (i *ManageSkillsIntent) transitionToListScreen() tea.Cmd
func (i *ManageSkillsIntent) transitionToDetailScreen() tea.Cmd
func (i *ManageSkillsIntent) transitionToFormScreen(skill *domain.Skill) tea.Cmd
func (i *ManageSkillsIntent) transitionToDeleteScreen(skill *domain.Skill) tea.Cmd
```

#### Why Not Extracted to Code

✅ Already consistent across intents that need it  
✅ Highly specific to each intent's state machine  
✅ Extraction would add abstraction overhead  
✅ Value is in convention and documentation, not reusable code

**Action**: Follow documented pattern when creating complex intents

---

### 5. Modal Lifecycle Management 📖

**Status**: Documented in Existing Guides
**Documentation**: `docs/MODAL_PATTERNS.md`
**Used By**: ManageSkills (3 modals), BrowseTimeline (5 modals)

#### Pattern Description

Modal lifecycle follows a consistent pattern:
1. Check if modal is visible
2. Update modal with message
3. Check if modal closed
4. If completed, process data
5. Clean up modal reference

#### Current Pattern (Already Consistent)

```go
if i.xModal != nil && i.xModal.IsVisible() {
    cmd, completed, data := i.xModal.Update(msg)
    if !i.xModal.IsVisible() {
        // Modal closed
        if completed && data != nil {
            // Process data
            i.processModalData(data)
        }
        // Cleanup
        i.xModal = nil
    }
    return cmd
}
```

#### Reference Implementations

**ManageSkills**: 3 modals (filter, sort, search) - lines 577-755  
**BrowseTimeline**: 5 modals (quickAdd, edit, delete, filter, viewDetail) - lines 184-342

#### Why Not Extracted to Code

✅ Already documented in `MODAL_PATTERNS.md`  
✅ Pattern is consistent across all 8 modals  
✅ Only ~8-12 common lines per modal (visibility check + cleanup)  
✅ Data processing logic varies significantly per modal  
✅ Extraction would obscure rather than clarify

**Action**: Follow documented pattern in `MODAL_PATTERNS.md`

---

### 6. Error Handling Pattern 📖

**Status**: Simple Pattern, No Extraction Needed
**Used By**: All intents

---

#### Pattern Description

Error handling in intents typically involves:
1. Creating an `IntentError` with code, message, and cause
2. Setting the error on the intent result or state
3. Optionally transitioning to error state or recovery state

#### Current Patterns (Already Simple)

**Pattern A: Simple setFailed helper** (1 line)
```go
func (i *Intent) setFailed(code, message string, cause error) {
    i.result = NewFailedResult[*IntentResult](code, message, cause)
    i.active = false
}

// Usage
i.setFailed("VALIDATION_ERROR", "Invalid input", err)
```

**Pattern B: Direct error assignment**
```go
i.state.error = &IntentError{
    Code:    "FORM_ERROR",
    Message: err.Error(),
    Cause:   err,
}
```

**Pattern C: Error with recovery** (ManageSkills)
```go
func (i *ManageSkillsIntent) handleErrorInternal(err error) tea.Cmd {
    i.result = &IntentResult[*ManageSkillsResult]{
        Status: Failed,
        Error: &IntentError{
            Code:    "SCREEN_ERROR",
            Message: err.Error(),
            Cause:   err,
        },
    }
    return i.transitionToListScreen()  // Recovery
}
```

#### Why Not Extracted to Code

✅ Already concise (1-10 lines per intent)  
✅ `NewFailedResult` helper already exists  
✅ Recovery strategy varies significantly per intent  
✅ Pattern is straightforward and well-understood

**Action**: Use `NewFailedResult` helper or direct assignment as appropriate

---

## Candidate Patterns for Future Extraction

No patterns currently identified for extraction. New patterns should be added here as they emerge and meet the Rule of Three (used in 3+ places).

---

## Pattern Usage Matrix

| Pattern | Browse Timeline | Manage Skills | Generate CV | Capture Event | Others |
|---------|-----------------|---------------|-------------|---------------|--------|
| **FilterBehavior** | ✅ Implemented | ✅ Implemented | ❌ N/A | ❌ N/A | ❌ N/A |
| **MessageInterceptor** | ⚠️ Partial | ✅ Full (7 handlers) | ⚠️ Partial | ⚠️ Partial | ⚠️ Partial |
| **Screen Result** | ✅ Implemented | ✅ Implemented | ✅ Implemented | ❌ N/A | ⚠️ Some |
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
| **Screen Result Handling** | **3** | **~30 per intent** | **High** |
| **Total** | - | **~240 lines** | - |

### Analysis Results: Documented vs Extracted

After thorough analysis of remaining candidate patterns, we determined that:

| Pattern | Initial Estimate | Actual Analysis | Decision |
|---------|------------------|-----------------|----------|
| Modal Lifecycle | 192 lines (8 modals) | ~8-12 lines common code per modal | ✅ Document (already consistent) |
| State Transitions | 200 lines (8 intents) | Highly intent-specific | ✅ Document (best practice guide) |
| Error Handling | 132 lines (11 intents) | Already concise (1-10 lines) | ✅ Document (use existing helpers) |

**Key Insight**: Not all duplication requires code extraction. Sometimes **consistency through documentation and convention** provides more value than abstraction.

### Final Extraction Summary

| Category | Patterns | Lines Saved | Approach |
|----------|----------|-------------|----------|
| **Extracted to Code** | 3 (FilterBehavior, MessageInterceptor, ScreenResultHandler) | ~240 lines | Reusable interfaces + helpers |
| **Documented as Convention** | 3 (State Transitions, Modal Lifecycle, Error Handling) | N/A | Best practice guides + references |
| **Total Patterns Addressed** | **6** | **~240 lines** | **Mixed approach (optimal)** |

---

## References

### Extracted Patterns (Code)

- **FilterBehavior Interface**: `internal/cli/intents/filter_behavior.go`
- **MessageInterceptor**: `internal/cli/intents/view_helpers.go` lines 446-632
- **ScreenResultHandler Interface**: `internal/cli/intents/screen_result_behavior.go`

### Documented Patterns (Convention)

- **State Transition Helpers**: `docs/development/STATE_TRANSITION_PATTERNS.md`
- **Modal Lifecycle**: `docs/MODAL_PATTERNS.md`
- **Error Handling**: Use `NewFailedResult` helper (see intents for examples)

### Reference Implementations

- **BrowseTimeline**: `internal/cli/intents/browse_timeline_intent.go` (all 12 patterns)
- **ManageSkills**: `internal/cli/intents/manage_skills_intent.go` (10/12 patterns + transition helpers)
- **GenerateCV**: `internal/cli/intents/generate_cv_intent.go` (screen result handling)

### Related Documentation

- **Intent Patterns Library**: `docs/development/INTENT_PATTERNS_LIBRARY.md`
- **Pattern Extraction Summary**: `docs/development/PATTERN_EXTRACTION_SUMMARY.md`
- **TUI Developer Guide**: `docs/TUI_DEVELOPER_GUIDE.md`

---

**Last Updated**: 2026-01-14  
**Next Review**: After new pattern emergence (3+ uses)  
**Owner**: KaRiya TUI Architecture Team
