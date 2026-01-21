# TASK-50: Centralize constants and configurable scoring

## Summary

Create a central constants package (`internal/constants/`) for all hardcoded values (roles, categories, audiences, tags, skill levels, focus areas) and make scoring weights/thresholds runtime configurable via `config.yaml`. This ensures a single source of truth throughout the application.

## Acceptance Criteria

- [x] Create `internal/constants/constants.go` with typed constants for roles, audiences, competency categories, skill categories, skill levels, event tags, focus areas, CV section types, and impact levels
- [x] Provide helper functions: `All*()`, `IsValid*()`, `CVTargetRoles()`, `SuggestedSkillCategories()`
- [x] Add `ScoringConfig` to `internal/config/config.go` with bullet score weights, confidence thresholds, and role-specific settings
- [x] Auto-migrate existing configs (missing scoring section uses defaults)
- [x] Update domain layer to use constants: `event.go`, `fact.go`, `cv.go`, `skill.go`
- [x] Update service layer to use constants and config for scoring
- [x] Update CLI layer to use constants (including skill_form.go)
- [x] Remove all duplicate constant definitions (no deprecation period)
- [x] Fix "mobile" category discrepancy between domain and forms
- [x] Update all test files and fixtures

## Technical Notes

### Files to Modify

- `internal/constants/constants.go` - New file with all type definitions and constants
- `internal/constants/constants_test.go` - New file with validation tests
- `internal/config/config.go` - Add ScoringConfig with auto-migration
- `internal/domain/career/event.go` - Remove AllowedCategories, AllowedTags; use constants
- `internal/domain/career/fact.go` - Remove RoleFit, AllowedRoleFits, AllowedAudienceRelevance; use constants
- `internal/domain/career/cv.go` - Remove all Allowed* maps; use constants
- `internal/domain/career/skill.go` - Remove CommonSkillCategories, AllowedSkillLevels; use constants
- `internal/service/career/classification/classifier.go` - Use constants.CompetencyCategory
- `internal/service/career/technology/focus_area.go` - Use constants.FocusArea
- `internal/service/career/cv/enhanced_bullet_generator.go` - Use config for weights/thresholds
- `internal/service/career/cv/bullet_generator.go` - Use config for scoring
- `internal/service/career/cv/section_builder.go` - Use config for role settings
- `internal/service/career/cv/length_format.go` - Use config for thresholds
- `internal/service/career/cv/variants.go` - Use config for variant thresholds
- `internal/service/career/cv/role_emphasis.go` - Use config for role categories
- `internal/service/career/burst_fact/classifier.go` - Use constants.Role
- `internal/service/career/burst_fact/detector.go` - Use config for min confidence
- `internal/cli/forms/skill_form.go` - Use constants
- `internal/cli/forms/fact_form.go` - Use constants
- `internal/cli/forms/capture_event_form.go` - Use constants
- `internal/cli/components/category_selector.go` - Remove duplicate; use constants
- `internal/cli/components/tag_selector.go` - Use constants
- `internal/cli/intents/configure_system.go` - Use constants
- `internal/cli/intents/generate_cv_intent.go` - Use constants
- `internal/cli/intents/manage_skills_intent.go` - Use constants
- `internal/cli/screens/skills/detail.go` - Use constants for color mapping
- `internal/cli/importer/mapper.go` - Use constants
- `internal/cli/importer/parser.go` - Use constants

### Dependencies

- No external dependencies
- `constants` is a leaf package (no internal imports)
- Import hierarchy: `constants` <- `config` <- `domain` <- `service` <- `cli`

### Patterns to Use

| Need | Use |
|------|-----|
| Type-safe enums | `type Role string` with const block |
| Validation | `IsValid*(s string) bool` helper functions |
| List all values | `All*() []Type` helper functions |
| Config access | Pass `*config.Config` to service constructors |

Run `make what-to-use NEED="keyword"` for details.

## Testing Requirements

### Unit Tests

- [x] `IsValidRoleFit()` returns true for valid roles, false for invalid
- [x] `AllRoleFits()` returns complete list
- [x] `CVTargetRoleFits()` returns only principal, staff, em, senior_ic
- [x] `IsValidCompetencyCategory()` validates correctly
- [x] `IsValidSkillLevel()` validates correctly
- [x] `IsValidEventTag()` validates correctly
- [x] `SuggestedSkillCategories()` includes all suggestions including "mobile"
- [x] `ScoringConfig` defaults applied when section missing
- [x] Auto-migration preserves existing config values
- [x] Weight validation rejects weights not summing to 1.0
- [x] Domain validation still works with constants
- [x] Service scoring uses config weights and thresholds

### E2E Tests (if new intent)

- [x] Happy Paths: Existing E2E tests pass with constants
- [x] Sad Paths: Existing validation tests pass

## Definition of Done

- [x] All acceptance criteria met
- [x] Tests written FIRST (TDD)
- [x] Tests pass with >= 95% coverage
- [x] No pattern violations (`make check-patterns`)
- [x] Compliance check passes (`make check-compliance`)
- [x] Documentation updated (PR description)
- [x] Committed with `make ai-commit`
- [x] PR created: #107
- [x] PR review feedback addressed (commit 479b66a)

## Completion Notes

**Status**: ✅ COMPLETE - Awaiting manual review thread resolution and merge

**PR**: https://github.com/baphled/KaRiya/pull/107

**Commits**:
1. `9df93d0` - feat(domain): add centralized constants package (TASK-50)
2. `8d7a413` - feat(config): add ScoringConfig for configurable CV scoring (TASK-50)
3. `bbeae3e` - refactor(domain): use centralized constants package (TASK-50)
4. `f30b04f` - refactor(service): use centralized constants package (TASK-50)
5. `b4bd839` - fix(config): remove redundant nil check for map (staticcheck)
6. `e049a04` - refactor(domain): fix naming to follow {Concept}{Value} pattern
7. `479b66a` - fix(cli): address PR review feedback

**Review Feedback Addressed**:
- Comment #1: Updated skill_form.go to use constants.SuggestedSkillCategories() (resolves "mobile" category missing from UI)
- Comment #2: Improved config auto-migration to check entire scoring section, avoiding overwriting user-provided 0.0 values

**CI Status**: ✅ All tests passing (Ubuntu, macOS, Windows)

**Manual Action Required**: Resolve review comment threads on GitHub UI before merge
