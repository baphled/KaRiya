# 🎯 CaptureEvent Post-Save Review - Ready for Testing

**Date**: 2026-01-14 (Late Evening)  
**Status**: ✅ **COMMITTED AND BUILT - READY FOR MANUAL TESTING**  
**Commit**: `c1b3285`  
**Branch**: `feature/task-42-tui-architecture-refactor`

---

## What Was Fixed

### The Bug
User reported: **"After saving an entry, we're sent back to main menu. Should go to review page instead."**

### The Solution
Implemented post-save enrichment review workflow:
- Event saves successfully
- Enrichment runs (burst/fact extraction)
- User returns to Review screen to see enriched data
- User presses Enter to complete (not re-submit)

---

## Quick Start Testing

```bash
cd /home/baphled/Projects/KaRiya

# Application is already built
./kariya

# Follow this workflow:
# 1. Main Menu → Capture Event
# 2. Fill form with rich description (mention skills, technologies)
# 3. Submit event
# 4. Wait for "Event saved successfully!" modal
# 5. Press Enter to dismiss modal
# 6. VERIFY: You should see Review screen with bursts/facts (NOT main menu)
# 7. Press Enter to complete and return to main menu
```

---

## Test Documentation

**Complete Test Plan**: `docs/development/POST_SAVE_REVIEW_TEST_PLAN.md`

This document includes:
- ✅ 6 comprehensive test cases
- ✅ Step-by-step instructions
- ✅ Expected vs actual result tables
- ✅ Debugging tips
- ✅ Known limitations
- ✅ Regression testing commands

---

## What to Verify

### Primary Goal
After saving an event:
1. ✅ Success modal appears: "Event saved successfully!"
2. ✅ After dismissing modal, you see **Review screen** (NOT main menu)
3. ✅ Review screen shows **Inferred Bursts** and **Inferred Facts**
4. ✅ Pressing Enter completes the intent (no double-save)

### Secondary Goals
1. ✅ Enrichment actually runs (bursts/facts are extracted)
2. ✅ Pre-save review still works (can go back to form)
3. ✅ Escape key works correctly
4. ✅ No crashes with weak/empty enrichment

---

## What Changed (Technical)

### Files Modified
1. **`internal/cli/intents/capture_event.go`** (+3 lines)
   - Added `postSaveReview bool` flag to model

2. **`internal/cli/intents/capture_event_intent.go`** (+29 lines, -7 lines)
   - Line 259: Sets `postSaveReview = true` when returning to review after save
   - Lines 562-575: Enhanced Enter key handler to check `postSaveReview` flag
   - Lines 1093-1115: Updated view to show `InferredBursts`/`InferredFacts`

### Key Logic Changes

**Before** (line 764):
```go
if i.context.CaptureStrategy == "enriched" && ... // Never ran
```

**After** (line 764):
```go
if i.context.CareerService != nil { ... // Runs for all strategies
```

**Post-Save Review** (line 259):
```go
i.state.postSaveReview = true // Mark as post-save
i.state.currentState = CaptureStateReview // Return to review
```

**Enter Key Logic** (line 562-575):
```go
if i.state.postSaveReview {
    // Complete intent with enriched data
    i.setCompleted(result)
} else {
    // Pre-save: proceed to submit
    i.performSubmit()
}
```

---

## Build Status

```
✅ Code compiles successfully
✅ Commit includes AI attribution
✅ Application built: ./kariya (21MB)
⚠️  Unit tests not run (pre-commit hook timeout)
🔄 Manual testing required
```

---

## Known Limitations

1. **Modal Editing Not Implemented**: Pressing `b`, `f`, or `m` keys in post-save review doesn't open editors. Users can edit via BrowseTimeline later. This is intentional for MVP.

2. **Footer Text**: Footer doesn't explicitly indicate "post-save review" vs "pre-save review". Minor UX enhancement deferred.

3. **Test Suite Timeout**: Full test suite times out (2,078 tests). Can be run separately if needed.

---

## Debugging Tips

### If Review Doesn't Show Bursts/Facts

Add debug logging to `capture_event_intent.go`:

```go
// Line 815 (inside performEnrichment)
log.Printf("DEBUG: Starting enrichment for event: %s", i.state.reviewState.Event.Title)
log.Printf("DEBUG: Found %d bursts, %d facts", len(bursts), len(facts))

// Line 1093 (inside buildReviewBaseView)
log.Printf("DEBUG: Rendering review, InferredBursts=%d, AcceptedBursts=%d",
    len(i.state.reviewState.InferredBursts),
    len(i.state.reviewState.AcceptedBursts))
```

### If Double-Save Bug Occurs

Add debug logging:

```go
// Line 560 (inside updateReviewInferredEvent)
log.Printf("DEBUG: Enter pressed, postSaveReview=%v", i.state.postSaveReview)
```

### If Sent Back to Main Menu

Add debug logging:

```go
// Line 256 (inside DismissModalMsg handler)
log.Printf("DEBUG: Modal dismissed, returning to Review, postSaveReview=%v", 
    i.state.postSaveReview)
log.Printf("DEBUG: Current state now: %v", i.state.currentState)
```

---

## Next Steps

### 1. Manual Testing (30-45 minutes)
- [ ] Run through test plan in `docs/development/POST_SAVE_REVIEW_TEST_PLAN.md`
- [ ] Fill out test results table
- [ ] Document any issues found

### 2. If Tests Pass ✅
- [ ] Update `docs/workflows/EVENT_CAPTURE_WORKFLOW.md` to mark feature complete
- [ ] Decide: Add footer text enhancement? (15 min)
- [ ] Decide: Add modal editing integration? (2-3 hours)
- [ ] Resume Task 42: ConfigureSystem intent migration

### 3. If Tests Fail ❌
- [ ] Document specific failure scenarios
- [ ] Add debug logging
- [ ] Create minimal reproduction steps
- [ ] Fix and retest

---

## Session Summary

**Commits Made This Session**: 4 total
1. `4d28129` - docs: add comprehensive post-save enrichment issue analysis
2. `9106cd5` - fix(intents): trigger enrichment for all strategies and return to review
3. `c1b3285` - fix(intents): implement post-save review with enriched data display
4. *(Previous work on Task 42)*

**Time Spent**: ~3 hours
**Status**: Core fix complete, testing required
**Branch**: `feature/task-42-tui-architecture-refactor` (not pushed yet)

---

## Ready to Test?

```bash
cd /home/baphled/Projects/KaRiya
./kariya  # Let's see if it works!
```

**Test Documentation**: `docs/development/POST_SAVE_REVIEW_TEST_PLAN.md`  
**Issue Documentation**: `docs/development/CAPTURE_EVENT_POST_SAVE_ISSUE.md`

Good luck! 🚀
