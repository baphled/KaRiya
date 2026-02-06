---
name: create-bug
description: Create and document bug reports with proper structure for tracking and fixing
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Help create well-structured bug reports that enable efficient debugging and resolution.

## When to use me

Use this skill when:
- Encountering unexpected behavior
- A test reveals a defect
- Users report issues
- Reviewing code and finding bugs

## Quick Start

```bash
make new-bug BUG="description of the issue"
```

This creates a new bug file in `docs/tasks/` with auto-numbered ID.

## Bug Report Structure

```markdown
# BUG-XXX: [Descriptive Title]

## Summary

Brief description of the bug (1-2 sentences).
What's broken and when does it happen?

## Steps to Reproduce

1. Start with a clean state
2. Perform specific action
3. Observe the failure

## Expected Behavior

What should happen when following these steps.

## Actual Behavior

What actually happens (include error messages).

## Environment

- OS: Linux/macOS/Windows
- Go version: go version output
- Branch: current branch name

## Screenshots/Logs

```
Paste relevant error messages, stack traces, or logs
```

## Severity

- [ ] Critical - Application crash or data loss
- [ ] High - Major feature completely broken
- [ ] Medium - Feature partially broken, workaround exists
- [ ] Low - Minor issue, cosmetic problem

## Notes

Additional context, suspected cause, or initial investigation.
```

## Severity Guidelines

| Severity | Definition | Response Time |
|----------|------------|---------------|
| **Critical** | App crashes, data loss, security issue | Immediate |
| **High** | Major feature broken, no workaround | Within 24 hours |
| **Medium** | Feature degraded, workaround exists | Within sprint |
| **Low** | Minor/cosmetic, doesn't affect function | Backlog |

## Good Bug Titles

```
# BAD - Vague
BUG-42: Something is broken
BUG-43: Error on click

# GOOD - Specific and actionable
BUG-42: Timeline crashes when filtering by empty date range
BUG-43: Event form loses data when pressing Escape twice
```

## Bug Fix Workflow

1. **Create Bug Report** - Document the issue
2. **Write Regression Test** - Test that fails due to the bug
3. **Fix the Bug** - Minimal change to pass the test
4. **Verify** - All tests pass, no regressions
5. **Close Bug** - Update status in bug report

## Linking to Commits

When fixing a bug, reference it in the commit:

```
fix(cli): prevent crash on empty date range

Timeline view now validates date range before filtering.

Fixes: BUG-42
```

## Investigation Template

When investigating a bug, document your findings:

```markdown
## Investigation Log

### Initial Analysis
- [ ] Reproduced locally
- [ ] Identified affected component
- [ ] Found root cause

### Root Cause
[Describe the actual cause]

### Proposed Fix
[Describe the solution approach]

### Impact Assessment
- [ ] Other areas affected?
- [ ] Breaking changes?
- [ ] Migration needed?
```

## Related skills

- `debug-test` - Debug failing tests
- `tdd-workflow` - Write regression test first
- `create-task` - When fix requires significant work
