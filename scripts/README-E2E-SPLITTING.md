# E2E Test Splitting Scripts

## Overview

These scripts automate the splitting of MIXED E2E workflow test files into:
1. **E2E tests** (SQLite) - stays in `internal/testutil/e2e/`
2. **Navigation tests** (Memory) - moves to `internal/cli/intents/`

## Files

| Script | Purpose |
|--------|---------|
| `split-e2e-tests.py` | Python script to split a single workflow test file |
| `split-e2e-tests.sh` | Bash wrapper (deprecated - use Python version) |
| `split-all-e2e-tests.sh` | Batch process all remaining workflow test files |

## Usage

### Option 1: Split One File at a Time

```bash
# Split a specific file
python3 scripts/split-e2e-tests.py browse_workflow_test.go

# Review the output
cat internal/testutil/e2e/browse_e2e_test.go
cat internal/cli/intents/browse_navigation_test.go

# Delete original and commit
rm internal/testutil/e2e/browse_workflow_test.go
git add -A
make ai-commit MSG="refactor(tests): split browse_workflow into e2e and navigation tests"
```

### Option 2: Batch Process All Files

```bash
# Process all remaining *_workflow_test.go files
./scripts/split-all-e2e-tests.sh
```

This will:
1. Find all `*_workflow_test.go` files in `internal/testutil/e2e/`
2. Split each file into E2E and Navigation tests
3. Delete the original file
4. Create a commit for each split

## How It Works

The Python script (`split-e2e-tests.py`):

1. **Reads the source file** - Parses the workflow test file
2. **Extracts Describe blocks**:
   - **E2E tests**: Blocks containing `e2e.Setup(GinkgoT())`
   - **Navigation tests**: Blocks containing `e2e.SetupWithMemory(GinkgoT())`
3. **Generates two files**:
   - `<base>_e2e_test.go` - E2E tests (package: `e2e_test`)
   - `<base>_navigation_test.go` - Navigation tests (package: `intents_test`)

### Detection Logic

**E2E Tests** (uses SQLite):
- Uses `e2e.Setup(GinkgoT())`
- Calls `PopulateTestData()`, `AddEvent()`, `AddBurst()`, `AddFact()`
- Asserts with `AssertEventCount()`, `AssertBurstCount()`, `AssertFactCount()`
- Tests `SimulateRestart()` for persistence verification

**Navigation Tests** (uses memory):
- Uses `e2e.SetupWithMemory(GinkgoT())`
- Only tests key presses, screen navigation, view rendering
- No data persistence assertions

## Remaining Files

As of last check, the following files need splitting:

```bash
browse_workflow_test.go         (342 lines)
burst_management_workflow_test.go (348 lines)
capture_workflow_test.go        (345 lines)
configure_workflow_test.go      (310 lines)
error_recovery_test.go          (378 lines - not a workflow file but needs splitting)
export_workflow_test.go         (378 lines)
generate_cv_workflow_test.go    (325 lines)
```

## Expected Output

After splitting all files, you should have:

### E2E Directory (`internal/testutil/e2e/`)
```
e2e/
├── browse_e2e_test.go
├── burst_management_e2e_test.go
├── capture_e2e_test.go
├── chained_workflows_e2e_test.go (already renamed)
├── configure_e2e_test.go
├── error_recovery_e2e_test.go
├── export_e2e_test.go
├── fact_management_e2e_test.go (already done)
├── generate_cv_e2e_test.go
├── e2e_suite_test.go
├── helpers.go
├── helpers_test.go
└── fixtures.go
```

### Intents Directory (`internal/cli/intents/`)
```
intents/
├── ... (existing test files)
├── browse_timeline_navigation_test.go
├── bulk_operations_navigation_test.go (already done)
├── burst_management_navigation_test.go
├── capture_event_navigation_test.go
├── configure_system_navigation_test.go
├── error_recovery_navigation_test.go
├── export_artifact_navigation_test.go
├── fact_management_navigation_test.go (already done)
├── generate_cv_navigation_test.go
├── import_wizard_navigation_test.go (already done)
└── metadata_editor_navigation_test.go (already done)
```

## Verification

After splitting, verify the tests still pass:

```bash
# Run E2E tests
ginkgo internal/testutil/e2e/

# Run navigation tests
ginkgo --focus="Navigation" internal/cli/intents/

# Run all tests
make test
```

## Troubleshooting

### Script fails to extract blocks

**Problem**: The Python script couldn't find Describe blocks with the expected setup calls.

**Solution**: Check that the file actually contains both `e2e.Setup(GinkgoT())` and `e2e.SetupWithMemory(GinkgoT())` calls. If it only has one type, manually handle it.

### Generated files have syntax errors

**Problem**: The extracted blocks don't compile.

**Solution**: Review the generated files and manually fix indentation or missing imports. The script does basic indentation but may need manual adjustment.

### Tests fail after splitting

**Problem**: Tests that passed before now fail.

**Solution**: 
1. Check that the correct blocks were extracted (E2E vs Navigation)
2. Verify all `BeforeEach` and `AfterEach` blocks are included
3. Ensure imports are correct (`tea` package for navigation tests)

## Next Steps

After all files are split:

1. **Update documentation**:
   - `AGENTS.md` - Add E2E test structure section
   - `docs/E2E_IMPLEMENTATION_GAPS.md` - Update file references
   - `docs/rules/go-guidelines.md` - Add E2E naming convention

2. **Create `docs/TESTING_GUIDE.md`**:
   - Comprehensive guide covering unit, integration, navigation, and E2E tests

3. **Final commit**:
   ```bash
   make ai-commit MSG="docs: update documentation for E2E test restructuring"
   ```

## Status Tracking

| File | Status | Commit |
|------|--------|--------|
| bulk_operations | ✅ Moved | `620bfbe` |
| import_wizard | ✅ Moved | `620bfbe` |
| metadata_editor | ✅ Moved | `620bfbe` |
| chained_workflows | ✅ Renamed | `eeb00b1` |
| fact_management | ✅ Split | `7aef762` |
| browse | ⏳ Pending | |
| burst_management | ⏳ Pending | |
| capture | ⏳ Pending | |
| configure | ⏳ Pending | |
| error_recovery | ⏳ Pending | |
| export | ⏳ Pending | |
| generate_cv | ⏳ Pending | |
