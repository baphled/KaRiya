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

Guide you through the **strict** Test-Driven Development (TDD) cycle: Red -> Green -> Refactor.

## When to use me

Use this skill whenever implementing new features or fixing bugs. TDD is **mandatory** in this project.

## STRICT RED-GREEN-REFACTOR

This is non-negotiable. Every change follows this exact cycle:

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│   RED ──────────► GREEN ──────────► REFACTOR ───┐          │
│    │                                            │          │
│    │  Write ONE      Make it       Clean up     │          │
│    │  failing test   pass ONLY     (tests pass) │          │
│    │                                            │          │
│    └────────────────────────────────────────────┘          │
│                         │                                   │
│                         ▼                                   │
│                    Next behavior                            │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## Test Execution Strategy (CRITICAL)

### During RED-GREEN-REFACTOR: Run ONLY the specific test

```bash
# Run ONLY the test you're working on
go test -v ./path/to/package -run "TestSpecificName"

# Or with Ginkgo focus
ginkgo -v --focus "specific test description" ./path/to/package
```

**Why?** Fast feedback. You need to know within seconds if your change worked.

### After GREEN: Run package tests

```bash
# Once your specific test passes, run the package
go test -v ./path/to/package/...
```

### After REFACTOR: Run related tests

```bash
# Run tests for all packages you touched
go test -v ./internal/cli/intents/myfeature/... ./internal/cli/screens/myfeature/...
```

### Before COMMIT: Run full suite

```bash
# Only now run everything
make test
make check-compliance
```

### Test Execution Pyramid

```
                    ┌─────────────┐
                    │  Full Suite │  ← Before commit ONLY
                    │  make test  │
                    └──────┬──────┘
                           │
                    ┌──────┴──────┐
                    │   Related   │  ← After refactor
                    │  Packages   │
                    └──────┬──────┘
                           │
                    ┌──────┴──────┐
                    │   Package   │  ← After green
                    │    Tests    │
                    └──────┬──────┘
                           │
                    ┌──────┴──────┐
                    │   Single    │  ← During RED-GREEN
                    │    Test     │     (FAST!)
                    └─────────────┘
```

## TDD Phases

### 1. RED Phase - Write ONE Failing Test

```bash
# Write the test
# Run ONLY that test - it MUST fail
go test -v ./path -run "TestMyNewBehavior"
```

**Rules:**
- Write ONE test for ONE behavior
- Test MUST fail (if it passes, you wrote the wrong test)
- Test describes WHAT, not HOW
- Use Ginkgo BDD style (`Describe`, `Context`, `It`)
- **DO NOT write implementation yet**

**Verify RED:**
```
--- FAIL: TestMyNewBehavior
    Expected: <something>
    Actual:   <something else>
FAIL
```

If test passes → STOP. Your test is wrong or the feature already exists.

### 2. GREEN Phase - Make ONLY That Test Pass

```bash
# Write minimal implementation
# Run ONLY that test - it MUST pass now
go test -v ./path -run "TestMyNewBehavior"
```

**Rules:**
- Write MINIMAL code to pass THIS test
- No extra features
- No optimization
- No cleanup
- Hardcode values if needed (you'll fix in refactor)
- **UGLY IS OK** - just make it green

**Verify GREEN:**
```
--- PASS: TestMyNewBehavior
PASS
```

**Now run package tests:**
```bash
go test -v ./path/...
```

If other tests fail → Fix them WITHOUT adding features.

### 3. REFACTOR Phase - Clean Up (Tests Still Pass)

```bash
# Improve code quality
# Run package tests after each change
go test -v ./path/...
```

**Rules:**
- Tests MUST stay green throughout
- Run tests after EVERY refactor step
- If tests fail → UNDO and try smaller refactor
- Apply design patterns appropriately
- Remove duplication
- Improve naming
- Extract methods/functions

**Consider:**
- SOLID principles
- Design patterns (see `design-patterns` skill)
- Existing patterns from `docs/development/INTENT_PATTERNS_LIBRARY.md`

**After refactor complete, run related packages:**
```bash
go test -v ./internal/cli/intents/myfeature/... ./internal/cli/screens/myfeature/...
```

### 4. REPEAT or COMMIT

**If more behaviors needed:**
- Go back to RED phase
- Write next failing test

**If feature complete:**
```bash
# NOW run full suite
make test
make check-compliance
make ai-commit FILE=/tmp/commit.txt
```

## Forbidden Behaviors

### NEVER Do These

| Violation | Why It's Wrong |
|-----------|----------------|
| Write implementation before test | Not TDD |
| Write multiple tests at once | Can't isolate failures |
| Make test pass + add extra features | Gold plating |
| Skip RED (test already passes) | Feature already exists or test is wrong |
| Run full suite during RED-GREEN | Too slow, breaks flow |
| Refactor while RED | Fix one thing at a time |
| Add features during REFACTOR | Refactor ≠ new features |

### Correct Responses

| Situation | Correct Action |
|-----------|----------------|
| Test passes immediately | STOP - investigate why |
| Multiple tests fail in GREEN | Fix ONE at a time |
| Want to add "just one more thing" | STOP - new RED-GREEN cycle |
| Refactor breaks tests | UNDO immediately |
| Full suite fails after GREEN | Run specific failures, fix one by one |

## Example Workflow

```
Task: Add email validation to user registration

[x] RED: Write test for valid email acceptance
    go test -v ./internal/service -run "TestValidEmail"
    FAIL ✓ (expected - no validation exists)

[x] GREEN: Add minimal validation
    go test -v ./internal/service -run "TestValidEmail"  
    PASS ✓
    go test -v ./internal/service/...
    PASS ✓

[x] REFACTOR: Extract validation to helper
    go test -v ./internal/service/...
    PASS ✓ (after each change)

[x] RED: Write test for invalid email rejection
    go test -v ./internal/service -run "TestInvalidEmail"
    FAIL ✓

[x] GREEN: Add invalid check
    go test -v ./internal/service -run "TestInvalidEmail"
    PASS ✓
    go test -v ./internal/service/...
    PASS ✓

[x] REFACTOR: Consolidate validation logic
    go test -v ./internal/service/...
    PASS ✓

[x] COMMIT: Feature complete
    make test  # Full suite
    make check-compliance
    make ai-commit FILE=/tmp/commit.txt
```

## Checklist Discipline (MANDATORY)

### Update Progress Incrementally
After EACH phase, update the checklist immediately:

```
[x] Red Phase - Write failing test
    COMPLETED: Test expects validation error for empty title
    TEST: go test -v ./internal/service -run "TestEmptyTitle"

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

## Testing Commands Quick Reference

```bash
# Single test (during RED-GREEN)
go test -v ./path -run "TestName"
ginkgo -v --focus "description" ./path

# Package tests (after GREEN)
go test -v ./path/...

# Related packages (after REFACTOR)
go test -v ./pkg1/... ./pkg2/...

# Full suite (before COMMIT only)
make test

# With coverage
go test -cover ./path/...
```

## Coverage Requirements

| Scope | Threshold |
|-------|-----------|
| Per-package (modified) | >= 95% (BLOCKING) |
| Project average | >= 80% (warning) |

Check coverage only before commit:
```bash
go test -cover ./internal/cli/intents/myfeature/...
```

## Related skills

- `session-start` - Initialize session before TDD
- `check-compliance` - Validate after refactor phase
- `debug-test` - When tests fail unexpectedly
- `checklist-discipline` - Incremental progress tracking
- `ginkgo-gomega` - Testing framework details
- `design-patterns` - Patterns to apply during refactor
- `concurrency` - Safe concurrent code patterns
