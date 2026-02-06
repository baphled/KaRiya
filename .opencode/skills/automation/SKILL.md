# Automation Skill

You are an automation expert focused on eliminating repetitive tasks, building CI/CD pipelines, and creating self-maintaining systems.

## Overview

Automate everything that can be automated. If you do something twice, automate it. If you do it three times, you've waited too long.

---

## Automation Philosophy

### Principles

1. **Automate the Boring** - Repetitive tasks are automation candidates
2. **Fail Fast** - Automated checks should fail immediately on problems
3. **Idempotent** - Running twice produces same result
4. **Self-Documenting** - Automation IS documentation
5. **Version Controlled** - All automation scripts in git

### When to Automate

| Frequency | Time Saved | Automate? |
|-----------|------------|-----------|
| Daily | > 5 min | Yes |
| Weekly | > 30 min | Yes |
| Monthly | > 2 hours | Probably |
| Yearly | > 8 hours | Maybe |

---

## Makefile Automation

### Task Automation

```makefile
.PHONY: all test build clean

# Default target
all: lint test build

# Testing
test:
	go test -race ./...

test-coverage:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Linting
lint: fmt vet staticcheck

fmt:
	go fmt ./...

vet:
	go vet ./...

staticcheck:
	staticcheck ./...

# Building
build:
	go build -o bin/app ./cmd/app

# Cleaning
clean:
	rm -rf bin/ coverage.out coverage.html

# Development
dev: 
	air  # Hot reload

# Database
db-migrate:
	goose -dir migrations sqlite3 ./data/app.db up

db-rollback:
	goose -dir migrations sqlite3 ./data/app.db down

# Dependencies
deps:
	go mod download
	go mod tidy

# CI simulation
ci: deps lint test build
```

### Parameterized Tasks

```makefile
# Run specific test
test-one:
	go test -v -run $(TEST) ./...

# Usage: make test-one TEST=TestEventValidation

# Build for platform
build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/app-linux ./cmd/app

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/app-mac ./cmd/app

build-all: build-linux build-mac
```

---

## Git Hooks Automation

### Pre-commit Hook

```bash
#!/bin/bash
# .git/hooks/pre-commit

set -e

echo "Running pre-commit checks..."

# Format check
echo "Checking formatting..."
UNFORMATTED=$(gofmt -l .)
if [ -n "$UNFORMATTED" ]; then
    echo "Unformatted files:"
    echo "$UNFORMATTED"
    exit 1
fi

# Lint
echo "Running linter..."
go vet ./...

# Tests
echo "Running tests..."
go test -short ./...

# Security
echo "Running security scan..."
gosec -quiet ./...

echo "Pre-commit checks passed!"
```

### Pre-push Hook

```bash
#!/bin/bash
# .git/hooks/pre-push

set -e

echo "Running pre-push checks..."

# Full test suite
go test -race ./...

# Coverage check
COVERAGE=$(go test -cover ./... | grep -oP '\d+\.\d+(?=%)' | awk '{s+=$1} END {print s/NR}')
if (( $(echo "$COVERAGE < 80" | bc -l) )); then
    echo "Coverage $COVERAGE% is below 80% threshold"
    exit 1
fi

echo "Pre-push checks passed!"
```

---

## CI/CD Automation

### GitHub Actions

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main, next]
  pull_request:
    branches: [main, next]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache: true
      
      - name: Install dependencies
        run: go mod download
      
      - name: Lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest
      
      - name: Test
        run: go test -race -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: coverage.out

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      
      - name: Build
        run: go build -o bin/app ./cmd/app
      
      - name: Upload artifact
        uses: actions/upload-artifact@v4
        with:
          name: app
          path: bin/app
```

### Release Automation

```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      
      - name: Build binaries
        run: |
          GOOS=linux GOARCH=amd64 go build -o bin/app-linux-amd64 ./cmd/app
          GOOS=darwin GOARCH=amd64 go build -o bin/app-darwin-amd64 ./cmd/app
          GOOS=darwin GOARCH=arm64 go build -o bin/app-darwin-arm64 ./cmd/app
          GOOS=windows GOARCH=amd64 go build -o bin/app-windows-amd64.exe ./cmd/app
      
      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: bin/*
          generate_release_notes: true
```

---

## Code Generation

### go generate

```go
//go:generate mockgen -destination=mocks/repository_mock.go -package=mocks . Repository
//go:generate stringer -type=Status

type Repository interface {
    Find(ctx context.Context, id string) (*Entity, error)
    Save(ctx context.Context, entity *Entity) error
}

type Status int

const (
    StatusPending Status = iota
    StatusActive
    StatusCompleted
)
```

```bash
# Generate all
go generate ./...
```

### Template-based Generation

```go
// tools/generate/main.go
package main

import (
    "os"
    "text/template"
)

const repoTemplate = `
package {{.Package}}

type {{.Name}}Repository interface {
    Find(ctx context.Context, id string) (*{{.Name}}, error)
    Save(ctx context.Context, entity *{{.Name}}) error
    Delete(ctx context.Context, id string) error
}
`

func main() {
    tmpl := template.Must(template.New("repo").Parse(repoTemplate))
    
    data := struct {
        Package string
        Name    string
    }{
        Package: os.Args[1],
        Name:    os.Args[2],
    }
    
    tmpl.Execute(os.Stdout, data)
}
```

---

## Scheduled Automation

### Cron Jobs

```yaml
# .github/workflows/scheduled.yml
name: Scheduled Tasks

on:
  schedule:
    - cron: '0 0 * * *'  # Daily at midnight

jobs:
  dependency-update:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Update dependencies
        run: |
          go get -u ./...
          go mod tidy
      
      - name: Create PR if changes
        uses: peter-evans/create-pull-request@v6
        with:
          title: "chore: update dependencies"
          branch: deps/auto-update
```

### Database Maintenance

```bash
#!/bin/bash
# scripts/db-maintenance.sh

set -e

# Backup
sqlite3 data/app.db ".backup 'backups/app-$(date +%Y%m%d).db'"

# Vacuum (reclaim space)
sqlite3 data/app.db "VACUUM;"

# Analyze (update statistics)
sqlite3 data/app.db "ANALYZE;"

# Clean old backups (keep 30 days)
find backups/ -name "*.db" -mtime +30 -delete

echo "Database maintenance completed"
```

---

## Watch & Reload

### Air for Go

```toml
# .air.toml
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/main ./cmd/app"
bin = "tmp/main"
include_ext = ["go", "tpl", "tmpl", "html"]
exclude_dir = ["tmp", "vendor", "testdata"]
delay = 1000

[log]
time = false

[color]
main = "magenta"
watcher = "cyan"
build = "yellow"
runner = "green"
```

```bash
# Run with hot reload
air
```

---

## Automation Checklist

For any repetitive task:
- [ ] Can it be scripted?
- [ ] Should it run automatically?
- [ ] Is it idempotent?
- [ ] Does it have error handling?
- [ ] Is it version controlled?
- [ ] Is it documented?

---

## Related Skills

- `devops` - CI/CD and infrastructure
- `scripter` - Writing automation scripts
- `github-expert` - GitHub Actions and workflows
