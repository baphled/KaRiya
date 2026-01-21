# BUG-009: Category Mismatch Between Code and Data

## Summary

The role filtering code in `enhanced_bullet_generator.go` expects event categories (`architecture`, `strategy`) that don't exist in the database, while the actual data uses different categories (`technical`, `leadership`, etc.), making role-based filtering ineffective.

## Steps to Reproduce

1. Query database categories:
   ```sql
   sqlite3 ~/.kariya/events.db "SELECT category, COUNT(*) FROM events GROUP BY category"
   ```
2. Inspect `getRoleFilter()` in `enhanced_bullet_generator.go`
3. Compare expected vs actual categories

## Expected Behavior

Role filters should reference categories that actually exist in the event data, allowing meaningful filtering and prioritization.

## Actual Behavior

### Database Categories (Actual Data)

```
technical:   199 events
leadership:   39 events
mentoring:    21 events
product:      ~30 events
research:     ~26 events
consulting:   ~6 events
```

### Code Expectations (getRoleFilter)

```go
// For staff role:
PrimaryCategories: []string{"technical", "architecture"}  // architecture = 0 events!
SecondaryCategories: []string{"leadership", "strategy"}   // strategy = 0 events!
```

### Result

- `architecture`: **0 events** - filter expects this but it doesn't exist
- `strategy`: **0 events** - filter expects this but it doesn't exist
- Role filters effectively match only `technical` and `leadership`
- All roles end up with similar event pools

## Root Cause

The category names in `getRoleFilter()` were designed with expected categories that were never populated in the actual event data. Either:

1. The categorization schema changed after the code was written
2. Events were imported with a different categorization scheme
3. The code was written speculatively without verifying actual data

## Files Affected

| File | Issue |
|------|-------|
| `internal/service/career/cv/enhanced_bullet_generator.go` | `getRoleFilter()` references non-existent categories |
| `~/.kariya/events.db` | Events use different category names |

## Current getRoleFilter Implementation

```go
func getRoleFilter(role career.TargetRole) *RoleFilter {
    switch role {
    case career.TargetRoleSeniorIC:
        return &RoleFilter{
            PrimaryCategories:   []string{"technical", "architecture"},  // architecture = 0!
            SecondaryCategories: []string{"implementation", "debugging"},
            MinConfidence:       0.65,
        }
    case career.TargetRoleStaff:
        return &RoleFilter{
            PrimaryCategories:   []string{"technical", "leadership"},
            SecondaryCategories: []string{"architecture", "mentoring"},  // architecture = 0!
            MinConfidence:       0.60,
        }
    case career.TargetRolePrincipal:
        return &RoleFilter{
            PrimaryCategories:   []string{"leadership", "strategy"},     // strategy = 0!
            SecondaryCategories: []string{"architecture", "technical"},  // architecture = 0!
            MinConfidence:       0.55,
        }
    // ...
    }
}
```

## Fix Approach

### Option A: Update Code to Use Existing Categories (Recommended)

Update `getRoleFilter()` to use categories that exist in the data:

```go
func getRoleFilter(role career.TargetRole) *RoleFilter {
    switch role {
    case career.TargetRoleSeniorIC:
        return &RoleFilter{
            PrimaryCategories:   []string{"technical"},
            SecondaryCategories: []string{"research"},
            ExcludeCategories:   []string{"consulting"},
            MinConfidence:       0.65,
        }
    case career.TargetRoleStaff:
        return &RoleFilter{
            PrimaryCategories:   []string{"technical", "leadership"},
            SecondaryCategories: []string{"mentoring"},
            MinConfidence:       0.60,
        }
    case career.TargetRolePrincipal:
        return &RoleFilter{
            PrimaryCategories:   []string{"leadership"},
            SecondaryCategories: []string{"technical", "product"},
            MinConfidence:       0.55,
        }
    case career.TargetRoleEM:
        return &RoleFilter{
            PrimaryCategories:   []string{"leadership", "mentoring"},
            SecondaryCategories: []string{"product"},
            ExcludeCategories:   []string{"research"},
            MinConfidence:       0.50,
        }
    }
}
```

### Option B: Re-categorize Events (More Effort)

Add missing categories to relevant events:
- Tag complex technical events as `architecture`
- Tag leadership events about direction/vision as `strategy`

**Not recommended** - requires data migration and affects existing workflows.

### Option C: Add Category Mapping Layer

Create a mapping from expected categories to actual categories:

```go
var categoryMapping = map[string][]string{
    "architecture": {"technical"},  // Map architecture -> technical
    "strategy":     {"leadership"}, // Map strategy -> leadership
}
```

**Acceptable** - but adds complexity without fixing root issue.

## Severity

- [ ] Critical - Application crash/data loss
- [x] High - Major feature broken
- [ ] Medium - Feature partially broken
- [ ] Low - Minor issue/cosmetic

## Regression Test

```go
var _ = Describe("BUG-009: Category alignment", func() {
    Describe("getRoleFilter", func() {
        It("should only reference categories that exist in event data", func() {
            validCategories := []string{
                "technical", "leadership", "mentoring", 
                "product", "research", "consulting",
            }
            
            roles := []career.TargetRole{
                career.TargetRoleSeniorIC,
                career.TargetRoleStaff,
                career.TargetRolePrincipal,
                career.TargetRoleEM,
            }
            
            for _, role := range roles {
                filter := getRoleFilter(role)
                
                for _, cat := range filter.PrimaryCategories {
                    Expect(validCategories).To(ContainElement(cat),
                        "Role %s has invalid primary category: %s", role, cat)
                }
                
                for _, cat := range filter.SecondaryCategories {
                    Expect(validCategories).To(ContainElement(cat),
                        "Role %s has invalid secondary category: %s", role, cat)
                }
            }
        })
        
        It("should NOT reference 'architecture' category (0 events exist)", func() {
            for _, role := range allRoles {
                filter := getRoleFilter(role)
                allCats := append(filter.PrimaryCategories, filter.SecondaryCategories...)
                Expect(allCats).NotTo(ContainElement("architecture"))
            }
        })
        
        It("should NOT reference 'strategy' category (0 events exist)", func() {
            for _, role := range allRoles {
                filter := getRoleFilter(role)
                allCats := append(filter.PrimaryCategories, filter.SecondaryCategories...)
                Expect(allCats).NotTo(ContainElement("strategy"))
            }
        })
    })
})
```

## Definition of Done

- [ ] Audit actual categories in production database
- [ ] Regression test written FIRST (TDD)
- [ ] `getRoleFilter()` updated to use existing categories
- [ ] All existing tests pass
- [ ] New category alignment tests pass
- [ ] `make check-compliance` passes
- [ ] Committed with `make ai-commit`

## Related

- **BUG-008**: Role Scoring Not Implemented (scoring ignores role entirely)
- **Task 40**: Role Emphasis Redesign (defined the current filter structure)

## Notes

This bug compounds with BUG-008. Even if role scoring were implemented correctly, the category mismatch means the filters would still be ineffective because they reference non-existent categories.

The combination of these two bugs means:
1. Role scoring ignores the role parameter (BUG-008)
2. Category filters reference non-existent categories (BUG-009)
3. Result: All roles produce nearly identical CV output

## Database Query for Verification

```sql
-- Run this to see actual category distribution
SELECT 
    category,
    COUNT(*) as event_count,
    ROUND(COUNT(*) * 100.0 / (SELECT COUNT(*) FROM events), 1) as percentage
FROM events 
GROUP BY category 
ORDER BY event_count DESC;

-- Verify architecture and strategy don't exist
SELECT COUNT(*) FROM events WHERE category = 'architecture';  -- Expected: 0
SELECT COUNT(*) FROM events WHERE category = 'strategy';      -- Expected: 0
```

## Investigation Session

**Date**: 2026-01-21
**Session**: Technology-Focused CV Generation Investigation

Found during analysis of why senior_ic and staff CVs produce nearly identical output.
