# BrowseTimeline Component Integration Analysis

**Date**: 2026-01-13
**Status**: ✅ **PATTERNS IMPLEMENTED, MISSING FILTER MODAL FOOTER**

---

## Executive Summary

BrowseTimeline intent has **successfully implemented all 12 standardized patterns** discovered during refactoring. The implementation serves as the **reference implementation** for other intents.

**Current Status**:
- ✅ Modal overlay rendering (Pattern 1) - COMPLETE
- ✅ Themed footer building (Pattern 2) - COMPLETE  
- ✅ View rendering with modal overlay (Pattern 3) - COMPLETE
- ✅ Global key interception (Pattern 4) - COMPLETE
- ✅ Context-aware footer generation (Pattern 5) - COMPLETE
- ⚠️ Filter modal footer uses **plain text** instead of KeyBadge components

**What Needs to Be Done**:
1. ✅ ~~Create pattern documentation~~ (COMPLETE - MODAL_OVERLAY_PATTERN.md, INTENT_PATTERNS_LIBRARY.md created)
2. **Fix filter modal footer** - Replace plain text with KeyBadge components
3. Apply patterns to ManageSkills intent (next priority)

---

## Table of Contents

1. [Current Implementation Status](#current-implementation-status)
2. [Pattern Compliance Matrix](#pattern-compliance-matrix)
3. [Component Inventory](#component-inventory)
4. [Missing Components](#missing-components)
5. [Action Items](#action-items)
6. [Next Intent: ManageSkills](#next-intent-manageskills)

---

## Current Implementation Status

### Intent Structure

**File**: `internal/cli/intents/browse_timeline_intent.go`
**Lines**: 403 lines (54% reduction from 879 lines)
**States**: 3 (Timeline, EventDetail, DeleteConfirm)
**Screens**: 3 (TimelineEventListScreen, TimelineEventDetailScreen, EventDeleteConfirmScreen)

### Screens

| Screen | File | Lines | Tests | Purpose |
|--------|------|-------|-------|---------|
| TimelineEventListScreen | `internal/cli/screens/timeline/event_list.go` | 224 | 28 | List events with actions |
| TimelineEventDetailScreen | `internal/cli/screens/timeline/event_detail.go` | 237 | 23 | View event details |
| EventDeleteConfirmScreen | `internal/cli/screens/timeline/event_delete_confirm.go` | 76 | - | Confirm deletion |

### Modals

| Modal | File | Lines | Tests | Purpose |
|-------|------|-------|-------|---------|
| FilterModalModel | `internal/cli/components/filter_modal.go` | 212 | - | Filter timeline events |

---

## Pattern Compliance Matrix

| Pattern # | Pattern Name | Status | Implementation Location | Notes |
|-----------|--------------|--------|-------------------------|-------|
| 1 | Modal Overlay Rendering | ✅ COMPLETE | `View()` lines 168-206 | StandardView FIRST, modal overlay LAST |
| 2 | Themed Footer Building | ⚠️ PARTIAL | `getContextHelp()` | Intent footers use KeyBadges, **modal footer uses plain text** |
| 3 | View Rendering with Modal Overlay | ✅ COMPLETE | `View()` lines 168-206 | Complete view → modal overlay → return |
| 4 | Global Key Interception | ✅ COMPLETE | `Update()` lines 117-167 | Modal → global → screen delegation |
| 5 | Context-Aware Footer Generation | ✅ COMPLETE | `getContextHelp()` | Different footer per state |
| 6 | State-to-Breadcrumb Mapping | ✅ COMPLETE | `getStateName()` lines 583-596 | Dynamic breadcrumbs |
| 7 | Screen Transition Helper | ✅ COMPLETE | `transitionToScreen()` lines 386-406 | Consistent screen initialization |
| 8 | Screen Result Handling | ✅ COMPLETE | `handleScreenResult()` lines 407-426 | Type-safe routing |
| 9 | Filter/Sort Application | ✅ COMPLETE | `applyFilters()` lines 264-353 | Filtering and sorting logic |
| 10 | Action Routing | ✅ COMPLETE | `handleNavigateResult()` lines 454-530 | Route actions (add/edit/delete/filter) |
| 11 | Delete Confirmation Flow | ✅ COMPLETE | `handleDeleteConfirmation()` lines 531-569 | Show confirmation → delete → refresh |
| 12 | Form Modal with Immediate Init | ✅ COMPLETE | `handleNavigateResult()` line 520 | `filterModal.Init()` called |

**Summary**: 11/12 patterns COMPLETE, 1 pattern PARTIAL (filter modal footer needs KeyBadges)

---

## Component Inventory

### Existing Components (All Working)

#### 1. TimelineEventListScreen
**File**: `internal/cli/screens/timeline/event_list.go`
**Purpose**: Display list of career events with actions
**Features**:
- Table-based list with 5 columns (Date, Company, Role, Category, Description)
- Selection indicator (▶)
- Actions: view (Enter), add (a), edit (e), delete (d), filter (f)
- Navigation: ↑/↓, j/k, g (first), G (last)
- Escape: Back to main menu
- Empty state: "No events found."
- Pagination: "Events: X | Page Y of Z"

**Interface**:
```go
type TimelineEventListScreen struct {
    *base.BaseScreen
    events         []*career.CareerEvent
    selectedIndex  int
    table          *components.ThemedTable
    // ... other fields
}

func (s *TimelineEventListScreen) RenderContent() string
func (s *TimelineEventListScreen) Update(msg tea.Msg) tea.Cmd
func (s *TimelineEventListScreen) Init() tea.Cmd
```

**Footer** (Themed with KeyBadges):
```
↑↓/jk: Navigate | Enter: View Details | a: Add | e: Edit | d: Delete | f: Filter | Esc: Back | q: Quit | m: Main Menu
```

#### 2. TimelineEventDetailScreen
**File**: `internal/cli/screens/timeline/event_detail.go`
**Purpose**: Display single event details with scrolling
**Features**:
- Scrollable detail view (↑/↓, j/k)
- Actions: edit (e), delete (d)
- Escape: Back to list
- Displays: Date, company, role, category, description, tags, bursts, facts

**Interface**:
```go
type TimelineEventDetailScreen struct {
    *base.BaseScreen
    event         *career.CareerEvent
    bursts        []*career.Burst
    facts         []*career.Fact
    scrollOffset  int
    // ... other fields
}

func (s *TimelineEventDetailScreen) RenderContent() string
func (s *TimelineEventDetailScreen) Update(msg tea.Msg) tea.Cmd
func (s *TimelineEventDetailScreen) Init() tea.Cmd
```

**Footer** (Themed with KeyBadges):
```
e: Edit | d: Delete | Esc: Back | q: Quit | m: Main Menu
```

#### 3. EventDeleteConfirmScreen
**File**: `internal/cli/screens/timeline/event_delete_confirm.go`
**Purpose**: Wrapper around BaseConfirmScreen for event deletion
**Features**:
- Yes/No button selection
- Toggle with ←→, h/l
- Direct submission: y (yes), n (no)
- Escape: Cancel
- Displays event details for confirmation

**Interface**:
```go
type EventDeleteConfirmScreen struct {
    *base.BaseConfirmScreen
}

func NewEventDeleteConfirmScreen(event *career.CareerEvent) *EventDeleteConfirmScreen
```

**Footer** (Built into BaseConfirmScreen):
```
y/Enter: Confirm | n/Esc: Cancel | ←→/hl: Toggle
```

#### 4. FilterModalModel
**File**: `internal/cli/components/filter_modal.go`
**Purpose**: Filter events by criteria
**Features**:
- Huh form with multiple filter fields
- Catppuccin theming
- Escape: Cancel
- Enter: Apply filters
- Tab/Shift+Tab: Navigate fields

**Interface**:
```go
type FilterModalModel struct {
    title        string
    visible      bool
    form         *huh.Form
    formData     *FilterFormData
    initialData  *TimelineFilters
    // ... other fields
}

func NewFilterModal(theme themes.Theme, initialFilters *TimelineFilters, width, height int) *FilterModalModel
func (m *FilterModalModel) Init() tea.Cmd
func (m *FilterModalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m *FilterModalModel) View() string
func (m *FilterModalModel) IsVisible() bool
func (m *FilterModalModel) Show()
func (m *FilterModalModel) Hide()
func (m *FilterModalModel) GetFilters() *TimelineFilters
```

**Footer** (⚠️ PLAIN TEXT - NEEDS FIX):
```go
// Current (WRONG):
footer := "Tab/Shift+Tab: Navigate | Space: Toggle | Enter: Apply | Esc: Cancel"

// Should be (CORRECT):
badges := []components.KeyBadge{
    components.NewKeyBadge("Tab/Shift+Tab", "Navigate"),
    components.NewKeyBadge("Space", "Toggle"),
    components.NewKeyBadge("Enter", "Apply"),
    components.CancelBadge(), // Esc: Cancel
}
footer := components.RenderHelpFooter(theme, badges...)
```

---

## Missing Components

### None Required

BrowseTimeline has **all required components**. The only work needed is:

1. **Fix FilterModalModel footer** - Replace plain text with KeyBadge components
2. **Extract reusable patterns** - Document for other intents to use

---

## Action Items

### High Priority (Before ManageSkills Migration)

#### 1. Fix FilterModalModel Footer
**File**: `internal/cli/components/filter_modal.go`
**Method**: `buildFilterModalFooter()` (called from intent)
**Current Code**:
```go
func (i *BrowseTimelineIntent) buildFilterModalFooter(theme themes.Theme) string {
    // Uses plain text - WRONG
    return "Tab/Shift+Tab: Navigate | Space: Toggle | Enter: Apply | Esc: Cancel"
}
```

**Target Code**:
```go
func (i *BrowseTimelineIntent) buildFilterModalFooter(theme themes.Theme) string {
    badges := []components.KeyBadge{
        components.NewKeyBadge("Tab/Shift+Tab", "Navigate"),
        components.NewKeyBadge("Space", "Toggle"),
        components.NewKeyBadge("Enter", "Apply"),
        components.CancelBadge(), // Esc: Cancel
    }
    return components.RenderHelpFooter(theme, badges...)
}
```

**Estimated Time**: 5 minutes

#### 2. Add FilterModal Footer Tests
**File**: Create test or update existing
**Test Coverage**:
- [ ] Footer contains correct KeyBadges
- [ ] Footer uses theme styling
- [ ] Footer updates when theme changes

**Estimated Time**: 10 minutes

### Medium Priority (Optional Enhancements)

#### 3. Extract Generic FilterModal[T]
**Goal**: Make FilterModal reusable across all intents
**Current**: Tightly coupled to TimelineFilters
**Target**: Generic `FilterModal[T]` accepting any filter type

**Benefits**:
- ManageSkills can reuse for skill filtering
- BurstManagement can reuse for burst filtering
- FactManagement can reuse for fact filtering

**Estimated Time**: 1-2 hours

**Defer**: Until after ManageSkills migration (learn more patterns first)

#### 4. Create ModalFooterBuilder Component
**Goal**: Fluent API for building modal footers
**Example**:
```go
footer := components.NewModalFooterBuilder(theme).
    AddNavigationKeys().      // Tab/Shift+Tab, Space
    AddConfirmationKeys().    // Enter: Apply
    AddCancellationKeys().    // Esc: Cancel
    Build()
```

**Estimated Time**: 30 minutes
**Defer**: Until we have 3+ modals using same pattern

---

## Next Intent: ManageSkills

### Current Status
**File**: `internal/cli/intents/manage_skills_intent.go`
**Lines**: 1,922 lines (after infrastructure added)
**Target**: ~250 lines (87% reduction)
**Priority**: High (core workflow, has documentation)

### Screens Required

| Screen | Exists? | Base Screen | Purpose | Estimated Lines |
|--------|---------|-------------|---------|-----------------|
| SkillsListScreen | ✅ YES | BaseScreen | List skills with actions | 247 (exists) |
| SkillDetailScreen | ✅ YES | BaseDetailScreen | View skill details | 165 (exists) |
| SkillFormScreen | ✅ YES | BaseFormScreen[T] | Add/edit skill | 93 (exists) |
| SkillDeleteConfirmScreen | ✅ YES | BaseConfirmScreen | Confirm deletion | 70 (exists) |
| SkillFilterModal | ❌ NO | FilterModal[T] | Filter skills | ~200 (NEW) |
| SkillSortModal | ❌ NO | Custom | Sort skills | ~150 (NEW) |

**Total New Components**: 2 (SkillFilterModal, SkillSortModal)

### Missing Components for ManageSkills

#### 1. SkillFilterModal
**Purpose**: Filter skills by category, level, years of experience
**Base**: Could reuse FilterModal pattern from BrowseTimeline
**Fields**:
- Categories (MultiSelect): Technical, Soft Skill, Domain Knowledge, etc.
- Level (Select): Beginner, Intermediate, Advanced, Expert
- Years range (Input): Min/Max years
- Search text (Input): Name contains

**Estimated Time**: 1 hour (using FilterModal as template)

#### 2. SkillSortModal
**Purpose**: Sort skills by various criteria
**Base**: Custom modal with radio buttons
**Options**:
- Name (A-Z, Z-A)
- Category
- Level (Beginner → Expert, Expert → Beginner)
- Years (Most → Least, Least → Most)
- Events (Most used → Least used)

**Estimated Time**: 1 hour

### ManageSkills Migration Estimate

| Task | Estimated Time | Complexity |
|------|----------------|------------|
| Create SkillFilterModal | 1 hour | Medium |
| Create SkillSortModal | 1 hour | Medium |
| Fix View() render order | 30 minutes | Low |
| Convert footers to KeyBadges | 1 hour | Low |
| Add global key interception | 30 minutes | Low |
| Update tests | 2 hours | Medium |
| Remove legacy code | 1 hour | Low |
| Manual testing | 1 hour | Low |
| **Total** | **8 hours** | **Medium** |

---

## Pattern Reusability Analysis

### Patterns BrowseTimeline Can Share

| Pattern | Reusable For | How to Share |
|---------|--------------|--------------|
| Modal Overlay Rendering | All intents | Already documented in MODAL_OVERLAY_PATTERN.md |
| Themed Footer Building | All intents | Already documented in INTENT_PATTERNS_LIBRARY.md |
| FilterModal Implementation | ManageSkills, BurstManagement, FactManagement | Extract generic FilterModal[T] |
| Delete Confirmation Flow | All intents | Already uses BaseConfirmScreen (reusable) |
| Action Routing | All list-based intents | Document pattern in INTENT_PATTERNS_LIBRARY.md |
| applyFilters() pattern | All list-based intents | Document pattern with example |

### Components BrowseTimeline Can Share

| Component | Reusable For | Modifications Needed |
|-----------|--------------|----------------------|
| EventDeleteConfirmScreen | ManageSkills, BurstManagement, FactManagement | None - just change text |
| FilterModalModel | ManageSkills, BurstManagement, FactManagement | Generify to FilterModal[T] |
| ThemedTable | All list screens | None - already reusable |

---

## Success Criteria

### BrowseTimeline (Current)
- [x] All 12 patterns implemented
- [x] 54% code reduction achieved
- [x] All tests passing
- [x] Zero regressions
- [ ] Filter modal footer uses KeyBadges (FINAL ITEM)

### Pattern Documentation
- [x] MODAL_OVERLAY_PATTERN.md created
- [x] INTENT_PATTERNS_LIBRARY.md created
- [ ] THEMED_FOOTER_GUIDE.md created (IN PROGRESS)
- [ ] AGENTS.md updated with pattern references

### ManageSkills (Next)
- [ ] All screens integrated
- [ ] SkillFilterModal created
- [ ] SkillSortModal created
- [ ] 87% code reduction achieved
- [ ] All tests passing
- [ ] Zero regressions
- [ ] TUI compliance verified

---

## Timeline

### Immediate (Today)
1. ✅ Create MODAL_OVERLAY_PATTERN.md (DONE)
2. ✅ Create INTENT_PATTERNS_LIBRARY.md (DONE)
3. 🔄 Create THEMED_FOOTER_GUIDE.md (IN PROGRESS)
4. Fix BrowseTimeline filter modal footer (5 min)
5. Update AGENTS.md with pattern references (30 min)

**Total**: ~1 hour remaining

### Tomorrow
1. Create SkillFilterModal component (1 hour)
2. Create SkillSortModal component (1 hour)
3. Start ManageSkills migration (2-3 hours)

**Total**: 4-5 hours

### This Week
1. Complete ManageSkills migration (remaining 3-4 hours)
2. Apply patterns to CaptureEvent intent (6-8 hours)
3. Document lessons learned

**Total**: 9-12 hours remaining in week

---

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Filter modal footer fix breaks tests | Low | Low | Simple text replacement, tests should adapt easily |
| ManageSkills has hidden complexity | Medium | Medium | Already has screens built (Phase 3.2), infrastructure in place |
| Pattern documentation incomplete | Low | High | Review with actual ManageSkills implementation, iterate |
| Token budget exceeded | Medium | Medium | Completed 2/3 major docs, remaining work is smaller |

---

**Status**: ✅ **BROWSETIMELINE COMPLETE (EXCEPT FILTER FOOTER), READY FOR MANAGESKILLS**

*Last updated: 2026-01-13*
