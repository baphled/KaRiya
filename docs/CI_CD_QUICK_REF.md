# CI/CD Quick Reference

## Setup (One-Time)

```bash
# Install Node.js dependencies
npm install

# Install git hooks
make install-git-hooks

# Verify setup
npx commitlint --help
```

---

## Commit Message Format

### Template

```
<type>(<scope>): <subject>

<body>

AI-Generated-By: <Assistant> (<Model>) [if AI-generated]
Reviewed-By: <Your Name> [if AI-generated]

<footer>
```

### Valid Types

| Type | Description | Release |
|------|-------------|---------|
| `feat` | New feature | Minor (0.X.0) |
| `fix` | Bug fix | Patch (0.0.X) |
| `perf` | Performance | Patch (0.0.X) |
| `refactor` | Refactoring | Patch (0.0.X) |
| `docs` | Documentation | None |
| `style` | Formatting | None |
| `test` | Tests | None |
| `chore` | Maintenance | None |
| `ci` | CI/CD | None |
| `build` | Build system | Patch (0.0.X) |
| `revert` | Revert commit | Patch (0.0.X) |

### Valid Scopes

- `domain` - Domain layer
- `service` - Service layer
- `repo` - Repository layer
- `cli` - CLI interface
- `logger` - Logging system
- `classification` - Classification service
- `deps` - Dependencies
- `release` - Release-related

---

## Examples

### Feature (Minor Release)

```
feat(service): add event filtering by date range

Implement date range filtering for timeline views.
Improves performance for large event collections.

Closes #56
```

### Bug Fix (Patch Release)

```
fix(domain): prevent duplicate tags in validation

Tag deduplication was not working correctly when
tags had different casing.

Fixes #87
```

### Breaking Change (Major Release)

```
feat(api): redesign event structure

Complete redesign of event data structure for
better performance and flexibility.

BREAKING CHANGE: Event structure has changed.
Migration guide available in docs/MIGRATION.md
```

### Documentation (No Release)

```
docs(readme): update installation instructions

Add Node.js requirement for CI/CD setup.
```

### With AI Attribution

```
feat(cli): add interactive event capture form

Implement BubbleTea form with real-time validation.

Human review confirmed:
- All tests pass (14/14)
- No security issues
- Follows project conventions

AI-Generated-By: Avante (Claude 3.5 Sonnet)
Reviewed-By: John Doe
Closes #42
```

---

## Versioning

### How Versions are Determined

```
v1.2.3
│ │ │
│ │ └─ Patch: fix, perf, refactor, build, revert
│ └─── Minor: feat
└───── Major: BREAKING CHANGE
```

### Examples

```
v0.1.0 → feat(service): add feature → v0.2.0
v0.2.0 → fix(domain): fix bug    → v0.2.1
v0.2.1 → feat + BREAKING CHANGE  → v1.0.0
```

---

## Workflows

### Feature Branch → Main

```bash
# 1. Create feature branch
git checkout -b feature/event-filtering

# 2. Make changes and commit
git add <files>
git commit
# Write message following format

# 3. Push to remote
git push origin feature/event-filtering

# 4. Create PR
gh pr create --title "feat(service): add event filtering"

# 5. PR Validation runs:
#    - Commitlint
#    - Tests
#    - Security scan
#    - AI attribution check

# 6. Merge PR to main
gh pr merge --squash

# 7. Release workflow runs automatically:
#    - Determines version
#    - Generates CHANGELOG
#    - Creates GitHub release
#    - Uploads binaries
```

---

## Local Commands

```bash
# Validate last commit
npx commitlint --from HEAD~1

# Validate commit range
npx commitlint --from origin/main

# Validate specific message
echo "feat(service): add feature" | npx commitlint

# Dry-run semantic-release
npx semantic-release --dry-run

# Check what would be released
git log $(git describe --tags --abbrev=0)..HEAD --oneline
```

---

## Troubleshooting

### ❌ Commitlint Fails

```bash
# Check format
✅ feat(service): add feature
❌ feat: add feature          # Missing scope
❌ Feature: add feature       # Wrong case
❌ feat(invalid): add feature # Invalid scope

# Check length
# Header max: 100 chars
# Subject max: 72 chars
```

### ❌ No Release Created

```bash
# Reasons:
# 1. Only non-release commits (docs, test, style)
# 2. No new commits since last release
# 3. [skip ci] in commit message

# Check releasable commits:
git log $(git describe --tags --abbrev=0)..HEAD --grep="^feat\|^fix"
```

### ❌ PR Validation Fails

```bash
# Fix PR title
# Must follow: <type>(<scope>): <subject>

# Fix commit messages
git rebase -i origin/main
# Edit messages to follow format

# Add breaking change if needed
git commit --amend
# Add: BREAKING CHANGE: description
```

---

## GitHub Actions Status

### CI Workflow (All Branches)

- ✅ Commitlint
- ✅ Lint & Format
- ✅ Tests (Ubuntu, macOS, Windows)
- ✅ Build
- ✅ Security Scan

### Release Workflow (Main Only)

- ✅ CI Checks
- ✅ Build Binaries
- ✅ Semantic Release
- ✅ Create GitHub Release
- ✅ Upload Binaries

### PR Validation

- ✅ PR Title Validation
- ✅ Conventional Commits
- ⚠️  AI Attribution Check (warning)
- ✅ Breaking Changes Detection
- ✅ Size Label

---

## Quick Rules

1. ✅ **Use conventional commits** for all commits
2. ✅ **Include scope** when applicable
3. ✅ **Add AI attribution** for AI-generated code
4. ✅ **Document breaking changes** in footer
5. ✅ **Only merge to main** when ready to release
6. ✅ **Let semantic-release** handle versioning

---

## Cheat Sheet

```bash
# Daily workflow
git checkout -b feature/my-feature
# ... make changes ...
git commit -m "feat(service): add new feature"
git push origin feature/my-feature
gh pr create
# ... wait for CI ...
gh pr merge --squash
# 🎉 Automatic release!

# Check version
cat VERSION

# View changelog
cat CHANGELOG.md

# Download binaries
gh release download
```

---

## Resources

- **Full Documentation**: `docs/CI_CD_PIPELINE.md`
- **Conventional Commits**: https://www.conventionalcommits.org/
- **Commitlint**: https://commitlint.js.org/
- **Semantic Release**: https://semantic-release.gitbook.io/

---

**Print this and keep it handy!**

