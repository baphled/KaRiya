# KaRiya CLI Guide

## Overview

KaRiya CLI is a terminal-based career journaling tool that helps you capture, organize, and manage your professional career events and milestones.

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

**Best for**: One-off events, special accomplishments, flexible entry

**Features**:
- Complete date flexibility
- No constraints on event dates
- Perfect for adding events anytime

### 2. Event Fields

When capturing an event, you can provide:

- **Description** (required): What you did
  - 1-2000 characters
  - Be specific about the event
  - Include context and impact

- **Date** (optional): When it happened
  - Format: YYYY-MM-DD (e.g., 2025-12-24)
  - Relative: "1 week ago", "3 months ago"
  - Default: Today

- **Company** (optional): Where it happened
  - Company name or organization
  - Helps with organization and filtering

- **Project** (optional): Related project
  - Project or initiative name
  - Additional context

- **Tags** (optional): Categorize the event
  - Up to 8 tags per event
  - Available tags:
    - `project` - Major project work
    - `achievement` - Notable accomplishments
    - `leadership` - Leadership activities
    - `technical` - Technical work
    - `consulting` - Advisory/consulting
    - `research` - Research activities
    - `product` - Product management
    - `mentoring` - Mentoring & coaching

### 3. Browsing Events

View your captured events with multiple options:

- **List**: See all recent events with pagination
- **Filter**: Filter by date range, tags, company
- **Search**: Keyword search across all fields
- **Sort**: Order by date, creation time, or text
- **Details**: View full event information

### 4. Help System

Access comprehensive help directly from the CLI:

```
Press 'h' to open help at any time
```

Help sections:
1. Overview - What is KaRiya
2. Event Capture - Capture modes explained
3. Tagging - Tag system and best practices
4. Listing - How to browse events
5. Shortcuts - Keyboard reference
6. Search - Search tips and examples
7. Tips - Best practices and recommendations

## Keyboard Shortcuts

### Navigation

| Key | Action |
|-----|--------|
| `h` | Open help |
| `c` | Capture new event |
| `l` | List events |
| `q` | Quit |
| `backspace` | Go back |

### Form Navigation

| Key | Action |
|-----|--------|
| `tab` | Next field |
| `shift+tab` | Previous field |
| `up/down` | Navigate dropdowns |
| `enter` | Confirm/Submit |
| `esc` | Cancel |

### List Navigation

| Key | Action |
|-----|--------|
| `up/down` | Move between events |
| `left/right` | Previous/Next page |
| `pgup/pgdn` | Page navigation |
| `enter` | View details |

## Examples

### Example 1: Capture a Project Achievement

```
1. Press 'c' to capture
2. Enter description: "Led migration of legacy system to microservices architecture"
3. Set date: "2025-11-15"
4. Company: "TechCorp Inc."
5. Project: "System Architecture Modernization"
6. Tags: leadership, technical, achievement
7. Submit
```

### Example 2: Backfill Historical Events

```
./kariya-cli --mode backfill

1. Enter historical achievement from 2023
2. Use relative dates: "2 years ago" or ISO format
3. Add company and project context
4. Use appropriate tags for categorization
5. Multiple tags help with searching later
```

### Example 3: Quick Manual Entry

```
./kariya-cli --mode manual

1. Enter recent accomplishment
2. Leave date as today
3. Add company if applicable
4. Tag appropriately
5. Submit quickly
```

## Best Practices

### Event Capture

- **Be Specific**: Include concrete details, not generic statements
- **Add Context**: Company and project help with organization
- **Use Tags**: Multiple tags make finding events easier
- **Regular Updates**: Capture events regularly, not just retroactively
- **Consistent Format**: Try to maintain consistent style for easy reading

### Organization

- **Create Habits**: Capture events weekly or as they happen
- **Review Regularly**: Look through events to refresh memory
- **Use All Fields**: Complete entries are easier to search
- **Tag Consistently**: Use same tags for similar events
- **Combine Features**: Use search + filter + sort together

### CV Building

- **Export Filtered**: Show only relevant accomplishments
- **Highlight Impact**: Ensure descriptions show value
- **Chronological Order**: Sort by date for CV format
- **Company Context**: Include company names for clarity
- **Achievement Tags**: Use "achievement" tag for CV-worthy events

## Configuration

### Command-Line Flags

```bash
# Custom database location
./kariya-cli --db /home/user/my_events.db

# Start in specific capture mode
./kariya-cli --mode timeline
./kariya-cli --mode backfill
./kariya-cli --mode manual

# Show recent events on startup
./kariya-cli --list

# Display help
./kariya-cli --help

# Display version
./kariya-cli --version
```

### Default Behavior

- Database: In-memory (data lost on exit)
- Capture Mode: Manual (flexible entry)
- Startup Screen: Home menu
- Storage: MemoryRepository (configurable)

## Troubleshooting

### Issue: Events Not Showing Up

**Causes**:
- Using in-memory database (data not persisted)
- Events filtered out by active filters

**Solution**:
1. Use `--db` flag to specify SQLite database
2. Check active filters with 'f' key
3. Clear filters to see all events

### Issue: Can't Find an Event

**Solutions**:
1. Try search feature ('s' key) with keywords
2. Use filters to narrow down by date or tags
3. Sort by different fields to reorder
4. Check company/project names

### Issue: Date Not Accepted

**Check**:
- Timeline Journaling mode limits to 30 days
- Use CV Backfill or Manual Entry for older dates
- Date format must be YYYY-MM-DD or relative

### Issue: Application Crashes

**Steps**:
1. Check terminal size is adequate (80+ columns)
2. Try `./kariya-cli --help` to verify installation
3. Check Go version: `go version` (1.24.0+)
4. Review error message for hints

## Advanced Usage

### Importing Events from CSV

```bash
# Prepare CSV with: description, date, company, tags
# Then use Manual Entry to add events (currently manual process)
```

### Exporting Events

```bash
# Events can be exported to JSON/CSV from list view
# Use filter + search to select specific events
# Then export for CV building or backup
```

### Database Management

```bash
# Backup database
cp kariya.db kariya.db.backup

# Use custom database
./kariya-cli --db /path/to/my/events.db

# Multiple databases
./kariya-cli --db ~/career/work_events.db
./kariya-cli --db ~/career/volunteer_events.db
```

## Keyboard Reference

For quick reference while in the app, press 'h' to open the help system which includes complete keyboard shortcuts.

## Tips & Tricks

1. **Tag Strategy**: Use overlapping tags (e.g., "technical" + "achievement" for good technical wins)
2. **Date Entry**: Use relative dates like "2 weeks ago" instead of calculating exact dates
3. **Search First**: When looking for events, try search before filtering
4. **Export for CV**: Filter by "achievement" tag, then export for CV building
5. **Regular Backup**: Periodically backup your database file
6. **Company Consistency**: Use exact same company names for better filtering

## File Locations

- **Default Database**: Current directory (memory-based by default)
- **Custom Database**: Specify with `--db` flag
- **Config**: Command-line flags only (currently)

## Performance

- **Startup Time**: < 1 second
- **Event Listing**: Instant for 1000+ events
- **Search**: Real-time with minimal lag
- **Filter**: Immediate response
- **Database**: SQLite supports 100,000+ events

## Limitations

Current version:
- No event editing after creation (reimport with new capture)
- No bulk operations
- Manual entry only (no CSV import)
- Terminal UI only (no web interface yet)

Future versions will address these limitations.

## Getting Help

1. **In-App Help**: Press 'h' while running
2. **This Guide**: `docs/CLI_GUIDE.md`
3. **Help Flag**: `./kariya-cli --help`
4. **Version**: `./kariya-cli --version`

## Support

For issues or suggestions:
1. Check the troubleshooting section above
2. Review the help system in the app
3. Check test coverage in `internal/cli/models/*_test.go`
4. Review code in `internal/cli/`

## License

KaRiya is part of the career journaling project.

