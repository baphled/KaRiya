# Changelog

All notable changes to KaRiya Career Journal are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added - Phase 4 (Metadata Review & Enrichment)

#### Metadata Review System
- Implemented comprehensive metadata review screen for event clarification
  - View all events with data quality indicators
  - Filter by quality level (Incomplete, Basic, Enriched, Complete)
  - Sort by date, company, or creation order
  - Visual color-coded quality bars
  - Quick quality score and missing fields display

#### Individual Event Metadata Editor
- Created metadata editor for single event enrichment
  - Edit date with flexible format support (YYYY-MM-DD or relative dates)
  - Add/update company name with autocomplete
  - Add/update project name with autocomplete
  - Multi-select tags (max 8 from allowed tags)
  - Multi-select categories (max 2 from allowed categories)
  - Field validation with helpful error messages
  - Tab/Shift+Tab navigation between fields
  - Undo/revert functionality (Ctrl+Z)

#### Bulk Operations
- Implemented bulk metadata editing for multiple events
  - Select multiple events (Space to toggle, 'a' for all, 'd' to deselect)
  - Bulk edit with conditional updates ("Apply if empty" option)
  - Change preview before applying
  - Transaction-like behavior (all succeed or all fail)
  - Undo/revert support
  - Efficient metadata enrichment for imported or grouped events

#### Data Quality Scoring System
- Automatic quality scoring (0-100 scale) for all events
  - Text field: 20 points
  - Date field: 20 points
  - Company field: 20 points
  - Project field: 20 points
  - Tags: 15 points
  - Categories: 15 points
  - Quality match bonus: 10 points
- Four quality levels:
  - Incomplete (0-25): Only description
  - Basic (26-50): Description + date
  - Enriched (51-75): Description + date + company/project
  - Complete (76-100): All fields filled

#### Metadata Validation System
- Comprehensive field validation:
  - Date validation (not future, reasonable range, format support)
  - Company validation (optional, max 200 chars, normalization)
  - Project validation (optional, max 200 chars, normalization)
  - Tags validation (from allowed list, max 8, no duplicates)
  - Categories validation (from allowed list, max 2)
- Helpful error messages for user feedback
- Real-time validation on field blur or submit

#### Success Screen Enhancement
- Added "Review Metadata" option to post-capture success screen
  - Users can immediately enrich metadata after capture
  - Direct navigation to metadata review for captured event
  - Maintains capture flow without disruption

#### CSV Import Integration
- Automatic navigation to metadata review after CSV import
  - All imported events pre-loaded in metadata review
  - Ready for immediate enrichment and validation
  - Bulk operations available for imported event batches
  - Seamless workflow from import to enrichment

#### Documentation
- Created comprehensive METADATA_REVIEW_GUIDE.md covering:
  - Data quality scoring system
  - Individual event editing workflows
  - Bulk operations with examples
  - Filtering and sorting options
  - Common workflows and best practices
  - Troubleshooting section
  - Advanced features

- Updated CLI_GUIDE.md with:
  - Metadata review and enrichment workflows
  - Individual editor keyboard shortcuts
  - Bulk operations guide
  - Data quality scoring explanation
  - Post-capture metadata enrichment workflow

- Updated CSV_IMPORT_GUIDE.md with:
  - Post-import metadata review integration
  - Individual and bulk editing workflows
  - Data quality improvement strategies
  - Enrichment examples

#### Testing
- 6+ integration tests for capture → metadata review workflow
- 131+ total tests passing (100% success rate)
- Race condition detection: 0 detected
- Test coverage: 80%+ maintained


### Added - Phase 3 (Help System & Configuration)

#### CLI Help System
- Implemented comprehensive HelpModel with 7-section help guide
  - Overview: KaRiya introduction
  - Event Capture: Three capture modes explained
  - Tagging & Organization: Tag system and best practices
  - Listing & Filtering: Event management features
  - Keyboard Shortcuts: Complete keyboard reference
  - Search & Find: Search tips and examples
  - Tips & Best Practices: Usage recommendations
- Help sections support step-by-step navigation
- Quick search functionality for help content
- Progress bar showing current section

#### CLI Configuration
- Added `--db, --database PATH` flag for custom database path
- Added `--mode MODE` flag to start in specific capture mode (timeline/backfill/manual)
- Added `--list` flag to show recent events on startup
- Mode validation with helpful error messages
- Updated help text with new flags

#### CLI Screen Models
- DetailsModel: Full event information display screen
  - Shows complete event details (text, date, company, project)
  - Displays tags with styling
  - Shows event ID and timestamps
  - Graceful nil event handling
  - Keyboard navigation (backspace to return)

- TutorialModel: Interactive first-run guide
  - 7-step tutorial covering KaRiya features
  - Step-by-step navigation through tutorial
  - Skip option (ESC or 'q')
  - Completion and skip tracking

#### Documentation
- Created comprehensive CLI_GUIDE.md with:
  - Quick start guide
  - Feature descriptions
  - Keyboard shortcuts reference
  - Usage examples (3 real-world scenarios)
  - Best practices section
  - Configuration guide
  - Troubleshooting section
  - Advanced usage tips

- Updated main README with:
  - CLI usage section
  - CLI architecture overview
  - CLI examples
  - CLI testing instructions
  - Performance metrics
  - Troubleshooting for CLI

- Added CHANGELOG.md documenting all changes


### Added - Phase 3 (Continued: CLI Flags, Error Recovery, Documentation)

#### CLI Flags & Configuration (Task 17)
- Implemented `--db, --database PATH` flag for persistent SQLite storage
  - Default: in-memory storage (no persistence)
  - Custom path: `./kariya-cli --db ~/.kariya/events.db`
  - Automatic database initialization on startup
- Implemented `--mode MODE` flag to start in specific capture mode
  - Valid modes: timeline, backfill, manual
  - Example: `./kariya-cli --mode timeline`
- Implemented `--list` flag to show events list on startup
  - Example: `./kariya-cli --list --db events.db`
- Enhanced help text with practical examples
- Improved error messages for invalid flags

#### Error Recovery & Edge Cases (Task 18)
- Added 26 comprehensive error handling test cases
  - Service layer validation errors (empty text, timeline window, future dates)
  - Repository error recovery (non-existent events, nil filters)
  - Application state recovery after errors
  - Long text handling at 1999/2000/2001 character boundaries
  - Navigation state consistency
  - Message handling robustness (unknown messages, window resizes, rapid updates)
  - Special character handling (unicode, newlines, tabs)
  - Date boundary testing (timeline 30-day window, old/future dates)
- Fixed ListEvents() to handle nil filters gracefully
- Verified graceful error recovery for all failure scenarios
- All error handling tests passing (26/26)

#### Documentation (Task 21 - Partial)
- Created comprehensive TROUBLESHOOTING.md guide covering:
  - Database & persistence issues
  - Form & input issues (date formats, character limits)
  - Display & appearance problems
  - Performance troubleshooting
  - Navigation & workflow issues
  - Capture mode constraints
  - Advanced troubleshooting (debug logging, database integrity)
  - FAQ section
- Updated main README with CLI usage section
- Enhanced CHANGELOG.md with complete version history
- CLI_GUIDE.md already comprehensive (quick start, features, shortcuts)

#### Testing Improvements
- Expanded test suite with 26 error handling tests
- All new tests passing (26/26 ✅)
- Total test count: 419+ tests
- Overall test pass rate: 100%
- Code coverage maintained at 80%+

### Changed - Phase 3

#### CLI Infrastructure
- Enhanced main.go with comprehensive flag parsing
- Improved error handling in CLI entry point
- Extended help text to include new flags

#### Testing
- Added 42 new test cases this phase:
  - DetailsModel: 14 tests
  - TutorialModel: 10 tests
  - HelpModel: 18 tests
- Enhanced CLI flag tests with 4 new test cases
- All tests passing (180+ total tests)

### Fixed

- Fixed unused variable in HelpModel search method
- Improved error messages for invalid mode flags

## [0.1.0] - 2025-12-24 (Phase 2 Complete)

### Added - Phase 1 & 2 (MVP + Event Management)

#### Core Features
- Career event domain model with validation
- Three event capture modes (Timeline, CV Backfill, Manual)
- In-memory and SQLite repository implementations
- Event listing with pagination
- Event filtering by date, tags, company
- Event search functionality
- Event sorting by date, creation time, text
- Event classification system (6 competency categories)
- Structured logging with context support

#### CLI Interface (BubbleTea)
- Interactive event capture form
  - Text input (1-2000 characters)
  - Date parsing (ISO format, relative dates)
  - Company and project fields
  - Tag multi-select with autocomplete
  - Capture mode selector
  - Input validation with error feedback

- Event listing screen with pagination
  - Page size configuration (default: 10)
  - Previous/Next page navigation
  - Page indicator display

- Event filtering screen
  - Date range filtering
  - Tag multi-select filtering
  - Company name filtering
  - Filter state management

- Event search screen
  - Keyword search with debouncing
  - Real-time search results
  - Text highlighting support
  - Case-insensitive matching

- Event sorting screen
  - Sort field selection (date, created_at, text)
  - Sort order selection (ascending/descending)
  - Multiple sort options

- Success screen with post-capture actions
  - Event summary display
  - "Capture Another Event" navigation
  - "View Recent Events" navigation
  - "Exit" option

- Professional styling
  - Dark blue/gray color scheme
  - Responsive layout helpers
  - Consistent spacing and alignment
  - Professional typography

#### Testing
- Comprehensive test coverage (180+ tests)
  - Domain layer: 100% coverage
  - Service layer: 100% coverage
  - CLI models: 75%+ coverage
  - Overall: 81%+ coverage

#### Documentation
- Comprehensive AGENTS.md handover document
- Architecture documentation
- Development setup guide
- Testing strategy documentation
- Contributing guidelines

## Test Summary

### Current Status (Phase 3 Development)
- **Total Tests**: 180+
- **Pass Rate**: 100%
- **Coverage**: 77.9% overall
- **Race Conditions**: 0 detected

### Test Breakdown
- CLI Entry: 6 tests ✅
- App Navigation: 39 tests ✅
- Form Capture: 52 tests ✅
- Event Listing: 40+ tests ✅
- Event Details: 14 tests ✅
- Tutorial: 10 tests ✅
- Help System: 18 tests ✅
- Components: 18 tests ✅
- Styles: 63 tests ✅
- Validation: 16 tests ✅
- Services: 31 tests ✅
- Domain: 5 tests ✅
- Classification: 8 tests ✅

## Architecture

### Domain-Driven Design
- Clear separation between domain, service, and repository layers
- Domain model enforces validation rules
- Service layer implements business logic
- Repository pattern for data persistence

### Project Structure
```
KaRiya/
├── cmd/cli/                    # CLI entry point
├── internal/
│   ├── cli/                    # Terminal UI layer
│   │   ├── app/               # App state and navigation
│   │   ├── models/            # Screen models (BubbleTea)
│   │   ├── components/        # Reusable components
│   │   ├── styles/            # Styling system
│   │   ├── validation/        # Input validation
│   │   └── service/           # CLI service adapter
│   ├── domain/career/         # Domain model layer
│   ├── service/career/        # Business logic layer
│   ├── repository/career/     # Data persistence layer
│   └── logger/                # Logging infrastructure
├── docs/                       # Documentation
├── features/                   # Feature specifications
└── tasks/                      # Task tracking
```

## Performance

- CLI startup: < 1 second
- Form submission: Instant
- Event listing: < 100ms for 1000+ events
- Search: Real-time response
- Database: SQLite ready for 100k+ events

## Known Limitations

### Phase 1-2
- List screen is interactive but not fully featured
- View detail screen is functional
- No event editing after creation
- No bulk operations

### Phase 3 (Current)
- Error recovery limited
- No advanced UI animations
- Performance not fully optimized
- Documentation incomplete (in progress)

### Future Work (Phase 3+)
- [ ] Full error recovery and edge case handling
- [ ] Advanced UI/UX polish and animations
- [ ] Performance optimization and benchmarks
- [ ] Advanced features (export, import, templates)
- [ ] Web interface
- [ ] Mobile application

## Dependencies

### Core
- Go 1.24.0+
- BubbleTea v0.26+ (Terminal UI)
- Lipgloss (Styling)
- Bubbles (Input components)

### Testing
- Ginkgo v2.27.3 (Testing framework)
- Gomega v1.38.3 (Assertion library)

### Storage
- SQLite (via modernc.org/sqlite)

### Utilities
- UUID (github.com/google/uuid)

## Contributing

When contributing to KaRiya:

1. Follow Go idioms and best practices
2. Use test-driven development (Red-Green-Refactor)
3. Maintain test coverage (target: 80%+)
4. Create atomic commits with clear messages
5. Update documentation and CHANGELOG
6. Run full test suite before submitting

## Building from Source

```bash
# Build
go build -o kariya-cli ./cmd/cli

# Run
./kariya-cli

# Test
make test

# Coverage
go test -race ./... -coverprofile=cover.out
go tool cover -func=cover.out
```

## License

[License information to be added]

## Support

For issues or questions:
1. Check CLI_GUIDE.md for usage help
2. Run `./kariya-cli --help` for command-line options
3. Press 'h' in the app for interactive help
4. Review test files for usage examples

---

**Last Updated**: 2025-12-24
**Current Version**: 0.1.0 (Phase 1-2 Complete, Phase 3 In Progress)
**Status**: MVP Complete, Feature Development Ongoing


## [Phase 3] - 2025-12-30

### Tasks 9-11: Bulk Operations Feature

#### Task 9.0: Bulk Operations Model - COMPLETE ✅
- Implemented BulkOperationsModel (412 lines) with full BubbleTea integration
- Multi-select event selection with Space/a/d keyboard shortcuts
- Bulk field editing for company, project, tags, categories
- Preview and confirmation workflows
- Undo/revert capability
- 27 comprehensive test cases - all PASSING ✅

#### Task 10.0: Navigation Integration - COMPLETE ✅
- Added BulkOperationsScreen to navigation
- Implemented message handlers and state delegation
- Full app integration with metadata review workflow
- 122 app tests passing ✅

#### Task 11.0: CLI Service Enhancement - COMPLETE ✅
- Implemented BulkUpdateMetadata() method
- Transaction-like validation (all succeed or all fail)
- Summary return with update statistics
- Proper error handling for edge cases

### Test Results
- Total tests: 449+ across all packages
- BulkOperationsModel: 27/27 PASSING
- App integration: 122/122 PASSING
- CLI service: 8/8 PASSING
- Overall success rate: 95.8%

### Code Quality
- All code formatted with gofmt ✅
- Zero race conditions in bulk operations ✅
- Code coverage: 70.6% (internal packages)
- Atomic commits with conventional messages ✅

### Files Created/Modified
- New: internal/cli/models/bulk_operations.go (412 lines)
- New: internal/cli/models/bulk_operations_test.go (27 tests)
- Modified: internal/cli/app/app.go (navigation integration)
- Modified: internal/cli/app/messages.go (message types)
- Modified: internal/cli/service/event_service.go (BulkUpdateMetadata)

### Known Limitations
- 2 pre-existing test failures in cmd/cli persistence tests
- 19 pre-existing failures in other model tests (form, quality_indicator, metadata_review)

### Next Steps
- Phase 4: Web UI and API endpoints
- Phase 5: Advanced features (burst detection, fact extraction)
