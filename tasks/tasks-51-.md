# Task 51: Standardize Wizard Pattern with WizardBehavior[T]

## Status: ✅ Complete

## Overview
- **Goal**: Create a generic `WizardBehavior[T]` in `behaviors/` and fix existing architecture violations in wizard modals
- **Time Estimate**: 6-8 hours
- **Prerequisites**: BUG-005 fix merged (form selection persistence)
- **Priority**: HIGH - Fixes architecture violations + reduces code duplication

---

## Problem Description

### Current Architecture Violations
Both wizard modals **violate** the Forms/Huh Architecture rules:
- `internal/cli/components/cv_config_wizard_modal.go` - imports `huh` directly
- `internal/cli/components/onboarding_wizard_modal.go` - imports `huh` directly

Per AGENTS.md:
> **Imports `huh` outside of `forms/` package** = Architecture violation to REFUSE

| Package | Can Import `huh`? | Current Status |
|---------|-------------------|----------------|
| `forms/` | **YES** (only place) | OK |
| `models/` | **NO** | OK |
| `components/` | **NO** | **FIXED** |
| `behaviors/` | **NO** | OK |

### Duplicated Patterns
Both wizards also share duplicated code:
- Multi-step form handling
- Data struct for collected values
- State tracking (visible/completed/cancelled)
- UIKit rendering (title, step indicator, content, footer)
- Keyboard handling (Tab, Enter, Esc)

---

## Session Contract Acknowledgment
- [x] Ran `make session-start` and it passed
- [x] Acknowledge and commit to following all workflow rules
- [x] Token count: within limits

## Pre-Task Checklist (MUST COMPLETE BEFORE STARTING)
- [x] `make check-compliance` passes
- [x] Reviewed AGENTS.md Forms/Huh Architecture section
- [x] Reviewed existing patterns in: `internal/cli/behaviors/table_behavior.go`
- [x] Reviewed docs/FORMS_GUIDE.md and docs/rules/FORMS_WORKFLOW_GUIDE.md
- [x] Confirmed understanding of layer hierarchy

---

## Correct Architecture

### Layer Responsibilities (IMPLEMENTED)

```
forms/                          # CAN import huh
├── cv_config_form.go           # CV form builder (pre-existing)
├── onboarding_wizard_form.go   # NEW: Onboarding form builder
└── wizard_form_adapter.go      # NEW: WizardFormAdapter wrapping *huh.Form

behaviors/                      # CANNOT import huh
└── wizard_behavior.go          # NEW: WizardForm interface + WizardBehavior[T]

components/                     # CANNOT import huh (FIXED)
├── cv_config_wizard_modal.go   # MIGRATED: Uses forms/ + behaviors/
└── onboarding_wizard_modal.go  # MIGRATED: Uses forms/ + behaviors/
```

### Design Decision: No models/ wrappers
The task originally specified models/ wrappers, but models/ is DEPRECATED per AGENTS.md.
Instead, WizardFormAdapter in forms/ satisfies the behaviors.WizardForm interface via
duck typing, keeping both packages decoupled with no circular dependencies.

### WizardForm Interface (in behaviors/)
```go
type WizardForm interface {
    Init() tea.Cmd
    Update(msg tea.Msg) tea.Cmd
    View() string
    IsCompleted() bool
    IsAborted() bool
    CurrentStep() int
    TotalSteps() int
    SetDimensions(width, height int)
}
```

### WizardBehavior (in behaviors/)
```go
type WizardBehavior[T any] struct {
    form      WizardForm  // Interface, NOT *huh.Form
    data      *T
    visible   bool
    completed bool
    cancelled bool
    skipped   bool
}
```

---

## Files Created

### forms/ package
- [x] `internal/cli/forms/onboarding_wizard_form.go` - Onboarding wizard form builder
- [x] `internal/cli/forms/onboarding_wizard_form_test.go` - Tests
- [x] `internal/cli/forms/wizard_form_adapter.go` - WizardFormAdapter wrapping *huh.Form
- [x] `internal/cli/forms/wizard_form_adapter_test.go` - Tests

### behaviors/ package
- [x] `internal/cli/behaviors/wizard_behavior.go` - WizardForm interface + WizardBehavior[T]
- [x] `internal/cli/behaviors/wizard_behavior_test.go` - Tests with mock WizardForm

## Files Modified
- [x] `internal/cli/components/cv_config_wizard_modal.go` - Remove huh import, use forms/ + behaviors/
- [x] `internal/cli/components/onboarding_wizard_modal.go` - Remove huh import, use forms/ + behaviors/

---

## TDD Checklist (COMPLETED)

### Phase 3 (done first - foundation): WizardBehavior in behaviors/
- [x] Write failing test: WizardBehavior creation
- [x] Implement `WizardBehavior` struct with `WizardForm` interface
- [x] Write failing test: state management (complete, cancel, skip, reset)
- [x] Implement state tracking
- [x] Write failing test: form delegation
- [x] Implement form delegation

### Phase 1: Form Builders in forms/
- [x] CV wizard form already existed in `forms/cv_config_form.go`
- [x] Write failing test: Onboarding wizard form creation
- [x] Implement `forms.NewOnboardingWizardForm()`

### Phase 2: WizardFormAdapter in forms/
- [x] Write failing test: WizardFormAdapter creation and methods
- [x] Implement `forms.WizardFormAdapter` wrapping *huh.Form
- [x] Verify duck typing satisfies behaviors.WizardForm

### Phase 4: Migrate Wizard Modals
- [x] Ensure existing tests pass before changes
- [x] Refactor `onboarding_wizard_modal.go` to use forms/ + behaviors/
- [x] **Verify NO huh import in onboarding modal**
- [x] Refactor `cv_config_wizard_modal.go` to use forms/ + behaviors/
- [x] **Verify NO huh import in CV wizard modal**
- [x] All existing tests still pass (192/192 component specs)

---

## Post-Task Checklist
- [x] `make check-compliance` passes (0 violations, 25 pre-existing warnings)
- [x] `grep -r "charmbracelet/huh" internal/cli/components/cv_config_wizard_modal.go` returns **nothing**
- [x] `grep -r "charmbracelet/huh" internal/cli/components/onboarding_wizard_modal.go` returns **nothing**
- [x] `grep -r "charmbracelet/huh" internal/cli/behaviors/` returns **nothing**
- [x] All checkboxes above completed
- [x] Task marked complete

---

## Acceptance Criteria
- [x] **NO huh imports in wizard modals or behaviors/** (architecture compliance)
- [x] Form builders exist in `forms/` package
- [x] `WizardFormAdapter` in forms/ replaces models/ wrappers (models/ is deprecated)
- [x] `WizardBehavior[T]` created with `WizardForm` interface
- [x] CV config wizard uses new architecture
- [x] Onboarding wizard uses new architecture
- [x] All existing wizard tests pass (192/192 component specs, 151 behavior specs, 190 forms specs)
- [x] New unit tests included
- [x] Full project builds cleanly

## Note
- `export_options_modal.go` in components/ still imports huh - this is a separate modal,
  not a wizard, and is outside the scope of this task.

---

## Commits
1. `chore(intents): remove migrated pending test from consistency suite`
2. `feat(behaviors): add WizardForm interface and WizardBehavior[T] generic`
3. `feat(components): extract onboarding wizard form builder to forms package`
4. `feat(components): add WizardFormAdapter wrapping huh.Form in forms package`
5. `refactor(components): remove huh import from onboarding wizard modal`
6. `refactor(components): remove huh import from CV config wizard modal`
7. `docs(behaviors): add required godoc sections to WizardBehavior methods`
