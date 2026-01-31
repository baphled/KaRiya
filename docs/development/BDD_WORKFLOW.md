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

## Testing Tools

KaRiya uses a comprehensive testing toolkit. Use the right tool for each scenario:

### Core Framework

| Tool | Package | Purpose |
|------|---------|---------|
| **Ginkgo v2** | `github.com/onsi/ginkgo/v2` | BDD-style test structure (`Describe`, `Context`, `It`) |
| **Gomega** | `github.com/onsi/gomega` | Expressive matchers (`Expect(...).To(...)`) |

### Fixtures & Test Data

| Tool | Package | Purpose |
|------|---------|---------|
| **fixtures** | `internal/testutil/fixtures` | Factory-based test data (`fixtures.Event()`, `fixtures.Burst()`) |
| **gofakeit** | `github.com/brianvoe/gofakeit/v7` | Realistic fake data generation |
| **factory-go** | `github.com/bluele/factory-go` | Factory pattern for complex objects |

### E2E Testing

| Tool | Package | Purpose |
|------|---------|---------|
| **TestEnv** | `internal/testutil/e2e` | Full application testing environment |
| **e2e fixtures** | `internal/testutil/e2e/fixtures.go` | Sample data for E2E tests |

### Mocking

| Tool | Package | Purpose |
|------|---------|---------|
| **Memory repositories** | `internal/repository/career` | Fast in-memory implementations |
| **Mock repositories** | `internal/repository/career/mocks` | Behavior-based mocking |

### Database

| Tool | Package | Purpose |
|------|---------|---------|
| **testutil.SetupTestDB** | `internal/testutil` | SQLite test database setup |
| **SQLite** | `modernc.org/sqlite` | In-process database for integration tests |

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
- **Fixtures first** - create fixture data before writing assertions
- **AAA Pattern** - structure test as Arrange-Act-Assert (see [AAA Pattern](#aaa-pattern-arrange-act-assert---critical))
- **One assertion per test** (unit tests only - see [One Assertion Per Test](#one-assertion-per-test-unit-tests---critical))

**Example:**
```go
It("returns error when profile not found", func() {
    // Arrange
    // (no fixture needed - testing not-found case)

    // Act
    _, err := service.GetProfile(ctx, "nonexistent")

    // Assert
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

## Test Best Practices

### Fixtures First

**Always create fixture data before writing assertions.** This ensures tests are readable and maintainable.

**Why fixtures first?**
- Separates test data setup from test logic
- Makes tests self-documenting
- Enables reuse of test data patterns
- Simplifies debugging when tests fail

**Create fixtures if they don't exist:**
```go
// If you need a new fixture pattern, add it to internal/testutil/fixtures/
// Example: Adding a new helper in event_factory.go
func EventWithTags(id string, tags []string) *career.CareerEvent {
    event := Event(id)
    event.Tags = tags
    return event
}
```

### Use Available Fixture Helpers

```go
// Quick helpers for minimal valid objects
event := fixtures.Event("my-id")
burst := fixtures.Burst("burst-id", "evt-1", "evt-2")
fact := fixtures.Fact("fact-id", "evt-1")

// Helpers with custom fields
event := fixtures.EventWith("evt-1", "Led team", "TechCorp", "Platform")

// Factory pattern for complex customization
event := fixtures.EventFactory.MustCreateWithOption(map[string]interface{}{
    "Company": "TechCorp",
    "Project": "Platform",
}).(*career.CareerEvent)

// Batch creation
events := fixtures.Events(5)  // Creates 5 events with sequential IDs

// Reproducible data
fixtures.SetSeed(12345)
event := fixtures.EventFactory.MustCreate().(*career.CareerEvent)
```

### Keep Tests Simple

**One behavior per test.** Complex tests are hard to debug and maintain.

**Do:**
```go
It("returns error when event ID is empty", func() {
    // Arrange
    // (no fixture needed)

    // Act
    _, err := service.GetEvent(ctx, "")

    // Assert
    Expect(err).To(HaveOccurred())
})
```

**Don't:**
```go
It("validates all event fields", func() {
    // Testing multiple behaviors in one test - hard to debug
    _, err := service.GetEvent(ctx, "")
    Expect(err).To(HaveOccurred())
    
    _, err = service.GetEvent(ctx, "nonexistent")
    Expect(err).To(MatchError(ErrNotFound))
    
    event, err := service.GetEvent(ctx, "valid-id")
    Expect(err).NotTo(HaveOccurred())
    Expect(event.ID).To(Equal("valid-id"))
})
```

### AAA Pattern (Arrange-Act-Assert) - CRITICAL

> **NON-NEGOTIABLE**: Every test MUST follow the AAA pattern. Tests that do not follow this structure will be rejected in code review.

The AAA pattern is **mandatory** for all tests in this codebase. It ensures tests are:
- **Readable**: Anyone can understand what's being tested at a glance
- **Maintainable**: Clear separation makes updates straightforward
- **Debuggable**: When tests fail, the structure reveals exactly what went wrong

```go
It("extracts facts from event with leadership keywords", func() {
    // Arrange - Set up test data using fixtures
    event := fixtures.EventWith(
        uuid.New().String(),
        "Led development team through critical project",
        "TechCorp",
        "Platform",
    )

    // Act - Execute the behavior under test
    facts, err := service.ExtractFactsFromEvent(ctx, event)

    // Assert - Verify the expected outcome
    Expect(err).NotTo(HaveOccurred())
})
```

**Mandatory Rules:**

| Phase | Rule | Violation |
|-------|------|-----------|
| **Arrange** | Use fixtures, not inline object creation | Inline `&career.CareerEvent{...}` in test body |
| **Arrange** | Set up ALL test data before acting | Creating data after the Act phase |
| **Act** | Single function call or operation | Multiple operations mixed together |
| **Act** | No assertions in the Act phase | `Expect()` calls before Act completes |
| **Assert** | **One assertion per test** (unit tests) | Multiple `Expect()` calls testing different behaviors |
| **Assert** | Assert AFTER acting, never before | Assertions scattered throughout the test |

**Section Comments Required:**

For clarity, include `// Arrange`, `// Act`, and `// Assert` comments in every test. This makes the structure explicit and helps during code review.

```go
It("returns error when event is nil", func() {
    // Arrange
    // (no fixture needed for nil input test)

    // Act
    _, err := service.ProcessEvent(ctx, nil)

    // Assert
    Expect(err).To(HaveOccurred())
})
```

**When Arrange is Empty:**

Some tests don't need fixture setup (e.g., nil input tests). Still include the comment to show the pattern was considered:

```go
// Arrange
// (no setup needed - testing nil input)
```

### One Assertion Per Test (Unit Tests) - CRITICAL

> **NON-NEGOTIABLE**: Unit tests MUST have exactly ONE assertion. If you need to verify multiple things, write multiple tests. E2E tests are exempt from this rule.

**Why one assertion?**
- When a test fails, you know exactly what broke
- Tests become documentation - each test name describes one behavior
- Easier to maintain and refactor
- Forces you to think about what you're actually testing

**Do - Separate tests for each assertion:**
```go
It("returns error when event is nil", func() {
    // Arrange
    // (no fixture needed)

    // Act
    _, err := service.ProcessEvent(ctx, nil)

    // Assert
    Expect(err).To(HaveOccurred())
})

It("returns error message indicating nil event", func() {
    // Arrange
    // (no fixture needed)

    // Act
    _, err := service.ProcessEvent(ctx, nil)

    // Assert
    Expect(err.Error()).To(ContainSubstring("event cannot be nil"))
})

It("extracts facts from event", func() {
    // Arrange
    event := fixtures.EventWith(uuid.New().String(), "Led team", "TechCorp", "Platform")

    // Act
    facts, err := service.ExtractFactsFromEvent(ctx, event)

    // Assert - only check error here
    Expect(err).NotTo(HaveOccurred())
})

It("extracts non-empty facts from event with content", func() {
    // Arrange
    event := fixtures.EventWith(uuid.New().String(), "Led team", "TechCorp", "Platform")

    // Act
    facts, _ := service.ExtractFactsFromEvent(ctx, event)

    // Assert - only check facts here
    Expect(facts).NotTo(BeEmpty())
})
```

**Don't - Multiple assertions in one test:**
```go
It("returns error when event is nil", func() {
    // Arrange
    // (no fixture needed)

    // Act
    _, err := service.ProcessEvent(ctx, nil)

    // Assert - VIOLATION: Multiple assertions
    Expect(err).To(HaveOccurred())
    Expect(err.Error()).To(ContainSubstring("event cannot be nil"))
})
```

**Exception - E2E Tests:**

E2E tests (`*_e2e_test.go`) may have multiple assertions because they test complete workflows where verifying intermediate states is necessary for understanding the full user journey:

```go
// E2E tests can have multiple assertions
It("completes the capture event workflow", func() {
    // Arrange
    env.PopulateTestData(5, 0, 0)

    // Act
    env.SelectIntentByName("capture_event")
    env.TypeText("Led team through migration")
    env.Confirm()

    // Assert - Multiple assertions OK in E2E
    env.AssertEventCount(6)
    env.AssertViewContains("Success")
    env.AssertViewNotContains("Error")
})
```

### Test Structure Template

```go
var _ = Describe("ComponentName", func() {
    var (
        // Shared dependencies
        repo    careerrepo.Repository
        service *Service
        ctx     context.Context
    )

    BeforeEach(func() {
        // Common setup - use memory repositories for speed
        repo = careerrepo.NewMemoryRepository()
        service = NewService(repo)
        ctx = context.Background()
    })

    Describe("MethodName", func() {
        Context("when condition is met", func() {
            It("does expected behavior", func() {
                // Arrange
                event := fixtures.Event("test-id")

                // Act
                result, err := service.DoSomething(ctx, event)

                // Assert
                Expect(err).NotTo(HaveOccurred())
            })

            It("returns correct result", func() {
                // Arrange
                event := fixtures.Event("test-id")

                // Act
                result, _ := service.DoSomething(ctx, event)

                // Assert
                Expect(result).To(Equal(expected))
            })
        })

        Context("when condition is not met", func() {
            It("returns appropriate error", func() {
                // Arrange
                // (no fixture needed for nil input)

                // Act
                _, err := service.DoSomething(ctx, nil)

                // Assert
                Expect(err).To(HaveOccurred())
            })
        })
    })
})
```

### E2E Test Setup

For integration tests, use the E2E test environment:

```go
var _ = Describe("Feature E2E", func() {
    var env *e2e.TestEnv

    BeforeEach(func() {
        // Arrange - Set up complete test environment
        env = e2e.Setup(GinkgoT())
    })

    AfterEach(func() {
        env.Cleanup()
    })

    It("completes the workflow successfully", func() {
        // Arrange - Populate with test data
        env.PopulateTestData(10, 3, 5)  // events, bursts, facts

        // Act - Simulate user interaction
        env.SelectIntentByName("capture_event")
        env.TypeText("Test event text")
        env.Confirm()

        // Assert - Multiple assertions OK in E2E tests
        env.AssertEventCount(11)
        env.AssertViewContains("Success")
    })
})
```

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

### Code Review Requirements

These are enforced during code review:
1. **AAA Pattern** - Every test must follow Arrange-Act-Assert structure
2. **Fixtures First** - Test data must use fixtures, not inline creation
3. **One Assertion Per Test** - Unit tests must have exactly one assertion (E2E exempt)
4. **Section Comments** - `// Arrange`, `// Act`, `// Assert` comments required

### Blocking Violations

These are BLOCKING (commit rejected):
- Skipped or pending tests
- Missing E2E tests for new intents
- Missing regression tests for bug fixes
- Coverage below 95%
- Tests not following AAA pattern (code review)
- Inline object creation instead of fixtures (code review)
- Multiple assertions in unit tests (code review)

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
