---
name: question-resolver
description: Systematically resolve questions - determine if answerable, gather evidence, test assumptions, reach correct conclusions
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

When a question arises, I systematically determine whether it can be answered, gather evidence, test assumptions, and reach correct conclusions. I prevent blind assumptions and ensure we actually know what we think we know.

## When to use me

**Automatically triggered when:**
- A question arises during development
- Uncertainty about how something works
- Making decisions based on assumptions
- "I think..." or "probably..." statements appear
- Conflicting information is encountered

## Core Principles

### 1. Questions Have Different Types

| Type | Can Answer? | How |
|------|-------------|-----|
| **Empirical** - "Does X do Y?" | Yes | Test it, read code, run experiment |
| **Historical** - "Why was X built this way?" | Maybe | Git history, docs, comments |
| **Predictive** - "Will X work for Y?" | Test first | Prototype, benchmark, spike |
| **Subjective** - "Is X better than Y?" | Depends | Define criteria, then measure |
| **Undefined** - "Should we do X?" | Need context | Clarify requirements first |

### 2. Assumptions Are Hypotheses

Every assumption is a hypothesis that needs testing:

```
ASSUMPTION: "The form validates email format"

Is this KNOWN or ASSUMED?
├── KNOWN (verified) → Proceed with confidence
└── ASSUMED (not verified) → TEST IT

Test: Write a test, check the code, or try invalid input
Result: Either confirm or correct the assumption
```

### 3. "I Don't Know" Is Valid

It's better to say "I don't know, let me find out" than to guess and be wrong.

## Question Resolution Process

### Step 1: Classify the Question

```
Question: "Does the timeline filter persist across navigation?"

Type: Empirical (can be tested)
Scope: Specific (timeline filter behavior)
Answerable: Yes, by examining code or testing
```

### Step 2: Identify What We Know vs Assume

```
KNOWN (verified):
- Timeline has a filter feature (seen in code)
- Navigation exists between screens (seen in code)

ASSUMED (not verified):
- Filter state is stored somewhere
- State persists across navigation

UNKNOWN:
- WHERE state is stored
- IF it actually persists
```

### Step 3: Gather Evidence

**For code questions:**
```bash
# Search for filter state
grep -r "filter" internal/cli/intents/browsetimeline/

# Check state management
grep -r "FilterState\|filterState" internal/cli/

# Look at navigation handlers
cat internal/cli/intents/browsetimeline/handlers.go | grep -A 20 "Navigate"
```

**For behavior questions:**
```bash
# Write a test
go test -v ./internal/cli/intents/browsetimeline -run "TestFilterPersistence"

# Or run the app and try it
go run ./cmd/kariya
```

### Step 4: Test Assumptions

```go
// Don't assume - verify
It("persists filter across navigation", func() {
    // Arrange
    intent := createTimelineIntent()
    intent.ApplyFilter(myFilter)
    
    // Act - navigate away and back
    intent.NavigateTo("detail")
    intent.NavigateTo("list")
    
    // Assert - filter should still be applied
    Expect(intent.CurrentFilter()).To(Equal(myFilter))
})
```

### Step 5: Conclude with Evidence

```
CONCLUSION: Filter does NOT persist across navigation

EVIDENCE:
1. Code review: FilterState is local to ListScreen, not Intent
2. Test result: Filter resets when returning to list
3. Git history: No persistence was ever implemented

IMPLICATION: If persistence is needed, it must be added
```

## Assumption Testing Patterns

### Pattern 1: Code Verification

```
Assumption: "Service validates input before saving"

Verification:
1. Find the Save method
2. Look for validation call
3. Check if validation errors are returned

Result: [CONFIRMED/REFUTED with evidence]
```

### Pattern 2: Behavioral Testing

```
Assumption: "Escape key always returns to previous screen"

Verification:
1. Write test for each screen
2. Send Escape key message
3. Assert navigation result

Result: [CONFIRMED/REFUTED with evidence]
```

### Pattern 3: Documentation Check

```
Assumption: "API returns paginated results"

Verification:
1. Check API documentation
2. Check actual response structure
3. Test with large dataset

Result: [CONFIRMED/REFUTED with evidence]
```

### Pattern 4: Historical Investigation

```
Assumption: "This was done for performance reasons"

Verification:
1. git log --oneline -p path/to/file
2. git blame path/to/file
3. Check commit messages and PR descriptions
4. Look for benchmarks or performance tests

Result: [CONFIRMED/REFUTED/UNKNOWN with evidence]
```

## Red Flags - When to Stop and Verify

| Red Flag | Action |
|----------|--------|
| "I think..." | Stop. Verify before proceeding |
| "Probably..." | Stop. Test the assumption |
| "Should be..." | Stop. Confirm it actually is |
| "Usually..." | Stop. Check this specific case |
| "I assume..." | Stop. Make it explicit and test |
| Copying code without understanding | Stop. Understand first |
| "It worked before" | Stop. Verify it still works |

## Question Templates

### For Code Behavior
```
QUESTION: Does [component] do [behavior]?

INVESTIGATION:
1. Read: [relevant file/function]
2. Test: [command to run]
3. Evidence: [what I found]

ANSWER: [Yes/No/Partially] because [evidence]
```

### For Design Decisions
```
QUESTION: Should we use [approach A] or [approach B]?

CRITERIA:
1. [criterion 1]: A scores [X], B scores [Y]
2. [criterion 2]: A scores [X], B scores [Y]

CONSTRAINTS:
- [relevant constraint]

RECOMMENDATION: [A/B] because [reasoning based on evidence]
```

### For Unknown Territory
```
QUESTION: How does [unfamiliar thing] work?

EXPLORATION:
1. Read documentation: [link/file]
2. Find examples: [where]
3. Create minimal test: [what I tried]

UNDERSTANDING: [summary of how it works]
CONFIDENCE: [High/Medium/Low]
GAPS: [what I still don't understand]
```

## Integration with Development

### During TDD
```
RED phase:
- Question: "What should this test assert?"
- Resolve: Check requirements, existing patterns, user expectations
- Don't assume the obvious answer is correct

GREEN phase:
- Question: "Is this the minimal implementation?"
- Resolve: Could it be simpler? Am I adding unnecessary complexity?

REFACTOR phase:
- Question: "Is this pattern appropriate here?"
- Resolve: What problem does the pattern solve? Do we have that problem?
```

### During Code Review
```
For each change, ask:
- What assumption does this make?
- Is that assumption verified?
- What happens if the assumption is wrong?
```

### During Debugging
```
1. State the assumption about what's happening
2. Design test to verify/refute assumption
3. Run test
4. Update understanding based on evidence
5. Repeat until root cause found
```

## Anti-Patterns

| Anti-Pattern | Problem | Correct Approach |
|--------------|---------|------------------|
| Assuming without testing | May be wrong | Test first |
| Asking without investigating | Wastes time | Try to find out first |
| Accepting first answer | May be incomplete | Verify independently |
| Ignoring contradicting evidence | Confirmation bias | Update beliefs with evidence |
| "It's obvious" | Often wrong | Make it explicit, test it |
| Copy-paste solution | May not fit | Understand, then adapt |

## Related Skills

- `critical-thinking` - Rigorous analysis
- `research` - Systematic investigation
- `systems-thinker` - Understanding interconnections
- `pragmatic-problem-solving` - Practical solutions
- `assumption-tracker` - Explicit assumption management
