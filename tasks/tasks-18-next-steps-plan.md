---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Task 18: Next Steps Plan - Post Form Refactoring

**Created**: 2026-01-07
**Status**: Planning
**Priority**: High
**Estimated Time**: 8-12 hours

---

## Overview

Following the successful CaptureEvent form refactoring (commits 7046248-ecca5af), this document outlines the next steps to complete the cleanup, fix failing tests, and address technical debt.

---

## Current State

### ✅ Completed (Task 17)
- Replaced capture mode system with strategy system (quick/manual)
- Fixed side-by-side field alignment (Company/Project, Tags/Categories)
- Removed view duplication (StandardView now handles chrome)
- Increased default terminal width (120 → 140)
- Created 5 atomic commits with clear documentation

### ⚠️ Known Issues
1. **52 failing form model tests** - Expect old capture mode validation
2. **28 staticcheck issues** - Code quality warnings (mostly in other intents)
3. **Deprecated fields** - `modeIndex`, `modes` still present in FormModel
4. **Missing documentation** - Strategy system not documented in user guides
5. **Test compilation errors** - 17+ errors in generate_cv_test.go, cv_test.go, cv_config_manager_test.go (unrelated to our changes)

### 🎯 Project Status
- **Total tests**: 980 (927 passing, 52 failing, 1 skipped)
- **Build status**: ✅ Successful
- **CaptureEvent tests**: ✅ 13/13 passing
- **Race conditions**: 0 detected

---

## Priority Tasks

### Priority 0: Fix Test Compilation Errors (CRITICAL)
**Impact**: CRITICAL - Tests won't run
**Time**: 30 min - 1 hour

#### Test Files with Compilation Errors
1. `internal/cli/intents/generate_cv_test.go` (14 errors)
   - `selectedAudiences` field doesn't exist
   - Type mismatches with string vs []string

2. `internal/cli/intents/generate_cv_integration_test.go` (3 errors)
   - Similar type mismatches

3. `internal/domain/career/cv_test.go` (12 errors)
   - Type mismatches with TargetAudience field
   - Section content structure issues

4. `internal/cli/models/cv_config_manager_test.go` (3 errors)
   - Type mismatches

#### Root Cause
These appear to be pre-existing issues from a previous CV-related refactoring that changed `TargetAudience` from `string` to `[]string` but didn't update all tests.

#### Resolution
- Fix type mismatches in test fixtures
- Update struct field references
- Verify all tests compile before proceeding

**Note**: These are NOT related to our form refactoring, but must be fixed before running tests.

---

### Priority 1: Fix Failing Tests (CRITICAL)
**Impact**: High - Blocking CI/CD pipeline
**Time**: 3-4 hours (after Priority 0)

#### Failing Test Categories
1. **Form submission with modes** (12 tests)
   - Timeline Journaling mode validation
   - CV Backfill mode validation
   - Manual Entry mode validation

2. **Date validation tests** (8 tests)
   - Events older than 30 days (Timeline Journaling)
   - Events from any date (CV Backfill)
   - Flexible dates (Manual Entry)

3. **Mode-specific behavior** (15 tests)
   - Mode selection tests
   - Mode transition tests
   - Mode-specific error messages

4. **Integration tests** (17 tests)
   - End-to-end capture workflows
   - Tag/category integration with modes
   - Form data validation with modes

#### Approach
**Option A: Delete obsolete tests** (Fastest - 1 hour)
- Mark 52 tests as `Pending` or delete them
- Rationale: They test functionality that no longer exists
- Risk: Loss of test coverage

**Option B: Rewrite tests for strategy system** (Recommended - 3-4 hours)
- Rewrite tests to validate quick/manual strategy behavior
- Add tests for optional field toggle ('t' key)
- Add tests for field visibility logic
- Maintain test coverage at current level

**Recommendation**: Option B - Maintains quality standards

---

### Priority 2: Remove Deprecated Code (MEDIUM)
**Impact**: Medium - Technical debt
**Time**: 1 hour

#### Deprecated Fields to Remove
```go
// internal/cli/models/form.go
type FormModel struct {
    // ... other fields
    modeIndex int                              // DEPRECATED - Remove
    modes     []careerservice.EventCaptureMode // DEPRECATED - Remove
}
```

#### Cleanup Tasks
- [ ] Remove `modeIndex` field from FormModel
- [ ] Remove `modes` field from FormModel
- [ ] Remove any unused mode-related helper functions
- [ ] Update comments referencing modes
- [ ] Search codebase for any remaining mode references

#### Verification
```bash
# Search for remaining mode references
rg "modeIndex|EventCaptureMode|Timeline.*Journaling|CV.*Backfill" --type go
```

---

### Priority 3: Fix Staticcheck Issues (MEDIUM)
**Impact**: Medium - Code quality
**Time**: 2-3 hours

#### Current Issues (28 total)
```
internal/cli/intents/bulk_operations_intent.go:79:14: error strings should not be capitalized (ST1005)
internal/cli/intents/bulk_operations_intent.go:117:10: strings.Title is deprecated (SA1019)
internal/cli/intents/burst_management_intent.go:123:2: field burstEditor is unused (U1000)
internal/cli/intents/fact_management_intent.go:174:14: error strings should not be capitalized (ST1005)
internal/cli/intents/generate_cv_intent.go:30:2: field cvPreview is unused (U1000)
... and 23 more
```

#### Categories
1. **Error capitalization** (8 issues) - Quick fixes
2. **Deprecated APIs** (3 issues) - Replace with modern equivalents
3. **Unused fields** (12 issues) - Remove or implement usage
4. **Unnecessary code** (5 issues) - Simplify

#### Approach
- Fix by category (batched commits)
- Run staticcheck after each category
- Ensure no regressions in tests

---

### Priority 4: Update Documentation (LOW)
**Impact**: Low - User experience
**Time**: 2 hours

#### Documents to Update
1. **User Guides**
   - [ ] `docs/CLI_GUIDE.md` - Add strategy system explanation
   - [ ] `docs/KEYBOARD_REFERENCE.md` - Document 't' key toggle
   - [ ] `docs/BURST_FACT_EXTRACTION_GUIDE.md` - Update capture workflow

2. **Developer Guides**
   - [ ] `docs/TUI_DEVELOPER_GUIDE.md` - Add strategy pattern example
   - [ ] `docs/TUI_STANDARDS.md` - Update form patterns
   - [ ] `AGENTS.md` - Update recent fixes section

3. **Architecture Docs**
   - [ ] `docs/SCREEN_TO_INTENT_MAPPING.md` - Update CaptureEvent section
   - [ ] `docs/WORKFLOW_DIAGRAM.md` - Update capture flow

#### New Documentation Needed
- [ ] Strategy system design document
- [ ] Migration guide (modes → strategies)
- [ ] Form field visibility specification

---

### Priority 5: Add Strategy System Tests (LOW)
**Impact**: Low - Test coverage enhancement
**Time**: 2-3 hours

#### Test Coverage Gaps
1. **Strategy selection tests**
   - [ ] Quick strategy selection
   - [ ] Manual strategy selection
   - [ ] Strategy persistence across form lifecycle

2. **Field visibility tests**
   - [ ] Quick mode: only event text visible
   - [ ] Manual mode: all fields visible
   - [ ] Toggle optional fields behavior

3. **Integration tests**
   - [ ] Quick capture end-to-end
   - [ ] Manual capture with optional fields
   - [ ] Edit mode always uses manual strategy

4. **Edge cases**
   - [ ] Toggle when focused on hidden field
   - [ ] Navigation skips hidden fields
   - [ ] Form submission with minimal fields

#### Target Coverage
- Maintain 87%+ overall coverage
- 100% coverage for strategy logic
- Add benchmark tests for field visibility checks

---

## Implementation Phases

### Phase 1: Critical Path (Priority 1-2)
**Goal**: Get all tests passing and remove technical debt
**Time**: 4-5 hours

**Tasks**:
1. ✅ Analyze 52 failing tests (30 min)
2. ✅ Rewrite tests for strategy system (3 hours)
3. ✅ Remove deprecated fields (30 min)
4. ✅ Verify all tests passing (30 min)

**Deliverables**:
- All 980 tests passing
- Zero deprecated fields
- Clean git history (atomic commits)

### Phase 2: Code Quality (Priority 3)
**Goal**: Fix all staticcheck issues
**Time**: 2-3 hours

**Tasks**:
1. ✅ Fix error capitalization (30 min)
2. ✅ Replace deprecated APIs (1 hour)
3. ✅ Remove unused fields (1 hour)
4. ✅ Simplify unnecessary code (30 min)

**Deliverables**:
- Zero staticcheck warnings
- Cleaner, more maintainable code
- Updated to modern Go idioms

### Phase 3: Documentation (Priority 4)
**Goal**: Update all documentation
**Time**: 2 hours

**Tasks**:
1. ✅ Update user guides (1 hour)
2. ✅ Update developer guides (30 min)
3. ✅ Update architecture docs (30 min)

**Deliverables**:
- Accurate, up-to-date documentation
- Strategy system fully documented
- Migration guide for future reference

### Phase 4: Test Enhancement (Priority 5)
**Goal**: Comprehensive strategy system testing
**Time**: 2-3 hours

**Tasks**:
1. ✅ Add strategy selection tests (1 hour)
2. ✅ Add field visibility tests (1 hour)
3. ✅ Add integration tests (1 hour)

**Deliverables**:
- 100% strategy system coverage
- Comprehensive edge case testing
- Performance benchmarks

---

## Success Criteria

### Must Have (Phase 1-2)
- [ ] All 980 tests passing (100% pass rate)
- [ ] Zero deprecated fields in codebase
- [ ] Zero staticcheck warnings
- [ ] Build successful (`go build ./cmd/cli`)
- [ ] Zero race conditions (`go test -race ./...`)

### Should Have (Phase 3)
- [ ] All documentation updated
- [ ] Strategy system fully documented
- [ ] Migration guide available

### Nice to Have (Phase 4)
- [ ] Strategy system tests at 100% coverage
- [ ] Performance benchmarks documented
- [ ] Edge cases fully tested

---

## Risk Assessment

### High Risk
1. **Test rewrites break functionality** 
   - Mitigation: TDD approach, run after each change
   
2. **Staticcheck fixes introduce bugs**
   - Mitigation: Small atomic commits, test after each

### Medium Risk
3. **Removing deprecated fields breaks something**
   - Mitigation: Comprehensive search before removal
   
4. **Documentation updates are time-consuming**
   - Mitigation: Focus on critical docs first

### Low Risk
5. **New tests reveal edge case bugs**
   - Mitigation: Fix as discovered, improves quality

---

## Timeline Estimate

### Optimistic (8 hours)
- Phase 1: 4 hours
- Phase 2: 2 hours
- Phase 3: 1 hour
- Phase 4: 1 hour

### Realistic (12 hours)
- Phase 1: 5 hours
- Phase 2: 3 hours
- Phase 3: 2 hours
- Phase 4: 2 hours

### Pessimistic (16 hours)
- Phase 1: 6 hours (test rewrites more complex)
- Phase 2: 4 hours (API replacements tricky)
- Phase 3: 3 hours (comprehensive docs)
- Phase 4: 3 hours (thorough testing)

**Recommended**: Plan for 12 hours (realistic timeline)

---

## Quick Start: Next Session

```bash
# 1. Start with highest priority
cd /home/baphled/Projects/KaRiya

# 2. Analyze failing tests
ginkgo -v --focus="Timeline Journaling|CV Backfill|Manual Entry" ./internal/cli/models/

# 3. Create task branch
git checkout -b fix/form-tests-and-cleanup

# 4. Begin Phase 1
# See "Phase 1: Critical Path" section above
```

---

## Follow-Up Tasks

After completing this task, consider:

1. **Performance optimization** - Benchmark form rendering
2. **Accessibility improvements** - Screen reader support
3. **Enhanced validation** - Field-level validation as user types
4. **Auto-save** - Persist form state across sessions
5. **Form themes** - Support for dark/light themes

---

## References

### Related Documents
- `AGENTS.md` - Project overview and standards
- `docs/rules/master-task-prompt.md` - Task execution workflow
- `docs/TUI_DEVELOPER_GUIDE.md` - TUI development patterns
- `docs/TUI_STANDARDS.md` - TUI design standards

### Related Tasks
- Task 17: CaptureEvent Form Refactoring (COMPLETE)
- Task 16: StandardView Testing & Polish (COMPLETE)
- Task 12-15: StandardView Implementation (COMPLETE)

### Recent Commits
- 7046248: feat(ui): increase default terminal width
- a969f40: refactor(form): replace capture mode system
- d041307: refactor(capture-event): integrate strategy system
- 5e653d5: test: update tests to reflect capture mode removal
- ecca5af: refactor(generate-cv): remove breadcrumbs

---

## Notes for Next Developer

### Key Decisions Made
1. **Strategy system over capture modes** - Simpler, more flexible
2. **StandardView handles all chrome** - No duplication
3. **Side-by-side alignment** - Equal column widths for consistency
4. **Default to manual strategy** - Show all fields by default

### Things to Watch Out For
1. **52 failing tests** - Need complete rewrite, not quick fixes
2. **Deprecated fields** - Remove carefully, check all usages
3. **Staticcheck warnings** - Some are in unrelated intents
4. **Documentation** - Keep user/dev guides in sync

### Testing Strategy
- TDD for test rewrites (Red-Green-Refactor)
- Run tests after each small change
- Use `go test -race` frequently
- Verify with `make check-compliance` before committing

---

**Last Updated**: 2026-01-07
**Author**: AI Assistant (via OpenCode)
**Status**: Ready for implementation
