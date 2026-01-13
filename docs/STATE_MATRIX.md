# KaRiya TUI State Matrix

> **Auto-generated**: Do not edit manually. Run `make generate-state-matrix` to update.  
> **Generated**: 2026-01-12T18:00:55Z  
> **Source**: `cmd/generate-state-matrix/main.go`

## Summary

| Metric | Value |
|--------|-------|
| Total Intents | 10 |
| Total States | 64 |

## State Type Definitions

| Type | Description | Escape Behavior |
|------|-------------|-----------------|
| **ROOT** | Entry state of intent | Cancel intent → Main Menu |
| **Intermediate** | Middle navigation state | Back → Previous state |
| **Confirmation** | User confirmation required | Cancel → Parent state |
| **Async** | Background operation | Varies (some block, some allow back) |
| **Modal** | Overlay editor/form | Close modal → Parent state |
| **Final** | Operation complete | Deactivate intent |
| **Error** | Operation failed | Deactivate (retry often available) |

## Complete State Matrix


### BrowseTimeline

**File**: `/home/baphled/Projects/KaRiya/internal/cli/intents/browse_timeline.go`  
**States**: 3

| State | Type | Escape Behavior |
|-------|------|-----------------|
| `BrowseStateTimeline` | ROOT | Cancel intent → Main Menu |
| `BrowseStateEventDetail` | Intermediate | Back → Previous state |
| `BrowseStateDeleteConfirm` | Confirmation | Cancel → Parent state |


### BulkOperations

**File**: `/home/baphled/Projects/KaRiya/internal/cli/intents/bulk_operations.go`  
**States**: 4

| State | Type | Escape Behavior |
|-------|------|-----------------|
| `BulkSelectOpState` | ROOT | Cancel intent → Main Menu |
| `BulkConfigureState` | Intermediate | Back → Previous state |
| `BulkExecuteState` | Intermediate | Back → Previous state |
| `BulkCompleteState` | Final | Deactivate intent |


### BurstManagement

**File**: `/home/baphled/Projects/KaRiya/internal/cli/intents/burst_management_intent.go`  
**States**: 14

| State | Type | Escape Behavior |
|-------|------|-----------------|
| `BurstListState` | ROOT | Cancel intent → Main Menu |
| `BurstViewState` | Intermediate | Back → Previous state |
| `BurstEditorState` | Intermediate | Back → Previous state |
| `BurstDeleteConfirmState` | Confirmation | Cancel → Parent state |
| `BurstSuggestState` | Intermediate | Back → Previous state |
| `BurstCompletedState` | Final | Deactivate intent |
| `BurstStateList` | ROOT | Cancel intent → Main Menu |
| `BurstStateDetail` | Intermediate | Back → Previous state |
| `BurstStateDetailEvents` | Intermediate | Back → Previous state |
| `BurstStateDetailFacts` | Intermediate | Back → Previous state |
| `BurstStateEdit` | Intermediate | Back → Previous state |
| `BurstStateDeleteConfirm` | Confirmation | Cancel → Parent state |
| `BurstStateConfirm` | Confirmation | Cancel → Parent state |
| `BurstStateExtractingFacts` | Async | Varies (let complete or allow back) |


### CaptureEvent

**File**: `/home/baphled/Projects/KaRiya/internal/cli/intents/capture_event.go`  
**States**: 4

| State | Type | Escape Behavior |
|-------|------|-----------------|
| `CaptureStateChooseStrategy` | ROOT | Cancel intent → Main Menu |
| `CaptureStateForm` | Intermediate | Back → Previous state |
| `CaptureStateReview` | Intermediate | Back → Previous state |
| `CaptureStateSubmit` | Intermediate | Back → Previous state |


### ConfigureSystem

**File**: `/home/baphled/Projects/KaRiya/internal/cli/intents/configure_system.go`  
**States**: 7

| State | Type | Escape Behavior |
|-------|------|-----------------|
| `ConfigStateSelectDomain` | ROOT | Cancel intent → Main Menu |
| `ConfigStateEditSettings` | Intermediate | Back → Previous state |
| `ConfigStateReviewChanges` | Intermediate | Back → Previous state |
| `ConfigStateConfirm` | Confirmation | Cancel → Parent state |
| `ConfigStateSaving` | Async | Varies (let complete or allow back) |
| `ConfigStateComplete` | Final | Deactivate intent |
| `ConfigStateFailed` | Error | Deactivate intent (retry available) |


### ExportArtifact

**File**: `/home/baphled/Projects/KaRiya/internal/cli/intents/export_artifact.go`  
**States**: 9

| State | Type | Escape Behavior |
|-------|------|-----------------|
| `ExportStateSelectType` | ROOT | Cancel intent → Main Menu |
| `ExportStateSelectFormat` | Intermediate | Back → Previous state |
| `ExportStateSelectDest` | Intermediate | Back → Previous state |
| `ExportStateConfigure` | Intermediate | Back → Previous state |
| `ExportStatePreview` | Intermediate | Back → Previous state |
| `ExportStateConfirm` | Confirmation | Cancel → Parent state |
| `ExportStateInProgress` | Async | Varies (let complete or allow back) |
| `ExportStateComplete` | Final | Deactivate intent |
| `ExportStateFailed` | Error | Deactivate intent (retry available) |


### FactManagement

**File**: `/home/baphled/Projects/KaRiya/internal/cli/intents/fact_management.go`  
**States**: 6

| State | Type | Escape Behavior |
|-------|------|-----------------|
| `FactListState` | ROOT | Cancel intent → Main Menu |
| `FactViewState` | Intermediate | Back → Previous state |
| `FactEditorState` | Intermediate | Back → Previous state |
| `FactDeleteConfirmState` | Confirmation | Cancel → Parent state |
| `FactResultsState` | Intermediate | Back → Previous state |
| `FactCompletedState` | Final | Deactivate intent |


### GenerateCv

**File**: `/home/baphled/Projects/KaRiya/internal/cli/intents/generate_cv.go`  
**States**: 10

| State | Type | Escape Behavior |
|-------|------|-----------------|
| `GenerateCVStateSelectProfile` | ROOT | Cancel intent → Main Menu |
| `GenerateCVStateSelectAudience` | Intermediate | Back → Previous state |
| `GenerateCVStateGenerating` | Async | Varies (let complete or allow back) |
| `GenerateCVStatePreview` | Intermediate | Back → Previous state |
| `GenerateCVStateReview` | Intermediate | Back → Previous state |
| `GenerateCVStateConfirm` | Confirmation | Cancel → Parent state |
| `GenerateCVStateExportSelectFormat` | Intermediate | Back → Previous state |
| `GenerateCVStateExportSelectLocation` | Intermediate | Back → Previous state |
| `GenerateCVStateExporting` | Async | Varies (let complete or allow back) |
| `GenerateCVStateExportComplete` | Final | Deactivate intent |


### ImportWizard

**File**: `/home/baphled/Projects/KaRiya/internal/cli/intents/import_wizard.go`  
**States**: 4

| State | Type | Escape Behavior |
|-------|------|-----------------|
| `ImportFileSelectState` | ROOT | Cancel intent → Main Menu |
| `ImportPreviewState` | Intermediate | Back → Previous state |
| `ImportProgressState` | Async | Varies (let complete or allow back) |
| `ImportCompleteState` | Final | Deactivate intent |


### MetadataEditor

**File**: `/home/baphled/Projects/KaRiya/internal/cli/intents/metadata_editor.go`  
**States**: 3

| State | Type | Escape Behavior |
|-------|------|-----------------|
| `MetadataReviewState` | Intermediate | Back → Previous state |
| `MetadataEditState` | Intermediate | Back → Previous state |
| `MetadataConfirmState` | Confirmation | Cancel → Parent state |



---

*Generated by `cmd/generate-state-matrix/main.go`*
