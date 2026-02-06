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

## Amending Commits Without Attribution

If a commit was made without proper AI attribution (e.g., using `git commit` directly), you can add attribution to HEAD:

```bash
make ai-commit AMEND=1
```

This will:
1. **Verify HEAD is unpushed** - Cannot amend pushed commits
2. **Extract existing commit message** - Preserves original message
3. **Strip old attribution** - Removes any existing AI trailers
4. **Add correct attribution** - Adds AI-Generated-By and Reviewed-By

### When to Amend

Use `AMEND=1` when:
- You accidentally used `git commit` instead of `make ai-commit`
- A commit was created without AI attribution
- You need to update attribution on an unpushed commit

### Safety Checks

The amend operation includes safety checks:
- **Unpushed verification** - Refuses to amend if HEAD is already on remote
- **Existing attribution warning** - Prompts before replacing existing attribution
- **Branch protection** - Same rules as normal commits

### Combined Options

```bash
make ai-commit AMEND=1 NO_VERIFY=1  # Amend without running hooks
```

**Warning**: After amending, if you had already pushed, you'll need `git push --force-with-lease`.

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
- Use `git commit --amend` directly (use `make ai-commit AMEND=1` instead)
- Use `git push --force` (use `--force-with-lease` if absolutely necessary)
- Skip hooks without explicit user request

**Amend is ONLY safe when ALL conditions are met:**
1. HEAD commit was created in this session (not by someone else)
2. HEAD commit has NOT been pushed to remote
3. You're adding attribution, not changing the actual content

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

### Missing AI Attribution on Commit

If you accidentally committed without `make ai-commit`:

```bash
# Check if commit is unpushed
git status  # Should show "Your branch is ahead of..."

# Add attribution to HEAD
make ai-commit AMEND=1
```

### Already Pushed Without Attribution

If the commit was already pushed:

```bash
# Option 1: Create a follow-up commit (safer)
# Just continue with proper attribution on next commit

# Option 2: Amend and force push (use with caution)
make ai-commit AMEND=1
git push --force-with-lease
```

**Warning**: Force pushing rewrites history. Only do this on feature branches where you're the only contributor.

## Related skills

- `check-compliance` - Run before committing
- `create-pr` - Create PR after commits
- `tdd-workflow` - Follow TDD for coverage
