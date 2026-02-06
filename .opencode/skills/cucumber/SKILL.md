---
name: cucumber
description: BDD with Cucumber/Gherkin syntax for writing executable specifications and acceptance tests using smallest-change development
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Guide BDD-style development: start with a scenario, implement with the smallest possible changes, and iterate until the scenario passes.

## When to use me

Use this skill when:
- Starting a new feature (scenario first!)
- Writing acceptance criteria as executable specs
- Creating E2E test scenarios
- Implementing features incrementally
- Bridging communication between technical and non-technical stakeholders

---

## BDD Smallest-Change Workflow

### The Core Loop

```
1. Write Scenario (RED)
   └─> Scenario fails (not implemented)
   
2. Smallest Change (GREEN)
   └─> Make ONE small change to pass ONE step
   └─> Run scenario
   └─> Repeat until scenario passes
   
3. Refactor (CLEAN)
   └─> Clean code, scenario still passes
   
4. Next Scenario
   └─> Repeat
```

### Step-by-Step Process

#### Step 1: Write the Scenario First

```gherkin
# features/capture_event.feature
Feature: Capture Career Event
  As a user
  I want to capture career events
  So that I can build my timeline

  Scenario: Quick capture an event
    Given I am on the main menu
    When I select "Capture Event"
    And I choose "Quick Capture"
    And I enter "Led project kickoff meeting"
    And I submit the form
    Then I should see "Event captured successfully"
    And the event should appear in my timeline
```

#### Step 2: Run and See It Fail

```bash
go test ./features/... -v
# All steps undefined or failing - this is expected!
```

#### Step 3: Implement ONE Step Definition

```go
// Start with just the first step
func (c *captureContext) iAmOnTheMainMenu() error {
    c.app = app.New()
    c.app.Init()
    return nil
}

// Register it
ctx.Step(`^I am on the main menu$`, c.iAmOnTheMainMenu)
```

#### Step 4: Run Again - First Step Passes

```bash
go test ./features/... -v
# Step 1: PASS
# Step 2: undefined
# ...
```

#### Step 5: Implement Next Step (Smallest Change)

```go
func (c *captureContext) iSelect(option string) error {
    // Find and select the menu item
    c.app.SelectByName(option)
    return nil
}
```

#### Step 6: Repeat Until Scenario Passes

Each iteration:
1. Run tests
2. See which step fails
3. Make the SMALLEST change to pass that step
4. Run tests again
5. Refactor if needed (keeping tests green)

### What "Smallest Change" Means

| Situation | Smallest Change | NOT Smallest Change |
|-----------|-----------------|---------------------|
| Need a function | Create empty function returning nil | Implement full logic |
| Need a struct | Create struct with needed field | Add all possible fields |
| Need validation | Add one validation rule | Add all validations |
| Need UI element | Add element with minimal styling | Full styled component |

### Example: Incremental Implementation

**Scenario Step:** `When I enter "Led project kickoff meeting"`

**Change 1:** Create empty input handler
```go
func (s *Screen) HandleInput(text string) {
    // Empty - just make it compile
}
```

**Change 2:** Store the input
```go
func (s *Screen) HandleInput(text string) {
    s.inputText = text
}
```

**Change 3:** Update UI to show input (if step checks display)
```go
func (s *Screen) HandleInput(text string) {
    s.inputText = text
    s.Refresh()
}
```

Each change is a separate commit or at least a separate test run.

## Gherkin Syntax

### Feature Files

```gherkin
Feature: Timeline Event Filtering
  As a user
  I want to filter timeline events by date
  So that I can focus on a specific time period

  Background:
    Given I have the following events:
      | title          | date       |
      | Started job    | 2024-01-15 |
      | Got promotion  | 2024-06-01 |
      | Left company   | 2024-12-31 |

  Scenario: Filter events by date range
    When I filter events from "2024-01-01" to "2024-06-30"
    Then I should see 2 events
    And I should see "Started job"
    And I should see "Got promotion"
    But I should not see "Left company"

  Scenario: Clear filter shows all events
    Given I have filtered events from "2024-01-01" to "2024-06-30"
    When I clear the filter
    Then I should see 3 events
```

### Keywords

| Keyword | Purpose | Example |
|---------|---------|---------|
| `Feature:` | Describes the feature | `Feature: User Login` |
| `Background:` | Steps run before each scenario | Common setup |
| `Scenario:` | A single test case | `Scenario: Valid login` |
| `Scenario Outline:` | Parameterized scenario | With Examples table |
| `Given` | Precondition/context | `Given I am logged in` |
| `When` | Action/trigger | `When I click submit` |
| `Then` | Expected outcome | `Then I see success` |
| `And` | Additional step | `And the form is cleared` |
| `But` | Negative assertion | `But I don't see errors` |

### Scenario Outline (Parameterized)

```gherkin
Scenario Outline: Validate event title length
  Given I am creating an event
  When I enter a title with <length> characters
  Then I should see <result>

  Examples:
    | length | result           |
    | 0      | validation error |
    | 5      | success          |
    | 100    | success          |
    | 501    | validation error |
```

### Data Tables

```gherkin
Scenario: Import multiple events
  When I import the following events:
    | title        | date       | category   |
    | Event One    | 2024-01-01 | work       |
    | Event Two    | 2024-02-01 | education  |
    | Event Three  | 2024-03-01 | personal   |
  Then I should have 3 events in my timeline
```

### Doc Strings

```gherkin
Scenario: Create event with long description
  When I create an event with description:
    """
    This is a multi-line description
    that spans several lines and contains
    detailed information about the event.
    """
  Then the event should be saved successfully
```

## Writing Good Scenarios

### Focus on Behavior, Not Implementation

```gherkin
# BAD - Implementation details
Scenario: Save event to database
  Given a database connection is established
  When I insert a record into the events table
  Then the record count should increase by 1

# GOOD - Business behavior
Scenario: Create a new career event
  Given I am on the timeline
  When I create an event titled "Started new role"
  Then I should see "Started new role" in my timeline
```

### Use Domain Language

```gherkin
# BAD - Technical language
Scenario: POST request creates resource
  When I send POST to /api/events with JSON payload
  Then response status should be 201

# GOOD - Domain language
Scenario: Record a career milestone
  When I record a promotion to "Senior Engineer"
  Then my career timeline shows the promotion
```

### One Behavior Per Scenario

```gherkin
# BAD - Multiple behaviors
Scenario: Event management
  When I create an event
  Then it appears in timeline
  When I edit the event
  Then changes are saved
  When I delete the event
  Then it's removed

# GOOD - Single behavior each
Scenario: Create event
  When I create an event titled "New role"
  Then I see "New role" in my timeline

Scenario: Edit event
  Given an event titled "New role" exists
  When I change the title to "Senior role"
  Then I see "Senior role" in my timeline

Scenario: Delete event
  Given an event titled "New role" exists
  When I delete the event
  Then I no longer see "New role" in my timeline
```

### Declarative Over Imperative

```gherkin
# BAD - Imperative (how)
Scenario: Login
  Given I am on the login page
  When I enter "user@example.com" in the email field
  And I enter "password123" in the password field
  And I click the login button
  Then I am redirected to the dashboard

# GOOD - Declarative (what)
Scenario: Successful login
  Given I am a registered user
  When I login with valid credentials
  Then I see my dashboard
```

## Godog (Go Cucumber)

### Step Definitions

```go
package main

import (
    "github.com/cucumber/godog"
)

type timelineContext struct {
    events []Event
    filter *DateFilter
}

func (tc *timelineContext) iHaveTheFollowingEvents(table *godog.Table) error {
    for _, row := range table.Rows[1:] { // Skip header
        tc.events = append(tc.events, Event{
            Title: row.Cells[0].Value,
            Date:  parseDate(row.Cells[1].Value),
        })
    }
    return nil
}

func (tc *timelineContext) iFilterEventsFromTo(from, to string) error {
    tc.filter = &DateFilter{
        From: parseDate(from),
        To:   parseDate(to),
    }
    return nil
}

func (tc *timelineContext) iShouldSeeNEvents(count int) error {
    filtered := tc.applyFilter()
    if len(filtered) != count {
        return fmt.Errorf("expected %d events, got %d", count, len(filtered))
    }
    return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
    tc := &timelineContext{}
    
    ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
        tc.events = nil
        tc.filter = nil
        return ctx, nil
    })
    
    ctx.Step(`^I have the following events:$`, tc.iHaveTheFollowingEvents)
    ctx.Step(`^I filter events from "([^"]*)" to "([^"]*)"$`, tc.iFilterEventsFromTo)
    ctx.Step(`^I should see (\d+) events$`, tc.iShouldSeeNEvents)
    ctx.Step(`^I should see "([^"]*)"$`, tc.iShouldSeeEvent)
    ctx.Step(`^I should not see "([^"]*)"$`, tc.iShouldNotSeeEvent)
}
```

### Running Tests

```bash
# Run all features
godog

# Run specific feature
godog features/timeline.feature

# Run with tags
godog --tags=@wip

# Generate step snippets
godog --format=snippets
```

## Tags

```gherkin
@timeline @filtering
Feature: Timeline Filtering

  @smoke
  Scenario: Basic filter works
    ...

  @slow @integration
  Scenario: Filter with large dataset
    ...

  @wip
  Scenario: New filter feature
    ...
```

```bash
# Run only smoke tests
godog --tags=@smoke

# Exclude slow tests
godog --tags="not @slow"

# Combine tags
godog --tags="@timeline and not @wip"
```

## KaRiya E2E Pattern

```gherkin
Feature: Browse Timeline Intent
  As a user
  I want to browse my career timeline
  So that I can review my career history

  Background:
    Given I have career events in my timeline

  @navigation
  Scenario: Navigate timeline with keyboard
    Given I am viewing the timeline
    When I press "j"
    Then the next event is selected
    When I press "k"
    Then the previous event is selected

  @modal
  Scenario: View event details
    Given I am viewing the timeline
    When I press "enter" on an event
    Then I see the event detail modal
    When I press "escape"
    Then the modal closes

  @filtering
  Scenario: Filter events by category
    Given I am viewing the timeline
    When I open the filter modal
    And I select category "work"
    Then I only see work events
```

## Best Practices

### File Organization

```
features/
├── timeline/
│   ├── browse.feature
│   ├── filter.feature
│   └── search.feature
├── events/
│   ├── create.feature
│   ├── edit.feature
│   └── delete.feature
└── support/
    └── steps/
        ├── timeline_steps.go
        └── event_steps.go
```

### Step Reusability

```go
// Reusable steps across features
ctx.Step(`^I am logged in$`, tc.iAmLoggedIn)
ctx.Step(`^I am on the timeline$`, tc.iAmOnTheTimeline)
ctx.Step(`^I see "([^"]*)"$`, tc.iSee)
ctx.Step(`^I don't see "([^"]*)"$`, tc.iDontSee)
```

### Avoid These

| Anti-Pattern | Problem | Better Approach |
|--------------|---------|-----------------|
| Too many steps | Hard to read | Combine into meaningful actions |
| Technical language | Not stakeholder-readable | Use domain language |
| Coupled scenarios | Order-dependent | Independent scenarios |
| UI details | Brittle tests | Declarative behavior |
| No Background | Repetitive setup | Use Background for common Given |

---

## BDD with Ginkgo (KaRiya Pattern)

Since KaRiya uses Ginkgo, BDD scenarios translate to Ginkgo specs:

### Scenario to Ginkgo Translation

```gherkin
Scenario: Quick capture an event
  Given I am on the main menu
  When I select "Capture Event"
  And I choose "Quick Capture"
  And I enter "Led project kickoff meeting"
  And I submit the form
  Then I should see "Event captured successfully"
```

Becomes:

```go
var _ = Describe("Capture Event", func() {
    var env *e2e.TestEnv
    
    BeforeEach(func() {
        env = e2e.GetSharedEnv()
        env.Reset()
    })
    
    Describe("Quick Capture", func() {
        It("captures an event successfully", func() {
            // Given I am on the main menu (implicit - env starts at menu)
            
            // When I select "Capture Event"
            env.SelectIntentByName("Capture Event")
            
            // And I choose "Quick Capture"
            env.NavigateDown()  // Go to Quick Capture
            env.Confirm()
            
            // And I enter "Led project kickoff meeting"
            env.TypeText("Led project kickoff meeting")
            
            // And I submit the form
            env.Confirm()
            
            // Then I should see "Event captured successfully"
            env.AssertViewContains("Event captured")
        })
    })
})
```

### Smallest Change with Ginkgo

```go
// Step 1: Write failing test
It("captures an event successfully", func() {
    env.SelectIntentByName("Capture Event")
    env.Confirm()
    env.TypeText("Test event")
    env.Confirm()
    env.AssertViewContains("Event captured")
})

// Step 2: Run - fails because intent doesn't exist
// Smallest change: Create empty intent

// Step 3: Run - fails because form doesn't exist
// Smallest change: Add form screen

// Step 4: Run - fails because submission doesn't work
// Smallest change: Handle form submission

// Step 5: Run - fails because success message missing
// Smallest change: Show success message

// Step 6: PASS - Refactor if needed
```

---

## The BDD Mindset

### Outside-In Development

```
Scenario (Acceptance Test)
    │
    ▼ drives
Integration Test
    │
    ▼ drives
Unit Test
    │
    ▼ drives
Implementation
```

### Questions Before Coding

1. **What behavior am I implementing?** (Write scenario)
2. **What's the smallest step?** (One Given/When/Then)
3. **What's the minimal code?** (Just enough to pass)
4. **Is it clean?** (Refactor)

### Commit Cadence

```
feat(capture): add scenario for quick capture [RED]
feat(capture): create capture intent stub [GREEN step 1]
feat(capture): add form screen [GREEN step 2]
feat(capture): handle form submission [GREEN step 3]
feat(capture): show success message [GREEN all steps]
refactor(capture): extract form validation [REFACTOR]
```

---

## Common Patterns

### Background for Shared Setup

```gherkin
Background:
  Given I have the following events:
    | title         | date       |
    | Started job   | 2024-01-15 |
    | Got promotion | 2024-06-01 |
```

### Scenario Outline for Variations

```gherkin
Scenario Outline: Validate event text
  When I enter event text "<text>"
  Then I should see "<result>"

  Examples:
    | text                    | result           |
    |                         | Text is required |
    | Valid event description | Event captured   |
```

### Tags for Organization

```gherkin
@wip
Scenario: Feature in progress
  ...

@slow
Scenario: Long-running test
  ...
```

```bash
# Run only WIP
go test ./... -tags=wip

# Skip slow tests
go test ./... -tags="!slow"
```

---

## Anti-Patterns

### DON'T: Implement Everything at Once

```go
// WRONG - Implementing full feature before running tests
func NewCaptureIntent() *CaptureIntent {
    return &CaptureIntent{
        formScreen: NewFormScreen(),
        reviewScreen: NewReviewScreen(),
        successScreen: NewSuccessScreen(),
        validation: NewValidator(),
        // ... 200 more lines
    }
}

// RIGHT - Smallest change to pass current step
func NewCaptureIntent() *CaptureIntent {
    return &CaptureIntent{}
}
```

### DON'T: Skip the RED Phase

```go
// WRONG - Writing implementation without failing test
func (s *Service) SaveEvent(e *Event) error {
    return s.repo.Save(e)
}

// RIGHT - Write scenario/test first, see it fail, then implement
```

### DON'T: Big Refactors During GREEN

```go
// WRONG - Major refactor while trying to pass test
// (refactors should happen AFTER green, not during)

// RIGHT - Minimal change to pass, refactor later
```

---

## Related skills

- `tdd-workflow` - TDD complements BDD (same red-green-refactor)
- `ginkgo-gomega` - KaRiya's testing framework
- `e2e-testing` - E2E patterns for scenarios
- `clean-code` - Refactor phase
- `create-task` - Acceptance criteria as Gherkin
