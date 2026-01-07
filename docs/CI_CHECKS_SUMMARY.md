# CI Checks - Complete Summary

This document provides a complete mapping between GitHub Actions CI checks and local commands.

## Overview

All CI checks can now be run locally using:

```bash
make ci-local
```

This ensures your code will pass CI before you push.

## CI Workflows

### 1. Main CI Workflow (`.github/workflows/ci.yml`)

**Triggers**: Push to `main`, `develop`, `feature/**`, `fix/**`, `refactor/**` branches and PRs

#### Job: commitlint (PR only)
**Purpose**: Validates commit message format

**CI Command**:
```bash
npx commitlint --from $BASE_SHA --to $HEAD_SHA --verbose
```

**Local Command**:
```bash
# Check last commit
git log -1 --pretty=%B | npx commitlint --config .commitlintrc.json

# Check range of commits (e.g., last 5)
git log --format=%B -5 | npx commitlint --config .commitlintrc.json
```

**Configuration**: `.commitlintrc.json`

---

#### Job: lint
**Purpose**: Code formatting and static analysis

**Steps**:

1. **go fmt check**
   - **CI**: `gofmt -l .` (fails if any files need formatting)
   - **Local**: `make fmt` or `gofmt -l .`
   - **Fix**: `go fmt ./...`

2. **go vet**
   - **CI**: `go vet ./...`
   - **Local**: `make vet` or `go vet ./...`

3. **staticcheck**
   - **CI**: `staticcheck ./...`
   - **Local**: `make staticcheck` or `staticcheck ./...`
   - **Install**: `go install honnef.co/go/tools/cmd/staticcheck@latest`

---

#### Job: test
**Purpose**: Run tests with race detector and coverage

**Platforms**: Ubuntu, macOS, Windows

**CI Command**:
```bash
ginkgo -v --race --cover --coverprofile=coverage.out ./...
```

**Local Command**:
```bash
make test
# OR
ginkgo -v --race --cover --coverprofile=coverage.out ./...
```

**Coverage Upload**: Only on Ubuntu (to Codecov)

---

#### Job: build
**Purpose**: Multi-platform binary builds

**CI Commands**:
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o kariya-linux-amd64 ./cmd/cli

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o kariya-darwin-amd64 ./cmd/cli

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o kariya-darwin-arm64 ./cmd/cli

# Windows
GOOS=windows GOARCH=amd64 go build -o kariya-windows-amd64.exe ./cmd/cli
```

**Local Command**:
```bash
# Current platform only
make build

# All platforms (run ci-local or manually)
GOOS=linux GOARCH=amd64 go build -o kariya-linux-amd64 ./cmd/cli
GOOS=darwin GOARCH=amd64 go build -o kariya-darwin-amd64 ./cmd/cli
GOOS=darwin GOARCH=arm64 go build -o kariya-darwin-arm64 ./cmd/cli
GOOS=windows GOARCH=amd64 go build -o kariya-windows-amd64.exe ./cmd/cli
```

**Artifacts**: Uploaded to GitHub with 7-day retention

---

#### Job: security
**Purpose**: Security vulnerability scanning

**CI Command**:
```bash
gosec -no-fail -fmt sarif -out gosec.sarif ./...
```

**Local Command**:
```bash
make gosec
# OR
gosec -no-fail -fmt text ./...
```

**Install**:
```bash
go install github.com/securego/gosec/v2/cmd/gosec@latest
```

**Upload**: SARIF file to GitHub Code Scanning

---

### 2. PR Validation Workflow (`.github/workflows/pr-validation.yml`)

**Triggers**: PR opened, synchronized, or reopened

#### Job: validate-pr-title
**Purpose**: Ensure PR title follows conventional commits

**CI Command**:
```bash
echo "$PR_TITLE" | npx commitlint --config .commitlintrc.json
```

**Local Command**:
```bash
# Check if a string follows commit conventions
echo "feat(cli): add new feature" | npx commitlint --config .commitlintrc.json
```

---

#### Job: check-ai-attribution
**Purpose**: Verify commits have AI attribution (if applicable)

**What it checks**:
- All commits with code files (`.go`, `.js`, `.ts`, `.py`, etc.)
- Looks for `AI-Generated-By:` in commit message

**CI Command**:
```bash
# For each commit in PR
git log --format=%B -n 1 $commit | grep -q "AI-Generated-By:"
```

**Local Command**:
```bash
make check-ai-attribution
# OR
git log -1 --pretty=%B | grep "AI-Generated-By:"
```

**Expected Format**:
```
AI-Generated-By: Claude 3.7 Sonnet (claude-3-5-sonnet-20241022)
Reviewed-By: Human
```

---

#### Job: conventional-commits
**Purpose**: Validate all commits in PR follow conventions

**CI Command**:
```bash
npx commitlint --from $BASE_SHA --to $HEAD_SHA --verbose
```

**Local Command**:
```bash
# Check commits in current branch vs main
git log --format=%B origin/main..HEAD | npx commitlint --config .commitlintrc.json
```

---

#### Job: breaking-changes
**Purpose**: Detect and flag breaking changes

**What it checks**: Commit messages for `BREAKING CHANGE:`

**CI Command**:
```bash
git log $BASE_SHA..$HEAD_SHA --format=%B | grep -qi "BREAKING CHANGE"
```

**Local Command**:
```bash
# Check if any recent commits have breaking changes
git log -10 --format=%B | grep -i "BREAKING CHANGE"
```

---

#### Job: size-label
**Purpose**: Auto-label PR by size

**Sizes**:
- `size/xs`: 0-10 lines
- `size/s`: 11-100 lines
- `size/m`: 101-500 lines
- `size/l`: 501-1000 lines
- `size/xl`: 1000+ lines

**No local equivalent** (GitHub Action only)

---

## Quick Reference

### Install All Tools

```bash
make ci-install-tools
```

Installs:
- `ginkgo` - Test runner
- `staticcheck` - Advanced linter
- `gosec` - Security scanner
- `npm` dependencies - For commitlint

### Run All CI Checks

```bash
make ci-local
```

Runs (in order):
1. ✅ Commitlint (last commit)
2. ✅ AI attribution check
3. ✅ go fmt
4. ✅ go vet
5. ✅ staticcheck
6. ✅ Tests with race detector & coverage
7. ✅ Multi-platform builds
8. ✅ Security scan (gosec)

### Individual Checks

| Check | Command |
|-------|---------|
| Format | `make fmt` |
| Vet | `make vet` |
| Lint | `make staticcheck` |
| Security | `make gosec` |
| Tests | `make test` |
| Coverage | `make coverage` |
| Build | `make build` |
| All Checks | `make ci-local` |

### Pre-Push Workflow

```bash
# 1. Format code
go fmt ./...

# 2. Run all CI checks
make ci-local

# 3. If all pass, push!
git push
```

## Configuration Files

| File | Purpose |
|------|---------|
| `.github/workflows/ci.yml` | Main CI workflow |
| `.github/workflows/pr-validation.yml` | PR validation workflow |
| `.commitlintrc.json` | Commitlint configuration |
| `scripts/ci-local.sh` | Local CI check script |
| `Makefile` | Convenient make targets |

## Commit Message Format

**Format**: `<type>(<scope>): <subject>`

**Types**:
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation
- `style` - Formatting
- `refactor` - Code restructuring
- `test` - Tests
- `chore` - Maintenance
- `perf` - Performance
- `ci` - CI/CD
- `build` - Build system
- `revert` - Revert commit

**Scopes**:
- `domain` - Domain models
- `service` - Service layer
- `repo` - Repository layer
- `cli` - CLI/TUI
- `logger` - Logging
- `classification` - Classification
- `deps` - Dependencies
- `release` - Release
- `lint` - Linting
- `security` - Security
- `workflow` - Workflows
- `config` - Configuration
- `docs` - Documentation
- `tests` - Tests
- `intents` - Intent system
- `models` - UI models
- `components` - UI components

**Examples**:
```
feat(cli): add burst detection to event capture
fix(service): handle nil pointer in CV generation
docs(readme): update installation instructions
chore(deps): bump ginkgo to v2.14.0
```

## Troubleshooting

### "ginkgo: command not found"
```bash
go install github.com/onsi/ginkgo/v2/ginkgo@latest
```

### "staticcheck: command not found"
```bash
go install honnef.co/go/tools/cmd/staticcheck@latest
```

### "gosec: command not found"
```bash
go install github.com/securego/gosec/v2/cmd/gosec@latest
```

### "commitlint: command not found"
```bash
npm ci
```

### Tests pass locally but fail in CI
- Check Go version (must be 1.24)
- Run `go mod tidy`
- Run with race detector: `go test -race ./...`
- Check platform-specific code (Windows/macOS differences)

### Staticcheck issues
```bash
# View all issues
staticcheck ./...

# Ignore specific checks (add to code)
//nolint:staticcheck
```

### Commitlint failures
```bash
# Check commit message format
git log -1 --pretty=%B | npx commitlint --config .commitlintrc.json

# Fix with interactive rebase
git rebase -i HEAD~1
```

## See Also

- [CI Local Guide](CI_LOCAL_GUIDE.md) - Detailed guide for running CI locally
- [Commit Quick Reference](rules/COMMIT_QUICK_REFERENCE.md) - Commit message templates
- [AI Commit Attribution](rules/AI_COMMIT_ATTRIBUTION.md) - AI attribution requirements
- [Master Task Prompt](rules/master-task-prompt.md) - Development workflow
