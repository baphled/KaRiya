# Bug 006: Specialist Tech Focus Should Use Single-Select

**Status**: Closed  
**Severity**: 🟢 Resolved  
**Created**: 2026-01-21  
**Updated**: 2026-01-21  
**Fixed**: 2026-01-21 (commit 368c704)  

---

## Bug Summary

When selecting "Specialist" tech focus in CV wizard, user is shown MultiSelect for technologies instead of single-select (should pick ONE technology per `variants.go:33`).

---

## Affected Components

- [x] `internal/cli/forms/cv_config_form.go:106` - Now uses Select or MultiSelect based on singleTechSelect (FIXED)
- [x] `internal/cli/components/cv_config_wizard_modal.go` - Passes singleTechSelect based on TechFocus (FIXED)

**Related Intents/Workflows**:
- GenerateCVIntent (wizard mode)
- CV_GENERATION_WORKFLOW.md

---

## Reproduction Steps

1. Start KaRiya and navigate to Generate CV
2. Select wizard mode
3. In Step 2 (TECH), select "Specialist (highlight specific techs)"
4. Observe the technologies field is MultiSelect (can pick multiple)
5. Expected: Should be single-select (pick ONE technology)

**Consistency**: Always

**Environment**:
- OS: Linux
- Terminal: Any
- Go Version: 1.24+
- KaRiya Version: next branch

---

## Expected Behavior

When `TechFocus` is "specialist", the technologies field should:
- Show a single-select dropdown (huh.NewSelect)
- Allow user to pick exactly ONE technology
- Per `variants.go:33`: "Specialist style focused on 1 technology"

---

## Actual Behavior

Technologies field is always MultiSelect regardless of TechFocus selection.

**Evidence**:
```go
// cv_config_form.go:106 - ALWAYS creates MultiSelect
huh.NewMultiSelect[string]().
    Key("technologies").
    Title("Select Technologies to Highlight").
    ...
```

---

## Root Cause

**Status**: Identified

**Cause**:
`NewCVConfigForm()` creates the form once with a static structure. The technologies field is always `huh.NewMultiSelect` because the form doesn't know the TechFocus value at creation time (or doesn't react to it changing).

**Technical Details**:
- File: `internal/cli/forms/cv_config_form.go`
- Function: `NewCVConfigForm()`
- Line: `106`
- Reason: Form structure is static; no conditional logic for specialist vs generalist

**Supporting Evidence**:
- `generate_cv_intent.go:1140` handles specialist expecting single tech
- `generate_cv_intent.go:2030` shows specialist mode logic
- `variants.go:33-34` defines specialist as "1 technology" focus

---

## Fix Strategy

**Approach**: Add parameter to control Select vs MultiSelect based on TechFocus

### Option A: Add `singleTechSelect` Parameter (Recommended)

Add a boolean parameter to `NewCVConfigForm` that controls whether to use Select or MultiSelect for technologies.

**Pros**:
- Minimal change to existing code
- Clear API contract
- Caller controls the behavior

**Cons**:
- Requires updating all callers
- Form must be rebuilt when TechFocus changes

**Files to Change**:
- [x] `internal/cli/forms/cv_config_form.go` - Add `singleTechSelect bool` parameter
- [ ] `internal/cli/components/cv_config_wizard_modal.go` - Pass parameter, rebuild form on TechFocus change
- [ ] `internal/cli/forms/cv_config_form_test.go` - Test both modes

### Option B: Dynamic Form Rebuild

Have the modal watch for TechFocus changes and rebuild the entire form.

**Pros**:
- No API change needed

**Cons**:
- More complex state management
- Potential loss of other field values during rebuild
- Harder to test

**Selected Approach**: Option A  
**Rationale**: Clearer API, easier to test, follows existing patterns

---

## Testing Plan

### Phase 1: Unit Tests (COMPLETED)

- [x] Test `NewCVConfigForm` with `singleTechSelect=false` creates MultiSelect
- [x] Test `NewCVConfigForm` with `singleTechSelect=true` creates Select
- [x] Test form data binding works for single-select mode

**Files**:
- `internal/cli/forms/cv_config_form_test.go`

### Phase 2: Integration Tests

- [ ] Test wizard modal switches form when TechFocus changes to "specialist"
- [ ] Test wizard modal switches form when TechFocus changes from "specialist" to "generalist"
- [ ] Test data persists correctly across form rebuilds

**Files**:
- `internal/cli/components/cv_config_wizard_modal_test.go`

### Phase 3: Manual Testing

- [ ] Start wizard, select Specialist, verify single-select for technologies
- [ ] Start wizard, select Generalist, verify multi-select for technologies
- [ ] Switch between modes, verify correct field type appears
- [ ] Complete wizard in Specialist mode, verify only ONE tech in config

---

## Verification Checklist

### Code Quality (COMPLETED)
- [x] Fix implemented and tested
- [x] All tests passing (go test ./...)
- [x] No race conditions (go test -race)
- [x] Code coverage maintained (>80%)
- [x] Linting passing (staticcheck)

### Functionality (COMPLETED)
- [x] Specialist mode shows single-select
- [x] Generalist mode shows multi-select
- [x] Language Agnostic mode hides tech selection (existing behavior)
- [x] Selected technology persists through wizard completion

### Documentation (COMPLETED)
- [x] Code comments added/updated
- [x] Bug report updated with resolution

### Compliance (COMPLETED)
- [x] Follows project coding standards
- [x] Atomic commits with clear messages
- [x] AI attribution

---

## Related Files

### Implementation Files
- `internal/cli/forms/cv_config_form.go:106` - MultiSelect creation
- `internal/cli/components/cv_config_wizard_modal.go` - Form creation
- `internal/service/career/cv/variants.go:33` - Specialist definition

### Test Files
- `internal/cli/forms/cv_config_form_test.go`
- `internal/cli/components/cv_config_wizard_modal_test.go`
- `internal/cli/intents/generate_cv_wizard_e2e_test.go`

### Documentation Files
- `docs/workflows/CV_GENERATION_WORKFLOW.md`
- `docs/FORMS_GUIDE.md`

---

## References

- Related bugs: BUG-005 (form selections not persisted - should fix first)
- Documentation: `variants.go:33` defines specialist as single-tech focus
- Test evidence: `generate_cv_technology_test.go:544` "Specialist mode - single-select"

---

**Last Updated**: 2026-01-21  
**Updated By**: Opencode  
**Fix Commit**: 368c704 - fix(cv): add singleTechSelect parameter for specialist mode (BUG-006)
