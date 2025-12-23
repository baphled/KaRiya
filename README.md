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

