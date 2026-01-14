---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Atomic Tasks: Phase 3 Completion (Tasks 9-15)

**Status**: 80% Complete - Fix 8 failing tests, complete integration, then finalize documentation

**Master Task Prompt**: Follow strictly - Red→Green→Refactor, atomic commits, token efficiency

---

## Tasks Checklist

### ❌ BLOCKING: Fix Failing Tests (8 failures)

#### Task A.1: Fix metadata_review_test.go Header Rendering
- [ ] **RED**: Verify test expectation for header rendering in View()
- [ ] **GREEN**: Fix MetadataReviewModel.View() to render header correctly
- [ ] **REFACTOR**: Ensure code clarity and DRY principles
- [ ] **COMMIT**: `fix(cli): fix metadata review header rendering in View()`
- [ ] **VERIFY**: Single test passes, no new failures
- **Token Check**: < 50k?

#### Task A.2: Fix form_test.go Capture Mode Integration (5 failures)
- [ ] **RED**: Review failing tests - capture mode integration with service
- [ ] **GREEN**: Verify FormModel passes correct capture mode to service
- [ ] **GREEN**: Verify service receives mode parameter in CaptureEvent()
- [ ] **REFACTOR**: Check mode validation and error handling
- [ ] **COMMIT**: `fix(cli): fix form capture mode integration with service`
- [ ] **VERIFY**: All 5 form tests pass, no new failures
- **Token Check**: < 50k?

#### Task A.3: Fix quality_indicator_test.go Icon Rendering
- [ ] **RED**: Verify RenderCompact test for Incomplete level icon
- [ ] **GREEN**: Fix QualityIndicator.RenderCompact() icon selection
- [ ] **REFACTOR**: Ensure icon mappings are consistent across methods
- [ ] **COMMIT**: `fix(cli): fix quality indicator icon rendering for Incomplete level`
- [ ] **VERIFY**: Single test passes, no new failures
- **Token Check**: < 50k?

---

### ✅ COMPLETED: Phase 3 Core Tasks (9-11)

**Status**: Models and service methods implemented, integration tests passing

- ✅ Task 9.0: Bulk Operations Model (33/33 tests)
- ✅ Task 10.0 Part 1: Bulk Operations Model Integration (3/3 integration tests)
- ✅ Task 11.0: Bulk Service Enhancement (9/9 tests)

---

### ⏳ IN PROGRESS: Task 10.0 Part 2 - Bulk Operations App Integration

#### Task B.1: Add BulkOperationsScreen to App Navigation
- [ ] **RED**: Write test for app.go to handle BulkOperationsScreen navigation
- [ ] **GREEN**: Add BulkOperationsScreen constant to Screen type
- [ ] **GREEN**: Add bulkOperationsModel field to Model struct
- [ ] **GREEN**: Add BulkOperationsScreen case to View() method
- [ ] **GREEN**: Add BulkOperationsScreen case to Update() method with message handling
- [ ] **REFACTOR**: Ensure navigation state management is consistent
- [ ] **COMMIT**: `feat(cli): add bulk operations screen to app navigation`
- [ ] **VERIFY**: New tests pass, existing tests still pass
- **Token Check**: < 50k?

#### Task B.2: Implement Bulk Operations Navigation from Metadata Review
- [ ] **RED**: Write test for metadata review to trigger bulk operations
- [ ] **GREEN**: Add 'b' keyboard shortcut in metadata review to start bulk operations
- [ ] **GREEN**: Pass selected events or full event list to bulk operations model
- [ ] **GREEN**: Initialize BulkOperationsModel with correct event list
- [ ] **REFACTOR**: Ensure message passing is clean and type-safe
- [ ] **COMMIT**: `feat(cli): add navigation from metadata review to bulk operations`
- [ ] **VERIFY**: Navigation tests pass, app integration tests pass
- **Token Check**: < 50k?

#### Task B.3: Implement Bulk Operations Return to Metadata Review
- [ ] **RED**: Write test for bulk operations completion to return to review
- [ ] **GREEN**: Handle submit in app.go to apply bulk changes via service
- [ ] **GREEN**: Handle cancel in app.go to discard changes and return to review
- [ ] **GREEN**: Refresh metadata review list after bulk operation completes
- [ ] **REFACTOR**: Ensure error handling for bulk update failures
- [ ] **COMMIT**: `feat(cli): implement bulk operations completion and return flow`
- [ ] **VERIFY**: All bulk operation tests pass, metadata review tests pass
- **Token Check**: < 50k?

---

### ⏳ PENDING: Phase 4 - Integration with Existing Features

#### Task C.1: Task 12.0 - CSV Import Integration (Part 1)
- [ ] **RED**: Write test for import_review.go to show metadata review after import
- [ ] **GREEN**: Modify ImportReviewModel to trigger metadata review on completion
- [ ] **GREEN**: Pass imported events to metadata review screen
- [ ] **GREEN**: Add navigation handler from import review to metadata review
- [ ] **REFACTOR**: Ensure error handling for import failures
- [ ] **COMMIT**: `feat(cli): trigger metadata review after CSV import`
- [ ] **VERIFY**: Import integration tests pass, metadata review tests pass
- **Token Check**: < 50k?

#### Task C.2: Task 12.0 - CSV Import Integration (Part 2)
- [ ] **RED**: Write test for metadata review to show field origins (CSV vs. default)
- [ ] **GREEN**: Add origin tracking to metadata review display
- [ ] **GREEN**: Display which fields came from CSV vs. defaults
- [ ] **GREEN**: Add visual indicators for imported vs. enriched fields
- [ ] **REFACTOR**: Ensure styling consistency with existing components
- [ ] **COMMIT**: `feat(cli): display field origins in metadata review`
- [ ] **VERIFY**: Field origin tests pass, metadata review rendering tests pass
- **Token Check**: < 50k?

#### Task C.3: Task 13.0 - Manual Capture Integration (Part 1)
- [ ] **RED**: Write test for form.go to show quick metadata enrichment after capture
- [ ] **GREEN**: Modify FormModel success flow to offer metadata enrichment
- [ ] **GREEN**: Add option to edit metadata before saving event
- [ ] **GREEN**: Pass captured event to metadata editor on selection
- [ ] **REFACTOR**: Ensure form state is properly reset after enrichment
- [ ] **COMMIT**: `feat(cli): add quick metadata enrichment after manual capture`
- [ ] **VERIFY**: Form integration tests pass, metadata editor tests pass
- **Token Check**: < 50k?

#### Task C.4: Task 13.0 - Manual Capture Integration (Part 2)
- [ ] **RED**: Write test for success.go to include metadata review option
- [ ] **GREEN**: Add "Review Metadata" button to success screen
- [ ] **GREEN**: Implement navigation to metadata review from success screen
- [ ] **GREEN**: Pass captured event to metadata review for enrichment
- [ ] **REFACTOR**: Ensure button placement and styling are consistent
- [ ] **COMMIT**: `feat(cli): add metadata review option to success screen`
- [ ] **VERIFY**: Success screen tests pass, navigation tests pass
- **Token Check**: < 50k?

---

### ⏳ PENDING: Phase 5 - Testing & Documentation

#### Task D.1: Task 14.0 - End-to-End Testing
- [ ] **RED**: Write comprehensive e2e test for complete metadata workflow
- [ ] **GREEN**: Test: import → metadata review → bulk operations → save
- [ ] **GREEN**: Test: capture → metadata enrichment → metadata review
- [ ] **GREEN**: Test: metadata review → individual edit → save → list update
- [ ] **REFACTOR**: Extract common test helpers for reuse
- [ ] **COMMIT**: `test(cli): add comprehensive end-to-end metadata workflow tests`
- [ ] **VERIFY**: All e2e tests pass, coverage ≥ 80%
- **Token Check**: < 50k?

#### Task D.2: Task 14.0 - Edge Cases & Error Handling
- [ ] **RED**: Write tests for edge cases (empty lists, large datasets, etc.)
- [ ] **GREEN**: Test: empty event list in metadata review
- [ ] **GREEN**: Test: bulk operations with no selections
- [ ] **GREEN**: Test: validation errors in metadata editor
- [ ] **GREEN**: Test: undo/revert functionality
- [ ] **REFACTOR**: Ensure all error paths return helpful messages
- [ ] **COMMIT**: `test(cli): add edge case and error handling tests`
- [ ] **VERIFY**: All edge case tests pass, no regressions
- **Token Check**: < 50k?

#### Task D.3: Task 14.0 - Race Condition Testing
- [ ] **RUN**: `go test -race ./...` for all packages
- [ ] **FIX**: Any race conditions detected
- [ ] **VERIFY**: All tests pass with race detector enabled
- [ ] **COMMIT**: `test(cli): verify no race conditions in metadata features`
- [ ] **Token Check**: < 50k?

#### Task D.4: Task 15.0 - Documentation Updates
- [ ] **UPDATE**: README.md with metadata review workflow description
- [ ] **UPDATE**: CLI_GUIDE.md with keyboard shortcuts for metadata features
- [ ] **UPDATE**: TROUBLESHOOTING.md with common metadata issues
- [ ] **CREATE**: Examples of metadata workflows in docs/
- [ ] **UPDATE**: CHANGELOG.md with feature descriptions
- [ ] **COMMIT**: `docs: add metadata feature documentation`
- [ ] **VERIFY**: All docs are clear and complete
- **Token Check**: < 50k?

---

## Execution Order

**Priority 1 (BLOCKING - Fix Tests)**:
1. Task A.1: Fix header rendering
2. Task A.2: Fix capture mode integration
3. Task A.3: Fix icon rendering

**Priority 2 (Complete Phase 3)**:
4. Task B.1: Add BulkOperationsScreen to app
5. Task B.2: Implement bulk operations navigation
6. Task B.3: Implement bulk operations completion

**Priority 3 (Phase 4 Integration)**:
7. Task C.1: CSV import integration (part 1)
8. Task C.2: CSV import integration (part 2)
9. Task C.3: Manual capture integration (part 1)
10. Task C.4: Manual capture integration (part 2)

**Priority 4 (Phase 5 Testing & Docs)**:
11. Task D.1: End-to-end testing
12. Task D.2: Edge cases & error handling
13. Task D.3: Race condition testing
14. Task D.4: Documentation updates

---

## Master Task Prompt Compliance

### For Each Task:

**Phase 1: Preparation**
- [ ] Token count checked (< 50k)
- [ ] `make check-compliance` passed
- [ ] Task is clearly understood
- [ ] Task is atomic (ONE change)
- [ ] Existing patterns reviewed

**Phase 2: Red-Green-Refactor**
- [ ] Test written and failing
- [ ] Implementation minimal and correct
- [ ] Code refactored if needed
- [ ] All tests passing
- [ ] No race conditions

**Phase 3: Compliance Verification**
- [ ] Code formatted (`make fmt`)
- [ ] No vet warnings (`make vet`)
- [ ] Coverage ≥ 80%
- [ ] All tests pass
- [ ] No race conditions

**Phase 4: Final Verification**
- [ ] All commits are atomic
- [ ] Commit messages follow conventions
- [ ] Messages explain WHY
- [ ] No generated files
- [ ] No debug code

**Phase 5: Task Completion**
- [ ] Summary documented
- [ ] Next steps identified
- [ ] Token count reasonable
- [ ] Ready for handoff

---

## Key Rules

### Atomic Commits
- One logical change per commit
- Type: feat, fix, test, docs, refactor
- Scope: cli, models, service, domain, repo
- Message: "type(scope): description"
- Body: Explains WHY

### Code Quality
- Go idioms and conventions
- Interface-based design
- Proper error handling
- Structured logging
- SOLID principles

### Testing (TDD)
- RED: Write failing test first
- GREEN: Minimal implementation
- REFACTOR: Improve code
- One expectation per It block
- Descriptive test names

### Token Efficiency
- Use tools (view, grep, ls)
- Be concise and specific
- Batch operations
- Reference context
- Focus on deltas

---

## Test Coverage Targets

- **Overall**: ≥ 80%
- **Models**: ≥ 85%
- **Service**: 100%
- **Domain**: 100%
- **CLI**: ≥ 80%

---

## Success Criteria

- ✅ All 8 failing tests fixed
- ✅ Task 10.0 Part 2 complete (bulk app integration)
- ✅ Task 12.0 complete (CSV import integration)
- ✅ Task 13.0 complete (manual capture integration)
- ✅ Task 14.0 complete (comprehensive testing)
- ✅ Task 15.0 complete (documentation)
- ✅ All tests passing (305+ tests)
- ✅ Coverage ≥ 80%
- ✅ No race conditions
- ✅ All commits atomic

---

**Document Version**: 1.0
**Last Updated**: 2025-12-30
**Status**: Ready for execution
**Prepared By**: Development Assistant

