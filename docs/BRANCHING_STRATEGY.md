# KaRiya Branching Strategy

## Overview

KaRiya uses a **dual-branch workflow** with `next` as the integration branch and `main` as the release branch. This strategy ensures that all features are tested together before releasing to production.

---

## Branch Structure

### `main` Branch
- **Purpose**: Production releases only
- **Protected**: Yes
- **Triggers**: Automatic semantic releases via GitHub Actions
- **Merge Source**: Only from `next` branch
- **Direct Commits**: Not allowed (except emergency hotfixes with approval)

### `next` Branch
- **Purpose**: Integration and testing of completed features
- **Protected**: Yes
- **Triggers**: Full CI suite (lint, test, build, security)
- **Merge Source**: Feature branches, fix branches, refactor branches
- **Release Workflow**: Manual merge to `main` when ready for release

### Feature Branches
- **Naming Convention**: `feature/description`, `fix/description`, `refactor/description`
- **Purpose**: Development of individual features or fixes
- **Merge Target**: `next` branch only
- **Lifetime**: Temporary (deleted after merge)

---

## Workflow Diagrams

### Standard Feature Development

```
feature/my-feature → (PR) → next → (PR when ready) → main → (auto) → GitHub Release
                       ↓              ↓                        ↓
                      CI             CI                   Semantic Release
```

### Release Process

```
Step 1: Feature Development
  feature/authentication → PR → next (merge & delete feature branch)
  feature/dashboard      → PR → next (merge & delete feature branch)
  fix/login-bug          → PR → next (merge & delete fix branch)

Step 2: Testing on next
  - All features integrated on next
  - Full CI passes
  - Manual testing completed
  - Stakeholder approval

Step 3: Release
  next → PR → main (manual merge when ready for release)
         ↓
    Automatic semantic-release creates:
      - Version tag (based on conventional commits)
      - GitHub release with binaries
      - CHANGELOG.md update
      - VERSION file update
```

---

## Developer Workflow

### 1. Starting New Work

```bash
# Ensure next is up to date
git checkout next
git pull origin next

# Create feature branch from next
git checkout -b feature/my-awesome-feature

# Make changes and commit (use conventional commits)
git add .
git commit -m "feat: add awesome feature"

# Push to remote
git push -u origin feature/my-awesome-feature
```

### 2. Creating a Pull Request

1. **Push your branch** to GitHub
2. **Open PR** targeting `next` branch (NOT `main`)
3. **Wait for CI** to pass (lint, tests, build, security)
4. **Address review comments** if any
5. **Squash and merge** when approved

```bash
# PR Title Examples (conventional commits format):
feat: add user authentication
fix: resolve login timeout issue
refactor: simplify event handling
perf: optimize database queries
docs: update API documentation
```

### 3. After Merge

```bash
# Switch back to next
git checkout next

# Pull latest changes
git pull origin next

# Delete local feature branch
git branch -d feature/my-awesome-feature

# Delete remote feature branch (if not auto-deleted)
git push origin --delete feature/my-awesome-feature
```

---

## Release Process

### Creating a Release (Maintainers Only)

**When to Release:**
- Multiple features ready for production
- All tests passing on `next`
- Manual testing completed
- Stakeholder approval received

**How to Release:**

```bash
# 1. Ensure next is ready
git checkout next
git pull origin next

# 2. Run full compliance check
make ci-local

# 3. Create PR from next to main
git checkout main
git pull origin main
gh pr create --base main --head next --title "release: prepare next release" --body "Release PR from next to main"

# 4. Wait for CI to pass and merge
# This will trigger automatic semantic-release

# 5. After merge, semantic-release will:
#    - Analyze commits (conventional commits)
#    - Determine version bump (major/minor/patch)
#    - Create GitHub release with binaries
#    - Update CHANGELOG.md
#    - Create git tag
```

**Manual Merge Alternative:**

```bash
# If you prefer manual merge instead of PR
git checkout main
git pull origin main
git merge next --no-ff -m "release: merge next into main for release"
git push origin main

# Semantic release will automatically trigger
```

---

## Branch Protection Rules

### Recommended Settings for `main`

- ✅ Require pull request before merging
- ✅ Require approvals: 1+
- ✅ Require status checks to pass
  - `Lint & Format`
  - `Test (ubuntu-latest, 1.24)`
  - `Build`
  - `Security Scan`
- ✅ Require conversation resolution before merging
- ✅ Do not allow bypassing the above settings
- ✅ Restrict who can push to matching branches (maintainers only)

### Recommended Settings for `next`

- ✅ Require pull request before merging
- ✅ Require approvals: 1+ (can be lower than main)
- ✅ Require status checks to pass
  - `Lint & Format`
  - `Test (ubuntu-latest, 1.24)`
- ✅ Require conversation resolution before merging
- ✅ Allow force pushes: No
- ✅ Allow deletions: No

---

## Commit Message Format

All commits must follow **Conventional Commits** format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types (determines version bump)

| Type | Version Bump | Purpose |
|------|--------------|---------|
| `feat` | Minor (0.X.0) | New feature |
| `fix` | Patch (0.0.X) | Bug fix |
| `perf` | Patch (0.0.X) | Performance improvement |
| `refactor` | Patch (0.0.X) | Code refactoring |
| `build` | Patch (0.0.X) | Build system changes |
| `revert` | Patch (0.0.X) | Revert previous commit |
| `docs` | None | Documentation only |
| `style` | None | Code style (formatting) |
| `test` | None | Tests only |
| `ci` | None | CI/CD changes |
| `chore` | None | Maintenance tasks |
| `BREAKING CHANGE` | Major (X.0.0) | Breaking API changes |

### Examples

```bash
# Feature (minor version bump)
feat(auth): add OAuth2 authentication

# Bug fix (patch version bump)
fix(export): resolve CSV export encoding issue

# Breaking change (major version bump)
feat(api)!: redesign REST API endpoints

BREAKING CHANGE: API endpoints have been restructured.
Migration guide available in docs/MIGRATION.md

# Performance improvement (patch version bump)
perf(cv): optimize bullet point generation algorithm

# Refactoring (patch version bump)
refactor(tui): simplify intent state machines

# Documentation (no version bump)
docs: update branching strategy documentation

# Tests (no version bump)
test(intents): add edge case tests for CaptureEvent
```

---

## Hotfix Process (Emergency Only)

For critical production bugs that cannot wait for the next release:

```bash
# 1. Create hotfix branch from main
git checkout main
git pull origin main
git checkout -b fix/critical-security-issue

# 2. Make minimal fix
git add .
git commit -m "fix: patch critical security vulnerability"

# 3. Create PR targeting main (exception to the rule)
gh pr create --base main --head fix/critical-security-issue \
  --title "fix: patch critical security vulnerability" \
  --body "EMERGENCY HOTFIX - Critical security issue requires immediate release"

# 4. After merge to main, backport to next
git checkout next
git pull origin next
git merge main
git push origin next
```

⚠️ **Hotfixes should be rare** - Most fixes should go through the normal `next` → `main` workflow.

---

## CI/CD Pipeline Behavior

### On Feature Branch Push
- ✅ Lint & format checks
- ✅ Tests (all platforms)
- ✅ Build verification
- ✅ Security scan

### On PR to `next`
- ✅ All CI checks (above)
- ✅ Commit message validation
- ✅ AI attribution check (if applicable)
- ✅ Breaking changes detection
- ✅ PR title validation
- ✅ Size labeling

### On `next` Branch Merge
- ✅ All CI checks run
- ❌ No release triggered

### On PR to `main`
- ✅ All CI checks (strict)
- ✅ Full compliance verification
- ⚠️ Review breaking changes carefully

### On `main` Branch Merge
- ✅ All CI checks run
- ✅ Build release binaries
- ✅ Semantic version analysis
- ✅ Create GitHub release
- ✅ Upload binaries to release
- ✅ Update CHANGELOG.md
- ✅ Update VERSION file
- ✅ Create git tag

---

## FAQ

### Q: Can I create a PR directly to `main`?

**A:** Only for emergency hotfixes with maintainer approval. All regular development should target `next`.

### Q: How often should we release?

**A:** Release when a meaningful set of features is ready and tested on `next`. This could be weekly, bi-weekly, or when specific milestones are reached.

### Q: What if CI fails on `next`?

**A:** Fix the issue immediately with a new PR to `next`. The `next` branch should always be in a releasable state.

### Q: Can I merge multiple feature branches at once?

**A:** Yes! That's the purpose of `next` - to integrate multiple features together before release.

### Q: What if my feature depends on another feature?

**A:** Wait for the dependency to be merged to `next`, then rebase your feature branch on the latest `next`.

```bash
git checkout feature/my-feature
git fetch origin
git rebase origin/next
git push --force-with-lease
```

### Q: How do I know what version will be released?

**A:** Semantic-release analyzes conventional commits:
- Any `feat` commits → minor version (0.X.0)
- Any `BREAKING CHANGE` → major version (X.0.0)
- Only `fix`, `perf`, `refactor` → patch version (0.0.X)

### Q: Can I test the release process locally?

**A:** Yes! Use dry-run mode:

```bash
npm install
npx semantic-release --dry-run
```

This shows what version would be released without actually releasing.

---

## Tools and Commands

### Useful Git Aliases

Add to `~/.gitconfig`:

```ini
[alias]
  # Create feature branch from next
  feature = "!f() { git checkout next && git pull && git checkout -b feature/$1; }; f"
  
  # Create fix branch from next
  fix = "!f() { git checkout next && git pull && git checkout -b fix/$1; }; f"
  
  # Update current branch with latest next
  sync = "!git fetch origin && git rebase origin/next"
  
  # Clean up merged branches
  cleanup = "!git branch --merged next | grep -v '\\*\\|main\\|next' | xargs -n 1 git branch -d"
```

### GitHub CLI Commands

```bash
# Create PR to next (default)
gh pr create --base next --fill

# View PR status
gh pr status

# List all PRs to next
gh pr list --base next

# Create release PR (next → main)
gh pr create --base main --head next --title "release: prepare next release"

# View release workflow status
gh run list --workflow=release.yml
```

### Make Commands

```bash
# Run all CI checks locally (before pushing)
make ci-local

# Run compliance check
make check-compliance

# Review commits before pushing
make review-commit
```

---

## Migration Notes

**Previous Workflow:**
- Feature branches merged directly to `main`
- Every merge to `main` triggered a release

**New Workflow:**
- Feature branches merge to `next`
- `next` merges to `main` only when ready for release
- Releases are deliberate and controlled

**Benefits:**
- 🎯 **Controlled releases**: Release when ready, not on every merge
- 🧪 **Integration testing**: Test features together before release
- 📦 **Batched releases**: Group related features into meaningful releases
- 🐛 **Reduced risk**: Catch integration issues before production
- 📝 **Better changelogs**: Meaningful release notes with grouped features

---

## References

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)
- [Semantic Release](https://semantic-release.gitbook.io/)
- [GitHub Flow](https://guides.github.com/introduction/flow/)
- [Git Feature Branch Workflow](https://www.atlassian.com/git/tutorials/comparing-workflows/feature-branch-workflow)

---

**Last Updated**: 2026-01-08  
**Status**: ✅ Active
