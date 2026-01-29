# BUG-015: CareerEvent accepts events with no company or project association

## Summary

`CareerEvent.Validate()` does not enforce that an event is associated with either
a company or a project. Events with both `Company` and `Project` as empty strings
pass validation and enter the system, causing downstream logic (deduplication,
date calculation, tenure detection) to handle an invalid state defensively rather
than rejecting it at the domain boundary.

## Steps to Reproduce

1. Create a `CareerEvent` with `Company: ""` and `Project: ""`
2. Call `Validate()` on the event
3. Observe it returns `nil` (no error)

## Expected Behavior

`Validate()` should return an error when both `Company` and `Project` are empty.
Every career event must be associated with at least one organisational context.

## Actual Behavior

`Validate()` only checks `Text`, `Date`, `Tags`, and `Categories`. The `Company`
and `Project` fields are completely unchecked. Events with no association pass
through the service layer (`CaptureEvent`, `UpdateEvent`) without error.

## Root Cause

**File**: `internal/domain/career/event.go:26-48`

```go
func (ce *CareerEvent) Validate() error {
    if err := ce.validateText(); err != nil { return err }
    if err := ce.validateDate(); err != nil { return err }
    if err := ce.validateTags(); err != nil { return err }
    if err := ce.validateCategories(); err != nil { return err }
    return nil
    // No validateAssociation() call
}
```

There is no `validateAssociation()` method requiring `Company != ""` or
`Project != ""`.

## Fix Plan

Add a `validateAssociation()` method to `CareerEvent`:

```go
func (ce *CareerEvent) validateAssociation() error {
    if strings.TrimSpace(ce.Company) == "" && strings.TrimSpace(ce.Project) == "" {
        return errors.New("event must be associated with a company or project")
    }
    return nil
}
```

Call it from `Validate()` alongside the existing validators.

## Impact

Existing tests that construct events with both fields empty will need updating to
provide at least one association. The database layer (`gorm`) does not enforce
`NOT NULL` on either column, so a migration may be needed for data integrity.

## Regression Tests

- Event with Company set and Project empty passes validation
- Event with Project set and Company empty passes validation
- Event with both Company and Project set passes validation
- Event with both Company and Project empty fails validation
- `CaptureEvent` rejects events with no association
- `UpdateEvent` rejects events with no association

## Related

- BUG-013: Bullet deduplication merges across companies (empty company causes incorrect merging)
- BUG-014: Date calculation ignores primary company

## Severity

- [x] Medium - Invalid data enters the system silently; downstream code must handle defensively
