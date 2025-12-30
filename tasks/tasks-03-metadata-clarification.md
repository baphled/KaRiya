# Task List: Metadata Clarification Feature

**PRD Reference**: `docs/features/02-metadata-clarification.md`

**Purpose**: Enable users to review, validate, and enrich event metadata (dates, companies, projects, tags, categories) before automated processing (burst detection and fact inference).

**Status**: Phase 1 & Phase 2 Tasks 1-6 Complete (70% overall)

---

## Completed Work Summary

### Phase 1: Foundation & Core Components (100% COMPLETE) ✅
- ✅ Task 1.0: Data Quality Scoring System (35 tests passing)
- ✅ Task 2.0: Metadata Validation System (49 tests passing)
- ✅ Task 3.0: Quality Indicator Component (16 tests passing)

### Phase 2: Metadata Review & Management (70% COMPLETE)
- ✅ Task 4.0: Metadata Review Screen Model (16 tests passing)
- ✅ Task 5.0: Metadata Review Navigation Integration (4 integration tests)
- ✅ Task 6.0: Individual Event Metadata Editor (16 tests passing) - **COMPLETED 2025-12-30**
- ⏳ Task 7.0: Metadata Editor Navigation (TODO)
- ⏳ Task 8.0: CLI Service Enhancement (TODO)
- ⏳ Task 9.0: Bulk Operations (TODO)
- ⏳ Task 10.0: Bulk Operations Integration (TODO)
- ⏳ Task 11.0: Bulk Service Enhancement (TODO)
- ⏳ Task 12.0: CSV Import Integration (TODO)
- ⏳ Task 13.0: Manual Capture Integration (TODO)
- ⏳ Task 14.0: Comprehensive Testing (TODO)
- ⏳ Task 15.0: Documentation (TODO)

---

## Relevant Files

### Files Created ✅

- ✅ `internal/service/career/data_quality.go` - Data quality scoring system
- ✅ `internal/service/career/data_quality_test.go` - 35 quality tests
- ✅ `internal/cli/validation/metadata_validator_test.go` - 21 validator tests
- ✅ `internal/cli/models/quality_indicator.go` - Visual quality indicator
- ✅ `internal/cli/models/quality_indicator_test.go` - 16 indicator tests
- ✅ `internal/cli/models/metadata_review.go` - Metadata review screen
- ✅ `internal/cli/models/metadata_review_test.go` - 16 review tests
- ✅ `internal/cli/models/metadata_editor.go` - Individual event metadata editor
- ✅ `internal/cli/models/metadata_editor_test.go` - 16 editor tests

### Files Modified ✅

- ✅ `internal/cli/validation/validator.go` - Extended with MetadataValidator
- ✅ `internal/cli/app/app.go` - Added MetadataReviewScreen navigation
- ✅ `internal/cli/app/app_test.go` - Added 4 integration tests

### Files to Create (Remaining)

- `internal/cli/models/bulk_operations.go` - Bulk operations model
- `internal/cli/models/bulk_operations_test.go` - Bulk operations tests
- `internal/cli/service/metadata_service.go` - Enhanced metadata operations

### Files to Modify (Remaining)

- `internal/cli/models/import_review.go` - Trigger metadata review after import
- `internal/cli/models/form.go` - Post-capture metadata enrichment
- `internal/cli/models/success.go` - Quick metadata review option
- `internal/cli/app/app.go` - Add editor and bulk screens (additional changes)

---

## Task Details

### Phase 1: Foundation & Core Components ✅

#### 1.0 Create Data Quality Scoring System ✅
- ✅ 1.1 Implement `data_quality.go` with quality score calculation (0-100 scale)
- ✅ 1.2 Create quality level constants (Incomplete, Basic, Enriched, Complete)
- ✅ 1.3 Implement scoring logic: Text(+20), Date(+20), Company/Project(+20), Tags(+15), Categories(+15), Match(+10)
- ✅ 1.4 Add method to determine quality status from score
- ✅ 1.5 Write comprehensive unit tests covering all scoring scenarios
- ✅ 1.6 Verify quality scoring integrates with CareerEvent domain model

#### 2.0 Create Metadata Validation System ✅
- ✅ 2.1 Implement `metadata_validator.go` with field-specific validators
- ✅ 2.2 Create date validator (not future, reasonable range, format validation)
- ✅ 2.3 Create company validator (optional, max 200 chars, normalization)
- ✅ 2.4 Create project validator (optional, max 200 chars, normalization)
- ✅ 2.5 Create tags validator (from AllowedTags, max 8, no duplicates, case-insensitive)
- ✅ 2.6 Create categories validator (from AllowedCategories, match validation)
- ✅ 2.7 Write comprehensive unit tests for all validators with edge cases
- ✅ 2.8 Ensure validators return helpful error messages for user feedback

#### 3.0 Create Quality Indicator Component ✅
- ✅ 3.1 Implement `quality_indicator.go` as visual component
- ✅ 3.2 Design visual representation (color-coded, icon-based, percentage display)
- ✅ 3.3 Use existing styling system from `internal/cli/styles/`
- ✅ 3.4 Show quality level (Incomplete/Basic/Enriched/Complete)
- ✅ 3.5 Display missing fields suggestion
- ✅ 3.6 Write unit tests for quality indicator rendering
- ✅ 3.7 Verify integration with metadata review screen

### Phase 2: Metadata Review & Management ✅

#### 4.0 Create Metadata Review Screen Model ✅
- ✅ 4.1 Implement `metadata_review.go` as BubbleTea Model
- ✅ 4.2 Fetch events from service and calculate quality scores
- ✅ 4.3 Display events in list with quality indicators
- ✅ 4.4 Implement filtering by quality level (All/Incomplete/Basic/Enriched/Complete)
- ✅ 4.5 Implement sorting (by quality/date/company)
- ✅ 4.6 Add navigation with Up/Down arrows, Space to expand, Enter to edit
- ✅ 4.7 Show quality details and suggestions for improvement
- ✅ 4.8 Write comprehensive unit tests for review screen

#### 5.0 Integrate Metadata Review with Navigation ✅
- ✅ 5.1 Add MetadataReviewScreen constant to app.go
- ✅ 5.2 Add MetadataReviewModel field to app Model
- ✅ 5.3 Implement navigation: 'm' key opens metadata review
- ✅ 5.4 Implement back navigation from review screen
- ✅ 5.5 Add integration tests for screen transitions
- ✅ 5.6 Verify message routing from review to app

#### 6.0 Create Individual Event Metadata Editor ✅ **COMPLETED 2025-12-30**
- ✅ 6.1 Implement `metadata_editor.go` as BubbleTea Model
- ✅ 6.2 Create form fields for: date, company, project, tags, categories
- ✅ 6.3 Implement Tab/Shift+Tab navigation between fields
- ✅ 6.4 Implement Up/Down arrows for tag/category selection
- ✅ 6.5 Implement Space key to toggle tag/category selection
- ✅ 6.6 Add Save/Cancel buttons with Enter key handling
- ✅ 6.7 Implement date parsing (ISO format, "today", relative dates)
- ✅ 6.8 Add field validation with error display
- ✅ 6.9 Implement revert functionality to restore original
- ✅ 6.10 Write 16 comprehensive unit tests
- ✅ 6.11 Verify integration with metadata service

#### 7.0 Integrate Metadata Editor with Navigation (TODO)
- ⏳ 7.1 Add MetadataEditorScreen constant to app.go
- ⏳ 7.2 Add MetadataEditorModel field to app Model
- ⏳ 7.3 Implement navigation: Enter from review screen opens editor
- ⏳ 7.4 Handle editor submission: update event and return to review
- ⏳ 7.5 Handle editor cancellation: return to review without changes
- ⏳ 7.6 Add integration tests for editor transitions

#### 8.0 Enhance CLI Service with Metadata Operations (TODO)
- ⏳ 8.1 Add UpdateEventMetadata() method to CLIEventService
- ⏳ 8.2 Implement validation before persistence
- ⏳ 8.3 Add error handling with helpful messages
- ⏳ 8.4 Implement transaction support for atomic updates
- ⏳ 8.5 Add UpdateEventMetadata tests
- ⏳ 8.6 Verify service integration with repository

---

## Phase 3: Bulk Operations (TODO)

#### 9.0 Create Bulk Operations Model (TODO)
- ⏳ 9.1 Implement `bulk_operations.go` as BubbleTea Model
- ⏳ 9.2 Add checkboxes for event selection
- ⏳ 9.3 Implement bulk edit form (company, project, tags, categories)
- ⏳ 9.4 Add confirmation dialog for bulk changes
- ⏳ 9.5 Implement undo/revert functionality
- ⏳ 9.6 Write unit tests for bulk operations

#### 10.0 Integrate Bulk Operations with Metadata Review (TODO)
- ⏳ 10.1 Add checkbox UI to metadata review screen
- ⏳ 10.2 Implement multi-select mode (Shift+Space)
- ⏳ 10.3 Add "Bulk Edit" button when items selected
- ⏳ 10.4 Navigate to bulk edit screen on button press
- ⏳ 10.5 Add integration tests

#### 11.0 Enhance CLI Service with Bulk Operations (TODO)
- ⏳ 11.1 Add UpdateEventsBulk() method
- ⏳ 11.2 Implement batch validation
- ⏳ 11.3 Add transaction support
- ⏳ 11.4 Implement rollback on partial failure
- ⏳ 11.5 Write comprehensive tests

---

## Phase 4: Integration (TODO)

#### 12.0 Enhance CSV Import Integration (TODO)
- ⏳ 12.1 Trigger metadata review after CSV import
- ⏳ 12.2 Show quality scores for imported events
- ⏳ 12.3 Allow immediate metadata enrichment
- ⏳ 12.4 Add integration tests

#### 13.0 Enhance Manual Capture Integration (TODO)
- ⏳ 13.1 Show metadata quality after form submission
- ⏳ 13.2 Offer quick metadata review from success screen
- ⏳ 13.3 Allow returning to editor for quick fixes
- ⏳ 13.4 Add integration tests

---

## Phase 5: Testing & Documentation (TODO)

#### 14.0 Comprehensive Testing Suite (TODO)
- ⏳ 14.1 End-to-end workflow tests
- ⏳ 14.2 Error scenario testing
- ⏳ 14.3 Performance testing
- ⏳ 14.4 Race condition testing
- ⏳ 14.5 Integration test coverage

#### 15.0 Documentation & User Guidance (TODO)
- ⏳ 15.1 Write user guide for metadata review
- ⏳ 15.2 Create tutorial for metadata editing
- ⏳ 15.3 Document quality scoring system
- ⏳ 15.4 Add troubleshooting guide
- ⏳ 15.5 Update main README

---

## Test Summary

### Phase 1 Tests: 100 tests passing ✅
- Data Quality: 35 tests
- Metadata Validation: 49 tests
- Quality Indicator: 16 tests

### Phase 2 Tests: 36 tests passing ✅
- Metadata Review: 16 tests
- Navigation Integration: 4 tests
- Metadata Editor: 16 tests

### Total: 136 tests passing ✅

---

## Recent Changes (2025-12-30)

### Task 6.0 Completion
- ✅ Created `metadata_editor.go` with full editing functionality
- ✅ Implemented form fields for date, company, project, tags, categories
- ✅ Added Tab/Shift+Tab navigation and Arrow key support
- ✅ Implemented Space key for tag/category selection
- ✅ Created 16 comprehensive tests
- ✅ Fixed test suite organization (moved to models_test package)
- ✅ All tests passing (100% success rate)

---

## Next Session Tasks

1. **Task 7.0**: Integrate metadata editor with navigation
   - Add MetadataEditorScreen to app.go
   - Connect from metadata review screen
   - Implement save/cancel handling
   - Estimated: 2-3 hours

2. **Task 8.0**: Enhance CLI service with metadata operations
   - Add UpdateEventMetadata() method
   - Implement validation and persistence
   - Write tests
   - Estimated: 2-3 hours

3. **Task 9.0**: Create bulk operations model
   - Implement checkbox selection
   - Create bulk edit form
   - Add confirmation dialog
   - Estimated: 3-4 hours

---

**Last Updated**: 2025-12-30 15:30 UTC
**Completed By**: Senior Development Engineer
**Status**: Ready for Task 7.0 (Navigation Integration)
