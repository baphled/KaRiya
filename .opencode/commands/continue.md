---
description: Continue work from where you left off - assess state and resume
agent: build
---

Continue work from where I left off.

## Skills to Load (MANDATORY)

- `pre-action` - Decision framework before ANY action (always active)
- `software-engineer` - Orchestrates technical work
- `task-completer` - Check what's done and what remains

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

4. **Apply Pre-Action Framework**
   Before taking any action based on findings:
   - STOP: What am I about to do?
   - THINK: What do I KNOW vs ASSUME?
   - INVESTIGATE: If tests fail, WHY? Is it intentional?
   - CONFIDENCE: Am I VERIFIED or just ASSUMED?
   - ASK if uncertain

5. **Identify Next Step**
   Based on state:
   - **Uncommitted changes** → Review, test, commit
   - **Failing tests** → Investigate WHY before fixing
   - **Partial implementation** → Continue TDD cycle
   - **Ready for commit** → Run compliance, commit
   - **Ready for PR** → Create PR

6. **Resume Work**
   Load appropriate skills based on what needs to be done next.

## State Detection

| State | Indicator | Action |
|-------|-----------|--------|
| Mid-implementation | Uncommitted code, tests failing | Investigate, then continue TDD |
| Tests passing | Green tests, uncommitted | Refactor or commit |
| Ready to commit | Clean, tested code | Run compliance, commit |
| Needs review | Committed, not pushed | Self-review, push |
| Ready for PR | Pushed to remote | Create PR |
| Blocked | Failing checks | Investigate cause first |

## Output

Report:
- Current state summary
- What was being worked on
- What's completed
- **Pre-action assessment** (KNOW/ASSUME/UNKNOWN)
- Recommended next action
- Then proceed (or ASK if uncertain)
