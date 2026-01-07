# KaRiya Career Event Capture

## Project Setup

### Prerequisites
- Go 1.24 or higher
- Ginkgo v2
- Make (optional, for task automation)
- Node.js 18+ and npm (for commitlint and CI/CD)

### Installation
1. Clone the repository
2. Run `go mod tidy` to install dependencies
3. Run `npm install` to install Node.js dependencies (commitlint, semantic-release)
4. Run `make install-git-hooks` to setup git hooks

## Running Tests

### Run All Tests
```bash
make test
```

### Generate Code Coverage
```bash
make coverage
```

### Clean Coverage Reports
```bash
make clean-coverage
```

## Development Workflow
1. Review tasks in project documentation
2. Follow Red-Green-Refactor methodology
3. Ensure high test coverage
4. Run linters before committing
5. **Follow AI commit attribution rules** (see below)

## AI Commit Attribution 🤖

**IMPORTANT**: All commits created with AI assistance MUST include proper attribution.

### Quick Setup

```bash
make install-git-hooks
```

### Required Format

For any AI-generated code, include in commit message:
```
AI-Generated-By: <Assistant Name> (<Model Version>)
Reviewed-By: <Your Name>
```

### Examples

```
AI-Generated-By: Avante (Claude 3.5 Sonnet)
AI-Generated-By: Claude (Claude 3.7 Sonnet)
AI-Generated-By: GitHub Copilot (GPT-4)
```

### Documentation

- **Comprehensive Guide**: [docs/rules/AI_COMMIT_ATTRIBUTION.md](docs/rules/AI_COMMIT_ATTRIBUTION.md)
- **Quick Setup**: [docs/setup/AI_COMMIT_SETUP.md](docs/setup/AI_COMMIT_SETUP.md)
- **Implementation Summary**: [docs/setup/AI_COMMIT_ATTRIBUTION_SETUP.md](docs/setup/AI_COMMIT_ATTRIBUTION_SETUP.md)

### Verification

```bash
make check-ai-attribution   # Check latest commit
make list-ai-commits        # List all AI commits
make audit-ai-commits       # Full audit with statistics
```

## CI/CD Pipeline 🚀

**IMPORTANT**: We use GitHub Actions for CI/CD with automated releases.

### Commit Message Format

All commits **MUST** follow conventional commits format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Examples**:
```
feat(service): add event filtering
fix(domain): prevent duplicate tags
docs(readme): update installation
```

### Valid Types

- `feat` - New feature (triggers minor release)
- `fix` - Bug fix (triggers patch release)
- `perf` - Performance improvement (triggers patch release)
- `refactor` - Code refactoring (triggers patch release)
- `docs` - Documentation (no release)
- `test` - Tests (no release)
- `chore` - Maintenance (no release)

### Automated Releases

Releases are **automatically created** when you push to `main`:

1. Commits are analyzed
2. Version is determined (major/minor/patch)
3. CHANGELOG is generated
4. GitHub release is created
5. Binaries are uploaded

### Documentation

- **Full Guide**: [docs/CI_CD_PIPELINE.md](docs/CI_CD_PIPELINE.md)
- **Quick Reference**: [docs/CI_CD_QUICK_REF.md](docs/CI_CD_QUICK_REF.md)

### Running CI Checks Locally

**Before pushing**, run all CI checks locally to ensure they pass:

```bash
# Install all CI tools (one-time setup)
make ci-install-tools

# Run ALL CI checks (mirrors GitHub Actions)
make ci-local
```

This runs:
- ✅ Commitlint validation
- ✅ AI attribution check  
- ✅ Code formatting (go fmt)
- ✅ Static analysis (go vet, staticcheck)
- ✅ Tests with race detector & coverage
- ✅ Multi-platform builds
- ✅ Security scanning (gosec)

**Individual checks**:
```bash
make fmt          # Format code
make vet          # Static analysis
make staticcheck  # Advanced linting
make gosec        # Security scan
make test         # Run tests
```

**Documentation**:
- [CI Checks Summary](docs/CI_CHECKS_SUMMARY.md) - Complete CI/local command mapping
- [CI Local Guide](docs/CI_LOCAL_GUIDE.md) - Detailed guide

### Validation

```bash
# Validate commit message
echo "feat(service): add feature" | npx commitlint

# Check what would be released
npx semantic-release --dry-run
```

## Test Coverage
- Coverage reports are generated in the `coverage` directory
- HTML report provides detailed code coverage visualization

## 📚 Documentation

Comprehensive documentation is organized in the `docs/` directory:

### Product Requirements
- **[docs/PRD_MASTER.md](docs/PRD_MASTER.md)** - Complete system architecture and requirements
- **[docs/PRD_USER_STORIES.md](docs/PRD_USER_STORIES.md)** - User stories and functional requirements
- **[docs/PRD_CLI.md](docs/PRD_CLI.md)** - CLI interface specifications

### Setup Guides
- **[docs/setup/](docs/setup/)** - Installation and configuration guides
  - AI commit attribution setup
  - Atomic commits setup
  - CI/CD pipeline setup
  - Master task prompt setup

### Development Rules & Guidelines
- **[docs/rules/](docs/rules/)** - Development standards (16 documents)
  - AI commit attribution rules
  - Atomic commit guidelines
  - Go coding standards
  - Task execution workflows
  - Quick reference guides

### Feature Specifications
- **[docs/features/](docs/features/)** - Detailed feature specs (8 features)
  - Career event capture
  - Burst & fact extraction
  - CV generation
  - Data model schema

### Tools & Editor Setup
- **[docs/tools/editor-setup/](docs/tools/editor-setup/)** - Editor configuration
  - Neovim/neotest setup guide

### Reports & Analysis
- **[docs/reports/](docs/reports/)** - Test reports and verification
  - Latest test status and coverage

### Other Documentation
- **[AGENTS.md](AGENTS.md)** - Comprehensive handover document
- **[DOCUMENTATION_REVIEW_SUMMARY.md](DOCUMENTATION_REVIEW_SUMMARY.md)** - Documentation reorganization plan
- **[docs/integration-test-strategy.md](docs/integration-test-strategy.md)** - Testing strategy

## Task Tracking
Track current development status in project documentation.


## CLI Usage

KaRiya includes an interactive terminal user interface built with BubbleTea. For detailed CLI usage, see [CLI_GUIDE.md](docs/CLI_GUIDE.md).

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

# View help
./kariya-cli --help
```

### Features

- **Event Capture**: Three modes (Timeline Journaling, CV Backfill, Manual Entry)
- **Event Browsing**: List, filter, search, sort, and view event details
- **Metadata Review & Enrichment**: Review, validate, and enrich event metadata with data quality scoring
- **Individual Event Editing**: Edit date, company, project, tags, categories for single events
- **Bulk Operations**: Edit metadata for multiple events at once with conditional updates
- **CSV Import**: Import events from CSV files with automatic metadata review
- **Data Quality Scoring**: Automatic scoring (0-100) to track metadata completeness
- **Tags**: Organize events with up to 8 tags per event
- **Categories**: Competency-based categorization with up to 2 categories per event
- **Burst Detection**: Automatically group related events into thematic bursts with confidence scoring
  - Triggered automatically after CSV import
  - Re-run on demand with `--detect-bursts` flag
  - View existing bursts with `--show-bursts` flag
  - Interactive review and confirmation in CLI
- **Fact Extraction**: Automatically extract grounded facts from events and bursts (no aspirational language)
  - Triggered automatically after CSV import
  - Re-run on demand with `--extract-facts` flag
  - View existing facts with `--show-facts` flag
  - Grounded statements only (no aspirational/speculative language)
- **Role Fit Classification**: Infer career level (Principal, EM, Staff Engineer, Senior IC)
- **Audience Relevance**: Determine who cares about each fact (Hiring Manager, Recruiter, Peer)
- **CV Generation**: Transform career events into role-specific, audience-tailored CVs
  - Role-specific generation (Principal, Staff, EM, Senior IC)
  - Audience-specific filtering (Hiring Manager, Recruiter, Peer)
  - Intelligent bullet ranking with confidence scoring
  - Full traceability of sources for each bullet
  - YAML-based configuration for reusable CV templates
  - Export to multiple formats (text, markdown, clipboard)
  - Compression logic respects role-specific bullet caps
  - No aspirational language, inferred metrics, or role inflation
- **Interactive Help**: 7-section help system with keyboard shortcuts
- **First-Run Tutorial**: Interactive guide for new users
- **CLI Flags**: Configuration via command-line arguments

### CV Generation

KaRiya can transform your career events into professional CVs tailored to specific roles and audiences.

**Quick Start**:
```bash
# Manage CV configurations
./kariya-cli --manage-cv

# Generate a CV from existing config
./kariya-cli --generate-cv "My Staff Engineer CV"

# List all CV configurations
./kariya-cli --list-cv-configs
```

**Key Features**:
- **Role-Specific CVs**: Generate different CVs for Principal, Staff, EM, or Senior IC roles
- **Audience Targeting**: Tailor CVs for Hiring Managers, Recruiters, or Peers
- **YAML Configurations**: Store reusable CV templates in `~/.kariya/cv_configs/`
- **Intelligent Ranking**: Bullets ranked by ownership, contribution, strategy, execution, and outcomes
- **Bullet Caps**: Automatic compression respects role-specific limits (Principal: 3-4, Staff: 4-5, etc.)
- **Source Traceability**: Every bullet traces back to source events and facts
- **Multiple Exports**: Export to text, markdown, or copy to clipboard
- **Quality Filters**: Excludes aspirational language, inferred metrics, and role inflation

For comprehensive CV generation documentation, see [CV_GENERATION_GUIDE.md](docs/guides/CV_GENERATION_GUIDE.md) and [CV_EXAMPLES.md](docs/guides/CV_EXAMPLES.md).

### Keyboard Shortcuts

- `c` - Capture new event
- `i` - Import from CSV
- `l` - List events
- `m` - Open metadata review
- `u` - View burst suggestions (from metadata review)
- `v` - Manage CV configurations
- `g` - Generate CV
- `h` - Help system
- `q` - Quit
- `tab`/`shift+tab` - Navigate form fields
- `up`/`down` - Move through lists
- `space` - Select/deselect (bulk operations, burst editing)
- `a` - Select all / `d` - Deselect all
- `y`/`n` - Confirm/reject burst suggestions

For complete keyboard reference, see [CLI_GUIDE.md](docs/CLI_GUIDE.md) or [BURST_FACT_EXTRACTION_GUIDE.md](docs/BURST_FACT_EXTRACTION_GUIDE.md) or press 'h' in the app.

## CLI Architecture

The CLI is organized into logical layers:

```
cmd/cli/main.go                 # Entry point and flag parsing
├── internal/cli/app/           # Application state and navigation
│   ├── app.go                  # Main app model
│   └── messages.go             # Message types
├── internal/cli/models/        # Screen models (BubbleTea)
│   ├── form.go                 # Event capture form
│   ├── list.go                 # Event listing
│   ├── details.go              # Event detail view
│   ├── tutorial.go             # First-run tutorial
│   ├── help.go                 # Help system
│   └── support models          # Filter, Search, Sort
├── internal/cli/components/    # Reusable components
│   ├── tag_selector.go         # Multi-select tags
│   ├── date_picker.go          # Date input
│   └── inputs.go               # Input fields
├── internal/cli/styles/        # Lipgloss styling
│   └── styles.go               # Color scheme and layout
├── internal/cli/validation/    # Input validation
│   └── validator.go            # Validation rules
└── internal/cli/service/       # Service adapter layer
    └── event_service.go        # Event service wrapper
```

## CLI Testing

```bash
# Run all CLI tests
make test

# Run specific test suite
ginkgo -v ./cmd/cli
ginkgo -v ./internal/cli/models

# Run with race detection
go test -race ./...

# Generate coverage
go test -race ./... -coverprofile=cover.out
go tool cover -func=cover.out
```

**Current Test Status**: 131+ tests, 100% passing (including metadata review and bulk operations)

### CLI Test Breakdown

- CLI Entry Point: 6 tests (version, help, flags)
- App Model: 39 tests (navigation, screen management)
- Form Model: 52 tests (capture, validation)
- List Model: 40+ tests (pagination, display)
- Details Model: 14 tests (event display)
- Tutorial Model: 10 tests (step navigation)
- Help Model: 18 tests (section management)
- Components: 18 tests (tag selector, inputs)
- Styles: 63 tests (styling, layout)
- Validation: 16 tests (input validation)

## CLI Examples

### Example 1: Capture Timeline Event

```
1. Run: ./kariya-cli --mode timeline
2. Press 'c' to capture
3. Enter: "Led API redesign for performance improvement"
4. Date: "today" (or leave blank)
5. Company: "TechCorp"
6. Project: "API Modernization"
7. Tags: technical, achievement, leadership
8. Submit
```

### Example 2: Backfill CV Event

```
1. Run: ./kariya-cli --mode backfill
2. Press 'c' to capture
3. Enter: "Architected microservices migration"
4. Date: "2023-06-15"
5. Company: "StartupXYZ"
6. Project: "System Architecture"
7. Tags: technical, leadership, achievement
8. Submit
```

### Example 3: Export with Custom Database

```bash
./kariya-cli --db ~/events/career.db

# Use UI to filter events
# (Example: filter by "achievement" tag)
# Then export to CV format
```

## CLI Configuration

### Environment Variables

Currently no environment variables. Use command-line flags instead.

### Database Configuration

```bash
# In-memory (default)
./kariya-cli

# SQLite database
./kariya-cli --db /path/to/events.db

# Note: SQLite support ready, use flag to enable
```

### Capture Mode

```bash
# Timeline mode (30-day window)
./kariya-cli --mode timeline

# CV Backfill (any past date)
./kariya-cli --mode backfill

# Manual (full flexibility)
./kariya-cli --mode manual
```

## Performance

- **Startup**: < 1 second
- **Event Listing**: < 100ms for 1000 events
- **Search**: Real-time
- **Filtering**: Instant
- **Memory**: Efficient for 10,000+ events

## Troubleshooting

### Events disappear after restart

**Cause**: Using default in-memory database

**Solution**: Use `--db` flag with SQLite:
```bash
./kariya-cli --db ~/.kariya/events.db
```

### Form field navigation issues

**Cause**: Terminal size too small

**Solution**: Increase terminal window width (minimum 80 columns)

### Special characters not displaying

**Cause**: Terminal doesn't support UTF-8

**Solution**: Ensure terminal is set to UTF-8 encoding

### Help system not showing

**Cause**: Terminal height too small

**Solution**: Maximize terminal window vertically

For more help, see [CLI_GUIDE.md](docs/CLI_GUIDE.md).

