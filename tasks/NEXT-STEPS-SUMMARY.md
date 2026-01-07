# Next Steps - Quick Reference

**Created**: 2026-01-07
**Last Session**: CaptureEvent Form Refactoring (Complete)
**Next Task**: Fix Tests & Cleanup

---

## 🎯 Immediate Priorities

### 1. Fix Test Compilation (CRITICAL - 30 min)
```bash
# Files with compilation errors (32 errors total):
internal/cli/intents/generate_cv_test.go            (14 errors)
internal/cli/intents/generate_cv_integration_test.go (3 errors)
internal/domain/career/cv_test.go                    (12 errors)
internal/cli/models/cv_config_manager_test.go        (3 errors)

# Issue: TargetAudience type changed (string → []string) but tests not updated
# Action: Fix type mismatches in test fixtures
```

### 2. Fix 52 Failing Form Tests (HIGH - 3-4 hours)
```bash
# Tests expecting old capture mode system
# Categories:
# - Timeline Journaling mode (12 tests)
# - CV Backfill mode (8 tests)  
# - Manual Entry mode (15 tests)
# - Integration tests (17 tests)

# Action: Rewrite tests for quick/manual strategy system
```

### 3. Remove Deprecated Code (MEDIUM - 1 hour)
```bash
# Fields to remove from FormModel:
# - modeIndex int
# - modes []careerservice.EventCaptureMode

# Search for references:
rg "modeIndex|EventCaptureMode|Timeline.*Journaling|CV.*Backfill" --type go
```

### 4. Fix 28 Staticcheck Issues (MEDIUM - 2-3 hours)
```bash
make staticcheck

# Categories:
# - Error capitalization (8 issues)
# - Deprecated APIs (3 issues)
# - Unused fields (12 issues)
# - Unnecessary code (5 issues)
```

---

## 📊 Current Status

### Build & Tests
- ✅ **Main code builds**: `go build ./cmd/cli` SUCCESS
- ⚠️ **Test compilation**: 32 errors (pre-existing, CV-related)
- ⚠️ **Test pass rate**: 927/979 passing (52 failing from mode removal)
- ✅ **CaptureEvent tests**: 13/13 passing
- ⚠️ **Staticcheck**: 28 warnings

### What Just Shipped (5 commits)
- ✅ Terminal width: 120 → 140
- ✅ Capture modes → Strategy system (quick/manual)
- ✅ Side-by-side field alignment fixed
- ✅ View duplication removed (StandardView integration)
- ✅ Tests updated for mode removal

---

## 🚀 Quick Start Next Session

```bash
# 1. Check out repo
cd /home/baphled/Projects/KaRiya

# 2. Create task branch
git checkout -b fix/tests-and-cleanup

# 3. Start with Priority 0 (compilation errors)
# Fix test fixtures in generate_cv_test.go first

# 4. Verify tests compile
go test ./... --build-only

# 5. Move to Priority 1 (rewrite 52 failing tests)
ginkgo -v ./internal/cli/models/ | grep FAIL

# 6. Follow TDD: Red → Green → Refactor
```

---

## 📝 Detailed Plan

See `tasks/tasks-18-next-steps-plan.md` for:
- Complete task breakdown
- Risk assessment
- Timeline estimates (8-12 hours)
- Phase-by-phase implementation guide
- Success criteria
- Documentation updates needed

---

## 🎓 Context for AI Assistant

### Key Files Modified (Last Session)
```
internal/cli/components/standard_view.go       (width increase)
internal/cli/models/form.go                    (strategy system)
internal/cli/intents/capture_event.go          (strategy types)
internal/cli/intents/capture_event_intent.go   (integration)
+ 5 test files (updated assertions)
```

### Architectural Decisions
1. **Strategy over Modes**: Simplified from 3 modes to 2 strategies
2. **StandardView Handles Chrome**: No duplication in form views
3. **Equal Column Widths**: Side-by-side fields use calculated widths
4. **Default to Manual**: Show all fields by default

### Testing Strategy
- Use TDD for test rewrites (Red-Green-Refactor)
- Run tests after every small change
- Use `go test -race` frequently
- Run `make check-compliance` before committing

---

## ⚠️ Important Notes

### Don't Touch (Working Code)
- StandardView implementation (complete, tested)
- Intent router (stable)
- CV generation service (tested, working)
- Repository layer (stable)

### Safe to Modify
- Form model tests (need complete rewrite)
- Deprecated fields (can remove after verification)
- Staticcheck warnings (fix category by category)
- Documentation (update as you go)

### Watch Out For
- Test compilation errors (fix FIRST)
- Race conditions (run with `-race` flag)
- Deprecated field removal (search thoroughly first)
- Breaking changes (maintain backward compatibility where possible)

---

## 📚 Reference Documents

### Must Read
- `docs/rules/master-task-prompt.md` - Complete workflow
- `docs/rules/TASK_QUICK_REF.md` - Quick checklist
- `AGENTS.md` - Project overview

### For Context
- `docs/TUI_STANDARDS.md` - TUI design standards
- `docs/TUI_DEVELOPER_GUIDE.md` - Development patterns
- `docs/rules/atomic-commits.md` - Commit standards

---

## ✅ Success Criteria (Exit Conditions)

### Must Have
- [ ] All tests compile (`go test ./... --build-only`)
- [ ] All 980 tests pass (100% pass rate)
- [ ] Zero deprecated fields
- [ ] Zero staticcheck warnings
- [ ] Build successful
- [ ] Zero race conditions

### Should Have
- [ ] Documentation updated
- [ ] Strategy system documented
- [ ] Commit messages follow conventions

### Nice to Have
- [ ] 100% strategy system test coverage
- [ ] Performance benchmarks
- [ ] Migration guide

---

**Last Updated**: 2026-01-07
**Status**: Ready for implementation
**Estimated Time**: 8-12 hours total
