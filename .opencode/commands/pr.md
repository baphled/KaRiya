---
description: Create a pull request targeting next branch
agent: build
---

Create a pull request for the current branch.

Load the `create-pr` skill for PR process.

## Process

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
   gh pr create --base next --title "..." --body "..."
   ```

5. **Return PR URL**

## Additional Context
$ARGUMENTS
