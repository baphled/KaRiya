---
name: e2e-testing
description: KaRiya end-to-end testing patterns using the TestEnv harness
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

# E2E Testing Skill

You are an expert in KaRiya's end-to-end testing patterns using the TestEnv harness.

## Overview

KaRiya has a custom E2E test framework in `internal/testutil/e2e/` that simulates TUI interactions via Bubble Tea message injection.

## Test File Naming

```
*_e2e_test.go           # E2E workflow tests
*_navigation_test.go    # Navigation-focused tests
*_escape_test.go        # Escape/cancel behavior tests
*_state_machine_test.go # State transition tests
```

## TestEnv Setup

### Shared Environment Pattern

```go
// e2e_suite_test.go
package e2e_test

import (
    "testing"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "github.com/baphled/kariya/internal/testutil/e2e"
)

var sharedEnv *e2e.TestEnv

func TestE2E(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "E2E Suite")
}

var _ = BeforeSuite(func() {
    var err error
    sharedEnv, err = e2e.Setup()
    Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
    if sharedEnv != nil {
        sharedEnv.Cleanup()
    }
})
```

### Getting Environment in Tests

```go
var _ = Describe("Capture Workflow", func() {
    var env *e2e.TestEnv
    
    BeforeEach(func() {
        env = e2e.GetSharedEnv()
        env.Reset()  // Clean state for each test
    })
    
    It("captures an event", func() {
        // Test using env
    })
})
```

---

## Navigation Helpers

### Menu Navigation

```go
// Select intent from main menu by name
env.SelectIntentByName("Capture Event")

// Navigate down in list
env.NavigateDown()
env.NavigateDown()

// Navigate up in list
env.NavigateUp()

// Confirm selection (Enter)
env.Confirm()

// Cancel/Go back (Escape)
env.Cancel()
```

### Keyboard Simulation

```go
// Press special key
env.PressKey(tea.KeyEsc)
env.PressKey(tea.KeyEnter)
env.PressKey(tea.KeyUp)
env.PressKey(tea.KeyDown)
env.PressKey(tea.KeyTab)

// Press rune key
env.PressKeyRune('a')  // Action key
env.PressKeyRune('d')  // Delete
env.PressKeyRune('e')  // Edit
env.PressKeyRune('q')  // Quit
env.PressKeyRune('/')  // Search

// Type text (for forms)
env.TypeText("My event description")
```

### Message Injection

```go
// Send any Bubble Tea message
env.SendMessage(tea.WindowSizeMsg{Width: 120, Height: 40})

// Send custom message
env.SendMessage(MyCustomMsg{Data: "test"})
```

---

## View Assertions

### Content Checks

```go
// Assert view contains text
env.AssertViewContains("Capture Event")
env.AssertViewContains("Events (15)")

// Assert view does NOT contain text
env.AssertViewNotContains("Error")
env.AssertViewNotContains("Loading")

// Get current view for custom assertions
view := env.GetView()
Expect(view).To(ContainSubstring("expected text"))
```

### State Checks

```go
// Check current intent
Expect(env.GetCurrentIntentName()).To(Equal("capture"))

// Check if modal is visible
Expect(env.IsModalVisible()).To(BeTrue())

// Check active screen type
Expect(env.GetActiveScreenType()).To(Equal("list"))
```

---

## Data Setup & Verification

### Adding Test Data

```go
// Add event directly
event := fixtures.Event(1)
env.AddEvent(event)

// Add multiple events
events := fixtures.Events(5)
for _, e := range events {
    env.AddEvent(e)
}

// Add with specific attributes
event := fixtures.EventWith(
    fixtures.WithText("Led project kickoff"),
    fixtures.WithCompany("Acme Corp"),
    fixtures.WithDate(time.Now()),
)
env.AddEvent(event)
```

### Verifying Data

```go
// Check event count
env.AssertEventCount(5)

// Get all events
events := env.GetEvents()
Expect(events).To(HaveLen(5))

// Find specific event
event := env.FindEventByText("Led project kickoff")
Expect(event).ToNot(BeNil())
```

---

## Common E2E Patterns

### Testing Capture Workflow

```go
var _ = Describe("Capture Event", func() {
    var env *e2e.TestEnv
    
    BeforeEach(func() {
        env = e2e.GetSharedEnv()
        env.Reset()
    })
    
    It("captures a new event via Quick Capture", func() {
        // Navigate to Capture intent
        env.SelectIntentByName("Capture Event")
        
        // Select Quick Capture mode
        env.NavigateDown()  // Go to Quick Capture option
        env.Confirm()
        
        // Fill form
        env.TypeText("Led successful project delivery")
        env.PressKey(tea.KeyTab)  // Next field
        env.TypeText("Acme Corp")
        env.PressKey(tea.KeyTab)
        env.TypeText("Project X")
        
        // Submit form
        env.Confirm()
        
        // Verify event was created
        env.AssertEventCount(1)
        env.AssertViewContains("Event captured")
    })
})
```

### Testing Navigation

```go
var _ = Describe("Timeline Navigation", func() {
    var env *e2e.TestEnv
    
    BeforeEach(func() {
        env = e2e.GetSharedEnv()
        env.Reset()
        
        // Setup: Add test events
        for i := 1; i <= 10; i++ {
            env.AddEvent(fixtures.Event(i))
        }
    })
    
    It("navigates through event list", func() {
        env.SelectIntentByName("Browse Timeline")
        
        // Navigate down through list
        env.NavigateDown()
        env.NavigateDown()
        env.NavigateDown()
        
        // Verify selection changed
        env.AssertViewContains("Event 4")  // Or check selection indicator
        
        // Navigate up
        env.NavigateUp()
        env.AssertViewContains("Event 3")
    })
    
    It("opens event detail on Enter", func() {
        env.SelectIntentByName("Browse Timeline")
        env.NavigateDown()
        env.Confirm()
        
        // Should show detail view
        env.AssertViewContains("Event Details")
    })
})
```

### Testing Escape Behavior

```go
var _ = Describe("Escape Behavior", func() {
    var env *e2e.TestEnv
    
    BeforeEach(func() {
        env = e2e.GetSharedEnv()
        env.Reset()
    })
    
    It("returns to main menu from intent list", func() {
        env.SelectIntentByName("Browse Timeline")
        env.AssertViewContains("Timeline")
        
        env.Cancel()  // Press Escape
        
        env.AssertViewContains("Main Menu")
    })
    
    It("closes modal without losing data", func() {
        env.SelectIntentByName("Browse Timeline")
        env.AddEvent(fixtures.Event(1))
        
        // Open detail modal
        env.Confirm()
        env.AssertViewContains("Event Details")
        
        // Close modal
        env.Cancel()
        
        // Should return to list, not main menu
        env.AssertViewContains("Timeline")
        env.AssertViewNotContains("Main Menu")
    })
})
```

### Testing State Machine

```go
var _ = Describe("Capture State Machine", func() {
    var env *e2e.TestEnv
    
    BeforeEach(func() {
        env = e2e.GetSharedEnv()
        env.Reset()
    })
    
    It("transitions through states correctly", func() {
        env.SelectIntentByName("Capture Event")
        
        // State: SelectMode
        env.AssertViewContains("Select capture mode")
        
        env.Confirm()  // Select default mode
        
        // State: Form
        env.AssertViewContains("Event text")
        
        env.TypeText("Test event")
        env.Confirm()
        
        // State: Review (if applicable)
        // or State: Complete
        env.AssertViewContains("Event captured")
    })
})
```

---

## Session Simulation

### Restart Simulation

```go
It("persists data across sessions", func() {
    // Create event
    env.SelectIntentByName("Capture Event")
    // ... capture event ...
    
    // Simulate app restart
    env.SimulateRestart()
    
    // Verify data persisted
    env.SelectIntentByName("Browse Timeline")
    env.AssertEventCount(1)
    env.AssertViewContains("My event")
})
```

---

## Inline Comments in E2E Tests

**EXCEPTION**: E2E tests are the ONLY place where inline comments are allowed for readability:

```go
It("completes full capture workflow", func() {
    env.SelectIntentByName("Capture Event") // Navigate to Capture
    env.NavigateDown()                       // Select Manual mode
    env.Confirm()                            // Confirm selection
    
    env.TypeText("Led project kickoff")      // Enter event text
    env.PressKey(tea.KeyTab)                 // Move to company
    env.TypeText("Acme Corp")                // Enter company
    env.Confirm()                            // Submit form
    
    env.AssertViewContains("Event captured") // Verify success
})
```

**Why**: E2E tests describe complex multi-step user interactions where inline comments clarify intent.

---

## Test Organization

### By Feature

```
internal/testutil/e2e/
├── e2e_suite_test.go           # Suite setup
├── helpers.go                  # TestEnv implementation
├── capture_e2e_test.go         # Capture workflow tests
├── capture_navigation_test.go  # Capture navigation tests
├── browse_e2e_test.go          # Browse workflow tests
├── browse_navigation_test.go   # Browse navigation tests
└── timeline_escape_test.go     # Escape behavior tests
```

### Test Naming

```go
// Pattern: Feature + Behavior
Describe("Capture Workflow", func() {
    Context("Quick Capture mode", func() {
        It("captures event with minimal input", func() { ... })
        It("validates required fields", func() { ... })
    })
    
    Context("Manual mode", func() {
        It("allows detailed event entry", func() { ... })
    })
})
```

---

## Debugging E2E Tests

### Print View

```go
It("debug test", func() {
    env.SelectIntentByName("Browse Timeline")
    
    // Print current view for debugging
    fmt.Println("=== CURRENT VIEW ===")
    fmt.Println(env.GetView())
    fmt.Println("=== END VIEW ===")
    
    // Continue with assertions
})
```

### Step Through

```go
It("step through workflow", func() {
    env.SelectIntentByName("Capture Event")
    env.AssertViewContains("Select mode")  // Checkpoint 1
    
    env.Confirm()
    env.AssertViewContains("Event text")   // Checkpoint 2
    
    env.TypeText("Test")
    env.AssertViewContains("Test")         // Checkpoint 3
    
    // Each checkpoint helps identify where failure occurs
})
```

---

## Anti-Patterns

### DON'T: Share State Between Tests

```go
// WRONG - Tests depend on each other
var globalEvent *Event

It("creates event", func() {
    globalEvent = createEvent()
})

It("uses event", func() {
    // Fails if first test didn't run
    Expect(globalEvent).ToNot(BeNil())
})

// CORRECT - Each test sets up its own state
BeforeEach(func() {
    env.Reset()
})

It("creates and uses event", func() {
    env.AddEvent(fixtures.Event(1))
    // ... test ...
})
```

### DON'T: Skip Reset

```go
// WRONG - Leftover state from previous test
It("test without reset", func() {
    // May have events from previous test!
    env.AssertEventCount(0)  // Might fail
})

// CORRECT - Always reset
BeforeEach(func() {
    env.Reset()  // Clean slate
})
```

### DON'T: Hardcode Timing

```go
// WRONG - Flaky due to timing
time.Sleep(100 * time.Millisecond)
env.AssertViewContains("Loaded")

// CORRECT - Use Eventually
Eventually(func() string {
    return env.GetView()
}).Should(ContainSubstring("Loaded"))
```

---

## Running E2E Tests

```bash
# Run all E2E tests
go test ./internal/testutil/e2e/... -v

# Run specific E2E test file
go test ./internal/testutil/e2e/capture_e2e_test.go -v

# Run with focus
go test ./internal/testutil/e2e/... -ginkgo.focus="Quick Capture"

# Skip E2E tests (if tagged)
go test ./... -short
```

---

## Related Skills

- `ginkgo-gomega` - Testing framework basics
- `bubble-tea-expert` - TUI patterns being tested
- `test-fixtures` - Creating test data
- `cucumber` - BDD scenarios (alternative)
