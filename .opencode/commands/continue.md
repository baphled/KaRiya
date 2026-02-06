---
description: Continue work from where you left off - assess state and resume
agent: build
---

Continue work from where I left off.

Load the `software-engineer` skill to assess and resume appropriately.
Load the `task-completer` skill to check what's done and what remains.

## Context
$ARGUMENTS

## Process

1. **Assess Current State**
   ```bash
   # What's changed?
   git status
   
   # What was I working on?
   git log -3 --oneline
   
   # Any uncommitted work?
   git diff --stat
   ```

2. **Check Task Progress**
   - What was the original task?
   - What's been completed?
   - What remains?

3. **Verify Environment**
   ```bash
   # Quick health check
   make test
   ```

4. **Identify Next Step**
   Based on state:
   - **Uncommitted changes** → Review, test, commit
   - **Failing tests** → Debug and fix
   - **Partial implementation** → Continue TDD cycle
   - **Ready for commit** → Run compliance, commit
   - **Ready for PR** → Create PR

5. **Resume Work**
   Load appropriate skills based on what needs to be done next.

## State Detection

| State | Indicator | Action |
|-------|-----------|--------|
| Mid-implementation | Uncommitted code, tests failing | Continue TDD green phase |
| Tests passing | Green tests, uncommitted | Refactor or commit |
| Ready to commit | Clean, tested code | Run compliance, commit |
| Needs review | Committed, not pushed | Self-review, push |
| Ready for PR | Pushed to remote | Create PR |
| Blocked | Failing checks | Fix issues first |

## Output

Report:
- Current state summary
- What was being worked on
- What's completed
- Recommended next action
- Then proceed with that action
