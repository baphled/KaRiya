# KaRiya Codebase Audit - Implementation Tasks Summary

**Created**: 2026-01-08
**Audit Date**: 2026-01-08
**Status**: Ready for Implementation
**Total Tasks Created**: 6 (Tasks 29-34)

---

## Overview

A comprehensive codebase audit identified **90+ partially implemented features** across all layers of the KaRiya application. This document summarizes the 6 task files created to systematically address these gaps.

**Audit Findings**:
- **Critical**: 2 completely non-functional intents (ExportArtifact, ConfigureSystem)
- **High**: 4 major feature gaps (modals, filters, audience filtering, review states)
- **Medium**: 1 cleanup task (consistency, dead code)
- **Total Estimated Time**: 21-30 hours
- **Current Project Status**: All 2,078 tests passing, zero staticcheck warnings, production ready (but with incomplete features)

---

## Task Files Created

| Task | Title | Priority | Time | Status |
|------|-------|----------|------|--------|
| [Task 29](tasks-29-export-artifact-critical-fixes.md) | ExportArtifact Critical Fixes | CRITICAL | 4-6h | Ready |
| [Task 30](tasks-30-configure-system-critical-fixes.md) | ConfigureSystem Critical Fixes | CRITICAL | 4-5h | Ready |
| [Task 31](tasks-31-capture-event-modal-review-fixes.md) | CaptureEvent Modal & Review Fixes | HIGH | 3-4h | Ready |
| [Task 32](tasks-32-browse-timeline-filter-implementation.md) | BrowseTimeline Filter Implementation | HIGH | 3-4h | Ready |
| [Task 33](tasks-33-service-layer-audience-filtering.md) | Service Layer Audience Filtering | HIGH | 2-3h | Ready |
| [Task 34](tasks-34-consistency-and-cleanup.md) | Consistency & Cleanup | MEDIUM | 2-3h | Ready |

**Total**: 6 tasks, 18-25 hours estimated

---

## Task 29: ExportArtifact Critical Fixes (CRITICAL)

**Priority**: CRITICAL | **Time**: 4-6 hours

### Issues
- Export operation is completely stubbed (always returns `/tmp/export.*` with 1024 bytes)
- All preview data is static/hardcoded (no real user data shown)
- `formatBytes()` bug: displays garbage characters instead of "1 KB", "1 MB"
- Scroll percentage display has same bug
- ExportService exists and works but not integrated

### Phases
1. **Bug Fixes** (1h): Fix formatBytes() and scroll percentage display
2. **Service Integration** (2-3h): Connect to ExportService, implement real export
3. **Real Preview Data** (1-2h): Fetch and display actual user data
4. **Consistency** (30min): Add vim j/k navigation, use loadingRotator

### Impact
Users cannot export their CVs, events, or artifacts. Feature advertised but completely non-functional.

---

## Task 30: ConfigureSystem Critical Fixes (CRITICAL)

**Priority**: CRITICAL | **Time**: 4-5 hours

### Issues
- Save operation is stubbed (always returns success, nothing persisted)
- EditSettings state has no input components (can't edit values)
- Configuration values are hardcoded placeholders
- No configuration file persistence

### Phases
1. **Configuration Persistence** (1-2h): Create config file handler, load/save to `~/.kariya/config.yaml`
2. **Settings Editing UI** (2-3h): Add text inputs, implement value editing
3. **Real Save Operation** (1h): Persist changes to config file
4. **Additional Fixes** (30min): Vim navigation, loadingRotator

### Impact
Users cannot configure the system. All settings changes are lost immediately.

---

## Task 31: CaptureEvent Modal & Review Fixes (HIGH)

**Priority**: HIGH | **Time**: 3-4 hours

### Issues
- Review view reads from wrong state path (always shows empty data)
- Modals are never displayed even when EditingMode is set
- Accept/reject handlers missing despite being advertised in help text

### Phases
1. **Fix Review State Path** (30min): Read from `reviewState.Event` instead of `result.Event`
2. **Implement Modal Display** (2h): Show metadata/burst/fact edit modals when mode is set
3. **Implement Accept/Reject** (1-1.5h): Add handlers for accepting/rejecting inferred items
4. **Fix Help Text** (15min): Update footer to match actual functionality

### Impact
Users cannot review or edit captured events properly. Inferred bursts/facts cannot be managed.

---

## Task 32: BrowseTimeline Filter Implementation (HIGH)

**Priority**: HIGH | **Time**: 3-4 hours

### Issues
- Footer advertises "f Filter" but pressing 'f' does nothing
- Footer advertises "/ Search" but pressing '/' does nothing
- 4 filter types defined but never implemented (Companies, Categories, DateFrom, DateTo)
- Pagination test failures - shows ALL events instead of current page

### Phases
1. **Implement Filter UI** (2h): Add 'f' key handler, implement all 4 filters
2. **Implement Search UI** (1h): Add '/' key handler, real-time search
3. **Fix Pagination** (1h): Display correct page of events, fix failing tests
4. **Vim Navigation** (30min): Add g/G for jump to top/bottom

### Impact
Users cannot filter or search their event timeline. Pagination is broken.

---

## Task 33: Service Layer Audience Filtering (HIGH)

**Priority**: HIGH | **Time**: 2-3 hours

### Issues
- `FilterByAudience()` returns all bullets for all audiences (no filtering)
- `isEventRelevantToAudience()` always returns true
- `isFactRelevantToAudience()` always returns true
- `customizeForRole()` returns text unchanged

### Phases
1. **FilterByAudience** (1h): Implement real audience-based bullet filtering
2. **isEventRelevantToAudience** (30min): Check categories/tags/text for keywords
3. **isFactRelevantToAudience** (30min): Check AudienceRelevance field
4. **customizeForRole** (30min): Adapt language for role seniority

### Impact
Generated CVs are not tailored to target audience. All bullets appear regardless of relevance.

---

## Task 34: Consistency & Cleanup (MEDIUM)

**Priority**: MEDIUM | **Time**: 2-3 hours

### Issues
- Help screen ('?' key) does nothing
- CV generation errors never shown to user
- 3 unused loadingRotator components
- 12+ orphaned message types
- 3 .bak files (3,039 lines of dead code)
- Stub functions that should be removed/implemented

### Phases
1. **Help Screen** (30min): Implement keyboard reference overlay
2. **Display Errors** (30min): Show CV generation errors to user
3. **Remove Dead Code** (1h): Delete .bak files, remove/use loadingRotators, remove orphaned messages
4. **Stub Cleanup** (30min): Remove or implement stub functions and unused CLI flags

### Impact
Improved user experience, cleaner codebase, better maintainability.

---

## Recommended Implementation Order

### Week 1: Critical Fixes (Must Have)
**Goal**: Make broken features functional

```
Day 1-2: Task 29 - ExportArtifact (4-6 hours)
Day 3-4: Task 30 - ConfigureSystem (4-5 hours)
```

**Deliverables**:
- Users can export CVs/events/artifacts to files
- Users can configure system settings (persisted)
- 2 completely broken intents are now functional

---

### Week 2: High Priority (Should Have)
**Goal**: Complete partially implemented features

```
Day 1: Task 31 - CaptureEvent Modals (3-4 hours)
Day 2: Task 32 - BrowseTimeline Filters (3-4 hours)
Day 3: Task 33 - Audience Filtering (2-3 hours)
```

**Deliverables**:
- Event review and editing works properly
- Timeline filtering and search functional
- CVs tailored to target audience

---

### Week 3: Polish (Nice to Have)
**Goal**: Cleanup and consistency

```
Day 1: Task 34 - Consistency & Cleanup (2-3 hours)
```

**Deliverables**:
- Help screen functional
- Dead code removed
- Cleaner, more maintainable codebase

---

## Success Criteria

### After Tasks 23-24 (Critical)
- [ ] Users can export artifacts to real files
- [ ] Users can configure and persist settings
- [ ] No completely non-functional features
- [ ] All tests still passing (2,078/2,078)

### After Tasks 25-27 (High Priority)
- [ ] Event review and editing works end-to-end
- [ ] Timeline can be filtered and searched
- [ ] CVs are tailored to audience
- [ ] All advertised features work
- [ ] All tests still passing

### After Task 34 (Polish)
- [ ] Help screen available from any screen
- [ ] No dead code (0 .bak files)
- [ ] Error messages displayed to users
- [ ] Zero staticcheck warnings
- [ ] Production ready with all features functional

---

## Testing Standards (All Tasks)

Each task follows TDD (Red-Green-Refactor):

1. **Red**: Write failing test for the fix
2. **Green**: Implement minimal code to pass
3. **Refactor**: Clean up implementation
4. **Verify**: Run full test suite

**Commands**:
```bash
# Before starting
make check-compliance

# After each phase
go test ./... -v
go test -race ./...

# Before finishing
make check-compliance
```

---

## Related Documentation

### Task Execution
- [`docs/rules/master-task-prompt.md`](../docs/rules/master-task-prompt.md) - Complete 5-phase workflow
- [`docs/rules/TASK_QUICK_REF.md`](../docs/rules/TASK_QUICK_REF.md) - Quick reference checklist
- [`docs/rules/process-task-list.md`](../docs/rules/process-task-list.md) - Task processing rules

### Standards
- [`docs/TUI_STANDARDS.md`](../docs/TUI_STANDARDS.md) - TUI design standards
- [`docs/rules/go-guidelines.md`](../docs/rules/go-guidelines.md) - Go coding standards
- [`docs/rules/atomic-commits.md`](../docs/rules/atomic-commits.md) - Commit standards

### Architecture
- [`docs/TUI_INTENT_DIAGRAM.md`](../docs/TUI_INTENT_DIAGRAM.md) - Intent architecture
- [`docs/TUI_DEVELOPER_GUIDE.md`](../docs/TUI_DEVELOPER_GUIDE.md) - TUI development
- [`AGENTS.md`](../AGENTS.md) - Project overview

---

## Audit Methodology

The audit was conducted on 2026-01-08 using the following approach:

1. **Intent Layer Exploration**: Examined all 5 primary intents (CaptureEvent, BrowseTimeline, GenerateCV, ExportArtifact, ConfigureSystem)
2. **Service Layer Analysis**: Reviewed CV generation, bullet generation, and data processing services
3. **Repository Layer Check**: Verified CRUD operations and data access patterns
4. **Domain Layer Review**: Checked domain models for incomplete validation/methods
5. **App/Router Layer**: Examined global context, routing, and message handling
6. **Component/Model Layer**: Reviewed TUI components and legacy models
7. **CLI Commands**: Checked main entry points and command-line flags

**Tools Used**:
- Codebase exploration agents (thorough mode)
- Pattern matching for TODO/FIXME/HACK/stub keywords
- Manual code review of key files
- Test coverage analysis

**Findings**:
- 90+ partially implemented features
- 6 major categories of issues (2 critical, 4 high priority)
- No security vulnerabilities identified
- Architecture is sound, implementation incomplete

---

## Notes for Implementers

### Before Starting
1. Read the task file thoroughly
2. Run `make check-compliance` to ensure clean baseline
3. Create a feature branch for the task
4. Review related documentation

### During Implementation
1. Follow TDD workflow (Red-Green-Refactor)
2. One phase at a time
3. Test after each phase
4. Commit atomically (one logical change per commit)
5. Use AI attribution in commits (automated via git hooks)

### Before Finishing
1. Run full test suite
2. Run with race detector
3. Run `make check-compliance`
4. Verify all acceptance criteria met
5. Update AGENTS.md if needed

---

## Questions or Issues?

If you encounter any issues or have questions about these tasks:

1. **Check Documentation**: Review the referenced docs in each task file
2. **Review Audit Notes**: This summary and individual task files have detailed context
3. **Ask for Clarification**: Open an issue or discussion if unclear
4. **Propose Changes**: Task files can be updated if requirements change

---

**Last Updated**: 2026-01-08
**Created By**: AI Assistant (via OpenCode)
**Audit Session**: 2026-01-08 Comprehensive Codebase Audit
**Next Task Number**: Task 29
