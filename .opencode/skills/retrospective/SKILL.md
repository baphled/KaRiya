---
name: retrospective
description: Learning from failures and successes, post-mortems, continuous improvement
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
  always_active: false
---

## What I do

Guide retrospective practices for learning from failures and successes, conducting effective post-mortems, and driving continuous improvement.

## When to use me

Use this skill when:
- After an incident or outage
- At the end of a feature/project
- After a significant bug is discovered
- Periodically for team improvement
- When patterns of problems emerge

## Blameless Post-Mortem

### Core Principles

1. **Focus on systems, not people** - "The process allowed this" not "Person X caused this"
2. **Assume good intentions** - Everyone was trying to do the right thing
3. **Seek understanding** - Why did the action seem reasonable at the time?
4. **Improve the system** - How do we make the right thing easier?

### Post-Mortem Template

```markdown
# Post-Mortem: [Incident Title]

**Date:** YYYY-MM-DD
**Author:** [Name]
**Status:** Draft | Final

## Executive Summary

[2-3 sentences: What happened, impact, resolution]

## Timeline

All times in UTC.

| Time | Event |
|------|-------|
| HH:MM | [First sign of problem] |
| HH:MM | [Alert triggered / issue reported] |
| HH:MM | [Investigation began] |
| HH:MM | [Root cause identified] |
| HH:MM | [Fix deployed] |
| HH:MM | [Issue resolved] |

## Impact

- **Duration:** X hours Y minutes
- **Users affected:** [number/percentage]
- **Features affected:** [list]
- **Data loss:** None / [description]
- **Revenue impact:** [if applicable]

## Root Cause Analysis

### What Happened

[Detailed technical explanation]

### Why It Happened

Using the "5 Whys" technique:

1. **Why** did [symptom] occur?
   - Because [cause 1]

2. **Why** did [cause 1] happen?
   - Because [cause 2]

3. **Why** did [cause 2] happen?
   - Because [cause 3]

4. **Why** did [cause 3] happen?
   - Because [cause 4]

5. **Why** did [cause 4] happen?
   - Because [root cause]

### Contributing Factors

- [Factor 1] - [How it contributed]
- [Factor 2] - [How it contributed]

## What Went Well

- [Positive 1] - [Why it helped]
- [Positive 2] - [Why it helped]

## What Could Have Gone Better

- [Improvement 1] - [What we'd do differently]
- [Improvement 2] - [What we'd do differently]

## Action Items

| ID | Action | Owner | Priority | Due Date | Status |
|----|--------|-------|----------|----------|--------|
| 1 | [Preventive action] | @name | High | YYYY-MM-DD | Open |
| 2 | [Detective action] | @name | Medium | YYYY-MM-DD | Open |
| 3 | [Process improvement] | @name | Low | YYYY-MM-DD | Open |

## Lessons Learned

### Technical

- [Lesson with supporting evidence]

### Process

- [Lesson with supporting evidence]

### Communication

- [Lesson with supporting evidence]

## Appendix

### Relevant Logs

```
[Relevant log excerpts]
```

### Metrics/Graphs

[Screenshots or links to relevant dashboards]

### Related Incidents

- [Link to similar past incidents]
```

## Sprint/Project Retrospective

### Format: Start, Stop, Continue

```markdown
# Retrospective: [Sprint/Project Name]

**Date:** YYYY-MM-DD
**Facilitator:** [Name]
**Participants:** [Names]

## Start Doing

[Things we should begin doing]

| Suggestion | Votes | Action |
|------------|-------|--------|
| [Idea 1] | 5 | [Assigned action] |
| [Idea 2] | 3 | [Assigned action] |

## Stop Doing

[Things we should stop doing]

| Suggestion | Votes | Action |
|------------|-------|--------|
| [Idea 1] | 4 | [Assigned action] |

## Continue Doing

[Things working well that we should keep]

| Item | Why It's Working |
|------|------------------|
| [Practice 1] | [Reason] |
| [Practice 2] | [Reason] |

## Action Items

| Action | Owner | Due |
|--------|-------|-----|
| [Action 1] | @name | YYYY-MM-DD |
```

### Format: 4 Ls

```markdown
# Retrospective: [Sprint/Project Name]

## Liked

[What did we enjoy?]

- [Item 1]
- [Item 2]

## Learned

[What did we learn?]

- [Item 1]
- [Item 2]

## Lacked

[What was missing?]

- [Item 1]
- [Item 2]

## Longed For

[What do we wish we had?]

- [Item 1]
- [Item 2]

## Actions

[Concrete improvements based on the above]
```

## Personal Retrospective

For individual reflection:

```markdown
# Personal Retrospective: [Week/Month/Project]

## Achievements

- [What I accomplished]
- [Impact made]

## Challenges

- [Difficulty faced]
- [How I handled it]
- [What I'd do differently]

## Learnings

### Technical

- [Technical skill gained]

### Process

- [Process improvement learned]

### Interpersonal

- [Communication/collaboration learning]

## Feedback Received

- [Feedback and how I'm addressing it]

## Goals for Next Period

- [ ] [Specific, measurable goal]
- [ ] [Specific, measurable goal]
```

## Root Cause Analysis Techniques

### 5 Whys

```
Problem: Tests failed in CI but passed locally

Why 1: Test assumed specific file order
  ↓
Why 2: File order differs between OS/filesystem
  ↓
Why 3: Test used os.ReadDir without sorting
  ↓
Why 4: Developer didn't know ReadDir order is undefined
  ↓
Why 5: No linter check for this pattern

Root Cause: Missing static analysis for filesystem ordering assumptions
Action: Add custom lint rule, document in contributor guide
```

### Fishbone Diagram (Ishikawa)

```
                    ┌──────────────────────────────────────┐
                    │                                      │
    People          │     Process        Equipment         │
        │           │         │              │             │
        ├─ Training │         ├─ No review   ├─ Old tools  │
        │           │         │              │             │
        ├─ Workload ├─────────┴──────────────┴─────────────┤ PROBLEM
        │           │                                      │
        │           │                                      │
    Environment     │     Materials                        │
        │           │         │                            │
        ├─ Pressure │         ├─ Bad data                  │
                    │                                      │
                    └──────────────────────────────────────┘
```

### Contributing Factor Analysis

| Factor | Type | Contribution | Fixable? |
|--------|------|--------------|----------|
| Missing test | Technical | High | Yes |
| Time pressure | Environment | Medium | Partially |
| Unclear requirements | Process | High | Yes |
| Legacy code | Technical | Medium | Long-term |

## Metrics for Improvement

Track over time:

| Metric | Description | Target |
|--------|-------------|--------|
| MTTR | Mean time to recovery | < 1 hour |
| Incident frequency | Incidents per month | Decreasing |
| Test coverage | Code coverage % | > 95% |
| Lead time | Idea to production | < 1 week |
| Change failure rate | Failed deployments | < 5% |

## Action Item Best Practices

Good action items are:

- **Specific:** "Add validation to ProcessEvent" not "Improve validation"
- **Measurable:** "Reduce MTTR to < 30 minutes"
- **Achievable:** Realistic scope
- **Relevant:** Directly addresses root cause
- **Time-bound:** Has a due date

```markdown
### Action Item Template

**Action:** [Specific task]
**Owner:** @name
**Due:** YYYY-MM-DD
**Success Criteria:** [How we'll know it's done]
**Priority:** P0/P1/P2
**Status:** Open | In Progress | Done
```

## Anti-Patterns

**DON'T:**
- Blame individuals
- Skip the retrospective when things went well
- Create action items with no owner
- Let action items go stale
- Focus only on negatives
- Hold retrospectives without the people involved

**DO:**
- Focus on systems and processes
- Celebrate successes
- Assign owners and due dates
- Follow up on action items
- Balance positive and negative
- Include all relevant voices

## Related Skills

- `incident-response` - Handling incidents in real-time
- `monitoring` - Detecting issues early
- `documentation-writing` - Writing post-mortems
- `critical-thinking` - Root cause analysis
