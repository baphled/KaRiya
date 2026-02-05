# E2E Test Helpers Guide

**Building self-documenting test utilities for behavior-driven E2E tests**

## Overview

KaRiya uses Cucumber-style self-documenting helper functions for E2E tests. These helpers allow tests to read like natural language descriptions of user behavior, making them both executable specifications and living documentation.

## Philosophy

### ✅ Good: User-Centric Actions
```go
It("should allow editing event metadata", func() {
    env.SelectIntentByName("capture_event")
    env.Confirm()
    env.SubmitEvent(testEvent)
    env.PressKeyRune('e')  // Edit metadata
    env.TypeText("Updated company")
    env.Tab()
    env.Confirm()
    
    env.AssertViewContains("Updated company")
})
```

**Why this works:**
- Reads like user behavior ("select, confirm, submit, edit")
- No implementation details exposed
- Self-documenting - anyone can understand the test flow
- Matches real user actions exactly

### ❌ Bad: Implementation-Focused
```go
It("should allow editing event metadata", func() {
    // Don't do this - too low-level
    intent := captureevent.NewIntent(ctx, nil)
    intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
    form := intent.GetForm()
    form.SetValue("text", "test")
    intent.Update(models.FormCompletedMsg{})
    
    // Complex setup that doesn't match user actions
    screen := intent.GetActiveScreen()
    modal := screen.GetModal()
    modal.SetVisible(true)
})
```

**Why this fails:**
- Exposes internal implementation
- Breaks when refactoring
- Doesn't match user behavior
- Hard to understand intent

## Core Principles

### 1. Actions Match User Behavior

Every helper function should correspond to something a **real user** would do:

| User Action | Helper Method | Example |
|-------------|---------------|---------|
| Select menu item | `env.SelectIntentByName("capture_event")` | User navigates to "Capture Event" |
| Press Enter | `env.Confirm()` | User confirms selection |
| Type text | `env.TypeText("Build API")` | User types in form field |
| Press Tab | `env.Tab()` | User moves to next field |
| Press Escape | `env.Cancel()` | User cancels operation |
| Press key | `env.PressKeyRune('e')` | User presses 'e' key |
| Navigate down | `env.NavigateDown()` | User moves selection down |

### 2. Self-Documenting Names

Function names should read like English sentences:

```go
// ✅ Good - reads naturally
env.SelectIntentByName("browse_timeline")
env.AssertEventCount(5)
env.AssertViewContains("No events")
env.PopulateTestData(5, 0, 0)

// ❌ Bad - requires mental translation
env.SetIntent(2)
env.CheckCount(5)
env.VerifyString("No events")
env.Seed(5, 0, 0)
```

### 3. Chainable Fluent API

Helpers return `*TestEnv` to enable method chaining:

```go
env.SelectIntentByName("capture_event").
    Confirm().
    TypeText("Built REST API").
    Tab().
    TypeText("Acme Corp").
    Confirm().
    AssertViewContains("Review Enrichment")
```

This reads like a user story: "Select capture event, confirm, type text, tab, type company, confirm, check review screen shows."

### 4. Complete End-to-End Flows

Tests should exercise the **entire user journey**, not isolated components:

```go
// ✅ Good - Complete user journey
It("should capture and display event in timeline", func() {
    // User captures event
    env.SelectIntentByName("capture_event")
    env.Confirm()
    env.SubmitEvent(fixtures.EventWith("", "Built API", "", ""))
    env.Confirm()  // Save event
    
    // User views timeline
    env.SelectIntentByName("browse_timeline")
    env.AssertViewContains("Built API")
    
    // User can navigate back
    env.Cancel()
    env.AssertViewContains("Capture Event")
})

// ❌ Bad - Tests isolated component
It("should show event in list", func() {
    // Direct database insertion - not what user does
    env.AddEvent(testEvent)
    
    // Directly navigate to intent - skips menu
    intent := browsetimeline.NewIntent(ctx, events)
    view := intent.View()
    Expect(view).To(ContainSubstring("Built API"))
})
```

## Writing New Helpers

### Step 1: Identify User Action

Ask: "What would a user do in the application?"

Examples:
- "User filters timeline by date"
- "User searches for events"
- "User deletes a burst"
- "User exports CV"

### Step 2: Write Helper Signature

Use clear, intention-revealing names:

```go
// Navigation helpers
func (e *TestEnv) NavigateToMenuItem(itemName string) *TestEnv
func (e *TestEnv) OpenFilterModal() *TestEnv
func (e *TestEnv) SelectFilter(filterName string) *TestEnv

// Data interaction helpers
func (e *TestEnv) SearchFor(query string) *TestEnv
func (e *TestEnv) SelectFirstItem() *TestEnv
func (e *TestEnv) DeleteSelectedItem() *TestEnv

// Assertion helpers
func (e *TestEnv) AssertFilterApplied(filterName string) *TestEnv
func (e *TestEnv) AssertSearchResultsCount(expected int) *TestEnv
func (e *TestEnv) AssertNoResultsShown() *TestEnv
```

### Step 3: Implement Using Basic Primitives

Build complex actions from basic ones:

```go
// SearchFor combines multiple low-level actions
func (e *TestEnv) SearchFor(query string) *TestEnv {
    e.T.Helper()
    
    // Open search (pressing 's' key as user would)
    e.PressKeyRune('s')
    
    // Type search query
    e.TypeText(query)
    
    // Submit search
    e.Confirm()
    
    return e
}

// DeleteSelectedItem is a complete deletion flow
func (e *TestEnv) DeleteSelectedItem() *TestEnv {
    e.T.Helper()
    
    // User presses 'd' to delete
    e.PressKeyRune('d')
    
    // Confirm deletion in modal
    e.Confirm()
    
    return e
}
```

### Step 4: Document Expected Behavior

Add godoc explaining what the helper does from **user perspective**:

```go
// SearchFor opens the search modal, enters the query, and submits it.
// This simulates a user pressing 's', typing their search term, and pressing Enter.
//
// Expected:
//   - Application must be in a state that supports search (timeline, fact list, etc.)
//   - query must not be empty
//
// Returns:
//   - Self for method chaining
//
// Side effects:
//   - Application state transitions to show search results
//
// Example:
//   env.SelectIntentByName("browse_timeline").
//       SearchFor("API").
//       AssertSearchResultsCount(5)
func (e *TestEnv) SearchFor(query string) *TestEnv {
    e.T.Helper()
    // implementation...
}
```

## Example: Adding Filter Support

Let's walk through adding filter support to E2E tests.

### Identify User Journey

```
1. User opens browse timeline
2. User presses 'f' to open filter modal
3. User selects "Last 30 days" filter
4. User confirms filter
5. User sees filtered results
```

### Write Test First

```go
It("should filter timeline by date range", func() {
    // Arrange - populate data
    env.PopulateTestData(20, 0, 0)  // 20 events spanning different dates
    
    // Act - User filters timeline
    env.SelectIntentByName("browse_timeline").
        OpenFilterModal().
        SelectDateFilter("Last 30 days").
        ConfirmFilter().
        AssertFilterApplied("Last 30 days").
        AssertEventCountInView(5)  // Only recent events
    
    // Act - User clears filter
    env.OpenFilterModal().
        ClearFilters().
        ConfirmFilter().
        AssertEventCountInView(20)  // All events restored
})
```

### Implement Helpers

```go
// OpenFilterModal simulates user pressing 'f' to open filters.
func (e *TestEnv) OpenFilterModal() *TestEnv {
    e.T.Helper()
    e.PressKeyRune('f')
    e.AssertViewContains("Filter")  // Verify modal opened
    return e
}

// SelectDateFilter chooses a date filter option.
func (e *TestEnv) SelectDateFilter(filterName string) *TestEnv {
    e.T.Helper()
    
    // Navigate to filter option
    // (Implementation depends on how UI presents filters)
    switch filterName {
    case "Last 30 days":
        e.PressKeyRune('1')  // Assuming numbered options
    case "Last 90 days":
        e.PressKeyRune('2')
    default:
        e.T.Fatalf("Unknown filter: %s", filterName)
    }
    
    return e
}

// ConfirmFilter submits the filter selection.
func (e *TestEnv) ConfirmFilter() *TestEnv {
    e.T.Helper()
    e.Confirm()
    return e
}

// AssertFilterApplied verifies filter indicator is shown.
func (e *TestEnv) AssertFilterApplied(filterName string) *TestEnv {
    e.T.Helper()
    view := e.GetView()
    if !strings.Contains(view, filterName) {
        e.T.Fatalf("Expected filter '%s' to be shown in view", filterName)
    }
    return e
}

// AssertEventCountInView counts events shown in current view.
func (e *TestEnv) AssertEventCountInView(expected int) *TestEnv {
    e.T.Helper()
    
    // Parse view to count visible events
    // (Implementation depends on view format)
    view := e.GetView()
    count := countEventsInView(view)
    
    if count != expected {
        e.T.Fatalf("Expected %d events in view, got %d", expected, count)
    }
    
    return e
}
```

## Existing Helper Functions

### Navigation
| Helper | User Action |
|--------|-------------|
| `SelectIntentByName(name)` | Navigate to and select menu item |
| `NavigateDown()` | Move selection down (arrow/j key) |
| `NavigateUp()` | Move selection up (arrow/k key) |
| `Confirm()` | Press Enter to confirm |
| `Cancel()` | Press Escape to cancel |
| `GoBack()` | Navigate back to previous screen |
| `Tab()` | Move to next form field |

### Input
| Helper | User Action |
|--------|-------------|
| `TypeText(text)` | Type text in current field |
| `PressKeyRune(r)` | Press specific character key |
| `PressKey(key)` | Press special key (arrows, etc) |
| `PressKeys(...keys)` | Press sequence of keys |

### Data Setup
| Helper | Purpose |
|--------|---------|
| `PopulateTestData(events, bursts, facts)` | Pre-populate database with test data |
| `AddEvent(event)` | Insert single event |
| `AddBurst(burst)` | Insert single burst |
| `AddFact(fact)` | Insert single fact |
| `AddSkill(skill)` | Insert single skill |

### Assertions
| Helper | Verification |
|--------|-------------|
| `AssertViewContains(substr)` | Check view contains text |
| `AssertViewNotContains(substr)` | Check view doesn't contain text |
| `AssertViewContainsAny(...substrs)` | Check view contains any of the texts |
| `AssertEventCount(n)` | Verify database has n events |
| `AssertBurstCount(n)` | Verify database has n bursts |
| `AssertFactCount(n)` | Verify database has n facts |
| `AssertSkillCount(n)` | Verify database has n skills |

### Workflow
| Helper | Purpose |
|--------|---------|
| `SubmitEvent(event)` | Complete event capture form with data |
| `SubmitHuhForm()` | Submit current huh form |
| `SimulateRestart()` | Simulate app restart (tests persistence) |

## Best Practices

### 1. One Helper = One User Action

```go
// ✅ Good - each helper is one action
env.OpenDeleteModal()
env.ConfirmDeletion()

// ❌ Bad - helper does too much
env.DeleteItemWithConfirmation()  // Hides multiple steps
```

### 2. Use Fixtures for Test Data

```go
// ✅ Good - use fixtures
testEvent := fixtures.EventWith("", "Built API", "Acme Corp", "")
env.SubmitEvent(testEvent)

// ❌ Bad - inline data construction
event := &career.Event{
    ID: uuid.New().String(),
    Text: "Built API",
    Company: "Acme Corp",
    CreatedAt: time.Now(),
    // ... many fields
}
env.AddEvent(event)
```

### 3. Assert Frequently

```go
// ✅ Good - assert after each major action
env.SelectIntentByName("capture_event")
env.AssertViewContains("Capture Event")

env.Confirm()
env.AssertViewContains("Event text")

env.SubmitEvent(testEvent)
env.AssertViewContains("Review")

// ❌ Bad - no intermediate assertions
env.SelectIntentByName("capture_event")
env.Confirm()
env.SubmitEvent(testEvent)
env.AssertViewContains("Review")  // If this fails, where did it break?
```

### 4. Test Complete Journeys

```go
// ✅ Good - complete user journey
It("should allow user to capture event, edit metadata, and view in timeline", func() {
    // Capture
    env.SelectIntentByName("capture_event").
        Confirm().
        SubmitEvent(fixtures.EventWith("", "Built API", "", ""))
    
    // Edit metadata
    env.PressKeyRune('e').
        TypeText("Acme Corp").
        Tab().
        TypeText("Backend").
        Confirm()
    
    // Save and view
    env.Confirm()
    env.SelectIntentByName("browse_timeline")
    env.AssertViewContains("Built API")
    env.AssertViewContains("Acme Corp")
})
```

## Anti-Patterns to Avoid

### ❌ Testing Implementation Details

```go
// Bad - testing internal structure
It("should have listScreen initialized", func() {
    intent := env.GetIntent()
    Expect(intent.listScreen).ToNot(BeNil())
})
```

### ❌ Direct Database Assertions for User-Visible State

```go
// Bad - checking database instead of what user sees
env.SelectIntentByName("browse_timeline")
events := env.GetEvents()
Expect(events).To(HaveLen(5))

// Good - check what user sees
env.SelectIntentByName("browse_timeline")
env.AssertViewContainsAny("5 events", "Events: 5")
```

### ❌ Mocking in E2E Tests

```go
// Bad - mocking in E2E test
mockRepo := mocks.NewEventRepository()
mockRepo.EXPECT().List().Return(testEvents, nil)
env.SetEventRepo(mockRepo)

// Good - use real database
env.PopulateTestData(5, 0, 0)
```

## Summary

**E2E test helpers should:**
- ✅ Match real user actions
- ✅ Have self-documenting names
- ✅ Support method chaining
- ✅ Test complete user journeys
- ✅ Be composable (build complex from simple)
- ✅ Assert frequently
- ✅ Use fixtures for test data

**E2E test helpers should NOT:**
- ❌ Expose implementation details
- ❌ Mock dependencies
- ❌ Test components in isolation
- ❌ Skip user interaction steps
- ❌ Use database assertions for UI state

**Remember:** If you can't describe the helper as a user action, it probably doesn't belong in E2E tests.

---

**See also:**
- [BDD_WORKFLOW.md](BDD_WORKFLOW.md) - BDD cycle and test structure
- [NAVIGATION_TESTING_GUIDE.md](NAVIGATION_TESTING_GUIDE.md) - Testing navigation flows
- [integration-test-strategy.md](../integration-test-strategy.md) - When to use E2E vs integration tests
