# E2E Test Location Guide

**Centralized location and consolidation rules for E2E tests**

## Critical Rule: One Location, One File Per Workflow

### ✅ CORRECT Location
```
internal/testutil/e2e/{workflow}_e2e_test.go
```

**All E2E tests MUST live in `internal/testutil/e2e/` directory.**

### ❌ WRONG Locations
```
❌ internal/cli/intents/{feature}/{feature}_e2e_test.go
❌ internal/cli/intents/{feature}_e2e_test.go
❌ internal/cli/uikit/{component}/{component}_e2e_test.go
❌ internal/domain/{entity}/{entity}_e2e_test.go
```

**Rationale:**
- E2E tests verify **complete user workflows**, not individual components
- Centralizing in `internal/testutil/e2e/` makes them easy to find and run
- Prevents duplication and fragmentation
- Clear separation: unit tests live with components, E2E tests are centralized

## Naming Convention

### Pattern
```
{workflow_name}_e2e_test.go
```

### Examples
| Workflow | E2E Test File |
|----------|---------------|
| Capture Event | `capture_event_e2e_test.go` |
| Browse Timeline | `browse_timeline_e2e_test.go` |
| Burst Management | `burst_management_e2e_test.go` |
| Generate CV | `generate_cv_e2e_test.go` |
| Manage Skills | `skills_management_e2e_test.go` |
| Manage Facts | `fact_management_e2e_test.go` |
| Onboarding | `onboarding_e2e_test.go` |
| Configure | `configure_e2e_test.go` |

### One Comprehensive File Per Workflow

**Each workflow should have ONE E2E test file** containing all scenarios for that workflow:

```go
// ✅ GOOD: internal/testutil/e2e/capture_event_e2e_test.go
var _ = Describe("E2E Capture Event Workflow", func() {
    Describe("Happy Path", func() {
        It("should capture event with manual enrichment", func() { ... })
        It("should capture event with auto enrichment", func() { ... })
    })
    
    Describe("Metadata Editing", func() {
        It("should allow editing metadata before saving", func() { ... })
        It("should validate metadata fields", func() { ... })
    })
    
    Describe("Navigation", func() {
        It("should navigate back to menu", func() { ... })
        It("should handle escape key correctly", func() { ... })
    })
    
    Describe("Edge Cases", func() {
        It("should handle empty input", func() { ... })
        It("should recover from errors", func() { ... })
    })
})

// ❌ BAD: Multiple files for same workflow
// capture_e2e_test.go
// capture_workflow_e2e_test.go
// capture_navigation_e2e_test.go
```

## Current State Analysis

### Files in Wrong Locations

| File | Current Location | Action Required |
|------|------------------|-----------------|
| `fact_management_e2e_test.go` | `internal/cli/intents/factmanagement/` | **MOVE** to `internal/testutil/e2e/` |
| `generate_cv_wizard_e2e_test.go` | `internal/cli/intents/` | **CONSOLIDATE** with testutil version |
| `escape_bug_e2e_test.go` | `internal/cli/intents/` | **MOVE** or **CONSOLIDATE** into workflow files |
| `global_keys_enforcement_e2e_test.go` | `internal/cli/intents/` | **MOVE** or **CONSOLIDATE** |
| `info_modal_e2e_test.go` | `internal/cli/uikit/feedback/` | **MOVE** or convert to integration test |
| `selectors_e2e_test.go` | `internal/cli/uikit/selectors/` | **MOVE** or convert to integration test |

### Duplicate Files (Same Workflow)

#### Fact Management
- ✅ `internal/testutil/e2e/fact_management_e2e_test.go` (KEEP)
- ❌ `internal/cli/intents/factmanagement/fact_management_e2e_test.go` (REMOVE after merging)

**Action:** Consolidate into `internal/testutil/e2e/fact_management_e2e_test.go`

#### Generate CV Wizard
- ✅ `internal/testutil/e2e/generate_cv_wizard_e2e_test.go` (KEEP)
- ❌ `internal/cli/intents/generate_cv_wizard_e2e_test.go` (REMOVE after merging)

**Action:** Consolidate into `internal/testutil/e2e/generate_cv_wizard_e2e_test.go`

### Multiple Files for Same Workflow (Consolidation Needed)

#### Capture Event (3 files)
- `capture_e2e_test.go` (134 lines)
- `capture_workflow_e2e_test.go` (246 lines)
- `capture_navigation_e2e_test.go` (289 lines)

**Action:** Consolidate into ONE file: `capture_event_e2e_test.go`

#### Generate CV (4 files)
- `generate_cv_e2e_test.go` (55 lines)
- `generate_cv_wizard_e2e_test.go` (506 lines)
- `generate_cv_baseline_e2e_test.go` (373 lines)
- `generate_cv_empty_state_e2e_test.go` (145 lines)

**Action:** Consolidate into ONE file: `generate_cv_e2e_test.go`

#### Burst Management (2 files)
- `burst_management_e2e_test.go` (668 lines)
- `burst_management_intent_e2e_test.go` (3931 lines - VERY LARGE)

**Action:** 
- **Option 1:** Consolidate into `burst_management_e2e_test.go` (if combined size is reasonable)
- **Option 2:** Split by major feature areas if too large:
  - `burst_detection_e2e_test.go` (burst detection workflow)
  - `burst_management_e2e_test.go` (CRUD operations workflow)

## Consolidation Workflow

### Step 1: Identify Duplicates
```bash
# Find all E2E test files
find . -name "*_e2e_test.go" -type f | sort

# Check for duplicates by workflow
ls internal/testutil/e2e/*_e2e_test.go
ls internal/cli/intents/**/*_e2e_test.go
```

### Step 2: Compare and Merge
```bash
# Compare duplicate files
diff internal/testutil/e2e/fact_management_e2e_test.go \
     internal/cli/intents/factmanagement/fact_management_e2e_test.go

# Check for unique test cases in each file
grep "It(" internal/testutil/e2e/fact_management_e2e_test.go
grep "It(" internal/cli/intents/factmanagement/fact_management_e2e_test.go
```

### Step 3: Consolidate Tests

**If both files have unique tests:**
```go
// Merge into internal/testutil/e2e/fact_management_e2e_test.go
var _ = Describe("E2E Fact Management Workflow", func() {
    // Tests from testutil/e2e version
    Describe("Empty Fact List", func() {
        It("should show empty state", func() { ... })
    })
    
    // Tests from intents/factmanagement version (if unique)
    Describe("Fact CRUD Operations", func() {
        It("should create new fact", func() { ... })
    })
})
```

**If one file is superset of other:**
```bash
# Keep the more comprehensive version
# Delete the duplicate
rm internal/cli/intents/factmanagement/fact_management_e2e_test.go
```

### Step 4: Move Misplaced Files
```bash
# Move file to correct location
git mv internal/cli/intents/factmanagement/fact_management_e2e_test.go \
       internal/testutil/e2e/fact_management_e2e_test.go

# Update import paths in moved file (if needed)
# Run tests to verify
make test-suite SUITE=./internal/testutil/e2e/...
```

### Step 5: Update References
```bash
# Search for references to old file location
grep -r "internal/cli/intents/factmanagement" .

# Update documentation
# Update CI configuration if needed
```

### Step 6: Run Full Test Suite
```bash
# Verify all E2E tests pass
make test-suite SUITE=./internal/testutil/e2e/...

# Run full compliance check
make check-compliance
```

## Special Cases

### Cross-Workflow Tests

For tests that verify interactions between multiple workflows:

```
internal/testutil/e2e/chained_workflows_e2e_test.go
```

**Example:**
- Capture event → Browse timeline → Generate CV
- Import CSV → Browse timeline → Burst detection

### Component E2E Tests in Wrong Location

Components in `internal/cli/uikit/` should NOT have E2E tests:

```
❌ internal/cli/uikit/feedback/info_modal_e2e_test.go
❌ internal/cli/uikit/selectors/selectors_e2e_test.go
```

**Action:**
1. **If testing component in isolation:** Convert to **integration test** (remove `_e2e` suffix)
2. **If testing component in workflow:** Move scenarios into appropriate **workflow E2E test**

Example:
```go
// ❌ BAD: internal/cli/uikit/feedback/info_modal_e2e_test.go
It("should show info modal", func() {
    modal.Show()
    // ...
})

// ✅ GOOD: internal/testutil/e2e/capture_event_e2e_test.go
It("should show help modal during capture", func() {
    env.SelectIntentByName("capture_event")
    env.PressKeyRune('?')  // User presses help key
    env.AssertViewContains("Help")
    // ...
})
```

### Global Behavior Tests

Tests for global behaviors (keyboard shortcuts, escape handling) should be:

**Option 1:** Consolidated into workflow files
```go
// internal/testutil/e2e/capture_event_e2e_test.go
Describe("Global Keys", func() {
    It("should handle escape key correctly", func() { ... })
    It("should show help with '?'", func() { ... })
})
```

**Option 2:** Separate file if testing across ALL workflows
```
internal/testutil/e2e/global_behavior_e2e_test.go
```

## Pre-Commit Checklist

Before committing new E2E tests, verify:

- [ ] E2E test is in `internal/testutil/e2e/` directory
- [ ] No duplicate E2E test files for same workflow
- [ ] File follows naming pattern: `{workflow}_e2e_test.go`
- [ ] No E2E tests in `internal/cli/intents/` or `internal/cli/uikit/`
- [ ] All E2E tests pass: `make test-suite SUITE=./internal/testutil/e2e/...`
- [ ] Tests use self-documenting helpers from `e2e.TestEnv`

## Automated Check (Future)

Add to `make check-compliance`:

```bash
# Check for E2E tests in wrong locations
find internal/cli/intents -name "*_e2e_test.go" -type f
find internal/cli/uikit -name "*_e2e_test.go" -type f
find internal/domain -name "*_e2e_test.go" -type f

# If any found, fail with error message
```

## Summary

**Golden Rules:**
1. ✅ **ONE location:** `internal/testutil/e2e/`
2. ✅ **ONE file per workflow:** `{workflow}_e2e_test.go`
3. ✅ **Complete workflows:** Test full user journeys, not components
4. ✅ **Self-documenting helpers:** Use `env.SelectIntentByName()`, etc.
5. ❌ **NO E2E tests in component directories**
6. ❌ **NO duplicate files for same workflow**

**When adding new E2E tests:**
1. Check if file exists: `ls internal/testutil/e2e/{workflow}_e2e_test.go`
2. If exists: Add tests to existing file
3. If not: Create new file in correct location
4. Use self-documenting helpers
5. Test complete user journey

**When finding misplaced E2E tests:**
1. **FIRST:** Move to correct location
2. **SECOND:** Consolidate with existing file (if duplicate)
3. **THIRD:** Implement the feature/fix
