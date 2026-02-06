---
description: Manage Git worktrees for parallel development
agent: build
---

Manage Git worktrees - create, switch, list, or remove worktrees without thinking about git commands.

Load the `git-worktree` skill for reference.

## Usage

```
/worktree list              # Show all worktrees and available branches
/worktree new <name>        # Create worktree for feature/name branch
/worktree pr <number>       # Create worktree to review PR #number
/worktree switch <name>     # Switch to existing worktree
/worktree remove <name>     # Remove worktree and clean up
/worktree status            # Show current worktree info
```

## Arguments
$ARGUMENTS

## Process

### 1. Parse the Command

Determine action from arguments:
- No args or `list` → List worktrees
- `new <name>` → Create new feature worktree
- `pr <number>` → Create PR review worktree
- `switch <name>` → Switch to worktree
- `remove <name>` → Remove worktree
- `status` → Show current state

### 2. Execute Action

#### List Worktrees
```bash
echo "=== Active Worktrees ==="
git worktree list

echo ""
echo "=== Feature Branches (local) ==="
git branch --list 'feature/*' --format='%(refname:short)'

echo ""
echo "=== Feature Branches (remote) ==="
git branch -r --list 'origin/feature/*' --format='%(refname:short)' | sed 's|origin/||'
```

#### Create New Feature Worktree
```bash
BRANCH_NAME="feature/$NAME"
WORKTREE_PATH="../KaRiya-$NAME"

# Check if branch exists
if git show-ref --verify --quiet "refs/heads/$BRANCH_NAME"; then
    echo "Branch $BRANCH_NAME exists, creating worktree..."
    git worktree add "$WORKTREE_PATH" "$BRANCH_NAME"
else
    echo "Creating new branch $BRANCH_NAME from origin/next..."
    git fetch origin next
    git worktree add -b "$BRANCH_NAME" "$WORKTREE_PATH" origin/next
fi

echo ""
echo "Worktree ready at: $WORKTREE_PATH"
echo "Run: cd $WORKTREE_PATH && make session-start"
```

#### Create PR Review Worktree
```bash
PR_NUM="$NUMBER"
WORKTREE_PATH="../KaRiya-pr-$PR_NUM"

# Fetch PR branch
gh pr checkout $PR_NUM --detach 2>/dev/null || true
BRANCH=$(gh pr view $PR_NUM --json headRefName -q '.headRefName')

git fetch origin "pull/$PR_NUM/head:pr-$PR_NUM"
git worktree add "$WORKTREE_PATH" "pr-$PR_NUM"

echo ""
echo "PR #$PR_NUM worktree ready at: $WORKTREE_PATH"
echo "Original branch: $BRANCH"
echo "Run: cd $WORKTREE_PATH"
```

#### Switch to Worktree
```bash
# Find worktree path
WORKTREE_PATH=$(git worktree list | grep -E "(feature/$NAME|$NAME)" | awk '{print $1}')

if [ -z "$WORKTREE_PATH" ]; then
    # Try with KaRiya prefix
    WORKTREE_PATH="../KaRiya-$NAME"
fi

if [ -d "$WORKTREE_PATH" ]; then
    echo "Switching to worktree: $WORKTREE_PATH"
    # Report the path for the user to cd into
else
    echo "Worktree not found for: $NAME"
    echo "Available worktrees:"
    git worktree list
fi
```

#### Remove Worktree
```bash
WORKTREE_PATH="../KaRiya-$NAME"

# Check if it exists
if git worktree list | grep -q "$WORKTREE_PATH"; then
    git worktree remove "$WORKTREE_PATH"
    echo "Worktree removed: $WORKTREE_PATH"
    
    # Optionally remove branch
    BRANCH_NAME="feature/$NAME"
    if git show-ref --verify --quiet "refs/heads/$BRANCH_NAME"; then
        echo "Branch $BRANCH_NAME still exists. Remove it? (run: git branch -d $BRANCH_NAME)"
    fi
else
    echo "Worktree not found: $WORKTREE_PATH"
fi

git worktree prune
```

#### Status
```bash
echo "=== Current Directory ==="
pwd

echo ""
echo "=== Current Branch ==="
git branch --show-current

echo ""
echo "=== All Worktrees ==="
git worktree list
```

### 3. Post-Action Guidance

After creating or switching worktrees, tell the user:

1. **The exact path** to cd into
2. **Run `make session-start`** in new worktrees
3. **Dependencies** - remind about `go mod download` if needed

### 4. Change Working Directory

When switching or creating, use the `workdir` parameter for subsequent commands:
- Set workdir to the new worktree path
- Confirm the switch worked

## Output Format

Always output:
1. Action taken
2. Worktree path (absolute)
3. Next steps

Example:
```
Created worktree: /home/user/Projects/KaRiya-my-feature
Branch: feature/my-feature
Base: origin/next

Next steps:
  cd /home/user/Projects/KaRiya-my-feature
  make session-start
```

## Directory Convention

```
~/Projects/
├── KaRiya/                  # Main repo (current)
├── KaRiya-feature-name/     # Feature worktrees
├── KaRiya-pr-123/           # PR review worktrees
└── KaRiya-hotfix-name/      # Hotfix worktrees
```
