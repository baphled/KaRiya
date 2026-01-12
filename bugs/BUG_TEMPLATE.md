# Bug XXX: [Short Description]

**Status**: Reported | Investigation | Root Cause | Fix | Testing | Verification | Closed  
**Severity**: 🔴 Critical | 🟠 High | 🟡 Medium | 🟢 Low  
**Created**: YYYY-MM-DD  
**Updated**: YYYY-MM-DD  
**Resolved**: YYYY-MM-DD (if closed)

---

## Bug Summary

[One-line description of the issue]

---

## Affected Components

- [ ] Component 1 (`path/to/file.go`)
- [ ] Component 2 (`path/to/file.go`)
- [ ] Component 3 (`path/to/file.go`)

**Related Intents/Workflows**:
- Intent name
- Workflow name

---

## Reproduction Steps

1. Step 1
2. Step 2
3. Step 3
4. Observe issue

**Consistency**: Always | Sometimes | Rarely

**Environment**:
- OS: [Linux/macOS/Windows]
- Terminal: [gnome-terminal/iTerm/etc]
- Go Version: [1.24+]
- KaRiya Version: [commit hash]

---

## Expected Behavior

[What should happen]

---

## Actual Behavior

[What actually happens]

**Evidence**:
- Screenshots: [if applicable]
- Error messages: [if any]
- Logs: [relevant logs]

---

## Investigation Log

### [YYYY-MM-DD HH:MM] - Initial Investigation
- Finding 1
- Finding 2
- Finding 3

### [YYYY-MM-DD HH:MM] - Code Review
- Reviewed `file.go:line`
- Found potential issue in `function()`
- Related to [pattern/concept]

### [YYYY-MM-DD HH:MM] - Test Review
- Existing tests: [pass/fail]
- Test coverage: [X%]
- Gap identified: [what's not tested]

---

## Root Cause

**Status**: Identified | Suspected | Unknown

**Cause**:
[Detailed explanation of why the bug occurs]

**Technical Details**:
- File: `path/to/file.go`
- Function: `functionName()`
- Line: `123`
- Reason: [technical explanation]

---

## Fix Strategy

**Approach**: [High-level approach]

### Option A: [Description] (Recommended/Alternative)
**Pros**:
- Pro 1
- Pro 2

**Cons**:
- Con 1
- Con 2

**Files to Change**:
- [ ] `path/to/file.go` - [description of change]
- [ ] `path/to/test.go` - [description of change]

### Option B: [Description] (Alternative)
[Similar structure as Option A]

**Selected Approach**: Option A/B  
**Rationale**: [Why this approach]

---

## Testing Plan

### Phase 1: Unit Tests
- [ ] Test case 1: [description]
- [ ] Test case 2: [description]
- [ ] Test case 3: [description]

**Files**:
- `path/to/test.go`

### Phase 2: Integration Tests
- [ ] Integration test 1: [description]
- [ ] Integration test 2: [description]

**Files**:
- `path/to/integration_test.go`

### Phase 3: Manual Testing
- [ ] Manual test 1: [steps]
- [ ] Manual test 2: [steps]
- [ ] Manual test 3: [steps]

### Phase 4: Regression Testing
- [ ] Verify no regressions in [feature]
- [ ] Run full test suite
- [ ] Check performance impact

---

## Verification Checklist

### Code Quality
- [ ] Fix implemented and tested
- [ ] All tests passing (go test ./...)
- [ ] No race conditions (go test -race)
- [ ] Code coverage maintained (>80%)
- [ ] Linting passing (staticcheck)

### Functionality
- [ ] Issue no longer reproduces
- [ ] Expected behavior confirmed
- [ ] Edge cases handled
- [ ] Error messages clear (if applicable)

### Documentation
- [ ] Code comments added/updated
- [ ] User-facing docs updated (if needed)
- [ ] Bug report updated with resolution
- [ ] Related issues cross-referenced

### Compliance
- [ ] Follows project coding standards
- [ ] Atomic commits with clear messages
- [ ] AI attribution (if applicable)
- [ ] No breaking changes (or documented)

---

## Resolution Summary

**Fix Description**:
[Summary of what was changed and why]

**Commits**:
- `commit-hash` - [commit message]
- `commit-hash` - [commit message]

**Pull Request**: #XXX (if applicable)

**Related Tasks**: tasks-XX-name.md (if created follow-up work)

---

## Follow-Up Actions

- [ ] Action 1
- [ ] Action 2
- [ ] Action 3

---

## Related Files

### Implementation Files
- `path/to/file.go:line` - [description]
- `path/to/other.go:line` - [description]

### Test Files
- `path/to/test.go:line` - [description]
- `path/to/integration_test.go:line` - [description]

### Documentation Files
- `docs/path/to/doc.md` - [description]

---

## References

- Related bugs: bug-XXX
- Related tasks: tasks-XX
- Documentation: docs/path
- External references: [URLs if applicable]

---

**Last Updated**: YYYY-MM-DD  
**Updated By**: [Name]
