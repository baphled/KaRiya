# BDD Workflow Guide

This guide defines the Behavior-Driven Development (BDD) workflow for KaRiya development.

## When TDD Applies

TDD is **required** for:
- New features (Go code)
- Bug fixes (Go code)
- Refactoring that changes behavior

TDD is **not required** for:
- Documentation changes (`.md` files)
- Configuration changes (`.yaml`, `.json`, etc.)
- Pure refactoring (no behavior change)
- Script updates (`.sh` files)

The pre-commit hook automatically detects whether Go code is being committed and skips TDD checks for non-code changes.

## TDD Cycle

Every **code** change follows the Red-Green-Refactor cycle:

```
make tdd-red → make tdd-green → make tdd-refactor → make tdd-document
```

### Red Phase (`make tdd-red`)

Write a failing test that describes the desired behavior.

**Rules:**
- Test must FAIL initially
- Test describes WHAT, not HOW
- One behavior per test
- Use descriptive test names

**Example:**
```go
It("returns error when profile not found", func() {
    _, err := service.GetProfile(ctx, "nonexistent")
    Expect(err).To(MatchError(domain.ErrProfileNotFound))
})
```

### Green Phase (`make tdd-green`)

Write minimal code to make the test pass.

**Rules:**
- Just enough code to pass
- No extra features
- No optimization
- No refactoring

### Refactor Phase (`make tdd-refactor`)

Improve code quality while keeping tests green.

**Consider:**
- Extract methods/functions
- Remove duplication
- Improve naming
- Apply design patterns
- Ensure self-documenting code

### Document Phase (`make tdd-document`)

Finalize documentation and prepare commit.

**Checklist:**
- Go doc comments on all exports
- Update relevant docs if behavior changed
- Run `make check-compliance`
- Commit with `make ai-commit FILE=/tmp/commit.txt`

## E2E Test Requirements

### New Intents

Every new intent requires an E2E test file: `*_e2e_test.go`

**Required Sections:**

```go
var _ = Describe("FeatureName E2E", func() {
    Describe("Happy Paths", func() {
        // Test all successful state transitions
        // Cover: initial → intermediate → final states
    })

    Describe("Sad Paths", func() {
        // Test error conditions and recovery
        // Cover: validation errors, service failures, edge cases
    })
})
```

### Happy Paths

Test all successful state transitions:
- Initial state setup
- User actions (key presses, form submissions)
- State transitions
- Final state verification

**Calculate coverage:** Every state in the state machine needs a happy path.

### Sad Paths

Test error conditions and recovery:
- Validation failures
- Service errors
- Edge cases (empty data, boundary conditions)
- Error modal display
- Recovery to valid state

**Calculate minimum:** Count failure points (validation rules, service calls, edge cases).

### Bug Fixes

Every bug fix requires a regression test:

```go
Describe("Bug Regressions", func() {
    It("BUG-123: prevents duplicate form submission", func() {
        // Test that verifies the bug is fixed
    })
})
```

## Coverage Requirements

### Changed Code Coverage

All changed code must have >= 95% coverage.

The pre-commit hook calculates coverage per-package for staged changes.

### How to Check

```bash
go test -cover ./internal/cli/intents/...
```

## Test Organization

### File Naming

| Type | Pattern | Example |
|------|---------|---------|
| Unit tests | `*_test.go` | `profile_service_test.go` |
| E2E tests | `*_e2e_test.go` | `profile_intent_e2e_test.go` |
| Integration | `*_integration_test.go` | `db_integration_test.go` |

### Test Structure (Ginkgo)

```go
var _ = Describe("ComponentName", func() {
    var (
        // Shared test fixtures
    )

    BeforeEach(func() {
        // Setup
    })

    AfterEach(func() {
        // Cleanup
    })

    Describe("MethodName", func() {
        Context("when condition", func() {
            It("does expected behavior", func() {
                // Test
            })
        })
    })
})
```

## Enforcement

### Pre-commit Checks

The pre-commit hook enforces:
1. No skipped tests (`XIt`, `Skip`)
2. No pending tests (`PIt`)
3. Tests pass for changed packages
4. Coverage >= 95% on changed code

### Blocking Violations

These are BLOCKING (commit rejected):
- Skipped or pending tests
- Missing E2E tests for new intents
- Missing regression tests for bug fixes
- Coverage below 95%

## Quick Reference

```bash
# Start TDD cycle
make tdd-red

# After writing failing test
make tdd-green

# After test passes
make tdd-refactor

# After refactoring
make tdd-document

# Then commit
make ai-commit FILE=/tmp/commit.txt
```
