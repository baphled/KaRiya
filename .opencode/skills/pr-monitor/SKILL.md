---
name: pr-monitor
description: Monitor PR for CI status, reviews, and coordinate response workflow
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
  always_active: false
---

## What I do

Orchestrate the PR monitoring and review response workflow. I watch for CI failures, collect review comments, and coordinate the response process.

## When to use me

Use this skill after creating a PR to:
- Wait for and check CI status
- Collect and categorise review comments
- Coordinate evidence-based responses

## PR Monitoring Workflow

### 1. Check CI Status

```bash
# Get PR number from current branch
PR_NUM=$(gh pr view --json number -q '.number')

# Check CI status
gh pr checks $PR_NUM
```

Wait for checks to complete:
```bash
# Wait for all checks (with timeout)
gh pr checks $PR_NUM --watch --fail-fast
```

### 2. Handle CI Failures

If CI fails:
1. **Identify failing check** - Which job failed?
2. **Get failure details** - `gh run view <run-id> --log-failed`
3. **Fix locally** - Run same checks: `make ci-local`
4. **Push fix** - Don't create new PR, push to same branch

```bash
# View failed run logs
gh run list --limit 5
gh run view <run-id> --log-failed

# Fix and push
make ci-local  # Verify fix locally
git push
```

### 3. Collect Review Comments

```bash
# Get all review comments
gh api repos/:owner/:repo/pulls/$PR_NUM/comments

# Get review threads
gh pr view $PR_NUM --comments

# Get specific review
gh api repos/:owner/:repo/pulls/$PR_NUM/reviews
```

### 4. Categorise Each Comment

For each review comment, determine:

| Category | Indicators | Response |
|----------|------------|----------|
| **Bug/Defect** | "This will break...", "Missing null check" | Accept immediately |
| **Security** | "SQL injection", "XSS", "secrets" | Accept immediately |
| **Architecture** | "Violates layer...", "Wrong pattern" | Evaluate against docs |
| **Performance** | "O(n^2)", "Memory leak", "Slow" | Require benchmarks |
| **Style** | "Rename to...", "Move this..." | Follow project standard |
| **Unclear** | Vague feedback | Ask clarifying questions |

### 5. Response Strategy

```
For each comment:
    │
    ├─ Use `evaluate-change-request` skill
    │   └─ Determine validity and category
    │
    ├─ If accepting:
    │   ├─ Implement the change
    │   ├─ Write regression test if bug
    │   └─ Use `respond-to-review` to confirm
    │
    ├─ If challenging:
    │   ├─ Use `prove-correctness` skill
    │   ├─ Use `justify-decision` skill
    │   └─ Use `respond-to-review` with evidence
    │
    └─ If unclear:
        └─ Use `respond-to-review` to ask questions
```

## Monitoring Commands

### Quick Status Check
```bash
# PR status summary
gh pr view --json state,reviews,statusCheckRollup

# Just CI status
gh pr checks

# Review status
gh pr view --json reviews -q '.reviews[] | "\(.author.login): \(.state)"'
```

### Watch for Updates
```bash
# Poll for changes (every 30 seconds)
watch -n 30 'gh pr checks && gh pr view --json reviews'
```

## Integration with Other Skills

| Skill | When to Use |
|-------|-------------|
| `evaluate-change-request` | For each review comment |
| `prove-correctness` | When challenging feedback with tests |
| `justify-decision` | When explaining architectural choices |
| `respond-to-review` | For all review responses |
| `trade-off-analysis` | When reviewer suggests alternatives |
| `critical-thinking` | Throughout evaluation process |
| `epistemic-rigor` | Distinguish facts from assumptions |

## Anti-Patterns

**DO NOT:**
- Blindly accept all feedback without evaluation
- Dismiss feedback without evidence
- Respond emotionally to criticism
- Delay responses (aim for same-day)
- Make changes without understanding why

**DO:**
- Evaluate each comment critically
- Provide evidence for disagreements
- Ask clarifying questions when unclear
- Thank reviewers for valid catches
- Learn from feedback for future PRs

## Progress Reporting

**IMPORTANT: Keep the user informed throughout the process.**

### Initial Status Report

When starting review response:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
PR STATUS: #[NUM] - [Title]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Branch: [branch] -> next
CI Status: [Passing/Failing/Pending]

## Review Comments: [Total]
| # | Reviewer | Category | Preview |
|---|----------|----------|---------|
| 1 | @alice | Bug | "Missing null check..." |
| 2 | @bob | Style | "Consider renaming..." |
| 3 | @copilot | Security | "Potential injection..." |

## Plan
Will address comments in order of severity:
1. Security issues first
2. Bugs second  
3. Architecture/Style last

Starting with comment #3 (Security)...
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### Per-Comment Progress

After addressing each comment:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
COMMENT #[N] COMPLETE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Reviewer: @[username]
Category: [Bug/Security/Style/etc.]
Feedback: "[Brief quote]"

Decision: [Accepted/Challenged/Clarified/Deferred]
Reasoning: [Why this decision]
Action: [What was done]
Commit: [SHA or N/A]
Thread: [Resolved/Pending reviewer response]

Progress: [N]/[Total] comments addressed
Remaining: [List remaining comments]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### Resolving Threads

After responding, resolve the GitHub thread:

```bash
# Get unresolved thread IDs
gh api graphql -f query='
  query($owner: String!, $repo: String!, $pr: Int!) {
    repository(owner: $owner, name: $repo) {
      pullRequest(number: $pr) {
        reviewThreads(first: 100) {
          nodes {
            id
            isResolved
            comments(first: 1) { nodes { body author { login } } }
          }
        }
      }
    }
  }
' -f owner=OWNER -f repo=REPO -f pr=$PR_NUM

# Resolve a thread
gh api graphql -f query='
  mutation($threadId: ID!) {
    resolveReviewThread(input: {threadId: $threadId}) {
      thread { isResolved }
    }
  }
' -f threadId="$THREAD_ID"
```

## Merge Readiness Summary

When all comments are addressed, provide final summary:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
PR READY FOR MERGE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PR #[NUM]: [Title]
URL: [PR URL]
Branch: [branch] -> next

## Review Summary

| Metric | Count |
|--------|-------|
| Total comments | [N] |
| Accepted | [X] |
| Challenged (with evidence) | [Y] |
| Clarified | [Z] |
| Deferred to issues | [W] |

## Resolution Details

| # | Reviewer | Decision | Resolution |
|---|----------|----------|------------|
| 1 | @alice | Accept | Fixed in abc123 |
| 2 | @bob | Challenge | Kept current (see evidence) |
| 3 | @copilot | Accept | Fixed in def456 |

## Commits Added

| SHA | Description |
|-----|-------------|
| abc123 | Fix null check (comment #1) |
| def456 | Add input validation (comment #3) |

## Thread Status

- Total threads: [N]
- Resolved: [X]
- Pending reviewer response: [Y]
- Unresolved (explain): [Z]

## CI Status

| Check | Status |
|-------|--------|
| Build | [Pass/Fail] |
| Tests | [Pass/Fail] |
| Coverage | [XX%] |
| Lint | [Pass/Fail] |

## Pre-Merge Checklist

- [x] All comments addressed and responded to
- [x] Conversations resolved (or awaiting reviewer)
- [x] CI passing
- [x] Coverage maintained/improved
- [x] Re-review requested from reviewers
- [ ] Approved by reviewer (waiting)

## Next Steps

[One of:]
- Ready for merge pending approval
- Waiting for @[reviewer] response on comment #[N]
- CI failing - needs fix before merge
- Coverage dropped - add tests before merge

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## Example Session

```bash
# 1. Create PR
gh pr create --base next --title "feat(intents): add skill management"

# 2. Request Copilot review
PR_NUM=$(gh pr view --json number -q '.number')
gh api repos/:owner/:repo/pulls/$PR_NUM/requested_reviewers -X POST -f "reviewers[]=copilot"

# 3. Wait for CI
gh pr checks $PR_NUM --watch

# 4. Check for reviews
gh pr view $PR_NUM --comments

# 5. For each comment, evaluate and respond
# (Use evaluate-change-request, prove-correctness, respond-to-review skills)
# Report progress after each comment

# 6. Resolve threads as comments are addressed
# Use GraphQL API to resolve review threads

# 7. Provide merge readiness summary
# Use template above

# 8. Request re-review
gh pr edit $PR_NUM --add-reviewer <reviewer>
```

## Related Skills

- `create-pr` - Create the PR initially
- `evaluate-change-request` - Critically assess feedback
- `prove-correctness` - Write tests as evidence
- `justify-decision` - Explain architectural choices
- `respond-to-review` - Craft responses
- `critical-thinking` - Rigorous analysis
