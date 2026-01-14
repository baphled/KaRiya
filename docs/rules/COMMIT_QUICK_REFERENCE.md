# Atomic Commits - Quick Reference

## Before Every Commit

```bash
# 1. Stage your changes
git add -p <file>

# 2. Review what you're committing
git diff --cached --stat

# 3. REQUIRED: Run compliance check before commit
make check-compliance

# 4. Commit with AI attribution (REQUIRED for AI-generated code)
make ai-commit MSG="feat(scope): description"

# OR: Manual commit (NOT recommended for AI-generated code)
git commit
```

---

## Commit Message Template

```
<type>(<scope>): <subject line - 50 chars max>

<body - explain WHY this change is needed>
<wrap at 72 characters>

AI-Generated-By: <Assistant> (<Model>) [IF AI-GENERATED]
Reviewed-By: <Your Name> [IF AI-GENERATED]
<footer - issue references>
```

**🤖 IMPORTANT**: If ANY code was AI-generated, you MUST include:
- `AI-Generated-By: <Assistant Name> (<Model Version>)`
- `Reviewed-By: <Your Name>`

See [AI Commit Attribution Rules](./AI_COMMIT_ATTRIBUTION.md) for details.

### Types
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `refactor`: Code restructuring (no behavior change)
- `test`: Adding/updating tests
- `chore`: Maintenance (deps, tooling)
- `style`: Formatting only
- `perf`: Performance improvement

### Scopes (optional)
- `(domain)`: Domain layer
- `(service)`: Service layer
- `(repo)`: Repository layer
- `(cli)`: CLI interface
- `(logger)`: Logging system

---

## Examples

### ✅ Good Commits

```
feat(service): add event filtering by date range

Implement date range filtering for timeline views and
historical analysis. Improves performance for large
event collections by limiting result sets.

Closes #56
```

```
fix(domain): prevent duplicate tags in event validation

Tag deduplication was not working correctly when tags
had different casing. Now normalize to lowercase before
checking for duplicates.

Fixes #87
```

```
refactor(repo): extract query building into helper method

Consolidate duplicate query construction code from List()
and Count() methods. No behavior changes.
```

```
test(classification): add coverage for edge cases

Add tests for multi-category classification and empty
text scenarios. Increases coverage to 100%.
```

### ✅ AI-Generated Commits (with proper attribution)

```
feat(cli): add interactive event capture form

Implement BubbleTea form with real-time validation,
character counter, and date parsing.

Features:
- Text input with 2000 character limit
- Smart date parsing (ISO, relative, special keywords)
- Visual focus indicators
- Inline validation with error display

Human review confirmed:
- All tests pass (14/14 passing)
- No security issues
- Follows project conventions

AI-Generated-By: Avante (Claude 3.5 Sonnet)
Reviewed-By: John Doe
Closes #42
```

---

## ❌ Bad Commits (Don't Do This)

```
update code              ❌ Too vague
fix stuff                ❌ No context
WIP                      ❌ Work in progress
Add feature and fix bug  ❌ Multiple changes
```

---

## Quick Checklist

Before committing, ask:

- [ ] **One logical change?** Can I describe this in one sentence?
- [ ] **Builds?** Does `go build ./...` succeed?
- [ ] **Tests pass?** Does `go test ./...` pass?
- [ ] **No generated files?** (coverage.out, *.exe, etc.)
- [ ] **Clear message?** Does it explain WHY, not just WHAT?
- [ ] **Tests included?** (If changing behavior)

---

## Common Patterns

### New Feature (Layer by Layer)

```bash
# Domain layer
git add internal/domain/career/event.go internal/domain/career/event_test.go
git commit -m "feat(domain): add tags field to career event"

# Repository layer  
git add internal/repository/career/*.go
git commit -m "feat(repo): add tag filtering support"

# Service layer
git add internal/service/career/*.go  
git commit -m "feat(service): expose tag filtering in service"
```

### Bug Fix with Test

```bash
# Option 1: Test first (recommended)
git add internal/service/career/service_test.go
git commit -m "test(service): add failing test for nil event handling"

git add internal/service/career/service.go
git commit -m "fix(service): prevent nil pointer in CaptureEvent"

# Option 2: Combined (if small)
git add internal/service/career/service.go internal/service/career/service_test.go
git commit -m "fix(service): prevent nil pointer in CaptureEvent

Added nil check before accessing event properties.
Includes regression test.

Fixes #87"
```

### Refactoring

```bash
# Safe incremental refactoring
git add internal/repository/career/repository.go
git commit -m "refactor(repo): add Filter interface (no callers yet)"

git add internal/repository/career/filters.go
git commit -m "refactor(repo): implement Filter with builder pattern"

git add internal/service/career/service.go
git commit -m "refactor(service): migrate to Filter interface"
```

---

## Tools

### Makefile Targets

```bash
make fmt            # Format code
make vet            # Run go vet
make staticcheck    # Run staticcheck
make test           # Run all tests
make pre-commit     # Quick checks
make review-commit  # Comprehensive review
```

### Manual Commands

```bash
# Stage interactively (for mixed changes in one file)
git add -p filename.go

# Review staged changes
git diff --cached

# Unstage everything
git reset

# Amend last commit (if not pushed)
git commit --amend
```

---

## Recovery

### Split Large Commit

```bash
# Undo commit, keep changes staged
git reset --soft HEAD~1

# Unstage everything
git reset

# Commit in logical groups
git add <files-for-change-1>
git commit -m "First logical change"

git add <files-for-change-2>  
git commit -m "Second logical change"
```

### Add Forgotten Files

```bash
# If not pushed yet
git add forgotten-file.go
git commit --amend --no-edit
```

### Fix Commit Message

```bash
# If not pushed yet
git commit --amend
# Edit message, save, close
```

---

## What NOT to Commit

### Generated Files
- `coverage.out`, `coverage.html`
- `*.test`, `*.exe`, `*.dll`, `*.so`
- `dist/`, `build/`, `bin/`

### Temporary Files  
- `*.tmp`, `*.swp`, `*~`
- `.DS_Store`

### IDE Files
- `.vscode/`, `.idea/`

### Secrets
- API keys, passwords, tokens
- Use environment variables instead

**Always check `.gitignore` is configured correctly!**

---

## Reference

- Full guide: [docs/rules/atomic-commits.md](./atomic-commits.md)
- Review prompt: [docs/rules/review-commit-prompt.md](./review-commit-prompt.md)
- Review script: `scripts/review-commit.sh`

---

**Remember:** Good commits are a gift to your future self and teammates!

