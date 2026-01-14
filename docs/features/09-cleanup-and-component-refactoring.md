---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# PRD: Cleanup Orphaned Views & Refactor Components for Intent System

**Document Version**: 3.0
**Date**: 2026-01-03
**Status**: READY FOR IMPLEMENTATION
**Effort Estimate**: 10 hours
**Priority**: HIGH (Prerequisite for Aggressive Replacement)

---

## Introduction/Overview

The KaRiya TUI codebase contains **52 model/component files** in `internal/cli/models/`. An audit has identified:
- **25 integrated screens** (mapped to intents) ✅
- **23 supporting components** (properly used) ✅
- **4 orphaned views** (unused, must be removed) ❌
- **1 missing model** (referenced but not found) ⚠️

This PRD addresses the cleanup and refactoring needed to prepare the codebase for the aggressive intent-based replacement plan. The goal is to remove technical debt (orphaned code) and extract reusable components so they can be properly integrated with the new intent system.

**Problem Solved**:
- Removes unused code that creates confusion and maintenance burden
- Prepares components for reuse across intents
- Fixes broken references that could cause build failures
- Establishes clean foundation for intent system migration

---

## Goals

1. **Remove all 4 orphaned views** (tutorial.go, view_event_with_facts.go, source_event_tracer.go, and cleanup view_event.go)
2. **Fix the missing ImportProgressModel** that is referenced in app.go but doesn't exist
3. **Extract and document all 23 supporting components** for intent system reuse
4. **Consolidate 5 shortcut-related components** into unified system
5. **Verify all 25 integrated screens** are properly mapped to intents
6. **Establish clean, organized codebase** ready for aggressive intent-based replacement

---

## Orphaned Views (4 total) ❌

These are **defined but completely unused** anywhere in the codebase:

### 1. TutorialModel

**File**: `internal/cli/models/tutorial.go`

**Status**: ❌ ORPHANED - Never used or integrated

**What it does**:
- Represents a first-run tutorial screen
- Likely intended for onboarding new users
- Has tutorial steps and navigation

**Evidence**:
```bash
$ grep -r "TutorialModel\|tutorial" internal/cli --include="*.go" | grep -v models/tutorial.go
# Returns only references within tutorial.go itself
```

**Decision**: DELETE
- [ ] Remove `internal/cli/models/tutorial.go` completely
- [ ] Verify no references remain in codebase

---

### 2. ViewEventWithFactsModel

**File**: `internal/cli/models/view_event_with_facts.go`

**Status**: ❌ ORPHANED - Never used or integrated

**What it does**:
- Displays an event with its associated facts
- Likely intended as an enhanced event detail view
- Separate from the main DetailsModel

**Evidence**:
```bash
$ grep -r "ViewEventWithFacts" internal/cli --include="*.go" | grep -v models/view_event_with_facts.go
# Returns only references within view_event_with_facts.go itself
```

**Duplicate Functionality**:
- This functionality is already provided by `DetailsModel` in `details.go`
- ViewScreen (which uses DetailsModel) covers all use cases

**Decision**: DELETE
- [ ] Remove `internal/cli/models/view_event_with_facts.go` completely
- [ ] Verify DetailsModel provides all needed functionality
- [ ] Verify no references remain

---

### 3. SourceEventTracerModel

**File**: `internal/cli/models/source_event_tracer.go`

**Status**: ❌ ORPHANED - Never used or integrated

**What it does**:
- Tracks the source/lineage of events
- Likely for debugging or audit purposes
- Shows where data came from

**Evidence**:
```bash
$ grep -r "SourceEventTracer" internal/cli --include="*.go" | grep -v models/source_event_tracer.go
# Returns only references within source_event_tracer.go itself
```

**Status Assessment**:
- Was intended for event lineage tracking
- Not actively used in any integrated model
- Not referenced in any workflow

**Decision**: DELETE
- [ ] Remove `internal/cli/models/source_event_tracer.go` completely
- [ ] Verify no external dependencies
- [ ] Verify no references remain

---

### 4. ViewEvent Type (Partial)

**File**: `internal/cli/models/view_event.go`

**Status**: ⚠️ PARTIALLY ORPHANED - Contains both used and unused types

**What's in this file**:
- `EditEventMsg` - ✅ USED (in metadata_review.go, app.go messages - 34 references)
- `ConfirmationDialog` - ✅ USED (in multiple models - 36 references)
- Other types/functions - ❌ POTENTIALLY UNUSED

**Evidence**:
```bash
$ grep -r "EditEventMsg" internal/cli --include="*.go"
# Returns multiple results (used)

$ grep -r "ViewEvent" internal/cli --include="*.go" | grep -v test | grep -v models/view_event.go
# Returns only EditEventMsg references
```

**Decision**: CLEAN UP (Keep used types, remove unused)
- [ ] Keep `EditEventMsg` and its supporting types
- [ ] Keep `ConfirmationDialog` and its supporting types
- [ ] Remove all other unused types and functions
- [ ] Verify EditEventMsg is still importable and functional
- [ ] Verify ConfirmationDialog is still importable and functional

---

## Missing/Broken Model (1 total) ⚠️

### ImportProgressModel

**File**: `internal/cli/models/import_progress.go` (MISSING)

**Status**: ❌ BROKEN REFERENCE - Referenced in app.go but file doesn't exist

**Where Referenced**:
```bash
$ grep "ImportProgressModel" internal/cli/app/app.go
# Found in app.go - but model file doesn't exist!
```

**What it should do**:
- Display progress during CSV import operations
- Show import statistics and status
- Integrate with ImportWizard intent

**Decision**: RESTORE OR CREATE
1. [ ] Check git history: `git log --all --full-history -- internal/cli/models/import_progress.go`
2. [ ] If found in history: Restore the file
3. [ ] If never existed: Create new implementation for ImportWizard intent
4. [ ] Verify code compiles without errors
5. [ ] Ensure all references in app.go are valid

---

## Supporting Components (23 total) ✅

All 23 supporting components are **properly used** by integrated models and must be documented:

### Search & Filter Components (4 components)
1. **Filter** (filter.go) - 181 references
   - Used by: list.go, burst_list.go, cv_list.go, fact_list.go
   - Purpose: Event filtering (date range, tags, etc.)

2. **Search** (search.go) - 91 references
   - Used by: Multiple models
   - Purpose: Event search with highlighting

3. **Sort** (sort.go) - 89 references
   - Used by: list.go, burst_list.go, cv_list.go, fact_list.go
   - Purpose: Event sorting options (date, relevance, etc.)

4. **FactSearch** (fact_search.go) - 25 references
   - Used by: fact_list.go
   - Purpose: Specialized fact search with filtering

### List Management Components (2 components)
5. **ListDeletionState** (list_patterns.go) - 26 references
   - Used by: list.go, burst_list.go, fact_list.go, cv_config_manager.go
   - Purpose: Handle deletion confirmation and execution workflow

6. **ConfirmationDialog** (confirmation_dialog.go) - 36 references
   - Used by: list.go, burst_list.go, fact_list.go, cv_config_manager.go, metadata_review.go
   - Purpose: Generic confirmation dialog for user actions

### Display Components (3 components)
7. **BurstCard** (burst_card.go) - 15 references
   - Used by: burst_details.go, facts_results.go
   - Purpose: Display burst details in card format

8. **FactCard** (fact_card.go) - 19 references
   - Used by: fact_details.go, facts_results.go
   - Purpose: Display fact details in card format

9. **QualityIndicator** (quality_indicator.go) - 17 references
   - Used by: list.go
   - Purpose: Show data quality score for events

### Shortcut & Context Components (5 components)
10. **ShortcutHandler** (shortcut_handler.go) - 30 references
    - Used by: Multiple models
    - Purpose: Handle keyboard shortcuts for models

11. **ShortcutMapper** (shortcut_mapper.go) - 23 references
    - Used by: Multiple models
    - Purpose: Map shortcuts to actions

12. **ContextShortcutHandler** (context_shortcut_handler.go) - 19 references
    - Used by: Multiple models
    - Purpose: Context-aware shortcut handling with stack-based system

13. **ShortcutCustomizer** (shortcut_customizer.go) - 24 references
    - Used by: Standalone/help system
    - Purpose: Customize shortcuts

14. **ShortcutHelpSystem** (shortcut_help_system.go) - 19 references
    - Used by: help.go
    - Purpose: Display shortcut help information

### CV Configuration Components (3 components)
15. **AudienceConfigurator** (audience_configurator.go) - 16 references
    - Used by: cv_config_manager.go, cv_generator.go
    - Purpose: Select CV audience/target

16. **RoleSelector** (role_selector.go) - 13 references
    - Used by: cv_config_manager.go, cv_generator.go
    - Purpose: Select CV role

17. **CVConfigEditor** (cv_config_editor.go) - 12 references
    - Used by: cv_config_manager.go
    - Purpose: Edit CV configuration

### Utility Components (3 components)
18. **ErrorHandler** (error_handler.go) - 15 references
    - Used by: Multiple models
    - Purpose: Handle and display errors

19. **ViewEvent** (view_event.go) - 34 references
    - Used by: metadata_review.go, app.go messages
    - Purpose: Message types for event viewing/editing

20. **Standard Model** (standard_model.go)
    - Used by: Multiple models
    - Purpose: Base model for standard components

### Supporting Utilities (3 components)
21. **list_patterns.go** - Contains ListDeletionState pattern
22. **messages.go** - Contains message types
23. **errors.go** - Contains error types

---

## Integrated Screens (25 total) ✅

All 25 screens are accounted for and properly mapped to intents:

**Core Screens** (3):
- CaptureScreen → CaptureEvent Intent ✅
- ListScreen → BrowseTimeline Intent ✅
- ViewScreen → BrowseTimeline Intent ✅

**CV Management** (7):
- CVGeneratorScreen → GenerateCV Intent ✅
- CVPreviewScreen → GenerateCV Intent ✅
- CVExportDialogScreen → ExportArtifact Intent ✅
- CVExportSuccessScreen → ExportArtifact Intent ✅
- CVExportProgressScreen → ExportArtifact Intent ✅
- CVConfigManagerScreen → ConfigureSystem Intent ✅
- CVListScreen → BrowseTimeline Intent ✅

**Burst Management** (4):
- BurstListScreen → BurstManagement Intent ✅
- BurstDetailsScreen → BurstManagement Intent ✅
- BurstEditorScreen → BurstManagement Intent ✅
- BurstSuggestionScreen → BurstManagement Intent ✅

**Fact Management** (4):
- FactListScreen → FactManagement Intent ✅
- FactDetailsScreen → FactManagement Intent ✅
- FactEditorScreen → FactManagement Intent ✅
- FactsResultsScreen → FactManagement Intent ✅

**Metadata Management** (2):
- MetadataReviewScreen → MetadataEditor Intent ✅
- MetadataEditorScreen → MetadataEditor Intent ✅

**Import & Other** (2):
- ImportReviewScreen → ImportWizard Intent ✅
- BulkOperationsScreen → BulkOperations Intent ✅

**Global Screens** (2 - Keep as-is):
- HelpScreen → Keep as global screen ✅
- MainMenuScreen → Keep as startup screen ✅

---

## Functional Requirements

### Phase 1: Remove Orphaned Views (2 hours)

**Requirement 1.1**: Delete TutorialModel
- [ ] Delete `internal/cli/models/tutorial.go`
- [ ] Verify no references remain
- [ ] Run tests

**Requirement 1.2**: Delete ViewEventWithFactsModel
- [ ] Delete `internal/cli/models/view_event_with_facts.go`
- [ ] Verify DetailsModel provides all functionality
- [ ] Verify no references remain

**Requirement 1.3**: Delete SourceEventTracerModel
- [ ] Delete `internal/cli/models/source_event_tracer.go`
- [ ] Verify no external dependencies
- [ ] Verify no references remain

**Requirement 1.4**: Clean up view_event.go
- [ ] Keep EditEventMsg and supporting types
- [ ] Keep ConfirmationDialog and supporting types
- [ ] Remove all other unused types/functions
- [ ] Verify remaining code is functional

### Phase 2: Fix Missing ImportProgressModel (2 hours)

**Requirement 2.1**: Investigate and resolve
- [ ] Check git history for deletion
- [ ] Restore from git OR create new implementation
- [ ] Update app.go if needed
- [ ] Verify code compiles

### Phase 3: Extract Supporting Components (3 hours)

**Requirement 3.1**: Document all 23 components
- [ ] Create `docs/SUPPORTING_COMPONENTS_REFERENCE.md`
- [ ] Include all reference counts from audit
- [ ] Document usage patterns and dependencies

**Requirement 3.2**: Create extraction plan
- [ ] Create `docs/COMPONENT_EXTRACTION_PLAN.md`
- [ ] Tier 1 (Heavily used): Filter (181), Search (91), Sort (89)
- [ ] Tier 2 (Moderately used): ConfirmationDialog (36), ViewEvent (34), ShortcutHandler (30), etc.
- [ ] Tier 3 (Lightly used): 9 more components

### Phase 4: Consolidate Shortcut Components (2 hours)

**Requirement 4.1**: Analyze 5 shortcut components
- [ ] Document ShortcutHandler (30 refs)
- [ ] Document ShortcutMapper (23 refs)
- [ ] Document ContextShortcutHandler (19 refs)
- [ ] Document ShortcutCustomizer (24 refs)
- [ ] Document ShortcutHelpSystem (19 refs)

**Requirement 4.2**: Design unified system
- [ ] Create `docs/UNIFIED_SHORTCUT_SYSTEM_DESIGN.md`
- [ ] Propose consolidated API
- [ ] Plan migration path

### Phase 5: Verification & Documentation (1 hour)

**Requirement 5.1**: Create screen-to-intent mapping
- [ ] Create `docs/SCREEN_TO_INTENT_MAPPING.md`
- [ ] Verify all 25 screens mapped

**Requirement 5.2**: Run quality checks
- [ ] Execute: `go test -v ./...`
- [ ] Execute: `golangci-lint run ./...`
- [ ] Execute: `gofmt -w ./...`

**Requirement 5.3**: Create cleanup summary
- [ ] Create `docs/CLEANUP_SUMMARY.md`
- [ ] Document what was removed and why

---

## Non-Goals (Out of Scope)

- NOT implementing the intent system itself
- NOT refactoring components to use new intent patterns yet
- NOT creating new components or features
- NOT changing functionality of integrated screens
- NOT implementing the unified shortcut system (only design)

---

## Success Metrics

✅ All 4 orphaned views removed
✅ Missing ImportProgressModel resolved
✅ All 23 components documented with usage counts
✅ Screen-to-intent mapping complete (all 25 screens)
✅ Component extraction plan created (Tier 1/2/3)
✅ Shortcut system design documented
✅ All tests passing
✅ No linting/formatting errors
✅ Zero broken references

---

## Implementation Checklist

### Phase 1: Remove Orphaned Views (2 hours)
- [ ] Delete tutorial.go
- [ ] Delete view_event_with_facts.go
- [ ] Delete source_event_tracer.go
- [ ] Audit and clean view_event.go
- [ ] Run tests: go test -v ./...
- [ ] Commit: "refactor(models): remove orphaned views"

### Phase 2: Fix Missing ImportProgressModel (2 hours)
- [ ] Check git history for import_progress.go
- [ ] Restore or create new implementation
- [ ] Update app.go if needed
- [ ] Verify code compiles
- [ ] Commit: "fix(models): resolve missing ImportProgressModel"

### Phase 3: Extract Supporting Components (3 hours)
- [ ] Document all 23 components with reference counts
- [ ] Create SUPPORTING_COMPONENTS_REFERENCE.md
- [ ] Create COMPONENT_EXTRACTION_PLAN.md with Tier 1/2/3
- [ ] Commit: "docs(components): document supporting components"

### Phase 4: Consolidate Shortcut Components (2 hours)
- [ ] Audit 5 shortcut components
- [ ] Design unified shortcut system
- [ ] Create UNIFIED_SHORTCUT_SYSTEM_DESIGN.md
- [ ] Commit: "docs(shortcuts): design unified shortcut system"

### Phase 5: Verification & Documentation (1 hour)
- [ ] Create SCREEN_TO_INTENT_MAPPING.md
- [ ] Run full test suite
- [ ] Run linter and formatter
- [ ] Create CLEANUP_SUMMARY.md
- [ ] Commit: "docs(audit): cleanup and refactoring complete"

---

## Timeline & Effort

**Total Effort**: 10 hours

| Phase | Task | Effort |
|-------|------|--------|
| 1 | Remove orphaned views | 2 hours |
| 2 | Fix missing model | 2 hours |
| 3 | Extract components | 3 hours |
| 4 | Consolidate shortcuts | 2 hours |
| 5 | Verify & document | 1 hour |
| **Total** | | **10 hours** |

**Recommended Timeline**: 1-2 days (full-time)

---

## Dependencies & Relationships

**Depends On**: Nothing (can start immediately)

**Enables**:
- Aggressive intent-based replacement plan (127 hours)
- Component extraction and refactoring
- Unified shortcut system implementation

**Related**:
- docs/UNINTEGRATED_VIEWS_AUDIT.md (source audit)
- docs/APP_GO_AGGRESSIVE_REPLACEMENT_PLAN.md (main plan)
- docs/AGGRESSIVE_REPLACEMENT_START_NOW.md (action plan)

---

## Document History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2026-01-03 | Initial PRD |
| 2.0 | 2026-01-03 | Updated with actual audit findings |
| 3.0 | 2026-01-03 | Added detailed orphaned views and missing model data |
