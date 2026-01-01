# KaRiya Project Handover Document

## Project Overview
KaRiya is a Career Journal CLI tool designed to help professionals track, manage, and reflect on their career events and progression. It provides an interactive terminal interface for capturing, organizing, and analyzing career milestones.

## Technical Specifications

### Technology Stack
- **Language**: Go (1.24+)
- **CLI Framework**: BubbleTea (Charmbracelet)
- **Testing**: Ginkgo v2
- **Database**: SQLite (modernc.org/sqlite)
- **Version Control**: Semantic Release
- **Commit Management**: Conventional Commits, Commitlint

### Key Dependencies
- github.com/charmbracelet/bubbles
- github.com/charmbracelet/bubbletea
- github.com/onsi/ginkgo/v2
- modernc.org/sqlite

## Development Workflow

### Prerequisites
- Go 1.24 or higher
- Node.js 18+ with npm
- Ginkgo v2 for testing
- Make (for task automation)

### Setup
1. Clone the repository
2. Run `go mod tidy` to install Go dependencies
3. Run `npm install` for Node.js dependencies
4. Run `make install-git-hooks` to setup git hooks

### Key Make Commands
- `make test`: Run all tests
- `make coverage`: Generate code coverage report
- `make install-git-hooks`: Setup git hooks
- `make check-ai-attribution`: Verify AI commit attribution

## Project Structure

### Main Directories
- `cmd/cli/`: CLI entry point and main application
- `internal/cli/`: Core CLI implementation
  - `app/`: Application state and navigation
  - `models/`: Screen models (BubbleTea)
  - `components/`: Reusable UI components
  - `styles/`: Styling and layout
  - `validation/`: Input validation
  - `service/`: Service layer adapters

### Key Configuration Files
- `go.mod`: Go module dependencies
- `package.json`: Node.js dependencies and scripts
- `.commitlintrc.json`: Commit message validation
- `.releaserc.json`: Semantic Release configuration
- `Makefile`: Development task automation

## Development Guidelines

### Commit Message Convention
Use conventional commits format:
```
<type>(<scope>): <subject>

<body>

<footer>
```

### AI Commit Attribution
- All AI-generated code must include attribution
- Format:
  ```
  AI-Generated-By: <Assistant Name> (<Model Version>)
  Reviewed-By: <Your Name>
  ```

### Testing
- 131+ tests across various components
- 100% passing test suite
- Use Ginkgo for testing
- Aim for comprehensive test coverage

## Deployment & Release
- Automated releases via GitHub Actions
- Semantic versioning
- Automatic CHANGELOG generation
- Binaries uploaded with each release

## Troubleshooting
- Refer to README.md for detailed troubleshooting
- Common issues include:
  - Database persistence
  - Terminal compatibility
  - Input navigation

## Recent Work Summary

### Commits Made (Session: Component & Model Refactoring)
The following commits were created to organize related changes into logical chunks:

1. **feat(components): add form_container and table_list_container components** (4d3fee1)
   - Introduced FormContainer component for consistent form layout and styling
   - Added TableListContainer for table-based list display
   - Includes comprehensive test coverage

2. **feat(models): add new detail and editor models for burst and fact screens** (60d57d1)
   - Added BurstCard model for displaying burst information in card format
   - Added BurstDetails model for detailed burst view
   - Added BurstEditor model for editing burst information
   - Added FactDetails model for detailed fact view
   - Added FactSearch model for searching facts
   - Includes fact_list_indicator tests for improved coverage

3. **feat(models): add action_menu model for contextual actions** (e94534b)
   - Implemented ActionMenu model for displaying context-specific actions
   - Supports keyboard navigation through menu items
   - Enables dynamic action selection based on screen context

4. **feat(navigation): enhance navigation system with constants and help integration** (5450b1f)
   - Added centralized navigation constants for consistent screen identification
   - Extended constants with new screen types and navigation flows
   - Updated help system integration for improved user guidance
   - Comprehensive test coverage for navigation constants

5. **refactor(models): restructure core models for improved consistency and functionality** (9f1dfa8)
   - Refactored List model with enhanced navigation and rendering logic
   - Updated FactList with improved layout and filtering capabilities
   - Restructured BurstList for better state management and UI consistency
   - Enhanced Form model with better field handling and validation
   - Updated FactEditor with streamlined editing workflow
   - Improved test coverage for form layout integration

6. **feat(models): enhance menu, messages, and event display models** (ff2db72)
   - Updated Menu model with additional navigation support
   - Expanded Messages model for new message types
   - Refined BurstSuggestion model for better recommendations
   - Updated ViewEventWithFacts model for improved event display

7. **feat(app): integrate new models and components into application flow** (bc8e6cb)
   - Updated App model with support for new screens and navigation flows
   - Enhanced service layer to support new detail and editor models
   - Updated app integration tests for new screen states
   - Integrated burst editor and fact editor into app navigation
   - Added support for detail views (BurstDetails, FactDetails)
   - Updated e2e tests to reflect new navigation patterns

8. **docs(agents): update project handover document** (1fabdd3)
   - Removed duplicate handover documentation
   - Ensured single source of truth for project documentation

9. **build: rebuild cli binary with latest changes** (aeaea56)
   - Rebuilt cli binary to include all latest code changes
