---
name: create-task
description: Create well-structured development tasks with clear acceptance criteria and technical guidance
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Help create clear, actionable development tasks with proper structure for implementation.

## When to use me

Use this skill when:
- Starting a new feature
- Breaking down a large piece of work
- Planning a refactoring effort
- Creating work items for the backlog

## Quick Start

```bash
make new-feature TASK="description of the feature"
```

This creates a new task file in `docs/tasks/` with auto-numbered ID.

## Task Structure

```markdown
# TASK-XXX: [Descriptive Title]

## Summary

Brief description of the task (1-2 sentences).
What needs to be done and why?

## Acceptance Criteria

- [ ] Criterion 1: Specific, measurable outcome
- [ ] Criterion 2: Another measurable outcome
- [ ] Criterion 3: And another

## Technical Notes

### Files to Modify

- `internal/cli/intents/feature/intent.go` - Add new handler
- `internal/cli/screens/feature/list.go` - Update view

### Dependencies

- Requires TASK-XXX to be completed first
- Needs new database migration

### Patterns to Use

| Need | Use |
|------|-----|
| Table view | `behaviors.TableBehavior[T]` |
| Form | `forms.NewInput()`, etc. |
| Modal | `feedback.Modal` |

Run `make what-to-use NEED="keyword"` for details.

## Testing Requirements

### Unit Tests

- [ ] Test handler returns correct result
- [ ] Test edge case: empty input
- [ ] Test error case: invalid data

### E2E Tests (if new workflow)

- [ ] Happy path: complete flow
- [ ] Cancel flow: user backs out
- [ ] Error flow: service failure

## Definition of Done

- [ ] All acceptance criteria met
- [ ] Tests written FIRST (TDD)
- [ ] Tests pass with >= 95% coverage
- [ ] No pattern violations (`make check-patterns`)
- [ ] Architecture check passes (`make check-intent-architecture`)
- [ ] Compliance check passes (`make check-compliance`)
- [ ] Documentation updated (if needed)
- [ ] Committed with `make ai-commit`
```

## Writing Good Acceptance Criteria

Use the **SMART** framework:
- **S**pecific - Clear and unambiguous
- **M**easurable - Can verify completion
- **A**chievable - Reasonable scope
- **R**elevant - Relates to the goal
- **T**ime-bound - Part of this task

```markdown
# BAD - Vague
- [ ] The feature works
- [ ] Users can do things
- [ ] It looks good

# GOOD - Specific and testable
- [ ] User can filter events by date range using the filter modal
- [ ] Empty date range shows validation error "Please select a date range"
- [ ] Filtered list updates within 100ms of applying filter
```

## Task Sizing

| Size | Description | Time Estimate |
|------|-------------|---------------|
| **XS** | Trivial change, one file | < 1 hour |
| **S** | Small feature, few files | 1-4 hours |
| **M** | Medium feature, multiple files | 4-8 hours |
| **L** | Large feature, multiple components | 1-2 days |
| **XL** | Epic, needs breakdown | 3+ days |

**Rule:** If a task is XL, break it down into smaller tasks.

## Breaking Down Large Tasks

```markdown
# EPIC: User Skill Management

## Sub-tasks:
1. TASK-101: Create skill repository and migration
2. TASK-102: Create skill service layer
3. TASK-103: Create skill list screen
4. TASK-104: Create skill edit modal
5. TASK-105: Integrate skill management into timeline
```

## Task Dependencies

Document dependencies clearly:

```markdown
## Dependencies

### Blocks
- TASK-102 (needs repository from this task)

### Blocked By
- TASK-100 (needs database schema changes)

### Related
- BUG-45 (may be affected by this change)
```

## Task States

| State | Meaning |
|-------|---------|
| **Open** | Ready to be worked on |
| **In Progress** | Currently being implemented |
| **In Review** | PR created, awaiting review |
| **Done** | Merged and deployed |
| **Blocked** | Waiting on dependency |
| **Cancelled** | No longer needed |

## Related skills

- `tdd-workflow` - Implement with TDD
- `create-intent` - For new workflow tasks
- `create-screen` - For new UI tasks
- `tech-debt` - For refactoring tasks
