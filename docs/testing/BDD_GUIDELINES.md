# BDD Testing Guidelines for KaRiya

**Last Updated**: February 12, 2026  
**Status**: Active Guidelines  
**Applies To**: All Cucumber/Gherkin feature files  

---

## Table of Contents

1. [Philosophy](#philosophy)
2. [Core Principle](#core-principle)
3. [What BDD Tests](#what-bdd-tests)
4. [What BDD Does NOT Test](#what-bdd-does-not-test)
5. [Test Pyramid](#test-pyramid)
6. [Anti-Patterns to Avoid](#anti-patterns-to-avoid)
7. [Writing Good Scenarios](#writing-good-scenarios)
8. [Examples](#examples)
9. [Decision Framework](#decision-framework)

---

## Philosophy

BDD exists to bridge the gap between business people and technical people by:

1. **Describing business outcomes**, not implementation details
2. **Using concrete examples** that everyone understands
3. **Creating living documentation** automatically verified by tests
4. **Collaborating early** (Three Amigos: PO, Tester, Developer)

**Source**: Official Cucumber.io best practices

---

## Core Principle

### ✅ WHAT We Test in BDD

**Write about business value and user workflows**

- User capabilities and workflows
- Business rules and constraints
- Data transformations and outcomes
- Integration across features
- Error handling and edge cases
- Critical user paths

### ❌ WHAT We Do NOT Test in BDD

**Avoid implementation details and UI mechanics**

- How users interact with UI (keyboard keys, clicks, modals)
- Form field navigation (Tab, Shift+Tab, focus movement)
- Modal opening/closing cycles
- Dropdown/list scrolling
- Button styling or visibility
- Animation timing
- Screen transitions
- Element presence verification

**These belong in UNIT TESTS**, not BDD.

---

## What BDD Tests

### Business Logic Scenarios

**Example: Event Capture**
```gherkin
Scenario: A professional captures a career event
  Given I want to record a recent achievement
  When I capture the event "Built microservices architecture in Go"
  And I set company to "Acme Corp"
  And I submit the event
  Then the system should create a new career event
  And the event should contain description "Built microservices architecture in Go"
  And the event should be linked to company "Acme Corp"
```

**Why This is Good**:
- Describes WHAT happens (business outcome), not HOW user clicks
- Tests the business rule: "Can I create an event with metadata?"
- Survives UI changes (keyboard, mouse, GUI, CLI)
- Reads as business documentation
- Validates critical path

### Workflow Scenarios

**Example: Multi-Step Enrichment**
```gherkin
Scenario: Professional reviews AI-suggested event enrichment
  Given I captured an event "Designed database schema"
  When I select the enrichment review screen
  And I accept the suggested burst "Data Architecture"
  And I confirm the category "Technical Leadership"
  Then the system should link the event to burst "Data Architecture"
  And the event metadata should include category "Technical Leadership"
```

**Why This is Good**:
- Tests end-to-end workflow
- Business value: automated event organization
- Hidden implementation (how you navigate internally)
- Multiple systems interact (event storage, burst service, categorization)

### Error Handling Scenarios

**Example: Validation Rules**
```gherkin
Scenario: System rejects invalid event metadata
  Given I am creating a new event
  When I attempt to save with empty event description
  Then the system should reject the save
  And should display validation error "Description required"

Scenario: System validates company names
  Given I am editing an event
  When I attempt to set company to "   " (whitespace)
  Then the system should reject the change
  And should preserve the original company value
```

**Why This is Good**:
- Tests business constraints
- Ensures data quality
- Covers edge cases
- Validates error messages users see

---

## What BDD Does NOT Test

### ❌ Anti-Pattern: Modal Mechanics

**WRONG** ❌
```gherkin
Scenario: User opens filter modal
  When I press "f" to open filter modal
  Then I should see the filter modal
  And the modal should contain checkboxes
  And the modal should display company options
```

**Why It's Wrong**:
- Tests UI implementation ("modal appearance"), not business outcome
- Breaks when you change modal library
- Provides no business value documentation
- Duplicates unit tests

**CORRECT** ✅
```gherkin
Scenario: Professional filters events by company
  Given I have events from multiple companies
  When I filter to show only "Acme Corp" events
  Then the timeline should display only "Acme Corp" events
  And other companies should not be visible
```

### ❌ Anti-Pattern: Keyboard Navigation

**WRONG** ❌
```gherkin
Scenario: User navigates with j key
  When I press "j" key
  Then the selection should move down

Scenario: User navigates with arrow keys
  When I press Down arrow key
  Then the cursor should move to next item
```

**Why It's Wrong**:
- Tests keyboard implementation, not business outcome
- These shortcuts are UI mechanics, not business value
- Belong in unit tests (UIKit, handlers)
- Makes scenarios brittle

**CORRECT** ✅
```gherkin
Scenario: Professional browses events in timeline
  Given I have multiple events
  When I navigate through the timeline
  Then I should see events in chronological order
```

### ❌ Anti-Pattern: Form Field Mechanics

**WRONG** ❌
```gherkin
Scenario: Form field navigation
  When I press Tab to move between fields
  And I navigate to the company field
  And I press Shift+Tab to go back
  Then the focus should be on the description field
```

**Why It's Wrong**:
- Tests form library behavior, not business logic
- Changes with UI framework updates
- No business value being tested
- Unit tests already cover form interactions

**CORRECT** ✅
```gherkin
Scenario: Professional completes event metadata
  Given I am editing an event
  When I fill in all required metadata fields
  And I submit the form
  Then the event should be updated with new metadata
```

### ❌ Anti-Pattern: Focus Management

**WRONG** ❌
```gherkin
Scenario: Modal captures focus
  When I open a modal
  Then typing should go to the modal input
  And pressing escape should close the modal
```

**Why It's Wrong**:
- Tests input implementation detail
- Terminal/GUI specific
- Not a business requirement
- Unit tests already cover this

**CORRECT** ✅
```gherkin
Scenario: Professional adds event within workflow
  Given I'm in the capture event workflow
  When I add required information and submit
  Then the event should be created
  And the workflow should confirm success
```

---

## Test Pyramid

KaRiya follows the industry-standard test pyramid:

```
       ┌─────────────────────┐
       │    BDD / E2E        │   20% of tests
       │  Business Workflows │   - Critical user paths
       │  & Integration      │   - Multi-feature scenarios
       │                     │   - User-visible value
       ├─────────────────────┤
       │   Integration Tests │   40% of tests
       │  Feature Interactions│   - Service interactions
       │  API Contracts      │   - Repository behavior
       │                     │   - Domain logic
       ├─────────────────────┤
       │    Unit Tests       │   40% of tests
       │ Implementation      │   - UI mechanics (keyboard, modals)
       │ Details             │   - Business rule calculations
       │                     │   - Utility functions
       │                     │   - Data transformations
       └─────────────────────┘
```

### How It Applies

**BDD (20% - ~260 scenarios in KaRiya)**
- Career event capture workflow
- Event enrichment and review
- Burst suggestion and acceptance
- Skill inference and management
- Multi-intent workflows
- Error recovery scenarios

**Integration (40%)**
- Event repository operations
- Skill service interactions
- Burst suggestion algorithm
- Event filtering and sorting logic
- Database transactions

**Unit (40%)**
- Keyboard shortcut handlers (j/k, arrows, /f/s)
- Modal opening/closing
- Form field tab navigation
- List scrolling mechanics
- Focus management
- Validation functions
- Utility calculations

**This allocation ensures**:
- ✅ Business value is tested (BDD)
- ✅ Features integrate correctly (Integration)
- ✅ Implementation is solid (Unit)
- ✅ Tests run fast (most are unit)
- ✅ Scenarios remain maintainable

---

## Anti-Patterns to Avoid

### 1. Modal Mechanics Testing

| Aspect | ❌ Anti-Pattern | ✅ Correct |
|--------|---|---|
| **What** | Testing modal open/close | Testing business outcome |
| **Example** | "When I press f Then modal appears" | "When I filter by company Then results show only that company" |
| **Where** | ❌ BDD feature files | ✅ Unit tests (UIKit modal_test.go) |
| **Breaks When** | UI library changes | Business logic changes |

### 2. Keyboard Navigation Testing

| Aspect | ❌ Anti-Pattern | ✅ Correct |
|--------|---|---|
| **What** | Testing key handlers | Testing business navigation |
| **Example** | "When I press j Then selection moves down" | "When I browse events Then I see them in order" |
| **Where** | ❌ BDD feature files | ✅ Unit tests (handlers_test.go) |
| **Breaks When** | Shortcut keys change | Business workflows change |

### 3. Form Field Navigation

| Aspect | ❌ Anti-Pattern | ✅ Correct |
|--------|---|---|
| **What** | Testing tab order | Testing data entry |
| **Example** | "When I tab Then focus moves to next field" | "When I enter event data Then event is created" |
| **Where** | ❌ BDD feature files | ✅ Unit tests (form_test.go) |
| **Breaks When** | Form library changes | Business rules change |

### 4. Screen Transition Testing

| Aspect | ❌ Anti-Pattern | ✅ Correct |
|--------|---|---|
| **What** | Testing screen changes | Testing business transitions |
| **Example** | "When I press escape Then I return to menu" | "When I cancel event capture Then data is not saved" |
| **Where** | ❌ BDD feature files | ✅ Integration tests |
| **Breaks When** | Navigation changes | Business workflows change |

### 5. Element Visibility Checking

| Aspect | ❌ Anti-Pattern | ✅ Correct |
|--------|---|---|
| **What** | "I should see button X" | "Button X should be enabled after I complete form" |
| **Example** | Then I should see a save button | Then the save operation should succeed |
| **Where** | ❌ BDD feature files | ✅ Unit tests |
| **Breaks When** | UI layout changes | Business rules change |

---

## Writing Good Scenarios

### Template: Business Outcome Scenario

```gherkin
Scenario: [Business goal that user cares about]
  Given [precondition - user has setup state they care about]
  When [user takes action that matters to business]
  And [potentially multi-step business action]
  Then [observable business outcome happens]
  And [related business outcome is also true]
```

### Quality Checklist

**✅ Good Scenario Should**:
- [ ] Describe a business workflow or outcome
- [ ] Use business language (not technical)
- [ ] Hide implementation details (how - in step definitions)
- [ ] Be <15 steps (sign of too many details)
- [ ] Survive UI changes (HTML, keyboard, mouse, CLI)
- [ ] Have clear Given-When-Then structure
- [ ] Test one business rule/workflow
- [ ] Provide business value documentation
- [ ] Not duplicate unit test coverage
- [ ] Not test purely technical concerns

**❌ Red Flags**:
- [ ] Mentions specific keyboard keys (j, k, /, f, s)
- [ ] Tests "I should see [UI element]"
- [ ] Includes Tab, Shift+Tab navigation
- [ ] Tests modal appearance
- [ ] Focuses on form field order
- [ ] Tests animation or timing
- [ ] Tests internal screen structure
- [ ] Has >15 steps
- [ ] Uses UI terminology (click, type, enter)

---

## Examples

### ✅ GOOD: Business Logic Focus

```gherkin
Feature: Event Enrichment
  As a professional
  I want the system to intelligently enrich my events
  So I can quickly build my career portfolio

  Scenario: System suggests related events for burst creation
    Given I captured event "Designed authentication system"
    And I captured event "Implemented OAuth provider"
    And I captured event "Added JWT token validation"
    When I select enrichment review
    Then the system should suggest grouping them into burst "Security Implementation"
    And the burst should contain 3 events
    And each event should link to the burst

  Scenario: Professional accepts AI-suggested burst
    Given the system suggested burst "Database Architecture"
    When I review and accept the suggestion
    Then the burst should be created
    And all suggested events should be linked
    And the burst should be marked as professional-approved

  Scenario: System validates burst has minimum content
    Given I'm creating a burst
    When I attempt to save with only 1 event
    Then the system should reject the burst
    And should require minimum 2 events
    And should preserve the entered data
```

**Why These Are Good**:
- Focus on business outcomes
- No mention of keys, clicks, or modals
- Test business rules (minimum events, suggestion logic)
- Read as documentation for business users
- Survive UI changes

### ❌ WRONG: Implementation Focus

```gherkin
Feature: Event Enrichment UI
  # ANTI-PATTERN: Focus on UI mechanics

  Scenario: Open enrichment modal
    When I press "/" to open search
    Then the enrichment modal should appear
    And the modal should be centered
    And I should see text input field

  Scenario: Navigate modal with keyboard
    When I press "j" key
    Then the selection should move down
    And the item should be highlighted
    And I should see focus border

  Scenario: Tab through form fields
    When I press Tab key
    Then focus should move to next field
    And the previous field should blur

  Scenario: Close modal with escape
    When I press Escape key
    Then the modal should close
    And I should return to main screen
```

**Why These Are Wrong**:
- Focus on UI mechanics, not business value
- Test keyboard implementation, not workflows
- Duplicate unit test coverage
- Break when UI changes
- Provide no business documentation

---

## Decision Framework

### Use BDD When

✅ **DO use BDD (feature files)** for:
- **Business workflows**: Event capture, review, enrichment
- **Multi-feature integration**: Events + Skills + Bursts together
- **Critical user paths**: Main functionality users depend on
- **Business rules**: Validation, constraints, outcomes
- **Error scenarios**: What happens when things go wrong
- **Data verification**: Event created with correct metadata

### Use Unit Tests When

✅ **DO use unit tests** for:
- **Keyboard handlers**: j/k navigation, shortcuts
- **Form mechanics**: Tab order, focus management
- **Modal display**: Opening, closing, visibility
- **Scroll behavior**: List pagination, viewport handling
- **UI calculations**: Layout measurements, positioning
- **Algorithms**: Business rule calculations, validation logic

### Use Integration Tests When

✅ **DO use integration tests** for:
- **API contracts**: Service interactions
- **Repository operations**: Database transactions
- **Service layer**: Multi-service workflows
- **State management**: Object lifecycle
- **Event publishing**: Message passing between systems

---

## Key Takeaways

| Concept | Remember |
|---------|----------|
| **BDD Purpose** | Document business outcomes, not implementation |
| **UI Mechanics** | Belong in unit tests, NOT BDD |
| **Modal Testing** | Is not business logic - test in unit tests |
| **Keyboard Keys** | Are implementation details - test in unit tests |
| **Test Pyramid** | BDD 20%, Integration 40%, Unit 40% |
| **Scenario Length** | Should be <15 steps (sign of too many details) |
| **Business Language** | Use terms business users understand |
| **Hide How** | Step definitions hide implementation details |
| **Test Outcomes** | What the system does, not how it does it |
| **Maintainability** | Good BDD survives UI framework changes |

---

## Resources

- **Official Cucumber.io**: https://cucumber.io/docs/bdd/
- **Gherkin Reference**: https://cucumber.io/docs/gherkin/
- **BDD Best Practices**: https://cucumber.io/docs/guides/bdd-in-action/
- **Anti-Patterns**: https://cucumber.io/docs/bdd/anti-patterns/
- **KaRiya BDD Examples**: See `features/` directory

---

## Questions?

For questions about BDD guidelines in KaRiya:
1. Review this document and examples
2. Check `features/` for good scenario examples
3. Look at related unit tests for comparison
4. Ask the team in #testing channel

**Last Updated**: February 12, 2026  
**Maintained By**: Testing Working Group
