# BUG-009: CV Generator Merges Separate Company Tenures

**Status**: FIXED  
**Severity**: Medium  
**Created**: 2026-01-22  
**Updated**: 2026-01-22  
**Fixed**: 2026-01-22  
**Assigned To**: Claude Code

---

## Bug Summary

The CV generator merges all events for a company into a single experience entry, even when someone worked at that company during separate periods with other jobs in between. This results in incorrect date ranges and forced use of year suffixes in company names to work around the issue.

**Root Cause**: `GroupEventsByCompany` uses company name as the sole grouping key, with no tenure detection logic.

---

## Affected Components

### Files with Grouping Logic

| File | Function | Issue |
|------|----------|-------|
| `internal/service/career/cv/data_processing_service.go` | `GroupEventsByCompany` | Groups by company name only |
| `internal/service/career/cv/section_builder.go` | `groupBulletsByCompany` | Groups by company name only |

### Data Workaround

**File**: `career_entries_with_skills.csv`

Currently uses year suffixes to force separate groupings:
- `Mindful Chef (2022)` vs `Mindful Chef (2023-2024)`
- `Spice Rack (2017)` vs `Spice Rack (2022)`
- `We Are Friday (2012)` vs `We Are Friday (2015)`
- `Freelance (Early Career)` vs `Freelance (2021)` vs `Freelance / Contract (2008-2012)`

---

## Reproduction Steps

1. Have career events at Company A (e.g., Mar 2022 - Jul 2022)
2. Have career events at Company B (e.g., Aug 2022 - Jun 2023)
3. Have career events at Company A again (e.g., Jul 2023 - Apr 2024)
4. Generate a CV using the CV generator
5. **Observe**: Company A shows as single entry spanning Mar 2022 - Apr 2024

---

## Expected Behavior

The CV generator should detect separate tenures based on intervening employment:

```
Timeline: [Company A: Mar-Jul 2022] [Company B: Aug 2022-Jun 2023] [Company A: Jul 2023-Apr 2024]

Expected Output:
### Company A - _Jul 2023 - Apr 2024_
- bullets from second tenure...

### Company A - _Mar 2022 - Jul 2022_
- bullets from first tenure...
```

**Tenure Detection Rule**: A new tenure is detected when events at OTHER companies (including Freelance) exist between two date ranges at the same company.

---

## Actual Behavior

All events for a company are merged into a single group:

```
### Company A - _Mar 2022 - Apr 2024_
- bullets from BOTH tenures merged together...
```

Or with current workaround (redundant date display):

```
### Company A (2023-2024) - _Jul 2023 - Apr 2024_
- bullets...

### Company A (2022) - _Mar 2022 - Jul 2022_
- bullets...
```

---

## Root Cause Analysis

### Current Implementation

**File**: `internal/service/career/cv/data_processing_service.go:122-130`

```go
// Group events by company
eventsByCompany := make(map[string][]*career.CareerEvent)
for _, event := range events {
    company := event.Company
    if company == "" {
        company = "Other"
    }
    eventsByCompany[company] = append(eventsByCompany[company], event)
}
```

**Problem**: Uses `company` as the sole map key. No logic to detect gaps or separate tenures.

### Missing Logic

```go
// NEEDED: Detect tenure boundaries
// A boundary exists when events at OTHER companies fall between 
// two events at this company
func detectTenures(companyEvents, allEvents []*CareerEvent) [][]*CareerEvent
```

---

## Fix Strategy

### 1. Add Tenure Detection Function

```go
// detectTenures splits a company's events into separate tenure groups
// based on whether events at OTHER companies exist between them.
func (svc *DefaultDataProcessingService) detectTenures(
    companyEvents []*career.CareerEvent,
    allEventsSorted []*career.CareerEvent,
) [][]*career.CareerEvent {
    // Sort company events by date
    // For each pair of consecutive company events:
    //   Check if any events from OTHER companies fall between them
    //   If yes, create a tenure boundary
    // Return slice of event slices (each is a tenure)
}
```

### 2. Update GroupEventsByCompany

- Sort all events chronologically first
- For each company, call `detectTenures`
- Create separate `CompanyGroup` for each tenure
- Sort final groups reverse-chronologically (newest first)

### 3. Update groupBulletsByCompany

Apply same tenure-aware logic in `section_builder.go`.

### 4. Clean CSV Data

Remove year suffixes from company names:
- `Mindful Chef (2022)` -> `Mindful Chef`
- `Mindful Chef (2023-2024)` -> `Mindful Chef`
- etc.

---

## Testing Plan

### Unit Tests (data_processing_service_test.go)

| Test Case | Scenario |
|-----------|----------|
| `should detect single tenure when no gaps` | Company A events only |
| `should detect multiple tenures with intervening company` | A -> B -> A |
| `should treat freelance as tenure separator` | A -> Freelance -> A |
| `should order tenures reverse-chronologically` | Newest first |
| `should handle adjacent same-company events as single tenure` | No gap = one tenure |

### Integration Tests

| Test Case | Scenario |
|-----------|----------|
| `should render separate experience entries per tenure` | CV output verification |
| `should display clean company names without year suffixes` | Header verification |

---

## Files to Modify

| File | Change |
|------|--------|
| `internal/service/career/cv/data_processing_service.go` | Add `detectTenures`, update `GroupEventsByCompany` |
| `internal/service/career/cv/section_builder.go` | Update `groupBulletsByCompany` |
| `internal/service/career/cv/data_processing_service_test.go` | Add tenure detection tests |
| `internal/service/career/cv/section_builder_test.go` | Add tenure-aware grouping tests |
| `career_entries_with_skills.csv` | Remove year suffixes from company names |

---

## Design Decisions

1. **Gap threshold**: Any gap with other companies triggers a new tenure (no minimum duration)
2. **Freelance handling**: Freelance periods count as tenure separators
3. **Display order**: Reverse-chronological (newest tenure first)

---

## Verification Checklist

- [x] Tenure detection correctly identifies separate engagements
- [x] Freelance periods properly separate tenures
- [x] CV output shows separate entries per tenure
- [ ] Company names display without year suffixes (CSV cleanup pending)
- [x] Reverse-chronological ordering works
- [x] All existing tests still pass
- [x] New tests provide adequate coverage
- [ ] CSV data cleaned of year suffixes (manual task pending)

---

## Fix Summary

**Implemented**: 2026-01-22

### Commits

| Commit | Description |
|--------|-------------|
| `3e3f770` | Initial fix: add tenure detection for separate company engagements |
| `c7f1867` | Cleanup: consolidate tenure detection utilities and remove dead code |
| `6f5e52b` | Refactor: replace hardcoded strings with constants |

### Changes Made

1. **`internal/constants/constants.go`**:
   - Added `DefaultCompanyName` constant (`"Other"`) for empty company fallback
   - Added `TenureSeparator` constant (`"#"`) for tenure key formatting

2. **`data_processing_service.go`**:
   - Added `detectTenures()` function to split events into separate tenure groups
   - Added shared `hasInterveningCompanyEvents()` helper (package-level function)
   - Updated `GroupEventsByCompany()` to use tenure detection
   - Multiple tenures use key format: `"Company A#1"`, `"Company A#2"`
   - Replaced hardcoded strings with constants

3. **`section_builder.go`**:
   - Added `bulletInfo` type for tenure tracking
   - Added `detectBulletTenures()` function for bullet-level tenure detection
   - Updated `groupBulletsByCompany()` to create separate groups per tenure
   - Uses shared `hasInterveningCompanyEvents()` from data_processing_service.go

### Cleanup Applied

- Removed duplicate `hasInterveningCompanyEvents()` from section_builder.go
- Removed unused `companyEventIDs` map in `detectTenures()`
- Removed unused `sourceEventID` field from `bulletInfo` struct
- Consistent empty company handling (treats `""` as `constants.DefaultCompanyName`)

### Test Coverage

- Added 5 tenure detection tests in `data_processing_service_test.go`
- Added 3 tenure-aware bullet grouping tests in `section_builder_test.go`
- All 342 CV tests pass
- Tests cover: single tenure, multiple tenures, Freelance separation, chronological ordering
- Tests use constants instead of hardcoded strings

### PR

- **PR #109**: https://github.com/baphled/KaRiya/pull/109
- **Branch**: `fix/bug-009-cv-tenure-detection`

---

**Last Updated**: 2026-01-22  
**Updated By**: Claude Code
