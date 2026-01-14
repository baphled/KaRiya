---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Feature: Metadata Clarification

## Purpose

Enable users to review, validate, and enrich event metadata (dates, companies,
projects, tags, categories) before automated processing (burst detection and
fact inference). Ensure data quality as the foundation for reliable burst
grouping and fact extraction.

## Core Concept

After capturing events (manually or via CSV import), users need a dedicated workflow to:
- Review event metadata for accuracy
- Correct or enhance metadata (dates, company, project, tags)
- Confirm or assign categories (competencies)
- Validate data quality before automated processing
- Handle bulk metadata operations for imported events

## User Stories

1. As a user importing 50 events from a CSV, I want to review and confirm the
   metadata is correct before the system creates bursts and facts, so that
   automated processing is based on clean data.

2. As a user who manually captured an event, I want to review the captured
   metadata and add missing information (company, project, tags) before it's
   grouped into bursts.

3. As a user with events from multiple companies, I want to bulk-confirm or
   bulk-adjust metadata (e.g., all events from "TechCorp" are correct) to speed
   up the review process.

4. As a user uncertain about event dates, I want to see a calendar or date
   picker to correct dates before bursts are created.

5. As a user with inconsistent tagging, I want to see which tags are applied and
   adjust them across multiple events before automated processing.

## Functional Requirements

### 4.1 Metadata Review Screen

1. The CLI must display a "Metadata Review" screen showing:
   - List of events awaiting metadata clarification
   - For each event:
     - Event text (truncated if long)
     - Current date
     - Current company
     - Current project
     - Current tags
     - Current categories
     - Data quality indicator (complete/incomplete)

2. The screen must support:
   - Scrolling through events (up/down arrows)
   - Viewing full event text (expand/collapse)
   - Filtering by data quality (incomplete only, all, etc.)
   - Sorting by date, company, or creation order

### 4.2 Individual Event Metadata Editing

1. Users must be able to edit metadata for individual events:
   - **Date**: Calendar picker or text input with format validation
   - **Company**: Text input with autocomplete from previous entries
   - **Project**: Text input with autocomplete from previous entries
   - **Tags**: Multi-select from AllowedTags with visual feedback
   - **Categories**: Multi-select from AllowedCategories with visual feedback

2. Editing must support:
   - Tab navigation between fields
   - Clear save/cancel buttons
   - Validation with helpful error messages
   - Undo/revert to original values

3. Field-specific validation:
   - **Date**: Not in future, valid format, reasonable range (not 50 years ago)
   - **Company**: Non-empty if project is provided, max 200 chars
   - **Project**: Non-empty if company is provided, max 200 chars
   - **Tags**: From AllowedTags, max 8 per event, no duplicates
   - **Categories**: From AllowedCategories, should match event content

### 4.3 Bulk Metadata Operations

1. For CSV imports or batch operations, users must be able to:
   - Select multiple events (checkbox per event, select all/none)
   - Bulk-edit common fields:
     - Apply same company to multiple events
     - Apply same project to multiple events
     - Add tags to multiple events
     - Add categories to multiple events
   - Preview bulk changes before confirming

2. Bulk operations must:
   - Show affected event count
   - Allow partial application (e.g., only if field is empty)
   - Support undo/revert
   - Provide clear confirmation before applying

### 4.4 Data Quality Indicators

1. Each event must show data quality status:
   - **Complete**: Text + date + (company OR standalone is ok)
   - **Incomplete**: Missing critical fields
   - **Enriched**: Has tags and/or categories
   - **Unreviewed**: Newly captured, not yet reviewed

2. Quality indicators must:
   - Use visual cues (color, icons)
   - Suggest what's missing
   - Provide quick actions to add missing data

### 4.5 CSV Import Integration

1. After CSV import, immediately show metadata review screen with:
   - All imported events pre-loaded
   - Indication of which fields came from CSV
   - Indication of any parsing issues or warnings
   - Duplicate detection status

2. Allow users to:
   - Review all imported metadata at once
   - Correct parsing errors (e.g., date format issues)
   - Confirm or adjust auto-detected categories
   - Add missing metadata before proceeding

3. Support batch operations for imports:
   - Bulk-confirm metadata for all events
   - Bulk-adjust company for events from same source
   - Bulk-add tags for events from same project

### 4.6 Post-Capture Metadata Review

1. After capturing an event, show a quick metadata review:
   - Display captured event with current metadata
   - Offer to add optional metadata (company, project, tags, categories)
   - Allow editing before saving
   - Option to "Add another event" or "Review metadata"

2. After N events captured, offer to review all at once:
   - "You've captured 5 events. Review metadata?"
   - Show all 5 events for batch operations

### 4.7 Keyboard Navigation

1. Support efficient keyboard navigation:
   - `↑` / `↓`: Navigate between events
   - `Enter`: Edit selected event
   - `e`: Edit current event
   - `Space`: Select/deselect event (for bulk operations)
   - `a`: Select all events
   - `d`: Deselect all events
   - `s`: Save/confirm changes
   - `Esc`: Cancel editing
   - `q`: Quit metadata review (if all required fields complete)

### 4.8 Workflow Integration

1. **After Manual Capture**:
   - Event captured → Quick metadata review → Option to add more → Back to
     capture or review all

2. **After CSV Import**:
   - CSV parsed → Review screen with all events → Bulk edit as needed → Confirm
     all → Proceed to next step

3. **On-Demand**:
   - User can access "Review Metadata" from home screen anytime
   - Shows all events with incomplete metadata
   - Allows reviewing and enriching existing events

## Validation Rules

### Date Validation

- Cannot be in the future
- Should be within reasonable range (not more than 50 years ago, unless backfilling)
- Supports flexible input (ISO 8601, relative dates, natural language)

### Company Validation

- Optional field
- If provided: max 200 characters
- If project is provided, company should also be provided
- Normalized to consistent casing (e.g., "TechCorp" not "TECHCORP")

### Project Validation

- Optional field
- If provided: max 200 characters
- If company is provided, project should also be provided

### Tags Validation

- Must be from AllowedTags
- Maximum 8 tags per event
- No duplicate tags
- Case-insensitive (normalized to lowercase)

### Categories Validation

- Must be from AllowedCategories
- Should reflect competencies demonstrated in event text
- Can be auto-suggested based on text analysis
- User can confirm, modify, or reject suggestions

## Allowed Metadata

### AllowedTags

- project
- achievement
- leadership
- technical
- consulting
- research
- product
- mentoring

### AllowedCategories

- technical
- leadership
- product
- consulting
- research
- mentoring

## Data Quality Scoring

Each event should have a quality score (0-100):
- **Text present**: +20 points (required)
- **Date present**: +20 points (required)
- **Company OR Project**: +20 points (at least one)
- **Tags assigned**: +15 points (optional)
- **Categories assigned**: +15 points (optional)
- **Categories match content**: +10 points (bonus)

**Quality Levels**:
- 0-40: Incomplete (missing required fields)
- 41-60: Basic (has required fields)
- 61-85: Enriched (has tags or categories)
- 86-100: Complete (all fields, verified match)

## Acceptance Criteria

1. Users can view all events awaiting metadata clarification
2. Users can edit metadata for individual events with clear validation
3. Users can perform bulk metadata operations on multiple events
4. Data quality is visible with clear indicators
5. CSV import automatically triggers metadata review
6. Manual capture offers metadata enrichment options
7. All metadata changes are validated before saving
8. Users can confirm metadata is ready before proceeding to burst extraction
9. Keyboard navigation supports efficient bulk operations
10. Changes are persisted to the database

## Non-Functional Requirements

### Performance

- Metadata review screen loads in <500ms for 100 events
- Bulk operations preview in <1s for 1000 events
- Metadata save in <100ms per event
- Date picker renders in <200ms

### Usability

- Clear visual hierarchy showing data quality
- Helpful error messages with correction suggestions
- Keyboard shortcuts for power users
- Mouse support for casual users
- Responsive layout for various terminal sizes

### Data Integrity

- No metadata changes without explicit save
- Undo/revert support for accidental changes
- Validation prevents invalid data entry
- Audit trail of metadata changes (future enhancement)

## Out of Scope

- Automatic metadata correction (ML-based)
- Metadata suggestions from external sources
- Metadata merging/deduplication
- Advanced filtering/search (basic filtering only)
- Metadata templates or presets

## Success Metrics

1. **Adoption**: % of imported events reviewed before burst extraction
2. **Data Quality**: % of events with complete metadata before bursts created
3. **Efficiency**: Average time to review metadata per event (target: <30s)
4. **Accuracy**: % of user-confirmed metadata matches actual content

## Implementation Notes

### Phase 2 Enables Phase 3

Once metadata clarification is complete, users should have:
- Clean, validated event data
- Accurate dates for date-based grouping
- Correct company/project associations for burst detection
- Appropriate tags for burst matching
- Confirmed categories for fact inference

This clean data foundation makes Phase 3 (Burst & Fact Extraction) much more reliable and accurate.

### Integration with Existing Features

- Builds on Phase 1 event capture
- Uses existing AllowedTags and AllowedCategories
- Leverages existing CLI screens and styling
- Prepares data for Phase 3 bursts and facts
- CSV import review screen enhancement (already partially exists)

### Future Enhancements

- Metadata templates for common patterns
- Autocomplete from historical data
- ML-based category suggestions
- Bulk import from external sources (LinkedIn, etc.)
- Metadata merge/consolidation tools

