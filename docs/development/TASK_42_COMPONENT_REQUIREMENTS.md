# Task 42: Component Requirements & Integration Plan

**Date**: 2026-01-13
**Phase**: 4.2 (BrowseTimeline) → 4.1 (ManageSkills)
**Status**: BrowseTimeline complete, ManageSkills next

---

## Purpose

This document tracks **what components need to be created or modified** for each intent migration in Task 42. It ensures we don't miss required components and maintain consistency across intents.

---

## Component Creation Strategy

### When to Create New Components

✅ **Create new component when**:
1. **No suitable base screen exists** (e.g., async progress screen)
2. **Domain-specific behavior required** (e.g., SkillFilterModal with skill-specific fields)
3. **Reusability across multiple intents** (e.g., generic FilterModal[T])

❌ **Don't create new component when**:
1. **Base screen already handles it** (e.g., BaseSelectScreen for lists)
2. **Simple wrapper sufficient** (e.g., EventDeleteConfirmScreen wraps BaseConfirmScreen)
3. **Intent-specific logic only** (keep in intent, not component)

### Component Types

| Type | Purpose | Examples | Reusability |
|------|---------|----------|-------------|
| **Base Screen** | Generic UI pattern | BaseSelectScreen, BaseFormScreen, BaseConfirmScreen | High (all intents) |
| **Domain Screen** | Specific to domain entity | TimelineEventListScreen, SkillDetailScreen | Medium (one intent) |
| **Modal** | Overlay on top of screens | FilterModalModel, EditModalModel | High (multiple intents) |
| **Helper** | UI utilities | ThemedTable, KeyBadge, StandardView | High (all UI) |

---

## Phase 4.2: BrowseTimeline ✅ COMPLETE

### Components Created (Phase 3.3)

| Component | Type | File | Lines | Status |
|-----------|------|------|-------|--------|
| TimelineEventListScreen | Domain Screen | `internal/cli/screens/timeline/event_list.go` | 224 | ✅ Complete |
| TimelineEventDetailScreen | Domain Screen | `internal/cli/screens/timeline/event_detail.go` | 237 | ✅ Complete |
| EventDeleteConfirmScreen | Domain Screen | `internal/cli/screens/timeline/event_delete_confirm.go` | 76 | ✅ Complete |
| FilterModalModel | Modal | `internal/cli/components/filter_modal.go` | 212 | ⚠️ Footer needs KeyBadges |

### Components Reused from Existing

| Component | Type | Source | Used For |
|-----------|------|--------|----------|
| BaseScreen | Base Screen | `internal/cli/screens/base/base_screen.go` | All screens inherit |
| StandardView | Helper | `internal/cli/components/standard_view.go` | Layout with logo/breadcrumbs |
| ThemedTable | Helper | `internal/cli/components/themed_table.go` | Event list display |
| KeyBadge | Helper | `internal/cli/components/key_badge.go` | Themed footers |

### Outstanding Work

- [ ] **Fix FilterModalModel footer** - Replace plain text with KeyBadge components (5 min)
- [ ] **Add footer tests** - Verify KeyBadge usage (10 min)

---

## Phase 4.1: ManageSkills 🔄 IN PROGRESS

### Components Already Created (Phase 3.2)

| Component | Type | File | Lines | Status |
|-----------|------|------|-------|--------|
| SkillsListScreen | Domain Screen | `internal/cli/screens/skills/list.go` | 247 | ✅ Complete |
| SkillDetailScreen | Domain Screen | `internal/cli/screens/skills/detail.go` | 165 | ✅ Complete |
| SkillFormScreen | Domain Screen | `internal/cli/screens/skills/form.go` | 93 | ✅ Complete |
| SkillDeleteConfirmScreen | Domain Screen | `internal/cli/screens/skills/delete.go` | 70 | ✅ Complete |

### Components Needed (Phase 4.1)

| Component | Type | Base | Purpose | Est. Lines | Priority |
|-----------|------|------|---------|------------|----------|
| **SkillFilterModal** | Modal | FilterModal pattern | Filter by category/level/years | ~200 | HIGH |
| **SkillSortModal** | Modal | Custom | Sort by name/category/level/years | ~150 | MEDIUM |

### Component Details

#### SkillFilterModal
**Purpose**: Filter skills by category, level, years of experience, search text
**Base**: Reuse FilterModalModel pattern from BrowseTimeline
**Fields**:
- **Categories** (MultiSelect): Technical, Soft Skill, Domain Knowledge, Language, Tool, Framework, Other
- **Level** (Select): Beginner, Intermediate, Advanced, Expert
- **Years Range** (Inputs): Min years, Max years
- **Search Text** (Input): Name contains

**Interface**:
```go
type SkillFilterModal struct {
    title        string
    visible      bool
    form         *huh.Form
    formData     *SkillFilterFormData
    initialData  *SkillFilters
    theme        themes.Theme
    width        int
    height       int
}

type SkillFilterFormData struct {
    Categories      []string  // MultiSelect
    Level           string    // Select
    MinYears        string    // Input (validated as int)
    MaxYears        string    // Input (validated as int)
    SearchText      string    // Input
    SubmitConfirmed bool      // Confirm button
}

func NewSkillFilterModal(theme themes.Theme, initialFilters *SkillFilters, width, height int) *SkillFilterModal
func (m *SkillFilterModal) Init() tea.Cmd
func (m *SkillFilterModal) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m *SkillFilterModal) View() string
func (m *SkillFilterModal) IsVisible() bool
func (m *SkillFilterModal) GetFilters() *SkillFilters
```

**Footer** (Themed with KeyBadges):
```go
badges := []components.KeyBadge{
    components.NewKeyBadge("Tab/Shift+Tab", "Navigate"),
    components.NewKeyBadge("Space", "Toggle"),
    components.NewKeyBadge("Enter", "Apply"),
    components.CancelBadge(), // Esc: Cancel
}
return components.RenderHelpFooter(theme, badges...)
```

**Creation Time**: 1 hour

#### SkillSortModal
**Purpose**: Sort skills by name, category, level, years, or event usage
**Base**: Custom modal with radio buttons (huh.Form with Select field)
**Options**:
- Name (A→Z)
- Name (Z→A)
- Category
- Level (Beginner → Expert)
- Level (Expert → Beginner)
- Years (Most → Least)
- Years (Least → Most)
- Events (Most used → Least used)

**Interface**:
```go
type SkillSortModal struct {
    title        string
    visible      bool
    form         *huh.Form
    formData     *SkillSortFormData
    initialSort  string
    theme        themes.Theme
    width        int
    height       int
}

type SkillSortFormData struct {
    SortBy          string  // Select (radio buttons)
    SubmitConfirmed bool    // Confirm button
}

func NewSkillSortModal(theme themes.Theme, initialSort string, width, height int) *SkillSortModal
func (m *SkillSortModal) Init() tea.Cmd
func (m *SkillSortModal) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m *SkillSortModal) View() string
func (m *SkillSortModal) IsVisible() bool
func (m *SkillSortModal) GetSortBy() string
```

**Footer** (Themed with KeyBadges):
```go
badges := []components.KeyBadge{
    components.NewKeyBadge("↑↓/jk", "Navigate"),
    components.NewKeyBadge("Enter", "Apply"),
    components.CancelBadge(), // Esc: Cancel
}
return components.RenderHelpFooter(theme, badges...)
```

**Creation Time**: 1 hour

### Components Reused from Existing

| Component | Type | Source | Used For |
|-----------|------|--------|----------|
| BaseScreen | Base Screen | `internal/cli/screens/base/base_screen.go` | All screens inherit |
| BaseFormScreen[T] | Base Screen | `internal/cli/screens/base/form_screen.go` | SkillFormScreen |
| BaseConfirmScreen | Base Screen | `internal/cli/screens/base/confirm_screen.go` | SkillDeleteConfirmScreen |
| StandardView | Helper | `internal/cli/components/standard_view.go` | Layout with logo/breadcrumbs |
| ThemedTable | Helper | `internal/cli/components/themed_table.go` | Skill list display |
| KeyBadge | Helper | `internal/cli/components/key_badge.go` | Themed footers |

### Integration Checklist

#### Prerequisites
- [x] All 4 screens created (Phase 3.2)
- [x] Screen orchestration infrastructure added to intent
- [x] Intent tests passing (~97.5%)
- [ ] SkillFilterModal created
- [ ] SkillSortModal created

#### View() Method
- [ ] Implement Pattern 1 & 3: Modal Overlay Rendering
  - [ ] StandardView rendered FIRST
  - [ ] Modal overlaid LAST
  - [ ] Different render per screen type (list, detail, form, delete)
- [ ] Add `renderFilterModalOverlay()` helper
- [ ] Add `renderSortModalOverlay()` helper
- [ ] Remove direct screen.View() calls (use RenderContent())

#### Update() Method
- [ ] Implement Pattern 4: Global Key Interception
  - [ ] Priority 1: Modal updates (if modal visible)
  - [ ] Priority 2: Global keys (q, ?, m)
  - [ ] Priority 3: Screen delegation
- [ ] Handle filter modal results
- [ ] Handle sort modal results
- [ ] Remove direct keyboard handling from intent (let screens handle)

#### Footer Generation
- [ ] Implement Pattern 2 & 5: Themed Footer Building
  - [ ] Convert all plain text footers to KeyBadge components
  - [ ] Create `getContextHelp()` method (context-aware)
  - [ ] Different footer per screen (list, detail, form, delete)
- [ ] Add `getListScreenFooter()` helper
- [ ] Add `getDetailScreenFooter()` helper
- [ ] Add `getFormScreenFooter()` helper
- [ ] Add `getDeleteScreenFooter()` helper

#### State Management
- [ ] Implement Pattern 6: State-to-Breadcrumb Mapping
  - [ ] `GetState()` method returns dynamic breadcrumbs
  - [ ] Format: "Main Menu ▸ Manage Skills ▸ [State]"
- [ ] Implement Pattern 7: Screen Transition Helper
  - [ ] Use `transitionToScreen()` for all screen changes
  - [ ] Set terminal/theme/logo before Init()

#### Action Routing
- [ ] Implement Pattern 8 & 10: Screen Result Handling
  - [ ] Add `handleScreenResult()` router
  - [ ] Add `handleListAction()` for action routing
  - [ ] Route: add → form screen
  - [ ] Route: edit → form screen with data
  - [ ] Route: delete → confirm screen
  - [ ] Route: filter → filter modal
  - [ ] Route: sort → sort modal
  - [ ] Route: view → detail screen

#### Data Management
- [ ] Implement Pattern 9: Filter/Sort Application
  - [ ] Add `applyFilters()` method
  - [ ] Add `applySorting()` method
  - [ ] Refresh list after filter/sort changes
- [ ] Implement Pattern 11: Delete Confirmation Flow
  - [ ] Show confirmation screen with skill details
  - [ ] Delete on confirm
  - [ ] Refresh list after delete

#### Testing
- [ ] Update tests for new View() structure
- [ ] Add modal visibility tests
- [ ] Add global key interception tests
- [ ] Add footer rendering tests
- [ ] Visual test at multiple terminal sizes

#### Cleanup
- [ ] Remove legacy code (~1,400 lines)
- [ ] Verify line count reduced to ~250 lines
- [ ] Run `make check-compliance`
- [ ] Run `make generate-diagrams`

---

## Phase 4.3: CaptureEvent (Next After ManageSkills)

### Current Status
**Priority**: High (core workflow, has documentation)
**Workflow Doc**: `docs/workflows/EVENT_CAPTURE_WORKFLOW.md`
**Current Lines**: ~1,200 lines (4 states + 3 modals)
**Target Lines**: ~300 lines (75% reduction)

### Components Needed

| Component | Type | Base | Purpose | Est. Lines | Priority |
|-----------|------|------|---------|------------|----------|
| **EventCaptureStrategyScreen** | Domain Screen | BaseSelectScreen | Choose capture strategy (manual/burst) | ~100 | HIGH |
| **EventCaptureFormScreen** | Domain Screen | BaseFormScreen[T] | Capture event details | ~150 | HIGH |
| **EventReviewScreen** | Domain Screen | BaseDetailScreen | Review captured event | ~120 | HIGH |
| **MetadataEditModal** | Modal | Reuse existing | Edit event metadata | 0 (exists) | - |
| **BurstEditModal** | Modal | Reuse existing | Edit burst details | 0 (exists) | - |
| **FactEditModal** | Modal | Reuse existing | Edit fact details | 0 (exists) | - |

**Total New Components**: 3 (Strategy, Form, Review screens)
**Existing Modals**: 3 (Metadata, Burst, Fact - already exist)

### Notes
- CaptureEvent has **complex workflows** (burst extraction, fact extraction, metadata editing)
- Requires integration with **3 existing modals** (don't recreate)
- May benefit from **async progress screen** for burst/fact extraction
- **Defer until ManageSkills complete** to learn more patterns

---

## Component Reusability Matrix

### High Reusability (Use Everywhere)

| Component | Created In | Reusable For | How to Use |
|-----------|------------|--------------|------------|
| BaseScreen | Phase 1 | All screens | Inherit in domain screens |
| BaseSelectScreen[T] | Phase 1 | List screens | Inherit with domain type |
| BaseFormScreen[T] | Phase 3 | Form screens | Inherit with FormData type |
| BaseConfirmScreen | Phase 3 | Confirmation screens | Inherit or wrap |
| BaseDetailScreen | Phase 3 | Detail screens | Inherit for scrollable details |
| BaseProgressScreen | Phase 3 | Async operations | Use for loading/progress |
| StandardView | Always existed | All screens | Use for layout |
| ThemedTable | Always existed | List displays | Use in RenderContent() |
| KeyBadge | Always existed | Footers | Use in themed footers |

### Medium Reusability (Similar Intents)

| Component | Created In | Reusable For | Modifications Needed |
|-----------|------------|--------------|----------------------|
| FilterModalModel | Phase 4.2 (BrowseTimeline) | ManageSkills, BurstManagement, FactManagement | Generify to FilterModal[T] |
| EventDeleteConfirmScreen | Phase 4.2 | ManageSkills, BurstManagement, FactManagement | Change text, inherit from BaseConfirmScreen |
| SkillFilterModal | Phase 4.1 (ManageSkills) | None (skill-specific) | N/A |
| SkillSortModal | Phase 4.1 (ManageSkills) | Adapt for BurstManagement, FactManagement | Change sort options |

### Low Reusability (Intent-Specific)

| Component | Created In | Reusable For | Reason |
|-----------|------------|--------------|--------|
| TimelineEventListScreen | Phase 4.2 | None | Timeline-specific logic |
| TimelineEventDetailScreen | Phase 4.2 | None | Event-specific display |
| SkillsListScreen | Phase 4.1 | None | Skill-specific logic |
| SkillDetailScreen | Phase 4.1 | None | Skill-specific display |
| SkillFormScreen | Phase 4.1 | None | Skill-specific form fields |

---

## Generic Component Extraction Plan

### FilterModal[T] (High Priority)

**Goal**: Extract generic filter modal that works with any filter type
**Current**: FilterModalModel (timeline-specific)
**Target**: FilterModal[T] (generic)

**Interface**:
```go
type FilterConfig[T any] struct {
    Title          string
    BuildForm      func(data *T, width, height int) *huh.Form
    BuildFooter    func(theme themes.Theme) string
    GetInitialData func() *T
}

type FilterModal[T any] struct {
    config       FilterConfig[T]
    visible      bool
    form         *huh.Form
    formData     *T
    theme        themes.Theme
    width        int
    height       int
}

func NewFilterModal[T any](config FilterConfig[T], theme themes.Theme, width, height int) *FilterModal[T]
```

**When to Extract**: After ManageSkills complete (have 2 examples to learn from)

### SortModal (Medium Priority)

**Goal**: Generic sort modal for list-based intents
**Pattern**: Similar to FilterModal[T]

**When to Extract**: After BurstManagement complete (have 3 examples)

---

## Testing Strategy

### Component Tests (Unit)
- **What**: Individual component behavior
- **How**: Ginkgo/Gomega tests in `*_test.go` files
- **Coverage**: >85% per component

### Integration Tests
- **What**: Screen transitions within intent
- **How**: Mock service, test state machine
- **Coverage**: All state transitions

### E2E Tests
- **What**: Complete user workflows
- **How**: Simulate user input, verify output
- **Coverage**: Happy path + error cases

### Visual Tests
- **What**: Manual verification of appearance
- **How**: Run application, test at different terminal sizes
- **Coverage**: All screens, all themes

---

## Documentation Updates

### Per Intent Migration

- [ ] Update task document with completion status
- [ ] Update STATE_MATRIX.md (`make generate-diagrams`)
- [ ] Update workflow guide (if state machine changed)
- [ ] Add examples to pattern documentation (if novel patterns used)

### Phase 4 Complete

- [ ] Update AGENTS.md with complete pattern reference
- [ ] Create "Intent Migration Complete" section in task doc
- [ ] Update TUI_DEVELOPER_GUIDE.md with lessons learned
- [ ] Create "Component Library Reference" guide

---

## Risk Mitigation

### Component Creation Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Missing required component mid-migration | Low | Medium | This document + pre-migration analysis |
| Component doesn't fit base screen pattern | Low | High | Test with simple wrapper first, extract pattern if needed |
| Component not reusable as expected | Medium | Low | Start intent-specific, generify after 2-3 examples |
| Over-engineering components too early | Medium | Medium | Follow "Rule of Three" - generify after 3 uses |

### Integration Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Screen delegation breaks existing tests | High | Medium | Update tests incrementally, one screen at a time |
| Modal overlay causes visual glitches | Low | High | Follow MODAL_OVERLAY_PATTERN.md exactly |
| Footer theming inconsistent | Low | Medium | Use KeyBadge components everywhere, no plain text |
| State preservation lost | Medium | High | Test back navigation thoroughly, preserve metadata |

---

## Success Metrics

### Per Intent
- [ ] Line count reduced by 70%+ (target achieved)
- [ ] All tests passing (100% pass rate)
- [ ] Zero regressions in functionality
- [ ] Zero race conditions
- [ ] TUI compliance verified (`make check-compliance`)
- [ ] State matrix updated (`make generate-diagrams`)

### Overall Phase 4
- [ ] 10 intents migrated (1 complete, 9 remaining)
- [ ] ~7,200 lines removed (77% reduction)
- [ ] ~2,190 lines remaining (orchestration only)
- [ ] All patterns documented
- [ ] Component library established

---

## Timeline Estimate

### BrowseTimeline Completion
- Fix filter modal footer: 5 minutes
- Add footer tests: 10 minutes
- **Total**: 15 minutes

### ManageSkills Migration
- Create SkillFilterModal: 1 hour
- Create SkillSortModal: 1 hour
- Fix View() render order: 30 minutes
- Convert footers to KeyBadges: 1 hour
- Add global key interception: 30 minutes
- Update tests: 2 hours
- Remove legacy code: 1 hour
- Manual testing: 1 hour
- **Total**: 8 hours

### CaptureEvent Migration (Next)
- Create 3 screens: 3 hours
- Integration: 2 hours
- Testing: 2 hours
- **Total**: 7 hours

**Phase 4 Total Estimate**: ~50 hours (5-6 working days at 8 hours/day)

---

**Status**: ✅ **REQUIREMENTS DOCUMENTED, READY TO PROCEED**

*Last updated: 2026-01-13*
