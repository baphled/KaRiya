# Task List: KaRiya Career Entry CLI Implementation

**Based on PRD**: `tasks/prd-career-entry-cli.md`
**Status**: Phase 1 In Progress (Infrastructure Complete, Form Implementation Ongoing)
**Target Audience**: Go developers familiar with KaRiya architecture and BubbleTea

---

## Relevant Files

### CLI Package Structure
- `cmd/cli/main.go` - CLI entry point and command initialization
- `cmd/cli/main_test.go` - Tests for CLI initialization and command setup

### BubbleTea Models & Screens
- `internal/cli/models/form.go` - Event capture form model and logic
- `internal/cli/models/form_test.go` - Unit tests for form model
- `internal/cli/models/list.go` - Event list model and pagination logic
- `internal/cli/models/list_test.go` - Unit tests for list model
- `internal/cli/models/details.go` - Event details view model
- `internal/cli/models/details_test.go` - Unit tests for details model
- `internal/cli/models/help.go` - Help/tutorial screen model
- `internal/cli/models/help_test.go` - Unit tests for help model
- `internal/cli/models/success.go` - Success/confirmation screen model
- `internal/cli/models/success_test.go` - Unit tests for success model

### CLI Components & Utilities
- `internal/cli/components/inputs.go` - Reusable input field components
- `internal/cli/components/inputs_test.go` - Tests for input components
- `internal/cli/components/date_picker.go` - Date input/picker component
- `internal/cli/components/date_picker_test.go` - Tests for date picker
- `internal/cli/components/tag_selector.go` - Multi-select tag component
- `internal/cli/components/tag_selector_test.go` - Tests for tag selector
- `internal/cli/styles/styles.go` - Lipgloss style definitions and theme
- `internal/cli/styles/styles_test.go` - Tests for style application
- `internal/cli/validation/validator.go` - CLI-specific validation logic
- `internal/cli/validation/validator_test.go` - Tests for validation

### Application State & Navigation
- `internal/cli/app/app.go` - Main CLI application state and navigation
- `internal/cli/app/app_test.go` - Tests for app state management
- `internal/cli/app/suite_test.go` - Test suite setup for app tests

### Integration & Service
- `internal/cli/service/event_service.go` - Service wrapper for event operations
- `internal/cli/service/event_service_test.go` - Tests for event service wrapper

### Notes

- Unit tests should be placed alongside the code files they test (e.g., `form.go` and `form_test.go` in the same directory)
- Use `make test` or `ginkgo -v ./...` to run tests
- All CLI code should follow KaRiya's existing patterns and use the existing service/repository layers

---

## Tasks

### Phase 1: MVP - Core Event Capture

- [x] 1.0 Set Up CLI Project Structure & Dependencies
  - [x] 1.1 Create CLI package directories (`internal/cli/{models,components,styles,validation,app,service}`)
  - [x] 1.2 Add BubbleTea and Lipgloss to `go.mod` dependencies
  - [x] 1.3 Create `cmd/cli/main.go` as CLI entry point
  - [x] 1.4 Implement CLI command initialization with flag parsing (--help, --version)
  - [x] 1.5 Set up integration with existing `career.Service` and repository layer
  - [x] 1.6 Create base application state struct with navigation between screens
  - [x] 1.7 Write integration tests for CLI initialization

- [x] 2.0 Implement Lipgloss Theme & Styling System
  - [x] 2.1 Define color scheme (dark blue/gray background, muted teal/green/purple accents)
  - [x] 2.2 Create reusable Lipgloss style definitions for buttons, inputs, cards, headers
  - [x] 2.3 Implement error message styling (red/amber for warnings)
  - [x] 2.4 Create component-level styles for consistency across screens
  - [x] 2.5 Implement responsive layout helper functions
  - [x] 2.6 Write tests for style application and theming

- [x] 3.0 Implement Event Capture Form (Core MVP) - **MOSTLY COMPLETE** (needs tag integration)
  - [x] 3.1 Create BubbleTea form model with multi-step navigation
  - [x] 3.2 Implement text input field with:
    - [x] 3.2a Character counting and limit enforcement (max 2000 chars)
    - [x] 3.2b Real-time character count display
    - [x] 3.2c Input validation feedback
  - [x] 3.3 Implement date input with:
    - [x] 3.3a Text input for date (format: YYYY-MM-DD or relative like "today", "1 week ago")
    - [x] 3.3b Validation against future dates
    - [x] 3.3c Default to today if not provided
    - [x] 3.3d Mode-specific date constraints (Timeline: 30 days, CVBackfill: any date, ManualEntry: any date)
  - [x] 3.4 Implement optional company field with validation
  - [x] 3.5 Implement optional project field with validation
  - [x] 3.6 Implement capture mode selector (dropdown with 3 modes: CV Backfill, Timeline Journaling, Manual Entry)
  - [x] 3.7 Write comprehensive unit tests for form model and validation (86.9% coverage)
  - [x] 3.8 Write integration tests for form submission with service layer

- [ ] 4.0 Implement Tag Selection Component
  - [ ] 4.1 Create multi-select tag picker component with AllowedTags set
  - [ ] 4.2 Implement tag autocomplete/filtering as user types
  - [ ] 4.3 Implement duplicate tag prevention
  - [ ] 4.4 Implement max 8 tags per event enforcement
  - [ ] 4.5 Display selected tags with visual indication
  - [ ] 4.6 Write unit tests for tag component
  - [ ] 4.7 Write integration tests with form model

- [ ] 5.0 Implement Event Display & Success Screen
  - [ ] 5.1 Create success screen model with event summary display
  - [ ] 5.2 Implement formatted event card using Lipgloss with:
    - [ ] 5.2a Event text display
    - [ ] 5.2b Date display with formatting
    - [ ] 5.2c Company display (if provided)
    - [ ] 5.2d Project display (if provided)
    - [ ] 5.2e Tags display with styling
    - [ ] 5.2f Event ID display
    - [ ] 5.2g Timestamp display
  - [ ] 5.3 Implement post-capture options:
    - [ ] 5.3a "Capture Another Event" button (returns to form)
    - [ ] 5.3b "View Recent Events" button (navigates to list)
    - [ ] 5.3c "Exit" button (graceful shutdown)
  - [ ] 5.4 Write unit tests for success screen
  - [ ] 5.5 Write integration tests for event capture workflow

- [ ] 6.0 Implement Input Validation & Error Handling
  - [ ] 6.1 Create validation wrapper for all form inputs
  - [ ] 6.2 Implement error message display with:
    - [ ] 6.2a Clear explanation of what went wrong
    - [ ] 6.2b Suggestion for how to fix
    - [ ] 6.2c Prominent visual styling
  - [ ] 6.3 Implement field-level validation with inline feedback
  - [ ] 6.4 Handle edge cases:
    - [ ] 6.4a Empty/whitespace-only text
    - [ ] 6.4b Future dates
    - [ ] 6.4c Invalid date formats
    - [ ] 6.4d Duplicate tags
    - [ ] 6.4e Invalid tags not in AllowedTags
    - [ ] 6.4f Text exceeding 2000 characters
  - [ ] 6.5 Implement graceful handling of database errors
  - [ ] 6.6 Write comprehensive tests for all validation scenarios

- [ ] 7.0 Implement Keyboard Navigation & Shortcuts
  - [ ] 7.1 Implement Tab/Shift+Tab for field navigation
  - [ ] 7.2 Implement Arrow keys for selections and dropdowns
  - [ ] 7.3 Implement Enter to confirm, Escape to cancel
  - [ ] 7.4 Display keyboard shortcut hints on screen
  - [ ] 7.5 Add visual feedback for focused fields
  - [ ] 7.6 Write tests for keyboard interaction

- [ ] 8.0 Integration Testing & MVP Completion
  - [ ] 8.1 Write end-to-end integration tests for complete capture workflow
  - [ ] 8.2 Test all three capture modes (CV Backfill, Timeline Journaling, Manual Entry)
  - [ ] 8.3 Test error recovery and field correction
  - [ ] 8.4 Test database persistence with SQLite repository
  - [ ] 8.5 Verify all MVP acceptance criteria are met
  - [ ] 8.6 Run full test suite with `make test`
  - [ ] 8.7 Achieve minimum 80% code coverage for CLI package

---

### Phase 2: Event Management (Priority 2)

- [ ] 9.0 Implement Event Listing & Pagination
  - [ ] 9.1 Create event list model with BubbleTea
  - [ ] 9.2 Implement list display showing:
    - [ ] 9.2a Event text preview (truncated if > 100 chars)
    - [ ] 9.2b Event date
    - [ ] 9.2c Company name (if available)
    - [ ] 9.2d Tags (if available)
  - [ ] 9.3 Implement pagination with:
    - [ ] 9.3a Page size configuration (default: 10, max: 50)
    - [ ] 9.3b Previous/Next page navigation
    - [ ] 9.3c Current page indicator
    - [ ] 9.3d Jump to page functionality
  - [ ] 9.4 Implement keyboard navigation through list
  - [ ] 9.5 Write unit tests for list model
  - [ ] 9.6 Write integration tests with repository layer

- [ ] 10.0 Implement Event Filtering System
  - [ ] 10.1 Create filter UI screen with options for:
    - [ ] 10.1a Date range selection (start/end date)
    - [ ] 10.1b Tag multi-select filtering
    - [ ] 10.1c Company name filtering
    - [ ] 10.1d Clear/reset filters option
  - [ ] 10.2 Implement date range validation
  - [ ] 10.3 Implement filter application to repository queries
  - [ ] 10.4 Display active filters on list screen
  - [ ] 10.5 Write unit tests for filter logic
  - [ ] 10.6 Write integration tests with list and repository

- [ ] 11.0 Implement Event Search Functionality
  - [ ] 11.1 Create search input field on list screen
  - [ ] 11.2 Implement keyword search across event text
  - [ ] 11.3 Implement real-time search with debouncing
  - [ ] 11.4 Display search results with highlighting
  - [ ] 11.5 Implement clear search option
  - [ ] 11.6 Write unit tests for search logic
  - [ ] 11.7 Write integration tests with repository

- [ ] 12.0 Implement Event Sorting
  - [ ] 12.1 Create sort options menu with choices:
    - [ ] 12.1a Date (ascending/descending)
    - [ ] 12.1b Creation date (ascending/descending)
    - [ ] 12.1c Text (alphabetical A-Z/Z-A)
  - [ ] 12.2 Implement sort application to repository queries
  - [ ] 12.3 Display current sort order on list screen
  - [ ] 12.4 Make sorting interactive (allow change without re-querying)
  - [ ] 12.5 Write unit tests for sort logic
  - [ ] 12.6 Write integration tests with repository

- [ ] 13.0 Implement Event Details View
  - [ ] 13.1 Create details screen model
  - [ ] 13.2 Display full event information:
    - [ ] 13.2a Full event text
    - [ ] 13.2b Date with formatting
    - [ ] 13.2c Company (if provided)
    - [ ] 13.2d Project (if provided)
    - [ ] 13.2e All tags with styling
    - [ ] 13.2f Event ID
    - [ ] 13.2g Created/Updated timestamps
  - [ ] 13.3 Implement navigation from list to details
  - [ ] 13.4 Implement back button to return to list
  - [ ] 13.5 Write unit tests for details model
  - [ ] 13.6 Write integration tests with list and repository

- [ ] 14.0 Implement First-Run Interactive Tutorial
  - [ ] 14.1 Create tutorial screen with step-by-step guidance
  - [ ] 14.2 Implement tutorial steps covering:
    - [ ] 14.2a What is KaRiya and career journaling
    - [ ] 14.2b The three capture modes
    - [ ] 14.2c How to fill each field
    - [ ] 14.2d Tag selection and best practices
    - [ ] 14.2e Viewing and filtering events
  - [ ] 14.3 Implement "Skip Tutorial" option
  - [ ] 14.4 Implement "View Tutorial Again" option from help
  - [ ] 14.5 Store tutorial completion state (skip on future runs)
  - [ ] 14.6 Write tests for tutorial flow

---

### Phase 3: Help System & Polish

- [ ] 15.0 Implement Comprehensive Help System
  - [ ] 15.1 Create help screen with sections for:
    - [ ] 15.1a Overview of KaRiya
    - [ ] 15.1b Event capture modes explained
    - [ ] 15.1c Field descriptions and requirements
    - [ ] 15.1d Tag selection and best practices
    - [ ] 15.1e How to view and filter events
    - [ ] 15.1f Keyboard shortcuts reference
  - [ ] 15.2 Implement context-sensitive help for each field
  - [ ] 15.3 Implement inline tips and hints during capture
  - [ ] 15.4 Implement searchable help content
  - [ ] 15.5 Write tests for help content

- [ ] 16.0 Implement CLI Flags & Configuration
  - [ ] 16.1 Implement `--help` flag with command documentation
  - [ ] 16.2 Implement `--version` flag showing CLI version
  - [ ] 16.3 Implement `--db` flag for custom database path
  - [ ] 16.4 Implement `--mode` flag to start in specific capture mode
  - [ ] 16.5 Implement `--list` flag to show recent events on startup
  - [ ] 16.6 Write tests for flag parsing and handling

- [ ] 17.0 Implement Error Recovery & Edge Cases
  - [ ] 17.1 Handle database connection failures gracefully
    - [ ] 17.1a Display user-friendly error message
    - [ ] 17.1b Suggest troubleshooting steps
  - [ ] 17.2 Handle service layer errors
    - [ ] 17.2a Timeout errors with retry option
    - [ ] 17.2b Validation errors from service
  - [ ] 17.3 Handle very long event text in list displays (truncation)
  - [ ] 17.4 Handle large datasets (10,000+ events) without performance degradation
  - [ ] 17.5 Implement graceful shutdown on interrupt (Ctrl+C)
  - [ ] 17.6 Write tests for error scenarios

- [ ] 18.0 UI/UX Polish & Refinement
  - [ ] 18.1 Implement visual feedback mechanisms:
    - [ ] 18.1a Spinner/loader during database operations
    - [ ] 18.1b Success checkmarks for completed actions
    - [ ] 18.1c Progress indicator for multi-step form
  - [ ] 18.2 Implement smooth animations and transitions between screens
  - [ ] 18.3 Refine color scheme for professional appearance
  - [ ] 18.4 Optimize layout for various terminal sizes
  - [ ] 18.5 Implement consistent spacing and padding
  - [ ] 18.6 Test on different terminal emulators
  - [ ] 18.7 Gather feedback and iterate on UX

- [ ] 19.0 Performance Optimization
  - [ ] 19.1 Profile CLI startup time (target: < 500ms)
  - [ ] 19.2 Optimize form submission (target: < 2 seconds)
  - [ ] 19.3 Optimize event list loading (target: < 1 second)
  - [ ] 19.4 Optimize search/filter operations (target: < 2 seconds)
  - [ ] 19.5 Implement caching for frequently accessed data
  - [ ] 19.6 Add database indexing if needed
  - [ ] 19.7 Write benchmark tests for performance-critical code

- [ ] 20.0 Documentation & Testing Completion
  - [ ] 20.1 Write comprehensive README for CLI usage
  - [ ] 20.2 Create examples for each capture mode
  - [ ] 20.3 Create troubleshooting guide
  - [ ] 20.4 Document all keyboard shortcuts
  - [ ] 20.5 Document configuration options
  - [ ] 20.6 Run full test suite and achieve 80%+ coverage
  - [ ] 20.7 Run race detector tests (`go test -race ./...`)
  - [ ] 20.8 Verify all acceptance criteria are met
  - [ ] 20.9 Create CHANGELOG entry for CLI feature

---

## Implementation Notes

### Architecture Integration

1. **Service Layer Integration**: The CLI will inject the existing `career.Service` into all models that need it for event operations.

2. **Repository Layer Usage**: Models will use repository methods through the service layer for:
   - Event persistence (`Create`)
   - Event retrieval (`GetByID`)
   - Event listing with filters (`List`, `Count`)
   - No direct database access from CLI

3. **Domain Model Usage**: All events will use the existing `career.CareerEvent` domain model with its built-in validation.

4. **Logging Integration**: All CLI operations will use the existing `logger.Logger` for structured logging.

### BubbleTea Patterns

1. **Model-View-Update (MVU)**: Each screen (form, list, details, help) will be a separate BubbleTea model.

2. **State Management**: Main app model will manage navigation and screen transitions.

3. **Message Types**: Custom message types will be used for communication between models and the main app.

### Testing Strategy

1. **Unit Tests**: Test individual models, components, and validation logic in isolation.

2. **Integration Tests**: Test complete workflows (capture → display → list) with mock service layer.

3. **Manual Testing**: UI/UX polish and terminal compatibility testing.

### Phased Approach

- **Phase 1 (MVP)**: Core event capture form with success screen (tasks 1-8) - **IN PROGRESS** (tasks 1-2 complete, task 3 mostly complete, tasks 4-8 pending)
- **Phase 2**: Event management features - listing, filtering, search, sorting (tasks 9-14) - **NOT STARTED**
- **Phase 3**: Help system, polish, optimization (tasks 15-20) - **NOT STARTED**

---

**Document Version**: 1.1
**Created**: 2025-12-23
**Last Updated**: 2025-12-23
**Status**: Phase 1 Partially Complete (Tasks 1.0-2.0 ✓, 3.0 ⚠️ needs tags, 4.0-8.0 pending)
**Total Tasks**: 20 parent tasks, 80+ sub-tasks

