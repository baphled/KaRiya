---
description: Continuously monitor PR and handle tasks until cancelled
agent: build
---

Continuously monitor the current PR for CI status, review comments, and handle tasks automatically.

**This command runs until you cancel it.** Use `/pr-status` to check status without continuous monitoring.

Load skills:
- `pr-monitor` - Monitoring workflow
- `evaluate-change-request` - Assess feedback
- `prove-correctness` - Evidence through tests
- `respond-to-review` - Craft responses

## Behaviour

The poll loop:
1. Check CI status
2. Check for new/unaddressed review comments
3. Handle any tasks that arise
4. Report progress
5. Wait and repeat

**Automatic handling:**
- CI failures → Diagnose and suggest/apply fixes
- New review comments → Evaluate and respond
- Resolved threads → Track and report
- Approvals → Report merge readiness

## Process

### 1. Initial Status

```bash
PR_NUM=$(gh pr view --json number -q '.number')
PR_TITLE=$(gh pr view --json title -q '.title')
```

Report initial state:
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
PR POLL STARTED
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
PR #[NUM]: [Title]
Monitoring for: CI status, review comments

Press Ctrl+C or say "stop" to end polling.
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### 2. Poll Loop

Each iteration:

```
[HH:MM:SS] Polling PR #[NUM]...

CI Status:
  - Build: [PASS/FAIL/PENDING]
  - Tests: [PASS/FAIL/PENDING]
  - Lint: [PASS/FAIL/PENDING]

Review Status:
  - Comments: [X] total, [Y] unaddressed
  - Threads: [X] resolved, [Y] pending
  - Approvals: [N]

[Any new activity since last poll]
```

### 3. Handle CI Failures

When CI fails:
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
CI FAILURE DETECTED
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Check: [Check name]
Status: FAILED

Fetching logs...
[Diagnosis]

Suggested fix: [Description]

Applying fix...
[Or: "Needs manual intervention - pausing poll"]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### 4. Handle New Review Comments

When new comment detected:
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
NEW REVIEW COMMENT
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Reviewer: @[username]
File: [path]:[line]
Comment: "[text]"

Evaluating...
Category: [Bug/Security/Architecture/Style]
Decision: [Accept/Challenge/Clarify]

[Taking action...]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### 5. Handle Approval

When PR is approved:
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
PR APPROVED
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Approved by: @[username]

Generating merge readiness summary...
[Use /pr-ready output]

PR is ready for merge!
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### 6. Pause Conditions

Pause polling and ask for input when:
- CI failure requires manual intervention
- Review comment needs clarification from user
- Conflicting reviewer feedback
- Merge conflicts detected

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
POLL PAUSED - INPUT NEEDED
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Reason: [Why paused]

Options:
1. [Option description]
2. [Option description]
3. Continue polling (skip this issue)
4. Stop polling

What would you like to do?
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### 7. Stop Polling

When cancelled or complete:
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
PR POLL STOPPED
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Duration: [X] minutes
Actions taken: [N]

Final Status:
  CI: [Status]
  Comments: [X] addressed, [Y] pending
  Approvals: [N]

To resume: /pr-poll
To check status: /pr-status
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## State Tracking

Track between polls:
- Last CI status (to detect changes)
- Known comment IDs (to detect new ones)
- Resolved thread IDs
- Actions taken this session

## Poll Interval

- Initial: Check immediately
- After activity: Every 30 seconds
- Idle (no changes): Every 2 minutes
- After CI completion: Every 1 minute

## Skills Used

- `pr-monitor` - Overall monitoring
- `evaluate-change-request` - Assess new comments
- `prove-correctness` - Write tests if needed
- `justify-decision` - Explain choices
- `respond-to-review` - Post responses

## Resume After Cancel

Use `/pr-poll` again to resume. State is tracked per PR.

## PR Number (optional)
$ARGUMENTS

Examples:
- `/pr-poll` - Poll current branch's PR
- `/pr-poll 162` - Poll specific PR number
