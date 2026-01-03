# New Intent Implementation Plan

**Date**: 2026-01-03
**Status**: Phase 1.3 - Planning
**Purpose**: Define state machines, context/result structures, and component specifications for 5 new intents

---

## 1. BurstManagement Intent

### 1.1 State Machine Definition

```
Initial State: BurstListState

States:
├── BurstListState (List all bursts)
│   └── Actions: View, Edit, Delete, Suggest, Refresh
│       ├── View → BurstViewState
│       ├── Edit → BurstEditorState
│       ├── Delete → BurstDeleteConfirmState
│       ├── Suggest → BurstSuggestState
│       └── Refresh → BurstListState (reload)
│
├── BurstViewState (View burst details)
│   └── Actions: Edit, Delete, Back
│       ├── Edit → BurstEditorState
│       ├── Delete → BurstDeleteConfirmState
│       └── Back → BurstListState
│
├── BurstEditorState (Edit burst)
│   └── Actions: Save, Cancel
│       ├── Save → BurstListState (with success feedback)
│       └── Cancel → BurstViewState (if editing) or BurstListState (if new)
│
├── BurstDeleteConfirmState (Confirm deletion)
│   └── Actions: Confirm, Cancel
│       ├── Confirm → BurstListState (with success feedback)
│       └── Cancel → BurstViewState (if viewing) or BurstListState
│
└── BurstSuggestState (AI suggestions)
    └── Actions: Accept, Reject, Customize
        ├── Accept → BurstListState (with new burst created)
        ├── Reject → BurstListState
        └── Customize → BurstEditorState → BurstListState

Final State: IntentResult returned to root app
```

### 1.2 Context Structure

```go
type BurstManagementContext struct {
    // Current state
    CurrentState BurstManagementState

    // Burst data
    Bursts []*career.Burst
    SelectedBurst *career.Burst
    SelectedBurstIndex int

    // Filter/Search
    FilterCompetency string
    SearchText string
    SortBy string        // "name", "startDate", "endDate", "eventCount"
    SortOrder string     // "asc", "desc"

    // Pagination
    CurrentPage int
    PageSize int        // 20
    TotalBursts int

    // Editor
    EditingBurst *career.Burst
    FormErrors map[string]string

    // Suggestions
    Suggestions []*BurstSuggestion
    SelectedSuggestionIndex int

    // UI state
    ScrollPosition int
    ExpandedRows map[int]bool

    // Services
    Service *careerservice.Service

    // Metadata
    PreviousState BurstManagementState
    ScrollRestoreNeeded bool
}

type BurstSuggestion struct {
    Title string
    Events []*career.Event
    ConfidenceScore float64
    RecommendedStartDate time.Time
    RecommendedEndDate time.Time
    Skills []string
}
```

### 1.3 Result Structure

```go
type BurstManagementResult struct {
    Action string  // "created", "updated", "deleted", "none"
    Burst *career.Burst
    Bursts []*career.Burst
    Error error
    Message string
}
```

### 1.4 Reusable Components

| Component | Purpose | Customization |
|-----------|---------|---|
| ListModel | Display burst list | Item: BurstCard, Sort: 4 options, Filter: competency |
| DetailsModel | Display burst details | Content: burst info, Related: events in burst |
| FormModel | Burst editor | Fields: name, dates, description, skills, tags |
| ConfirmationDialog | Delete confirmation | Message: "Delete burst X?" |
| CardModel | Suggestion display | Show: title, events, confidence score |

### 1.5 Test Coverage Plan

**30+ tests covering:**

| Category | Tests | Examples |
|----------|-------|----------|
| List Operations | 8 | Pagination, sorting, filtering, search |
| View Operations | 4 | Display burst, display details, empty state |
| Edit Operations | 6 | Form validation, save, cancel, unsaved changes |
| Delete Operations | 4 | Confirmation, cascade handling, errors |
| Suggestion Operations | 5 | Generation, acceptance, customization |
| Edge Cases | 3 | Empty list, invalid data, concurrent edits |

### 1.6 Async Operations & Dependencies

| Operation | Type | Async | Dependencies |
|-----------|------|-------|---|
| Load bursts | Read | Yes | CareerService.GetBursts() |
| Create burst | Write | No | CareerService.CreateBurst() |
| Update burst | Write | No | CareerService.UpdateBurst() |
| Delete burst | Write | No | CareerService.DeleteBurst() |
| Generate suggestions | Compute | Yes | AI/ML service (TBD) |
| Validate burst | Validate | No | Business rules |

---

## 2. FactManagement Intent

### 2.1 State Machine Definition

```
Initial State: FactListState

States:
├── FactListState (List all facts)
│   └── Actions: View, Edit, Delete, Search, Filter, Sort
│       ├── View → FactViewState
│       ├── Edit → FactEditorState
│       ├── Delete → FactDeleteConfirmState
│       └── Search/Filter/Sort → FactListState (reload)
│
├── FactViewState (View fact details)
│   └── Actions: Edit, Delete, Back
│       ├── Edit → FactEditorState
│       ├── Delete → FactDeleteConfirmState
│       └── Back → FactListState
│
├── FactEditorState (Edit fact)
│   └── Actions: Save, Cancel
│       ├── Save → FactListState (with success feedback)
│       └── Cancel → FactViewState (if editing) or FactListState (if new)
│
├── FactDeleteConfirmState (Confirm deletion)
│   └── Actions: Confirm, Cancel
│       ├── Confirm → FactListState (with success feedback)
│       └── Cancel → FactViewState or FactListState
│
└── FactsResultsState (Import results review)
    └── Actions: Accept, Reject, Edit, AcceptAll, RejectAll
        ├── Accept → FactListState (add fact)
        ├── Reject → FactListState (skip fact)
        ├── Edit → FactEditorState → FactsResultsState
        ├── AcceptAll → FactListState (add all)
        └── RejectAll → FactListState (skip all)

Final State: IntentResult returned to root app
```

### 2.2 Context Structure

```go
type FactManagementContext struct {
    // Current state
    CurrentState FactManagementState

    // Fact data
    Facts []*career.Fact
    SelectedFact *career.Fact
    SelectedFactIndex int

    // Filter/Search
    SearchText string
    MinQualityScore float64  // 0-1
    MaxQualityScore float64
    FilterBySource string     // "event", "burst", "all"
    FilterSourceID string
    SortBy string             // "date", "quality", "relevance"
    SortOrder string          // "asc", "desc"

    // Pagination
    CurrentPage int
    PageSize int              // 20
    TotalFacts int

    // Editor
    EditingFact *career.Fact
    FormErrors map[string]string
    QualityScore float64      // Calculated as user edits

    // Import results
    ImportedFacts []*career.Fact
    FactActions map[string]string  // factID → "accept"/"reject"

    // UI state
    ScrollPosition int
    ExpandedRows map[int]bool

    // Services
    Service *careerservice.Service

    // Metadata
    PreviousState FactManagementState
    ScrollRestoreNeeded bool
}
```

### 2.3 Result Structure

```go
type FactManagementResult struct {
    Action string  // "created", "updated", "deleted", "imported", "none"
    Fact *career.Fact
    Facts []*career.Fact
    ImportedCount int
    SkippedCount int
    Error error
    Message string
}
```

### 2.4 Reusable Components

| Component | Purpose | Customization |
|-----------|---------|---|
| ListModel | Display fact list | Item: FactCard, Sort: 3 options, Filter: quality/source |
| DetailsModel | Display fact details | Content: fact info, Related: similar facts |
| FormModel | Fact editor | Fields: content, quality, tags, source |
| TableModel | Import results | Columns: content, quality, action, errors |
| CardModel | Fact card display | Show: content, quality, source |

### 2.5 Test Coverage Plan

**40+ tests covering:**

| Category | Tests | Examples |
|----------|-------|----------|
| List Operations | 10 | Pagination, sorting, filtering, search, quality ranges |
| View Operations | 4 | Display fact, display details, empty state |
| Edit Operations | 8 | Form validation, quality calc, save, cancel |
| Delete Operations | 4 | Confirmation, cascade handling, errors |
| Import Results | 8 | Accept, reject, bulk actions, validation |
| Quality Scoring | 4 | Multi-criteria scoring, edge cases |
| Edge Cases | 2 | Empty list, invalid data |

### 2.6 Async Operations & Dependencies

| Operation | Type | Async | Dependencies |
|-----------|------|-------|---|
| Load facts | Read | Yes | CareerService.GetFacts() |
| Create fact | Write | No | CareerService.CreateFact() |
| Update fact | Write | No | CareerService.UpdateFact() |
| Delete fact | Write | No | CareerService.DeleteFact() |
| Extract facts | Compute | Yes | NLP/ML service (TBD) |
| Calculate quality | Compute | No | Quality scoring algorithm |
| Detect duplicates | Compute | No | Fuzzy matching |

---

## 3. ImportWizard Intent

### 3.1 State Machine Definition

```
Initial State: FileSelectionState

States:
├── FileSelectionState (Select CSV file)
│   └── Actions: Select, Cancel
│       ├── Select → DataPreviewState
│       └── Cancel → IntentResult (none)
│
├── DataPreviewState (Preview CSV data)
│   └── Actions: Proceed, Back, Cancel
│       ├── Proceed → ConflictResolutionState
│       ├── Back → FileSelectionState
│       └── Cancel → IntentResult (none)
│
├── ConflictResolutionState (Handle conflicts)
│   └── Actions: Proceed, Back, Cancel
│       ├── Proceed → ImportProgressState
│       ├── Back → DataPreviewState
│       └── Cancel → IntentResult (none)
│
├── ImportProgressState (Track import progress)
│   └── Actions: Pause, Resume, Cancel
│       ├── Complete → ImportResultsState
│       ├── Pause → ImportProgressState (paused)
│       ├── Resume → ImportProgressState (resumed)
│       └── Cancel → IntentResult (partial import)
│
└── ImportResultsState (Show import results)
    └── Actions: Done, ViewErrors
        ├── Done → IntentResult (completed)
        └── ViewErrors → ImportResultsState (show details)

Final State: IntentResult returned to root app
```

### 3.2 Context Structure

```go
type ImportWizardContext struct {
    // Current state
    CurrentState ImportWizardState

    // File data
    FilePath string
    FileSize int64
    FileModifiedDate time.Time
    FileEncoding string  // "utf-8", "iso-8859-1", etc.

    // CSV data
    Headers []string
    PreviewRows [][]string  // First 20 rows
    TotalRows int

    // Validation
    ValidationErrors map[int][]string  // rowIndex → errors
    ValidRowCount int
    ErrorRowCount int

    // Conflict resolution
    Conflicts []*ImportConflict
    ConflictResolutions map[string]string  // conflictID → resolution
    MergeStrategy string  // "skip", "replace", "merge"

    // Import progress
    ImportedCount int
    SkippedCount int
    ErrorCount int
    CurrentRowIndex int
    IsImporting bool
    IsPaused bool
    ImportStartTime time.Time

    // Results
    ImportedEvents []*career.Event
    SkippedEvents []SkippedEventInfo
    ErroredEvents []ErroredEventInfo

    // UI state
    ScrollPosition int
    SelectedConflictIndex int

    // Services
    Service *careerservice.Service
    ImportService *importer.ImportService

    // Metadata
    PreviousState ImportWizardState
}

type ImportConflict struct {
    ID string
    RowIndex int
    EventTitle string
    ExistingEvent *career.Event
    NewEvent *career.Event
    ConflictType string  // "duplicate", "merge", etc.
}

type SkippedEventInfo struct {
    RowIndex int
    Reason string
}

type ErroredEventInfo struct {
    RowIndex int
    Error string
    Field string
}
```

### 3.3 Result Structure

```go
type ImportWizardResult struct {
    Action string  // "imported", "cancelled", "paused"
    ImportedCount int
    SkippedCount int
    ErrorCount int
    ImportedEvents []*career.Event
    Errors []string
    ErrorDetails map[int]string  // rowIndex → error message
}
```

### 3.4 Reusable Components

| Component | Purpose | Customization |
|-----------|---------|---|
| FilePicker | Select CSV file | Accept: .csv only |
| TableModel | Preview/results | Columns: dynamic from CSV |
| ProgressModel | Import progress | Show: rows/total, ETA, pause/resume |
| FormModel | Conflict resolution | Fields: merge strategy options |
| ConfirmationDialog | Conflict handling | Message: conflict details |

### 3.5 Test Coverage Plan

**25+ tests covering:**

| Category | Tests | Examples |
|----------|-------|----------|
| File Selection | 4 | Valid file, invalid file, file not found |
| Data Preview | 4 | Parse CSV, detect headers, show sample |
| Conflict Detection | 4 | Detect duplicates, show conflicts |
| Import Progress | 6 | Track progress, pause/resume, cancel |
| Error Handling | 4 | Parse errors, validation errors, save errors |
| Edge Cases | 3 | Empty file, large file, special characters |

### 3.6 Async Operations & Dependencies

| Operation | Type | Async | Dependencies |
|-----------|------|-------|---|
| Load file | Read | Yes | File system |
| Parse CSV | Compute | Yes | CSV parser |
| Validate data | Validate | No | Validation rules |
| Detect conflicts | Compute | Yes | Duplicate detection |
| Import data | Write | Yes | CareerService.ImportEvents() |
| Track progress | Read | Yes | Progress callbacks |

---

## 4. MetadataEditor Intent

### 4.1 State Machine Definition

```
Initial State: MetadataReviewState

States:
├── MetadataReviewState (Review metadata)
│   └── Actions: Edit, Back
│       ├── Edit → MetadataEditState
│       └── Back → IntentResult (none)
│
├── MetadataEditState (Edit metadata)
│   └── Actions: Save, Reset, Cancel
│       ├── Save → MetadataReviewState (with success feedback)
│       ├── Reset → MetadataEditState (revert to original)
│       └── Cancel → MetadataReviewState

Final State: IntentResult returned to root app
```

### 4.2 Context Structure

```go
type MetadataEditorContext struct {
    // Current state
    CurrentState MetadataEditorState

    // Target entity
    EntityType string  // "event", "burst", "fact"
    EntityID string
    Entity interface{}

    // Metadata
    OriginalMetadata map[string]interface{}
    EditedMetadata map[string]interface{}

    // Changes tracking
    ChangedFields map[string]bool
    FieldValues map[string]interface{}  // current edit values
    FormErrors map[string]string

    // Metadata categories
    EditableFields []MetadataField
    ReadOnlyFields []MetadataField
    ComputedFields []MetadataField

    // UI state
    SelectedFieldIndex int
    ScrollPosition int

    // Services
    Service *careerservice.Service

    // Metadata
    PreviousState MetadataEditorState
    SaveAttempted bool
}

type MetadataField struct {
    Name string
    Type string  // "string", "int", "date", "float", etc.
    Value interface{}
    IsEditable bool
    IsRequired bool
    ValidationRules []string
    HelpText string
    Options []string  // for select fields
}
```

### 4.3 Result Structure

```go
type MetadataEditorResult struct {
    Action string  // "saved", "cancelled"
    EntityType string
    EntityID string
    Changes map[string]interface{}
    OriginalValues map[string]interface{}
    Error error
    Message string
}
```

### 4.4 Reusable Components

| Component | Purpose | Customization |
|-----------|---------|---|
| FormModel | Metadata editor | Fields: dynamic per entity type |
| DetailsModel | Metadata review | Content: original values (read-only) |
| TableModel | Change summary | Columns: field, old value, new value |
| ConfirmationDialog | Save confirmation | Message: show changed fields |

### 4.5 Test Coverage Plan

**20+ tests covering:**

| Category | Tests | Examples |
|----------|-------|----------|
| Review Operations | 4 | Display metadata, categories, read-only |
| Edit Operations | 6 | Form validation, field updates, changes |
| Save Operations | 4 | Confirmation, success, error handling |
| Change Tracking | 4 | Track changes, highlight, rollback |
| Edge Cases | 2 | Invalid metadata, constraint violations |

### 4.6 Async Operations & Dependencies

| Operation | Type | Async | Dependencies |
|-----------|------|-------|---|
| Load metadata | Read | No | CareerService.GetMetadata() |
| Validate metadata | Validate | No | Validation rules |
| Save metadata | Write | No | CareerService.SaveMetadata() |
| Track changes | Compute | No | Change tracking logic |

---

## 5. BulkOperations Intent

### 5.1 State Machine Definition

```
Initial State: OperationSelectionState

States:
├── OperationSelectionState (Select operation)
│   └── Actions: Select, Cancel
│       ├── Select → OperationConfigState
│       └── Cancel → IntentResult (none)
│
├── OperationConfigState (Configure operation)
│   └── Actions: Proceed, Back, Cancel
│       ├── Proceed → ExecutionConfirmState
│       ├── Back → OperationSelectionState
│       └── Cancel → IntentResult (none)
│
├── ExecutionConfirmState (Confirm before execution)
│   └── Actions: Execute, Back, Cancel
│       ├── Execute → OperationProgressState
│       ├── Back → OperationConfigState
│       └── Cancel → IntentResult (none)
│
├── OperationProgressState (Track operation progress)
│   └── Actions: Pause, Resume, Cancel
│       ├── Complete → OperationResultsState
│       ├── Pause → OperationProgressState (paused)
│       ├── Resume → OperationProgressState (resumed)
│       └── Cancel → OperationResultsState (partial)
│
└── OperationResultsState (Show operation results)
    └── Actions: Done, ViewDetails
        ├── Done → IntentResult (completed)
        └── ViewDetails → OperationResultsState (expand errors)

Final State: IntentResult returned to root app
```

### 5.2 Context Structure

```go
type BulkOperationsContext struct {
    // Current state
    CurrentState BulkOperationsState

    // Operation selection
    AvailableOperations []BulkOperation
    SelectedOperation BulkOperation

    // Configuration
    OperationConfig map[string]interface{}
    ConfigErrors map[string]string

    // Scope selection
    ScopeType string  // "selected", "filtered", "all"
    SelectedItems []string  // item IDs
    FilterCriteria map[string]interface{}

    // Preview
    AffectedItemCount int
    AffectedItems []interface{}
    EstimatedDuration time.Duration

    // Execution
    IsExecuting bool
    IsPaused bool
    ExecutionStartTime time.Time
    ProcessedCount int
    SuccessCount int
    FailureCount int
    SkippedCount int
    CurrentItemIndex int
    CurrentItemID string

    // Results
    Results map[string]OperationResult  // itemID → result
    Errors []OperationError

    // UI state
    ScrollPosition int
    SelectedResultIndex int

    // Services
    Service *careerservice.Service

    // Metadata
    PreviousState BulkOperationsState
}

type BulkOperation struct {
    ID string
    Name string
    Description string
    Icon string
    EstimatedTime time.Duration
    ConfigFields []ConfigField
    SupportedScopes []string  // "selected", "filtered", "all"
}

type ConfigField struct {
    Name string
    Type string  // "string", "int", "select", "bool"
    Label string
    Required bool
    Options []string  // for select fields
    DefaultValue interface{}
    HelpText string
}

type OperationResult struct {
    ItemID string
    Status string  // "success", "error", "skipped"
    Message string
    Data interface{}
}

type OperationError struct {
    ItemID string
    Error string
    Details string
}
```

### 5.3 Result Structure

```go
type BulkOperationsResult struct {
    Operation string
    ProcessedCount int
    SuccessCount int
    FailureCount int
    SkippedCount int
    Results map[string]OperationResult
    Errors []OperationError
    Error error
    Message string
}
```

### 5.4 Reusable Components

| Component | Purpose | Customization |
|-----------|---------|---|
| GridModel | Operation selection | Items: operation cards |
| FormModel | Operation config | Fields: operation-specific |
| TableModel | Item list & results | Columns: dynamic per operation |
| ProgressModel | Operation progress | Show: items/total, pause/resume |
| ConfirmationDialog | Execution confirm | Message: affected items count |

### 5.5 Test Coverage Plan

**20+ tests covering:**

| Category | Tests | Examples |
|----------|-------|----------|
| Operation Selection | 3 | Select operation, show details |
| Configuration | 4 | Validate config, preview items |
| Execution | 6 | Execute, pause/resume, cancel |
| Results | 4 | Success, errors, partial failure |
| Edge Cases | 3 | Empty selection, large batch, errors |

### 5.6 Async Operations & Dependencies

| Operation | Type | Async | Dependencies |
|-----------|------|-------|---|
| Load operations | Read | No | Operation definitions |
| Preview items | Compute | No | Filter/query logic |
| Execute operation | Write | Yes | Operation handlers |
| Track progress | Read | Yes | Progress callbacks |

---

## Summary Table

| Intent | States | Context Fields | Components | Tests | Async Ops |
|--------|--------|---|---|---|---|
| BurstManagement | 5 | 18 | 5 | 30+ | 2 |
| FactManagement | 5 | 22 | 5 | 40+ | 3 |
| ImportWizard | 5 | 24 | 4 | 25+ | 4 |
| MetadataEditor | 2 | 14 | 4 | 20+ | 1 |
| BulkOperations | 5 | 20 | 5 | 20+ | 2 |
| **TOTAL** | **22** | **98** | **23** | **135+** | **12** |

---

## Component Reuse Analysis

### Shared Components Across Intents

| Component | Used By | Count |
|-----------|---------|-------|
| ListModel | BurstMgmt, FactMgmt, BulkOps | 3 |
| FormModel | BurstMgmt, FactMgmt, ImportWizard, MetadataEditor, BulkOps | 5 |
| DetailsModel | BurstMgmt, FactMgmt, MetadataEditor | 3 |
| TableModel | FactMgmt, ImportWizard, BulkOps | 3 |
| ProgressModel | ImportWizard, BulkOps | 2 |
| ConfirmationDialog | All 5 intents | 5 |
| CardModel | BurstMgmt, FactMgmt | 2 |

**Reuse Efficiency**: 80% of components are shared across multiple intents
**Custom Components Needed**: 5 (FactCard, BurstCard, FilePicker, GridModel, ConfigField)

---

## Framework Readiness Checklist

### IntentRouter ✅
- [x] RegisterIntent() method
- [x] RegisterResultHandler() method
- [x] ActivateIntent() method
- [x] GetCurrentIntent() method
- [x] HandleIntentResult() method
- [x] Navigation support

### IntentResult[T] ✅
- [x] Type-safe generic implementation
- [x] Status field (pending, success, error)
- [x] Data field (generic)
- [x] Error field (for errors)
- [x] Metadata support (for context preservation)
- [x] Helper methods (WithMetadata, GetMetadata)

### Testing Utilities ✅
- [x] Intent test harness
- [x] Mock service support
- [x] State transition testing
- [x] View rendering testing
- [x] Result validation
- [x] Benchmarking support

### Framework Enhancements Needed
- [ ] None identified - framework is complete

---

## Implementation Order & Dependencies

```
Phase 1: Preparation (CURRENT)
├── [x] Audit legacy screens
├── [x] Document features
├── [x] Plan state machines
└── [x] Define structures (THIS DOCUMENT)

Phase 2: Implementation (NEXT)
├── BurstManagement (16h)
│   ├── Context & Result structures
│   ├── Intent implementation
│   └── 30+ tests
├── FactManagement (16h)
│   ├── Context & Result structures
│   ├── Intent implementation
│   └── 40+ tests
├── ImportWizard (12h)
│   ├── Context & Result structures
│   ├── Intent implementation
│   └── 25+ tests
├── MetadataEditor (10h)
│   ├── Context & Result structures
│   ├── Intent implementation
│   └── 20+ tests
└── BulkOperations (10h)
    ├── Context & Result structures
    ├── Intent implementation
    └── 20+ tests

Phase 3: Root App Rebuild (14h)
├── Audit app.go
├── Design minimal structure
├── Implement new app.go
├── Register all intents
└── Remove legacy code

Phase 4: Testing & Validation (34h)
├── Unit tests
├── Integration tests
├── E2E workflows
├── Performance benchmarks
└── Code quality verification
```

---

**Document Version**: 1.0
**Last Updated**: 2026-01-03
**Status**: Complete - Ready for Phase 1.4 and Phase 2 Implementation

