---
name: git-worktree
description: Use Git worktrees for parallel development - test features in isolation, work on multiple branches simultaneously
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Guide the use of Git worktrees for parallel development. Worktrees allow you to have multiple branches checked out simultaneously in separate directories, enabling testing and development in isolation.

## When to use me

- Working on multiple features simultaneously
- Testing a branch while continuing work on another
- Reviewing PRs while keeping your work intact
- Running long tests on one branch while developing on another
- Comparing behaviour between branches

## Core Concept

```
Traditional Git:
repo/
├── .git/
├── src/
└── (only ONE branch at a time)

With Worktrees:
repo.git/                    # Bare repository (no working files)
├── main/                    # Worktree for main branch
│   └── src/
├── feature-auth/            # Worktree for feature branch
│   └── src/
└── bugfix-123/              # Worktree for bugfix
    └── src/

Each worktree is independent - different branch, same repo.
```

## Setting Up Bare Repository

### Converting Existing Repo to Bare + Worktrees

```bash
# From parent directory of your project
cd /path/to/projects

# Clone as bare repository
git clone --bare git@github.com:user/repo.git repo.git

# Or convert existing repo
mv repo/.git repo.git
rm -rf repo
cd repo.git
git config --bool core.bare true

# Create worktree for main development
git worktree add main main
# or for 'next' branch
git worktree add next next
```

### Fresh Setup

```bash
# Clone bare
git clone --bare git@github.com:user/kariya.git kariya.git
cd kariya.git

# Create primary worktree
git worktree add next next

# Your structure:
# kariya.git/           # Bare repo
# kariya.git/next/      # Working directory on 'next' branch
```

## Worktree Commands

### Create Worktree

```bash
# New worktree from existing branch
git worktree add ../feature-auth feature/auth

# New worktree with new branch (from current HEAD)
git worktree add -b feature/new-thing ../new-thing

# New worktree with new branch from specific base
git worktree add -b feature/new-thing ../new-thing origin/next
```

### List Worktrees

```bash
git worktree list

# Output:
# /path/to/kariya.git        (bare)
# /path/to/kariya.git/next   abc1234 [next]
# /path/to/feature-auth      def5678 [feature/auth]
```

### Remove Worktree

```bash
# Remove worktree (keeps branch)
git worktree remove ../feature-auth

# Force remove (if uncommitted changes)
git worktree remove --force ../feature-auth

# Clean up stale worktree references
git worktree prune
```

### Move Worktree

```bash
git worktree move ../old-path ../new-path
```

## Workflow Patterns

### Pattern 1: Feature Development

```bash
# You're working on 'next' branch
cd kariya.git/next

# Need to start new feature
git worktree add -b feature/user-auth ../user-auth origin/next

# Work on feature in separate directory
cd ../user-auth
# ... make changes, commit, push ...

# When done, merge and clean up
cd ../next
git merge feature/user-auth
git worktree remove ../user-auth
git branch -d feature/user-auth
```

### Pattern 2: Parallel Testing

```bash
# Testing branch A while working on branch B
git worktree add ../test-branch-a branch-a

# Terminal 1: Run tests on branch-a
cd ../test-branch-a
make test

# Terminal 2: Continue working on current branch
cd ../next
# ... keep coding ...
```

### Pattern 3: PR Review

```bash
# Review a PR without losing your work
git fetch origin pull/123/head:pr-123
git worktree add ../pr-123 pr-123

# Review in separate directory
cd ../pr-123
make test
# ... review code ...

# Clean up after review
cd ../next
git worktree remove ../pr-123
git branch -D pr-123
```

### Pattern 4: Comparing Behaviour

```bash
# Compare current vs main
git worktree add ../compare-main main

# Terminal 1: Run current version
cd ../next
go run ./cmd/kariya

# Terminal 2: Run main version
cd ../compare-main
go run ./cmd/kariya

# Side-by-side comparison
```

## Directory Structure Convention

```
~/Projects/
├── KaRiya.git/              # Bare repository
│   ├── next/                # Primary development worktree
│   ├── main/                # Production branch (optional)
│   └── hooks/               # Git hooks (shared)
├── KaRiya-feature-auth/     # Feature worktree (temporary)
├── KaRiya-bugfix-123/       # Bugfix worktree (temporary)
└── KaRiya-pr-456/           # PR review worktree (temporary)
```

**Naming convention:**
- Primary worktrees: Inside `.git` directory
- Temporary worktrees: Sibling directories with prefix

## Important Considerations

### Shared Resources

```bash
# All worktrees share:
# - .git objects (commits, blobs)
# - Remote configuration
# - Hooks

# Each worktree has its own:
# - Working directory
# - Index (staging area)
# - HEAD reference
```

### Branch Locking

```bash
# Cannot checkout same branch in multiple worktrees
# This will fail:
git worktree add ../another-next next
# fatal: 'next' is already checked out at '/path/to/next'

# Solution: Use different branches or detached HEAD
git worktree add --detach ../temp-next next
```

### Dependencies and Build Artifacts

```bash
# Each worktree needs its own:
# - node_modules/
# - vendor/
# - Build cache

# After creating worktree:
cd ../new-worktree
go mod download        # Go dependencies
npm install            # Node dependencies (if any)
```

### IDE Configuration

```bash
# Each worktree may need its own:
# - .idea/ (IntelliJ)
# - .vscode/ (VS Code settings)
# - .env (environment variables)

# Copy from primary worktree if needed:
cp ../next/.env .
```

## Integration with KaRiya Workflow

### Session Start in Worktree

```bash
# Always run session-start in each worktree
cd kariya.git/next
make session-start

# For feature worktree
cd ../feature-auth
make session-start
```

### Committing in Worktrees

```bash
# Same process - each worktree is independent
cd ../feature-auth
make ai-commit FILE=/tmp/commit.txt

# Push from worktree
git push -u origin feature/auth
```

### Creating PR from Worktree

```bash
cd ../feature-auth
make pre-pr
gh pr create --base next --title "Add user authentication"
```

## Troubleshooting

### "Branch is already checked out"

```bash
# Find where it's checked out
git worktree list

# Remove that worktree or use different branch
git worktree remove /path/to/worktree
```

### Stale Worktree References

```bash
# After manually deleting worktree directory
git worktree prune

# List should be clean now
git worktree list
```

### Worktree Corruption

```bash
# If worktree gets corrupted
git worktree remove --force ../broken-worktree
git worktree prune

# Recreate
git worktree add ../fresh-worktree branch-name
```

## Quick Reference

```bash
# Create worktree
git worktree add <path> <branch>
git worktree add -b <new-branch> <path> [<base>]

# List worktrees
git worktree list

# Remove worktree
git worktree remove <path>

# Clean stale references
git worktree prune

# Move worktree
git worktree move <old-path> <new-path>
```

## Related Skills

- `ai-commit` - Committing in worktrees
- `create-pr` - PRs from feature worktrees
- `session-start` - Session setup per worktree
- `github-expert` - GitHub CLI operations
