---
name: create-pr
description: Create a pull request targeting the next branch following KaRiya branching strategy
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Create a pull request following KaRiya's branching strategy.

## When to use me

Use this skill when your feature branch is ready for review.

## Pre-PR Validation

```bash
make pre-pr
```

This validates:
- You're not on `main` or `next` branch
- All CI checks pass
- Branch is ready for PR

## Branching Rules

| Rule | Details |
|------|---------|
| Feature branches | Always work on feature branches |
| PR target | Always target `next` (NEVER `main`) |
| `main` updates | Only via `next -> main` release PRs |

## Creating the PR

### 1. Ensure Branch is Up to Date

```bash
git fetch origin
git rebase origin/next
```

### 2. Push to Remote

```bash
git push -u origin HEAD
```

### 3. Create PR and Request Copilot Review

```bash
# Create PR and capture the URL
PR_URL=$(gh pr create --base next --title "type(scope): description" --body "$(cat <<'EOF'
## Summary

- Brief description of changes
- Key implementation details

## Changes

- List of specific changes made

## Testing

- How the changes were tested
- Coverage information

## Checklist

- [ ] Tests pass (`make test`)
- [ ] Coverage >= 95% for modified packages
- [ ] Compliance check passes (`make check-compliance`)
- [ ] Documentation updated if needed
EOF
)")

# Extract PR number and request Copilot review
PR_NUM=$(echo "$PR_URL" | grep -oE '[0-9]+$')
gh api repos/:owner/:repo/pulls/$PR_NUM/requested_reviewers -X POST -f "reviewers[]=copilot"

echo "PR created: $PR_URL"
echo "Copilot review requested"
```

**Note**: Copilot review is automatically requested for all PRs. This provides immediate AI-powered feedback on code quality, security, and best practices.

## PR Title Format

Same as commit message:
```
type(scope): short description
```

Examples:
- `feat(intents): add burst management workflow`
- `fix(screens): resolve modal overlay z-index issue`
- `refactor(forms): migrate to forms package`

## PR Body Template

```markdown
## Summary

1-3 bullet points describing the change

## Changes

- Detailed list of changes
- New files added
- Files modified

## Testing

- Test coverage percentage
- Key scenarios tested

## Screenshots (if UI changes)

Before/After screenshots if applicable
```

## Post-Creation: Monitor and Respond

After creating the PR, you MUST monitor for CI failures and reviews.

### 4. Wait for CI Checks

```bash
# Watch CI status (blocks until complete)
gh pr checks $PR_NUM --watch --fail-fast
```

If CI fails:
```bash
# View failure details
gh run list --limit 5
gh run view <run-id> --log-failed

# Fix locally
make ci-local

# Push fix (same branch, same PR)
git push
```

### 5. Monitor for Reviews

```bash
# Check review status
gh pr view $PR_NUM --json reviews,comments

# Get detailed review comments
gh api repos/:owner/:repo/pulls/$PR_NUM/comments
```

### 6. Respond to Review Feedback

**IMPORTANT: Do not blindly accept all feedback.**

For each review comment:

1. **Evaluate** - Use `evaluate-change-request` skill
   - Is this objectively correct (bug, security)?
   - Is there evidence provided?
   - Does it align with project standards?

2. **Decide** - Based on evaluation:
   - **Accept**: Implement and respond "Fixed in [SHA]"
   - **Challenge**: Use `prove-correctness` and `justify-decision`
   - **Clarify**: Ask specific questions
   - **Defer**: Create issue for out-of-scope work

3. **Respond** - Use `respond-to-review` skill
   - Always respond to every comment
   - Be professional and specific
   - Provide evidence when disagreeing

### 7. Request Re-review

After addressing all feedback:
```bash
# Summary comment
gh pr comment $PR_NUM --body "Addressed all feedback. Ready for re-review.

| Comment | Resolution |
|---------|------------|
| [Issue 1] | Fixed in abc123 |
| [Issue 2] | Discussed - keeping current |"

# Request re-review
gh pr edit $PR_NUM --add-reviewer <reviewer>
```

## Review Process

1. **Copilot review** - Automatically requested on PR creation
2. **CI checks must pass** - Fix any failures before review
3. **Evaluate feedback critically** - Don't blindly accept
4. **Respond to all comments** - Accept, challenge, or clarify
5. At least one human approval required
6. All comments resolved (Copilot + human)
7. Squash merge to `next`

### Review Response Workflow

```
Review Comment Received
    │
    ├─ Use `evaluate-change-request` skill
    │   └─ Categorise and assess validity
    │
    ├─ If Valid Bug/Security:
    │   └─ Accept immediately, fix, respond
    │
    ├─ If Questionable:
    │   ├─ Use `prove-correctness` (write tests)
    │   ├─ Use `justify-decision` (document reasoning)
    │   └─ Use `respond-to-review` (craft response)
    │
    └─ If Unclear:
        └─ Ask clarifying questions first
```

### Copilot Review Benefits

- Immediate feedback (no waiting for human availability)
- Security vulnerability detection
- Code quality suggestions
- Best practice recommendations
- Catches common issues before human review

## After Merge

```bash
git checkout next
git pull origin next
git branch -d feature/my-feature
```

## Related Skills

### Pre-PR
- `ai-commit` - Create commits before PR
- `check-compliance` - Validate before PR
- `session-start` - Start new feature work

### Post-PR (Review Response)
- `pr-monitor` - Orchestrate monitoring workflow
- `evaluate-change-request` - Critically assess feedback
- `prove-correctness` - Write tests as evidence
- `justify-decision` - Explain architectural choices
- `respond-to-review` - Craft professional responses
- `trade-off-analysis` - Compare alternatives
