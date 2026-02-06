# BDD Testing with Godog

KaRiya uses [Godog](https://github.com/cucumber/godog) (Go's official Cucumber BDD framework) for writing Gherkin-style feature specifications that test TUI workflows.

## Running BDD Tests

```bash
# Run all BDD tests
make bdd

# Run only @smoke tagged scenarios (quick sanity check)
make bdd-smoke

# Run @wip (work in progress) tagged scenarios
make bdd-wip

# Run a specific feature/scenario
make bdd-feature FEATURE="Complete_full_onboarding"
```

## Directory Structure

```
features/
├── support/
│   ├── doc.go              # Package documentation
│   ├── env.go              # Test environment bridge
│   └── hooks.go            # BeforeScenario/AfterScenario hooks
├── steps/
│   ├── doc.go              # Package documentation
│   ├── common_steps.go     # Shared assertion/navigation steps
│   ├── browse_steps.go     # Browse timeline workflow steps
│   ├── bursts_steps.go     # Burst management workflow steps
│   ├── capture_steps.go    # Capture event workflow steps
│   ├── cli_steps.go        # CLI commands workflow steps
│   ├── configure_steps.go  # Configure system workflow steps
│   ├── cv_steps.go         # CV generation workflow steps
│   ├── facts_steps.go      # Fact management workflow steps
│   ├── navigation_steps.go # Navigation workflow steps
│   ├── onboarding_steps.go # Onboarding workflow steps
│   └── skills_steps.go     # Skills management workflow steps
├── browse_timeline.feature     # Timeline browsing scenarios
├── burst_management.feature    # Burst management scenarios
├── capture_event.feature       # Event capture scenarios
├── chained_workflows.feature   # Cross-workflow scenarios
├── cli_commands.feature        # CLI command scenarios
├── configure_system.feature    # System configuration scenarios
├── fact_management.feature     # Fact management scenarios
├── generate_cv.feature         # CV generation scenarios
├── navigation.feature          # Application navigation scenarios
├── onboarding.feature          # Onboarding wizard scenarios
├── skills_management.feature   # Skills management scenarios
├── godog_test.go               # Test runner
└── README.md                   # This file
```

## Writing Feature Files

Feature files use Gherkin syntax:

```gherkin
@onboarding
Feature: User Onboarding
  As a new user
  I want to complete the onboarding wizard
  So that I can set up my profile

  Background:
    Given I start the onboarding wizard

  @smoke
  Scenario: Complete Step 1 with name
    When I enter "Test User" as my name
    And I press enter
    Then I should see "Step 2 of 3"
```

## Tag Conventions

### Workflow Tags (for VHS tape generation)

| Tag | Purpose | Description |
|-----|---------|-------------|
| `@happy` | Happy path | Complete successful workflow from start to finish |
| `@sad` | Sad path | Error cases, validation failures, recovery flows |

**Happy/Sad path scenarios are designed to:**
- Capture the complete user journey through a feature
- Serve as source material for VHS demo tape generation
- Drive new functionality development
- Document expected application behavior

### Execution Tags

| Tag | Purpose | Usage |
|-----|---------|-------|
| `@smoke` | Quick sanity tests | `make bdd-smoke` |
| `@wip` | Work in progress | `make bdd-wip` |

### Feature Tags

| Tag | Purpose |
|-----|---------|
| `@onboarding` | Onboarding workflow tests |
| `@capture` | Capture workflow tests |
| `@browse` | Browse timeline tests |
| `@cv` | CV generation tests |

## Available Step Patterns

### Common Steps (common_steps.go)

| Step Pattern | Description |
|--------------|-------------|
| `I should see "{text}"` | Assert view contains text |
| `I should not see "{text}"` | Assert view does not contain text |
| `I should see one of:` | Assert view contains at least one (table) |

### Onboarding Steps (onboarding_steps.go)

| Step Pattern | Description |
|--------------|-------------|
| `I start the onboarding wizard` | Initialize onboarding environment |
| `I am on step {n} of {total}` | Assert current step |
| `I enter "{text}" as my name` | Type name |
| `I enter "{text}" as my email` | Type email |
| `I enter "{text}" as my location` | Type location |
| `I press enter` | Press Enter key |
| `I press tab` | Press Tab key |
| `I skip the optional fields` | Tab through optional fields and submit |
| `the onboarding wizard should be complete` | Assert wizard completed |
| `my profile should have name "{name}"` | Assert profile name |
| `my profile should have email "{email}"` | Assert profile email |
| `my profile should have location "{loc}"` | Assert profile location |

## Adding New Steps

1. Create a new step definition file in `features/steps/`:

```go
package steps

import (
    "context"
    "github.com/baphled/kariya/features/support"
    "github.com/cucumber/godog"
    . "github.com/onsi/gomega"
)

func RegisterMyWorkflowSteps(sc *godog.ScenarioContext) {
    sc.Step(`^I do something with "([^"]*)"$`, iDoSomething)
}

func iDoSomething(ctx context.Context, value string) error {
    env := support.GetOnboardingEnv(ctx)
    if env == nil {
        return godog.ErrPending
    }
    // Perform action and assert
    Expect(env.View()).To(ContainSubstring(value))
    return nil
}
```

2. Register in `features/godog_test.go`:

```go
func InitializeScenario(sc *godog.ScenarioContext) {
    support.RegisterHooks(sc)
    steps.RegisterCommonSteps(sc)
    steps.RegisterOnboardingSteps(sc)
    steps.RegisterMyWorkflowSteps(sc)  // Add new registration
}
```

## Test Architecture

BDD tests have replaced the previous Ginkgo-based E2E tests:

- **BDD tests**: High-level, business-focused scenarios readable by non-technical stakeholders
- **Unit tests**: Technical unit tests for individual components remain in Ginkgo

BDD tests use the test infrastructure from `internal/testutil/e2e/` which provides:
- `TestEnv` for full application testing
- Navigation helpers (`SelectIntentByName`, `Confirm`, `Cancel`, etc.)
- Data management (`AddEvent`, `AddSkill`, `GetEvents`, etc.)
- Assertion helpers (`AssertViewContains`, `AssertEventCount`, etc.)

## Best Practices

1. **Keep steps high-level** - Abstract implementation details
2. **Use descriptive scenario names** - They appear in test output
3. **Tag appropriately** - Use `@smoke` for critical paths, `@wip` for development
4. **Reuse common steps** - Add to `common_steps.go` when patterns repeat
5. **Use Gomega matchers** - Consistent with existing test style
