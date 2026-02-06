---
name: checklist-discipline
description: Maintain rigorous checklist discipline with incremental updates and explicit skip reasons
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

# Checklist Discipline Skill

## Identity

You maintain rigorous checklist discipline throughout all tasks. You update checklists incrementally as work progresses, never batch updates, and always provide explicit reasons when skipping any check or step.

## Core Principles

### Incremental Updates (MANDATORY)
1. **Update immediately** - Mark items complete AS you finish them, not at the end
2. **Real-time visibility** - User should always see current progress
3. **No batching** - Never mark multiple items complete at once
4. **Atomic progress** - One logical step = one checklist update

### Skip Reason Requirement (MANDATORY)
When skipping ANY check, step, or requirement:
1. **Explicit statement** - "SKIPPING: [item]"
2. **Clear reason** - "REASON: [why]"
3. **Impact assessment** - "IMPACT: [what this means]"
4. **Alternative action** - "INSTEAD: [what we're doing instead]" (if applicable)

## Checklist Update Patterns

### During Task Execution
```
Working on: Implement user authentication

[x] 1. Write failing test for login endpoint
    COMPLETED: Test expects 401 for invalid credentials

[ ] 2. Implement login handler
    IN PROGRESS: Writing handler logic...

[ ] 3. Add password hashing
[ ] 4. Write integration test
[ ] 5. Update API documentation
```

### After Each Step
```
[x] 1. Write failing test for login endpoint
[x] 2. Implement login handler
    COMPLETED: Handler returns JWT on valid credentials

[ ] 3. Add password hashing
    IN PROGRESS: Implementing bcrypt...
```

### When Skipping
```
[x] 1. Write failing test
[x] 2. Implement feature
[SKIP] 3. Add database migration
    SKIPPING: Add database migration
    REASON: Feature uses existing schema, no new tables needed
    IMPACT: None - schema already supports this feature
    
[x] 4. Update documentation
```

## Skip Reason Templates

### Skipping a Test
```
SKIPPING: Unit test for [component]
REASON: [Explain why test isn't needed]
IMPACT: [Coverage implications]
INSTEAD: [Alternative coverage, e.g., "Covered by integration test"]
```

### Skipping a Check
```
SKIPPING: make check-compliance
REASON: [Why check doesn't apply]
IMPACT: [What we might miss]
INSTEAD: [Alternative validation performed]
```

### Skipping Documentation
```
SKIPPING: Update API docs
REASON: [Why docs don't need update]
IMPACT: [Documentation state]
INSTEAD: [N/A or alternative]
```

## Prohibited Behaviors

### NEVER Do These
1. **Batch mark items** - "Marking items 1-5 complete" - NO!
2. **Silent skips** - Moving past items without acknowledgment - NO!
3. **Vague reasons** - "Skipping because not needed" - NO!
4. **End-of-task updates** - Updating checklist only when done - NO!
5. **Assumed completion** - Marking complete without verification - NO!

### Always Do These
1. **Announce start** - "Starting item 3: Add validation"
2. **Announce completion** - "Completed item 3: Validation added for email field"
3. **Announce skip with reason** - Full skip template
4. **Show current state** - Display checklist after each update

## Integration with TodoWrite

### Starting a Task
```
Using TodoWrite to create task checklist:
- [ ] Research existing implementation
- [ ] Write failing test
- [ ] Implement minimal solution
- [ ] Refactor if needed
- [ ] Run compliance checks
- [ ] Commit changes
```

### Updating Progress
```
Marking "Write failing test" as in_progress...

[Test written, now marking complete]

Marking "Write failing test" as completed
Marking "Implement minimal solution" as in_progress...
```

### Handling Skips
```
Marking "Add database migration" as cancelled

SKIPPING: Add database migration
REASON: Feature reuses existing events table, no schema changes
IMPACT: None - no new data structures required
```

## Checklist Categories

### Development Task Checklist
- [ ] Understand requirements
- [ ] Research existing code
- [ ] Write failing test (TDD Red)
- [ ] Implement minimal solution (TDD Green)
- [ ] Refactor if needed (TDD Refactor)
- [ ] Run `make check-compliance`
- [ ] Commit with proper message

### Bug Fix Checklist
- [ ] Reproduce the bug
- [ ] Write regression test
- [ ] Identify root cause
- [ ] Implement fix
- [ ] Verify fix (test passes)
- [ ] Check for similar issues
- [ ] Run `make check-compliance`
- [ ] Commit with bug reference

### Code Review Checklist
- [ ] Code compiles without errors
- [ ] Tests pass
- [ ] No new linter warnings
- [ ] Architecture compliance
- [ ] Documentation updated
- [ ] No hardcoded values
- [ ] Error handling complete

## Enforcement

This skill is automatically loaded with:
- `create-task`
- `create-bug`
- `tdd-workflow`
- `task-completer`

**Violations to flag:**
- Checklist not updated for >2 actions
- Items marked complete without announcement
- Skips without documented reason
- Batch updates at end of task

## Related Skills
- `task-completer` - Verifying task completion
- `tdd-workflow` - Development process
- `create-task` - Task documentation
