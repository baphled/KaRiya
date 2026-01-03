# Legacy Screen-to-Intent Audit Document

**Date**: 2026-01-03
**Status**: Phase 1 - Preparation
**Purpose**: Document all legacy screens and their mapping to the new intent-based architecture

---

## Executive Summary

The current `app.go` (1,766 lines) contains 31 screen constants and manages 40+ model fields directly. This audit maps each legacy screen to its corresponding intent (5 existing + 5 new) to enable the aggressive replacement strategy.

**Current State**:
- 31 screen constants
- 40+ model fields in root Model struct
- 906 lines in Update() method
- 180 lines in View() method
- 30+ legacy model files
- Monolithic architecture

**Target State**:
- 10 intent-based screens
- 5-7 core model fields
- 50 lines in Update() method
- 30 lines in View() method
- 0 legacy model files (migrated to intents)
- Modular intent-driven architecture

---

## Screen Mapping

### Existing 5 Intents (Already Implemented)

#### 1. CaptureEvent Intent
**Status**: ✅ Implemented
**Screen Constants Handled**:
- `CaptureScreen` - Main capture form
- `ConfirmationScreen` - Confirmation dialog

**Legacy Model Files**:
- `internal/cli/models/form.go` - Form input handling (partially)
- `internal/cli/models/confirmation_dialog.go` - Confirmation UI

**Features**:
- Event capture with form validation
- Multiple input fields (title, description, date, etc.)
- Confirmation before saving
- Success feedback

**State Machine**:
- Initial → FormInput → Review → Confirm → Completed

**Test Coverage**: 30+ tests in `capture_event_views_test.go`

---

#### 2. BrowseTimeline Intent
**Status**: ✅ Implemented
**Screen Constants Handled**:
- `ListScreen` - Event list view
- `ViewScreen` - Event detail view
- `ActionMenuScreen` - Context menu for actions
- `CVListScreen` - CV browsing (partial)

**Legacy Model Files**:
- `internal/cli/models/list.go` - List rendering
- `internal/cli/models/details.go` - Detail view
- `internal/cli/models/action_menu.go` - Action menu
- `internal/cli/models/cv_list.go` - CV list (partial)

**Features**:
- Browse career events in timeline
- View event details
- Action menu for event operations
- Search and filter events
- Pagination for large lists

**State Machine**:
- List → View → ActionMenu → (back to List or execute action)

**Test Coverage**: 37 tests in `browse_timeline_test.go`

---

#### 3. GenerateCV Intent
**Status**: ✅ Implemented
**Screen Constants Handled**:
- `CVGeneratorScreen` - CV generation form
- `CVPreviewScreen` - CV preview
- `CVConfigManagerScreen` - Configuration

**Legacy Model Files**:
- `internal/cli/models/cv_generator.go` - CV generation form
- `internal/cli/models/cv_preview.go` - CV preview rendering
- `internal/cli/models/cv_config_manager.go` - CV configuration

**Features**:
- Configure CV generation parameters
- Select audience/template
- Preview generated CV
- Validate before export

**State Machine**:
- Config → Profile → Audience → Preview → Review → Confirm

**Test Coverage**: 41 tests in `generate_cv_test.go`

---

#### 4. ExportArtifact Intent
**Status**: ✅ Implemented
**Screen Constants Handled**:
- `CVExportDialogScreen` - Export configuration
- `CVExportProgressScreen` - Export progress
- `CVExportSuccessScreen` - Export success

**Legacy Model Files**:
- `internal/cli/models/cv_export_dialog.go` - Export options
- Progress tracking utilities

**Features**:
- Select export format (PDF, DOCX, TXT, JSON)
- Configure export options
- Track export progress
- Display success/error messages

**State Machine**:
- Select → Configure → Progress → Complete

**Test Coverage**: 400+ tests in `export_artifact_test.go`

---

#### 5. ConfigureSystem Intent
**Status**: ✅ Implemented
**Screen Constants Handled**:
- `MetadataReviewScreen` - Review system metadata (partial)
- System-wide configuration screens

**Legacy Model Files**:
- System configuration utilities

**Features**:
- Configure system settings
- Manage metadata
- Domain configuration
- Staged changes with review/confirm

**State Machine**:
- Review → Edit → StagedChanges → Confirm

**Test Coverage**: 400+ tests in `configure_system_test.go`

---

### New 5 Intents (To Be Implemented)

#### 6. BurstManagement Intent
**Status**: ❌ Not Implemented
**Screen Constants Handled**:
- `BurstListScreen` - Display existing bursts
- `BurstDetailsScreen` - View burst details
- `BurstEditorScreen` - Edit burst details
- `BurstSuggestionScreen` - AI suggestions for bursts

**Legacy Model Files**:
- `internal/cli/models/burst_list.go` - List of bursts
- `internal/cli/models/burst_details.go` - Detail view
- `internal/cli/models/burst_editor.go` - Burst editor form
- `internal/cli/models/burst_suggestion.go` - Burst suggestions

**Special Logic**:
- Burst creation/editing from events
- AI-powered burst suggestions
- Burst validation and constraints
- Relationship management with events

**Features**:
- List all bursts
- View burst details
- Create/edit bursts
- Generate burst suggestions
- Delete bursts
- Validate burst constraints

**State Machine**:
- List → View → (Edit | Suggest) → Confirm → Complete

**Dependencies**:
- CareerService for burst operations
- AI/ML service for suggestions
- Event service for relationships

**Test Coverage Target**: 30+ tests

**Edge Cases**:
- Empty burst list
- Invalid burst data
- Duplicate bursts
- Burst-event relationship constraints

---

#### 7. FactManagement Intent
**Status**: ❌ Not Implemented
**Screen Constants Handled**:
- `FactListScreen` - Display facts for events/bursts
- `FactDetailsScreen` - View individual fact details
- `FactEditorScreen` - Edit individual facts
- `FactsResultsScreen` - Review facts extracted after import
- `FactActionMenuScreen` - Context menu for fact operations

**Legacy Model Files**:
- `internal/cli/models/fact_list.go` - List of facts
- `internal/cli/models/fact_details.go` - Detail view
- `internal/cli/models/fact_editor.go` - Fact editor form
- `internal/cli/models/facts_results.go` - Import results
- `internal/cli/models/fact_card.go` - Fact card rendering

**Special Logic**:
- Fact extraction from events
- Quality assessment of facts
- Fact deduplication
- Relationship tracking (fact → event/burst)

**Features**:
- List facts for an event/burst
- View fact details
- Create/edit facts
- Review extracted facts
- Delete facts
- Search/filter facts
- Quality scoring

**State Machine**:
- List → View → (Edit | Review) → Confirm → Complete

**Dependencies**:
- CareerService for fact operations
- Event/Burst service for relationships
- Quality assessment service

**Test Coverage Target**: 40+ tests

**Edge Cases**:
- Empty fact list
- Invalid fact data
- Duplicate facts
- Quality assessment failures
- Orphaned facts (event deleted)

---

#### 8. ImportWizard Intent
**Status**: ❌ Not Implemented
**Screen Constants Handled**:
- `ImportReviewScreen` - Review import data
- `ImportProgressScreen` - Import progress tracking

**Legacy Model Files**:
- `internal/cli/models/import_review.go` - Import review screen

**Special Logic**:
- File selection and validation
- CSV parsing and preview
- Conflict resolution
- Progress tracking
- Error handling and recovery

**Features**:
- Select file to import
- Preview import data
- Validate file format
- Handle conflicts/duplicates
- Track import progress
- Show import results

**State Machine**:
- SelectFile → Review → Progress → Complete → Confirm

**Dependencies**:
- ImportService for file operations
- CareerService for data persistence
- File system access

**Test Coverage Target**: 25+ tests

**Edge Cases**:
- Invalid file format
- Empty file
- Duplicate data
- Missing required fields
- File access errors
- Large files (performance)

---

#### 9. MetadataEditor Intent
**Status**: ❌ Not Implemented
**Screen Constants Handled**:
- `MetadataReviewScreen` - Review metadata
- `MetadataEditorScreen` - Edit metadata

**Legacy Model Files**:
- `internal/cli/models/metadata_review.go` - Review screen
- `internal/cli/models/metadata_editor.go` - Editor screen

**Special Logic**:
- Metadata validation
- Change tracking for audit
- Batch metadata updates
- Relationship updates

**Features**:
- Review event/burst metadata
- Edit metadata fields
- Track changes for audit
- Validate metadata constraints
- Bulk metadata updates

**State Machine**:
- Review → Edit → Confirm → Complete

**Dependencies**:
- CareerService for metadata operations
- Validation service

**Test Coverage Target**: 20+ tests

**Edge Cases**:
- Invalid metadata
- Constraint violations
- Concurrent updates
- Orphaned relationships

---

#### 10. BulkOperations Intent
**Status**: ❌ Not Implemented
**Screen Constants Handled**:
- `BulkOperationsScreen` - Bulk operation selection and execution

**Legacy Model Files**:
- `internal/cli/models/bulk_operations.go` - Bulk operations

**Special Logic**:
- Operation selection
- Parameter configuration
- Batch processing
- Progress tracking
- Error handling and rollback

**Features**:
- Select bulk operation (delete, export, tag, etc.)
- Configure operation parameters
- Preview affected items
- Execute operation with progress
- Show results and errors

**State Machine**:
- SelectOp → Configure → Execute → Confirm → Complete

**Dependencies**:
- CareerService for bulk operations
- Progress tracking service
- Error handling

**Test Coverage Target**: 20+ tests

**Edge Cases**:
- Empty selection
- Invalid parameters
- Partial failure
- Rollback on error
- Large batch operations

---

## Screen Constant Inventory

### Total: 31 Screen Constants

| Screen Constant | Intent | Status | Notes |
|---|---|---|---|
| HomeScreen | Root/Menu | Keep | Main menu screen |
| MainMenuScreen | Root/Menu | Keep | Menu selection |
| CaptureScreen | CaptureEvent | Migrate | ✅ Ready |
| ListScreen | BrowseTimeline | Migrate | ✅ Ready |
| ViewScreen | BrowseTimeline | Migrate | ✅ Ready |
| QuitScreen | Root | Keep | Quit confirmation |
| SuccessScreen | Root | Keep | Generic success |
| ActionMenuScreen | BrowseTimeline | Migrate | ✅ Ready |
| ConfirmationScreen | CaptureEvent | Migrate | ✅ Ready |
| ImportReviewScreen | ImportWizard | Migrate | ❌ To implement |
| ImportProgressScreen | ImportWizard | Migrate | ❌ To implement |
| MetadataReviewScreen | MetadataEditor | Migrate | ❌ To implement |
| MetadataEditorScreen | MetadataEditor | Migrate | ❌ To implement |
| BulkOperationsScreen | BulkOperations | Migrate | ❌ To implement |
| BurstSuggestionScreen | BurstManagement | Migrate | ❌ To implement |
| BurstListScreen | BurstManagement | Migrate | ❌ To implement |
| FactsResultsScreen | FactManagement | Migrate | ❌ To implement |
| HelpScreen | Root | Keep | Help display |
| FactListScreen | FactManagement | Migrate | ❌ To implement |
| FactActionMenuScreen | FactManagement | Migrate | ❌ To implement |
| FactDetailsScreen | FactManagement | Migrate | ❌ To implement |
| FactEditorScreen | FactManagement | Migrate | ❌ To implement |
| BurstDetailsScreen | BurstManagement | Migrate | ❌ To implement |
| BurstEditorScreen | BurstManagement | Migrate | ❌ To implement |
| CVConfigManagerScreen | GenerateCV | Migrate | ✅ Ready |
| CVGeneratorScreen | GenerateCV | Migrate | ✅ Ready |
| CVPreviewScreen | GenerateCV | Migrate | ✅ Ready |
| CVListScreen | BrowseTimeline | Migrate | ✅ Ready (partial) |
| CVExportDialogScreen | ExportArtifact | Migrate | ✅ Ready |
| CVExportSuccessScreen | ExportArtifact | Migrate | ✅ Ready |
| CVExportProgressScreen | ExportArtifact | Migrate | ✅ Ready |

---

## Model Field Inventory

### Current app.go Model Fields (40+ fields)

| Field | Type | Intent | Status |
|---|---|---|---|
| cliService | *service.CLIEventService | Multiple | Keep/Inject |
| logger | *logger.Logger | Root | Keep |
| service | *careerservice.Service | Multiple | Keep/Inject |
| currentScreen | Screen | Root | Keep (simplified) |
| previousScreen | Screen | Root | Keep (simplified) |
| screenBeforeActionMenu | Screen | Root | Keep (simplified) |
| breadcrumbs | []string | Root | Remove |
| width | int | Root | Keep |
| height | int | Root | Keep |
| workflowState | *workflow.WorkflowState | Root | Remove |
| formModel | *models.FormModel | CaptureEvent | Move to intent |
| successModel | *models.SuccessModel | Root | Keep (generic) |
| listModel | *models.ListModel | BrowseTimeline | Move to intent |
| detailsModel | *models.DetailsModel | BrowseTimeline | Move to intent |
| actionMenuModel | *models.ActionMenuModel | BrowseTimeline | Move to intent |
| factActionMenuModel | *models.FactActionMenuModel | FactManagement | Move to intent |
| confirmationDialog | *models.ConfirmationDialog | CaptureEvent | Move to intent |
| deleteEventID | string | BrowseTimeline | Move to intent |
| importService | *importer.ImportService | ImportWizard | Move to intent |
| importReviewModel | *models.ImportReviewModel | ImportWizard | Move to intent |
| importProgressModel | *models.ImportProgressModel | ImportWizard | Move to intent |
| importFilePath | string | ImportWizard | Move to intent |
| importResult | *importer.ImportResult | ImportWizard | Move to intent |
| metadataReviewModel | *models.MetadataReviewModel | MetadataEditor | Move to intent |
| metadataEditorModel | *models.MetadataEditorModel | MetadataEditor | Move to intent |
| bulkOperationsModel | *models.BulkOperationsModel | BulkOperations | Move to intent |
| burstSuggestionModel | *models.BurstSuggestionModel | BurstManagement | Move to intent |
| burstListModel | *models.BurstListModel | BurstManagement | Move to intent |
| burstDetailsModel | *models.BurstDetailsModel | BurstManagement | Move to intent |
| burstEditorModel | *models.BurstEditorModel | BurstManagement | Move to intent |
| factListModel | *models.FactListModel | FactManagement | Move to intent |
| factsResultsModel | *models.FactsResultsModel | FactManagement | Move to intent |
| factEditorModel | *models.FactEditorModel | FactManagement | Move to intent |
| factDetailsModel | *models.FactDetailsModel | FactManagement | Move to intent |
| menuModel | *models.MenuModel | Root | Keep (simplified) |
| helpModel | *models.HelpModel | Root | Keep |
| cvConfigManagerModel | *models.CVConfigManagerModel | GenerateCV | Move to intent |
| cvGeneratorModel | *models.CVGeneratorModel | GenerateCV | Move to intent |
| cvPreviewModel | *models.CVPreviewModel | GenerateCV | Move to intent |
| cvListModel | *models.CVListModel | BrowseTimeline | Move to intent |
| cvExportDialogModel | *models.CVExportDialogModel | ExportArtifact | Move to intent |
| cvExportSuccessModel | *models.CVExportSuccessModel | ExportArtifact | Move to intent |

**Summary**:
- Keep in root: 7 fields (logger, service, currentScreen, width, height, menuModel, helpModel)
- Move to intents: 33 fields
- Remove: 5 fields

---

## Cross-Model Dependencies

### Shared Utilities
- `FormModel` - Used by CaptureEvent, BurstManagement, FactManagement, MetadataEditor
- `ListModel` - Used by BrowseTimeline, BurstManagement, FactManagement
- `DetailsModel` - Used by BrowseTimeline, FactManagement
- `ActionMenuModel` - Used by BrowseTimeline, FactManagement
- `ConfirmationDialog` - Used by multiple intents
- `SuccessModel` - Used by multiple intents

### Service Dependencies
- `CareerService` - Used by all intents
- `CLIEventService` - Used by multiple intents
- `ImportService` - Used by ImportWizard
- `Logger` - Used globally

### Styling Dependencies
- `internal/cli/styles/styles.go` - Centralized styling
- Consistent color scheme across all intents
- Uniform border styles
- Standardized spacing

---

## Reusable Components

### Existing Components (in `internal/cli/components/`)
- Card - Display data in card format
- List - Render scrollable lists
- Form - Handle form input
- Modal - Display modal dialogs
- Progress - Show progress indicators
- Menu - Display menu options

### Components Needed for New Intents
- BurstCard - Display burst information
- FactCard - Display fact information (exists: `fact_card.go`)
- ImportPreview - Preview import data
- ProgressBar - Enhanced progress display
- BulkOperationPreview - Preview bulk operation results

---

## Custom Styling & UI Patterns

### CaptureEvent
- Form with validation feedback
- Field-level error messages
- Focus indicators for accessibility
- Tab navigation between fields

### BrowseTimeline
- Scrollable list with highlight
- Detail panel with rich formatting
- Action menu positioning
- Breadcrumb navigation (to be removed)

### GenerateCV
- Multi-step wizard UI
- Preview with syntax highlighting
- Configuration panel
- Export progress indicator

### ExportArtifact
- Format selection with icons
- Configuration options per format
- Real-time progress updates
- Success/error feedback

### BurstManagement
- Card-based list layout
- Suggestion cards with AI indicators
- Editor form with validation
- Relationship visualization

### FactManagement
- Fact cards with quality indicators
- Searchable list view
- Detailed fact editor
- Import results table

### ImportWizard
- File picker UI
- Data preview table
- Progress bar with ETA
- Conflict resolution dialog

### MetadataEditor
- Two-column layout (original/edited)
- Change highlighting
- Validation feedback
- Batch editor UI

### BulkOperations
- Operation selection grid
- Configuration form
- Progress with item count
- Results summary table

---

## Implementation Priority

### Phase 1: Preparation ✅ (This Document)
- [x] Audit legacy screens
- [x] Map to intents
- [ ] Document cross-dependencies
- [ ] Identify reusable components

### Phase 2: New Intent Implementation (64 hours)
Priority order by complexity and dependencies:
1. **BurstManagement** (16h) - Medium complexity, moderate dependencies
2. **FactManagement** (16h) - Medium-high complexity, good component reuse
3. **ImportWizard** (12h) - Medium complexity, file system operations
4. **MetadataEditor** (10h) - Low-medium complexity, simple state machine
5. **BulkOperations** (10h) - Medium complexity, async operations

### Phase 3: Root app.go Rebuild (14 hours)
- Minimal model with intent router
- Menu-based intent selection
- Result handling and navigation
- Global shortcuts

### Phase 4: Testing & Validation (34 hours)
- Unit tests for all 5 new intents
- Integration tests
- End-to-end workflows
- Performance benchmarks

---

## Success Criteria

### Code Reduction
- [ ] app.go: 1,766 → < 250 lines (85% reduction)
- [ ] Update(): 906 → < 50 lines (95% reduction)
- [ ] View(): 180 → < 30 lines (83% reduction)
- [ ] Model fields: 40+ → 7 (82% reduction)
- [ ] Screen constants: 31 → 0 (100% removal)

### Feature Preservation
- [ ] All 31 screens migrated to intents
- [ ] All CRUD operations working
- [ ] All menu items functional
- [ ] All workflows completable
- [ ] All data operations unchanged

### Quality Metrics
- [ ] 164+ tests passing
- [ ] 87%+ code coverage
- [ ] 0 race conditions
- [ ] All lint checks passing
- [ ] All performance benchmarks met

---

## Next Steps

1. **Task 1.1.2**: Map each legacy screen to intent (using this document)
2. **Task 1.1.3**: Document special logic for each screen
3. **Task 1.1.4**: Identify cross-model dependencies
4. **Task 1.1.5**: Document custom styling patterns
5. **Task 1.2.x**: Review features for each new intent
6. **Task 1.3.x**: Plan new intent implementations
7. **Task 1.4.x**: Verify intent framework readiness

---

**Document Version**: 1.0
**Last Updated**: 2026-01-03
**Status**: Complete - Ready for Phase 1.1.2

