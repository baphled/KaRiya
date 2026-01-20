# Task 49: Fact Management UX Improvements

## Overview

- **Goal**: Enable users to manage and review facts with full visibility into role coverage, filtering by role and company, sorting, and the ability to re-classify facts (single or batch) to fix incorrect role assignments.
- **Time Estimate**: 8-10 hours
- **Prerequisites**: 
  - Understanding of TableBehavior[T] (internal/cli/behaviors/table.go)
  - Understanding of FilterBehavior interface pattern (internal/cli/intents/filter_behavior.go)
  - Understanding of burst_fact.Classifier (internal/service/career/burst_fact/classifier.go)

## Context

**Problem**: When generating CVs, users discover that companies are missing because facts have incorrect `role_fit` values. The current fact classifier only used keyword matching, ignoring source event categories. We fixed the classifier (`ClassifyRoleFitWithCategories`), but existing facts still have incorrect values.

**Solution**: Improve the Fact Management intent with:
1. Role + Company columns for visibility
2. Role coverage header showing distribution
3. Filter by Role AND Company
4. Sort by multiple fields
5. Single fact re-classify (view state)
6. Batch re-classify all facts (with confirm + progress + summary)
7. FIFO filter clearing

**Related Work**: 
- Classifier fix: `ClassifyRoleFitWithCategories()` in `internal/service/career/burst_fact/classifier.go`
- Pattern reference: `ManageSkillsIntent` filter/sort implementation

## Session Contract Acknowledgment

- [ ] Ran `make session-start` and it passed
- [ ] Acknowledge and commit to following all workflow rules
- [ ] Token count: _____ (must be < 50k to start)

## Pre-Task Checklist (MUST COMPLETE BEFORE STARTING)

- [ ] `make check-compliance` passes
- [ ] Reviewed existing patterns in:
  - [ ] `internal/cli/intents/fact_management.go` (context)
  - [ ] `internal/cli/intents/fact_management_intent.go` (intent)
  - [ ] `internal/cli/intents/manage_skills_intent.go` (FilterBehavior pattern)
  - [ ] `internal/cli/components/skill_filter_modal.go` (filter modal pattern)
  - [ ] `internal/cli/components/skill_sort_modal.go` (sort modal pattern)
  - [ ] `internal/cli/behaviors/table.go` (TableBehavior API)
- [ ] Confirmed this is ONE atomic task per phase
- [ ] Identified which test files will be created/modified

## Files to Create

- [ ] `internal/cli/components/fact_filter_modal.go`
- [ ] `internal/cli/components/fact_filter_modal_test.go`
- [ ] `internal/cli/components/fact_sort_modal.go`
- [ ] `internal/cli/components/fact_sort_modal_test.go`

## Files to Modify

- [ ] `internal/cli/intents/fact_management.go` (~80 LOC)
- [ ] `internal/cli/intents/fact_management_intent.go` (~350 LOC)
- [ ] `internal/cli/intents/fact_management_test.go` (new tests)
- [ ] `internal/cli/uikit/primitives/badge.go` (+5 LOC for SortBadge)

## Implementation Checklist

---

### Phase 1: Data Infrastructure (~80 LOC, ~12 tests)

**Goal**: Add event lookup capability and role coverage statistics to FactManagementContext.

#### Files to Modify
- `internal/cli/intents/fact_management.go`

#### Changes Required

1. **Add new fields to `FactManagementContext`**:
   ```go
   EventRepository     careerrepo.Repository
   EventMap            map[string]*domain.CareerEvent  // eventID -> event
   RoleCoverage        *RoleCoverageStats
   UniqueCompanies     []string
   ActiveRoleFilter    string  // "" = all, or "senior_ic", "staff", "em", "principal"
   ActiveCompanyFilter string  // "" = all, or company name
   ```

2. **Add `RoleCoverageStats` struct**:
   ```go
   type RoleCoverageStats struct {
       SeniorIC  int
       Staff     int
       EM        int
       Principal int
       Total     int
   }
   ```

3. **Implement new methods**:
   - `LoadEvents() error` - Populate EventMap from facts' SourceEventID
   - `CalculateRoleCoverage()` - Count facts per role
   - `GetCompanyForFact(fact *domain.Fact) string` - Lookup company via EventMap
   - `ExtractUniqueCompanies()` - Extract unique companies for filter modal

4. **Update constructor** to accept `EventRepository`:
   ```go
   func NewFactManagementContext(
       factRepo careerrepo.FactRepository,
       eventRepo careerrepo.Repository,  // NEW
       ctx context.Context,
   ) *FactManagementContext
   ```

#### TDD Checklist

##### RED Phase
- [ ] Test file created/modified: `internal/cli/intents/fact_management_test.go`
- [ ] Tests written and **FAIL**:
  - [ ] TestLoadEvents_PopulatesEventMap
  - [ ] TestLoadEvents_HandlesEmptySourceEventID
  - [ ] TestLoadEvents_HandlesNilEventRepository
  - [ ] TestCalculateRoleCoverage_CountsCorrectly
  - [ ] TestCalculateRoleCoverage_EmptyFacts
  - [ ] TestGetCompanyForFact_ReturnsCompany
  - [ ] TestGetCompanyForFact_ReturnsDashForUnknown
  - [ ] TestExtractUniqueCompanies_ExtractsUnique
  - [ ] TestExtractUniqueCompanies_SortsAlphabetically
  - [ ] TestNewFactManagementContext_WithEventRepo

##### GREEN Phase
- [ ] Minimal implementation written
- [ ] All tests pass

##### REFACTOR Phase
- [ ] Code refactored for clarity/DRY
- [ ] Tests still pass

#### Pre-Commit Checklist
- [ ] `make check-compliance` passes
- [ ] Use `make ai-commit MSG="feat(fact-management): add event lookup and role coverage infrastructure"`

---

### Phase 2: Add Role + Company Columns (~40 LOC, ~6 tests)

**Goal**: Update table to display Role and Company columns.

#### Files to Modify
- `internal/cli/intents/fact_management_intent.go`

#### Changes Required

1. **Update column definitions** (5 columns instead of 3):
   ```go
   columns := []behaviors.ColumnDef{
       {Title: "Fact", Width: 35},
       {Title: "Role", Width: 10},
       {Title: "Company", Width: 18},
       {Title: "Strength", Width: 12},
       {Title: "Categories", Width: 20},
   }
   ```

2. **Create `createRowFormatter()` closure method** that has access to EventMap

3. **Add `formatRoleFit()` helper**:
   ```go
   func formatRoleFit(role domain.RoleFit) string {
       switch role {
       case domain.RoleFitSeniorIC:  return "SrIC"
       case domain.RoleFitStaff:     return "Staff"
       case domain.RoleFitEM:        return "EM"
       case domain.RoleFitPrincipal: return "Princ"
       default:                      return "-"
       }
   }
   ```

#### TDD Checklist

##### RED Phase
- [ ] Tests written and **FAIL**:
  - [ ] TestFactRowFormatter_IncludesRoleColumn
  - [ ] TestFactRowFormatter_IncludesCompanyColumn
  - [ ] TestFormatRoleFit_AllRoles
  - [ ] TestFactTable_HasFiveColumns

##### GREEN Phase
- [ ] Implementation written
- [ ] All tests pass

#### Pre-Commit Checklist
- [ ] `make check-compliance` passes
- [ ] Use `make ai-commit MSG="feat(fact-management): add role and company columns to fact table"`

---

### Phase 3: Role Coverage Header (~50 LOC, ~5 tests)

**Goal**: Show role distribution above the fact table.

#### Files to Modify
- `internal/cli/intents/fact_management_intent.go`

#### Changes Required

1. **Implement `renderRoleCoverageHeader()` method**:
   - Format: "Roles: SrIC 265 │ Staff 23 │ EM 32 │ Princ 41"

2. **Implement `renderActiveFiltersIndicator()` method**:
   - Format: "🔍 Role: senior_ic │ Company: Acme │ x to clear"

3. **Update `viewList()`** to include header and indicator

#### TDD Checklist

##### RED Phase
- [ ] Tests written and **FAIL**:
  - [ ] TestRenderRoleCoverageHeader_ShowsAllRoles
  - [ ] TestRenderRoleCoverageHeader_NilStatsReturnsEmpty
  - [ ] TestRenderActiveFiltersIndicator_NoFilters
  - [ ] TestRenderActiveFiltersIndicator_RoleFilter
  - [ ] TestRenderActiveFiltersIndicator_BothFilters

##### GREEN Phase
- [ ] Implementation written
- [ ] All tests pass

#### Pre-Commit Checklist
- [ ] `make check-compliance` passes
- [ ] Use `make ai-commit MSG="feat(fact-management): add role coverage header and filter indicator"`

---

### Phase 4: Filter Modal (~220 LOC, ~18 tests)

**Goal**: Create filter modal for filtering by Role and Company.

#### Files to Create
- `internal/cli/components/fact_filter_modal.go`
- `internal/cli/components/fact_filter_modal_test.go`

#### Changes Required

1. **Define structs**:
   - `FactFilters` (RoleFit string, Company string)
   - `FactFilterFormData` (form field values)
   - `FactFilterModal` (form, theme, counts)

2. **Implement modal methods**:
   - `NewFactFilterModal()` constructor
   - `buildForm()` with role + company selects (showing counts)
   - `Init()`, `Update()`, `View()` methods
   - `ToFactFilters()` conversion
   - `RenderOverlay()` method

3. **Integrate into intent**:
   - Add `filterModal` field to `FactManagementModel`
   - Add `"f"` key handler in `handleListState()`
   - Implement `openFilterModal()`, `handleFilterModalUpdate()`, `applyFilters()`
   - Update `View()` to render filter modal overlay

#### TDD Checklist

##### RED Phase
- [ ] Tests written and **FAIL**:
  - [ ] TestNewFactFilterModal_CreatesModal
  - [ ] TestFactFilterModal_Init
  - [ ] TestFactFilterModal_BuildForm_RoleOptions
  - [ ] TestFactFilterModal_BuildForm_CompanyOptions
  - [ ] TestFactFilterModal_BuildForm_ShowsCounts
  - [ ] TestFactFilterModal_Update_AppliesOnEnter
  - [ ] TestFactFilterModal_Update_CancelsOnEsc
  - [ ] TestFactFilterModal_ToFactFilters
  - [ ] TestFactFilterModal_RenderOverlay
  - [ ] TestFactFilterModal_PrePopulatesCurrentFilters
  - [ ] TestFactManagement_FKeyOpensFilterModal
  - [ ] TestFactManagement_ApplyFilters_RoleOnly
  - [ ] TestFactManagement_ApplyFilters_CompanyOnly
  - [ ] TestFactManagement_ApplyFilters_Both
  - [ ] TestFactManagement_ApplyFilters_ClearsOnEmpty

##### GREEN Phase
- [ ] Implementation written
- [ ] All tests pass

#### Pre-Commit Checklist
- [ ] `make check-compliance` passes
- [ ] Use `make ai-commit MSG="feat(fact-management): add filter modal for role and company filtering"`

---

### Phase 5: Sort Modal (~160 LOC, ~14 tests)

**Goal**: Create sort modal for sorting facts by multiple fields.

#### Files to Create
- `internal/cli/components/fact_sort_modal.go`
- `internal/cli/components/fact_sort_modal_test.go`

#### Changes Required

1. **Define structs**:
   - `FactSortConfig` (SortBy string, SortOrder string)
   - `FactSortFormData`
   - `FactSortModal`

2. **Sort options**: text, role, company, created

3. **Implement sort comparators** using TableBehavior.SetSort()

4. **Add `roleOrder()` helper** for role sorting

#### TDD Checklist

##### RED Phase
- [ ] Tests written and **FAIL**:
  - [ ] TestNewFactSortModal_CreatesModal
  - [ ] TestFactSortModal_BuildForm_SortByOptions
  - [ ] TestFactSortModal_BuildForm_SortOrderOptions
  - [ ] TestFactSortModal_Update_AppliesOnEnter
  - [ ] TestFactSortModal_ToSortConfig
  - [ ] TestFactManagement_SKeyOpensSortModal
  - [ ] TestFactManagement_ApplySortConfig_ByText
  - [ ] TestFactManagement_ApplySortConfig_ByRole
  - [ ] TestFactManagement_ApplySortConfig_ByCompany
  - [ ] TestFactManagement_ApplySortConfig_ByCreated
  - [ ] TestFactManagement_ApplySortConfig_Descending
  - [ ] TestRoleOrder_ReturnsCorrectOrder

##### GREEN Phase
- [ ] Implementation written
- [ ] All tests pass

#### Pre-Commit Checklist
- [ ] `make check-compliance` passes
- [ ] Use `make ai-commit MSG="feat(fact-management): add sort modal with multiple sort options"`

---

### Phase 6: Single Fact Re-classify (~60 LOC, ~12 tests)

**Goal**: Add ability to re-classify a single fact using the category-aware classifier.

#### Files to Modify
- `internal/cli/intents/fact_management_intent.go`

#### Changes Required

1. **Add classifier to model**:
   ```go
   classifier *burst_fact.Classifier
   ```

2. **Initialize in constructor**

3. **Add `"c"` key handler** in `handleViewState()`

4. **Implement `reclassifyFact()` method**:
   - Get categories from source event via EventMap
   - Call `ClassifyRoleFitWithCategories()`
   - Update fact if role changed
   - Refresh stats and table

5. **Update FactViewState footer** with Re-classify badge

#### TDD Checklist

##### RED Phase
- [ ] Tests written and **FAIL**:
  - [ ] TestReclassifyFact_ChangesRole
  - [ ] TestReclassifyFact_UsesEventCategories
  - [ ] TestReclassifyFact_NoChangeWhenSameRole
  - [ ] TestReclassifyFact_HandlesNoSourceEvent
  - [ ] TestReclassifyFact_HandlesUpdateError
  - [ ] TestReclassifyFact_RefreshesStats
  - [ ] TestFactViewState_CKeyReclassifies
  - [ ] TestFactViewState_Footer_ShowsReclassifyBadge

##### GREEN Phase
- [ ] Implementation written
- [ ] All tests pass

#### Pre-Commit Checklist
- [ ] `make check-compliance` passes
- [ ] Use `make ai-commit MSG="feat(fact-management): add single fact re-classify action"`

---

### Phase 7: Batch Re-classify with Confirm + Progress + Summary (~180 LOC, ~18 tests)

**Goal**: Add batch re-classify with confirmation dialog, progress indicator, and summary modal.

#### Files to Modify
- `internal/cli/intents/fact_management.go` (new states)
- `internal/cli/intents/fact_management_intent.go`

#### Changes Required

1. **Add new states**:
   ```go
   FactReclassifyConfirmState FactManagementState = "reclassify_confirm"
   FactReclassifyingState     FactManagementState = "reclassifying"
   FactReclassifySummaryState FactManagementState = "reclassify_summary"
   ```

2. **Add `ReclassifyProgress` struct and field**

3. **Define message types**:
   - `ReclassifyProgressMsg`
   - `ReclassifyCompleteMsg`

4. **Add `"C"` key handler** in `handleListState()` (capital C)

5. **Implement state handlers**:
   - `handleReclassifyConfirmState()` (y/n/esc)
   - `handleReclassifyingState()` (progress + cancel)
   - `handleReclassifySummaryState()` (enter/esc)

6. **Implement batch processing**:
   - `startBatchReclassify()`
   - `processNextFact()` command

7. **Implement view methods**:
   - `getReclassifyConfirmContent()`
   - `getReclassifyingContent()`
   - `getReclassifySummaryContent()`

8. **Update footers** for new states

#### TDD Checklist

##### RED Phase
- [ ] Tests written and **FAIL**:
  - [ ] TestFactListState_CKeyShowsConfirmDialog
  - [ ] TestReclassifyConfirm_YStartsBatch
  - [ ] TestReclassifyConfirm_NCancels
  - [ ] TestReclassifyConfirm_EscCancels
  - [ ] TestBatchReclassify_ProcessesAllFacts
  - [ ] TestBatchReclassify_TracksProgress
  - [ ] TestBatchReclassify_TracksChanged
  - [ ] TestBatchReclassify_EscCancelsEarly
  - [ ] TestBatchReclassify_ShowsSummaryOnComplete
  - [ ] TestReclassifySummary_ShowsStats
  - [ ] TestReclassifySummary_EnterReturnsToList
  - [ ] TestReclassifySummary_EscReturnsToList
  - [ ] TestGetReclassifyConfirmContent
  - [ ] TestGetReclassifyingContent
  - [ ] TestGetReclassifySummaryContent

##### GREEN Phase
- [ ] Implementation written
- [ ] All tests pass

#### Pre-Commit Checklist
- [ ] `make check-compliance` passes
- [ ] Use `make ai-commit MSG="feat(fact-management): add batch re-classify with progress and summary"`

---

### Phase 8: FilterBehavior Interface + Clear Filters (~50 LOC, ~10 tests)

**Goal**: Implement FilterBehavior interface with FIFO clearing.

#### Files to Modify
- `internal/cli/intents/fact_management_intent.go`

#### Changes Required

1. **Add interface assertion**:
   ```go
   var _ FilterBehavior = (*FactManagementModel)(nil)
   ```

2. **Implement interface methods**:
   - `HasActiveFilters()` - check role, company, sort
   - `ClearFilters()` - FIFO: role → company → sort
   - `ApplyFilters()` - reapply current filters
   - `RefreshData()` - reload from repository

3. **Add `"x"` key handler** for clearing

4. **Add `reapplyFilters()` helper**

#### TDD Checklist

##### RED Phase
- [ ] Tests written and **FAIL**:
  - [ ] TestHasActiveFilters_NoFilters
  - [ ] TestHasActiveFilters_RoleFilter
  - [ ] TestHasActiveFilters_CompanyFilter
  - [ ] TestHasActiveFilters_SortOnly
  - [ ] TestClearFilters_FIFO_RoleFirst
  - [ ] TestClearFilters_FIFO_CompanySecond
  - [ ] TestClearFilters_FIFO_SortLast
  - [ ] TestFactListState_XKeyClearsFilters
  - [ ] TestRefreshData_ReloadsFromRepo

##### GREEN Phase
- [ ] Implementation written
- [ ] All tests pass

#### Pre-Commit Checklist
- [ ] `make check-compliance` passes
- [ ] Use `make ai-commit MSG="feat(fact-management): implement FilterBehavior interface with FIFO clearing"`

---

### Phase 9: SortBadge + Final Polish (~5 LOC, ~2 tests)

**Goal**: Add SortBadge helper and finalize footers.

#### Files to Modify
- `internal/cli/uikit/primitives/badge.go`
- `internal/cli/intents/fact_management_intent.go`

#### Changes Required

1. **Add `SortBadge()` function**:
   ```go
   func SortBadge(th Theme) *Badge {
       return HelpKeyBadge("s", "Sort", th)
   }
   ```

2. **Update `FactListState` footer** to use all badges conditionally

#### TDD Checklist

##### RED Phase
- [ ] Tests written and **FAIL**:
  - [ ] TestSortBadge_ReturnsCorrectBadge
  - [ ] TestFactListState_Footer_ShowsSortBadge

##### GREEN Phase
- [ ] Implementation written
- [ ] All tests pass

#### Pre-Commit Checklist
- [ ] `make check-compliance` passes
- [ ] Use `make ai-commit MSG="feat(uikit): add SortBadge helper and finalize fact management footers"`

---

## Post-Task Checklist (MUST COMPLETE BEFORE MARKING DONE)

- [ ] `make check-compliance` passes
- [ ] All checkboxes above completed
- [ ] All tests pass: `go test -v ./internal/cli/intents/... ./internal/cli/components/...`
- [ ] No race conditions: `go test -race ./...`
- [ ] Manual testing completed:
  - [ ] Filter by role works
  - [ ] Filter by company works
  - [ ] Sort by all fields works
  - [ ] Clear filters (x) works with FIFO
  - [ ] Single re-classify (c) works
  - [ ] Batch re-classify (C) with confirm/progress/summary works
  - [ ] Role coverage header shows correct counts
  - [ ] Columns display correctly
- [ ] Task marked complete `[x]` in task file

## Acceptance Criteria

- [ ] Fact table shows 5 columns: Fact, Role, Company, Strength, Categories
- [ ] Role coverage header displays above table with correct counts
- [ ] Filter modal (`f`) allows filtering by Role AND Company with counts
- [ ] Sort modal (`s`) allows sorting by text, role, company, created date
- [ ] Clear filters (`x`) clears in FIFO order: role → company → sort
- [ ] Re-classify single fact (`c` in view) uses category-aware classifier
- [ ] Batch re-classify (`C` in list) shows confirm → progress → summary
- [ ] Batch re-classify can be cancelled mid-operation
- [ ] All tests pass (target: ~97 new tests)
- [ ] Coverage maintained ≥ 80%
- [ ] No regressions in existing functionality

## Rollback Plan

1. **If Phase 1-3 breaks existing functionality**:
   - Revert context changes
   - Remove new fields/methods
   - Existing 3-column table continues working

2. **If filter/sort modals have issues**:
   - Remove modal imports
   - Remove keyboard handlers
   - Table still works without filter/sort

3. **If re-classify has issues**:
   - Remove classifier import
   - Remove c/C handlers
   - Manual editing via `e` key still works

## Summary

| Metric | Target |
|--------|--------|
| **Total New LOC** | ~845 |
| **Total New Tests** | ~97 |
| **Files Created** | 4 |
| **Files Modified** | 4 |
| **New States** | 3 |
| **New Keyboard Shortcuts** | 5 (`f`, `s`, `x`, `c`, `C`) |

## Keyboard Shortcuts Reference

| Key | State | Action |
|-----|-------|--------|
| `f` | FactListState | Open filter modal (role + company) |
| `s` | FactListState | Open sort modal |
| `x` | FactListState | Clear filters (FIFO) |
| `c` | FactViewState | Re-classify single fact |
| `C` | FactListState | Re-classify ALL facts (batch) |
