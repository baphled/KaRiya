# Legacy Screen Special Logic, Edge Cases & Styling

**Date**: 2026-01-03
**Status**: Phase 1.1.3-1.1.5 - Documentation
**Purpose**: Document complex logic, edge cases, and UI patterns for migration to intents

---

## CaptureEvent Intent

### Special Logic
- **Form Validation**: Multi-field validation with field-level error messages
  - Title: Required, max 255 chars
  - Description: Optional, max 2000 chars
  - Date: Required, must be valid date, not future
  - Duration: Required, must be positive integer
  - Skills: Optional list, each max 100 chars
  - Tags: Optional list, auto-complete from existing tags

- **Tab Navigation**: Users can Tab/Shift+Tab between form fields
- **Auto-save**: Drafts are auto-saved to prevent data loss
- **Field Dependencies**: Some fields enable/disable based on other fields
  - If "Event Type" = "Project", show "Project Details" section
  - If "Outcome" selected, show "Metrics" fields

### Edge Cases
- Empty form submission (validation prevents)
- Rapid successive saves (debouncing)
- Editing existing event (pre-populate form)
- Network errors during save (retry mechanism)
- Concurrent edits (last-write-wins)
- Extremely long input (graceful truncation)
- Special characters in input (escaping)
- Copy-paste with formatting (strip formatting)

### Custom Styling
- **Form Container**: Dark background with subtle border
- **Field Labels**: Bold, left-aligned, with asterisk for required fields
- **Input Fields**:
  - White text on dark background
  - Focus state: Colored border (primary color)
  - Error state: Red border with error message below
  - Success state: Green checkmark on right side
- **Character Counter**: Bottom-right, gray text, "current/max"
- **Submit Button**:
  - Primary color background
  - Disabled state: Gray, no cursor
  - Hover state: Lighter shade
- **Cancel Button**: Secondary style, outlined

### Components Used
- FormModel (custom multi-field form)
- ConfirmationDialog (for save confirmation)
- SuccessModel (for success feedback)

### Data Transformations
- Form input → CareerEvent domain model
- Date string → time.Time (validation: not future, valid format)
- Duration string → int (validation: positive, reasonable range)
- Tags string → []string (split by comma, trim whitespace)
- Skills string → []Skill (validate against skill database)

---

## BrowseTimeline Intent

### Special Logic
- **List Pagination**: Large timelines split into pages of 20 items
- **Scroll Position Memory**: Remember scroll position when navigating back
- **Search/Filter State**: Persist filter state across navigation
  - Filter by date range
  - Filter by skill tags
  - Filter by event type
  - Search by title/description (full-text)
- **Selection Highlight**: Current item highlighted, survives navigation back
- **Action Menu Positioning**: Menu appears near selected item, doesn't go off-screen
- **Breadcrumb Trail**: Home → Timeline → Event Details → Action Menu

### Edge Cases
- Empty timeline (no events)
- Single event (no pagination)
- Very large timeline (100+ events, performance optimization)
- Deleted event during viewing (refresh, show error)
- Concurrent edits (show stale data warning)
- Filter returns no results (show empty state message)
- Search with special characters (escape for database query)
- Very long event titles (truncate with ellipsis)
- Very old events (date formatting edge cases)

### Custom Styling
- **List Items**:
  - Height: 4 lines (title, date, description, tags)
  - Selection: Inverted colors (light background, dark text)
  - Hover: Slightly lighter background
  - Separator: Thin line between items
- **Detail Panel**:
  - Two-column layout (left: content, right: metadata)
  - Rich text rendering for description
  - Tag badges with colors
  - Date formatted as "Monday, January 3, 2026"
- **Action Menu**:
  - Dropdown style, appears below selected item
  - Options: View, Edit, Delete, Export, Tag, Archive
  - Icons next to each option
  - Highlighted option on hover
- **Search Bar**:
  - Top of list, dark background
  - Placeholder text: "Search events..."
  - Clear button (X) on right when text entered
- **Filter Indicators**:
  - Show active filters above list
  - Each filter as removable badge

### Components Used
- ListModel (scrollable list with selection)
- DetailsModel (detail view with rich formatting)
- ActionMenuModel (context menu)
- SearchModel (search input with auto-complete)
- FilterModel (filter configuration)

### Data Transformations
- Database query → CareerEvent list
- Filter criteria → SQL WHERE clause
- Search text → Full-text search query
- Sort order → SQL ORDER BY
- Pagination offset → LIMIT/OFFSET

---

## GenerateCV Intent

### Special Logic
- **Multi-Step Wizard**: 5 steps with progress indicator
  1. Profile: Select/create profile
  2. Audience: Select target audience
  3. Content: Select which events/skills to include
  4. Template: Choose CV template/format
  5. Preview: Review before export
- **Conditional Steps**: Some steps appear only if conditions met
  - Template selection only if multiple templates available
  - Content selection only if multiple events available
- **State Preservation**: Resume from any step if interrupted
- **Preview Caching**: Don't regenerate preview unless settings changed
- **Template Variables**: Different templates have different variables

### Edge Cases
- No events (show message, can't proceed)
- No profile (offer to create)
- No template available (show error)
- Very large CV content (pagination in preview)
- Special characters in content (proper escaping)
- Encoding issues (UTF-8 validation)
- PDF generation timeout (show progress, allow retry)
- Disk space issues (graceful error)
- Concurrent CV generation (queue requests)

### Custom Styling
- **Progress Bar**:
  - Shows current step and total steps (e.g., "Step 2 of 5")
  - Colored segments for completed/current/remaining steps
  - Clickable to jump to step (if allowed)
- **Step Content**:
  - Left sidebar: Step list with numbers
  - Center: Step content form/selection
  - Right: Preview pane (if available)
- **Buttons**:
  - Previous/Next buttons at bottom
  - Previous disabled on first step
  - Next disabled until current step valid
  - Finish button on last step
- **Preview Pane**:
  - Scrollable, shows live preview
  - Syntax highlighting for code sections
  - Page break indicators
  - Zoom controls (if applicable)

### Components Used
- FormModel (for profile/audience selection)
- ListModel (for event/skill selection)
- CVPreviewModel (for preview rendering)
- ProgressModel (for step progress)

### Data Transformations
- Profile selection → CV configuration
- Event list → CV content sections
- Skills list → CV skills section
- Template selection → Template engine
- CV data → PDF/DOCX/TXT output

---

## ExportArtifact Intent

### Special Logic
- **Format Selection**: Different formats have different options
  - PDF: Font selection, margin size, page size
  - DOCX: Template selection, styling options
  - TXT: Line width, encoding
  - JSON: Pretty-print, include metadata
- **Export Progress**: Real-time progress updates
  - File generation
  - Compression (if applicable)
  - Upload (if cloud storage)
- **File Naming**: Auto-generate or user-specified
  - Auto: "CV_[Profile]_[Date].pdf"
  - User: Custom name with validation
- **Destination Selection**: Local file or cloud storage
- **Overwrite Confirmation**: If file exists, ask user

### Edge Cases
- File system full (graceful error)
- Invalid filename characters (auto-sanitize)
- Permission denied (show error)
- File already open in another program (retry or alternative)
- Network timeout (retry mechanism)
- Large file generation (show progress)
- Corrupted output file (validation check)
- Encoding issues (UTF-8 conversion)

### Custom Styling
- **Format Selection Grid**:
  - 4 format options displayed as cards
  - Each card shows format icon, name, and description
  - Selected card has highlighted border
  - Hover shows additional info
- **Options Panel**:
  - Format-specific options displayed below grid
  - Grouped by category (e.g., "Appearance", "Content")
  - Inline help text for each option
  - Live preview updates as options change
- **Progress Bar**:
  - Shows percentage complete (0-100%)
  - Current step description below bar
  - Elapsed time and estimated time remaining
  - Cancel button (if long operation)
- **Success Screen**:
  - Large checkmark icon
  - File path/name displayed
  - "Open File" button (if applicable)
  - "Done" button to return

### Components Used
- FormModel (for format options)
- ProgressModel (for export progress)
- SuccessModel (for success feedback)

### Data Transformations
- CV data → Format-specific output
- Options → Format-specific parameters
- Template → Output structure
- Encoding → File encoding

---

## BurstManagement Intent (New)

### Special Logic
- **Burst Creation**: From selected events
  - Validate events can be grouped (date proximity, skill overlap)
  - Auto-suggest burst title from event titles
  - Calculate burst metrics (total duration, skills, events)
- **Burst Suggestion**: AI-powered suggestions
  - Analyze events for patterns
  - Suggest grouping based on skills/dates/outcomes
  - Show confidence score for each suggestion
  - Allow accept/reject/customize
- **Relationship Management**: Track burst ↔ events
  - Add/remove events from burst
  - Cascade delete handling
  - Maintain referential integrity
- **Validation Rules**:
  - Burst must have at least 2 events
  - Burst title required, max 255 chars
  - Date range must be valid
  - No overlapping burst periods (optional constraint)

### Edge Cases
- Single event (can't create burst)
- No events with overlapping dates
- Circular dependencies (prevent)
- Deleting event in burst (cascade or prevent)
- Empty burst after edits (prevent save)
- Duplicate bursts (detect and warn)
- Very large bursts (100+ events, performance)
- Concurrent burst creation (race condition prevention)

### Custom Styling
- **Burst Card List**:
  - Cards show: title, date range, event count, skill tags
  - Selection: inverted colors
  - Hover: shadow effect
  - Icon for burst type (if applicable)
- **Burst Details View**:
  - Header: title, date range, metrics
  - Content: list of events in burst
  - Sidebar: actions, metadata
  - Edit button to modify
- **Suggestion Cards**:
  - Show suggested grouping
  - Confidence score as percentage
  - Preview of events in group
  - Accept/Reject buttons
  - Customize button to modify before accepting

### Components Used
- ListModel (for burst list)
- DetailsModel (for burst details)
- FormModel (for burst editor)
- CardModel (for suggestion cards)

### Data Transformations
- Event list → Burst suggestions (AI)
- Selected events → Burst entity
- Burst entity → Database record
- Database record → UI display

---

## FactManagement Intent (New)

### Special Logic
- **Fact Extraction**: From events/bursts
  - Parse event descriptions for facts
  - Extract structured data (metrics, outcomes, skills)
  - Quality scoring (confidence level)
- **Fact Deduplication**: Detect and merge duplicate facts
  - Fuzzy matching on fact content
  - Merge with higher confidence score
  - Track merge history
- **Quality Assessment**: Multi-criteria scoring
  - Specificity (concrete vs. vague)
  - Measurability (quantifiable vs. qualitative)
  - Relevance (to career goals)
  - Completeness (all fields filled)
- **Search & Filter**: Powerful fact discovery
  - Full-text search on fact content
  - Filter by quality score
  - Filter by source (event/burst)
  - Filter by date range
  - Filter by tags/skills

### Edge Cases
- Event with no extractable facts
- Fact too generic (quality score low)
- Fact with missing required fields
- Orphaned facts (event deleted)
- Duplicate facts (merge strategy)
- Very long fact content (truncation)
- Special characters in facts (escaping)
- Concurrent fact edits (conflict resolution)
- Fact quality changes over time (history tracking)

### Custom Styling
- **Fact Card**:
  - Compact view: content, quality score, source
  - Quality score as colored bar (red/yellow/green)
  - Source tag (event/burst name)
  - Click to expand for full content
- **Fact List**:
  - Sortable by: date, quality, relevance
  - Filter indicators at top
  - Search bar with auto-complete
  - Pagination for large lists
- **Fact Editor**:
  - Multi-field form
  - Quality score updates as you edit
  - Suggestion dropdown for similar facts
  - Validation feedback
- **Results View** (from import):
  - Table format: fact, quality, action
  - Bulk actions: accept all, reject all, review
  - Expandable rows for details

### Components Used
- ListModel (for fact list)
- DetailsModel (for fact details)
- FormModel (for fact editor)
- CardModel (for fact cards)
- TableModel (for import results)

### Data Transformations
- Event description → Extracted facts (NLP)
- Facts → Quality scores (multi-criteria)
- User edits → Fact entity
- Fact entity → Database record

---

## ImportWizard Intent (New)

### Special Logic
- **File Selection**: Browse and select CSV file
  - Validate file exists and readable
  - Check file size (warn if > 100MB)
  - Preview first few rows
- **Data Preview**: Show import preview
  - Display headers and sample rows
  - Detect data issues (missing fields, wrong format)
  - Show import statistics (total rows, valid rows, errors)
- **Conflict Resolution**: Handle existing data
  - Duplicate detection (compare with existing events)
  - Merge strategy selection (skip, replace, merge)
  - Manual conflict resolution for edge cases
- **Progress Tracking**: Real-time import progress
  - Rows processed / total rows
  - Estimated time remaining
  - Current operation (parsing, validating, saving)
  - Pause/resume capability
- **Error Handling**: Comprehensive error reporting
  - Row-level errors (specific field, error message)
  - Rollback on critical error (or continue with warnings)
  - Error log for review after import

### Edge Cases
- Empty CSV file
- CSV with wrong headers (validation)
- CSV with special characters (encoding)
- CSV with dates in wrong format (detection/conversion)
- CSV with duplicate rows (handling)
- CSV larger than memory (streaming)
- Partial import failure (recovery)
- Network interruption (resume capability)
- Concurrent imports (queue/prevent)
- Invalid data types (conversion/rejection)

### Custom Styling
- **File Picker**:
  - Current directory path
  - File list with size and date
  - File type filter (CSV only)
  - Selected file highlighted
- **Data Preview**:
  - Header row bold/different color
  - Sample data in table format
  - Row numbers on left
  - Scroll if many columns
  - Error indicators (red X) for problem rows
- **Progress Screen**:
  - Large progress bar with percentage
  - Step indicator (parsing/validating/importing)
  - Row count display (X of Y processed)
  - Time remaining estimate
  - Pause/Resume buttons
  - Cancel button
- **Results Screen**:
  - Summary: total imported, skipped, errors
  - Error log (if any)
  - "View Details" button for each error
  - "Done" button to return

### Components Used
- FilePicker (for file selection)
- TableModel (for data preview)
- ProgressModel (for import progress)
- FormModel (for conflict resolution)

### Data Transformations
- CSV file → Row data
- Row data → CareerEvent domain models
- Validation rules → Validation errors
- Merge strategy → Database operations

---

## MetadataEditor Intent (New)

### Special Logic
- **Metadata Review**: Show current metadata
  - Event metadata (date, duration, location, etc.)
  - System metadata (created date, modified date, version)
  - Computed metadata (quality score, relevance, etc.)
- **Change Tracking**: Track all edits
  - Original vs. modified values
  - Highlight changed fields
  - Change history/audit log
  - Rollback capability
- **Validation**: Metadata constraints
  - Required fields validation
  - Format validation (dates, numbers, etc.)
  - Relationship validation (foreign keys)
  - Business rule validation
- **Batch Editing**: Edit multiple events' metadata
  - Apply changes to selection
  - Preview affected items
  - Confirmation before save

### Edge Cases
- Metadata with missing required fields
- Circular relationships (prevent)
- Stale metadata (cache invalidation)
- Concurrent metadata edits (conflict resolution)
- Large batch operations (performance)
- Invalid relationships (validation)
- Orphaned records (cleanup)

### Custom Styling
- **Two-Column Layout**:
  - Left: Original values (read-only, gray)
  - Right: Edited values (editable, white)
  - Changed fields highlighted with border
  - Separator line between columns
- **Field Display**:
  - Label on left, value on right
  - Type-specific input (date picker, number input, etc.)
  - Validation feedback (red border on error)
  - Help text below field (gray, smaller)
- **Change Summary**:
  - List of changed fields at top
  - Each item shows: field name, old value → new value
  - Removable (undo individual changes)
  - "Undo All" button
- **Buttons**:
  - Save: Primary color, enabled only if changes
  - Reset: Secondary color, resets to original
  - Cancel: Outlined, returns without saving

### Components Used
- FormModel (for metadata editing)
- DetailsModel (for metadata review)
- TableModel (for change history)

### Data Transformations
- Database metadata → UI display
- User edits → Metadata changes
- Changes → Audit log entries
- Updated metadata → Database save

---

## BulkOperations Intent (New)

### Special Logic
- **Operation Selection**: Choose bulk operation
  - Available operations: delete, tag, archive, export, duplicate
  - Each operation shows: name, description, icon, estimated time
- **Configuration**: Set operation parameters
  - Operation-specific configuration
  - Scope selection (selected items, filtered items, all items)
  - Preview affected items
- **Execution**: Run operation with progress
  - Real-time progress updates
  - Item-by-item feedback (success/error)
  - Pause/resume capability
  - Rollback on critical error (or continue with warnings)
- **Results**: Summary of operation results
  - Items processed: X
  - Items succeeded: Y
  - Items failed: Z
  - Error details for failures

### Edge Cases
- Empty selection (show message)
- Very large selection (100+ items, performance)
- Operation timeout (retry or cancel)
- Partial failure (rollback or continue)
- Concurrent operations (queue or prevent)
- Permission denied (show error)
- Data conflicts (conflict resolution)
- Insufficient disk space (for export)

### Custom Styling
- **Operation Selection**:
  - Grid of operation cards
  - Each card: icon, name, description
  - Selected card: highlighted border
  - Hover: shadow effect
- **Configuration Panel**:
  - Form with operation-specific fields
  - Preview section showing affected items
  - Item count (e.g., "Affects 42 events")
  - Confirm button (large, primary color)
- **Progress Screen**:
  - Large progress bar with percentage
  - Item list (current item highlighted)
  - Status for each item (pending, processing, done, error)
  - Pause/Resume/Cancel buttons
- **Results Screen**:
  - Large checkmark (success) or warning icon
  - Summary statistics
  - Error log (if any)
  - "View Details" button for each error
  - "Done" button to return

### Components Used
- GridModel (for operation selection)
- FormModel (for configuration)
- TableModel (for item list and results)
- ProgressModel (for operation progress)

### Data Transformations
- Selected items → Operation input
- Configuration → Operation parameters
- Operation execution → Result data
- Result data → Summary statistics

---

## Cross-Model Dependencies

### FormModel Dependencies
Used by:
- CaptureEvent (event form)
- BurstManagement (burst editor)
- FactManagement (fact editor)
- MetadataEditor (metadata form)
- ImportWizard (conflict resolution)

**Shared Features**:
- Multi-field form with validation
- Tab navigation
- Character counters
- Error messages
- Submit/Cancel buttons

**Customization Points**:
- Field definitions (types, validation rules)
- Layout (number of columns, field order)
- Styling (colors, spacing)
- Behavior (auto-save, validation timing)

### ListModel Dependencies
Used by:
- BrowseTimeline (event list)
- BurstManagement (burst list)
- FactManagement (fact list)
- BulkOperations (item list)

**Shared Features**:
- Scrollable list with selection
- Keyboard navigation (j/k or arrow keys)
- Search/filter
- Pagination
- Sorting

**Customization Points**:
- Item rendering (card vs. row)
- Selection behavior (single vs. multiple)
- Sort options
- Filter options
- Item height

### DetailsModel Dependencies
Used by:
- BrowseTimeline (event details)
- BurstManagement (burst details)
- FactManagement (fact details)

**Shared Features**:
- Rich content display
- Metadata section
- Action buttons
- Related items

**Customization Points**:
- Content sections
- Metadata fields
- Action buttons
- Related item display

---

## Service Dependencies

### CareerService
Used by all intents for:
- Event CRUD operations
- Burst CRUD operations
- Fact CRUD operations
- Query/filtering operations
- Validation operations

### CLIEventService
Used by:
- CaptureEvent (event creation)
- BrowseTimeline (event display)
- BulkOperations (bulk operations)

### ImportService
Used by:
- ImportWizard (file import, parsing, conflict resolution)

### Logger
Used globally for:
- Error logging
- Debug logging
- Audit logging

---

## Styling Standards

### Color Palette
- **Primary**: #0099FF (blue)
- **Secondary**: #FF9900 (orange)
- **Success**: #00CC00 (green)
- **Error**: #FF0000 (red)
- **Warning**: #FFFF00 (yellow)
- **Background**: #1E1E1E (dark gray)
- **Text**: #FFFFFF (white)
- **Border**: #333333 (darker gray)

### Typography
- **Headers**: Bold, size 16-20
- **Body**: Regular, size 12-14
- **Labels**: Bold, size 12
- **Help Text**: Gray, size 10-11
- **Monospace**: For code/data

### Spacing
- **Padding**: 1-2 units (2-4 characters)
- **Margin**: 1-3 units (2-6 characters)
- **Gap**: 1 unit (2 characters) between items

### Borders
- **Standard**: Single line, border color
- **Focus**: Double line, primary color
- **Error**: Single line, error color
- **Success**: Single line, success color

### Icons
- Used consistently across intents
- Standard icons from Lipgloss library
- Custom icons for domain-specific items

---

## Accessibility Patterns

### Keyboard Navigation
- Tab/Shift+Tab to move between fields
- Arrow keys to navigate lists
- Enter to select/confirm
- Esc to cancel/go back
- j/k as vim-style navigation alternatives

### Focus Indicators
- Clear focus state (colored border or highlight)
- Focus order follows logical flow
- Focus visible in all states

### Error Messages
- Clear, specific error descriptions
- Positioned near the error source
- Colored for visibility

### Help Text
- Available for complex fields
- Accessible via '?' or help button
- Context-sensitive

---

## Performance Considerations

### List Rendering
- Lazy load large lists (100+ items)
- Virtual scrolling for very large lists
- Pagination (20 items per page)
- Caching of rendered items

### Form Rendering
- Debounce validation (avoid excessive checks)
- Lazy load complex fields
- Cache field options

### Preview Generation
- Cache generated previews
- Regenerate only on setting changes
- Async generation for large previews

### Database Queries
- Optimize queries for common patterns
- Add indexes for frequently searched fields
- Use pagination for large result sets
- Eager load related entities

---

**Document Version**: 1.0
**Last Updated**: 2026-01-03
**Status**: Complete - Ready for Task 1.2.x

