---
name: git-advanced
description: Advanced Git operations - rebasing, cherry-picking, bisect, history management, recovery
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Guide advanced Git operations for maintaining clean history, finding bugs, and recovering from mistakes.

## When to use me

- Cleaning up commit history before PR
- Finding which commit introduced a bug
- Moving commits between branches
- Recovering from Git mistakes
- Understanding complex history

## Core Principles

1. **Clean history helps everyone** - Future you will thank present you
2. **Don't rewrite shared history** - Only rebase unpushed commits
3. **Commit early, clean up later** - Work messy, publish clean
4. **Git rarely loses data** - Recovery is usually possible

## Interactive Rebase

### Cleaning Up Before PR

```bash
# Rebase last N commits interactively
git rebase -i HEAD~5

# Rebase all commits on branch
git rebase -i main
```

### Rebase Commands

```
pick   = use commit as-is
reword = use commit, but edit message
edit   = use commit, but stop to amend
squash = meld into previous commit (keep message)
fixup  = meld into previous commit (discard message)
drop   = remove commit entirely
```

### Common Cleanups

```bash
# Example: Clean up work-in-progress commits

# Before:
# abc123 WIP
# def456 fix typo
# ghi789 Add user filter
# jkl012 WIP filter
# mno345 Complete filter feature

# git rebase -i HEAD~5
# Change to:
pick mno345 Complete filter feature
fixup jkl012 WIP filter
fixup ghi789 Add user filter
fixup def456 fix typo
fixup abc123 WIP

# After: Single clean commit
# xyz999 Complete filter feature
```

### Reword Commit Message

```bash
git rebase -i HEAD~3

# Change 'pick' to 'reword' for the commit
# Save and close
# Editor opens for new message
```

### Split a Commit

```bash
git rebase -i HEAD~3

# Change 'pick' to 'edit' for the commit to split
# Save and close

# Now at that commit:
git reset HEAD~1           # Undo commit, keep changes
git add file1.go
git commit -m "First part"
git add file2.go  
git commit -m "Second part"
git rebase --continue
```

## Cherry-Pick

### Moving Commits Between Branches

```bash
# Copy single commit to current branch
git cherry-pick abc1234

# Copy range of commits
git cherry-pick abc1234..def5678

# Copy multiple specific commits
git cherry-pick abc1234 def5678 ghi9012

# Cherry-pick without committing (stage only)
git cherry-pick -n abc1234
```

### Cherry-Pick Workflow

```bash
# Scenario: Need hotfix from feature branch on main

# 1. Find the commit hash
git log feature/auth --oneline
# abc1234 Fix critical auth bug

# 2. Switch to main
git checkout main

# 3. Cherry-pick the fix
git cherry-pick abc1234

# 4. Push hotfix
git push origin main
```

### Handling Conflicts

```bash
git cherry-pick abc1234
# CONFLICT!

# Resolve conflicts in files
git add resolved_file.go

# Continue
git cherry-pick --continue

# Or abort
git cherry-pick --abort
```

## Git Bisect

### Finding Bug-Introducing Commit

```bash
# Start bisect
git bisect start

# Mark current (broken) as bad
git bisect bad

# Mark known good commit
git bisect good v1.0.0

# Git checks out middle commit
# Test it, then:
git bisect good  # if working
git bisect bad   # if broken

# Repeat until Git finds the culprit
# "abc1234 is the first bad commit"

# End bisect
git bisect reset
```

### Automated Bisect

```bash
# Run test automatically
git bisect start HEAD v1.0.0
git bisect run make test

# Or with custom script
git bisect run ./test_for_bug.sh

# Script should:
# - Exit 0 if good
# - Exit 1 if bad
# - Exit 125 to skip (can't test this commit)
```

### Bisect Example

```bash
#!/bin/bash
# test_for_bug.sh

# Build
go build ./... || exit 125  # Skip if doesn't build

# Run specific test
go test -run TestBrokenFeature ./...
# Exit code indicates good/bad
```

## History Management

### View History

```bash
# Compact log
git log --oneline -20

# With graph
git log --oneline --graph --all

# Find commits by message
git log --grep="filter"

# Find commits changing file
git log --follow -- path/to/file.go

# Find commits by author
git log --author="name"

# Find commits in date range
git log --since="2024-01-01" --until="2024-02-01"
```

### Amending Commits

```bash
# Amend last commit message
git commit --amend -m "New message"

# Amend last commit with staged changes
git add forgotten_file.go
git commit --amend --no-edit

# WARNING: Never amend pushed commits!
```

### Reset vs Revert

```bash
# RESET: Move branch pointer (rewrites history)
git reset --soft HEAD~1   # Undo commit, keep staged
git reset --mixed HEAD~1  # Undo commit, unstage (default)
git reset --hard HEAD~1   # Undo commit, discard changes

# REVERT: Create new commit that undoes (safe for shared)
git revert abc1234        # Creates new commit undoing abc1234
git revert HEAD~3..HEAD   # Revert last 3 commits
```

## Recovery

### Recovering Deleted Commits

```bash
# Find lost commits in reflog
git reflog

# Example output:
# abc1234 HEAD@{0}: reset: moving to HEAD~3
# def5678 HEAD@{1}: commit: Important work  <-- Lost!
# ghi9012 HEAD@{2}: commit: Previous work

# Recover
git checkout def5678           # Detached HEAD at lost commit
git checkout -b recovered      # Create branch to keep it

# Or cherry-pick back
git cherry-pick def5678
```

### Recovering Deleted Branch

```bash
# Find where branch was
git reflog | grep "branch-name"
# Or
git reflog --all | grep "branch-name"

# Recreate branch at that commit
git branch recovered-branch abc1234
```

### Recovering Staged but Not Committed

```bash
# Files were staged but commit was reset
git fsck --lost-found

# Look in .git/lost-found/other/
# Files are there as blobs

# Recover specific blob
git show <blob-hash> > recovered_file.go
```

### Recovering from Bad Rebase

```bash
# Find pre-rebase state
git reflog
# abc1234 HEAD@{5}: rebase -i (start): checkout main
# def5678 HEAD@{6}: commit: Last commit before rebase <-- Here!

# Reset to pre-rebase state
git reset --hard def5678
```

## Branch Management

### Cleaning Up Branches

```bash
# Delete merged local branches
git branch --merged main | grep -v "main\|next" | xargs git branch -d

# Delete remote tracking branches for deleted remotes
git fetch --prune

# Delete remote branch
git push origin --delete feature/old-branch
```

### Renaming Branches

```bash
# Rename current branch
git branch -m new-name

# Rename other branch
git branch -m old-name new-name

# Update remote
git push origin -u new-name
git push origin --delete old-name
```

## Stashing

### Advanced Stash Usage

```bash
# Stash with message
git stash push -m "Work in progress on filter"

# Stash specific files
git stash push -m "Partial work" -- file1.go file2.go

# Stash including untracked
git stash push -u -m "Including new files"

# List stashes
git stash list

# Apply specific stash
git stash apply stash@{2}

# Pop specific stash
git stash pop stash@{2}

# Show stash contents
git stash show -p stash@{0}

# Create branch from stash
git stash branch new-branch stash@{0}
```

## Safety Guidelines

### Golden Rules

```markdown
## NEVER on shared branches (main, next):
- git push --force
- git reset --hard (then push)
- git rebase (then push)
- Amend pushed commits

## SAFE to do on feature branches (unpushed):
- Interactive rebase
- Amend commits
- Reset
- Force push YOUR OWN branch (with lease)
```

### Force Push Safely

```bash
# WRONG: Overwrites everything
git push --force  # DANGEROUS

# BETTER: Fails if remote has new commits
git push --force-with-lease  # Safer

# BEST: Only force push your own branches
# And communicate with team first
```

## Related Skills

- `ai-commit` - Commit message standards
- `git-worktree` - Parallel development
- `create-pr` - Pull request workflow
