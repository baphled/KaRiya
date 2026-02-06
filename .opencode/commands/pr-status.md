---
description: Check PR status with interactive options for next actions
agent: build
---

Check the current PR status and provide interactive options for what to do next.

Load the `pr-monitor` skill for monitoring workflow.

## Output Format

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
PR STATUS: #[NUM] - [Title]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Branch: [branch] -> next
URL: [PR URL]

## CI Status
| Check | Status | Details |
|-------|--------|---------|
| Build | [PASS/FAIL/PENDING] | |
| Tests | [PASS/FAIL/PENDING] | |
| Lint | [PASS/FAIL/PENDING] | |

## Review Status
- Approvals: [N]
- Changes requested: [N]
- Comments: [X] total, [Y] unaddressed
- Threads: [X] resolved, [Y] pending

## Unaddressed Comments
[List each with brief preview]

## Blockers
[List any blocking issues]

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## Interactive Options

After showing status, present options based on current state:

### If CI Failing
```
What would you like to do?
1. View failure details and diagnose
2. Fix the failing check
3. Ignore and continue (not recommended)
```

### If Unaddressed Comments
```
What would you like to do?
1. Address all comments (/respond-review)
2. Address specific comment (#N)
3. View comment details
4. Start continuous monitoring (/pr-poll)
```

### If All Good (Ready for Merge)
```
What would you like to do?
1. Generate merge readiness summary (/pr-ready)
2. Start continuous monitoring (/pr-poll)
3. Nothing - just checking
```

### If Waiting for Reviewer
```
What would you like to do?
1. Ping reviewer for response
2. Start continuous monitoring (/pr-poll)
3. Nothing - just checking
```

## Process

1. **Get PR Data**
   ```bash
   PR_NUM=$(gh pr view --json number -q '.number')
   gh pr view $PR_NUM --json title,url,state,reviews,statusCheckRollup
   ```

2. **Check CI Status**
   ```bash
   gh pr checks $PR_NUM
   ```

3. **Get Review Comments**
   ```bash
   gh api repos/:owner/:repo/pulls/$PR_NUM/comments
   ```

4. **Get Thread Status**
   ```bash
   gh api graphql -f query='...' # Get resolved/unresolved threads
   ```

5. **Display Status Summary**
   Use format above.

6. **Present Options**
   Based on current state, show relevant options.

7. **Handle Selection**
   Execute chosen action or hand off to appropriate command.

## Quick Actions

For common follow-ups:

| State | Quick Action |
|-------|--------------|
| CI failing | `/pr-status fix` - Jump to fix |
| Comments pending | `/pr-status respond` - Jump to respond |
| Ready for merge | `/pr-status ready` - Show merge summary |
| Monitor | `/pr-status poll` - Start polling |

## Skills Used
- `pr-monitor` - Monitoring workflow
- `evaluate-change-request` - Assess feedback

## Arguments
$ARGUMENTS

Examples:
- `/pr-status` - Show status with options
- `/pr-status 162` - Check specific PR
- `/pr-status fix` - Jump to fixing CI
- `/pr-status respond` - Jump to responding to comments
- `/pr-status poll` - Start continuous monitoring
