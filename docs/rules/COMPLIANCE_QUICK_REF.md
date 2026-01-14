---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Rules Compliance & Token Efficiency Check

## Purpose

This prompt ensures ALL project rules are followed and token usage stays efficient during development. Use this regularly to verify compliance and maintain focus.

---

## Quick Compliance Check (5 minutes)

### ✅ Code Quality
- [ ] `go fmt ./...` - Code formatted
- [ ] `go vet ./...` - No static analysis warnings
- [ ] `staticcheck ./...` - No advanced static analysis warnings
- [ ] `go test ./...` - All tests pass
- [ ] `go test -race ./...` - No race conditions
- [ ] Coverage ≥ 80%

### ✅ Commit Rules
- [ ] ONE logical change per commit
- [ ] Commit message: `<type>(<scope>): <subject>`
- [ ] No generated files (*.out, *.exe, coverage.*)
- [ ] No debug code (TODO, FIXME, fmt.Println)
- [ ] `make review-commit` passes

### ✅ Task Processing
- [ ] Working on ONE task at a time
- [ ] Following locked task checklist
- [ ] Using tools to verify (not assumptions)
- [ ] Following existing patterns

### ✅ Token Efficiency
- [ ] Token count < 50,000
- [ ] Using tools over text (view, grep, ls)
- [ ] Concise responses (no verbose explanations)
- [ ] No repeated information

---

## Automated Check

Run this before commits and PRs:

```bash
make check-compliance
```

---

## Common Violations & Fixes

### Large Commit (>10 files)
```bash
git reset
# Commit by logical groups
```

### No Tests with Code
```bash
# Write tests first, then commit together
```

### Generated Files Staged
```bash
git reset coverage.out *.exe
```

### Token Count High (>50k)
```bash
# Start fresh conversation
# Focus on current task only
```

---

## Daily Workflow

```bash
# Morning
make check-compliance

# Before each commit
make check-compliance   # REQUIRED before every commit
make ai-commit MSG="type(scope): description"

# After significant work
make check-compliance
```

---

**Full Guide:** [docs/rules/rules-compliance-check.md](./rules-compliance-check.md)

