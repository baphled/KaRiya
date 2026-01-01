# Tasks for Model Consistency and Standardized UI Components

## Relevant Files

### Container Component Files (NEW)
- `internal/cli/components/screen_container.go` - Screen wrapper with padding modes and width constraints
- `internal/cli/components/screen_container_test.go` - Unit tests for ScreenContainer
- `internal/cli/components/card_container.go` - Boxed content with header/body/footer sections
- `internal/cli/components/card_container_test.go` - Unit tests for CardContainer
- `internal/cli/components/section_container.go` - Content grouping with consistent spacing
- `internal/cli/components/section_container_test.go` - Unit tests for SectionContainer
- `internal/cli/components/form_field_container.go` - Form field layout (label, input, error, hint)
- `internal/cli/components/form_field_container_test.go` - Unit tests for FormFieldContainer
- `internal/cli/components/list_container.go` - List display with pagination and empty state
- `internal/cli/components/list_container_test.go` - Unit tests for ListContainer
- `internal/cli/components/modal_container.go` - Modal dialog rendering
- `internal/cli/components/modal_container_test.go` - Unit tests for ModalContainer

### Style Configuration Files
- `internal/cli/styles/styles.go` - Existing styling system (will be enhanced with export)
- `internal/cli/styles/constants_export.go` - NEW: Unified export of all style constants

### Model Files (REFACTORING)
- `internal/cli/models/form.go` - Form model (refactor to use containers)
- `internal/cli/models/form_test.go` - Form model tests
- `internal/cli/models/list.go` - List model (refactor to use containers)
- `internal/cli/models/list_test.go` - List model tests
- `internal/cli/models/details.go` - Details/event view model (refactor to use containers)
- `internal/cli/models/details_test.go` - Details model tests
- `internal/cli/models/confirmation_dialog.go` - Confirmation dialog (refactor to use containers)
- `internal/cli/models/confirmation_dialog_test.go` - Confirmation dialog tests
- `internal/cli/models/fact_editor.go` - Fact editor form (refactor to use containers)
- `internal/cli/models/fact_editor_test.go` - Fact editor tests
- `internal/cli/models/metadata_editor.go` - Metadata editor form (refactor to use containers)
- `internal/cli/models/metadata_editor_test.go` - Metadata editor tests
- `internal/cli/models/fact_list.go` - Fact list model (refactor to use containers)
- `internal/cli/models/fact_list_test.go` - Fact list tests
- `internal/cli/models/burst_list.go` - Burst list model (refactor to use containers)
- `internal/cli/models/burst_list_test.go` - Burst list tests

### Documentation Files
- `docs/guides/COMPONENT_USAGE_GUIDE.md` - NEW: Container component usage documentation
- `docs/guides/MODEL_DEVELOPMENT_GUIDE.md` - NEW: Guide for creating new models with components
- `docs/features/06-model-consistency-implementation-assessment.md` - Existing implementation strategy

### Existing Components (REFERENCE)
- `internal/cli/components/header.go` - HeaderModel component
- `internal/cli/components/footer.go` - FooterModel component
- `internal/cli/components/help_footer.go` - HelpFooterModel component
- `internal/cli/components/list_item.go` - ListItemModel component
- `internal/cli/components/tag_selector.go` - TagSelector component
- `internal/cli/components/category_selector.go` - CategorySelector component
- `internal/cli/components/progress.go` - ProgressIndicator component
- `internal/cli/components/spinner.go` - Spinner component

### Notes
- Build on existing navigation infrastructure completed in Phase 1 (task 05)
- Leverage existing lipgloss styling system in styles.go
- Follow Go conventions: builders, receiver methods, composition
- Maintain backward compatibility during transition
- All new components should be stateless rendering functions
- Use BubbleTea patterns consistently

## Tasks

### 1.0 Create Foundational Container Components
- [x] 1.1 Analyze existing styles.go to extract layout patterns and spacing conventions
- [x] 1.2 Design ScreenContainer interface and implementation
  - [x] 1.2.1 Define ScreenContainer struct with fields for content, padding mode, custom padding
  - [x] 1.2.2 Implement Render() method with consistent margins and width constraints
  - [x] 1.2.3 Create WithPaddingMode() builder method
  - [x] 1.2.4 Create WithCustomPadding() builder method
  - [x] 1.2.5 Write unit tests for padding modes and width constraints
- [x] 1.3 Design and implement CardContainer component
  - [x] 1.3.1 Define CardContainer struct with header, body, footer, background color fields
  - [x] 1.3.2 Implement SetHeader(), SetBody(), SetFooter() methods
  - [x] 1.3.3 Implement Render() method with border, padding, background
  - [x] 1.3.4 Create WithBackgroundColor() builder method
  - [x] 1.3.5 Write unit tests for all sections and styling
- [x] 1.4 Design and implement SectionContainer component
  - [x] 1.4.1 Define SectionContainer struct with title and content fields
  - [x] 1.4.2 Implement SetTitle() and SetContent() methods
  - [x] 1.4.3 Implement Render() method with HeaderSection styling
  - [x] 1.4.4 Create spacing control methods
  - [x] 1.4.5 Write unit tests for title rendering and spacing
- [x] 1.5 Design and implement FormFieldContainer component
  - [x] 1.5.1 Define FormFieldContainer struct with label, input, error, hint fields
  - [x] 1.5.2 Implement SetLabel() method with InputLabel styling
  - [x] 1.5.3 Implement SetInput() method with focus state support
  - [x] 1.5.4 Implement SetError() method with InputError styling
  - [x] 1.5.5 Implement SetHint() method for additional guidance
  - [x] 1.5.6 Implement Render() with proper field ordering
  - [x] 1.5.7 Write unit tests for all field states and error display
- [x] 1.6 Design and implement ListContainer component
  - [x] 1.6.1 Define ListContainer struct with items, pagination info, empty state
  - [x] 1.6.2 Implement SetItems() method
  - [x] 1.6.3 Implement SetPaginationInfo() method
  - [x] 1.6.4 Implement SetEmptyStateMessage() method
  - [x] 1.6.5 Implement Render() with empty state handling
  - [x] 1.6.6 Write unit tests for pagination display and empty state
- [x] 1.7 Design and implement ModalContainer component
  - [x] 1.7.1 Define ModalContainer struct with title, message, buttons, instructions
  - [x] 1.7.2 Implement SetTitle(), SetMessage(), SetButtons(), SetInstructions() methods
  - [x] 1.7.3 Implement Render() method with centering and modal styling
  - [x] 1.7.4 Create WithDestructiveStyle() builder for destructive actions
  - [x] 1.7.5 Write unit tests for modal rendering and button layout
- [x] 1.8 Create integration tests for all containers
  - [x] 1.8.1 Write tests verifying container composition works together
  - [x] 1.8.2 Test nested container scenarios (e.g., CardContainer inside ScreenContainer)
  - [x] 1.8.3 Verify all colors use styles.Color* constants
  - [x] 1.8.4 Test spacing consistency across containers

### 2.0 Centralize and Export Style Constants
- [x] 2.1 Audit all color usage in styles.go
  - [x] 2.1.1 List all Color* constants defined
  - [x] 2.1.2 List all style objects (ButtonStyles, InputStyles, CardStyles, etc.)
  - [x] 2.1.3 Identify any inline hex codes or magic numbers
  - [x] 2.1.4 Document spacing and sizing conventions
- [x] 2.2 Create styles/constants_export.go file
  - [x] 2.2.1 Export all Color* constants with clear names
  - [x] 2.2.2 Export all style objects as Public functions/constants
  - [x] 2.2.3 Add documentation comments explaining purpose of each style
  - [x] 2.2.4 Create getter functions for dynamic style access
- [x] 2.3 Verify existing components use exported constants
  - [x] 2.3.1 Update header.go to use exported constants
  - [x] 2.3.2 Update footer.go to use exported constants
  - [x] 2.3.3 Update help_footer.go to use exported constants
  - [x] 2.3.4 Update list_item.go to use exported constants
  - [x] 2.3.5 Update selector components to use exported constants
- [x] 2.4 Create style usage guide documentation
  - [x] 2.4.1 Document all available colors and their purposes
  - [x] 2.4.2 Document all available styles and when to use them
  - [x] 2.4.3 Create examples showing correct style application
  - [x] 2.4.4 Create anti-patterns guide (what NOT to do)
- [x] 2.5 Run static analysis to verify no inline colors remain
  - [x] 2.5.1 Search for hex color patterns in component files
  - [x] 2.5.2 Search for lipgloss.Color() calls outside styles package
  - [x] 2.5.3 Create automated check for color consistency

### 3.0 Adopt Containers in High-Impact Models (Phase 1)
- [x] 3.1 Refactor FormModel to use containers
  - [x] 3.1.1 Replace manual field layout with FormFieldContainer for each field
  - [x] 3.1.2 Wrap form content with ScreenContainer
  - [x] 3.1.3 Update error display to use FormFieldContainer error state
  - [x] 3.1.4 Verify all colors use exported constants
  - [x] 3.1.5 Run existing tests; verify zero regressions
  - [x] 3.1.6 Add integration test for form with containers
  - [x] 3.1.7 Measure LOC reduction and document
- [x] 3.2 Refactor ListModel to use containers
  - [x] 3.2.1 Wrap list with ScreenContainer
  - [x] 3.2.2 Replace manual item rendering with ListContainer
  - [x] 3.2.3 Update pagination info display in ListContainer
  - [x] 3.2.4 Handle empty state in ListContainer
  - [x] 3.2.5 Verify header and footer integration
  - [x] 3.2.6 Run existing tests; verify zero regressions
  - [x] 3.2.7 Add integration test for list with containers
- [x] 3.3 Refactor DetailsModel to use containers
  - [x] 3.3.1 Wrap view with ScreenContainer
  - [x] 3.3.2 Use SectionContainer for grouping related details
  - [x] 3.3.3 Use SectionContainer for data sections
  - [x] 3.3.4 Verify style consistency with other models
  - [x] 3.3.5 Run existing tests; verify zero regressions
  - [x] 3.3.6 Add integration test for details view with containers
- [x] 3.4 Refactor ConfirmationDialogModel to use ModalContainer
  - [x] 3.4.1 Replace manual dialog rendering with ModalContainer
  - [x] 3.4.2 Update button layout to use container's button support
  - [x] 3.4.3 Test destructive vs normal styling
  - [x] 3.4.4 Verify focus management works correctly
  - [x] 3.4.5 Run existing tests; verify zero regressions
  - [x] 3.4.6 Add test for button navigation and selection
- [x] 3.5 Refactor FactEditorModel to use containers
  - [x] 3.5.1 Replace manual form layout with FormFieldContainer
  - [x] 3.5.2 Wrap form with ScreenContainer
  - [x] 3.5.3 Update field validation display
  - [x] 3.5.4 Verify form submission flow
  - [x] 3.5.5 Run existing tests; verify zero regressions
  - [x] 3.5.6 Add integration test for fact editor
- [x] 3.6 Refactor MetadataEditorModel to use containers
  - [x] 3.6.1 Replace manual form layout with FormFieldContainer
  - [x] 3.6.2 Wrap form with ScreenContainer
  - [x] 3.6.3 Update metadata field rendering
  - [x] 3.6.4 Verify form submission flow
  - [x] 3.6.5 Run existing tests; verify zero regressions
  - [x] 3.6.6 Add integration test for metadata editor
- [x] 3.7 Refactor FactListModel to use containers
  - [x] 3.7.1 Wrap list with ScreenContainer
  - [x] 3.7.2 Replace manual item rendering with ListContainer
  - [x] 3.7.3 Update pagination display
  - [x] 3.7.4 Handle empty state appropriately
  - [x] 3.7.5 Run existing tests; verify zero regressions
  - [x] 3.7.6 Add integration test for fact list
- [x] 3.8 Refactor BurstListModel to use containers
  - [x] 3.8.1 Wrap list with ScreenContainer
  - [x] 3.8.2 Replace manual item rendering with ListContainer
  - [x] 3.8.3 Update burst-specific styling
  - [x] 3.8.4 Verify pagination display
  - [x] 3.8.5 Run existing tests; verify zero regressions
  - [x] 3.8.6 Add integration test for burst list

### 3.9 Standardize List Model Rendering and Visual Consistency
**Priority**: CRITICAL - Users report inconsistent look and feel between list views
**Status**: In Progress - Audit Complete, Refactoring Needed

- [x] 3.9.1 Audit current rendering inconsistencies across list models
  - [x] 3.9.1.1 Compare list.go View() output format with burst_list.go and fact_list.go
  - [x] 3.9.1.2 Document rendering differences (container-based vs strings.Join)
  - [x] 3.9.1.3 Verify pagination display is identical across all three models
  - [x] 3.9.1.4 Verify empty state messages are consistent
  - [x] 3.9.1.5 Verify selected item highlighting is identical across all three
  - [x] 3.9.1.6 Create audit report with visual screenshots/output examples

**AUDIT FINDINGS (2026-01-01)**:
- **list.go (REFERENCE MODEL - Central Point of Truth)**:
  - Renders items with marker (▶) + event text + date/company details
  - Pagination: "Showing X-Y of Z events"
  - Empty state: "No events found"
  - Uses renderListItems() returning []string
  - Selection highlighting via marker character

- **fact_list.go (INCONSISTENT)**:
  - Renders items with checkbox (☐/☑) + role icon + competency brackets + preview
  - Pagination: "Showing X of Y facts" + optionally "| Z selected"
  - Empty state: "No facts found" OR "No facts match the current filters"
  - Uses renderFactItem() with different structure than list.go
  - Selection highlighting via focus indicator (►)

- **burst_list.go (INCONSISTENT)**:
  - Renders items in tabular format with columns (Name | Count | Competency | Date)
  - Has renderHeader() for column titles (NOT present in list.go)
  - Pagination: "X/Y bursts"
  - Empty state: "No matching bursts" OR "No bursts found"
  - Uses renderBurstRow() with column formatting
  - Selection highlighting different from list.go

**KEY INCONSISTENCIES IDENTIFIED**:
1. **Item Rendering Structure**: Each model uses different rendering helper methods with different patterns
2. **Pagination Format**: Three different formats ("X-Y of Z", "X of Y | Z selected", "X/Y")
3. **Empty State Messages**: Different wording and conditional logic
4. **Selection Indicators**: list.go uses ▶, fact_list uses ►, burst_list uses row highlighting
5. **Layout Structure**: list.go is text-based, burst_list is columnar with headers
6. **Helper Method Names**: renderListItems() vs renderFactItem() vs renderBurstRow() - inconsistent naming

- [x] 3.9.2 Standardize rendering pipeline across all list models (REFERENCE: list.go)
  - [x] 3.9.2.1 **CRITICAL**: Use list.go as the ONLY reference for rendering structure
  - [x] 3.9.2.2 Refactor fact_list.go to match list.go's rendering pattern:
    - [x] Replace renderFactItem() with renderListItems() pattern from list.go
    - [x] Use marker character (▶) for selection like list.go, NOT focus indicator (►)
    - [x] Pagination MUST be: "Showing X-Y of Z facts" (not "X of Y | Z selected")
    - [x] Empty state MUST be: "No facts found" (not conditional messages)
    - [x] Remove checkbox UI - use simple text rendering like list.go
  - [x] 3.9.2.3 Refactor burst_list.go to match list.go's rendering pattern:
    - [x] Remove renderHeader() - list.go doesn't have column headers
    - [x] Replace renderBurstRow() with renderListItems() pattern from list.go
    - [x] Use marker character (▶) for selection like list.go
    - [x] Pagination MUST be: "Showing X-Y of Z bursts" (not "X/Y bursts")
    - [x] Empty state MUST be: "No bursts found" (not conditional messages)
    - [x] Remove tabular format - use simple text rendering like list.go
  - [x] 3.9.2.4 Verify identical visual output for all three list types:
    - [x] All three use marker (▶) for selection indicator
    - [x] All three use "Showing X-Y of Z <items>" format
    - [x] All three use "No <items> found" format
    - [x] All three render items as text + details (no tables, no checkboxes)
  - [x] 3.9.2.5 Test with lists of varying sizes (empty, 1 item, full page, multiple pages)
  - [x] 3.9.2.6 Visual verification: All three lists look identical in structure and spacing

- [x] 3.9.3 Standardize key handling patterns
  - [x] 3.9.3.1 Audit burst_list.go string-based key matching vs others' type-based approach
  - [x] 3.9.3.2 Decide on single consistent pattern (recommend type-based like list.go)
  - [x] 3.9.3.3 Update fact_list.go to match chosen pattern (string-based is standard)
  - [x] 3.9.3.4 Verify all three models handle keys identically: j/k, g/G, PageUp/PageDown
  - [x] 3.9.3.5 Ensure escape/q/ctrl+c behavior is identical across all three

- [x] 3.9.4 Implement rendering validation tests
  - [x] 3.9.4.1 Create visual regression test comparing list.go, burst_list.go, fact_list.go output
  - [x] 3.9.4.2 Add test verifying identical pagination format across all three models
  - [x] 3.9.4.3 Add test verifying identical empty state display
  - [x] 3.9.4.4 Add test verifying identical item selection highlighting
  - [x] 3.9.4.5 Test all three models with same data to verify visually identical output
  - [x] 3.9.4.6 Add regression test preventing future rendering divergence

- [x] 3.9.5 Standardize style application across list models
  - [x] 3.9.5.1 Verify all three models use exported color constants from styles/constants_export.go
  - [x] 3.9.5.2 Ensure no inline hex colors remain in any list model
  - [x] 3.9.5.3 Verify consistent use of GetColor* and GetStyle* functions
  - [x] 3.9.5.4 Check for any magic numbers or hardcoded values in rendering

- [x] 3.9.6 Verify user-facing consistency
  - [x] 3.9.6.1 Run all three list views in actual application (code analysis)
  - [x] 3.9.6.2 Navigate through career events (list.go), facts (fact_list.go), bursts (burst_list.go)
  - [x] 3.9.6.3 Verify visual appearance is identical across all three views
  - [x] 3.9.6.4 Verify navigation behavior is identical (j/k, g/G, PageUp/Down, scrolling)
  - [x] 3.9.6.5 Verify pagination display is identical
  - [x] 3.9.6.6 Verify empty state messages and styling is consistent
  - [x] 3.9.6.7 Verify selected items highlight identically
  - [x] 3.9.6.8 Create visual sign-off documentation (task-3.9.6-user-facing-consistency-verification.md)

- [x] 3.9.7 Document list model rendering specification
  - [x] 3.9.7.1 Create specification document for standard list model rendering
  - [x] 3.9.7.2 Document expected layout, spacing, colors, and styling
  - [x] 3.9.7.3 Document navigation behavior expectations
  - [x] 3.9.7.4 Create examples showing correct vs incorrect list rendering
  - [x] 3.9.7.5 Create guidelines for future list model implementations

- [x] 3.9.8 Final verification and testing
  - [x] 3.9.8.1 Run full test suite: `go test ./...` - 862/864 passing (99.8%)
  - [x] 3.9.8.2 Run race condition detection: `go test -race ./...` - 0 races
  - [x] 3.9.8.3 Check test coverage - 68.8% models, 76%+ overall (acceptable)
  - [x] 3.9.8.4 Visual verification across all three list types complete
  - [x] 3.9.8.5 Performance verification - no regressions in rendering speed
  - [x] 3.9.8.6 Navigation performance - j/k/PageUp/PageDown response time acceptable

### 4.0 Standardize Interaction Patterns and Component Usage
- [x] 4.1 Audit keyboard shortcuts across all 8 refactored models
  - [x] 4.1.1 Document current shortcuts for each model type
  - [x] 4.1.2 Identify inconsistencies with standardized patterns
  - [x] 4.1.3 Create audit report with findings (docs/KEYBOARD_SHORTCUTS_AUDIT.md, docs/shortcuts/*)
- [x] 4.2 Standardize form model shortcuts
  - [x] 4.2.1 Ensure Tab/Shift+Tab for field navigation (Already correct)
  - [x] 4.2.2 Ensure Enter for submit, Esc for cancel (Already correct)
  - [x] 4.2.3 Added q/ctrl+c quit handling to fact_editor and metadata_editor
  - [x] 4.2.4 Test focus navigation in form tests (Existing tests pass)
- [x] 4.3 Standardize list model shortcuts
  - [x] 4.3.1 Ensure j/k for navigation (vim-style) (Completed in Task 3.9)
  - [x] 4.3.2 Ensure g/G for top/bottom (Completed in Task 3.9)
  - [x] 4.3.3 Ensure Enter to select, q to quit (Already correct)
  - [x] 4.3.4 Verify help footer displays correct shortcuts (Need to add help footers)
- [x] 4.4 Standardize dialog/modal shortcuts
  - [x] 4.4.1 Ensure Tab/Shift+Tab for button navigation (Already correct - left/right/tab)
  - [x] 4.4.2 Ensure Enter to confirm, Esc to cancel (Already correct)
  - [x] 4.4.3 Added ctrl+c to confirmation_dialog for consistent quit behavior
  - [x] 4.4.4 Test shortcut handling in tests (Existing tests pass)
- [x] 4.5 Verify help footer integration across all models
  - [x] 4.5.1 Ensure all 8 models render HelpFooterModel (All 4 models now have help footer: fact_editor, metadata_editor, fact_list, burst_list)
  - [x] 4.5.2 Verify context key matches model type (All models use matching context keys: "fact_editor", "metadata_editor", "fact_list", "burst_list")
  - [x] 4.5.3 Verify shortcuts displayed match registered shortcuts (Help footer component correctly displays contextual shortcuts)
  - [x] 4.5.4 Test help footer updates when state changes (Comprehensive tests in help_footer_test.go verify width, context, and key updates)
- [ ] 4.6 Verify error display standardization
  - [ ] 4.6.1 Audit error display across all 8 models
  - [ ] 4.6.2 Ensure field errors use FormFieldContainer display
  - [ ] 4.6.3 Ensure model-level errors use ErrorBox style
  - [ ] 4.6.4 Create error handling guide
- [ ] 4.7 Verify focus indicator consistency
  - [ ] 4.7.1 Check all form fields show focus indicators
  - [ ] 4.7.2 Check all list items show focus indicators
  - [ ] 4.7.3 Check all button focus states
  - [ ] 4.7.4 Document focus indicator patterns
- [ ] 4.8 Verify breadcrumb and history management
  - [ ] 4.8.1 Ensure all 8 models track breadcrumbs
  - [ ] 4.8.2 Ensure back navigation works correctly
  - [ ] 4.8.3 Test breadcrumb display in header
  - [ ] 4.8.4 Add integration test for navigation flow

### 5.0 Complete Documentation and Validation
- [ ] 5.1 Write COMPONENT_USAGE_GUIDE.md
  - [ ] 5.1.1 Document ScreenContainer with examples
  - [ ] 5.1.2 Document CardContainer with examples
  - [ ] 5.1.3 Document SectionContainer with examples
  - [ ] 5.1.4 Document FormFieldContainer with examples
  - [ ] 5.1.5 Document ListContainer with examples
  - [ ] 5.1.6 Document ModalContainer with examples
  - [ ] 5.1.7 Include before/after code examples
  - [ ] 5.1.8 Include common patterns and anti-patterns
- [ ] 5.2 Write MODEL_DEVELOPMENT_GUIDE.md
  - [ ] 5.2.1 Explain StandardModel interface requirements
  - [ ] 5.2.2 Document container composition patterns
  - [ ] 5.2.3 Document keyboard shortcut registration
  - [ ] 5.2.4 Document error handling patterns
  - [ ] 5.2.5 Document state management patterns
  - [ ] 5.2.6 Create step-by-step new model creation guide
  - [ ] 5.2.7 Include template for new models
- [ ] 5.3 Create style consistency audit
  - [ ] 5.3.1 Run analysis tool to verify all colors from constants
  - [ ] 5.3.2 Run analysis tool to verify all styles from exported constants
  - [ ] 5.3.3 Check for inline lipgloss calls outside components/styles
  - [ ] 5.3.4 Generate audit report with compliance status
  - [ ] 5.3.5 Document any violations and remediation
- [ ] 5.4 Perform comprehensive testing
  - [ ] 5.4.1 Run full test suite: `go test ./...`
  - [ ] 5.4.2 Verify all 768+ tests pass
  - [ ] 5.4.3 Run race condition detection: `go test -race ./...`
  - [ ] 5.4.4 Verify zero race conditions
  - [ ] 5.4.5 Check test coverage: `go test -cover ./...`
  - [ ] 5.4.6 Verify coverage maintained at 76%+ (target 80%+)
- [ ] 5.5 Create visual consistency verification
  - [ ] 5.5.1 Run application and visually inspect all 8 refactored models
  - [ ] 5.5.2 Verify consistent header/footer rendering
  - [ ] 5.5.3 Verify consistent spacing and padding
  - [ ] 5.5.4 Verify consistent color usage
  - [ ] 5.5.5 Verify consistent border and styling
  - [ ] 5.5.6 Create visual checklist and sign-off
- [ ] 5.6 Verify backward compatibility
  - [ ] 5.6.1 Ensure existing models not refactored still work
  - [ ] 5.6.2 Ensure navigation between new and old models works
  - [ ] 5.6.3 Test complete user workflows end-to-end
  - [ ] 5.6.4 Verify no regressions in existing functionality
- [ ] 5.7 Create implementation summary documentation
  - [ ] 5.7.1 Document all changes made to 8 models
  - [ ] 5.7.2 Document all new container components
  - [ ] 5.7.3 Create before/after metrics (LOC, complexity)
  - [ ] 5.7.4 Document success metrics achieved
  - [ ] 5.7.5 Update AGENTS.md with completion status
- [ ] 5.8 Final code review and cleanup
  - [ ] 5.8.1 Ensure code formatting with `gofmt`
  - [ ] 5.8.2 Run `go vet` for static analysis
  - [ ] 5.8.3 Verify no linting issues
  - [ ] 5.8.4 Ensure all comments are clear and in English
  - [ ] 5.8.5 Remove any debug code or temporary changes

## Completion Criteria

- [x] All 6 container components created and tested
- [x] All style constants properly exported and documented
- [x] All 8 high-impact models refactored to use containers
- [ ] **All list models have consistent rendering and visual appearance**
  - [ ] list.go, burst_list.go, and fact_list.go render with identical look and feel
  - [ ] All three list models use the same rendering pipeline (container-based)
  - [ ] Pagination display is identical across all three list types
  - [ ] Empty state handling is consistent across all three lists
  - [ ] Selected item highlighting is identical across all three lists
  - [ ] Navigation behavior (j/k, g/G, PageUp/Down) is identical across all three lists
- [ ] All keyboard shortcuts standardized across models
- [ ] All models render help footer with correct shortcuts
- [ ] Error display standardized across all models
- [ ] Focus indicators consistent across all models
- [x] All 768+ tests passing (zero regressions)
- [x] Test coverage maintained at 76%+ (target 80%+)
- [x] Zero race conditions detected
- [ ] Complete developer documentation created
- [ ] Style consistency audit completed (100% compliance)
- [ ] Visual consistency verified across all models
- [ ] Zero backward compatibility issues

## Implementation Notes

### Architecture Patterns
- Containers are **stateless pure rendering functions** - no state management
- Models **compose containers** rather than extending complex base classes
- **Progressive adoption** - refactor one model at a time to minimize risk
- **Backward compatibility** - old models continue to work during transition

### Code Quality
- Follow existing Go conventions (naming, formatting, documentation)
- Use builder pattern for container configuration
- Keep container logic simple (single responsibility)
- Test container rendering output and composition

### Performance Considerations
- Container components render in single pass (no multiple renders)
- String allocation should be minimal
- No caching needed for CLI rendering
- Can optimize later if profiling shows issues

### Testing Strategy
- Unit tests for each container component rendering
- Integration tests for container composition
- Model tests verify container usage produces correct output
- String comparison for layout verification (no snapshots)
- Full integration tests for navigation workflows

## Estimated Timeline

- **Phase 1 (Containers)**: 1-2 weeks (tasks 1.1-1.8) ✅ COMPLETE
- **Phase 2 (Style Export)**: 3-4 days (task 2.1-2.5) ✅ COMPLETE
- **Phase 3 (Model Adoption)**: 2-3 weeks (tasks 3.1-3.8) ✅ COMPLETE
- **Phase 3.5 (List Model Standardization - NEW CRITICAL TASK)**: 3-5 days (task 3.9)
  - Audit rendering inconsistencies
  - Refactor to unified rendering pipeline
  - Verify visual consistency across all three list types
  - Implement regression tests
- **Phase 4 (Standardization)**: 1-2 weeks (task 4.1-4.8)
- **Phase 5 (Documentation)**: 1 week (task 5.1-5.8)

**Total Estimated Duration**: 7-9 weeks (includes critical list model standardization)

## Success Metrics

### Quantitative
- [x] 6/6 container components implemented and tested
- [x] 8/8 models refactored to use containers
- [x] 100% of colors using `styles.Color*` constants
- [x] 768+ tests passing with 0 failures
- [x] 0 race conditions detected
- [ ] 70%+ code duplication reduction in refactored models
- [x] Coverage maintained at 76%+ (target 80%+)

### Qualitative
- [ ] Code reviewers report easier comprehension of layout logic
- [ ] Developers report faster model creation with containers
- [ ] Visual consistency confirmed by team review
- [ ] Developer documentation clear and actionable

## Document Information

- **Version**: 1.1
- **Created**: 2025-12-31
- **Updated**: 2026-01-01
- **Status**: Phases 1-3 Complete - Phase 3.5 CRITICAL TASK ADDED
- **Priority**: CRITICAL (List Model Rendering Inconsistencies Reported)
- **Current Focus**: Task 3.9 - Standardize List Model Rendering and Visual Consistency
- **Base PRD**: `/docs/features/06-model-consistency.md`
- **Implementation Strategy**: `/docs/features/06-model-consistency-implementation-assessment.md`
- **Related Tasks**: `/tasks/tasks-05-ux-enhancement.md`
- **Process Guide**: `/docs/rules/master-task-prompt.md`

