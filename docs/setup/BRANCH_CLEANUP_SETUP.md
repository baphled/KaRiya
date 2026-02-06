# Branch Cleanup Setup

**Automatic deletion of merged branches for cleaner repository management**

## Overview

This guide documents the setup for automatic deletion of merged branches in the KaRiya repository. This includes both remote (GitHub) and local branch cleanup.

**Setup Date**: 2026-02-06  
**Status**: ✅ Configured and Active

---

## What Was Configured

### 1. GitHub Auto-Delete (Remote)

GitHub is now configured to automatically delete head branches after pull requests are merged.

**Enabled Features**:
- Auto-merge capability
- Delete branch on merge

**How It Works**:
- When a PR is merged into `next` or `main`
- GitHub automatically deletes the feature branch on the remote
- No manual cleanup needed for merged branches

**Verification**:
```bash
# Check GitHub repository settings
gh repo view --json deleteBranchOnMerge,hasAutoMerge
```

---

### 2. Local Branch Cleanup (Script)

A cleanup script was created to remove local branches that have been merged.

**Location**: `/home/baphled/Projects/KaRiya/scripts/cleanup-merged-branches.sh`

**Features**:
- Identifies branches merged into the default branch (main)
- Excludes protected branches (main, next)
- Confirmation prompt before deletion
- Colour-coded output for clarity
- Force mode available with `--force` flag

**Usage**:
```bash
# Interactive mode (asks for confirmation)
bash /home/baphled/Projects/KaRiya/scripts/cleanup-merged-branches.sh

# Force mode (no confirmation)
bash /home/baphled/Projects/KaRiya/scripts/cleanup-merged-branches.sh --force
```

---

### 3. Git Aliases (Quick Access)

Two Git aliases were added for easy access to the cleanup script:

```bash
# Interactive cleanup (with confirmation)
git cleanup

# Force cleanup (no confirmation)
git cleanup-force
```

**How They Work**:
- Use `$(git rev-parse --git-dir)` to locate the script
- Work from any worktree in the repository
- Compatible with bare repository setup

---

## When to Use

### Automatic (GitHub)

**No action needed** - branches are automatically deleted when PRs are merged.

**Applies to**:
- Feature branches (`feature/*`)
- Bug fix branches (`fix/*`)
- Any branch merged via PR

**Does NOT apply to**:
- Protected branches (`main`, `next`)
- Branches not merged via PR

---

### Manual (Local Script)

Use the cleanup script to remove local branches that were deleted remotely:

**Common Scenarios**:
```bash
# After merging several PRs
git cleanup

# Weekly cleanup of merged branches
git cleanup-force

# After pulling latest changes
git fetch --prune && git cleanup
```

**What It Removes**:
- Local branches merged into `main`
- Branches already deleted on remote
- Stale feature branches

**What It Preserves**:
- Current branch
- Protected branches (`main`, `next`)
- Unmerged branches

---

## Integration with KaRiya Workflow

### Standard Feature Workflow

```bash
# 1. Create feature branch
git worktree add ../KaRiya-feature-123 -b feature/task-123

# 2. Work on feature
cd ../KaRiya-feature-123
# ... make changes, test, commit ...

# 3. Create PR targeting 'next'
gh pr create --base next --title "feat: ..."

# 4. Merge PR (via GitHub UI or CLI)
gh pr merge --merge

# 5. GitHub automatically deletes remote branch
# (No manual action needed)

# 6. Clean up local branch
cd /home/baphled/Projects/KaRiya
git worktree remove ../KaRiya-feature-123
git cleanup  # Remove local merged branch
```

---

### Weekly Maintenance

Add to your weekly workflow:

```bash
# Start of week: update and clean
git fetch --prune          # Remove remote-tracking refs
git cleanup-force          # Delete merged local branches
git worktree prune         # Clean up worktree references
```

---

## Technical Details

### GitHub Configuration

**API Command Used**:
```bash
gh repo edit --enable-auto-merge --delete-branch-on-merge
```

**Repository Settings Changed**:
- `auto_merge_enabled`: `true`
- `delete_branch_on_merge`: `true`

**Scope**: Applies to all PRs in the repository

---

### Script Implementation

**Key Features**:

1. **Default Branch Detection**:
   ```bash
   DEFAULT_BRANCH=$(git symbolic-ref refs/remotes/origin/HEAD 2>/dev/null | \
       sed 's@^refs/remotes/origin/@@' || echo "main")
   ```

2. **Merged Branch Detection**:
   ```bash
   git branch --merged "$DEFAULT_BRANCH" | \
       grep -v "^\*" | \
       grep -v "^  $DEFAULT_BRANCH$" | \
       grep -v "^  next$"
   ```

3. **Safe Deletion**:
   ```bash
   git branch -d "$branch"  # Only deletes if merged
   ```

**Error Handling**:
- Fails safely if branch can't be deleted
- Shows which branches succeeded/failed
- Non-zero exit on any failure

---

### Git Alias Configuration

**Aliases Added**:
```bash
git config alias.cleanup \
    '!bash "$(git rev-parse --git-dir)/scripts/cleanup-merged-branches.sh"'

git config alias.cleanup-force \
    '!bash "$(git rev-parse --git-dir)/scripts/cleanup-merged-branches.sh" --force'
```

**Stored In**: `.git/config` (repository-specific)

**Compatibility**: Works with bare repository and worktrees

---

## Verification

### Verify GitHub Configuration

```bash
# Check repository settings
gh repo view --json deleteBranchOnMerge

# Expected output:
# {
#   "deleteBranchOnMerge": true
# }
```

---

### Verify Script

```bash
# Check script exists and is executable
ls -l /home/baphled/Projects/KaRiya/scripts/cleanup-merged-branches.sh

# Expected output:
# -rwxr-xr-x ... cleanup-merged-branches.sh
```

---

### Verify Aliases

```bash
# List Git aliases
git config --get alias.cleanup
git config --get alias.cleanup-force

# Expected output:
# !bash "$(git rev-parse --git-dir)/scripts/cleanup-merged-branches.sh"
# !bash "$(git rev-parse --git-dir)/scripts/cleanup-merged-branches.sh" --force
```

---

### Test Cleanup (Dry Run)

```bash
# See what would be deleted (no actual deletion)
git branch --merged main | grep -v "^\*" | grep -v "main" | grep -v "next"
```

---

## Troubleshooting

### Issue: Script Not Found

**Symptom**: `git cleanup` fails with "script not found"

**Solution**:
```bash
# Verify script location
ls /home/baphled/Projects/KaRiya/scripts/cleanup-merged-branches.sh

# If missing, recreate aliases
git config alias.cleanup \
    '!bash "$(git rev-parse --git-dir)/scripts/cleanup-merged-branches.sh"'
```

---

### Issue: Branch Not Deleted

**Symptom**: Branch still exists after merge

**Possible Causes**:
1. Branch not fully merged
2. Protected branch
3. Current branch

**Solution**:
```bash
# Check if branch is merged
git branch --merged main | grep "branch-name"

# Force delete (use with caution)
git branch -D branch-name
```

---

### Issue: Worktree References Remain

**Symptom**: Worktree directories still exist after branch deletion

**Solution**:
```bash
# List worktrees
git worktree list

# Remove worktree
git worktree remove /path/to/worktree

# Prune stale worktrees
git worktree prune
```

---

## Best Practices

### Do

- ✅ Run `git cleanup` regularly (weekly/monthly)
- ✅ Use `git fetch --prune` before cleanup
- ✅ Review branches before force cleanup
- ✅ Keep feature branches short-lived
- ✅ Merge PRs promptly after approval

### Don't

- ❌ Delete branches that aren't merged
- ❌ Force delete without verification
- ❌ Skip `git fetch --prune`
- ❌ Leave long-lived feature branches
- ❌ Bypass the cleanup script

---

## Related Documentation

### Branch Strategy
- **[BRANCHING_STRATEGY.md](../BRANCHING_STRATEGY.md)** - KaRiya branch workflow
- **[Git Advanced Skill](/.claude/skills/git-advanced/SKILL.md)** - Advanced Git operations

### Workflows
- **[DEVELOPMENT_WORKFLOW.md](../development/DEVELOPMENT_WORKFLOW.md)** - Complete development flow
- **[COMMON_TASKS.md](../development/COMMON_TASKS.md)** - Common Git tasks

### CI/CD
- **[CI_CD_SETUP.md](CI_CD_SETUP.md)** - CI/CD pipeline setup
- **[CI_CD_PIPELINE.md](../CI_CD_PIPELINE.md)** - Pipeline details

---

## Maintenance

### Script Updates

If the cleanup script needs modification:

1. Edit: `/home/baphled/Projects/KaRiya/scripts/cleanup-merged-branches.sh`
2. Test manually before committing
3. Update this documentation if behaviour changes
4. Verify aliases still work

---

### Configuration Changes

To modify GitHub settings:

```bash
# Disable auto-delete (not recommended)
gh repo edit --no-delete-branch-on-merge

# Re-enable
gh repo edit --delete-branch-on-merge
```

---

## Summary

**What's Configured**:
- ✅ GitHub auto-delete merged branches
- ✅ Local cleanup script
- ✅ Git aliases (`git cleanup`, `git cleanup-force`)

**Benefits**:
- Cleaner repository
- Less manual maintenance
- Fewer stale branches
- Reduced confusion

**Usage**:
```bash
# GitHub handles remote automatically
# Local cleanup as needed
git cleanup
```

---

**Last Updated**: 2026-02-06  
**Status**: ✅ Active and Tested  
**Location**: `docs/setup/BRANCH_CLEANUP_SETUP.md`
