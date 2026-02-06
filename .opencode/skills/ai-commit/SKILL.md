---
name: ai-commit
description: Create properly attributed commits for AI-generated code in KaRiya
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Create commits with proper AI attribution using the project's commit standards.

## When to use me

Use this skill when ready to commit changes. **NEVER use `git commit` directly.**

## Commit Procedure

### 1. Prepare Commit Message

Create a commit message file:

```bash
cat > /tmp/commit.txt << 'EOF'
type(scope): short description

Optional longer explanation of WHY the change was made.

Optional footer (issue refs)
EOF
```

### 2. Create Commit

```bash
make ai-commit FILE=/tmp/commit.txt
```

If pre-commit hooks need to be skipped (NOT recommended):
```bash
make ai-commit FILE=/tmp/commit.txt NO_VERIFY=1
```

## Commit Message Format

### Types
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation only
- `refactor` - Code restructuring
- `test` - Adding/updating tests
- `chore` - Maintenance tasks

### Scopes
- `domain` - Domain models
- `service` - Service layer
- `cli` - CLI application
- `intents` - Intent workflows
- `uikit` - UIKit components
- `forms` - Form components

### Examples

```
feat(intents): add skill management workflow

Implements the skill management intent with list, detail,
and edit screens following the subdirectory structure.

Closes #123
```

```
fix(cli): prevent crash on empty timeline

The timeline view now handles empty event lists gracefully
by showing an appropriate message instead of panicking.
```

```
refactor(screens): extract modal to dedicated package

Moves the filter modal from inline definition to
screens/timeline/modals/ for better organization.
```

## Pre-Commit Requirements

Before committing:

1. **Run compliance check:**
   ```bash
   make check-compliance
   ```

2. **Verify tests pass:**
   ```bash
   make test
   ```

3. **Check coverage:**
   ```bash
   go test -cover ./path/to/modified/packages/...
   ```
   Coverage must be >= 95% for modified packages.

## AI Attribution

The `ai-commit` script automatically adds:

```
AI-Generated-By: OpenCode
AI-Model: [detected model]
```

## Git Safety Rules

**NEVER:**
- Commit directly to `main` or `next`
- Use `git commit --amend` (unless fixing a just-created commit)
- Use `git push --force`
- Skip hooks without explicit user request

## Troubleshooting

### Pre-commit Hook Failed

Fix the issues, then create a NEW commit (don't amend):

```bash
# Fix issues reported by hook
# ...

# Create new commit
make ai-commit FILE=/tmp/commit.txt
```

### Coverage Too Low

```bash
# Find low-coverage functions
go test -coverprofile=/tmp/cover.out ./path/... && \
  go tool cover -func=/tmp/cover.out | grep -v "100.0%"
```

Add tests for uncovered code, then commit.

## Related skills

- `check-compliance` - Run before committing
- `create-pr` - Create PR after commits
- `tdd-workflow` - Follow TDD for coverage
