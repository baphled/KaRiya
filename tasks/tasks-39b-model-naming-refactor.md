---
created: 2026-01-14T02:31
modified: 2026-01-14T02:39
---
# Task 39B: Model Naming Refactor - Remove Library Prefixes

## Overview
- **Goal**: Refactor model file names and types to remove library-specific prefixes (e.g., `huh_`) for better abstraction and maintainability
- **Time Estimate**: 1-2 hours
- **Prerequisites**: Task 39 (User-Defined Skills) complete

## Session Contract Acknowledgment
- [x] Ran `make session-start` and it passed
- [x] Acknowledge and commit to following all workflow rules
- [x] Token count: N/A (task complete)

## Pre-Task Checklist (MUST COMPLETE BEFORE STARTING)
- [x] `make check-compliance` passes
- [x] Reviewed existing patterns in:
  - `internal/cli/models/huh_skill_form.go` (now `skill_form.go`)
  - `internal/cli/models/huh_capture_form.go` (now `capture_form.go`)
- [x] Confirmed this is ONE atomic task (renaming refactor)
- [x] Identified which test files will be created/modified (none needed - pure rename)

## Context

During Task 39 (User-Defined Skills), wrapper models were created for huh forms to handle form alignment issues. These models were named with the `huh_` prefix (e.g., `huh_skill_form.go`, `HuhSkillForm`), which couples the naming to the specific library implementation.

**Problem**: Library-specific naming in model files creates several issues:
1. **Tight Coupling**: Names reveal implementation details that consumers shouldn't need to know
2. **Refactoring Friction**: If we ever switch form libraries, all references need updating
3. **Inconsistency**: Other models don't have library prefixes (e.g., `form.go`, not `bubbles_form.go`)
4. **Documentation Overhead**: Requires explaining "huh" in documentation

**Solution**: Rename files and types to use domain-focused names that describe *what* they do, not *how* they do it.

## Current State

### Files to Rename

| Current File | Current Type | New File | New Type |
|--------------|--------------|----------|----------|
| `internal/cli/models/huh_skill_form.go` | `HuhSkillForm` | `internal/cli/models/skill_form.go` | `SkillForm` |
| `internal/cli/models/huh_capture_form.go` | `HuhCaptureForm` | `internal/cli/models/capture_form.go` | `CaptureForm` |

### References to Update

After renaming, all references to these types need updating:

1. **Intent files** that use these models
2. **Test files** that reference these types
3. **Documentation** mentioning `HuhSkillForm` or `HuhCaptureForm`

## Implementation Plan

### Phase 1: Rename SkillForm

**Goal**: Rename `huh_skill_form.go` to `skill_form.go` and `HuhSkillForm` to `SkillForm`

#### TDD Checklist - Phase 1:
- [x] Verify existing tests pass before refactor
- [x] Rename file: `huh_skill_form.go` -> `skill_form.go`
- [x] Rename type: `HuhSkillForm` -> `SkillForm`
- [x] Rename constructor: `NewHuhSkillForm` -> `NewSkillForm`
- [x] Update all references in:
  - [x] `internal/cli/intents/manage_skills_intent.go`
  - [x] `internal/cli/intents/manage_skills_test.go` (no direct references found)
  - [x] Any other files referencing `HuhSkillForm`
- [x] Run tests to verify refactor success
- [x] Commit: `refactor(models): rename Huh-prefixed form wrappers to domain-focused names` (6f42fe7)

### Phase 2: Rename CaptureForm

**Goal**: Rename `huh_capture_form.go` to `capture_form.go` and `HuhCaptureForm` to `CaptureForm`

#### TDD Checklist - Phase 2:
- [x] Verify existing tests pass before refactor
- [x] Rename file: `huh_capture_form.go` -> `capture_form.go`
- [x] Rename type: `HuhCaptureForm` -> `CaptureForm`
- [x] Rename constructor: `NewHuhCaptureForm` -> `NewCaptureForm`
- [x] Update all references in:
  - [x] `internal/cli/intents/capture_event_intent.go`
  - [x] `internal/cli/intents/capture_event.go`
  - [x] `internal/cli/models/form.go` (not referenced)
  - [x] Any other files referencing `HuhCaptureForm`
- [x] Run tests to verify refactor success
- [x] Commit: Combined with Phase 1 in single commit (6f42fe7)

### Phase 3: Documentation Update

**Goal**: Update documentation to reflect new naming

#### TDD Checklist - Phase 3:
- [x] Update `docs/FORMS_GUIDE.md`:
  - [x] Change references from `HuhSkillForm` to `SkillForm`
  - [x] Change references from `HuhCaptureForm` to `CaptureForm`
  - [x] Update code examples
- [x] Update `docs/rules/FORMS_WORKFLOW_GUIDE.md` (not present)
- [x] Update `AGENTS.md` if wrapper models are mentioned
- [x] Commit: Combined with code changes in single commit (6f42fe7)

## Pre-Commit Checklist (BEFORE EACH COMMIT)
- [ ] `make check-compliance` passes (REQUIRED before commit)
- [ ] Use `make ai-commit MSG="type(scope): description"` for AI-generated code
- [ ] Commit message explains **WHY**, not just WHAT
- [ ] Commit is atomic (ONE logical change)

## Post-Task Checklist (MUST COMPLETE BEFORE NEXT TASK)
- [x] `make check-compliance` passes
- [x] All checkboxes above completed
- [x] Task marked complete `[x]` in task file
- [x] Token count: N/A (task complete)

## Acceptance Criteria

- [x] No files in `internal/cli/models/` have `huh_` prefix
- [x] No types have `Huh` prefix (unless genuinely huh-specific internal details)
- [x] All tests pass (100% pass rate maintained)
- [x] All references updated (no compilation errors)
- [x] Documentation reflects new naming
- [x] `make check-compliance` passes

## Rollback Plan

This is a pure rename refactor with no logic changes:
- Git history preserves old names
- `git revert` will restore original naming if needed
- No database or state changes involved

## Design Principles

### Naming Convention Going Forward

For future form wrapper models, use this pattern:

```go
// File: internal/cli/models/{domain}_form.go
// Type: {Domain}Form
// Constructor: New{Domain}Form

// Example: Skill management form wrapper
// File: skill_form.go
// Type: SkillForm
// Constructor: NewSkillForm()

// Example: Event capture form wrapper  
// File: capture_form.go
// Type: CaptureForm
// Constructor: NewCaptureForm()
```

### When Library Prefixes ARE Appropriate

Library prefixes should only be used when:
1. Multiple implementations exist (e.g., `SqliteRepository`, `MemoryRepository`)
2. The library name IS the domain (e.g., `BubbleTeaModel` if wrapping BubbleTea)
3. Internal implementation details not exposed to consumers

### When Library Prefixes Should Be Avoided

Avoid library prefixes when:
1. There's only one implementation
2. The consumer doesn't need to know the underlying library
3. The name describes what it does, not how it's implemented

## Notes

### Why This Matters

Good naming is a form of documentation. When a developer sees `SkillForm`, they understand it's a form for skills. When they see `HuhSkillForm`, they need to:
1. Know what "huh" is
2. Wonder if there are other skill form implementations
3. Question if they should be using a different one

### Impact on Existing Code

This refactor has minimal risk because:
- All changes are compile-time verifiable (Go compiler catches missing references)
- No runtime behavior changes
- Tests validate all functionality remains intact
- IDE refactoring tools can assist with bulk renaming

### Future Considerations

If additional form libraries are introduced in the future and multiple implementations are needed, we can:
1. Use interface-based design with a common `Form` interface
2. Use factory pattern to create appropriate implementation
3. Keep library-specific details internal to the implementation
