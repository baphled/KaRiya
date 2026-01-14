---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Task 20: Coverage Improvement - Reach 80% Threshold

## Overview
- **Goal**: Improve test coverage from 60.3% to ≥80% to meet compliance requirements
- **Time Estimate**: 2-3 days
- **Prerequisites**: Understanding of existing test patterns (Ginkgo/Gomega)
- **Blocker For**: Task 21 (TUI Visual Overhaul) - must complete before proceeding

## Session Contract Acknowledgment
- [x] Ran `make session-start` and it passed (or coverage is the only blocker)
- [x] Acknowledge and commit to following all workflow rules
- [x] Token count: 95,720 (session complete)

## Pre-Task Checklist (MUST COMPLETE BEFORE STARTING)
- [x] All checks pass EXCEPT coverage in `make check-compliance`
- [x] Reviewed existing test patterns in:
  - [x] `internal/cli/intents/*_test.go` (Ginkgo pattern examples)
  - [x] `internal/service/career/*_test.go` (Service test patterns)
  - [x] `internal/repository/career/*_test.go` (Repository test patterns)
- [x] Confirmed this is ONE atomic task (coverage improvement only)
- [x] Identified specific packages below 80% coverage

## Current Coverage Analysis

**Packages Below 80% (from most critical to least):**

### Priority 1: Core Packages (Essential for Task 21)
- **`internal/cli/styles`**: 32.4% → Target: 80%
  - Critical for TUI visual overhaul (Task 21)
  - Need tests for: style generation, theme application, color schemes

### Priority 2: Intent Layer
- **`internal/cli/intents`**: 53.7% → Target: 80%
  - Need tests for: edge cases, error handling, state transitions

### Priority 3: Repository Layer
- **`internal/repository/career`**: 52.3% → Target: 80%
  - Need tests for: CRUD operations, error handling, query edge cases

### Priority 4: Configuration
- **`internal/config`**: 63.6% → Target: 80%
  - Need tests for: config loading, validation, defaults

### Lower Priority (Can Skip if Time Constrained)
- **`internal/logger`**: 70.0% → Target: 80%
- **`internal/cli/layout`**: 79.1% → Target: 80% (already close!)

## Implementation Plan

### ✅ PHASE 1: Styles Package Coverage (Day 1) - COMPLETE

**Result**: 32.4% → 100% (+67.6%)
**Commits**: `beeef98`
**Specs**: 85 new tests, 135 total passing

### PHASE 1: Styles Package Coverage (Day 1)

#### Step 1.1: Analyze Current Coverage
- [ ] **Run coverage analysis**
  ```bash
  go test -coverprofile=coverage.out ./internal/cli/styles/
  go tool cover -html=coverage.out -o coverage.html
  ```
- [ ] **Identify uncovered functions**
  - Open `coverage.html` in browser
  - List uncovered functions in notes below

#### Step 1.2: TDD - Style Generation Tests
- [ ] **Write failing tests** (`internal/cli/styles/styles_test.go`)
  - [ ] Test: Color scheme generation
  - [ ] Test: Border style generation
  - [ ] Test: Text style generation
  - [ ] Test: Layout style generation
  - [ ] Run tests and confirm they FAIL
  - [ ] Commit: `test(styles): add failing tests for style generation`

- [ ] **Implement missing test coverage**
  - [ ] Add tests to reach 80% coverage
  - [ ] Run tests and confirm they PASS
  - [ ] Commit: `test(styles): achieve 80% coverage for styles package`

#### Step 1.3: Verify Coverage
- [x] **Run coverage check**
  ```bash
  go test -cover ./internal/cli/styles/
  ```
- [x] **Verify ≥80% coverage**
- [x] **Commit if needed**: `test(styles): final coverage improvements`

### PHASE 2: Intents Package Coverage (Day 2)

#### Step 2.1: Identify Coverage Gaps
- [x] **Run coverage analysis**
  ```bash
  go test -coverprofile=coverage.out ./internal/cli/intents/
  go tool cover -html=coverage.out -o coverage.html
  ```
- [x] **Identify most critical gaps** (error handling, edge cases)

#### Step 2.2: TDD - Intent Edge Cases
- [ ] **Write failing tests for uncovered paths**
  - [ ] Test: Error state handling
  - [ ] Test: Cancel operations
  - [ ] Test: Back navigation edge cases
  - [ ] Test: Invalid state transitions
  - [ ] Run tests and confirm they FAIL
  - [ ] Commit: `test(intents): add tests for error and edge cases`

- [ ] **Add coverage to reach 80%**
  - [ ] Focus on error paths and edge cases
  - [ ] Run tests and confirm they PASS
  - [ ] Commit: `test(intents): achieve 80% coverage for intents package`

#### Step 2.3: Verify Coverage
- [ ] **Run coverage check**
  ```bash
  go test -cover ./internal/cli/intents/
  ```
- [ ] **Verify ≥80% coverage**

### PHASE 3: Repository Package Coverage (Day 2-3)

#### Step 3.1: Identify Coverage Gaps
- [ ] **Run coverage analysis**
  ```bash
  go test -coverprofile=coverage.out ./internal/repository/career/
  go tool cover -html=coverage.out -o coverage.html
  ```
- [ ] **Identify uncovered database operations**

#### Step 3.2: TDD - Repository Error Cases
- [ ] **Write failing tests for uncovered paths**
  - [ ] Test: Database connection errors
  - [ ] Test: Constraint violations
  - [ ] Test: Query edge cases (empty results, large datasets)
  - [ ] Test: Transaction rollback scenarios
  - [ ] Run tests and confirm they FAIL
  - [ ] Commit: `test(repository): add tests for database edge cases`

- [ ] **Add coverage to reach 80%**
  - [ ] Focus on error handling and edge cases
  - [ ] Run tests and confirm they PASS
  - [ ] Commit: `test(repository): achieve 80% coverage for repository package`

#### Step 3.3: Verify Coverage
- [ ] **Run coverage check**
  ```bash
  go test -cover ./internal/repository/career/
  ```
- [ ] **Verify ≥80% coverage**

### PHASE 4: Config Package Coverage (Day 3)

#### Step 4.1: Identify Coverage Gaps
- [ ] **Run coverage analysis**
  ```bash
  go test -coverprofile=coverage.out ./internal/config/
  go tool cover -html=coverage.out -o coverage.html
  ```

#### Step 4.2: TDD - Config Validation Tests
- [ ] **Write failing tests for uncovered paths**
  - [ ] Test: Invalid config loading
  - [ ] Test: Missing config file handling
  - [ ] Test: Default value application
  - [ ] Test: Config validation errors
  - [ ] Run tests and confirm they FAIL
  - [ ] Commit: `test(config): add tests for validation and defaults`

- [ ] **Add coverage to reach 80%**
  - [ ] Run tests and confirm they PASS
  - [ ] Commit: `test(config): achieve 80% coverage for config package`

#### Step 4.3: Verify Coverage
- [ ] **Run coverage check**
  ```bash
  go test -cover ./internal/config/
  ```
- [ ] **Verify ≥80% coverage**

### PHASE 5: Verification & Cleanup (Day 3)

#### Step 5.1: Full Coverage Check
- [ ] **Run full coverage analysis**
  ```bash
  go test -coverprofile=coverage.out ./...
  go tool cover -func=coverage.out | grep total
  ```
- [ ] **Verify overall coverage ≥80%**

#### Step 5.2: Compliance Check
- [ ] **Run full compliance check**
  ```bash
  make check-compliance
  ```
- [ ] **Verify ALL checks pass** (including coverage)

#### Step 5.3: Documentation
- [ ] **Update AGENTS.md if needed**
  - Update coverage statistics
  - Document any new test patterns used
  - Commit: `docs: update coverage statistics`

## Pre-Commit Checklist (BEFORE EACH COMMIT)
- [ ] `make review-commit` passes
- [ ] AI attribution included (if AI-generated code)
- [ ] Commit message explains **WHY**, not just WHAT
- [ ] Commit is atomic (ONE logical change per package/area)
- [ ] Tests pass for committed code

## Post-Task Checklist (MUST COMPLETE BEFORE NEXT TASK)
- [ ] `make check-compliance` passes (INCLUDING coverage ≥80%)
- [ ] All checkboxes above completed
- [ ] Task marked complete `[x]` in task file
- [ ] Token count: _____ (< 100k to continue)

## Acceptance Criteria
- [ ] Overall test coverage ≥80%
- [ ] `internal/cli/styles` coverage ≥80%
- [ ] `internal/cli/intents` coverage ≥80%
- [ ] `internal/repository/career` coverage ≥80%
- [ ] `internal/config` coverage ≥80%
- [ ] All existing tests still pass
- [ ] `make check-compliance` passes completely
- [ ] Zero race conditions detected
- [ ] No new staticcheck warnings

## Rollback Plan
1. All test additions are non-breaking
2. Can revert individual commits if needed
3. Tests are additive (don't modify existing functionality)
4. If coverage target proves unrealistic, discuss with team before proceeding

## Notes

### Coverage Improvement Strategies
1. **Focus on error paths**: Most uncovered code is error handling
2. **Test edge cases**: Boundary conditions, empty inputs, nil values
3. **Use table-driven tests**: Efficient way to add coverage
4. **Mock external dependencies**: Isolate units under test
5. **Don't test for coverage sake**: Ensure tests are meaningful

### Test Patterns to Follow
- **Ginkgo/Gomega** for all tests
- **Table-driven tests** for multiple scenarios
- **Test helpers** in `internal/testutil/` for common setup
- **Mock interfaces** where needed (avoid excessive mocking)

### What NOT to Test
- **Generated code** (mocks, etc.)
- **Third-party library code**
- **Trivial getters/setters** (unless they have logic)
- **Test utilities themselves** (except testutil package)

### Success Metrics
- [ ] Coverage ≥80% overall
- [ ] Task 21 can proceed without compliance blockers
- [ ] No performance regression (test suite < 10s total)
- [ ] All tests meaningful and maintainable
- [ ] Zero flaky tests introduced

## Priority Reminder

**This task is a BLOCKER for Task 21 (TUI Visual Overhaul).**

The visual overhaul is a major change that requires a solid foundation. Completing this coverage task ensures:
1. We meet compliance requirements
2. We have confidence in existing functionality before major changes
3. We can catch regressions during the visual overhaul
4. The codebase is in good health for the next phase

**Once this task is complete and `make session-start` passes, we can proceed with Task 21.**

---

## Session 1 Summary (2026-01-09)

### Achievements
- **Overall Coverage**: 60.3% → 63.7% (+3.4 percentage points)
- **Duration**: ~2 hours
- **Branch**: `test/coverage-improvements`
- **PR**: #43 (https://github.com/baphled/KaRiya/pull/43)

### Packages Improved
1. **Styles**: 32.4% → 100% ✅ (CRITICAL FOR TASK 21)
2. **Config**: 63.6% → 81.8% ✅
3. **Layout**: 79.1% → 88.4% ✅  
4. **Repository**: 52.3% → 52.6%

### Commits
- `beeef98` - test(cli): add comprehensive getter function tests for styles
- `c0f7d39` - test(config): add tests for LoadConfig and SaveConfig wrappers
- `bb1cd71` - test(cli): add tests for layout manager edge cases
- `fee2166` - test(repo): add tests for NewSQLiteRepositoryWithDB and GetDB

### Next Steps
1. Review and merge PR #43 to `next`
2. Decide on follow-up approach:
   - **Option A**: Continue with Task 21 (styles at 100%, sufficient for TUI overhaul)
   - **Option B**: Create follow-up task for remaining packages (intents, repository)

### Remaining Work (to reach 80%)
- **Intents package**: 53.7% → Need +26.3% (~1-2 days)
- **Repository package**: 52.6% → Need +27.4% (~1 day)
- **Logger package**: 70.0% → Need +10% (~2-3 hours)

**Estimated time to 80%**: 2-3 additional days
