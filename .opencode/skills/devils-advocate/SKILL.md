---
name: devils-advocate
description: Challenge ideas, find weaknesses, and stress-test solutions before implementation
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Challenge ideas, find weaknesses, and stress-test solutions before implementation to avoid costly mistakes.

## When to use me

Use this skill when:
- Evaluating a proposed solution
- Before committing to an approach
- When something seems too easy
- When everyone agrees (suspicious!)
- Before major architectural decisions
- Reviewing your own ideas

## Devil's Advocate Principles

1. **Attack ideas, not people** - The goal is better solutions
2. **Assume failure** - What would make this fail?
3. **Seek disconfirming evidence** - Look for what's wrong
4. **Be specific** - Vague concerns don't help
5. **Offer alternatives** - Don't just criticize

## Challenge Framework

### 1. Question the Problem

Before solving, question whether we're solving the right thing:

- Is this actually a problem?
- Whose problem is it?
- How do we know it's a problem?
- What if we did nothing?
- Are we solving the symptom or the cause?

### 2. Challenge Assumptions

Every solution rests on assumptions. Attack them:

```markdown
## Assumption Attack

| Assumption | Challenge | Evidence Needed |
|------------|-----------|-----------------|
| "Users want filters" | Do they? How many asked? | User research data |
| "It needs to be fast" | How fast? What's acceptable? | Performance requirements |
| "We need a database" | Do we? Could we use files? | Data requirements |
| "This is the only way" | Is it? What else could work? | Alternative analysis |
```

### 3. Find the Failure Modes

For every solution, ask "How will this fail?":

```markdown
## Failure Mode Analysis

### Solution: Add caching layer

| Failure Mode | Likelihood | Impact | Mitigation |
|--------------|------------|--------|------------|
| Cache gets stale | High | Medium | TTL + invalidation |
| Cache miss storm | Medium | High | Warm-up strategy |
| Memory exhaustion | Low | High | Size limits |
| Debugging harder | High | Medium | Cache bypass flag |
| Inconsistent state | Medium | High | Single source of truth |
```

### 4. Play the Skeptic

Ask uncomfortable questions:

**On Complexity:**
- "Do we really need this abstraction?"
- "Could we do this more simply?"
- "What are we optimizing for?"

**On Timing:**
- "Why now? What's the urgency?"
- "What if we waited?"
- "Are we ready for this?"

**On Scope:**
- "Is this the minimum viable solution?"
- "What can we cut?"
- "Are we gold-plating?"

**On Confidence:**
- "How do we know this will work?"
- "What's our evidence?"
- "What if we're wrong?"

### 5. Stress Test Edge Cases

Push the solution to its limits:

```markdown
## Edge Case Stress Test

### Feature: Event List with Filters

| Edge Case | Expected Behavior | Potential Problem |
|-----------|-------------------|-------------------|
| 0 events | Empty state message | nil pointer? |
| 10,000 events | Paginated, responsive | Memory? Performance? |
| Filter returns 0 | "No matches" message | User confusion? |
| Invalid date range | Validation error | Edge of time range? |
| Concurrent updates | Consistent view | Race condition? |
| Network failure | Graceful degradation | Error shown? |
```

### 6. Consider Second-Order Effects

What happens after the immediate effect?

```markdown
## Second-Order Analysis

### Decision: Add user preferences table

**First Order:**
- Users can save preferences
- More personalized experience

**Second Order:**
- Need migration strategy
- Need preference sync logic
- Need default handling
- Increases database size

**Third Order:**
- Preference conflicts across devices?
- Performance impact of preference checks?
- Privacy implications?
```

## Challenging Questions Toolkit

### On Requirements
- "Who asked for this?"
- "What happens if we don't do it?"
- "Is this a 'nice to have' or 'must have'?"

### On Design
- "Why this approach over alternatives?"
- "What are we trading off?"
- "Where's the complexity hiding?"

### On Implementation
- "How will this fail?"
- "What's the worst case?"
- "How will we debug this?"

### On Testing
- "How will we know it works?"
- "What tests are missing?"
- "Can this be gamed?"

### On Maintenance
- "Who will maintain this?"
- "How will we know it's broken?"
- "What happens when X changes?"

## Red Flags to Challenge

| When You Hear | Ask |
|---------------|-----|
| "It's obvious that..." | "What makes it obvious?" |
| "Everyone does it this way" | "Is that the right reason?" |
| "We don't have time to..." | "What's the cost of not doing it?" |
| "It should just work" | "How do we verify that?" |
| "We can fix it later" | "Will we? What's the real cost?" |
| "It's only temporary" | "When does it become permanent?" |
| "Trust me" | "Show me the evidence" |

## Constructive Challenging

Don't just criticize - contribute:

```markdown
## Challenge: Solution X

### Concern
I'm worried about [specific issue] because [evidence/reasoning].

### Questions
1. How would we handle [edge case]?
2. What happens when [failure mode]?

### Alternative Consideration
Have we considered [alternative approach]? It might address [concern] by [mechanism].

### Suggestion
If we proceed with X, we should at least [mitigation] to reduce the risk of [failure mode].
```

## Self-Devil's Advocacy

Apply this to your own ideas:

1. **Step back** - Take a break, come back fresh
2. **Assume it's wrong** - Look for problems, not confirmation
3. **Ask a colleague** - "Tell me what's wrong with this"
4. **Write the post-mortem** - Imagine it failed, why?
5. **List the risks** - Be honest about what could go wrong

## Related skills

- `critical-thinking` - Rigorous analysis
- `research` - Gather evidence
- `systems-thinker` - See broader impacts
- `code-reviewer` - Apply to code review
