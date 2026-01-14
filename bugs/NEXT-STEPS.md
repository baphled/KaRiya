# Bug 001: Escape Key Navigation - Next Steps

**Date**: 2026-01-12  
**Status**: Investigation Complete ✅, Ready for Implementation  
**Bug Report**: `bugs/bug-001-escape-key-navigation.md`  
**Audit Summary**: `bugs/AUDIT-SUMMARY.md`

---

## Quick Summary

**What We Found**:
- ❌ Issue isolated to **CaptureEvent intent only** (4 functions)
- ✅ All other 9 intents already correct
- 📊 Lower scope than initially thought

**What Needs Fixing**:
1. `updateCaptureForm()` - Line 298
2. `updateReviewInferredEvent()` - MetadataModal (Line 383)
3. `updateReviewInferredEvent()` - BurstModal (Line 403)
4. `updateReviewInferredEvent()` - FactModal (Line 418)

**Estimated Time**: 6-7 hours total

---

## Implementation Phases

### Phase 1: Fix Code (2-3 hours)

**File**: `internal/cli/intents/capture_event_intent.go`

#### Fix 1: updateCaptureForm() (30 min)

**Current (Line 298-316)**:
```go
func (i *CaptureEventIntent) updateCaptureForm(msg tea.Msg) tea.Cmd {
    _, formCmd := i.state.captureForm.Update(msg)  // ❌ WRONG ORDER

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch HandleGlobalKeys(msg) {  // ❌ TOO LATE
```

**Fixed**:
```go
func (i *CaptureEventIntent) updateCaptureForm(msg tea.Msg) tea.Cmd {
    // ✅ Check global keys FIRST
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch HandleGlobalKeys(msg) {
        case KeyQuit:
            return tea.Quit
        case KeyHelp:
            i.ToggleHelp()
            return nil
        case KeyBack:
            i.state.currentState = CaptureStateChooseStrategy
            return nil
        }
        
        switch msg.String() {
        case "ctrl+s":
            return i.state.captureForm.SubmitForm()
        }
    }

    // ✅ NOW delegate to form
    _, formCmd := i.state.captureForm.Update(msg)

    // ✅ Check completion messages
    switch msg := msg.(type) {
    case models.SubmitMsg:
        // ... existing handling
    }

    return formCmd
}
```

#### Fix 2-4: updateReviewInferredEvent() (1 hour)

**Add at START of function** (before modal routing):

```go
func (i *CaptureEventIntent) updateReviewInferredEvent(msg tea.Msg) tea.Cmd {
    // ✅ Check global keys FIRST
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch HandleGlobalKeys(msg) {
        case KeyQuit:
            return tea.Quit
        case KeyHelp:
            i.ToggleHelp()
            return nil
        case KeyBack:
            // If modal is active, close it
            if i.state.reviewState.EditingMode != EditingModeNone {
                i.state.reviewState.metadataModal = nil
                i.state.reviewState.burstModal = nil
                i.state.reviewState.factModal = nil
                i.state.reviewState.EditingMode = EditingModeNone
                return nil
            }
            // Otherwise go back to form
            i.state.currentState = CaptureStateForm
            return nil
        }
    }

    // NOW check modal routing (existing code continues)
    switch i.state.reviewState.EditingMode {
    case EditingModeMetadata:
        if i.state.reviewState.metadataModal != nil {
            modal, cmd := i.state.reviewState.metadataModal.Update(msg)
            // ... rest of existing code
```

**Also Remove**: Manual escape check in BurstModal case (Line 409) - no longer needed

---

### Phase 2: Add Unit Tests (1 hour)

**File**: `internal/cli/intents/capture_event_escape_test.go`

**Add New Tests**:

```go
Describe("Form State - Escape with Huh Form Active", func() {
    BeforeEach(func() {
        intent.state.currentState = CaptureStateForm
        intent.state.captureForm.Init()  // Initialize form
    })

    It("should handle escape even when form has focus", func() {
        intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
        Expect(intent.state.currentState).To(Equal(CaptureStateChooseStrategy))
        Expect(intent.active).To(BeTrue())
    })
})

Describe("Review State - Escape with Active Modals", func() {
    BeforeEach(func() {
        intent.state.currentState = CaptureStateReview
        intent.state.reviewState.Event = &career.CareerEvent{...}
    })

    Context("MetadataModal active", func() {
        BeforeEach(func() {
            intent.state.reviewState.EditingMode = EditingModeMetadata
            intent.state.reviewState.metadataModal = models.NewMetadataEditorModelNew(...)
        })

        It("should close modal when escape pressed", func() {
            intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
            Expect(intent.state.reviewState.EditingMode).To(Equal(EditingModeNone))
            Expect(intent.state.reviewState.metadataModal).To(BeNil())
        })
    })

    // Similar tests for BurstModal and FactModal
})
```

---

### Phase 3: Add E2E Tests (2 hours)

**New File**: `internal/testutil/e2e/capture_event_escape_e2e_test.go`

```go
var _ = Describe("E2E - CaptureEvent Escape Navigation", func() {
    var env *e2e.TestEnv

    BeforeEach(func() {
        env = e2e.SetupWithMemory(GinkgoT())
    })

    AfterEach(func() {
        env.Cleanup()
    })

    It("should navigate back from form to strategy", func() {
        env.SelectIntentByName("capture_event")
        env.Confirm()  // Select Quick strategy
        env.AssertViewContains("Event")  // In form

        env.Cancel()  // Press Escape

        env.AssertViewContainsAny("Quick", "Manual", "Strategy")
        env.AssertViewNotContains("Event Text")
    })

    It("should close metadata modal with escape", func() {
        // Setup: Create event and enter review state
        // ... (test helper to populate and navigate to review)

        // Press 'e' to open metadata modal
        env.PressKeyRune('e')
        env.AssertViewContains("Edit Metadata")

        // Press Escape to close modal
        env.Cancel()
        env.AssertViewNotContains("Edit Metadata")
        env.AssertViewContains("Review")  // Still in review state
    })

    // Similar tests for burst and fact modals
})
```

---

### Phase 4: Documentation (30 min)

#### Create New Guide

**File**: `docs/development/ESCAPE_KEY_HANDLING.md`

**Content**:
- Explanation of correct pattern
- Code examples (correct vs incorrect)
- Why it matters for UX
- How to test
- Common mistakes
- Reference to TUI_STANDARDS.md

#### Update TUI Standards

**File**: `docs/TUI_STANDARDS.md`

**Add Prominent Warning**:

```markdown
### ⚠️ CRITICAL: Global Key Handling Order

**Pattern**: Check global keys BEFORE delegating to sub-components

**CORRECT**:
```go
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
    // ✅ Step 1: Check global keys FIRST
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch HandleGlobalKeys(msg) {
        case KeyBack:
            // Handle escape key
            return nil
        }
    }
    
    // ✅ Step 2: THEN delegate to sub-component
    subComponent.Update(msg)
}
```

**INCORRECT** (causes users to get trapped):
```go
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
    // ❌ WRONG: Delegate first
    subComponent.Update(msg)
    
    // ❌ TOO LATE: Global keys checked after
    switch HandleGlobalKeys(msg) {
        // This never runs - escape already consumed
    }
}
```

**Violation**: User gets trapped in screens, escape key doesn't work
**Reference**: See bug-001 and ESCAPE_KEY_HANDLING.md
```

---

### Phase 5: Manual Verification (30 min)

**Test Matrix**:

| Test | Expected Result | Status |
|------|-----------------|--------|
| 1. Start app → Capture Event → Quick → Press Esc | Back to strategy selection | ⬜ |
| 2. From strategy → Press Esc | Back to main menu | ⬜ |
| 3. In review → Press 'e' → Press Esc | Metadata modal closes | ⬜ |
| 4. In review → Press 'b' → Press Esc | Burst modal closes | ⬜ |
| 5. In review → Press 'f' → Press Esc | Fact modal closes | ⬜ |
| 6. From any state → Press 'm' | Return to main menu | ⬜ |
| 7. From any state → Press 'q' | Quit application | ⬜ |
| 8. Complete full capture workflow | All states navigate correctly | ⬜ |

**Command**:
```bash
go build -o kariya ./cmd/cli && ./kariya
```

---

## Pre-Implementation Checklist

Before starting the fix:

- [ ] Read bug report: `bugs/bug-001-escape-key-navigation.md`
- [ ] Read audit summary: `bugs/AUDIT-SUMMARY.md`
- [ ] Understand the 4 functions that need fixing
- [ ] Review correct pattern in other intents (e.g., BurstManagement)
- [ ] Ensure development environment ready

---

## Implementation Checklist

### Code Changes
- [ ] Fix `updateCaptureForm()` (Line 298)
- [ ] Fix `updateReviewInferredEvent()` - Add global key check at start
- [ ] Fix MetadataModal case (Line 383)
- [ ] Fix BurstModal case (Line 403) - Remove manual escape check
- [ ] Fix FactModal case (Line 418)
- [ ] Verify code compiles: `go build ./internal/cli/intents`

### Testing
- [ ] Add unit tests to `capture_event_escape_test.go`
- [ ] Run unit tests: `go test -v ./internal/cli/intents/capture_event_escape_test.go`
- [ ] Create E2E test file: `capture_event_escape_e2e_test.go`
- [ ] Run E2E tests: `ginkgo -r ./internal/testutil/e2e`
- [ ] Run full test suite: `go test ./...`
- [ ] Run with race detector: `go test -race ./...`

### Manual Testing
- [ ] Complete manual test matrix (8 tests above)
- [ ] Test form submission still works (Ctrl+S)
- [ ] Test all modal completion flows
- [ ] Test 'm' key from all states
- [ ] Test 'q' key from all states

### Documentation
- [ ] Create `docs/development/ESCAPE_KEY_HANDLING.md`
- [ ] Update `docs/TUI_STANDARDS.md` with warning
- [ ] Add code comments to fixed functions
- [ ] Update bug-001 with resolution summary

### Final Verification
- [ ] All tests passing
- [ ] No race conditions
- [ ] Code coverage maintained (>87%)
- [ ] Linting passing: `staticcheck ./...`
- [ ] Compliance check: `make check-compliance`

---

## Success Criteria

**Code Quality** ✅:
- All 4 functions follow correct pattern
- HandleGlobalKeys checked BEFORE delegation
- Consistent with other 9 intents
- No manual escape checks

**Testing** ✅:
- Unit tests cover all 4 functions
- E2E tests cover complete workflows
- All tests passing
- No regressions

**User Experience** ✅:
- Escape works from all screens
- User never gets trapped
- Consistent behavior across app
- Professional UX maintained

**Documentation** ✅:
- Pattern documented clearly
- Warning in TUI_STANDARDS.md
- Developer guide created
- Bug report updated

---

## Getting Help

If stuck or unsure:

1. **Review Examples**: Look at BurstManagement or FactManagement intents (correct pattern)
2. **Read Guides**: 
   - `bugs/DEBUGGING_GUIDE.md` - Systematic debugging
   - `docs/development/NAVIGATION_TESTING_GUIDE.md` - Navigation testing
3. **Check Tests**: Existing escape tests in other intents
4. **Bug Report**: Full investigation in `bug-001-escape-key-navigation.md`

---

## Timeline Estimate

| Phase | Duration | Cumulative |
|-------|----------|------------|
| Phase 1: Code Fixes | 2-3 hours | 2-3 hours |
| Phase 2: Unit Tests | 1 hour | 3-4 hours |
| Phase 3: E2E Tests | 2 hours | 5-6 hours |
| Phase 4: Documentation | 30 min | 5.5-6.5 hours |
| Phase 5: Manual Testing | 30 min | 6-7 hours |

**Total**: 6-7 hours

---

## After Completion

When all work is done:

1. **Update Bug Report**:
   - Change status to "Closed"
   - Fill in "Resolution Summary"
   - Add commit hashes
   - Document verification results

2. **Update README**:
   - Move bug from "Active" to "Closed" in `bugs/README.md`
   - Add one-line summary to closed bugs table

3. **Commit Changes**:
   ```bash
   git add internal/cli/intents/capture_event_intent.go
   git add internal/cli/intents/capture_event_escape_test.go
   git add internal/testutil/e2e/capture_event_escape_e2e_test.go
   git add docs/development/ESCAPE_KEY_HANDLING.md
   git add docs/TUI_STANDARDS.md
   git add bugs/bug-001-escape-key-navigation.md
   git add bugs/README.md
   
   git commit -m "fix(intents): handle escape key before form/modal delegation

Fixes escape key navigation in CaptureEvent intent by checking
HandleGlobalKeys before delegating to forms and modals.

Root cause: Message delegation happened before global key checks,
allowing Huh forms and modals to consume escape key events.

Solution: Reordered message handling to check global keys first
in all 4 affected functions.

Functions fixed:
- updateCaptureForm() (Line 298)
- updateReviewInferredEvent() - MetadataModal (Line 383)
- updateReviewInferredEvent() - BurstModal (Line 403)
- updateReviewInferredEvent() - FactModal (Line 418)

Testing:
- Added 4 unit tests for form/modal escape handling
- Added E2E tests for complete workflow navigation
- Manual testing: All 8 test scenarios passing

Impact: Users can now escape from all CaptureEvent screens
consistently using escape key, 'm' key, and 'q' key.

Fixes: bug-001

AI-Generated-By: OpenCode (Claude)
Reviewed-By: [Your Name]
"
   ```

4. **Create Follow-up Tasks** (if needed):
   - Document pattern in intent template
   - Create linter rule to catch this
   - Add to code review checklist

---

**Ready to Start**: All investigation complete, clear path forward

**Next Action**: Begin Phase 1 - Code Fixes

**Questions**: Refer to bug report and audit summary for details
