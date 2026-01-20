# CaptureEvent Post-Save Review Issue

**Discovered**: 2026-01-14  
**Severity**: High - Missing core workflow step  
**Status**: Documented, needs implementation  
**Related Task**: Should be Task 43 (Post-ConfigureSystem)

---

## Problem

After successfully saving an event, CaptureEvent immediately completes and returns to main menu. This skips the **post-save enrichment and review step** that is documented in the workflow guide.

---

## Expected Workflow

From `docs/workflows/EVENT_CAPTURE_WORKFLOW.md`:

1. **Form** → **Review** (preview event before save)
2. **Submit** → Save to DB ✅ **WORKS**
3. **Enrichment** → Extract bursts/facts from saved event ❌ **MISSING**
4. **Post-Save Review** → User reviews and accepts/rejects inferred bursts/facts ❌ **MISSING**
5. **Complete** → Navigate to BrowseTimeline or Main Menu

### Key Quote from Documentation

> "Enrichment (burst/fact extraction) happens **AFTER** event save, not before"  
> — `docs/workflows/EVENT_CAPTURE_WORKFLOW.md` line 442

---

## Current Workflow (Buggy)

1. **Form** → **Review** (preview event)
2. **Submit** → Save to DB ✅ **WORKS**
3. **Success modal** (3s auto-dismiss) ✅ **WORKS**
4. **Immediately complete** and return to main menu ❌ **WRONG**

**User Impact**: Users cannot review or edit inferred bursts and facts immediately after creating an event. They must manually navigate to BrowseTimeline and edit the event later.

---

## What's Missing

### 1. Post-Save Enrichment Trigger

After `eventService.CaptureEvent()` succeeds, we need to:
- Call `careerService.ExtractBursts(eventID)` to infer bursts
- Call `careerService.ExtractFacts(eventID)` to extract facts
- Handle async enrichment with loading state

### 2. Enrichment Message Handling

Create new message types:
```go
type EnrichmentStartMsg struct {}
type EnrichmentCompleteMsg struct {
    Bursts []*career.Burst
    Facts  []*career.Fact
}
type EnrichmentErrorMsg struct {
    Code    string
    Message string
    Cause   error
}
```

### 3. Return to Review State

After enrichment completes:
- Populate `reviewState.InferredBursts` with results
- Populate `reviewState.InferredFacts` with results
- Transition back to `CaptureStateReview`
- Update EventReviewScreen to display enriched data

### 4. Modal Integration

The 3 modals already exist but aren't integrated:
- **BurstEditModal** (`internal/cli/models/burst_modal.go`)
- **FactEditModal** (`internal/cli/models/fact_modal.go`)
- **MetadataEditModal** (`internal/cli/models/metadata_modal.go`)

These need to be wired up in the Review state:
- Press `b` → Show BurstEditModal
- Press `f` → Show FactEditModal
- Press `m` → Show MetadataEditModal

### 5. Final Completion

Only complete the intent after user:
- Reviews enriched data
- Edits/accepts/rejects bursts and facts
- Presses Enter or similar to confirm final review

---

## Technical Implementation Plan

### Step 1: Add Enrichment Trigger (2 hours)

**File**: `internal/cli/intents/capture_event_intent.go`

**After SubmitSuccessMsg** (around line 240-253):

```go
case components.ModalDismissedMsg:
    // Success modal auto-dismissed
    if i.state.submitModal != nil {
        i.state.submitModal = nil
        
        // Instead of completing immediately, trigger enrichment
        return i.performEnrichment()
    }
    return nil
```

**New method**:
```go
func (i *CaptureEventIntent) performEnrichment() tea.Cmd {
    return func() tea.Msg {
        if i.state.reviewState.Event == nil || i.state.reviewState.Event.ID == "" {
            return EnrichmentErrorMsg{
                Code:    "MISSING_EVENT_ID",
                Message: "Cannot enrich event without ID",
            }
        }
        
        eventID := i.state.reviewState.Event.ID
        
        // Extract bursts
        bursts, err := i.careerService.ExtractBursts(context.Background(), eventID)
        if err != nil {
            // Log error but continue (enrichment is optional)
            i.logger.Error("Failed to extract bursts: %v", err)
        }
        
        // Extract facts
        facts, err := i.careerService.ExtractFacts(context.Background(), eventID)
        if err != nil {
            // Log error but continue (enrichment is optional)
            i.logger.Error("Failed to extract facts: %v", err)
        }
        
        return EnrichmentCompleteMsg{
            Bursts: bursts,
            Facts:  facts,
        }
    }
}
```

### Step 2: Handle Enrichment Result (1 hour)

**In Update method**:
```go
case EnrichmentCompleteMsg:
    // Populate review state with enriched data
    i.state.reviewState.InferredBursts = msg.Bursts
    i.state.reviewState.InferredFacts = msg.Facts
    
    // Return to review state for user to review enriched data
    i.state.currentState = CaptureStateReview
    
    // Dismiss loading modal
    i.state.submitModal = nil
    return nil

case EnrichmentErrorMsg:
    // Show error but still return to review (enrichment is optional)
    i.state.error = &IntentError{
        Code:    msg.Code,
        Message: msg.Message,
        Cause:   msg.Cause,
    }
    i.state.currentState = CaptureStateReview
    i.state.submitModal = nil
    return nil
```

### Step 3: Integrate Modals in Review State (2 hours)

**In updateReview method**:
```go
func (i *CaptureEventIntent) updateReview(msg tea.Msg) tea.Cmd {
    // ... existing code ...
    
    switch msg.String() {
    case "b":
        // Show burst edit modal
        i.state.burstModal = models.NewBurstEditModal(
            i.state.reviewState.Event,
            i.state.reviewState.InferredBursts,
        )
        return i.state.burstModal.Init()
        
    case "f":
        // Show fact edit modal
        i.state.factModal = models.NewFactEditModal(
            i.state.reviewState.Event,
            i.state.reviewState.InferredFacts,
        )
        return i.state.factModal.Init()
        
    case "m":
        // Show metadata edit modal
        i.state.metadataModal = models.NewMetadataEditModal(
            i.state.reviewState.Event,
        )
        return i.state.metadataModal.Init()
        
    case "enter", "ctrl+s":
        // User finished reviewing - complete intent
        result := &CaptureEventResult{
            Event:          i.state.reviewState.Event,
            Bursts:         i.state.reviewState.AcceptedBursts,
            Facts:          i.state.reviewState.AcceptedFacts,
            AcceptedFields: make(map[string]bool),
            RejectedFields: i.state.reviewState.RejectedItems,
        }
        i.setCompleted(result)
        return nil
    }
    
    // ... handle modal updates ...
}
```

### Step 4: Update EventReviewScreen (1 hour)

**File**: `internal/cli/screens/capture/event_review_screen.go`

Update view to show enriched data:
```go
func (s *EventReviewScreen) View() string {
    // ... existing event display ...
    
    // Show inferred bursts
    if len(s.reviewState.InferredBursts) > 0 {
        content += "\n\n" + s.renderBursts()
    }
    
    // Show inferred facts
    if len(s.reviewState.InferredFacts) > 0 {
        content += "\n\n" + s.renderFacts()
    }
    
    // ... rest of view ...
}
```

### Step 5: Testing (2 hours)

**Test scenarios**:
1. Event capture → Enrichment succeeds → Review bursts/facts → Accept → Complete
2. Event capture → Enrichment fails → Review without enrichment → Complete
3. Event capture → Review → Edit bursts → Save → Review → Complete
4. Event capture → Review → Edit facts → Save → Review → Complete
5. Event capture → Review → Edit metadata → Save → Review → Complete

**Test files**:
- Update `internal/cli/intents/capture_event_test.go` with enrichment tests
- Add integration tests for modal workflows

---

## Related Files

### Intent Files
- `internal/cli/intents/capture_event_intent.go` - Main intent logic (needs enrichment trigger)
- `internal/cli/intents/capture_event.go` - Data structures (ReviewState has burst/fact fields)

### Modal Files (Already Exist)
- `internal/cli/models/burst_modal.go` - Edit bursts modal
- `internal/cli/models/fact_modal.go` - Edit facts modal
- `internal/cli/models/metadata_modal.go` - Edit metadata modal

### Screen Files
- `internal/cli/screens/capture/event_review_screen.go` - Review screen (needs enrichment display)

### Service Files
- `internal/service/career/service.go` - CareerService (has ExtractBursts/ExtractFacts methods)

---

## Documentation References

1. **Workflow Guide**: `docs/workflows/EVENT_CAPTURE_WORKFLOW.md`
   - Line 442: "Enrichment (burst/fact extraction) happens AFTER event save, not before"
   - Lines 432-436: Full path shows Edit Bursts/Facts in Review step

2. **Task File**: `tasks/tasks-42-tui-architecture-refactor.md`
   - Lines 1567-1577: Notes modals exist but not integrated
   - Lines 1569: "Review screen exists but currently bypassed (form → submit)"

3. **Burst & Fact Guide**: `docs/BURST_FACT_EXTRACTION_GUIDE.md`
   - Explains burst and fact concepts
   - Shows enrichment workflow

---

## Estimated Effort

**Total**: 8-10 hours

| Phase | Time | Description |
|-------|------|-------------|
| Enrichment trigger | 2h | Add performEnrichment() and async handling |
| Message handling | 1h | Handle EnrichmentCompleteMsg/ErrorMsg |
| Modal integration | 2h | Wire up 3 existing modals in Review state |
| Screen updates | 1h | Update EventReviewScreen to show enriched data |
| Testing | 2-4h | Unit tests + integration tests + manual testing |

---

## Recommendation

**Create separate task**: Task 43 - CaptureEvent Post-Save Enrichment & Review

**Priority**: High (core workflow feature)

**Dependencies**: None (can be done independently)

**When**: After ConfigureSystem migration is complete (Task 42)

**Why separate task**:
- Significant async workflow changes
- Requires modal integration (3 modals)
- Needs comprehensive testing
- Could take a full day (8-10 hours)
- ConfigureSystem work is already 10+ hours

---

## Workaround (Current Behavior)

Users can manually review/enrich events later:

1. Complete event capture (goes to main menu)
2. Navigate to **BrowseTimeline**
3. Find the newly created event
4. Press **Enter** to view details
5. Use event detail modals to edit metadata/bursts/facts

**This works but is less convenient** than the documented workflow where enrichment happens immediately after save.

---

## Questions for Product Owner

1. **Is enrichment mandatory** or can users skip it?
   - Current: Enrichment is optional (user presses b/f keys)
   - Proposed: Automatic after save, user reviews results

2. **Should enrichment block completion**?
   - Option A: Show loading modal, user waits
   - Option B: Run in background, allow immediate completion

3. **What if enrichment fails**?
   - Option A: Show error, allow retry
   - Option B: Log error, continue without enrichment

4. **Should enrichment happen for Quick mode**?
   - Quick mode skips review currently
   - Should it also skip enrichment?

---

**Created**: 2026-01-14  
**Author**: AI (Claude Sonnet 4.5)  
**Reviewed By**: Pending
