# Screen and Modal Extraction Guide

**Purpose**: Define patterns for extracting views from legacy intents based on `browse_timeline` reference  
**Created**: 2026-01-24  
**Reference**: `intents/browse_timeline/` and `screens/timeline/`

---

## Core Pattern: Screens vs Modals

### From browse_timeline Reference

The reference implementation demonstrates clear separation:

| Type | Purpose | When to Use | Example |
|------|---------|-------------|---------|
| **Screen** | Full-page view with state transition | Major workflow step | `event_list.go`, `event_detail.go` |
| **Modal** | Overlay on current screen | Quick action, no state change | `filter_modal.go`, `search_modal.go` |

### Decision Tree

```
Is this view a major workflow step that occupies the full terminal?
  ├─ YES → SCREEN (in screens/{feature}/*.go)
  │         Examples: List, Detail, Form, Review, Results
  │
  └─ NO → Is it an overlay that returns to the current view?
           ├─ YES → MODAL (in screens/{feature}/modals/*.go)
           │         Examples: Filter, Sort, Search, Quick Add, Confirm
           │
           └─ NO → Is it a generic pattern?
                    ├─ YES → UIKIT MODAL (feedback.ConfirmModal, etc.)
                    │         Examples: Delete confirm, Error, Success
                    │
                    └─ NO → SCREEN (probably needs its own state)
```

### browse_timeline Structure

```
intents/browse_timeline/
├── types.go           # Intent struct with screen/modal fields
│   ├── activeScreen   # Current screen (generic)
│   ├── filterModal    # Feature-specific modal
│   ├── searchModal    # Feature-specific modal
│   ├── sortModal      # Feature-specific modal
│   ├── quickAddModal  # Feature-specific modal
│   ├── editModal      # Feature-specific modal
│   ├── viewDetailModal# Feature-specific modal
│   ├── deleteModal    # feedback.ConfirmModal (generic)
│   └── errorModal     # feedback.Modal (generic)

screens/timeline/
├── event_list.go      # SCREEN: List view (TableBehavior)
├── event_detail.go    # SCREEN: Detail view
├── event_delete_confirm.go  # SCREEN: Legacy (should be modal)
└── modals/
    ├── filter_modal.go    # MODAL: Multi-select filters
    ├── search_modal.go    # MODAL: Text search
    ├── sort_modal.go      # MODAL: Sort options
    ├── quick_add_modal.go # MODAL: Quick event creation
    ├── edit_modal.go      # MODAL: Event editing form
    ├── detail_modal.go    # MODAL: Event detail viewer
    └── helpers.go         # Shared modal helpers
```

---

## Legacy Intent Extraction Plans

### 1. fact_management (610 lines)

**Current States**:
```go
FactListState          = "list"           // List all facts
FactViewState          = "view"           // View single fact
FactEditorState        = "editor"         // Edit/create fact
FactDeleteConfirmState = "delete_confirm" // Confirm deletion
FactResultsState       = "results"        // Show operation results
FactCompletedState     = "completed"      // Intent complete
```

**Extraction Plan**:

| State | Type | Target File | Notes |
|-------|------|-------------|-------|
| `list` | SCREEN | `screens/facts/list_screen.go` | TableBehavior for facts |
| `view` | MODAL | `screens/facts/modals/detail_modal.go` | Overlay on list |
| `editor` | MODAL | `screens/facts/modals/edit_modal.go` | Form overlay |
| `delete_confirm` | MODAL | `feedback.ConfirmModal` | Use UIKit |
| `results` | MODAL | `feedback.SuccessModal` | Use UIKit |
| `completed` | N/A | Intent completion | No view needed |

**Required Files**:
```
screens/facts/
├── list_screen.go           # NEW: Main list view
└── modals/
    ├── detail_modal.go      # NEW: Fact detail viewer
    └── edit_modal.go        # NEW: Fact editor form
```

**Screens**: 1  
**Modals**: 2 custom + 2 UIKit  
**Estimated Reduction**: 610 → ~300 lines (50%)

---

### 2. configure_system (657 lines + 1,211 context)

**Current States**:
```go
ConfigStateSelectDomain  = "select_domain"   // Choose config domain
ConfigStateEditSettings  = "edit_settings"   // Edit settings form
ConfigStateReviewChanges = "review_changes"  // Review changes
ConfigStateConfirm       = "confirm"         // Confirm changes
ConfigStateSaving        = "saving"          // Saving state
ConfigStateComplete      = "complete"        // Success
ConfigStateFailed        = "failed"          // Error
```

**Existing Screens**: `screens/configure/` has 9 files

| State | Type | Existing File | Action |
|-------|------|---------------|--------|
| `select_domain` | SCREEN | `domain_select.go` | Keep |
| `edit_settings` | SCREEN | `edit_settings.go` | Keep |
| `review_changes` | SCREEN | `review_changes.go` | Keep |
| `confirm` | SCREEN | `confirm_screen.go` | Keep |
| `saving` | MODAL | `feedback.LoadingModal` | Use UIKit |
| `complete` | SCREEN | `complete_screen.go` | Keep |
| `failed` | SCREEN | `failed_screen.go` | Keep |

**Modals to Migrate**:
- `confirm_modal.go` → `screens/configure/modals/confirm_modal.go`
- `edit_settings_modal.go` → `screens/configure/modals/edit_settings_modal.go`
- `review_changes_modal.go` → `screens/configure/modals/review_changes_modal.go`

**Required Changes**:
```
screens/configure/
├── domain_select.go        # EXISTS: Keep
├── edit_settings.go        # EXISTS: Keep
├── review_changes.go       # EXISTS: Keep
├── confirm_screen.go       # EXISTS: Keep
├── complete_screen.go      # EXISTS: Keep
├── failed_screen.go        # EXISTS: Keep
└── modals/
    ├── confirm_modal.go    # MOVE from configure/
    ├── edit_settings_modal.go # MOVE from configure/
    └── review_changes_modal.go # MOVE from configure/
```

**Screens**: 6 (existing)  
**Modals**: 3 custom (move to modals/) + UIKit  
**Estimated Reduction**: 657 → ~350 lines (47%)

> **Note on Complexity Tradeoffs**: While line counts decrease with extraction, the number of files increases. This improves maintainability and testability but requires clear documentation of screen/modal responsibilities. Always document the purpose of each extracted component in PR descriptions and code comments to help new contributors understand the workflow.

---

### 3. burst_management (1,488 lines)

**Current States**:
```go
BurstListState          = "list"           // List all bursts
BurstViewState          = "view"           // View single burst
BurstEditorState        = "editor"         // Edit/create burst
BurstDeleteConfirmState = "delete_confirm" // Confirm deletion
BurstSuggestState       = "suggest"        // AI suggestions
BurstCompletedState     = "completed"      // Intent complete
```

**Extraction Plan**:

| State | Type | Target File | Notes |
|-------|------|-------------|-------|
| `list` | SCREEN | `screens/bursts/list_screen.go` | TableBehavior |
| `view` | MODAL | `screens/bursts/modals/detail_modal.go` | Overlay on list |
| `editor` | MODAL | `screens/bursts/modals/edit_modal.go` | Form overlay |
| `delete_confirm` | MODAL | `feedback.ConfirmModal` | Use UIKit |
| `suggest` | SCREEN | `screens/bursts/suggest_screen.go` | Full-page suggestions |
| `completed` | N/A | Intent completion | No view needed |

**Required Files**:
```
screens/bursts/
├── list_screen.go           # NEW: Main list view
├── suggest_screen.go        # NEW: AI suggestions view
└── modals/
    ├── detail_modal.go      # NEW: Burst detail viewer
    └── edit_modal.go        # NEW: Burst editor form
```

**Screens**: 2  
**Modals**: 2 custom + 1 UIKit  
**Estimated Reduction**: 1,488 → ~400 lines (73%)

---

### 4. capture_event (1,884 lines)

**Current States**:
```go
CaptureStateChooseStrategy = "choose_strategy" // Select capture method
CaptureStateForm           = "form"            // Event capture form
CaptureStateReview         = "review"          // Review captured event
CaptureStateSubmit         = "submit"          // Submit confirmation
```

**Existing Screens**: `screens/capture/` has 4 files

| State | Type | Existing File | Action |
|-------|------|---------------|--------|
| `choose_strategy` | SCREEN | `strategy_select.go` | Keep |
| `form` | SCREEN | `event_form_screen.go` | Keep |
| `review` | SCREEN | `event_review_screen.go` | Keep |
| `submit` | SCREEN | `event_submit_screen.go` | Keep |

**Required Changes**:
```
screens/capture/
├── strategy_select.go      # EXISTS: Keep
├── event_form_screen.go    # EXISTS: Keep
├── event_review_screen.go  # EXISTS: Keep
├── event_submit_screen.go  # EXISTS: Keep
└── modals/                 # NEW directory
    └── .gitkeep           # No custom modals needed
```

**Screens**: 4 (existing)  
**Modals**: 0 custom, use UIKit for errors  
**Estimated Reduction**: 1,884 → ~400 lines (79%)

---

### 5. manage_skills (2,415 lines)

**Current States**:
```go
SkillsStateList              = "list"                // View all skills
SkillsStateDetail            = "detail"              // View skill details
SkillsStateDetailEvents      = "detail_events"       // View events for skill
SkillsStateDetailEventDetail = "detail_event_detail" // View event from skill
SkillsStateAdd               = "add"                 // Add new skill
SkillsStateEdit              = "edit"                // Edit skill
SkillsStateDelete            = "delete"              // Confirm deletion
SkillsStateFilter            = "filter"              // Filter menu
SkillsStateSort              = "sort"                // Sort menu
```

**Existing Screens**: `screens/skills/` has 4 files

| State | Type | Existing/Target File | Action |
|-------|------|----------------------|--------|
| `list` | SCREEN | `list.go` | EXISTS |
| `detail` | SCREEN | `detail.go` | EXISTS |
| `detail_events` | SCREEN | `events_screen.go` | NEW |
| `detail_event_detail` | MODAL | `modals/event_detail_modal.go` | NEW |
| `add` | MODAL | `modals/form_modal.go` | Combine with edit |
| `edit` | MODAL | `modals/form_modal.go` | Combine with add |
| `delete` | MODAL | `feedback.ConfirmModal` | Use UIKit |
| `filter` | MODAL | `modals/filter_modal.go` | NEW |
| `sort` | MODAL | `modals/sort_modal.go` | NEW |

**Required Changes**:
```
screens/skills/
├── list.go                 # EXISTS: Keep
├── detail.go               # EXISTS: Keep
├── form.go                 # EXISTS: Review for modal conversion
├── delete.go               # EXISTS: Review for UIKit conversion
├── events_screen.go        # NEW: Events for skill
└── modals/
    ├── form_modal.go       # NEW: Add/Edit form overlay
    ├── filter_modal.go     # NEW: Filter menu
    ├── sort_modal.go       # NEW: Sort menu
    └── event_detail_modal.go # NEW: Event detail overlay
```

**Screens**: 3 (2 existing + 1 new)  
**Modals**: 4 custom + 1 UIKit  
**Estimated Reduction**: 2,415 → ~450 lines (81%)

---

### 6. generate_cv (2,746 lines)

**Current States** (grouped by phase):

**Configuration Phase**:
```go
GenerateCVStateSelectProfile        = "select_profile"
GenerateCVStateSelectAudience       = "select_audience"
GenerateCVStateSelectTechnologyFocus = "select_technology_focus"
GenerateCVStateSelectTechnologies   = "select_technologies"
GenerateCVStateSelectFocusArea      = "select_focus_area"
GenerateCVStateSelectSkillsConfig   = "select_skills_config"
GenerateCVStateSelectLengthFormat   = "select_length_format"
```

**Generation Phase**:
```go
GenerateCVStateExtractingTechnologies = "extracting_technologies"
GenerateCVStateGenerating             = "generating"
```

**Review Phase**:
```go
GenerateCVStatePreview = "preview"
GenerateCVStateReview  = "review"
GenerateCVStateConfirm = "confirm"
```

**Export Phase**:
```go
GenerateCVStateExportSelectFormat   = "export_select_format"
GenerateCVStateExportSelectLocation = "export_select_location"
GenerateCVStateExporting            = "exporting"
GenerateCVStateExportComplete       = "export_complete"
```

**Existing Screens**: `screens/cv/` has 5 files

| State | Type | Existing/Target File | Action |
|-------|------|----------------------|--------|
| `select_profile` | SCREEN | `profile_select.go` | EXISTS |
| `select_audience` | SCREEN | `audience_select.go` | EXISTS |
| Config steps | MODAL | `modals/config_wizard_modal.go` | NEW (combine 5 states) |
| `extracting` | MODAL | `feedback.LoadingModal` | Use UIKit |
| `generating` | SCREEN | `generating.go` | EXISTS |
| `preview` | SCREEN | `preview.go` | EXISTS |
| `review` | SCREEN | `review.go` | EXISTS |
| `confirm` | MODAL | `feedback.ConfirmModal` | Use UIKit |
| Export states | MODAL | `modals/export_modal.go` | NEW (combine 4 states) |

**Key Insight**: Many configuration states can be combined into a multi-step wizard modal.

**Required Changes**:
```
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

**Screens**: 5 (existing)  
**Modals**: 2 custom wizards + UIKit  
**Estimated Reduction**: 2,746 → ~500 lines (82%)

---

## Modal Patterns from browse_timeline

### Feature-Specific Modal Template

```go
// screens/{feature}/modals/filter_modal.go
package modals

import (
    "github.com/baphled/kariya/internal/cli/forms"
    "github.com/baphled/kariya/internal/cli/themes"
    tea "github.com/charmbracelet/bubbletea"
)

// FilterModal provides filtering UI for {feature}.
type FilterModal struct {
    theme   themes.Theme
    form    forms.Form
    visible bool
    result  *FilterResult
}

// FilterResult holds the filter configuration.
type FilterResult struct {
    Categories []string
    Tags       []string
    DateFrom   string
    DateTo     string
}

// NewFilterModal creates a new filter modal.
func NewFilterModal(theme themes.Theme) *FilterModal {
    return &FilterModal{
        theme: theme,
    }
}

// Show displays the modal.
func (m *FilterModal) Show() {
    m.visible = true
    m.buildForm()
}

// Hide hides the modal.
func (m *FilterModal) Hide() {
    m.visible = false
}

// IsVisible returns whether modal is visible.
func (m *FilterModal) IsVisible() bool {
    return m.visible
}

// Update handles updates.
func (m *FilterModal) Update(msg tea.Msg) tea.Cmd {
    if !m.visible {
        return nil
    }
    
    // Handle escape
    if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
        m.Hide()
        return nil
    }
    
    // Check form completion
    if forms.IsCompleted(m.form) {
        m.extractResult()
        m.Hide()
        return nil
    }
    
    // Delegate to form
    var cmd tea.Cmd
    m.form, cmd = m.form.Update(msg)
    return cmd
}

// View renders the modal.
func (m *FilterModal) View() string {
    if !m.visible {
        return ""
    }
    return m.form.View()
}

// Result returns the filter result.
func (m *FilterModal) Result() *FilterResult {
    return m.result
}

// Private methods...
func (m *FilterModal) buildForm() { ... }
func (m *FilterModal) extractResult() { ... }
```

### Using feedback.ConfirmModal (UIKit)

```go
// In intent
import "github.com/baphled/kariya/internal/cli/uikit/feedback"

type Intent struct {
    deleteModal *feedback.ConfirmModal
}

func (i *Intent) openDeleteConfirm(item *Item) tea.Cmd {
    i.deleteModal = feedback.NewConfirmModal(
        "Delete Item",
        fmt.Sprintf("Delete '%s'? This cannot be undone.", item.Name),
    ).WithVariant(feedback.ConfirmDestructive)
    
    return i.deleteModal.Init()
}

func (i *Intent) Update(msg tea.Msg) tea.Cmd {
    if i.deleteModal != nil && i.deleteModal.IsVisible() {
        cmd, confirmed := i.deleteModal.Update(msg)
        if confirmed {
            return i.handleDeleteConfirmed()
        }
        return cmd
    }
    // ...
}

func (i *Intent) View() string {
    baseView := i.activeScreen.View()
    
    if i.deleteModal != nil && i.deleteModal.IsVisible() {
        return behaviors.RenderModalOverlay(i.deleteModal, baseView)
    }
    
    return baseView
}
```

---

## Intent Types.go Pattern

Each migrated intent should have explicit screen/modal fields:

```go
// types.go
package my_feature

type Intent struct {
    *intents.BaseIntent
    
    // Context
    context *Context
    
    // State machine
    state  State
    active bool
    result *intents.IntentResult[*Result]
    
    // --- Screen Orchestration ---
    
    // Screens (one per major state)
    listScreen   *screens.ListScreen
    detailScreen *screens.DetailScreen
    formScreen   *screens.FormScreen
    
    // Active screen pointer
    activeScreen screens.Screen
    
    // --- Modals ---
    
    // Feature-specific modals (from screens/{feature}/modals/)
    filterModal *modals.FilterModal
    sortModal   *modals.SortModal
    editModal   *modals.EditModal
    
    // UIKit modals (generic patterns)
    deleteModal  *feedback.ConfirmModal
    errorModal   *feedback.Modal
    successModal *feedback.SuccessModal
    loadingModal *feedback.LoadingModal
    
    // Modal registry for unified handling
    modalRegistry *intents.ModalRegistry
}
```

---

## Summary: Extraction by Intent

| Intent | Screens | Custom Modals | UIKit Modals | Est. Lines |
|--------|---------|---------------|--------------|------------|
| fact_management | 1 | 2 | 2 | ~300 |
| configure_system | 6 (existing) | 3 (move) | 1 | ~350 |
| burst_management | 2 | 2 | 1 | ~400 |
| capture_event | 4 (existing) | 0 | 1 | ~400 |
| manage_skills | 3 | 4 | 1 | ~450 |
| generate_cv | 5 (existing) | 2 (wizards) | 2 | ~500 |

**Total Estimated Reduction**: 9,800 lines → ~2,400 lines (75% reduction)

---

## Checklist for Each Migration

### Phase 1: Analyze States
- [ ] List all states in the legacy intent
- [ ] Classify each state as Screen or Modal
- [ ] Identify existing screens that can be reused
- [ ] Identify UIKit modals that can replace custom modals

### Phase 2: Create/Move Screens
- [ ] Create `screens/{feature}/` directory if missing
- [ ] Create `screens/{feature}/modals/` directory
- [ ] Create new screens for missing states
- [ ] Verify screens use TableBehavior, UIKit primitives

### Phase 3: Create/Move Modals
- [ ] Create custom modals in `screens/{feature}/modals/`
- [ ] Replace custom delete confirmation with `feedback.ConfirmModal`
- [ ] Replace custom error display with `feedback.Modal`
- [ ] Replace custom loading with `feedback.LoadingModal`

### Phase 4: Update Intent Types
- [ ] Add explicit screen fields to intent struct
- [ ] Add explicit modal fields (custom + UIKit)
- [ ] Add modalRegistry for unified handling
- [ ] Remove inline render methods

### Phase 5: Update Intent Logic
- [ ] Update Init() to create screens
- [ ] Update Update() to delegate to screens/modals
- [ ] Update View() to use modalRegistry.RenderOverlay()
- [ ] Implement ScreenResultHandler

---

**Last Updated**: 2026-01-24  
**Reference**: `intents/browse_timeline/` + `screens/timeline/`
