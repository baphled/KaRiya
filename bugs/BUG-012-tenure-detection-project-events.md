# BUG-012: Tenure detection treats project-only events as intervening company events

## Summary

`hasInterveningCompanyEvents()` assigns `DefaultCompanyName` ("Other") to events
with empty Company fields (project-only entries like n-vyro.io, QuikCV, KaRiya).
These are then treated as employment at a different company, causing false tenure
splits within a single employment stint when there is a month gap in entries.

## Steps to Reproduce

1. Import career events where a company has a month gap (e.g., Mindful Chef has
   entries for Feb 2024 and Apr 2024 but not Mar 2024)
2. Ensure project-only entries exist during that gap (e.g., n-vyro.io at Mar 2024)
3. Generate a CV via `kariya` TUI
4. Observe the company appears as multiple separate tenures

## Expected Behavior

Mindful Chef second stint should appear as one tenure: Jul 2023 - Apr 2024.

## Actual Behavior

Mindful Chef second stint is split into two tenures:

- Mindful Chef: Jul 2023 - Feb 2024
- Mindful Chef: Apr 2024

The n-vyro.io entry at Mar 2024 (Company="", Project="n-vyro.io") is assigned
`DefaultCompanyName = "Other"` and counted as an intervening event between the
Feb 2024 and Apr 2024 Mindful Chef entries.

## Root Cause

**File**: `internal/service/career/cv/data_processing_service.go:254-258`

```go
eventCompany := event.Company
if eventCompany == "" {
    eventCompany = constants.DefaultCompanyName  // "Other"
}
```

Project-only events (personal/side projects with no Company) should not influence
employment tenure detection. They run concurrently with employment and are not
"intervening" work at another employer.

## Fix Plan

In `hasInterveningCompanyEvents`, skip events with empty Company instead of
assigning a default:

```go
if event.Company == "" {
    continue  // Project-only events don't affect employment tenures
}
```

## Regression Tests

- `hasInterveningCompanyEvents` returns false when only project-only events exist between dates
- `hasInterveningCompanyEvents` still returns true when real company events intervene
- Mindful Chef with month gap and concurrent n-vyro.io entries produces 2 tenures (not 3)
- Existing BUG-009 tenure detection still correctly splits genuine separate stints

## Related

- BUG-009: CV tenure detection (the feature that introduced `hasInterveningCompanyEvents`)
- BUG-013: Bullet deduplication merges across companies
- BUG-014: Date calculation ignores primary company

## Severity

- [x] High - CV displays incorrect employment history
