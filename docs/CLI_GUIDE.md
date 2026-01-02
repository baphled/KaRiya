# KaRiya CLI Guide

## Overview

KaRiya CLI is a terminal-based career journaling tool that helps you capture, organize, manage, and enrich your professional career events and milestones.

## Installation

```bash
# Build from source
go build -o kariya-cli ./cmd/cli

# Run
./kariya-cli
```

## Quick Start

### Basic Usage

```bash
# Start the application
./kariya-cli

# View help
./kariya-cli --help

# Check version
./kariya-cli --version
```

### Using Command-Line Flags

```bash
# Start with custom database path
./kariya-cli --db /path/to/events.db

# Start in specific capture mode
./kariya-cli --mode timeline

# Show recent events on startup
./kariya-cli --list

# Import CSV with automatic burst/fact detection
./kariya-cli --import events.csv

# Re-run burst detection on all events
./kariya-cli --detect-bursts

# Re-run fact extraction on all events
./kariya-cli --extract-facts

# View all existing bursts
./kariya-cli --show-bursts

# View all existing facts
./kariya-cli --show-facts
```

## Features

### 1. Event Capture

Capture career events in three ways:

#### Timeline Journaling
Use this to log events as they happen (limited to last 30 days)

```bash
./kariya-cli --mode timeline
```

**Best for**: Real-time event logging, building career journals

**Constraints**:
- Event dates must be within last 30 days
- Recent captures are fresher in memory
- Good for immediate reflection

#### CV Backfill
Import historical events from your CV or career history

```bash
./kariya-cli --mode backfill
```

**Best for**: Building comprehensive career history, importing from CV

**Advantages**:
- Accepts any date in the past
- No time constraints
- Perfect for filling gaps in career timeline

#### Manual Entry
Flexible entry for any event at any time

```bash
./kariya-cli --mode manual
```

**Best for**: Flexible event entry, capturing events from any time period

**Features**:
- No date restrictions
- Can be used as default capture mode
- Quick metadata enrichment after capture

### 2. Metadata Review & Enrichment

After capturing an event, you can immediately review and enrich its metadata:

#### Success Screen Options

After capturing an event, you'll see the success screen with the following options:

- **Capture Another**: Continue capturing more events
- **Review Metadata**: Enrich the captured event with additional details
- **View Recent**: See recent events
- **Exit**: Exit the application

#### Metadata Review Screen

Access the metadata review screen to:
- View all events with data quality indicators
- See which events need metadata enrichment
- Filter by quality level (Incomplete, Basic, Enriched, Complete)
- Sort by date, company, or creation order

**Keyboard shortcuts**:
- `↑/↓`: Navigate through events
- `Enter`: Edit selected event
- `Space`: Select event for bulk operations
- `a`: Select all events
- `d`: Deselect all events
- `e`: Enter bulk operations mode
- `m`: Open metadata review from home screen

#### Individual Event Editor

Edit metadata for a single event:
- **Date**: Change event date (format: YYYY-MM-DD or relative dates like "2 weeks ago")
- **Company**: Add or update company name
- **Project**: Add or update project name
- **Tags**: Select from available tags (max 8)
- **Categories**: Select event category

**Keyboard shortcuts**:
- `Tab/Shift+Tab`: Navigate between fields
- `Space`: Toggle tag/category selection
- `Enter`: Save changes
- `Esc`: Cancel without saving
- `Ctrl+Z`: Undo changes

#### Bulk Operations

Edit metadata for multiple events at once:

1. Select events (Space to toggle, 'a' to select all, 'd' to deselect all)
2. Press 'e' to enter bulk edit mode
3. Choose which fields to update
4. Preview changes before confirming
5. Confirm to apply changes to all selected events

**Features**:
- Conditional updates: "Apply if field is empty"
- Preview changes before applying
- Undo/revert functionality
- Transaction-like behavior (all succeed or all fail)

### 3. Event Listing & Filtering

View and filter your events:

- **List View**: See all captured events
- **Filter by Company**: Find events from specific companies
- **Filter by Date Range**: View events from specific time periods
- **Search**: Full-text search across event descriptions
- **Sort**: Order by date, company, or creation time

### 4. Event Details

View detailed information for a single event:

- Event description
- Date captured
- Company and project
- Tags and categories
- Data quality score
- Event metadata completeness

### 5. CSV Import

Import events from CSV files:

1. Prepare CSV with columns: `description`, `date`, `company`, `project`, `tags`
2. From home screen, select "Import from CSV"
3. Choose CSV file
4. Review parsed events
5. Select which events to import
6. Automatically navigate to metadata review for imported events
7. Review and enrich imported event metadata

**CSV Format Example**:
```csv
description,date,company,project,tags
"Led team through migration",2024-06-15,TechCorp,Platform Migration,"technical,leadership"
"Implemented critical fix",2024-06-10,TechCorp,Product,"technical,achievement"
```


### 6. Burst Detection & Grouping

Bursts are automatically detected groupings of 2+ related career events:

1. Complete metadata clarification for your events
2. From metadata review screen, press 'u' to view burst suggestions
3. Review suggested bursts with confidence scores:
   - **Green (>0.7)**: High confidence - events clearly related
   - **Yellow (0.4-0.7)**: Medium confidence - review before confirming
   - **Red (<0.4)**: Low confidence - usually not suggested
4. For each burst:
   - Review related events
   - Confirm or reject suggestion
   - Optionally edit burst name/description (press 'e')
5. Confirmed bursts persist to database

**How Burst Detection Works**:
- Text similarity (word overlap and keywords)
- Company/project matching
- Tags matching
- Temporal proximity (events within 6 months)

**Example Burst**:
```
Platform Modernization Initiative (4 events)
- "Designed microservices architecture"
- "Led migration planning and execution"
- "Implemented deployment automation"
- "Mentored team on microservices patterns"
Competency Focus: Technical
Confidence: 85%
```

### 7. Fact Extraction & Enrichment

Facts are automatically extracted inferences about your competencies and achievements:

1. After confirming burst suggestions, facts are automatically extracted
2. For each event/burst, review extracted facts:
   - **Text**: What you accomplished (no aspirational language)
   - **Competencies**: 6 categories (Technical, Leadership, Product, Consulting, Research, Mentoring)
   - **Role Fit**: Career level (Principal, EM, Staff Engineer, Senior IC)
   - **Audience**: Who cares (Hiring Manager, Recruiter, Peer)
3. Confirm or edit each fact
4. Facts are used for CV generation and portfolio creation

**Fact Validation Rules**:
- No aspirational language (will, should, could, might, may, want, wish, hope, plan, intend, attempt, try, would)
- Grounded metrics (measured, not speculative)
- 1-2000 characters
- ≥1 competency category
- ≥1 audience type

**Example Fact**:
```
Text: "Architected microservices platform supporting 50M+ requests daily"
Competencies: [Technical, Architecture]
Role Fit: Principal
Audience: [Hiring Manager, Peer]
Strength Signal: "Architected system handling massive scale"
```

## Keyboard Reference

### Home Screen
- `c`: Capture new event
- `i`: Import from CSV
- `l`: List recent events
- `m`: Open metadata review
- `q`: Quit application

### List View
- `↑/↓`: Navigate through events
- `Enter`: View event details
- `/`: Search events
- `f`: Filter events
- `Esc`: Return to home

### Capture Screen
- `Tab/Shift+Tab`: Navigate between fields
- `Enter`: Submit event
- `Esc`: Cancel capture

### Success Screen
- `←/→`: Navigate between options
- `Enter`: Select option
- `Esc`: Return to home

### Metadata Review Screen
- `↑/↓`: Navigate through events
- `Enter`: Edit selected event
- `Space`: Select/deselect event
- `a`: Select all events
- `d`: Deselect all events
- `e`: Enter bulk operations mode
- `f`: Filter by quality level
- `s`: Change sort order
- `Esc`: Return to previous screen

### Metadata Editor
- `Tab/Shift+Tab`: Navigate between fields
- `Space`: Toggle tag/category selection
- `Enter`: Save changes
- `Esc`: Cancel without saving
- `Ctrl+Z`: Undo changes

### Bulk Operations
- `↑/↓`: Navigate through events
- `Space`: Toggle selection
- `a`: Select all
- `d`: Deselect all
- `e`: Enter edit mode
- `Enter`: Confirm changes
- `Esc`: Cancel

For quick reference while in the app, press 'h' to open the help system which includes complete keyboard shortcuts.

## Workflows

### Quick Event Capture with Metadata Enrichment

1. Start application: `./kariya-cli`
2. Press `c` to capture
3. Enter event description
4. (Optional) Add date, company, project
5. Press `Enter` to submit
6. On success screen, press `→` to select "Review Metadata"
7. Press `Enter` to go to metadata editor
8. Add company, project, tags, categories
9. Press `Enter` to save
10. Back at metadata review, press `Esc` to return home

### Bulk Metadata Enrichment

1. Press `m` to open metadata review
2. See events with incomplete metadata (color-coded quality indicators)
3. Press `Space` to select multiple events
4. Press `a` to select all events needing enrichment
5. Press `e` to enter bulk operations
6. Select fields to update (e.g., Company, Tags)
7. Enter values
8. Press `Enter` to preview changes
9. Press `Enter` to confirm and apply

### CSV Import with Metadata Review

1. From home screen, press `i` for import
2. Select CSV file with event data
3. Review parsed rows
4. Select which rows to import
5. Automatically navigates to metadata review
6. Review imported events
7. Edit individual events or bulk edit multiple events
8. Complete metadata enrichment

### Export for CV

1. Open list view (press `l`)
2. Filter by "achievement" tag (press `f`)
3. Select events you want to export
4. (Future) Export to formatted document

## Data Quality Scoring

Events are automatically scored for data completeness:

- **Incomplete** (0-25): Only event description
- **Basic** (26-50): Description + date
- **Enriched** (51-75): Description + date + company/project
- **Complete** (76-100): All fields filled with tags/categories

Quality indicators show:
- Current score as percentage
- Quality level badge
- Missing fields that would improve the score
- Visual color coding (red → yellow → green)

## Tips & Tricks

1. **Quick Metadata**: Use "Review Metadata" option immediately after capture
2. **Tag Strategy**: Use overlapping tags (e.g., "technical" + "achievement" for good technical wins)
3. **Date Entry**: Use relative dates like "2 weeks ago" instead of calculating exact dates
4. **Bulk Operations**: For similar events, capture one fully, then bulk-copy metadata to others
5. **Company Consistency**: Use exact same company names for better filtering
6. **Regular Backups**: Periodically backup your database file
7. **CSV Templates**: Create reusable CSV templates for regular import workflows

## Troubleshooting

### Issue: Can't capture events

**Check**:
- Ensure you're in one of the three capture modes (Timeline, Backfill, Manual)
- Timeline mode limits dates to last 30 days
- Use CV Backfill or Manual Entry for older dates

### Issue: Metadata changes not saved

**Check**:
- Ensure you press `Enter` to save changes (not just `Esc`)
- Check that all required fields are filled
- Verify validation errors are resolved (shown in red)

### Issue: Bulk operations failed

**Check**:
- Ensure at least one event is selected
- Check that metadata values are valid
- Review validation errors for specific fields
- Try updating fewer fields at once

### Issue: CSV import has parsing errors

**Check**:
- CSV format matches expected columns
- Date format is YYYY-MM-DD or recognizable
- No special characters in field values that might confuse parser
- Try with smaller CSV file first

### Issue: Search is slow

**Check**:
- For large databases (10,000+ events), search may take a moment
- Use filters first to narrow down results
- Try more specific search terms

### Issue: Application Crashes

**Steps**:
1. Check terminal size is adequate (80+ columns)
2. Try `./kariya-cli --help` to verify installation
3. Check Go version: `go version` (1.24.0+)
4. Review error message for hints

## File Locations

- **Default Database**: Current directory (memory-based by default)
- **Custom Database**: Specify with `--db` flag
- **Config**: Command-line flags only (currently)

## Performance

- **Startup Time**: < 1 second
- **Event Listing**: Instant for 1000+ events
- **Search**: Real-time with minimal lag
- **Filter**: Immediate response
- **Metadata Review**: Loads 100+ events in < 500ms
- **Database**: SQLite supports 100,000+ events

## CV Generation Workflow

KaRiya can generate professional CVs from your career events, tailored to specific roles and audiences.

### Quick Start

1. **Create a CV Configuration**
   ```bash
   # From the main menu, select "Manage CV Configs"
   # Press 'n' to create new configuration
   # Fill in: CV Name, Target Role, Target Audience
   # Optionally add filters (date range, companies, tags)
   # Press Enter to save
   ```

2. **Generate a CV**
   ```bash
   # Select your configuration from the list
   # Press Enter to generate
   # KaRiya will create your CV in seconds
   ```

3. **Preview and Export**
   ```bash
   # View the generated CV with bullets and sections
   # Press 'e' to export (text, markdown, or clipboard)
   # Exports are saved to ~/.kariya/cv_exports/
   ```

### CV Configuration

CV configurations are stored as YAML files in `~/.kariya/cv_configs/`

**Example Configuration**:
```yaml
name: "Staff Engineer CV"
targetRole: "Staff"
targetAudience:
  - "HiringManager"
  - "Peer"
eventFilters:
  startDate: "2023-01-01"
  endDate: "2024-12-31"
  companies:
    - "TechCorp"
  tags:
    - "technical"
    - "leadership"
```

### Target Roles

- **Principal**: 3-4 bullets max, emphasizes strategic impact
- **Staff**: 4-5 bullets max, emphasizes technical depth
- **EM**: 3-4 bullets max, emphasizes team growth and delivery
- **SeniorIC**: 4-5 bullets max, emphasizes technical leadership

### Target Audiences

- **HiringManager**: Focus on business impact and outcomes
- **Recruiter**: Focus on skills, competencies, and career growth
- **Peer**: Focus on technical depth and collaboration

### Keyboard Shortcuts for CV Workflow

| Screen | Key | Action |
|--------|-----|--------|
| Config Manager | `j`/`k` | Navigate list |
| Config Manager | `Enter` | Generate CV |
| Config Manager | `n` | New config |
| Config Manager | `e` | Edit config |
| Config Manager | `d` | Delete config |
| Config Editor | `Tab` | Next field |
| Config Editor | `Shift+Tab` | Previous field |
| Config Editor | `Enter` | Save config |
| Config Editor | `Esc` | Cancel |
| CV Preview | `j`/`k` | Navigate bullets |
| CV Preview | `→`/`←` | Navigate sections |
| CV Preview | `Enter` | View sources |
| CV Preview | `e` | Export CV |

### Workflow Examples

**Example 1: Create Staff Engineer CV**
```
1. Select "Manage CV Configs" from main menu
2. Press 'n' to create new config
3. Name: "Staff Engineer CV"
4. Role: "Staff"
5. Audience: "HiringManager"
6. Leave filters blank (include all events)
7. Press Enter
8. Select config and press Enter to generate
9. Review CV in preview mode
10. Press 'e' to export
```

**Example 2: Create Role-Specific CVs**
```
Create multiple configs:
- "Principal CV" (target: Principal)
- "Staff CV" (target: Staff)
- "Manager CV" (target: EM)

Generate each and compare how the same events
produce different bullets for different roles
```

**Example 3: Audience-Specific CVs**
```
Create two configs with same role, different audiences:
- "Staff for HiringManager"
- "Staff for Recruiter"

Generate both to see how the same events
are presented differently for different audiences
```

### Bullet Generation Rules

Bullets are automatically generated from your events following these rules:

**Inclusion Criteria**:
- Must trace to at least one event
- Single claim only (not multiple unrelated achievements)
- No aspirational language (no "will", "aims to", etc.)
- No inferred metrics (only metrics from source events)
- No role inflation (accurate ownership representation)

**Ranking Priority**:
1. Ownership (highest) - "Designed and built X"
2. Contribution - "Contributed to X"
3. Strategy - "Defined strategy for X"
4. Execution - "Implemented X"
5. Outcome - "Improved X by Y%"
6. Activity (lowest) - "Worked on X"

**Compression**:
- When exceeding role bullet cap, lower-priority bullets are removed
- Older work is deprioritized
- At least one bullet per section is preserved

### Traceability

Every bullet in your CV traces back to source events:
1. Navigate to a bullet in CV preview
2. Press Enter to view sources
3. See which events/facts contributed
4. Check confidence score and inclusion reason

### Export Formats

**Text Export**:
- Plain text format
- Suitable for pasting into email/messaging
- File: `~/.kariya/cv_exports/cv_name_YYYY-MM-DD_HH-MM-SS.txt`

**Markdown Export**:
- Markdown format with formatting
- Suitable for GitHub, documentation
- File: `~/.kariya/cv_exports/cv_name_YYYY-MM-DD_HH-MM-SS.md`

**Clipboard**:
- Copy directly to clipboard
- No file saved
- Quick sharing option

### Tips & Best Practices

1. **Event Quality Matters**: Better events produce better CV bullets
2. **Use Tags Strategically**: Tag events to enable focused CVs
3. **Create Multiple Configs**: Don't try one CV for all purposes
4. **Review Generated Bullets**: Verify accuracy and traceability
5. **Understand Compression**: Know why certain bullets were excluded
6. **Use Filters**: Create specialized CVs for specific roles/companies

### For More Information

- **Comprehensive Guide**: [CV_GENERATION_GUIDE.md](../guides/CV_GENERATION_GUIDE.md)
- **Examples**: [CV_EXAMPLES.md](../guides/CV_EXAMPLES.md)
- **In-App Help**: Press 'h' while in CV screens


## Features & Capabilities

### Current Features ✅
- Event capture (Timeline, Backfill, Manual)
- Metadata review and enrichment
- Individual event editing
- Bulk metadata operations
- CSV import with automatic metadata review
- Data quality scoring
- Event filtering and search
- Keyboard-driven navigation
- Multi-screen navigation

### Future Features 🚀
- Web interface
- Export to multiple formats (JSON, CSV, PDF)
- Advanced analytics and reporting
- Burst detection and grouping
- Fact extraction from events
- Integration with external services
- Cloud synchronization

## Getting Help

1. **In-App Help**: Press 'h' while running
2. **This Guide**: `docs/CLI_GUIDE.md`
3. **Help Flag**: `./kariya-cli --help`
4. **Version**: `./kariya-cli --version`
5. **Troubleshooting**: See section above
6. **Code Examples**: Check `internal/cli/models/*_test.go`

## Support

For issues or suggestions:
1. Check the troubleshooting section above
2. Review the help system in the app
3. Check test coverage in `internal/cli/models/*_test.go`
4. Review code in `internal/cli/`

## License

KaRiya is part of the career journaling project.
