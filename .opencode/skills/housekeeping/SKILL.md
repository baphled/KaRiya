---
name: housekeeping
description: Maintain codebase health through regular cleanup, organization, and hygiene tasks
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Maintain codebase health through regular cleanup, organization, and hygiene tasks.

## When to use me

Use this skill when:
- Starting or ending a work session
- Codebase feels cluttered
- Build times are slow
- Tests are flaky
- Dependencies are outdated
- Before major releases

## Housekeeping Principles

1. **Regular maintenance prevents emergencies**
2. **Small frequent cleanups beat big infrequent ones**
3. **Automate what you can**
4. **Track what needs attention**

## Housekeeping Tasks

### Daily Hygiene

```bash
# Format code
make fmt

# Run quick checks
make pre-commit

# Clean build artifacts
go clean ./...
```

### Weekly Maintenance

```bash
# Update dependencies (review changes)
go get -u ./...
go mod tidy

# Check for security issues
make gosec

# Review test coverage
make coverage

# Clean up old branches
git branch --merged | grep -v "main\|next" | xargs git branch -d
```

### Monthly Deep Clean

```bash
# Full compliance audit
make check-compliance

# Find unused code
staticcheck -unused ./...

# Check for outdated patterns
make check-patterns

# Review TODOs
grep -rn "TODO\|FIXME" internal/

# Audit tech debt
# Create/update tech debt inventory
```

## Cleanup Categories

### Dead Code

```bash
# Find unused exports
go vet ./...

# Find unreferenced files
# Look for files not imported anywhere

# Find commented-out code
grep -rn "^[[:space:]]*//" internal/cli/*.go
```

**Action:** Remove dead code. Git has history.

### Obsolete Dependencies

```bash
# Find unused dependencies
go mod tidy

# Check for vulnerabilities
go list -m -json all | go run golang.org/x/vuln/cmd/govulncheck@latest

# Update outdated
go list -m -u all
```

**Action:** Remove unused, update outdated, fix vulnerabilities.

### Test Hygiene

```bash
# Find slow tests
go test -v ./... 2>&1 | grep -E "^---.*[0-9]+\.[0-9]+s"

# Find flaky tests (run multiple times)
for i in {1..5}; do make test || echo "Flaky on run $i"; done

# Check test coverage gaps
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep -v "100.0%"
```

**Action:** Speed up slow tests, fix flaky tests, improve coverage.

### Documentation Hygiene

```bash
# Check for missing docs
make check-docblocks

# Find outdated docs
# Compare doc dates to code changes

# Validate links
# Check that referenced files exist
```

**Action:** Add missing docs, update outdated ones, fix broken links.

### Git Hygiene

```bash
# Clean up merged branches
git branch --merged origin/next | grep -v "next\|main" | xargs git branch -d

# Find large files in history
git rev-list --objects --all | \
  git cat-file --batch-check='%(objecttype) %(objectname) %(objectsize) %(rest)' | \
  sort -k3 -n -r | head -20

# Verify hooks installed
make verify-hooks
```

**Action:** Delete merged branches, consider cleaning large files, ensure hooks work.

### Build Hygiene

```bash
# Clean build cache
go clean -cache

# Verify clean build
rm -rf bin/ && make build

# Check build time
time make build
```

**Action:** Keep builds fast and reproducible.

## Housekeeping Checklist

### Quick (5 min)
- [ ] `make fmt`
- [ ] `make pre-commit`
- [ ] No uncommitted changes left

### Standard (30 min)
- [ ] All tests pass
- [ ] No new linter warnings
- [ ] Dependencies up to date
- [ ] Merged branches deleted
- [ ] No debug code left

### Deep (2 hours)
- [ ] Full compliance check
- [ ] Coverage reviewed
- [ ] Security scan clean
- [ ] Documentation current
- [ ] Tech debt inventory updated
- [ ] Dependencies audited
- [ ] Test performance reviewed

## Automated Hygiene

### Pre-commit Hook
Already in place via `.git-hooks/pre-commit`:
- Format check
- Vet
- Tests
- Coverage

### CI Pipeline
- All checks run on PR
- Coverage reported
- Security scanning

### Add to Routine

```markdown
## Session Start
- [ ] `make session-start` (includes checks)

## Before Commit
- [ ] `make check-compliance`

## Session End
- [ ] No uncommitted changes
- [ ] No dangling TODOs
- [ ] Session notes written
```

## Warning Signs

| Sign | Meaning | Action |
|------|---------|--------|
| Build times increasing | Cache issues, bloat | Clean cache, review deps |
| Tests getting flaky | Race conditions, coupling | Fix isolation |
| Coverage dropping | New code untested | Enforce coverage gates |
| Linter warnings growing | Standards slipping | Address warnings |
| Merge conflicts frequent | Poor modularization | Refactor boundaries |

## Housekeeping Calendar

| Frequency | Tasks |
|-----------|-------|
| Every commit | Format, pre-commit checks |
| Daily | Quick hygiene |
| Weekly | Dependency updates, branch cleanup |
| Monthly | Deep clean, tech debt review |
| Quarterly | Major updates, architecture review |

## Related skills

- `clean-code` - Code quality standards
- `tech-debt` - Track maintenance needs
- `task-completer` - Finish cleanups properly
- `refactor` - Systematic improvements
