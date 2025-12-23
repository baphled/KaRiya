# CI/CD Pipeline Documentation

## Overview

The KaRiya project uses **GitHub Actions** for continuous integration and deployment, with **commitlint** for commit message validation and **semantic-release** for automated versioning and releases.

---

## Table of Contents

1. [Architecture](#architecture)
2. [Workflows](#workflows)
3. [Commitlint Configuration](#commitlint-configuration)
4. [Semantic Release Configuration](#semantic-release-configuration)
5. [Release Process](#release-process)
6. [Local Development](#local-development)
7. [Troubleshooting](#troubleshooting)

---

## Architecture

### CI/CD Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                       Developer Workflow                         │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
                    ┌───────────────────────┐
                    │  Commit with          │
                    │  Conventional Format  │
                    └───────────────────────┘
                                │
                                ▼
                    ┌───────────────────────┐
                    │  Local Hooks          │
                    │  - AI Attribution     │
                    │  - Commitlint         │
                    └───────────────────────┘
                                │
                                ▼
                    ┌───────────────────────┐
                    │  Push to Branch       │
                    └───────────────────────┘
                                │
                ┌───────────────┴───────────────┐
                ▼                               ▼
┌───────────────────────┐         ┌───────────────────────┐
│   Feature Branch      │         │   Main Branch         │
│   CI Workflow         │         │   Release Workflow    │
├───────────────────────┤         ├───────────────────────┤
│ - Commitlint          │         │ - CI Checks           │
│ - Lint & Format       │         │ - Build Binaries      │
│ - Tests (All OS)      │         │ - Semantic Release    │
│ - Build               │         │   - Analyze Commits   │
│ - Security Scan       │         │   - Generate Version  │
└───────────────────────┘         │   - Create CHANGELOG  │
                                  │   - Create GitHub     │
                                  │     Release           │
                                  │   - Upload Binaries   │
                                  └───────────────────────┘
```

### Key Components

1. **Commitlint**: Validates commit messages follow conventional commits format
2. **Semantic Release**: Automates versioning based on commit messages
3. **GitHub Actions**: Runs CI/CD workflows
4. **Local Hooks**: Pre-commit validation for faster feedback

---

## Workflows

### 1. CI Workflow (`.github/workflows/ci.yml`)

**Triggers**: Push to any branch, Pull requests

**Jobs**:
- **Commitlint**: Validates PR commits follow conventional format
- **Lint**: Runs go fmt, go vet, staticcheck
- **Test**: Runs tests on Ubuntu, macOS, Windows with Go 1.24
- **Build**: Builds binaries for all platforms
- **Security**: Runs Gosec security scanner

**Status**: Required for PR merge

### 2. Release Workflow (`.github/workflows/release.yml`)

**Triggers**: Push to `main` branch only

**Jobs**:
- **CI Checks**: Runs full test suite
- **Build**: Creates release binaries with version info
- **Release**: Runs semantic-release to:
  - Analyze commits since last release
  - Determine new version (major/minor/patch)
  - Generate CHANGELOG.md
  - Create GitHub release
  - Upload binaries as release assets
  - Update VERSION file

**Permissions**: `contents: write`, `issues: write`, `pull-requests: write`

### 3. PR Validation Workflow (`.github/workflows/pr-validation.yml`)

**Triggers**: Pull request events (opened, synchronize, reopened)

**Jobs**:
- **Validate PR Title**: Ensures PR title follows conventional commits
- **Check AI Attribution**: Scans commits for AI attribution (warning only)
- **Conventional Commits**: Validates all commits in PR
- **Breaking Changes**: Detects breaking changes and comments on PR
- **Size Label**: Adds size label (xs/s/m/l/xl) based on changes

---

## Commitlint Configuration

### File: `.commitlintrc.json`

```json
{
  "extends": ["@commitlint/config-conventional"],
  "rules": {
    "type-enum": ["feat", "fix", "docs", "style", "refactor", "test", "chore", "perf", "ci", "build", "revert"],
    "scope-enum": ["domain", "service", "repo", "cli", "logger", "classification", "deps", "release"],
    "subject-max-length": [2, "always", 72],
    "header-max-length": [2, "always", 100]
  }
}
```

### Valid Commit Types

| Type | Description | Release |
|------|-------------|---------|
| `feat` | New feature | Minor |
| `fix` | Bug fix | Patch |
| `perf` | Performance improvement | Patch |
| `refactor` | Code refactoring | Patch |
| `docs` | Documentation only | None |
| `style` | Formatting, no code change | None |
| `test` | Adding tests | None |
| `build` | Build system changes | Patch |
| `ci` | CI/CD changes | None |
| `chore` | Maintenance tasks | None |
| `revert` | Revert previous commit | Patch |

### Valid Scopes

- `domain`: Domain layer changes
- `service`: Service layer changes
- `repo`: Repository layer changes
- `cli`: CLI interface changes
- `logger`: Logging system changes
- `classification`: Classification service changes
- `deps`: Dependency updates
- `release`: Release-related changes

### Commit Message Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Examples**:
```
feat(service): add event filtering by date range

Implement date range filtering for timeline views.
Improves performance for large event collections.

Closes #56
```

```
fix(domain): prevent duplicate tags in validation

Tag deduplication was not working correctly.

BREAKING CHANGE: Tags are now case-insensitive
```

---

## Semantic Release Configuration

### File: `.releaserc.json`

### Version Determination

Semantic release analyzes commit messages to determine version bump:

- **Major (X.0.0)**: Commits with `BREAKING CHANGE` in body/footer
- **Minor (0.X.0)**: Commits with `feat` type
- **Patch (0.0.X)**: Commits with `fix`, `perf`, `refactor`, `build`, `revert` types
- **No Release**: Commits with `docs`, `style`, `test`, `ci`, `chore` types

### Release Assets

For each release, semantic-release uploads:
- `kariya-linux-amd64` - Linux AMD64 binary
- `kariya-darwin-amd64` - macOS AMD64 binary
- `kariya-darwin-arm64` - macOS ARM64 binary
- `kariya-windows-amd64.exe` - Windows AMD64 binary

### Generated Files

- `CHANGELOG.md` - Auto-generated changelog
- `VERSION` - Current version number
- Git tag (e.g., `v1.2.3`)
- GitHub Release with binaries

---

## Release Process

### Automatic Release (Main Branch)

1. **Merge PR to main**:
   ```bash
   gh pr merge <pr-number> --squash
   ```

2. **Release workflow triggers automatically**:
   - Runs all CI checks
   - Builds release binaries
   - Analyzes commits since last release
   - Determines new version
   - Creates CHANGELOG entry
   - Creates git tag
   - Creates GitHub release
   - Uploads binaries

3. **Notification**:
   - GitHub release created with notes
   - CHANGELOG.md updated
   - VERSION file updated

### Manual Release (If Needed)

```bash
# Trigger manually from GitHub Actions UI
# Go to Actions → Release → Run workflow
```

### Version Examples

**Starting version**: `v0.1.0`

**Scenario 1: Feature Addition**
```bash
git commit -m "feat(service): add event search functionality"
# New version: v0.2.0 (minor bump)
```

**Scenario 2: Bug Fix**
```bash
git commit -m "fix(domain): correct date validation logic"
# New version: v0.1.1 (patch bump)
```

**Scenario 3: Breaking Change**
```bash
git commit -m "feat(api): redesign event structure

BREAKING CHANGE: Event structure has changed"
# New version: v1.0.0 (major bump)
```

### Pre-release Versions

Currently, we only release from `main`. For pre-releases:

1. Update `.releaserc.json` to include `develop` branch
2. Use `prerelease` configuration for beta/alpha versions

---

## Local Development

### Setup

1. **Install Node.js dependencies**:
   ```bash
   npm install
   ```

2. **Install git hooks**:
   ```bash
   make install-git-hooks
   ```

3. **Verify setup**:
   ```bash
   npm run commitlint -- --help
   ```

### Local Commit Validation

When you commit locally:

1. **prepare-commit-msg hook**: Adds AI attribution template (if code files)
2. **commit-msg hook**: Validates AI attribution format
3. **commit-msg-lint hook**: Validates conventional commits format (if Node.js installed)

**Example workflow**:
```bash
# Make changes
git add <files>

# Commit (hooks run automatically)
git commit

# Write message in editor:
feat(service): add event filtering

Implement date range filtering for timeline views.

AI-Generated-By: Avante (Claude 3.5 Sonnet)
Reviewed-By: John Doe
Closes #56
```

### Manual Validation

```bash
# Validate a commit message
echo "feat(service): add feature" | npx commitlint

# Validate last commit
npx commitlint --from HEAD~1

# Validate commit range
npx commitlint --from origin/main
```

---

## Troubleshooting

### Commitlint Fails

**Problem**: Commit rejected by commitlint

**Solutions**:

1. **Check format**:
   ```
   ✅ feat(service): add event filtering
   ❌ feat: add event filtering (missing scope)
   ❌ Feature: add event filtering (wrong case)
   ❌ feat(invalid): add feature (invalid scope)
   ```

2. **Check length**:
   - Header max: 100 characters
   - Subject max: 72 characters
   - Body lines max: 100 characters

3. **Valid types**: feat, fix, docs, style, refactor, test, chore, perf, ci, build, revert

4. **Valid scopes**: domain, service, repo, cli, logger, classification, deps, release

### Release Workflow Fails

**Problem**: Release workflow fails on main

**Check**:

1. **Commits are valid**: All commits follow conventional format
2. **Tests pass**: All tests must pass
3. **Permissions**: GitHub token has required permissions
4. **Branch protection**: Main branch rules don't block bot commits

**View logs**:
```bash
# Go to GitHub Actions
# Click on failed workflow
# Review job logs
```

### No Release Created

**Problem**: Push to main doesn't trigger release

**Reasons**:

1. **No releasable commits**: Only docs/style/test commits
2. **Already released**: No new commits since last release
3. **Workflow skipped**: Commit message contains `[skip ci]`

**Check**:
```bash
# View commits since last release
git log $(git describe --tags --abbrev=0)..HEAD --oneline

# Check for releasable commits (feat, fix, perf, refactor, build, revert)
git log $(git describe --tags --abbrev=0)..HEAD --grep="^feat\|^fix\|^perf"
```

### Commitlint Not Running Locally

**Problem**: Commits accepted locally but fail in CI

**Solutions**:

1. **Install dependencies**:
   ```bash
   npm install
   ```

2. **Reinstall hooks**:
   ```bash
   make install-git-hooks
   ```

3. **Verify hook exists**:
   ```bash
   ls -la .git/hooks/commit-msg-lint
   ```

4. **Test manually**:
   ```bash
   npx commitlint --from HEAD~1
   ```

### PR Validation Fails

**Problem**: PR blocked by validation failures

**Solutions**:

1. **Invalid PR title**:
   ```bash
   # Update PR title to follow conventional commits
   # Example: feat(service): add new feature
   ```

2. **Invalid commits**:
   ```bash
   # Rebase and fix commit messages
   git rebase -i origin/main
   # Edit commit messages to follow format
   ```

3. **Breaking changes not documented**:
   ```bash
   # Add BREAKING CHANGE to commit body
   git commit --amend
   # Add: BREAKING CHANGE: description
   ```

---

## CI/CD Optimization

### Caching

All workflows use caching for:
- Go modules (`go.sum`)
- Go build cache
- npm packages (`package-lock.json`)

### Parallelization

- Tests run in parallel across Ubuntu, macOS, Windows
- CI jobs run independently where possible
- Build job depends on test completion

### Resource Usage

- **CI Workflow**: ~5-10 minutes per run
- **Release Workflow**: ~10-15 minutes per run
- **PR Validation**: ~2-5 minutes per run

---

## Best Practices

### Commit Messages

1. **Use conventional format** for all commits
2. **Include scope** when applicable
3. **Write clear subjects** (imperative mood)
4. **Add body** for non-trivial changes
5. **Document breaking changes** in footer
6. **Include AI attribution** for AI-generated code

### Pull Requests

1. **Validate locally** before pushing
2. **Keep PRs focused** (one feature/fix)
3. **Update PR title** to match conventions
4. **Squash commits** on merge (if using squash merge)
5. **Document breaking changes** in PR description

### Releases

1. **Only merge to main** when ready to release
2. **Review CHANGELOG** after release
3. **Test binaries** from GitHub releases
4. **Document breaking changes** in release notes
5. **Notify team** of major releases

---

## Makefile Commands

```bash
# Install Node.js dependencies
npm install

# Install git hooks
make install-git-hooks

# Validate last commit
npx commitlint --from HEAD~1

# Run semantic-release locally (dry run)
npx semantic-release --dry-run

# Check what would be released
npx semantic-release --dry-run --branches main
```

---

## Reference

### Links

- **Conventional Commits**: https://www.conventionalcommits.org/
- **Commitlint**: https://commitlint.js.org/
- **Semantic Release**: https://semantic-release.gitbook.io/
- **GitHub Actions**: https://docs.github.com/actions

### Files

- `.commitlintrc.json` - Commitlint configuration
- `.releaserc.json` - Semantic release configuration
- `package.json` - Node.js dependencies
- `.github/workflows/ci.yml` - CI workflow
- `.github/workflows/release.yml` - Release workflow
- `.github/workflows/pr-validation.yml` - PR validation workflow

---

## Summary

✅ **Commitlint** enforces conventional commit format
✅ **Semantic Release** automates versioning and releases
✅ **GitHub Actions** runs CI/CD pipelines
✅ **Release only from main** branch
✅ **Binaries auto-uploaded** to GitHub releases
✅ **CHANGELOG auto-generated** from commits

**Start using**: Commit with conventional format, push to main to release!

