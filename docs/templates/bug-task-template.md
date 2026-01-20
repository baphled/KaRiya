# TASK-XXX: Fix BUG-YYY - [Title]

## Bug Reference

See: `bugs/BUG-YYY.md`

## Summary

Brief description of the fix (1-2 sentences).

## Root Cause

Explanation of why the bug occurred.

## Fix Approach

How the bug will be fixed.

## Files to Modify

- `path/to/file.go` - Description of change

## Testing Requirements

### Regression Test (REQUIRED)

```go
Describe("Bug Regressions", func() {
    It("BUG-YYY: prevents [bug behavior]", func() {
        // Test that verifies the bug is fixed
    })
})
```

### Additional Tests

- [ ] Test case covering the fix
- [ ] Test edge cases if applicable

## Definition of Done

- [ ] Root cause identified
- [ ] Regression test written FIRST (TDD)
- [ ] Fix implemented
- [ ] All tests pass
- [ ] No pattern violations (`make check-patterns`)
- [ ] Compliance check passes (`make check-compliance`)
- [ ] Committed with `make ai-commit`
- [ ] Bug file updated with resolution
