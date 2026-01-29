# BUG-014: Date calculation in groupBulletsByCompany ignores primary company

## Summary

`groupBulletsByCompany()` computes `earliestDate` and `latestDate` for each bullet
by iterating over ALL `SourceEventIDs` regardless of which company those events
belong to. When a bullet has cross-company SourceEventIDs (caused by BUG-013), the
date range spans multiple companies and corrupts the company section dates.

## Steps to Reproduce

1. Have a bullet with SourceEventIDs from multiple companies (see BUG-013)
2. Generate a CV via `kariya` TUI
3. Observe company sections with impossible date ranges (e.g., BEIS: Nov 2012 - Dec 2019)

## Expected Behavior

A bullet assigned to BEIS should compute its date range using only BEIS events
from its SourceEventIDs, producing Jun 2019 - Dec 2019.

## Actual Behavior

The date range is computed across ALL source events. If the bullet's SourceEventIDs
include a We Are Friday event from Nov 2012 and a BEIS event from Oct 2019, the
bullet gets `earliestDate = Nov 2012`, which propagates to the BEIS company
section start date.

## Root Cause

**File**: `internal/service/career/cv/section_builder.go:421-435`

```go
for _, eventID := range bullet.SourceEventIDs {
    if event, exists := eventMap[eventID]; exists {
        if event.Company != "" {
            companyCounts[event.Company]++
            // Dates computed across ALL companies - no filtering
            if earliestDate.IsZero() || event.Date.Before(earliestDate) {
                earliestDate = event.Date
            }
            if latestDate.IsZero() || event.Date.After(latestDate) {
                latestDate = event.Date
            }
        }
    }
}
```

The code determines `primaryCompany` as the most frequent company in
`companyCounts`, but computes dates across all companies in the same pass. It
should recompute dates using only events matching `primaryCompany`.

## Fix Plan

Split the loop into two passes, or recompute dates after determining the primary
company:

```go
// After determining primaryCompany, filter dates
earliestDate = time.Time{}
latestDate = time.Time{}
for _, eventID := range bullet.SourceEventIDs {
    if event, exists := eventMap[eventID]; exists {
        if event.Company == primaryCompany {
            if earliestDate.IsZero() || event.Date.Before(earliestDate) {
                earliestDate = event.Date
            }
            if latestDate.IsZero() || event.Date.After(latestDate) {
                latestDate = event.Date
            }
        }
    }
}
```

This is a defence-in-depth fix. Even after BUG-013 is resolved, this protects
against any future scenario where SourceEventIDs span multiple companies.

## Regression Tests

- Date range computed using only events from the primary company
- Bullet with SourceEventIDs from BEIS (2019) and We Are Friday (2012) assigned to
  BEIS produces dates Jun-Dec 2019 (not Nov 2012 - Dec 2019)
- Cross-cutting bullet with 16 source companies produces dates for winning company only

## Related

- BUG-013: Bullet deduplication merges across companies (root cause of cross-company SourceEventIDs)
- BUG-012: Tenure detection project events

## Severity

- [x] High - CV displays incorrect employment dates
