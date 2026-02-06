---
name: check-compliance
description: Run full KaRiya compliance checks before and after making changes
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Run comprehensive compliance checks to validate code quality, architecture, and documentation.

## When to use me

**BEFORE and AFTER every task.** This is mandatory.

## Quick Check

```bash
make check-compliance
```

This runs:
- `staticcheck` - Static analysis
- `check-intent-architecture` - Architecture validation
- `validate-documentation` - Doc requirements
- `check-fixtures` - Test fixture enforcement
- Full compliance script

## Individual Checks

### Architecture Validation
```bash
make check-intent-architecture
```
Validates:
- Intent structure (5 core files)
- Type locations
- Handler organization
- File size limits

### Documentation Check
```bash
make check-docblocks
```
Validates:
- `doc.go` files exist
- Exported functions have Required/Returns/Side effects sections
- Types have godoc comments

### Pattern Enforcement
```bash
make check-patterns        # Quick check
make check-patterns-strict # Blocking check
```
Checks:
- Form wrapper pattern
- BaseIntent embedding
- Theme consistency
- UIKit usage

### Static Analysis
```bash
make staticcheck
make vet
make gosec          # Security scanning
make golangci-lint  # Comprehensive linting
```

## CI Checks Locally

```bash
make ci-local
```

Runs ALL CI checks locally, including:
- Format check
- Vet
- Staticcheck
- Gosec
- All tests
- Coverage

## Pre-PR Validation

```bash
make pre-pr
```

Validates:
- Not on `main` or `next` branch
- All CI checks pass
- Ready for PR to `next`

## Coverage Requirements

| Scope | Threshold | Enforcement |
|-------|-----------|-------------|
| Per-package (modified) | >= 95% | Pre-commit (BLOCKING) |
| Project average | >= 80% | Warning |

Check coverage:
```bash
go test -cover ./internal/cli/intents/myfeature/...
```

## Documentation Requirements

Every exported symbol needs:

**Functions:**
```go
// MyFunction does something.
//
// Expected: param1 (required), param2 (optional)
// Returns: result description
// Side effects: None (or describe effects)
func MyFunction(param1 string, param2 int) error {
```

**Types:**
```go
// MyType represents something important.
//
// It maintains state for...
type MyType struct {
```

**Packages:**
```go
// Package mypackage provides...
//
// # Overview
//
// Description...
package mypackage
```

## Fixing Documentation

```bash
make fix-docs FILE=path/to/file.go
```

Or for all files:
```bash
make fix-all-documentation
```

## Skip Documentation (MANDATORY)

If ANY check must be skipped, document with full reasoning:

```
SKIPPING: make check-patterns
REASON: Changes are markdown-only, no Go code modified
IMPACT: None - pattern checks only apply to Go code
```

### Required Format
```
SKIPPING: [Check being skipped]
REASON: [Specific reason why skip is justified]
IMPACT: [What risks this introduces]
INSTEAD: [Alternative validation performed, if any]
```

### Acceptable Skip Reasons
| Check | Valid Skip Reason |
|-------|-------------------|
| `check-patterns` | Documentation-only changes |
| `check-docblocks` | Test files only (but rare) |
| `staticcheck` | NEVER skip |
| `check-intent-architecture` | Non-intent changes |
| `gosec` | NEVER skip |
| `coverage` | Trivial getters/setters only |

### Unacceptable Skips
- "Takes too long"
- "Probably fine"
- "Will fix later"
- "Not important"
- No reason given (VIOLATION)

**Silent skips are violations.** Every skip must be documented.

## Related skills

- `session-start` - Initialize session with compliance check
- `fix-architecture` - Fix detected violations
- `ai-commit` - Commit after compliance passes
- `checklist-discipline` - Progress tracking and skip documentation
