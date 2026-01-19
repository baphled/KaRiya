# Task 48: Fix Flaky Burst Modal Escape Test on Windows

## Overview
- **Goal**: Fix the intermittently failing test `should preserve original burst data when escape is pressed` that causes CI failures on Windows
- **Time Estimate**: 15-30 minutes
- **Prerequisites**: Understanding of the test environment helpers and burst list ordering
- **Priority**: HIGH - blocking PR #81 from merging

---

## Problem Description

### The Failing Test
- **File**: `internal/cli/intents/burst_management_modal_escape_test.go:91-109`
- **Test Name**: `BurstManagement Modal Escape Handling > Edit Modal Escape - Existing Burst Edit > should preserve original burst data when escape is pressed`
- **Platform**: Windows only (passes on Linux and macOS)

### Root Cause
The test at lines 91-109 uses a simple string check to find "Team Mentoring" in the view:

```go
view := env.GetView()
if !strings.Contains(view, "Team Mentoring") {
    env.NavigateDown() // Move to second burst
}
```

This is problematic because:
1. Burst ordering in the list is non-deterministic across platforms
2. The string "Team Mentoring" might appear in breadcrumbs or other UI elements
3. After `NavigateDown()`, if "Team Mentoring" wasn't the cursor target, the test confirms the wrong burst

### Correct Pattern
The test at lines 69-89 in the same file uses the proper cursor marker detection:

```go
view := env.GetView()
if !strings.Contains(view, "▶ Team Mentoring") && !strings.Contains(view, "> Team Mentoring") {
    env.NavigateDown() // Move to second burst
}
```

This checks for the **cursor marker** (`▶` or `>`) before the burst name, ensuring we're checking if that specific burst is currently selected.

---

## Session Contract Acknowledgment
- [ ] Ran `make session-start` and it passed
- [ ] Acknowledge and commit to following all workflow rules
- [ ] Token count: _____ (must be < 50k to start)

## Pre-Task Checklist (MUST COMPLETE BEFORE STARTING)
- [ ] `make check-compliance` passes
- [ ] Reviewed existing patterns in: `internal/cli/intents/burst_management_modal_escape_test.go`
- [ ] Confirmed this is ONE atomic task (not multiple changes)
- [ ] Identified which test files will be created/modified

## Files to Modify
- [ ] `internal/cli/intents/burst_management_modal_escape_test.go` (lines 94-98)

---

## TDD Checklist (MUST COMPLETE IN ORDER)

### Note on TDD for Test Fixes
This task fixes a flaky test, not production code. The "test" for this fix is running the actual test suite and confirming:
1. The test passes consistently on all platforms
2. The test still validates the intended behavior

### Implementation Phase
- [ ] Modify line 94-98 to use cursor marker detection pattern
- [ ] Run tests locally to verify fix: `go test -v ./internal/cli/intents/... -run "should preserve original burst data"`
- [ ] Run full burst management test suite: `go test -v ./internal/cli/intents/... -run "BurstManagement"`

---

## Implementation Details

### Current Code (lines 91-109)
```go
It("should preserve original burst data when escape is pressed", func() {
    env.SelectIntentByName("burst_management")

    // Navigate to find "Team Mentoring" burst (could be first or second depending on OS)
    view := env.GetView()
    if !strings.Contains(view, "Team Mentoring") {
        env.NavigateDown() // Move to second burst
    }

    env.Confirm()
    env.AssertViewContainsAny("Team Mentoring", "Focused mentoring")

    env.PressKeyRune('e')
    env.AssertViewContainsAny("Burst Name", "Description")

    env.Cancel()

    env.AssertViewContainsAny("Team Mentoring", "Focused mentoring")
})
```

### Fixed Code
```go
It("should preserve original burst data when escape is pressed", func() {
    env.SelectIntentByName("burst_management")

    // Navigate to find "Team Mentoring" burst (could be first or second depending on OS)
    view := env.GetView()
    if !strings.Contains(view, "▶ Team Mentoring") && !strings.Contains(view, "> Team Mentoring") {
        env.NavigateDown() // Move to second burst
    }

    env.Confirm()
    env.AssertViewContainsAny("Team Mentoring", "Focused mentoring")

    env.PressKeyRune('e')
    env.AssertViewContainsAny("Burst Name", "Description")

    env.Cancel()

    env.AssertViewContainsAny("Team Mentoring", "Focused mentoring")
})
```

---

## Pre-Commit Checklist (BEFORE EACH COMMIT)
- [ ] `make check-compliance` passes (REQUIRED before commit)
- [ ] Use `make ai-commit MSG="fix(tests): use cursor marker detection in burst modal escape test"` for AI-generated code
- [ ] Commit message explains **WHY**, not just WHAT
- [ ] Commit is atomic (ONE logical change)

## Post-Task Checklist (MUST COMPLETE BEFORE NEXT TASK)
- [ ] `make check-compliance` passes
- [ ] All checkboxes above completed
- [ ] Task marked complete `[x]` in task file
- [ ] Token count: _____ (< 100k to continue)

---

## Acceptance Criteria
- [ ] Test passes consistently on Windows (verified via CI)
- [ ] Test passes on Linux and macOS (no regressions)
- [ ] Test still validates the original intended behavior (escape preserves burst data)
- [ ] PR #81 CI passes after this fix is merged

## Verification Steps
1. Push fix to the feature branch
2. Verify CI passes on all platforms (Windows, Linux, macOS)
3. Confirm PR #81 checks turn green

## Rollback Plan
- Revert the single line change if any unexpected issues occur
- The change is minimal and low-risk

---

## Context

### Related PR
- **PR #81**: feat: CV Generation wizard modal
- **URL**: https://github.com/baphled/KaRiya/pull/81
- **Status**: Blocked by this flaky test failure

### Evidence of Pre-existing Bug
- Same test fails on `next` branch (CI run ID: `21045570272`)
- Test file was not modified in PR #81 (last changed in commits `89be814` and `87d6e1c`)
- Failure is Windows-specific due to platform-dependent ordering

### Test Data
When `PopulateTestData(5, 2, 0)` is called, it creates:
- Burst 0: "Authentication System Overhaul" (Confirmed: true)
- Burst 1: "Team Mentoring Initiative" (Confirmed: false)

The order these appear in the UI list varies by platform.
