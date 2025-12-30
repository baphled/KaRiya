# Task List: TUI Standardization and Experience Enhancement

**PRD Reference**: `docs/features/04-tui-standarization.md`

**Purpose**: Standardize the Terminal User Interface (TUI) experience across all models by implementing consistent navigation patterns, unified keyboard shortcuts, standardized menu layouts, and improved visual coherence.

**Status**: 🔄 **IN PLANNING** (Ready for Phase 1 execution)

---

## Tasks

### Phase 1: Navigation and Keyboard Standardization

#### 1.0 Create Navigation Constants and Shortcuts System
- [ ] 1.1 Create `internal/cli/navigation/constants.go` with unified keyboard shortcut definitions
- [ ] 1.2 Define NavigationKey type with constants:
  - [ ] `KeyBack = "Esc"` (replace Backspace with Escape)
  - [ ] `KeyUp = "↑/k"` (navigate up)
  - [ ] `KeyDown = "↓/j"` (navigate down)
  - [ ] `KeyLeft = "←/h"` (navigate left)
  - [ ] `KeyRight = "→/l"` (navigate right)
  - [ ] `KeySelect = "Enter"` (confirm/select)
  - [ ] `KeyToggle = "Space"` (toggle/expand)
  - [ ] `KeyFilter = "f"` (show filters)
  - [ ] `KeySort = "s"` (show sort)
  - [ ] `KeySearch = "/"` (show search)
  - [ ] `KeyEdit = "e"` (edit item)
  - [ ] `KeyDelete = "d"` (delete item)
  - [ ] `KeyHelp = "?"` (show help)
  - [ ] `KeyHome = "h"` (go home)
  - [ ] `KeyQuit = "q"` (quit application)
  - [ ] `KeyBulk = "b"` (bulk operations)
  - [ ] `KeyCapture = "c"` (capture event)
  - [ ] `KeyList = "l"` (list events)
  - [ ] `KeyMetadata = "m"` (metadata review)
- [ ] 1.3 Create `HelpText` struct with formatted help strings for each navigation key
- [ ] 1.4 Implement `GetHelpText(keys []NavigationKey) -> string` function
- [ ] 1.5 Write unit tests for navigation constants (verify all keys defined)
- [ ] 1.6 Write unit tests for help text generation

#### 2.0 Create Navigation Menu Component
- [ ] 2.1 Create `internal/cli/components/navigation_menu.go` as reusable component
- [ ] 2.2 Implement NavigationMenuModel with configurable menu items
- [ ] 2.3 Support menu item structure: Label, Shortcut, Description, Action
- [ ] 2.4 Implement rendering with consistent styling (title, items, footer)
- [ ] 2.5 Support horizontal and vertical menu layouts
- [ ] 2.6 Implement keyboard navigation (↑/↓/←/→ for selection, Enter to select)
- [ ] 2.7 Add visual indicators for selected menu item (highlight, arrow)
- [ ] 2.8 Implement mouse support (click to select)
- [ ] 2.9 Write comprehensive unit tests (layout, navigation, selection)
- [ ] 2.10 Test edge cases (single item, many items, responsive sizing)

#### 3.0 Create Unified Help/Instructions Component
- [ ] 3.1 Create `internal/cli/components/help_footer.go` for consistent help display
- [ ] 3.2 Implement HelpFooterModel with configurable shortcuts
- [ ] 3.3 Support dynamic help text based on current context
- [ ] 3.4 Implement compact and full help formats
- [ ] 3.5 Use consistent styling from styles package
- [ ] 3.6 Handle responsive sizing (truncate on small terminals)
- [ ] 3.7 Write unit tests for help footer rendering
- [ ] 3.8 Test with various shortcut combinations

#### 4.0 Replace Backspace with Escape Key Globally
- [ ] 4.1 Update `internal/cli/models/form.go` to use Escape instead of Backspace
- [ ] 4.2 Update `internal/cli/models/metadata_editor.go` to use Escape
- [ ] 4.3 Update `internal/cli/models/metadata_review.go` to use Escape
- [ ] 4.4 Update `internal/cli/models/bulk_operations.go` to use Escape
- [ ] 4.5 Update `internal/cli/models/help.go` to use Escape
- [ ] 4.6 Update `internal/cli/models/import_review.go` to use Escape
- [ ] 4.7 Update `internal/cli/models/view_event.go` to use Escape
- [ ] 4.8 Update `internal/cli/models/action_menu.go` to use Escape
- [ ] 4.9 Update `internal/cli/models/details.go` to use Escape
- [ ] 4.10 Update all help text strings to show "Esc: Back" instead of "Backspace: Back"
- [ ] 4.11 Update all test files to use Escape key instead of Backspace
- [ ] 4.12 Verify app.go navigation handles Escape correctly

#### 5.0 Standardize Navigation Keys Across All Models
- [ ] 5.1 Update all models to use consistent j/k for up/down navigation
- [ ] 5.2 Update all models to use consistent h/l for left/right navigation
- [ ] 5.3 Update all models to use consistent shortcuts from navigation/constants.go
- [ ] 5.4 Verify consistency in form.go, list.go, metadata_review.go, etc.
- [ ] 5.5 Update all help text to reflect standardized shortcuts
- [ ] 5.6 Write integration tests for keyboard consistency across models
- [ ] 5.7 Test vim-style navigation (hjkl) works in all models

### Phase 2: Visual Consistency and Layout Standardization

#### 6.0 Create Unified Header Component
- [ ] 6.1 Create `internal/cli/components/header.go` for consistent screen headers
- [ ] 6.2 Implement HeaderModel with title, subtitle, status indicators
- [ ] 6.3 Support breadcrumb navigation display
- [ ] 6.4 Implement consistent styling using styles package
- [ ] 6.5 Support responsive sizing for various terminal widths
- [ ] 6.6 Write unit tests for header rendering
- [ ] 6.7 Test breadcrumb display and navigation context

#### 7.0 Create Unified Footer Component
- [ ] 7.1 Create `internal/cli/components/footer.go` for consistent screen footers
- [ ] 7.2 Implement FooterModel with status, mode, and help text
- [ ] 7.3 Display current mode/context (e.g., "Capture Mode: Timeline")
- [ ] 7.4 Display status messages (e.g., "3/10 events")
- [ ] 7.5 Integrate help footer component
- [ ] 7.6 Use consistent styling and layout
- [ ] 7.7 Write unit tests for footer rendering
- [ ] 7.8 Test with various status/help combinations

#### 8.0 Standardize Form Layout and Styling
- [ ] 8.1 Review form.go for consistent input field styling
- [ ] 8.2 Ensure all form fields use consistent label and input styling
- [ ] 8.3 Implement consistent error message display
- [ ] 8.4 Ensure consistent focus indicators (highlight, cursor, color)
- [ ] 8.5 Implement consistent field spacing and alignment
- [ ] 8.6 Add consistent character counter display (if applicable)
- [ ] 8.7 Test form rendering on various terminal sizes
- [ ] 8.8 Write unit tests for form layout consistency

#### 9.0 Standardize List Item Display
- [ ] 9.1 Create `internal/cli/components/list_item.go` for consistent list item rendering
- [ ] 9.2 Implement ListItemModel with title, subtitle, metadata fields
- [ ] 9.3 Support icons/indicators for item status or type
- [ ] 9.4 Implement consistent selected/focused item styling
- [ ] 9.5 Support truncation for long text (with ellipsis)
- [ ] 9.6 Implement consistent spacing and alignment
- [ ] 9.7 Write unit tests for list item rendering
- [ ] 9.8 Test with various data lengths and terminal sizes

#### 10.0 Standardize Modal/Dialog Styling
- [ ] 10.1 Review confirmation_dialog.go for consistent modal styling
- [ ] 10.2 Ensure consistent border, padding, and background styling
- [ ] 10.3 Implement consistent button styling within modals
- [ ] 10.4 Ensure consistent text alignment and spacing
- [ ] 10.5 Test modal rendering with various content lengths
- [ ] 10.6 Write unit tests for modal styling consistency
- [ ] 10.7 Test modals on various terminal sizes

### Phase 3: Model Integration and Refactoring

#### 11.0 Integrate Navigation System into Form Model
- [ ] 11.1 Update form.go to use NavigationKey constants
- [ ] 11.2 Implement help footer with standardized shortcuts
- [ ] 11.3 Add header component with title and context
- [ ] 11.4 Ensure consistent keyboard handling
- [ ] 11.5 Update all help text to use standardized format
- [ ] 11.6 Write integration tests for form with new navigation
- [ ] 11.7 Verify backward compatibility with existing tests

#### 12.0 Integrate Navigation System into List Model
- [ ] 12.1 Update list.go to use NavigationKey constants
- [ ] 12.2 Implement help footer with list-specific shortcuts
- [ ] 12.3 Add header component with title and status
- [ ] 12.4 Integrate list item component for consistent display
- [ ] 12.5 Ensure vim-style navigation (hjkl) works
- [ ] 12.6 Write integration tests for list with new navigation
- [ ] 12.7 Verify backward compatibility with existing tests

#### 13.0 Integrate Navigation System into Metadata Review Model
- [ ] 13.1 Update metadata_review.go to use NavigationKey constants
- [ ] 13.2 Implement help footer with metadata review shortcuts
- [ ] 13.3 Add header component with title and progress
- [ ] 13.4 Integrate list item component for event display
- [ ] 13.5 Ensure consistent keyboard handling
- [ ] 13.6 Write integration tests for metadata review with new navigation
- [ ] 13.7 Verify backward compatibility with existing tests

#### 14.0 Integrate Navigation System into Other Models
- [ ] 14.1 Update metadata_editor.go with navigation constants and help footer
- [ ] 14.2 Update bulk_operations.go with navigation constants and help footer
- [ ] 14.3 Update view_event.go with navigation constants and help footer
- [ ] 14.4 Update action_menu.go with navigation constants and help footer
- [ ] 14.5 Update help.go with navigation constants and help footer
- [ ] 14.6 Update import_review.go with navigation constants and help footer
- [ ] 14.7 Update details.go with navigation constants and help footer
- [ ] 14.8 Update success.go with navigation constants and help footer
- [ ] 14.9 Write integration tests for all updated models
- [ ] 14.10 Verify backward compatibility with existing tests

#### 15.0 Refactor Common UI Patterns into Reusable Components
- [ ] 15.1 Identify common UI patterns across models (scrollable list, form input, menu)
- [ ] 15.2 Extract scrollable list pattern into reusable component
- [ ] 15.3 Extract form input pattern into reusable component
- [ ] 15.4 Extract menu pattern into reusable component
- [ ] 15.5 Update models to use new reusable components
- [ ] 15.6 Ensure consistent behavior across all models using components
- [ ] 15.7 Write unit tests for new components
- [ ] 15.8 Verify all models work correctly with new components

### Phase 4: Visual Enhancements and Polish

#### 16.0 Enhance Breadcrumb and Navigation Context
- [ ] 16.1 Implement breadcrumb display in header component
- [ ] 16.2 Show current navigation path (e.g., "Home > List > Event Details")
- [ ] 16.3 Allow breadcrumb navigation (click to go back)
- [ ] 16.4 Implement consistent breadcrumb styling
- [ ] 16.5 Write unit tests for breadcrumb display and navigation
- [ ] 16.6 Test with various navigation paths

#### 17.0 Add Status and Progress Indicators
- [ ] 17.1 Add progress indicator for multi-step workflows
- [ ] 17.2 Display current step (e.g., "Step 2 of 5")
- [ ] 17.3 Implement progress bar visualization
- [ ] 17.4 Show status messages in footer
- [ ] 17.5 Implement status colors (success, warning, error, info)
- [ ] 17.6 Write unit tests for progress indicators
- [ ] 17.7 Test with various workflow lengths

#### 18.0 Implement Consistent Color Scheme and Theming
- [ ] 18.1 Review styles.go for color consistency
- [ ] 18.2 Ensure all colors follow professional dark theme
- [ ] 18.3 Verify adequate contrast for accessibility
- [ ] 18.4 Test colors on various terminal backgrounds
- [ ] 18.5 Document color scheme and usage guidelines
- [ ] 18.6 Write unit tests for color consistency
- [ ] 18.7 Create color scheme reference documentation

#### 19.0 Add Visual Feedback for User Actions
- [ ] 19.1 Implement loading spinners for long operations
- [ ] 19.2 Add success/error messages with visual indicators
- [ ] 19.3 Implement animations for state transitions (if applicable)
- [ ] 19.4 Add visual feedback for button presses
- [ ] 19.5 Implement hover/focus effects where applicable
- [ ] 19.6 Write unit tests for visual feedback
- [ ] 19.7 Test feedback on various terminal types

### Phase 5: Testing and Documentation

#### 20.0 Comprehensive Navigation Testing
- [ ] 20.1 Write tests for all keyboard shortcuts across all models
- [ ] 20.2 Test Escape key works as back button everywhere
- [ ] 20.3 Test vim-style navigation (hjkl) in all navigation contexts
- [ ] 20.4 Test navigation consistency between models
- [ ] 20.5 Test navigation edge cases (first/last item, empty lists)
- [ ] 20.6 Test rapid key presses don't cause issues
- [ ] 20.7 Write integration tests for complete navigation workflows

#### 21.0 Visual Consistency Testing
- [ ] 21.1 Test header rendering on various terminal sizes
- [ ] 21.2 Test footer rendering on various terminal sizes
- [ ] 21.3 Test list items render consistently
- [ ] 21.4 Test forms render consistently
- [ ] 21.5 Test modals render consistently
- [ ] 21.6 Test responsive layout on narrow terminals (80 cols)
- [ ] 21.7 Test responsive layout on wide terminals (200+ cols)
- [ ] 21.8 Take screenshots for visual regression testing

#### 22.0 Performance and Stability Testing
- [ ] 22.1 Run race detector: `go test -race ./...`
- [ ] 22.2 Verify no memory leaks with large event lists
- [ ] 22.3 Test rendering performance with many items
- [ ] 22.4 Verify smooth scrolling through large lists
- [ ] 22.5 Test rapid navigation doesn't cause lag
- [ ] 22.6 Verify code coverage meets 80%+ threshold
- [ ] 22.7 Test error handling and recovery

#### 23.0 Documentation and User Guidance
- [ ] 23.1 Create TUI_STANDARDS.md documenting design standards
- [ ] 23.2 Document all keyboard shortcuts in standardized table format
- [ ] 23.3 Create keyboard reference card (printable)
- [ ] 23.4 Update README.md with standardized navigation overview
- [ ] 23.5 Update CLI_GUIDE.md with consistent shortcut documentation
- [ ] 23.6 Create developer guide for adding new screens/models
- [ ] 23.7 Document component reuse patterns and guidelines
- [ ] 23.8 Update CHANGELOG.md with standardization changes
- [ ] 23.9 Add in-app help (?) showing current context shortcuts
- [ ] 23.10 Document accessibility features and considerations

---

## Implementation Guidelines

### Architecture Patterns

1. **Separation of Concerns**
   - Navigation constants in `internal/cli/navigation/`
   - Reusable components in `internal/cli/components/`
   - Model-specific logic in `internal/cli/models/`
   - Styling in `internal/cli/styles/`

2. **Component Reusability**
   - Create small, focused components (header, footer, list item)
   - Compose larger screens from smaller components
   - Share common patterns across models
   - DRY principle for styling and layout

3. **Backward Compatibility**
   - Maintain existing model interfaces
   - Support both Escape and Backspace during transition
   - Gradual migration to new navigation system
   - Tests should pass before and after refactoring

4. **Consistency**
   - All models follow same keyboard shortcuts
   - All models use same help text format
   - All models use same header/footer style
   - All models use same color scheme

### Key Files to Create

**Navigation System**:
- `internal/cli/navigation/constants.go` (keyboard shortcuts)
- `internal/cli/navigation/help.go` (help text generation)

**Reusable Components**:
- `internal/cli/components/header.go` (screen header)
- `internal/cli/components/footer.go` (screen footer)
- `internal/cli/components/help_footer.go` (help text footer)
- `internal/cli/components/navigation_menu.go` (menu component)
- `internal/cli/components/list_item.go` (list item display)

**Documentation**:
- `docs/TUI_STANDARDS.md` (design standards)
- `docs/KEYBOARD_REFERENCE.md` (keyboard shortcuts)
- `docs/TUI_DEVELOPER_GUIDE.md` (for developers)

### Success Criteria (All Must Be Met)

- [x] Escape key works as back button everywhere
- [x] vim-style navigation (hjkl) available in all models
- [x] All models use consistent keyboard shortcuts
- [x] All models have consistent help footer
- [x] All models have consistent header with context
- [x] All models use consistent styling and colors
- [x] List items render consistently across all lists
- [x] Forms render consistently across all forms
- [x] Modals render consistently across all modals
- [x] Navigation is intuitive and discoverable
- [x] Help text is clear and helpful
- [x] Code is well-organized and reusable
- [x] All tests passing (100% pass rate)
- [x] Race detector passes (0 conditions)
- [x] Code coverage maintained at 80%+
- [x] Documentation is comprehensive
- [x] Performance is smooth on all terminal sizes

---

## Phase Dependencies

This feature builds on:
- **Phases 1-3**: Existing CLI implementation ✅
- **Phase 3**: Metadata review and bulk operations ✅

This feature enables:
- **Burst/Fact Extraction**: Cleaner UI for new features
- **Portfolio Generation**: Consistent UI for new workflows
- **Future Features**: Strong foundation for new screens

---

## Estimated Effort

- Phase 1: 6-8 hours (navigation constants, keyboard standardization)
- Phase 2: 8-10 hours (visual components, layout standardization)
- Phase 3: 10-12 hours (model integration, refactoring)
- Phase 4: 6-8 hours (visual enhancements, polish)
- Phase 5: 8-10 hours (testing, documentation)

**Total**: 38-48 hours

---

## Completion Tracking

- **Phase 1**: ⏳ Ready for execution
- **Phase 2**: ⏳ Awaiting Phase 1 completion
- **Phase 3**: ⏳ Awaiting Phase 2 completion
- **Phase 4**: ⏳ Awaiting Phase 3 completion
- **Phase 5**: ⏳ Awaiting Phase 4 completion

---

## Key Improvements

### Before Standardization
- Different models use Backspace, Escape, or both
- Inconsistent keyboard shortcuts across models
- Help text varies in format and completeness
- No consistent header/footer pattern
- Visual styling varies by model
- Navigation patterns differ by screen
- No centralized navigation configuration

### After Standardization
- Escape key consistently used for back navigation
- Unified keyboard shortcuts (hjkl, vim-style)
- Help text follows standardized format
- All screens have consistent header/footer
- Unified visual styling and color scheme
- Navigation patterns consistent across all screens
- Centralized navigation constants and help system
- Reusable components reduce code duplication
- Improved user experience and discoverability
- Better developer experience for adding new screens

---

**Document Version**: 1.0
**Created**: 2025-12-30
**Status**: Ready for Phase 1 Execution
**Template Source**: tasks-03-metadata-clarification.md
**Process Guide**: docs/rules/process-task-list.md

