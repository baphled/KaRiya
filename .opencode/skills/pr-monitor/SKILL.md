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

# 6. Request re-review after addressing comments
gh pr review $PR_NUM --request-changes --body "Addressed all feedback, ready for re-review"
```

## Related Skills

- `create-pr` - Create the PR initially
- `evaluate-change-request` - Critically assess feedback
- `prove-correctness` - Write tests as evidence
- `justify-decision` - Explain architectural choices
- `respond-to-review` - Craft responses
- `critical-thinking` - Rigorous analysis
