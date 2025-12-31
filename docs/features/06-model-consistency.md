# Product Requirements Document: Comprehensive Model Consistency and Standardized UI Components

## 1. Introduction/Overview

KaRiya is a Career Journal CLI tool with multiple model types (forms, lists, data viewers, dialogs, menus) that currently have inconsistent styling, layout, and component integration patterns. This PRD establishes a comprehensive standardization framework that leverages the existing `lipgloss` styling system and `components` library to create a unified, maintainable UI across all model types.

The standardization effort will ensure that all models—regardless of type—follow consistent patterns for layout composition, component usage, styling application, and interaction behavior.

## 2. Goals

1. **Eliminate style duplication** by centralizing all UI layout patterns into reusable lipgloss-based container components
2. **Establish consistent interaction patterns** across all model types (forms, lists, dialogs, menus, utility models)
3. **Create a comprehensive component library** that models can compose rather than implement styling directly
4. **Improve code maintainability** by standardizing how models apply layouts, spacing, and visual hierarchy
5. **Ensure visual consistency** across the entire application through unified lipgloss-based styling
6. **Reduce cognitive load** for developers maintaining and extending models by providing clear, composable patterns

## 3. User Stories

1. As a user, I want every screen in the application to follow consistent layout patterns so I can quickly navigate and understand the information architecture.
2. As a user, I want consistent styling (colors, spacing, borders) across all screens to create a cohesive visual experience.
3. As a user, I want predictable keyboard interactions across forms, lists, dialogs, and menus so I don't need to relearn navigation patterns.
4. As a power user, I want consistent help text, error displays, and status indicators across all screens.
5. As a developer, I want to create new models quickly by composing existing UI components rather than implementing layouts from scratch.
6. As a developer, I want to update styling globally without touching individual model implementations.
7. As a developer, I want clear patterns for handling focus states, error states, and disabled states consistently across components.

## 4. Current Implementation Strengths

### Existing Styling Infrastructure (`internal/cli/styles/styles.go`)
- **Comprehensive color palette**: Professional dark theme with muted accents (teal, green, purple)
- **Predefined style objects**: Button styles (Primary, Secondary, Focused, Disabled), Input styles (Base, Focused, Error), Card styles, Modal styles, Header styles, Error/Warning/Success styles
- **Consistent spacing and borders**: All styles use rounded borders and consistent padding
- **Text hierarchy support**: Primary, Secondary, and Muted text colors; Bold, Italic, Faint variants
- **Status color coding**: Error (red), Warning (amber), Success (green), Info (blue)

### Existing Component Library (`internal/cli/components/`)
- **HeaderModel**: Renders title, subtitle, and breadcrumbs with consistent styling
- **FooterModel**: Consistent footer with status information
- **HelpFooterModel**: Keyboard shortcut hints with context-aware display
- **ListItemModel**: Reusable list item component with selection states
- **TagSelector & CategorySelector**: Specialized input components for multi-select
- **AudienceRelevanceSelector**: Domain-specific selector component
- **NavigationMenu**: Menu component with focus management
- **ProgressIndicator & Spinner**: Status display components

### Existing Model Infrastructure (`BaseStandardModel`)
- **Context and metadata management**: Screen ID, Previous state, Custom data storage
- **Enhanced breadcrumb management**: Add, pop, peek, navigate operations with metadata support
- **Navigation history tracking**: Undo/redo capability with labeled state snapshots
- **Keyboard shortcut registration and handling**: Standardized shortcut management
- **Error handling**: Error storage, clearing, and retrieval
- **State management**: Generic state storage and retrieval
- **Validation support**: Standard validate/reset interface

### Existing Model Implementations
- **FormModel**: Advanced input handling with field-level validation, dynamic character count, multiple input types, focus indicators
- **ListModel**: Pagination, filtering, searching, sorting with consistent component integration
- **FactListModel, BurstListModel**: Specialized list variants with domain-specific features
- **DetailsModel**: Event details display with consistent header/footer structure
- **MetadataEditorModel, FactEditorModel**: Specialized form variants
- **ActionMenuModel, ConfirmationDialogModel**: Dialog and menu components

## 5. Functional Requirements

### 5.1 Standardized Layout Containers
The application must provide reusable lipgloss-based layout components that all models compose:

**Requirement 1: ScreenContainer Component**
- Wraps entire model view with consistent margins and max-width constraints
- Takes content slice and applies consistent outer spacing
- Supports optional padding modes (compact, standard, spacious)
- Uses `ColorBackground` as base and applies consistent width constraints
- Methods: `Render(content []string, padMode PaddingMode) string`

**Requirement 2: CardContainer Component**
- Renders boxed content with consistent border, padding, and background
- Supports header, body, and footer sections
- Applies `CardBase` style with optional custom background color
- Methods: `SetHeader()`, `SetBody()`, `SetFooter()`, `Render() string`

**Requirement 3: SectionContainer Component**
- Groups related content within a screen with consistent spacing
- Applies `HeaderSection` style to title
- Provides consistent spacing above and below content
- Methods: `SetTitle()`, `SetContent()`, `Render() string`

**Requirement 4: FormFieldContainer Component**
- Standardizes form field layout (label, input, error, hint)
- Applies consistent label styling via `InputLabel`
- Handles focus state via `InputFocused` style
- Handles error state via `InputError` style
- Methods: `SetLabel()`, `SetInput()`, `SetError()`, `SetHint()`, `Render() string`

**Requirement 5: ListContainer Component**
- Standardizes list display with consistent spacing and sizing
- Provides pagination info display
- Handles empty state rendering
- Methods: `SetItems()`, `SetPaginationInfo()`, `Render() string`

**Requirement 6: ModalContainer Component**
- Renders modal dialogs with consistent sizing, centering, and styling
- Supports title, message, buttons, and instructions sections
- Uses `ModalBase` or `ModalDestructive` styles appropriately
- Methods: `SetTitle()`, `SetMessage()`, `SetButtons()`, `SetInstructions()`, `Render() string`

### 5.2 Standardized Component Usage Patterns
All models must follow consistent patterns when using components:

**Requirement 7: Header/Footer Integration Standard**
- Every screen must use `HeaderModel` at the top with title and breadcrumbs
- Every screen must use `FooterModel` at the bottom with status/pagination info
- Header and footer must have consistent width matching screen content
- Header and footer state must be updated via dedicated setter methods in model

**Requirement 8: Help Text Standard**
- All models must render `HelpFooterModel` above the footer
- Help footer context key must match model type (e.g., "form", "list", "details")
- Help footer must display only keyboard shortcuts relevant to current state
- Help footer must use consistent styling from `styles.HelpText`

**Requirement 9: Error Display Standard**
- Field-level errors must use `FormFieldContainer` error display
- Model-level errors must use `ErrorBox` style
- Error messages must be brief, actionable, and human-readable
- Error dismissal must be consistent across all models

**Requirement 10: Focus Indicator Standard**
- Focused fields must use `InputFocused` or equivalent model-specific focus style
- Focus indicators must be clear and consistent (e.g., border color, bold text)
- Focus position must be rendered in consistent location (e.g., left margin indicator)
- All models with multiple focusable elements must support consistent focus navigation

### 5.3 Interaction Pattern Standardization
All models must implement consistent keyboard and mouse interaction patterns:

**Requirement 11: Navigation Standard**
- Form models: Tab/Shift+Tab for field navigation, Enter to submit, Esc to cancel
- List models: j/k for line navigation, g for top, G for bottom, Enter to select, q to quit
- Dialog models: Tab/Shift+Tab for button navigation, Enter to confirm, Esc to cancel
- Menu models: j/k for menu navigation, Enter to select, Esc to dismiss
- All models must support these patterns consistently

**Requirement 12: Keyboard Shortcut Registration**
- All models must register available shortcuts via `BaseStandardModel.RegisterShortcuts()`
- Shortcuts must be discoverable via `GetAvailableShortcuts()` method
- Help footer must dynamically display registered shortcuts
- Shortcut handling must be delegated to `HandleShortcut()` method where available

**Requirement 13: State Preservation Standard**
- All models must preserve breadcrumbs before navigation
- All models must track navigation history using `PushNavigationHistory()`
- All models must support back navigation via breadcrumb or history
- Focus position must be preserved when navigating back to a model

### 5.4 Style Application Standardization
All models must apply styles consistently through the standardized patterns:

**Requirement 14: Text Styling Consistency**
- Titles must use `HeaderMain` style
- Section titles must use `HeaderSection` style
- Primary text must use `ColorTextPrimary` (light gray)
- Secondary text must use `ColorTextSecondary` (medium gray)
- Muted text must use `ColorTextMuted` (dark gray)
- All text colors must come from `styles.Color*` constants

**Requirement 15: Border and Background Consistency**
- All containers must use `ColorBorder` for inactive borders
- All active/focused elements must use `ColorBorderActive` (teal)
- Error elements must use `ColorBorderError` (red)
- All backgrounds must use `ColorBackground` or `ColorBackgroundCard`
- No inline hex color codes; all colors from `styles.Color*` constants

**Requirement 16: Spacing Consistency**
- All padding must use consistent increments (e.g., 1, 2, 3 units)
- Margins between sections must follow consistent patterns
- All lipgloss styles must use `Padding()` and `Margin()` methods consistently
- Max width constraints must use `styles.MaxWidth()` utility

**Requirement 17: Status Indicator Consistency**
- Success states must use `ColorSuccess` with `SuccessBox` style
- Error states must use `ColorError` with `ErrorBox` style
- Warning states must use `ColorWarning` with `WarningBox` style
- Info states must use `ColorInfo` with `InfoBox` style

### 5.5 Model-Specific Standardization

**Requirement 18: Form Model Standardization**
- Must use `FormFieldContainer` for each input field
- Must render fields in consistent order: input → validation → character count → hint
- Must use `InputLabel` for field labels
- Must support error display via `FormFieldContainer`
- Must apply `InputFocused` when field has focus
- Must use consistent button styling for Submit/Cancel/Reset buttons

**Requirement 19: List Model Standardization**
- Must use `ScreenContainer` to wrap list content
- Must render `HeaderModel` with list title and breadcrumbs
- Must render list items using consistent `ListItemModel` component
- Must display pagination info in footer
- Must show empty state message with consistent styling when no items
- Must render help footer with list-specific shortcuts

**Requirement 20: Dialog Model Standardization**
- Must use `ModalContainer` component for rendering
- Must apply `ModalBase` or `ModalDestructive` style appropriately
- Must render title, message, buttons, and instructions in consistent order
- Must center modal within available space
- Must support consistent button arrangement (OK/Cancel in standard positions)

**Requirement 21: Menu Model Standardization**
- Must use consistent menu styling via lipgloss
- Must render menu items with consistent spacing and separators
- Must show focus indicators for selected menu item
- Must support consistent navigation shortcuts (j/k or arrow keys)

**Requirement 22: Data View Model Standardization**
- Must use `CardContainer` for grouping related data
- Must use `SectionContainer` for sectioning data
- Must apply consistent text styling via `HeaderSection`, `ColorTextPrimary`, etc.
- Must render metadata in consistent format and layout

## 6. Non-Goals (Out of Scope)

1. Changing the BubbleTea framework or fundamental architecture
2. Redesigning the color palette (use existing professional dark theme)
3. Adding new interaction paradigms (stay keyboard-first with mouse support)
4. Building web version or alternate UI rendering
5. Changing the data model or domain entities
6. Rewriting all tests in a single effort (tests updated alongside component adoption)
7. Creating animation/transition effects
8. Supporting responsive layouts below 80 character width

## 7. Design Considerations

- **Composition over inheritance**: Models compose components rather than extending complex base classes
- **Immutable component configuration**: Component states created fresh rather than mutated
- **Progressive adoption**: Existing models gradually migrated rather than all-or-nothing refactor
- **Lipgloss as single source of truth**: All styling goes through lipgloss styles in `styles.go`
- **Performance**: Container components avoid unnecessary string allocations and rendering passes
- **Accessibility**: Consistent focus indicators and keyboard navigation support screen reader compatibility
- **Future extensibility**: Container components support theming via style parameter overrides

## 8. Technical Considerations

### 8.1 Implementation Approach
- New container components created in `internal/cli/components/` package
- Each container is a simple struct with builder-pattern methods
- Containers use `lipgloss` exclusively for styling (no raw ANSI codes)
- Existing `styles.go` provides all style constants; no new style definitions in containers

### 8.2 Backward Compatibility
- Existing models continue to work during transition
- New container components used alongside existing components initially
- Gradual migration path: old models → adopt containers → adopt new patterns
- `BaseStandardModel` remains unchanged; new functionality layered on top

### 8.3 Performance
- Container components render to string efficiently (single pass)
- No caching or memoization needed initially
- Inline string building acceptable for CLI (not web rendering)
- Can optimize later if profiling shows bottlenecks

### 8.4 Testing Strategy
- Each container component has unit tests
- Integration tests verify layout output matches expected string format
- Model tests verify container usage produces correct View() output
- No snapshot testing; use string comparison for layout verification

## 9. Success Metrics

### Quantitative
1. **Code reusability**: All models implement at least 80% of layout via container components (measurable via code analysis)
2. **Duplication reduction**: Eliminate 70%+ of direct `lipgloss` style application in models (compare before/after LOC)
3. **Component reuse**: Every container component used in 3+ models within 6 months
4. **Style consistency**: 100% of colors from `styles.Color*` constants (verified via static analysis)
5. **Development velocity**: New models created in 40% less time when using containers vs. without

### Qualitative
1. **Visual consistency**: Users perceive unified appearance across all screens
2. **Developer experience**: Developers report easier model creation and maintenance
3. **Code clarity**: Code reviewers report easier comprehension of layout logic
4. **Consistency satisfaction**: Internal team agreement that models follow consistent patterns (survey)

## 10. Open Questions and Considerations

1. **Container Component Granularity**: Should `FormFieldContainer` be one component or separate Label/Input/Error components?
   - **Answer**: Single component with SetLabel/SetInput/SetError methods for ease of use

2. **Width Management**: How should max-width be enforced—in each container or screen-wide?
   - **Answer**: `ScreenContainer` handles outer width; individual containers inherit

3. **Spacing Presets**: Should containers support predefined spacing modes or allow custom values?
   - **Answer**: Support both via `PaddingMode` enum and optional `CustomPadding` parameter

4. **State Management**: Should containers be stateless rendering or stateful components?
   - **Answer**: Stateless pure rendering functions; state lives in models only

5. **Backward Compatibility Timeline**: How long should old-style models be supported?
   - **Answer**: Migrate all models as a requirement; no deprecation period
     needed

6. **Testing Approach**: How to test layout output without brittle string comparisons?
   - **Answer**: Test logical structure (sections, heights) not exact spacing; use visual inspection for final review

7. **Theme Support**: How to handle potential future dark/light theme switching?
   - **Answer**: Design containers to accept style parameters; `styles.go` becomes theme-agnostic

8. **Performance Baseline**: What's the acceptable rendering time per model?
   - **Answer**: <50ms per render acceptable for CLI; measure and optimize if needed

## 11. Implementation Phases

### Phase 1: Foundation (2-3 weeks)
**Create Core Container Components**
- Implement `ScreenContainer` with padding modes
- Implement `CardContainer` with header/body/footer sections
- Implement `SectionContainer` for content grouping
- Implement `FormFieldContainer` for form field layout
- Create 20+ unit tests for container rendering
- Update `styles.go` documentation to reference container usage

**Deliverables**:
- 5 new container component files in `internal/cli/components/`
- 100+ lines of unit tests
- Update `DESIGN.md` with container API documentation
- Existing styles.go unchanged

### Phase 2: Adoption—Data Display Models (2-3 weeks)
**Migrate Core Data Models to Containers**
- Refactor `ListModel` to use `ScreenContainer`, `ListContainer`, `CardContainer`
- Refactor `DetailsModel` to use containers
- Refactor `FactListModel` and `BurstListModel` to use `ListContainer`
- Refactor `FactsResultsModel` to use containers
- Update associated tests
- Verify visual consistency via manual testing

**Deliverables**:
- 5 model files refactored
- All tests passing
- Before/after LOC comparison showing duplication reduction
- Visual verification checklist completed

### Phase 3: Adoption—Form and Dialog Models (2-3 weeks)
**Migrate Form and Dialog Models to Containers**
- Refactor `FormModel` to use `ScreenContainer` and `FormFieldContainer`
- Refactor `MetadataEditorModel` to use containers
- Refactor `FactEditorModel` to use containers
- Refactor `ConfirmationDialogModel` to use `ModalContainer`
- Refactor other dialog/modal models to use containers
- Update associated tests

**Deliverables**:
- 5+ model files refactored
- All tests passing
- Form field consistency verified
- Dialog styling consistency verified

### Phase 4: Adoption—Menu and Utility Models (1-2 weeks)
**Migrate Remaining Model Types**
- Refactor `ActionMenuModel` and `NavigationMenuModel`
- Refactor menu-style models to use consistent styling
- Refactor utility models (progress, spinner, etc.) as appropriate
- Update all remaining models to follow patterns

**Deliverables**:
- All models using container components
- All tests passing
- Final consistency audit completed

### Phase 5: Pattern Standardization—Interaction (2-3 weeks)
**Standardize Keyboard and Mouse Interactions**
- Audit all model keyboard handling
- Ensure navigation shortcuts follow Requirement 11 (Tab, j/k, Enter, Esc)
- Audit all models for breadcrumb/history usage
- Ensure all models register shortcuts via BaseStandardModel
- Update help footers to reflect standardized shortcuts
- Add missing shortcut handlers

**Deliverables**:
- All models follow standardized interaction patterns
- Help footer reflects accurate shortcuts per model type
- Navigation consistency document
- Integration tests for common navigation workflows

### Phase 6: Documentation and Validation (1-2 weeks)
**Document Patterns and Complete Validation**
- Write `COMPONENT_USAGE_GUIDE.md` with examples for each container
- Write `MODEL_DEVELOPMENT_GUIDE.md` for creating new models
- Create style consistency checklist and audit results
- Verify color consistency via static analysis
- Performance testing and optimization if needed
- Final visual walkthrough and sign-off

**Deliverables**:
- Complete developer documentation
- Style audit report (100% compliance verification)
- Performance benchmark results
- Team sign-off checklist

## 12. Success Criteria Validation

### At Completion
1. ✅ All data display models use container components (ListModel, DetailsModel, FactListModel, etc.)
2. ✅ All form models use FormFieldContainer (FormModel, MetadataEditorModel, FactEditorModel)
3. ✅ All dialog/modal models use ModalContainer (ConfirmationDialogModel, etc.)
4. ✅ All keyboard shortcuts follow standardized patterns (Tab, j/k, Enter, Esc)
5. ✅ All text colors use `styles.Color*` constants (static analysis verification)
6. ✅ All borders use `ColorBorder` or `ColorBorderActive` constants
7. ✅ Help footers accurately reflect model's keyboard shortcuts
8. ✅ Breadcrumb and history management consistent across all models
9. ✅ 768+ tests passing (same as baseline)
10. ✅ Code coverage maintained at 76%+ (same as baseline)
11. ✅ Visual consistency audit completed and approved

## Document Information

- **Version**: 1.0
- **Created**: 2025-12-31
- **Status**: Ready for Implementation
- **Target Start**: 2026-Q1
- **Estimated Duration**: 10-14 weeks (6 phases)
- **Priority**: High (Foundational for Phase 2 of UX Enhancement)
- **Dependencies**: None (builds on existing infrastructure)
- **Related Documents**:
  - `/docs/features/05-ux-enhancement.md` (Phase 1: Navigation - Completed)
  - `/internal/cli/styles/styles.go` (Styling baseline)
  - `/internal/cli/models/standard_model.go` (BaseStandardModel reference)
  - `/internal/cli/components/` (Existing component library)

