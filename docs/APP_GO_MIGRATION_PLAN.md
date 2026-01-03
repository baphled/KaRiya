# app.go Migration Plan: Screen-Based to Intent-Based Architecture

**Document Version**: 1.0
**Date Created**: 2026-01-03
**Status**: PLANNING PHASE
**Target Completion**: 4-6 weeks (phased approach)
**Complexity**: HIGH (906-line Update method, 37+ model fields, 31 screens)

---

## Executive Summary

This document outlines a comprehensive, phased approach to migrate the KaRiya TUI application from a **screen-based navigation architecture** to the new **intent-driven architecture** defined in tasks-09-tui-intent-refactoring.md.

### Current State
- **Architecture**: Screen-based with 31 screens
- **Update() Method**: 906 lines of complex logic
- **Model Fields**: 37+ model instances (one per screen)
- **Code Quality**: Monolithic, difficult to maintain
- **Status**: Functional but needs refactoring

### Target State
- **Architecture**: Intent-driven with 5 core intents
- **Update() Method**: <100 lines delegating to router
- **Model Fields**: Minimal, focused on app-level state
- **Code Quality**: Modular, maintainable, testable
- **Status**: Production-ready with clear separation of concerns

### Key Benefits
✅ **Reduced Complexity** - From 906 lines to <100 lines in Update()
✅ **Better Maintainability** - Each intent is self-contained
✅ **Improved Testability** - Intent logic separated from app logic
✅ **Type Safety** - Compile-time safety via IntentResult[T]
✅ **Clearer Navigation** - Explicit state machines per intent
✅ **Reusability** - Intents can be composed and reused

---

## Architecture Overview

### Current Screen-Based Architecture

```
┌─────────────────────────────────────────────────┐
│              Main Application Model             │
│  (1766 lines, 906-line Update, 37+ fields)    │
├─────────────────────────────────────────────────┤
│  FormModel    ListModel    DetailsModel  ...   │
│  (31 screens, each with its own model)         │
├─────────────────────────────────────────────────┤
│  Navigation Logic (currentScreen, previousScreen)
│  Message Handling (31 case statements)          │
│  Breadcrumb Management                          │
│  Workflow State Tracking                        │
└─────────────────────────────────────────────────┘
```

### Target Intent-Based Architecture

```
┌────────────────────────────────────────────────┐
│      Main Application Model (Minimal)          │
│  (<200 lines, <50-line Update, <5 fields)     │
├────────────────────────────────────────────────┤
│         IntentRouter (Message Dispatcher)      │
├────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐           │
│  │CaptureEvent  │  │BrowseTimeline│  ...     │
│  │  Intent      │  │   Intent     │           │
│  └──────────────┘  └──────────────┘           │
│  (Self-contained state machines)               │
├────────────────────────────────────────────────┤
│  Global Shortcuts (Quit, Help, Back, Menu)    │
│  App-Level State (Workflow, Config)            │
└────────────────────────────────────────────────┘
```

---

## Current State Analysis

### Screen Inventory

**31 Screens Currently Defined**:

| Category | Screens | Model Fields | Status |
|----------|---------|--------------|--------|
| **Core Workflows** | CaptureScreen, ListScreen, ViewScreen | formModel, listModel, detailsModel | High Priority |
| **CV Management** | CVConfigManagerScreen, CVGeneratorScreen, CVPreviewScreen, CVListScreen, CVExportDialogScreen, CVExportSuccessScreen, CVExportProgressScreen | cvConfigManagerModel, cvGeneratorModel, cvPreviewModel, cvListModel, cvExportDialogModel, cvExportSuccessModel, cvExportProgressModel | High Priority |
| **Burst Management** | BurstListScreen, BurstDetailsScreen, BurstEditorScreen, BurstSuggestionScreen | burstListModel, burstDetailsModel, burstEditorModel, burstSuggestionModel | Medium Priority |
| **Fact Management** | FactListScreen, FactDetailsScreen, FactEditorScreen, FactActionMenuScreen, FactsResultsScreen | factListModel, factDetailsModel, factEditorModel, factActionMenuModel, factsResultsModel | Medium Priority |
| **Import** | ImportReviewScreen, ImportProgressScreen | importReviewModel, importProgressModel, importFilePath, importResult, importService | Medium Priority |
| **Metadata** | MetadataReviewScreen, MetadataEditorScreen | metadataReviewModel, metadataEditorModel | Medium Priority |
| **Menu/UI** | MainMenuScreen, HomeScreen, HelpScreen, ActionMenuScreen, FactActionMenuScreen, ConfirmationScreen, SuccessScreen, QuitScreen | menuModel, helpModel, actionMenuModel, confirmationDialog, successModel | Low Priority |
| **Other** | BulkOperationsScreen | bulkOperationsModel | Low Priority |

### Update() Method Breakdown

**906 lines total**, organized as:

1. **Global Shortcuts** (50 lines)
   - BackMsg, QuitMsg, HelpMsg, MainMenuMsg
   - ⚠️ Don't check `inIntentMode`

2. **Message Type Handlers** (150 lines)
   - ViewEventMsg, EditEventMsg, BulkOperationsMsg
   - BurstSuggestionsTriggeredMsg, etc.

3. **Screen-Specific Logic** (700+ lines)
   - CaptureScreen updates (100+ lines)
   - ListScreen updates (100+ lines)
   - CVConfigManagerScreen updates (100+ lines)
   - And 28 more screens...

4. **Model Update Delegation** (50 lines)
   - Calls Update() on current screen's model
   - Handles resize messages

### Model Fields Breakdown

**37+ fields** organized as:

1. **Services** (3 fields)
   - cliService, service (careerService), importService

2. **Configuration** (3 fields)
   - configManager, cvGenerationService, cvExportService

3. **Navigation** (3 fields)
   - currentScreen, previousScreen, screenBeforeActionMenu
   - breadcrumbs

4. **Workflow** (1 field)
   - workflowState

5. **Screen Models** (27+ fields)
   - One field for each screen's UI component

6. **Temporary State** (3+ fields)
   - deleteEventID, importFilePath, importResult

7. **Intent System** (2 fields)
   - intentRouter, inIntentMode

8. **Rendering** (2 fields)
   - width, height

9. **Logging** (1 field)
   - logger

---

## Intent Implementation Status

### Implemented Intents ✅

All 5 core intents are **fully implemented** with state machines, views, and tests:

| Intent | File | Status | Tests | Coverage |
|--------|------|--------|-------|----------|
| **CaptureEvent** | capture_event_intent.go | ✅ COMPLETE | 30+ specs | 88.3% |
| **BrowseTimeline** | browse_timeline_intent.go | ✅ COMPLETE | 37 tests | >90% |
| **GenerateCV** | generate_cv_intent.go | ✅ COMPLETE | 41 tests | >90% |
| **ExportArtifact** | export_artifact_intent.go | ✅ COMPLETE | 411 tests | >90% |
| **ConfigureSystem** | configure_system_intent.go | ✅ COMPLETE | 400+ tests | >90% |

### Router Integration Status

| Component | Status | Details |
|-----------|--------|---------|
| IntentRouter | ✅ READY | Initialized, all intents registered |
| Intent Activation | ⚠️ NOT INTEGRATED | activateIntent() method exists but not called |
| Message Routing | ⚠️ NOT INTEGRATED | handleIntentMessage() exists but not called |
| View Delegation | ⚠️ NOT INTEGRATED | View() doesn't delegate to router |
| Global Shortcuts | ⚠️ NOT INTEGRATED | Shortcuts don't check inIntentMode |

---

## Migration Strategy

### Phase Overview

```
Phase 1: Foundation (1 week)
├─ Implement Priority 1 fixes (Update/View delegation)
├─ Activate intent mode in Update() and View()
├─ Update global shortcuts to check inIntentMode
└─ Establish dual-mode operation

Phase 2: High-Value Screens (2 weeks)
├─ Migrate CaptureScreen → CaptureEvent Intent
├─ Migrate ListScreen/ViewScreen → BrowseTimeline Intent
├─ Migrate CV screens → GenerateCV/ExportArtifact Intents
├─ Add menu items to activate intents
└─ Test end-to-end workflows

Phase 3: Remaining Screens (1.5 weeks)
├─ Migrate Burst screens → New BurstManagement Intent
├─ Migrate Fact screens → New FactManagement Intent
├─ Migrate Import screens → New ImportWizard Intent
├─ Migrate Metadata screens → New MetadataEditor Intent
└─ Consolidate menu and navigation

Phase 4: Cleanup & Optimization (1 week)
├─ Remove legacy screen code
├─ Remove unused model fields
├─ Simplify Update() and View() methods
├─ Optimize performance
└─ Final testing and documentation
```

---

## Detailed Migration Plan

### Phase 1: Foundation (1 week) ⚠️ CRITICAL

**Objective**: Enable the intent system to work end-to-end

#### 1.1 Implement Update() Delegation

**File**: `internal/cli/app/app.go` (line 368)

**Current Code**:
```go
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle quit and back messages first
    switch msg.(type) {
    case models.BackMsg:
        // 30+ lines of screen-specific back handling
        ...
    // 900+ lines of other logic
}
```

**New Code**:
```go
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Check intent mode FIRST - delegate to router if active
    if m.inIntentMode {
        return m.handleIntentMessage(msg)
    }

    // Handle global shortcuts
    switch msg.(type) {
    case models.QuitMsg:
        return m, tea.Quit
    case models.HelpMsg:
        m.previousScreen = m.currentScreen
        m.currentScreen = HelpScreen
        return m, nil
    case models.MainMenuMsg:
        m.previousScreen = m.currentScreen
        m.currentScreen = MainMenuScreen
        return m, nil
    case models.BackMsg:
        // Screen-based back navigation
        ...
    }

    // Rest of screen-based logic (unchanged for now)
    ...
}
```

**Effort**: 30 minutes
**Testing**: Run existing tests, verify no regression

#### 1.2 Implement View() Delegation

**File**: `internal/cli/app/app.go` (line 1273)

**Current Code**:
```go
func (m *Model) View() string {
    switch m.currentScreen {
    case MainMenuScreen:
        if m.menuModel != nil {
            return m.menuModel.View()
        }
        // ... 31 case statements
    }
}
```

**New Code**:
```go
func (m *Model) View() string {
    // Check intent mode FIRST - delegate to router if active
    if m.inIntentMode && m.intentRouter != nil {
        return m.intentRouter.View()
    }

    // Screen-based rendering (unchanged for now)
    switch m.currentScreen {
    case MainMenuScreen:
        if m.menuModel != nil {
            return m.menuModel.View()
        }
        // ... rest of screens
    }
}
```

**Effort**: 15 minutes
**Testing**: Run existing tests, verify no regression

#### 1.3 Update Global Shortcuts

**File**: `internal/cli/app/app.go` (line 400+)

**Current Code**:
```go
case models.BackMsg:
    // 30+ lines of screen-specific handling
    // Doesn't check inIntentMode
```

**New Code**:
```go
case models.BackMsg:
    if m.inIntentMode {
        // Delegate to router's back navigation
        cmd, err := m.intentRouter.Back()
        if err != nil {
            m.deactivateIntent()
        }
        return m, cmd
    }

    // Screen-based back handling (unchanged)
    if m.currentScreen == ListScreen {
        // ...
    }
```

**Effort**: 30 minutes
**Testing**: Test back navigation in intent mode

#### 1.4 Add Menu Items to Activate Intents

**File**: `internal/cli/app/app.go` (handleMenuItemSelection method, line 1611)

**Current Code**:
```go
func (m *Model) handleMenuItemSelection(key string) (tea.Model, tea.Cmd) {
    switch key {
    case "c":
        // Capture Career Event (screen-based)
        m.previousScreen = m.currentScreen
        m.currentScreen = CaptureScreen
        m.formModel = models.NewFormModel(m.cliService)
        return m, nil
    // ... other cases
    }
}
```

**New Code**:
```go
func (m *Model) handleMenuItemSelection(key string) (tea.Model, tea.Cmd) {
    switch key {
    case "c":
        // Capture Career Event (intent-based)
        return m, m.activateIntent("capture_event", make(map[string]interface{}))
    case "l":
        // Browse Timeline (intent-based)
        return m, m.activateIntent("browse_timeline", make(map[string]interface{}))
    case "g":
        // Generate CV (intent-based)
        return m, m.activateIntent("generate_cv", make(map[string]interface{}))
    case "x":
        // Export Artifact (intent-based)
        return m, m.activateIntent("export_artifact", make(map[string]interface{}))
    case "s":
        // Configure System (intent-based)
        return m, m.activateIntent("configure_system", make(map[string]interface{}))
    // ... legacy screens for now
    }
}
```

**Effort**: 1 hour
**Testing**: Test menu activation of each intent

#### 1.5 Establish Dual-Mode Operation

**Goals**:
- ✅ Intents work when activated
- ✅ Legacy screens still work
- ✅ Can navigate between intent mode and screen mode
- ✅ No breaking changes

**Testing**:
- Run full test suite
- Verify no regressions
- Test intent activation
- Test back navigation

**Effort**: 2 hours
**Total Phase 1 Effort**: ~5 hours

---

### Phase 2: High-Value Screens (2 weeks)

**Objective**: Migrate the 5 most-used screens to intents

#### 2.1 Verify CaptureEvent Intent

**Status**: ✅ COMPLETE - capture_event_intent.go fully implemented

**Tasks**:
1. ✅ Review capture_event_intent.go implementation
2. ✅ Verify all states implemented
3. ✅ Verify views implemented
4. ✅ Verify result handling
5. ✅ Run tests (30+ specs passing)

**Action**: Activate in menu with Phase 1 changes
**Testing**: Test capture workflow end-to-end
**Effort**: 2 hours (testing and verification)

#### 2.2 Verify BrowseTimeline Intent

**Status**: ✅ COMPLETE - browse_timeline_intent.go fully implemented

**Tasks**:
1. ✅ Review browse_timeline_intent.go implementation
2. ✅ Verify all states implemented
3. ✅ Verify views implemented
4. ✅ Verify result handling
5. ✅ Run tests (37 tests passing)

**Action**: Activate in menu with Phase 1 changes
**Testing**: Test timeline browsing end-to-end
**Effort**: 2 hours (testing and verification)

#### 2.3 Verify GenerateCV Intent

**Status**: ✅ COMPLETE - generate_cv_intent.go fully implemented

**Tasks**:
1. ✅ Review generate_cv_intent.go implementation
2. ✅ Verify all states implemented
3. ✅ Verify views implemented
4. ✅ Verify result handling
5. ✅ Run tests (41 tests passing)

**Action**: Activate in menu with Phase 1 changes
**Testing**: Test CV generation end-to-end
**Effort**: 2 hours (testing and verification)

#### 2.4 Verify ExportArtifact Intent

**Status**: ✅ COMPLETE - export_artifact_intent.go fully implemented

**Tasks**:
1. ✅ Review export_artifact_intent.go implementation
2. ✅ Verify all states implemented
3. ✅ Verify views implemented
4. ✅ Verify result handling
5. ✅ Run tests (411 tests passing)

**Action**: Activate in menu with Phase 1 changes
**Testing**: Test artifact export end-to-end
**Effort**: 2 hours (testing and verification)

#### 2.5 Verify ConfigureSystem Intent

**Status**: ✅ COMPLETE - configure_system_intent.go fully implemented

**Tasks**:
1. ✅ Review configure_system_intent.go implementation
2. ✅ Verify all states implemented
3. ✅ Verify views implemented
4. ✅ Verify result handling
5. ✅ Run tests (400+ tests passing)

**Action**: Activate in menu with Phase 1 changes
**Testing**: Test system configuration end-to-end
**Effort**: 2 hours (testing and verification)

#### 2.6 Enhance Result Handlers

**File**: `internal/cli/app/app.go` (lines 238-304)

**Current Code**:
```go
router.RegisterResultHandler("browse_timeline", func(result *intents.IntentResult[interface{}]) tea.Cmd {
    if result.Status == intents.Completed {
        log.Info("BrowseTimeline intent completed")
    }
    return func() tea.Msg {
        return models.BackMsg{}
    }
})
```

**New Code**:
```go
router.RegisterResultHandler("browse_timeline", func(result *intents.IntentResult[interface{}]) tea.Cmd {
    if result.Status == intents.Completed {
        if timelineResult, ok := result.Data.(*intents.BrowseTimelineResult); ok {
            log.Info("BrowseTimeline intent completed with event: %s", timelineResult.SelectedEventID)
            // Update app state based on result
            if timelineResult.SelectedEventID != "" {
                // Could navigate to event details or trigger another action
                return func() tea.Msg {
                    return models.ViewEventMsg{EventID: timelineResult.SelectedEventID}
                }
            }
        }
    } else if result.Status == intents.Cancelled {
        log.Info("BrowseTimeline intent cancelled")
    } else if result.Status == intents.Failed {
        log.Error("BrowseTimeline intent failed: %v", result.Error)
    }

    // Return to screen mode
    m.deactivateIntent()
    return nil
})
```

**Effort**: 2 hours
**Testing**: Test result handling for each intent

#### 2.7 Remove Legacy Screen Activation

**File**: `internal/cli/app/app.go` (handleMenuItemSelection method)

**Current Code**:
```go
case "l":
    // List Events
    ctx := context.Background()
    m.listModel = models.NewListModel(m.service, ctx)
    m.previousScreen = m.currentScreen
    m.currentScreen = ListScreen
    return m, nil
```

**New Code**:
```go
case "l":
    // List Events (now via intent)
    return m, m.activateIntent("browse_timeline", make(map[string]interface{}))
```

**Effort**: 1 hour
**Testing**: Verify menu items activate intents correctly

#### 2.8 Integration Testing

**Tasks**:
1. Test each intent activation from menu
2. Test navigation within each intent
3. Test back navigation from intents
4. Test result handling
5. Test switching between intents

**Effort**: 4 hours

**Total Phase 2 Effort**: ~20 hours (2.5 weeks)

---

### Phase 3: Remaining Screens (1.5 weeks)

**Objective**: Migrate remaining screens to intents or consolidate

#### 3.1 Create BurstManagement Intent

**Screens to Migrate**:
- BurstListScreen
- BurstDetailsScreen
- BurstEditorScreen
- BurstSuggestionScreen

**Implementation**:
1. Create `internal/cli/intents/burst_management.go` - state machine
2. Create `internal/cli/intents/burst_management_intent.go` - implementation
3. Create `internal/cli/intents/burst_management_test.go` - tests
4. Register with router in app.go
5. Add menu item

**Effort**: 8 hours
**Testing**: 3 hours

#### 3.2 Create FactManagement Intent

**Screens to Migrate**:
- FactListScreen
- FactDetailsScreen
- FactEditorScreen
- FactsResultsScreen
- FactActionMenuScreen

**Implementation**:
1. Create `internal/cli/intents/fact_management.go`
2. Create `internal/cli/intents/fact_management_intent.go`
3. Create `internal/cli/intents/fact_management_test.go`
4. Register with router
5. Add menu item

**Effort**: 10 hours
**Testing**: 3 hours

#### 3.3 Create ImportWizard Intent

**Screens to Migrate**:
- ImportReviewScreen
- ImportProgressScreen

**Implementation**:
1. Create `internal/cli/intents/import_wizard.go`
2. Create `internal/cli/intents/import_wizard_intent.go`
3. Create `internal/cli/intents/import_wizard_test.go`
4. Register with router
5. Add menu item

**Effort**: 8 hours
**Testing**: 3 hours

#### 3.4 Create MetadataEditor Intent

**Screens to Migrate**:
- MetadataReviewScreen
- MetadataEditorScreen

**Implementation**:
1. Create `internal/cli/intents/metadata_editor.go`
2. Create `internal/cli/intents/metadata_editor_intent.go`
3. Create `internal/cli/intents/metadata_editor_test.go`
4. Register with router
5. Add menu item

**Effort**: 6 hours
**Testing**: 2 hours

#### 3.5 Consolidate Menu and Navigation

**Tasks**:
1. Update menu model to show intent-based options
2. Remove screen-based menu items
3. Simplify navigation logic
4. Update breadcrumbs for intent mode

**Effort**: 4 hours
**Testing**: 2 hours

**Total Phase 3 Effort**: ~50 hours (1.5 weeks)

---

### Phase 4: Cleanup & Optimization (1 week)

**Objective**: Remove legacy code and optimize

#### 4.1 Remove Legacy Screen Code

**Tasks**:
1. Remove screen constants (keep only HomeScreen, MainMenuScreen, HelpScreen, QuitScreen)
2. Remove screen-specific model fields
3. Remove screen-specific Update() logic
4. Remove screen-specific View() cases
5. Remove unused models

**Effort**: 4 hours
**Testing**: 2 hours

#### 4.2 Simplify app.go

**Current Stats**:
- 1766 lines total
- 906 lines in Update()
- 37+ model fields

**Target Stats**:
- ~400 lines total
- ~50 lines in Update()
- 5-7 model fields

**Tasks**:
1. Remove 900+ lines of screen logic
2. Remove 30+ model field declarations
3. Simplify NewModel() initialization
4. Consolidate navigation logic

**Effort**: 4 hours
**Testing**: 2 hours

#### 4.3 Performance Optimization

**Tasks**:
1. Profile rendering performance
2. Optimize hot paths
3. Reduce allocations
4. Cache expensive computations

**Effort**: 3 hours
**Testing**: 2 hours

#### 4.4 Final Testing

**Tasks**:
1. Run full test suite
2. Test all workflows end-to-end
3. Performance benchmarking
4. User acceptance testing

**Effort**: 4 hours

#### 4.5 Documentation

**Tasks**:
1. Update README.md
2. Update CLI_GUIDE.md
3. Update architecture documentation
4. Create migration guide

**Effort**: 3 hours

**Total Phase 4 Effort**: ~24 hours (1 week)

---

## Implementation Timeline

### Recommended Schedule

```
Week 1: Phase 1 Foundation
  Mon-Tue: Implement Update() and View() delegation (5 hours)
  Wed:     Update global shortcuts (3 hours)
  Thu:     Add menu items (2 hours)
  Fri:     Testing and verification (4 hours)

Week 2-3: Phase 2 High-Value Screens
  Week 2: Verify and test 5 core intents (12 hours)
  Week 3: Enhance result handlers and remove legacy (8 hours)

Week 4: Phase 3 Remaining Screens (Part 1)
  Create BurstManagement intent (11 hours)

Week 5: Phase 3 Remaining Screens (Part 2)
  Create FactManagement intent (13 hours)
  Create ImportWizard intent (11 hours)

Week 6: Phase 3 Remaining Screens (Part 3)
  Create MetadataEditor intent (8 hours)
  Consolidate menu (6 hours)

Week 7: Phase 4 Cleanup
  Remove legacy code (6 hours)
  Simplify app.go (6 hours)
  Performance optimization (5 hours)
  Final testing (4 hours)
```

**Total Effort**: ~120 hours (3 weeks full-time, 6 weeks part-time)

---

## Dual-Mode Operation Strategy

### During Migration Period

The application will support **both** screen-based and intent-based navigation:

```
┌─────────────────────────────────────┐
│   Update(msg tea.Msg)               │
├─────────────────────────────────────┤
│ if inIntentMode:                    │
│   return handleIntentMessage(msg)   │
│ else:                               │
│   return handleScreenMessage(msg)   │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│   View()                            │
├─────────────────────────────────────┤
│ if inIntentMode:                    │
│   return intentRouter.View()        │
│ else:                               │
│   return renderScreen()             │
└─────────────────────────────────────┘
```

### Benefits
✅ **Zero Breaking Changes** - Existing screen-based code continues to work
✅ **Gradual Migration** - Can migrate one screen at a time
✅ **Testing Flexibility** - Can test new intents in parallel
✅ **Risk Mitigation** - Can fall back to screens if needed

### Risks to Monitor
⚠️ **Code Duplication** - Some functionality may be duplicated temporarily
⚠️ **Navigation Complexity** - Switching between modes can be confusing
⚠️ **Testing Burden** - Need to test both modes

### Mitigation Strategies
✅ **Clear Separation** - Keep intent mode and screen mode completely separate
✅ **Documentation** - Document dual-mode behavior clearly
✅ **Automated Testing** - Test both paths automatically
✅ **Gradual Removal** - Remove screen code in phases, not all at once

---

## Risk Analysis

### High-Risk Areas

| Risk | Severity | Mitigation |
|------|----------|-----------|
| Breaking existing workflows | HIGH | Keep dual-mode operation, extensive testing |
| Migration complexity | HIGH | Phase approach, clear checkpoints |
| Performance regression | MEDIUM | Benchmark before/after, optimize hot paths |
| Intent registration issues | MEDIUM | Comprehensive testing, clear error handling |
| Navigation confusion | MEDIUM | Clear documentation, visual indicators |

### Testing Strategy

**Before Phase 1**:
- ✅ Review all tests pass
- ✅ Benchmark current performance

**During Each Phase**:
- ✅ Run full test suite after each change
- ✅ Test new functionality
- ✅ Test existing workflows not affected

**After Phase 4**:
- ✅ Full regression testing
- ✅ Performance benchmarking
- ✅ User acceptance testing

---

## Success Criteria

### Phase 1 Success
- ✅ Intent mode Update() delegation works
- ✅ Intent mode View() delegation works
- ✅ Back navigation works in intent mode
- ✅ Menu items activate intents
- ✅ All existing tests pass
- ✅ No performance regression

### Phase 2 Success
- ✅ All 5 core intents activate from menu
- ✅ Each intent workflow completes end-to-end
- ✅ Result handlers work correctly
- ✅ Legacy screens still accessible
- ✅ All tests pass
- ✅ No breaking changes

### Phase 3 Success
- ✅ All remaining intents created and tested
- ✅ Menu consolidated
- ✅ Navigation simplified
- ✅ All workflows accessible via intents
- ✅ All tests pass

### Phase 4 Success
- ✅ Legacy screen code removed
- ✅ app.go simplified to <400 lines
- ✅ Update() method <50 lines
- ✅ Model fields reduced to <10
- ✅ All tests pass
- ✅ Performance benchmarks met
- ✅ Documentation updated

---

## Rollback Plan

If critical issues are discovered:

1. **During Phase 1**:
   - Revert Update() and View() changes
   - Keep legacy screen code untouched
   - No data loss possible

2. **During Phase 2-3**:
   - Keep dual-mode operation
   - Remove problematic intent
   - Continue with other intents

3. **During Phase 4**:
   - Restore removed code from git history
   - Revert to Phase 3 state
   - Fix issue and re-test

---

## Appendix: File Changes Summary

### Files to Modify

| File | Changes | Lines | Effort |
|------|---------|-------|--------|
| `internal/cli/app/app.go` | Add Update/View delegation, update menu, remove legacy | 1766 | 20 hours |
| `internal/cli/intents/router.go` | Already complete | 0 | 0 hours |
| `internal/cli/intents/contract.go` | Already complete | 0 | 0 hours |

### Files to Create

| File | Purpose | Lines | Effort |
|------|---------|-------|--------|
| `internal/cli/intents/burst_management.go` | Burst intent model | 200 | 8 hours |
| `internal/cli/intents/burst_management_intent.go` | Burst intent implementation | 400 | |
| `internal/cli/intents/burst_management_test.go` | Burst intent tests | 300 | |
| `internal/cli/intents/fact_management.go` | Fact intent model | 200 | 10 hours |
| `internal/cli/intents/fact_management_intent.go` | Fact intent implementation | 400 | |
| `internal/cli/intents/fact_management_test.go` | Fact intent tests | 300 | |
| `internal/cli/intents/import_wizard.go` | Import intent model | 150 | 8 hours |
| `internal/cli/intents/import_wizard_intent.go` | Import intent implementation | 350 | |
| `internal/cli/intents/import_wizard_test.go` | Import intent tests | 250 | |
| `internal/cli/intents/metadata_editor.go` | Metadata intent model | 150 | 6 hours |
| `internal/cli/intents/metadata_editor_intent.go` | Metadata intent implementation | 300 | |
| `internal/cli/intents/metadata_editor_test.go` | Metadata intent tests | 200 | |

### Files to Delete

| File | Reason |
|------|--------|
| Various legacy screen models | Consolidated into intents |

---

## Conclusion

This migration plan provides a **structured, phased approach** to transform KaRiya from a monolithic screen-based architecture to a clean, modular intent-driven architecture.

**Key Advantages**:
- ✅ **Phased Approach** - Can be done incrementally without breaking changes
- ✅ **Dual-Mode Operation** - Allows parallel development and testing
- ✅ **Clear Checkpoints** - Each phase has well-defined success criteria
- ✅ **Risk Mitigation** - Multiple safeguards and rollback options
- ✅ **Comprehensive** - Covers all 31 screens and all logic

**Expected Outcome**:
- Reduced app.go from 1766 to ~400 lines
- Reduced Update() from 906 to ~50 lines
- Reduced model fields from 37+ to ~7
- Improved maintainability and testability
- No breaking changes during migration

**Timeline**: 4-6 weeks (phased, part-time compatible)

---

**Document Version**: 1.0
**Status**: PLANNING PHASE - Ready for Implementation
**Next Step**: Begin Phase 1 (Priority 1 fixes in Update/View)


