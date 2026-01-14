---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Task 17: Remove All Legacy Views (internal/cli/models/)

## Overview

- **Goal**: Complete migration of all functionality from legacy `internal/cli/models/` to the new intent-based architecture, then remove the entire legacy directory
- **Time Estimate**: 15-20 days (full migration to removal)
- **Prerequisites**: All 5 primary intents functional (CaptureEvent, BrowseTimeline, GenerateCV, ExportArtifact, ConfigureSystem)
- **Status**: In Progress

## Context

The legacy view system in `internal/cli/models/` (47 source files, 45 test files, ~20,000 lines) predates the intent-based architecture. While the new intent system provides excellent structure, several critical features remain only in the legacy system.

### Key Decisions Made

| Question | Decision | Rationale |
|----------|----------|-----------|
| CSV Import | Use CLI version, re-integrate later | ImportWizard intent is incomplete; CLI `--import` works |
| Burst Suggestions | Add to PRD, not priority | AI-assisted burst review is future work |
| Bulk Operations | Add metadata updates to existing BulkOperationsIntent | Extend existing intent rather than create new one |
| Quality Scoring | Yes, preserve in new system | Documented feature in PRD and METADATA_REVIEW_GUIDE |
| Timeline | No deadline, focus on completeness | Ensure full feature parity before removal |

### Implementation Decisions

| Question | Decision |
|----------|----------|
| CaptureForm Approach | **Full feature parity** - Replicate all legacy features |
| Quality Scoring Location | **All of the above** - Event list, metadata review, detail views |
| Tag/Category Selectors | **Create visual components** in `internal/cli/components/` |
| Test Migration | **Migrate scenarios** (rewrite tests for new architecture) |
| Removal Strategy | **Phased removal** (remove each file as replacement is verified) |

---

## Files to Modify/Create

### Phase 1: Utility Migration
- [ ] `internal/cli/components/pagination.go` - NEW (migrate PaginationHelper)
- [ ] `internal/cli/components/text_utils.go` - NEW (migrate truncateText)
- [ ] `internal/cli/app/messages.go` - MODIFY (consolidate messages)

### Phase 2: CaptureForm Component
- [ ] `internal/cli/intents/capture_form.go` - NEW (full form implementation)
- [ ] `internal/cli/intents/capture_form_test.go` - NEW (form validation tests)
- [ ] `internal/cli/components/tag_selector.go` - NEW (visual tag selector)
- [ ] `internal/cli/components/category_selector.go` - NEW (visual category selector)
- [ ] `internal/cli/intents/capture_event.go` - MODIFY (change captureForm type)
- [ ] `internal/cli/intents/capture_event_intent.go` - MODIFY (update form creation/handling)

### Phase 3: Burst/Fact Editor Integration
- [ ] `internal/cli/intents/burst_management_intent.go` - MODIFY (integrate EditBurstModal)
- [ ] `internal/cli/intents/fact_management_intent.go` - MODIFY (integrate EditFactModal)
- [ ] `internal/cli/intents/burst_management_test.go` - MODIFY (add editor tests)
- [ ] `internal/cli/intents/fact_management_test.go` - MODIFY (add editor tests)

### Phase 4: Metadata Review Enhancement
- [ ] `internal/cli/intents/metadata_editor_intent.go` - MODIFY (add quality scoring, filters, sorting)
- [ ] `internal/cli/components/quality_indicator.go` - NEW (quality score display)
- [ ] `internal/cli/intents/metadata_editor_test.go` - MODIFY (add quality/filter/sort tests)

### Phase 5: Bulk Metadata Operations
- [ ] `internal/cli/intents/bulk_operations_intent.go` - MODIFY (add metadata update operations)
- [ ] `internal/cli/intents/bulk_operations_test.go` - MODIFY (add metadata operation tests)

### Phase 6: Quality Scoring Integration
- [ ] `internal/cli/intents/browse_timeline_intent.go` - MODIFY (add quality score column)
- [ ] `internal/cli/intents/browse_timeline_test.go` - MODIFY (test quality display)

### Phase 7: Legacy Removal (Phased)
- [ ] Remove legacy models one-by-one as replacements are verified
- [ ] Update all imports after removal
- [ ] Final verification and documentation update

---

## Implementation Checklist

### Phase 1: Utility Migration (0.5 days)

**Goal**: Migrate reusable utility code to appropriate locations

- [ ] Create `internal/cli/components/pagination.go`
  - [ ] Port `PaginationHelper` struct from `models/list_patterns.go`
  - [ ] Port methods: `NextPage()`, `PrevPage()`, `SetPage()`, `TotalPages()`, `StartIndex()`, `EndIndex()`
  - [ ] Write unit tests for pagination logic
  - [ ] Verify: All pagination math is correct (off-by-one errors)

- [ ] Create `internal/cli/components/text_utils.go`
  - [ ] Port `truncateText(text string, maxLen int) string` from `models/list_patterns.go`
  - [ ] Add optional ellipsis parameter
  - [ ] Write unit tests for edge cases (empty, exact length, unicode)

- [ ] Update `internal/cli/app/messages.go`
  - [ ] Remove duplicate `BackMsg` from `models/messages.go`
  - [ ] Remove duplicate `EditEventMsg` from `models/messages.go`
  - [ ] Consolidate to single definition in `app/messages.go`

- [ ] Verification
  - [ ] All tests pass: `go test ./internal/cli/components/...`
  - [ ] No regressions in existing code

---

### Phase 2: CaptureForm Component (3-4 days)

**Goal**: Create new intent-local form to replace `models.FormModel`

#### Step 1: Create Visual Selector Components (1 day)

- [ ] Create `internal/cli/components/tag_selector.go`
  - [ ] Implement visual tag selector with allowed tags (8 max)
  - [ ] Features:
    - [ ] Arrow/j/k navigation
    - [ ] Space to toggle selection
    - [ ] Visual indication of selected tags (checkmarks)
    - [ ] Max 8 tag limit with error message
    - [ ] Allowed tags: project, achievement, leadership, technical, consulting, research, product, mentoring
  - [ ] Methods: `Init()`, `Update(msg)`, `View()`, `SelectedTags() []string`
  - [ ] Write unit tests: selection, deselection, navigation, max limit

- [ ] Create `internal/cli/components/category_selector.go`
  - [ ] Implement visual category selector with allowed categories (2 max)
  - [ ] Features:
    - [ ] Arrow/j/k navigation
    - [ ] Space to toggle selection
    - [ ] Visual indication of selected categories
    - [ ] Max 2 category limit with error message
    - [ ] Allowed categories: technical, leadership, product, consulting, research, mentoring
  - [ ] Methods: `Init()`, `Update(msg)`, `View()`, `SelectedCategories() []string`
  - [ ] Write unit tests: selection, deselection, navigation, max limit

#### Step 2: Create CaptureForm Component (2 days)

- [ ] Create `internal/cli/intents/capture_form.go` (NEW - ~800-1000 lines)
  - [ ] Define `CaptureFormModel` struct:
    ```go
    type CaptureFormModel struct {
        textInput     textinput.Model  // Event text (required, 2000 char max)
        dateInput     textinput.Model  // Date (optional, defaults to today)
        companyInput  textinput.Model  // Company (optional)
        projectInput  textinput.Model  // Project (optional)
        tagSelector   *components.TagSelector
        categorySelector *components.CategorySelector
        modeSelector  int  // 0=Timeline, 1=CV Backfill, 2=Manual
        focusedField  int  // Current field index
        submitted     bool
        err           error
        charCounter   map[string]int  // Track char counts per field
    }
    ```

  - [ ] Implement BubbleTea interface:
    - [ ] `Init() tea.Cmd` - Initialize all inputs
    - [ ] `Update(msg tea.Msg) (tea.Model, tea.Cmd)` - Handle:
      - [ ] Tab/Shift+Tab for field navigation
      - [ ] Enter to submit (validate first)
      - [ ] Esc to cancel
      - [ ] Character input for text fields
      - [ ] Delegate to tag/category selectors when focused
    - [ ] `View() string` - Render form with:
      - [ ] Field labels with focus indicators
      - [ ] Character counters (e.g., "1850/2000")
      - [ ] Field-level error messages
      - [ ] Mode selector with current mode highlighted
      - [ ] Help text showing keyboard shortcuts

  - [ ] Implement validation:
    - [ ] `ValidateText()` - Non-empty, max 2000 chars
    - [ ] `ValidateDate()` - Not future, format YYYY-MM-DD or relative ("today", "1 week ago")
    - [ ] `ValidateMode()` - Timeline mode: date within 30 days
    - [ ] `ValidateTags()` - Max 8, all from allowed set
    - [ ] `ValidateCategories()` - Max 2, all from allowed set
    - [ ] `Validate() error` - Validate all fields, return first error

  - [ ] Implement accessors:
    - [ ] `GetEvent() *career.CareerEvent` - Build event from form state
    - [ ] `IsSubmitted() bool`
    - [ ] `GetError() error`

  - [ ] Implement reset:
    - [ ] `Reset()` - Clear all fields, reset state

- [ ] Create `internal/cli/intents/capture_form_test.go` (NEW - ~1500 lines)
  - [ ] Test field navigation:
    - [ ] Tab moves forward through fields
    - [ ] Shift+Tab moves backward
    - [ ] Focus wraps around at boundaries
  - [ ] Test text input:
    - [ ] Accepts valid text
    - [ ] Rejects empty text on submit
    - [ ] Rejects text > 2000 chars
    - [ ] Character counter updates
  - [ ] Test date validation:
    - [ ] Accepts today
    - [ ] Accepts past dates
    - [ ] Rejects future dates
    - [ ] Accepts relative dates ("2 weeks ago")
    - [ ] Accepts YYYY-MM-DD format
  - [ ] Test mode-specific validation:
    - [ ] Timeline mode: rejects dates > 30 days old
    - [ ] CV Backfill mode: accepts any past date
    - [ ] Manual mode: accepts any past date
  - [ ] Test tag selection:
    - [ ] Accepts up to 8 tags
    - [ ] Rejects 9th tag
    - [ ] Accepts only from allowed set
  - [ ] Test category selection:
    - [ ] Accepts up to 2 categories
    - [ ] Rejects 3rd category
    - [ ] Accepts only from allowed set
  - [ ] Test form submission:
    - [ ] Success with valid data
    - [ ] Error with invalid data
    - [ ] Error message displayed
    - [ ] Can correct and resubmit
  - [ ] Test view rendering:
    - [ ] All fields visible
    - [ ] Focus indicators correct
    - [ ] Character counters accurate
    - [ ] Error messages displayed

#### Step 3: Integrate CaptureForm into CaptureEventIntent (1 day)

- [ ] Update `internal/cli/intents/capture_event.go`
  - [ ] Change `captureForm *models.FormModel` to `captureForm *CaptureFormModel`
  - [ ] Remove import of `models` package

- [ ] Update `internal/cli/intents/capture_event_intent.go`
  - [ ] Replace `models.NewFormModel()` with `NewCaptureFormModel()` (line 94)
  - [ ] Replace `models.SubmitMsg` with local `FormSubmittedMsg` (line 278)
  - [ ] Update `GetForm()` return type (line 864)
  - [ ] Update form delegation in `Update()`
  - [ ] Update form view rendering in `View()`

- [ ] Create local message type:
  ```go
  type FormSubmittedMsg struct {
      Event *career.CareerEvent
      Err   error
  }
  ```

- [ ] Verification:
  - [ ] Application builds: `go build ./cmd/kariya`
  - [ ] CaptureEvent intent works end-to-end
  - [ ] Form validation works
  - [ ] Event submission works
  - [ ] No imports from `models` package remain

---

### Phase 3: Burst/Fact Editor Integration (2-3 days)

**Goal**: Wire up existing `EditBurstModal` and `EditFactModal` to management intents

#### Step 1: Integrate EditBurstModal (1 day)

- [ ] Update `internal/cli/intents/burst_management_intent.go`
  - [ ] Import `EditBurstModal` from `modals.go`
  - [ ] Add field to state: `editModal *EditBurstModal`
  - [ ] In `updateEdit()` state (currently stub at line 1082):
    - [ ] Create `EditBurstModal` with current burst
    - [ ] Initialize modal
    - [ ] Store in state
  - [ ] Update `Update()` to delegate to modal when in Edit state:
    ```go
    if i.state.currentState == BurstManagementStateEdit && i.state.editModal != nil {
        updatedModal, cmd := i.state.editModal.Update(msg)
        i.state.editModal = updatedModal.(*EditBurstModal)
        if i.state.editModal.IsDone() {
            if i.state.editModal.WasAccepted() {
                // Save changes to repository
                result := i.state.editModal.GetResult()
                burst := result.Modified
                err := i.context.BurstRepository.Update(burst)
                if err != nil {
                    i.SetError(err)
                } else {
                    // Transition back to Detail state
                    i.state.currentState = BurstManagementStateDetail
                    // Refresh burst
                }
            } else {
                // Cancelled - go back to Detail
                i.state.currentState = BurstManagementStateDetail
            }
        }
        return cmd
    }
    ```
  - [ ] Update `viewEdit()` to render modal:
    ```go
    if i.state.editModal != nil {
        return i.state.editModal.View()
    }
    ```

- [ ] Update `internal/cli/intents/burst_management_test.go`
  - [ ] Add test for Edit state:
    - [ ] Test modal initialization
    - [ ] Test modal cancellation (Esc)
    - [ ] Test modal acceptance (Enter after editing)
    - [ ] Test save success
    - [ ] Test save error
    - [ ] Test burst refresh after save

- [ ] Verification:
  - [ ] Tests pass: `go test ./internal/cli/intents/ -run BurstManagement`
  - [ ] Manual test: Navigate to burst, press 'e' to edit, modify, save
  - [ ] Verify changes persisted to database

#### Step 2: Integrate EditFactModal (1-1.5 days)

- [ ] Update `internal/cli/intents/fact_management_intent.go`
  - [ ] Import `EditFactModal` from `modals.go`
  - [ ] Add field to state: `editModal *EditFactModal`
  - [ ] In `updateEditor()` state (currently stub at lines 524-543):
    - [ ] Create `EditFactModal` with current fact
    - [ ] Initialize modal
    - [ ] Store in state
  - [ ] Update `Update()` to delegate to modal when in Editor state:
    ```go
    if i.state.currentState == FactManagementStateEditor && i.state.editModal != nil {
        updatedModal, cmd := i.state.editModal.Update(msg)
        i.state.editModal = updatedModal.(*EditFactModal)
        if i.state.editModal.IsDone() {
            if i.state.editModal.WasAccepted() {
                result := i.state.editModal.GetResult()
                fact := result.Modified
                err := i.context.FactRepository.Update(fact)
                if err != nil {
                    i.SetError(err)
                } else {
                    i.state.currentState = FactManagementStateView
                    // Refresh fact
                }
            } else {
                i.state.currentState = FactManagementStateView
            }
        }
        return cmd
    }
    ```
  - [ ] Update `viewEditor()` to render modal:
    ```go
    if i.state.editModal != nil {
        return i.state.editModal.View()
    }
    ```

- [ ] Update `internal/cli/intents/fact_management_test.go`
  - [ ] Add test for Editor state:
    - [ ] Test modal initialization with fact data
    - [ ] Test modal cancellation
    - [ ] Test modal acceptance
    - [ ] Test validation (aspirational language check)
    - [ ] Test competency selection
    - [ ] Test role fit selection
    - [ ] Test audience selection
    - [ ] Test save success
    - [ ] Test save error

- [ ] Verification:
  - [ ] Tests pass: `go test ./internal/cli/intents/ -run FactManagement`
  - [ ] Manual test: Navigate to fact, press 'e' to edit, modify all fields, save
  - [ ] Verify aspirational language validation works

---

### Phase 4: Metadata Review Enhancement (2-3 days)

**Goal**: Add quality scoring, filtering, and sorting to MetadataEditorIntent

#### Step 1: Create Quality Indicator Component (0.5 day)

- [ ] Create `internal/cli/components/quality_indicator.go`
  - [ ] Define `QualityIndicator` struct:
    ```go
    type QualityIndicator struct {
        Score int  // 0-100
    }
    ```
  - [ ] Implement rendering:
    - [ ] `Render() string` - Returns colored bar like: `████████░░ 75%`
    - [ ] Color coding:
      - 0-25: Red (Incomplete)
      - 26-50: Yellow (Basic)
      - 51-75: Green (Enriched)
      - 76-100: Bright Green (Complete)
    - [ ] `GetLevel() string` - Returns "Incomplete", "Basic", "Enriched", "Complete"
    - [ ] `GetMissingFields(*career.CareerEvent) []string` - Returns list of missing fields

- [ ] Write tests:
  - [ ] Test score to color mapping
  - [ ] Test bar rendering at different scores
  - [ ] Test level categorization
  - [ ] Test missing fields detection

#### Step 2: Calculate Quality Scores (0.5 day)

- [ ] Add quality calculation to `internal/domain/career/event.go` (or create utility)
  - [ ] `CalculateQualityScore(event *CareerEvent) int`:
    - Text: 20 points (always present)
    - Date: 20 points
    - Company: 20 points
    - Project: 20 points
    - Tags: 15 points (scaled by count/8)
    - Categories: 15 points (scaled by count/2)
    - Bonus: 10 points for high match quality (if applicable)
  - [ ] Write unit tests for various event completeness levels

#### Step 3: Enhance MetadataEditorIntent (1-2 days)

- [ ] Update `internal/cli/intents/metadata_editor_intent.go`
  - [ ] Add filtering support:
    - [ ] Add state: `filterMode string` (values: "all", "incomplete", "basic")
    - [ ] Add method: `filterEvents() []*career.CareerEvent`
    - [ ] Add key handler: 'f' toggles filter mode
    - [ ] Update view to show current filter in header

  - [ ] Add sorting support:
    - [ ] Add state: `sortBy string` (values: "date", "quality", "company")
    - [ ] Add method: `sortEvents([]*career.CareerEvent)`
    - [ ] Add key handler: 's' cycles through sort modes
    - [ ] Update view to show current sort in header

  - [ ] Add quality display:
    - [ ] Calculate quality score for each event
    - [ ] Render quality indicator in event list
    - [ ] Show missing fields when event is focused
    - [ ] Add quality column to table (if using table view)

  - [ ] Update event list view:
    - [ ] Use `components.TableListContainer` or custom rendering
    - [ ] Columns: [Checkbox] Event Text | Quality | Company | Date
    - [ ] Quality column uses `QualityIndicator.Render()`
    - [ ] Pagination using `components.PaginationHelper`

- [ ] Update `internal/cli/intents/metadata_editor_test.go`
  - [ ] Test filtering:
    - [ ] Filter by incomplete (0-25)
    - [ ] Filter by basic (26-50)
    - [ ] Filter shows correct events
  - [ ] Test sorting:
    - [ ] Sort by quality (ascending/descending)
    - [ ] Sort by date
    - [ ] Sort by company
  - [ ] Test quality display:
    - [ ] Quality score calculated correctly
    - [ ] Quality indicator renders
    - [ ] Missing fields shown

- [ ] Verification:
  - [ ] Tests pass: `go test ./internal/cli/intents/ -run MetadataEditor`
  - [ ] Manual test: View events, toggle filters, sort, verify quality scores

---

### Phase 5: Bulk Metadata Operations (2 days)

**Goal**: Add bulk metadata update operations to BulkOperationsIntent

- [ ] Update `internal/cli/intents/bulk_operations_intent.go`
  - [ ] Add new operation types to enum:
    - `BulkOpUpdateCompany` - Set company for all selected
    - `BulkOpUpdateProject` - Set project for all selected
    - `BulkOpUpdateTags` - Add tags to all selected
    - `BulkOpUpdateCategories` - Add categories to all selected

  - [ ] Add configuration state for metadata operations:
    ```go
    type BulkMetadataConfig struct {
        Company      string
        Project      string
        Tags         []string
        Categories   []string
        ApplyIfEmpty bool  // Only update if field is empty
    }
    ```

  - [ ] In `BulkConfigureState`:
    - [ ] Show input fields based on operation type
    - [ ] For company/project: text input
    - [ ] For tags: use `components.TagSelector`
    - [ ] For categories: use `components.CategorySelector`
    - [ ] Add checkbox for "Apply if empty"

  - [ ] In `BulkExecuteState`:
    - [ ] Implement `executeBulkMetadataUpdate()`:
      ```go
      for _, item := range selectedItems {
          if config.ApplyIfEmpty && item.Company != "" {
              continue  // Skip if field already has value
          }
          item.Company = config.Company
          err := i.context.EventRepository.Update(item)
          // Track success/failure
      }
      ```
    - [ ] Update progress modal during execution
    - [ ] Track counts: updated, skipped, failed

  - [ ] Update view rendering:
    - [ ] `viewSelectOp()`: Show new operation types in list
    - [ ] `viewConfigure()`: Render metadata input fields
    - [ ] `viewComplete()`: Show summary (X updated, Y skipped)

- [ ] Update `internal/cli/intents/bulk_operations_test.go`
  - [ ] Test bulk company update:
    - [ ] Updates all selected events
    - [ ] Respects "apply if empty" flag
    - [ ] Reports correct counts
  - [ ] Test bulk tags update:
    - [ ] Adds tags to all selected
    - [ ] Validates tags are in allowed set
    - [ ] Respects max 8 tags limit
  - [ ] Test bulk categories update:
    - [ ] Adds categories to all selected
    - [ ] Validates categories are in allowed set
    - [ ] Respects max 2 categories limit
  - [ ] Test error handling:
    - [ ] Continues on individual item failure
    - [ ] Reports failed items
    - [ ] Allows retry

- [ ] Verification:
  - [ ] Tests pass: `go test ./internal/cli/intents/ -run BulkOperations`
  - [ ] Manual test: Select events, bulk update company, verify changes
  - [ ] Manual test: Select events, bulk add tags, verify max limit respected

---

### Phase 6: Quality Scoring Integration (1 day)

**Goal**: Display quality scores in BrowseTimelineIntent

- [ ] Update `internal/cli/intents/browse_timeline_intent.go`
  - [ ] Add quality score calculation to event loading
  - [ ] Add quality column to table view:
    - Columns: Event Text | Quality | Company | Date
    - Quality column uses `QualityIndicator.Render()`
  - [ ] Update `viewTimeline()`:
    - [ ] Calculate quality for each visible event
    - [ ] Render quality indicator in table
    - [ ] Consider color-coding row backgrounds by quality

- [ ] Update `internal/cli/intents/browse_timeline_test.go`
  - [ ] Test quality column rendering
  - [ ] Test quality scores accurate for various events

- [ ] Verification:
  - [ ] Tests pass: `go test ./internal/cli/intents/ -run BrowseTimeline`
  - [ ] Manual test: View event list, verify quality scores shown

---

### Phase 7: Legacy Removal (1-2 days, phased)

**Goal**: Remove legacy `internal/cli/models/` directory in phases

#### Removal Order (Remove as Replacements are Verified)

**Round 1: Remove Unused/Dead Code**
- [ ] Remove `internal/cli/models/cv_preview.go` (unused `cvPreview` field in GenerateCVIntent)
- [ ] Remove associated test file
- [ ] Verify: Build succeeds, no imports broken

**Round 2: Remove Form Components (After Phase 2 Complete)**
- [ ] Remove `internal/cli/models/form.go`
- [ ] Remove `internal/cli/models/form_*.go` (all 12 form-related files)
- [ ] Remove form test files
- [ ] Update imports in `capture_event_intent.go`
- [ ] Verify: CaptureEvent intent works, all tests pass

**Round 3: Remove Burst/Fact Components (After Phase 3 Complete)**
- [ ] Remove `internal/cli/models/burst_editor.go`
- [ ] Remove `internal/cli/models/fact_editor.go`
- [ ] Remove `internal/cli/models/burst_details.go`
- [ ] Remove `internal/cli/models/fact_details.go`
- [ ] Remove `internal/cli/models/burst_card.go`
- [ ] Remove `internal/cli/models/fact_card.go`
- [ ] Remove associated test files
- [ ] Verify: BurstManagement and FactManagement intents work

**Round 4: Remove Metadata Components (After Phase 4-5 Complete)**
- [ ] Remove `internal/cli/models/metadata_review.go`
- [ ] Remove `internal/cli/models/metadata_editor.go`
- [ ] Remove `internal/cli/models/bulk_operations.go`
- [ ] Remove associated test files
- [ ] Verify: MetadataEditor and BulkOperations intents work

**Round 5: Remove List Components**
- [ ] Remove `internal/cli/models/list.go`
- [ ] Remove `internal/cli/models/list_patterns.go` (after PaginationHelper migrated)
- [ ] Remove `internal/cli/models/*_list.go` files
- [ ] Remove associated test files
- [ ] Verify: All list-based intents work

**Round 6: Remove Infrastructure Files**
- [ ] Remove `internal/cli/models/standard_model.go`
- [ ] Remove `internal/cli/models/messages.go` (after consolidation)
- [ ] Remove `internal/cli/models/errors.go`
- [ ] Remove `internal/cli/models/shortcut_*.go` files
- [ ] Remove associated test files

**Round 7: Remove Remaining Files**
- [ ] Remove `internal/cli/models/menu.go`
- [ ] Remove `internal/cli/models/help.go`
- [ ] Remove `internal/cli/models/success.go`
- [ ] Remove `internal/cli/models/confirmation_dialog.go`
- [ ] Remove all remaining files in `internal/cli/models/`
- [ ] Remove test files

**Round 8: Final Cleanup**
- [ ] Remove entire `internal/cli/models/` directory
- [ ] Search for any remaining imports: `grep -r "internal/cli/models" .`
- [ ] Fix any missed imports
- [ ] Run full test suite: `go test ./...`
- [ ] Run with race detector: `go test -race ./...`
- [ ] Verify coverage maintained: `go test -cover ./...`

---

## Testing Instructions

### After Each Phase

1. **Build verification**:
   ```bash
   go build ./cmd/kariya
   ```

2. **Run affected tests**:
   ```bash
   # Phase 2 (CaptureForm)
   go test ./internal/cli/intents/ -run CaptureEvent -v
   
   # Phase 3 (Burst/Fact Editors)
   go test ./internal/cli/intents/ -run BurstManagement -v
   go test ./internal/cli/intents/ -run FactManagement -v
   
   # Phase 4-5 (Metadata Review & Bulk Ops)
   go test ./internal/cli/intents/ -run MetadataEditor -v
   go test ./internal/cli/intents/ -run BulkOperations -v
   
   # Phase 6 (Quality Scoring)
   go test ./internal/cli/intents/ -run BrowseTimeline -v
   ```

3. **Run with race detector**:
   ```bash
   go test -race ./internal/cli/intents/...
   ```

4. **Manual testing workflow**:
   ```bash
   # Build application
   go build -o kariya ./cmd/kariya
   
   # Test workflow for current phase
   ./kariya
   # Navigate to affected intent
   # Test all new functionality
   # Verify no regressions
   ```

### Before Each Legacy File Removal

1. **Verify replacement exists and works**:
   - Run tests for replacement component
   - Manual test of workflow

2. **Check for imports**:
   ```bash
   grep -r "models/filename_without_extension" .
   ```

3. **Remove file**:
   ```bash
   git rm internal/cli/models/filename.go
   git rm internal/cli/models/filename_test.go
   ```

4. **Verify build**:
   ```bash
   go build ./cmd/kariya
   ```

5. **Run full test suite**:
   ```bash
   go test ./...
   ```

### Final Verification (After Phase 7 Complete)

1. **Full test suite**:
   ```bash
   go test ./...
   go test -race ./...
   go test -cover ./...
   ```

2. **Verify no imports remain**:
   ```bash
   grep -r "internal/cli/models" . | grep -v ".git" | grep -v "tasks-17"
   # Should return empty
   ```

3. **Manual testing of all intents**:
   - CaptureEvent: Capture event with all fields, all modes
   - BrowseTimeline: Browse events, view details, check quality scores
   - GenerateCV: Generate CV, preview, export
   - ExportArtifact: Export different formats
   - ConfigureSystem: Modify settings
   - BurstManagement: View, edit, delete bursts
   - FactManagement: View, edit, delete facts
   - MetadataEditor: Review, filter, sort, edit events
   - BulkOperations: Bulk update metadata

4. **Coverage verification**:
   ```bash
   go test -coverprofile=coverage.out ./...
   go tool cover -html=coverage.out
   # Verify coverage >= 87%
   ```

---

## Acceptance Criteria

### Phase 1: Utility Migration
- [ ] `PaginationHelper` works identically in `components` package
- [ ] `truncateText` handles all edge cases
- [ ] No duplicate messages in codebase
- [ ] All tests pass

### Phase 2: CaptureForm Component
- [ ] CaptureForm has full feature parity with legacy FormModel
- [ ] Visual tag/category selectors work (arrow/j/k navigation, space toggle)
- [ ] All validation rules enforced (text, date, mode-specific)
- [ ] Character counter updates in real-time
- [ ] Field-level error messages displayed
- [ ] Tab/Shift+Tab navigation works
- [ ] Form submission creates valid CareerEvent
- [ ] All form tests migrated (~2000 lines of test scenarios)
- [ ] CaptureEvent intent uses new form (no `models` import)

### Phase 3: Burst/Fact Editor Integration
- [ ] EditBurstModal integrated into BurstManagementIntent
- [ ] EditFactModal integrated into FactManagementIntent
- [ ] Burst editing saves changes to repository
- [ ] Fact editing validates aspirational language
- [ ] Cancel works (reverts changes)
- [ ] All editor tests passing

### Phase 4: Metadata Review Enhancement
- [ ] Quality scores displayed for all events
- [ ] Quality levels: Incomplete (0-25), Basic (26-50), Enriched (51-75), Complete (76-100)
- [ ] Quality indicator shows colored bar and percentage
- [ ] Missing fields shown when event focused
- [ ] Filtering works (all, incomplete, basic)
- [ ] Sorting works (date, quality, company)
- [ ] Table-based event list with quality column
- [ ] Pagination works

### Phase 5: Bulk Metadata Operations
- [ ] Bulk company update works
- [ ] Bulk project update works
- [ ] Bulk tags update works (max 8, allowed set)
- [ ] Bulk categories update works (max 2, allowed set)
- [ ] "Apply if empty" flag respected
- [ ] Progress tracking during bulk operations
- [ ] Summary shows updated/skipped/failed counts
- [ ] Error handling for individual failures

### Phase 6: Quality Scoring Integration
- [ ] Quality scores shown in BrowseTimeline event list
- [ ] Quality column in table view
- [ ] Color-coded quality indicators
- [ ] Accurate score calculation

### Phase 7: Legacy Removal
- [ ] Entire `internal/cli/models/` directory removed
- [ ] No imports from `models` package anywhere
- [ ] All tests pass (980+ tests)
- [ ] No regressions in functionality
- [ ] Code coverage >= 87%
- [ ] All 5 primary intents functional
- [ ] All management intents functional

### Overall Acceptance
- [ ] **Zero imports** from `internal/cli/models/` package
- [ ] **Full feature parity** with legacy system
- [ ] **Test coverage maintained** (>87%)
- [ ] **All workflows functional** (capture, browse, generate, export, configure, manage)
- [ ] **Documentation updated** (AGENTS.md, guides)
- [ ] **Performance maintained** (no regressions)

---

## Rollback Plan

### If Phase Fails

1. **Revert code changes**:
   ```bash
   git checkout -- <affected files>
   ```

2. **Restore legacy dependency** (if removed prematurely):
   ```bash
   git checkout -- internal/cli/models/<removed file>
   ```

3. **Verify build and tests**:
   ```bash
   go build ./cmd/kariya
   go test ./...
   ```

### Safety Measures

- Commit after each phase completion
- Tag successful milestones: `git tag phase-N-complete`
- Keep legacy files until replacement verified
- Phased removal (not big bang)
- Comprehensive testing before each removal

---

## Dependencies

### External
- BubbleTea framework (for interactive components)
- Lipgloss (for styling)
- Charmbracelet bubbles (textinput, table)

### Internal
- `internal/domain/career` - Domain models
- `internal/service/career` - Business logic
- `internal/repository/career` - Data access
- `internal/cli/components` - Reusable UI components
- `internal/cli/intents` - Intent framework

---

## Notes

### Out of Scope (Future Work)

These features exist in legacy but are deprioritized:

1. **CSV Import UI** - Use CLI `--import` command; re-integrate UI later
2. **Burst Suggestions** - AI-suggested burst review workflow (add to PRD)
3. **Fact Search** - Search/filter facts by query (future enhancement)

### Quality Scoring Formula

Based on METADATA_REVIEW_GUIDE.md:

```
Quality Score = Text (20) + Date (20) + Company (20) + Project (20) 
                + Tags (15 * count/8) + Categories (15 * count/2) 
                + Bonus (10 if high quality)

Levels:
- Incomplete: 0-25   (Red)
- Basic:      26-50  (Yellow)
- Enriched:   51-75  (Green)
- Complete:   76-100 (Bright Green)
```

### PRD References

- **PRD_MASTER.md**: Core concepts, validation rules, CV generation
- **PRD_USER_STORIES.md**: User stories, acceptance criteria
- **PRD_CLI.md**: Event capture requirements, validation rules
- **METADATA_REVIEW_GUIDE.md**: Quality scoring spec, filtering, bulk operations

---

## Completion Checklist

- [ ] Phase 1: Utility Migration (0.5 days)
- [ ] Phase 2: CaptureForm Component (3-4 days)
- [ ] Phase 3: Burst/Fact Editor Integration (2-3 days)
- [ ] Phase 4: Metadata Review Enhancement (2-3 days)
- [ ] Phase 5: Bulk Metadata Operations (2 days)
- [ ] Phase 6: Quality Scoring Integration (1 day)
- [ ] Phase 7: Legacy Removal (1-2 days, phased)
- [ ] All tests passing (980+ tests)
- [ ] Coverage >= 87%
- [ ] No `models` package imports
- [ ] Documentation updated
- [ ] AGENTS.md updated (remove references to legacy models)

**Total Estimated Effort**: 15-20 days
**Status**: Ready to begin
