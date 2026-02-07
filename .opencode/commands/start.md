---
description: Start a KaRiya development session with environment validation
agent: build
---

Start a new KaRiya development session.

## Skills to Load (MANDATORY)

Load these skills - they govern ALL subsequent actions:
- `pre-action` - Decision framework before ANY action (always active)
- `token-efficiency` - Concise, efficient communication (always active)
- `software-engineer` - Orchestrates technical work
- `session-start` - Session protocol
- `task-tracker` - Task file synchronization

## Process

1. Run environment validation:
```bash
make session-start
```

2. If it fails, help me fix the issues before proceeding.

3. Update task file with branch information:
   - Detect current task from worktree directory name
   - Find matching `tasks/tasks-{N}-*.md` file
   - Add `Branch: <current-branch>` after `## Status` if missing
   - Verify task file is synchronized

4. Acknowledge the critical rules:
   - Feature branches only (never commit to next/main)
   - TDD workflow (test first)
   - Use `make ai-commit` for commits
   - Run `make check-compliance` before and after tasks
   - **Pre-action framework before ANY action**
   - **Update task file after completing work**

$ARGUMENTS
