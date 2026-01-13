# Task 21: TUI Architecture Refactoring - Intent/Screen Pattern

## Overview
- **Goal**: Refactor TUI from monolithic intents to Intent/Screen/Component architecture
- **Time Estimate**: 6-8 weeks
- **Prerequisites**: Stable main branch, all tests passing

## Session Contract Acknowledgment
- [ ] Ran `make session-start` and it passed
- [ ] Acknowledge and commit to following all workflow rules
- [ ] Token count: _____ (must be < 50k to start)

---

## Phase 1: Foundation (Week 1)

### 1.1 Create screens/ Directory Structure

**Files to Create**:
- [ ] `internal/cli/screens/contract.go` - Screen interface, ScreenResult types
- [ ] `internal/cli/screens/base/base_screen.go` - BaseScreen with UI capabilities
- [ ] `internal/cli/screens/base/select_screen.go` - BaseSelectScreen[T]
- [ ] `internal/cli/screens/base/list_screen.go` - BaseListScreen[T]

**TDD Checklist**:
- [ ] RED: Write tests for Screen interface contract
- [ ] RED: Write tests for BaseScreen methods
- [ ] GREEN: Implement contract.go
- [ ] GREEN: Implement base_screen.go
- [ ] REFACTOR: Extract common patterns

**Acceptance Criteria**:
- [ ] Screen interface compiles
- [ ] BaseScreen has SetTerminalInfo, SetTheme, CreateView methods
- [ ] ScreenResult types defined (NavigateResult, CancelResult, SubmitResult, ErrorResult)
- [ ] Tests pass with race detector

### 1.2 Create Base Select Screen

**Files to Create**:
- [ ] `internal/cli/screens/base/select_screen.go`
- [ ] `internal/cli/screens/base/select_screen_test.go`

**TDD Checklist**:
- [ ] RED: Test navigation (up/down, vim keys)
- [ ] RED: Test selection (enter returns NavigateResult)
- [ ] RED: Test cancellation (esc returns CancelResult)
- [ ] GREEN: Implement BaseSelectScreen[T]
- [ ] REFACTOR: Extract itemRenderer pattern

**Acceptance Criteria**:
- [ ] Generic type parameter works with any type
- [ ] Navigation with ↑/↓/j/k/g/G
- [ ] Selection returns NavigateResult with data
- [ ] Escape returns CancelResult
- [ ] View renders items with selection indicator

---

## Phase 2: Pilot Intent - GenerateCV (Week 2-3)

### 2.1 Create CV-Specific Screens

**Files to Create**:
- [ ] `internal/cli/screens/cv/profile_select.go` - CVProfileSelect
- [ ] `internal/cli/screens/cv/audience_select.go` - CVAudienceSelect
- [ ] `internal/cli/screens/cv/role_emphasis_select.go` - CVRoleEmphasisSelect
- [ ] `internal/cli/screens/cv/length_format_select.go` - CVLengthFormatSelect
- [ ] `internal/cli/screens/cv/generating.go` - CVGenerating
- [ ] `internal/cli/screens/cv/preview.go` - CVPreview

**TDD Checklist for Each Screen**:
- [ ] RED: Test screen creation with config
- [ ] RED: Test Update returns correct ScreenResult
- [ ] RED: Test View renders domain content
- [ ] GREEN: Implement screen
- [ ] REFACTOR: Share common patterns

### 2.2 Refactor GenerateCVIntent

**Files to Modify**:
- [ ] `internal/cli/intents/generate_cv.go` - Simplify to orchestration only
- [ ] `internal/cli/intents/generate_cv_intent.go` - Use screens instead of views

**Changes**:
```go
// Before: ~900 lines with updateXxx() and viewXxx() for each state
// After: ~200 lines with transitionTo() and handleScreenResult()

type GenerateCVIntent struct {
    state        CVState
    context      *GenerateCVContext
    result       *IntentResult[*GenerateCVResult]
    activeScreen screens.Screen
    // ... collected data fields
}

func (i *GenerateCVIntent) transitionTo(state CVState) {
    switch state {
    case CVSelectProfile:
        i.activeScreen = cv.NewCVProfileSelect(...)
    // ...
    }
}

func (i *GenerateCVIntent) Update(msg tea.Msg) tea.Cmd {
    cmd, result := i.activeScreen.Update(msg)
    if result != nil {
        return i.handleScreenResult(result)
    }
    return cmd
}

func (i *GenerateCVIntent) View() string {
    return i.activeScreen.View()
}
```

**Acceptance Criteria**:
- [ ] Intent reduced to <300 lines
- [ ] All existing tests still pass
- [ ] Workflow unchanged from user perspective
- [ ] No regressions in CV generation

---

## Phase 3: Build Screen Library (Week 3-4)

### 3.1 Base Screens

**Files to Create**:
- [ ] `internal/cli/screens/base/form_screen.go` - BaseFormScreen
- [ ] `internal/cli/screens/base/detail_screen.go` - BaseDetailScreen
- [ ] `internal/cli/screens/base/confirm_screen.go` - BaseConfirmScreen
- [ ] `internal/cli/screens/base/progress_screen.go` - BaseProgressScreen

### 3.2 Skills Screens

**Files to Create**:
- [ ] `internal/cli/screens/skills/list.go` - SkillsList
- [ ] `internal/cli/screens/skills/detail.go` - SkillDetail
- [ ] `internal/cli/screens/skills/form.go` - SkillForm
- [ ] `internal/cli/screens/skills/delete.go` - SkillDeleteConfirm
- [ ] `internal/cli/screens/skills/filter.go` - SkillFilterMenu
- [ ] `internal/cli/screens/skills/sort.go` - SkillSortMenu

### 3.3 Timeline Screens

**Files to Create**:
- [ ] `internal/cli/screens/timeline/event_list.go` - TimelineEventList
- [ ] `internal/cli/screens/timeline/event_detail.go` - TimelineEventDetail

---

## Phase 4: Intent Migration (Week 5-7)

Migrate remaining intents one at a time:

### 4.1 ManageSkillsIntent
- [ ] Create skills screens (Phase 3.2)
- [ ] Refactor ManageSkillsIntent to use screens
- [ ] Verify all tests pass
- [ ] Expected reduction: 1,647 → ~250 lines

### 4.2 BrowseTimelineIntent
- [ ] Create timeline screens (Phase 3.3)
- [ ] Refactor BrowseTimelineIntent
- [ ] Verify all tests pass

### 4.3 CaptureEventIntent
- [ ] Create capture screens
- [ ] Refactor CaptureEventIntent
- [ ] Verify all tests pass

### 4.4 ExportArtifactIntent
- [ ] Create export screens
- [ ] Refactor ExportArtifactIntent
- [ ] Verify all tests pass

### 4.5 Remaining Intents
- [ ] BurstManagementIntent
- [ ] FactManagementIntent
- [ ] ImportWizardIntent
- [ ] BulkOperationsIntent
- [ ] ConfigureSystemIntent

---

## Phase 5: Cleanup & Documentation (Week 8)

### 5.1 Deprecate models/ Directory
- [ ] Migrate remaining wrappers to screens/
- [ ] Remove deprecated files
- [ ] Update imports

### 5.2 Update Documentation
- [ ] Update TUI_DEVELOPER_GUIDE.md
- [ ] Update TUI_INTENT_DIAGRAM.md
- [ ] Create screens/ package documentation

### 5.3 Final Verification
- [ ] All 2,000+ tests pass
- [ ] No race conditions
- [ ] Coverage maintained >85%
- [ ] Staticcheck warnings: 0

---

## Pre-Commit Checklist (BEFORE EACH COMMIT)
- [ ] `make check-compliance` passes
- [ ] Use `make ai-commit MSG="type(scope): description"` for AI-generated code
- [ ] Commit is atomic (ONE logical change)

## Acceptance Criteria

### Code Quality
- [ ] All tests pass (2,000+)
- [ ] Coverage >85%
- [ ] Zero staticcheck warnings
- [ ] Zero race conditions

### Architecture
- [ ] Intent files reduced by 70%+ average
- [ ] Screens are reusable across intents
- [ ] Clear separation: Intent (orchestration) / Screen (view) / Component (primitive)

### User Experience
- [ ] Zero regressions in workflows
- [ ] Keyboard shortcuts unchanged
- [ ] Navigation patterns preserved

### Documentation
- [ ] Architecture doc updated
- [ ] Naming conventions documented
- [ ] Screen patterns documented

---

## Rollback Plan

Each phase can be rolled back independently:
1. **Phase 1**: Delete screens/ directory
2. **Phase 2-4**: Revert intent changes, screens remain
3. **Phase 5**: No rollback needed (documentation only)

Git tags at each phase completion:
- `v2.0.0-arch-phase1`
- `v2.0.0-arch-phase2`
- `v2.0.0-arch-phase3`
- `v2.0.0-arch-phase4`
- `v2.0.0-arch-complete`

---

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Test failures during migration | Medium | High | Migrate one intent at a time, full test suite between |
| Performance regression | Low | Medium | Benchmark before/after each phase |
| User workflow changes | Low | High | Manual testing of all workflows |
| Scope creep | Medium | Medium | Strict phase boundaries, no feature additions |
