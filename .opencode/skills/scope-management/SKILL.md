---
name: scope-management
description: Manage scope effectively - say no appropriately, prevent scope creep, break down oversized tasks
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Help manage scope effectively - knowing when to say no, preventing scope creep, recognising when tasks are too big, and keeping work focused on what matters.

## When to use me

- Task feels unbounded or unclear
- Requirements keep expanding
- "While we're at it..." requests appear
- Work is taking longer than expected
- Unsure if something is in scope

## Core Principles

1. **Scope is a feature** - Constraints enable focus and delivery
2. **Say no to say yes** - Declining one thing enables another
3. **Small batches** - Smaller scope = faster feedback = less risk
4. **Explicit boundaries** - Unclear scope expands indefinitely
5. **Done is better than perfect** - Ship, then iterate

## Recognising Scope Problems

### Signs a Task is Too Big

```
RED FLAGS:
- "This will take 2-3 weeks"
- Can't explain it in one sentence
- Multiple unrelated changes bundled
- "And then we'll also need to..."
- No clear definition of done
- Touches many parts of the codebase
- Requirements keep being discovered
```

### Signs of Scope Creep

```
CREEP PATTERNS:
- "While you're in there, could you also..."
- "It would be nice if..."
- "What about edge case X?"
- "Users might also want..."
- Requirements changing mid-implementation
- Gold plating (adding unrequested features)
```

## Saying No Effectively

### The No Framework

```
1. ACKNOWLEDGE: Show you understand the request
2. EXPLAIN: Give a clear reason (not excuse)
3. OFFER ALTERNATIVE: Suggest what you CAN do
4. DOCUMENT: Track for future consideration
```

### Examples

```markdown
## Request: "Add sorting to all list views"

ACKNOWLEDGE: "That would improve usability across the app."

EXPLAIN: "Adding to all views would delay the current feature 
by 2+ weeks and increase testing complexity significantly."

ALTERNATIVE: "I can add sorting to the timeline view now (1 day),
and we can add it to other views in a follow-up task."

DOCUMENT: Created task #123 for "Add sorting to remaining views"
```

```markdown
## Request: "Handle this edge case"

ACKNOWLEDGE: "Good catch - that is a potential issue."

EXPLAIN: "This edge case affects <1% of users and would add
significant complexity to handle properly."

ALTERNATIVE: "I'll add input validation that prevents this case
for now, and we can revisit if users actually hit it."

DOCUMENT: Added to known limitations in docs.
```

### When to Say No

| Situation | Response |
|-----------|----------|
| Unrelated to current task | "Let's track that separately" |
| Nice-to-have, not need-to-have | "Future enhancement" |
| Premature optimisation | "Let's measure first" |
| Edge case unlikely to occur | "Document and monitor" |
| Would delay critical path | "After we ship this" |
| Unclear requirements | "Let's clarify before building" |

## Breaking Down Large Tasks

### The Slicing Technique

```markdown
## Original: "Implement event filtering"

TOO BIG - Slice by:

### Option 1: By Filter Type
1. Date filter only
2. Add type filter
3. Add text search
4. Add combined filters

### Option 2: By Layer
1. Filter UI (hardcoded options)
2. Filter logic (in-memory)
3. Filter persistence
4. Server-side filtering

### Option 3: By User Journey
1. Basic filter (single criterion)
2. Clear filter
3. Multiple criteria
4. Save filter presets

PICK: Option 1 - delivers value incrementally
```

### The INVEST Criteria

Each task should be:

```
I - Independent: Can be done without other tasks
N - Negotiable: Details can be discussed
V - Valuable: Delivers user/business value
E - Estimable: Can reasonably estimate size
S - Small: Fits in a few days max
T - Testable: Clear acceptance criteria
```

### Splitting Patterns

| Pattern | When to Use |
|---------|-------------|
| **By workflow step** | Multi-step processes |
| **By user type** | Different user needs |
| **By data type** | Multiple entity types |
| **By operation** | CRUD operations |
| **By platform** | Cross-platform features |
| **By quality** | Basic → Enhanced → Polished |

## Scope Boundaries

### Defining Done

```markdown
## Task: Add date filter to timeline

### In Scope
- [ ] Date range picker UI
- [ ] Filter events by date range
- [ ] Clear filter button
- [ ] Filter persists during session

### Explicitly Out of Scope
- Server-side filtering (separate task)
- Saved filter presets (future enhancement)
- Filter by multiple criteria (next iteration)
- Custom date formats (use system default)

### Acceptance Criteria
- User can select start and end date
- Only events within range are shown
- Clearing filter shows all events
- Filter state maintained when navigating to detail and back
```

### Scope Change Process

```markdown
## Scope Change Request

**Original scope:** Date filter with single date range
**Requested change:** Add preset ranges (last week, last month)

### Assessment
- Effort: +0.5 days
- Value: High (common use case)
- Risk: Low (additive change)
- Deadline impact: Minor

### Decision
ACCEPT - Low effort, high value, doesn't change core scope

### If Rejected
"Good idea - let's add it as a fast-follow after we ship
the basic filter. Created task #124."
```

## Preventing Scope Creep

### At Task Start

```markdown
## Pre-Work Checklist

- [ ] Scope documented and agreed
- [ ] Out-of-scope explicitly listed
- [ ] Acceptance criteria defined
- [ ] Estimate based on defined scope
- [ ] Stakeholders aligned
```

### During Implementation

```markdown
## When New Requirements Appear

1. PAUSE - Don't just add it
2. ASK - Is this in original scope?
3. ASSESS - What's the impact?
4. DECIDE - Add now, later, or never
5. DOCUMENT - Track the decision

RESPONSE TEMPLATE:
"This wasn't in the original scope. Let me assess the impact
and we can decide whether to include it now or track it
for later."
```

### Scope Creep Responses

| Request | Response |
|---------|----------|
| "Quick addition" | "Let me check the impact first" |
| "Easy change" | "Even easy changes need testing" |
| "Users expect this" | "Do we have evidence?" |
| "Obvious requirement" | "Let's document it explicitly" |
| "Just one more thing" | "Let's finish this first" |

## When Tasks Grow

### Recognising Growth

```
ORIGINAL ESTIMATE: 2 days
CURRENT STATUS: Day 3, ~50% done
CAUSE: Discovered complexity / scope grew

OPTIONS:
1. Continue (accept delay)
2. Cut scope (ship smaller)
3. Split task (ship part now)
4. Get help (pair/mob)
```

### The Re-Scoping Conversation

```markdown
"I'm on day 3 of a 2-day estimate. Here's what I've learned:

DISCOVERED:
- Need to refactor X first (unplanned)
- Edge case Y is more complex than expected
- Dependency on Z wasn't known

OPTIONS:
1. Continue: 3 more days, full scope
2. Cut: 1 more day, without feature A
3. Split: Ship core today, feature A tomorrow

RECOMMENDATION: Option 3 - delivers value sooner,
reduces risk, maintains momentum.

What do you think?"
```

## Managing Yourself

### Avoiding Self-Inflicted Scope Creep

```markdown
WATCH FOR:
- "I'll just clean this up while I'm here"
- "This would be better if..."
- "Let me add tests for this unrelated thing"
- "I should refactor this first"
- Perfectionism

DISCIPLINE:
- Stick to the task
- Note improvements for later
- Boy Scout Rule: Small cleanups only
- If it's bigger than 10 minutes, it's a new task
```

### Timeboxing

```markdown
## Timebox Rules

INVESTIGATION: Max 2 hours, then decide
SPIKE: Max 4 hours, then report
RABBIT HOLE: Max 30 minutes, then surface

When timebox expires:
1. Stop
2. Document what you learned
3. Decide: Continue, pivot, or stop
4. Communicate status
```

## Related Skills

- `estimation` - Sizing work accurately
- `pragmatic-problem-solving` - Practical focus
- `create-task` - Documenting scope
- `checklist-discipline` - Tracking progress
