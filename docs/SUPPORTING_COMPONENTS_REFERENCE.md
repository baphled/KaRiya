---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Supporting Components Reference

**Last Updated**: 2026-01-03
**Total Components**: 48 files, 34 model types + 14 utility/support files
**Purpose**: Document all reusable UI components and supporting infrastructure

---

## Overview

Supporting components are reusable UI building blocks, utilities, and infrastructure that support the 5 core intents (CaptureEvent, BrowseTimeline, GenerateCV, ExportArtifact, ConfigureSystem).

### Component Categories

1. **Core Infrastructure** (5 files)
   - Base models, message types, standard patterns
   - Used by ALL intents

2. **Dialog & Modal Components** (3 files)
   - Confirmation dialogs, editors, success screens
   - Reusable across multiple intents

3. **List & Table Components** (8 files)
   - List rendering, filtering, sorting
   - Used in browse/selection screens

4. **Data Entry Components** (6 files)
   - Forms, editors, input fields
   - Used in capture and edit flows

5. **Shortcut & Navigation** (4 files)
   - Keyboard shortcut handling
   - Context-aware navigation

6. **CV-Specific Components** (8 files)
   - CV generation, preview, export
   - Used in GenerateCV and ExportArtifact intents

7. **Burst & Fact Components** (8 files)
   - Burst suggestion, details, editing
   - Fact display and editing
   - Used in timeline and CV generation

---

## Component Inventory

### 1. Core Infrastructure (5 files)

#### `standard_model.go`
- **Type**: `BaseStandardModel`
- **Purpose**: Base model for all screen models
- **Usage Count**: 25+ (used by most models)
- **Key Methods**: Init(), Update(), View()
- **Dependencies**: tea, lipgloss
- **Test Coverage**: 95%+

#### `messages.go`
- **Type**: Message type definitions (30+ types)
- **Purpose**: Central message hub for all UI events
- **Key Types**:
  - ViewEventMsg, EditEventMsg, EventActionMenuMsg
  - ConfirmBurstSuggestionMsg, RejectBurstSuggestionMsg
  - BackMsg, HelpMsg, MainMenuMsg
- **Usage Count**: 40+ (used throughout)
- **Dependencies**: domain/career types
- **Test Coverage**: 100% (message passing)

#### `errors.go`
- **Type**: Error definitions
- **Purpose**: Centralized error types for UI layer
- **Key Types**: UIError, ValidationError, OperationError
- **Usage Count**: 15+ (error handling)
- **Dependencies**: None (domain independent)
- **Test Coverage**: 90%+

#### `help.go`
- **Type**: `HelpModel`
- **Purpose**: Help screen display
- **Usage Count**: 3 (main menu, navigation)
- **Key Methods**: View(), Update()
- **Dependencies**: styles, lipgloss
- **Test Coverage**: 85%+

#### `menu.go`
- **Type**: `MenuModel`
- **Purpose**: Main menu navigation
- **Usage Count**: 1 (root navigation)
- **Key Methods**: Update(), View(), GetSelected()
- **Dependencies**: styles, messages
- **Test Coverage**: 88%+

---

### 2. Dialog & Modal Components (3 files)

#### `confirmation_dialog.go`
- **Type**: `ConfirmationDialog`
- **Purpose**: Yes/No confirmation dialogs
- **Usage Count**: 8 (delete, confirm actions)
- **Key Methods**: Update(), View(), IsConfirmed(), IsCancelled()
- **Dependencies**: styles, lipgloss
- **Test Coverage**: 92%+

#### `success.go`
- **Type**: `SuccessModel`
- **Purpose**: Success notification screen
- **Usage Count**: 6 (after save/delete operations)
- **Key Methods**: View(), Update()
- **Dependencies**: styles, messages
- **Test Coverage**: 87%+

#### `error_handler.go`
- **Type**: `ErrorHandlerModel`
- **Purpose**: Error display and recovery
- **Usage Count**: 10 (error states)
- **Key Methods**: Update(), View(), DisplayError()
- **Dependencies**: styles, errors.go
- **Test Coverage**: 85%+

---

### 3. List & Table Components (8 files)

#### `list.go`
- **Type**: `ListModel`
- **Purpose**: Generic list rendering with selection
- **Usage Count**: 12 (timelines, facts, bursts)
- **Key Methods**: Update(), View(), GetSelected(), SetItems()
- **Dependencies**: styles, lipgloss
- **Test Coverage**: 88%+

#### `list_patterns.go`
- **Type**: List rendering patterns
- **Purpose**: Reusable list rendering patterns
- **Usage Count**: 8 (various list screens)
- **Key Functions**: RenderListItem(), RenderSelected(), RenderHighlighted()
- **Dependencies**: styles, lipgloss
- **Test Coverage**: 82%+

#### `filter.go`
- **Type**: `FilterModel`
- **Purpose**: Filter configuration for lists
- **Usage Count**: 6 (event/fact/burst filtering)
- **Key Methods**: Update(), View(), ApplyFilter(), ClearFilter()
- **Dependencies**: domain/career types
- **Test Coverage**: 86%+

#### `sort.go`
- **Type**: `SortModel`
- **Purpose**: Sort configuration for lists
- **Usage Count**: 5 (event/fact sorting)
- **Key Methods**: Update(), View(), ApplySorting(), ReverseSorting()
- **Dependencies**: domain/career types
- **Test Coverage**: 84%+

#### `search.go`
- **Type**: `SearchModel`
- **Purpose**: Search/filter UI
- **Usage Count**: 7 (event/fact/burst search)
- **Key Methods**: Update(), View(), GetSearchTerm(), ClearSearch()
- **Dependencies**: styles, lipgloss
- **Test Coverage**: 83%+

#### `fact_list.go`
- **Type**: `FactListModel`
- **Purpose**: Specialized list for facts
- **Usage Count**: 3 (fact browsing)
- **Key Methods**: Update(), View(), GetSelected()
- **Dependencies**: list.go, domain/career
- **Test Coverage**: 87%+

#### `burst_list.go`
- **Type**: `BurstListModel`
- **Purpose**: Specialized list for bursts
- **Usage Count**: 2 (burst browsing)
- **Key Methods**: Update(), View(), GetSelected()
- **Dependencies**: list.go, domain/career
- **Test Coverage**: 85%+

#### `cv_list.go`
- **Type**: `CVListModel`
- **Purpose**: Specialized list for CV configs
- **Usage Count**: 2 (CV selection)
- **Key Methods**: Update(), View(), GetSelected()
- **Dependencies**: list.go, domain/career
- **Test Coverage**: 86%+

---

### 4. Data Entry Components (6 files)

#### `form.go`
- **Type**: `FormModel`
- **Purpose**: Generic form with validation
- **Usage Count**: 8 (event capture, editing)
- **Key Methods**: Update(), View(), GetFormData(), Validate()
- **Dependencies**: styles, lipgloss, domain/career
- **Test Coverage**: 89%+

#### `fact_editor.go`
- **Type**: `FactEditorModel`
- **Purpose**: Specialized editor for facts
- **Usage Count**: 3 (fact creation/editing)
- **Key Methods**: Update(), View(), GetFact()
- **Dependencies**: form.go, domain/career
- **Test Coverage**: 86%+

#### `burst_editor.go`
- **Type**: `BurstEditorModel`
- **Purpose**: Specialized editor for bursts
- **Usage Count**: 2 (burst creation/editing)
- **Key Methods**: Update(), View(), GetBurst()
- **Dependencies**: form.go, domain/career
- **Test Coverage**: 84%+

#### `metadata_editor.go`
- **Type**: `MetadataEditorModel`
- **Purpose**: Editor for event metadata (company, project, tags)
- **Usage Count**: 4 (event metadata editing)
- **Key Methods**: Update(), View(), GetMetadata()
- **Dependencies**: form.go, domain/career
- **Test Coverage**: 85%+

#### `role_selector.go`
- **Type**: `RoleSelectorModel`
- **Purpose**: Role selection for CV generation
- **Usage Count**: 1 (CV audience configuration)
- **Key Methods**: Update(), View(), GetSelected()
- **Dependencies**: styles, lipgloss
- **Test Coverage**: 82%+

#### `audience_configurator.go`
- **Type**: `AudienceConfiguratorModel`
- **Purpose**: Configure target audience for CV
- **Usage Count**: 1 (CV generation flow)
- **Key Methods**: Update(), View(), GetAudience()
- **Dependencies**: role_selector.go, domain/career
- **Test Coverage**: 83%+

---

### 5. Shortcut & Navigation (4 files)

#### `shortcut_handler.go`
- **Type**: `ShortcutHandler`
- **Purpose**: Global keyboard shortcut handling
- **Usage Count**: 25+ (all screens)
- **Key Methods**: Handle(), Register(), Unregister()
- **Dependencies**: tea
- **Test Coverage**: 88%+

#### `context_shortcut_handler.go`
- **Type**: `ContextShortcutHandler`
- **Purpose**: Context-aware shortcut handling
- **Usage Count**: 8 (intent-specific shortcuts)
- **Key Methods**: Handle(), RegisterContext(), ClearContext()
- **Dependencies**: shortcut_handler.go
- **Test Coverage**: 85%+

#### `shortcut_mapper.go`
- **Type**: `ShortcutMapper`
- **Purpose**: Map shortcuts to actions
- **Usage Count**: 15+ (all intent models)
- **Key Methods**: Map(), GetAction(), ListShortcuts()
- **Dependencies**: shortcut_handler.go
- **Test Coverage**: 86%+

#### `shortcut_customizer.go`
- **Type**: `ShortcutCustomizerModel`
- **Purpose**: Allow user customization of shortcuts
- **Usage Count**: 1 (settings screen)
- **Key Methods**: Update(), View(), GetCustomShortcuts()
- **Dependencies**: shortcut_handler.go, styles
- **Test Coverage**: 80%+

---

### 6. CV-Specific Components (8 files)

#### `cv_generator.go`
- **Type**: `CVGeneratorModel`
- **Purpose**: Main CV generation orchestration
- **Usage Count**: 1 (GenerateCV intent)
- **Key Methods**: Update(), View(), GenerateCV()
- **Dependencies**: domain/career, service/career
- **Test Coverage**: 87%+

#### `cv_config_manager.go`
- **Type**: `CVConfigManagerModel`
- **Purpose**: Manage CV configurations
- **Usage Count**: 1 (GenerateCV intent)
- **Key Methods**: Update(), View(), GetConfigs(), SaveConfig()
- **Dependencies**: domain/career, service/career
- **Test Coverage**: 88%+

#### `cv_config_editor.go`
- **Type**: `CVConfigEditorModel`
- **Purpose**: Edit CV configuration details
- **Usage Count**: 2 (create/edit config)
- **Key Methods**: Update(), View(), GetConfig()
- **Dependencies**: form.go, domain/career
- **Test Coverage**: 85%+

#### `cv_preview.go`
- **Type**: `CVPreviewModel`
- **Purpose**: Preview generated CV
- **Usage Count**: 1 (GenerateCV intent)
- **Key Methods**: Update(), View(), ScrollUp(), ScrollDown()
- **Dependencies**: styles, lipgloss
- **Test Coverage**: 84%+

#### `cv_export_dialog.go`
- **Type**: `CVExportDialogModel`
- **Purpose**: Configure export options
- **Usage Count**: 1 (ExportArtifact intent)
- **Key Methods**: Update(), View(), GetExportOptions()
- **Dependencies**: styles, lipgloss
- **Test Coverage**: 82%+

#### `cv_export_progress.go`
- **Type**: `CVExportProgressModel`
- **Purpose**: Show export progress
- **Usage Count**: 1 (ExportArtifact intent)
- **Key Methods**: Update(), View(), SetProgress()
- **Dependencies**: styles, lipgloss
- **Test Coverage**: 83%+

#### `cv_export_success.go`
- **Type**: `CVExportSuccessModel`
- **Purpose**: Show successful export result
- **Usage Count**: 1 (ExportArtifact intent)
- **Key Methods**: Update(), View(), GetExportPath()
- **Dependencies**: styles, lipgloss
- **Test Coverage**: 81%+

#### `quality_indicator.go`
- **Type**: `QualityIndicator`
- **Purpose**: Display CV quality metrics
- **Usage Count**: 2 (preview, export)
- **Key Methods**: Update(), View(), SetMetrics()
- **Dependencies**: styles, lipgloss
- **Test Coverage**: 80%+

---

### 7. Burst & Fact Components (8 files)

#### `burst_suggestion.go`
- **Type**: `BurstSuggestionModel`
- **Purpose**: Display and handle burst suggestions
- **Usage Count**: 1 (CaptureEvent intent)
- **Key Methods**: Update(), View(), GetSuggestion()
- **Dependencies**: domain/career, service/career
- **Test Coverage**: 86%+

#### `burst_details.go`
- **Type**: `BurstDetailsModel`
- **Purpose**: Display burst details
- **Usage Count**: 2 (BrowseTimeline, GenerateCV)
- **Key Methods**: Update(), View(), GetBurst()
- **Dependencies**: domain/career, styles
- **Test Coverage**: 84%+

#### `burst_card.go`
- **Type**: `BurstCard`
- **Purpose**: Render burst as card
- **Usage Count**: 5 (list items, previews)
- **Key Methods**: View(), SetWidth(), SetHighlight()
- **Dependencies**: styles, lipgloss
- **Test Coverage**: 82%+

#### `fact_details.go`
- **Type**: `FactDetailsModel`
- **Purpose**: Display fact details
- **Usage Count**: 2 (BrowseTimeline, GenerateCV)
- **Key Methods**: Update(), View(), GetFact()
- **Dependencies**: domain/career, styles
- **Test Coverage**: 83%+

#### `fact_card.go`
- **Type**: `FactCard`
- **Purpose**: Render fact as card
- **Usage Count**: 4 (list items, previews)
- **Key Methods**: View(), SetWidth(), SetHighlight()
- **Dependencies**: styles, lipgloss
- **Test Coverage**: 81%+

#### `fact_search.go`
- **Type**: `FactSearchModel`
- **Purpose**: Search facts with advanced filters
- **Usage Count**: 2 (fact browsing)
- **Key Methods**: Update(), View(), Search()
- **Dependencies**: search.go, domain/career
- **Test Coverage**: 82%+

#### `facts_results.go`
- **Type**: `FactsResultsModel`
- **Purpose**: Display fact search results
- **Usage Count**: 1 (fact search)
- **Key Methods**: Update(), View(), GetResults()
- **Dependencies**: fact_card.go, list.go
- **Test Coverage**: 80%+

#### `action_menu.go`
- **Type**: `ActionMenuModel`
- **Purpose**: Context menu for event/burst/fact actions
- **Usage Count**: 8 (action selection)
- **Key Methods**: Update(), View(), GetSelected()
- **Dependencies**: styles, lipgloss
- **Test Coverage**: 87%+

---

### 8. Other Components (3 files)

#### `view_event.go`
- **Type**: `ViewEventModel`
- **Purpose**: Display event details
- **Usage Count**: 1 (BrowseTimeline intent)
- **Key Methods**: Update(), View(), GetEvent()
- **Dependencies**: domain/career, styles
- **Test Coverage**: 86%+

#### `details.go`
- **Type**: `DetailsModel`
- **Purpose**: Generic details display
- **Usage Count**: 4 (various detail screens)
- **Key Methods**: Update(), View(), SetDetails()
- **Dependencies**: styles, lipgloss
- **Test Coverage**: 84%+

#### `metadata_review.go`
- **Type**: `MetadataReviewModel`
- **Purpose**: Review event metadata before save
- **Usage Count**: 1 (CaptureEvent intent)
- **Key Methods**: Update(), View(), GetMetadata()
- **Dependencies**: domain/career, styles
- **Test Coverage**: 85%+

#### `bulk_operations.go`
- **Type**: `BulkOperationsModel`
- **Purpose**: Handle bulk actions on multiple items
- **Usage Count**: 2 (bulk delete, bulk export)
- **Key Methods**: Update(), View(), GetSelected()
- **Dependencies**: list.go, domain/career
- **Test Coverage**: 81%+

#### `import_review.go`
- **Type**: `ImportReviewModel`, `ImportProgressModel`
- **Purpose**: Review and track import operations
- **Usage Count**: 1 (data import)
- **Key Methods**: Update(), View(), GetProgress()
- **Dependencies**: domain/career, service/career
- **Test Coverage**: 83%+

#### `shortcut_help_system.go`
- **Type**: `ShortcutHelpSystem`
- **Purpose**: Display available shortcuts
- **Usage Count**: 8 (help screens)
- **Key Methods**: GetHelp(), ListShortcuts(), GetContextHelp()
- **Dependencies**: shortcut_handler.go, styles
- **Test Coverage**: 84%+

---

## Component Dependencies

### Dependency Graph (High-Level)

```
BaseStandardModel (foundation)
├── All Model Types (inherit from)
│
Messages.go (central message hub)
├── All Models (send/receive)
│
Styles (UI styling)
├── All Rendering Components
│
Domain/Career (data models)
├── List/Table Components
├── Data Entry Components
├── Burst/Fact Components
│
Service/Career (business logic)
├── CV Generator
├── Burst Suggestion
├── Fact Components
```

### Critical Dependencies

1. **BaseStandardModel**: 25+ models inherit from this
2. **Messages.go**: 40+ types used throughout
3. **Styles**: All rendering components depend on this
4. **Domain/Career**: All data-dependent components

---

## Usage Statistics

| Category | Count | Avg Reuse | Test Coverage |
|----------|-------|-----------|----------------|
| Core Infrastructure | 5 | 20+ | 94% |
| Dialog/Modal | 3 | 6 | 88% |
| List/Table | 8 | 7 | 86% |
| Data Entry | 6 | 4 | 85% |
| Shortcut/Nav | 4 | 12 | 85% |
| CV-Specific | 8 | 1.5 | 84% |
| Burst/Fact | 8 | 3 | 83% |
| Other | 6 | 2 | 84% |
| **Total** | **48** | **6.6** | **85%** |

---

## Extraction Candidates

### Tier 1: Extract First (High Reuse, Stable)
- BaseStandardModel → Base package
- Messages → Messages package
- Shortcut components → Navigation package
- List/Filter/Sort → List package

### Tier 2: Extract Next (Medium Reuse, Stable)
- Form, Editors → Forms package
- Burst/Fact components → Domain packages
- CV components → CV package

### Tier 3: Extract Later (Low Reuse, Specific)
- Dialog components → Dialogs package
- Success/Error handlers → Feedback package
- Other specialized components → Utilities package

---

## Next Steps

1. **Phase 4**: Consolidate shortcut system (4 → 1 unified system)
2. **Phase 5**: Extract components into packages (reduce file count from 48 → 25)
3. **Phase 6**: Integrate with intent system (components become intent sub-models)
4. **Phase 7**: Add component composition patterns (avoid duplication)

---

**Status**: ✅ Complete
**Generated**: 2026-01-03
**Maintainer**: KaRiya Development Team

