---
name: task-completer
description: Ensure tasks are fully completed with all requirements met and no loose ends
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Ensure tasks are fully completed with all requirements met, no loose ends, and proper closure.

## When to use me

Use this skill when:
- Finishing a task
- Checking if work is truly done
- Reviewing before commit/PR
- Validating completeness

## The "Done" Problem

Tasks often appear done but aren't:
- Tests pass but edge cases missed
- Code works but docs outdated
- Feature complete but not integrated
- Changes made but not committed

**"Almost done" is not done.**

## Definition of Done

A task is DONE when:

```markdown
## Completion Checklist

### Code Quality
- [ ] Code compiles without warnings
- [ ] All tests pass
- [ ] Coverage >= 95% for modified code
- [ ] No linter warnings
- [ ] No TODO/FIXME left behind

### Architecture
- [ ] Code in correct layer
- [ ] Follows existing patterns
- [ ] No new dependencies on wrong layers
- [ ] Architecture check passes

### Testing
- [ ] Happy path tested
- [ ] Error cases tested
- [ ] Edge cases tested
- [ ] Integration points tested

### Documentation
- [ ] Exports have godoc
- [ ] Expected/Returns/Side effects sections
- [ ] README updated (if applicable)
- [ ] CHANGELOG updated (if applicable)

### Cleanup
- [ ] No debug code left
- [ ] No commented-out code
- [ ] Boy Scout Rule applied
- [ ] Temporary files removed

### Git
- [ ] Changes committed
- [ ] Commit message descriptive
- [ ] Branch up to date with base

### Verification
- [ ] `make check-compliance` passes
- [ ] Manual testing done (if UI)
- [ ] Works in clean environment
```

## Task Completion Process

### 1. Verify Requirements

Go back to the original task:

```markdown
## Task: Add date filter to timeline

### Acceptance Criteria
- [x] Filter modal shows date range picker
- [x] Empty range shows validation error
- [x] Filter applies to event list
- [ ] Filter persists across navigation  <-- MISSED!
```

**Check every criterion explicitly.**

### 2. Run All Checks

```bash
# Full compliance
make check-compliance

# Architecture
make check-intent-architecture

# Tests with coverage
go test -cover ./path/to/modified/...

# Patterns
make check-patterns
```

### 3. Self-Review

Before calling it done, review your own changes:

```bash
# See all changes
git diff

# Staged changes
git diff --cached

# Files changed
git status
```

Ask yourself:
- Would I approve this PR?
- Is anything unclear?
- Did I take shortcuts?
- What could break?

### 4. Clean Up

```bash
# Remove debug statements
grep -rn "fmt.Print" internal/cli/ | grep -v "_test.go"

# Find TODOs (should be none)
grep -rn "TODO\|FIXME" internal/cli/

# Check for commented code
grep -rn "^[[:space:]]*//" internal/cli/*.go | head -20
```

### 5. Final Verification

```bash
# Clean build
go build ./...

# Full test suite
make test

# One more compliance check
make check-compliance
```

### 6. Commit Properly

```bash
# Create descriptive commit message
cat > /tmp/commit.txt << 'EOF'
feat(timeline): add date range filter

Users can now filter timeline events by date range.
Validation prevents invalid ranges. Filter state
persists across navigation.

Closes #123
EOF

make ai-commit FILE=/tmp/commit.txt
```

## Loose End Detection

### Common Loose Ends

| Loose End | How to Find | Fix |
|-----------|-------------|-----|
| Unfinished code | TODOs, FIXMEs | Complete or create task |
| Missing tests | Coverage report | Add tests |
| Dead code | Unused functions | Remove |
| Debug code | Print statements | Remove |
| Commented code | // blocks | Remove |
| Missing docs | Linter warnings | Add godoc |
| Hardcoded values | Magic numbers | Extract constants |

### Check for Orphans

```bash
# Unused exports
go vet ./...

# Unreachable code
staticcheck ./...

# Missing test files
for f in internal/cli/intents/*/*.go; do
  [[ "$f" != *"_test.go" ]] && [[ ! -f "${f%.go}_test.go" ]] && echo "No test: $f"
done
```

## Task Handoff

If you can't complete a task, document clearly:

```markdown
## Handoff Notes

### Completed
- [x] Basic filter implementation
- [x] Unit tests for filter logic

### Remaining
- [ ] Integration with navigation state
- [ ] E2E tests

### Blockers
- Need decision on filter persistence approach

### Notes
- See `timeline_filter.go:45` for the integration point
- Consider using context for state
```

## Preventing Incomplete Tasks

### Start with End in Mind
Before coding, define "done":
- What are the acceptance criteria?
- What tests need to pass?
- What documentation is needed?

### Track Progress (MANDATORY)
Use task checklist as you work - update IMMEDIATELY after each step:

```
[x] Criterion 1
    COMPLETED: Implemented filter modal with date picker

[ ] Criterion 2
    IN PROGRESS: Adding validation...

[ ] Criterion 3
```

**NEVER batch updates.** Mark each item complete as you finish it.

### Don't Context Switch
Finish current task before starting new one.
Half-done tasks create debt.

### Time Box
If taking too long:
1. Complete what you can
2. Document remaining work
3. Create follow-up task

## Skip Reason Requirement (MANDATORY)

When skipping ANY checklist item, document explicitly:

```
[SKIP] Update CHANGELOG
    SKIPPING: Update CHANGELOG
    REASON: Internal refactoring only, no user-visible changes
    IMPACT: None - no external API changes
```

**Format required:**
- `SKIPPING:` What is being skipped
- `REASON:` Why it's being skipped
- `IMPACT:` What the consequences are

**NEVER silently skip items.** Undocumented skips are violations.

## Related skills

- `check-compliance` - Verification checks
- `clean-code` - Quality standards
- `code-reviewer` - Self-review
- `create-task` - Document remaining work
- `checklist-discipline` - Incremental progress tracking
