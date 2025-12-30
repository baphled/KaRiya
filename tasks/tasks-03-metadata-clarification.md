# Task List: Metadata Clarification Feature

**PRD Reference**: `docs/features/02-metadata-clarification.md`

**Purpose**: Enable users to review, validate, and enrich event metadata (dates, companies, projects, tags, categories) before automated processing (burst detection and fact inference).

**Status**: Not Started

---

## Relevant Files

### New Files to Create

- `internal/cli/models/metadata_review.go` - Main metadata review screen model showing list of events awaiting clarification
- `internal/cli/models/metadata_review_test.go` - Unit tests for metadata review screen
- `internal/cli/models/metadata_editor.go` - Individual event metadata editor with field-level editing
- `internal/cli/models/metadata_editor_test.go` - Unit tests for metadata editor
- `internal/cli/models/bulk_operations.go` - Bulk metadata operations (select multiple, batch edit)
- `internal/cli/models/bulk_operations_test.go` - Unit tests for bulk operations
- `internal/cli/validation/metadata_validator.go` - Metadata validation rules (date, company, project, tags, categories)
- `internal/cli/validation/metadata_validator_test.go` - Unit tests for metadata validator
- `internal/service/career/data_quality.go` - Data quality scoring and calculation logic
- `internal/service/career/data_quality_test.go` - Unit tests for data quality scoring
- `internal/cli/models/quality_indicator.go` - Visual data quality indicator component
- `internal/cli/models/quality_indicator_test.go` - Unit tests for quality indicator

### Files to Modify

- `internal/cli/app/app.go` - Add MetadataReviewScreen, MetadataEditorScreen, BulkOperationsScreen states and navigation
- `internal/cli/app/messages.go` - Add new message types for metadata operations
- `internal/cli/models/import_review.go` - Enhance to trigger metadata review after CSV import
- `internal/cli/models/form.go` - Add post-capture metadata enrichment option
- `internal/cli/models/success.go` - Add quick metadata review option after manual capture
- `internal/cli/service/event_service.go` - Add metadata update and bulk update methods
- `internal/domain/career/event.go` - Add or verify data quality calculation method

### Test Files (Existing - May Need Updates)

- `internal/cli/models/suite_test.go` - Update test suite if needed for new models
- `internal/cli/app/app_test.go` - Add navigation tests for new screens

### Notes

- Unit tests should be placed alongside the code files they test (e.g., `metadata_review.go` and `metadata_review_test.go`)
- Use `make test` to run all tests or `ginkgo -v ./internal/cli/models` to run specific package tests
- Follow existing BubbleTea model patterns used in `form.go`, `list.go`, and `details.go`
- Leverage existing validation patterns and styling from `internal/cli/styles/`

---

## Tasks

### Phase 1: Foundation & Core Components

- [ ] 1.0 Create Data Quality Scoring System
  - [ ] 1.1 Implement `data_quality.go` with quality score calculation (0-100 scale)
  - [ ] 1.2 Create quality level constants (Incomplete, Basic, Enriched, Complete)
  - [ ] 1.3 Implement scoring logic: Text(+20), Date(+20), Company/Project(+20), Tags(+15), Categories(+15), Match(+10)
  - [ ] 1.4 Add method to determine quality status from score
  - [ ] 1.5 Write comprehensive unit tests covering all scoring scenarios
  - [ ] 1.6 Verify quality scoring integrates with CareerEvent domain model

- [ ] 2.0 Create Metadata Validation System
  - [ ] 2.1 Implement `metadata_validator.go` with field-specific validators
  - [ ] 2.2 Create date validator (not future, reasonable range, format validation)
  - [ ] 2.3 Create company validator (optional, max 200 chars, normalization)
  - [ ] 2.4 Create project validator (optional, max 200 chars, normalization)
  - [ ] 2.5 Create tags validator (from AllowedTags, max 8, no duplicates, case-insensitive)
  - [ ] 2.6 Create categories validator (from AllowedCategories, match validation)
  - [ ] 2.7 Write comprehensive unit tests for all validators with edge cases
  - [ ] 2.8 Ensure validators return helpful error messages for user feedback

- [ ] 3.0 Create Quality Indicator Component
  - [ ] 3.1 Implement `quality_indicator.go` as visual component
  - [ ] 3.2 Design visual representation (color-coded, icon-based, percentage display)
  - [ ] 3.3 Use existing styling system from `internal/cli/styles/`
  - [ ] 3.4 Show quality level (Incomplete/Basic/Enriched/Complete)
  - [ ] 3.5 Display missing fields suggestion
  - [ ] 3.6 Write unit tests for quality indicator rendering
  - [ ] 3.7 Verify integration with metadata review screen

### Phase 2: Metadata Review Screen

- [ ] 4.0 Create Metadata Review Screen Model
  - [ ] 4.1 Implement `metadata_review.go` with BubbleTea Model interface
  - [ ] 4.2 Create list view showing events awaiting clarification
  - [ ] 4.3 Display event text (truncated), date, company, project, tags, categories for each event
  - [ ] 4.4 Integrate quality indicator for each event
  - [ ] 4.5 Implement scrolling (up/down arrows) through event list
  - [ ] 4.6 Add expand/collapse for full event text viewing
  - [ ] 4.7 Implement filtering by data quality (incomplete only, all, etc.)
  - [ ] 4.8 Implement sorting by date, company, or creation order
  - [ ] 4.9 Implement keyboard navigation (↑/↓ for events, Enter to edit, Space for select)
  - [ ] 4.10 Write comprehensive unit tests covering all interactions
  - [ ] 4.11 Test edge cases (empty list, single event, large event list)

- [ ] 5.0 Integrate Metadata Review with Navigation
  - [ ] 5.1 Add MetadataReviewScreen constant to `app.go`
  - [ ] 5.2 Add state management for metadata review screen
  - [ ] 5.3 Add navigation trigger from home screen (keyboard shortcut or menu)
  - [ ] 5.4 Add navigation trigger from post-capture success screen
  - [ ] 5.5 Add navigation trigger after CSV import completion
  - [ ] 5.6 Implement back/exit from metadata review screen
  - [ ] 5.7 Write app integration tests for screen navigation

### Phase 3: Metadata Editor

- [ ] 6.0 Create Individual Event Metadata Editor
  - [ ] 6.1 Implement `metadata_editor.go` as BubbleTea Model
  - [ ] 6.2 Create form with fields: Date, Company, Project, Tags, Categories
  - [ ] 6.3 Implement date field with calendar picker or text input
  - [ ] 6.4 Implement company field with autocomplete from previous entries
  - [ ] 6.5 Implement project field with autocomplete from previous entries
  - [ ] 6.6 Implement tags field as multi-select from AllowedTags
  - [ ] 6.7 Implement categories field as multi-select from AllowedCategories
  - [ ] 6.8 Add visual feedback for focused fields (highlight, cursor)
  - [ ] 6.9 Implement Tab/Shift+Tab navigation between fields
  - [ ] 6.10 Add Save and Cancel buttons with clear visual indication
  - [ ] 6.11 Implement validation on field blur or submit
  - [ ] 6.12 Add undo/revert to original values functionality
  - [ ] 6.13 Display helpful error messages for validation failures
  - [ ] 6.14 Write comprehensive unit tests for editor interactions
  - [ ] 6.15 Test all validation scenarios and error messages

- [ ] 7.0 Integrate Metadata Editor with Navigation
  - [ ] 7.1 Add MetadataEditorScreen constant to `app.go`
  - [ ] 7.2 Add state management for editor screen
  - [ ] 7.3 Add navigation from metadata review screen (Enter key)
  - [ ] 7.4 Pass selected event to editor on navigation
  - [ ] 7.5 Implement save/cancel handling from editor
  - [ ] 7.6 Refresh metadata review list after successful save
  - [ ] 7.7 Return to metadata review on cancel
  - [ ] 7.8 Write app integration tests for editor flow

- [ ] 8.0 Enhance CLI Service with Metadata Operations
  - [ ] 8.1 Add `UpdateEventMetadata()` method to CLI service
  - [ ] 8.2 Accept event ID and metadata fields to update
  - [ ] 8.3 Call validation before persistence
  - [ ] 8.4 Return validation errors with helpful messages
  - [ ] 8.5 Update event in repository on success
  - [ ] 8.6 Write unit tests for metadata update operations

### Phase 4: Bulk Metadata Operations

- [ ] 9.0 Create Bulk Operations Model
  - [ ] 9.1 Implement `bulk_operations.go` as BubbleTea Model
  - [ ] 9.2 Add checkbox per event for selection
  - [ ] 9.3 Implement select all/none functionality
  - [ ] 9.4 Display selected event count
  - [ ] 9.5 Implement bulk field editing (company, project, tags, categories)
  - [ ] 9.6 Add "apply only if field empty" option for bulk operations
  - [ ] 9.7 Show preview of bulk changes before confirmation
  - [ ] 9.8 Implement keyboard shortcuts (Space for select, 'a' for all, 'd' for none)
  - [ ] 9.9 Add confirmation dialog before applying bulk changes
  - [ ] 9.10 Support undo/revert of bulk operations
  - [ ] 9.11 Write comprehensive unit tests for bulk operations
  - [ ] 9.12 Test edge cases (no events selected, partial application, etc.)

- [ ] 10.0 Integrate Bulk Operations with Metadata Review
  - [ ] 10.1 Add bulk operations mode toggle to metadata review screen
  - [ ] 10.2 Show checkboxes and bulk action buttons when in bulk mode
  - [ ] 10.3 Add navigation from metadata review to bulk operations
  - [ ] 10.4 Implement bulk operation confirmation and application
  - [ ] 10.5 Refresh metadata review after bulk operations complete
  - [ ] 10.6 Write app integration tests for bulk operations flow

- [ ] 11.0 Enhance CLI Service with Bulk Operations
  - [ ] 11.1 Add `BulkUpdateMetadata()` method to CLI service
  - [ ] 11.2 Accept event IDs and metadata updates
  - [ ] 11.3 Apply conditional updates (if field empty, etc.)
  - [ ] 11.4 Validate all events before any updates
  - [ ] 11.5 Support transaction-like behavior (all succeed or all fail)
  - [ ] 11.6 Return summary of applied changes
  - [ ] 11.7 Write unit tests for bulk update operations

### Phase 5: CSV Import & Manual Capture Integration

- [ ] 12.0 Enhance CSV Import Integration
  - [ ] 12.1 Modify `import_review.go` to show metadata review after import
  - [ ] 12.2 Display all imported events pre-loaded in metadata review
  - [ ] 12.3 Indicate which fields came from CSV vs. default values
  - [ ] 12.4 Show parsing issues or warnings for problematic imports
  - [ ] 12.5 Add duplicate detection status display
  - [ ] 12.6 Implement bulk operations for imported events
  - [ ] 12.7 Allow users to "Bulk Confirm" all metadata for imported events
  - [ ] 12.8 Write integration tests for import → metadata review flow

- [ ] 13.0 Enhance Manual Capture Integration
  - [ ] 13.1 Modify `form.go` to show quick metadata review after capture
  - [ ] 13.2 Display captured event with current metadata
  - [ ] 13.3 Offer option to add optional metadata (company, project, tags, categories)
  - [ ] 13.4 Allow editing before saving to database
  - [ ] 13.5 Provide option to "Add another event" or "Review all metadata"
  - [ ] 13.6 After N events (configurable), offer to review all at once
  - [ ] 13.7 Modify `success.go` to include metadata review option
  - [ ] 13.8 Write integration tests for capture → metadata review flow

### Phase 6: Testing & Validation

- [ ] 14.0 Comprehensive Testing Suite
  - [ ] 14.1 Write end-to-end tests for complete metadata review workflow
  - [ ] 14.2 Test metadata review → editor → save → review updated list
  - [ ] 14.3 Test bulk operations workflow (select → edit → preview → confirm)
  - [ ] 14.4 Test import → metadata review → bulk confirm workflow
  - [ ] 14.5 Test manual capture → metadata enrichment → metadata review
  - [ ] 14.6 Test all keyboard navigation shortcuts
  - [ ] 14.7 Test edge cases (empty lists, single items, large datasets)
  - [ ] 14.8 Verify all validation rules work correctly
  - [ ] 14.9 Test undo/revert functionality
  - [ ] 14.10 Run race detector: `go test -race ./...`
  - [ ] 14.11 Verify code coverage meets 80%+ threshold
  - [ ] 14.12 Performance test: metadata review loads in <500ms for 100 events

- [ ] 15.0 Documentation & User Guidance
  - [ ] 15.1 Update README.md with metadata review workflow description
  - [ ] 15.2 Update CLI_GUIDE.md with metadata review keyboard shortcuts
  - [ ] 15.3 Create examples of metadata editing workflows
  - [ ] 15.4 Document bulk operations with examples
  - [ ] 15.5 Update CHANGELOG.md with feature description
  - [ ] 15.6 Update troubleshooting guide with common metadata issues
  - [ ] 15.7 Document data quality scoring system for users

---

## Implementation Notes

### Architecture Decisions

1. **Separation of Concerns**:
   - Data quality calculation in service layer (`data_quality.go`)
   - Validation logic in validation package (`metadata_validator.go`)
   - UI components in models package (BubbleTea models)
   - Service adapter in CLI service layer

2. **Reuse Existing Patterns**:
   - Follow BubbleTea Model pattern from `form.go`, `list.go`, `details.go`
   - Use existing styling system from `internal/cli/styles/`
   - Leverage existing validation patterns
   - Build on existing service layer methods

3. **Quality Indicator Integration**:
   - Quality score calculated in domain/service layer
   - Visual indicator component in CLI models
   - Displayed in metadata review list and editor

### Dependencies

- Existing `career.CareerEvent` domain model
- Existing `CareerService` for event operations
- Existing `CLIEventService` for CLI-specific operations
- Existing BubbleTea components and styling
- Existing validation infrastructure

### Testing Strategy

- Unit tests for each component (validators, quality calculator, models)
- Integration tests for workflows (capture → metadata → review)
- End-to-end tests for complete user journeys
- Edge case testing for validation and data handling
- Performance testing for large event lists

### Success Criteria

1. ✓ Users can view events awaiting metadata clarification
2. ✓ Users can edit metadata for individual events
3. ✓ Users can perform bulk metadata operations
4. ✓ Data quality is visible with clear indicators
5. ✓ CSV import triggers metadata review
6. ✓ Manual capture offers metadata enrichment
7. ✓ All metadata changes are validated
8. ✓ Keyboard navigation supports efficient workflows
9. ✓ All changes persisted to database
10. ✓ Code coverage ≥ 80%
11. ✓ All tests passing (100% pass rate)

---

## Phase Dependencies

This feature builds on:
- **Phase 1**: Event capture (manual entry)
- **Phase 1.5**: CSV import (import_review.go)
- **Phase 2**: Event listing and filtering

This feature enables:
- **Phase 3**: Burst detection and grouping (requires clean metadata)
- **Phase 3**: Fact extraction (requires accurate categories)

---

## Estimated Effort

- **Phase 1** (Foundation): 3-4 hours
- **Phase 2** (Review Screen): 4-5 hours
- **Phase 3** (Editor): 4-5 hours
- **Phase 4** (Bulk Operations): 3-4 hours
- **Phase 5** (Integration): 3-4 hours
- **Phase 6** (Testing & Docs): 3-4 hours

**Total Estimated**: 20-26 hours

---

**Last Updated**: 2025-12-30
**Status**: Ready for Phase 1 Implementation

