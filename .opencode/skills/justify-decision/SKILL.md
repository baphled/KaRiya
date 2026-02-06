---
name: justify-decision
description: Provide evidence-based justification for architectural and design decisions
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
  always_active: false
---

## What I do

Provide clear, evidence-based justifications for technical decisions. I help explain WHY a decision was made, not just defend it.

## When to use me

Use this skill when:
- Reviewer questions an architectural decision
- Explaining why a pattern was chosen
- Defending against "why didn't you..." questions
- Documenting decision rationale

## Core Principle

**Good justification explains the reasoning, not just the conclusion.**

A justification should answer:
- What problem were we solving?
- What alternatives did we consider?
- Why did we choose this approach?
- What trade-offs did we accept?

## Justification Framework

### 1. Reference Documentation First

Always check if a documented decision exists:

```markdown
This follows our documented architecture pattern:

> "Screens must not import from intents package"
> - [AGENTS.md](../../../AGENTS.md#dependency-rules)

The pattern exists because [reason from docs].
```

**Documentation sources:**
- `AGENTS.md` - Architecture rules
- `docs/ARCHITECTURE.md` - System design
- `docs/guides/*.md` - Feature-specific guidance
- ADRs (Architecture Decision Records) - Historical decisions

### 2. Explain the Problem

Before justifying the solution, establish the problem:

```markdown
## Problem
We needed to [requirement] while maintaining [constraint].

The challenge was [specific difficulty].
```

### 3. Present Alternatives Considered

Show you thought about options:

```markdown
## Alternatives Considered

### Option A: [Name]
- **Pros**: [benefits]
- **Cons**: [drawbacks]
- **Rejected because**: [specific reason]

### Option B: [Name] (Chosen)
- **Pros**: [benefits]
- **Cons**: [drawbacks]
- **Chosen because**: [specific reason]

### Option C: [Name]
- **Pros**: [benefits]
- **Cons**: [drawbacks]
- **Rejected because**: [specific reason]
```

### 4. Explain Trade-offs Accepted

Be honest about downsides:

```markdown
## Trade-offs

We accepted these trade-offs:

| Trade-off | Mitigation |
|-----------|------------|
| Slightly more verbose | Improved readability |
| Extra interface | Better testability |
| More files | Clear separation of concerns |
```

### 5. Provide Evidence

Support claims with concrete evidence:

```markdown
## Evidence

### Performance
Benchmarks show this approach is [X]% faster:
[benchmark results]

### Maintainability
Similar pattern used successfully in:
- `internal/cli/intents/browse_timeline/`
- `internal/cli/intents/burst_management/`

### Correctness
100% test coverage on this component:
[coverage output]
```

## Common Justification Scenarios

### "Why not use X library?"

```markdown
We considered [library] but chose not to use it because:

1. **Dependency weight**: Adds [X]KB to binary for [limited feature]
2. **Maintenance**: Last commit [date], [issues] open issues
3. **Fit**: Designed for [use case], we need [different use case]
4. **Consistency**: Project uses [standard approach] elsewhere

Our implementation is [X] lines and gives us control over [specific aspect].
```

### "Why this architecture pattern?"

```markdown
This follows the [Pattern Name] pattern because:

1. **Separation of concerns**: [Component A] handles [X], [Component B] handles [Y]
2. **Testability**: Each component can be tested in isolation
3. **Project consistency**: Same pattern used in [other areas]
4. **Documented standard**: See [AGENTS.md#architecture]

The alternative ([other pattern]) would require [problematic thing].
```

### "Why not refactor to X?"

```markdown
The suggested refactor would:

1. **Scope**: Require changes to [N] files beyond this PR
2. **Risk**: Touch [critical area] without full regression testing
3. **Value**: Provide [benefit] but at cost of [downside]

I suggest we:
- Complete this PR with current approach
- Create issue #[X] to track the refactor
- Plan refactor with proper scope and testing

This keeps the PR focused and reduces risk.
```

### "Why this naming convention?"

```markdown
The name follows our project conventions:

- **Pattern**: `[Entity][Type]Screen` (see [SCREEN_NAMING.md])
- **Consistency**: Matches `FactListScreen`, `BurstDetailScreen`
- **Discoverability**: Easy to find via glob: `*ListScreen*`

Alternative names considered:
- `ListView` - Conflicts with UIKit naming
- `FactScreen` - Ambiguous (list vs detail?)
- `ListOfFacts` - Doesn't match established pattern
```

## Evidence Types

| Claim Type | Evidence Needed |
|------------|-----------------|
| Performance | Benchmarks, profiling |
| Correctness | Tests, formal verification |
| Maintainability | Similar examples, metrics |
| Security | Threat model, security review |
| Usability | User feedback, UX testing |
| Consistency | Documentation, existing code |

## Justification Quality Checklist

Before submitting justification:

- [ ] References documentation where applicable
- [ ] Explains the problem being solved
- [ ] Lists alternatives considered
- [ ] States trade-offs honestly
- [ ] Provides concrete evidence
- [ ] Acknowledges valid points in feedback
- [ ] Offers path forward (even if disagreeing)

## Anti-Patterns

**Avoid these weak justifications:**

| Weak | Better |
|------|--------|
| "It's best practice" | "This pattern improves [specific thing] because..." |
| "I prefer it this way" | "This approach has [measurable benefit]" |
| "It's always done like this" | "Project standard (see [doc]) for consistency" |
| "Trust me" | "Evidence: [tests/benchmarks/docs]" |
| "The reviewer is wrong" | "The feedback raises [valid point], however..." |

## Template Response

```markdown
Thanks for the feedback on [area].

## Why This Approach

[Brief explanation of the decision]

## Problem Being Solved

[What requirement/constraint drove this]

## Alternatives Considered

1. **[Alternative A]**: [Why not chosen]
2. **[Alternative B]**: [Why not chosen]
3. **[Chosen approach]**: [Why this works]

## Evidence

[Tests/benchmarks/documentation supporting the choice]

## Trade-offs Accepted

[Honest acknowledgment of downsides and mitigations]

## Path Forward

[What you propose - accept their feedback, keep current, or compromise]

Let me know if you'd like to discuss further or if I've missed something in my analysis.
```

## Related Skills

- `evaluate-change-request` - Determines if justification needed
- `prove-correctness` - Provides evidence through tests
- `trade-off-analysis` - Structured comparison of options
- `critical-thinking` - Rigorous analysis
- `architecture` - Architectural patterns reference
- `respond-to-review` - Crafting the response
