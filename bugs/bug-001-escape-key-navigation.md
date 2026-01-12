# Bug 001: Escape Key Navigation - Resolution & Enforcement

**Status**: ✅ Issues Resolved | 🚧 Enforcement In Progress  
**Severity**: 🟠 High  
**Created**: 2026-01-12  
**Updated**: 2026-01-12 (Consolidated with Bug 002, added enforcement plan)

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Part 1: Original Issues (Resolved)](#part-1-original-issues-resolved)
3. [Part 2: Enforcement Strategy](#part-2-enforcement-strategy)
4. [Part 3: Migration Progress](#part-3-migration-progress)
5. [Part 4: CI Enforcement](#part-4-ci-enforcement)
6. [Part 5: Verification Checklist](#part-5-verification-checklist)
7. [Appendix: Historical Context](#appendix-historical-context)

---

## Executive Summary

**Problem**: Users were unable to use the escape key to navigate back through workflows due to message delegation occurring before global key handling in forms and modals.

**Resolution**: Fixed 7 locations across intents and models where child components were receiving messages before the parent could check for global keys (esc, q, ?, m).

**Current Status**: All issues resolved, 88 escape tests passing (100% success rate), but pattern is duplicated ~53 times across 10 intents.

**Next Steps**: Migrate all intents to use the `MessageInterceptor` pattern to consolidate global key handling and prevent future regressions. Add CI enforcement to block PRs with incorrect patterns.

**Timeline**: 
- Issues Discovered: 2026-01-12
- Fixes Applied: 2026-01-12
- Enforcement Plan: 2026-01-12 (this document)
- Enforcement Completion: In Progress

---

## Part 1: Original Issues (Resolved)

### 1.1 Bug Summary

**Root Cause**: Message delegation happened BEFORE global key checking in multiple locations, allowing child components (forms, modals) to consume keyboard events before the parent intent could process them.

**Impact**:
- Users trapped in forms/modals with no way to navigate back
- Escape key had no effect
- Violated TUI standards for universal keyboard shortcuts
- Required Ctrl+C to force quit

**Affected Components** (7 issues fixed):

#### Intent Layer (4 issues)
1. **CaptureEvent** - `updateCaptureForm()` - Form delegation before HandleGlobalKeys
2. **CaptureEvent** - `updateReviewInferredEvent()` - 3 modal delegations before HandleGlobalKeys
3. **FactManagement** - `handleEditorState()` - Modal delegation before HandleGlobalKeys
4. **BurstManagement** - `updateEditView()` - Modal delegation before HandleGlobalKeys (Bug 002)

#### Model Layer (3 issues)
5. **metadata_editor_new** - Missing escape handling before huh form delegation
6. **fact_editor_new** - Missing escape handling before huh form delegation
7. **huh_capture_form** - NO escape handling at all (critical)

### 1.2 Root Cause Analysis

**The Pattern Violation**:

```go
// ❌ INCORRECT (allowed child to consume keys)
func (i *Intent) updateState(msg tea.Msg) tea.Cmd {
    // WRONG: Delegate to child FIRST
    _, cmd := i.form.Update(msg)  // Child consumes escape
    
    // Too late - key already consumed
    switch HandleGlobalKeys(msg) {
    case KeyBack:
        // Never reached
    }
    
    return cmd
}
```

**The Correct Pattern**:

```go
// ✅ CORRECT (global keys checked first)
func (i *Intent) updateState(msg tea.Msg) tea.Cmd {
    // RIGHT: Check global keys FIRST
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch HandleGlobalKeys(msg) {
        case KeyBack:
            i.state = PreviousState
            return nil
        case KeyQuit:
            return tea.Quit
        case KeyHelp:
            i.ToggleHelp()
            return nil
        }
    }
    
    // NOW delegate to child (after global keys handled)
    _, cmd := i.form.Update(msg)
    return cmd
}
```

### 1.3 Modal Delegation Issues (from Bug 002)

**Problem**: Same issue but in modal contexts - modals were receiving messages before global keys could be checked.

**Affected Modals**:
- CaptureEvent: metadataModal, burstModal, factModal
- BurstManagement: editModal

**Fix Pattern**: Add HandleGlobalKeys check BEFORE modal delegation, close modal on escape.

### 1.4 Resolution Summary

**Files Fixed**:
1. `internal/cli/intents/capture_event_intent.go` - 4 functions fixed
2. `internal/cli/intents/fact_management_intent.go` - 1 function fixed
3. `internal/cli/intents/burst_management_intent.go` - 1 function fixed
4. `internal/cli/models/metadata_editor_new.go` - Escape handling added
5. `internal/cli/models/fact_editor_new.go` - Escape handling added
6. `internal/cli/models/huh_capture_form.go` - Escape handling added

**Test Results**:
- ✅ 88/88 escape tests passing (100% success rate)
- ✅ 1,022/1,022 total tests passing
- ✅ 0 race conditions detected
- ✅ Code coverage maintained (>87%)
- ✅ Build successful

**Resolution Date**: 2026-01-12

---

## Part 2: Enforcement Strategy

### 2.1 Problem Statement

**Current Situation**:
- Issues are fixed, but pattern is duplicated ~53 times across codebase
- No enforcement mechanism to prevent future regressions
- Developer must manually remember to check global keys before delegation
- Easy to introduce the anti-pattern in new code

**Goal**: 
- Consolidate global key handling into reusable pattern
- Prevent future delegation-before-global-keys issues
- Add automated enforcement via CI

### 2.2 Reference Implementation: GenerateCV

**Why GenerateCV?**
- 10 states (most complex intent)
- 16 escape tests (comprehensive coverage)
- All tests passing
- Already follows correct pattern
- Will serve as gold standard for enforcement tests

**Test Coverage**:
```
GenerateCV Escape Tests (16 specs):
✅ SelectProfile (root state) - cancel on escape
✅ SelectAudience - go back to SelectProfile
✅ Generating (async) - go back to SelectAudience  
✅ Preview - go back to SelectAudience
✅ Review - go back to Preview
✅ Confirm - go back to Review
✅ ExportSelectFormat - go back to Confirm
✅ ExportSelectLocation - go back to ExportSelectFormat
✅ Exporting (async) - go back to ExportSelectLocation
✅ ExportComplete - go back to ExportSelectLocation
```

### 2.3 MessageInterceptor Pattern

**Existing Infrastructure**: The `MessageInterceptor` pattern already exists in `view_helpers.go` but is not currently used by any intents.

**Pattern Overview**:
```go
func (i *MyIntent) updateSomeState(msg tea.Msg) tea.Cmd {
    return NewMessageInterceptor().
        OnQuit(StandardQuitHandler()).
        OnHelp(StandardHelpHandler(i.BaseIntent)).
        OnBack(func() tea.Cmd {
            i.state = PreviousState
            return nil
        }).
        InterceptOr(msg, func() tea.Cmd {
            // Handle intent-specific keys and delegate
            if keyMsg, ok := msg.(tea.KeyMsg); ok {
                switch keyMsg.String() {
                case "enter":
                    // Intent-specific handling
                }
            }
            return i.childComponent.Update(msg)
        })
}
```

**Benefits**:
- ✅ Fluent, declarative API
- ✅ Global keys always checked first (enforced by pattern)
- ✅ Reusable across all intents
- ✅ Self-documenting code
- ✅ Testable pattern

### 2.4 New Helper Methods

We'll add convenience helpers to reduce boilerplate:

#### StandardQuitHandler
```go
// StandardQuitHandler returns tea.Quit
func StandardQuitHandler() GlobalKeyHandler {
    return func() tea.Cmd { return tea.Quit }
}
```

#### StandardHelpHandler
```go
// StandardHelpHandler creates a handler that toggles help on a BaseIntent
func StandardHelpHandler(intent *BaseIntent) GlobalKeyHandler {
    return func() tea.Cmd {
        intent.ToggleHelp()
        return nil
    }
}
```

#### OnContextAwareBack
```go
// OnContextAwareBack handles the common pattern where:
// - Edit mode: cancel intent and return to caller
// - New mode: go back to previous state
func (m *MessageInterceptor) OnContextAwareBack(
    isEditMode func() bool,
    goBack func() tea.Cmd,
    cancel func() tea.Cmd,
) *MessageInterceptor {
    return m.OnBack(func() tea.Cmd {
        if isEditMode() {
            return cancel()
        }
        return goBack()
    })
}
```

#### OnModalAwareBack
```go
// OnModalAwareBack handles the pattern where:
// - Modal active: close modal
// - No modal: go back to previous state
func (m *MessageInterceptor) OnModalAwareBack(
    hasActiveModal func() bool,
    closeModal func() tea.Cmd,
    goBack func() tea.Cmd,
) *MessageInterceptor {
    return m.OnBack(func() tea.Cmd {
        if hasActiveModal() {
            return closeModal()
        }
        return goBack()
    })
}
```

### 2.5 Migration Plan

**Approach**: Migrate intents one at a time, simplest first, GenerateCV last (it's the reference).

**Order** (by complexity):

| Order | Intent | States | Switches | Complexity | Notes |
|-------|--------|--------|----------|------------|-------|
| 1 | BrowseTimeline | 3 | 2 | ⭐ Simple | Linear flow, no context-aware |
| 2 | MetadataEditor | 3 | 3 | ⭐ Simple | Linear flow |
| 3 | ImportWizard | 4 | 3 | ⭐⭐ Medium | Simple async |
| 4 | BulkOperations | 4 | 3 | ⭐⭐ Medium | Simple async |
| 5 | ExportArtifact | 9 | 9 | ⭐⭐⭐ High | Many states |
| 6 | ConfigureSystem | 7 | 6 | ⭐⭐⭐ High | Editing sub-state |
| 7 | FactManagement | 6 | 4 | ⭐⭐⭐ High | Context-aware (IsNewFact) |
| 8 | BurstManagement | 8 | 9 | ⭐⭐⭐⭐ High | Context-aware (IsNewBurst) |
| 9 | CaptureEvent | 4 | 4 | ⭐⭐⭐⭐ High | Context-aware (PreviousEvent), modals |
| 10 | GenerateCV | 10 | 10 | ⭐⭐⭐⭐⭐ Reference | **LAST** - reference implementation |

**Total**: 53 switch statements to migrate

**Per-Intent Process**:
1. Convert first switch statement to MessageInterceptor
2. Run tests - verify intent still works
3. Convert remaining switches in that intent
4. Run full test suite - verify no regressions
5. Commit changes
6. Move to next intent

---

## Part 3: Migration Progress

### 3.1 Intent Migration Status

| Intent | States | Switches | Status | Tests | Notes |
|--------|--------|----------|--------|-------|-------|
| GenerateCV | 10 | 10 | ✅ Reference | 16 | Gold standard - verify enforcement tests pass |
| BrowseTimeline | 3 | 2 | ⏳ Pending | 5 | First to migrate |
| MetadataEditor | 3 | 3 | ⏳ Pending | 0 | Simple linear flow |
| ImportWizard | 4 | 3 | ⏳ Pending | 0 | Simple async |
| BulkOperations | 4 | 3 | ⏳ Pending | 0 | Simple async |
| ExportArtifact | 9 | 9 | ⏳ Pending | 7 | Many states |
| ConfigureSystem | 7 | 6 | ⏳ Pending | 7 | Editing sub-state |
| FactManagement | 6 | 4 | ⏳ Pending | 0 | Context-aware (IsNewFact) |
| BurstManagement | 8 | 9 | ⏳ Pending | 11 | Context-aware (IsNewBurst) |
| CaptureEvent | 4 | 4 | ⏳ Pending | 9 | Context-aware (PreviousEvent) |

**Overall Progress**: 0/53 switch statements migrated (0%)

### 3.2 Test Coverage Matrix

All intents must pass enforcement tests for:
- Root state escape (cancel intent)
- Intermediate state escape (go back)
- Quit key (q) from all states
- Help key (?) from all states

| Intent | Root | Intermediate | Async | Modal | Quit | Help | Total |
|--------|------|--------------|-------|-------|------|------|-------|
| GenerateCV | ✅ | ✅ | ✅ | N/A | ✅ | ✅ | 16 |
| BrowseTimeline | ⏳ | ⏳ | N/A | N/A | ⏳ | ⏳ | 5 |
| MetadataEditor | ⏳ | ⏳ | N/A | N/A | ⏳ | ⏳ | TBD |
| ImportWizard | ⏳ | ⏳ | ⏳ | N/A | ⏳ | ⏳ | TBD |
| BulkOperations | ⏳ | ⏳ | ⏳ | N/A | ⏳ | ⏳ | TBD |
| ExportArtifact | ⏳ | ⏳ | ⏳ | N/A | ⏳ | ⏳ | 7 |
| ConfigureSystem | ⏳ | ⏳ | ⏳ | N/A | ⏳ | ⏳ | 7 |
| FactManagement | ⏳ | ⏳ | N/A | ⏳ | ⏳ | ⏳ | TBD |
| BurstManagement | ⏳ | ⏳ | N/A | ⏳ | ⏳ | ⏳ | 11 |
| CaptureEvent | ⏳ | ⏳ | N/A | ⏳ | ⏳ | ⏳ | 9 |

**Target**: ~200 enforcement test specs across all intents

---

## Part 4: CI Enforcement

### 4.1 Enforcement Test Suite

**File**: `internal/cli/intents/global_keys_enforcement_e2e_test.go`

**Purpose**: Exhaustive E2E tests that verify all intents handle global keys correctly.

**Structure**:
```go
var _ = Describe("Global Keys Enforcement", func() {
    Describe("Reference: GenerateCV Intent", func() {
        // Comprehensive tests for all 10 GenerateCV states
        // These MUST pass - GenerateCV is the gold standard
    })
    
    Describe("All Intents Contract Compliance", func() {
        DescribeTable("root state escape should cancel",
            func(name string, setup func() (Intent, func())) {
                // Test all 10 intents
            },
            Entry("GenerateCV", ...),
            Entry("BrowseTimeline", ...),
            // ... all intents
        )
        
        DescribeTable("intermediate states escape should go back", ...)
        DescribeTable("all states should handle quit key", ...)
        DescribeTable("all states should handle help key", ...)
    })
    
    Describe("Context-Aware Navigation", func() {
        // CaptureEvent (PreviousEvent), BurstManagement (IsNewBurst), FactManagement (IsNewFact)
    })
})
```

**Test Count**: ~200 specs covering all states in all intents

### 4.2 PR Validation Job

**File**: `.github/workflows/pr-validation.yml`

Add new job:

```yaml
# Check that intent files use MessageInterceptor pattern
check-message-interceptor:
  name: Check MessageInterceptor Pattern
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    
    - name: Verify MessageInterceptor usage in intents
      run: |
        echo "Checking intent files for MessageInterceptor pattern..."
        
        INTENT_FILES=$(find internal/cli/intents -name "*_intent.go" -not -name "*_test.go")
        INTENT_FILES="$INTENT_FILES $(find internal/cli/intents -name "export_artifact.go" -o -name "configure_system.go")"
        
        VIOLATIONS=()
        
        for file in $INTENT_FILES; do
          # Skip view_helpers.go (contains the pattern definition)
          [[ "$file" == *"view_helpers.go"* ]] && continue
          
          # Check if file has Update methods
          if grep -q "func.*Update.*tea.Msg" "$file"; then
            # Must use MessageInterceptor (not raw HandleGlobalKeys switch)
            INTERCEPTOR_COUNT=$(grep -c "NewMessageInterceptor\|MessageInterceptor()" "$file" 2>/dev/null || echo "0")
            RAW_SWITCH_COUNT=$(grep -c "switch HandleGlobalKeys" "$file" 2>/dev/null || echo "0")
            
            if [ "$RAW_SWITCH_COUNT" -gt 0 ]; then
              VIOLATIONS+=("$file: Found $RAW_SWITCH_COUNT raw HandleGlobalKeys switches (should use MessageInterceptor)")
            fi
            
            if [ "$INTERCEPTOR_COUNT" -eq 0 ]; then
              VIOLATIONS+=("$file: No MessageInterceptor usage found")
            fi
          fi
        done
        
        if [ ${#VIOLATIONS[@]} -gt 0 ]; then
          echo ""
          echo "❌ MessageInterceptor Pattern Violations:"
          printf '  - %s\n' "${VIOLATIONS[@]}"
          echo ""
          echo "All intents MUST use MessageInterceptor for global key handling."
          echo "See: docs/TUI_DEVELOPER_GUIDE.md#message-interceptor-pattern"
          exit 1
        fi
        
        echo "✅ All intent files use MessageInterceptor pattern correctly"
```

**Behavior**: Blocking - fails PR if violations found

### 4.3 Local Verification

**File**: `scripts/check-message-interceptor.sh`

Same logic as CI job, for local use before pushing:

```bash
#!/bin/bash
# Check that all intent files use MessageInterceptor pattern

echo "Checking intent files for MessageInterceptor pattern..."

INTENT_FILES=$(find internal/cli/intents -name "*_intent.go" -not -name "*_test.go")
INTENT_FILES="$INTENT_FILES $(find internal/cli/intents -name "export_artifact.go" -o -name "configure_system.go")"

VIOLATIONS=()

for file in $INTENT_FILES; do
  [[ "$file" == *"view_helpers.go"* ]] && continue
  
  if grep -q "func.*Update.*tea.Msg" "$file"; then
    INTERCEPTOR_COUNT=$(grep -c "NewMessageInterceptor\|MessageInterceptor()" "$file" 2>/dev/null || echo "0")
    RAW_SWITCH_COUNT=$(grep -c "switch HandleGlobalKeys" "$file" 2>/dev/null || echo "0")
    
    if [ "$RAW_SWITCH_COUNT" -gt 0 ]; then
      VIOLATIONS+=("$file: $RAW_SWITCH_COUNT raw HandleGlobalKeys switches")
    fi
  fi
done

if [ ${#VIOLATIONS[@]} -gt 0 ]; then
  echo "❌ Violations found:"
  printf '  - %s\n' "${VIOLATIONS[@]}"
  exit 1
fi

echo "✅ All files compliant"
```

**Makefile Target**:
```makefile
# Check MessageInterceptor pattern usage
check-interceptor:
	@echo "Checking MessageInterceptor pattern..."
	@./scripts/check-message-interceptor.sh

# Add to check-compliance
check-compliance: check-interceptor ...existing targets...
```

---

## Part 5: Verification Checklist

### Infrastructure Tasks
- [ ] Add `OnContextAwareBack()` to `view_helpers.go`
- [ ] Add `OnModalAwareBack()` to `view_helpers.go`
- [ ] Add `StandardQuitHandler()` to `view_helpers.go`
- [ ] Add `StandardHelpHandler()` to `view_helpers.go`
- [ ] Add tests for new helpers to `message_interceptor_test.go`

### Test Suite Tasks
- [ ] Create `global_keys_enforcement_e2e_test.go`
- [ ] Add GenerateCV reference tests (all 10 states)
- [ ] Verify GenerateCV tests pass (baseline)
- [ ] Add root state escape tests (all 10 intents)
- [ ] Add intermediate state escape tests (~35 states)
- [ ] Add quit key tests (all ~58 states)
- [ ] Add help key tests (all ~58 states)
- [ ] Add context-aware navigation tests (3 intents)

### Intent Migration Tasks
- [ ] Migrate BrowseTimeline (2 switches) → verify tests
- [ ] Migrate MetadataEditor (3 switches) → verify tests
- [ ] Migrate ImportWizard (3 switches) → verify tests
- [ ] Migrate BulkOperations (3 switches) → verify tests
- [ ] Migrate ExportArtifact (9 switches) → verify tests
- [ ] Migrate ConfigureSystem (6 switches) → verify tests
- [ ] Migrate FactManagement (4 switches) → verify tests
- [ ] Migrate BurstManagement (9 switches) → verify tests
- [ ] Migrate CaptureEvent (4 switches) → verify tests
- [ ] Migrate GenerateCV (10 switches) → verify tests

### CI Enforcement Tasks
- [ ] Add `check-message-interceptor` job to `pr-validation.yml`
- [ ] Create `scripts/check-message-interceptor.sh`
- [ ] Make script executable (`chmod +x`)
- [ ] Add `check-interceptor` target to `Makefile`
- [ ] Verify local enforcement works (`make check-interceptor`)
- [ ] Verify CI enforcement works (push to PR)

### Documentation Tasks
- [ ] Update `docs/TUI_DEVELOPER_GUIDE.md` with MessageInterceptor section
- [ ] Update `docs/INTENT_DEVELOPMENT_CHECKLIST.md` with global key requirements
- [ ] Add examples to documentation
- [ ] Delete `bugs/bug-002-escape-key-modal-delegation.md`

### Final Verification Tasks
- [ ] Run `make test` - all tests pass
- [ ] Run `make test-race` - no race conditions
- [ ] Run `make check-compliance` - passes
- [ ] Run `make check-interceptor` - passes
- [ ] Run `make ci-local` - all CI checks pass
- [ ] Manual smoke test of all workflows

---

## Appendix: Historical Context

### Investigation Timeline

**2026-01-12 02:00** - Initial investigation
- User reported escape key not working
- Found 88 existing escape tests passing
- Confusion: tests pass but user reports issues

**2026-01-12 02:15** - Code review
- Discovered message delegation order issue
- Form.Update(msg) called BEFORE HandleGlobalKeys(msg)
- Huh library consuming escape before parent could check

**2026-01-12 02:30** - Test analysis
- Unit tests bypassed the delegation path
- E2E tests would have caught the issue
- Added E2E tests to reproduce bug

**2026-01-12 03:30** - Initial audit (CaptureEvent)
- Found 4 issues in CaptureEvent intent
- Applied fixes, all tests passing

**2026-01-12 13:00** - Comprehensive intent audit
- Audited all 10 intents systematically
- Found 1 additional issue in FactManagement
- Verified 9 intents already correct

**2026-01-12 13:30** - Application-wide audit
- Audited all 27 files with Update() methods
- Found 3 issues in model layer
- 100% application compliance achieved

**2026-01-12 14:00** - Bug 002 (Modal delegation)
- Identified modal-specific delegation issues
- Fixed 4 modal delegations
- All modal escape tests passing

**2026-01-12 15:00** - Enforcement planning
- Recognized pattern duplication (53 switches)
- Designed MessageInterceptor migration strategy
- Created this consolidated document

### Key Lessons Learned

1. **E2E Tests are Essential**: Unit tests alone missed this issue because they bypassed the actual message routing
2. **Message Order Matters**: Global keys must ALWAYS be checked before delegation
3. **Pattern Duplication is Risky**: Same mistake repeated 53 times makes enforcement critical
4. **Context-Awareness is Key**: Don't hard-code state transitions, check context (new vs edit)
5. **Documentation Gaps**: Standards claimed "100% complete" but lacked ordering requirements

### Related Documentation

- **TUI Standards**: `docs/TUI_STANDARDS.md` (lines 87-217)
- **Navigation Testing**: `docs/development/NAVIGATION_TESTING_GUIDE.md`
- **Keyboard Guide**: `docs/KEYBOARD_SHORTCUTS_GUIDE.md`
- **Developer Guide**: `docs/TUI_DEVELOPER_GUIDE.md` (to be updated)

### Test Files

**Existing Escape Tests** (88 total, 55 specs):
- `capture_event_escape_test.go` - 9 specs
- `browse_timeline_escape_test.go` - 5 specs
- `burst_management_modal_escape_test.go` - 11 specs
- `configure_system_escape_test.go` - 7 specs
- `export_artifact_escape_test.go` - 7 specs
- `generate_cv_escape_test.go` - 16 specs

**New Enforcement Tests** (to be created):
- `global_keys_enforcement_e2e_test.go` - ~200 specs

### Commits

**Bug Fixes** (2026-01-12):
1. `fix(intents): context-aware escape navigation in CaptureEvent`
2. `test(intents): add context-aware escape tests for CaptureEvent`
3. `fix(intents): escape key handling in FactManagement modal`
4. `fix(models): escape key handling in huh form models`
5. `docs(bugs): update bug-001 with complete application audit`

**Enforcement** (In Progress):
- Will be added as work progresses

---

**Last Updated**: 2026-01-12  
**Status**: Issues Resolved ✅ | Enforcement In Progress 🚧  
**Next Steps**: Begin Phase 2 (Infrastructure) - Add helper methods to view_helpers.go
