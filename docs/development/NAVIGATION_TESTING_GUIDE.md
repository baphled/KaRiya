# Navigation Testing Guide

**Last Updated**: 2026-01-12  
**Target Audience**: Developers implementing or modifying TUI intents  
**Related**: [Navigation Testing Checklist](NAVIGATION_TESTING_CHECKLIST.md), [TUI Standards](../TUI_STANDARDS.md)

---

## Table of Contents

1. [Overview](#overview)
2. [Why Navigation Testing Matters](#why-navigation-testing-matters)
3. [Navigation Requirements](#navigation-requirements)
4. [E2E Test Framework](#e2e-test-framework)
5. [Navigation Test Patterns](#navigation-test-patterns)
6. [Escape Key Testing](#escape-key-testing)
7. [Common Navigation Bugs](#common-navigation-bugs)
8. [State Transition Testing](#state-transition-testing)
9. [Integration vs Unit Testing](#integration-vs-unit-testing)
10. [Debugging Navigation Issues](#debugging-navigation-issues)
11. [Best Practices](#best-practices)

---

## Overview

Navigation is the **backbone of the TUI user experience**. Users must be able to:
- Navigate forward through workflows
- Navigate backward (Escape key)
- Return to main menu from anywhere ('m' key)
- Quit the application ('q' key)
- Recover from errors gracefully

**Critical principle**: If navigation breaks, **the entire TUI is unusable**. Navigation testing is non-negotiable.

### What This Guide Covers

- **E2E test framework** (`internal/testutil/e2e/`)
- **Navigation test patterns** (real examples from codebase)
- **Escape key behavior** (per state type: root, intermediate, async, final)
- **Common bugs** and how to prevent them
- **Debugging techniques** for navigation issues

### What This Guide Does NOT Cover

- UI styling (see [TUI Standards](../TUI_STANDARDS.md))
- Intent architecture (see [TUI Intent Diagram](../TUI_INTENT_DIAGRAM.md))
- Workflow documentation (see [Workflow Guides](../workflows/README.md))

---

## Why Navigation Testing Matters

### Real-World Impact

**Without Navigation Tests**:
- Users get stuck in workflows (cannot go back)
- Escape key does nothing or crashes
- State machine errors cause panics
- Back navigation loses user context
- Application becomes unusable

**With Navigation Tests**:
- ✅ All states have clear escape behavior
- ✅ State transitions are verified
- ✅ Back navigation preserves context
- ✅ Error states allow recovery
- ✅ Users can always return to main menu

### KaRiya's Navigation Testing

**Current Status** (as of 2026-01-12):
- **11 navigation test files** (~2,500 lines)
- **10 E2E workflow test files** (~1,800 lines)
- **Error recovery tests** (380 lines)
- **Chained workflow tests** (280 lines)
- **Total**: ~4,960 lines of navigation testing

**Result**: Zero known navigation bugs in production

---

## Navigation Requirements

All KaRiya TUI intents must follow these navigation standards (from [TUI_STANDARDS.md](../TUI_STANDARDS.md)):

### Universal Shortcuts (Must Work Everywhere)

| Key | Behavior | Context |
|-----|----------|---------|
| **Esc** | Go back one state (or cancel if root) | All states |
| **m** | Return to main menu | All states |
| **q** / **Ctrl+C** | Quit application | All states |
| **Tab** / **Shift+Tab** | Move between fields | Forms |
| **↑↓** / **jk** | Navigate lists | Lists |
| **Enter** | Confirm/Select | Lists, Forms, Confirmations |

### State Type Navigation Rules

KaRiya intents use **4 state types** with different escape behavior:

#### 1. Root State
**Definition**: First state in a workflow (entry point)  
**Escape behavior**: Cancel intent, return to main menu  
**Example**: `CaptureStateChooseStrategy`, `GenerateCVStateSelectProfile`

```go
// Root state - escape cancels the entire intent
case CaptureStateChooseStrategy:
    if key.Type == tea.KeyEsc {
        c.active = false
        c.result = &IntentResult[*CaptureEventResult]{
            Status: Cancelled,
        }
        return tea.Quit
    }
```

#### 2. Intermediate State
**Definition**: Has a previous state to go back to  
**Escape behavior**: Go back to previous state (preserve context)  
**Example**: `CaptureStateForm`, `GenerateCVStateSelectAudience`

```go
// Intermediate state - escape goes back
case CaptureStateForm:
    if key.Type == tea.KeyEsc {
        c.state.currentState = CaptureStateChooseStrategy
        return nil
    }
```

#### 3. Async Operation State
**Definition**: Background work in progress (cannot interrupt immediately)  
**Escape behavior**: Let operation complete, then navigate back  
**Example**: `GenerateCVStateGenerating`, `ExportArtifactStateExporting`

```go
// Async state - escape allowed but operation completes first
case GenerateCVStateGenerating:
    if key.Type == tea.KeyEsc {
        // Mark for cancellation but let current operation finish
        g.state.currentState = GenerateCVStateSelectAudience
        return nil
    }
```

#### 4. Final State
**Definition**: Workflow complete (success or error)  
**Escape behavior**: Retry (if error) or go back (if success)  
**Example**: `CaptureStateSubmit` (error), `GenerateCVStateExportComplete` (success)

```go
// Final state - escape retries or goes back
case CaptureStateSubmit:
    if c.state.error != nil && key.Type == tea.KeyEsc {
        // Error state - go back to review to retry
        c.state.currentState = CaptureStateReview
        // Preserve error for user context
        return nil
    }
```

---

## E2E Test Framework

KaRiya provides a comprehensive E2E test framework at `internal/testutil/e2e/`.

### TestEnv Structure

```go
type TestEnv struct {
    T            TestingT              // Works with *testing.T and GinkgoT()
    Model        *app.Model            // Root application model
    DB           *sql.DB               // SQLite database (if using Setup)
    EventRepo    *careerrepo.SQLiteRepository
    BurstRepo    *careerrepo.SQLiteBurstRepository
    FactRepo     *careerrepo.SQLiteFactRepository
    Service      *careerservice.Service
    Ctx          context.Context
}
```

### Setup Functions

**`Setup(t TestingT) *TestEnv`** - Full persistence with SQLite
```go
env := e2e.Setup(GinkgoT())
defer env.Cleanup()
// Use for tests that verify database persistence
```

**`SetupWithMemory(t TestingT) *TestEnv`** - Fast in-memory testing
```go
env := e2e.SetupWithMemory(GinkgoT())
defer env.Cleanup()
// Use for navigation tests (faster, no I/O)
```

### Navigation Helper Methods

#### Key Press Methods

```go
// Single key press
env.PressKey(tea.KeyEnter)           // Press Enter
env.PressKey(tea.KeyEsc)             // Press Escape
env.PressKeyRune('j')                // Press 'j' (vim down)
env.PressKeyRune('m')                // Press 'm' (main menu)

// Multiple keys
env.PressKeys(tea.KeyDown, tea.KeyDown, tea.KeyEnter)
env.PressKeys('j', 'j', tea.KeyEnter)

// Type text
env.TypeText("Led team standup")
```

#### Navigation Shortcuts

```go
env.NavigateDown()        // Press 'j' (or down arrow)
env.NavigateUp()          // Press 'k' (or up arrow)
env.Confirm()             // Press Enter
env.Cancel()              // Press Escape
env.GoBack()              // Press Escape (alias)
env.Quit()                // Press 'q'
env.Tab()                 // Press Tab
```

#### Intent Selection

```go
// By index (0-based)
env.SelectIntent(2)       // Select 3rd menu item

// By name (recommended - more readable)
env.SelectIntentByName("capture_event")
env.SelectIntentByName("browse_timeline")
env.SelectIntentByName("generate_cv")
env.SelectIntentByName("export_artifact")
env.SelectIntentByName("configure_system")
```

#### View Assertions

```go
// Check view contains text
env.AssertViewContains("Capture Event")

// Check view does NOT contain text
env.AssertViewNotContains("Error")

// Check view contains any of the strings (OR logic)
env.AssertViewContainsAny("Quick", "Manual", "Strategy")

// Get raw view for custom assertions
view := env.GetView()
Expect(view).To(ContainSubstring("Timeline"))
```

#### Data Verification

```go
// Verify counts
env.AssertEventCount(5)
env.AssertBurstCount(2)
env.AssertFactCount(3)

// Get data
events := env.GetEvents()
bursts := env.GetBursts()
facts := env.GetFacts()
```

#### Test Data Population

```go
// Populate test data (events, bursts, facts)
env.PopulateTestData(5, 2, 3)  // 5 events, 2 bursts, 3 facts

// Add individual items
env.AddEvent(&career.CareerEvent{...})
env.AddBurst(&career.Burst{...})
env.AddFact(&career.Fact{...})

// Use fixtures
events := e2e.CreateSampleEvents(10)
bursts := e2e.CreateSampleBursts(3, events)
facts := e2e.CreateSampleFacts(5, events)
```

#### Session Simulation

```go
// Simulate application restart (SQLite only)
env.SimulateRestart()
env.AssertEventCount(5)  // Data persists across restart
```

---

## Navigation Test Patterns

### Pattern 1: Forward Navigation Test

**Purpose**: Verify users can navigate forward through workflow states

**Example**: CaptureEvent workflow (Strategy → Form → Review → Submit)

```go
var _ = Describe("Forward Navigation", func() {
    var env *e2e.TestEnv

    BeforeEach(func() {
        env = e2e.SetupWithMemory(GinkgoT())
    })

    AfterEach(func() {
        env.Cleanup()
    })

    It("should navigate through complete workflow", func() {
        // Start workflow
        env.SelectIntentByName("capture_event")
        env.AssertViewContainsAny("Quick", "Manual", "Strategy")

        // Select strategy
        env.Confirm()  // Select Quick
        env.AssertViewContains("Event")
        env.AssertViewContains("Text")

        // Fill form (simplified - huh forms maintain internal state)
        env.TypeText("Led team standup")
        env.Tab()
        // Form submission would happen here (not shown for brevity)
    })
})
```

**Key Principles**:
- ✅ Test complete workflow end-to-end
- ✅ Verify each state transition
- ✅ Use `AssertViewContains` to verify state changes
- ✅ Chain method calls for readability

### Pattern 2: Back Navigation Test

**Purpose**: Verify Escape key navigates back through states

**Example**: GenerateCV workflow (backward navigation)

```go
var _ = Describe("Back Navigation", func() {
    var env *e2e.TestEnv

    BeforeEach(func() {
        env = e2e.SetupWithMemory(GinkgoT())
        env.PopulateTestData(5, 0, 3)  // Need events and facts for CV generation
    })

    AfterEach(func() {
        env.Cleanup()
    })

    It("should navigate back through all states to main menu", func() {
        // Navigate forward to deep state
        env.SelectIntentByName("generate_cv")
        env.Confirm()  // Select profile
        env.AssertViewContainsAny("Audience", "hiring", "recruiter")

        // Navigate backward
        env.Cancel()  // Back to profile selection
        env.AssertViewContainsAny("Profile", "Senior", "Staff")

        env.Cancel()  // Back to main menu
        env.AssertViewContains("Capture Event")
    })

    It("should preserve error context when navigating back", func() {
        // ... trigger error state ...
        
        // Navigate back - error should be preserved
        env.Cancel()
        env.AssertViewContains("Error")  // Error message still visible
    })
})
```

**Key Principles**:
- ✅ Test backward navigation from every state
- ✅ Verify state transitions are correct
- ✅ Verify context preservation (error messages, selections, etc.)
- ✅ Verify root state cancels to main menu

### Pattern 3: Cancel at Each State Test

**Purpose**: Verify Escape key works correctly at every state

**Example**: CaptureEvent (systematic Escape key testing)

```go
var _ = Describe("Cancel at Each State", func() {
    var env *e2e.TestEnv

    BeforeEach(func() {
        env = e2e.SetupWithMemory(GinkgoT())
    })

    AfterEach(func() {
        env.Cleanup()
    })

    It("should cancel from strategy selection and return to menu", func() {
        env.SelectIntentByName("capture_event")
        env.Cancel()
        env.AssertViewContains("Capture Event")
        env.AssertViewNotContains("Select Capture Strategy")
    })

    It("should navigate back when cancelling from form", func() {
        env.SelectIntentByName("capture_event")
        env.Confirm()  // Select Quick strategy
        env.Cancel()   // Cancel from form
        // Should navigate back - either to strategy or menu
        env.AssertViewNotContains("Enter Details")
    })

    It("should navigate back when cancelling from review", func() {
        env.SelectIntentByName("capture_event")
        env.Confirm()  // Select Quick strategy
        // ... fill form and submit ...
        env.Cancel()   // Cancel from review
        env.AssertViewContains("Form")  // Back to form
    })
})
```

**Key Principles**:
- ✅ Test Escape key at EVERY state
- ✅ Verify correct state transition for each
- ✅ Verify root state cancels to menu
- ✅ Verify intermediate states go back one step

### Pattern 4: Multi-Intent Navigation Test

**Purpose**: Verify navigation between different intents

**Example**: Chained workflow testing

```go
var _ = Describe("Multi-Intent Navigation", func() {
    var env *e2e.TestEnv

    BeforeEach(func() {
        env = e2e.Setup(GinkgoT())  // Use SQLite for persistence
        env.PopulateTestData(5, 2, 3)
    })

    AfterEach(func() {
        env.Cleanup()
    })

    It("should navigate through multiple intents in sequence", func() {
        // Browse
        env.SelectIntentByName("browse_timeline")
        env.AssertViewContainsAny("Timeline", "Events")
        env.Cancel()

        // Generate CV
        env.SelectIntentByName("generate_cv")
        env.AssertViewContainsAny("Profile", "Select")
        env.Cancel()

        // Export
        env.SelectIntentByName("export_artifact")
        env.AssertViewContainsAny("Export", "Type")
        env.Cancel()

        // Should be back at menu
        env.AssertViewContains("Capture Event")
    })

    It("should maintain data across intent navigation", func() {
        // Navigate to Browse and verify data
        env.SelectIntentByName("browse_timeline")
        env.AssertViewContains("Acme Corp")
        env.Cancel()

        // Navigate to FactManagement and verify same data
        env.SelectIntentByName("fact_management")
        env.AssertViewContainsAny("Reduced", "Mentored")
        env.Cancel()

        // Data should still be there
        env.AssertEventCount(5)
        env.AssertFactCount(3)
    })
})
```

**Key Principles**:
- ✅ Test navigation between multiple intents
- ✅ Verify data persistence across intent switches
- ✅ Verify return to menu between intents
- ✅ Test realistic user workflows

### Pattern 5: Error Recovery Navigation Test

**Purpose**: Verify navigation works during error states

**Example**: Error state escape behavior

```go
var _ = Describe("Error Recovery", func() {
    var env *e2e.TestEnv

    BeforeEach(func() {
        env = e2e.SetupWithMemory(GinkgoT())
    })

    AfterEach(func() {
        env.Cleanup()
    })

    It("should allow navigation back from error state", func() {
        // ... trigger error ...
        
        env.Cancel()  // Escape from error
        // Should go back to previous valid state
        env.AssertViewNotContains("Error")
    })

    It("should preserve error message when going back", func() {
        // ... trigger error ...
        
        originalView := env.GetView()
        Expect(originalView).To(ContainSubstring("Error"))

        env.Cancel()  // Go back
        // Error context should be preserved
        view := env.GetView()
        Expect(view).To(ContainSubstring("Error"))
    })
})
```

**Key Principles**:
- ✅ Test navigation during error states
- ✅ Verify error messages are preserved
- ✅ Verify user can retry or go back
- ✅ Test recovery to valid state

### Pattern 6: Boundary Condition Navigation Test

**Purpose**: Verify navigation at list boundaries (prevent crashes)

**Example**: List navigation edge cases

```go
var _ = Describe("Boundary Navigation", func() {
    var env *e2e.TestEnv

    BeforeEach(func() {
        env = e2e.Setup(GinkgoT())
        env.PopulateTestData(3, 0, 3)
    })

    AfterEach(func() {
        env.Cleanup()
    })

    It("should handle navigating past list start", func() {
        env.SelectIntentByName("browse_timeline")
        // Press up multiple times at the start
        for i := 0; i < 5; i++ {
            env.PressKey(tea.KeyUp)
        }
        view := env.GetView()
        Expect(view).NotTo(ContainSubstring("panic"))
        Expect(view).NotTo(BeEmpty())
    })

    It("should handle navigating past list end", func() {
        env.SelectIntentByName("browse_timeline")
        // Press down more times than items
        for i := 0; i < 10; i++ {
            env.PressKey(tea.KeyDown)
        }
        view := env.GetView()
        Expect(view).NotTo(ContainSubstring("panic"))
        Expect(view).NotTo(BeEmpty())
    })

    It("should handle rapid navigation", func() {
        env.SelectIntentByName("browse_timeline")
        for i := 0; i < 20; i++ {
            env.PressKeyRune('j')
            env.PressKeyRune('k')
        }
        view := env.GetView()
        Expect(view).NotTo(ContainSubstring("panic"))
    })
})
```

**Key Principles**:
- ✅ Test navigation at list boundaries
- ✅ Test rapid key presses
- ✅ Verify no panics or crashes
- ✅ Verify view remains functional

---

## Escape Key Testing

Escape key behavior is **the most critical navigation test**. Every state must handle Escape correctly.

### Escape Key Test Matrix

Use this matrix to ensure complete Escape key coverage:

| State Type | Escape Behavior | Go Back To | Test Priority |
|------------|-----------------|------------|---------------|
| **Root** | Cancel intent | Main menu | **CRITICAL** |
| **Intermediate** | Go back | Previous state | **CRITICAL** |
| **Async** | Let complete, then back | Previous state | **HIGH** |
| **Final (Success)** | Go back | Previous state | **MEDIUM** |
| **Final (Error)** | Retry/Go back | Previous state | **HIGH** |

### Escape Key Test Template

Use this template for every intent:

```go
var _ = Describe("IntentName - Escape Key Behavior", func() {
    var intent *YourIntent

    BeforeEach(func() {
        // Setup intent
    })

    Describe("RootState (Cancel to Menu)", func() {
        It("should cancel intent when esc is pressed", func() {
            intent.state.currentState = StateRoot
            intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

            Expect(intent.active).To(BeFalse())
            Expect(intent.result).NotTo(BeNil())
            Expect(intent.result.Status).To(Equal(Cancelled))
        })
    })

    Describe("IntermediateState (Go Back)", func() {
        BeforeEach(func() {
            intent.state.currentState = StateIntermediate
        })

        It("should go back to RootState when esc is pressed", func() {
            intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

            Expect(intent.state.currentState).To(Equal(StateRoot))
            Expect(intent.active).To(BeTrue())
        })
    })

    Describe("AsyncState (Let Complete)", func() {
        BeforeEach(func() {
            intent.state.currentState = StateAsync
            intent.state.isProcessing = true
        })

        It("should go back to IntermediateState when esc is pressed", func() {
            intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

            Expect(intent.state.currentState).To(Equal(StateIntermediate))
            Expect(intent.active).To(BeTrue())
        })
    })

    Describe("FinalState (Error Recovery)", func() {
        BeforeEach(func() {
            intent.state.currentState = StateFinal
            intent.state.error = &IntentError{
                Code:    "TEST_ERROR",
                Message: "Operation failed",
            }
        })

        It("should go back to IntermediateState when esc is pressed", func() {
            intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

            Expect(intent.state.currentState).To(Equal(StateIntermediate))
            Expect(intent.active).To(BeTrue())
        })

        It("should preserve error when going back", func() {
            originalError := intent.state.error
            intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

            Expect(intent.state.error).To(Equal(originalError))
        })
    })
})
```

### Real-World Examples

**Example 1**: CaptureEvent Escape Behavior (4 states)
- `ChooseStrategy` (root) → Cancel intent
- `Form` (intermediate) → Go back to `ChooseStrategy`
- `Review` (intermediate) → Go back to `Form`
- `Submit` (final) → Go back to `Review` (preserve error if exists)

**Example 2**: GenerateCV Escape Behavior (10 states)
- `SelectProfile` (root) → Cancel intent
- `SelectAudience` (intermediate) → Go back to `SelectProfile`
- `Generating` (async) → Let complete, go back to `SelectAudience`
- `Preview` (intermediate) → Go back to `SelectAudience`
- `Review` (intermediate) → Go back to `Preview`
- `Confirm` (intermediate) → Go back to `Review`
- `ExportSelectFormat` (intermediate) → Go back to `Confirm`
- `ExportSelectLocation` (intermediate) → Go back to `ExportSelectFormat`
- `Exporting` (async) → Let complete, go back to `ExportSelectLocation`
- `ExportComplete` (final) → Go back to `ExportSelectLocation`

**See Also**: Escape key test files for all 5 primary intents:
- `internal/cli/intents/capture_event_escape_test.go` (151 lines)
- `internal/cli/intents/browse_timeline_escape_test.go` (131 lines)
- `internal/cli/intents/generate_cv_escape_test.go` (258 lines)
- `internal/cli/intents/export_artifact_escape_test.go`
- `internal/cli/intents/configure_system_escape_test.go`

---

## Common Navigation Bugs

### Bug 1: Escape Key Does Nothing

**Symptom**: Pressing Escape has no effect

**Root Cause**: Missing Escape key handler in `Update()` method

**Fix**:
```go
func (i *YourIntent) Update(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch i.state.currentState {
        case StateYourState:
            switch msg.Type {
            case tea.KeyEsc:  // ← ADD THIS
                i.state.currentState = StatePrevious
                return nil
            }
        }
    }
    return nil
}
```

**Test**:
```go
It("should handle escape key", func() {
    intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
    Expect(intent.state.currentState).To(Equal(StatePrevious))
})
```

### Bug 2: Escape Key Cancels Instead of Going Back

**Symptom**: Escape cancels the entire intent instead of going back one state

**Root Cause**: Missing state type check (treating intermediate state as root state)

**Fix**:
```go
// WRONG - cancels from intermediate state
case StateIntermediate:
    switch msg.Type {
    case tea.KeyEsc:
        i.active = false  // ← WRONG! Should go back, not cancel
        i.result = &IntentResult{Status: Cancelled}
    }

// CORRECT - goes back from intermediate state
case StateIntermediate:
    switch msg.Type {
    case tea.KeyEsc:
        i.state.currentState = StatePrevious  // ← CORRECT
        return nil
    }
```

**Test**:
```go
It("should go back when esc is pressed (not cancel)", func() {
    intent.state.currentState = StateIntermediate
    intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

    Expect(intent.active).To(BeTrue())  // Still active
    Expect(intent.state.currentState).To(Equal(StatePrevious))
})
```

### Bug 3: Lost Context on Back Navigation

**Symptom**: Navigating back loses user's selections, errors, or form data

**Root Cause**: Not preserving state when transitioning

**Fix**:
```go
// WRONG - loses error context
case StateError:
    switch msg.Type {
    case tea.KeyEsc:
        i.state.error = nil  // ← WRONG! Error context lost
        i.state.currentState = StatePrevious
    }

// CORRECT - preserves error context
case StateError:
    switch msg.Type {
    case tea.KeyEsc:
        // Keep error so user knows what went wrong
        i.state.currentState = StatePrevious
        // i.state.error is preserved
    }
```

**Test**:
```go
It("should preserve error when going back", func() {
    originalError := &IntentError{Code: "TEST", Message: "Test error"}
    intent.state.error = originalError

    intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

    Expect(intent.state.error).To(Equal(originalError))
})
```

### Bug 4: Panic on Rapid Escape Presses

**Symptom**: Pressing Escape multiple times rapidly causes panic

**Root Cause**: State transition happens before view update, causing nil pointer dereference

**Fix**:
```go
// WRONG - can cause nil pointer panic
case StateView:
    switch msg.Type {
    case tea.KeyEsc:
        i.state.currentState = StatePrevious
        i.state.viewData = nil  // ← WRONG! View() might access this
    }

// CORRECT - preserve data until next state is stable
case StateView:
    switch msg.Type {
    case tea.KeyEsc:
        i.state.currentState = StatePrevious
        // Keep viewData - it will be replaced in next Init()
    }
```

**Test**:
```go
It("should handle rapid escape presses", func() {
    for i := 0; i < 5; i++ {
        intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
    }
    view := intent.View()
    Expect(view).NotTo(ContainSubstring("panic"))
})
```

### Bug 5: List Navigation Past Boundaries

**Symptom**: Navigating past list start/end causes panic or index out of range

**Root Cause**: No boundary checks on list navigation

**Fix**:
```go
// WRONG - no boundary check
case tea.KeyDown:
    i.selectedIndex++  // ← Can go out of bounds

// CORRECT - boundary check
case tea.KeyDown:
    if i.selectedIndex < len(i.items)-1 {
        i.selectedIndex++
    }

case tea.KeyUp:
    if i.selectedIndex > 0 {
        i.selectedIndex--
    }
```

**Test**:
```go
It("should not go below first item", func() {
    intent.selectedIndex = 0
    for i := 0; i < 5; i++ {
        intent.Update(tea.KeyMsg{Type: tea.KeyUp})
    }
    Expect(intent.selectedIndex).To(Equal(0))
})

It("should not go above last item", func() {
    intent.selectedIndex = len(intent.items) - 1
    for i := 0; i < 5; i++ {
        intent.Update(tea.KeyMsg{Type: tea.KeyDown})
    }
    Expect(intent.selectedIndex).To(Equal(len(intent.items) - 1))
})
```

### Bug 6: Main Menu Key ('m') Not Working

**Symptom**: Pressing 'm' doesn't return to main menu

**Root Cause**: 'm' key handler missing or only works in some states

**Fix**:
```go
// CORRECT - handle 'm' at the top level (all states)
func (i *YourIntent) Update(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle 'm' globally (before state-specific logic)
        if msg.Type == tea.KeyRunes && len(msg.Runes) > 0 && msg.Runes[0] == 'm' {
            i.active = false
            i.result = &IntentResult{Status: Cancelled}
            return tea.Quit
        }

        // Then handle state-specific keys
        switch i.state.currentState {
        // ... state handlers ...
        }
    }
    return nil
}
```

**Test**:
```go
It("should return to main menu when 'm' is pressed", func() {
    // Test from multiple states
    states := []YourState{StateRoot, StateIntermediate, StateFinal}
    for _, state := range states {
        intent.state.currentState = state
        intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
        Expect(intent.active).To(BeFalse())
    }
})
```

---

## State Transition Testing

### State Machine Verification

Every intent has a state machine. Navigation tests verify:
1. ✅ All states are reachable
2. ✅ All transitions are valid
3. ✅ Invalid transitions are rejected
4. ✅ State data is preserved correctly

### State Transition Test Pattern

```go
var _ = Describe("State Transitions", func() {
    var intent *YourIntent

    BeforeEach(func() {
        // Setup intent
    })

    Describe("Valid Transitions", func() {
        It("should transition from StateA to StateB", func() {
            intent.state.currentState = StateA
            intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
            Expect(intent.state.currentState).To(Equal(StateB))
        })

        It("should transition from StateB to StateC", func() {
            intent.state.currentState = StateB
            intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
            Expect(intent.state.currentState).To(Equal(StateC))
        })
    })

    Describe("Invalid Transitions", func() {
        It("should not transition from StateA to StateC directly", func() {
            intent.state.currentState = StateA
            // No valid action to go from A to C directly
            intent.Update(tea.KeyMsg{Type: tea.KeySpace})
            Expect(intent.state.currentState).To(Equal(StateA))  // Unchanged
        })
    })

    Describe("Backward Transitions", func() {
        It("should allow going back from StateB to StateA", func() {
            intent.state.currentState = StateB
            intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
            Expect(intent.state.currentState).To(Equal(StateA))
        })
    })
})
```

### State Coverage Matrix

Create a state coverage matrix for your intent:

| From State | To State | Key | Test Status |
|------------|----------|-----|-------------|
| Root | Intermediate | Enter | ✅ |
| Intermediate | Root | Esc | ✅ |
| Intermediate | Final | Enter | ✅ |
| Final | Intermediate | Esc | ✅ |
| Any | Cancelled | 'm' | ✅ |
| Any | Quit | 'q' | ✅ |

**Goal**: 100% coverage of all valid transitions

---

## Integration vs Unit Testing

### When to Use Integration Tests (E2E)

**Use E2E tests** (`internal/testutil/e2e/`) for:
- ✅ Complete workflow navigation (start to finish)
- ✅ Multi-intent navigation (Browse → Generate CV → Export)
- ✅ Data persistence across navigation
- ✅ User-facing scenarios (real workflows)

**Example**:
```go
// E2E test - verifies complete workflow
It("should capture event and then browse it", func() {
    env := e2e.Setup(GinkgoT())
    defer env.Cleanup()

    // Capture event
    env.SelectIntentByName("capture_event")
    env.Confirm()
    env.TypeText("Test event")
    // ... complete form submission ...

    // Browse timeline
    env.SelectIntentByName("browse_timeline")
    env.AssertViewContains("Test event")
})
```

### When to Use Unit Tests

**Use unit tests** (`*_escape_test.go`, `*_navigation_test.go`) for:
- ✅ Escape key behavior per state
- ✅ State transition logic
- ✅ Edge cases (boundary conditions, error states)
- ✅ Fast feedback during development

**Example**:
```go
// Unit test - verifies single state behavior
It("should go back when esc is pressed", func() {
    intent.state.currentState = StateForm
    intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
    Expect(intent.state.currentState).To(Equal(StateChooseStrategy))
})
```

### Test Pyramid for Navigation

```
     E2E Tests (Slow, Comprehensive)
    /                                  \
   /  Multi-Intent Workflows           \
  /   Complete User Journeys            \
 /                                       \
/─────────────────────────────────────────\
\                                         /
 \ Integration Tests (Medium Speed)     /
  \  Single Intent Workflows           /
   \ Forward + Backward Navigation    /
    \                                /
     \──────────────────────────────/
      \                            /
       \ Unit Tests (Fast)        /
        \ Escape Key per State   /
         \ State Transitions    /
          \                    /
           \──────────────────/
```

**Recommended Split**:
- **70% Unit Tests**: Fast, focused on single states
- **20% Integration Tests**: Complete workflows per intent
- **10% E2E Tests**: Multi-intent user journeys

---

## Debugging Navigation Issues

### Debug Technique 1: View Snapshots

When navigation fails, capture view snapshots to see what the user sees:

```go
It("should show strategy selection", func() {
    env.SelectIntentByName("capture_event")
    view := env.GetView()
    fmt.Println("=== VIEW SNAPSHOT ===")
    fmt.Println(view)
    fmt.Println("=== END SNAPSHOT ===")
    env.AssertViewContains("Quick")
})
```

### Debug Technique 2: State Inspection

Print current state to understand state machine issues:

```go
It("should transition to form state", func() {
    intent.state.currentState = StateChooseStrategy
    fmt.Printf("Before: %v\n", intent.state.currentState)
    
    intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
    fmt.Printf("After: %v\n", intent.state.currentState)
    
    Expect(intent.state.currentState).To(Equal(StateForm))
})
```

### Debug Technique 3: Event Tracing

Trace all key events to see what the intent receives:

```go
It("should handle key press", func() {
    originalUpdate := intent.Update
    intent.Update = func(msg tea.Msg) tea.Cmd {
        fmt.Printf("Received msg: %#v\n", msg)
        return originalUpdate(msg)
    }
    
    env.PressKey(tea.KeyEsc)
})
```

### Debug Technique 4: Ginkgo Focus

Use `FDescribe` and `FIt` to focus on failing tests:

```go
// Only run this test
FIt("should go back when esc is pressed", func() {
    intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
    Expect(intent.state.currentState).To(Equal(StatePrevious))
})
```

### Debug Technique 5: Race Detector

Run tests with race detector to find concurrency issues:

```bash
go test -race ./internal/cli/intents/...
```

### Common Debug Scenarios

**Scenario 1**: Test passes but UI doesn't work
- **Cause**: Test mocks too much, doesn't reflect real behavior
- **Fix**: Use E2E tests with full app model

**Scenario 2**: Escape key works in test but not in app
- **Cause**: Test uses `Update()` directly, bypassing BubbleTea routing
- **Fix**: Use E2E test framework with full message routing

**Scenario 3**: State transitions work but view doesn't update
- **Cause**: View() method not checking current state
- **Fix**: Add state switch in View() method

---

## Best Practices

### 1. Test Escape Key EVERYWHERE

**Rule**: Every state must have an Escape key test

```go
// Create escape test file: intent_name_escape_test.go
var _ = Describe("IntentName - Escape Key Behavior", func() {
    // Test Escape in EVERY state
})
```

### 2. Use Descriptive Test Names

**Good**:
```go
It("should go back to profile selection when escape is pressed from audience selection", func() {
    // Clear what state we're in and what we expect
})
```

**Bad**:
```go
It("should work", func() {
    // What does "work" mean?
})
```

### 3. Test One Thing Per Test

**Good**:
```go
It("should transition to form state when strategy is selected", func() {
    intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
    Expect(intent.state.currentState).To(Equal(StateForm))
})

It("should show form fields after strategy selection", func() {
    intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
    view := intent.View()
    Expect(view).To(ContainSubstring("Event"))
})
```

**Bad**:
```go
It("should transition and show form", func() {
    // Testing two things - hard to debug when it fails
})
```

### 4. Use BeforeEach for Common Setup

**Good**:
```go
Describe("Form State", func() {
    BeforeEach(func() {
        intent.state.currentState = StateForm
        // Common setup for all tests in this describe block
    })

    It("should show form fields", func() {
        // Test-specific logic only
    })
})
```

### 5. Always Use defer cleanup()

**Good**:
```go
BeforeEach(func() {
    env = e2e.Setup(GinkgoT())
})

AfterEach(func() {
    env.Cleanup()  // Always cleanup
})
```

### 6. Test Error Paths

**Don't just test happy path**:
```go
It("should handle missing data gracefully", func() {
    intent.state.data = nil
    view := intent.View()
    Expect(view).NotTo(ContainSubstring("panic"))
})
```

### 7. Use Table-Driven Tests for Similar Scenarios

**Good**:
```go
DescribeTable("should handle escape from all intermediate states",
    func(fromState, toState YourState) {
        intent.state.currentState = fromState
        intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
        Expect(intent.state.currentState).To(Equal(toState))
    },
    Entry("Form → ChooseStrategy", StateForm, StateChooseStrategy),
    Entry("Review → Form", StateReview, StateForm),
    Entry("Submit → Review", StateSubmit, StateReview),
)
```

### 8. Test Both Arrow Keys and Vim Keys

**Good**:
```go
It("should navigate down with down arrow", func() {
    env.PressKey(tea.KeyDown)
    // ... assertions ...
})

It("should navigate down with 'j' key", func() {
    env.PressKeyRune('j')
    // ... assertions ...
})
```

### 9. Verify No Panics

**Always check for panics in edge cases**:
```go
It("should not panic with nil data", func() {
    intent.state.data = nil
    Expect(func() { intent.View() }).NotTo(Panic())
})
```

### 10. Use Method Chaining for Readability

**Good**:
```go
env.SelectIntentByName("capture_event").
    Confirm().
    TypeText("Test event").
    Tab().
    AssertViewContains("Date")
```

**vs**:
```go
env.SelectIntentByName("capture_event")
env.Confirm()
env.TypeText("Test event")
env.Tab()
env.AssertViewContains("Date")
```

---

## Next Steps

After reading this guide:

1. **Read**: [Navigation Testing Checklist](NAVIGATION_TESTING_CHECKLIST.md) for quick reference
2. **Review**: Existing navigation test files:
   - `internal/cli/intents/*_navigation_test.go` (11 files)
   - `internal/cli/intents/*_escape_test.go` (5 files)
   - `internal/testutil/e2e/*_e2e_test.go` (10 files)
3. **Practice**: Write navigation tests for your next feature
4. **Verify**: Run `make test` to ensure all navigation tests pass

## Related Documentation

- [Navigation Testing Checklist](NAVIGATION_TESTING_CHECKLIST.md) - Quick reference
- [TUI Standards](../TUI_STANDARDS.md) - UI/UX requirements
- [TUI Developer Guide](../TUI_DEVELOPER_GUIDE.md) - Creating TUI components
- [TUI Intent Diagram](../TUI_INTENT_DIAGRAM.md) - Intent architecture
- [Workflow Guides](../workflows/README.md) - User workflow documentation
- [Keyboard Shortcuts Guide](../KEYBOARD_SHORTCUTS_GUIDE.md) - User shortcut reference
- [Keyboard System Guide](KEYBOARD_SYSTEM_GUIDE.md) - Developer keyboard implementation

---

*Last Updated: 2026-01-12*  
*Maintainer: KaRiya Development Team*
