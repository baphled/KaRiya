---
description: Create a pull request targeting next branch
agent: build
---

Create a pull request for the current branch.

Load the `create-pr` skill for full PR workflow including post-creation monitoring.

## Process

### Phase 1: Create PR

1. **Pre-PR Validation**
   ```bash
   make pre-pr
   ```
   This validates branch name and runs all CI checks.

2. **Review Changes Since Branch**
   ```bash
   git log origin/next..HEAD --oneline
   git diff origin/next...HEAD --stat
   ```

3. **Push Branch**
   ```bash
   git push -u origin HEAD
   ```

4. **Create PR**
   Using the changes, create a PR with:
   - Title: `type(scope): description`
   - Summary of changes
   - Testing notes
   - Checklist

   ```bash
   PR_URL=$(gh pr create --base next --title "..." --body "...")
   PR_NUM=$(echo "$PR_URL" | grep -oE '[0-9]+$')
   ```

5. **Request Copilot Review**
   ```bash
   gh api repos/:owner/:repo/pulls/$PR_NUM/requested_reviewers -X POST -f "reviewers[]=copilot"
   ```

### Phase 2: Monitor and Respond

6. **Wait for CI Checks**
   ```bash
   gh pr checks $PR_NUM --watch --fail-fast
   ```
   Fix any failures before proceeding.

7. **Monitor for Reviews**
   ```bash
   gh pr view $PR_NUM --json reviews,comments,statusCheckRollup
   ```

8. **Respond to Feedback**
   - Load `evaluate-change-request` skill for each comment
   - Do NOT blindly accept all feedback
   - Use `prove-correctness` and `justify-decision` when challenging
   - Use `respond-to-review` for all responses

9. **Return PR URL**

## Skills Used
- `create-pr` - Full PR workflow
- `pr-monitor` - Post-creation monitoring
- `evaluate-change-request` - Critical feedback analysis
- `respond-to-review` - Professional responses

## Additional Context
$ARGUMENTS
