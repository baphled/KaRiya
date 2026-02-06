---
name: release-management
description: Versioning, changelogs, release notes, and release branch management
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
  always_active: false
---

## What I do

Guide the release process including versioning, changelog generation, release notes, and managing the `next` to `main` promotion.

## When to use me

Use this skill when:
- Preparing a release from `next` to `main`
- Creating version tags
- Writing release notes
- Managing semantic versioning

## KaRiya Release Flow

```
feature/* -> next (integration) -> main (production)
                                      |
                                      v
                                   v1.x.x (tag)
```

## Semantic Versioning

KaRiya follows [SemVer](https://semver.org/):

| Change Type | Version Bump | Examples |
|-------------|--------------|----------|
| Breaking API changes | MAJOR (X.0.0) | Remove public function, change signature |
| New features (backwards compatible) | MINOR (0.X.0) | Add new intent, new command |
| Bug fixes, patches | PATCH (0.0.X) | Fix crash, correct behaviour |

### Determining Version Bump

```bash
# View commits since last release
git log $(git describe --tags --abbrev=0)..HEAD --oneline

# Check for breaking changes
git log $(git describe --tags --abbrev=0)..HEAD --oneline | grep -i "BREAKING"

# Check commit types
git log $(git describe --tags --abbrev=0)..HEAD --format="%s" | cut -d: -f1 | sort | uniq -c
```

**Decision Matrix:**

| Commits Include | Version Bump |
|-----------------|--------------|
| `feat!:` or `BREAKING CHANGE:` | MAJOR |
| `feat:` (no breaking) | MINOR |
| Only `fix:`, `docs:`, `chore:` | PATCH |

## Release Checklist

### 1. Pre-Release Validation

```bash
# Ensure next is stable
git checkout next
git pull origin next

# Run full validation
make check-compliance
make test
make coverage

# Check CI status on next
gh run list --branch next --limit 5
```

### 2. Changelog Generation

Create or update `CHANGELOG.md`:

```markdown
# Changelog

All notable changes to KaRiya will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.2.0] - 2026-02-06

### Added
- New skill management intent (#142)
- OpenCode commands for PR monitoring (#155)

### Changed
- Improved table navigation performance (#148)

### Fixed
- Modal overlay rendering on small terminals (#151)
- Form validation for empty fields (#153)

### Security
- Updated dependencies for CVE-XXXX-YYYY (#150)

## [1.1.0] - 2026-01-15
...
```

**Generate changelog from commits:**

```bash
# Get commits since last tag
LAST_TAG=$(git describe --tags --abbrev=0)
git log $LAST_TAG..HEAD --pretty=format:"- %s (%h)" | sort
```

### 3. Create Release PR

```bash
# Create release branch (optional for complex releases)
git checkout -b release/v1.2.0

# Or merge directly from next
git checkout main
git pull origin main
```

**Release PR:**

```bash
gh pr create \
  --base main \
  --head next \
  --title "Release v1.2.0" \
  --body "$(cat <<'EOF'
## Release v1.2.0

### Summary
[Brief description of this release]

### Changes
See [CHANGELOG.md](CHANGELOG.md) for full details.

### Pre-Release Checklist
- [ ] All CI checks pass on `next`
- [ ] CHANGELOG.md updated
- [ ] Version bumped in relevant files
- [ ] No breaking changes without migration guide
- [ ] Security vulnerabilities addressed

### Post-Merge Actions
- [ ] Create git tag v1.2.0
- [ ] Create GitHub release
- [ ] Update documentation
EOF
)"
```

### 4. Version Tagging

After PR is merged to main:

```bash
git checkout main
git pull origin main

# Create annotated tag
git tag -a v1.2.0 -m "Release v1.2.0

Changes:
- New skill management intent
- OpenCode commands for PR monitoring
- Fixed modal overlay rendering
"

# Push tag
git push origin v1.2.0
```

### 5. GitHub Release

```bash
# Create release from tag
gh release create v1.2.0 \
  --title "v1.2.0" \
  --notes "$(cat <<'EOF'
## What's New

### Features
- **Skill Management**: New intent for managing OpenCode skills (#142)
- **PR Monitoring**: Commands for tracking PR status and reviews (#155)

### Improvements
- Table navigation is now 40% faster (#148)

### Bug Fixes
- Fixed modal overlay rendering on small terminals (#151)
- Fixed form validation for empty fields (#153)

### Security
- Updated go-yaml to address CVE-XXXX-YYYY (#150)

## Upgrade Guide

This release is backwards compatible. Update with:
\`\`\`bash
go get github.com/baphled/kariya@v1.2.0
\`\`\`

## Contributors
- @baphled
- Claude Code (AI-assisted development)

**Full Changelog**: https://github.com/baphled/KaRiya/compare/v1.1.0...v1.2.0
EOF
)"
```

## Release Notes Template

```markdown
## v[X.Y.Z] Release Notes

**Release Date:** [YYYY-MM-DD]

### Highlights

[2-3 sentences about the most important changes]

### Breaking Changes

[List any breaking changes with migration instructions]

### New Features

- **[Feature Name]**: [Description] (#PR)

### Improvements

- [Description] (#PR)

### Bug Fixes

- Fixed [issue description] (#PR)

### Security

- [Security fix description] (#PR)

### Deprecations

- `[function/feature]` is deprecated, use `[replacement]` instead

### Dependencies

- Updated [package] from vX.X.X to vY.Y.Y

### Contributors

- @username
- Claude Code (AI-assisted)

### Upgrade Instructions

[Step-by-step upgrade guide if needed]

**Full Changelog:** [compare link]
```

## Hotfix Process

For critical fixes that can't wait for next release:

```bash
# Create hotfix branch from main
git checkout main
git checkout -b hotfix/critical-fix

# Make fix
# ...

# Create PR to main (exception to normal flow)
gh pr create --base main --title "hotfix: critical security fix"

# After merge, tag immediately
git checkout main
git pull
git tag -a v1.2.1 -m "Hotfix: critical security fix"
git push origin v1.2.1

# Backport to next
git checkout next
git cherry-pick <hotfix-commit>
git push origin next
```

## Version Files

If KaRiya has version embedded in code:

```go
// internal/version/version.go
package version

const (
    Version = "1.2.0"
    Commit  = "" // Set at build time
)
```

Update before release:
```bash
# Update version
sed -i 's/Version = "[^"]*"/Version = "1.2.0"/' internal/version/version.go

# Commit version bump
make ai-commit FILE=/tmp/version-bump.txt
```

## Related Skills

- `pre-merge` - Final validation before merge
- `breaking-changes` - Managing backwards compatibility
- `documentation-writing` - Release documentation
- `git-advanced` - Tagging and branch management
