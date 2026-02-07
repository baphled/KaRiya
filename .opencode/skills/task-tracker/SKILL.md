# Task Tracker Skill

Keep task files synchronized with actual progress throughout development.

## When to Use

**Always active** - This skill is part of the foundation workflow for task-based development.

Load automatically in `software-engineer` skill.

## Core Responsibilities

### 1. Task File Location

Task files follow the pattern: `tasks/tasks-{N}-*.md`

Detect current task by:
1. Get current branch: `git branch --show-current`
2. Look for `Branch:` field in task files
3. Match current branch to task file
4. Fallback: Use worktree directory name (e.g., `task-57`)

### 2. Update Task File on Progress

**TRIGGER: Whenever you complete a checkbox item**

Update the task file immediately:
- [ ] → [x] when item is done
- Add completion notes if relevant
- Update statistics/counters
- Keep `## Status` section current

**Example:**
```markdown
## Status

### Completed
- [x] Godog framework setup (v0.15.0)
- [x] Skills management feature (16/16 passing) ← Just completed

### Statistics
- **Total scenarios**: 351
- **Passing (non-@wip)**: 134 (38%) ← Update count
```

### 3. Validate Task File Accuracy

**TRIGGER: Before committing, before creating PR, on request**

Check that task file matches reality:

1. **Verify checkboxes reflect actual state**
   - Run tests/checks to confirm completed items
   - Flag items marked done but not actually working
   
2. **Update statistics**
   - Test pass counts
   - Coverage percentages
   - File counts
   
3. **Check Definition of Done**
   - List remaining DoD items
   - Estimate completion percentage

### 4. Prompt for Updates

If you notice a mismatch between task file and reality:

```
⚠️  Task file may be out of sync:
- Task file shows 16/16 skills tests passing
- Actual: 14/16 passing (2 failures detected)

Update task file? (Y/n)
```

## Task File Format

### Required Fields

```markdown
# Task {N}: {Title}

## Summary
Brief description of the task

## Status

### Branch
`feature/task-{N}-{slug}`

### Completed
- [x] Item 1
- [x] Item 2

### In Progress
- [ ] Item 3

### Statistics
- **Metric 1**: X/Y (Z%)
- **Metric 2**: A/B

## Definition of Done
- [ ] Criteria 1
- [ ] Criteria 2
```

### Branch Field

When starting work on a task:
1. Create feature branch
2. Add `### Branch` field under `## Status`
3. Record the branch name

**Example:**
```markdown
## Status

### Branch
`feature/task-57-bdd-godog-coverage`

### Completed
...
```

## Workflow Integration

### At Task Start
1. Read task file
2. Check for `### Branch` field
3. If missing, add it with current branch
4. Commit update

### During Development
1. After completing any checkbox item, update task file
2. Update statistics if changed
3. Commit task file with related code changes

### Before Committing
1. Review task file accuracy
2. Update any stale checkboxes
3. Refresh statistics
4. Include task file in commit if changed

### Before Creating PR
1. Run full validation
2. Ensure all completed items are checked
3. Update Definition of Done progress
4. Confirm statistics are current

## Commands

### `/checklist`
Show current task progress (reads task file)

### Validation Check
```bash
# Check task file matches reality
# (This would be a make target or script)
make validate-task
```

## Anti-Patterns

❌ **DON'T:**
- Mark items complete before they actually work
- Let statistics drift from reality
- Forget to update task file after progress
- Assume task file is current without checking

✅ **DO:**
- Update task file immediately when completing items
- Verify claims with actual test runs
- Keep statistics fresh
- Commit task file updates with related code

## Integration with Other Skills

- **pre-action**: Check task file before starting work
- **memory-keeper**: Store task-specific discoveries in task file notes
- **ai-commit**: Include task file updates in commits
- **pr-ready**: Validate task file before creating PR

## Examples

### Example 1: Completing a Feature

```markdown
Before:
### In Progress
- [ ] Skills management feature

After running tests and confirming 16/16 pass:
### Completed
- [x] Skills management feature (16/16 passing)

### Statistics
- **Passing (non-@wip)**: 134 (38%) ← Updated from 118
```

### Example 2: Starting New Work

```markdown
User: Continue working on task-57

Agent:
1. Read tasks/tasks-57-godog-bdd-testing.md
2. Check Branch field → `feature/task-57-bdd-godog-coverage`
3. Verify current branch matches → ✓
4. Load context from Completed/In Progress sections
5. Continue from last checkpoint
```

### Example 3: Detecting Drift

```markdown
Agent notices:
- Task file: "- [x] All tests passing"
- Reality: 2 test failures in CI

Action:
1. Update task file to reflect actual state
2. Move item back to "In Progress" 
3. Add note about failures
4. Commit correction
```

## Notes

- Task files are the source of truth for what's done
- Keep them honest - don't claim completion prematurely  
- Update immediately, don't batch updates
- Statistics should be verifiable (run a command to confirm)
