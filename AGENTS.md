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

## Future Improvements
- Expand metadata enrichment
- Enhance burst and fact detection
- Improve export capabilities
- Add more comprehensive reporting

## Ongoing Feature Development

### Active Feature Tracks
1. **TUI Standardization** (Phase 1 Complete)
   - Standardizing model rendering
   - Ensuring component consistency
   - Implementing display validation system

2. **UX Enhancement and Model Standardization** (Phase 1: Navigation Completed ✅)
   - **Navigation State Management** (COMPLETED)
     - Centralized NavigationRegistry with screen definitions
     - Context preservation between screens
     - Intelligent breadcrumb generation (Default and Hierarchical strategies)
     - Universal back/forward navigation with explicit parent support
     - Undo/redo infrastructure with future stack
     - 122 comprehensive tests (99 unit + 23 integration)
     - Full thread-safety with mutex protection
     - Reuses existing NavigationKey constants and help system
   - Comprehensive model reorganization (in progress)
   - Standardized component library (in progress)
   - Predictable navigation patterns (in progress)
   - Elimination of legacy interaction remnants (in progress)
   - Consistent layout and interaction design (in progress)

## Contact & Support
- Project Repository: https://github.com/baphled/kariya
- Issue Tracker: https://github.com/baphled/kariya/issues

## Final Notes
This project represents a comprehensive career tracking solution. Maintain the focus on user experience, data quality, and continuous improvement.

**Handover Date**: 2025-12-31
**Prepared By**: Senior Development Engineer


## Recent Compliance and Quality Improvements

### Code Quality Compliance (Completed - 2025-12-31)
As part of ensuring the master-task-prompt workflow is being followed, the following improvements were made:

1. **Code Formatting**
   - Applied gofmt to 45+ unformatted Go files
   - All code now follows Go formatting standards

2. **Test Suite Consolidation**
   - Removed duplicate `suite_test.go` from models package that caused RunSpecs to be called twice
   - Consolidated navigation tests into proper structure:
     - `constants_test.go` - White-box tests for NavigationKey constants
     - `registry_test.go` - Black-box tests for NavigationRegistry
     - `registry_integration_test.go` - Integration tests for navigation state management
   - Ensured each package has only one TestFunction entry point

3. **Test Results**
   - 768 tests passing with zero failures
   - Zero race conditions detected
   - Test coverage at 76.57% (needs improvement to 80%+)
   - All Ginkgo test suites properly structured

4. **Compliance Status**
   - ✅ Code formatting: PASS
   - ✅ Build: PASS
   - ✅ Tests: PASS (768/768)
   - ✅ Go Vet: PASS
   - ✅ Race Detection: PASS
   - ⚠️ Coverage: 76.57% (Target: 80%)
   - ✅ Architectural compliance: PASS
   - ✅ Documentation: PASS
   - ✅ Git health: PASS

5. **Next Steps**
   - Improve test coverage to 80%+ by adding tests for newly implemented features
   - Continue following master-task-prompt workflow for all future work
   - Consider splitting large changesets into atomic commits for better maintainability

### Master Task Prompt Adherence
The project now has infrastructure to support the master-task-prompt workflow:
- Proper test structure with Ginkgo v2
- Atomic commit support with Make commands
- Compliance checking via `make check-compliance`
- Clear separation of concerns in test organization

### Style System Centralization (Completed - 2025-12-31)

Task 2.0 from `tasks/tasks-07-model-consistency.md` completed successfully:

#### Achievements
1. **Color Audit & Documentation**
   - Audited all 16 color constants in styles.go
   - Documented 42+ style variables and their purposes
   - Identified all spacing conventions (padding, margins, widths)
   - Verified zero inline hex colors exist

2. **Constants Export System**
   - Created `internal/cli/styles/constants_export.go` with:
     - 71 getter functions for all colors and styles
     - SpacingConstants struct for consistent spacing
     - Comprehensive documentation comments for each export
     - Clear organization by category (Colors, Buttons, Inputs, Cards, etc.)

3. **Style Usage Guide**
   - Created `docs/guides/STYLE_USAGE_GUIDE.md` with:
     - Complete color palette reference with use cases
     - Pre-built style examples and usage patterns
     - Spacing constants and layout helpers
     - Best practices and anti-patterns
     - Migration guide from inline styles to constants
     - Common patterns and troubleshooting

4. **Verification & Compliance**
   - Verified all 16 components use exported constants
   - Static analysis: 0 violations found
   - All 768+ tests passing
   - No inline hex colors or magic numbers in components

#### Files Created/Modified
- ✅ Created: `internal/cli/styles/constants_export.go` (400+ lines)
- ✅ Created: `docs/guides/STYLE_USAGE_GUIDE.md` (600+ lines)
- ✅ Modified: Task checklist in `tasks/tasks-07-model-consistency.md`

#### Quality Metrics
- 71 getter functions exported
- 16 color constants centralized
- 42+ style objects documented
- 100% of colors use constants (0 inline hex)
- 100% of components verified
- All tests passing (768/768)

#### Next Steps
- Task 2.0 Complete ✅
- Ready for Task 3.0: Adopt Containers in High-Impact Models (Phase 1)
  - Refactor FormModel to use containers
  - Refactor ListModel to use containers
  - Refactor DetailsModel to use containers
  - And 5 more model refactorings
