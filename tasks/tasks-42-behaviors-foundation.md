# Task 42: Behaviors Foundation - Table & CRUD Components

## STATUS: IN PROGRESS

**Started**: 2026-01-16
**Prerequisites**: Phase 1 (UIKit Foundation) complete ✓
**Estimated Duration**: 3-4 days
**Approach**: Incremental migration, TDD Red-Green-Refactor

---

## Overview

- **Goal**: Create embeddable behavior components (`internal/cli/behaviors/`) that eliminate 600+ lines of duplicated table/CRUD code across 5 intents while maintaining type safety and composability
- **Time Estimate**: 3-4 days
- **Prerequisites**: 
  - Phase 1 (UIKit Foundation) complete ✓
  - All existing tests passing
  - Understanding of existing table patterns in BrowseTimeline, ManageSkills, BurstManagement, FactManagement intents

## Session Contract Acknowledgment

- [ ] Ran `make session-start` and it passed
- [ ] Acknowledge and commit to following all workflow rules
- [ ] Token count: 112,593 (started well above 50k - continuing from planning session)

---

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| **Architecture** | Embeddable behaviors (like BaseIntent) | Familiar pattern, composable, explicit dependencies |
| **Package Location** | `internal/cli/behaviors/` | Separate from uikit (different abstraction level) |
| **Composition Pattern** | Behaviors reference Table (one-way) | CRUD/Filter/Sort need Table data, not vice versa |
| **Type Safety** | Generics (`TableBehavior[T]`) | Type-safe item selection, no casting |
| **Migration** | Incremental per intent | Low risk, can test each migration independently |
| **State Ownership** | Table is source of truth | Behaviors delegate to table, never cache data |
| **Testing** | Unit tests + integration tests per behavior | Each behavior testable in isolation |

---

## Package Structure

```
internal/cli/
├── behaviors/                    # NEW PACKAGE
│   ├── types.go                  # Shared types (ColumnDef, CRUDMode, MenuOption, etc.)
│   ├── types_test.go
│   ├── table.go                  # TableBehavior[T] - core table functionality
│   ├── table_test.go
│   ├── crud.go                   # CRUDBehavior[T] - create/edit/delete handling
│   ├── crud_test.go
│   ├── filter_menu.go            # FilterMenuBehavior[T] - filter UI
│   ├── filter_menu_test.go
│   ├── sort_menu.go              # SortMenuBehavior[T] - sort UI
│   ├── sort_menu_test.go
│   └── doc.go                    # Package documentation
└── uikit/containers/             # Updated in parallel
    ├── box.go                    # For modal frames (Phase 3 prep)
    ├── box_test.go
    ├── overlay.go                # For modal centering (Phase 3 prep)
    ├── overlay_test.go
    └── doc.go
```

---

## Pre-Task Checklist (MUST COMPLETE BEFORE STARTING)

- [x] `make check-compliance` passes (tests pass, 82.76% coverage)
- [x] Reviewed existing table patterns in:
  - [x] `internal/cli/intents/browse_timeline_intent.go` (~150 lines table code)
  - [x] `internal/cli/intents/manage_skills_intent.go` (~300 lines table + CRUD + menus)
  - [x] `internal/cli/intents/burst_management_intent.go` (~150 lines table code)
  - [x] `internal/cli/intents/fact_management_intent.go` (~150 lines table code)
- [x] Reviewed existing components:
  - [x] `internal/cli/components/table_list_container.go` (315 lines - will be replaced)
  - [x] `internal/cli/navigation/list_navigator.go` (105 lines - interface remains)
- [x] Confirmed this is ONE atomic task per sub-phase
- [x] Identified test coverage targets (>85% per component)

## Current Progress

**Phase 2.5.1**: Shared Types - ✅ COMPLETE (16 tests passing)
**Phase 2.5.2**: TableBehavior[T] - ✅ COMPLETE (87 tests passing)
**Phase 2.5.3**: CRUDBehavior[T] - ✅ COMPLETE (32 tests passing, 119 total in suite)
**Phase 2.5.4**: FilterMenuBehavior[T] - ✅ COMPLETE (26 tests passing, 145 total in suite)
**Phase 2.5.5**: SortMenuBehavior[T] - ✅ COMPLETE (27 tests passing, 172 total in suite)

---

## Phase 2.5.1: Shared Types

### Files to Create
- [x] `internal/cli/behaviors/types.go`
- [x] `internal/cli/behaviors/types_test.go`
- [x] `internal/cli/behaviors/doc.go`

### TDD Checklist

#### RED Phase
- [x] Test file created: `internal/cli/behaviors/types_test.go`
- [x] Tests written for:
  - [x] `ColumnDef` validation (title required, width > 0)
  - [x] `CRUDMode` enum values (List, Create, Edit, Delete)
  - [x] `MenuOption` construction
  - [x] `MenuSection` construction
  - [x] `RowFormatter[T]` function signature (compile-time check)
  - [x] `FilterPredicate[T]` function signature
  - [x] `SortComparator[T]` function signature
- [x] Tests **FAIL** with error: `undefined: behaviors.ColumnDef` (and other types)
- [x] Test committed: `c99aa8d - test(cli): add failing tests for behaviors shared types`

#### GREEN Phase
- [x] `types.go` implemented with:
  ```go
  // ColumnDef defines table column configuration
  type ColumnDef struct {
      Title string
      Width int
  }
  
  // CRUDMode represents CRUD operation state
  type CRUDMode string
  const (
      ModeList   CRUDMode = "list"
      ModeCreate CRUDMode = "create"
      ModeEdit   CRUDMode = "edit"
      ModeDelete CRUDMode = "delete"
  )
  
  // MenuOption represents a single menu option
  type MenuOption struct {
      Label      string
      Value      interface{}
      IsSelected bool
      IsDisabled bool
  }
  
  // MenuSection groups related options
  type MenuSection struct {
      Title   string
      Options []MenuOption
  }
  
  // RowFormatter transforms item into table row cells
  type RowFormatter[T any] func(item T, index int) []string
  
  // FilterPredicate returns true if item passes filter
  type FilterPredicate[T any] func(item T) bool
  
  // SortComparator returns <0 if a<b, 0 if a==b, >0 if a>b
  type SortComparator[T any] func(a, b T) int
  ```
- [x] Tests now **PASS** (16/16 specs)
- [x] Implementation committed: `3004f81 - feat(cli): implement behaviors shared types`

#### REFACTOR Phase
- [x] Add doc comments for all exported types
- [x] Package documentation added (`doc.go`)
- [x] Refactored variable declarations per staticcheck warnings
- [x] Refactoring committed: `1077c24 - refactor(cli): merge variable declarations in behaviors tests`

**Acceptance Criteria**:
- [x] All types are exported and documented
- [x] Generic function types compile correctly
- [x] Tests pass with race detector
- [x] Coverage ≥ 85% (100% for types.go)

**Actual LOC**: ~60 source, ~200 test (exceeded estimates due to comprehensive function type tests)

---

## Phase 2.5.2: TableBehavior[T]

### Files to Create
- [x] `internal/cli/behaviors/table.go`
- [x] `internal/cli/behaviors/table_test.go`
- [x] `internal/cli/behaviors/behaviors_suite_test.go` (unified test suite)

### Design

`TableBehavior[T]` provides:
- **Data Binding**: Type-safe storage of items of type `T`
- **Pagination**: Automatic page calculation (no manual `updateTableRows()`)
- **Navigation**: Implements `ListNavigator` interface
- **Filtering**: Client-side filter predicates
- **Sorting**: Client-side sort comparators
- **Rendering**: Uses `table.Model` from bubbles, adds pagination info

### TDD Checklist

#### RED Phase
- [x] Test file created: `internal/cli/behaviors/table_test.go`
- [ ] Tests written for:
  - [ ] **Construction**:
    - [ ] `NewTableBehavior()` creates valid instance
    - [ ] Default page size is 15
    - [ ] Empty message defaults to "No items to display"
  - [ ] **Configuration (fluent API)**:
    - [ ] `PageSize()` sets page size and returns self
    - [ ] `EmptyMessage()` sets message and returns self
    - [ ] `Dimensions()` updates width/height
  - [ ] **Data Management**:
    - [ ] `SetItems()` replaces items and resets selection to 0
    - [ ] `GetItems()` returns current items
    - [ ] `GetSelectedItem()` returns pointer to selected item
    - [ ] `GetSelectedItem()` returns nil when empty
    - [ ] `GetSelectedIndex()` returns current index
    - [ ] `IsEmpty()` returns true when no items
    - [ ] `Count()` returns number of items
  - [ ] **ListNavigator Interface**:
    - [ ] `GetTotalItems()` returns item count
    - [ ] `SetSelectedIndex()` validates bounds (0 to len-1)
    - [ ] `SetSelectedIndex()` clamps negative to 0
    - [ ] `SetSelectedIndex()` clamps too-large to len-1
    - [ ] `GetPageSize()` returns configured page size
  - [ ] **Navigation**:
    - [ ] `HandleNavigation("up")` moves up by 1
    - [ ] `HandleNavigation("down")` moves down by 1
    - [ ] `HandleNavigation("k")` moves up (vim style)
    - [ ] `HandleNavigation("j")` moves down (vim style)
    - [ ] `HandleNavigation("pgup")` moves up by page size
    - [ ] `HandleNavigation("pgdn")` moves down by page size
    - [ ] `HandleNavigation("home")` goes to first
    - [ ] `HandleNavigation("end")` goes to last
    - [ ] `HandleNavigation("g")` goes to first (vim style)
    - [ ] `HandleNavigation("G")` goes to last (vim style)
    - [ ] Navigation doesn't go below 0
    - [ ] Navigation doesn't go above last item
    - [ ] Unknown keys return false
  - [ ] **Filtering**:
    - [ ] `SetFilter()` applies predicate
    - [ ] `SetFilter()` recalculates display items
    - [ ] `SetFilter(nil)` clears filter
    - [ ] `ClearFilter()` removes active filter
    - [ ] `HasFilter()` returns true when filter active
    - [ ] Filtered items preserve selection if possible
  - [ ] **Sorting**:
    - [ ] `SetSort()` applies comparator
    - [ ] `SetSort()` with reverse flag works
    - [ ] `SetSort(nil)` clears sort
    - [ ] `ClearSort()` removes active sort
    - [ ] `HasSort()` returns true when sort active
  - [ ] **Rendering**:
    - [ ] `Render()` returns table view when items exist
    - [ ] `Render()` shows empty message when no items
    - [ ] `Render()` includes pagination info
    - [ ] `RenderPaginationInfo()` formats "Items: X | Page Y of Z"
    - [ ] Pagination calculates correct page numbers
  - [ ] **Row Formatting**:
    - [ ] RowFormatter is called for each visible item
    - [ ] Selection indicator (▶) added to first column automatically
    - [ ] Non-selected rows have spacing (  ) in first column
- [x] Tests **FAIL**
- [x] Test committed:
  ```bash
  make ai-commit MSG="test(behaviors): add failing tests for TableBehavior"
  ```

#### GREEN Phase
- [x] `table.go` implemented with:
  ```go
  type TableBehavior[T any] struct {
      theme.Aware
      
      // Configuration
      columns      []ColumnDef
      rowFormatter RowFormatter[T]
      pageSize     int
      emptyMessage string
      paginationPrefix string
      showPagination bool
      
      // Data
      allItems      []T
      displayItems  []T  // After filter/sort
      selectedIndex int
      
      // Filter/Sort
      filterPredicate FilterPredicate[T]
      sortComparator  SortComparator[T]
      sortReverse     bool
      
      // Internal
      table         table.Model
      width, height int
      needsRefresh  bool
  }
  ```
- [x] Constructor creates bubbles table with columns
- [x] Fluent configuration methods implemented
- [x] Data methods implemented with bounds checking
- [x] ListNavigator interface implemented
- [x] Navigation methods with wrap prevention
- [x] Filter/sort logic with display item recalculation
- [x] Render logic with pagination calculation
- [x] Tests now **PASS**
- [x] Implementation committed:
  ```bash
  make ai-commit MSG="feat(behaviors): implement TableBehavior with pagination and navigation"
  ```

#### REFACTOR Phase
- [x] Extract pagination calculation to helper method
- [x] Extract row generation to helper method
- [x] Add inline documentation
- [x] Ensure no code duplication
- [x] Refactoring committed (if changes made)

**Acceptance Criteria**:
- [x] All navigation keys work correctly
- [x] Pagination calculates correctly for all item counts
- [x] Filter and sort work independently and together
- [x] Selection is preserved when possible during filter/sort
- [x] No panics on empty list
- [x] Tests pass with race detector
- [x] Coverage ≥ 85%

**Estimated LOC**: ~150 source, ~100 test

**Actual LOC**: ~224 source, ~338 test (27 comprehensive specs)

**Key Implementation Notes**:
- Simple flat list of options (no sections - simpler than FilterMenu)
- Each option stores: label, comparator function, reverse flag
- selectedIndex tracks currently applied sort (-1 if none initially)
- Navigation with wrapping (up/down, j/k)
- Applies sort directly to table via `table.SetSort(comparator, reverse)`
- Themed rendering with Catppuccin Macchiato colors
- Focus indicator (▶) and selection checkmark (✓)
- Show/Hide state management (only renders when active)
- OnApply callback invoked after applying sort
- All 172 tests passing (100% pass rate)
- Zero staticcheck warnings

**Design Simplicity**: SortMenu is simpler than FilterMenu because:
- No sections (just a flat list)
- Single selected option at a time
- Direct comparator application (no predicate building)

**Commits**:
- `6aad41d` - RED: Add failing tests for SortMenuBehavior
- `f2d71a7` - GREEN: Implement SortMenuBehavior with comparators

---

## Phase 2.5.6: Box Container (Phase 3 Prep)

### Files to Create
- [ ] `internal/cli/uikit/containers/box.go`
- [ ] `internal/cli/uikit/containers/box_test.go`

### Design

`Box` provides bordered container for modal frames, cards, panels.

### TDD Checklist

#### RED Phase
- [ ] Test file created: `internal/cli/uikit/containers/box_test.go`
- [ ] Tests written for:
  - [ ] `NewBox()` creates box with theme
  - [ ] `Content()` sets content
  - [ ] `Title()` sets title
  - [ ] `Variant()` sets variant (Default, Emphasized, Destructive, Subtle)
  - [ ] `Width()` sets width (0 = auto)
  - [ ] `Height()` sets height (0 = auto)
  - [ ] `Padding()` sets padding
  - [ ] `WithShadow()` enables shadow
  - [ ] `Render()` produces bordered box
  - [ ] Variants use correct theme colors
  - [ ] Emphasized uses thick border
  - [ ] Destructive uses error color
  - [ ] Shadow is rendered when enabled
- [ ] Tests **FAIL**
- [ ] Test committed

#### GREEN Phase
- [ ] `box.go` implemented with lipgloss styling
- [ ] Variants implemented with theme colors
- [ ] Border styles applied per variant
- [ ] Tests **PASS**
- [ ] Implementation committed

#### REFACTOR Phase
- [ ] Extract border style selection
- [ ] Add inline documentation

**Acceptance Criteria**:
- [ ] All variants render correctly
- [ ] Width/height constraints work
- [ ] Shadow rendering works
- [ ] Theme integration works
- [ ] Coverage ≥ 85%

**Estimated LOC**: ~150 source, ~100 test

---

## Phase 2.5.7: Overlay Container (Phase 3 Prep)

### Files to Create
- [ ] `internal/cli/uikit/containers/overlay.go`
- [ ] `internal/cli/uikit/containers/overlay_test.go`

### Design

`Overlay` provides centered modal overlay with optional background dimming.

### TDD Checklist

#### RED Phase
- [ ] Test file created: `internal/cli/uikit/containers/overlay_test.go`
- [ ] Tests written for:
  - [ ] `NewOverlay()` creates overlay with dimensions
  - [ ] `Content()` sets content to center
  - [ ] `Dimmed()` enables background dimming
  - [ ] `DimmedWith()` sets custom dim character
  - [ ] `Render()` centers content horizontally
  - [ ] `Render()` centers content vertically
  - [ ] `Render()` dims background when enabled
  - [ ] Dimensions are respected
- [ ] Tests **FAIL**
- [ ] Test committed

#### GREEN Phase
- [ ] `overlay.go` implemented with lipgloss.Place
- [ ] Dimming implemented with background fill
- [ ] Centering implemented with lipgloss positioning
- [ ] Tests **PASS**
- [ ] Implementation committed

#### REFACTOR Phase
- [ ] Add inline documentation

**Acceptance Criteria**:
- [ ] Content is centered in terminal
- [ ] Background dimming works
- [ ] Dimensions are correct
- [ ] Coverage ≥ 85%

**Estimated LOC**: ~100 source, ~80 test

---

## Post-Task Checklist (MUST COMPLETE BEFORE NEXT TASK)

- [ ] `make check-compliance` passes
- [ ] All behaviors tests passing (target: ~165 tests)
- [ ] All container tests passing (target: ~20 tests)
- [ ] Total coverage ≥ 85%
- [ ] Zero staticcheck warnings
- [ ] Zero race conditions
- [ ] All checkboxes above completed
- [ ] Task marked complete `[x]` in this file
- [ ] Token count: _____ (< 100k to continue)

---

## Acceptance Criteria

### Code Quality
- [ ] All tests pass (target: ~185 total)
- [ ] Coverage ≥ 85% for behaviors package
- [ ] Zero staticcheck warnings
- [ ] Zero race conditions
- [ ] All exported types documented

### Behaviors
- [ ] `TableBehavior[T]` handles data, pagination, navigation, filter, sort
- [ ] `CRUDBehavior[T]` handles create/edit/delete with confirmation
- [ ] `FilterMenuBehavior[T]` provides sectioned filter menu
- [ ] `SortMenuBehavior[T]` provides sort options menu
- [ ] All behaviors embeddable in intents
- [ ] One-way references (behaviors → table)

### Containers (Phase 3 Prep)
- [ ] `Box` provides themed bordered containers
- [ ] `Overlay` provides centered modal placement

### API Design
- [ ] All behaviors use embeddable pattern (like BaseIntent)
- [ ] Fluent configuration APIs
- [ ] Type-safe generics throughout
- [ ] Clear documentation

---

## Migration Plan (Phase 2.6 - Separate Task)

After Phase 2.5 completes, a follow-up task will migrate existing intents:

1. **BrowseTimeline** (low risk, straightforward table)
2. **FactManagement** (low risk, similar to timeline)
3. **BurstManagement** (medium risk, adds filtering)
4. **ManageSkills** (high complexity, full CRUD + filter + sort)

Each migration will:
- Keep old code commented for comparison
- Add tests to verify behavioral equivalence
- Be a separate atomic commit

---

## Rollback Plan

If issues arise:
1. `behaviors/` package is new - can be deleted without affecting existing code
2. `uikit/containers/` is new - safe to delete
3. No existing code depends on these packages yet
4. Migration happens in Phase 2.6 (separate task)

---

## Estimated Totals

| Phase | Component | Est. LOC | Est. Tests |
|-------|-----------|----------|------------|
| 2.5.1 | types.go | 80 | 10 |
| 2.5.2 | table.go | 350 | 45 |
| 2.5.3 | crud.go | 200 | 30 |
| 2.5.4 | filter_menu.go | 180 | 25 |
| 2.5.5 | sort_menu.go | 150 | 20 |
| 2.5.6 | box.go | 150 | 20 |
| 2.5.7 | overlay.go | 100 | 15 |
| **Total** | | **~1,210** | **~165** |

---

## Benefits Summary

| Benefit | Impact |
|---------|--------|
| **Code Reduction** | ~600 lines saved across 5 intents |
| **Type Safety** | `GetSelectedItem()` returns `*T` not `interface{}` |
| **Consistency** | All tables behave identically |
| **Testability** | Behaviors tested in isolation |
| **Composability** | Mix and match behaviors per intent |
| **Maintainability** | Fix bug once, affects all usages |
