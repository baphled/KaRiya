# Task 11: Burst Management Table Enhancement

## Overview
- **Goal**: Improve the burst management table layout to display more useful information (Description, Confirmed status, Created date) with color-coded competencies
- **Time Estimate**: 75 minutes
- **Prerequisites**: Understanding of BubbleTea table component, Lipgloss styling, intent architecture
- **Status**: In Progress

## Motivation
Users need to efficiently manage bursts with better visibility into burst details without navigating to detail view. Adding Description, Confirmed status, and Created date columns will improve productivity.

## Files to Modify
- [x] `tasks/tasks-11-burst-table-enhancement.md` (this file)
- [ ] `internal/cli/intents/burst_management_intent.go` (primary implementation)
- [ ] `internal/cli/intents/burst_management_test.go` (test updates)

## Implementation Checklist

### Phase 1: Preparation (5 min)
- [x] Create task file
- [ ] Run compliance check (baseline)
- [ ] Review current table implementation
- [ ] Verify test suite passes

### Phase 2: Update Column Definitions (5 min)
- [ ] Update table column definitions from 3 to 6 columns
- [ ] Adjust column widths: Name(22), Description(25), Confirmed(8), Competency(15), Events(8), Created(12)
- [ ] Commit: `feat(burst): update burst table column definitions`

### Phase 3: Add Helper Functions (15 min)
- [ ] Add `formatConfirmedStatus(confirmed bool) string` - Returns styled "✓ Yes" or "✗ No"
- [ ] Add `formatCompetency(competency string) string` - Returns color-coded competency text
- [ ] Add `formatDescription(description string) string` - Returns truncated description (max 25 chars)
- [ ] Add `formatCreatedDate(createdAt time.Time) string` - Returns YYYY-MM-DD formatted date
- [ ] Add necessary imports (strings, time)
- [ ] Commit: `feat(burst): add formatting helpers for table columns`

### Phase 4: Update Row Generation Logic (10 min)
- [ ] Update `updateTableRows()` method to use new 6-column structure
- [ ] Integrate all 4 helper functions
- [ ] Update name truncation to 19 chars (to fit 22 with indicator)
- [ ] Commit: `feat(burst): implement 6-column table row generation`

### Phase 5: Update Existing Tests (15 min)
- [ ] Find tests checking column count (update 3 → 6)
- [ ] Find tests checking column headers (add new headers)
- [ ] Update row rendering assertions for 6 columns
- [ ] Verify all existing tests pass
- [ ] Commit: `test(burst): update tests for 6-column table structure`

### Phase 6: Add New Helper Tests (20 min)
- [ ] Add test suite for `formatConfirmedStatus()`
  - [ ] Test confirmed=true returns "✓" and "Yes"
  - [ ] Test confirmed=false returns "✗" and "No"
  - [ ] Verify colors (green for Yes, gray for No)
- [ ] Add test suite for `formatCompetency()`
  - [ ] Test all 6 competency categories have colors
  - [ ] Test empty competency returns "-"
  - [ ] Test unknown competency uses default color
- [ ] Add test suite for `formatDescription()`
  - [ ] Test long descriptions truncate at 22 chars + "..."
  - [ ] Test short descriptions render fully
  - [ ] Test empty descriptions return "-"
  - [ ] Test newlines are removed
- [ ] Add test suite for `formatCreatedDate()`
  - [ ] Test date formats as YYYY-MM-DD
  - [ ] Test various dates render correctly
- [ ] Commit: `test(burst): add helper function tests`

### Phase 7: Integration Testing (10 min)
- [ ] Run full test suite: `go test ./internal/cli/intents/...`
- [ ] Run with race detector: `go test -race ./internal/cli/intents/...`
- [ ] Fix any failing tests
- [ ] Verify 100% pass rate
- [ ] Commit if fixes needed: `fix(burst): resolve test failures`

### Phase 8: Manual Testing (10 min)
- [ ] Build application: `go build -o kariya ./cmd/kariya`
- [ ] Run and navigate to burst management
- [ ] Verify all 6 columns display
- [ ] Check confirmed status shows ✓/✗ with colors
- [ ] Verify competency colors render correctly
- [ ] Test with long/short/empty descriptions
- [ ] Verify date formatting
- [ ] Test navigation and pagination
- [ ] Document any issues

### Phase 9: Final Verification (5 min)
- [ ] Run compliance check: `make check-compliance`
- [ ] Run linter: `golangci-lint run ./internal/cli/intents/...`
- [ ] Format code: `go fmt ./internal/cli/intents/...`
- [ ] Review all commits for atomicity
- [ ] Verify commit messages follow conventional format

## Testing Instructions

### Automated Tests
```bash
# Run all intent tests
go test -v ./internal/cli/intents/

# Run with coverage
go test -cover ./internal/cli/intents/

# Run with race detector
go test -race ./internal/cli/intents/

# Run specific test
go test -v ./internal/cli/intents/ -run TestBurstManagement
```

### Manual Tests
1. Build: `go build -o kariya ./cmd/kariya`
2. Run: `./kariya`
3. Navigate to "Manage Bursts" (m key)
4. Verify table layout:
   - 6 columns visible
   - Headers: Name, Description, Confirmed, Competency, Events, Created
   - Confirmed shows ✓ Yes (green) or ✗ No (gray)
   - Competencies are color-coded
   - Descriptions are truncated
   - Dates show as YYYY-MM-DD
5. Test edge cases:
   - Empty descriptions
   - Empty competencies
   - Long names
   - Various confirmed states

## Acceptance Criteria
- [ ] Table displays 6 columns with correct headers
- [ ] Name column: 22 chars, truncated with "...", focus indicator "▶"
- [ ] Description column: 25 chars, truncated with "...", shows "-" if empty
- [ ] Confirmed column: Shows `✓ Yes` (green) or `✗ No` (gray)
- [ ] Competency column: Color-coded (blue, purple, green, orange, teal, pink)
- [ ] Events column: Shows count
- [ ] Created column: Formatted as YYYY-MM-DD
- [ ] All existing tests pass (100%)
- [ ] New helper tests added (16+ new test cases)
- [ ] Table width ~100 chars (balanced layout)
- [ ] Colors render correctly in terminal
- [ ] Navigation and pagination work
- [ ] Code passes linting and formatting
- [ ] Compliance check passes

## Rollback Plan
If issues are discovered:
1. Identify problematic commit(s)
2. Run: `git revert <commit-hash>`
3. Alternative: `git reset --hard <previous-working-commit>`
4. Re-run tests to verify working state
5. Create new implementation with fixes

## Technical Details

### Column Layout (Total: ~100 chars)
```
Name (22) | Description (25) | Confirmed (8) | Competency (15) | Events (8) | Created (12)
```

### Color Mapping
- **Confirmed**: 
  - ✓ Yes: ColorSuccess (#6cb56c)
  - ✗ No: ColorTextMuted (#5e6673)
- **Competencies**:
  - technical: ColorInfo (#6ab0d3)
  - leadership: ColorAccentPurple (#a99bd1)
  - product: ColorAccentGreen (#6cb56c)
  - consulting: ColorWarning (#d9a66c)
  - research: ColorAccentTeal (#5fb3b3)
  - mentoring: Custom Pink (#d99bd1)

### Helper Function Signatures
```go
func (i *BurstManagementIntent) formatConfirmedStatus(confirmed bool) string
func (i *BurstManagementIntent) formatCompetency(competency string) string
func (i *BurstManagementIntent) formatDescription(description string) string
func (i *BurstManagementIntent) formatCreatedDate(createdAt time.Time) string
```

## Progress Log

### 2026-01-06 - Initial Creation
- Created task file with comprehensive implementation plan
- Defined 9 phases with detailed checklists
- Established acceptance criteria and testing strategy
- Ready to begin implementation

---

**Next Step**: Run baseline compliance check, then begin Phase 2 (Update Column Definitions)
