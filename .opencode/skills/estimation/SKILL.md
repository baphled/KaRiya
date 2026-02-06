---
name: estimation
description: Estimate work effectively - break down tasks, account for uncertainty, communicate ranges
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Guide effective work estimation - breaking down tasks, accounting for uncertainty, learning from past estimates, and communicating appropriately.

## When to use me

- Planning new features
- Estimating bug fixes
- Sprint planning
- Roadmap discussions
- When asked "how long will this take?"

## Core Principles

1. **Estimates are ranges, not points** - Uncertainty is inherent
2. **Break down the work** - Smaller pieces are easier to estimate
3. **Account for the unknown** - Include investigation, testing, review
4. **Learn from history** - Track actual vs estimated
5. **Communicate uncertainty** - Don't false precision

## The Estimation Problem

### Why Estimates Are Hard

```
What we estimate:    Writing the code
What actually happens:
- Understanding requirements (often unclear)
- Investigating existing code
- Discovering edge cases
- Writing tests
- Code review feedback
- Bug fixes
- Documentation
- Deployment issues
```

### Cone of Uncertainty

```
Early project:    0.25x to 4x actual
After design:     0.5x to 2x actual
After detailed design: 0.8x to 1.25x actual
```

**Implication:** Early estimates should be wide ranges.

## Estimation Process

### Step 1: Clarify Requirements

Before estimating, ask:
- What problem are we solving?
- What does "done" look like?
- What are the acceptance criteria?
- What's explicitly out of scope?

```markdown
## Task: Add date filter to timeline

### In Scope
- Date range picker UI
- Filter logic in list view
- Persist filter across navigation

### Out of Scope
- Server-side filtering
- Saved filter presets
- Date validation beyond basic range

### Acceptance Criteria
- [ ] User can select start and end date
- [ ] List shows only events in range
- [ ] Filter persists when returning to list
```

### Step 2: Break Down the Work

```markdown
## Breakdown: Add date filter

### Tasks
1. Add filter modal UI
   - Date picker component
   - Modal layout
   - Apply/cancel buttons
   
2. Implement filter logic
   - Filter function
   - State management
   - Screen integration

3. Add persistence
   - Store filter in intent state
   - Restore on navigation return

4. Testing
   - Unit tests for filter
   - E2E test for workflow

5. Code review + fixes
```

### Step 3: Estimate Each Piece

```markdown
| Task | Optimistic | Likely | Pessimistic |
|------|------------|--------|-------------|
| Filter modal UI | 2h | 4h | 8h |
| Filter logic | 1h | 2h | 4h |
| Persistence | 1h | 2h | 4h |
| Testing | 2h | 4h | 6h |
| Review + fixes | 1h | 2h | 4h |
| **Total** | **7h** | **14h** | **26h** |

Estimate: 1-3 days (likely ~2 days)
```

### Step 4: Add Contingency

```
Base estimate: 2 days

Contingency factors:
- New to date picker library: +20%
- Requirements might change: +20%
- Unknown unknowns: +20%

Adjusted estimate: 2-3 days
```

## Estimation Techniques

### T-Shirt Sizing

Quick relative estimation:

| Size | Meaning | Typical Duration |
|------|---------|------------------|
| XS | Trivial, well-understood | < 2 hours |
| S | Small, clear scope | 2-4 hours |
| M | Medium, some unknowns | 1-2 days |
| L | Large, needs breakdown | 3-5 days |
| XL | Epic, must be broken down | > 1 week |

### Reference-Based Estimation

Compare to similar past work:

```markdown
## Estimate: Add skill filter

### Similar past work
- Date filter: Estimated 2 days, took 3 days
- Type filter: Estimated 1 day, took 1.5 days

### Comparison
- Complexity similar to date filter
- Less UI work (simpler picker)
- More backend (skill relationships)

### Estimate
2-3 days (accounting for past underestimation)
```

### Three-Point Estimation

```
Optimistic (O): Everything goes perfectly
Likely (L): Normal conditions
Pessimistic (P): Significant obstacles

Expected = (O + 4L + P) / 6
Standard Deviation = (P - O) / 6

Example:
O = 2 days, L = 4 days, P = 10 days
Expected = (2 + 16 + 10) / 6 = 4.7 days
SD = (10 - 2) / 6 = 1.3 days

Estimate: 4-6 days (95% confidence: 2-8 days)
```

## Communicating Estimates

### Use Ranges

```
# GOOD - Range communicates uncertainty
"This will take 2-4 days, most likely around 3 days."

# BAD - False precision
"This will take 2.5 days."
```

### Qualify Confidence

```
# High confidence (well-understood, done before)
"2-3 days, high confidence - similar to the filter we built last sprint."

# Medium confidence (some unknowns)
"3-5 days, medium confidence - depends on how complex the API integration is."

# Low confidence (many unknowns)
"1-2 weeks rough estimate - need spike to understand the problem better."
```

### Call Out Assumptions

```
"My estimate assumes:
- Requirements are final
- No major changes to existing code needed
- Code review within 1 day
- No blocking dependencies

If any of these change, estimate will change."
```

## When Estimates Change

### Re-estimate When

- Requirements change
- Discover unexpected complexity
- Dependencies change
- Significant time has passed

### How to Communicate

```
"Original estimate was 3 days. After investigating, I found we need
to refactor the existing filter system first. New estimate is 5-7 days.

Breakdown:
- Refactor existing filters: 2-3 days (unplanned)
- Original work: 3-4 days (unchanged)
"
```

## Learning from Estimates

### Track Actuals

```markdown
## Estimation Log

| Task | Estimated | Actual | Variance | Notes |
|------|-----------|--------|----------|-------|
| Date filter | 2d | 3d | +50% | Requirements changed |
| Skill filter | 2d | 2d | 0% | On target |
| Export feature | 3d | 5d | +67% | Underestimated testing |

### Patterns
- Consistently underestimate testing (+50% on average)
- Requirements changes add ~1 day typically
- Well-understood work estimated accurately
```

### Adjust Future Estimates

Based on patterns:
- Add buffer for testing
- Add buffer for requirements uncertainty
- Keep estimates for familiar work

## Common Estimation Mistakes

| Mistake | Impact | Fix |
|---------|--------|-----|
| Forgetting testing | Underestimate | Include in breakdown |
| Ignoring review | Underestimate | Add review time |
| Optimistic assumptions | Underestimate | Ask "what could go wrong?" |
| Not clarifying scope | Wrong estimate | Clarify before estimating |
| Point estimates | False confidence | Use ranges |
| Anchoring | Biased estimate | Estimate independently first |

## Quick Estimation Checklist

Before giving an estimate:
- [ ] Understood what "done" means?
- [ ] Broken down into tasks?
- [ ] Included testing?
- [ ] Included code review?
- [ ] Accounted for unknowns?
- [ ] Considered similar past work?
- [ ] Expressed as a range?
- [ ] Stated assumptions?

## Related Skills

- `create-task` - Task documentation
- `pragmatic-problem-solving` - Practical approaches
- `critical-thinking` - Analysing requirements
