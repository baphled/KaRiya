# TASK-56: Refactor Generic Modals to Shared Package

## Summary

Create `internal/cli/modals/` package and move generic, cross-feature modals from feature-specific packages to a shared location for better code organization and discoverability.

## Acceptance Criteria

- [ ] `internal/cli/modals/` package created with doc.go
- [ ] `SuggestionReviewModal` moved from `screens/burst_management/modals/` to `internal/cli/modals/`
- [ ] `AddEditModal` moved from `screens/skills/modals/` to `internal/cli/modals/` and renamed to `SkillFormModal`
- [ ] `EventsModal` moved from `screens/skills/modals/` to `internal/cli/modals/`
- [ ] All consuming packages updated with new import paths
- [ ] All tests passing (no regressions)
- [ ] Documentation updated (architecture docs, modal patterns)
- [ ] PR created and reviewed

## Technical Notes

### Problem Statement

Generic modals used by multiple intents are scattered across feature-specific packages:

**SuggestionReviewModal** - In `screens/burst_management/modals/` but used by:
- `intents/browsetimeline` (skill suggestions)
- `intents/skillsmanagement` (skill suggestions)
- `intents/burst_management` (burst + skill suggestions)

**AddEditModal** - In `screens/skills/modals/` but used by:
- `intents/browsetimeline` (create new skill)
- `intents/skillsmanagement` (manage skills)
- `intents/burst_management` (create skills from suggestions)

**EventsModal** - In `screens/skills/modals/` but used by:
- `intents/skillsmanagement` (view events for skill)
- `intents/burst_management` (view events for suggestion)

### Files to Modify

**New Package Structure:**
```
internal/cli/modals/
├── doc.go
├── suggestion_review_modal.go      # Move from burst_management/modals/
├── suggestion_review_modal_test.go
├── skill_form_modal.go              # Move AddEditModal from skills/modals/ (RENAMED)
├── skill_form_modal_test.go
├── events_modal.go                  # Move from skills/modals/
└── events_modal_test.go
```

**Consuming Packages to Update (9 files):**

browsetimeline:
- `internal/cli/intents/browsetimeline/types.go`
- `internal/cli/intents/browsetimeline/helpers.go`

skillsmanagement:
- `internal/cli/intents/skillsmanagement/types.go`
- `internal/cli/intents/skillsmanagement/helpers.go`

burst_management:
- `internal/cli/intents/burst_management/types.go`
- `internal/cli/intents/burst_management/helpers.go`
- `internal/cli/intents/burst_management/handlers.go`

Tests:
- Any integration tests importing these modals

### Dependencies

**Prerequisites:**
- Task 55 (Event Skill Management) merged
- All tests passing before migration starts

**Blocks:**
- Any new features using these generic modals

### Patterns to Use

| Need | Use |
|------|-----|
| Package docs | `doc.go` with package-level godoc |
| Import path | `github.com/baphled/kariya/internal/cli/modals` |
| Modal suffix | Keep `*Modal` naming convention |
| Rename clarity | `AddEditModal` → `SkillFormModal` |

## Testing Requirements

### Unit Tests

- [ ] All existing modal tests pass after move
- [ ] Package path updates in test imports
- [ ] No new tests needed (pure refactor)

### Integration Tests

- [ ] browsetimeline E2E tests: 101/101 passing
- [ ] skillsmanagement tests passing
- [ ] burst_management tests passing
- [ ] No test behavior changes

### Regression Testing

Run after each modal move:
```bash
go test ./internal/cli/modals/... -v
go test ./internal/cli/intents/browsetimeline/... -v -tags=e2e
go test ./internal/cli/intents/skillsmanagement/... -v -tags=e2e
go test ./internal/cli/intents/burst_management/... -v -tags=e2e
```

## Definition of Done

- [ ] All acceptance criteria met
- [ ] Tests written FIRST (N/A - pure refactor)
- [ ] Tests pass with no regressions
- [ ] No pattern violations (`make check-patterns`)
- [ ] Compliance check passes (`make check-compliance`)
- [ ] Documentation updated (architecture guide, modal patterns)
- [ ] Committed with `make ai-commit`

## Migration Plan

### Phase 1: Create Package and Doc (15 min)

1. Create directory:
   ```bash
   mkdir -p internal/cli/modals
   ```

2. Create `doc.go`:
   ```go
   // Package modals provides generic, cross-feature modal components.
   //
   // This package contains modals that are reused across multiple intents
   // and features. Feature-specific modals should remain in their respective
   // screens/{feature}/modals/ directories.
   //
   // # Generic Modals
   //
   // - SuggestionReviewModal: Review AI-generated suggestions (skills, bursts)
   // - SkillFormModal: Create/edit skills
   // - EventsModal: View events associated with a skill
   //
   // # Usage
   //
   //	import "github.com/baphled/kariya/internal/cli/modals"
   //	
   //	modal := modals.NewSuggestionReviewModal(suggestions)
   //
   package modals
   ```

### Phase 2: Move SuggestionReviewModal (45 min)

1. Move files:
   ```bash
   git mv internal/cli/screens/burst_management/modals/suggestion_review_modal.go \
          internal/cli/modals/suggestion_review_modal.go
   git mv internal/cli/screens/burst_management/modals/suggestion_review_modal_test.go \
          internal/cli/modals/suggestion_review_modal_test.go
   ```

2. Update package declaration in both files:
   ```go
   package modals
   ```

3. Update imports in consuming files:
   ```go
   // OLD
   import burstModals "github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
   
   // NEW
   import "github.com/baphled/kariya/internal/cli/modals"
   ```

4. Update struct references:
   ```go
   // OLD
   suggestionModal *burstModals.SuggestionReviewModal
   
   // NEW
   suggestionModal *modals.SuggestionReviewModal
   ```

5. Run tests:
   ```bash
   go test ./internal/cli/modals/... -v
   go test ./internal/cli/intents/browsetimeline/... -v -tags=e2e
   go test ./internal/cli/intents/skillsmanagement/... -v
   go test ./internal/cli/intents/burst_management/... -v
   ```

### Phase 3: Move and Rename AddEditModal (45 min)

1. Move and rename files:
   ```bash
   git mv internal/cli/screens/skills/modals/add_edit_modal.go \
          internal/cli/modals/skill_form_modal.go
   git mv internal/cli/screens/skills/modals/add_edit_modal_test.go \
          internal/cli/modals/skill_form_modal_test.go
   ```

2. Update package and rename struct:
   ```go
   package modals
   
   // OLD: type AddEditModal struct
   // NEW: type SkillFormModal struct
   
   // OLD: func NewAddEditModal
   // NEW: func NewSkillFormModal
   ```

3. Update imports and references in consuming files:
   ```go
   // OLD
   import skillModals "github.com/baphled/kariya/internal/cli/screens/skills/modals"
   addEditModal *skillModals.AddEditModal
   
   // NEW
   import "github.com/baphled/kariya/internal/cli/modals"
   skillFormModal *modals.SkillFormModal
   ```

4. Run tests

### Phase 4: Move EventsModal (30 min)

1. Move files:
   ```bash
   git mv internal/cli/screens/skills/modals/events_modal.go \
          internal/cli/modals/events_modal.go
   git mv internal/cli/screens/skills/modals/events_modal_test.go \
          internal/cli/modals/events_modal_test.go
   ```

2. Update package, imports, references

3. Run tests

### Phase 5: Cleanup and Documentation (15 min)

1. Remove empty `screens/burst_management/modals/` if empty
2. Update architecture docs
3. Update modal patterns guide
4. Create PR

## Risks & Mitigation

| Risk | Impact | Mitigation |
|------|--------|------------|
| Breaking imports | High | Move one modal at a time, test after each |
| Test failures | Medium | Run full suite after each move |
| Circular dependencies | Low | Modals are leaf nodes, no intent dependencies |
| Merge conflicts | Medium | Complete quickly, coordinate with team |

## Estimated Effort

- Package creation + doc: 15 min
- Move SuggestionReviewModal: 45 min
- Move + rename SkillFormModal: 45 min
- Move EventsModal: 30 min
- Test + fix issues: 30 min
- Documentation: 15 min

**Total**: 2-3 hours

## Benefits

1. **Discoverability**: Developers know where to find reusable modals
2. **Consistency**: Generic modals in one place encourages reuse
3. **Maintenance**: Easier to find and update shared components
4. **Architecture**: Clear separation between feature-specific and generic UI

## References

### Related Tasks
- TASK-55: Event Skill Management (PR #155)

### Documentation
- [Modal Patterns](../MODAL_PATTERNS.md)
- [Intent Architecture Guide](../INTENT_ARCHITECTURE_GUIDE.md)
- [Modal Naming Conventions](../conventions/MODAL_NAMING.md)

### Code Locations
- Current: `internal/cli/screens/*/modals/`
- Target: `internal/cli/modals/`
