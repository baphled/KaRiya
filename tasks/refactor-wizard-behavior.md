# TASK-XX: Standardize Wizard Pattern with WizardBehavior[T]

## Summary

Create a generic `WizardBehavior[T]` in `behaviors/` to standardize multi-step wizard modals, similar to `TableBehavior[T]`.

## Acceptance Criteria

- [ ] `WizardBehavior[T]` created with generic wizard functionality
- [ ] CV config wizard migrated to use `WizardBehavior`
- [ ] Onboarding wizard migrated to use `WizardBehavior`
- [ ] All existing wizard tests pass
- [ ] New unit tests for `WizardBehavior`

## Technical Notes

### Files to Modify

- `internal/cli/behaviors/wizard_behavior.go` - NEW: Generic wizard behavior
- `internal/cli/behaviors/wizard_behavior_test.go` - NEW: Tests
- `internal/cli/components/cv_config_wizard_modal.go` - Refactor to use behavior
- `internal/cli/components/onboarding_wizard_modal.go` - Refactor to use behavior

### Dependencies

- BUG-005 fix merged (form selection persistence)

### Patterns to Use

| Need | Use |
|------|-----|
| Generic behavior | `WizardBehavior[T any]` (like `TableBehavior[T]`) |
| Form handling | `*huh.Form` with pointer bindings |
| Rendering | `uikit/containers`, `uikit/primitives` |

Run `make what-to-use NEED="behavior"` for details.

### Common Wizard Patterns to Extract

Both existing wizards share:
- Multi-step `huh.Form`
- Data struct for collected values
- State tracking (visible/completed/cancelled)
- Width/height sizing
- `buildForm()` method
- UIKit rendering (title, step indicator, content, footer)

## Testing Requirements

### Unit Tests

- [ ] `WizardBehavior` step navigation (next/previous)
- [ ] `WizardBehavior` completion tracking
- [ ] `WizardBehavior` cancellation handling
- [ ] `WizardBehavior` keyboard shortcuts (Tab, Enter, Esc)
- [ ] Migrated wizards retain existing behavior

### E2E Tests (if applicable)

- [ ] CV config wizard e2e tests still pass
- [ ] Onboarding wizard e2e tests still pass

## Definition of Done

- [ ] All acceptance criteria met
- [ ] Tests written FIRST (TDD)
- [ ] Tests pass with >= 95% coverage
- [ ] No pattern violations (`make check-patterns`)
- [ ] Compliance check passes (`make check-compliance`)
- [ ] Documentation updated (if needed)
- [ ] Committed with `make ai-commit`
