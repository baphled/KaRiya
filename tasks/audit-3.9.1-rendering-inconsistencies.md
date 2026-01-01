# Task 3.9.1 - Audit Report: List Model Rendering Inconsistencies

**Date**: 2026-01-01
**Status**: IN PROGRESS
**Priority**: CRITICAL

## Executive Summary

Audit of three list models (list.go, burst_list.go, fact_list.go) reveals significant inconsistencies in rendering pipelines, pagination handling, and visual presentation despite all using container-based architecture.

## 1. Rendering Pipeline Comparison

### list.go
- Uses dedicated `renderListWithContainer()` helper method
- Items rendered via `renderListItems()` helper
- Pagination format: `"Page %d of %d (%d total events)"`
- Container setup:
  ```go
  listContainer := components.NewListContainer().
      SetItems(items).
      SetEmptyStateMessage("No events found. Start capturing...").
      SetPaginationInfo(paginationText)
  ```
- Wrapping: Uses `CardBase` style before header/footer
- View flow:
  1. Build content array
  2. Check errors
  3. Render with ListContainer
  4. Add help footer
  5. Wrap with CardBase
  6. Add header/footer
  7. Join vertically

### burst_list.go
- Direct loop with `getDisplayedBursts()`
- Items built with `renderBurstRow()` and `renderExpandedEvents()`
- Pagination format: Not explicitly shown (may be missing)
- Container setup:
  ```go
  listContainer := components.NewListContainer().
      SetItems(items).
      SetEmptyStateMessage(emptyStateMessage)
  ```
- Wrapping: Uses `ScreenContainer` with PaddingNormal
- Custom header rendering with explicit column widths
- View flow:
  1. Get displayed bursts
  2. Loop and build items
  3. Setup ListContainer
  4. Add ScreenContainer wrapper
  5. Add header/footer
  6. Join vertically
- **Issue**: No `SetPaginationInfo()` call - pagination info may be missing from UI

### fact_list.go
- Early return for empty state via `renderEmpty()`
- Manual pagination logic with scrollOffset + height calculation
- Items built with `renderFactItem()`
- Pagination format: `"%d/%d facts"` with selected count
- Container setup:
  ```go
  listContainer := components.NewListContainer().
      SetItems(items).
      SetEmptyStateMessage("No facts found").
      SetPaginationInfo(paginationInfo)
  ```
- Wrapping: Uses `ScreenContainer` with PaddingNormal
- View flow:
  1. Check if empty (different pattern)
  2. Manual pagination calculation
  3. Loop and build items
  4. Create pagination info
  5. Setup ListContainer
  6. Add ScreenContainer wrapper
  7. Add header/footer
  8. Join vertically

## 2. Key Inconsistencies Identified

### 2.1 Rendering Pipeline Differences
- **list.go**: Uses abstracted helper methods (renderListWithContainer, renderListItems)
- **burst_list.go**: Direct inline rendering with custom header
- **fact_list.go**: Mixed approach with early return and manual pagination

**Problem**: Users see different visual patterns and layouts across list types

### 2.2 Pagination Display
| Model | Format | SetPaginationInfo() | Selected Count |
|-------|--------|-------------------|-----------------|
| list.go | Page X of Y (Z total) | ✅ Yes | Not shown |
| burst_list.go | (MISSING?) | ❌ No | N/A |
| fact_list.go | X/Y facts \| N selected | ✅ Yes | ✅ Shown |

**Problem**: Inconsistent pagination information displayed to users

### 2.3 Header Rendering
- **list.go**: Uses Header component indirectly
- **burst_list.go**: Custom `renderHeader()` with hardcoded column widths
- **fact_list.go**: Uses Header component with dynamic title

**Problem**: Burst list headers may not align with other lists

### 2.4 Container Wrapping
- **list.go**: `CardBase` style + header/footer arrangement
- **burst_list.go**: `ScreenContainer` with PaddingNormal
- **fact_list.go**: `ScreenContainer` with PaddingNormal

**Problem**: list.go uses different wrapper style than the others

### 2.5 Empty State Handling
- **list.go**: Checked in ListContainer via SetEmptyStateMessage
- **burst_list.go**: Checked in ListContainer via SetEmptyStateMessage
- **fact_list.go**: **Early return** via renderEmpty() - completely different pattern

**Problem**: Empty state messages and styling may be inconsistent

### 2.6 Item Rendering Helpers
| Model | Helper Method | Pattern |
|-------|--------------|---------|
| list.go | renderListItems() | Returns string array |
| burst_list.go | renderBurstRow() | Returns single string |
| fact_list.go | renderFactItem() | Returns single string |

**Problem**: Different helper patterns make code less consistent

## 3. Visual Output Comparison

### Expected vs Actual Structure

**Standard expected structure** (should be identical for all):
```
┌─ Header Component ─────────────────────┐
├────────────────────────────────────────┤
│ [List Container with items]            │
│  - Item 1                              │
│  - Item 2 (selected)                   │
│  - Item 3                              │
│ Pagination Info: X of Y                │
├────────────────────────────────────────┤
│ Footer Component                       │
└────────────────────────────────────────┘
```

**Current issues**:
1. list.go uses CardBase wrapper not used by others
2. burst_list.go may be missing pagination info
3. fact_list.go has different empty state handling
4. Column headers may not align (burst_list custom header)

## 4. Key Handling Pattern Comparison

From recent standardization work:
- All three models NOW support: j/k, g/G, PageUp/PageDown, space, esc, q, ctrl+c
- **BUT**: burst_list.go uses string-based key matching (older pattern)
- **BUT**: list.go and fact_list.go use type-based approach

### Update() Method Key Handling
| Model | Pattern | Status |
|-------|---------|--------|
| list.go | Type-based (tea.KeyType) | ✅ Current |
| burst_list.go | String-based matching | ⚠️ Legacy |
| fact_list.go | Type-based (tea.KeyType) | ✅ Current |

**Problem**: Inconsistent approach makes maintenance harder

## 5. Test Coverage Assessment

**Current test files**:
- list_test.go - 853 tests passing
- burst_list_test.go - Tests updated for string key format
- fact_list_test.go - Tests present

**Missing**:
- Visual regression tests comparing output
- Rendering consistency tests across all three models
- Pagination format validation tests
- Empty state consistency tests

## 6. Style Usage Verification

**Status Check**:
- All models use exported constants ✅
- No inline hex colors detected ✅
- Consistent color scheme usage ✅

## 7. Root Causes

1. **Container adoption was gradual**: Models were refactored at different times
2. **No standardized rendering template**: Each model evolved independently
3. **Manual pagination logic**: fact_list and burst_list handle pagination differently
4. **Custom headers**: burst_list has special column handling not in others
5. **Different wrapper patterns**: list.go diverged with CardBase usage
6. **Helper method variations**: Different abstraction levels across models

## 8. Impact on Users

Users navigating between list views will observe:
1. Inconsistent pagination display format
2. Potentially missing pagination info on bursts
3. Different visual wrapping/padding
4. Inconsistent column alignment (bursts)
5. Different empty state presentation (facts vs others)
6. Inconsistent navigation feel

## 9. Recommendations for Fix

### Phase 1: Standardize Rendering (next subtask)
- Extract common rendering pipeline to shared helper
- Ensure identical View() method structure across all three
- Standardize pagination display format
- Use consistent container wrapping

### Phase 2: Standardize Key Handling
- Update burst_list.go to use type-based key matching
- Ensure all three use identical key handling patterns

### Phase 3: Visual Verification
- Run all three list views side-by-side
- Verify pixel-perfect visual consistency
- Create visual regression tests

### Phase 4: Documentation
- Document standard list rendering specification
- Create template for future list implementations

## 10. Files to Modify

- `internal/cli/models/list.go` - May need minor adjustments
- `internal/cli/models/burst_list.go` - Major refactoring needed
- `internal/cli/models/fact_list.go` - Moderate refactoring needed
- `internal/cli/models/*_test.go` - Add regression tests

## Verification Points

- [ ] All three models produce identical visual output
- [ ] Pagination displayed consistently
- [ ] Headers aligned consistently
- [ ] Empty states consistent
- [ ] Navigation behavior identical
- [ ] All 853+ tests passing
- [ ] No race conditions
- [ ] Coverage maintained 76%+

---

**Next Steps**: Proceed with task 3.9.2 to implement standardized rendering pipeline

