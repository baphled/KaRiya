---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Task 31: CaptureEvent Modal & Review Fixes

**Created**: 2026-01-08
**Status**: Ready for Implementation
**Priority**: HIGH
**Estimated Time**: 3-4 hours
**Related**: Codebase Audit (2026-01-08)

---

## Overview

CaptureEvent intent has broken review functionality: modals are never displayed when EditingMode is set, the review view reads from wrong state path, and accept/reject handlers are missing despite being advertised in help text.

**Issues**:
1. Review view checks `i.state.result.Event` (always nil) instead of `i.state.reviewState.Event`
2. EditingMode is set (e/b/f keys) but modals are never instantiated or displayed
3. Accept/reject handlers (a/r keys) are advertised but don't exist
4. Help text shows shortcuts that don't work

---

## Files to Modify

- [ ] `internal/cli/intents/capture_event_intent.go`
- [ ] `internal/cli/intents/capture_event_test.go`

---

## Implementation Plan

### Phase 1: Fix Review State Path (30 min)

**Location**: Lines 762-770, 815-826

**Change**:
```go
// BEFORE (WRONG)
if i.state.result != nil && i.state.result.Event != nil {
    title := i.state.result.Event.Text

// AFTER (CORRECT)
if i.state.reviewState.Event != nil {
    title := i.state.reviewState.Event.Text
```

**Tasks**:
- [ ] Replace `i.state.result.Event` with `i.state.reviewState.Event` in viewReviewInferredEvent()
- [ ] Replace `i.state.result.Event` with `i.state.reviewState.Event` in viewSubmit()
- [ ] Test review view displays event data correctly

---

### Phase 2: Implement Modal Display (2 hours)

**Location**: Lines 375-388, 756-805

**Current**: EditingMode is set but never checked
```go
case "e":
    i.state.reviewState.EditingMode = EditingModeMetadata
    return nil  // Does nothing!
```

**Fix**: Check EditingMode and display modal
```go
func (i *CaptureEventIntent) viewReviewInferredEvent() string {
    // Check if editing mode is active
    switch i.state.reviewState.EditingMode {
    case EditingModeMetadata:
        return i.viewMetadataEditModal()
    case EditingModeBursts:
        return i.viewBurstEditModal()
    case EditingModeFacts:
        return i.viewFactEditModal()
    }
    
    // Normal review view
    // ... existing code
}

func (i *CaptureEventIntent) viewMetadataEditModal() string {
    if i.state.reviewState.metadataModal == nil {
        i.state.reviewState.metadataModal = models.NewEditMetadataModal(
            i.state.reviewState.Event,
        )
    }
    return i.state.reviewState.metadataModal.View()
}

// Handle modal updates
func (i *CaptureEventIntent) updateReviewInferredEvent(msg tea.Msg) tea.Cmd {
    // If modal is active, pass updates to modal
    switch i.state.reviewState.EditingMode {
    case EditingModeMetadata:
        if i.state.reviewState.metadataModal != nil {
            modal, cmd := i.state.reviewState.metadataModal.Update(msg)
            i.state.reviewState.metadataModal = modal.(*models.EditMetadataModal)
            
            // Check if modal completed
            if modal.IsComplete() {
                if modal.WasAccepted() {
                    // Apply changes to event
                    i.state.reviewState.Event = modal.GetModifiedEvent()
                }
                // Clear modal and editing mode
                i.state.reviewState.metadataModal = nil
                i.state.reviewState.EditingMode = EditingModeNone
            }
            return cmd
        }
    // Similar for bursts and facts
    }
    
    // Normal review handling
    // ... existing code
}
```

**Add to ReviewState**:
```go
type ReviewState struct {
    Event         *career.CareerEvent
    InferredBursts []*career.Burst
    InferredFacts  []*career.Fact
    AcceptedBursts []*career.Burst
    AcceptedFacts  []*career.Fact
    EditingMode    EditingMode
    
    // ADD THESE
    metadataModal  *models.EditMetadataModal
    burstModal     *models.EditBurstModal
    factModal      *models.EditFactModal
}
```

**Tasks**:
- [ ] Add modal fields to ReviewState
- [ ] Check EditingMode in viewReviewInferredEvent()
- [ ] Display appropriate modal when mode is set
- [ ] Handle modal updates in updateReviewInferredEvent()
- [ ] Apply modal changes when accepted
- [ ] Clear modal and mode when complete/cancelled
- [ ] Test each modal (metadata, burst, fact)

---

### Phase 3: Implement Accept/Reject (1-1.5 hours)

**Location**: Lines 349-410

**Add handlers**:
```go
case "a":
    // Accept current item
    i.acceptCurrentItem()
    
case "r":
    // Reject current item
    i.rejectCurrentItem()

func (i *CaptureEventIntent) acceptCurrentItem() {
    // Determine what's selected (burst or fact)
    if i.state.reviewState.SelectedItemType == "burst" {
        idx := i.state.reviewState.SelectedIndex
        if idx < len(i.state.reviewState.InferredBursts) {
            burst := i.state.reviewState.InferredBursts[idx]
            i.state.reviewState.AcceptedBursts = append(
                i.state.reviewState.AcceptedBursts, burst)
            // Remove from inferred
            i.state.reviewState.InferredBursts = append(
                i.state.reviewState.InferredBursts[:idx],
                i.state.reviewState.InferredBursts[idx+1:]...)
        }
    } else if i.state.reviewState.SelectedItemType == "fact" {
        // Similar for facts
    }
    
    // Move to next item
    if i.state.reviewState.SelectedIndex >= len(i.state.reviewState.InferredBursts) {
        i.state.reviewState.SelectedIndex = 0
    }
}

func (i *CaptureEventIntent) rejectCurrentItem() {
    // Just remove from inferred, don't add to accepted
    if i.state.reviewState.SelectedItemType == "burst" {
        idx := i.state.reviewState.SelectedIndex
        if idx < len(i.state.reviewState.InferredBursts) {
            i.state.reviewState.InferredBursts = append(
                i.state.reviewState.InferredBursts[:idx],
                i.state.reviewState.InferredBursts[idx+1:]...)
        }
    }
    // Move to next
}
```

**Add to ReviewState**:
```go
type ReviewState struct {
    // ... existing fields
    SelectedItemType string  // "burst" or "fact"
    SelectedIndex    int     // Which item is selected
}
```

**Tasks**:
- [ ] Add SelectedItemType and SelectedIndex to ReviewState
- [ ] Implement acceptCurrentItem()
- [ ] Implement rejectCurrentItem()
- [ ] Add 'a' key handler
- [ ] Add 'r' key handler
- [ ] Add item navigation (j/k to move between items)
- [ ] Update view to show current selection
- [ ] Test accept/reject workflow

---

### Phase 4: Fix Help Text (15 min)

**Location**: Line 662

**Update footer**:
```go
case CaptureStateReview:
    if i.state.reviewState.EditingMode != EditingModeNone {
        return "Editing... | Esc: Cancel | Enter: Save"
    }
    return CombineFooters(NavigationFooter(), 
        "e Edit metadata | b Edit bursts | f Edit facts | a Accept | r Reject", base)
```

**Tasks**:
- [ ] Update help text to match actual handlers
- [ ] Show different help when modal is active
- [ ] Test help text displays correctly

---

## Acceptance Criteria

### Must Have
- [ ] Review state displays event data from correct state path
- [ ] Pressing 'e' shows metadata edit modal
- [ ] Pressing 'b' shows burst edit modal
- [ ] Pressing 'f' shows fact edit modal
- [ ] Modal changes are applied when accepted
- [ ] Modal changes are discarded when cancelled
- [ ] Pressing 'a' accepts current item (burst/fact)
- [ ] Pressing 'r' rejects current item
- [ ] Help text matches actual functionality
- [ ] All tests passing (maintain 2,078/2,078)

### Should Have
- [ ] Visual indicator shows which item is selected
- [ ] Can navigate between items with j/k
- [ ] Accepted items are visually distinct
- [ ] Rejected items are removed from view

---

## Testing Strategy

```go
It("should display event data from reviewState", func() {
    intent.state.reviewState.Event = testEvent
    view := intent.View()
    Expect(view).To(ContainSubstring(testEvent.Text))
})

It("should show metadata modal when e pressed", func() {
    intent.Update(tea.KeyMsg{String: "e"})
    view := intent.View()
    Expect(view).To(ContainSubstring("Edit Metadata"))
})

It("should accept burst when a pressed", func() {
    intent.state.reviewState.InferredBursts = []*career.Burst{testBurst}
    intent.Update(tea.KeyMsg{String: "a"})
    Expect(intent.state.reviewState.AcceptedBursts).To(HaveLen(1))
    Expect(intent.state.reviewState.InferredBursts).To(HaveLen(0))
})
```

---

## References

- `internal/cli/models/modals.go` - Existing modal implementations
- `internal/cli/intents/capture_event_intent.go` - Implementation file
- `docs/TUI_STANDARDS.md` - Modal patterns

---

**Last Updated**: 2026-01-08
**Status**: Ready for implementation
