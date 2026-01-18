# Task 45: Burst Auto-Detection UI

## Overview
- **Goal**: Expose `BurstDetector` in Burst Management TUI for automatic burst suggestions
- **Time Estimate**: 6-8 hours
- **Prerequisites**: Understanding of burst detection algorithm, intent architecture

## Session Contract Acknowledgment
- [ ] Ran `make session-start` and it passed
- [ ] Acknowledge and commit to following all workflow rules
- [ ] Token count: _____ (must be < 50k to start)

## Pre-Task Checklist (MUST COMPLETE BEFORE STARTING)
- [ ] `make check-compliance` passes
- [ ] Reviewed existing patterns in:
  - `internal/service/career/burst_fact/detector.go`
  - `internal/cli/intents/burst_management_intent.go`
  - `internal/cli/intents/burst_management.go`
- [ ] Confirmed this is ONE atomic task
- [ ] Identified which test files will be created/modified

## Context

The `BurstDetector` is fully implemented with:
- Temporal grouping (6-month window)
- Multi-factor similarity scoring (text, keywords, company, project)
- Cluster detection via BFS
- Confidence scoring

However, it's NOT exposed in the TUI. Users can only create bursts manually.

## Current UX Flow
```
Burst List → [n] New Burst → Manual Creation → Done
```

## Target UX Flow
```
Burst List → [s] Suggest Bursts → Detecting... →
Review Suggestions (accept/reject each) →
Accepted become new bursts
```

## Files to Modify

### Intent Layer
- [ ] `internal/cli/intents/burst_management.go` - Add new states and messages
- [ ] `internal/cli/intents/burst_management_intent.go` - Add detection UI and handlers
- [ ] `internal/cli/intents/burst_management_context.go` - Add DetectBursts() method
- [ ] `internal/cli/intents/burst_management_test.go` - Tests for new functionality

## New States

| State | Constant | Description |
|-------|----------|-------------|
| Suggesting | `BurstStateSuggesting` | Loading state during detection |
| Suggestion Review | `BurstStateSuggestionReview` | Review individual suggestions |

## New Messages

```go
// BurstSuggestionsLoadedMsg sent when detection completes
type BurstSuggestionsLoadedMsg struct {
    Suggestions []burst_fact.BurstSuggestion
    Error       error
}

// BurstSuggestionAcceptedMsg sent when user accepts a suggestion
type BurstSuggestionAcceptedMsg struct {
    Suggestion burst_fact.BurstSuggestion
    Burst      *domain.Burst
    Error      error
}
```

## New Key Bindings

### List View
| Key | Action |
|-----|--------|
| `s` | Suggest bursts (trigger detection) |

### Suggestion Review View
| Key | Action |
|-----|--------|
| `a` / `Enter` | Accept suggestion (create burst) |
| `r` | Reject suggestion (skip to next) |
| `A` | Accept all remaining suggestions |
| `Esc` | Cancel review, return to list |

## Implementation Checklist

### Phase 1: Add State Constants and Messages

#### RED Phase
- [ ] Test: New states exist in burst_management.go
- [ ] Test FAILS

#### GREEN Phase
- [ ] Add `BurstStateSuggesting = "suggesting"` constant
- [ ] Add `BurstStateSuggestionReview = "suggestion_review"` constant
- [ ] Add message types to burst_management.go
- [ ] Test PASSES

### Phase 2: Add Context Method for Detection

#### RED Phase
- [ ] Test: `BurstManagementContext.DetectBursts()` returns suggestions
- [ ] Test FAILS

#### GREEN Phase
- [ ] Add to `BurstManagementContext`:
  ```go
  func (ctx *BurstManagementContext) DetectBursts() ([]burst_fact.BurstSuggestion, error) {
      // Load all events
      events, err := ctx.Service.ListEvents(ctx.Context, nil)
      if err != nil {
          return nil, fmt.Errorf("failed to load events: %w", err)
      }

      // Convert to value slice for detector
      eventValues := make([]career.CareerEvent, len(events))
      for i, e := range events {
          eventValues[i] = *e
      }

      // Run detection
      detector := burst_fact.NewBurstDetector()
      return detector.DetectBursts(ctx.Context, eventValues, nil)
  }
  ```
- [ ] Test PASSES

### Phase 3: Add Model State for Suggestions

#### RED Phase
- [ ] Test: Model tracks suggestions and current suggestion index
- [ ] Test FAILS

#### GREEN Phase
- [ ] Add to `BurstManagementIntentModel`:
  ```go
  // Suggestion state
  suggestions           []burst_fact.BurstSuggestion
  currentSuggestionIdx  int
  suggestionsLoading    bool
  suggestionsError      error
  ```
- [ ] Test PASSES

### Phase 4: Add "s" Key Handler in List View

#### RED Phase
- [ ] Test: Pressing "s" in list view triggers detection
- [ ] Test FAILS

#### GREEN Phase
- [ ] In `updateListView()`, add:
  ```go
  case "s":
      i.state.currentState = BurstStateSuggesting
      i.state.suggestionsLoading = true
      return i.detectBursts()
  ```
- [ ] Add `detectBursts()` method returning tea.Cmd
- [ ] Test PASSES

### Phase 5: Add Suggestion Loading View

#### RED Phase
- [ ] Test: Suggesting state shows loading indicator
- [ ] Test FAILS

#### GREEN Phase
- [ ] Add `viewSuggesting()` method with loading spinner
- [ ] Handle `BurstSuggestionsLoadedMsg` to transition to review
- [ ] Test PASSES

### Phase 6: Add Suggestion Review View

#### RED Phase
- [ ] Test: Review view shows current suggestion details
- [ ] Test FAILS

#### GREEN Phase
- [ ] Add `viewSuggestionReview()` method showing:
  - Progress indicator (e.g., "Suggestion 1 of 5")
  - Confidence score with visual indicator
  - Event count
  - Suggested name and description
  - Preview of first 3 event texts
- [ ] Test PASSES

### Phase 7: Add Accept/Reject Handlers

#### RED Phase
- [ ] Test: Accepting creates burst and moves to next
- [ ] Test: Rejecting skips to next without creating
- [ ] Test FAILS

#### GREEN Phase
- [ ] In `updateSuggestionReview()`:
  ```go
  case "a", "enter":
      return i.acceptCurrentSuggestion()
  case "r":
      return i.rejectCurrentSuggestion()
  case "A":
      return i.acceptAllRemaining()
  ```
- [ ] Implement helper methods
- [ ] Test PASSES

### Phase 8: Handle Empty Suggestions

#### RED Phase
- [ ] Test: No suggestions shows appropriate message
- [ ] Test FAILS

#### GREEN Phase
- [ ] If `len(suggestions) == 0`, show:
  ```
  No burst suggestions found.

  Tips for better detection:
  • Add more events with similar projects or companies
  • Use consistent tagging across related events
  • Events within 6 months are grouped together
  ```
- [ ] Test PASSES

### Phase 9: Footer and Help Text

- [ ] Update `getContextHelp()` for new states
- [ ] Add appropriate key badges for suggestion review

## Pre-Commit Checklist (BEFORE EACH COMMIT)
- [ ] `make check-compliance` passes
- [ ] Use `make ai-commit MSG="type(scope): description"`
- [ ] Commit is atomic (ONE logical change)

## Post-Task Checklist
- [ ] `make check-compliance` passes
- [ ] All tests pass
- [ ] Can trigger detection from list view
- [ ] Can review and accept/reject suggestions
- [ ] Accepted suggestions create valid bursts

## Acceptance Criteria
- [ ] "s" key in list view triggers burst detection
- [ ] Detection shows loading state
- [ ] Suggestions display with confidence scores
- [ ] Accept creates burst, reject skips
- [ ] "Accept all" works for remaining suggestions
- [ ] Empty state handled gracefully
- [ ] All navigation (Esc, etc.) works correctly

## Rollback Plan
- Remove new states and handlers
- Remove context detection method
- Restore original list view behavior
