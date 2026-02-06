---
name: trade-off-analysis
description: Systematically evaluate trade-offs when comparing alternatives
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
  always_active: false
---

## What I do

Provide structured analysis of trade-offs when choosing between alternatives. I help make decisions explicit and defensible.

## When to use me

Use this skill when:
- Reviewer suggests an alternative approach
- Choosing between design options
- Evaluating refactoring suggestions
- Comparing libraries or tools

## Core Principle

**Every decision involves trade-offs. Make them explicit.**

There's rarely a "right" answer - only trade-offs. Good decisions:
- Acknowledge what's being given up
- Explain why the trade-off is acceptable
- Consider the specific context

## Trade-off Framework

### 1. Identify the Dimensions

Common dimensions to compare:

| Dimension | Questions |
|-----------|-----------|
| **Performance** | Speed? Memory? Throughput? |
| **Complexity** | Lines of code? Cognitive load? |
| **Maintainability** | Easy to change? Easy to understand? |
| **Testability** | Easy to test? Mockable? |
| **Correctness** | Handles edge cases? Type-safe? |
| **Consistency** | Matches existing patterns? |
| **Dependencies** | External deps? Coupling? |
| **Flexibility** | Easy to extend? Configurable? |

### 2. Score Each Alternative

Create a comparison matrix:

```markdown
| Dimension | Option A | Option B | Option C |
|-----------|----------|----------|----------|
| Performance | ++ | + | +++ |
| Complexity | + | ++ | - |
| Maintainability | ++ | + | - |
| Testability | ++ | ++ | + |
| Consistency | +++ | + | - |
| **Total** | 10 | 7 | 4 |
```

Legend: `+++` = excellent, `++` = good, `+` = acceptable, `-` = poor, `--` = very poor

### 3. Weight by Context

Not all dimensions are equally important:

```markdown
For this specific case, priorities are:
1. **Maintainability** (weight: 3x) - Long-lived code
2. **Consistency** (weight: 2x) - Team familiarity
3. **Performance** (weight: 1x) - Not on hot path
```

Weighted comparison:
```markdown
| Dimension (Weight) | Option A | Option B |
|--------------------|----------|----------|
| Performance (1x) | 2 | 3 |
| Maintainability (3x) | 6 | 3 |
| Consistency (2x) | 6 | 2 |
| **Weighted Total** | 14 | 8 |
```

### 4. Document the Decision

```markdown
## Decision: Option A

### Context
[What we're deciding and why it matters]

### Options Considered
1. **Option A**: [Brief description]
2. **Option B**: [Brief description]

### Trade-off Analysis
[Comparison matrix]

### Decision Rationale
Chose Option A because:
- [Primary reason]
- [Secondary reason]

Accepted trade-offs:
- [Downside 1]: Mitigated by [approach]
- [Downside 2]: Acceptable because [reason]

### Reversibility
This decision is [easily/moderately/difficult to] reverse because [reason].
```

## Common Trade-off Scenarios

### Performance vs Readability

```markdown
## Trade-off: Optimised Loop vs Readable Code

### Option A: Optimised
```go
for i := 0; i < len(items); i++ {
    // Avoids bounds checks, ~10% faster
}
```

### Option B: Readable
```go
for _, item := range items {
    // Clearer intent, idiomatic Go
}
```

### Analysis
| Dimension | Optimised | Readable |
|-----------|-----------|----------|
| Performance | ++ | + |
| Readability | - | ++ |
| Idiomatic | - | ++ |

### Decision
**Readable** - This code runs once per user action (not hot path).
The 10% micro-optimisation isn't worth the readability cost.
```

### Abstraction vs Simplicity

```markdown
## Trade-off: Interface vs Concrete Type

### Option A: Interface
```go
type Storage interface {
    Save(data []byte) error
    Load() ([]byte, error)
}
```

### Option B: Concrete
```go
func SaveToFile(path string, data []byte) error
func LoadFromFile(path string) ([]byte, error)
```

### Analysis
| Dimension | Interface | Concrete |
|-----------|-----------|----------|
| Testability | ++ | - |
| Flexibility | ++ | - |
| Simplicity | - | ++ |
| Current needs | + | ++ |

### Decision
**Interface** - We need to mock storage in tests.
The abstraction cost is justified by testability benefit.
```

### Library vs Custom

```markdown
## Trade-off: External Library vs Custom Code

### Option A: Use library X
- 15KB added to binary
- Well-tested, maintained
- More features than we need

### Option B: Custom implementation
- ~50 lines of code
- Exactly what we need
- We maintain it

### Analysis
| Dimension | Library | Custom |
|-----------|---------|--------|
| Effort | ++ | - |
| Binary size | - | ++ |
| Control | - | ++ |
| Maintenance | ++ | - |

### Decision
**Custom** - Our needs are simple (50 lines).
Library adds 15KB for features we won't use.
We can always switch to library if needs grow.
```

### Consistency vs Improvement

```markdown
## Trade-off: Match Existing Pattern vs Better Approach

### Option A: Match existing
- Consistent with codebase
- Team knows this pattern
- Perpetuates technical debt

### Option B: New approach
- Better design
- Learning curve
- Inconsistent until refactored

### Analysis
| Dimension | Existing | New |
|-----------|----------|-----|
| Consistency | ++ | - |
| Quality | - | ++ |
| Team velocity | ++ | - |
| Long-term | - | ++ |

### Decision
**Existing** for this PR, with plan:
1. Complete this PR consistently
2. Create RFC for new pattern
3. Migrate in dedicated refactor PR
```

## Template for Review Responses

When reviewer suggests alternative:

```markdown
Thanks for suggesting [alternative]. I evaluated both approaches:

## Comparison

| Dimension | Current | Suggested |
|-----------|---------|-----------|
| [Dim 1] | [score] | [score] |
| [Dim 2] | [score] | [score] |
| [Dim 3] | [score] | [score] |

## Context
For this specific case:
- [Relevant constraint 1]
- [Relevant constraint 2]

## Conclusion
I believe [current/suggested] is better here because [reasons].

[If disagreeing]: Happy to discuss if I've missed something in my analysis.
[If agreeing]: I'll make the change.
```

## Reversibility Assessment

Always consider how hard it is to change later:

| Reversibility | Characteristics | Approach |
|---------------|-----------------|----------|
| **Easy** | Internal, no API change | Decide quickly, change if wrong |
| **Moderate** | Some refactoring needed | Analyse carefully, document |
| **Difficult** | Public API, data format | Extensive analysis, get consensus |
| **Irreversible** | Database migration, security | Maximum diligence, formal review |

## Quality Checklist

Before presenting trade-off analysis:

- [ ] All reasonable alternatives listed
- [ ] Dimensions relevant to context
- [ ] Weights justified by context
- [ ] Both pros AND cons for each option
- [ ] Decision clearly stated
- [ ] Accepted trade-offs acknowledged
- [ ] Reversibility assessed
- [ ] Path forward clear

## Related Skills

- `evaluate-change-request` - Triggers trade-off analysis
- `justify-decision` - Uses trade-off analysis as evidence
- `critical-thinking` - Rigorous evaluation
- `systems-thinker` - Consider broader implications
- `devils-advocate` - Challenge your own analysis
