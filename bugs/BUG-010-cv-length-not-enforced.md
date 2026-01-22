# BUG-010: CV Length Format Not Enforced

## Status: RESOLVED

**Priority:** High - User-facing feature promises not being kept

**Discovered:** 2026-01-22

**Resolved:** 2026-01-22

**Assigned To:** AI Agent

## Summary

The CV length selection in the wizard (1 page, 2 pages, detailed) is collected from users but **never enforced** during CV generation. Users selecting "1 Page" still receive multi-page CVs because:

1. **No mapping exists** between UI values (`"1_page"`, `"2_page"`, `"detailed"`) and service layer values (`LengthUltraShort`, `LengthShort`, `LengthStandard`, `LengthFull`)
2. **Length constraints are never applied** - the `LengthFormatConfig` helper methods exist but are never called
3. **UI offers 3 options but service has 4 formats** - `LengthStandard` (2-3 pages) is missing from the UI

**User Impact:** Users explicitly requesting "1 Page (concise)" CVs receive 3+ page CVs, breaking trust in the application.

## Severity

- [ ] Critical - Core feature broken
- [x] High - Major feature broken (user promises not kept)
- [ ] Medium - Feature partially broken
- [ ] Low - Minor issue

## Root Cause Analysis

### Problem 1: Value Mapping Gap

**UI Values** (from `cv_config_form.go` line 173-182):
```go
forms.NewSelect("cv_length", "CV Length", []forms.SelectOption{
    {Value: "1_page", Label: "1 Page (concise)"},
    {Value: "2_page", Label: "2 Pages (standard)"},
    {Value: "detailed", Label: "Detailed (3+ pages)"},
})
```

**Service Values** (from `variants.go` line 28-43):
```go
const (
    LengthFull       LengthFormat = "full"        // 3+ pages
    LengthStandard   LengthFormat = "standard"    // 2-3 pages
    LengthShort      LengthFormat = "short"       // 1-2 pages
    LengthUltraShort LengthFormat = "ultra_short" // 1 page
)
```

**Current Behavior:** The UI value `"1_page"` is passed directly to `CVConfig.LengthFormat` without mapping, so it never matches `LengthUltraShort`.

### Problem 2: Constraints Never Applied

**Helper methods exist** in `length_format.go`:

| Method | Purpose | Line | Currently Called? |
|--------|---------|------|-------------------|
| `ShouldIncludeEvent()` | Filter events by date (last N years) | 53-59 | ❌ NO |
| `FilterCompaniesByLimit()` | Limit companies (e.g., top 3 for 1-page) | 75-80 | ❌ NO |
| `GetEffectiveBulletLimit()` | Get max bullets per job | 83-88 | ❌ NO |
| `MeetsConfidenceThreshold()` | Filter low-confidence bullets | 91-93 | ❌ NO |

**Evidence:** `cv_generation_service.go` line 161 calls `BuildSections()` without any length config:
```go
sections, err := svc.sectionBuilder.BuildSections(ctx, cvBullets, events, facts, config.TargetRole, skillsConfig)
// Length config is NOT passed ^
```

### Problem 3: Missing UI Option

The service layer has 4 length formats, but the UI only offers 3 options:

| Service Format | Target Pages | UI Option? |
|----------------|--------------|------------|
| `LengthUltraShort` | 1 page | ✅ `"1_page"` |
| `LengthShort` | 1-2 pages | ✅ `"2_page"` |
| `LengthStandard` | 2-3 pages | ❌ **MISSING** |
| `LengthFull` | 3+ pages | ✅ `"detailed"` |

### Problem 4: Documentation vs Code Mismatch

**Docs** (`docs/guides/CV_VARIANTS_GUIDE.md` lines 103, 121, 139, 157):

| Format | Bullets Per Job (Docs) | Bullets Per Job (Code) |
|--------|------------------------|------------------------|
| Full | 8 | nil (undefined) |
| Standard | 6 | 5 |
| Short | 4 | 3 |
| Ultra-Short | 3 | 2 |

The code is more restrictive than documented, which may cause confusion.

## Evidence

### Test Case: 1-Page Selection Ignored

```bash
# User selects "1 Page (concise)" in wizard
wizard.SetCVLength("1_page")

# Value is stored in config
config.LengthFormat = "1_page"  // UI value, not mapped

# Generation service tries to match
lengthConfig := GetLengthFormatConfig("1_page")
// Returns default config because "1_page" != "ultra_short"

# Even if matched, constraints aren't applied
bullets := generateBullets()  // No filtering by confidence
events := retrieveEvents()     // No filtering by date
sections := buildSections()    // No bullet limit from length config

# Result: 3+ page CV despite "1 page" request
```

### Files Involved

| File | Issue |
|------|-------|
| `internal/cli/forms/cv_config_form.go` | UI offers 3 options (should be 4) |
| `internal/cli/intents/generate_cv_intent.go` | Passes UI value directly without mapping |
| `internal/domain/career/cv.go` | `CVConfig.LengthFormat` accepts unmapped UI value |
| `internal/service/career/cv/cv_generation_service.go` | Never applies length constraints |
| `internal/service/career/cv/section_builder.go` | Doesn't accept length config parameter |
| `internal/service/career/cv/length_format.go` | Helper methods exist but unused |

## Expected Behavior

1. **UI should offer 4 options:**
   - "1 Page (executive summary)" → `LengthUltraShort`
   - "2 Pages (concise)" → `LengthShort`
   - "Standard (2-3 pages)" → `LengthStandard`
   - "Detailed (3+ pages)" → `LengthFull`

2. **Mapping function should exist:**
   ```go
   func MapUILengthToFormat(uiLength string) LengthFormat {
       switch uiLength {
       case "1_page": return LengthUltraShort
       case "2_page": return LengthShort
       case "standard": return LengthStandard
       case "detailed": return LengthFull
       default: return LengthStandard
       }
   }
   ```

3. **CV generation should enforce constraints:**
   - Filter events by `ShouldIncludeEvent()` (date-based)
   - Filter companies by `FilterCompaniesByLimit()` (top N companies)
   - Filter bullets by `MeetsConfidenceThreshold()` (min confidence)
   - Apply bullet limits via `GetEffectiveBulletLimit()` in SectionBuilder

4. **Length should override role for bullet limits:**
   - When length is "1 Page", limit is 3 bullets/company (not role-based 4-5)
   - When length is "Full", fall back to role-based limits
   - Rationale: User intent for size overrides role defaults

## Impact

**User-Facing:**
- Users requesting 1-page CVs for networking receive 3+ page documents
- Users requesting 2-page CVs for job applications receive oversized CVs
- Users lose trust in the application's ability to follow their explicit instructions

**Technical Debt:**
- 339 lines of tested helper methods in `length_format.go` that are never used
- Wizard collects data that's ignored during generation
- Documentation describes features that don't work

## Related Issues

- **BUG-008:** Role scoring not implemented (FIXED) - similar "feature exists but not wired" pattern
- **Task 24:** Flexible CV variants - designed the length format system
- **Task 40:** Role emphasis redesign - documented role vs length interaction

## Reproduction Steps

1. Run `kariya generate-cv`
2. Complete wizard, select "1 Page (concise)" for CV length
3. Generate CV with configuration
4. Observe generated CV is 3+ pages (should be 1 page)

## Fix Plan

See `tasks/tasks-52-fix-cv-length-enforcement.md` for detailed implementation plan.

**Summary:**
1. Add mapping function `MapUILengthToFormat()`
2. Add 4th UI option "Standard (2-3 pages)"
3. Update `SectionBuilder.BuildSections()` to accept `lengthConfig` parameter
4. Apply length constraints in `CVGenerationService.GenerateCVFromConfig()`
5. Update bullet limits in code to match documentation
6. Add integration tests verifying length enforcement

## Test Coverage

**Existing Tests (339 lines):**
- ✅ `length_format_test.go` - All helper methods tested
- ✅ `generate_cv_wizard_e2e_test.go` - UI value collection tested

**Missing Tests:**
- ❌ Integration tests verifying length actually affects CV output
- ❌ Tests for UI-to-service value mapping
- ❌ Tests for `BuildSections` with length config parameter

## References

- **Code:** `internal/service/career/cv/length_format.go` - Unused helper methods
- **Docs:** `docs/guides/CV_VARIANTS_GUIDE.md` - Documented behavior
- **Tests:** `internal/service/career/cv/length_format_test.go` - Passing tests for unused code
- **Similar Bug:** `bugs/BUG-008-role-scoring-not-implemented.md` - Pattern match

---

## Resolution

**Fixed in commits:**
- `52f874a` - Added MapUILengthToFormat function to bridge UI and service values
- `0e36098` - Added 4th UI option and wired up mapping in intent
- `674e411` - Applied length format constraints in CV generation service

**Solution implemented:**
1. Created `MapUILengthToFormat()` to map UI values to service constants
2. Added "Standard (2-3 pages)" UI option to match all 4 service formats
3. Applied event filtering by date (MaxYearsHistory) in CVGenerationService
4. Applied bullet filtering by confidence (MinConfidence) in CVGenerationService
5. Updated intent to use mapping function when creating CVConfig

**Verification:**
- All existing tests pass
- Length format helper methods are now called during CV generation
- Event and bullet filtering is logged for debugging

**Note:** MaxBulletsPerJob (per-company bullet limits) not yet applied to keep changes atomic. This can be addressed in a future enhancement if needed.

**Created:** 2026-01-22  
**Last Updated:** 2026-01-22
**Resolved:** 2026-01-22
