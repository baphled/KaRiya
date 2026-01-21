# BUG-008: Role Scoring Not Implemented in Enhanced Bullet Generator

## Summary

The `calculateRoleScore()` function in `enhanced_bullet_generator.go` receives a `role` parameter but never uses it, causing all CV roles (senior_ic, staff, principal, em) to produce identical bullet scoring and selection.

## Steps to Reproduce

1. Generate CV with role = "senior_ic" (Senior Engineer)
2. Generate CV with role = "staff" (Staff Engineer)
3. Compare bullet selection between the two CVs
4. Observe that bullet rankings are nearly identical

## Expected Behavior

Different roles should prioritize different types of achievements:

| Role | Primary Focus | Secondary Focus |
|------|---------------|-----------------|
| senior_ic | Hands-on technical delivery, implementation | Code quality, debugging |
| staff | Technical leadership, architecture decisions | Mentoring, cross-team work |
| principal | Strategy, cross-team influence | Technical vision, standards |
| em | People management, team building | Process improvement, hiring |

## Actual Behavior

All roles produce the same bullet rankings because:

1. `calculateRoleScore()` (lines 452-468) receives `role` but **never uses it**
2. `customizeForRole()` (lines 602-607) is a stub that returns text unchanged
3. `getRoleFilter()` defines `PrimaryCategories`/`SecondaryCategories` but these are **never used** for scoring

### Code Evidence

```go
// enhanced_bullet_generator.go:452-468
func (g *EnhancedBulletGenerator) calculateRoleScore(
    bullet *Bullet,
    role career.TargetRole,  // <-- Parameter received but NEVER USED
    filter *RoleFilter,
) float64 {
    // ... scoring logic that ignores role entirely ...
}
```

```go
// enhanced_bullet_generator.go:602-607
func (g *EnhancedBulletGenerator) customizeForRole(text string, role career.TargetRole) string {
    // Future: apply role-specific wording preferences
    return text  // <-- Stub: returns unchanged
}
```

## Root Cause

Task 40 (Role Emphasis Redesign) was marked complete on 2026-01-14, but the actual role differentiation logic was left as stubs with TODO comments. The UI and workflow were implemented, but the underlying scoring and filtering don't differentiate between roles.

## Files Affected

| File | Lines | Issue |
|------|-------|-------|
| `internal/service/career/cv/enhanced_bullet_generator.go` | 452-468 | `calculateRoleScore()` ignores role |
| `internal/service/career/cv/enhanced_bullet_generator.go` | 602-607 | `customizeForRole()` is stub |
| `internal/service/career/cv/bullet_generator.go` | varies | `FilterByRole()` only filters by MinConfidence |

## Fix Approach

Implement actual role-based scoring in `calculateRoleScore()`:

```go
func (g *EnhancedBulletGenerator) calculateRoleScore(
    bullet *Bullet,
    role career.TargetRole,
    filter *RoleFilter,
) float64 {
    baseScore := bullet.Confidence
    
    // Apply role-specific category weights
    categoryBonus := g.getRoleCategoryBonus(bullet.Category, role)
    
    return math.Min(1.0, baseScore + categoryBonus)
}

func (g *EnhancedBulletGenerator) getRoleCategoryBonus(category string, role career.TargetRole) float64 {
    weights := map[career.TargetRole]map[string]float64{
        career.TargetRoleSeniorIC: {
            "technical":  +0.20,
            "leadership": -0.05,
            "mentoring":  +0.05,
        },
        career.TargetRoleStaff: {
            "technical":  +0.15,
            "leadership": +0.15,
            "mentoring":  +0.10,
        },
        career.TargetRolePrincipal: {
            "technical":  +0.10,
            "leadership": +0.20,
            "mentoring":  +0.05,
        },
        career.TargetRoleEM: {
            "technical":  +0.00,
            "leadership": +0.20,
            "mentoring":  +0.15,
        },
    }
    
    if roleWeights, ok := weights[role]; ok {
        if bonus, ok := roleWeights[category]; ok {
            return bonus
        }
    }
    return 0.0
}
```

## Severity

- [ ] Critical - Application crash/data loss
- [x] High - Major feature broken
- [ ] Medium - Feature partially broken
- [ ] Low - Minor issue/cosmetic

## Regression Test

```go
var _ = Describe("BUG-008: Role scoring differentiation", func() {
    var generator *EnhancedBulletGenerator
    
    BeforeEach(func() {
        generator = NewEnhancedBulletGenerator(/* deps */)
    })
    
    Describe("calculateRoleScore", func() {
        It("should score technical bullets higher for senior_ic than staff", func() {
            bullet := &Bullet{
                Category:   "technical",
                Confidence: 0.70,
            }
            
            seniorScore := generator.calculateRoleScore(bullet, career.TargetRoleSeniorIC, nil)
            staffScore := generator.calculateRoleScore(bullet, career.TargetRoleStaff, nil)
            
            Expect(seniorScore).To(BeNumerically(">", staffScore))
        })
        
        It("should score leadership bullets higher for staff than senior_ic", func() {
            bullet := &Bullet{
                Category:   "leadership",
                Confidence: 0.70,
            }
            
            staffScore := generator.calculateRoleScore(bullet, career.TargetRoleStaff, nil)
            seniorScore := generator.calculateRoleScore(bullet, career.TargetRoleSeniorIC, nil)
            
            Expect(staffScore).To(BeNumerically(">", seniorScore))
        })
        
        It("should score mentoring bullets higher for EM than senior_ic", func() {
            bullet := &Bullet{
                Category:   "mentoring",
                Confidence: 0.70,
            }
            
            emScore := generator.calculateRoleScore(bullet, career.TargetRoleEM, nil)
            seniorScore := generator.calculateRoleScore(bullet, career.TargetRoleSeniorIC, nil)
            
            Expect(emScore).To(BeNumerically(">", seniorScore))
        })
    })
})
```

## Definition of Done

- [ ] Root cause confirmed via code inspection
- [ ] Regression test written FIRST (TDD)
- [ ] `calculateRoleScore()` uses role parameter
- [ ] `customizeForRole()` implemented or removed
- [ ] All existing tests pass
- [ ] New role differentiation tests pass
- [ ] `make check-compliance` passes
- [ ] Committed with `make ai-commit`

## Related

- **Task 40**: Role Emphasis Redesign (marked complete but left stubs)
- **BUG-009**: Category Mismatch (related - wrong categories in filters)

## Notes

This bug explains why CVs "look like a staff engineer" regardless of role selection. The entire role-based filtering system is effectively a no-op because the scoring function ignores the role parameter.

## Investigation Session

**Date**: 2026-01-21
**Session**: Technology-Focused CV Generation Investigation

Found during analysis of why senior_ic and staff CVs produce nearly identical output.
