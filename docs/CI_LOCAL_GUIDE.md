---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# CI Local Checks Guide

This guide explains how to run all CI checks locally before pushing to GitHub.

## Quick Start

**Run all CI checks (mirrors GitHub Actions exactly):**

```bash
make ci-local
```

This single command runs:
- ✅ Commitlint validation
- ✅ AI attribution check
- ✅ Code formatting (go fmt)
- ✅ Static analysis (go vet)
- ✅ Advanced linting (staticcheck)
- ✅ Tests with race detector and coverage
- ✅ Multi-platform builds (Linux, macOS, Windows)
- ✅ Security scanning (gosec)

## Installation

**Install all required tools:**

```bash
make ci-install-tools
```

This installs:
- `ginkgo` - Test framework
- `staticcheck` - Advanced Go linter
- `gosec` - Security scanner
- `npm` dependencies - For commitlint

## Individual CI Checks

Run specific checks individually:

### 1. Format Check
```bash
make fmt
```
Checks if all Go files are properly formatted.

**CI equivalent**: `.github/workflows/ci.yml` - lint job - go fmt step

### 2. Static Analysis
```bash
make vet
```
Runs `go vet` for static analysis.

**CI equivalent**: `.github/workflows/ci.yml` - lint job - go vet step

### 3. Advanced Linting
```bash
make staticcheck
```
Runs `staticcheck` for advanced linting.

**CI equivalent**: `.github/workflows/ci.yml` - lint job - staticcheck step

### 4. Security Scanning
```bash
make gosec
```
Runs security scanner to detect vulnerabilities.

**CI equivalent**: `.github/workflows/ci.yml` - security job

### 5. Tests
```bash
make test
```
Runs all tests with race detector.

**CI equivalent**: `.github/workflows/ci.yml` - test job

**With coverage:**
```bash
ginkgo -v --race --cover --coverprofile=coverage.out ./...
```

### 6. Build
```bash
make build
```
Builds for current platform.

**Multi-platform builds:**
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

**CI equivalent**: `.github/workflows/ci.yml` - build job

### 7. Commit Message Validation
```bash
# Check last commit message
git log -1 --pretty=%B | npx commitlint --config .commitlintrc.json
```

**CI equivalent**: `.github/workflows/pr-validation.yml` - conventional-commits job

### 8. AI Attribution Check
```bash
make check-ai-attribution
```

**CI equivalent**: `.github/workflows/pr-validation.yml` - check-ai-attribution job

## CI Workflows

### Main CI Workflow (`.github/workflows/ci.yml`)

Runs on push to `main`, `develop`, feature/fix/refactor branches:

1. **commitlint** (PR only) - Validates commit messages
2. **lint** - Code formatting and linting
3. **test** - Multi-platform testing (Linux, macOS, Windows)
4. **build** - Multi-platform builds
5. **security** - Security scanning with gosec

### PR Validation Workflow (`.github/workflows/pr-validation.yml`)

Additional checks for pull requests:

1. **validate-pr-title** - PR title follows conventional commits
2. **check-ai-attribution** - Commits have AI attribution (if applicable)
3. **conventional-commits** - All commits follow conventions
4. **breaking-changes** - Detects breaking changes
5. **size-label** - Auto-labels PR size

## Pre-Push Checklist

Before pushing code, ensure:

```bash
# 1. Run all CI checks locally
make ci-local

# 2. If all pass, you're good to push!
git push
```

## Quick Commands Comparison

| What you want | Local command | CI job |
|---------------|---------------|--------|
| All checks | `make ci-local` | All workflows |
| Format check | `make fmt` | lint job |
| Linting | `make vet && make staticcheck` | lint job |
| Security | `make gosec` | security job |
| Tests | `make test` | test job |
| Build | `make build` | build job |
| Commit check | `git log -1 \| npx commitlint` | commitlint job |

## Troubleshooting

### "ginkgo: command not found"
```bash
make ci-install-tools
```

### "staticcheck: command not found"
```bash
go install honnef.co/go/tools/cmd/staticcheck@latest
```

### "gosec: command not found"
```bash
go install github.com/securego/gosec/v2/cmd/gosec@latest
```

### "npx: command not found"
Install Node.js 18+ and npm, then:
```bash
npm ci
```

### Tests failing locally but not in CI
- Ensure you're using Go 1.24 (same as CI)
- Check `go.mod` is up to date: `go mod tidy`
- Clear build cache: `go clean -cache`

### Build failing for specific platform
- Verify GOOS/GOARCH values
- Check for platform-specific code
- Ensure CGO is disabled for cross-compilation: `CGO_ENABLED=0`

## Advanced Usage

### Run CI checks with verbose output
```bash
bash -x scripts/ci-local.sh
```

### Run only failed checks
The `ci-local.sh` script continues on failure and shows a summary at the end.

### Skip specific checks
Edit `scripts/ci-local.sh` and comment out the `run_check` calls you want to skip.

### Custom test timeout
```bash
ginkgo -v --race --timeout=10m ./...
```

## CI Configuration Files

- **Main CI**: `.github/workflows/ci.yml`
- **PR Validation**: `.github/workflows/pr-validation.yml`
- **Commitlint Config**: `.commitlintrc.json`
- **Local CI Script**: `scripts/ci-local.sh`

## See Also

- [Commit Guidelines](rules/COMMIT_QUICK_REFERENCE.md)
- [AI Attribution](rules/AI_COMMIT_ATTRIBUTION.md)
- [Master Task Prompt](rules/master-task-prompt.md)
- [Check Compliance](rules/rules-compliance-check.md)
