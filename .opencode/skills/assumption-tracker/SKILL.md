---
name: assumption-tracker
description: Explicitly track, test, and validate assumptions - prevent blind spots and reduce risk from untested beliefs
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Track assumptions explicitly, ensure they get tested, and update understanding when assumptions prove wrong. I prevent the silent accumulation of untested beliefs that lead to bugs and wrong decisions.

## When to use me

**Always.** Assumptions are made constantly - this skill ensures they're visible and validated.

## Core Principle

**Every assumption is a risk until verified.**

Untracked assumptions become:
- Bugs when code assumes wrong things
- Wasted work when requirements were misunderstood
- Technical debt when design assumptions change
- Security vulnerabilities when trust assumptions fail

## Assumption Categories

### 1. Technical Assumptions

```
"The database connection will always be available"
"This function never returns nil"
"The API response format won't change"
"This operation is fast enough"
```

**Test by:** Code inspection, unit tests, integration tests, load tests

### 2. Requirements Assumptions

```
"Users want to filter by date"
"The list won't have more than 100 items"
"This field is always required"
"Users understand what 'burst' means"
```

**Test by:** Clarify with stakeholder, user testing, analytics

### 3. Environmental Assumptions

```
"We'll always run on Linux"
"Terminal supports 256 colors"
"Network latency is low"
"Disk space is sufficient"
```

**Test by:** Test in different environments, add runtime checks

### 4. Design Assumptions

```
"This pattern will scale"
"Components are loosely coupled"
"State management is simple enough"
"We won't need to undo this"
```

**Test by:** Architecture review, spike, prototype

## Assumption Tracking Format

### When Making an Assumption

```
ASSUMPTION: [What we're assuming]
CATEGORY: [Technical/Requirements/Environmental/Design]
CONFIDENCE: [High/Medium/Low]
RISK IF WRONG: [Impact description]
VALIDATION: [How we will/did verify]
STATUS: [Untested/Testing/Verified/Refuted]
```

### Example

```
ASSUMPTION: Form validation runs before submit handler
CATEGORY: Technical
CONFIDENCE: Medium (haven't verified)
RISK IF WRONG: Invalid data saved to database
VALIDATION: Read form code, write test
STATUS: Untested → Testing

[After investigation]

STATUS: Verified
EVIDENCE: forms/validation.go:45 - Validate() called in Submit()
```

## Tracking During Development

### At Task Start

```
## Assumptions for: Add date filter to timeline

ASSUMPTION 1: Events have a date field
- Confidence: High (seen in code)
- Status: Verified (domain/event.go:12)

ASSUMPTION 2: Date filtering is client-side
- Confidence: Medium (not sure)
- Status: Untested
- Action: Check if server supports date query

ASSUMPTION 3: Users want to filter by range, not single date
- Confidence: Low (guessing)
- Status: Untested  
- Action: Clarify with requirements
```

### During Implementation

As you code, notice and record assumptions:

```go
func (s *Screen) filterByDate(events []*Event, date time.Time) []*Event {
    // ASSUMPTION: events slice is never nil
    // ASSUMPTION: Event.Date is always populated
    // ASSUMPTION: Date comparison uses local timezone
    
    var filtered []*Event
    for _, e := range events {
        if e.Date.Equal(date) {
            filtered = append(filtered, e)
        }
    }
    return filtered
}
```

Convert to defensive code or tests:

```go
func (s *Screen) filterByDate(events []*Event, date time.Time) []*Event {
    if events == nil {
        return nil  // Handle nil case explicitly
    }
    
    var filtered []*Event
    for _, e := range events {
        if e.Date.IsZero() {
            continue  // Skip events without dates
        }
        if e.Date.Truncate(24*time.Hour).Equal(date.Truncate(24*time.Hour)) {
            filtered = append(filtered, e)
        }
    }
    return filtered
}
```

### At Task End

Review and close out assumptions:

```
## Assumption Review: Add date filter to timeline

ASSUMPTION 1: Events have a date field
- Status: VERIFIED

ASSUMPTION 2: Date filtering is client-side
- Status: VERIFIED (server has no date param)
- Note: May want server-side for large datasets

ASSUMPTION 3: Users want range, not single date
- Status: CHANGED
- Resolution: Requirements confirmed single date for v1, range for v2
- Action: Created task for v2 range filter
```

## Testing Assumptions

### Converting Assumptions to Tests

```
ASSUMPTION: "Empty list shows helpful message"

TEST:
It("shows empty state message when no events", func() {
    screen := NewListScreen(theme, []*Event{})  // empty
    view := screen.View()
    Expect(view).To(ContainSubstring("No events"))
})
```

### Asserting Assumptions in Code

```go
// For critical assumptions, make them explicit
func processEvent(event *Event) error {
    // Assert our assumption
    if event == nil {
        return errors.New("processEvent: event cannot be nil")
    }
    if event.ID == "" {
        return errors.New("processEvent: event must have ID")
    }
    
    // Now safe to proceed
    // ...
}
```

### Documenting Assumed Preconditions

```go
// ProcessEvents processes a batch of events.
//
// Expected:
//   - events: non-nil slice (may be empty)
//   - each event has valid ID and Date
//
// Assumes:
//   - Events are already validated by caller
//   - Caller handles returned errors
//
// Returns: processed count, or error if processing fails
// Side effects: Updates database via repository
func (s *Service) ProcessEvents(events []*Event) (int, error) {
```

## Assumption Red Flags

| Red Flag | Hidden Assumption | Action |
|----------|-------------------|--------|
| No nil checks | "This is never nil" | Verify or add check |
| No error handling | "This never fails" | Check error paths |
| Magic numbers | "This value is correct" | Make constant, document why |
| Empty catch blocks | "Errors don't matter" | Handle or log errors |
| "Obviously..." | "Everyone knows this" | Make explicit, verify |
| No validation | "Input is always valid" | Add validation |
| Hardcoded config | "This won't change" | Extract to config |

## Communicating Assumptions

### In Code Reviews

```
PR Comment:
"This assumes the service is always available. What happens if it's down?
Should we add retry logic or graceful degradation?"
```

### In Discussions

```
"I'm assuming we want server-side filtering. Is that correct, or should 
this be client-side? The trade-off is latency vs data transfer."
```

### In Documentation

```
## Design Decisions

### Client-side Filtering
We chose client-side filtering assuming:
- Event lists are small (<1000 items)
- Network latency is more costly than processing
- Users filter frequently

If these assumptions change, consider server-side filtering.
```

## Assumption Debt

Like technical debt, assumption debt accumulates:

```
HIGH DEBT (address soon):
- Untested security assumptions
- Unverified performance assumptions
- Requirements assumptions from months ago

MEDIUM DEBT (track):
- Assumptions about "typical" usage
- Assumptions about data quality
- Assumptions about third-party behavior

LOW DEBT (acceptable):
- Well-tested technical assumptions
- Documented design trade-offs
- Explicitly accepted risks
```

## Related Skills

- `question-resolver` - Answering questions systematically
- `pragmatic-problem-solving` - Practical decision making
- `critical-thinking` - Rigorous analysis
- `research` - Investigation techniques
- `systems-thinker` - Understanding dependencies
