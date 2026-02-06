---
name: pragmatic-problem-solving
description: Focus on practical solutions - balance ideal with achievable, validate approaches early, ship working software
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Guide practical problem-solving that balances ideal solutions with achievable outcomes. I prevent analysis paralysis, over-engineering, and ivory tower thinking while ensuring we actually solve the real problem.

## When to use me

- Facing a complex problem with multiple possible solutions
- Risk of over-engineering or premature optimization
- Need to make progress despite uncertainty
- Balancing "right way" vs "ship it"
- Time constraints require practical trade-offs

## Core Principles

### 1. Solve the Actual Problem

```
WRONG: "We need a generic event system that handles all possible future requirements"
RIGHT: "We need to notify users when their data is saved"

Ask: What is the ACTUAL problem we're solving TODAY?
```

### 2. Validate Early, Validate Often

```
Don't: Spend 2 days designing the perfect solution
Do: Spend 2 hours on a spike to validate the approach works

Validation beats speculation.
```

### 3. Good Enough Today > Perfect Never

```
Progression:
1. Make it work (correct)
2. Make it right (clean)
3. Make it fast (optimized) - only if needed

Ship working software. Iterate.
```

### 4. Reversibility Matters

```
Reversible decisions: Decide quickly, change if wrong
- Function names
- Internal implementations  
- Local refactoring

Irreversible decisions: Decide carefully
- Public APIs
- Database schemas (harder to change)
- External contracts
```

## Problem-Solving Framework

### Step 1: Define the Real Problem

```
PROBLEM STATEMENT:

What: [specific thing that's wrong or needed]
Who: [who is affected]
Impact: [what happens if unsolved]
Constraint: [time/resources/technical limits]

NOT the problem:
- Symptoms (those point to the problem)
- Solutions (those solve the problem)
```

### Step 2: Explore the Solution Space (Timeboxed)

```
TIMEBOX: [30 min / 1 hour / half day]

Options identified:
1. [Option A] - [quick summary]
2. [Option B] - [quick summary]
3. [Option C] - [quick summary]

Quick assessment:
- Feasibility: Can we actually do this?
- Fit: Does it solve the real problem?
- Effort: How much work?
```

### Step 3: Validate Before Committing

```
VALIDATION APPROACH:

For technical uncertainty:
→ Spike: Build minimal proof-of-concept
→ Timebox: [X hours] max
→ Goal: Prove approach works, not build feature

For requirements uncertainty:
→ Ask: Clarify with stakeholder
→ Prototype: Show something tangible
→ Iterate: Refine based on feedback
```

### Step 4: Implement Incrementally

```
INCREMENTAL DELIVERY:

Slice 1: [minimal useful thing] - delivers [value]
Slice 2: [next enhancement] - delivers [additional value]
Slice 3: [polish/optimization] - delivers [improvement]

Each slice is:
- Complete (works end-to-end)
- Tested (covered by tests)
- Shippable (could deploy if needed)
```

### Step 5: Reflect and Adjust

```
AFTER IMPLEMENTATION:

What worked:
- [thing that went well]

What didn't:
- [thing that was harder than expected]

What we learned:
- [insight for next time]

Adjust approach for next problem.
```

## Decision-Making Heuristics

### When to Choose Simple Over Clever

```
Choose SIMPLE when:
- Team needs to maintain it
- Requirements might change
- You're not sure clever is needed
- Time is limited

Choose CLEVER when:
- Performance is proven bottleneck (measured!)
- Complexity is essential, not accidental
- You'll maintain it yourself forever (rarely true)
```

### When to Spike vs Design

```
SPIKE when:
- Technology is unfamiliar
- Approach is uncertain
- Risk is high
- You're guessing

DESIGN when:
- Problem is well-understood
- Similar solutions exist
- Team has experience
- Requirements are stable
```

### When to Ask vs Figure Out

```
ASK when:
- It's a business/requirements question
- Someone else is the expert
- It would take hours to figure out, minutes to ask
- Getting it wrong has high cost

FIGURE OUT when:
- It's a technical question you should know
- Documentation exists
- Code tells the answer
- Learning is valuable
```

## Pragmatic Trade-offs

### Speed vs Quality

```
Context determines balance:

Prototype/Spike: Speed >> Quality
- Goal is learning, not shipping
- Will be thrown away

Production Code: Quality >> Speed  
- Will be maintained
- Bugs have real cost

Emergency Fix: Speed > Quality (temporarily)
- Fix the bleeding
- Clean up immediately after
```

### Generalization vs Specificity

```
Rule of Three:
- First time: Just do it (specific)
- Second time: Note the duplication
- Third time: NOW generalize

Premature generalization creates:
- Unused flexibility
- Unnecessary complexity
- Harder maintenance
```

### Now vs Later

```
Do NOW:
- Core functionality
- Critical bugs
- Security issues
- Things blocking others

Do LATER:
- Optimizations (unless proven needed)
- Nice-to-haves
- Speculative features
- Non-critical refactoring

Track LATER items so they don't get lost.
```

## Common Traps

| Trap | Symptom | Escape |
|------|---------|--------|
| Analysis Paralysis | Endless discussion, no code | Timebox, then decide |
| Gold Plating | Adding unrequested features | Stick to requirements |
| Premature Optimization | Optimizing unproven bottlenecks | Measure first |
| Not Invented Here | Rebuilding existing solutions | Use libraries |
| Cargo Culting | Copying without understanding | Understand first |
| Perfect is Enemy of Good | Never shipping | Define "good enough" |

## Pragmatic Questions

Ask these to stay grounded:

1. **What problem are we actually solving?**
2. **Who needs this and when?**
3. **What's the simplest thing that could work?**
4. **How will we know if it's working?**
5. **What's the cost of being wrong?**
6. **Can we validate this before building it all?**
7. **What would we cut if we had half the time?**

## Integration with TDD

```
Pragmatic TDD:

RED: Write test for ACTUAL requirement (not imagined future need)
GREEN: Write MINIMAL code (resist urge to generalize)
REFACTOR: Clean up ONLY what you touched (Boy Scout, not renovation)

Don't:
- Test hypothetical scenarios
- Build for requirements you don't have
- Refactor unrelated code
```

## Related Skills

- `question-resolver` - Systematic question resolution
- `critical-thinking` - Rigorous analysis
- `systems-thinker` - Understanding impact
- `assumption-tracker` - Managing assumptions
- `tdd-workflow` - Incremental development
