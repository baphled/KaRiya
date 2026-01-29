# BUG-013: Bullet deduplication merges entries across companies

## Summary

`deduplicateBullets()` uses normalized text as the sole deduplication key, with no
company awareness. When two career events at different companies have identical
descriptions, their bullets are merged into a single bullet with `SourceEventIDs`
spanning multiple companies and wildly different date ranges.

## Steps to Reproduce

1. Import career events where identical descriptions exist at different companies
   (e.g., "Acted as senior stabilising engineer during late-stage delivery pressure"
   at both We Are Friday (Nov 2012) and BEIS (Oct 2019))
2. Generate a CV via `kariya` TUI
3. Observe company date ranges are corrupted

## Expected Behavior

Each company's entry should retain its own bullet with correct dates. BEIS should
show Jun 2019 - Dec 2019.

## Actual Behavior

BEIS shows "Nov 2012 - Dec 2019" because the deduplication merged the We Are Friday
bullet (Nov 2012) with the BEIS bullet (Oct 2019) into a single bullet with
`SourceEventIDs` from both events. The date calculation then takes the min/max
across all source events regardless of company.

This affects any company with text matching another company's entries. It is
especially damaging for recurring cross-cutting entries that appear identically
across many companies (e.g., "Designed and delivered scalable backend services
using Ruby on Rails..." appears at 16 different companies). After deduplication,
these become a single bullet with 16 SourceEventIDs spanning 18 years.

## Root Cause

**File**: `internal/service/career/cv/bullet_generator.go:401-440`

```go
normalizedText := strings.ToLower(strings.TrimSpace(bullet.Text))

if existing, exists := seen[normalizedText]; exists {
    // MERGES SourceEventIDs from different companies
    existing.SourceEventIDs = mergeUniqueStrings(existing.SourceEventIDs, bullet.SourceEventIDs)
}
```

The deduplication key is text-only. Identical work descriptions at different
employers are legitimate separate career events, not duplicates.

## Fix Plan

Make the deduplication key company-aware by including the primary company in the
key. This requires resolving the primary company for each bullet before
deduplication, or passing event context into the deduplication function.

Option A - composite key using source event company:

```go
primaryCompany := resolvePrimaryCompany(bullet, eventMap)
key := normalizedText + "|" + primaryCompany
```

Option B - deduplicate per-company instead of globally:

Group bullets by company first, then deduplicate within each group.

## Regression Tests

- Bullets with identical text but different source companies are NOT merged
- Bullets with identical text AND same source company ARE merged
- Cross-cutting entries (16 copies across companies) produce 16 separate bullets
- SourceEventIDs after dedup only contain events from the same company

## Related

- BUG-014: Date calculation ignores primary company (defence-in-depth fix)
- BUG-012: Tenure detection project events

## Severity

- [x] High - CV displays incorrect employment dates and corrupted company sections
