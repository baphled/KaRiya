# KaRiya TUI Architecture

**Version**: 2.0
**Last Updated**: 2026-01-20
**Status**: Proposed

> **UIKit Components**: For UI component usage, see [UIKIT_GUIDE.md](../UIKIT_GUIDE.md).
> Use `uikit/` components instead of `components/` (legacy).

---

## Overview

KaRiya's TUI follows an **Intent → Screen → Component** architecture that separates workflow orchestration from view rendering, enabling high reusability and low boilerplate.

### Architecture Layers

```
┌─────────────────────────────────────────────────────────────────────┐
│                              INTENT                                  │
│  • Owns workflow state machine                                       │
│  • Orchestrates screen transitions                                   │
│  • Handles service calls and business logic                          │
│  • Returns IntentResult[T]                                           │
└──────────────────────────────┬──────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────────┐
│                              SCREEN                                  │
│  • Renders a single view within an intent                            │
│  • Handles user input for that view                                  │
│  • Manages UI state (terminal, theme, modals)                        │
│  • Returns ScreenResult to intent                                    │
└──────────────────────────────┬──────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────────┐
│                            COMPONENT                                 │
│  • UI primitives (StandardView, Modal, KeyBadge, etc.)              │
│  • Stateless or minimal state                                        │
│  • Highly reusable across screens                                    │
└─────────────────────────────────────────────────────────────────────┘
```

### Layer Responsibilities

| Layer | Responsibility | State Owned | Reusability |
|-------|---------------|-------------|-------------|
| **Intent** | Workflow orchestration, service calls | Workflow state, domain data | Per-workflow |
| **Screen** | Single view rendering, input handling | UI state, selection indices | High (shared base + domain-specific) |
| **Component** | Visual primitives | Minimal/none | Very High |

---

## Directory Structure

```
internal/cli/
├── intents/                          # Workflow orchestration
│   ├── contract.go                   # Intent, IntentResult interfaces
│   ├── router.go                     # IntentRouter implementation
│   │
│   ├── capture_event.go              # CaptureEventIntent
│   ├── generate_cv.go                # GenerateCVIntent
│   ├── browse_timeline.go            # BrowseTimelineIntent
│   ├── manage_skills.go              # ManageSkillsIntent
│   ├── burst_management.go           # BurstManagementIntent
│   ├── fact_management.go            # FactManagementIntent
│   └── configure_system.go           # ConfigureSystemIntent
│
├── screens/                          # View layer (hybrid: base + domain)
│   ├── contract.go                   # Screen interface, ScreenResult types
│   │
│   ├── base/                         # Reusable screen primitives
│   │   ├── base_screen.go            # BaseScreen with UI capabilities
│   │   ├── select_screen.go          # BaseSelectScreen[T]
│   │   ├── list_screen.go            # BaseListScreen[T]
│   │   ├── form_screen.go            # BaseFormScreen
│   │   ├── detail_screen.go          # BaseDetailScreen
│   │   ├── confirm_screen.go         # BaseConfirmScreen
│   │   └── progress_screen.go        # BaseProgressScreen
│   │
│   ├── capture/                      # CaptureEvent screens
│   │   ├── strategy_select.go        # CaptureStrategySelect
│   │   ├── event_form.go             # CaptureEventForm
│   │   └── event_review.go           # CaptureEventReview
│   │
│   ├── cv/                           # GenerateCV screens
│   │   ├── profile_select.go         # CVProfileSelect
│   │   ├── audience_select.go        # CVAudienceSelect
│   │   ├── role_emphasis_select.go   # CVRoleEmphasisSelect
│   │   ├── length_format_select.go   # CVLengthFormatSelect
│   │   ├── generating.go             # CVGenerating
│   │   └── preview.go                # CVPreview
│   │
│   ├── skills/                       # ManageSkills screens
│   │   ├── list.go                   # SkillsList (grouped by category)
│   │   ├── detail.go                 # SkillDetail
│   │   ├── events.go                 # SkillEvents (events using skill)
│   │   ├── form.go                   # SkillForm (add/edit)
│   │   ├── delete.go                 # SkillDeleteConfirm
│   │   ├── filter.go                 # SkillFilterMenu
│   │   └── sort.go                   # SkillSortMenu
│   │
│   ├── timeline/                     # BrowseTimeline screens
│   │   ├── event_list.go             # TimelineEventList
│   │   └── event_detail.go           # TimelineEventDetail
│   │
│   ├── facts/                        # FactManagement screens
│   │   ├── list.go                   # FactsList
│   │   └── detail.go                 # FactDetail
│   │
│   ├── bursts/                       # BurstManagement screens
│   │   ├── list.go                   # BurstsList
│   │   └── detail.go                 # BurstDetail
│   │
│   └── export/                       # ExportArtifact screens
│       ├── type_select.go            # ExportTypeSelect
│       ├── format_select.go          # ExportFormatSelect
│       └── complete.go               # ExportComplete
│
├── components/                       # UI primitives (existing)
│   ├── standard_view.go
│   ├── modal.go
│   ├── ascii_logo.go
│   ├── breadcrumb_bar.go
│   ├── key_badge.go
│   └── ...
│
├── forms/                            # Huh form configurations
│   ├── skill_form.go                 # SkillFormData, validators
│   ├── capture_event_form.go
│   └── ...
│
└── models/                           # DEPRECATED - phase out
    └── (migrate wrappers to screens/)
```

---

## Intent Design

### Intent Interface (Unchanged)

```go
type Intent interface {
    Init() tea.Cmd
    Update(msg tea.Msg) tea.Cmd
    View() string
    Result() *IntentResult[interface{}]
}
```

### Minimal Intent Structure

```go
type GenerateCVIntent struct {
    // Workflow state
    state   CVState
    context *CVContext
    result  *IntentResult[*CVResult]
    
    // Collected data during workflow
    selectedProfile  *cv.ProfileConfig
    selectedAudience string
    generatedCV      *career.CVView
    
    // Current active screen (owns all UI concerns)
    activeScreen screens.Screen
}
```

### Intent Responsibilities

**MUST:**
- Own workflow state machine
- Orchestrate screen transitions
- Call domain services
- Return IntentResult

**MUST NOT:**
- Render views directly (delegate to screens)
- Handle terminal/theme management (screens do this)
- Own UI state like scroll positions (screens do this)

---

## Screen Design

### Screen Interface

```go
type Screen interface {
    // Init initializes the screen
    Init() tea.Cmd
    
    // Update handles input, returns command and optional result
    Update(msg tea.Msg) (tea.Cmd, ScreenResult)
    
    // View renders the screen
    View() string
    
    // SetTerminalInfo updates terminal dimensions
    SetTerminalInfo(info *terminal.Info)
    
    // SetTheme sets the active theme
    SetTheme(theme themes.Theme)
}
```

### ScreenResult Types

```go
// ScreenResult communicates outcomes to the intent
type ScreenResult interface {
    isScreenResult()
}

// NavigateResult - user made a selection, ready to proceed
type NavigateResult struct {
    Data interface{}  // Selected data to pass forward
}

// CancelResult - user wants to go back
type CancelResult struct {
    Reason string  // "back", "escape", "cancel"
}

// SubmitResult - form was submitted
type SubmitResult struct {
    Data   interface{}
    Errors []ValidationError
}

// ErrorResult - an error occurred
type ErrorResult struct {
    Err error
}
```

---

## Real-World Examples

### Example 1: ManageSkillsIntent (Current Implementation)

**Current State** (1,647 lines in `manage_skills_intent.go`):

```go
type ManageSkillsIntent struct {
    *BaseIntent
    context *ManageSkillsContext
    state   SkillsState  // 9 states
    // ... 20+ fields for list, detail, form, filter, sort state
}

// 9 states, each with updateXxx() and viewXxx() methods
func (i *ManageSkillsIntent) Update(msg tea.Msg) tea.Cmd {
    switch i.state {
    case SkillsStateList:        return i.updateList(msg)
    case SkillsStateDetail:      return i.updateDetail(msg)
    case SkillsStateDetailEvents: return i.updateDetailEvents(msg)
    case SkillsStateAdd:         return i.updateAdd(msg)
    case SkillsStateEdit:        return i.updateEdit(msg)
    case SkillsStateDelete:      return i.updateDelete(msg)
    case SkillsStateFilter:      return i.updateFilter(msg)
    case SkillsStateSort:        return i.updateSort(msg)
    // ...
    }
}
```

**Proposed Refactor** (~250 lines):

```go
type ManageSkillsIntent struct {
    // Workflow state
    state        SkillsState
    context      *ManageSkillsContext
    result       *IntentResult[*ManageSkillsResult]
    
    // Collected data
    skills       []*career.Skill
    selectedSkill *career.Skill
    
    // Active screen (owns all UI)
    activeScreen screens.Screen
}

func (i *ManageSkillsIntent) transitionTo(state SkillsState) {
    i.state = state
    
    switch state {
    case SkillsList:
        i.activeScreen = skills.NewSkillsList(skills.SkillsListConfig{
            Skills:      i.skills,
            Breadcrumbs: []string{"Main Menu", "Manage Skills"},
        })
        
    case SkillsDetail:
        i.activeScreen = skills.NewSkillDetail(skills.SkillDetailConfig{
            Skill:       i.selectedSkill,
            Breadcrumbs: []string{"Main Menu", "Manage Skills", i.selectedSkill.Name},
        })
        
    case SkillsForm:
        i.activeScreen = skills.NewSkillForm(skills.SkillFormConfig{
            Skill:       i.selectedSkill, // nil for add
            Breadcrumbs: []string{"Main Menu", "Manage Skills", "Edit"},
        })
        
    case SkillsDelete:
        i.activeScreen = skills.NewSkillDeleteConfirm(skills.ConfirmConfig{
            Title:   "Delete Skill",
            Message: fmt.Sprintf("Delete %q?", i.selectedSkill.Name),
        })
    }
}

func (i *ManageSkillsIntent) Update(msg tea.Msg) tea.Cmd {
    cmd, result := i.activeScreen.Update(msg)
    if result != nil {
        return i.handleScreenResult(result)
    }
    return cmd
}
```

---

### Example 2: GenerateCVIntent with Variants

**Current Flow** (10+ states):
```
SelectProfile → SelectAudience → SelectRoleEmphasis → SelectLengthFormat → 
Generating → Preview → Confirm → Export
```

**Proposed with Screens** (~200 lines):

```go
type GenerateCVIntent struct {
    state   CVState
    context *GenerateCVContext
    result  *IntentResult[*GenerateCVResult]
    
    // Collected data through workflow
    selectedProfile      *cv.ProfileConfig
    selectedAudience     string
    selectedRoleEmphasis cv.RoleEmphasis
    selectedLengthFormat cv.LengthFormat
    selectedVariant      *cv.CVVariant
    generatedCV          *career.CVView
    
    activeScreen screens.Screen
}

func (i *GenerateCVIntent) transitionTo(state CVState) {
    i.state = state
    
    switch state {
    case CVSelectProfile:
        i.activeScreen = cv.NewCVProfileSelect(i.context.AvailableProfiles)
        
    case CVSelectAudience:
        i.activeScreen = cv.NewCVAudienceSelect([]string{
            "hiring_manager", "recruiter", "peer",
        })
        
    case CVSelectRoleEmphasis:
        emphases := cv.NewVariantService().ListRoleEmphases()
        i.activeScreen = cv.NewCVRoleEmphasisSelect(emphases)
        
    case CVSelectLengthFormat:
        lengths := cv.NewVariantService().ListLengthFormats()
        i.activeScreen = cv.NewCVLengthFormatSelect(lengths)
        
    case CVGenerating:
        i.activeScreen = cv.NewCVGenerating(cv.ProgressConfig{
            Title:   "Generating CV",
            Message: fmt.Sprintf("Creating %s CV...", i.selectedVariant.Name),
        })
        
    case CVPreview:
        i.activeScreen = cv.NewCVPreview(cv.CVPreviewConfig{
            CV:        i.generatedCV,
            Variant:   i.selectedVariant,
            Structure: i.selectedVariant.BaseStructure,
        })
    }
}

func (i *GenerateCVIntent) handleScreenResult(result screens.ScreenResult) tea.Cmd {
    switch r := result.(type) {
    case *screens.NavigateResult:
        switch i.state {
        case CVSelectProfile:
            i.selectedProfile = r.Data.(*cv.ProfileConfig)
            i.transitionTo(CVSelectAudience)
            
        case CVSelectAudience:
            i.selectedAudience = r.Data.(string)
            i.transitionTo(CVSelectRoleEmphasis)
            
        case CVSelectRoleEmphasis:
            i.selectedRoleEmphasis = r.Data.(cv.RoleEmphasis)
            i.transitionTo(CVSelectLengthFormat)
            
        case CVSelectLengthFormat:
            i.selectedLengthFormat = r.Data.(cv.LengthFormat)
            i.selectedVariant = i.lookupVariant()
            i.transitionTo(CVGenerating)
            return i.generateCVAsync()
        }
        
    case *screens.CancelResult:
        i.goBack()
    }
    return nil
}
```

---

## Benefits Summary

| Aspect | Before (Current) | After (Screens) |
|--------|------------------|-----------------|
| **ManageSkillsIntent** | 1,647 lines | ~250 lines |
| **GenerateCVIntent** | ~900 lines | ~200 lines |
| **Adding new screen** | Write update+view (~100-150 lines) | Configure screen (~20-30 lines) |
| **Testing** | Test whole intent | Test screens independently |
| **List navigation** | Duplicated per intent | Shared in BaseListScreen |
| **Form handling** | Per-intent wrappers | Shared in BaseFormScreen |
| **Code reduction** | - | 70%+ average |

---

## Data Flow

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              INTENT                                      │
│                                                                          │
│  1. User triggers intent from Main Menu                                  │
│  2. Intent creates first Screen with initial data                        │
│  3. Intent delegates Update/View to active Screen                        │
│                                                                          │
│     ┌─────────────────────────────────────────────────────────────┐     │
│     │                         SCREEN                               │     │
│     │                                                              │     │
│     │  4. Screen handles key input                                 │     │
│     │  5. Screen renders view using Components                     │     │
│     │  6. Screen returns ScreenResult on completion                │     │
│     │                                                              │     │
│     │     ┌─────────────────────────────────────────────────┐     │     │
│     │     │                  COMPONENTS                      │     │     │
│     │     │  StandardView, Modal, KeyBadge, etc.            │     │     │
│     │     └─────────────────────────────────────────────────┘     │     │
│     └─────────────────────────────────────────────────────────────┘     │
│                                                                          │
│  7. Intent receives ScreenResult                                         │
│  8. Intent transitions state and creates next Screen                     │
│  9. Repeat until workflow complete                                       │
│ 10. Intent returns IntentResult to Router                                │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Migration Strategy

See `tasks/tasks-21-tui-architecture-refactor.md` for detailed migration plan.

**Summary:**
1. **Phase 1**: Create screens/ foundation with base screens
2. **Phase 2**: Pilot with GenerateCVIntent
3. **Phase 3**: Build screen library
4. **Phase 4**: Migrate remaining intents
5. **Phase 5**: Deprecate models/ directory

**Timeline**: 6-8 weeks

---

## References

- **Naming Conventions**: `docs/rules/NAMING_CONVENTIONS.md`
- **Screen Specification**: `docs/architecture/SCREEN_SPECIFICATION.md`
- **Implementation Task**: `tasks/tasks-21-tui-architecture-refactor.md`
- **User Journeys**: `docs/PRD_MASTER.md` (Section 16)
