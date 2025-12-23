# Atomic Commits Guidelines

---

## Overview

Atomic commits are the cornerstone of a clean, maintainable Git history. Each commit should represent a **single logical change** that makes sense in isolation. This practice makes code reviews easier, debugging more efficient, and rollbacks safer.

This document outlines our project's expectations for creating atomic commits and provides practical guidance for maintaining high-quality commit practices.

---

## Table of Contents

1. [What is an Atomic Commit?](#what-is-an-atomic-commit)
2. [Why Atomic Commits Matter](#why-atomic-commits-matter)
3. [Core Principles](#core-principles)
4. [Commit Message Guidelines](#commit-message-guidelines)
5. [Breaking Down Work](#breaking-down-work)
6. [What NOT to Do](#what-not-to-do)
7. [Practical Examples](#practical-examples)
8. [Pre-Commit Checklist](#pre-commit-checklist)
9. [Recovery Strategies](#recovery-strategies)
10. [Integration with Workflow](#integration-with-workflow)

---

## What is an Atomic Commit?

An **atomic commit** is a commit that:
- Contains exactly **one logical change**
- Is **complete and functional** (builds successfully, tests pass)
- Can be **understood independently** without requiring other commits
- Can be **reverted safely** without breaking unrelated functionality

### Atomic ≠ Small

Atomic doesn't mean "tiny." A commit can span multiple files and hundreds of lines if they all contribute to **one coherent change**. The key is **logical cohesion**, not file count or line count.

---

## Why Atomic Commits Matter

### 1. **Easier Code Reviews**
- Reviewers can understand each change in isolation
- Focused feedback on specific functionality
- Faster review cycles

### 2. **Better Debugging**
- `git bisect` works effectively to pinpoint bugs
- Clear history shows when and why features were introduced
- Easy to identify which commit introduced a regression

### 3. **Safer Rollbacks**
- Revert specific functionality without undoing unrelated changes
- Reduced risk when backing out problematic code
- Granular control over production deployments

### 4. **Clearer Documentation**
- Commit history serves as project documentation
- Future developers (including your future self) understand intent
- Easier onboarding for new team members

### 5. **Better Collaboration**
- Reduced merge conflicts
- Clearer understanding of parallel work
- Easier to cherry-pick fixes across branches

---

## Core Principles

### 1. One Logical Change Per Commit

✅ **Good:**
```
commit: Add user authentication middleware
- Create JWT validation middleware
- Add error handling for expired tokens
- Include unit tests for middleware
```

❌ **Bad:**
```
commit: Add authentication and fix bug in logger and update README
- Create JWT validation middleware
- Fix logger timestamp format
- Update installation instructions
- Refactor user service
```

### 2. Each Commit Must Be Buildable

Every commit should:
- Compile successfully
- Pass all existing tests
- Not break existing functionality
- Include tests for new functionality (when applicable)

This ensures `git bisect` remains a reliable debugging tool.

### 3. Commit Messages Tell the Story

A commit message should answer:
- **What** changed? (summary line)
- **Why** did it change? (body)
- **How** does it work? (body, if complex)

### 4. Complete the Change

Don't split a feature across commits if it leaves the codebase in a broken state. For example:

❌ **Bad:**
```
commit 1: Add user model (references missing table)
commit 2: Add database migration
```

✅ **Good:**
```
commit: Add user model with database migration
```

### 5. Separate Refactoring from Features

Refactoring and feature work should be separate commits:

✅ **Good:**
```
commit 1: Refactor event validation into separate function
commit 2: Add support for multi-tag filtering
```

❌ **Bad:**
```
commit: Add multi-tag filtering and refactor validation
```

---

## Commit Message Guidelines

### Structure

Follow the **Conventional Commits** format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Components

#### 1. Type (Required)
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style/formatting (no logic change)
- `refactor`: Code restructuring (no behavior change)
- `test`: Adding or updating tests
- `chore`: Maintenance tasks (dependencies, tooling)
- `perf`: Performance improvements

#### 2. Scope (Optional)
The area of the codebase affected:
- `(cli)`: CLI interface
- `(domain)`: Domain layer
- `(service)`: Service layer
- `(repo)`: Repository layer
- `(logger)`: Logging system

#### 3. Subject (Required)
- Max 50 characters
- Imperative mood ("Add feature" not "Added feature")
- No period at the end
- Lowercase after type/scope

#### 4. Body (Recommended)
- Wrap at 72 characters
- Explain **why** the change was needed
- Describe any side effects or implications
- Reference related issues/tickets

#### 5. Footer (Optional)
- Breaking changes: `BREAKING CHANGE: <description>`
- Issue references: `Closes #123`, `Fixes #456`

### Examples

#### Good Commit Messages

```
feat(cli): add event capture form with validation

Implement interactive form for capturing career events with:
- Real-time character counter (2000 char limit)
- Date parsing supporting relative dates
- Mode selection (Timeline/CV Backfill/Manual)
- Inline validation with error display

The form enforces domain validation rules before submission
and provides immediate feedback to users.

Closes #42
```

```
fix(service): enforce 30-day window for timeline journaling

Timeline journaling mode should only accept events from the
past 30 days. Previous implementation allowed any past date,
violating product requirements.

This fix adds date validation in CaptureEvent() and returns
a descriptive error when dates exceed the 30-day window.

Fixes #67
```

```
refactor(repo): extract common query building logic

Consolidate duplicate query construction code from List() and
Count() methods into buildFilterQuery() helper function.

This reduces duplication and makes it easier to add new filter
types in the future. No behavior changes.
```

```
test(classification): add coverage for multi-category events

Add test cases for events that match multiple competency
categories to ensure priority-based classification works
correctly.

Increases classification test coverage to 100%.
```

#### Bad Commit Messages

❌ **Too vague:**
```
fix: update code
```

❌ **Multiple unrelated changes:**
```
feat: add CLI form and fix logger bug and update docs
```

❌ **Not imperative mood:**
```
feat: added new feature for users
```

❌ **Missing context:**
```
fix: change date validation
```
(Why? What was broken?)

---

## Breaking Down Work

### Strategy 1: Bottom-Up (Foundation First)

Build from domain layer upward:

```
1. feat(domain): add CareerEvent model with validation
2. feat(repo): implement in-memory repository for events
3. feat(service): add event capture service with modes
4. test(service): add comprehensive service layer tests
5. feat(cli): create event capture form UI
6. test(cli): add form validation tests
```

### Strategy 2: Feature Slicing (Vertical Slice)

Implement minimal end-to-end functionality, then enhance:

```
1. feat: add basic event capture (happy path only)
2. feat(domain): add tag validation to events
3. feat(service): add multi-mode event capture
4. feat(cli): add interactive form with validation
5. feat(classification): add automatic competency tagging
```

### Strategy 3: Test-Driven (Red-Green-Refactor)

```
1. test(domain): add failing test for event validation
2. feat(domain): implement event validation (make test pass)
3. refactor(domain): extract validation logic into methods
4. test(service): add failing test for event capture
5. feat(service): implement event capture (make test pass)
```

### When to Combine Changes

It's acceptable to combine changes in ONE commit when:
- Adding a new file with its corresponding test file
- Creating a model with its validation logic
- Adding a feature with necessary configuration
- Including documentation updates for new code

✅ **Good combined commit:**
```
feat(logger): add structured logging with context support

Add logger package with:
- Logger interface and implementation
- Log levels (Debug, Info, Warn, Error, Fatal)
- Context enrichment with field chaining
- Comprehensive unit tests
- Usage documentation
```

---

## What NOT to Do

### ❌ 1. Don't Commit Work-in-Progress

**Bad:**
```
commit: WIP - halfway through feature
```

**Why:** Breaks the build, confuses reviewers, clutters history

**Instead:** Use git stash or local branches for in-progress work

### ❌ 2. Don't Commit Generated Files

**Bad:**
```
Changes to be committed:
  new file:   coverage.out
  new file:   coverage.html
  modified:   go.sum
```

**Why:** Binary/generated files bloat repository and cause merge conflicts

**Instead:** Add to `.gitignore`:
```gitignore
# Test coverage
coverage.out
coverage.html
*.test

# Build artifacts
*.exe
*.dll
*.so
dist/
build/
```

### ❌ 3. Don't Mix Concerns

**Bad:**
```
commit: Add new feature and fix unrelated bug and update deps
```

**Why:** Makes it impossible to review, revert, or cherry-pick changes

**Instead:** Three separate commits:
```
1. fix(logger): correct timestamp format in log output
2. feat(cli): add event listing screen
3. chore(deps): update ginkgo to v2.27.3
```

### ❌ 4. Don't Skip Tests

**Bad:**
```
commit: Add new service method
(no tests included)
```

**Why:** Reduces test coverage, increases bug risk, makes refactoring dangerous

**Instead:**
```
commit: feat(service): add GetEventsByTags method with tests

Implement GetEventsByTags() to filter events by multiple tags.
Includes comprehensive unit tests covering edge cases.
```

### ❌ 5. Don't Use Vague Messages

**Bad:**
```
commit: fix stuff
commit: update
commit: changes
```

**Why:** Future developers (including you) won't understand what changed or why

**Instead:** Be specific and descriptive (see [Commit Message Guidelines](#commit-message-guidelines))

### ❌ 6. Don't Commit Formatting with Logic

**Bad:**
```
commit: Add validation and reformat entire file
```

**Why:** Makes it hard to review actual logic changes

**Instead:** Two commits:
```
1. style(domain): run gofmt on event.go
2. feat(domain): add tag deduplication to event validation
```

---

## Practical Examples

### Example 1: Current Staging Area Analysis

**Problem:** You have 48 files staged for initial commit including:
- Domain models
- Service layer
- Repository implementations
- CLI interface
- Tests
- Documentation
- Configuration
- **Coverage reports (should not be committed)**

**Solution:** Break into logical atomic commits:

```bash
# 1. Project initialization
git reset
git add go.mod go.sum ginkgo.yml README.md
git commit -m "chore: initialize Go project with dependencies

Add Go modules configuration and Ginkgo test framework setup.
Include project README with basic project overview."

# 2. Domain layer
git add internal/domain/career/*.go
git commit -m "feat(domain): add CareerEvent domain model

Implement core domain model with:
- CareerEvent struct with validation
- Tag validation against allowed set
- Date validation (no future dates)
- Maximum text length enforcement (2000 chars)
- Comprehensive unit tests (100% coverage)

The domain model is the foundation for career event tracking
and enforces all business rules at the domain level."

# 3. Logger infrastructure
git add internal/logger/*.go
git commit -m "feat(logger): add structured logging with context

Implement logger package with:
- Multiple log levels (Debug, Info, Warn, Error, Fatal)
- Context enrichment via field chaining
- Caller information tracking
- Thread-safe operations
- Comprehensive unit tests

Provides consistent logging interface across the application."

# 4. Repository interface and memory implementation
git add internal/repository/career/repository.go
git add internal/repository/career/memory_repository.go
git add internal/repository/career/repository_test.go
git commit -m "feat(repo): add repository interface with in-memory implementation

Define Repository interface for career event persistence with:
- CRUD operations (Create, GetByID, Update, Delete)
- List with filtering (tags, date ranges, pagination, sorting)
- Count with filtering

Include in-memory implementation for testing and development.
All repository tests passing with high coverage."

# 5. SQLite repository
git add internal/repository/career/sqlite_repository.go
git add internal/repository/career/sqlite_repository_test.go
git add internal/repository/career/integration_test.go
git commit -m "feat(repo): add SQLite repository implementation

Implement production-ready SQLite persistence layer:
- Schema creation with proper indexes
- Parameterized queries (SQL injection prevention)
- Transaction support
- Comprehensive integration tests

Provides durable storage for career events."

# 6. Classification service
git add internal/service/career/classification/*.go
git commit -m "feat(classification): add competency classification system

Implement automatic event classification into 6 competency categories:
- Technical (backend, frontend, system, architecture)
- Leadership (management, strategy, vision)
- Product (feature design, customer focus)
- Consulting (advisory, transformation)
- Research (investigation, data analysis)
- Mentoring (coaching, training)

Uses keyword matching with priority-based ordering.
Includes comprehensive unit tests (8 test cases, 100% coverage)."

# 7. Service layer
git add internal/service/career/service.go
git add internal/service/career/service_test.go
git commit -m "feat(service): add career event service with multi-mode capture

Implement business logic layer with:
- Event capture with 3 modes (Timeline, CV Backfill, Manual)
- Mode-specific validation (30-day window for Timeline)
- CRUD operations through repository abstraction
- Structured logging for all operations
- UUID and timestamp management

Includes comprehensive unit tests with 100% coverage."

# 8. CLI styles foundation
git add internal/cli/styles/*.go
git commit -m "feat(cli): add comprehensive TUI styling system

Create lipgloss-based styling for terminal UI:
- Professional dark theme color scheme
- Component styles (buttons, inputs, cards, lists)
- Message styles (error, warning, success, info)
- Responsive layout helpers
- Comprehensive style tests

Provides consistent visual language for CLI interface."

# 9. CLI service adapter
git add internal/cli/service/*.go
git commit -m "feat(cli): add CLI service adapter layer

Implement adapter between CLI and core domain services:
- Functional options pattern for flexibility
- Event capture with optional fields
- Event retrieval with filtering
- Clean separation from domain logic

Includes unit tests with 88% coverage."

# 10. CLI form model
git add internal/cli/models/*.go
git commit -m "feat(cli): add interactive event capture form

Implement BubbleTea form with:
- Real-time character counter (2000 char limit)
- Date parsing (ISO format, relative dates)
- Mode selection interface
- Inline validation with error display
- Tab navigation between fields

Provides user-friendly event capture interface."

# 11. CLI app integration
git add internal/cli/app/*.go
git commit -m "feat(cli): add BubbleTea app with screen navigation

Implement main CLI application model:
- Screen management (Home, Capture, List, View)
- Keyboard navigation
- Window resize handling
- State management following Elm Architecture

Includes app tests with 88% coverage."

# 12. CLI entry point
git add cmd/cli/*.go
git commit -m "feat(cli): add main entry point with dependency injection

Create application bootstrap:
- Command-line flag parsing (--version, --help)
- Dependency wiring (repository, services)
- BubbleTea program initialization

Provides executable entry point for KaRiya CLI."

# 13. Documentation
git add docs/*.md features/*.md tasks/*.md
git commit -m "docs: add comprehensive project documentation

Include:
- Product requirements (PRD)
- Architecture overview
- Feature specifications
- Task tracking documents
- Integration test strategy
- Development guidelines

Provides context for development and onboarding."

# 14. Development guidelines
git add docs/rules/*.md
git commit -m "docs(rules): add development guidelines and standards

Add guidelines for:
- Go coding standards
- Senior engineer workflows
- Task processing
- PRD and task generation
- Code review expectations

Establishes team standards and expectations."

# 15. Project handover
git add AGENTS.md
git commit -m "docs: add comprehensive handover document

Create AGENTS.md with:
- Complete project overview
- Architecture documentation
- Component reference guide
- Development setup instructions
- Testing strategy
- Troubleshooting guide
- Future enhancement roadmap

Serves as primary handover document for new developers."

# 16. Test configuration
git add Makefile
git commit -m "chore: add Makefile for development tasks

Add targets for:
- Running tests (test, test-suite, individual-test)
- Coverage generation (coverage, clean-coverage)
- Build automation

Standardizes development workflows."
```

### Example 2: Fixing the Current Situation

If you've already staged everything:

```bash
# Unstage everything
git reset

# Add .gitignore first to prevent committing generated files
cat > .gitignore << 'EOF'
# Test coverage
coverage.out
coverage/
*.test

# Build artifacts
*.exe
*.dll
*.so
dist/
build/
kariya

# IDE
.vscode/
.idea/

# OS
.DS_Store
Thumbs.db

# Temporary files
*.swp
*.swo
*~
EOF

git add .gitignore
git commit -m "chore: add .gitignore for generated files"

# Now follow the breakdown from Example 1
```

### Example 3: Using Interactive Staging

For files with mixed changes:

```bash
# Stage specific hunks interactively
git add -p internal/service/career/service.go

# Review what's staged
git diff --cached

# Commit staged changes
git commit -m "feat(service): add event filtering by date range"

# Stage and commit remaining changes
git add internal/service/career/service.go
git commit -m "refactor(service): extract validation into helper method"
```

---

## Pre-Commit Checklist

Before committing, verify:

### ✅ Commit Quality Checklist

- [ ] **Single Logical Change**: Does this commit represent exactly one thing?
- [ ] **Builds Successfully**: Does `go build` succeed?
- [ ] **Tests Pass**: Does `go test ./...` pass?
- [ ] **No Race Conditions**: Does `go test -race ./...` pass?
- [ ] **Proper Scope**: Are only relevant files included?
- [ ] **No Generated Files**: Are coverage reports/binaries excluded?
- [ ] **Clear Message**: Is the commit message descriptive and follows conventions?
- [ ] **Complete Change**: Is the functionality complete and working?
- [ ] **Tests Included**: Are new features/fixes covered by tests?
- [ ] **Documentation Updated**: Are relevant docs updated (if applicable)?

### Running the Checklist

```bash
# Create a pre-commit script
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash

echo "🔍 Running pre-commit checks..."

# Check build
echo "Building..."
if ! go build ./...; then
    echo "❌ Build failed"
    exit 1
fi

# Run tests
echo "Running tests..."
if ! go test ./...; then
    echo "❌ Tests failed"
    exit 1
fi

# Run race detector
echo "Checking for race conditions..."
if ! go test -race ./...; then
    echo "❌ Race conditions detected"
    exit 1
fi

# Check formatting
echo "Checking formatting..."
if [ -n "$(gofmt -l .)" ]; then
    echo "❌ Code not formatted. Run: go fmt ./..."
    exit 1
fi

echo "✅ All checks passed"
EOF

chmod +x .git/hooks/pre-commit
```

---

## Recovery Strategies

### Scenario 1: Already Committed Too Much

**Problem:** You made a large commit with multiple unrelated changes.

**Solution:** Split the commit using interactive rebase:

```bash
# Undo the last commit but keep changes staged
git reset --soft HEAD~1

# Unstage everything
git reset

# Now commit changes in logical groups
git add <files-for-change-1>
git commit -m "First logical change"

git add <files-for-change-2>
git commit -m "Second logical change"
```

### Scenario 2: Forgot to Include Files

**Problem:** Committed a change but forgot to include necessary files.

**Solution:** Amend the commit (if not pushed):

```bash
git add forgotten-file.go
git commit --amend --no-edit
```

**Note:** Never amend commits that have been pushed to shared branches.

### Scenario 3: Wrong Commit Message

**Problem:** Committed with a typo or unclear message.

**Solution:** Rewrite the commit message (if not pushed):

```bash
git commit --amend
# Edit message in editor, save, and close
```

### Scenario 4: Mixed Changes in Working Directory

**Problem:** Multiple unrelated changes in files, need to commit separately.

**Solution:** Use interactive staging:

```bash
# Stage specific changes interactively
git add -p filename.go

# Or use git gui for visual selection
git gui
```

### Scenario 5: Need to Insert a Commit in History

**Problem:** Realized a commit should come before recent commits.

**Solution:** Use interactive rebase:

```bash
# Rebase last 3 commits interactively
git rebase -i HEAD~3

# In the editor, reorder commits or mark for edit
# Save and follow instructions
```

---

## Integration with Workflow

### Red-Green-Refactor Cycle

Atomic commits align naturally with TDD:

```
1. Write failing test
   → commit: "test(service): add failing test for tag validation"

2. Implement minimal code to pass
   → commit: "feat(service): add basic tag validation"

3. Refactor
   → commit: "refactor(service): extract tag validation into helper"
```

### Feature Branch Workflow

```bash
# Create feature branch
git checkout -b feature/event-filtering

# Make atomic commits
git commit -m "feat(repo): add tag filtering to repository"
git commit -m "feat(service): add tag filtering to service layer"
git commit -m "test(service): add comprehensive filtering tests"
git commit -m "feat(cli): add tag filter UI component"

# Push feature branch
git push origin feature/event-filtering

# Create pull request for review
```

### Code Review Process

Reviewers can:
- Review each commit individually
- Provide targeted feedback per logical change
- Approve incrementally
- Request changes to specific commits

```bash
# Review specific commit
git show <commit-hash>

# Review commit by commit
git log -p
```

---

## Tools and Automation

### Git Aliases for Atomic Commits

Add to `~/.gitconfig`:

```ini
[alias]
    # Show files in last commit
    last = log -1 HEAD --stat

    # Amend without editing message
    amend = commit --amend --no-edit

    # Interactive staging
    stage = add -p

    # Show diff of staged changes
    staged = diff --cached

    # Undo last commit (keep changes staged)
    undo = reset --soft HEAD~1

    # Undo last commit (keep changes unstaged)
    uncommit = reset HEAD~1

    # Show commit graph
    graph = log --graph --oneline --all
```

### Makefile Targets

Add to your `Makefile`:

```makefile
.PHONY: pre-commit check-commit

pre-commit:
	@echo "Running pre-commit checks..."
	@go fmt ./...
	@go build ./...
	@go test ./...
	@go test -race ./...
	@echo "✅ Ready to commit"

check-commit:
	@echo "Checking last commit..."
	@git diff --stat HEAD~1 HEAD
	@echo ""
	@git log -1 --pretty=format:"%B"
```

---

## References and Further Reading

### Internal Documentation
- [Senior Engineer Guidelines](./senior-engineer-guidelines.md)
- [Go Guidelines](./go-guidelines.md)
- [Task Processing](./process-task-list.md)

### External Resources
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Git Book - Commit Guidelines](https://git-scm.com/book/en/v2/Distributed-Git-Contributing-to-a-Project)
- [Thoughtbot - 5 Useful Tips For A Better Commit Message](https://thoughtbot.com/blog/5-useful-tips-for-a-better-commit-message)
- [Chris Beams - How to Write a Git Commit Message](https://chris.beams.io/posts/git-commit/)

---

## Summary

**Key Takeaways:**

1. **One commit = One logical change**
2. **Every commit must build and pass tests**
3. **Write clear, descriptive commit messages**
4. **Separate refactoring from features**
5. **Don't commit generated files**
6. **Use the pre-commit checklist**
7. **Review your commits before pushing**

By following these guidelines, you'll create a clean, maintainable Git history that serves as valuable documentation and makes collaboration smoother.

**Remember:** Good commits are a gift to your future self and your teammates. Take the time to get them right.

---

*Last Updated: 2025-12-23*
*Version: 1.0*

