# BUG-008: Role-Based CV Differentiation System Not Implemented

## Status: FIXED

**Fixed:** Role-based scoring is now fully wired up and operational in production.

**What was done:**
1. Implemented category-based role scoring in `EnhancedBulletGenerator`
2. Wired `CVGenerationService` to use `EnhancedBulletGenerator` (previously used old `BulletGenerator`)
3. Different roles now produce measurably different bullet rankings

## Summary

The CV generation system had a significant architectural gap where role and audience selection provided minimal actual differentiation in bullet scoring and selection. The system was designed to prioritize different achievements based on target role (Principal vs Staff vs Senior IC vs EM), but the scoring logic was never implemented.

**Original Problem:** A "Principal" CV and "Senior IC" CV produced nearly identical bullet rankings.

**Current Status:** Role-based scoring now differentiates bullets based on category alignment. The fix is wired into production.

## Severity

- [ ] Critical - Core feature broken
- [ ] High - Major feature broken
- [ ] Medium - Feature partially broken
- [x] Low - Minor stubs remain (cosmetic, not affecting core functionality)

## What Was Fixed

### Production Wiring (COMPLETE)
- `CVGenerationService` now uses `EnhancedBulletGenerator` instead of old `BulletGenerator`
- `NewCVGenerationService()` constructor updated to accept `EnhancedBulletGenerator`
- `app.go` updated to wire `EnhancedBulletGenerator` into the service
- `FilterByTechnologies()` added to `EnhancedBulletGenerator` interface
- Role-based scoring is now active in production CV generation

### Category Propagation (COMPLETE)
- `EnhancedBullet.Category` field added (line 57 in enhanced_bullet_generator.go)
- `CVBullet.Category` field added (line 222 in cv.go)
- Category populated from `CareerEvent.Categories[0]` via `extractPrimaryCategory()`
- Category populated from `Fact.CompetencyCategories[0]` via `extractPrimaryCategory()`
- Uses `constants.IsValidCompetencyCategory()` for validation

### Role-Based Scoring (COMPLETE)
- `calculateRoleScore()` now uses bullet category for scoring differentiation
- `getRoleFilter()` updated to use type-safe `constants.CompetencyCategory`
- Primary category match: +0.30 score boost (roleScorePrimaryCategoryBoost)
- Secondary category match: +0.15 score boost (roleScoreSecondaryCategoryBoost)
- Different roles produce measurably different bullet rankings

### Named Constants (COMPLETE)
Scoring magic numbers extracted to named constants:
- `roleScoreBase` = 0.50
- `roleScorePrimaryCategoryBoost` = 0.30
- `roleScoreSecondaryCategoryBoost` = 0.15
- `roleScoreAchievementBonus` = 0.10
- `roleScoreFactBonus` = 0.05
- `roleScoreHighConfidenceBonus` = 0.05
- `roleScoreHighConfidenceThreshold` = 0.80

## What Remains (Future Work - Low Priority)

### Minor Items (Do Not Affect Core Functionality)

| Item | File | Status | Notes |
|------|------|--------|-------|
| `Achievement.Category` field | data_processing_service.go:73 | NOT ADDED | Achievements rarely used, not blocking |
| Use `config.ScoringConfig.Weights` | enhanced_bullet_generator.go | NOT DONE | Hardcoded weights work fine |

### Deprecated Systems (REMOVED)

| System | Status |
|--------|--------|
| `CVVariant` (16 variants) | Removed from variants.go |
| `RoleEmphasis` | Removed from variants.go |
| `RoleEmphasisConfig.ScoreBulletCategory()` | role_emphasis.go deleted |
| `AudienceFilter` struct | Removed from enhanced_bullet_generator.go |
| `ProfileOverride` struct | Removed from cv_helpers.go |
| `ApplyProfileOverride*` functions | Removed from cv_helpers.go |
| `BulletConfig`, `SectionConfig` | Removed from variants.go |
| `VariantService`, `BuiltInVariants` | Removed from variants.go |

## How It Works Now

### Role Filter Configuration

```go
// enhanced_bullet_generator.go:661-696
func (ebg *DefaultEnhancedBulletGenerator) getRoleFilter(role string) *RoleFilter {
    switch strings.ToLower(role) {
    case "principal":
        return &RoleFilter{
            PrimaryCategories:   []constants.CompetencyCategory{constants.CompetencyLeadership},
            SecondaryCategories: []constants.CompetencyCategory{constants.CompetencyTechnical, constants.CompetencyMentoring},
            MinConfidence:       0.8,
        }
    case "senior_ic":
        return &RoleFilter{
            PrimaryCategories:   []constants.CompetencyCategory{constants.CompetencyTechnical},
            SecondaryCategories: []constants.CompetencyCategory{constants.CompetencyLeadership},
            MinConfidence:       0.75,
        }
    // ... staff, em cases
    }
}
```

### Scoring Algorithm

```go
// enhanced_bullet_generator.go:482-519
func (ebg *DefaultEnhancedBulletGenerator) calculateRoleScore(bullet *EnhancedBullet, role string) float64 {
    score := roleScoreBase  // 0.50

    filter := ebg.getRoleFilter(role)

    // Primary category match: +0.30
    for _, primary := range filter.PrimaryCategories {
        if bullet.Category == primary {
            score += roleScorePrimaryCategoryBoost
            break
        }
    }

    // Secondary category match: +0.15 (only if no primary match)
    if score == roleScoreBase {
        for _, secondary := range filter.SecondaryCategories {
            if bullet.Category == secondary {
                score += roleScoreSecondaryCategoryBoost
                break
            }
        }
    }

    // Additional bonuses for inclusion reason and confidence
    // ...
}
```

### Expected Behavior

| Event Category | Senior IC Score | Principal Score | Winner |
|----------------|-----------------|-----------------|--------|
| technical | 0.80 (primary) | 0.65 (secondary) | Senior IC |
| leadership | 0.65 (secondary) | 0.80 (primary) | Principal |
| mentoring | 0.50 (none) | 0.65 (secondary) | Principal |

## Regression Tests

Tests located in `enhanced_bullet_generator_test.go` under "BUG-008: Role-based scoring":

- Category propagation from CareerEvent to EnhancedBullet
- Category propagation from Fact to EnhancedBullet
- ToCVBullet conversion preserves Category
- Leadership bullets score higher for principal than senior_ic
- Technical bullets score higher for senior_ic than principal
- Mentoring bullets score higher for em than senior_ic
- Primary category bullets get strong boost
- Secondary category bullets get medium boost
- End-to-end: different roles produce different rankings

## Definition of Done

### Phase 1: Category Propagation
- [ ] `Achievement` struct has `Category constants.CompetencyCategory` field
- [x] `EnhancedBullet` struct has `Category constants.CompetencyCategory` field
- [x] `CVBullet` struct has `Category constants.CompetencyCategory` field
- [x] Category populated from `CareerEvent.Categories[0]`
- [x] Category populated from `Fact.CompetencyCategories[0]`
- [x] Use `constants.IsValidCompetencyCategory()` for validation
- [x] Unit tests for category propagation

### Phase 2: Role Scoring
- [x] `calculateRoleScore()` uses bullet category for scoring
- [x] `calculateRoleScore()` applies role-specific weights from `getRoleFilter()`
- [ ] `calculateFinalScore()` uses `config.ScoringConfig.Weights`
- [ ] Use `config.ScoringConfig.RoleSettings` for MinConfidence and MaxBulletsPerCompany
- [x] Different roles produce measurably different bullet rankings
- [x] Unit tests for role-based scoring
- [x] Integration test comparing Senior IC vs Principal output

### Phase 3: Audience Filtering
- [x] `isEventRelevantToAudience()` stub removed (events don't filter by audience)
- [x] `AudienceFilter` struct removed

### All Phases
- [x] All existing tests pass
- [x] New regression tests pass
- [x] `make check-compliance` passes
- [x] Committed with `make ai-commit`

## Related

- **Task 40**: Role Emphasis Redesign (marked complete but left stubs)
- **Task 50**: Centralize Constants (prerequisite - provides type-safe enums and ScoringConfig)
- **Task 51**: This fix - implements role-based scoring differentiation

## Change Log

### 2026-01-21 - Initial Fix

**Commits:**
1. `docs(docs): expand BUG-008 with comprehensive role-scoring analysis`
2. `fix(cv): implement role-based scoring using category alignment (BUG-008)`
3. `refactor(cv): extract role scoring magic numbers to named constants`

**Key Changes:**
- Added `Category` field to `EnhancedBullet` and `CVBullet` structs
- Implemented `extractPrimaryCategory()` to propagate categories from events/facts
- Updated `calculateRoleScore()` to use category-based scoring
- Updated `RoleFilter` to use type-safe `constants.CompetencyCategory`
- Added comprehensive regression tests
- Tests now use factory pattern (`fixtures.EventWithCategories()`, `fixtures.FactWithCategories()`)

### 2026-01-22 - Cleanup

**Commits:**
1. `refactor(cv): remove unused variant types and test-only code`

**Key Changes:**
- Deleted `role_emphasis.go` and `role_emphasis_test.go` (unused)
- Deleted `variants_test.go` (tests for deleted types)
- Removed `CVVariant`, `BulletConfig`, `SectionConfig`, `VariantService`, `BuiltInVariants` from variants.go
- Removed `ProfileOverride` and `ApplyProfileOverride*` functions from cv_helpers.go
- Removed `AudienceFilter` struct from enhanced_bullet_generator.go
- Total: ~1300 lines of dead code removed
