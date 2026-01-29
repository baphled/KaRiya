# BUG-011: Import creates unconfirmed bursts with no fact linking

## Summary

CSV import auto-saves burst suggestions as unconfirmed bursts and extracts facts
per-event only, never linking facts to their parent bursts. Additionally, the TUI
confirm burst flow promises fact extraction in its dialog text but does not deliver.

## Steps to Reproduce

1. Import career events via `kariya --import events.csv`
2. Open the TUI and navigate to Burst Management
3. Observe all bursts show `Confirmed: false`
4. Select a burst and view its facts (press `f`)
5. Observe 0 facts despite facts existing in the system
6. Press `c` to confirm the burst, accept the dialog
7. Observe the burst is confirmed but still has 0 facts

## Expected Behavior

- Imported bursts should either be confirmed or clearly documented as requiring review.
- Facts extracted during import should be linked to their parent bursts via `SourceBurstID`.
- The confirm burst flow should extract burst-level facts when confirming a burst with 0 facts.

## Actual Behavior

- All 134/135 imported bursts have `Confirmed = false` because `SaveBurstSuggestions`
  calls `ConfirmBurst()` which does NOT set the `Confirmed` field to `true`.
- All 860 facts have `source_event_id` set but `source_burst_id` is empty for all of them.
- The confirm dialog says "Mark this burst as confirmed and extract facts?" but only
  sets `Confirmed = true` without triggering `extractFactsForBurst()`.

## Root Causes

### Issue 1: `SaveBurstSuggestions` creates unconfirmed bursts

**File**: `internal/service/career/service.go:373-413`

`SaveBurstSuggestions` creates bursts without `Confirmed: true` and calls `ConfirmBurst()`
which is misleadingly named - it only validates and persists, without setting `Confirmed`.

### Issue 2: Import links facts to events, not bursts

**File**: `internal/cli/importer/service.go:151-200`

`ImportRows()` calls `ExtractFactsFromEvent()` which sets `SourceEventID` but never
`SourceBurstID`. Despite having just created bursts from these events, the facts are
not linked to those bursts.

### Issue 3: Confirm burst flow does not extract facts

**File**: `internal/cli/intents/burst_management/helpers.go:261-289`

`confirmBurst()` only sets `Confirmed = true`. The comment explicitly states:
"Note: Fact extraction happens when accepting suggestions, not here."
But the dialog text promises fact extraction.

## Database Evidence

```
Total bursts: 135 (134 unconfirmed, 1 confirmed)
Total facts: 860 (ALL with source_event_id, NONE with source_burst_id)
Total events: 861
```

## Environment

- OS: Linux
- Go version: See go.mod
- Branch: fix/burst-import-facts-linking

## Severity

- [x] High - Major feature broken

## Fix Plan

1. **Rename `ConfirmBurst` to `SaveBurst`** in the service layer. Create a proper
   `ConfirmBurst` that sets `Confirmed = true` and `ConfirmedAt`. Update
   `SaveBurstSuggestions` to call `SaveBurst`.

2. **Link facts to bursts during import**. After saving burst suggestions in
   `ImportRows()`, update facts whose `SourceEventID` matches burst event IDs to
   also set `SourceBurstID`.

3. **Confirm burst flow should extract facts**. Modify `confirmBurst()` in the TUI
   intent to trigger `extractFactsForBurst()` after setting `Confirmed = true`.

## Regression Tests

- Service: `SaveBurst` validates and saves without confirming
- Service: `ConfirmBurst` sets `Confirmed=true` and `ConfirmedAt`
- Service: `SaveBurstSuggestions` creates unconfirmed bursts via `SaveBurst`
- Importer: `ImportRows` links facts to parent bursts via `SourceBurstID`
- TUI: Confirming a burst triggers fact extraction
