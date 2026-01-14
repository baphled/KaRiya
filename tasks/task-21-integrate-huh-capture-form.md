---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Task 21: Integrate Huh-Based Capture Form into CaptureEventIntent

## Overview
- **Goal**: Replace textinput-based FormModel with huh-based form for consistent UX
- **Time Estimate**: 2-3 hours
- **Prerequisites**: Huh form already exists at `internal/cli/forms/capture_event_form.go`

## Session Contract Acknowledgment
- [x] Ran `make session-start` and it passed (previously verified)
- [x] Acknowledge and commit to following all workflow rules
- [x] Token count: 112994 (< 50k to start) ⚠️ HIGH

## Problem Statement
The huh-based capture event form exists but is not integrated into CaptureEventIntent.
The CaptureEventIntent still uses the old textinput-based FormModel.

User reports:
- "Missing the capture event description (text)"
- "New capture event does not allow users to enter the description"

## Root Cause
CaptureEventIntent uses `*models.FormModel` (old textinput-based) instead of the new
huh-based form. The form at `internal/cli/forms/capture_event_form.go` exists but
is never instantiated or used.

## TDD Checklist (MUST COMPLETE IN ORDER)

### RED Phase - Write Failing Tests
- [ ] Test file: `internal/cli/intents/capture_event_test.go`
- [ ] Test: "CaptureEventIntent should use huh form"
  - Assert form is created with strategy
  - Assert form handles text input
  - Assert form handles date input
  - Assert form handles Ctrl+S submission
- [ ] Run tests and confirm they FAIL
- [ ] Commit failing tests:
  ```
  git commit -m "test(intents): add tests for huh form integration in CaptureEventIntent"
  ```

### GREEN Phase - Minimal Implementation
- [ ] Create wrapper model: `internal/cli/models/capture_form_huh.go`
  - Wraps `*huh.Form` from `forms.NewCaptureEventForm`
  - Implements BubbleTea model interface (Init, Update, View)
  - Handles form submission via SubmitMsg
  - Handles WindowSizeMsg for responsive forms
- [ ] Update CaptureEventModel: Change `captureForm` from `*models.FormModel` to interface or new type
- [ ] Update CaptureEventIntent: Use new form in `NewCaptureEventIntent`
- [ ] Update form initialization: Call `forms.NewCaptureEventForm` with strategy
- [ ] Run tests and confirm they PASS
- [ ] Commit implementation:
  ```
  git commit -m "feat(intents): integrate huh-based capture form"
  ```

### REFACTOR Phase
- [ ] Extract form interface if needed for testability
- [ ] Cleanup any duplicate code
- [ ] Update documentation/comments
- [ ] Run tests and confirm still passing
- [ ] Commit refactoring (if any):
  ```
  git commit -m "refactor(intents): cleanup capture form integration"
  ```

## Acceptance Criteria
- [ ] CaptureEventIntent uses huh-based form (not FormModel)
- [ ] Text field is visible and editable
- [ ] Date field is visible and editable  
- [ ] Company/Project fields visible in manual mode
- [ ] Tags/Categories MultiSelect visible in manual mode
- [ ] Ctrl+S submits form
- [ ] Form is properly styled with Catppuccin theme
- [ ] Form is scrollable when content exceeds height
- [ ] All tests pass

## Rollback Plan
If integration fails:
1. `git revert HEAD` to undo integration commit
2. Investigate test failures
3. Fix and retry

## Post-Task Checklist
- [ ] `make check-compliance` passes
- [ ] All tests pass (including E2E)
- [ ] Manual testing confirms form works
- [ ] Task marked complete [x]
