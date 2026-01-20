# BUG-004: Generate CV shows stub data when no events exist

**Status**: FIXED  
**Severity**: HIGH  
**Reported**: 2026-01-20  
**Fixed**: 2026-01-20  
**Component**: `internal/cli/app/app.go`  
**Affected Versions**: All versions with CV generation prior to fix  
**Reporter**: User (via onboarding flow testing)

---

## Summary

When a user deletes their database (or has a fresh install) and completes the onboarding wizard, selecting "Generate CV" from the menu shows fake stub/test data ("Test event for navigation integration") instead of a proper empty state message informing the user they need to add career events first.

---

## Steps to Reproduce

1. Delete the database file: `rm ~/.kariya/events.db`
2. Run KaRiya: `go run cmd/cli/main.go`
3. Complete the onboarding wizard (enter name, email)
4. From the main menu, select "Generate CV"
5. Observe the CV generation wizard

---

## Expected Behavior

A warning/info modal should appear with a message like:

> **No Career Events**
> 
> You need to add career events before generating a CV.
> 
> Use 'Capture Event' from the main menu to record your achievements, projects, and career milestones.

The user should be returned to the main menu after dismissing the modal (via Enter, Space, or Esc).

---

## Actual Behavior

The CV generation wizard starts and shows fake stub data:
- Event: "Test event for navigation integration" (ID: `ev-stub`)
- Fact: "Test fact for navigation integration" (ID: `fact-stub`)

This is misleading as users see fake data that doesn't belong to them.

---

## Root Cause

In `internal/cli/app/app.go` lines 718-725, when registering the `generate_cv` intent, stub data is injected if no events/facts exist:

```go
events, err := careerService.GetEventRepository().List(ctx, careerrepo.ListFilters{Limit: 100})
if err != nil || len(events) == 0 {
    events = []*career.CareerEvent{{ID: "ev-stub", Text: "Test event for navigation integration", Date: time.Now()}}
}
facts, err := careerService.GetFactRepository().List(ctx, careerrepo.FactListFilters{Limit: 100})
if err != nil || len(facts) == 0 {
    facts = []*career.Fact{{ID: "fact-stub", Text: "Test fact for navigation integration", ...}}
}
```

This was added in commit `bf5f3e6` (CV generation wizard modal implementation) as a workaround for navigation testing, but it now affects real users.

---

## Environment

- OS: Linux (any)
- Go version: 1.24+
- Branch: `next`

---

## Severity

- [ ] Critical - Application crash/data loss
- [x] High - Major feature broken
- [ ] Medium - Feature partially broken
- [ ] Low - Minor issue/cosmetic

**Rationale**: Users see fake data that could be confusing or misleading. The Generate CV feature doesn't properly handle the empty state, which is a core UX failure.

---

## Impact

| Impact Area | Description |
|-------------|-------------|
| **User Confusion** | Users see fake "test" data they didn't create |
| **Trust** | Users may not trust the CV generation if it shows bogus data |
| **Onboarding UX** | New users completing onboarding are immediately shown broken behavior |
| **Feature Accessibility** | No guidance on what to do to actually use CV generation |

---

## Proposed Fix

1. **Remove stub data fallback** from `app.go` lines 718-725
2. **Create `InfoModal` component** following the `DeleteConfirmModal` pattern
3. **Add empty state check** in `handleMenuInput` before activating `generate_cv` intent
4. **Show warning modal** when user has no career events, with guidance to use "Capture Event"
5. **Dismiss modal** on Enter, Space, or Esc, returning user to main menu

### Files to Modify

| File | Change |
|------|--------|
| `internal/cli/components/info_modal.go` | New - Info/warning modal component |
| `internal/cli/components/info_modal_test.go` | New - Tests for info modal |
| `internal/cli/app/app.go` | Remove stub data, add modal field, add empty state check |
| `internal/testutil/e2e/generate_cv_empty_state_e2e_test.go` | New - E2E regression test |

---

## Testing Checklist

- [x] E2E test: Selecting "Generate CV" with no events shows warning modal
- [x] E2E test: Modal does NOT show stub data
- [x] E2E test: Esc dismisses modal and returns to menu
- [x] E2E test: Enter dismisses modal and returns to menu
- [x] E2E test: With events, Generate CV works normally (no modal)
- [x] Unit test: InfoModal displays correct title and message
- [x] Unit test: InfoModal dismisses on Enter, Space, Esc
- [ ] Manual test: Delete database, complete onboarding, select Generate CV

---

## Related

- Commit `bf5f3e6` - Added stub data as workaround
- `internal/cli/components/delete_confirm_modal.go` - Pattern for modal component
- `internal/cli/app/app.go:718-725` - Location of bug

---

## Notes

The stub data was originally added for navigation testing during development. The proper fix is to handle the empty state gracefully with user feedback, not to inject fake data.

---

**Last Updated**: 2026-01-20
