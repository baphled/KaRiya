---
name: technical-leadership
description: Leading technical decisions, writing RFCs, building consensus, and driving architecture
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
  always_active: false
---

## What I do

Guide technical leadership practices including writing RFCs, making architectural decisions, building consensus, and driving technical direction.

## When to use me

Use this skill when:
- Proposing significant architectural changes
- Writing technical proposals (RFCs)
- Making decisions that affect multiple teams/components
- Building consensus on technical direction
- Documenting architectural decisions (ADRs)

## RFC (Request for Comments) Process

### When to Write an RFC

| Change Type | RFC Required? |
|-------------|---------------|
| New major feature | Yes |
| Architectural change | Yes |
| Breaking API change | Yes |
| New external dependency | Consider |
| Internal refactoring | No (unless significant) |
| Bug fix | No |
| Small feature | No |

### RFC Template

```markdown
# RFC: [Title]

**Author:** [Name]
**Status:** Draft | Under Review | Accepted | Rejected | Superseded
**Created:** YYYY-MM-DD
**Last Updated:** YYYY-MM-DD

## Summary

[One paragraph summary of the proposal]

## Motivation

### Problem Statement

[What problem are we solving?]

### Goals

- [Goal 1]
- [Goal 2]

### Non-Goals

- [Explicitly out of scope]

## Background

[Context needed to understand the proposal]

### Current State

[How things work today]

### Prior Art

[How others have solved this]

## Proposal

### Overview

[High-level description of the solution]

### Detailed Design

[Technical details, interfaces, data structures]

```go
// Example code showing the proposed interface
type NewThing interface {
    DoSomething(ctx context.Context) error
}
```

### Migration Path

[How to get from current state to proposed state]

## Alternatives Considered

### Alternative 1: [Name]

**Description:** [What is it?]
**Pros:** [Benefits]
**Cons:** [Drawbacks]
**Why not:** [Why we didn't choose this]

### Alternative 2: [Name]

[Same structure]

## Trade-offs

| Aspect | Benefit | Cost |
|--------|---------|------|
| Performance | Faster queries | More memory |
| Complexity | Simpler API | More internal complexity |
| Flexibility | Easier to extend | Harder to optimize |

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| [Risk 1] | Medium | High | [Mitigation] |

## Success Criteria

[How will we know this succeeded?]

- [ ] [Measurable criterion 1]
- [ ] [Measurable criterion 2]

## Timeline

| Phase | Duration | Deliverable |
|-------|----------|-------------|
| Design | 1 week | Approved RFC |
| Implementation | 2 weeks | Feature complete |
| Testing | 1 week | All tests pass |
| Rollout | 1 week | In production |

## Open Questions

- [ ] [Question that needs resolution]

## References

- [Link to related doc]
- [Link to prior discussion]

## Appendix

[Additional details, diagrams, data]
```

## ADR (Architecture Decision Record)

### ADR Template

```markdown
# ADR-[NUMBER]: [Title]

**Date:** YYYY-MM-DD
**Status:** Proposed | Accepted | Deprecated | Superseded by ADR-XXX

## Context

[What is the issue that we're seeing that is motivating this decision?]

## Decision

[What is the change that we're proposing/have agreed to implement?]

## Consequences

### Positive

- [Benefit 1]
- [Benefit 2]

### Negative

- [Drawback 1]
- [Drawback 2]

### Neutral

- [Observation 1]

## Compliance

[How will we ensure this decision is followed?]

## References

- [RFC if applicable]
- [Discussion links]
```

### Example ADR

```markdown
# ADR-001: Use Subdirectory Structure for Intents

**Date:** 2026-01-15
**Status:** Accepted

## Context

Intent files are growing too large (some exceed 1,500 lines). Large files:
- Are hard to navigate and understand
- Make code review difficult
- Indicate poor separation of concerns

## Decision

All intents must use a subdirectory structure with 5+ files:
- `context.go` - Input parameters
- `constants.go` - State enum
- `messages.go` - Message types
- `result.go` - Output type
- `intent.go` - Implementation (max 400 lines)

## Consequences

### Positive
- Clear file responsibilities
- Smaller, focused files
- Easier code review
- Enforced by automated check

### Negative
- More files to navigate
- Migration effort for existing intents
- Learning curve for new patterns

## Compliance

Check #17-29 in `check-intent-architecture.sh` enforce this structure.
Pre-commit hook blocks violations.
```

## Building Consensus

### Stakeholder Matrix

| Stakeholder | Interest | Influence | Approach |
|-------------|----------|-----------|----------|
| Architecture team | High | High | Early collaboration |
| Feature teams | Medium | Medium | Review and feedback |
| QA | Medium | Low | Inform of testing impact |
| DevOps | Low | Medium | Consult on deployment |

### Consensus Techniques

**1. Write It Down First**

```
Draft RFC → Share for async feedback → Incorporate feedback → Sync discussion → Decision
```

**2. Decision Meeting Structure**

```markdown
## Decision Meeting: [Topic]

**Goal:** Reach decision on [specific question]
**Time:** 30 minutes max
**Attendees:** [Decision makers only]

### Agenda

1. Summary of options (5 min)
2. Key trade-offs (5 min)
3. Discussion (15 min)
4. Decision (5 min)

### Pre-Read

- RFC: [link]
- Discussion thread: [link]
```

**3. Disagree and Commit**

When consensus isn't possible:

```markdown
## Decision Record

**Decision:** We will proceed with Option A
**Rationale:** [Why this option]

**Dissenting View:** @engineer preferred Option B because [reasons]
**Commitment:** Despite disagreement, team commits to fully supporting Option A

**Review Date:** [Date to revisit if needed]
```

## Technical Debt Management

### Debt Tracking

```markdown
## Technical Debt Register

| ID | Description | Impact | Effort | Priority | Owner |
|----|-------------|--------|--------|----------|-------|
| TD-001 | Legacy form models | Medium | High | Medium | @dev |
| TD-002 | Inconsistent error handling | High | Medium | High | @dev |
```

### Prioritization Framework

```
Priority = (Impact × Frequency) / Effort

Where:
- Impact: 1 (low) to 5 (critical)
- Frequency: How often this causes pain
- Effort: 1 (quick fix) to 5 (major project)
```

## Communication Patterns

### Technical Announcements

```markdown
## [TECH] [Topic] - [One Line Summary]

**TL;DR:** [One sentence summary]

### What's Changing

[Brief description]

### Why

[Motivation]

### Impact

- **Teams affected:** [List]
- **Migration required:** Yes/No
- **Timeline:** [Dates]

### Action Required

- [ ] [Specific action by date]

### Resources

- RFC: [link]
- Migration guide: [link]
- Questions: [channel/contact]
```

### Status Updates

```markdown
## [Project] Status Update - Week [N]

### Summary
[Traffic light: 🟢 On Track | 🟡 At Risk | 🔴 Blocked]

### Progress This Week
- [Completed item]
- [Completed item]

### Next Week
- [Planned item]
- [Planned item]

### Risks/Blockers
- [Risk with mitigation]

### Metrics
- Coverage: X%
- Tests passing: Y/Z
- PRs merged: N
```

## Leading Technical Reviews

### Code Review Leadership

```markdown
## Review Priorities

1. **Architecture** - Does it fit our patterns?
2. **Correctness** - Does it work correctly?
3. **Security** - Any vulnerabilities?
4. **Performance** - Any concerns?
5. **Maintainability** - Is it readable/testable?
6. **Style** - Follows conventions?

Focus reviews on 1-3; automate 4-6 where possible.
```

### Design Review Checklist

```markdown
## Design Review: [Feature]

### Alignment
- [ ] Fits existing architecture
- [ ] Follows established patterns
- [ ] No unnecessary new dependencies

### Scalability
- [ ] Handles edge cases
- [ ] Performance considered
- [ ] Resource usage acceptable

### Maintainability
- [ ] Clear interfaces
- [ ] Testable design
- [ ] Documented decisions

### Security
- [ ] Input validation
- [ ] No sensitive data exposure
- [ ] Access control considered

### Operability
- [ ] Observable (logging, metrics)
- [ ] Debuggable
- [ ] Failure modes handled
```

## Related Skills

- `trade-off-analysis` - Evaluating alternatives
- `justify-decision` - Explaining decisions
- `documentation-writing` - Technical writing
- `systems-thinker` - Understanding interconnections
