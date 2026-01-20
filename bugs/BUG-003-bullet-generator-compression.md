# BUG-003: BulletGenerator Role-Based Compression Causes Missing Companies in CV

**Status**: FIXED  
**Severity**: HIGH  
**Reported**: 2026-01-20  
**Fixed**: 2026-01-20  
**Component**: `internal/service/career/cv/bullet_generator.go`  
**Affected Versions**: All versions prior to fix  
**Reporter**: User (via CV generation testing)

---

## Summary

The `DefaultBulletGenerator` applied a total bullet cap based on role (e.g., 40 bullets for "staff") **before** section building, causing many companies to be completely excluded from generated CVs. This contradicts the intended design where compression should only happen:

1. Per-company/project in `SectionBuilder` (correct behavior)
2. Based on audience relevance (correct behavior)

---

## Symptoms

Users reported that generated CVs were missing companies from their career history. When generating a CV with the following settings:
- Role: Senior Engineer / Staff
- Audience: Peer/Colleague  
- Technology Focus: Language Agnostic
- Structure: Grouped by Category
- Bullets: 5, Detailed

Only **14 out of 25 companies** appeared in the output.

---

## Impact

| Impact Area | Description |
|-------------|-------------|
| **Data Loss** | 11 out of 25 companies (44%) completely excluded from CV |
| **User Experience** | Users see incomplete career history |
| **CV Quality** | Severely degraded - missing significant work experience |
| **Scalability** | Problem worsens with more career events |
| **Trust** | Users cannot rely on generated CVs for job applications |

### Affected User Workflows

- CV Generation wizard
- CV Export (all formats: markdown, YAML, text)
- Any workflow that uses `CVGenerationService`

---

## Root Cause Analysis

### Direct Cause

In `bullet_generator.go`, the `GenerateBullets` function called `compressByRole()` which applied a hardcoded total bullet cap:

```go
// Line 62-63 in bullet_generator.go (BEFORE FIX)
compressedBullets := bg.compressByRole(rankedBullets, targetRole)
```

```go
// Lines 256-269 in bullet_generator.go (BEFORE FIX)
func (bg *DefaultBulletGenerator) getBulletCapForRole(targetRole string) int {
    switch strings.ToLower(targetRole) {
    case "principal":
        return 50
    case "staff":
        return 40  // <-- THIS WAS THE PROBLEM
    case "em":
        return 40
    case "senior_ic":
        return 40
    default:
        return 30
    }
}
```

### Why This Was Wrong

The compression happened **before** bullets were grouped by company. This meant:

1. 690 events in database
2. 639 bullets generated after filtering (correct)
3. **Compressed to 40 bullets** (bug) - arbitrary cut based on rank alone
4. Those 40 bullets only covered 14 companies
5. 11 companies had zero representation

### Architectural Issue

Two bullet generators exist with different designs:

| Generator | Total Cap | Per-Company Cap | Used By |
|-----------|-----------|-----------------|---------|
| `BulletGenerator` | Yes (buggy) | No | `CVGenerationService` |
| `EnhancedBulletGenerator` | No (correct) | Deferred to SectionBuilder | Intent context only |

The `CVGenerationService` was wired to use the basic `BulletGenerator` instead of the correctly-designed `EnhancedBulletGenerator`.

---

## Evidence

### Log Output (Before Fix)

```
Compressed bullets from 639 to 40 for role staff
Generated 40 bullets from 690 events and 689 facts
groupBulletsByCompany: created 14 groups, skipped 13 bullets without company
buildExperienceSection: found 14 company groups from 40 bullets
```

### Missing Companies

The following companies were excluded from generated CVs:

1. BEIS
2. CrowdVision
3. Mindful Chef (2022)
4. Mindful Chef (2023-2024)
5. Money Advice Service
6. NTTData
7. Nature Publishing Group
8. RWDMag
9. Spice Rack (2017)
10. Spice Rack (2022)
11. mGage

### Database vs Output Comparison

| Metric | Database | CV Output | Difference |
|--------|----------|-----------|------------|
| Total Events | 690 | N/A | - |
| Total Facts | 689 | N/A | - |
| Bullets After Filtering | 639 | 40 | -599 (93.7% loss) |
| Companies Represented | 25 | 14 | -11 (44% loss) |

---

## Intended Design

The `EnhancedBulletGenerator` documents the correct design:

```go
// Line 173 in enhanced_bullet_generator.go
// Note: Per-company/project caps are applied by SectionBuilder based on audience
```

### Correct Flow

```
Events/Facts
    │
    ▼
BulletGenerator
    ├── Filter by role relevance
    ├── Filter by audience relevance
    ├── Deduplicate
    ├── Rank by relevance
    └── Return ALL bullets (no total cap)
    │
    ▼
SectionBuilder
    ├── Group by company/project
    ├── Apply per-company cap (4-5 bullets)
    └── Build CV sections
```

### Buggy Flow (Before Fix)

```
Events/Facts
    │
    ▼
BulletGenerator
    ├── Filter by role relevance
    ├── Filter by audience relevance
    ├── Deduplicate
    ├── Rank by relevance
    └── TRUNCATE to 40 bullets  ◄── BUG HERE
    │
    ▼
SectionBuilder
    ├── Group by company (only 14 companies have bullets!)
    ├── Apply per-company cap
    └── Build incomplete CV sections
```

---

## Why This Happened

### Timeline

1. **Initial Implementation**: `BulletGenerator` created with naive total cap approach
2. **Later**: `EnhancedBulletGenerator` created with correct design (no total cap)
3. **Wiring Issue**: `CVGenerationService` continued using old `BulletGenerator`
4. **Undetected**: Tests only verified cap was applied, not that all companies were represented

### Contributing Factors

1. **Two generators with different designs** - technical debt
2. **Tests verified buggy behavior** - tests checked `<= 40` instead of company coverage
3. **No integration test** for "all companies should appear in CV"
4. **Lack of logging** showing which companies were excluded

---

## Resolution

### Fix Applied: Option 1 (Minimal Change)

Removed the total bullet cap from `BulletGenerator`. The `SectionBuilder` already handles per-company caps correctly.

### Code Changes

#### `bullet_generator.go`

**Removed Functions:**
- `compressByRole()`
- `getBulletCapForRole()`

**Modified `GenerateBullets()`:**

```go
// GenerateBullets generates ranked CV bullets from events and facts
// Note: Per-company/project bullet caps are applied by SectionBuilder, not here.
// This function filters by role/audience relevance and ranks bullets, but does not
// apply any total bullet cap. See BUG-003 for rationale.
func (bg *DefaultBulletGenerator) GenerateBullets(ctx context.Context, events []*career.CareerEvent, facts []*career.Fact, targetRole string, targetAudience string) ([]*career.CVBullet, error) {
    // ... validation ...

    // Generate initial bullets from events and facts
    bullets := bg.generateInitialBullets(events, facts, targetRole, targetAudience)

    // Apply inclusion criteria filtering
    filteredBullets := bg.filterByInclusionCriteria(bullets, targetRole)

    // Rank the bullets
    rankedBullets := bg.rankBullets(filteredBullets, events)

    // NOTE: We intentionally do NOT apply a total bullet cap here.
    // Per-company/project caps are correctly applied by SectionBuilder.getBulletsPerCompanyForRole()
    // A total cap here would exclude entire companies from the CV. (See BUG-003)

    bg.logger.Info("Generated %d bullets from %d events and %d facts for role %s", 
        len(rankedBullets), len(events), len(facts), targetRole)
    return rankedBullets, nil
}
```

#### `bullet_generator_test.go`

**Updated Test Section:**

Changed from "Role-Specific Compression" to "Bullet Generation Without Total Cap (BUG-003 Fix)"

Tests now verify:
- All events passing inclusion criteria are returned
- No total bullet cap is applied
- Bullets are properly ranked for SectionBuilder to use

### Files Changed

| File | Change Type | Description |
|------|-------------|-------------|
| `internal/service/career/cv/bullet_generator.go` | Modified | Removed compression logic |
| `internal/service/career/cv/bullet_generator_test.go` | Modified | Updated 5 tests for correct behavior |
| `bugs/BUG-003-bullet-generator-compression.md` | Added | This document |

---

## Verification

### Test Results

```
Ran 375 of 375 Specs in 0.163 seconds
SUCCESS! -- 375 Passed | 0 Failed | 0 Pending | 0 Skipped
```

### Build Status

- Application builds successfully
- No compilation errors
- No new warnings

### Expected Behavior After Fix

| Metric | Before Fix | After Fix |
|--------|------------|-----------|
| Bullets to SectionBuilder | 40 | 639 |
| Companies in CV | 14 | 25 |
| Per-company cap | 4-5 | 4-5 (unchanged) |
| CV length | ~40 bullets | ~100-125 bullets (25 companies x 4-5) |

---

## Testing Checklist

- [x] All 375 CV service tests pass
- [x] Application builds successfully
- [ ] Manual test: Generate CV and verify all 25 companies appear
- [ ] Manual test: Verify per-company bullet caps still work (4-5 per company)
- [ ] Manual test: CV export shows complete career history
- [ ] Manual test: CV length is reasonable (not excessively long)

---

## Follow-Up Tasks

### Recommended

1. **Migrate to EnhancedBulletGenerator** - The correctly-designed generator should be used by `CVGenerationService`
2. **Add integration test** - Verify all companies with events appear in generated CV
3. **Add company coverage logging** - Log which companies are included/excluded

### Optional

1. **Deprecate BulletGenerator** - Mark as deprecated in favor of `EnhancedBulletGenerator`
2. **Consolidate generators** - Merge functionality into single generator

---

## Lessons Learned

1. **Test behavior, not implementation** - Tests should verify "all companies appear" not "cap is 40"
2. **Integration tests catch architectural issues** - Unit tests passed but integration was broken
3. **Document design decisions** - `EnhancedBulletGenerator` had the correct design but it wasn't used
4. **Avoid parallel implementations** - Two generators with different designs caused confusion
5. **Log filtering decisions** - Would have caught this earlier if we logged "excluding company X"

---

## Related Documentation

- `internal/service/career/cv/enhanced_bullet_generator.go` - Correctly designed alternative
- `internal/service/career/cv/section_builder.go` - Contains correct per-company capping
- `SectionBuilder.getBulletsPerCompanyForRole()` - The correct place for bullet limits

---

## Appendix: Git History

### Commit that introduced the bug

```
ba34f74 - feat(cv-generation): implement bullet generator service with filtering and ranking
```

The `compressByRole` function was added in the initial implementation of `BulletGenerator`.

### Fix Commit

To be created with message:
```
fix(cv): remove total bullet cap causing missing companies in CV (BUG-003)
```

---

## Reproduction Steps

1. Import career events (690+ events across 25 companies)
2. Open KaRiya TUI
3. Navigate to "Generate CV"
4. Select: Role: Staff/Senior Engineer, Audience: Peer/Colleague
5. Select: Technology Focus: Language Agnostic
6. Select: Structure: Grouped by Category, 5 bullets, Detailed
7. Generate CV
8. Observe: Only 14 companies appear instead of 25

**Consistency**: Always (100% reproducible with sufficient career data)

**Environment**:
- OS: Linux (any)
- Terminal: Any
- Go Version: 1.24+
- KaRiya Version: All versions prior to fix

---

**Last Updated**: 2026-01-20  
**Updated By**: AI Assistant (Claude)
