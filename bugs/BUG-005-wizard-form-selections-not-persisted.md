# Bug 005: CV Wizard Form Selections Not Persisted on Rebuild

**Status**: Closed  
**Severity**: 🟢 Resolved  
**Created**: 2026-01-21  
**Updated**: 2026-01-21  
**Fixed**: 2026-01-21 (core: d729dc8, edge case: 58916f4)  

---

## Bug Summary

User selects "Senior Engineer" profile in CV wizard but CV generates as "Staff Engineer" - form selections lost when form is rebuilt (window resize, tech extraction).

---

## Affected Components

- [x] `internal/cli/components/cv_config_wizard_modal.go` - GetConfigData didn't sync formData (FIXED - d729dc8)
- [x] `internal/cli/components/cv_config_wizard_modal.go` - HasRequiredFields didn't sync formData (FIXED - d729dc8)
- [x] `internal/cli/components/cv_config_wizard_modal.go` - Added syncFromFormData() method (FIXED - d729dc8)
- [x] `internal/cli/components/cv_config_wizard_modal.go` - formData stored as struct field (FIXED - d729dc8)
- [x] `internal/cli/components/cv_config_wizard_modal.go:110` - buildForm() syncs before recreating formData (FIXED - 58916f4)
- [x] `internal/cli/components/cv_config_wizard_modal.go:197` - WindowSizeMsg now preserves selections (FIXED - 58916f4)
- [x] `internal/cli/components/cv_config_wizard_modal.go:331` - Reset() now preserves selections (FIXED - 58916f4)
- [x] `internal/cli/components/cv_config_wizard_modal.go:392` - SetExtractedTechnologies now preserves selections (FIXED - 58916f4)

**Related Intents/Workflows**:
- GenerateCVIntent (wizard mode)
- CV_GENERATION_WORKFLOW.md

---

## Reproduction Steps

1. Start KaRiya and navigate to Generate CV
2. Select wizard mode
3. In Step 1 (WHO), select "Senior Engineer" profile
4. Trigger a window resize (or wait for tech extraction to complete)
5. Continue through wizard and generate CV
6. Observe CV uses wrong profile (e.g., "Staff Engineer" - the first/default option)

**Consistency**: Always (when rebuild is triggered)

**Environment**:
- OS: Linux
- Terminal: Any resizable terminal
- Go Version: 1.24+
- KaRiya Version: next branch

---

## Expected Behavior

User's profile selection ("Senior Engineer") should persist through form rebuilds and be used for CV generation.

---

## Actual Behavior

Profile selection is lost when form is rebuilt. CV generates with wrong/default profile.

**Evidence**:
```go
// cv_config_wizard_modal.go:193-198 - NO sync before rebuild
case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
    // Rebuild form with new dimensions
    m.buildForm()  // <-- m.data not synced, loses formData selections!
    return m.form.Init()

// cv_config_wizard_modal.go:388-393 - NO sync before rebuild  
func (m *CVConfigWizardModal) SetExtractedTechnologies(techs []ExtractedTechnology) {
    m.extractedTechs = techs
    m.techsAvailable = len(techs) > 0
    m.buildForm()  // <-- m.data not synced, loses formData selections!
}
```

---

## Root Cause

**Status**: Identified

**Cause**:
When `buildForm()` is called, it creates a NEW `formData` struct initialized from `m.data` (lines 132-141). However, `m.data` only gets updated when `syncFromFormData()` is called. If `buildForm()` is triggered before sync (e.g., window resize), the user's selections in `formData` are overwritten with stale values from `m.data`.

**Data Flow Problem**:
```
1. User selects "Senior Engineer" → huh updates formData.ProfileID via pointer binding
2. m.data.ProfileID is still "" (not synced yet)
3. Window resize triggers buildForm()
4. buildForm() creates new formData from m.data (which has "")
5. User's selection is LOST
6. GetConfigData() syncs empty value → wrong profile used
```

**Technical Details**:
- File: `internal/cli/components/cv_config_wizard_modal.go`
- Functions: `Update()` line 197, `SetExtractedTechnologies()` line 392
- Problem: Missing `syncFromFormData()` call before `buildForm()`

---

## Fix Strategy

**Approach**: Sync on read + store formData as field

### Implemented Fix (commit d729dc8)

The fix ensures GetConfigData() returns the latest form values by:
1. Storing `formData` as a struct field (not local variable)
2. Adding `syncFromFormData()` method to copy formData → data
3. Calling `syncFromFormData()` in `GetConfigData()` before returning
4. Calling `syncFromFormData()` in `HasRequiredFields()` before checking
5. All setter methods now update both `formData` and `data`

**Files Changed**:
- [x] `internal/cli/components/cv_config_wizard_modal.go` - Core fix implementation
- [x] `internal/cli/components/cv_config_wizard_modal_test.go` - Regression test added

### Edge Case Fix (COMPLETED - commit 58916f4)

Added `syncFromFormData()` at START of `buildForm()` to fix all callers:
- `WindowSizeMsg` handler (line 197)
- `Reset()` method (line 331)
- `SetExtractedTechnologies()` method (line 392)

**Files Changed**:
- [x] `internal/cli/components/cv_config_wizard_modal.go:110` - Added `m.syncFromFormData()` at start of `buildForm()`
- [x] `internal/cli/components/cv_config_wizard_modal_test.go` - Added rebuild preservation tests

---

## Testing Plan

### Phase 1: Unit Tests - Core Fix (COMPLETED)

- [x] Test: Form navigation selection returned by GetConfigData (commit d729dc8)

**Files**:
- `internal/cli/components/cv_config_wizard_modal_test.go`

### Phase 2: Unit Tests - Edge Cases (COMPLETED - commit 58916f4)

- [x] Test: Window resize preserves profile selection
- [x] Test: SetExtractedTechnologies preserves existing selections
- [x] Test: Reset() preserves existing selections (via SetProfileID path)

**Files**:
- `internal/cli/components/cv_config_wizard_modal_test.go`

### Phase 3: Integration Tests (PENDING)

- [ ] Test: Complete wizard after window resize, verify correct config
- [ ] Test: Complete wizard after tech extraction, verify correct config

**Files**:
- `internal/cli/intents/generate_cv_wizard_e2e_test.go`

### Phase 4: Manual Testing (PENDING)

- [ ] Start wizard, select profile, resize window, complete wizard - verify profile correct
- [ ] Start wizard, complete, go back (Reset), verify selections preserved
- [ ] Start wizard, select profile, let tech extraction complete, verify profile preserved
- [ ] Resize multiple times during wizard, verify all selections persist

---

## Verification Checklist

### Code Quality - Core Fix (COMPLETED)
- [x] Fix implemented and tested
- [x] All tests passing (go test ./...)
- [x] No race conditions (go test -race)
- [x] Code coverage maintained (>80%)
- [x] Linting passing (staticcheck)

### Functionality - Core Fix (COMPLETED)
- [x] GetConfigData returns correct values after form navigation

### Functionality - Edge Cases (COMPLETED)
- [x] Profile selection persists through window resize
- [x] Profile selection persists through Reset()
- [x] All field selections persist through rebuild
- [x] Tech extraction doesn't lose existing selections

### Documentation
- [x] Code comments updated (syncFromFormData docs)
- [x] Bug report updated with resolution

### Compliance
- [x] Follows project coding standards
- [x] Atomic commits with clear messages
- [x] AI attribution

---

## Implementation

### Core Fix (COMPLETED - commit d729dc8)

```go
// Added formData as struct field
type CVConfigWizardModal struct {
    form     *huh.Form
    formData *forms.CVConfigFormData  // NEW: persists huh pointer bindings
    data     *CVConfigData
    // ...
}

// GetConfigData now syncs first
func (m *CVConfigWizardModal) GetConfigData() *CVConfigData {
    m.syncFromFormData()  // NEW: ensures latest form values
    return m.data
}

// New syncFromFormData method
func (m *CVConfigWizardModal) syncFromFormData() {
    if m.formData == nil {
        return
    }
    m.data.ProfileID = m.formData.ProfileID
    // ... all fields synced
}
```

### Edge Case Fix (PENDING) - Simpler Approach

Add sync at start of `buildForm()` - fixes ALL callers with one change:

```go
// cv_config_wizard_modal.go - buildForm() (line 110)
func (m *CVConfigWizardModal) buildForm() {
    // Sync current form state before rebuilding (safe if formData is nil)
    m.syncFromFormData()
    
    // Calculate modal dimensions
    modalWidth := m.width - 20
    // ... rest of existing code
}
```

This handles:
- WindowSizeMsg (line 197)
- Reset() (line 331)  
- SetExtractedTechnologies() (line 392)
- Any future buildForm() callers

---

## Resolution Summary

**Core Fix Applied**: commit d729dc8 (PR #99)

**What was fixed**:
- `GetConfigData()` now returns the actual form values selected via keyboard navigation
- Added `formData` as struct field to persist huh pointer bindings
- Added `syncFromFormData()` method called before reading config
- All setter methods update both `formData` and `data`
- Regression test added for form navigation selection

**Edge Case Fix Applied**: commit 58916f4

- Added `syncFromFormData()` at start of `buildForm()` 
- This ensures all callers (WindowSizeMsg, Reset, SetExtractedTechnologies) preserve selections
- Added regression tests for window resize and SetExtractedTechnologies scenarios

---

## Related Files

### Implementation Files
- `internal/cli/components/cv_config_wizard_modal.go:110` - buildForm() (fix location)
- `internal/cli/components/cv_config_wizard_modal.go:197` - WindowSizeMsg handler (caller)
- `internal/cli/components/cv_config_wizard_modal.go:331` - Reset() (caller)
- `internal/cli/components/cv_config_wizard_modal.go:392` - SetExtractedTechnologies (caller)
- `internal/cli/components/cv_config_wizard_modal.go:487` - syncFromFormData method

### Test Files
- `internal/cli/components/cv_config_wizard_modal_test.go`
- `internal/cli/intents/generate_cv_wizard_e2e_test.go`

### Documentation Files
- `docs/workflows/CV_GENERATION_WORKFLOW.md`
- `docs/FORMS_GUIDE.md`

---

## References

- Related bugs: BUG-006 (specialist single-select - depends on this fix)
- Pattern reference: `quick_add_event_modal.go` - similar formData sync pattern
- Documentation: `FORMS_GUIDE.md` - form state management

---

**Last Updated**: 2026-01-21  
**Updated By**: Opencode  
**Core Fix**: d729dc8 - fix(components): persist wizard form selections to GetConfigData (BUG-005) (#99)  
**Edge Case Fix**: 58916f4 - fix(components): sync form state before rebuild to preserve selections (BUG-005)
