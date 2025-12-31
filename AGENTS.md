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

2. **UX Enhancement and Model Standardization** (Deep Analysis Phase)
   - Comprehensive model reorganization
   - Standardized component library
   - Predictable navigation patterns
   - Elimination of legacy interaction remnants
   - Consistent layout and interaction design

## Contact & Support
- Project Repository: https://github.com/baphled/kariya
- Issue Tracker: https://github.com/baphled/kariya/issues

## Final Notes
This project represents a comprehensive career tracking solution. Maintain the focus on user experience, data quality, and continuous improvement.

**Handover Date**: 2025-12-31
**Prepared By**: Senior Development Engineer

