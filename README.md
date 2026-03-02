# KaRiya Career Event Capture

KaRiya is a career event capture terminal user interface built with Go and Bubble Tea. It helps you track achievements, manage career growth, and generate tailored CVs directly from your terminal.

## Prerequisites

- Go 1.25.4 or higher
- Ginkgo v2
- Make (optional, for task automation)
- Node.js 18+ and npm (for commitlint)

## Installation

1. Clone the repository
2. Run `go mod tidy` to install dependencies
3. Run `npm install` to install Node.js dependencies
4. Run `make install-git-hooks` to set up git hooks

## CLI Usage

KaRiya includes an interactive terminal user interface. Use `./kariya-cli --help` for full command information.

### Quick Start

```bash
# Build the CLI
go build -o kariya-cli ./cmd/cli

# Run with defaults
./kariya-cli

# Run with custom database
./kariya-cli --db ~/.kariya/events.db

# Start in specific capture mode
./kariya-cli --mode timeline

# Import CSV with automatic burst and fact detection
./kariya-cli --import events.csv

# View help
./kariya-cli --help
```

## Features

- **Event Capture**: Timeline Journaling, CV Backfill, and Manual Entry modes
- **Event Management**: List, filter, search, sort, and edit event details
- **Metadata Enrichment**: Review and validate event metadata with quality scoring
- **Bulk Operations**: Update metadata for multiple events with conditional logic
- **CSV Import**: Import events with automatic metadata review and processing
- **Intelligent Processing**: Automatic burst detection and fact extraction from events
- **CV Generation**: Transform events into role-specific, audience-tailored CVs
- **Interactive TUI**: 7-section help system and first-run tutorial for new users

## CV Generation

KaRiya transforms your career events into professional CVs tailored to specific roles (Principal, Staff, EM, Senior IC) and audiences (Hiring Manager, Recruiter, Peer).

### Quick Start

```bash
# Manage CV configurations
./kariya-cli --manage-cv

# Generate a CV from existing config
./kariya-cli --generate-cv "My Staff Engineer CV"

# List all CV configurations
./kariya-cli --list-cv-configs
```

Features include intelligent bullet ranking, automatic compression to respect role caps, and full source traceability for every bullet point.

## Keyboard Shortcuts

- `c` - Capture new event
- `i` - Import from CSV
- `l` - List events
- `m` - Open metadata review
- `u` - View burst suggestions
- `v` - Manage CV configurations
- `g` - Generate CV
- `h` - Help system
- `q` - Quit
- `tab` / `shift+tab` - Navigate form fields
- `up` / `down` - Move through lists
- `space` - Select or deselect items
- `a` - Select all / `d` - Deselect all
- `y` / `n` - Confirm or reject suggestions

## CLI Examples

### Example 1: Capture Timeline Event
1. Run: `./kariya-cli --mode timeline`
2. Press `c` to capture
3. Enter: "Led API redesign for performance improvement"
4. Date: "today" (or leave blank)
5. Company: "TechCorp"
6. Project: "API Modernisation"
7. Tags: technical, achievement, leadership
8. Submit

### Example 2: Backfill CV Event
1. Run: `./kariya-cli --mode backfill`
2. Press `c` to capture
3. Enter: "Architected microservices migration"
4. Date: "2023-06-15"
5. Company: "StartupXYZ"
6. Project: "System Architecture"
7. Tags: technical, leadership, achievement
8. Submit

## Performance

- **Startup**: Less than 1 second
- **Event Listing**: Under 100ms for 1,000 events
- **Search & Filtering**: Real-time and instant
- **Memory**: Efficiently handles 10,000+ events

## Troubleshooting

### Events disappear after restart
**Cause**: Using the default in-memory database.
**Solution**: Use the `--db` flag with a persistent SQLite path:
```bash
./kariya-cli --db ~/.kariya/events.db
```

### Form field navigation issues
**Cause**: Terminal window is too small.
**Solution**: Increase terminal width to at least 80 columns.

### Special characters not displaying
**Cause**: Terminal does not support UTF-8.
**Solution**: Ensure your terminal is set to UTF-8 encoding.

### Help system not showing
**Cause**: Terminal height is insufficient.
**Solution**: Increase terminal window height.

## Development

See [AGENTS.md](AGENTS.md) for development guidelines, testing workflow, architecture, and contribution rules.
