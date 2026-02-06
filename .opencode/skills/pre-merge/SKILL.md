---
name: pre-merge
description: Final validation checklist before merging PRs to ensure quality and completeness
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
  always_active: false
---

## What I do

Provide a comprehensive final validation checklist before merging PRs. I ensure nothing is missed and the codebase remains healthy after merge.

## When to use me

Use this skill when:
- PR has been approved
- All review comments addressed
- CI is passing
- About to click "Merge"

## Pre-Merge Checklist

### 1. Code Quality Gates

```bash
# Run full compliance check
make check-compliance

# Architecture validation
make check-intent-architecture

# Strict pattern check
make check-patterns-strict

# Documentation validation
make check-docblocks
```

**All must pass before merge.**

### 2. Test Coverage Verification

```bash
# Check coverage hasn't dropped
make coverage

# For modified packages specifically
go test -cover ./internal/cli/intents/your_feature/...
```

| Requirement | Threshold |
|-------------|-----------|
| Per-package (modified) | >= 95% |
| Project average | >= 80% |

### 3. Review Status Check

```bash
PR_NUM=$(gh pr view --json number -q '.number')

# Verify approval
gh pr view $PR_NUM --json reviews -q '.reviews[] | "\(.author.login): \(.state)"'

# Check all threads resolved
gh api graphql -f query='
  query($owner: String!, $repo: String!, $pr: Int!) {
    repository(owner: $owner, name: $repo) {
      pullRequest(number: $pr) {
        reviewThreads(first: 100) {
          nodes { isResolved }
        }
      }
    }
  }
' -f owner=OWNER -f repo=REPO -f pr=$PR_NUM | jq '.data.repository.pullRequest.reviewThreads.nodes | map(select(.isResolved == false)) | length'
```

### 4. CI Status Verification

```bash
# All checks must pass
gh pr checks $PR_NUM

# No pending checks
gh pr view $PR_NUM --json statusCheckRollup -q '.statusCheckRollup[] | select(.status != "COMPLETED")'
```

### 5. Branch Hygiene

```bash
# Ensure targeting correct branch
gh pr view $PR_NUM --json baseRefName -q '.baseRefName'
# Should be: next (NEVER main for feature PRs)

# Check for merge conflicts
gh pr view $PR_NUM --json mergeable -q '.mergeable'
# Should be: MERGEABLE

# Rebase if needed (prefer rebase over merge commits)
git fetch origin next
git rebase origin/next
git push --force-with-lease
```

### 6. Commit History Review

```bash
# Review commits (should be atomic and well-described)
gh pr view $PR_NUM --json commits -q '.commits[] | "\(.oid[0:7]) \(.messageHeadline)"'

# Check for AI attribution (all AI commits should have Co-authored-by)
git log origin/next..HEAD --format="%h %s" | while read commit; do
  if ! git log -1 --format="%b" ${commit%% *} | grep -q "Co-authored-by:"; then
    echo "WARNING: $commit may be missing AI attribution"
  fi
done
```

### 7. Documentation Completeness

Check that documentation was updated if:
- [ ] New public API added -> needs godoc
- [ ] New feature -> needs user-facing docs
- [ ] Breaking change -> needs migration guide
- [ ] Configuration change -> needs updated examples

### 8. Security Scan

```bash
# Run security scanner
make gosec

# Check for secrets
git diff origin/next..HEAD | grep -i -E "(password|secret|key|token|api_key)" && echo "WARNING: Possible secrets"
```

## Final Validation Template

Before clicking merge, copy and verify:

```markdown
## Pre-Merge Validation

PR: #[NUM] - [Title]
Branch: [feature-branch] -> next

### Code Quality
- [ ] `make check-compliance` passes
- [ ] `make check-intent-architecture` passes
- [ ] `make check-docblocks` passes
- [ ] No forbidden comment markers (TODO, FIXME, HACK)

### Testing
- [ ] All tests pass (`make test`)
- [ ] Coverage >= 95% for modified packages
- [ ] No race conditions (`go test -race ./...`)

### Review
- [ ] Approved by at least one reviewer
- [ ] All review threads resolved
- [ ] All comments addressed with evidence

### CI/CD
- [ ] All CI checks pass
- [ ] No pending checks
- [ ] No merge conflicts

### Branch
- [ ] Targeting `next` branch (NOT `main`)
- [ ] Rebased on latest `next`
- [ ] Commits are atomic and well-described
- [ ] AI commits have proper attribution

### Documentation
- [ ] Godoc updated for new public APIs
- [ ] User docs updated if needed
- [ ] No undocumented breaking changes

### Security
- [ ] `make gosec` passes
- [ ] No secrets in diff
- [ ] No security vulnerabilities introduced

**Ready to merge: [YES/NO]**
```

## Common Pre-Merge Issues

| Issue | Solution |
|-------|----------|
| Unresolved threads | Reply and resolve, or note why unresolved |
| CI still running | Wait - never merge with pending checks |
| Coverage dropped | Add tests before merging |
| Merge conflict | Rebase on latest `next` |
| Missing approval | Request re-review |
| Targeting `main` | Change base branch to `next` |

## Merge Methods

**Prefer squash merge** for feature branches:
```bash
gh pr merge $PR_NUM --squash
```

**Use merge commit** for release branches (next -> main):
```bash
gh pr merge $PR_NUM --merge
```

**Never use rebase merge** for KaRiya (loses co-author info).

## Post-Merge Actions

```bash
# Delete local feature branch
git checkout next
git branch -d feature/your-feature

# Update local next
git pull origin next

# Verify merge
gh pr view $PR_NUM --json state -q '.state'
# Should be: MERGED
```

## Related Skills

- `pr-monitor` - Track PR status and reviews
- `check-compliance` - Run all quality checks
- `ai-commit` - Proper commit attribution
- `release-management` - For next -> main merges
