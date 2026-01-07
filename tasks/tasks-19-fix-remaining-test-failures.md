# Task 19: Fix Remaining Test Failures After Form Refactoring

**Created**: 2026-01-07
**Status**: In Progress
**Priority**: CRITICAL
**Estimated Time**: 1-2 hours
**Related**: Task 17 (Form Refactoring), Task 18 (Planning)

---

## Overview

After the CaptureEvent form refactoring (Task 17), only **3 test failures** remain instead of the 52 anticipated. This task will fix these specific failures and ensure all tests pass.

### Current Status
- **Total tests**: 980
- **Passing**: 977  
- **Failing**: 3
- **Build status**: ✅ Successful
- **Test compilation**: ✅ All tests compile

---

## Test Failures to Fix

### 1. cmd/cli/persistence_test.go - Line 183
**Test**: "Form Submission Persistence > when user submits form > should return event data from form submission"

**Issue**: Test expects `cmd` to be non-nil when Enter is pressed on Submit button, but receives nil

**Root Cause**: Test navigation assumes all fields are visible (mode system), but with the strategy system, field visibility depends on strategy. The test is tabbing through fields but may not be landing on the SubmitButton due to incorrect field count.

**Fix**:
1. Update field navigation to account for strategy system
2. Ensure test uses manual strategy (all fields visible) OR correctly navigates through visible fields only
3. Verify SubmitButton focus before pressing Enter

---

### 2. cmd/cli/persistence_test.go - Line 253  
**Test**: "Form Submission Debug > debugging form persistence > should show what mode is being used in form submission"

**Issue**: Test expects `cmd` to be non-nil when Enter is pressed on Submit button, but receives nil

**Root Cause**: Same as #1 - test navigation assumes old mode system with ModeField at index 6

**Fix**:
1. Remove reference to ModeField (line 248: `form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab}) // ModeField (6)`)
2. Update field navigation count
3. Verify manual strategy is used (or update navigation for quick strategy)

---

### 3. internal/cli/intents/capture_event_escape_test.go - Line 182
**Test**: "CaptureEvent - Escape Key Behavior > View Methods > should show 'esc', 'm', and 'r' in Submit footer"

**Issue**: Test expects footer to contain "Retry" but actual footer text differs

**Root Cause**: Footer text was updated but test expectations weren't updated to match

**Fix**:
1. Check actual footer text in CaptureStateSubmit state
2. Update test expectations to match current footer text
3. Verify "Esc", "Main Menu", and "Retry" (or equivalent) are present

---

## Files to Modify

- [ ] `cmd/cli/persistence_test.go` (lines 134-202, 206-270)
- [ ] `internal/cli/intents/capture_event_escape_test.go` (lines 170-183)

---

## Implementation Plan

### Phase 1: Analyze Current Behavior (15 min)

- [x] Run failing tests to see exact error messages
- [x] Identify root causes
- [ ] Check FormModel field visibility logic
- [ ] Check SubmitButton focus conditions
- [ ] Check footer text in Submit state

### Phase 2: Fix Persistence Tests (30 min)

#### Test 1 & 2 Fix Strategy

Option A: **Use Manual Strategy** (Recommended - 15 min)
```go
// Create form model with manual strategy
form := models.NewFormModel(cliSvc)
form.SetStrategy("manual")  // Ensure all fields visible
form.SetShowOptionalFields(true)

// Then navigate as before
form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab}) // DateField
form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab}) // CompanyField
form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab}) // ProjectField
form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab}) // TagsField
form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab}) // CategoriesField
form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab}) // SubmitButton
```

Option B: **Navigate Using Quick Strategy** (Alternative - 20 min)
```go
// In quick mode, only TextField and SubmitButton are visible
form, _ = updateForm(form, tea.KeyMsg{Type: tea.KeyTab}) // SubmitButton (skips hidden fields)
```

**Tasks**:
- [ ] Choose fix strategy (recommend Option A for thorough testing)
- [ ] Update test at line 134-202 (first persistence test)
- [ ] Update test at line 206-270 (debug test)
- [ ] Remove reference to ModeField (line 248)
- [ ] Add strategy setup to both tests
- [ ] Verify tests pass

### Phase 3: Fix Escape Test (15 min)

**Tasks**:
- [ ] Read actual footer text from CaptureStateSubmit view
- [ ] Update test expectations to match
- [ ] Ensure test checks for correct shortcut keys
- [ ] Verify test passes

### Phase 4: Verification (15 min)

**Tasks**:
- [ ] Run all cmd/cli tests: `go test ./cmd/cli -v`
- [ ] Run all intent tests: `go test ./internal/cli/intents -v`
- [ ] Run full test suite: `go test ./... -v`
- [ ] Verify race detector: `go test -race ./...`
- [ ] Check test coverage maintained

---

## Acceptance Criteria

### Must Have
- [ ] All 3 failing tests now pass
- [ ] Total test pass rate: 980/980 (100%)
- [ ] No new test failures introduced
- [ ] All tests compile without errors
- [ ] Zero race conditions detected

### Should Have
- [ ] Tests use strategy system correctly
- [ ] Test navigation logic is clear and maintainable
- [ ] Footer text tests match actual implementation

### Nice to Have
- [ ] Add comments explaining field navigation logic
- [ ] Add helper function for form field navigation in tests
- [ ] Document strategy behavior in test comments

---

## Testing Strategy

### TDD Approach
1. **Red**: Run failing test, confirm exact failure
2. **Green**: Make minimal change to fix test
3. **Refactor**: Clean up test code if needed
4. **Verify**: Run full suite to ensure no regressions

### Verification Commands
```bash
# Individual test files
go test ./cmd/cli -v -run="Form Submission"
go test ./internal/cli/intents -v -run="Escape Key"

# Full test suite
go test ./... -v

# With race detector
go test -race ./...

# Check compliance
make check-compliance
```

---

## Risk Assessment

### Low Risk
1. **Tests are isolated** - Changes only affect test code, not implementation
2. **Small scope** - Only 3 tests to fix
3. **Clear root causes** - Navigation logic issue, not complex bugs

### Mitigation
- Fix one test at a time
- Run tests after each fix
- Use TDD approach (Red-Green-Refactor)

---

## Success Metrics

- ✅ All 980 tests passing (100% pass rate)
- ✅ Zero test compilation errors
- ✅ Zero race conditions
- ✅ Build successful
- ✅ Ready for Phase 2 (staticcheck fixes)

---

## Next Steps (After Task 19)

Once all tests pass:
1. **Task 20**: Fix staticcheck issues (28 warnings)
2. **Task 21**: Remove deprecated code (modeIndex, modes fields)
3. **Task 22**: Update documentation
4. **Task 23**: Add strategy system tests

---

## References

### Related Documents
- `tasks/tasks-17-legacy-views-removal.md` - Form refactoring (COMPLETE)
- `tasks/tasks-18-next-steps-plan.md` - Overall plan
- `docs/rules/master-task-prompt.md` - Task workflow
- `docs/TUI_STANDARDS.md` - TUI design standards

### Related Files
- `internal/cli/models/form.go` - Form implementation (strategy system)
- `cmd/cli/persistence_test.go` - Persistence tests (2 failures)
- `internal/cli/intents/capture_event_escape_test.go` - Escape test (1 failure)

---

## Notes

### Key Decisions
1. **Use manual strategy in tests** - Ensures all fields visible for thorough testing
2. **Fix tests, not implementation** - Form code is correct, tests need updating
3. **One test at a time** - Verify each fix before moving to next

### Things to Watch
- Field navigation logic - must account for hidden fields
- Strategy initialization - ensure tests use correct strategy
- Footer text - must match actual implementation

---

**Last Updated**: 2026-01-07
**Author**: AI Assistant (via OpenCode)
**Status**: Ready for implementation
