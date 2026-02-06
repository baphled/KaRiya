---
description: Generate merge readiness summary for current PR
agent: build
---

Generate a comprehensive merge readiness summary for the current PR.

Use this after addressing all review comments to provide a clear status report.

Load skills:
- `pr-monitor` - For merge readiness template
- `respond-to-review` - For summary format

## Process

### 1. Gather PR Data

```bash
PR_NUM=$(gh pr view --json number -q '.number')
PR_TITLE=$(gh pr view --json title -q '.title')
PR_URL=$(gh pr view --json url -q '.url')
BRANCH=$(git rev-parse --abbrev-ref HEAD)
```

### 2. Check CI Status

```bash
gh pr checks $PR_NUM
```

### 3. Get Review Thread Status

```bash
# Count resolved vs unresolved threads
gh api graphql -f query='
  query($owner: String!, $repo: String!, $pr: Int!) {
    repository(owner: $owner, name: $repo) {
      pullRequest(number: $pr) {
        reviewThreads(first: 100) {
          nodes {
            id
            isResolved
            comments(first: 1) { 
              nodes { body author { login } } 
            }
          }
        }
      }
    }
  }
'
```

### 4. Get Commits Since PR Creation

```bash
git log origin/next..HEAD --oneline
```

### 5. Generate Summary

Output the following format:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
PR MERGE READINESS SUMMARY
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PR #[NUM]: [Title]
URL: [PR URL]
Branch: [branch] -> next

## Review Summary

| Metric | Count |
|--------|-------|
| Total comments received | [N] |
| Accepted (changes made) | [X] |
| Challenged (with evidence) | [Y] |
| Clarified (questions answered) | [Z] |
| Deferred (to future issues) | [W] |

## Resolution Details

| # | Reviewer | Category | Decision | Resolution |
|---|----------|----------|----------|------------|
| 1 | @xxx | Bug | Accept | Fixed in abc123 |
| 2 | @xxx | Style | Challenge | Kept (see evidence) |
| ... | ... | ... | ... | ... |

## Commits Added During Review

| SHA | Description |
|-----|-------------|
| abc123 | Fix null check |
| def456 | Add input validation |
| ... | ... |

## Thread Status

- Total threads: [N]
- Resolved: [X]
- Pending reviewer response: [Y] (list which)
- Unresolved: [Z] (explain why)

## CI Status

| Check | Status | Details |
|-------|--------|---------|
| Build | [PASS/FAIL] | |
| Tests | [PASS/FAIL] | [X] passed, [Y] failed |
| Coverage | [XX%] | [Maintained/Improved/Dropped] |
| Lint | [PASS/FAIL] | |

## Pre-Merge Checklist

- [ ] All comments addressed
- [ ] All accepted changes committed
- [ ] All threads resolved (or justified)
- [ ] CI passing
- [ ] Coverage maintained (>= 95% modified packages)
- [ ] Re-review requested

## Blockers (if any)

[List any issues preventing merge]

## Next Steps

[One of:]
- READY: Awaiting final approval from @[reviewer]
- BLOCKED: [Describe blocker and action needed]
- WAITING: Reviewer response needed on comment #[N]

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### 6. Post Summary to PR (optional)

If requested, post summary as PR comment:

```bash
gh pr comment $PR_NUM --body "[summary]"
```

## Skills Used
- `pr-monitor` - Merge readiness template
- `respond-to-review` - Summary format

## Options
$ARGUMENTS

Examples:
- `/pr-ready` - Generate summary only
- `/pr-ready post` - Generate and post to PR
- `/pr-ready verbose` - Include full comment details
