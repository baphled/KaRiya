# Token Efficiency Skill

You are an expert in optimizing AI interactions for token efficiency - maximizing value while minimizing token usage.

## Overview

Tokens cost money and time. Be concise, precise, and strategic in communications to maximize productivity per token.

---

## Core Principles

### 1. Be Concise

```
❌ WASTEFUL:
"I think that perhaps it would be a good idea if we were to consider 
the possibility of maybe implementing a function that could potentially 
handle the validation of user input in some way."

✅ EFFICIENT:
"Implement input validation function."
```

### 2. Be Specific

```
❌ VAGUE:
"Fix the bug in the code"

✅ SPECIFIC:
"Fix nil pointer in EventService.Create when repository is nil"
```

### 3. Be Structured

```
❌ WALL OF TEXT:
"I need you to create a new feature that allows users to filter 
events by date and also by company and we should probably add 
sorting too and make sure it works with the existing pagination..."

✅ STRUCTURED:
Create event filtering:
1. Filter by date range
2. Filter by company
3. Add sorting options
4. Integrate with pagination
```

---

## Communication Patterns

### Task Requests

```
❌ INEFFICIENT:
"Hey, so I was thinking that maybe we could add some functionality 
to our application. What I'm envisioning is something like a feature 
where users can see their events in a timeline format. I'm not sure 
exactly how it should work but maybe you can help me figure that out?"

✅ EFFICIENT:
Add timeline view:
- Display events chronologically
- Show date, title, company
- Navigate with j/k keys
- Press Enter for details
```

### Bug Reports

```
❌ INEFFICIENT:
"Something seems to be wrong with the application. When I try to do 
certain things it doesn't work properly. I'm not exactly sure what 
the issue is but it's definitely not behaving as expected."

✅ EFFICIENT:
Bug: Timeline crashes on empty events
- Steps: Open timeline with no events
- Expected: Show "No events" message  
- Actual: Panic at timeline.go:42
```

### Code Requests

```
❌ INEFFICIENT:
"Could you please help me write some code? I need a function that 
will take some input and process it somehow. The function should 
be able to handle various edge cases and return appropriate results."

✅ EFFICIENT:
Write ValidateEvent(e *Event) error:
- Check text non-empty
- Check date not future
- Return descriptive errors
```

---

## Response Optimization

### When Asking Questions

```
❌ MULTIPLE QUESTIONS:
"What does this code do? And why is it structured this way? Also, 
is this the best approach? What would you recommend instead?"

✅ ONE FOCUSED QUESTION:
"Why use goroutines here instead of sequential processing?"
```

### When Providing Context

```
❌ TOO MUCH CONTEXT:
[Entire 500-line file]
"What's wrong with this code?"

✅ RELEVANT CONTEXT:
```go
// Lines 42-50 of service.go
func (s *Service) Process() error {
    result := s.repo.Find()  // Returns nil, not error
    return result.Validate() // Panic here
}
```
"Fix nil check at line 48"
```

### When Reviewing Output

```
❌ VAGUE FEEDBACK:
"That's not quite right, can you try again?"

✅ SPECIFIC FEEDBACK:
"Wrong: used string. Need: EventStatus type. Update line 15."
```

---

## Code Comments (In Prompts)

### Minimal Context

```
❌ WASTEFUL:
```go
// This is the main function that handles the event processing
// It takes an event as input and validates it according to our
// business rules. If validation passes, it saves the event to
// the database. If validation fails, it returns an error.
func ProcessEvent(e *Event) error {
```

✅ EFFICIENT:
```go
// Save validated event to DB
func ProcessEvent(e *Event) error {
```
```

### Reference by Line

```
✅ EFFICIENT:
"service.go:42 - add nil check before calling Validate()"

Instead of copying entire function
```

---

## Iterative Refinement

### Start Small

```
Round 1: "Add date filter to events"
Round 2: "Also filter by company" (only if needed)
Round 3: "Add sorting" (only if needed)

Better than:
"Add filtering by date, company, project, status, 
sorting by all fields, pagination, search..."
```

### Verify Before Expanding

```
1. Request small change
2. Verify it works
3. Request next change

Not: Request everything at once, debug massive changes
```

---

## File Operations

### Reading Files

```
❌ INEFFICIENT:
"Show me the entire contents of service.go"

✅ EFFICIENT:
"Show lines 40-60 of service.go" (specific section)
"Find ProcessEvent function in service.go" (targeted search)
```

### Writing Files

```
❌ INEFFICIENT:
[Rewrite entire 200-line file for 1-line change]

✅ EFFICIENT:
"Line 45: change `return nil` to `return ErrNotFound`"
```

---

## Search Optimization

### Targeted Searches

```
❌ INEFFICIENT:
"Find all error handling in the codebase"

✅ EFFICIENT:
"Find ErrNotFound in repository package"
```

### Specific Patterns

```
❌ BROAD:
"Show me all tests"

✅ SPECIFIC:
"Show TestEventService_Create in service_test.go"
```

---

## Token-Saving Templates

### Bug Fix Request

```
Fix: [file:line] [issue]
Cause: [reason]
Expected: [behavior]
```

### Feature Request

```
Add [feature]:
1. [Requirement 1]
2. [Requirement 2]
Location: [where to add]
```

### Code Review

```
Review [file]:
Focus: [specific aspect]
```

### Question

```
Q: [Specific question]
Context: [Minimal relevant context]
```

---

## Anti-Patterns

### DON'T: Repeat Yourself

```
❌ "Can you help me? I need help with something. 
I'm having trouble and need assistance..."

✅ "Help: EventService nil pointer at line 42"
```

### DON'T: Over-Explain

```
❌ "I want to add a feature. The feature should allow users 
to do X. When users do X, the system should respond with Y.
The reason we need this feature is because..."

✅ "Add feature: X → returns Y"
```

### DON'T: Include Irrelevant Code

```
❌ [300 lines of unrelated code]
"Fix the bug"

✅ [10 relevant lines]
"Fix nil check at line 5"
```

### DON'T: Ask for Explanations You Don't Need

```
❌ "Explain what this does and then fix it"

✅ "Fix the nil pointer" (explanation if needed later)
```

---

## Efficiency Checklist

Before sending:
- [ ] Is every word necessary?
- [ ] Is the request specific?
- [ ] Is context minimal but sufficient?
- [ ] Could this be split into smaller requests?
- [ ] Am I asking for only what I need now?

---

## Metrics

Track your efficiency:
- Lines of response per task completed
- Tasks completed per session
- Rework rate (re-requests for same task)

---

## Related Skills

- `clean-code` - Concise code, concise prompts
- `research` - Efficient information gathering
- `critical-thinking` - Ask the right questions
