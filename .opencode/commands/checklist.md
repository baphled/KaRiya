---
description: Show current task progress and remaining work
agent: research
---

Show the current task's progress, completed items, and what's left to do.

## Process

1. **Detect Current Task**
   - Get current branch: `git branch --show-current`
   - Search for task file with matching `Branch:` field in `tasks/tasks-*.md`
   - Fallback: Use worktree directory name (e.g., `task-57`)
   
2. **Read Task File**
   - Parse `## Status` section
   - Extract completed items
   - Extract in-progress items
   - Extract statistics
   - Extract Definition of Done

3. **Display Summary**
   - Task number and title
   - Branch name
   - Completion statistics
   - Completed items (concise list)
   - In progress items
   - Definition of Done checklist
   - Estimated % complete

## Output Format

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📋 TASK 57: BDD Test Coverage
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Branch: feature/task-57-bdd-godog-coverage

📊 PROGRESS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Total scenarios: 351
Passing (non-@wip): 134/351 (38%)

✅ COMPLETED
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
• Godog framework setup (v0.15.0)
• Skills management feature (16/16 passing)
• Support infrastructure (hooks, env helpers)
... (show first 10, then "and X more")

🚧 IN PROGRESS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
• Implement remaining @wip scenarios incrementally

📝 DEFINITION OF DONE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
☐ All @happy paths passing
☐ All @sad paths passing
☐ All @inference scenarios passing
☐ All @enrichment scenarios passing
☐ make bdd 100% pass

Estimated completion: ~38%
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## Implementation Notes

- Use the `task-tracker` skill to locate and parse the task file
- If task file not found, show helpful error with search locations
- If statistics are stale (no update in last commit), warn the user
- Keep output concise but informative
