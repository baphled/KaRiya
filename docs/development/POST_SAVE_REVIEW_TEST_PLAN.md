# Post-Save Review Testing Plan

**Date**: 2026-01-14
**Feature**: CaptureEvent Post-Save Enrichment Review
**Commit**: `c1b3285`
**Branch**: `feature/task-42-tui-architecture-refactor`

---

## Overview

This test plan verifies that the CaptureEvent intent now correctly shows enriched bursts and facts after saving an event, instead of immediately returning to the main menu.

## What Was Fixed

### Before (Buggy Behavior)
```
Form → Review → Submit → Save → Success Modal → COMPLETE (back to main menu) ❌
```

### After (Fixed Behavior)
```
Form → Review → Submit → Save → Enrichment → Success Modal → 
  Review (post-save, shows bursts/facts) → Enter to Complete ✅
```

### Key Changes
1. **Enrichment now runs for ALL strategies** (not just "enriched" mode)
2. **Post-save returns to Review state** (not completion)
3. **Review shows InferredBursts/InferredFacts** (not empty AcceptedBursts/AcceptedFacts)
4. **postSaveReview flag** prevents double-save bug

---

## Prerequisites

```bash
cd /home/baphled/Projects/KaRiya
./kariya  # Application is built and ready
```

**Expected Database**: `~/.kariya/events.db` (or use `--in-memory` for testing)

---

## Test Cases

### Test 1: Basic Post-Save Review Flow

**Objective**: Verify user sees review screen after save with enriched data

**Steps**:
1. Launch application: `./kariya`
2. Select **"Capture Event"** from main menu
3. Fill out form with test data:
   - **Title**: "Implemented new authentication system"
   - **Date**: "2026-01-14" (or use "today")
   - **Company**: "Test Company"
   - **Description**: "Built a secure OAuth2 authentication flow with JWT tokens and role-based access control. Integrated with existing user management system."
4. Press **Tab** to advance through fields
5. Press **Enter** when form is complete
6. Should see **Review** screen with event details
7. Press **Enter** or **Ctrl+S** to submit
8. Wait for enrichment to complete
9. Should see **Success Modal**: "Event saved successfully!"
10. Press **Enter** to dismiss modal

**Expected Results**:
- ✅ After dismissing success modal, you should return to **Review** screen (NOT main menu)
- ✅ Review screen should show:
  - Event details (title, date, company, description)
  - **"Inferred Bursts"** section with 1+ bursts (e.g., "Authentication System", "OAuth2 Integration")
  - **"Inferred Facts"** section with 1+ facts (e.g., skills, technologies)
- ✅ Footer should show navigation options (exact text may vary)

**Actual Results**:
- [ ] Pass / [ ] Fail
- Notes: _______________________________________________

---

### Test 2: Complete Intent from Post-Save Review

**Objective**: Verify pressing Enter in post-save review completes the intent (doesn't re-submit)

**Steps**:
1. Continue from Test 1 (you should be at post-save Review screen)
2. Press **Enter** or **Ctrl+S**

**Expected Results**:
- ✅ Intent should complete and return to main menu
- ✅ Should NOT see duplicate save operations
- ✅ Should NOT see "Event saved successfully!" modal again

**Actual Results**:
- [ ] Pass / [ ] Fail
- Notes: _______________________________________________

---

### Test 3: Verify Enrichment Actually Ran

**Objective**: Confirm burst and fact extraction occurred

**Steps**:
1. From main menu, select **"Browse Timeline"**
2. Find the event you just created ("Implemented new authentication system")
3. Select the event to view details
4. Check for associated bursts and facts

**Expected Results**:
- ✅ Event should have 1+ associated bursts
- ✅ Event should have 1+ associated facts
- ✅ Bursts should be contextually relevant (e.g., about authentication, OAuth2, security)
- ✅ Facts should extract specific technologies/skills mentioned

**Actual Results**:
- [ ] Pass / [ ] Fail
- Bursts found: _______________________________________________
- Facts found: _______________________________________________

---

### Test 4: Pre-Save Review Still Works

**Objective**: Verify pre-save review (before enrichment) hasn't broken

**Steps**:
1. From main menu, select **"Capture Event"**
2. Fill out form with minimal data:
   - **Title**: "Quick test event"
   - **Date**: "today"
   - **Company**: "Test Co"
3. Press **Enter** to proceed to Review
4. **DO NOT SUBMIT YET** - verify you can navigate back
5. Press **Esc** to go back to form
6. Verify form still has your data
7. Press **Enter** again to return to Review
8. Now press **Enter** to submit

**Expected Results**:
- ✅ Pre-save review shows event details (no bursts/facts yet)
- ✅ Esc goes back to form without losing data
- ✅ Pressing Enter in pre-save review submits the event
- ✅ After submit, enrichment runs and shows post-save review

**Actual Results**:
- [ ] Pass / [ ] Fail
- Notes: _______________________________________________

---

### Test 5: Error Handling

**Objective**: Verify graceful handling if enrichment fails

**Steps**:
1. From main menu, select **"Capture Event"**
2. Fill out form with very short description (to test weak enrichment):
   - **Title**: "Meeting"
   - **Date**: "today"
   - **Company**: "Test"
   - **Description**: "Had a meeting"
3. Submit and wait for enrichment

**Expected Results**:
- ✅ Event should save successfully
- ✅ Post-save review should appear even if no bursts/facts were inferred
- ✅ Should show "No inferred bursts" or empty list (not crash)
- ✅ Should show "No inferred facts" or empty list (not crash)
- ✅ Can still press Enter to complete intent

**Actual Results**:
- [ ] Pass / [ ] Fail
- Notes: _______________________________________________

---

### Test 6: Escape Key Behavior

**Objective**: Verify escape key works correctly in post-save review

**Steps**:
1. Complete a full capture workflow (Form → Review → Submit → Save)
2. When post-save review appears, press **Esc**

**Expected Results**:
- ✅ Should go back one state (likely to Submit or Review confirmation)
- ✅ Should NOT complete the intent immediately
- ✅ Should NOT cause a crash or unexpected behavior

**Actual Results**:
- [ ] Pass / [ ] Fail
- Notes: _______________________________________________

---

## Known Limitations

1. **Modal Editing Not Implemented**: Pressing `b`, `f`, or `m` keys in post-save review is not yet hooked up to edit bursts/facts. This is intentional - users can edit via BrowseTimeline later.

2. **Footer Text**: Footer may not explicitly say "post-save review" vs "pre-save review". This is a minor UX enhancement deferred for now.

3. **Enrichment Service Availability**: If CareerService is not initialized, enrichment will be skipped silently. This is expected behavior for test/development modes.

---

## Debugging Tips

### If Review Doesn't Show Bursts/Facts

**Check**:
1. Are you using in-memory mode? (`--in-memory`)
2. Is CareerService initialized? (check logs)
3. Did enrichment actually run? (add debug prints in `performEnrichment()` at line 809)

**Workaround**: Add debug logging:
```go
// In capture_event_intent.go, line 815 (inside performEnrichment)
log.Printf("DEBUG: Starting enrichment for event: %s", i.state.reviewState.Event.Title)
log.Printf("DEBUG: Found %d bursts, %d facts", len(bursts), len(facts))
```

### If Double-Save Bug Occurs

**Symptom**: Pressing Enter in post-save review saves the event again

**Check**:
1. Is `postSaveReview` flag being set? (line 259)
2. Is the flag being checked in Enter handler? (line 562)

**Debug**:
```go
// In updateReviewInferredEvent, line 560
log.Printf("DEBUG: Enter pressed, postSaveReview=%v", i.state.postSaveReview)
```

### If Sent Back to Main Menu

**Symptom**: After success modal, goes to main menu instead of review

**Check**:
1. Is `DismissModalMsg` handler correctly setting state? (line 254-260)
2. Is `currentState` being set to `CaptureStateReview`? (line 259)

**Debug**:
```go
// In Update method, line 256
log.Printf("DEBUG: Modal dismissed, returning to Review state, postSaveReview=%v", i.state.postSaveReview)
```

---

## Success Criteria

All 6 test cases must pass:
- [ ] Test 1: Basic post-save review flow
- [ ] Test 2: Complete intent from post-save review
- [ ] Test 3: Verify enrichment actually ran
- [ ] Test 4: Pre-save review still works
- [ ] Test 5: Error handling (weak/no enrichment)
- [ ] Test 6: Escape key behavior

---

## Regression Testing

After confirming these tests pass, run full test suite:

```bash
cd /home/baphled/Projects/KaRiya

# Run all tests (may take 2-3 minutes)
go test ./... 2>&1 | tee test_results.txt

# Check for failures
grep -E "FAIL|PASS" test_results.txt | tail -20

# Look for specific CaptureEvent test results
grep -A5 "capture_event" test_results.txt
```

**Expected**: All existing tests should still pass. We did not modify test files, only implementation.

---

## Next Steps After Testing

### If All Tests Pass ✅

1. Update `docs/workflows/EVENT_CAPTURE_WORKFLOW.md` to mark post-save enrichment as implemented
2. Consider adding footer text enhancement (optional, 15 min)
3. Consider adding modal editing integration (optional, 2-3 hours)
4. Resume Task 42 work on ConfigureSystem intent

### If Tests Fail ❌

1. Document specific failure scenarios
2. Add debug logging to identify root cause
3. Create minimal reproduction steps
4. Fix identified issues
5. Retest

---

## Test Results Summary

**Tester**: _______________________________________________
**Date/Time**: _______________________________________________
**Build**: `./kariya` (commit `c1b3285`)

| Test | Pass/Fail | Notes |
|------|-----------|-------|
| Test 1: Basic post-save review | [ ] | |
| Test 2: Complete from post-save | [ ] | |
| Test 3: Enrichment verification | [ ] | |
| Test 4: Pre-save review | [ ] | |
| Test 5: Error handling | [ ] | |
| Test 6: Escape behavior | [ ] | |

**Overall Status**: [ ] PASS / [ ] FAIL

**Issues Found**:
- 
- 
- 

**Recommendations**:
- 
- 
- 
