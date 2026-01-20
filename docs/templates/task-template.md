# TASK-XXX: [Title]

## Summary

Brief description of the task (1-2 sentences).

## Acceptance Criteria

- [ ] Criterion 1
- [ ] Criterion 2
- [ ] Criterion 3

## Technical Notes

### Files to Modify

- `path/to/file.go` - Description of change

### Dependencies

- List any dependencies or prerequisites

### Patterns to Use

| Need | Use |
|------|-----|
| Table | `behaviors.TableBehavior[T]` |
| Form | `models.*Form` wrapper |
| Colors | `theme.Primary()` etc |

Run `make what-to-use NEED="keyword"` for details.

## Testing Requirements

### Unit Tests

- [ ] Test case 1
- [ ] Test case 2

### E2E Tests (if new intent)

- [ ] Happy Paths: All state transitions
- [ ] Sad Paths: Error conditions

## Definition of Done

- [ ] All acceptance criteria met
- [ ] Tests written FIRST (TDD)
- [ ] Tests pass with >= 95% coverage
- [ ] No pattern violations (`make check-patterns`)
- [ ] Compliance check passes (`make check-compliance`)
- [ ] Documentation updated (if needed)
- [ ] Committed with `make ai-commit`
