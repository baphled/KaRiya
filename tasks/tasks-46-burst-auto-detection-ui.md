# Task 46: Burst Auto-Detection UI

## Overview
- **Goal**: Expose `BurstDetector` in Burst Management TUI for automatic burst suggestions
- **Time Estimate**: 6-8 hours
- **Prerequisites**: Understanding of burst detection algorithm, intent architecture

## Session Contract Acknowledgment
- [x] Ran `make session-start` and it passed
- [x] Acknowledge and commit to following all workflow rules
- [x] Token count: N/A (task already implemented)

## Pre-Task Checklist (MUST COMPLETE BEFORE STARTING)
- [x] `make check-compliance` passes
- [x] Reviewed existing patterns in:
  - `internal/service/career/burst_fact/detector.go`
  - `internal/cli/intents/burst_management/` (package)
  - `internal/cli/screens/burst_management/modals/`
- [x] Confirmed this is ONE atomic task
- [x] Identified which test files will be created/modified

## Context

The `BurstDetector` is fully implemented with:
- Temporal grouping (6-month window)
- Multi-factor similarity scoring (text, keywords, company, project)
- Cluster detection via BFS
- Confidence scoring

**STATUS: IMPLEMENTED** - The detection is now exposed in the TUI via the "s" key.

## Current Architecture

The burst management intent follows **package-based organization**:

```
internal/cli/intents/burst_management/
├── constants.go           # State enum (type State string)
├── context.go            # IntentContext with repository operations
├── handlers.go           # ScreenResult handlers
├── helpers.go            # Modal helpers, view helpers
├── intent.go             # Init/Update/View (core Bubble Tea)
├── interfaces.go         # BurstService interface
├── messages.go           # All *Msg types
├── result.go             # Result struct
├── types.go              # Intent struct definition
└── *_test.go             # Tests

internal/cli/screens/burst_management/
├── burst_list_screen.go  # Main list screen
└── modals/
    ├── detail_modal.go
    ├── events_modal.go
    ├── facts_modal.go
    ├── edit_modal.go
    └── suggestion_review_modal.go  # For suggestion review
```

### Key Patterns
- **Single Screen + Multiple Modals**: `BurstListScreen` is always active; modals overlay for actions
- **Modal Registry**: Unified modal rendering with priority ordering
- **Screen Result Handler**: Intent implements `ScreenResultHandler` interface
- **Message-Based**: All async operations use typed messages

## Current UX Flow (IMPLEMENTED)
```
Burst List → [s] Suggest Bursts → Loading Modal →
Suggestion Review Modal (accept/reject each) →
Accepted become new bursts → Facts extracted → Return to List
```

## Files Modified

### Intent Package (`internal/cli/intents/burst_management/`)
- [x] `constants.go` - States `StateSuggesting` and `StateSuggestionReview` added
- [x] `messages.go` - `BurstSuggestionsLoadedMsg` and `SuggestionReviewCompleteMsg` added
- [x] `context.go` - Detection handled via service in helpers.go instead
- [x] `types.go` - `suggestionsLoading`, `suggestionsError`, `suggestionModal` added
- [x] `intent.go` - Handles `BurstSuggestionsLoadedMsg`, `SuggestionReviewCompleteMsg`, `FactExtractionCompleteMsg`
- [x] `helpers.go` - `startBurstDetection()`, `handleBurstSuggestionsLoaded()`, `saveAndExtractBurst()` implemented
- [x] `handlers.go` - `case "suggest"` in `handleActionData()` implemented

### Screens/Modals (`internal/cli/screens/burst_management/`)
- [x] `burst_list_screen.go` - "s" key binding returns `NavigateResult{action: "suggest"}`
- [x] `modals/suggestion_review_modal.go` - Full modal implementation

### Tests
- [x] `internal/cli/intents/burst_management/*_test.go`
- [x] `internal/cli/screens/burst_management/modals/suggestion_review_modal_test.go`

## States (in `constants.go`) - IMPLEMENTED

| State | Constant | Description |
|-------|----------|-------------|
| `list` | `StateList` | Main list view (existing) |
| `suggesting` | `StateSuggesting` | Loading state during detection |
| `suggestion_review` | `StateSuggestionReview` | Review individual suggestions |

## Messages (in `messages.go`) - IMPLEMENTED

```go
// BurstSuggestionsLoadedMsg sent when detection completes
type BurstSuggestionsLoadedMsg struct {
    Suggestions []burst_fact.BurstSuggestion
    Error       error
}

// SuggestionReviewCompleteMsg sent when user finishes reviewing
type SuggestionReviewCompleteMsg struct {
    AcceptedSuggestions []burst_fact.BurstSuggestion
    Cancelled           bool
}
```

## Key Bindings - IMPLEMENTED

### List Screen (`burst_list_screen.go`)
| Key | Action |
|-----|--------|
| `s` | Suggest bursts (trigger detection) |

### Suggestion Review Modal
| Key | Action |
|-----|--------|
| `a` | Accept suggestion (create burst + extract facts) |
| `r` | Reject suggestion (skip to next) |
| `n` / `→` | Next suggestion |
| `p` / `←` | Previous suggestion |
| `Esc` | Cancel review, return to list |

## Implementation Checklist

### Phase 1: Verify/Add State Constants - COMPLETE

#### RED Phase
- [x] Test: States `StateSuggesting` and `StateSuggestionReview` exist in `constants.go`
- [x] Test PASSES (already implemented)

#### GREEN Phase
- [x] States exist in `constants.go`:
  ```go
  StateSuggesting       State = "suggesting"
  StateSuggestionReview State = "suggestion_review"
  ```
- [x] Test PASSES

### Phase 2: Verify/Add Message Types - COMPLETE

#### RED Phase
- [x] Test: `BurstSuggestionsLoadedMsg` and `SuggestionReviewCompleteMsg` exist
- [x] Test PASSES (already implemented)

#### GREEN Phase
- [x] Messages exist in `messages.go`
- [x] Test PASSES

### Phase 3: Add Detection Method - COMPLETE (Alternative Implementation)

**Note**: Instead of adding `DetectBursts()` to context.go, detection is handled via
`startBurstDetection()` in helpers.go which calls `Service.SuggestBursts()` directly.
This is architecturally cleaner as it keeps async operations in the intent layer.

#### Implementation (helpers.go:557-610)
- [x] `startBurstDetection()` loads events and calls `Service.SuggestBursts()`
- [x] Shows loading modal during detection
- [x] Returns `BurstSuggestionsLoadedMsg` when complete

### Phase 4: Add Suggestion State to Intent - COMPLETE

- [x] `suggestionsLoading` field in types.go (line 73-74)
- [x] `suggestionsError` field in types.go (line 87-88)
- [x] `suggestionModal` field in types.go (line 118-119)

### Phase 5: Add "s" Key Handler in List Screen - COMPLETE

- [x] burst_list_screen.go (line 182-188):
  ```go
  case "s":
      return nil, &screens.NavigateResult{
          ResultData: map[string]interface{}{"action": "suggest"},
      }
  ```
- [x] Help footer includes "s" for Suggest

### Phase 6: Handle Suggest Action in Intent - COMPLETE

- [x] handlers.go (line 141-145): `case "suggest"` calls `startBurstDetection()`
- [x] helpers.go: `startBurstDetection()` implementation with loading modal

### Phase 7: Handle BurstSuggestionsLoadedMsg - COMPLETE

- [x] intent.go (line 48-49): handles `BurstSuggestionsLoadedMsg`
- [x] helpers.go (line 612-639): `handleBurstSuggestionsLoaded()` implementation
  - Clears loading modal
  - Shows error modal if error
  - Shows info modal if no suggestions
  - Opens suggestion review modal if suggestions found

### Phase 8: Implement Suggestion Review Modal - COMPLETE

- [x] `modals/suggestion_review_modal.go` fully implemented:
  - Progress indicator: "Burst Suggestion X of Y"
  - Confidence score: "Confidence: XX.X%"
  - Event count display
  - Name and description display
  - Navigation keys (n/p, arrows)
  - Accept/reject keys (a/r)
  - Footer badges for all actions

### Phase 9: Handle Accept/Reject in Modal - COMPLETE

- [x] Modal handles 'a' to accept (adds to accepted list, removes from remaining)
- [x] Modal handles 'r' to reject (removes from remaining)
- [x] Intent intercepts 'a' key to save burst immediately (helpers.go `saveAndExtractBurst`)
- [x] Fact extraction triggered on accept
- [x] When modal closes, remaining accepted suggestions are processed

### Phase 10: Handle Empty Suggestions - COMPLETE

- [x] helpers.go (line 624-628): Shows error modal with message:
  "No suggestions were generated from your events. Try adding more events or adjusting detection settings."

### Phase 11: Register Modal in Registry - COMPLETE

- [x] helpers.go (line 491-498): `suggestionModal` registered in `rebuildModalRegistry()`
- [x] Uses `NewViewModalAdapter` for proper integration

### Phase 12: Update Help Text - COMPLETE

- [x] helpers.go (line 75): "s" key badge in StateList help
- [x] helpers.go (line 55-57): `getStateName()` returns "Suggesting Bursts" and "Review Suggestion"

## Pre-Commit Checklist (BEFORE EACH COMMIT)
- [x] `make check-compliance` passes
- [x] Use `make ai-commit FILE=/tmp/commit.txt`
- [x] Commit is atomic (ONE logical change)

## Post-Task Checklist
- [x] `make check-compliance` passes
- [x] All tests pass
- [x] Can trigger detection from list view
- [x] Can review and accept/reject suggestions
- [x] Accepted suggestions create valid bursts

## Acceptance Criteria - ALL MET
- [x] "s" key in list view triggers burst detection
- [x] Detection shows loading modal
- [x] Suggestions display in review modal with confidence scores
- [x] Accept creates burst, reject skips
- [x] Empty state handled with info modal
- [x] All navigation (Esc, etc.) works correctly
- [x] Modal registry properly manages suggestion modal

## Notes

### Implementation Differences from Original Plan

1. **Detection Method**: Instead of `IntentContext.DetectBursts()`, detection is handled in
   `startBurstDetection()` in helpers.go calling `Service.SuggestBursts()`. This keeps async
   operations in the intent layer where they belong.

2. **Immediate Save on Accept**: When user presses 'a', the burst is saved immediately and
   fact extraction is triggered right away (not waiting for modal close). This provides
   instant feedback and ensures data isn't lost if app crashes.

3. **No "Accept All" Key**: The modal allows accept/reject one at a time with navigation.
   This was a design decision for more thoughtful review of suggestions.

4. **Fact Extraction Integration**: Accepted suggestions automatically trigger fact extraction,
   showing a loading modal during the process.

## Rollback Plan (If Needed)
- Remove suggestion-related handlers from handlers.go
- Remove `startBurstDetection()` and related methods from helpers.go
- Remove suggestion modal handling from intent.go
- Remove "s" key binding from burst_list_screen.go
- Keep but don't use suggestion_review_modal.go
