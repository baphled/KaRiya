# Task 51: Standardize Wizard Pattern with WizardBehavior[T]

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
| `components/` | **NO** | **VIOLATION** |
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
- [ ] Ran `make session-start` and it passed
- [ ] Acknowledge and commit to following all workflow rules
- [ ] Token count: _____ (must be < 50k to start)

## Pre-Task Checklist (MUST COMPLETE BEFORE STARTING)
- [ ] `make check-compliance` passes
- [ ] Reviewed AGENTS.md Forms/Huh Architecture section
- [ ] Reviewed existing patterns in: `internal/cli/behaviors/table_behavior.go`
- [ ] Reviewed docs/FORMS_GUIDE.md and docs/rules/FORMS_WORKFLOW_GUIDE.md
- [ ] Confirmed understanding of layer hierarchy

---

## Correct Architecture

### Layer Responsibilities

```
forms/                          # CAN import huh
├── cv_config_wizard_form.go    # Form builder, wraps huh
└── onboarding_wizard_form.go   # Form builder, wraps huh

models/                         # CANNOT import huh
├── cv_config_wizard_model.go   # Form wrapper, uses forms.Form
└── onboarding_wizard_model.go  # Form wrapper, uses forms.Form

behaviors/                      # CANNOT import huh
└── wizard_behavior.go          # Generic behavior, uses WizardForm interface

components/                     # CANNOT import huh
├── cv_config_wizard_modal.go   # Uses models + behaviors
└── onboarding_wizard_modal.go  # Uses models + behaviors
```

### WizardForm Interface (in behaviors/)
```go
// WizardForm abstracts form operations without exposing huh
type WizardForm interface {
    Update(msg tea.Msg) (tea.Cmd, error)
    View() string
    IsCompleted() bool
    IsAborted() bool
    CurrentStep() int
    TotalSteps() int
}
```

### WizardBehavior (in behaviors/)
```go
type WizardBehavior[T any] struct {
    form       WizardForm  // Interface, NOT *huh.Form
    data       *T
    visible    bool
    completed  bool
    cancelled  bool
    width      int
    height     int
}
```

---

## Files to Create

### forms/ package
- [ ] `internal/cli/forms/cv_config_wizard_form.go` - CV wizard form builder
- [ ] `internal/cli/forms/onboarding_wizard_form.go` - Onboarding wizard form builder

### models/ package
- [ ] `internal/cli/models/cv_config_wizard_model.go` - CV wizard form wrapper
- [ ] `internal/cli/models/onboarding_wizard_model.go` - Onboarding wizard form wrapper

### behaviors/ package
- [ ] `internal/cli/behaviors/wizard_behavior.go` - Generic wizard behavior
- [ ] `internal/cli/behaviors/wizard_behavior_test.go` - Unit tests

## Files to Modify
- [ ] `internal/cli/components/cv_config_wizard_modal.go` - Remove huh import, use models + behaviors
- [ ] `internal/cli/components/onboarding_wizard_modal.go` - Remove huh import, use models + behaviors
- [ ] Update corresponding test files

---

## TDD Checklist (MUST COMPLETE IN ORDER)

### Phase 1: Form Builders in forms/ (RED -> GREEN -> REFACTOR)
- [ ] Write failing test: CV wizard form creation
- [ ] Implement `forms.NewCVConfigWizardForm()`
- [ ] Write failing test: Onboarding wizard form creation
- [ ] Implement `forms.NewOnboardingWizardForm()`
- [ ] Refactor for consistency

### Phase 2: Form Wrappers in models/ (RED -> GREEN -> REFACTOR)
- [ ] Write failing test: CV wizard model wraps form
- [ ] Implement `models.CVConfigWizardModel`
- [ ] Write failing test: Onboarding wizard model wraps form
- [ ] Implement `models.OnboardingWizardModel`
- [ ] Verify `forms.IsCompleted()` / `forms.IsAborted()` used (NOT huh state checks)

### Phase 3: WizardBehavior in behaviors/ (RED -> GREEN -> REFACTOR)
- [ ] Write failing test: WizardBehavior creation
- [ ] Implement `WizardBehavior` struct with `WizardForm` interface
- [ ] Write failing test: step navigation
- [ ] Implement step navigation
- [ ] Write failing test: completion/cancellation tracking
- [ ] Implement state tracking
- [ ] Refactor

### Phase 4: Migrate Wizard Modals (RED -> GREEN -> REFACTOR)
- [ ] Ensure existing tests pass before changes
- [ ] Refactor `cv_config_wizard_modal.go` to use models + behaviors
- [ ] **Verify NO huh import in components/**
- [ ] Refactor `onboarding_wizard_modal.go` to use models + behaviors
- [ ] **Verify NO huh import in components/**
- [ ] All existing tests still pass

---

## Pre-Commit Checklist (BEFORE EACH COMMIT)
- [ ] `make check-compliance` passes (REQUIRED before commit)
- [ ] `make check-patterns` passes - **NO huh imports outside forms/**
- [ ] Use `make ai-commit FILE=<path>` for AI-generated code
- [ ] Commit message explains **WHY**, not just WHAT
- [ ] Commit is atomic (ONE logical change)

## Post-Task Checklist (MUST COMPLETE BEFORE NEXT TASK)
- [ ] `make check-compliance` passes
- [ ] `grep -r "charmbracelet/huh" internal/cli/components/` returns **nothing**
- [ ] `grep -r "charmbracelet/huh" internal/cli/behaviors/` returns **nothing**
- [ ] All checkboxes above completed
- [ ] Task marked complete `[x]` in task file

---

## Acceptance Criteria
- [ ] **NO huh imports in components/ or behaviors/** (architecture compliance)
- [ ] Form builders exist in `forms/` package
- [ ] Form wrappers exist in `models/` package
- [ ] `WizardBehavior[T]` created with `WizardForm` interface
- [ ] CV config wizard uses new architecture
- [ ] Onboarding wizard uses new architecture
- [ ] All existing wizard tests pass
- [ ] New unit tests with >= 95% coverage
- [ ] `make check-patterns` passes

## Verification Steps
1. `grep -r "charmbracelet/huh" internal/cli/components/` - must return empty
2. `grep -r "charmbracelet/huh" internal/cli/behaviors/` - must return empty
3. `make check-compliance`
4. `make check-patterns`
5. Run all wizard tests
6. Manual verification: Run app and test both wizards

## Rollback Plan
- Revert commits if behavior breaks existing functionality
- Architecture fixes are critical - do not merge partial solutions

---

## Context

### Architecture Rules Reference
From AGENTS.md:
- `forms/` - **YES** can import huh (only place)
- `models/` - **NO** - use `forms.Form`, `forms.IsCompleted()`
- `components/` - **NO** - use `forms.NewXXX()` builders
- `behaviors/` - **NO** - use uikit, themes only

### Related Work
- **BUG-005**: Fixed form selection persistence (prerequisite - PR #99)
- **FORMS_GUIDE.md**: Form architecture patterns
- **FORMS_WORKFLOW_GUIDE.md**: Form workflow patterns
- **TableBehavior**: Reference for generic behavior pattern

