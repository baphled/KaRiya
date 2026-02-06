---
description: Check PR status including CI and review feedback
agent: build
---

Check the current PR status and collect review feedback.

Load the `pr-monitor` skill for monitoring workflow.

## Process

1. **Get PR Number**
   ```bash
   PR_NUM=$(gh pr view --json number -q '.number')
   ```

2. **Check CI Status**
   ```bash
   gh pr checks $PR_NUM
   ```

3. **Get Review Status**
   ```bash
   gh pr view $PR_NUM --json state,reviews,statusCheckRollup
   ```

4. **Collect Review Comments**
   ```bash
   gh pr view $PR_NUM --comments
   gh api repos/:owner/:repo/pulls/$PR_NUM/comments
   ```

5. **Summarise Status**
   Report:
   - CI status (passing/failing/pending)
   - Number of reviews and their state
   - Pending comments requiring response
   - Blocking issues

6. **If CI Failing**
   ```bash
   gh run list --limit 5
   gh run view <run-id> --log-failed
   ```
   Identify the failure and suggest fix.

7. **If Reviews Pending**
   For each unresolved comment:
   - Quote the feedback
   - Suggest response strategy (accept/challenge/clarify)
   - Reference `evaluate-change-request` skill

## Skills Used
- `pr-monitor` - Monitoring workflow
- `evaluate-change-request` - Assess feedback

## PR Number (optional)
$ARGUMENTS
