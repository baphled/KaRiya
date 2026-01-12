# Escape Key Navigation - Comprehensive Audit Summary

**Date**: 2026-01-12  
**Bug**: bug-001-escape-key-navigation  
**Status**: Audit Complete, Ready for Fix Implementation

---

## Executive Summary

**Initial Report**: "Escape key not working application-wide"

**Audit Finding**: Issue isolated to **CaptureEvent intent only** (4 functions affected)

**Scope Clarification**:
- ✅ **9 of 10 intents CORRECT**: Already follow proper pattern
- ❌ **1 of 10 intents BROKEN**: CaptureEvent needs fixes

**Impact**: Lower than initially thought - isolated fix required

---

## Audit Methodology

### Automated Script

**Created**: `bugs/escape-key-audit.sh`

**Process**:
1. Parse all `*_intent.go` files
2. Find `.Update(msg)` delegation calls
3. Check if `HandleGlobalKeys` called BEFORE delegation
4. Report findings per intent

**Output**: `bugs/audit-results.txt`

### Manual Verification

For each affected function:
1. Read source code (lines 290-450)
2. Trace message flow
3. Identify exact delegation order
4. Confirm global key handling

---

## Detailed Findings

### ✅ CORRECT Intents (9 of 10)

These intents already implement the correct pattern (global keys checked BEFORE delegation):

| Intent | File | Pattern | Status |
|--------|------|---------|--------|
| **BrowseTimeline** | `browse_timeline_intent.go` | No problematic delegation | ✅ |
| **BulkOperations** | `bulk_operations_intent.go` | No problematic delegation | ✅ |
| **BurstManagement** | `burst_management_intent.go:620` | HandleGlobalKeys BEFORE editModal | ✅ |
| **ConfigureSystem** | `configure_system_intent.go:60` | HandleGlobalKeys BEFORE model | ✅ |
| **ExportArtifact** | `export_artifact_intent.go:49` | HandleGlobalKeys BEFORE model | ✅ |
| **FactManagement** | `fact_management_intent.go:448` | HandleGlobalKeys BEFORE editModal | ✅ |
| **GenerateCV** | `generate_cv_intent.go:339` | HandleGlobalKeys BEFORE viewport | ✅ |
| **ImportWizard** | `import_wizard_intent.go` | No problematic delegation | ✅ |
| **MetadataEditor** | `metadata_editor_intent.go` | No problematic delegation | ✅ |

**Example Correct Pattern** (from BurstManagement):

```go
func (i *BurstManagementIntent) updateEditView(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // ✅ Handle global keys FIRST
        switch HandleGlobalKeys(msg) {
        case KeyQuit:
            return tea.Quit
        case KeyHelp:
            i.ToggleHelp()
            return nil
        case KeyBack:
            i.state.currentState = BurstStateDetail
            i.state.editError = nil
            return nil
        }
    }

    // ✅ THEN delegate to modal
    cmd := i.state.editModal.Update(msg)
    // ... rest of function
}
```

---

### ❌ BROKEN Intent (1 of 10)

**Intent**: CaptureEvent  
**File**: `internal/cli/intents/capture_event_intent.go`  
**Functions Affected**: 4

#### Issue 1: updateCaptureForm() - Line 298

**Problem**:
```go
func (i *CaptureEventIntent) updateCaptureForm(msg tea.Msg) tea.Cmd {
    // ❌ WRONG: Delegate to form FIRST
    _, formCmd := i.state.captureForm.Update(msg)

    // ❌ TOO LATE: Check global keys AFTER
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch HandleGlobalKeys(msg) {
        case KeyBack:
            i.state.currentState = CaptureStateChooseStrategy
            return nil
        }
    }
    // ...
}
```

**Impact**: User trapped in form, escape key consumed by Huh library

**Fix**: Check `HandleGlobalKeys(msg)` BEFORE calling `form.Update(msg)`

---

#### Issue 2: updateReviewInferredEvent() - MetadataModal (Line 383)

**Problem**:
```go
func (i *CaptureEventIntent) updateReviewInferredEvent(msg tea.Msg) tea.Cmd {
    switch i.state.reviewState.EditingMode {
    case EditingModeMetadata:
        if i.state.reviewState.metadataModal != nil {
            // ❌ WRONG: Delegate to modal immediately
            modal, cmd := i.state.reviewState.metadataModal.Update(msg)
            // ... handle completion
            return cmd
        }
    }
    // ... global keys checked LATER (line 440)
}
```

**Impact**: User trapped in metadata editor modal

**Fix**: Check `HandleGlobalKeys(msg)` BEFORE modal delegation

---

#### Issue 3: updateReviewInferredEvent() - BurstModal (Line 403)

**Problem**:
```go
case EditingModeBursts:
    if i.state.reviewState.burstModal != nil {
        // ❌ WRONG: Delegate to modal immediately
        modal, cmd := i.state.reviewState.burstModal.Update(msg)
        
        // ⚠️ WORKAROUND: Manual escape check
        if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
            i.state.reviewState.burstModal = nil
            i.state.reviewState.EditingMode = EditingModeNone
        }
        return cmd
    }
```

**Impact**: 
- Inconsistent pattern (manual escape check)
- Modal may still consume escape before check
- Doesn't follow HandleGlobalKeys pattern

**Fix**: Check `HandleGlobalKeys(msg)` BEFORE modal delegation, remove manual check

---

#### Issue 4: updateReviewInferredEvent() - FactModal (Line 418)

**Problem**:
```go
case EditingModeFacts:
    if i.state.reviewState.factModal != nil {
        // ❌ WRONG: Delegate to modal immediately
        modal, cmd := i.state.reviewState.factModal.Update(msg)
        
        // Check completion
        if i.state.reviewState.factModal.IsSubmitted() {
            // ... handle submission
        }
        return cmd
    }
```

**Impact**: User trapped in fact editor modal

**Fix**: Check `HandleGlobalKeys(msg)` BEFORE modal delegation

---

## Root Cause Analysis

### Why This Happened

**1. Inconsistent Implementation**:
- 9 intents correctly implemented pattern
- 1 intent (CaptureEvent) implemented different pattern
- Likely due to:
  - Earlier implementation before pattern established
  - Different developer or session
  - Copy-paste from outdated example

**2. Test Gap**:
- Unit tests call `intent.Update()` directly
- Bypasses the delegation flow
- Tests pass even though real app fails
- Missing E2E tests for form interaction

**3. Pattern Not Enforced**:
- No linting rule to catch this
- No code review checklist
- Pattern documented but not enforced

### Prevention Measures

**Immediate**:
1. Fix CaptureEvent intent (4 functions)
2. Add E2E tests for form escape behavior
3. Update unit tests to reflect real flow

**Short Term**:
1. Create developer guide: `ESCAPE_KEY_HANDLING.md`
2. Add to TUI_STANDARDS.md as critical pattern
3. Create code review checklist

**Long Term**:
1. Consider custom linter rule
2. Create intent template with correct pattern
3. Automated pattern verification in CI

---

## Fix Plan

### Phase 1: Fix CaptureEvent Intent (2-3 hours)

**Files to Modify**: 1 file, 4 functions

**File**: `internal/cli/intents/capture_event_intent.go`

**Changes**:

1. **Function**: `updateCaptureForm()` (Line 298)
   - Move `HandleGlobalKeys(msg)` check to START of function
   - Move `form.Update(msg)` AFTER global key check
   - Estimated: 30 min

2. **Function**: `updateReviewInferredEvent()` - Add global key check at start
   - Check `HandleGlobalKeys(msg)` BEFORE modal routing
   - Handle escape to close active modal
   - Update MetadataModal case (Line 383) - 20 min
   - Update BurstModal case (Line 403) - 20 min (remove manual check)
   - Update FactModal case (Line 418) - 20 min

**Total Coding Time**: ~1.5 hours

### Phase 2: Add Unit Tests (1 hour)

**File**: `internal/cli/intents/capture_event_escape_test.go`

**New Tests**:
```go
Describe("Form State - With Active Huh Form", func() {
    It("should handle escape even when form has focus", func() { ... })
})

Describe("Review State - With Active Modals", func() {
    It("should handle escape in metadata modal", func() { ... })
    It("should handle escape in burst modal", func() { ... })
    It("should handle escape in fact modal", func() { ... })
})
```

### Phase 3: Add E2E Tests (2 hours)

**New File**: `internal/testutil/e2e/capture_event_escape_e2e_test.go`

**Tests**:
- Navigate into form, press escape, verify back to strategy
- Enter metadata modal, press escape, verify modal closes
- Enter burst modal, press escape, verify modal closes
- Enter fact modal, press escape, verify modal closes
- Complete end-to-end workflow with escape navigation

### Phase 4: Documentation (30 min)

**New File**: `docs/development/ESCAPE_KEY_HANDLING.md`

**Content**:
- Correct pattern explanation
- Code examples (correct vs incorrect)
- Why it matters
- How to test
- Common mistakes

**Update**: `docs/TUI_STANDARDS.md`
- Add prominent warning box
- Reference escape key handling guide
- Mark as CRITICAL pattern

### Phase 5: Manual Verification (30 min)

**Test Matrix**:
- [ ] Start app, enter capture form, press escape → back to strategy ✅
- [ ] From strategy, press escape → back to main menu ✅
- [ ] In review, press 'e' for metadata, press escape → modal closes ✅
- [ ] In review, press 'b' for bursts, press escape → modal closes ✅
- [ ] In review, press 'f' for facts, press escape → modal closes ✅
- [ ] Verify 'm' key works from all states ✅
- [ ] Verify 'q' key works from all states ✅

**Total Estimated Time**: 6-7 hours

---

## Risk Assessment

### Low Risk Items ✅

- Only 1 intent affected (isolated change)
- Pattern proven in 9 other intents
- No library modifications needed
- Preserves all existing functionality
- Can verify with manual testing

### Medium Risk Items ⚠️

- Modals might have internal escape handling
- Form might rely on receiving all keys
- Need thorough testing of form submission flow

### Mitigation

- Test each function individually
- Verify form submission still works
- Check modal completion detection
- Run full test suite after each fix
- Manual test each screen

---

## Success Criteria

### Code Changes ✅
- [ ] All 4 functions fixed in `capture_event_intent.go`
- [ ] HandleGlobalKeys checked BEFORE delegation in all cases
- [ ] Manual escape check removed from BurstModal
- [ ] Code follows pattern used in other 9 intents

### Testing ✅
- [ ] Unit tests updated (4+ new tests)
- [ ] E2E tests created (5+ new tests)
- [ ] All tests passing (no regressions)
- [ ] Manual test matrix 100% complete

### User Experience ✅
- [ ] Escape works from form screen
- [ ] Escape works from all modal screens
- [ ] 'm' key works from all screens
- [ ] 'q' key works from all screens
- [ ] No user can get trapped in any screen

### Documentation ✅
- [ ] ESCAPE_KEY_HANDLING.md created
- [ ] TUI_STANDARDS.md updated
- [ ] Bug report updated with resolution
- [ ] Code comments added to fixed functions

---

## Conclusion

**Initial Perception**: Application-wide problem requiring extensive fixes

**Actual Reality**: Isolated to CaptureEvent intent, 4 functions need fixing

**Impact**: Lower scope, faster fix, less risk

**Confidence**: High - Pattern proven in 9 other intents, surgical fix possible

**Recommendation**: Proceed with Phase 1 fix immediately, then validate with tests

---

**Audit Completed**: 2026-01-12 03:45  
**Auditor**: OpenCode AI Assistant  
**Files Reviewed**: 10 intent files, ~5,000 lines of code  
**Findings**: 4 issues in 1 intent, 9 intents verified correct
