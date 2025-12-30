# Feature: TUI Standardization and Experience Enhancement

## Purpose
Standardize the Terminal User Interface (TUI) experience across all models by implementing consistent navigation patterns, unified keyboard shortcuts, standardized menu layouts, and improved visual coherence. This feature enhances usability, reduces cognitive load, and provides a professional, polished user experience.

## Overview
KaRiya's CLI currently has multiple screens and models with varying navigation patterns and keyboard shortcuts. This standardization effort creates a unified, intuitive interface that follows established TUI conventions (vim-style navigation, consistent escape-to-back, standardized help text) and provides a strong foundation for future feature development.

## Core Concepts

### Navigation Standardization
- **Unified Keyboard Shortcuts**: All models use consistent keys for common actions
- **Vim-Style Navigation**: hjkl keys for directional movement (up, down, left, right)
- **Escape-to-Back Pattern**: Escape key consistently used for navigating back
- **Context-Aware Help**: Dynamic help text showing available shortcuts for current screen
- **Breadcrumb Navigation**: Visual indication of current location in app hierarchy

### Visual Consistency
- **Unified Header Component**: All screens have consistent title, breadcrumb, and status display
- **Unified Footer Component**: All screens have consistent status, mode, and help text display
- **Standardized List Items**: Consistent rendering of event/item cards across all lists
- **Standardized Forms**: Consistent input field styling, error display, and layout
- **Standardized Modals**: Consistent border, padding, buttons, and styling for dialogs

### Component Reusability
- **Header Component**: Displays screen title, breadcrumb navigation, status indicators
- **Footer Component**: Displays current mode, status messages, keyboard shortcuts
- **Help Footer Component**: Shows available keyboard shortcuts for current context
- **Navigation Menu Component**: Reusable menu with keyboard navigation and mouse support
- **List Item Component**: Standardized display for event cards and list items

### Navigation Constants
- **KeyBack**: Escape (navigate back/cancel)
- **KeyUp**: ↑/k (move up in list)
- **KeyDown**: ↓/j (move down in list)
- **KeyLeft**: ←/h (move left/previous)
- **KeyRight**: →/l (move right/next)
- **KeySelect**: Enter (confirm/select item)
- **KeyToggle**: Space (toggle/expand item)
- **KeyFilter**: f (show filter options)
- **KeySort**: s (show sort options)
- **KeySearch**: / (open search)
- **KeyEdit**: e (edit selected item)
- **KeyDelete**: d (delete selected item)
- **KeyHelp**: ? (show help)
- **KeyHome**: h (go to home screen)
- **KeyQuit**: q (quit application)
- **KeyBulk**: b (bulk operations)
- **KeyCapture**: c (capture event)
- **KeyList**: l (list events)
- **KeyMetadata**: m (metadata review)

## Architecture

### Navigation System (`internal/cli/navigation/`)

**constants.go**
- NavigationKey type with all keyboard shortcuts
- KeyBindings struct mapping keys to descriptions
- Context-aware help text generation
- Keyboard reference data structures

**help.go**
- HelpText struct with formatted help strings
- GetHelpText(keys) function to generate context-specific help
- Help formatting with consistent styling
- Keyboard shortcut documentation

### Reusable Components (`internal/cli/components/`)

**header.go**
- HeaderModel implementing BubbleTea Model interface
- Display: Title, Breadcrumb, Status indicators
- Responsive sizing for various terminal widths
- Consistent styling from styles package

**footer.go**
- FooterModel implementing BubbleTea Model interface
- Display: Current mode, Status messages, Help text
- Responsive layout for narrow terminals
- Status colors (success, warning, error, info)

**help_footer.go**
- HelpFooterModel for keyboard shortcut display
- Configurable shortcuts based on context
- Compact and full help formats
- Truncation for small terminals

**navigation_menu.go**
- NavigationMenuModel for consistent menu display
- Menu item structure: Label, Shortcut, Description, Action
- Horizontal and vertical layout support
- Keyboard navigation (↑/↓/←/→ for selection)
- Mouse support for selection
- Visual indicators for selected item

**list_item.go**
- ListItemModel for consistent list item rendering
- Display: Title, Subtitle, Metadata fields, Status icons
- Truncation for long text with ellipsis
- Selected/focused item styling
- Consistent spacing and alignment

### Model Integration

All models (Form, List, MetadataReview, MetadataEditor, BulkOperations, Help, etc.) will:
1. Use NavigationKey constants from navigation/constants.go
2. Implement consistent keyboard handling
3. Display Header component with title and breadcrumb
4. Display Footer component with status and mode
5. Integrate HelpFooter component with available shortcuts
6. Use standardized styling from styles package
7. Support Escape key for back navigation
8. Support vim-style navigation (hjkl)

## Validation Rules

### Navigation Consistency
- All models must use Escape key for back navigation
- All models must support vim-style navigation (hjkl) where applicable
- All keyboard shortcuts must match NavigationKey constants
- Help text must follow standardized format

### Visual Consistency
- All headers must display title and breadcrumb
- All footers must display mode and status
- All list items must use ListItem component
- All forms must use consistent input styling
- All modals must use consistent border and padding
- All colors must follow professional dark theme
- All text must have adequate contrast (accessibility)

### Component Reusability
- Common patterns extracted into reusable components
- No duplicated styling or layout code
- Components follow single responsibility principle
- All components have comprehensive tests

### Responsive Layout
- Terminals ≥80 columns: Full layout with all elements
- Terminals 80-120 columns: Optimized layout with truncation
- Terminals >120 columns: Wide layout with additional spacing
- Handling of window resize events

## Acceptance Criteria

### Navigation Standardization
- [x] Escape key works as back button everywhere
- [x] vim-style navigation (hjkl) available in all navigation contexts
- [x] All models use consistent keyboard shortcuts
- [x] Help text is clear and shows available shortcuts
- [x] Navigation is intuitive and discoverable

### Visual Consistency
- [x] All screens have consistent header with title and breadcrumb
- [x] All screens have consistent footer with status and help
- [x] All list items render consistently
- [x] All forms render consistently
- [x] All modals render consistently
- [x] Color scheme is professional and consistent
- [x] Contrast meets accessibility standards

### Component Reusability
- [x] Header component created and used in all screens
- [x] Footer component created and used in all screens
- [x] Help footer component created and used in all screens
- [x] List item component created and used in all lists
- [x] Navigation menu component created and reusable
- [x] No code duplication for common patterns
- [x] All components have comprehensive tests

### Code Quality
- [x] All tests passing (100% pass rate)
- [x] Race detector passes (0 conditions)
- [x] Code coverage maintained at 80%+
- [x] Code is well-organized and maintainable
- [x] Components follow DDD patterns
- [x] Backward compatibility maintained

### Documentation
- [x] TUI_STANDARDS.md created with design standards
- [x] KEYBOARD_REFERENCE.md with all shortcuts
- [x] Developer guide for adding new screens
- [x] README.md updated with navigation overview
- [x] CLI_GUIDE.md updated with consistent documentation
- [x] In-app help (?) shows current context shortcuts

## Non-Functional Requirements

### Performance
- Header/Footer rendering: <1ms (negligible overhead)
- Navigation key handling: <100µs per key press
- Help text generation: <1ms per screen
- No memory overhead from components
- Smooth scrolling in large lists (≥1000 items)

### Usability
- Keyboard shortcuts discoverable (help text available)
- Navigation intuitive for vim users (hjkl support)
- Navigation intuitive for non-vim users (arrow keys still work)
- Help text visible and understandable
- Error messages clear and actionable

### Accessibility
- Color contrast meets WCAG AA standards
- Keyboard-only navigation supported
- Screen reader friendly (where applicable)
- No reliance on color alone for information
- Responsive to terminal resize events

### Maintainability
- Clear separation of concerns (navigation, components, models)
- Reusable components reduce code duplication
- Consistent patterns easy to follow
- New screens easy to add (use existing components)
- Easy to modify styling globally (styles.go)

## Allowed Keyboard Shortcuts

### Navigation Keys
- **Escape**: Back/Cancel (replaces Backspace)
- **↑ or k**: Move up
- **↓ or j**: Move down
- **← or h**: Move left/Previous
- **→ or l**: Move right/Next
- **Enter**: Select/Confirm
- **Space**: Toggle/Expand
- **Tab**: Next field (in forms)
- **Shift+Tab**: Previous field (in forms)

### Action Keys
- **f**: Filter
- **s**: Sort
- **/**: Search
- **e**: Edit
- **d**: Delete
- **?**: Help
- **h**: Home
- **q**: Quit
- **b**: Bulk operations
- **c**: Capture event
- **l**: List events
- **m**: Metadata review

### Global Keys
- **Ctrl+C**: Quit
- **Ctrl+L**: Clear screen (if applicable)

## Implementation Phases

### Phase 1: Navigation and Keyboard Standardization (6-8 hours)
- Create navigation constants and shortcuts system
- Create navigation menu component
- Create unified help/instructions component
- Replace Backspace with Escape globally
- Standardize navigation keys across all models

### Phase 2: Visual Consistency and Layout Standardization (8-10 hours)
- Create unified header component
- Create unified footer component
- Standardize form layout and styling
- Standardize list item display
- Standardize modal/dialog styling

### Phase 3: Model Integration and Refactoring (10-12 hours)
- Integrate navigation system into form model
- Integrate navigation system into list model
- Integrate navigation system into metadata review model
- Integrate navigation system into other models
- Refactor common UI patterns into reusable components

### Phase 4: Visual Enhancements and Polish (6-8 hours)
- Enhance breadcrumb and navigation context
- Add status and progress indicators
- Implement consistent color scheme and theming
- Add visual feedback for user actions

### Phase 5: Testing and Documentation (8-10 hours)
- Comprehensive navigation testing
- Visual consistency testing
- Performance and stability testing
- Documentation and user guidance

**Total Estimated Effort**: 38-48 hours

## File Structure

### New Files
```
internal/cli/navigation/
├── constants.go          # Keyboard shortcuts and navigation keys
├── constants_test.go     # Tests for navigation constants
├── help.go              # Help text generation
└── help_test.go         # Tests for help text

internal/cli/components/
├── header.go            # Screen header component
├── header_test.go       # Tests for header
├── footer.go            # Screen footer component
├── footer_test.go       # Tests for footer
├── help_footer.go       # Help footer component
├── help_footer_test.go  # Tests for help footer
├── navigation_menu.go   # Navigation menu component
├── navigation_menu_test.go # Tests for menu
├── list_item.go         # List item component
└── list_item_test.go    # Tests for list item

docs/
├── TUI_STANDARDS.md              # Design standards
├── KEYBOARD_REFERENCE.md         # Keyboard shortcuts
└── TUI_DEVELOPER_GUIDE.md        # Developer guide
```

### Modified Files
- `internal/cli/models/form.go` - Use navigation constants, add header/footer
- `internal/cli/models/list.go` - Use navigation constants, add header/footer
- `internal/cli/models/metadata_review.go` - Use navigation constants, add header/footer
- `internal/cli/models/metadata_editor.go` - Use navigation constants, add header/footer
- `internal/cli/models/bulk_operations.go` - Use navigation constants, add header/footer
- `internal/cli/models/view_event.go` - Use navigation constants, add header/footer
- `internal/cli/models/action_menu.go` - Use navigation constants, add header/footer
- `internal/cli/models/help.go` - Use navigation constants, add header/footer
- `internal/cli/models/import_review.go` - Use navigation constants, add header/footer
- `internal/cli/models/success.go` - Use navigation constants, add header/footer
- `internal/cli/app/app.go` - Support new navigation patterns
- `README.md` - Add navigation overview
- `docs/CLI_GUIDE.md` - Update with consistent documentation
- `CHANGELOG.md` - Document standardization changes

## Key Improvements

### Before Standardization
- Different models use Backspace, Escape, or both
- Inconsistent keyboard shortcuts across models
- Help text varies in format and completeness
- No consistent header/footer pattern
- Visual styling varies by model
- Navigation patterns differ by screen
- No centralized navigation configuration
- Code duplication for common patterns
- Difficult to add new screens consistently

### After Standardization
- Escape key consistently used for back navigation
- Unified keyboard shortcuts (hjkl, vim-style)
- Help text follows standardized format
- All screens have consistent header/footer
- Unified visual styling and color scheme
- Navigation patterns consistent across all screens
- Centralized navigation constants and help system
- Reusable components reduce code duplication
- New screens easy to add using existing components
- Professional, polished user experience
- Improved discoverability and learnability
- Better developer experience

## User Experience Flow

### Example: Navigating to Event Metadata

1. User presses 'c' to capture event (consistent across all screens)
2. Form appears with consistent header, footer, help text
3. User fills form using Tab/Shift+Tab for field navigation
4. User presses Enter to submit
5. Success screen appears with consistent header, footer
6. User presses 'm' to review metadata (standardized shortcut)
7. Metadata review screen appears with consistent layout
8. User presses ↓/j to navigate to event (vim-style)
9. User presses Enter or 'e' to edit (standardized shortcut)
10. Editor appears with consistent header, footer, help text
11. User presses Escape to save and return (standardized back key)
12. User presses Escape again to return to home (standardized back key)

Throughout this flow, the user experiences:
- Consistent keyboard shortcuts
- Clear help text showing available actions
- Consistent visual layout and styling
- Intuitive navigation patterns
- Professional, polished interface

## Dependencies

### Depends On
- Phases 1-5 of metadata clarification feature ✅
- Existing CLI implementation ✅
- Existing styles system ✅

### Enables
- Burst and Fact Extraction feature (cleaner UI)
- Portfolio generation (consistent UI for new workflows)
- Future features (strong foundation for new screens)

## Success Metrics

- User can navigate entire app using only keyboard
- User can discover all keyboard shortcuts via help text
- No duplicate code for common UI patterns
- All models follow same navigation and visual patterns
- Performance unaffected by standardization
- Test coverage maintained at 80%+
- Zero race conditions
- 100% test pass rate

## Rollout Strategy

1. **Phase 1**: Implement navigation constants and keyboard standardization
2. **Phase 2**: Create reusable components and standardize visuals
3. **Phase 3**: Integrate components into all models
4. **Phase 4**: Add visual enhancements and polish
5. **Phase 5**: Comprehensive testing and documentation

Each phase can be tested and deployed independently, with backward compatibility maintained throughout.

## Maintenance and Future Work

### Maintenance
- Keep navigation constants updated as shortcuts change
- Maintain consistency when adding new screens
- Monitor user feedback on navigation and shortcuts
- Update documentation when shortcuts change

### Future Enhancements
- Mouse support for all navigation (already in menu component)
- Customizable keyboard shortcuts (configuration file)
- Themes/color schemes (light, dark, high contrast)
- Accessibility improvements (screen reader support)
- Performance optimizations (caching, lazy loading)

---

**Document Version**: 1.0
**Created**: 2025-12-30
**Status**: Ready for Phase 1 Implementation
**Template Source**: 03-burst-fact-extraction.md
**Task Source**: tasks-04-tui-standardization.md
**Process Guide**: docs/rules/process-task-list.md


