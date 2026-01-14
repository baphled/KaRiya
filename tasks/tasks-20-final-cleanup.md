---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Task 20: Final Cleanup - Staticcheck and Code Quality

**Created**: 2026-01-07
**Status**: In Progress
**Priority**: MEDIUM
**Estimated Time**: 30 minutes
**Related**: Task 18 (Planning), Task 19 (Test Fixes - COMPLETE)

---

## Overview

After completing Task 19 (all tests passing), we now have minimal cleanup remaining:
- **3 staticcheck warnings** (unused functions)
- All deprecated fields already removed
- All tests passing (2,078/2,078)

This task will address the final code quality issues before the project is 100% clean.

---

## Current Status

### ✅ Already Complete
- All 2,078 tests passing (100% pass rate)
- Zero race conditions
- Build successful
- Deprecated fields removed (modeIndex, modes)
- Test coverage >87%

### ⚠️ Remaining Issues
Only **3 staticcheck warnings** - all unused functions:

1. `internal/cli/intents/burst_management_intent.go:707` - `applyFilters` unused
2. `internal/cli/intents/burst_management_intent.go:1342` - `setFailed` unused
3. `internal/cli/intents/generate_cv_intent.go:364` - `getStateName` unused

---

## Files to Modify

- [ ] `internal/cli/intents/burst_management_intent.go` (2 unused functions)
- [ ] `internal/cli/intents/generate_cv_intent.go` (1 unused function)

---

## Implementation Plan

### Phase 1: Analyze Unused Functions (10 min)

**Tasks**:
- [ ] Check if `applyFilters` is actually needed
- [ ] Check if `setFailed` is actually needed
- [ ] Check if `getStateName` is actually needed
- [ ] Determine if they should be removed or used

### Phase 2: Remove Unused Functions (10 min)

#### Option A: Remove Functions (Recommended if truly unused)
```go
// Simply delete the unused functions
```

#### Option B: Mark as Internal/Helper (If needed for future use)
```go
// Add comment explaining why kept, or prefix with _ to indicate unused
func (i *Intent) _helperForFuture() { ... }
```

**Tasks**:
- [ ] Remove or mark unused functions
- [ ] Run staticcheck to verify warnings gone
- [ ] Run tests to ensure no breakage

### Phase 3: Verification (10 min)

**Tasks**:
- [ ] Run staticcheck: `staticcheck ./...`
- [ ] Run full test suite: `go test ./...`
- [ ] Run with race detector: `go test -race ./...`
- [ ] Verify build: `go build ./cmd/cli`
- [ ] Check compliance: `make check-compliance`

---

## Detailed Analysis

### 1. burst_management_intent.go:707 - applyFilters

**Location**: Line 707
**Issue**: Function is defined but never called

**Analysis Needed**:
- Was this intended for filtering bursts?
- Is there a future feature that will use it?
- Should it be removed or kept with a comment?

### 2. burst_management_intent.go:1342 - setFailed

**Location**: Line 1342
**Issue**: Function is defined but never called

**Analysis Needed**:
- Was this for error handling?
- Is there error handling that should use it?
- Should it be removed or integrated?

### 3. generate_cv_intent.go:364 - getStateName

**Location**: Line 364
**Issue**: Function is defined but never called

**Analysis Needed**:
- Was this for debugging?
- Is it used in View() or logging?
- Should it be removed or kept for debugging?

---

## Decision Criteria

For each unused function, ask:

1. **Is it referenced anywhere?**
   ```bash
   rg "applyFilters|setFailed|getStateName" --type go
   ```

2. **Was it part of a refactoring that removed its usage?**
   - Check git history: `git log -p -- <file> | grep -A 5 -B 5 <function>`

3. **Is it a helper that should be used but isn't?**
   - Check surrounding code for error handling gaps
   - Check if similar patterns exist elsewhere

4. **Decision**:
   - **Remove**: If truly unused and no clear future use
   - **Keep with comment**: If intended for future use, add `// TODO:` or `// Unused but kept for...`
   - **Use it**: If it should be called but isn't (rare)

---

## Acceptance Criteria

### Must Have
- [ ] Zero staticcheck warnings
- [ ] All 2,078 tests still passing (100% pass rate)
- [ ] No new test failures introduced
- [ ] Build successful
- [ ] Zero race conditions

### Should Have
- [ ] Clear git commit message explaining removals
- [ ] Code more maintainable (no dead code)

### Nice to Have
- [ ] Documentation updated if functions were part of public API

---

## Testing Strategy

### Before Making Changes
```bash
# Baseline
staticcheck ./... > /tmp/before_staticcheck.txt
go test ./... > /tmp/before_tests.txt
```

### After Each Change
```bash
# Incremental verification
staticcheck ./...
go test ./internal/cli/intents/... -v
```

### Final Verification
```bash
# Complete check
make check-compliance
go test -race ./...
```

---

## Risk Assessment

### Very Low Risk
1. **Removing unused functions** - No impact if truly unused
2. **Breaking tests** - All tests already passing, unlikely to break

### Mitigation
- Check git history before removing
- Run tests after each removal
- One function at a time
- Easy to revert if needed

---

## Success Metrics

- ✅ Zero staticcheck warnings (from 3)
- ✅ All tests passing (maintained at 2,078/2,078)
- ✅ Build successful
- ✅ Zero race conditions
- ✅ Clean codebase (no dead code)

---

## Next Steps (After Task 20)

Once staticcheck is clean:
1. **Optional**: Update documentation (if time permits)
2. **Optional**: Add strategy system tests (if time permits)
3. **Ready for production**: All critical cleanup complete

---

## References

### Related Documents
- `tasks/tasks-18-next-steps-plan.md` - Original cleanup plan
- `tasks/tasks-19-fix-remaining-test-failures.md` - Test fixes (COMPLETE)
- `docs/rules/master-task-prompt.md` - Task workflow

### Related Commands
```bash
# Check staticcheck issues
staticcheck ./...

# Check for usage
rg "functionName" --type go

# Git history
git log -p -- internal/cli/intents/burst_management_intent.go | grep -A 5 -B 5 "applyFilters"

# Full verification
make check-compliance
```

---

## Notes

### Key Decisions
1. **Remove dead code** - Improves maintainability
2. **One function at a time** - Safe, incremental approach
3. **Verify after each removal** - Catch issues early

### Things to Watch
- Make sure functions aren't called via reflection
- Check if functions are part of interface implementations
- Verify tests still compile

---

**Last Updated**: 2026-01-07
**Author**: AI Assistant (via OpenCode)
**Status**: Ready for implementation
