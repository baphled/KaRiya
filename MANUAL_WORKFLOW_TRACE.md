# Manual Workflow Trace - Complete Issue Analysis

**Date**: January 3, 2026
**Status**: INCOMPLETE - Multiple Critical Issues Found
**Severity**: CRITICAL - Production deployment blocked

---

## Executive Summary

Manual tracing through the application workflow reveals **multiple critical issues** that will cause panics and crashes in production. The aggressive app.go refactoring is **NOT production-ready** and requires significant work to complete.

**Critical Issues Found**: 7
**High Priority Issues**: 3
**Medium Priority Issues**: 2

---

## Workflow Trace 1: Application Startup

### Step 1a: NewModel() Called
```go
// internal/cli/app/app.go line 56
func NewModel(cliService *service.CLIEventService, careerService *careerservice.Service) *Model {
    // ...
    router := intents.NewDefaultIntentRouter()
    registerAllIntents(router, careerService, log, ctx)  // Line 70
    // ...
}
```

**Status**: ✅ OK - No issues here

---

### Step 1b: registerAllIntents() Called - CaptureEvent

```go
// internal/cli/app/app.go lines 220-232
router.RegisterIntent("capture_event", func() intents.Intent {
    captureCtx := &intents.CaptureEventContext{
        CaptureStrategy: "manual",
        Metadata:        make(map[string]string),
    }
    intent, err := intents.NewCaptureEventIntent(captureCtx)
    if err != nil {
        log.Error("Failed to create CaptureEvent intent: %v", err)
        return nil  // ⚠️ ISSUE #1: Returns nil on error
    }
    return intent
})
```

**Issue #1: Nil Return on Error - CaptureEvent**
- **Problem**: If intent creation fails, factory returns `nil`
- **When Called**: Router calls `intent.Init()` on nil pointer → **PANIC**
- **Impact**: CRITICAL - Application crashes if CaptureEvent intent fails to initialize
- **Root Cause**: No error handling for nil returns in router

---

### Step 1c: registerAllIntents() Called - BrowseTimeline

```go
// internal/cli/app/app.go lines 234-248
router.RegisterIntent("browse_timeline", func() intents.Intent {
    events, err := careerService.GetEventRepository().List(ctx, careerrepo.ListFilters{Limit: 1000})
    if err != nil {
        log.Error("Failed to load events: %v", err)
        events = make([]*career.CareerEvent, 0)
    }
    browserCtx := &intents.BrowseTimelineContext{
        Events: events,
    }
    intent, err := intents.NewBrowseTimelineIntent(browserCtx)
    if err != nil {
        log.Error("Failed to create BrowseTimeline intent: %v", err)
        return nil  // ⚠️ ISSUE #2: Returns nil on error
    }
    return intent
})
```

**Issue #2: Nil Return on Error - BrowseTimeline**
- **Problem**: If intent creation fails, factory returns `nil`
- **When Called**: Router calls `intent.Init()` on nil pointer → **PANIC**
- **Impact**: CRITICAL - Application crashes if BrowseTimeline intent fails
- **Root Cause**: Same as Issue #1

---

### Step 1d: registerAllIntents() Called - BurstManagement (NEW INTENT)

```go
// internal/cli/app/app.go lines 295-307
router.RegisterIntent("burst_management", func() intents.Intent {
    bursts, err := careerService.GetBurstRepository().List(ctx, careerrepo.BurstListFilters{Limit: 1000})
    if err != nil {
        log.Error("Failed to load bursts: %v", err)
        bursts = make([]*career.Burst, 0)
    }
    burstCtx := &intents.BurstManagementContext{
        Bursts:  bursts,
        Service: careerService,
    }
    return intents.NewBurstManagementIntent(burstCtx)  // ⚠️ ISSUE #3: No error handling
})
```

**Issue #3: Missing Context Initialization - BurstManagement**
- **Problem**: `BurstManagementContext` created but missing required fields
- **Missing Fields**:
  - `CurrentState` - not initialized
  - `Context` - not initialized (needed for database operations)
  - `BurstRepository` - not initialized
  - `FormErrors` - not initialized
  - `ExpandedRows` - not initialized
- **When Called**: Intent.Init() tries to access these fields → **PANIC or nil dereference**
- **Impact**: CRITICAL - Intent will crash when trying to load bursts

---

### Step 1e: registerAllIntents() Called - FactManagement (NEW INTENT)

```go
// internal/cli/app/app.go lines 311-323
router.RegisterIntent("fact_management", func() intents.Intent {
    facts, err := careerService.GetFactRepository().List(ctx, careerrepo.FactListFilters{Limit: 1000})
    if err != nil {
        log.Error("Failed to load facts: %v", err)
        facts = make([]*career.Fact, 0)
    }
    factCtx := &intents.FactManagementContext{
        Facts: facts,  // ⚠️ ISSUE #4: Missing required fields
    }
    return intents.NewFactManagementIntent(factCtx)
})
```

**Issue #4: Missing Context Initialization - FactManagement**
- **Problem**: `FactManagementContext` created but missing required fields
- **Missing Fields**:
  - `CurrentState` - not initialized
  - `Context` - not initialized
  - `Service` - not initialized
  - `FormErrors` - not initialized
  - `MinQuality`, `MaxQuality` - not initialized
- **When Called**: Intent.Init() tries to access these fields → **PANIC or nil dereference**
- **Impact**: CRITICAL - Intent will crash when trying to load facts

---

### Step 1f: registerAllIntents() Called - ImportWizard (NEW INTENT)

```go
// internal/cli/app/app.go lines 325-329
router.RegisterIntent("import_wizard", func() intents.Intent {
    importCtx := &intents.ImportWizardContext{}  // ⚠️ ISSUE #5: Empty context
    return intents.NewImportWizardIntent(importCtx)
})
```

**Issue #5: Empty Context Initialization - ImportWizard**
- **Problem**: `ImportWizardContext` created completely empty
- **Missing Fields**:
  - `CurrentState` - not initialized
  - `Context` - not initialized
  - All other fields empty
- **When Called**: Intent.Init() tries to use empty context → **PANIC**
- **Impact**: CRITICAL - Intent will crash immediately

---

### Step 1g: registerAllIntents() Called - MetadataEditor (NEW INTENT)

```go
// internal/cli/app/app.go lines 331-335
router.RegisterIntent("metadata_editor", func() intents.Intent {
    metaCtx := &intents.MetadataEditorContext{}  // ⚠️ ISSUE #6: Empty context
    return intents.NewMetadataEditorIntent(metaCtx)
})
```

**Issue #6: Empty Context Initialization - MetadataEditor**
- **Problem**: `MetadataEditorContext` created completely empty
- **Missing Fields**:
  - `CurrentState` - not initialized
  - `Context` - not initialized
  - All other fields empty
- **When Called**: Intent.Init() tries to use empty context → **PANIC**
- **Impact**: CRITICAL - Intent will crash immediately

---

### Step 1h: registerAllIntents() Called - BulkOperations (NEW INTENT)

```go
// internal/cli/app/app.go lines 337-341
router.RegisterIntent("bulk_operations", func() intents.Intent {
    bulkCtx := &intents.BulkOperationsContext{}  // ⚠️ ISSUE #7: Empty context
    return intents.NewBulkOperationsIntent(bulkCtx)
})
```

**Issue #7: Empty Context Initialization - BulkOperations**
- **Problem**: `BulkOperationsContext` created completely empty
- **Missing Fields**:
  - `CurrentState` - not initialized
  - `Context` - not initialized
  - All other fields empty
- **When Called**: Intent.Init() tries to use empty context → **PANIC**
- **Impact**: CRITICAL - Intent will crash immediately

---

## Workflow Trace 2: User Selects Menu Item

### Step 2a: User Presses Down Arrow

```go
// internal/cli/app/app.go line 159
func (m *Model) handleMenuInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "down", "j":
        if m.selectedMenuIndex < len(m.menuItems)-1 {
            m.selectedMenuIndex++
        }
    // ...
    }
    return m, nil
}
```

**Status**: ✅ OK - Menu navigation works

---

### Step 2b: User Presses Enter to Select "Manage Bursts"

```go
// internal/cli/app/app.go lines 163-168
case "enter", " ":
    selectedItem := m.menuItems[m.selectedMenuIndex]  // "burst_management"
    m.state = StateIntent
    return m, m.activateIntent(selectedItem.Intent)
```

**Status**: ✅ OK - Intent activation triggered

---

### Step 2c: activateIntent() Called

```go
// internal/cli/app/app.go lines 191-196
func (m *Model) activateIntent(intentName string) tea.Cmd {
    return func() tea.Msg {
        _, _ = m.intentRouter.ActivateIntent(intentName, make(map[string]interface{}))
        return nil
    }
}
```

**Status**: ✅ OK - Router activation called

---

### Step 2d: Router.ActivateIntent() Called

```go
// internal/cli/intents/router.go lines 75-93
func (r *DefaultIntentRouter) ActivateIntent(name string, context map[string]interface{}) (tea.Cmd, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    factory, exists := r.intents[name]
    if !exists {
        return nil, fmt.Errorf("intent %q not found", name)
    }

    // Push the current intent to history (if any).
    if r.activeIntent != nil {
        r.intentHistory = append(r.intentHistory, r.activeIntent)
    }

    // Create and activate the new intent.
    intent := factory()  // ⚠️ Calls the factory from registerAllIntents()
    r.activeIntent = intent

    // Call the intent's Init method to get any startup commands.
    return intent.Init(), nil  // ⚠️ PANIC if intent is nil!
}
```

**CRITICAL ISSUE**: If `factory()` returns nil (due to error), then `intent.Init()` will panic with:
```
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=...]
```

---

### Step 2e: BurstManagement Intent Factory Executes

Looking back at the factory code (lines 295-307), when the factory is called:

```go
bursts, err := careerService.GetBurstRepository().List(ctx, careerrepo.BurstListFilters{Limit: 1000})
if err != nil {
    log.Error("Failed to load bursts: %v", err)
    bursts = make([]*career.Burst, 0)  // ✅ Handles error
}
burstCtx := &intents.BurstManagementContext{
    Bursts:  bursts,
    Service: careerService,
    // ❌ Missing: CurrentState, Context, BurstRepository, FormErrors, ExpandedRows, etc.
}
return intents.NewBurstManagementIntent(burstCtx)
```

**Issue**: Context is incomplete. When `NewBurstManagementIntent()` returns and `Init()` is called:

```go
// internal/cli/intents/burst_management_intent.go
func (m *BurstManagementModel) Init() tea.Cmd {
    m.data.CurrentState = BurstListState  // ❌ m.data exists but is incomplete

    if err := m.data.LoadBursts(); err != nil {  // ❌ LoadBursts() might fail
        m.result = &IntentResult[*BurstManagementResult]{
            Status: Failed,
            // ...
        }
        return nil
    }
    // ...
}
```

**Potential Panics**:
1. `m.data` is not nil, but fields are uninitialized
2. `LoadBursts()` might panic if it tries to access `m.data.BurstRepository` (which is nil)
3. `LoadBursts()` might panic if it tries to access `m.data.Context` (which is nil)

---

## Summary of Critical Issues

| # | Issue | Location | Severity | Impact |
|---|-------|----------|----------|--------|
| 1 | Nil return on CaptureEvent error | app.go:230 | CRITICAL | Panic in router.Init() |
| 2 | Nil return on BrowseTimeline error | app.go:245 | CRITICAL | Panic in router.Init() |
| 3 | Incomplete BurstManagement context | app.go:298-303 | CRITICAL | Panic in Init() or LoadBursts() |
| 4 | Incomplete FactManagement context | app.go:314-320 | CRITICAL | Panic in Init() or LoadFacts() |
| 5 | Empty ImportWizard context | app.go:327 | CRITICAL | Panic in Init() |
| 6 | Empty MetadataEditor context | app.go:333 | CRITICAL | Panic in Init() |
| 7 | Empty BulkOperations context | app.go:339 | CRITICAL | Panic in Init() |

---

## Required Fixes

### Fix 1: Handle nil returns in router

```go
// internal/cli/intents/router.go - ActivateIntent()
intent := factory()
if intent == nil {
    return nil, fmt.Errorf("failed to create intent %q", name)
}
r.activeIntent = intent
return intent.Init(), nil
```

### Fix 2: Initialize all context fields properly

```go
// internal/cli/app/app.go - registerAllIntents()
burstCtx := &intents.BurstManagementContext{
    CurrentState:     intents.BurstListState,
    Bursts:           bursts,
    Service:          careerService,
    BurstRepository:  careerService.GetBurstRepository(),
    Context:          ctx,
    FormErrors:       make(map[string]string),
    ExpandedRows:     make(map[int]bool),
    // ... initialize all other fields
}
```

### Fix 3: Ensure all new intent constructors handle nil context

```go
// internal/cli/intents/burst_management_intent.go
func NewBurstManagementIntent(data *BurstManagementContext) (*BurstManagementModel, error) {
    if data == nil {
        return nil, fmt.Errorf("context is required")
    }
    if data.Context == nil {
        return nil, fmt.Errorf("context.Context is required")
    }
    if data.Service == nil {
        return nil, fmt.Errorf("context.Service is required")
    }
    return &BurstManagementModel{
        data: data,
        // ...
    }, nil
}
```

---

## Completion Status

**Work Completed**: ~30%
**Work Remaining**: ~70%
**Estimated Time to Complete**: 4-6 hours

### What's Done:
- ✅ app.go refactored to 386 lines
- ✅ Menu system implemented
- ✅ Intent registration structure in place
- ✅ Tests for intent framework passing

### What's NOT Done:
- ❌ Context initialization for all new intents
- ❌ Nil pointer protection in router
- ❌ Error handling for intent creation failures
- ❌ Integration tests for menu workflows
- ❌ End-to-end testing of all 10 intents

---

## Conclusion

**The application is NOT production-ready.** Multiple critical issues will cause panics when users interact with the application. The refactoring has created the foundation, but the implementation is incomplete.

**Recommendation**: Complete all fixes before any production deployment.


