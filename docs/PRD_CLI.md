# Product Requirements Document: KaRiya Career Entry CLI

## 1. Introduction/Overview

The KaRiya Career Entry CLI is a command-line interface tool designed to provide
an intuitive, visually polished entry point for capturing career events into the
KaRiya system. The CLI leverages BubbleTea and Lipgloss to deliver a
professional, interactive user experience that guides users through the event
capture process with minimal friction while maintaining data quality and
validation.

This tool serves as the primary interface for users to begin their career
journaling journey, making the initial capture experience smooth, encouraging,
and error-resistant.

## 2. Goals

1. Provide an intuitive, interactive CLI interface for capturing career events
2. Support multiple event capture modes (Timeline Journaling, CV Backfill, Manual Entry)
3. Deliver a polished, professional user experience using BubbleTea and Lipgloss
4. Ensure robust validation with helpful, corrective feedback
5. Enable users to view and manage their captured events
6. Make career event capture the primary entry point to the KaRiya system

## 3. User Stories

1. As a senior engineer, I want to quickly launch a CLI tool and capture a
   career event without navigating complex menus, so that I can log achievements
   while they're fresh in my mind.

2. As a user backfilling my CV, I want to import career events from different
   time periods with flexible date selection, so that I can build a
   comprehensive career history.

3. As a new user, I want clear guidance and visual feedback when entering event
   details, so that I understand what information is required and how to correct
   mistakes.

4. As a user exploring the system, I want to view my recently captured events
   directly from the CLI, so that I can verify my entries and see my growing
   career journal.

5. As a user managing my career data, I want to search and filter my events by
   various criteria, so that I can find and review specific career milestones.

## 4. Functional Requirements

### 4.1 Event Capture
1. The CLI must support three capture modes:
   - **CV Backfill Mode**: For importing events from existing CVs (allows older dates, no time constraints)
   - **Timeline Journaling Mode**: For real-time event logging (within 30 days of today)
   - **Manual Entry Mode**: For manually adding individual events (no date constraints)

2. The system must accept the following inputs for each event:
   - **Text** (mandatory): Event description (1-2000 characters)
   - **Date** (optional): Event date, defaults to today if not provided
   - **Company** (optional): Company or organization name
   - **Project** (optional): Project name associated with the event
   - **Tags** (optional): Multiple tags from the AllowedTags set

3. The CLI must validate all inputs with the following rules:
   - Text: Non-empty, maximum 2000 characters
   - Date: Not in the future (validate against current date)
   - Tags: Must be from the AllowedTags set (project, achievement, leadership, technical, consulting, research, product, mentoring)
   - No duplicate tags allowed
   - Maximum 8 tags per event

4. When validation fails, the CLI must:
   - Display a clear error message explaining what went wrong
   - Suggest how to fix the error
   - Allow the user to correct and re-submit the specific field without restarting

### 4.2 User Interface & Interaction
1. The CLI must use BubbleTea and Lipgloss to create a professional, polished interface with:
   - Rich colors and visual styling
   - Clear visual hierarchy
   - Responsive layout
   - Smooth animations and transitions
   - Professional/corporate aesthetic (clean, minimal colors)

2. The interface must guide users through a multi-step form with:
   - Clear field labels and descriptions
   - Input prompts with contextual help
   - Progress indication (e.g., "Step 1 of 5")
   - Ability to navigate between fields (back/forward)

3. Input collection must support:
   - Text input fields with character count feedback
   - Date picker with calendar or text input
   - Multi-select checklist for tags with autocomplete suggestions
   - Dropdown selection for capture mode
   - Inline validation feedback as user types

### 4.3 Event Confirmation & Display
1. After successful event capture, the CLI must:
   - Display a formatted summary of the captured event with all details
   - Show the assigned event ID
   - Present the event in a visually appealing card/box format using Lipgloss

2. The display must include:
   - Event text
   - Date captured
   - Company (if provided)
   - Project (if provided)
   - Tags (if provided)
   - Timestamp of capture
   - Event ID for reference

3. After displaying the event, the CLI must offer options to:
   - Capture another event (restart the form)
   - Exit the application
   - View recent events
   - Search/filter events

### 4.4 Event Listing & Browsing
1. The CLI must support viewing recent captured events with:
   - List of last N events (default: 10, configurable)
   - Event summary display (text preview, date, company, tags)
   - Pagination support for large event lists
   - Ability to view full details of any event

2. The CLI must support filtering events by:
   - Date range (start date, end date)
   - Tags (single or multiple)
   - Company name
   - Text search (keyword matching)

3. The CLI must support sorting by:
   - Date (ascending/descending)
   - Creation date (ascending/descending)
   - Text (alphabetical)

### 4.5 Help & Documentation
1. The CLI must include:
   - Interactive tutorial on first run explaining the capture process
   - Built-in --help flag with command documentation
   - Inline tips and hints for each input field
   - Context-sensitive help accessible during event capture

2. Help content must cover:
   - What each field represents
   - How to use different capture modes
   - Tag selection and best practices
   - How to view and filter events
   - Keyboard shortcuts and navigation tips

### 4.6 Error Handling & Validation Feedback
1. All validation errors must:
   - Display in a prominent, visually distinct way (using Lipgloss styling)
   - Explain what went wrong in plain language
   - Suggest how to correct the error
   - Highlight the problematic field

2. The CLI must handle edge cases:
   - Empty/whitespace-only text input
   - Dates in the future
   - Invalid date formats
   - Duplicate tags
   - Tags not in the AllowedTags set
   - Text exceeding 2000 characters
   - Can only have company or project, not both empty or both filled
   - Network/database errors (graceful error messages)

## 5. Non-Goals (Out of Scope)

1. Event editing or deletion from the CLI (future enhancement)
2. CV generation from the CLI (separate component)
3. Burst generation or fact extraction from the CLI
4. User authentication (assumes local/trusted environment for MVP)
5. Real-time sync with remote servers (MVP uses local storage)
6. Advanced analytics or reporting
7. Integration with external systems (LinkedIn, etc.) in this version
8. Mobile or web UI (CLI-only for this version)

## 6. Design Considerations

### 6.1 Visual Design
1. Use a professional, corporate color scheme:
   - Primary: Dark blue or dark gray backgrounds
   - Accent: Muted teal, green, or purple for highlights
   - Text: Light gray or white for readability
   - Errors: Red or amber for warnings/errors

2. Layout principles:
   - Center-aligned main content area
   - Clear visual separation between sections
   - Consistent spacing and padding
   - Use of borders, boxes, and dividers for structure

3. Typography:
   - Clear, readable fonts
   - Distinct sizes for headings, labels, and body text
   - Consistent styling across the interface

### 6.2 Interaction Design
1. Keyboard-first navigation:
   - Tab/Shift+Tab for field navigation
   - Arrow keys for selections/dropdowns
   - Enter to confirm, Escape to cancel
   - Clear keyboard shortcut hints

2. Feedback mechanisms:
   - Visual feedback on field focus
   - Spinner/loader during database operations
   - Success messages with checkmarks
   - Error messages with clear icons

3. User guidance:
   - Progress bar showing capture form completion
   - Helpful tooltips on hover (where applicable)
   - Clear call-to-action buttons
   - Confirmation prompts for destructive actions

## 7. Technical Considerations

### 7.1 Architecture
1. The CLI must integrate with the existing KaRiya service layer:
   - Use `internal/service/career/Service` for event capture
   - Use `internal/repository/career/Repository` for persistence
   - Leverage `internal/domain/career/CareerEvent` domain model
   - Integrate with `internal/service/career/classification/Classifier` (optional display)

2. Technology stack:
   - **Framework**: BubbleTea v0.x (latest stable)
   - **Styling**: Lipgloss v0.x (latest stable)
   - **Logging**: Existing KaRiya logger
   - **Data Storage**: SQLite via existing repository implementation

### 7.2 Performance Requirements
1. Event capture form must load in < 500ms
2. Form submission must complete in < 2 seconds
3. Event list display must load in < 1 second
4. Search/filter operations must complete in < 2 seconds
5. Support up to 10,000 events without performance degradation

### 7.3 Error Handling & Resilience
1. Graceful handling of database connection failures
2. Clear messaging when service is unavailable
3. Offline mode support (queue events for later sync if applicable)
4. Recovery options for interrupted captures

### 7.4 Dependencies
1. Must use existing domain models and services (no duplication)
2. Must follow KaRiya's DDD architecture patterns
3. Must integrate with existing logging infrastructure
4. Must use existing repository implementations (SQLite/Memory)

## 8. Success Metrics

1. **Primary Metric**: Number of career events captured via CLI
2. **Secondary Metrics**:
   - Average time to capture an event (target: < 2 minutes)
   - User completion rate (percentage of users who complete first capture)
   - Error rate (percentage of invalid submissions)
   - User satisfaction with UI/UX (if survey conducted)

3. **Quality Metrics**:
   - Data validation accuracy (100% of stored events pass validation)
   - CLI uptime/reliability (99%+)
   - Error message clarity (qualitative assessment)

## 9. Acceptance Criteria

1. ✓ Users can launch the CLI and be presented with a capture form
2. ✓ Users can select from three capture modes (CV Backfill, Timeline Journaling, Manual Entry)
3. ✓ Users can input event text, date, company, project, and tags
4. ✓ All validation rules are enforced with clear error messages
5. ✓ Users can correct invalid inputs and resubmit without losing data
6. ✓ After successful capture, event details are displayed in a formatted card
7. ✓ Users can capture another event or view recent events from the success screen
8. ✓ Users can list recent events with pagination
9. ✓ Users can filter events by date range, tags, company, and search text
10. ✓ Users can sort events by date, creation date, or text
11. ✓ Help documentation is accessible and comprehensive
12. ✓ First-run tutorial guides new users through the capture process
13. ✓ All UI elements use BubbleTea and Lipgloss for consistent styling
14. ✓ Error handling is graceful with helpful recovery suggestions
15. ✓ CLI integrates seamlessly with existing KaRiya service layer

## 10. Open Questions

1. Should the CLI support importing events from a CSV/JSON file in bulk?
2. Should captured events be automatically classified and displayed with competency categories?
3. What should be the default limit for listing recent events (10, 20, 50)?
4. Should the CLI support exporting events to JSON/YAML for backup?
5. Should there be a "quick capture" mode with minimal fields for rapid entry?
6. How should the CLI handle very long event text (> 500 chars) in list displays?

## 11. Phased Rollout Approach

### Phase 1: MVP - Core Event Capture (Priority 1)
- Basic event capture form with Text field (mandatory)
- Single capture mode selection (user chooses from 3 modes)
- Date input (optional, defaults to today)
- Company and Project fields (optional)
- Tag selection from AllowedTags
- Input validation with error messages
- Success screen with event display
- Basic help/documentation
- Capture another event or exit options

### Phase 2: Event Management (Priority 2)
- View recent events (list with pagination)
- View event details
- Filter events by tags, company, date range
- Sort events by date/creation date/text
- Search events by keyword
- First-run interactive tutorial
- Improved inline help and tips

### Phase 3: Polish & Enhancement (Priority 3)
- Bulk import from CSV/JSON
- Export events to JSON/YAML
- Quick capture mode (minimal fields)
- Event editing (if architecture supports)
- Event deletion with confirmation
- Advanced filtering and saved searches
- Performance optimizations

## 12. Implementation Notes

### 12.1 BubbleTea Integration
- Use BubbleTea's model-view-update (MVU) pattern
- Implement separate models for each screen (capture form, event list, details, etc.)
- Use BubbleTea's key handling for keyboard navigation
- Leverage BubbleTea's program state management

### 12.2 Lipgloss Usage
- Create reusable style definitions for consistent theming
- Use Lipgloss for all text styling, borders, and layout
- Implement color scheme consistently across screens
- Create component-level styles for buttons, cards, inputs

### 12.3 Service Integration
- Inject existing `career.Service` into CLI commands
- Use `career.Service.CaptureEvent()` for event creation
- Use repository methods for listing and filtering
- Leverage domain validation from `CareerEvent`

### 12.4 Testing Strategy
- Unit tests for CLI input validation logic
- Integration tests with mock service layer
- Manual testing for UI/UX polish
- Edge case testing for error scenarios

---

**Document Version**: 1.0
**Status**: Ready for Implementation
**Target Audience**: Go developers familiar with BubbleTea and KaRiya architecture

