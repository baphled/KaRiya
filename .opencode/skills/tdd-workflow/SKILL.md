---
name: tdd-workflow
description: Follow the TDD Red-Green-Refactor cycle for KaRiya development with proper phase tracking
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Guide you through the Test-Driven Development (TDD) cycle: Red -> Green -> Refactor -> Document.

## When to use me

Use this skill whenever implementing new features or fixing bugs. TDD is **mandatory** in this project.

## TDD Phases

### 1. Red Phase - Write Failing Test

```bash
make tdd-red
```

**Rules:**
- Test must FAIL initially
- Test describes WHAT, not HOW
- One behavior per test
- Use Ginkgo BDD style (`Describe`, `Context`, `It`)

**Example test structure:**
```go
var _ = Describe("MyFeature", func() {
    Context("when condition is met", func() {
        It("should do expected behavior", func() {
            // Arrange
            // Act
            // Assert
            Expect(result).To(Equal(expected))
        })
    })
})
```

### 2. Green Phase - Make Test Pass

```bash
make tdd-green
```

**Rules:**
- Write MINIMAL code to pass
- No extra features
- Don't optimize yet
- Just enough to make the test green

### 3. Refactor Phase - Improve Code

```bash
make tdd-refactor
```

**Consider:**
- Extract methods/functions
- Remove duplication
- Improve naming
- Apply SOLID principles
- Use existing patterns from `docs/development/INTENT_PATTERNS_LIBRARY.md`

### 4. Document Phase - Finalize

```bash
make tdd-document
```

**Checklist:**
- Go doc comments on all exports (Required sections: Expected, Returns, Side effects)
- Run `make check-compliance`
- Create commit: `make ai-commit FILE=/tmp/commit.txt`

## Testing Commands

```bash
make test                          # Run all tests
make test-suite SUITE=./path/...   # Run specific suite
make individual-test TEST="name"   # Run single test
make coverage                      # Generate coverage report
```

## Coverage Requirements

| Scope | Threshold |
|-------|-----------|
| Per-package (modified) | >= 95% (BLOCKING) |
| Project average | >= 80% (warning) |

## Checklist Discipline (MANDATORY)

### Update Progress Incrementally
After EACH phase, update the checklist immediately:

```
[x] Red Phase - Write failing test
    COMPLETED: Test expects validation error for empty title

[ ] Green Phase - Make test pass
    IN PROGRESS: Implementing validation...

[ ] Refactor Phase
[ ] Document Phase
```

### Skip Reason Requirement
If skipping ANY step, document explicitly:

```
[SKIP] Refactor Phase
    SKIPPING: Refactor phase
    REASON: Implementation was already clean, single function
    IMPACT: None - code meets quality standards
```

**NEVER silently skip phases.** All skips require documented reason.

## Related skills

- `session-start` - Initialize session before TDD
- `check-compliance` - Validate after refactor phase
- `debug-test` - When tests fail unexpectedly
- `checklist-discipline` - Incremental progress tracking
