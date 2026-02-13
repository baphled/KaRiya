# BDD: Good vs Bad Examples

**Purpose**: Reference guide showing concrete examples of good and bad BDD scenarios  
**Last Updated**: February 12, 2026  

---

## Table of Contents

1. [Event Capture Workflow](#event-capture-workflow)
2. [Event Filtering & Search](#event-filtering--search)
3. [Navigation Scenarios](#navigation-scenarios)
4. [Form Handling](#form-handling)
5. [Modal Interactions](#modal-interactions)
6. [Error Handling](#error-handling)

---

## Event Capture Workflow

### ✅ GOOD: Focus on Business Outcome

```gherkin
Scenario: Professional captures recent career event
  Given I want to record a recent achievement
  When I capture event with description "Designed microservices architecture"
  And I add company "Acme Corp"
  And I submit the event
  Then the system should create a new event
  And the event details should include "Designed microservices architecture"
  And the event should be associated with "Acme Corp"
```

**Why It's Good**:
- Describes what user WANTS to do (business outcome)
- Tests business rule: "Event can be created with metadata"
- Hides HOW user navigates (step definitions handle)
- Works with keyboard, mouse, or any UI
- Survives UI framework changes
- Provides business documentation

**Related Unit Tests** (Separate Tests):
```go
// Unit test for keyboard shortcut
func TestAddEventShortcut(t *testing.T) {
  // Tests that pressing 'a' opens add event form
}

// Unit test for form navigation
func TestFormTabOrder(t *testing.T) {
  // Tests that Tab moves between fields
}

// Unit test for text input
func TestEventDescriptionInput(t *testing.T) {
  // Tests typing into description field
}
```

### ❌ BAD: Focus on UI Mechanics

```gherkin
Scenario: User fills event form using keyboard
  Given I am on the main menu
  When I press "a" to add event
  Then the add event modal should appear
  And I should see a form with input fields
  
  When I click on the description field
  And I type "Designed microservices architecture"
  And I press Tab key
  Then focus should move to the company field
  
  When I see the company field is focused
  And I type "Acme Corp"
  And I press Tab key to go to submit button
  And I see the submit button is highlighted
  And I press Enter
  Then the modal should close
  And I should see the event in the timeline
```

**Why It's Bad**:
- Tests keyboard implementation (a, Tab, Enter)
- Tests modal mechanics (appear, close)
- Tests form field mechanics (focus, tab order)
- Tests UI state (button highlighted, field focused)
- Breaks when:
  - Keyboard shortcuts change (use 'c' instead of 'a')
  - UI framework changes (different modal library)
  - Form library changes (different tab order)
  - Style changes (different button highlight)
- Duplicates what unit tests already verify
- Provides no business documentation
- 20+ steps with implementation details

**What To Do Instead**: 
- This entire scenario should be split into:
  - 1 BDD scenario testing business outcome (event creation)
  - 5 unit tests covering: shortcut, modal, form, tab order, submit

---

## Event Filtering & Search

### ✅ GOOD: Business Logic

```gherkin
Scenario: Professional filters events by company
  Given I have events from multiple companies
    | company |
    | Acme Corp |
    | TechCorp |
    | StartupXYZ |
  When I filter to show only "Acme Corp" events
  Then the timeline should display 1 event
  And the event should be from "Acme Corp"
  And other company events should not be visible

Scenario: Professional searches events by keyword
  Given I have events with various descriptions
    | description |
    | Built REST API |
    | Deployed Kubernetes cluster |
    | Wrote Python script |
  When I search for "API"
  Then I should see "Built REST API" event
  And I should not see "Wrote Python script"
  And results should show 1 matching event
```

**Why It's Good**:
- Tests BUSINESS LOGIC: filtering works correctly
- Tests DATA OUTCOME: correct events shown/hidden
- Survives any UI change
- Works with list, table, modal, or any display
- Tests the business rule, not the UI

**Related Unit Tests**:
```go
// Test the actual filter logic (domain logic)
func TestEventFilterByCompany(t *testing.T) {
  // Test filtering algorithm returns correct events
}

// Test search implementation
func TestEventSearch(t *testing.T) {
  // Test search algorithm finds correct matches
}

// Test UI filter button (separate)
func TestFilterButtonOpensModal(t *testing.T) {
  // Test that pressing 'f' opens filter modal
}
```

### ❌ BAD: UI Mechanics

```gherkin
Scenario: User opens filter with keyboard
  When I press "f" to filter
  Then the filter modal should appear
  And the modal should be centered on screen
  And I should see checkbox options for companies
  And "Acme Corp" checkbox should be visible
  And the checkbox should be unchecked initially

Scenario: User selects filter option
  Given the filter modal is open
  When I navigate down with "j" key to "TechCorp" option
  And I press "x" to select it
  Then the "TechCorp" checkbox should be marked as checked
  And a checkmark should appear in the checkbox
  And the option should highlight in blue

Scenario: User confirms and applies filter
  Given I selected "TechCorp" in filter modal
  When I press "Enter" to confirm
  Then the modal should close with animation
  And the timeline should fade and redraw
  And only "TechCorp" events should appear
```

**Why It's Bad**:
- Tests keyboard implementation ('f', 'j', 'x', Enter)
- Tests modal appearance and mechanics
- Tests checkbox UI state (checked, highlighted, color)
- Tests animation ('fade and redraw')
- Breaks when UI changes (different display, different keys)
- No business value tested
- Duplicates unit test coverage
- Makes scenarios fragile

**What To Do Instead**:
- 1 BDD scenario: "Filter events by company" (business outcome)
- Multiple unit tests for: 'f' key handling, modal opening, checkbox interaction, animation timing

---

## Navigation Scenarios

### ✅ GOOD: User Workflow

```gherkin
Scenario: Professional browses career timeline
  Given I have multiple events in my timeline
  When I navigate through the timeline
  Then I should see events in chronological order
  And I should be able to view event details
  And I should be able to exit to main menu

Scenario: Professional enters intent and exits
  Given I'm on the main menu
  When I select "browse_timeline" intent
  Then I should see the timeline screen
  
  When I decide not to browse
  And I exit the timeline
  Then I should return to main menu
```

**Why It's Good**:
- Tests business workflow: "Can I access timeline and return?"
- Doesn't test HOW (keyboard shortcuts, keys)
- Survives navigation implementation changes
- Verifies user can do the business task
- Tests error recovery: "Can I get back?"

### ❌ BAD: Keyboard Mechanics

```gherkin
Scenario: User navigates list with j key
  When I press "j" key
  Then the selection should move down
  And the cursor should be on the next item

Scenario: User navigates with arrow keys
  When I press Down arrow key
  Then the cursor should move to next item

Scenario: User jumps to end with G key
  When I press "G" key
  Then the cursor should jump to last item
  And I should see the last event displayed

Scenario: User jumps to start with g key
  When I press "g" key twice
  Then the cursor should jump to first item
  And I should see the first event displayed
```

**Why It's Bad**:
- Tests keyboard implementation, not business outcome
- Each tests a different key ('j', down arrow, 'G', 'g')
- Breaks when shortcuts change
- These are pure UI mechanics
- Duplicate what unit handlers_test.go covers
- Provide no business value

**What To Do Instead**:
- REMOVE from BDD
- Each keyboard shortcut test belongs in unit tests:
  ```go
  func TestJKeyNavigation(t *testing.T)     // j key
  func TestArrowKeyNavigation(t *testing.T) // arrow keys
  func TestJumpToEnd(t *testing.T)          // G key
  func TestJumpToStart(t *testing.T)        // g key
  ```

---

## Form Handling

### ✅ GOOD: Data Entry Workflow

```gherkin
Scenario: Professional edits event metadata
  Given I have an existing event
  When I edit the event
  And I change the company from "OldCorp" to "NewCorp"
  And I change the category to "Leadership"
  And I confirm the changes
  Then the event should be updated with new company "NewCorp"
  And the event should have category "Leadership"
  And the old values should no longer apply
```

**Why It's Good**:
- Tests business outcome: event metadata changes
- Doesn't care HOW you navigate the form
- Tests business rule: "Can I update event data?"
- Works with any form UI (web, terminal, etc.)
- Survives form library changes

### ❌ BAD: Form Mechanics

```gherkin
Scenario: User navigates form with Tab key
  When I press Tab key
  Then focus should move to the next field
  And the previous field should blur
  And the new field should have focus border

Scenario: User navigates back with Shift+Tab
  When I press Shift+Tab
  Then focus should move to previous field
  And the current field should lose focus

Scenario: User types in form field
  Given the description field is focused
  When I type "New description"
  Then the text should appear in the field
  And the field content should match what I typed

Scenario: User clears form field
  Given the form field contains text
  When I press Ctrl+A to select all
  And I press Delete to clear
  Then the field should be empty
```

**Why It's Bad**:
- Tests form mechanics (Tab order, focus, blur)
- Tests keyboard input implementation
- Tests specific keyboard shortcuts (Ctrl+A, Delete)
- Each tests different form behavior
- Breaks when form library changes
- Not testing business outcome
- Duplicate unit test coverage

**What To Do Instead**:
- REMOVE from BDD
- Replace with 1 BDD scenario: "Professional updates event metadata"
- Create unit tests for each form behavior:
  ```go
  func TestFormTabNavigation(t *testing.T)      // Tab key
  func TestFormShiftTabNavigation(t *testing.T) // Shift+Tab
  func TestFormTextInput(t *testing.T)          // typing
  func TestFormFieldClear(t *testing.T)         // Ctrl+A, Delete
  ```

---

## Modal Interactions

### ✅ GOOD: Business Outcome

```gherkin
Scenario: Professional confirms event deletion
  Given I have an event to delete
  When I initiate deletion
  And I confirm the deletion
  Then the event should be removed
  And the timeline should no longer display it

Scenario: Professional cancels deletion
  Given I initiated event deletion
  When I cancel the confirmation
  Then the event should remain
  And the timeline should still display it
  And no changes should be made
```

**Why It's Good**:
- Tests business outcome: event deleted or preserved
- Doesn't test modal appearance/mechanics
- Tests business rule: "Can I safely delete?"
- Survives UI changes (modal → dialog → popup)
- Provides business documentation

### ❌ BAD: Modal Mechanics

```gherkin
Scenario: Delete confirmation modal appears
  When I select delete
  Then a confirmation modal should appear
  And the modal should be centered
  And the modal should display "Confirm deletion?"
  And I should see "Yes" and "No" buttons

Scenario: Modal closes after confirmation
  Given the delete modal is open
  And I see the confirm button
  When I press Enter to activate the button
  Then the modal should close with fade animation
  And the delete should complete
  And I should see success message

Scenario: Modal closes on cancel
  Given the delete modal is open
  When I press Escape key
  Then the modal should close
  And no deletion should occur
  And I should return to the timeline
```

**Why It's Bad**:
- Tests modal appearance ("should appear", "centered")
- Tests modal animation ("fade animation")
- Tests button display ("Yes" and "No" buttons")
- Tests keyboard handling (Enter, Escape)
- Tests UI state transitions (modal open/close)
- Breaks when: modal library changes, animation changes, buttons change
- Duplicates unit test coverage
- No business value

**What To Do Instead**:
- REMOVE from BDD
- Replace with 1 BDD scenario: "Professional confirms event deletion"
- Create unit tests:
  ```go
  func TestDeleteModalAppears(t *testing.T)  // modal display
  func TestDeleteModalClose(t *testing.T)     // modal closes
  func TestDeleteWithEnter(t *testing.T)      // Enter key
  func TestDeleteWithEscape(t *testing.T)     // Escape key
  func TestDeleteAnimation(t *testing.T)      // animation timing
  ```

---

## Error Handling

### ✅ GOOD: Business Rules

```gherkin
Scenario: System rejects empty event description
  When I attempt to create event with empty description
  Then the creation should fail
  And the system should show validation error
  And the event should not be created

Scenario: System requires company for some event types
  When I create a project-type event
  And I don't specify a company
  Then the save should fail
  And I should see error "Company required for project events"

Scenario: System validates event date is not in future
  When I attempt to create event dated tomorrow
  Then the system should reject it
  And should show error "Event date cannot be in future"
  And the form should retain my entered data
```

**Why It's Good**:
- Tests business rules (validation, constraints)
- Tests error outcomes (what user sees, what happens)
- Survives UI changes (different error display)
- Tests business logic, not UI

### ❌ BAD: UI Error Display

```gherkin
Scenario: Error message appears in red text
  When I submit invalid data
  Then an error message should appear
  And the text should be red
  And the text should use sans-serif font
  And the error box should have border

Scenario: Error message uses specific wording
  When I make this error
  Then I should see exact text "ERROR: Invalid input"
  And the error box should have an icon
  And the icon should be a red X
```

**Why It's Bad**:
- Tests UI styling (red text, font, border)
- Tests UI element styling (icon appearance)
- Tests exact error wording (brittle)
- Tests display mechanics, not business logic
- Breaks when: styling changes, wording changes, icon changes
- No business value
- Belongs in design/UI tests, not BDD

**What To Do Instead**:
- 1 BDD scenario: "System validates event has required fields"
- Unit test: "Error message includes field name"
- UI tests (separate): "Error displays in red"

---

## Summary: The Pattern

| Scenario Type | ✅ Good (BDD) | ❌ Bad (Move to Unit) |
|---|---|---|
| **Business outcome** | "Event is created" | "Modal appears" |
| **Data operations** | "Event metadata updated" | "Form field focused" |
| **Business rules** | "Minimum 2 events in burst" | "Tab moves between fields" |
| **Integration** | "Events + Skills linked" | "Animation completes" |
| **Workflows** | "Professional captures event" | "j key moves selection" |
| **Error handling** | "Validation prevents bad data" | "Error text is red" |

**Quick Decision**:
- If it tests **WHAT** the system does → ✅ BDD
- If it tests **HOW** the user does it → ❌ Unit Test
- If it tests **STYLING/ANIMATION** → ❌ UI/Visual Test

---

## Key Principles

1. **BDD tests outcomes, not mechanics**
2. **UI interactions belong in unit tests**
3. **Each layer tests different concerns**
4. **Good scenarios survive UI changes**
5. **Move 10+ step scenarios to unit tests**
6. **Hide HOW in step definitions**
7. **Test business value, not implementation**

---

**Last Updated**: February 12, 2026
