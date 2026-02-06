---
name: critical-thinking
description: Apply rigorous analysis to evaluate solutions, challenge assumptions, and make sound technical decisions
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Apply rigorous analysis to evaluate solutions, challenge assumptions, and make sound technical decisions.

## When to use me

Use this skill when:
- Evaluating proposed solutions
- Making architectural decisions
- Reviewing code or designs
- Debugging complex issues
- Assessing trade-offs
- Challenging assumptions

## Core Principles

1. **Question everything** - Assumptions are hypotheses until verified
2. **Evidence over intuition** - Data beats opinions
3. **Consider alternatives** - First idea isn't always best
4. **Think second-order** - What are the downstream effects?
5. **Steel man** - Understand the best version of opposing views

## Critical Thinking Framework

### 1. Clarify the Problem

Before solving, ensure you understand:

```markdown
## Problem Analysis

### What is actually being asked?
[Restate in your own words]

### What are the constraints?
- Must: [Non-negotiable requirements]
- Should: [Important but flexible]
- Could: [Nice to have]

### What does success look like?
[Measurable outcomes]

### What are we NOT solving?
[Explicit scope boundaries]
```

### 2. Challenge Assumptions

Every solution rests on assumptions. Identify and verify them:

```markdown
## Assumption Analysis

| Assumption | Evidence For | Evidence Against | Confidence |
|------------|--------------|------------------|------------|
| Users need X | User research | None | High |
| Y is fast enough | Benchmarks | None | Medium |
| Z won't change | No evidence | Business volatility | Low |
```

**Questions to ask:**
- Why do we believe this?
- What would prove this wrong?
- What if the opposite were true?
- Is this still true in edge cases?

### 3. Evaluate Alternatives

Never accept the first solution. Generate alternatives:

```markdown
## Option Analysis

### Option A: [Name]
**Approach:** Brief description
**Pros:**
- Pro 1
- Pro 2
**Cons:**
- Con 1
- Con 2
**Risks:**
- Risk 1

### Option B: [Name]
...

### Option C: Do Nothing
**Approach:** Accept current state
**Pros:**
- No development cost
- No risk of regression
**Cons:**
- Problem persists
```

### 4. Consider Trade-offs

Every decision has trade-offs. Make them explicit:

```markdown
## Trade-off Analysis

| Factor | Option A | Option B |
|--------|----------|----------|
| Complexity | Higher | Lower |
| Performance | Better | Worse |
| Maintainability | Worse | Better |
| Time to implement | Longer | Shorter |
| Risk | Lower | Higher |
```

**Key trade-off dimensions:**
- Speed vs. Quality
- Flexibility vs. Simplicity
- Short-term vs. Long-term
- Effort vs. Payoff

### 5. Think Second-Order

First-order: What happens immediately?
Second-order: What happens because of that?

```markdown
## Second-Order Effects

### Decision: Add caching layer

**First-order effects:**
- Faster response times
- Reduced database load

**Second-order effects:**
- Cache invalidation complexity
- Stale data possibilities
- More infrastructure to maintain
- Debugging becomes harder

**Third-order effects:**
- Team needs caching expertise
- Monitoring requirements increase
```

### 6. Apply Inversion

Instead of "How do I succeed?", ask "How would I fail?"

```markdown
## Inversion Analysis

### Goal: Build reliable event filtering

### How would we FAIL?
- Don't handle empty results
- Ignore invalid date ranges
- Skip error handling
- Don't test edge cases
- Assume all data is clean

### Therefore, we MUST:
- Handle empty results gracefully
- Validate date ranges
- Comprehensive error handling
- Test all edge cases
- Sanitize input data
```

## Decision-Making Tools

### The 5 Whys

Dig to root cause:

```
Problem: Tests are flaky

Why? They depend on timing
Why? They use real delays
Why? No mock for the timer
Why? Original author didn't know about mocks
Why? No documentation on testing patterns

Root cause: Missing testing documentation
Solution: Document testing patterns, add timer mock
```

### Pre-Mortem

Imagine the project failed. Why?

```markdown
## Pre-Mortem: Feature X

It's 3 months from now. Feature X was a disaster. Why?

### Possible Failure Modes
1. Performance was unacceptable because we didn't benchmark
2. Users hated the UX because we didn't test with them
3. It broke existing features because we skipped regression tests
4. Maintenance became nightmare because we cut corners

### Preventive Actions
1. Set performance budget, benchmark early
2. User testing before full implementation
3. Comprehensive test coverage required
4. No code review shortcuts
```

### Decision Matrix

Score options objectively:

```markdown
## Decision Matrix: State Management Approach

| Criteria (Weight) | Option A | Option B | Option C |
|-------------------|----------|----------|----------|
| Simplicity (3) | 4 (12) | 5 (15) | 2 (6) |
| Performance (2) | 5 (10) | 3 (6) | 5 (10) |
| Team familiarity (2) | 3 (6) | 5 (10) | 2 (4) |
| Maintainability (3) | 4 (12) | 4 (12) | 3 (9) |
| **Total** | **40** | **43** | **29** |

Winner: Option B
```

## Cognitive Biases to Guard Against

| Bias | Description | Counter |
|------|-------------|---------|
| **Confirmation** | Seeking evidence that confirms beliefs | Actively seek disconfirming evidence |
| **Anchoring** | Over-relying on first information | Consider multiple starting points |
| **Sunk Cost** | Continuing because of past investment | Evaluate based on future value only |
| **Availability** | Overweighting recent/memorable events | Look at base rates and data |
| **Bandwagon** | Following what others do | Evaluate on merits |
| **Not Invented Here** | Rejecting external solutions | Judge solutions on quality, not origin |
| **Optimism** | Underestimating difficulty | Use reference class forecasting |

## Red Flags in Reasoning

Watch for these in yourself and others:

- "We've always done it this way"
- "Everyone knows that..."
- "It's obvious that..."
- "We don't have time to consider alternatives"
- "This is the only way"
- "Trust me"
- Emotional attachment to a solution
- Dismissing concerns without addressing them

## Questions for Critical Evaluation

### For Solutions
- What problem does this actually solve?
- What are we assuming?
- What could go wrong?
- What are the alternatives?
- What are the trade-offs?
- How will we know if it works?

### For Code
- Why is this the right abstraction?
- What happens when requirements change?
- What are the edge cases?
- How will this fail?
- Is this the simplest solution?

### For Decisions
- What would change our mind?
- What are we optimizing for?
- What are we giving up?
- Who disagrees and why?
- What's the reversibility?

## Related skills

- `research` - Gather evidence for analysis
- `code-reviewer` - Apply critical thinking to code
- `architecture` - Architectural decisions
- `tech-debt` - Evaluate technical trade-offs
- `software-engineer` - Integrate into overall approach
