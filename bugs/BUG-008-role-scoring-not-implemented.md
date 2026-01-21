# BUG-008: Role-Based CV Differentiation System Not Implemented

## Summary

The CV generation system has a significant architectural gap where role and audience selection provides minimal actual differentiation in bullet scoring and selection. The system was designed to prioritize different achievements based on target role (Principal vs Staff vs Senior IC vs EM), but the scoring logic was never implemented.

**The result:** A "Principal" CV and "Senior IC" CV produce nearly identical bullet rankings, defeating the purpose of role selection.

## Severity

- [x] Critical - Core feature broken (role differentiation is the primary value proposition)
- [ ] High - Major feature broken
- [ ] Medium - Feature partially broken
- [ ] Low - Minor issue/cosmetic

## Prerequisites

**Task-50 (Centralize Constants)** must be merged before starting this fix. It provides:

| Component | Location | Use For |
|-----------|----------|---------|
| `constants.RoleFit` | `internal/constants/constants.go` | Type-safe role enums (`RoleFitPrincipal`, `RoleFitSeniorIC`, etc.) |
| `constants.CompetencyCategory` | `internal/constants/constants.go` | Type-safe category enums (`CompetencyTechnical`, `CompetencyLeadership`, etc.) |
| `constants.Audience` | `internal/constants/constants.go` | Type-safe audience enums (`AudienceHiringManager`, etc.) |
| `constants.InclusionReason` | `internal/constants/constants.go` | Type-safe inclusion reasons |
| `config.ScoringConfig` | `internal/config/config.go` | Configurable scoring weights and thresholds |
| `config.RoleScoringCfg` | `internal/config/config.go` | Per-role MinConfidence and MaxBulletsPerCompany |

**What Task-50 does NOT provide (this bug must implement):**
- Category-based role scoring weights (e.g., "technical" +0.20 for senior_ic)
- The `Category` field on bullet structs
- Wiring up category scoring in `calculateRoleScore()`

## How It Should Work

### The Core Concept

The same career history should produce **different CVs** based on target role:

| Role | Should Prioritize | Should De-prioritize |
|------|-------------------|---------------------|
| **Senior IC** | Hands-on technical work, implementation | Leadership, strategy |
| **Staff** | Technical leadership + deep technical work | Pure management |
| **Principal** | Strategy, architecture, cross-team influence | Hands-on implementation |
| **EM** | People management, team outcomes, mentoring | Deep technical work |

### Example: Same Event, Different Scores

```
Event: "Led migration of 3 microservices to Kubernetes"
Categories: [technical, leadership]

Expected Role Scores:
- Senior IC:  0.70 + 0.20(technical) - 0.05(leadership) = 0.85
- Staff:      0.70 + 0.15(technical) + 0.15(leadership) = 1.00
- Principal:  0.70 + 0.10(technical) + 0.20(leadership) = 1.00
- EM:         0.70 + 0.00(technical) + 0.20(leadership) = 0.90

Actual Role Scores (current broken behavior):
- Senior IC:  0.50 + 0.30(achievement) + 0.10(confidence) = 0.90
- Staff:      0.50 + 0.30(achievement) + 0.10(confidence) = 0.90
- Principal:  0.50 + 0.30(achievement) + 0.10(confidence) = 0.90
- EM:         0.50 + 0.30(achievement) + 0.10(confidence) = 0.90
              ↑ ALL IDENTICAL - role is ignored
```

## Root Cause Analysis

### Problem 1: Category Data Doesn't Flow Through Pipeline

Categories exist in source data but are lost when bullets are created:

| Struct | Has Category? | Used in Scoring? |
|--------|---------------|------------------|
| `CareerEvent.Categories` | ✅ YES | Filtering only (not scoring) |
| `Fact.CompetencyCategories` | ✅ YES | ❌ NOT USED |
| `Achievement` | ❌ **NO FIELD** | N/A |
| `EnhancedBullet` | ❌ **NO FIELD** | N/A |
| `CVBullet` | ❌ **NO FIELD** | N/A |

### Problem 2: `calculateRoleScore()` Ignores Role Parameter

**File:** `internal/service/career/cv/enhanced_bullet_generator.go` (lines 452-468)

```go
func (ebg *DefaultEnhancedBulletGenerator) calculateRoleScore(
    bullet *EnhancedBullet,
    role string,  // <-- RECEIVED BUT NEVER USED
) float64 {
    score := 0.5                                           // Base score
    if bullet.InclusionReason == "achievement_extraction" { score += 0.3 }
    if bullet.InclusionReason == "fact_extraction"        { score += 0.2 }
    if bullet.Confidence > 0.8                            { score += 0.1 }
    return math.Min(score, 1.0)
    // NO category-based scoring
    // NO role-based differentiation
}
```

### Problem 3: `getRoleFilter()` Categories Never Used

**File:** `internal/service/career/cv/enhanced_bullet_generator.go` (lines 609-645)

The function correctly defines `PrimaryCategories` and `SecondaryCategories` for each role:

```go
case "principal":
    return &RoleFilter{
        PrimaryCategories:   []string{"leadership", "strategy", "architecture"},
        SecondaryCategories: []string{"technical", "mentoring"},
        MinConfidence:       0.8,
    }
case "senior_ic":
    return &RoleFilter{
        PrimaryCategories:   []string{"technical", "architecture"},
        SecondaryCategories: []string{"leadership", "strategy"},
        MinConfidence:       0.75,
    }
```

**But these categories are never used for scoring** - only `MinConfidence` is applied.

### Problem 4: Event Audience Filtering Is a Stub

**File:** `internal/service/career/cv/bullet_generator.go` (lines 262-271)

```go
func (bg *DefaultBulletGenerator) isEventRelevantToAudience(event *career.CareerEvent, audience string) bool {
    if audience == "" {
        return true
    }
    // For now, accept all events for all audiences
    // In future, could implement audience-specific filtering
    return true  // <-- STUB: always returns true
}
```

### Problem 5: Unused Systems (For Reference)

The codebase contains additional systems that were designed but never connected:

| System | Location | Status |
|--------|----------|--------|
| `CVVariant` (16 variants) | `variants.go` | Defined but UI uses `CVProfile` instead |
| `RoleEmphasis` | `variants.go` | Different taxonomy (`senior_backend`, `staff_principal`, etc.) - not needed |
| `RoleEmphasisConfig.ScoreBulletCategory()` | `role_emphasis.go` | Working method but never called |
| `AudienceFilter` struct | `enhanced_bullet_generator.go` | Defined but never instantiated |

**Recommendation:** These can be deprecated in future cleanup. Focus fix on the core `TargetRole` system.

## Files Affected

| File | Lines | Issue |
|------|-------|-------|
| `internal/service/career/cv/enhanced_bullet_generator.go` | 38-56 | `EnhancedBullet` missing `Category` field |
| `internal/service/career/cv/enhanced_bullet_generator.go` | 452-468 | `calculateRoleScore()` ignores role |
| `internal/service/career/cv/enhanced_bullet_generator.go` | 602-607 | `customizeForRole()` is stub |
| `internal/service/career/cv/enhanced_bullet_generator.go` | 609-645 | `getRoleFilter()` categories unused |
| `internal/service/career/cv/bullet_generator.go` | 262-271 | `isEventRelevantToAudience()` is stub |
| `internal/service/career/cv/data_processing_service.go` | ~73 | `Achievement` missing `Category` field |
| `internal/domain/career/cv.go` | ~254 | `CVBullet` missing `Category` field |

### Files to Use (from Task-50)

| File | Purpose |
|------|---------|
| `internal/constants/constants.go` | Use `RoleFit`, `CompetencyCategory`, `Audience` types |
| `internal/config/config.go` | Use `ScoringConfig.Weights` for final score calculation |

## Fix Approach

### Phase 1: Add Category Field Propagation

1. Add `Category constants.CompetencyCategory` field to `Achievement` struct
2. Add `Category constants.CompetencyCategory` field to `EnhancedBullet` struct  
3. Add `Category constants.CompetencyCategory` field to `CVBullet` struct
4. Update `InclusionReason` field to use `constants.InclusionReason` type
5. Populate category from source:
   - From `CareerEvent.Categories[0]` → convert to `constants.CompetencyCategory`
   - From `Fact.CompetencyCategories[0]` → convert to `constants.CompetencyCategory`
6. Use `constants.IsValidCompetencyCategory()` for validation

**Example struct update:**
```go
// EnhancedBullet with type-safe fields
type EnhancedBullet struct {
    ID              string
    Text            string
    EnhancedText    string
    Category        constants.CompetencyCategory  // NEW: type-safe category
    InclusionReason constants.InclusionReason     // UPDATED: type-safe
    // ... other fields
}
```

### Phase 2: Implement Role-Based Scoring

Update `calculateRoleScore()` to use category + role with type-safe constants:

```go
import (
    "github.com/baphled/kariya/internal/constants"
    "github.com/baphled/kariya/internal/config"
)

func (ebg *DefaultEnhancedBulletGenerator) calculateRoleScore(
    bullet *EnhancedBullet,
    role constants.RoleFit,
) float64 {
    baseScore := 0.5
    
    // Get role filter with category definitions
    filter := ebg.getRoleFilter(role)
    
    // Apply category-based scoring using type-safe constants
    if filter != nil && bullet.Category != "" {
        // Check primary categories (strongest boost)
        for _, primary := range filter.PrimaryCategories {
            if bullet.Category == primary {
                baseScore += 0.30
                break
            }
        }
        // Check secondary categories (medium boost)
        for _, secondary := range filter.SecondaryCategories {
            if bullet.Category == secondary {
                baseScore += 0.15
                break
            }
        }
    }
    
    // Existing bonuses using type-safe inclusion reason constants
    if bullet.InclusionReason == constants.InclusionReasonAchievementExtraction {
        baseScore += 0.10
    }
    
    // Use config threshold for high confidence check
    cfg := config.DefaultConfig()
    if bullet.Confidence > cfg.Scoring.Thresholds.HighConfidence {
        baseScore += 0.10
    }
    
    return math.Min(baseScore, 1.0)
}

// Update RoleFilter to use type-safe constants
type RoleFilter struct {
    PrimaryCategories   []constants.CompetencyCategory
    SecondaryCategories []constants.CompetencyCategory
    MinConfidence       float64
    PreferredMetrics    []string
}

// Update getRoleFilter to return type-safe categories
func (ebg *DefaultEnhancedBulletGenerator) getRoleFilter(role constants.RoleFit) *RoleFilter {
    switch role {
    case constants.RoleFitPrincipal:
        return &RoleFilter{
            PrimaryCategories:   []constants.CompetencyCategory{
                constants.CompetencyLeadership,
                constants.CompetencyConsulting,  // strategy mapped to consulting
            },
            SecondaryCategories: []constants.CompetencyCategory{
                constants.CompetencyTechnical,
                constants.CompetencyMentoring,
            },
            MinConfidence: 0.8,
        }
    case constants.RoleFitSeniorIC:
        return &RoleFilter{
            PrimaryCategories:   []constants.CompetencyCategory{
                constants.CompetencyTechnical,
                constants.CompetencyResearch,
            },
            SecondaryCategories: []constants.CompetencyCategory{
                constants.CompetencyLeadership,
                constants.CompetencyProduct,
            },
            MinConfidence: 0.75,
        }
    // ... other roles
    default:
        return &RoleFilter{MinConfidence: 0.7}
    }
}
```

**Note:** The `config.ScoringConfig.Weights` should be used in `calculateFinalScore()` for the weighted combination:

```go
func (ebg *DefaultEnhancedBulletGenerator) calculateFinalScore(bullet *EnhancedBullet) float64 {
    cfg := config.DefaultConfig()
    w := cfg.Scoring.Weights
    
    return math.Min(
        (w.RoleScore * bullet.RoleScore) +
        (w.AudienceScore * bullet.AudienceScore) +
        (w.MetricScore * bullet.MetricScore) +
        (w.ImpactScore * bullet.ImpactScore) +
        (w.Confidence * bullet.Confidence),
        1.0,
    )
}
```

### Phase 3: Fix Event Audience Filtering

Implement `isEventRelevantToAudience()` using type-safe constants:

```go
func (bg *DefaultBulletGenerator) isEventRelevantToAudience(
    event *career.CareerEvent,
    audience constants.Audience,
) bool {
    if audience == "" {
        return true
    }
    // Implement actual filtering logic here
    // For now, could check event tags or categories for audience relevance
    return true
}
```

## Regression Tests

```go
import (
    "github.com/baphled/kariya/internal/constants"
)

var _ = Describe("BUG-008: Role-based CV differentiation", func() {
    Describe("Category propagation", func() {
        It("should propagate category from CareerEvent to EnhancedBullet", func() {
            // Note: CareerEvent.Categories uses []string for backward compatibility
            // but validates against constants.IsValidCompetencyCategory()
            event := &career.CareerEvent{
                ID:         "evt-1",
                Text:       "Led team migration",
                Categories: []string{"leadership", "technical"},
            }
            bullets := generator.createBulletsFromEvents([]*career.CareerEvent{event})
            // EnhancedBullet.Category uses the type-safe constant type
            Expect(bullets[0].Category).To(Equal(constants.CompetencyLeadership))
        })
        
        It("should propagate category from Fact to EnhancedBullet", func() {
            fact := &career.Fact{
                ID:                   "fact-1",
                Text:                 "Mentored junior engineers",
                CompetencyCategories: []string{"mentoring"},
            }
            bullets := generator.createBulletsFromFacts([]*career.Fact{fact})
            Expect(bullets[0].Category).To(Equal(constants.CompetencyMentoring))
        })
    })

    Describe("Role-based scoring", func() {
        It("should score leadership bullets higher for principal than senior_ic", func() {
            bullet := &EnhancedBullet{
                Category:   constants.CompetencyLeadership,
                Confidence: 0.7,
            }
            
            principalScore := generator.calculateRoleScore(bullet, constants.RoleFitPrincipal)
            seniorScore := generator.calculateRoleScore(bullet, constants.RoleFitSeniorIC)
            
            Expect(principalScore).To(BeNumerically(">", seniorScore))
        })

        It("should score technical bullets higher for senior_ic than principal", func() {
            bullet := &EnhancedBullet{
                Category:   constants.CompetencyTechnical,
                Confidence: 0.7,
            }
            
            seniorScore := generator.calculateRoleScore(bullet, constants.RoleFitSeniorIC)
            principalScore := generator.calculateRoleScore(bullet, constants.RoleFitPrincipal)
            
            Expect(seniorScore).To(BeNumerically(">", principalScore))
        })

        It("should score mentoring bullets higher for em than senior_ic", func() {
            bullet := &EnhancedBullet{
                Category:   constants.CompetencyMentoring,
                Confidence: 0.7,
            }
            
            emScore := generator.calculateRoleScore(bullet, constants.RoleFitEM)
            seniorScore := generator.calculateRoleScore(bullet, constants.RoleFitSeniorIC)
            
            Expect(emScore).To(BeNumerically(">", seniorScore))
        })
        
        It("should produce different final rankings for different roles", func() {
            events := []*career.CareerEvent{
                {ID: "1", Text: "Built pipeline", Categories: []string{"technical"}},
                {ID: "2", Text: "Led team", Categories: []string{"leadership"}},
                {ID: "3", Text: "Mentored juniors", Categories: []string{"mentoring"}},
            }
            
            seniorBullets, _ := generator.GenerateBullets(ctx, events, nil, nil,
                constants.RoleFitSeniorIC, "")
            principalBullets, _ := generator.GenerateBullets(ctx, events, nil, nil,
                constants.RoleFitPrincipal, "")
            
            // Technical bullet should rank higher for senior_ic
            // Leadership bullet should rank higher for principal
            Expect(seniorBullets[0].Category).To(Equal(constants.CompetencyTechnical))
            Expect(principalBullets[0].Category).To(Equal(constants.CompetencyLeadership))
        })
    })
    
    Describe("Scoring config integration", func() {
        It("should use config weights for final score calculation", func() {
            cfg := config.DefaultConfig()
            
            // Verify weights sum to 1.0
            Expect(cfg.Scoring.ValidateWeights()).To(Succeed())
            
            // Verify role settings exist
            Expect(cfg.Scoring.RoleSettings).To(HaveKey("principal"))
            Expect(cfg.Scoring.RoleSettings).To(HaveKey("senior_ic"))
            Expect(cfg.Scoring.RoleSettings["principal"].MinConfidence).To(
                BeNumerically(">", cfg.Scoring.RoleSettings["senior_ic"].MinConfidence))
        })
    })
})
```

## Definition of Done

### Phase 1: Category Propagation
- [ ] `Achievement` struct has `Category constants.CompetencyCategory` field
- [ ] `EnhancedBullet` struct has `Category constants.CompetencyCategory` field
- [ ] `CVBullet` struct has `Category constants.CompetencyCategory` field
- [ ] Category populated from `CareerEvent.Categories[0]`
- [ ] Category populated from `Fact.CompetencyCategories[0]`
- [ ] Use `constants.IsValidCompetencyCategory()` for validation
- [ ] Unit tests for category propagation

### Phase 2: Role Scoring
- [ ] `calculateRoleScore()` accepts `constants.RoleFit` parameter
- [ ] `calculateRoleScore()` uses bullet category for scoring
- [ ] `calculateRoleScore()` applies role-specific weights from `getRoleFilter()`
- [ ] `calculateFinalScore()` uses `config.ScoringConfig.Weights`
- [ ] Use `config.ScoringConfig.RoleSettings` for MinConfidence and MaxBulletsPerCompany
- [ ] Different roles produce measurably different bullet rankings
- [ ] Unit tests for role-based scoring
- [ ] Integration test comparing Senior IC vs Principal output

### Phase 3: Audience Filtering (Optional)
- [ ] `isEventRelevantToAudience()` accepts `constants.Audience` parameter
- [ ] `isEventRelevantToAudience()` implemented (not stub)
- [ ] Or: Remove stub and document that events don't filter by audience

### All Phases
- [ ] All existing tests pass
- [ ] New regression tests pass
- [ ] `make check-compliance` passes
- [ ] Committed with `make ai-commit`

## Related

- **Task 40**: Role Emphasis Redesign (marked complete but left stubs)
- **Task 50**: Centralize Constants (prerequisite - provides type-safe enums and ScoringConfig)
- **BUG-009**: Category Mismatch (related - wrong categories in filters)

## Investigation Session

**Date**: 2026-01-21
**Session**: Role-Based CV Differentiation Deep Dive

### Key Discoveries

1. **Category data exists but doesn't flow** - `CareerEvent.Categories` and `Fact.CompetencyCategories` exist but aren't propagated to bullet structs
2. **`calculateRoleScore()` ignores role parameter** - Returns same score regardless of role
3. **`getRoleFilter()` categories are correct but unused** - The category definitions are good, just never applied to scoring
4. **Event audience filtering is a stub** - Always returns true
5. **Unused parallel systems exist** - `CVVariant`, `RoleEmphasis`, `RoleEmphasisConfig` were designed but never connected (can be deprecated)

### Historical Context

Task 40 (Role Emphasis Redesign) appears to have created the UI workflow and configuration structures, but the actual scoring implementation was left as stubs with TODO comments.

### Simplified Scope

Focus on the **seniority-based role taxonomy** (`principal`, `staff`, `em`, `senior_ic`) which is what the UI actually uses. The `RoleEmphasis` system (`senior_backend`, `staff_principal`, `consulting`, `language_agnostic`) uses a different taxonomy and can be deprecated.

### Task-50 Integration (2026-01-21)

Reviewed `feature/task-50-centralize-constants` branch which provides:

1. **Type-safe constants** in `internal/constants/constants.go`:
   - `RoleFit` type with `RoleFitPrincipal`, `RoleFitStaff`, `RoleFitEM`, `RoleFitSeniorIC`
   - `CompetencyCategory` type with `CompetencyTechnical`, `CompetencyLeadership`, etc.
   - `Audience` type with `AudienceHiringManager`, `AudienceRecruiter`, `AudiencePeer`
   - Validation functions like `IsValidRoleFit()`, `IsValidCompetencyCategory()`

2. **ScoringConfig** in `internal/config/config.go`:
   - `ScoringWeights`: RoleScore (0.25), AudienceScore (0.20), MetricScore (0.20), ImpactScore (0.20), Confidence (0.15)
   - `ScoringThresholds`: FactDefaultConfidence (0.85), EventDefaultConfidence (0.80), etc.
   - `RoleSettings`: Per-role MinConfidence and MaxBulletsPerCompany

**Key insight:** Task-50 provides the infrastructure but NOT the category-based role scoring. This bug must implement the actual scoring logic that uses categories to differentiate roles.
