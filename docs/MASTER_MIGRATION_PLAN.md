# Master Migration Plan

**Purpose**: Single source of truth for migrating KaRiya to the new architecture  
**Created**: 2026-01-24  
**Status**: IN PROGRESS

---

## Executive Summary

### Current State

| Metric | Value | Target |
|--------|-------|--------|
| Intents migrated | **1 of 7** (14%) | 100% |
| Legacy intent lines | **9,800** | 0 |
| Direct `huh` imports | **27 files** | 0 |
| Incorrect form state checks | **22 files** | 0 |
| TODOs in code | **16** | 0 |
| Modals in deprecated `components/` | **15** | 0 |

### Migration Progress

```
[==========----------] 14% Complete (1/7 intents)

Migrated:
  [x] browse_timeline (reference implementation)

Pending:
  [ ] fact_management (610 lines) - QUICK WIN
  [ ] configure_system (657 lines)
  [ ] burst_management (1,488 lines)
  [ ] capture_event (1,884 lines)
  [ ] manage_skills (2,415 lines)
  [ ] generate_cv (2,746 lines)
```

---

## Table of Contents

1. [Migration Priorities](#migration-priorities)
2. [Phase 1: Foundation](#phase-1-foundation)
3. [Phase 2: Quick Wins](#phase-2-quick-wins)
4. [Phase 3: Medium Complexity](#phase-3-medium-complexity)
5. [Phase 4: High Complexity](#phase-4-high-complexity)
6. [Phase 5: Pattern Standardization](#phase-5-pattern-standardization)
7. [Phase 6: Cleanup](#phase-6-cleanup)
8. [Generator Script Specification](#generator-script-specification)
9. [Enforcement Strategy](#enforcement-strategy)
10. [Progress Tracking](#progress-tracking)

---

## Screen/Modal Extraction Overview

**Reference**: [Screen Modal Extraction Guide](guides/SCREEN_MODAL_EXTRACTION_GUIDE.md)

Based on `browse_timeline` patterns, each legacy intent maps to:

| Intent | Screens | Custom Modals | UIKit Modals | Target Lines |
|--------|---------|---------------|--------------|--------------|
| fact_management | 1 new | 2 new | 2 | ~300 |
| configure_system | 6 existing | 3 move | 1 | ~350 |
| burst_management | 2 new | 2 new | 1 | ~400 |
| capture_event | 4 existing | 0 | 1 | ~400 |
| manage_skills | 3 (2 exist) | 4 new | 1 | ~450 |
| generate_cv | 5 existing | 2 wizards | 2 | ~500 |

**Total Expected Reduction**: 9,800 lines → ~2,400 lines (75%)

### Key Patterns from browse_timeline

**Screens** = Full-page views (List, Detail, Form)
- Created in `screens/{feature}/*.go`
- Implement `screens.Screen` interface
- Return `ScreenResult` from `Update()`

**Modals** = Overlays on current screen (Filter, Sort, Edit, Confirm)
- Feature-specific: `screens/{feature}/modals/*.go`
- Generic: `uikit/feedback/*.go` (ConfirmModal, ErrorModal, etc.)
- Rendered via `behaviors.RenderModalOverlay()`

---

## Migration Priorities

### Priority Matrix

| Priority | Task | Effort | Impact | Status |
|----------|------|--------|--------|--------|
| **P0** | Generator script | 2h | HIGH | Pending |
| **P1** | Migrate `fact_management` | 2h | MEDIUM | Pending |
| **P2** | Migrate `configure_system` | 4h | MEDIUM | Pending |
| **P3** | Standardize form patterns | 4h | HIGH | Pending |
| **P4** | Migrate `burst_management` | 6h | HIGH | Pending |
| **P5** | Move modals from `components/` | 4h | MEDIUM | Pending |
| **P6** | Migrate `capture_event` | 8h | HIGH | Pending |
| **P7** | Migrate `manage_skills` | 10h | HIGH | Pending |
| **P8** | Migrate `generate_cv` | 12h | HIGH | Pending |

**Total Estimated Effort**: ~52 hours

---

## Phase 1: Foundation

**Goal**: Create tooling and establish patterns  
**Duration**: 4 hours  
**Deliverables**: Generator script, validation improvements

### 1.1 Generator Script (P0)

Create `make new-intent NAME=feature` command.

**Files to generate**:
```
intents/{feature}/
├── context.go
├── result.go
├── constants.go
├── messages.go
├── intent.go
└── intent_test.go

screens/{feature}/
├── list_screen.go
├── detail_screen.go
└── modals/
    └── .gitkeep
```

**See**: [Generator Script Specification](#generator-script-specification)

**Status**: [ ] Not Started

### 1.2 Enforcement Improvements

Add grandfather clause for legacy code:
- Existing >600 line files: WARN instead of BLOCK
- New violations: BLOCK immediately
- Add migration deadline tracking

**Status**: [ ] Not Started

---

## Phase 2: Quick Wins

**Goal**: Migrate smallest intents first  
**Duration**: 6 hours  
**Deliverables**: 2 intents migrated

### 2.1 Migrate `fact_management` (P1)

**Current State**:
- File: `intents/fact_management_intent.go`
- Lines: **610** (only 10 over limit)
- Context file: `fact_management.go` (253 lines)
- Associated screens: None dedicated

**States → Screen/Modal Mapping** (from browse_timeline pattern):
| State | Type | Target |
|-------|------|--------|
| `list` | SCREEN | `screens/facts/list_screen.go` |
| `view` | MODAL | `screens/facts/modals/detail_modal.go` |
| `editor` | MODAL | `screens/facts/modals/edit_modal.go` |
| `delete_confirm` | MODAL | `feedback.ConfirmModal` (UIKit) |
| `results` | MODAL | `feedback.SuccessModal` (UIKit) |

**Required Directory Structure**:
```
intents/fact_management/
├── context.go, result.go, constants.go, messages.go, intent.go

screens/facts/
├── list_screen.go         # NEW: TableBehavior list
└── modals/
    ├── detail_modal.go    # NEW: Fact detail viewer
    └── edit_modal.go      # NEW: Fact editor form
```

**Migration Steps**:
1. [ ] Create `intents/fact_management/` subdirectory
2. [ ] Extract `constants.go` (State enum)
3. [ ] Extract `messages.go` (*Msg types)
4. [ ] Extract `result.go` (output type)
5. [ ] Consolidate context into `context.go`
6. [ ] Create `screens/facts/list_screen.go` with TableBehavior
7. [ ] Create `screens/facts/modals/detail_modal.go`
8. [ ] Create `screens/facts/modals/edit_modal.go`
9. [ ] Refactor `intent.go` to broker pattern
10. [ ] Run `make check-intent-architecture`
9. [ ] Delete old files

**Estimated Effort**: 2 hours

**Status**: [ ] Not Started

### 2.2 Migrate `configure_system` (P2)

**Current State**:
- File: `intents/configure_system_intent.go`
- Lines: **657** (57 over limit)
- Context file: `configure_system.go` (**1,211 lines** - needs split!)
- Associated screens: `screens/configure/` (9 files) - **ALL EXIST**

**States → Screen/Modal Mapping**:
| State | Type | Existing File | Action |
|-------|------|---------------|--------|
| `select_domain` | SCREEN | `domain_select.go` | Keep |
| `edit_settings` | SCREEN | `edit_settings.go` | Keep |
| `review_changes` | SCREEN | `review_changes.go` | Keep |
| `confirm` | SCREEN | `confirm_screen.go` | Keep |
| `saving` | MODAL | `feedback.LoadingModal` | Use UIKit |
| `complete` | SCREEN | `complete_screen.go` | Keep |
| `failed` | SCREEN | `failed_screen.go` | Keep |

**Required Changes**:
```
screens/configure/
├── (6 screens exist - keep as-is)
└── modals/
    ├── confirm_modal.go        # MOVE from root
    ├── edit_settings_modal.go  # MOVE from root
    └── review_changes_modal.go # MOVE from root
```

**Migration Steps**:
1. [ ] Create `intents/configure_system/` subdirectory
2. [ ] Split large context.go (1,211 lines) into multiple files
3. [ ] Extract `constants.go` (State enum)
4. [ ] Extract `messages.go` (*Msg types)
5. [ ] Extract `result.go` (output type)
6. [ ] Create `screens/configure/modals/` and move modal files
7. [ ] Refactor `intent.go` to broker pattern
8. [ ] Verify existing screens work with new structure
9. [ ] Run `make check-intent-architecture`
10. [ ] Delete old files

**Estimated Effort**: 4 hours

**Status**: [ ] Not Started

---

## Phase 3: Medium Complexity

**Goal**: Migrate medium-sized intents  
**Duration**: 10 hours  
**Deliverables**: 1 intent migrated, pattern standardization

### 3.1 Standardize Form Patterns (P3)

**Current State**:
- 27 files with direct `huh` imports (should only be in `forms/`)
- 22 files using `form.State == huh.StateCompleted` (should use `forms.IsCompleted()`)

**Files Affected**:
```
screens/timeline/modals/*.go (5 files)
screens/configure/edit_settings*.go (2 files)
screens/base/form_screen.go
models/skill_form.go
models/capture_form.go
components/*.go (12 files)
intents/modals.go
```

**Migration Steps**:
1. [ ] Create list of all direct `huh` imports
2. [ ] For each file, replace with `forms` package usage
3. [ ] Replace `form.State == huh.StateCompleted` with `forms.IsCompleted(form)`
4. [ ] Replace `form.State == huh.StateAborted` with `forms.IsAborted(form)`
5. [ ] Run tests to verify behavior unchanged
6. [ ] Update any missing `forms` package wrappers

**Estimated Effort**: 4 hours

**Status**: [ ] Not Started

### 3.2 Migrate `burst_management` (P4)

**Current State**:
- File: `intents/burst_management_intent.go`
- Lines: **1,488** (2.5x over limit)
- Context file: `burst_management.go` (352 lines)
- Associated screens: None dedicated - **ALL NEW**

**States → Screen/Modal Mapping**:
| State | Type | Target | Notes |
|-------|------|--------|-------|
| `list` | SCREEN | `screens/bursts/list_screen.go` | TableBehavior |
| `view` | MODAL | `screens/bursts/modals/detail_modal.go` | Overlay |
| `editor` | MODAL | `screens/bursts/modals/edit_modal.go` | Form overlay |
| `delete_confirm` | MODAL | `feedback.ConfirmModal` | Use UIKit |
| `suggest` | SCREEN | `screens/bursts/suggest_screen.go` | Full-page |

**Required Directory Structure**:
```
intents/burst_management/
├── context.go, result.go, constants.go, messages.go, intent.go
├── handlers.go, helpers.go (optional)

screens/bursts/
├── list_screen.go         # NEW: TableBehavior list
├── suggest_screen.go      # NEW: AI suggestions
└── modals/
    ├── detail_modal.go    # NEW: Burst detail viewer
    └── edit_modal.go      # NEW: Burst editor form
```

**Migration Steps**:
1. [ ] Create `intents/burst_management/` subdirectory
2. [ ] Extract `constants.go` (State enum)
3. [ ] Extract `messages.go` (*Msg types)
4. [ ] Extract `result.go` (output type)
5. [ ] Extract `context.go` (business logic)
6. [ ] Create `screens/bursts/list_screen.go` with TableBehavior
7. [ ] Create `screens/bursts/suggest_screen.go`
8. [ ] Create `screens/bursts/modals/detail_modal.go`
9. [ ] Create `screens/bursts/modals/edit_modal.go`
10. [ ] Refactor `intent.go` to broker pattern (<400 lines)
11. [ ] Run `make check-intent-architecture`
12. [ ] Delete old files

**Estimated Effort**: 6 hours

**Status**: [ ] Not Started

---

## Phase 4: High Complexity

**Goal**: Migrate largest intents  
**Duration**: 30 hours  
**Deliverables**: 3 intents migrated

### 4.1 Move Modals from `components/` (P5)

**Current State**:
- 15 modals in deprecated `components/` location
- Should be in `uikit/feedback/` or `screens/*/modals/`

**Modal Inventory**:
```
components/
├── cv_config_wizard_modal.go (549 lines) → screens/cv/modals/
├── delete_modal.go → uikit/feedback/ (already have ConfirmModal)
├── error_modal.go → uikit/feedback/ (already have ErrorModal)
├── filter_modal.go → screens/*/modals/
├── help_modal.go → uikit/feedback/ (already have HelpModal)
├── info_modal.go → uikit/feedback/ (already have InfoModal)
├── loading_modal.go → uikit/feedback/ (already have LoadingModal)
├── success_modal.go → uikit/feedback/ (already have SuccessModal)
└── ... (7 more)
```

**Migration Steps**:
1. [ ] Audit all modals in `components/`
2. [ ] For generic modals: Verify `uikit/feedback/` equivalents exist
3. [ ] For feature modals: Move to `screens/*/modals/`
4. [ ] Update imports in dependent files
5. [ ] Run tests
6. [ ] Delete deprecated `components/` modals

**Estimated Effort**: 4 hours

**Status**: [ ] Not Started

### 4.2 Migrate `capture_event` (P6)

**Current State**:
- File: `intents/capture_event_intent.go`
- Lines: **1,884** (3.1x over limit)
- Context file: `capture_event.go` (161 lines)
- Uses deprecated `models.CaptureForm`
- Associated screens: `screens/capture/` (4 files) - **ALL EXIST**

**States → Screen/Modal Mapping**:
| State | Type | Existing File | Action |
|-------|------|---------------|--------|
| `choose_strategy` | SCREEN | `strategy_select.go` | Keep |
| `form` | SCREEN | `event_form_screen.go` | Keep |
| `review` | SCREEN | `event_review_screen.go` | Keep |
| `submit` | SCREEN | `event_submit_screen.go` | Keep |

**Required Changes**:
```
intents/capture_event/
├── context.go, result.go, constants.go, messages.go, intent.go
├── handlers.go, helpers.go, interfaces.go (optional)

screens/capture/
├── strategy_select.go       # EXISTS: Keep
├── event_form_screen.go     # EXISTS: Keep (uses forms directly)
├── event_review_screen.go   # EXISTS: Keep
├── event_submit_screen.go   # EXISTS: Keep
└── modals/
    └── .gitkeep            # No custom modals needed
```

**Migration Steps**:
1. [ ] Create `intents/capture_event/` subdirectory
2. [ ] Extract `constants.go` (State enum)
3. [ ] Extract `messages.go` (*Msg types)
4. [ ] Extract `result.go` (output type)
5. [ ] Extract `context.go` (business logic)
6. [ ] Extract `handlers.go`, `helpers.go`, `interfaces.go`
7. [ ] Replace `models.CaptureForm` usage with existing FormScreen
8. [ ] Verify existing screens work with new intent structure
9. [ ] Refactor `intent.go` to broker pattern (<400 lines)
10. [ ] Run `make check-intent-architecture`
11. [ ] Delete old files

**Estimated Effort**: 8 hours

**Status**: [ ] Not Started

### 4.3 Migrate `manage_skills` (P7)

**Current State**:
- File: `intents/manage_skills_intent.go`
- Lines: **2,415** (4x over limit)
- Context file: `manage_skills.go` (105 lines)
- Uses deprecated `models.SkillForm`
- **16 render methods** (should be max 2)
- Associated screens: `screens/skills/` (4 files) - **MOST EXIST**

**States → Screen/Modal Mapping**:
| State | Type | Existing/Target | Action |
|-------|------|-----------------|--------|
| `list` | SCREEN | `list.go` | EXISTS |
| `detail` | SCREEN | `detail.go` | EXISTS |
| `detail_events` | SCREEN | `events_screen.go` | NEW |
| `detail_event_detail` | MODAL | `modals/event_detail_modal.go` | NEW |
| `add` | MODAL | `modals/form_modal.go` | NEW (combine) |
| `edit` | MODAL | `modals/form_modal.go` | NEW (combine) |
| `delete` | MODAL | `feedback.ConfirmModal` | Use UIKit |
| `filter` | MODAL | `modals/filter_modal.go` | NEW |
| `sort` | MODAL | `modals/sort_modal.go` | NEW |

**Required Changes**:
```
intents/manage_skills/
├── context.go, result.go, constants.go, messages.go, intent.go
├── handlers.go, helpers.go, filters.go (optional)

screens/skills/
├── list.go                 # EXISTS: Keep
├── detail.go               # EXISTS: Keep
├── form.go                 # EXISTS: Review for modal conversion
├── delete.go               # EXISTS: Replace with UIKit
├── events_screen.go        # NEW: Events for skill
└── modals/
    ├── form_modal.go       # NEW: Add/Edit form overlay
    ├── filter_modal.go     # NEW: Filter menu
    ├── sort_modal.go       # NEW: Sort menu
    └── event_detail_modal.go # NEW: Event detail overlay
```

**Migration Steps**:
1. [ ] Create `intents/manage_skills/` subdirectory
2. [ ] Extract `constants.go` (State enum)
3. [ ] Extract `messages.go` (*Msg types)
4. [ ] Extract `result.go` (output type)
5. [ ] Extract `context.go` (business logic)
6. [ ] Create `screens/skills/events_screen.go`
7. [ ] Create `screens/skills/modals/` with 4 modals
8. [ ] Move ALL 16 render methods to screens/modals
9. [ ] Replace `models.SkillForm` with FormModal
10. [ ] Replace custom delete with `feedback.ConfirmModal`
11. [ ] Refactor `intent.go` to broker pattern (<400 lines)
12. [ ] Run `make check-intent-architecture`
13. [ ] Delete old files

**Estimated Effort**: 10 hours

**Status**: [ ] Not Started

### 4.4 Migrate `generate_cv` (P8)

**Current State**:
- File: `intents/generate_cv_intent.go`
- Lines: **2,746** (4.6x over limit - LARGEST)
- Context file: `generate_cv.go` (382 lines)
- **5 render methods** (should be max 2)
- Associated screens: `screens/cv/` (5 files) - **ALL EXIST**
- **20+ states** - Many can be combined into wizard modals

**States → Screen/Modal Mapping** (grouped by phase):

**Configuration Phase** (combine into wizard modal):
| State | Type | Action |
|-------|------|--------|
| `select_profile` | SCREEN | EXISTS: `profile_select.go` |
| `select_audience` | SCREEN | EXISTS: `audience_select.go` |
| `select_technology_focus` | MODAL | Combine into `config_wizard_modal.go` |
| `select_technologies` | MODAL | Combine into `config_wizard_modal.go` |
| `select_focus_area` | MODAL | Combine into `config_wizard_modal.go` |
| `select_skills_config` | MODAL | Combine into `config_wizard_modal.go` |
| `select_length_format` | MODAL | Combine into `config_wizard_modal.go` |

**Generation Phase**:
| State | Type | Action |
|-------|------|--------|
| `extracting` | MODAL | `feedback.LoadingModal` (UIKit) |
| `generating` | SCREEN | EXISTS: `generating.go` |

**Review Phase**:
| State | Type | Action |
|-------|------|--------|
| `preview` | SCREEN | EXISTS: `preview.go` |
| `review` | SCREEN | EXISTS: `review.go` |
| `confirm` | MODAL | `feedback.ConfirmModal` (UIKit) |

**Export Phase** (combine into wizard modal):
| State | Type | Action |
|-------|------|--------|
| `export_select_format` | MODAL | Combine into `export_modal.go` |
| `export_select_location` | MODAL | Combine into `export_modal.go` |
| `exporting` | MODAL | `feedback.LoadingModal` (UIKit) |
| `export_complete` | MODAL | `feedback.SuccessModal` (UIKit) |

**Required Changes**:
```
intents/generate_cv/
├── context.go, result.go, constants.go, messages.go, intent.go
├── handlers.go, helpers.go, interfaces.go (optional)

screens/cv/
├── profile_select.go       # EXISTS: Keep
├── audience_select.go      # EXISTS: Keep
├── generating.go           # EXISTS: Keep
├── preview.go              # EXISTS: Keep
├── review.go               # EXISTS: Keep
└── modals/
    ├── config_wizard_modal.go  # NEW: Multi-step config wizard
    └── export_modal.go         # NEW: Export wizard
```

**Key Insight**: 20+ states can be reduced to 5 screens + 2 wizard modals + UIKit modals

**Migration Steps**:
1. [ ] Create `intents/generate_cv/` subdirectory
2. [ ] Extract `constants.go` (consolidate 20+ states to ~10)
3. [ ] Extract `messages.go` (*Msg types)
4. [ ] Extract `result.go` (output type)
5. [ ] Extract `context.go` (business logic)
6. [ ] Create `screens/cv/modals/config_wizard_modal.go` (multi-step)
7. [ ] Create `screens/cv/modals/export_modal.go` (multi-step)
8. [ ] Move `cv_config_wizard_modal.go` from components/
9. [ ] Move ALL render methods to existing screens
10. [ ] Refactor `intent.go` to broker pattern (<400 lines)
11. [ ] Run `make check-intent-architecture`
12. [ ] Delete old files

**Estimated Effort**: 12 hours

**Status**: [ ] Not Started

---

## Phase 5: Pattern Standardization

**Goal**: Fix remaining pattern violations  
**Duration**: 4 hours  
**Deliverables**: All pattern checks passing

### 5.1 Fix `context.Background()` Usage

**Current State**:
- ~12 instances of direct `context.Background()` in intents
- Should use `i.getContext()` or `i.GetContext()`

**Files Affected**:
```
intents/browse_timeline/intent.go:412
intents/capture_event_intent.go:643, 658, 674, 841, 1235, 1253, 1276, 1595, 1614, 1631
intents/generate_cv_intent.go:1000, 2557
```

**Status**: [ ] Not Started

### 5.2 Remove TODO/FIXME Comments

**Current State**:
- 16 TODOs in production code
- Must be resolved before merge

**Files with TODOs**:
```
intents/manage_skills_intent.go:212
screens/contract.go:63
screens/capture/event_submit_screen.go:27
screens/base/detail_screen.go:176
screens/base/base_screen.go:68, 179
intents/generate_cv_intent.go:45
intents/capture_event_intent.go:588, 600, 1345, 1848
```

**Status**: [ ] Not Started

### 5.3 Fix Raw Lipgloss Usage

**Current State**:
- ~50 locations with excessive raw lipgloss
- Should use theme system

**Major Files**:
```
app/app.go (~25 occurrences)
intents/manage_skills_intent.go (7 occurrences)
models/form.go (10 occurrences)
uikit/feedback/info_modal.go (4 occurrences)
uikit/widgets/detail_view.go (4 occurrences)
```

**Status**: [ ] Not Started

---

## Phase 6: Cleanup

**Goal**: Final cleanup and validation  
**Duration**: 2 hours  
**Deliverables**: All checks passing, documentation updated

### 6.1 Cleanup Tasks

1. [ ] Delete all migrated flat files
2. [ ] Remove unused imports
3. [ ] Update all documentation references
4. [ ] Archive old migration docs
5. [ ] Run full test suite
6. [ ] Verify coverage >= 95%

### 6.2 Final Validation

```bash
# All must pass
make check-intent-architecture
make check-compliance
make golangci-lint
make test
make coverage
```

---

## Generator Script Specification

### Command

```bash
make new-intent NAME=my_feature
```

### Generated Structure

```
intents/my_feature/
├── context.go
├── result.go
├── constants.go
├── messages.go
├── intent.go
└── intent_test.go

screens/my_feature/
├── list_screen.go
├── detail_screen.go
└── modals/
    └── .gitkeep
```

### Template Files

#### context.go.template
```go
package {{.PackageName}}

import (
    "context"
)

// {{.ContextName}} holds input parameters and business logic.
type {{.ContextName}} struct {
    // Add input parameters here
}

// New{{.ContextName}} creates a new context.
func New{{.ContextName}}() *{{.ContextName}} {
    return &{{.ContextName}}{}
}

// Validate validates the context.
func (c *{{.ContextName}}) Validate() error {
    return nil
}
```

#### constants.go.template
```go
package {{.PackageName}}

// {{.StateName}} defines the state machine states.
type {{.StateName}} string

const (
    // StateInitial is the initial state.
    StateInitial {{.StateName}} = "initial"
    
    // StateList is the list view state.
    StateList {{.StateName}} = "list"
    
    // StateDetail is the detail view state.
    StateDetail {{.StateName}} = "detail"
)
```

#### messages.go.template
```go
package {{.PackageName}}

// ItemSelectedMsg is sent when an item is selected.
type ItemSelectedMsg struct {
    // Add fields here
}

// ErrorMsg is sent when an error occurs.
type ErrorMsg struct {
    Error error
}
```

#### result.go.template
```go
package {{.PackageName}}

// {{.ResultName}} is the output of the {{.FeatureName}} intent.
type {{.ResultName}} struct {
    // Add result fields here
}
```

#### intent.go.template
```go
package {{.PackageName}}

import (
    "github.com/baphled/kariya/internal/cli/behaviors"
    "github.com/baphled/kariya/internal/cli/intents"
    "github.com/baphled/kariya/internal/cli/screens"
    tea "github.com/charmbracelet/bubbletea"
)

// Ensure interface compliance.
var _ intents.ScreenResultHandler = (*{{.IntentName}})(nil)

// {{.IntentName}} orchestrates the {{.FeatureName}} workflow.
type {{.IntentName}} struct {
    *intents.BaseIntent
    
    // Context
    context *{{.ContextName}}
    
    // State machine
    state  {{.StateName}}
    active bool
    result *intents.IntentResult[*{{.ResultName}}]
    
    // Screens
    activeScreen screens.Screen
    
    // Modals
    modalRegistry *intents.ModalRegistry
}

// New{{.IntentName}} creates a new {{.FeatureName}} intent.
func New{{.IntentName}}(ctx *{{.ContextName}}) (*{{.IntentName}}, error) {
    if err := ctx.Validate(); err != nil {
        return nil, err
    }
    
    return &{{.IntentName}}{
        BaseIntent:    intents.NewBaseIntent(),
        context:       ctx,
        state:         StateInitial,
        active:        true,
        modalRegistry: intents.NewModalRegistry(),
    }, nil
}

// Init initializes the intent.
func (i *{{.IntentName}}) Init() tea.Cmd {
    // Initialize first screen
    return nil
}

// Update processes messages.
func (i *{{.IntentName}}) Update(msg tea.Msg) tea.Cmd {
    if !i.active {
        return nil
    }
    
    // Modal handling
    if cmd := i.modalRegistry.Update(msg); cmd != nil {
        return cmd
    }
    
    // Screen delegation
    if i.activeScreen != nil {
        cmd, result := i.activeScreen.Update(msg)
        if result != nil {
            return i.handleScreenResult(result)
        }
        return cmd
    }
    
    return nil
}

// View renders the intent.
func (i *{{.IntentName}}) View() string {
    if i.activeScreen == nil {
        return ""
    }
    
    baseView := i.activeScreen.View()
    return i.modalRegistry.RenderOverlay(baseView)
}

// Result returns the intent result.
func (i *{{.IntentName}}) Result() *intents.IntentResult[interface{}] {
    if i.result == nil {
        return nil
    }
    return &intents.IntentResult[interface{}]{
        Status: i.result.Status,
        Data:   i.result.Data,
    }
}

// ScreenResultHandler implementation

func (i *{{.IntentName}}) HandleCancel(result *screens.CancelResult) tea.Cmd {
    i.setCancelled()
    return nil
}

func (i *{{.IntentName}}) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
    // Handle navigation
    return nil
}

func (i *{{.IntentName}}) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
    // Handle form submission
    return nil
}

func (i *{{.IntentName}}) HandleError(result *screens.ErrorResult) tea.Cmd {
    // Show error modal
    return nil
}

// Private helpers

func (i *{{.IntentName}}) handleScreenResult(result screens.ScreenResult) tea.Cmd {
    return intents.NewScreenResultDispatcher(i).Dispatch(result)
}

func (i *{{.IntentName}}) setCancelled() {
    i.result = &intents.IntentResult[*{{.ResultName}}]{
        Status: intents.Cancelled,
    }
    i.active = false
}
```

### Implementation Script

```bash
#!/bin/bash
# scripts/new-intent.sh

set -e

NAME=$1
if [ -z "$NAME" ]; then
    echo "Usage: make new-intent NAME=feature_name"
    exit 1
fi

# Convert to various cases
PACKAGE_NAME=$(echo "$NAME" | tr '[:upper:]' '[:lower:]' | tr '-' '_')
FEATURE_NAME=$(echo "$NAME" | sed 's/_/ /g' | sed 's/\b\(.\)/\u\1/g' | tr -d ' ')
INTENT_NAME="${FEATURE_NAME}Intent"
CONTEXT_NAME="${FEATURE_NAME}Context"
RESULT_NAME="${FEATURE_NAME}Result"
STATE_NAME="${FEATURE_NAME}State"

# Create directories
mkdir -p "internal/cli/intents/${PACKAGE_NAME}"
mkdir -p "internal/cli/screens/${PACKAGE_NAME}"
mkdir -p "internal/cli/screens/${PACKAGE_NAME}/modals"

# Generate files from templates
for template in examples/intent_subdirectory_template/*.template; do
    filename=$(basename "$template" .template)
    output="internal/cli/intents/${PACKAGE_NAME}/${filename}"
    
    sed -e "s/{{.PackageName}}/${PACKAGE_NAME}/g" \
        -e "s/{{.FeatureName}}/${FEATURE_NAME}/g" \
        -e "s/{{.IntentName}}/${INTENT_NAME}/g" \
        -e "s/{{.ContextName}}/${CONTEXT_NAME}/g" \
        -e "s/{{.ResultName}}/${RESULT_NAME}/g" \
        -e "s/{{.StateName}}/${STATE_NAME}/g" \
        "$template" > "$output"
done

# Generate screen files
for template in examples/intent_subdirectory_template/screens/*.template; do
    filename=$(basename "$template" .template)
    output="internal/cli/screens/${PACKAGE_NAME}/${filename}"
    
    sed -e "s/{{.PackageName}}/${PACKAGE_NAME}/g" \
        -e "s/{{.FeatureName}}/${FEATURE_NAME}/g" \
        "$template" > "$output"
done

# Create .gitkeep for modals
touch "internal/cli/screens/${PACKAGE_NAME}/modals/.gitkeep"

echo "Created intent structure for: ${NAME}"
echo ""
echo "Files created:"
echo "  intents/${PACKAGE_NAME}/"
echo "    - context.go"
echo "    - result.go"
echo "    - constants.go"
echo "    - messages.go"
echo "    - intent.go"
echo "    - intent_test.go"
echo ""
echo "  screens/${PACKAGE_NAME}/"
echo "    - list_screen.go"
echo "    - detail_screen.go"
echo "    - modals/"
echo ""
echo "Next steps:"
echo "  1. Edit context.go to add input parameters"
echo "  2. Edit constants.go to define states"
echo "  3. Edit intent.go to implement logic"
echo "  4. Run: make check-intent-architecture"
```

---

## Enforcement Strategy

### Current Enforcement

| Check | What | Enforcement |
|-------|------|-------------|
| #17 | Subdirectory structure | WARN (should BLOCK) |
| #18 | File size >600 lines | BLOCK |
| #19 | >2 render methods | BLOCK |
| #20 | Types in wrong files | BLOCK |
| #21 | UIKit usage | WARN |
| #22 | models/ for forms | BLOCK |

### Proposed Changes

1. **Grandfather Legacy Code**
   - Existing violations: WARN
   - New violations: BLOCK
   - Track via `.legacy-intents` file

2. **Migration Deadline**
   - Add deadline to each legacy intent
   - WARN until deadline
   - BLOCK after deadline

3. **Progress Tracking**
   - Add `migration-status.json` to track progress
   - Display in CI output

### Legacy Intents File

Create `.legacy-intents`:
```
# Legacy intents grandfathered until migration
# Format: intent_name|deadline|status
fact_management|2026-02-01|pending
configure_system|2026-02-01|pending
burst_management|2026-02-15|pending
capture_event|2026-02-15|pending
manage_skills|2026-03-01|pending
generate_cv|2026-03-01|pending
```

---

## Progress Tracking

### Migration Status

| Intent | Lines | Priority | Assigned | Status | Completed |
|--------|-------|----------|----------|--------|-----------|
| browse_timeline | ~600 | - | - | DONE | 2026-01-XX |
| fact_management | 610 | P1 | - | Pending | - |
| configure_system | 657 | P2 | - | Pending | - |
| burst_management | 1,488 | P4 | - | Pending | - |
| capture_event | 1,884 | P6 | - | Pending | - |
| manage_skills | 2,415 | P7 | - | Pending | - |
| generate_cv | 2,746 | P8 | - | Pending | - |

### Pattern Standardization Status

| Pattern | Files Affected | Status | Completed |
|---------|----------------|--------|-----------|
| Direct `huh` imports | 27 | Pending | - |
| Incorrect form state | 22 | Pending | - |
| TODOs in code | 16 | Pending | - |
| Modals in `components/` | 15 | Pending | - |
| `context.Background()` | 12 | Pending | - |
| Raw lipgloss | ~50 | Pending | - |

### Weekly Progress Updates

**Week of 2026-01-24**:
- [ ] Created Master Migration Plan
- [ ] Generator script: Not started
- [ ] Intents migrated: 0

---

## Resources

### Documentation
- [Intent Migration Guide](guides/INTENT_MIGRATION_TO_SUBDIRECTORY.md)
- [Migration Rules](guides/MIGRATION_RULES.md)
- [Old to New Architecture](guides/OLD_TO_NEW_ARCHITECTURE_MIGRATION.md)
- [Screen Extraction Guide](guides/SCREEN_EXTRACTION_GUIDE.md)
- [Helper Extraction Guide](guides/HELPER_EXTRACTION_GUIDE.md)

### Templates
- `examples/intent_subdirectory_template/`

### Tools
- `make check-intent-architecture`
- `make check-compliance`
- `make what-to-use NEED="keyword"`

### Reference Implementation
- `intents/browse_timeline/`
- `screens/timeline/`

---

**Last Updated**: 2026-01-24  
**Next Review**: Weekly  
**Owner**: AI Agents + Human Review
