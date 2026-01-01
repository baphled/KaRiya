# KaRiya Project Handover Document

## Project Overview
KaRiya is a Career Journal CLI tool designed to help professionals track, manage, and reflect on their career events and progression. It provides an interactive terminal interface for capturing, organizing, and analyzing career milestones. The application supports multiple workflows including timeline journaling, CV backfill, and manual event entry, with advanced features for burst fact extraction and metadata management.

## Technical Specifications

### Technology Stack
- **Language**: Go (1.24+)
- **CLI Framework**: BubbleTea (Charmbracelet) - Modern TUI library
- **Testing**: Ginkgo v2 - BDD testing framework
- **Database**: SQLite (modernc.org/sqlite) - Embedded relational database
- **Version Control**: Semantic Release - Automated versioning
- **Commit Management**: Conventional Commits, Commitlint - Enforced commit message standards
- **Module System**: Go modules with dependency management

### Key Dependencies
- **github.com/charmbracelet/bubbles** - Pre-built BubbleTea components
- **github.com/charmbracelet/bubbletea** - Core TUI framework
- **github.com/charmbracelet/lipgloss** - Terminal styling and layout
- **github.com/onsi/ginkgo/v2** - BDD testing framework
- **github.com/onsi/gomega** - Assertion and matching library
- **modernc.org/sqlite** - Pure Go SQLite implementation
- **github.com/google/uuid** - UUID generation

## Development Workflow

### Prerequisites
- Go 1.24 or higher
- Node.js 18+ with npm
- Ginkgo v2 for testing
- Make (for task automation)
- Git with hooks configured

### Setup
1. Clone the repository
2. Run `go mod tidy` to install Go dependencies
3. Run `npm install` for Node.js dependencies
4. Run `make install-git-hooks` to setup git hooks
5. Run `make test` to verify everything works

### Key Make Commands
- `make test`: Run all tests (currently 176/178 passing)
- `make coverage`: Generate code coverage report
- `make install-git-hooks`: Setup git hooks for commit message validation
- `make check-ai-attribution`: Verify AI commit attribution
- `make list-ai-commits`: List all AI-assisted commits
- `make audit-ai-commits`: Full audit of AI commits with statistics

## Project Structure

### Main Directories

#### `cmd/cli/`
**CLI Entry Point**
- `main.go` - Application bootstrap with argument parsing
- Handles multiple modes: timeline, backfill, manual
- Supports various CLI flags: --import, --detect-bursts, --extract-facts, --list, etc.

#### `internal/`
Core application logic organized by layer:

**`internal/domain/career/`** - Domain Models
- `event.go` - CareerEvent model with validation
  - Tags: project, achievement, leadership, technical, consulting, research, product, mentoring
  - Categories: technical, leadership, product, consulting, research, mentoring
  - Comprehensive validation for text, date, tags, and categories
- `burst.go` - Burst model (grouping of related events)
- `fact.go` - Fact model (extracted insights from events)
- All models include validation methods and JSON serialization

**`internal/repository/career/`** - Data Persistence Layer (170+ files)
- `repository.go` - Repository interface definitions
- `sqlite_repository.go` - SQLite implementation for CareerEvent storage
- `sqlite_burst_repository.go` - Burst persistence layer
- `sqlite_fact_repository.go` - Fact persistence layer
- Uses prepared statements and proper transaction handling
- Comprehensive test coverage with integration tests

**`internal/service/career/`** - Business Logic Layer
- `service.go` - Main career event service
- `classification/` - Event classification and metadata extraction
- `burst_fact/` - Burst and fact extraction algorithms
- Event filtering, searching, and manipulation logic
- Supports three capture modes: timeline journaling, CV backfill, manual entry

**`internal/logger/`** - Logging Infrastructure
- Structured logging with context support
- JSON-formatted output for production

**`internal/cli/`** - CLI UI Layer (Core TUI Implementation)

##### `internal/cli/app/` - Application State Management
- `app.go` - Main Model struct coordinating entire application
- Manages 20+ screens with centralized state
- Screen types: HomeScreen, MainMenuScreen, CaptureScreen, ListScreen, ViewScreen, etc.
- Handles navigation, screen transitions, and message routing
- Integration with service and repository layers
- 25+ test files covering all application flows

##### `internal/cli/models/` - BubbleTea Screen Models (29 models, 23K+ lines)
Core UI models implementing BubbleTea's Model interface:

**List Models:**
- `list.go` - Event list with pagination, filtering, sorting
- `burst_list.go` - Burst listing and management
- `fact_list.go` - Fact listing and discovery

**Detail & Editor Models:**
- `details.go` - Generic detail view component
- `burst_details.go` - Burst detail display
- `burst_editor.go` - Burst creation/editing
- `fact_details.go` - Fact detail view
- `fact_editor.go` - Fact editing
- `form.go` - Generic form component with field validation
- `fact_editor.go` - Event metadata editing with multi-field forms

**Specialized Models:**
- `action_menu.go` - Context-sensitive action menus
- `confirmation_dialog.go` - User confirmation dialogs
- `bulk_operations.go` - Batch operations on multiple events
- `burst_suggestion.go` - AI-assisted burst recommendations
- `burst_card.go` - Card-based burst display
- `fact_card.go` - Card-based fact display
- `metadata_review.go` - Metadata verification screen
- `metadata_editor.go` - Metadata field editing
- `view_event.go` - Event detail view
- `view_event_with_facts.go` - Event view with associated facts
- `tutorial.go` - Interactive tutorial/help system

**Base & Support Models:**
- `base.go` - BaseStandardModel providing common functionality (error handling, help footer, etc.)
- `messages.go` - Message types for TUI communication
- `errors.go` - Error handling model

All models follow these patterns:
- Implement BubbleTea's Model interface (Update, View)
- Inherit from BaseStandardModel for consistency
- Include help footer integration
- Standardized error display
- Focus indicator consistency
- Comprehensive test coverage

##### `internal/cli/components/` - Reusable UI Components (36 files, 7.8K+ lines)
Smart, composable components for TUI rendering:

**Layout Components:**
- `form_container.go` - Intelligent form layout (single-column, two-column, responsive)
- `table_list_container.go` - Table-based list rendering with pagination
- `modal_container.go` - Modal dialog container
- `section_container.go` - Section-based content organization
- `header_model.go` - Header component with breadcrumbs

**UI Components:**
- `help_footer_model.go` - Standardized help footer showing keyboard shortcuts
- `tag_selector.go` - Tag selection component for event categorization
- `input.go` - Text input field with validation
- `list.go` - Generic list component

**Styling Infrastructure:**
- All components use the centralized styling system
- Support for focus indicators, error states, and interactive feedback

##### `internal/cli/navigation/` - Navigation System
- `constants.go` - Centralized navigation constants and screen IDs
- `help_system.go` - Context-aware help content
- `key_handler.go` - Centralized keyboard event handling

##### `internal/cli/styles/` - Styling & Theme System
- Centralized color scheme and styling
- Consistent UI appearance across all screens
- Supports dark mode and terminal compatibility

##### `internal/cli/validation/` - Input Validation
- Field validation logic
- Event text validation
- Date validation
- Tag validation

##### `internal/cli/service/` - CLI Service Layer
- `event_service.go` - High-level event management for UI
- Bridge between domain service and CLI models

##### `internal/cli/importer/` - Data Import functionality
- `importer.go` - CSV import logic
- Burst fact import capabilities

##### `internal/cli/workflow/` - Workflow Orchestration
- Workflow state machines for complex multi-step processes

### Key Configuration Files
- `go.mod` - Go module dependencies
- `go.sum` - Go module checksums
- `package.json` - Node.js dependencies (commitlint, semantic-release)
- `.commitlintrc.json` - Conventional commit validation rules
- `.releaserc.json` - Semantic Release configuration
- `Makefile` - Development task automation
- `ginkgo.yml` - Ginkgo test configuration
- `.gitmessage` - Commit message template
- `.gitignore` - Git exclusion patterns

## Architecture & Design Patterns

### Layered Architecture
The application follows a clean, layered architecture:

```
┌─────────────────────────────────────────┐
│   CLI Layer (TUI Models & Components)   │
│   - Presentation logic                  │
│   - User interaction handling           │
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│   Service Layer                         │
│   - Business logic                      │
│   - Event processing                    │
│   - Burst/Fact extraction               │
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│   Repository Layer                      │
│   - Data persistence                    │
│   - SQLite operations                   │
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│   Domain Layer                          │
│   - Business entities                   │
│   - Validation rules                    │
└─────────────────────────────────────────┘
```

### BubbleTea Application Pattern
The application uses the BubbleTea TUI framework with:
- **Model** - Application state and logic
- **Update(msg Msg) (Model, tea.Cmd)** - Handle messages and state changes
- **View() string** - Render current screen
- **Cmd** - Commands that produce messages (async operations)

### Component-Based UI
- Reusable components (FormContainer, TableListContainer, etc.)
- Consistent styling through styles package
- Smart layout adaptation based on terminal dimensions
- Help footer integration for discoverability

### Service-Oriented Design
- Repository pattern for data access
- Service layer for business logic
- Clear separation of concerns
- Dependency injection for testability

## Testing Strategy

### Test Coverage
- **Current Status**: 176 passing tests, 2 failing tests (target: 100% pass rate)
- **Test Framework**: Ginkgo v2 (BDD style testing)
- **Test Organization**: Parallel test suites across multiple packages

### Test Suites by Package
1. **Domain Tests** (`internal/domain/career/`)
   - Event validation tests
   - Burst validation tests
   - Fact extraction tests

2. **Repository Tests** (`internal/repository/career/`)
   - SQLite integration tests
   - Transaction handling tests
   - Query correctness tests

3. **Service Tests** (`internal/service/career/`)
   - Business logic validation
   - Event processing
   - Burst/fact extraction algorithms

4. **CLI Tests** (`internal/cli/`)
   - Model rendering tests
   - Navigation tests
   - User interaction simulation
   - Component integration tests
   - App-level integration tests

### Running Tests
```bash
# Run all tests
make test

# Run specific package tests
go test ./internal/cli/models -v

# Run with coverage
make coverage

# View coverage in browser
make coverage && open coverage.html
```

### Known Test Issues
1. **failing_test_1**: `should render list screen view` in `app_test.go`
   - Related to list screen rendering with updated components
   - Likely needs assertion update after model refactoring
2. **failing_test_2**: `should display Home > Events breadcrumbs in list header` in `breadcrumb_display_test.go`
   - Breadcrumb rendering issue in list screen header
   - May need update to match new TableListContainer implementation

## Development Guidelines

### Code Organization Principles
1. **Single Responsibility** - Each component/model has one reason to change
2. **Interface-Driven** - Use interfaces for dependencies (testability)
3. **Composition Over Inheritance** - Build with small, focused components
4. **Explicit Error Handling** - No silent failures
5. **Comprehensive Logging** - Structured logging for debugging

### Naming Conventions
- Models: `*Model` suffix (e.g., `ListModel`, `FormModel`)
- Components: Descriptive names (e.g., `FormContainer`, `TableListContainer`)
- Services: `*Service` suffix
- Repositories: `*Repository` suffix
- Interface: descriptive names without suffixes

### Commit Message Convention
Use conventional commits format:
```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types**: feat, fix, docs, style, refactor, test, chore, build
**Scopes**: app, models, components, styles, navigation, validation, service, etc.

### AI Commit Attribution (IMPORTANT)
All AI-generated code must include attribution in commit message:
```
AI-Generated-By: <Assistant Name> (<Model Version>)
Reviewed-By: <Your Name>
```

Example:
```
feat(components): add form_container component

Implements responsive form layout system with intelligent field arrangement.

AI-Generated-By: Claude (Claude 3.5 Sonnet)
Reviewed-By: John Doe
```

### Code Style
- Follow `gofmt` formatting
- Use `golangci-lint` for linting recommendations
- Keep functions focused and small
- Use meaningful variable names
- Add comments for non-obvious logic

## Testing & Quality Assurance

### Pre-Commit Checks
- Commit message format validation (commitlint)
- AI attribution verification
- Conventional commits enforcement

### Running Specific Tests
```bash
# Run tests for a specific package
go test ./internal/cli/models -v

# Run a specific test
go test ./internal/cli/models -run TestListModel -v

# Run with coverage threshold
go test -cover ./...
```

### Coverage Goals
- Target: > 80% overall coverage
- Current status: Actively improving with each feature
- Use `make coverage` to generate detailed reports

## Deployment & Release

### Automated Release Process
- Triggered by semantic release on version tags
- Automatic CHANGELOG generation
- GitHub Actions CI/CD pipeline
- Binary builds for multiple platforms

### Version Management
- Uses semantic versioning (MAJOR.MINOR.PATCH)
- Automated via `semantic-release`
- Triggered by conventional commits

### Build & Binary
- Compiled CLI binary in `cli/` directory
- Built with `go build`
- Current version: 0.1.0

## Database Schema

### Tables

#### `career_events`
```
- id (TEXT PRIMARY KEY) - UUID
- text (TEXT) - Event description
- date (DATETIME) - Event date
- company (TEXT) - Associated company
- project (TEXT) - Associated project
- tags (TEXT JSON) - Array of tags
- categories (TEXT JSON) - Array of categories
- created_at (DATETIME) - Creation timestamp
- updated_at (DATETIME) - Last update timestamp
```

#### `bursts`
```
- id (TEXT PRIMARY KEY) - UUID
- name (TEXT) - Burst name
- description (TEXT) - Burst description
- event_ids (TEXT JSON) - Array of related event IDs
- competency_focus (TEXT) - Primary competency
- created_at (DATETIME) - Creation timestamp
- updated_at (DATETIME) - Last update timestamp
```

#### `facts`
```
- id (TEXT PRIMARY KEY) - UUID
- burst_id (TEXT) - Associated burst
- content (TEXT) - Fact content
- confidence (REAL) - Confidence score
- created_at (DATETIME) - Creation timestamp
```

## Recent Work Summary

### Latest Session: Component & Model Refactoring
Comprehensive refactoring of UI components and models for improved consistency and maintainability.

**Key Commits:**
1. **Form and Container Components** (4d3fee1)
   - FormContainer: Responsive form layout system
   - TableListContainer: Table-based list rendering
   - Smart adaptation to terminal size

2. **Detail and Editor Models** (60d57d1)
   - BurstDetails, BurstEditor models
   - FactDetails, FactSearch models
   - Improved data display patterns

3. **Navigation System** (5450b1f)
   - Centralized navigation constants
   - Enhanced help system integration
   - Consistent screen identification

4. **Model Restructuring** (9f1dfa8)
   - Refactored List model with better navigation
   - Updated FactList and BurstList
   - Enhanced Form model with validation
   - Improved test coverage

5. **Menu & Messages** (ff2db72)
   - Enhanced Menu model
   - Expanded Messages for new types
   - Refined BurstSuggestion
   - Updated ViewEventWithFacts

6. **App Integration** (bc8e6cb)
   - Integrated new models and components
   - Updated app navigation flows
   - Enhanced service layer
   - Improved e2e tests

### Previous Sessions
- **View Patterns Guide** - Applied consistent view patterns across UI
- **Error Display Standardization** - Unified error handling across 8+ models
- **Help Footer Integration** - Added keyboard shortcut hints to major screens
- **Focus Indicator Consistency** - Standardized focus state visualization
- **List Navigation Standardization** - Unified pagination and list controls
- **TUI Standardization** - Applied consistent design patterns throughout

## Important Code Patterns & Best Practices

### Model Pattern (BubbleTea)
```go
type MyModel struct {
    *BaseStandardModel  // Inherit common functionality
    // specific fields
}

// Implement Model interface
func (m MyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle messages and state changes
}

func (m MyModel) View() string {
    // Render current state
}
```

### Service Layer Pattern
```go
type Service struct {
    repo repo.Repository
    logger *logger.Logger
    // other dependencies
}

// Business logic methods with validation
func (s *Service) ProcessEvent(ctx context.Context, event *Event) error {
    if err := event.Validate(); err != nil {
        return err
    }
    return s.repo.Save(ctx, event)
}
```

### Component Pattern
```go
type MyComponent struct {
    // Component state
}

// Render method for display logic
func (c *MyComponent) Render() string {
    // Return rendered output
}
```

## Known Issues & TODO Items

### Failing Tests (Action Required)
1. **List Screen Rendering** (`internal/cli/app/app_test.go:128`)
   - Test: `should render list screen view`
   - Status: Failing - likely needs assertion update after TableListContainer integration
   - Priority: High - affects core list functionality

2. **Breadcrumb Display** (`internal/cli/app/breadcrumb_display_test.go:45`)
   - Test: `should display Home > Events breadcrumbs in list header`
   - Status: Failing - breadcrumb rendering in list header
   - Priority: High - affects navigation UX

### Code TODO Items
1. **View Event with Facts** (`internal/cli/models/view_event_with_facts.go`)
   - TODO: Get facts from bursts containing this event
   - Impact: Complete fact associations in event details
   - Priority: Medium

2. **Test Debugging**
   - DEBUG marker in `internal/cli/models/form_test.go`
   - Should verify repository works independently

## Migration & Upgrade Notes

### Database Migrations
- SQLite schema managed in repository package
- Use prepared migrations for version upgrades
- Test migrations with integration tests

### Breaking Changes
- None documented in recent commits
- Backward compatibility maintained for event data

## Common Development Tasks

### Adding a New Screen/Model
1. Create model file in `internal/cli/models/`
2. Implement Model interface (Update, View)
3. Inherit from BaseStandardModel
4. Add to app.go Screen enum
5. Implement navigation in app.go
6. Add tests in corresponding `*_test.go` file
7. Follow existing model patterns

### Adding a New Component
1. Create component file in `internal/cli/components/`
2. Implement Render() method
3. Add configuration options in constructor
4. Test with component tests
5. Use in models as needed

### Adding Tests
1. Use Ginkgo's `Describe`, `Context`, `It` blocks
2. Follow BDD style: "should do X when Y"
3. Use GinkGo helper for setup/teardown
4. Test both happy path and error cases

### Debugging Tips
1. Use structured logging: `logger.Info("message", "key", value)`
2. Print TUI state with `.View()` output
3. Use breakpoints with Delve debugger
4. Check breadcrumb display for navigation issues
5. Verify component dimensions with width/height checks

## Resources & Documentation

### Key Documentation Files
- `docs/guides/VIEW_PATTERNS_GUIDE.md` - UI pattern specifications
- `docs/guides/ERROR_HANDLING_GUIDE.md` - Error handling patterns
- `docs/guides/FOCUS_INDICATOR_GUIDE.md` - Focus state visualization
- `docs/guides/STYLE_USAGE_GUIDE.md` - Styling best practices
- `docs/guides/LIST_MODEL_RENDERING_SPECIFICATION.md` - List component specs
- `docs/TUI_DEVELOPER_GUIDE.md` - TUI development principles
- `docs/TUI_STANDARDS.md` - TUI design standards
- `docs/KEYBOARD_REFERENCE.md` - Keyboard shortcuts
- `docs/CLI_GUIDE.md` - User-facing CLI documentation
- `docs/TROUBLESHOOTING.md` - Common issues and solutions

### External Resources
- [BubbleTea Documentation](https://github.com/charmbracelet/bubbletea)
- [Lipgloss Documentation](https://github.com/charmbracelet/lipgloss)
- [Ginkgo Testing Framework](https://onsi.github.io/ginkgo/)
- [Go Best Practices](https://golang.org/doc/effective_go)

## Troubleshooting Common Issues

### Tests Failing
1. Run `make test` to see full error messages
2. Check if database file needs cleanup: `rm -f *.db`
3. Verify all dependencies installed: `go mod tidy`
4. Check for race conditions: `go test -race ./...`

### Database Issues
1. Verify SQLite is properly initialized
2. Check file permissions on database
3. Ensure migrations have run
4. Review repository transaction handling

### TUI Rendering Issues
1. Verify terminal size constraints
2. Check component width/height calculations
3. Review lipgloss style application
4. Test with different terminal emulators

### Build Issues
1. Run `go mod tidy && go mod verify`
2. Check Go version with `go version`
3. Clear build cache: `go clean -cache`
4. Rebuild binary: `go build ./cmd/cli`

## Next Steps for New Developers

1. **Read Core Documentation**
   - Start with README.md
   - Review PROJECT_HANDOVER_DOCUMENT.md
   - Check TUI_DEVELOPER_GUIDE.md

2. **Understand Architecture**
   - Trace a complete event capture flow
   - Study app.go state management
   - Review BubbleTea model pattern

3. **Fix Known Issues**
   - Start with the 2 failing tests
   - Review test assertions
   - Verify component integration

4. **Run Tests & Build**
   - Run `make test` to establish baseline
   - Fix failing tests step by step
   - Run `make coverage` for coverage report

5. **Make Small Changes**
   - Start with bug fixes
   - Then add small features
   - Follow established patterns

6. **Set Up Development Environment**
   - Install all prerequisites
   - Run `make install-git-hooks`
   - Configure editor for Go development

## Summary

KaRiya is a well-structured, tested Go CLI application with clear layering and comprehensive testing. The recent refactoring focused on component reusability and UI consistency. The main challenges are:

1. **2 Failing Tests** - Need investigation and fixes
2. **Complex State Management** - App.go manages 20+ screens
3. **TUI Complexity** - BubbleTea learning curve for new developers
4. **Fact Extraction** - TODO item for burst fact associations

The codebase is in good shape with:
- ✅ Clear architecture and layering
- ✅ Comprehensive test coverage
- ✅ Well-organized components and models
- ✅ Consistent code patterns
- ✅ Good documentation
- ⚠️ 2 test failures to fix
- ⚠️ Some incomplete features (fact associations)

This handover provides the foundation for any new developer to understand and contribute to the project effectively.

