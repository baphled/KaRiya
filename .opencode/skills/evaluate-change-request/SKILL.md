---
name: evaluate-change-request
description: Critically evaluate review feedback before accepting or rejecting
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
  always_active: false
---

## What I do

Critically analyse review feedback to determine its validity, category, and appropriate response. I prevent blind acceptance of changes that may not improve the application.

## When to use me

Use this skill when receiving any review feedback to:
- Determine if the feedback is valid
- Categorise the type of feedback
- Decide whether to accept, challenge, or clarify

## Core Principle

**The goal is to improve the application, not to please reviewers or win arguments.**

A change should only be made if it demonstrably improves:
- Correctness (fewer bugs)
- Security (fewer vulnerabilities)
- Maintainability (easier to change)
- Performance (faster, less resource usage)
- Clarity (easier to understand)

## Evaluation Framework

### Step 1: Understand the Feedback

Before reacting, ensure you understand:

| Question | Why It Matters |
|----------|----------------|
| What exactly is being requested? | Avoid misinterpretation |
| What problem does this solve? | Understand the motivation |
| Is there evidence provided? | Facts vs opinions |
| Does it reference standards/docs? | Objective vs subjective |

**If unclear, ask clarifying questions before proceeding.**

### Step 2: Categorise the Feedback

| Category | Characteristics | Default Response |
|----------|-----------------|------------------|
| **Objective Bug** | Demonstrably incorrect behaviour | Accept |
| **Security Issue** | Vulnerability, exposure risk | Accept immediately |
| **Architecture Violation** | Breaks documented patterns | Evaluate against docs |
| **Performance Issue** | Measurable degradation | Require benchmarks |
| **Code Quality** | Complexity, duplication, naming | Evaluate impact |
| **Style Preference** | No project standard exists | Discuss briefly |
| **Misunderstanding** | Reviewer misread the code | Clarify |
| **Scope Creep** | Beyond original PR scope | Defer to new PR |

### Step 3: Assess Validity

For each piece of feedback, ask:

```
1. Is this objectively correct?
   ├─ YES → Strong evidence to accept
   └─ NO/MAYBE → Continue evaluation

2. Is there evidence provided?
   ├─ YES → Evaluate the evidence
   └─ NO → Request evidence or provide counter-evidence

3. Does this align with project standards?
   ├─ YES → Strong reason to accept
   ├─ NO → Strong reason to challenge
   └─ NO STANDARD → Discuss and potentially establish one

4. What are the trade-offs?
   ├─ Net positive → Accept
   ├─ Net negative → Challenge with evidence
   └─ Unclear → Discuss trade-offs explicitly

5. Is this within PR scope?
   ├─ YES → Address in this PR
   └─ NO → Suggest deferring to new PR/issue
```

### Step 4: Determine Response

| Verdict | Action | Skills to Use |
|---------|--------|---------------|
| **Accept** | Implement change, write test if bug | `respond-to-review` |
| **Challenge** | Provide counter-evidence | `prove-correctness`, `justify-decision` |
| **Discuss** | Present trade-offs | `trade-off-analysis`, `respond-to-review` |
| **Clarify** | Ask specific questions | `respond-to-review` |
| **Defer** | Create issue for later | `respond-to-review` |

## Challenge Criteria

You should challenge feedback when:

### 1. No Evidence Provided
```
Reviewer: "This is inefficient"
Response: "Could you clarify what inefficiency you're seeing? 
I can run benchmarks to measure the actual performance."
```

### 2. Contradicts Documented Standards
```
Reviewer: "Move this to the service layer"
Response: "Our architecture guide (docs/ARCHITECTURE.md) specifies 
this logic belongs in the domain layer because [reason]. 
Should we discuss updating the standard?"
```

### 3. Introduces New Problems
```
Reviewer: "Use a map instead of slice for O(1) lookup"
Response: "A map would improve lookup but we need ordered iteration 
for [reason]. The slice is only ~10 items so O(n) is negligible.
Happy to discuss trade-offs."
```

### 4. Subjective Without Project Standard
```
Reviewer: "Rename 'items' to 'entries'"
Response: "Both names seem reasonable. Is there a project naming 
convention I should follow? If not, I'd prefer 'items' for 
consistency with [other file]."
```

### 5. Out of Scope
```
Reviewer: "While you're here, could you also refactor X?"
Response: "Good suggestion! That's outside the scope of this PR 
though. I've created issue #123 to track it. Happy to address 
in a follow-up PR."
```

## Accept Criteria

You should accept feedback when:

1. **It's objectively correct** - The code has a bug, security issue, or clear defect
2. **Evidence is provided** - Benchmarks, test failures, documentation references
3. **It aligns with project standards** - Even if you disagree with the standard
4. **The trade-offs favour the change** - Net improvement to the application
5. **It improves clarity** - Makes code easier to understand for others

## Red Flags (Suspicious Feedback)

Be extra critical of:

| Red Flag | Why | Response |
|----------|-----|----------|
| "Best practice" without context | Often subjective | Ask for specific reasoning |
| "Always/Never do X" | Absolutes are rarely true | Ask about this specific case |
| Major refactor suggestion | May be scope creep | Discuss if warranted |
| Style-only changes | Low value, high churn | Follow project standard or discuss |
| "I would do it differently" | Preference, not improvement | Ask what problem it solves |

## Documentation

After evaluation, document your reasoning:

```markdown
## Feedback: [Quote the feedback]

### Category
[Bug/Security/Architecture/Performance/Style/etc.]

### Assessment
- **Validity**: [Valid/Invalid/Unclear]
- **Evidence provided**: [Yes/No]
- **Project standard exists**: [Yes/No/Partial]
- **Trade-offs**: [List pros and cons]

### Decision
[Accept/Challenge/Discuss/Clarify/Defer]

### Reasoning
[Why this decision was made]

### Action
[What will be done]
```

## Integration with Other Skills

| Skill | When to Use |
|-------|-------------|
| `critical-thinking` | Throughout evaluation |
| `epistemic-rigor` | Distinguish facts from assumptions |
| `devils-advocate` | Challenge your own position too |
| `prove-correctness` | When you need to provide evidence |
| `justify-decision` | When explaining architectural choices |
| `trade-off-analysis` | When comparing alternatives |

## Related Skills

- `pr-monitor` - Orchestrates the review process
- `prove-correctness` - Write tests as evidence
- `justify-decision` - Explain choices with evidence
- `respond-to-review` - Craft the actual response
- `critical-thinking` - Core analytical skill
- `devils-advocate` - Challenge assumptions
